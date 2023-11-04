package account

// HtTmConnectUserRaku2s らく通とskyticketを結びつけるアカウント管理用
type HtTmConnectUserRaku2s struct {
	ConnectUserRaku2ID  int64  `gorm:"primaryKey;autoIncrement:true" json:"connect_user_raku2_id"`
	UserIDEnc           string `json:"user_id_enc"`
	PasswordEnc         string `json:"password_enc"`
	PropertyID          int64  `json:"property_id"`
	StopFlag            bool   `json:"stop_flag"`
}

// IAccountRaku2Usecase アカウント関連のらく通usecaseのインターフェース
type IAccountRaku2Usecase interface {
	FetchConnectUser(request *CheckConnectInput) bool
}

// IAccountRaku2Repository アカウント関連のらく通repositoryのインターフェース
type IAccountRaku2Repository interface {
	// FetchConnectUser らく通の停止されていない接続用アカウントを1件取得
	FetchConnectUser(propertyID int64) (*HtTmConnectUserRaku2s, error)
	// UpsertConnectUser らく通連携用のアカウントを作成もしくは更新
	UpsertConnectUser(userIDEnc string, passwordEnc string, propertyID int64) error
}
