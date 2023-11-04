package facility

import "github.com/Adventureinc/hotel-hm-api/src/common"

// HtTmPropertyNeppans ねっぱん施設情報テーブル
type HtTmPropertyNeppans struct {
	PropNeppanID                 int64  `gorm:"primaryKey;autoIncrement:true" json:"prop_neppan_id,omitempty"`
	PropertyID                   int64  `json:"property_id"`
	NeppanPropertyCategoryID     int64  `json:"neppan_property_category_id,omitempty"`
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

// HtTmPropertyAmenityNeppans 施設単位で設定できるアメニティ群のテーブル
type HtTmPropertyAmenityNeppans struct {
	PropertyAmenityID int64  `gorm:"primaryKey;autoIncrement:true" json:"property_amenity_id"`
	LangCd            string `json:"lang_cd"`
	AmenityName       string `json:"amenity_name"`
	Description       string `json:"description"`
	common.Times      `gorm:"embedded"`
}

// HtTmPropertyNeppansUseAmenity 施設がどのアメニティを設定しているか紐付けるテーブル
type HtTmPropertyNeppansUseAmenity struct {
	UseAmenityID      int64  `gorm:"primaryKey;autoIncrement:true" json:"use_amenity_id"`
	PropertyAmenityID string `json:"property_amenity_id"`
	PropertyID        int64  `json:"property_id"`
	common.Times      `gorm:"embedded"`
}

// IFacilityNeppanRepository ねっぱん施設関連のrepositoryのインターフェース
type IFacilityNeppanRepository interface {
	// FetchAllFacilities 施設情報を複数件取得
	FetchAllFacilities(propertyIDs []int64, wholesalerID int64) ([]InitFacilityOutput, error)
	// FetchPropertyDetail 施設詳細情報を1件取得
	FetchPropertyDetail(propertyID int64) (*DetailOutput, error)
	// FetchAmenities 施設に紐づくアメニティを複数件取得
	FetchAmenities(propertyID int64) (*[]HtTmPropertyAmenityNeppans, error)
	// FirstOrCreate 施設情報の作成＆1件取得
	FirstOrCreate(propertyID int64) (*HtTmPropertyNeppans, error)
	// UpdateDispPriority サイト公開フラグの更新
	UpdateDispPriority(propertyID int64, dispPriority bool) error
	// UpsertPropertyNeppan 施設の詳細情報を更新・新規作成
	UpsertPropertyNeppan(upsertData *HtTmPropertyNeppans) error
	// ClearPropertyAmenity 施設に紐づくアメニティを全て削除
	ClearPropertyAmenity(propertyID int64) error
	// CreatePropertyAmenity 施設に紐づくアメニティを作成
	CreatePropertyAmenity(facilities []HtTmPropertyNeppansUseAmenity) error
	// FetchAllAmenities 施設アメニティを複数件取得
	FetchAllAmenities() (*[]HtTmPropertyAmenityNeppans, error)
	// 指定の連動IDが他の施設IDで紐づけている数を取得
	FetchCountOtherConnectedID(propertyID int64, userIDEnc string) (int, error)
}

// IFacilityUsecaseをねっぱん用に拡張したインターフェース
type IFacilityNeppanUsecase interface {
	IFacilityUsecase
	IsRegisteredConnect(request *SaveBaseInfoInput) (bool, error)
}