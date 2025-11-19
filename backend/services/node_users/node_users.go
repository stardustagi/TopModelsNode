package nodeUsers

import (
	"fmt"

	"github.com/stardustagi/TopModelsNode/models"
	"go.uber.org/zap"
)

func (nus *NodeUsersHttpService) getUserWalletBalance(userId int64, walletType string) (*models.UserWallet, error) {
	session := nus.dao.NewSession()
	defer session.Close()

	user := &models.UserWallet{
		UserId:     userId,
		WalletType: walletType,
	}
	ok, err := session.FindOne(user)
	if err != nil {
		nus.logger.Error("Failed to find user by ID", zap.Error(err), zap.Int64("userID", userId))
		return nil, err
	}
	if !ok {
		nus.logger.Debug("User not found", zap.Int64("userID", userId))
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
