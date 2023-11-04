package account

// TlApiAccpuntUrl TLのアカウント用
type TlApiAccpuntUrl struct {
	PropertySearch string `json:"property_search"`
	PlanSearch     string `json:"plan_search"`
	Reserve        string `json:"reserve"`
}

// IAccountTlRepository アカウント関連のてまAPIのインターフェース
type IAccountTlRepository interface {
	// FetchAPIAccount 削除されていないTLの接続用アカウントを1件取得
	FetchAPIAccount(propertyID int64) (*HtTmWholesalerApiAccounts, error)
}
