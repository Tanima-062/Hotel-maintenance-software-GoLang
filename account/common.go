package account

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmHotelManager hotel managerのアカウント管理用テーブル
type HtTmHotelManager struct {
	HotelManagerID        int64     `gorm:"primaryKey;autoIncrement:true" json:"hotel_manager_id"`
	ClientCompanyID       int64     `json:"client_company_id,omitempty"`
	PropertyID            int64     `json:"property_id"`
	WholesalerID          int64     `json:"wholesaler_id"`
	FirstNameEnc          string    `json:"first_name_enc"`
	LastNameEnc           string    `json:"last_name_enc"`
	EmailEnc              string    `json:"email_enc"`
	UsernameEnc           string    `json:"username_enc,omitempty"`
	PasswordEnc           string    `json:"password_enc,omitempty"`
	IsPrimary             bool      `json:"is_primary,omitempty"`
	MasterEditFlg         bool      `json:"master_edit_flg,omitempty"`
	SettlementNeedFlg     bool      `json:"settlement_need_flg,omitempty"`
	ClosingDate           int       `json:"closing_date,omitempty"`
	PaymentDateMonthLater int       `json:"payment_date_month_later,omitempty"`
	DelFlg                bool      `json:"del_flg,omitempty"`
	LoginedAt             time.Time `gorm:"type:time" json:"logined_at"`
	LoginFailedNum        int       `json:"login_failed_num,omitempty"`
	RememberToken         string    `json:"remember_token,omitempty"`
	common.Times          `gorm:"embedded"`
}

// HtTmHotelManagerTokens HMアカウントに紐づくトークン管理用テーブル
type HtTmHotelManagerTokens struct {
	TokenID        int64  `gorm:"primaryKey;autoIncrement:true" json:"token_id,omitempty"`
	HotelManagerID int64  `json:"hotel_manager_id,omitempty"`
	APIToken       string `json:"api_token,omitempty"`
	common.Times   `gorm:"embedded"`
}

// ClaimParam トークンの属性情報
type ClaimParam struct {
	HotelManagerID int64  `json:"hotel_manager_id,omitempty"`
	APIToken       string `json:"api_token,omitempty"`
}

// LoginInput Loginの入力
type LoginInput struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

// TokenOutput token返却用
type TokenOutput struct {
	APIToken string `json:"api_token,omitempty"`
}

// ChangePasswordInput パスワード変更の入力
type ChangePasswordInput struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

// CheckConnectInput 接続ユーザ確認の入力
type CheckConnectInput struct {
	PropertyID   int64 `json:"property_id" param:"propertyId"`
	WholesalerID int64 `json:"wholesaler_id" query:"wholesaler_id"`
}

// IsParentAccountInput 親アカウント確認の入力
type IsParentAccountInput struct {
	HotelManagerID int64 `json:"hotel_manager_id" param:"hotelManagerId"`
}

// IAccountUsecase アカウント関連のusecaseのインターフェース
type IAccountUsecase interface {
	Login(LoginInput *LoginInput) (string, error)
	Logout(claimParam *ClaimParam) error
	CheckToken(claimParam *ClaimParam) (string, error)
	FetchDetail(claimParam *ClaimParam) (*HtTmHotelManager, error)
	FetchHMUserByToken(claimParam *ClaimParam) (HtTmHotelManager, error)
	ChangePassword(request *ChangePasswordInput) error
	IsParentAccount(hotelManagerID int64) bool
}

// IAccountRepository アカウント関連のrepositoryのインターフェース
type IAccountRepository interface {
	// FetchHMUserByLoginInfo ログインユーザとパスワードに合致するアカウントを1件取得
	FetchHMUserByLoginInfo(hmUser *HtTmHotelManager) HtTmHotelManager
	// FetchHMUserByToken トークンからアカウントを1件取得
	FetchHMUserByToken(claimParam *ClaimParam) (HtTmHotelManager, error)
	// SaveLoginInfo ログイン日時とトークンを更新
	SaveLoginInfo(hmUser *HtTmHotelManager, claimParam *ClaimParam, newToken string) error
	// UpdatePassword パスワードを更新
	UpdatePassword(hotelManagerID int64, password string) error
	// DeleteAPIToken トークン削除
	DeleteAPIToken(claimParam *ClaimParam) error
	// FetchHMUserByPropertyID property_idでHMアカウントを1件取得
	FetchHMUserByPropertyID(propertyID int64, wholesalerID int64) (HtTmHotelManager, error)
	// FetchOne PRIMARY KEYでアカウントを1件取得
	FetchOne(hotelManagerID int64) (HtTmHotelManager, error)
}
