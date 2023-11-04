package infra

import (
	"os"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"gorm.io/gorm"
)

type accountTemaRepository struct {
	db *gorm.DB
}

// NewAccountTemaRepository インスタンス生成
func NewAccountTemaRepository(db *gorm.DB) account.IAccountTemaRepository {
	return &accountTemaRepository{
		db: db,
	}
}

// FetchConnectUser 削除されていないてまの接続用アカウントを1件取得
func (a *accountTemaRepository) FetchConnectUser(propertyID int64) (*account.HtTmWholesalerApiAccounts, error) {
	result := &account.HtTmWholesalerApiAccounts{}
	err := a.db.
		Model(&account.HtTmWholesalerApiAccounts{}).
		Where("property_id = ?", propertyID).
		Where("wholesaler_id = ?", utils.WholesalerIDTema).
		Where("app_env = ?", os.Getenv("APP_ENV")).
		Where("deleted_at IS NULL").
		First(&result).Error
	return result, err
}

// FetchConnectedUser 指定の連動IDが他の施設IDで紐づけている数を取得
func (a *accountTemaRepository) FetchCountOtherConnectedID(propertyID int64, username string) (int, error) {
	result := 0
	err := a.db.
		Model(&account.HtTmWholesalerApiAccounts{}).
		Select("count(*)").
		Where("property_id != ?", propertyID).
		Where("username = ?", username).
		Where("deleted_at IS NULL").
		Scan(&result).Error
	return result, err
}

// UpsertConnectUser tema連携用の情報を作成、もしくは更新
func (a *accountTemaRepository) UpsertConnectUser(apiAccount *account.HtTmWholesalerApiAccounts) error {
	assignData := map[string]interface{}{
		"wholesaler_id": utils.WholesalerIDTema,
		"name":          apiAccount.Name,
		"login_id":      apiAccount.LoginID,
		"login_pw_enc":  apiAccount.LoginPWEnc,
		"username":      apiAccount.Username,
		"password_enc":  apiAccount.PasswordEnc,
		"app_env":       os.Getenv("APP_ENV"),
		"urls":          apiAccount.Urls,
		"property_id":   apiAccount.PropertyID,
		"updated_at":    time.Now(),
	}
	if apiAccount.CreatedAt.IsZero() == false {
		assignData["createdAt"] = time.Now()
	}
	return a.db.Model(&account.HtTmWholesalerApiAccounts{}).
		Where("property_id = ?", apiAccount.PropertyID).
		Where("wholesaler_Id = ?", utils.WholesalerIDTema).
		Where("app_env = ?", os.Getenv("APP_ENV")).
		Where("deleted_at IS NULL").
		Assign(assignData).
		FirstOrCreate(&account.HtTmWholesalerApiAccounts{}).
		Error
}
