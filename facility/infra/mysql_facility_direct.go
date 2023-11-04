package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// facilityDirectRepository 直仕入れ施設関連repository
type facilityDirectRepository struct {
	db *gorm.DB
}

// NewFacilityDirectRepository インスタンス生成
func NewFacilityDirectRepository(db *gorm.DB) facility.IFacilityDirectRepository {
	return &facilityDirectRepository{
		db: db,
	}
}

// FetchAllFacilities 施設情報を複数件取得
func (f *facilityDirectRepository) FetchAllFacilities(propertyIDs []int64, wholesalerID int64) ([]facility.InitFacilityOutput, error) {
	result := &[]facility.InitFacilityOutput{}
	err := f.db.
		Select("properties.property_id, property_use_wholesalers.wholesaler_id, properties.name, property_directs.disp_priority, properties.state_province_name, properties.city, properties.line_1, properties.line_2, properties.line_3").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_langs AS property_langs ON properties.property_id = property_langs.property_id").
		Joins("LEFT JOIN ht_tm_property_directs AS property_directs ON properties.property_id = property_directs.property_id").
		Where("property_use_wholesalers.wholesaler_id = ?", wholesalerID).
		Where("property_langs.lang_cd = ?", "ja-JP").
		Where("properties.property_id IN ?", propertyIDs).
		Find(result).Error
	return *result, err
}

// FetchPropertyDetail 施設詳細情報を1件取得
func (f *facilityDirectRepository) FetchPropertyDetail(propertyID int64) (*facility.DetailOutput, error) {
	result := &facility.DetailOutput{}
	err := f.db.Select("property_directs.*").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_directs AS property_directs ON properties.property_id = property_directs.property_id").
		Where("property_use_wholesalers.wholesaler_id = ?", utils.WholesalerIDDirect).
		Where("property_directs.lang_cd = ?", "ja-JP").
		Where("properties.property_id = ?", propertyID).
		First(result).Error
	return result, err
}

// FetchAmenities 施設に紐づくアメニティを複数件取得
func (f *facilityDirectRepository) FetchAmenities(propertyID int64) (*[]facility.HtTmPropertyAmenityDirects, error) {
	result := &[]facility.HtTmPropertyAmenityDirects{}
	err := f.db.
		Table("ht_tm_property_amenity_directs AS amenity").
		Joins("INNER JOIN ht_tm_property_directs_use_amenities as use_amenity ON amenity.property_amenity_id = use_amenity.property_amenity_id").
		Where("amenity.lang_cd = ?", "ja-JP").
		Where("use_amenity.property_id = ?", propertyID).
		Find(result).Error
	return result, err
}

// FirstOrCreate 施設情報の作成＆1件取得
func (f *facilityDirectRepository) FirstOrCreate(propertyID int64) (*facility.HtTmPropertyDirects, error) {
	response := &facility.HtTmPropertyDirects{}
	err := f.db.
		FirstOrCreate(response, facility.HtTmPropertyDirects{
			PropertyID: propertyID,
			LangCd:     "ja-JP", // 元となったhotelリポジトリでハードコードだった箇所
		}).Error
	return response, err
}

// UpdateDispPriority サイト公開フラグの更新
func (f *facilityDirectRepository) UpdateDispPriority(propertyID int64, dispPriority bool) error {
	return f.db.Model(&facility.HtTmPropertyDirects{}).
		Where("property_id = ?", propertyID).
		Updates(map[string]interface{}{
			"disp_priority": dispPriority,
			"updated_at":    time.Now(),
		}).Error
}

// UpsertPropertyDirect 施設の詳細情報を更新・新規作成
func (f *facilityDirectRepository) UpsertPropertyDirect(upsertData *facility.HtTmPropertyDirects) error {
	assignData := map[string]interface{}{
		"property_id":                    upsertData.PropertyID,
		"lang_cd":                        "ja-JP",
		"fee_mandatory":                  upsertData.FeeMandatory,
		"fee_optional":                   upsertData.FeeOptional,
		"description_amenity":            upsertData.DescriptionAmenity,
		"description_attractions":        upsertData.DescriptionAttractions,
		"description_business_amenities": upsertData.DescriptionBusinessAmenities,
		"description_dining":             upsertData.DescriptionDining,
		"description_location":           upsertData.DescriptionLocation,
		"description_headline":           upsertData.DescriptionHeadline,
		"description_rooms":              upsertData.DescriptionRooms,
		"checkin_begin":                  upsertData.CheckinBegin,
		"checkin_end":                    upsertData.CheckinEnd,
		"checkout":                       upsertData.Checkout,
		"instructions":                   upsertData.Instructions,
		"special_instructions":           upsertData.SpecialInstructions,
		"policy_know_before_you_go":      upsertData.PolicyKnowBeforeYouGo,
	}
	return f.db.Model(&facility.HtTmPropertyDirects{}).
		Where("property_id = ?", upsertData.PropertyID).
		Where("lang_cd = ?", "ja-JP").
		Assign(assignData).
		FirstOrCreate(&facility.HtTmPropertyDirects{}).
		Error
}

// ClearPropertyAmenity 施設に紐づくアメニティを全て削除
func (f *facilityDirectRepository) ClearPropertyAmenity(propertyID int64) error {
	return f.db.Delete(&facility.HtTmPropertyDirectsUseAmenity{}, "property_id = ?", propertyID).Error
}

// CreatePropertyAmenity 施設に紐づくアメニティを作成
func (f *facilityDirectRepository) CreatePropertyAmenity(facilities []facility.HtTmPropertyDirectsUseAmenity) error {
	return f.db.Create(&facilities).Error
}

// FetchAllAmenities 施設アメニティを複数件取得
func (f *facilityDirectRepository) FetchAllAmenities() (*[]facility.HtTmPropertyAmenityDirects, error) {
	res := &[]facility.HtTmPropertyAmenityDirects{}
	err := f.db.Model(&facility.HtTmPropertyAmenityDirects{}).Where("lang_cd = ?", "ja-JP").Find(res).Error
	return res, err
}
