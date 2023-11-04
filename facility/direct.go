package facility

import "github.com/Adventureinc/hotel-hm-api/src/common"

// HtTmPropertyDirects 直仕入れ施設情報テーブル
type HtTmPropertyDirects struct {
	PropDirectID                 int64  `gorm:"primaryKey;autoIncrement:true" json:"prop_direct_id,omitempty"`
	PropertyID                   int64  `json:"property_id"`
	DirectPropertyCategoryID     int64  `json:"direct_property_category_id,omitempty"`
	LangCd                       string `json:"lang_cd,omitempty"`
	DescriptionRooms             string `json:"description_rooms,omitempty"`
	FeeMandatory                 string `json:"fee_mandatory,omitempty"`
	FeeOptional                  string `json:"fee_optional,omitempty"`
	DescriptionAmenity           string `json:"description_amenity,omitempty"`
	DescriptionDining            string `json:"description_dining,omitempty"`
	DescriptionLocation          string `json:"description_location,omitempty"`
	DescriptionHeadline          string `json:"description_headline,omitempty"`
	DescriptionBusinessAmenities string `json:"description_business_amenities,omitempty"`
	DescriptionAttractions       string `json:"description_attractions,omitempty"`
	DispFlag                     bool   `json:"disp_flag" gorm:"default:true"`
	DispPriority                 bool   `json:"disp_priority"`
	CheckinBegin                 string `json:"checkin_begin,omitempty"`
	CheckinEnd                   string `json:"checkin_end,omitempty"`
	Checkout                     string `json:"checkout,omitempty"`
	Instructions                 string `json:"instructions,omitempty"`
	SpecialInstructions          string `json:"special_instructions,omitempty"`
	PolicyKnowBeforeYouGo        string `json:"policy_know_before_you_go,omitempty"`
	CancelPenaltyJSON            string `json:"cancel_penalty_json"`
	common.Times                 `gorm:"embedded"`
}

// HtTmPropertyAmenityDirects 施設単位で設定できるアメニティ群のテーブル
type HtTmPropertyAmenityDirects struct {
	PropertyAmenityID int64  `gorm:"primaryKey;autoIncrement:true" json:"property_amenity_id"`
	LangCd            string `json:"lang_cd"`
	AmenityName       string `json:"amenity_name"`
	Description       string `json:"description"`
	common.Times      `gorm:"embedded"`
}

// HtTmPropertyDirectsUseAmenity 施設がどのアメニティを設定しているか紐付けるテーブル
type HtTmPropertyDirectsUseAmenity struct {
	UseAmenityID      int64  `gorm:"primaryKey;autoIncrement:true" json:"use_amenity_id"`
	PropertyAmenityID string `json:"property_amenity_id"`
	PropertyID        int64  `json:"property_id"`
	common.Times      `gorm:"embedded"`
}

// IFacilityDirectRepository 直仕入れ施設関連のrepositoryのインターフェース
type IFacilityDirectRepository interface {
	// FetchAllFacilities 施設情報を複数件取得
	FetchAllFacilities(propertyIDs []int64, wholesalerID int64) ([]InitFacilityOutput, error)
	// FetchPropertyDetail 施設詳細情報を1件取得
	FetchPropertyDetail(propertyID int64) (*DetailOutput, error)
	// FetchAmenities 施設に紐づくアメニティを複数件取得
	FetchAmenities(propertyID int64) (*[]HtTmPropertyAmenityDirects, error)
	// FirstOrCreate 施設情報の作成＆1件取得
	FirstOrCreate(propertyID int64) (*HtTmPropertyDirects, error)
	// UpdateDispPriority サイト公開フラグの更新
	UpdateDispPriority(propertyID int64, dispPriority bool) error
	// UpsertPropertyDirect 施設の詳細情報を更新・新規作成
	UpsertPropertyDirect(upsertData *HtTmPropertyDirects) error
	// ClearPropertyAmenity 施設に紐づくアメニティを全て削除
	ClearPropertyAmenity(propertyID int64) error
	// CreatePropertyAmenity 施設に紐づくアメニティを作成
	CreatePropertyAmenity(facilities []HtTmPropertyDirectsUseAmenity) error
	// FetchAllAmenities 施設アメニティを複数件取得
	FetchAllAmenities() (*[]HtTmPropertyAmenityDirects, error)
}
