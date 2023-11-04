package facility

import "github.com/Adventureinc/hotel-hm-api/src/common"

// HtTmPropertyTemas てま施設情報テーブル
type HtTmPropertyTemas struct {
	PropTemaID                   int64  `gorm:"primaryKey;autoIncrement:true" json:"prop_tema_id,omitempty"`
	PropertyID                   int64  `json:"property_id"`
	TemaPropertyID               string `json:"tema_property_id,omitempty"`
	LangCd                       string `json:"lang_cd,omitempty"`
	Rank                         int    `json:"rank"`
	DispFlag                     bool   `json:"disp_flag" gorm:"default:true"`
	DispPriority                 bool   `json:"disp_priority"`
	CheckinBegin                 string `json:"checkin_begin,omitempty"`
	CheckinEnd                   string `json:"checkin_end,omitempty"`
	Instructions                 string `json:"instructions,omitempty"`
	SpecialInstructions          string `json:"special_instructions,omitempty"`
	MinAge                       int    `json:"min_age"`
	Checkout                     string `json:"checkout,omitempty"`
	FeeMandatory                 string `json:"fee_mandatory,omitempty"`
	FeeOptional                  string `json:"fee_optional,omitempty"`
	PolicyKnowBeforeYouGo        string `json:"policy_know_before_you_go,omitempty"`
	HeroImage                    string `json:"hero_image,omitempty"`
	Category                     string `json:"category,omitempty"`
	DescriptionAmenity           string `json:"description_amenity,omitempty"`
	DescriptionAttractions       string `json:"description_attractions,omitempty"`
	DescriptionBusinessAmenities string `json:"description_business_amenities,omitempty"`
	DescriptionDining            string `json:"description_dining,omitempty"`
	DescriptionLocation          string `json:"description_location,omitempty"`
	DescriptionHeadline          string `json:"description_headline,omitempty"`
	DescriptionRooms             string `json:"description_rooms,omitempty"`
	// CancelPenaltyJSON            string `json:"cancel_penalty_json,omitempty"`
	common.Times `gorm:"embedded"`
}

// HtTmPropertyAmenityTemas 施設単位で設定できるアメニティ群のテーブル
type HtTmPropertyAmenityTemas struct {
	PropertyAmenityID string `json:"property_amenity_id"`
	LangCd            string `json:"lang_cd"`
	AmenityName       string `json:"amenity_name"`
	Description       string `json:"description"`
	common.Times      `gorm:"embedded"`
}

// HtTmPropertyTemaUseAmenity 施設がどのアメニティを設定しているか紐付けるテーブル
type HtTmPropertyTemaUseAmenity struct {
	UseAmenityID      int64  `gorm:"primaryKey;autoIncrement:true" json:"use_amenity_id"`
	PropertyAmenityID string `json:"property_amenity_id"`
	PropertyID        int64  `json:"property_id"`
	common.Times      `gorm:"embedded"`
}

// IFacilityTemaRepository てま施設関連のrepositoryのインターフェース
type IFacilityTemaRepository interface {
	// FetchAllFacilities 施設情報を複数件取得
	FetchAllFacilities(propertyIDs []int64, wholesalerID int64) ([]InitFacilityOutput, error)
	// FetchPropertyDetail 施設詳細情報を1件取得
	FetchPropertyDetail(propertyID int64) (*DetailOutput, error)
	// FetchAmenities 施設に紐づくアメニティを複数件取得
	FetchAmenities(propertyID int64) (*[]HtTmPropertyAmenityTemas, error)
	// UpdateDispPriority サイト公開フラグの更新
	UpdateDispPriority(propertyID int64, dispPriority bool) error
	// UpsertPropertyTema 施設の詳細情報を更新・新規作成
	UpsertPropertyTema(upsertData *HtTmPropertyTemas) error
	// ClearPropertyAmenity 施設に紐づくアメニティを全て削除
	ClearPropertyAmenity(propertyID int64) error
	// CreatePropertyAmenity 施設に紐づくアメニティを作成
	CreatePropertyAmenity(facilities []HtTmPropertyTemaUseAmenity) error
}

// IFacilityUsecaseを手間いらず用に拡張したインターフェース
type IFacilityTemaUsecase interface {
	IFacilityUsecase
	IsRegisteredConnect(request *SaveBaseInfoInput) (bool, error)
}