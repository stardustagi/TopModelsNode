package nodeUsers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stardustagi/TopLib/libs/jwt"
	"github.com/stardustagi/TopLib/libs/redis"
	"github.com/stardustagi/TopLib/libs/uuid"
	"github.com/stardustagi/TopModelsNode/constants"
	"github.com/stardustagi/TopModelsNode/models"
	"github.com/stardustagi/TopModelsNode/protocol/responses"
	"go.uber.org/zap"
)

const (
	// JWTSecretKey JWT密钥（生产环境应该从配置文件读取）
	JWTSecretKey = "your-secret-key-change-in-production"
)

// nodeUserMailEncodeToken 节点用户邮箱加密
func (nus *NodeUsersHttpService) nodeUserMailEncodeToken(mail, password string, fixedSalt string) (string, error) {
	return jwt.EncryptWithFixedSalt(password, constants.ObtenationIterations, mail, fixedSalt)
}

// nodeUserMailDecodeToken 节点用户邮箱解密
func (nus *NodeUsersHttpService) nodeUserMailDecodeToken(inPassword, dbPassword string, fixedSalt string) (string, error) {
	return jwt.DecryptWithFixedSalt(inPassword, constants.ObtenationIterations, dbPassword, fixedSalt)
}

// 辅助方法：生成盐值
func (nus *NodeUsersHttpService) generateSalt(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		nus.logger.Error("generate salt failed", zap.Error(err))
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

// 辅助方法：密码加密
func (nus *NodeUsersHttpService) hashPassword(password, salt string) string {
	hash := md5.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// 辅助方法：生成JWT Token
func (nus *NodeUsersHttpService) generateJWTToken(userId int64) (string, string, error) {
	// 生成JWT token
	token := uuid.GetUuidString()
	jwtKey := fmt.Sprintf("%s-%s-%d", constants.AppName, constants.AppVersion, userId)
	jwtString := jwt.JWTEncrypt(fmt.Sprintf("%d", userId), token, jwtKey)

	// 存储到Redis
	key := constants.NodeUserTokenKey(userId)
	err := nus.rds.Set(nus.ctx, key, []byte(token), constants.NodeUserTokenExpire)
	if err != nil {
		nus.logger.Error("Failed to store Token in Redis", zap.Error(err), zap.String("key", key))
		return "", "", fmt.Errorf("failed to store JWT: %w", err)
	}
	key = constants.NodeUserJwtKey(userId)
	err = nus.rds.Set(nus.ctx, key, []byte(jwtString), constants.NodeUserJwtExpire)
	if err != nil {
		nus.logger.Error("Failed to store JWT in Redis", zap.Error(err), zap.String("key", key))
		return "", "", fmt.Errorf("failed to store JWT: %w", err)
	}
	return token, jwtString, nil
}

// 辅助方法：缓存用户Token
func (nus *NodeUsersHttpService) cacheUserToken(userID int64, token, expireTime string) error {
	nus.mu.Lock()
	defer nus.mu.Unlock()

	key := constants.NodeUserTokenKey(userID)

	err := nus.rds.Set(nus.ctx, key, []byte(token), expireTime)
	if err != nil {
		nus.logger.Error("cache user token failed", zap.Error(err), zap.Int64("userID", userID))
		return err
	}
	return nil
}

// 辅助方法：缓存用户信息
func (nus *NodeUsersHttpService) cacheUserInfo(user *models.NodeUsers) error {
	nus.userInfosMu.Lock()
	defer nus.userInfosMu.Unlock()

	key := constants.NodeUserAccessTokenKey(user.Id)

	// 将用户信息序列化为JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		nus.logger.Error("marshal user info failed", zap.Error(err))
		return err
	}

	err = nus.rds.Set(nus.ctx, key, userJSON, constants.NodeUserTokenExpire)
	if err != nil {
		nus.logger.Error("cache user info failed", zap.Error(err), zap.Int64("userID", user.Id))
		return err
	}
	return nil
}

// 辅助方法：从Redis获取用户信息
func (nus *NodeUsersHttpService) getUserInfoFromCache(userID int64) (*models.NodeUsers, error) {
	nus.userInfosMu.RLock()
	defer nus.userInfosMu.RUnlock()

	key := constants.NodeUserAccessTokenKey(userID)
	data, err := nus.rds.Get(nus.ctx, key)
	if err != nil {
		nus.logger.Debug("get user info from cache failed", zap.Error(err), zap.Int64("userID", userID))
		return nil, err
	}

	var user models.NodeUsers
	if err := json.Unmarshal(data, &user); err != nil {
		nus.logger.Error("unmarshal user info failed", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

// 辅助方法：清除用户缓存
func (nus *NodeUsersHttpService) clearUserCache(userID int64) error {
	nus.userInfosMu.Lock()
	defer nus.userInfosMu.Unlock()

	// 清除Token缓存
	tokenKey := constants.NodeUserTokenKey(userID)
	if _, err := nus.rds.Del(nus.ctx, tokenKey); err != nil {
		nus.logger.Error("delete user token cache failed", zap.Error(err), zap.Int64("userID", userID))
		return err
	}

	// 清除JWT缓存
	jwtKey := constants.NodeUserJwtKey(userID)
	if _, err := nus.rds.Del(nus.ctx, jwtKey); err != nil {
		nus.logger.Error("delete user jwt cache failed", zap.Error(err), zap.Int64("userID", userID))
		return err
	}

	// 清除用户信息缓存
	infoKey := constants.NodeUserAccessTokenKey(userID)
	if _, err := nus.rds.Del(nus.ctx, infoKey); err != nil {
		nus.logger.Error("delete user info cache failed", zap.Error(err), zap.Int64("userID", userID))
		return err
	}

	nus.logger.Info("user cache cleared", zap.Int64("userID", userID))
	return nil
}

// 辅助方法：转换为响应格式
func (nus *NodeUsersHttpService) convertToUserInfoResponse(user *models.NodeUsers) responses.UserInfoResponse {
	return responses.UserInfoResponse{
		UserID:     user.Id,
		Username:   user.Email,
		Email:      user.Email,
		Status:     user.IsActive,
		CreateTime: time.Unix(user.CreatedAt, 0),
		UpdateTime: time.Unix(user.LastUpdate, 0),
	}
}

func (nus *NodeUsersHttpService) nodeUserActive(nodeUserId int64, activeCode string) (bool, error) {
	nus.logger.Info("NodeUserActive called", zap.Int64("nodeUserId", nodeUserId), zap.String("activeCode", activeCode))
	key := constants.NodeUserMailVerifyKey(nodeUserId)
	storedCode, err := nus.rds.Get(nus.ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		nus.logger.Info("获取激活码失败", zap.Int64("nodeUserId", nodeUserId), zap.Error(err))
		return false, err
	}
	if string(storedCode) != activeCode {
		nus.logger.Info("激活码不匹配", zap.Int64("nodeUserId", nodeUserId))
		return false, fmt.Errorf("激活码不匹配")
	}
	nodeUser := &models.NodeUsers{
		Id: nodeUserId,
	}
	session := nus.dao.NewSession()
	ok, err := session.FindOne(nodeUser)
	if err != nil {
		nus.logger.Error("nodeUser查询失败", zap.Error(err), zap.Int64("nodeUserId", nodeUserId))
		return false, err
	}
	if !ok {
		nus.logger.Info("节点用户不存在", zap.Int64("nodeUserId", nodeUserId))
		return false, fmt.Errorf("节点用户不存在: %d", nodeUserId)
	}
	// 更新节点用户状态为激活
	nodeUser.IsActive = 1
	nodeUser.LastUpdate = time.Now().Unix()

	_, err = session.UpdateById(nodeUserId, nodeUser)
	if err != nil {
		nus.logger.Error("更新节点用户状态失败", zap.Error(err), zap.Int64("nodeUserId", nodeUserId))
		return false, err
	}
	nus.logger.Info("节点用户激活成功", zap.Int64("nodeUserId", nodeUserId))
	return true, nil
}

func (nus *NodeUsersHttpService) getUserIdFromContext(ctx echo.Context) (int64, error) {
	_id := ctx.Request().Header.Get("Id")
	if _id == "" {
		return 0, fmt.Errorf("用户ID不能为空")
	}
	return strconv.ParseInt(_id, 10, 64)
}
