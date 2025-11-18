package nodeUsers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stardustagi/TopLib/libs/jwt"
	"github.com/stardustagi/TopModelsNode/constants"
	"github.com/stardustagi/TopModelsNode/models"
	"github.com/stardustagi/TopModelsNode/protocol/responses"
	"go.uber.org/zap"
)

const (
	// JWTSecretKey JWT密钥（生产环境应该从配置文件读取）
	JWTSecretKey = "your-secret-key-change-in-production"
)

// NodeUserMailEncodeToken 节点用户邮箱加密
func (nus *NodeUsersHttpService) nodeUserMailEncodeToken(mail, password string, fixedSalt string) (string, error) {
	return jwt.EncryptWithFixedSalt(password, constants.ObtenationIterations, mail, fixedSalt)
}

// NodeUserMailDecodeToken 节点用户邮箱解密
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
func (nus *NodeUsersHttpService) generateJWTToken(userID int64, email string) (string, time.Time, error) {
	// 解析过期时间
	duration, err := time.ParseDuration(constants.NodeUserJwtExpire)
	if err != nil {
		nus.logger.Error("parse jwt expire duration failed", zap.Error(err))
		duration = 720 * time.Hour // 默认30天
	}

	expireTime := time.Now().Add(duration)

	// 简单的token生成（生产环境应该使用标准的JWT库）
	tokenData := fmt.Sprintf("%d:%s:%d", userID, email, expireTime.Unix())
	hash := md5.New()
	hash.Write([]byte(tokenData))
	token := hex.EncodeToString(hash.Sum(nil))

	return token, expireTime, nil
}

// 辅助方法：缓存用户Token
func (nus *NodeUsersHttpService) cacheUserToken(userID int64, token string, expireTime time.Time) error {
	nus.mu.Lock()
	defer nus.mu.Unlock()

	key := constants.NodeUserTokenKey(userID)
	duration := time.Until(expireTime)
	durationStr := fmt.Sprintf("%ds", int(duration.Seconds()))

	err := nus.rds.Set(nus.ctx, key, []byte(token), durationStr)
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
