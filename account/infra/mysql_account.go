package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"gorm.io/gorm"
)

// accountRepository アカウント関連のDB
type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository インスタンス生成
func NewAccountRepository(db *gorm.DB) account.IAccountRepository {
	return &accountRepository{
		db: db,
	}
}

// FetchHMUserByLoginInfo ログインユーザとパスワードに合致するアカウントを1件取得
func (a *accountRepository) FetchHMUserByLoginInfo(hmUser *account.HtTmHotelManager) account.HtTmHotelManager {
	result := account.HtTmHotelManager{}
	a.db.Where("username_enc = ? AND password_enc = ? AND del_flg = 0", hmUser.UsernameEnc, hmUser.PasswordEnc).First(&result)
	return result
}

// FetchHMUserByToken トークンからアカウントを1件取得
func (a *accountRepository) FetchHMUserByToken(claimParam *account.ClaimParam) (account.HtTmHotelManager, error) {
	result := account.HtTmHotelManager{}
	err := a.db.Table("ht_tm_hotel_managers").
		Joins("INNER JOIN ht_tm_hotel_manager_tokens ON ht_tm_hotel_managers.hotel_manager_id = ht_tm_hotel_manager_tokens.hotel_manager_id").
		Where("ht_tm_hotel_manager_tokens.api_token = ?", claimParam.APIToken).
		Where("ht_tm_hotel_managers.hotel_manager_id = ? AND ht_tm_hotel_managers.del_flg = 0", claimParam.HotelManagerID).
		First(&result).Error
	return result, err
}

// SaveLoginInfo ログイン日時とトークンを更新
func (a *accountRepository) SaveLoginInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, newToken string) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		if err := a.db.Model(account.HtTmHotelManager{}).
			Where("hotel_manager_id = ?", hmUser.HotelManagerID).
			Updates(map[string]interface{}{
				"logined_at": hmUser.LoginedAt,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return err
		}
		return a.db.
			Where(account.HtTmHotelManagerTokens{HotelManagerID: claimParam.HotelManagerID, APIToken: claimParam.APIToken}).
			Assign(account.HtTmHotelManagerTokens{
				APIToken: newToken,
				Times:    common.Times{UpdatedAt: time.Now()},
			}).
			FirstOrCreate(&account.HtTmHotelManagerTokens{}).Error
	})
}

// UpdatePassword パスワードを更新
func (a *accountRepository) UpdatePassword(hotelManagerID int64, password string) error {
	return a.db.Model(account.HtTmHotelManager{}).
		Where("hotel_manager_id = ?", hotelManagerID).
		Updates(map[string]interface{}{
			"password_enc": password,
			"updated_at":   time.Now(),
		}).Error
}

// DeleteAPIToken トークン削除
func (a *accountRepository) DeleteAPIToken(claimParam *account.ClaimParam) error {
	return a.db.Delete(&account.HtTmHotelManagerTokens{}, "hotel_manager_id = ? AND api_token = ?", claimParam.HotelManagerID, claimParam.APIToken).Error
}

// FetchHMUserByPropertyID property_idでHMアカウントを1件取得
func (a *accountRepository) FetchHMUserByPropertyID(propertyID int64, wholesalerID int64) (account.HtTmHotelManager, error) {
	result := account.HtTmHotelManager{}
	err := a.db.
		Where("property_id = ? AND del_flg = 0", propertyID).
		Where("wholesaler_id = ?", wholesalerID).
		First(&result).Error
	return result, err
}

// FetchOne PRIMARY KEYでアカウントを1件取得
func (a *accountRepository) FetchOne(hotelManagerID int64) (account.HtTmHotelManager, error) {
	result := account.HtTmHotelManager{}
	err := a.db.
		Where("hotel_manager_id = ? AND del_flg = 0", hotelManagerID).
		First(&result).Error
	return result, err
}
