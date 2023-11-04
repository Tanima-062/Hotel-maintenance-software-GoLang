package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"gorm.io/gorm"
)

type accountTemaUsecase struct {
	ARepository account.IAccountTemaRepository
}

// NewAccountTemaUsecase インスタンス生成
func NewAccountTemaUsecase(db *gorm.DB) account.IAccountTemaUsecase {
	return &accountTemaUsecase{
		ARepository: infra.NewAccountTemaRepository(db),
	}
}

// FetchConnectUser てま接続用ユーザが登録済みかどうか
func (a *accountTemaUsecase) FetchConnectUser(request *account.CheckConnectInput) bool {
	res, err := a.ARepository.FetchConnectUser(request.PropertyID)
	if err != nil {
		return false
	}
	return res.PropertyID != 0
}
