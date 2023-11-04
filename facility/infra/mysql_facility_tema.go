package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// facilityTemaRepository てま施設関連repository
type facilityTemaRepository struct {
	db *gorm.DB
}

// NewFacilityTemaRepository インスタンス生成
func NewFacilityTemaRepository(db *gorm.DB) facility.IFacilityTemaRepository {
	return &facilityTemaRepository{
		db: db,
	}
}

// FetchAllFacilities 施設情報を複数件取得
func (f *facilityTemaRepository) FetchAllFacilities(propertyIDs []int64, wholesalerID int64) ([]facility.InitFacilityOutput, error) {
	result := &[]facility.InitFacilityOutput{}
	err := f.db.
		Select("properties.property_id, property_use_wholesalers.wholesaler_id, properties.name, property_temas.disp_priority, properties.state_province_name, properties.city, properties.line_1, properties.line_2, properties.line_3").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_langs AS property_langs ON properties.property_id = property_langs.property_id").
		Joins("LEFT JOIN ht_tm_property_temas AS property_temas ON properties.property_id = property_temas.property_id").
		Where("property_use_wholesalers.wholesaler_id = ?", wholesalerID).
		Where("property_langs.lang_cd = ?", "ja-JP").
		Where("properties.property_id IN ?", propertyIDs).
		Find(result).Error
	return *result, err
}

// FetchPropertyDetail 施設詳細情報を1件取得
func (f *facilityTemaRepository) FetchPropertyDetail(propertyID int64) (*facility.DetailOutput, error) {
	result := &facility.DetailOutput{}
	err := f.db.Select("property_temas.*").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_langs AS property_langs ON properties.property_id = property_langs.property_id").
		Joins("LEFT JOIN ht_tm_property_temas AS property_temas ON properties.property_id = property_temas.property_id").
		Where("property_use_wholesalers.wholesaler_id = ?", utils.WholesalerIDTema).
		Where("property_langs.lang_cd = ?", "ja-JP").
		Where("properties.property_id = ?", propertyID).
		First(result).Error
	return result, err
}

// FetchAmenities 施設に紐づくアメニティを複数件取得
func (f *facilityTemaRepository) FetchAmenities(propertyID int64) (*[]facility.HtTmPropertyAmenityTemas, error) {
	result := &[]facility.HtTmPropertyAmenityTemas{}
	err := f.db.
		Table("ht_tm_property_amenity_tls AS amenity").
		Joins("INNER JOIN ht_tm_property_tema_use_amenities as use_amenity ON amenity.property_amenity_id = use_amenity.property_amenity_id").
		Where("amenity.lang_cd = ?", "ja-JP").
		Where("use_amenity.property_id = ?", propertyID).
		Find(result).Error
	return result, err
}

// UpdateDispPriority サイト公開フラグの更新
func (f *facilityTemaRepository) UpdateDispPriority(propertyID int64, dispPriority bool) error {
	return f.db.Model(&facility.HtTmPropertyTemas{}).
		Where("property_id = ?", propertyID).
		Updates(map[string]interface{}{
			"disp_priority": dispPriority,
			"updated_at":    time.Now(),
		}).Error
}

// UpsertPropertyTema 施設の詳細情報を更新・新規作成
func (f *facilityTemaRepository) UpsertPropertyTema(upsertData *facility.HtTmPropertyTemas) error {
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
	return f.db.Model(&facility.HtTmPropertyTemas{}).
		Where("property_id = ?", upsertData.PropertyID).
		Assign(assignData).
		FirstOrCreate(&facility.HtTmPropertyTemas{}).
		Error
}

// ClearPropertyAmenity 施設に紐づくアメニティを全て削除
func (f *facilityTemaRepository) ClearPropertyAmenity(propertyID int64) error {
	return f.db.Delete(&facility.HtTmPropertyTemaUseAmenity{}, "property_id = ?", propertyID).Error
}

// CreatePropertyAmenity 施設に紐づくアメニティを作成
func (f *facilityTemaRepository) CreatePropertyAmenity(facilities []facility.HtTmPropertyTemaUseAmenity) error {
	return f.db.Create(&facilities).Error
}