package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"gorm.io/gorm"
)

type accountNeppanRepository struct {
	db *gorm.DB
}

// NewAccountNeppanRepository インスタンス生成
func NewAccountNeppanRepository(db *gorm.DB) account.IAccountNeppanRepository {
	return &accountNeppanRepository{
		db: db,
	}
}

// FetchConnectUser ねっぱんの停止されていない接続用アカウントを1件取得
func (a *accountNeppanRepository) FetchConnectUser(propertyID int64) (*account.HtTmConnectUserNeppans, error) {
	result := &account.HtTmConnectUserNeppans{}
	err := a.db.
		Model(&account.HtTmConnectUserNeppans{}).
		Where("property_id = ?", propertyID).
		Where("stop_flag = 0").
		First(&result).Error
	return result, err
}

// UpsertConnectUser ねっぱん連携用のアカウントを作成もしくは更新
func (a *accountNeppanRepository) UpsertConnectUser(userIDEnc string, passwordEnc string, propertyID int64) error {
	assignData := map[string]interface{}{
		"user_id_enc":  userIDEnc,
		"password_enc": passwordEnc,
		"property_id":  propertyID,
	}
	return a.db.Model(&account.HtTmConnectUserNeppans{}).
		Where("property_id = ?", propertyID).
		Assign(assignData).
		FirstOrCreate(&account.HtTmConnectUserNeppans{}).
		Error
}
