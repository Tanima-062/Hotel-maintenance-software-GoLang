package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// facilityRepository 施設関連repository
type facilityRepository struct {
	db *gorm.DB
}

// NewFacilityRepository インスタンス生成
func NewFacilityRepository(db *gorm.DB) facility.IFacilityRepository {
	return &facilityRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (f *facilityRepository) TxStart() (*gorm.DB, error) {
	tx := f.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (f *facilityRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (f *facilityRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchAllClientCompanies 子施設のPropertyIDを複数件取得
func (f *facilityRepository) FetchAllClientCompanies(hotelManagerID int64) ([]facility.HtTmProperties, error) {
	result := &[]facility.HtTmProperties{}
	err := f.db.Select("properties.property_id").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_hotel_managers AS hm ON properties.client_company_id = hm.client_company_id AND hm.hotel_manager_id = ?", hotelManagerID).
		Find(result).Error
	return *result, err
}

// FetchAllFacilitiesByPropertyID propetyIDに紐づく施設を取得
func (f *facilityRepository) FetchAllFacilitiesByPropertyID(propertyIDs []int64) ([]facility.InitFacilityOutput, error) {
	result := &[]facility.InitFacilityOutput{}
	err := f.db.
		Select("properties.property_id, property_use_wholesalers.wholesaler_id, properties.name, properties.state_province_name, properties.city, properties.line_1, properties.line_2, properties.line_3").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_langs AS property_langs ON properties.property_id = property_langs.property_id").
		Where("property_langs.lang_cd = ?", "ja-JP").
		Where("properties.property_id IN ?", propertyIDs).
		Find(result).Error
	return *result, err
}

// FetchProperty 施設情報を1件取得
func (f *facilityRepository) FetchProperty(propertyID int64) (*facility.HtTmProperties, error) {
	result := &facility.HtTmProperties{}
	err := f.db.
		Model(&facility.HtTmProperties{}).
		Where("ht_tm_properties.property_id = ?", propertyID).
		First(result).Error
	return result, err
}

// FetchCategory 施設のカテゴリー情報を1件取得
func (f *facilityRepository) FetchCategory(propertyID int64) (*facility.HtTmCategorise, error) {
	result := &facility.HtTmCategorise{}
	err := f.db.
		Select("ht_tm_categorise.name").
		Table("ht_tm_categorise").
		Joins("JOIN ht_tm_property_categorise ON ht_tm_categorise.category_id = ht_tm_property_categorise.category_id").
		Where("ht_tm_property_categorise.property_id = ?", propertyID).
		Where("ht_tm_categorise.lang_cd = ?", "ja-JP").
		First(result).Error
	return result, err
}

// UpdateProperty 施設の基本情報を更新
func (f *facilityRepository) UpdateProperty(property *facility.HtTmProperties) error {
	return f.db.Model(&facility.HtTmProperties{}).
		Where("property_id = ?", property.PropertyID).
		Updates(map[string]interface{}{
			"name":                property.Name,
			"postal_code":         property.PostalCode,
			"state_province_name": property.StateProvinceName,
			"city":                property.City,
			"line_1":              property.Line1,
			"line_2":              property.Line2,
			"line_3":              property.Line3,
			"phone":               property.Phone,
			"fax":                 property.Fax,
			"updated_at":          time.Now(),
		}).Error
}

// UpsertPropertyLangsBase 施設の基本情報を更新（Langs）
func (f *facilityRepository) UpsertPropertyLangsBase(property *facility.HtTmPropertyLangs) error {
	assignData := map[string]interface{}{
		"property_id":         property.PropertyID,
		"lang_cd":             "ja-JP",
		"name":                property.Name,
		"postal_code":         property.PostalCode,
		"state_province_name": property.StateProvinceName,
		"city":                property.City,
		"line_1":              property.Line1,
		"line_2":              "",
		"line_3":              "",
		"updated_at":          time.Now(),
	}
	if property.CreatedAt.IsZero() == false {
		assignData["createdAt"] = time.Now()
	}
	return f.db.Model(&facility.HtTmPropertyLangs{}).
		Where("property_id = ?", property.PropertyID).
		Where("lang_cd = ?", "ja-JP").
		Assign(assignData).
		FirstOrCreate(&facility.HtTmPropertyLangs{}).Error
}

// UpsertPropertyLangsDetail 施設の詳細情報を更新
func (f *facilityRepository) UpsertPropertyLangsDetail(upsertData *facility.HtTmPropertyLangs) error {
	assignData := map[string]interface{}{
		"property_id":               upsertData.PropertyID,
		"lang_cd":                   "ja-JP",
		"checkin_begin":             upsertData.CheckinBegin,
		"checkin_end":               upsertData.CheckinEnd,
		"checkout":                  upsertData.Checkout,
		"instructions":              upsertData.Instructions,
		"special_instructions":      upsertData.SpecialInstructions,
		"policy_know_before_you_go": upsertData.PolicyKnowBeforeYouGo,
	}
	return f.db.Model(&facility.HtTmPropertyLangs{}).
		Where("property_id = ?", upsertData.PropertyID).
		Where("lang_cd = ?", "ja-JP").
		Assign(assignData).
		FirstOrCreate(&facility.HtTmPropertyLangs{}).
		Error
}
