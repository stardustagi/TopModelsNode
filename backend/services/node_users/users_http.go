package nodeUsers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/databases"
	"github.com/stardustagi/TopLib/libs/errors"
	"github.com/stardustagi/TopLib/libs/logs"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/uuid"
	"github.com/stardustagi/TopLib/protocol"
	"github.com/stardustagi/TopModelsNode/backend"
	"github.com/stardustagi/TopModelsNode/backend/services/node_llm"
	"github.com/stardustagi/TopModelsNode/constants"
	"github.com/stardustagi/TopModelsNode/models"
	"github.com/stardustagi/TopModelsNode/protocol/requests"
	"github.com/stardustagi/TopModelsNode/protocol/responses"
	mailgateway "github.com/stardustagi/TopModelsNode/third_party/mail_gateway"

	"go.uber.org/zap"
)

type NodeUsersHttpService struct {
	logger      *zap.Logger
	ctx         context.Context
	cancelCtx   context.CancelFunc
	dao         databases.BaseDao
	rds         redis.RedisCli
	mu          sync.RWMutex // 读写锁保护用户数据
	userInfosMu sync.RWMutex // 读写保护用户信息缓存
	app         *backend.Application
}

var (
	nodeUsersHttpServiceInstance *NodeUsersHttpService
	nodeUsersHttpServiceOnce     sync.Once
)

// GetNodeUsersHttpServiceInstance 获取 HTTP 用户服务实例
func GetNodeUsersHttpServiceInstance() *NodeUsersHttpService {
	nodeUsersHttpServiceOnce.Do(func() {
		nodeUsersHttpServiceInstance = NewNodeUsersHttpService()
	})
	return nodeUsersHttpServiceInstance
}

// NewNodeUsersHttpService 创建新的 HTTP 用户服务
func NewNodeUsersHttpService() *NodeUsersHttpService {
	ctx, cancel := context.WithCancel(context.Background())
	return &NodeUsersHttpService{
		logger:    logs.GetLogger("NodeUsersHttpService"),
		ctx:       ctx,
		cancelCtx: cancel,
		dao:       databases.GetDao(),
		rds: redis.NewRedisView(redis.GetRedisDb(),
			constants.NodeUserKeyPrefix,
			logs.GetLogger("NodeUsersRedis")),
	}
}

func (nus *NodeUsersHttpService) Start(app *backend.Application) {
	if app == nil {
		panic("请设置后端应用")
	}
	nus.app = app
	nus.initialization()
	nus.logger.Info("Starting NodeUsersHttpService...")
}

func (nus *NodeUsersHttpService) Stop() {
	nus.logger.Info("Stopping NodeUsersHttpService...")
	nus.cancelCtx()
	nus.logger.Info("NodeUsersHttpService stopped.")
}

// NodeUserRegister 注册节点用户
// godoc 注册节点用户
// 注册节点用户
// @Summary 注册节点用户
// @Description 注册节点用户
// @Tags Node
// @Accept json
// @Produce json
// @Param request body requests.RegisterUserRequest true "注册节点用户请求"
// @Success 200 {object} responses.DefaultResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /users/system/register [post]
func (nus *NodeUsersHttpService) NodeUserRegister(c echo.Context,
	req requests.RegisterUserRequest,
	resp responses.DefaultResponse) error {
	nus.logger.Debug("Node user register request", zap.Any("request", req))
	// 检查用户是否已存在
	session := nus.dao.NewSession()
	nodeUser := &models.NodeUsers{Email: req.Email}
	ok, err := session.FindOne(nodeUser)
	if err != nil {
		nus.logger.Error("查询节点用户失败", zap.Error(err), zap.String("email", req.Email))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}
	if ok {
		nus.logger.Warn("节点用户已存在", zap.String("email", req.Email))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	// 生成盐值和密码哈希
	salt := uuid.GenString(8)
	hashedPwd, err := nus.nodeUserMailEncodeToken(req.Email, req.Password, salt)
	if err != nil {
		nus.logger.Error("密码加密失败", zap.Error(err), zap.String("email", req.Email))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}

	// 创建新的节点用户
	nodeUser.Password = hashedPwd
	nodeUser.Salt = salt
	nodeUser.CreatedAt = time.Now().Unix()
	nodeUser.IsActive = 1

	// 保存到数据库
	_, err = session.InsertOne(nodeUser)
	if err != nil {
		nus.logger.Error("保存节点用户到数据库失败", zap.Error(err), zap.String("email", req.Email))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}
	activeCode := uuid.GenString(8)
	key := constants.NodeUserMailVerifyKey(nodeUser.Id)
	// 保存激活码到Redis
	err = nus.rds.Set(nus.ctx, key, []byte(activeCode), constants.NodeUserMailVerifyExpireTimeString)

	link := fmt.Sprintf("%s/api/system/nodeUserActive?nodeUserId=%d&nodeUserActiveCode=%s", constants.Domain, nodeUser.Id, activeCode)
	content := fmt.Sprintf(
		"节点用户%s注册成功,用户ID:%d,激活码%s,有效时间为:%s,<br><a href=\"%s\">点击此处激活账号</a>",
		req.Email, nodeUser.Id, activeCode, constants.NodeUserMailVerifyExpireTimeString, link,
	)

	err = mailgateway.SendEmail(
		req.Email, "注册成功，请激活您的节点用户账号", content,
	)
	if err != nil {
		nus.logger.Error("发送注册成功邮件失败", zap.Error(err), zap.String("email", req.Email))
		protocol.Response(c, constants.ErrUserRegFailed, nil)
	}
	nus.logger.Info("节点用户注册成功", zap.String("email", req.Email), zap.String("code", activeCode))
	return protocol.Response(c, nil, "注册成功")
}

// NodeUserActive 激活节点用户
// godoc 激活节点用户
// 激活节点用户
// @Summary 激活节点用户
// @Description 激活节点用户
// @Tags Node
// @Accept json
// @Produce json
// @Param nodeUserId query string true "节点用户ID"
// @Param nodeUserActiveCode query string true "节点用户激活码"
// @Success 200 {object} responses.DefaultResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /users/system/nodeUserActive [get]
func (nus *NodeUsersHttpService) NodeUserActive(c echo.Context, req requests.DefaultRequest, resp responses.DefaultResponse) error {
	logger := logs.GetLogger("NodeUserActive called ")
	nodeUserId := c.QueryParams().Get("nodeUserId")
	nodeUserActiveCode := c.QueryParams().Get("nodeUserActiveCode")
	nodeUserIdNumber, err := strconv.ParseInt(nodeUserId, 10, 64)
	if err != nil {
		logger.Error("Failed to parse node user id", zap.String("nodeUserId", nodeUserId), zap.Error(err))
		return protocol.Response(c, constants.ErrInvalidParams, nil)
	}

	ok, err := nus.nodeUserActive(nodeUserIdNumber, nodeUserActiveCode)
	if err != nil || !ok {
		logger.Error("Failed to active node user", zap.Error(err))
		return protocol.Response(c, constants.ErrInternalServer, nil)
	}
	return protocol.Response(c, nil, "激活成功")
}

// LoginFromByEmail 用户登录
// godoc
// 登录用户
// @Summary 登录用户
// @Description 登录用户
// @Success 200 {object} responses.LoginUserResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /node/public/login [post]
func (nus *NodeUsersHttpService) LoginFromByEmail(c echo.Context, req requests.LoginUserRequest, resp responses.LoginUserResponse) error {
	nus.logger.Info("LoginUser called", zap.String("username", req.Email))

	// 1. 验证参数
	if req.Email == "" {
		return protocol.Response(c, errors.New("邮箱不能为空", http.StatusBadRequest), nil)
	}
	if req.Password == "" {
		return protocol.Response(c, errors.New("密码不能为空", http.StatusBadRequest), nil)
	}

	// 2. 查询用户
	user := &models.NodeUsers{Email: req.Email, Deleted: 0}
	session := nus.dao.NewSession()
	defer session.Close()
	has, err := session.FindOne(user)
	if err != nil {
		nus.logger.Error("query user failed", zap.Error(err))
		return protocol.Response(c, errors.New("查询用户失败", http.StatusInternalServerError), nil)
	}
	if !has {
		return protocol.Response(c, errors.New("用户不存在或密码错误", http.StatusUnauthorized), nil)
	}

	// 3. 验证密码
	decMail, err := nus.nodeUserMailDecodeToken(req.Password, user.Password, user.Salt)
	if decMail != user.Email {
		return protocol.Response(c, errors.New("用户不存在或密码错误", http.StatusUnauthorized), nil)
	}

	// 4. 检查用户状态
	if user.IsActive != 1 {
		return protocol.Response(c, errors.New("用户已被禁用", http.StatusForbidden), nil)
	}

	// 5. 生成JWT token
	token, jwtString, err := nus.generateJWTToken(user.Id)
	if err != nil {
		nus.logger.Error("generate jwt token failed", zap.Error(err))
		return protocol.Response(c, errors.New("生成Token失败", http.StatusInternalServerError), nil)
	}

	// 6. 缓存Token到Redis
	if err := nus.cacheUserToken(user.Id, token, constants.NodeUserTokenExpire); err != nil {
		nus.logger.Warn("cache user token failed", zap.Error(err))
	}

	// 7. 更新最后登录时间
	user.LastUpdate = time.Now().Unix()
	updateSession := nus.dao.NewSession()
	defer updateSession.Close()
	_, err = updateSession.UpdateById(user.Id, user)
	if err != nil {
		nus.logger.Warn("update last login time failed", zap.Error(err))
	}

	// 8. 缓存用户信息
	if err := nus.cacheUserInfo(user); err != nil {
		nus.logger.Warn("cache user info failed", zap.Error(err))
	}

	// 9. 返回响应
	c.Response().Header().Set("Access-Control-Expose-Headers", "id, jwt")
	c.Response().Header().Add("jwt", jwtString)
	c.Response().Header().Add("id", fmt.Sprintf("%d", user.Id))
	resp.UserInfo = nus.convertToUserInfoResponse(user)
	resp.Jwt = token
	resp.ExpireAt = constants.NodeUserTokenExpire

	nus.logger.Info("user logged in successfully", zap.Int64("userId", user.Id))
	return protocol.Response(c, nil, resp)
}

// LogoutUser 用户登出
// godoc
// @Summary 用户登出
// @Description 用户登出
// @Tags 用户
// @Accept json
// @Produce json
// @Param userId path int true "用户ID"
// @Success 200 {object} responses.LogoutUserResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /node/user/logout [post]
func (nus *NodeUsersHttpService) LogoutUser(c echo.Context, req requests.LogoutUserRequest, resp responses.LogoutUserResponse) error {
	nus.logger.Info("LogoutUser called", zap.Int64("userIdg", req.UserID))

	// 清除Redis中的Token缓存
	if err := nus.clearUserCache(req.UserID); err != nil {
		nus.logger.Error("clear user cache failed", zap.Error(err))
		return protocol.Response(c, errors.New("登出失败", http.StatusInternalServerError), nil)
	}

	resp.Success = true
	resp.Message = "登出成功"

	nus.logger.Info("user logged out successfully", zap.Int64("userId", req.UserID))
	return protocol.Response(c, nil, resp)
}

// GetUserInfo 获取用户信息
// godoc
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param userId path int true "用户ID"
// @Success 200 {object} responses.GetUserInfoResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /node/user/getUserInfo [post]
func (nus *NodeUsersHttpService) GetUserInfo(c echo.Context, req requests.GetUserInfoRequest, resp responses.GetUserInfoResponse) error {
	nus.logger.Info("GetUserInfo called", zap.Int64("userId", req.UserID))

	// 1. 先从缓存获取
	cachedUser, err := nus.getUserInfoFromCache(req.UserID)
	if err == nil && cachedUser != nil {
		resp.UserInfo = nus.convertToUserInfoResponse(cachedUser)
		return protocol.Response(c, nil, resp)
	}

	// 2. 从数据库查询
	user := &models.NodeUsers{}
	session := nus.dao.NewSession()
	defer session.Close()
	has, err := session.FindById(req.UserID, user)
	if err != nil {
		nus.logger.Error("query user failed", zap.Error(err))
		return protocol.Response(c, errors.New("查询用户失败", http.StatusInternalServerError), nil)
	}
	if !has || user.Deleted != 0 {
		return protocol.Response(c, errors.New("用户不存在", http.StatusNotFound), nil)
	}

	// 3. 更新缓存
	if err := nus.cacheUserInfo(user); err != nil {
		nus.logger.Warn("cache user info failed", zap.Error(err))
	}

	resp.UserInfo = nus.convertToUserInfoResponse(user)
	return protocol.Response(c, nil, resp)
}

// UpdateUserInfo 更新用户信息
// godoc
// @Summary 更新用户信息
// @Description 更新用户信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body requests.UpdateUserInfoRequest true "用户信息"
// @Success 200 {object} responses.UpdateUserInfoResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /node/user/updateUserInfo [post]
func (nus *NodeUsersHttpService) UpdateUserInfo(c echo.Context, req requests.UpdateUserInfoRequest, resp responses.UpdateUserInfoResponse) error {
	nus.logger.Info("UpdateUserInfo called", zap.Int64("userId", req.UserID))

	// 1. 查询用户
	user := &models.NodeUsers{}
	session := nus.dao.NewSession()
	defer session.Close()
	has, err := session.FindById(req.UserID, user)
	if err != nil {
		nus.logger.Error("query user failed", zap.Error(err))
		return protocol.Response(c, errors.New("查询用户失败", http.StatusInternalServerError), nil)
	}
	if !has || user.Deleted != 0 {
		return protocol.Response(c, errors.New("用户不存在", http.StatusNotFound), nil)
	}

	// 2. 更新字段
	if req.Email != "" {
		// 检查新邮箱是否已被其他用户使用
		checkUser := &models.NodeUsers{Email: req.Email, Deleted: 0}
		checkSession := nus.dao.NewSession()
		defer checkSession.Close()
		exists, err := checkSession.Exists(checkUser)
		if err != nil {
			nus.logger.Error("check email exists failed", zap.Error(err))
			return protocol.Response(c, errors.New("检查邮箱失败", http.StatusInternalServerError), nil)
		}
		if exists && checkUser.Id != req.UserID {
			return protocol.Response(c, errors.New("该邮箱已被使用", http.StatusBadRequest), nil)
		}
		user.Email = req.Email
	}

	user.LastUpdate = time.Now().Unix()

	// 3. 更新数据库
	updateSession := nus.dao.NewSession()
	defer updateSession.Close()
	affected, err := updateSession.UpdateById(user.Id, user)
	if err != nil {
		nus.logger.Error("update user failed", zap.Error(err))
		return protocol.Response(c, errors.New("更新用户失败", http.StatusInternalServerError), nil)
	}
	if affected == 0 {
		nus.logger.Warn("no rows affected", zap.Int64("userId", req.UserID))
	}

	// 4. 更新缓存
	if err := nus.cacheUserInfo(user); err != nil {
		nus.logger.Warn("cache user info failed", zap.Error(err))
	}

	resp.UserInfo = nus.convertToUserInfoResponse(user)
	nus.logger.Info("user info updated successfully", zap.Int64("userId", user.Id))
	return protocol.Response(c, nil, resp)
}

// ListUsers 用户列表
// godoc
// @Summary 用户列表
// @Description 用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body requests.ListUsersRequest true "用户列表"
// @Success 200 {object} responses.ListUsersResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /node/user/listUsers [post]
func (nus *NodeUsersHttpService) ListUsers(c echo.Context, req requests.ListUsersRequest, resp responses.ListUsersResponse) error {
	nus.logger.Info("ListUsers called")

	// 1. 构建查询条件
	condiBean := &models.NodeUsers{Deleted: 0}

	switch req.Status {
	case 1:
		condiBean.IsActive = 1
	case 2:
		condiBean.IsActive = 0
	default:
		condiBean.IsActive = -1
	}

	// 2. 分页查询
	var users []models.NodeUsers
	pageable := databases.NewPageable(req.Page.Skip, req.Page.Limit, "created_at DESC")
	session := nus.dao.NewSession()
	defer session.Close()
	total, err := session.FindAndCount(&users, pageable, condiBean)
	if err != nil {
		nus.logger.Error("query users failed", zap.Error(err))
		return protocol.Response(c, errors.New("查询用户列表失败", http.StatusInternalServerError), nil)
	}

	// 3. 转换响应
	var userList []responses.UserInfoResponse
	for _, user := range users {
		userList = append(userList, nus.convertToUserInfoResponse(&user))
	}

	resp.Total = total
	resp.List = userList

	return protocol.Response(c, nil, resp)
}

// NodeCheckUserBalanceHandler 节点查询用户余额接口
// @Summary 节点查询用户余额接口
// @Description 节点查询用户余额接口
// @Tags 节点管理
// @Accept json
// @Produce json
// @Param nodeId header string true "节点ID"
// @Param req body requests.UserBalanceReq true "请求参数"
// @Success 200 {object} responses.UserBalanceResp
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /node/user/checkUserBalance [post]
func (nus *NodeUsersHttpService) NodeCheckUserBalanceHandler(ctx echo.Context, req requests.UserBalanceReq, resp responses.UserBalanceResp) error {
	nus.logger.Info("NodeCheckUserBalanceHandler called", zap.Int64("userId", req.UserID), zap.String("walletType", req.WalletType))
	nodeId, err := nus.getUserIdFromContext(ctx)
	if err != nil {
		nus.logger.Error("invalid user id", zap.Int64("id", nodeId), zap.Error(err))
		return protocol.Response(ctx, constants.ErrInvalidParams, nil)
	}

	llmSrv := node_llm.GetNodeHttpServiceInstance()
	// 检查节点是否在线
	ok, err := llmSrv.NodeCheckin(nodeId)
	if err != nil || !ok {
		return protocol.Response(ctx, constants.ErrNodeNotRegister.AppendErrors(err), "")
	}
	wallet, err := nus.getUserWalletBalance(req.UserID, req.WalletType)
	if err != nil {
		return protocol.Response(ctx, constants.ErrGetUserBalance.AppendErrors(err), "")
	}
	resp.Balance = wallet.Balance
	resp.WalletType = req.WalletType
	resp.WalletAddress = wallet.WalletAddress
	return protocol.Response(ctx, nil, resp)
}

// UpsetModelsInfos @Tags node/llm
// @Accept json
// @Produce json
// @Param request body requests.AddModelsInfoRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Failure 400 {object} responses.DefaultResponse
// @Failure 500 {object} responses.DefaultResponse
// @Router /node/llm/addModelsInfos [post]
func (nus *NodeUsersHttpService) UpsetModelsInfos(ctx echo.Context,
	req requests.AddModelsInfoRequest,
	resp responses.DefaultResponse) error {
	session := nus.dao.NewSession()
	defer session.Close()
	modelsInfos := make([]*models.ModelsInfo, len(req.ModelsInfos))
	for _, modelInfo := range req.ModelsInfos {
		where := &models.ModelsInfo{
			Name:       modelInfo.Name,
			ApiVersion: modelInfo.ApiVersion,
		}
		bean := &models.ModelsInfo{
			Name:        modelInfo.Name,
			ApiVersion:  modelInfo.ApiVersion,
			DeployName:  modelInfo.DeployName,
			InputPrice:  modelInfo.InputPrice,
			OutputPrice: modelInfo.OutputPrice,
			CachePrice:  modelInfo.CachePrice,
			Status:      modelInfo.Status,
			IsPrivate:   modelInfo.IsPrivate,
			OwnerId:     modelInfo.OwnerId,
			Address:     modelInfo.Address,
			LastUpdate:  time.Now().Unix(),
		}
		_, err := session.Upsert(where, bean)
		if err != nil {
			nus.logger.Info("databases error：", zap.Error(err))
			continue
		}
		modelsInfos = append(modelsInfos, bean)
	}
	return protocol.Response(ctx, nil, modelsInfos)
}

func (nus *NodeUsersHttpService) ListModelsInfos(ctx echo.Context,
	req requests.PageReq,
	resp responses.DefaultResponse) error {
	nus.logger.Info("ListModelsInfos called")

	// 从 header 中获取用户 ID
	_id, err := nus.getUserIdFromContext(ctx)
	if err != nil {
		nus.logger.Error("invalid user id", zap.Int64("_id", _id), zap.Error(err))
		return protocol.Response(ctx, constants.ErrInvalidParams, nil)
	}
	session := nus.dao.NewSession()
	defer session.Close()

	var modelsInfos []models.ModelsInfo
	// 尝试查询所有模型信息
	if err = session.Native().OrderBy(req.Sort).Limit(req.Limit, req.Skip).Find(&modelsInfos); err != nil {
		nus.logger.Error("query models info failed", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, modelsInfos)
}

// UpsetModelsProvider 添加模型供应商
// godoc
// @Summary 添加模型供应商
// @Description 添加模型供应商
// @Tags NodeLLM
// @Accept json
// @Produce json
// @Param request body requests.AddModelsProviderInfoRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/addModelsProvider [post]
func (nus *NodeUsersHttpService) UpsetModelsProvider(ctx echo.Context,
	req requests.AddModelsProviderInfoRequest,
	resp responses.DefaultResponse) error {
	nus.logger.Info("UpsetModelsProvider is Called", zap.Any("req", req))
	session := nus.dao.NewSession()
	defer session.Close()
	// 1. 检查模型服务商是否存在
	modelsProviders := make([]*models.ModelsProvider, len(req.ModelsProviderInfo))
	for _, v := range req.ModelsProviderInfo {
		where := &models.ModelsProvider{
			Name:     v.Name,
			Type:     v.Type,
			Endpoint: v.Endpoint,
		}
		bean := &models.ModelsProvider{
			Name:        v.Name,
			Type:        v.Type,
			Endpoint:    v.Endpoint,
			ApiType:     v.ApiType,
			ModelName:   v.ModelName,
			InputPrice:  v.InputPrice,
			OutputPrice: v.OutputPrice,
			CachePrice:  v.CachePrice,
			ApiKeys:     v.ApiKeys,
			LastUpdate:  time.Now().Unix(),
		}
		modelsProvider := &models.ModelsProvider{}
		_, err := session.Upsert(where, bean)
		if err != nil {
			nus.logger.Info("databases error：", zap.Error(err))
			continue
		}
		modelsProviders = append(modelsProviders, modelsProvider)
	}
	return protocol.Response(ctx, nil, modelsProviders)
}

func (nus *NodeUsersHttpService) ListModelsProviderInfos(ctx echo.Context,
	req requests.PageReq,
	resp responses.DefaultResponse) error {
	nus.logger.Info("ListModelsProviderInfos called")

	// 从 header 中获取用户 ID

	// 转换用户 ID 为整数
	userId, err := nus.getUserIdFromContext(ctx)
	if err != nil {
		nus.logger.Error("invalid user id", zap.Int64("userId", userId), zap.Error(err))
		return protocol.Response(ctx, constants.ErrInvalidParams, nil)
	}

	session := nus.dao.NewSession()
	defer session.Close()

	// 查询属于该用户的模型供应商信息
	var modelsProviders []models.ModelsProvider
	condiBean := &models.ModelsProvider{
		OwnerId: userId,
	}
	if err := session.Native().OrderBy(req.Sort).Limit(req.Limit, req.Skip).Find(&modelsProviders, condiBean); err != nil {
		nus.logger.Error("query models provider info failed", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, modelsProviders)
}

// MapModelsProviderInfoToNode 映射模型和模型供应商
// godoc
// @Summary 映射模型和模型供应商
// @Description 映射模型和模型供应商
// @Tags NodeLLM
// @Accept json
// @Produce json
// @Param request body requests.MapModelsProviderInfoToNodeRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/mapModelsProviderInfoToNode [post]
func (nus *NodeUsersHttpService) MapModelsProviderInfoToNode(ctx echo.Context,
	req requests.MapModelsProviderInfoToNodeRequest,
	resp responses.DefaultResponse) error {
	nus.logger.Info("call MapModelsProviderInfoToNode", zap.Any("req", req))
	session := nus.dao.NewSession()
	defer session.Close()
	llmSrv := node_llm.GetNodeHttpServiceInstance()
	// 检查Node和mode的合法师
	if !(llmSrv.CheckModelIdExists(req.ModelId) && llmSrv.CheckNodeIdExists(req.NodeId)) {
		return protocol.Response(ctx, nil, "node or model not exists")
	}
	// 设置条件
	where := &models.NodeModelsInfoMaps{
		NodeId:  req.NodeId,
		ModelId: req.ModelId,
	}
	count := 0
	for _, providerIds := range req.ProviderIds {
		nodeModelsInfoMap := &models.NodeModelsInfoMaps{
			NodeId:          req.NodeId,
			ModelId:         req.ModelId,
			ModelProviderId: providerIds,
			CreatedAt:       time.Now().Unix(),
			UpdatedAt:       time.Now().Unix(),
		}
		if _, err := session.Upsert(where, nodeModelsInfoMap); err != nil {
			nus.logger.Info("databases error:", zap.Error(err))
			continue
		} else {
			count++
		}
	}
	return protocol.Response(ctx, nil, fmt.Sprintf("update success %d", count))
}

// ListNodeModelsProviderInfos 获取节点模型供应商信息
// godoc
// @Summary 获取节点模型供应商信息
// @Description 获取节点模型供应商信息
// @Tags node/llm
// @Accept json
// @Produce json
// @Param request body requests.PageReq true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/listNodeModelsProviderInfos [post]
func (nus *NodeUsersHttpService) ListNodeModelsProviderInfos(ctx echo.Context,
	req requests.PageReq,
	resp responses.DefaultResponse) error {
	nus.logger.Info("ListNodeModelsProviderInfos called", zap.Any("req", req))
	session := nus.dao.NewSession()
	defer session.Close()
	// 默认排序
	if req.Sort == "" {
		req.Sort = "id asc"
	}
	nodeUserId, err := nus.getUserIdFromContext(ctx)
	if err != nil {
		nus.logger.Error("invalid user id", zap.Int64("nodeUserId", nodeUserId), zap.Error(err))
		return protocol.Response(ctx, constants.ErrInvalidParams, nil)
	}
	result, err := session.CallProcedure("ListNodeModelsProviderInfos",
		nodeUserId, req.Limit, req.Skip, req.Sort)
	if err != nil {
		nus.logger.Error("ListNodeModelsProviderInfos error:", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, result)
}

// UpsetNodeInfos 添加/更新节点信息
// godoc
// @Summary 添加/更新节点信息
// @Description 添加/更新节点信息
// @Tags node
// @Accept json
// @Produce json
// @Param request body requests.UpsetNodeInfoRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/upsetNodeInfos [post]
func (nus *NodeUsersHttpService) UpsetNodeInfos(ctx echo.Context,
	req requests.UpsetNodeInfoRequest,
	resp responses.DefaultResponse) error {
	nus.logger.Info("call UpsetNodeInfos", zap.Any("req", req))
	session := nus.dao.NewSession()
	defer session.Close()
	if req.Name == "" {
		req.Name = uuid.GenString(12)
	}
	where := &models.Nodes{
		Name:    req.Name,
		OwnerId: req.NodeUserId,
	}
	bean := &models.Nodes{
		Name:         req.Name,
		OwnerId:      req.NodeUserId,
		AccessKey:    req.AccessKey,
		SecurityKey:  req.SecretKey,
		LastupdateAt: time.Now().Unix(),
		Domain:       req.Domain,
		CompanyId:    req.CompanyId,
	}
	if _, err := session.Upsert(where, bean); err != nil {
		nus.logger.Info("databases error:", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}
	return protocol.Response(ctx, nil, "success")
}

// ListNodeInfos 获取模型信息
// godoc
// @Summary 获取模型信息
// @Description 获取模型信息
// @Tags node
// @Accept json
// @Produce json
// @Param request body requests.ListNodeInfoRequest true "请求参数"
// @Success 200 {object} responses.DefaultResponse
// @Router /node/llm/listNodeInfos [post]
func (nus *NodeUsersHttpService) ListNodeInfos(ctx echo.Context,
	req requests.ListNodeInfoRequest,
	resp responses.DefaultResponse) error {
	nus.logger.Info("List Node Infos call", zap.Any("req", req))
	session := nus.dao.NewSession()
	defer session.Close()
	// 默认排序
	if req.PageInfo.Sort == "" {
		req.PageInfo.Sort = "node_id asc"
	}
	result, err := session.CallProcedure("ListNodeUserNodeInfos",
		req.NodeUserId, req.PageInfo.Limit, req.PageInfo.Skip, req.PageInfo.Sort)
	if err != nil {
		nus.logger.Error("ListNodeUserNodeInfos error:", zap.Error(err))
		return protocol.Response(ctx, constants.ErrInternalServer.AppendErrors(err), nil)
	}

	return protocol.Response(ctx, nil, result)
}
