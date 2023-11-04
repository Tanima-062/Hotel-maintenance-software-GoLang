package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"gorm.io/gorm"
)

type accountRaku2Repository struct {
	db *gorm.DB
}

// NewAccountRaku2Repository インスタンス生成
func NewAccountRaku2Repository(db *gorm.DB) account.IAccountRaku2Repository {
	return &accountRaku2Repository{
		db: db,
	}
}

// FetchConnectUser らく通の停止されていない接続用アカウントを1件取得
func (a *accountRaku2Repository) FetchConnectUser(propertyID int64) (*account.HtTmConnectUserRaku2s, error) {
	result := &account.HtTmConnectUserRaku2s{}
	err := a.db.
		Model(&account.HtTmConnectUserRaku2s{}).
		Where("property_id = ?", propertyID).
		Where("stop_flag = 0").
		First(&result).Error
	return result, err
}

// UpsertConnectUser らく通連携用のアカウントを作成もしくは更新
func (a *accountRaku2Repository) UpsertConnectUser(userIDEnc string, passwordEnc string, propertyID int64) error {
	assignData := map[string]interface{}{
		"user_id_enc":  userIDEnc,
		"password_enc": passwordEnc,
		"property_id":  propertyID,
	}
	return a.db.Model(&account.HtTmConnectUserRaku2s{}).
		Where("property_id = ?", propertyID).
		Assign(assignData).
		FirstOrCreate(&account.HtTmConnectUserRaku2s{}).
		Error
}
