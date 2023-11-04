package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"gorm.io/gorm"
)

type accountNeppanUsecase struct {
	ARepository account.IAccountNeppanRepository
}

// NewAccountNeppanUsecase インスタンス生成
func NewAccountNeppanUsecase(db *gorm.DB) account.IAccountNeppanUsecase {
	return &accountNeppanUsecase{
		ARepository: infra.NewAccountNeppanRepository(db),
	}
}

// FetchConnectUser ねっぱん連携用ユーザが登録済みかどうか
func (a *accountNeppanUsecase) FetchConnectUser(request *account.CheckConnectInput) bool {
	res, err := a.ARepository.FetchConnectUser(request.PropertyID)
	if err != nil {
		return false
	}
	return res.PropertyID != 0
}
