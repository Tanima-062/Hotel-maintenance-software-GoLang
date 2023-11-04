package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// facilityRaku2Repository らく通施設関連repository
type facilityRaku2Repository struct {
	db *gorm.DB
}

// NewFacilityRaku2Repository インスタンス生成
func NewFacilityRaku2Repository(db *gorm.DB) facility.IFacilityRaku2Repository {
	return &facilityRaku2Repository{
		db: db,
	}
}

// FetchAllFacilities 施設情報を複数件取得
func (f *facilityRaku2Repository) FetchAllFacilities(propertyIDs []int64, wholesalerID int64) ([]facility.InitFacilityOutput, error) {
	result := &[]facility.InitFacilityOutput{}
	err := f.db.
		Select("properties.property_id, property_use_wholesalers.wholesaler_id, properties.name, property_raku2s.disp_priority, properties.state_province_name, properties.city, properties.line_1, properties.line_2, properties.line_3").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_langs AS property_langs ON properties.property_id = property_langs.property_id").
		Joins("LEFT JOIN ht_tm_property_raku2s AS property_raku2s ON properties.property_id = property_raku2s.property_id").
		Where("property_use_wholesalers.wholesaler_id = ?", wholesalerID).
		Where("property_langs.lang_cd = ?", "ja-JP").
		Where("properties.property_id IN ?", propertyIDs).
		Find(result).Error
	return *result, err
}

// FetchPropertyDetail 施設詳細情報を1件取得
func (f *facilityRaku2Repository) FetchPropertyDetail(propertyID int64) (*facility.DetailOutput, error) {
	result := &facility.DetailOutput{}
	err := f.db.Select("property_raku2s.*").
		Table("ht_tm_properties as properties").
		Joins("INNER JOIN ht_tm_property_use_wholesalers AS property_use_wholesalers ON properties.property_id = property_use_wholesalers.property_id").
		Joins("LEFT JOIN ht_tm_property_raku2s AS property_raku2s ON properties.property_id = property_raku2s.property_id").
		Where("property_use_wholesalers.wholesaler_id = ?", utils.WholesalerIDRaku2).
		Where("property_raku2s.lang_cd = ?", "ja-JP").
		Where("properties.property_id = ?", propertyID).
		First(result).Error
	return result, err
}

// FetchAmenities 施設に紐づくアメニティを複数件取得
func (f *facilityRaku2Repository) FetchAmenities(propertyID int64) (*[]facility.HtTmPropertyAmenityRaku2s, error) {
	result := &[]facility.HtTmPropertyAmenityRaku2s{}
	err := f.db.
		Table("ht_tm_property_amenity_raku2s AS amenity").
		Joins("INNER JOIN ht_tm_property_raku2s_use_amenities as use_amenity ON amenity.property_amenity_id = use_amenity.property_amenity_id").
		Where("amenity.lang_cd = ?", "ja-JP").
		Where("use_amenity.property_id = ?", propertyID).
		Find(result).Error
	return result, err
}

// FirstOrCreate 施設情報の作成＆1件取得
func (f *facilityRaku2Repository) FirstOrCreate(propertyID int64) (*facility.HtTmPropertyRaku2s, error) {
	response := &facility.HtTmPropertyRaku2s{}
	err := f.db.
		FirstOrCreate(response, facility.HtTmPropertyRaku2s{
			PropertyID: propertyID,
			LangCd:     "ja-JP", // 元となったhotelリポジトリでハードコードだった箇所
		}).Error
	return response, err
}

// UpdateDispPriority サイト公開フラグの更新
func (f *facilityRaku2Repository) UpdateDispPriority(propertyID int64, dispPriority bool) error {
	return f.db.Model(&facility.HtTmPropertyRaku2s{}).
		Where("property_id = ?", propertyID).
		Updates(map[string]interface{}{
			"disp_priority": dispPriority,
			"updated_at":    time.Now(),
		}).Error
}

// UpsertPropertyRaku2 施設の詳細情報を更新・新規作成
func (f *facilityRaku2Repository) UpsertPropertyRaku2(upsertData *facility.HtTmPropertyRaku2s) error {
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
	return f.db.Model(&facility.HtTmPropertyRaku2s{}).
		Where("property_id = ?", upsertData.PropertyID).
		Where("lang_cd = ?", "ja-JP").
		Assign(assignData).
		FirstOrCreate(&facility.HtTmPropertyRaku2s{}).
		Error
}

// ClearPropertyAmenity 施設に紐づくアメニティを全て削除
func (f *facilityRaku2Repository) ClearPropertyAmenity(propertyID int64) error {
	return f.db.Delete(&facility.HtTmPropertyRaku2sUseAmenity{}, "property_id = ?", propertyID).Error
}

// CreatePropertyAmenity 施設に紐づくアメニティを作成
func (f *facilityRaku2Repository) CreatePropertyAmenity(facilities []facility.HtTmPropertyRaku2sUseAmenity) error {
	return f.db.Create(&facilities).Error
}

// FetchAllAmenities 施設アメニティを複数件取得
func (f *facilityRaku2Repository) FetchAllAmenities() (*[]facility.HtTmPropertyAmenityRaku2s, error) {
	res := &[]facility.HtTmPropertyAmenityRaku2s{}
	err := f.db.Model(&facility.HtTmPropertyAmenityRaku2s{}).Where("lang_cd = ?", "ja-JP").Find(res).Error
	return res, err
}

// FetchConnectedUser 指定の連動IDが他の施設IDで紐づけている数を取得
func (f *facilityRaku2Repository) FetchCountOtherConnectedID(propertyID int64, userIDEnc string) (int, error) {
	result := 0
	err := f.db.Select("count(*)").
		Table("ht_tm_connect_user_raku2s").
		Where("user_id_enc = ?", userIDEnc).
		Where("property_id != ?", propertyID).
		Where("stop_flag = ?", 0).
		Scan(&result).Error
	return result, err
}