package account

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"gorm.io/gorm"
)

// HtTmWholesalerApiAccounts ねっぱんとskyticketを結びつけるアカウント管理用
type HtTmWholesalerApiAccounts struct {
	WholesalerAccountID int64          `gorm:"primaryKey;autoIncrement:true" json:"wholesaler_account_id"`
	WholesalerID        int64          `json:"wholesaler_id"`
	Name                string         `json:"name"`
	LoginID             string         `json:"login_id"`
	LoginPWEnc          string         `json:"login_pw_enc"`
	Username            string         `json:"username"`
	PasswordEnc         string         `json:"password_enc"`
	AppEnv              string         `json:"app_env"`
	TokenKey            string         `json:"token_key"`
	Urls                string         `json:"urls"`
	PropertyID          int64          `json:"property_id"`
	CompanyCode         string         `json:"company_code"`
	CompanyName         string         `json:"company_name"`
	PropertyName        string         `json:"property_name"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at"`
	common.Times
}

// IAccountTemaUsecase アカウント関連のてまusecaseのインターフェース
type IAccountTemaUsecase interface {
	FetchConnectUser(request *CheckConnectInput) bool
}

// IAccountTemaRepository アカウント関連のてまrepositoryのインターフェース
type IAccountTemaRepository interface {
	// FetchConnectUser 削除されていないてまの接続用アカウントを1件取得
	FetchConnectUser(propertyID int64) (*HtTmWholesalerApiAccounts, error)
	// 指定の連動IDが他の施設IDで紐づけている数を取得
	FetchCountOtherConnectedID(propertyID int64, username string) (int, error)
	// UpsertConnectUser tema連携用の情報を作成、もしくは更新
	UpsertConnectUser(apiAccount *HtTmWholesalerApiAccounts) error
}
