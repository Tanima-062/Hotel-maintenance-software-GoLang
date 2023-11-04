package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// facilityTlRepository TL施設関連repository
type facilityTlRepository struct {
	db *gorm.DB
}

// NewFacilityTlRepository インスタンス生成
func NewFacilityTlRepository(db *gorm.DB) facility.IFacilityTlRepository {
	return &facilityTlRepository{
		db: db,
	}
}

// FetchAllFacilities 施設情報を複数件取得
func (f *facilityTlRepository) FetchAllFacilities(propertyIDs []int64, wholesalerID int64) ([]facility.InitFacilityOutput, error) {
	result := &[]facility.InitFacilityOutput{}
	err := f.db.
		Select("properties.property_id, property_use_wholesalers.wholesaler_id, properties.name, property_tls.disp_priority, properties.state_province_name, properties.city, properties.line_1, properties.line_2, properties.line_3").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_langs AS property_langs ON properties.property_id = property_langs.property_id").
		Joins("LEFT JOIN ht_tm_property_tls AS property_tls ON properties.property_id = property_tls.property_id").
		Where("property_use_wholesalers.wholesaler_id = ?", wholesalerID).
		Where("property_langs.lang_cd = ?", "ja-JP").
		Where("properties.property_id IN ?", propertyIDs).
		Find(result).Error
	return *result, err
}

// FetchPropertyByPropertyID 施設情報を1件取得
func (f *facilityTlRepository) FetchPropertyByPropertyID(propertyID int64) (*facility.HtTmPropertyTls, error) {
	result := &facility.HtTmPropertyTls{}
	err := f.db.Model(&facility.HtTmPropertyTls{}).
		Where("property_id = ?", propertyID).
		First(result).Error
	return result, err
}

// FetchPropertyDetail 施設詳細情報を1件取得
func (f *facilityTlRepository) FetchPropertyDetail(propertyID int64) (*facility.DetailOutput, error) {
	result := &facility.DetailOutput{}
	err := f.db.Select("property_tls.*").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_langs AS property_langs ON properties.property_id = property_langs.property_id").
		Joins("LEFT JOIN ht_tm_property_tls AS property_tls ON properties.property_id = property_tls.property_id").
		Where("property_use_wholesalers.wholesaler_id = ?", utils.WholesalerIDTl).
		Where("property_langs.lang_cd = ?", "ja-JP").
		Where("properties.property_id = ?", propertyID).
		First(result).Error
	return result, err
}

// FetchAmenities 施設に紐づくアメニティを複数件取得
func (f *facilityTlRepository) FetchAmenities(propertyID int64) (*[]facility.HtTmPropertyAmenityTls, error) {
	result := &[]facility.HtTmPropertyAmenityTls{}
	err := f.db.
		Table("ht_tm_property_amenity_tls AS amenity").
		Joins("INNER JOIN ht_tm_property_tls_use_amenities as use_amenity ON amenity.property_amenity_id = use_amenity.property_amenity_id").
		Where("amenity.lang_cd = ?", "ja-JP").
		Where("use_amenity.property_id = ?", propertyID).
		Find(result).Error
	return result, err
}

// UpdateDispPriority サイト公開フラグの更新
func (f *facilityTlRepository) UpdateDispPriority(propertyID int64, dispPriority bool) error {
	return f.db.Model(&facility.HtTmPropertyTls{}).
		Where("property_id = ?", propertyID).
		Updates(map[string]interface{}{
			"disp_priority": dispPriority,
			"updated_at":    time.Now(),
		}).Error
}

// UpsertProperty 施設の詳細情報を更新・新規作成
func (f *facilityTlRepository) UpsertPropertyTl(upsertData *facility.HtTmPropertyTls) error {
	assignData := map[string]interface{}{
		"property_id":                    upsertData.PropertyID,
		"lang_cd":                        "ja-JP",
		"checkin_begin":                  upsertData.CheckinBegin,
		"checkin_end":                    upsertData.CheckinEnd,
		"checkout":                       upsertData.Checkout,
		"instructions":                   upsertData.Instructions,
		"special_instructions":           upsertData.SpecialInstructions,
		"fee_mandatory":                  upsertData.FeeMandatory,
		"fee_optional":                   upsertData.FeeOptional,
		"policy_know_before_you_go":      upsertData.PolicyKnowBeforeYouGo,
		"description_amenity":            upsertData.DescriptionAmenity,
		"description_attractions":        upsertData.DescriptionAttractions,
		"description_business_amenities": upsertData.DescriptionBusinessAmenities,
		"description_dining":             upsertData.DescriptionDining,
		"description_location":           upsertData.DescriptionLocation,
		"description_headline":           upsertData.DescriptionHeadline,
		"description_rooms":              upsertData.DescriptionRooms,
	}
	return f.db.Model(&facility.HtTmPropertyTls{}).
		Where("property_id = ?", upsertData.PropertyID).
		Assign(assignData).
		FirstOrCreate(&facility.HtTmPropertyTls{}).
		Error
}

// ClearPropertyAmenity 施設に紐づくアメニティを全て削除
func (f *facilityTlRepository) ClearPropertyAmenity(propertyID int64) error {
	return f.db.Delete(&facility.HtTmPropertyTlsUseAmenity{}, "property_id = ?", propertyID).Error
}

// CreatePropertyAmenity 施設に紐づくアメニティを作成
func (f *facilityTlRepository) CreatePropertyAmenity(facilities []facility.HtTmPropertyTlsUseAmenity) error {
	return f.db.Create(&facilities).Error
}

// FetchAllAmenities 施設アメニティを複数件取得
func (f *facilityTlRepository) FetchAllAmenities() (*[]facility.HtTmPropertyAmenityTls, error) {
	res := &[]facility.HtTmPropertyAmenityTls{}
	err := f.db.Model(&facility.HtTmPropertyAmenityTls{}).Where("lang_cd = ?", "ja-JP").Find(res).Error
	return res, err
}
