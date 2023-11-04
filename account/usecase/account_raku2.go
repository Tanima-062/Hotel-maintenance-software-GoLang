package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"gorm.io/gorm"
)

type accountRaku2Usecase struct {
	ARepository account.IAccountRaku2Repository
}

// NewAccountRaku2Usecase インスタンス生成
func NewAccountRaku2Usecase(db *gorm.DB) account.IAccountRaku2Usecase {
	return &accountRaku2Usecase{
		ARepository: infra.NewAccountRaku2Repository(db),
	}
}

// FetchConnectUser らく通連携用ユーザが登録済みかどうか
func (a *accountRaku2Usecase) FetchConnectUser(request *account.CheckConnectInput) bool {
	res, err := a.ARepository.FetchConnectUser(request.PropertyID)
	if err != nil {
		return false
	}
	return res.PropertyID != 0
}
