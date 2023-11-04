package infra

import (
	"os"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"gorm.io/gorm"
)

type accountTlRepository struct {
	db *gorm.DB
}

// NewAccountTlRepository インスタンス生成
func NewAccountTlRepository(db *gorm.DB) account.IAccountTlRepository {
	return &accountTlRepository{
		db: db,
	}
}

// FetchAPIAccount 削除されていないTLの接続用アカウントを1件取得
func (a *accountTlRepository) FetchAPIAccount(propertyID int64) (*account.HtTmWholesalerApiAccounts, error) {
	result := &account.HtTmWholesalerApiAccounts{}
	err := a.db.
		Model(&account.HtTmWholesalerApiAccounts{}).
		Where("property_id = ?", propertyID).
		Where("wholesaler_id = ?", utils.WholesalerIDTl).
		Where("app_env = ?", os.Getenv("APP_ENV")).
		First(&result).Error
	return result, err
}
