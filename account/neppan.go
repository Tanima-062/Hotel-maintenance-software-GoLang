package account

// HtTmConnectUserNeppans ねっぱんとskyticketを結びつけるアカウント管理用
type HtTmConnectUserNeppans struct {
	ConnectUserNeppanID int64  `gorm:"primaryKey;autoIncrement:true" json:"connect_user_neppan_id"`
	UserIDEnc           string `json:"user_id_enc"`
	PasswordEnc         string `json:"password_enc"`
	PropertyID          int64  `json:"property_id"`
	StopFlag            bool   `json:"stop_flag"`
}

// IAccountNeppanUsecase アカウント関連のねっぱんusecaseのインターフェース
type IAccountNeppanUsecase interface {
	FetchConnectUser(request *CheckConnectInput) bool
}

// IAccountNeppanRepository アカウント関連のねっぱんrepositoryのインターフェース
type IAccountNeppanRepository interface {
	// FetchConnectUser ねっぱんの停止されていない接続用アカウントを1件取得
	FetchConnectUser(propertyID int64) (*HtTmConnectUserNeppans, error)
	// UpsertConnectUser ねっぱん連携用のアカウントを作成もしくは更新
	UpsertConnectUser(userIDEnc string, passwordEnc string, propertyID int64) error
}
