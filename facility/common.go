package facility

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

const (
	ParentPropertyId = 0
)

// HtTmProperties 施設情報テーブル
type HtTmProperties struct {
	PropertyID          int64     `gorm:"primaryKey;autoIncrement:true" json:"property_id,omitempty"`
	ClientCompanyID     int64     `json:"client_company_id,omitempty"`
	Name                string    `json:"name,omitempty"`
	BilCountryCode      string    `json:"bil_country_code,omitempty"`
	PostalCode          string    `json:"postal_code,omitempty"`
	StateProvinceCode   string    `json:"state_province_code,omitempty"`
	StateProvinceName   string    `json:"state_province_name,omitempty"`
	City                string    `json:"city,omitempty"`
	Line1               string    `json:"line_1,omitempty" gorm:"column:line_1"`
	Line2               string    `json:"line_2,omitempty" gorm:"column:line_2"`
	Line3               string    `json:"line_3,omitempty" gorm:"column:line_3"`
	Lat                 float64   `json:"lat"`
	Lng                 float64   `json:"lng"`
	PlaceID             string    `json:"place_id,omitempty"`
	Phone               string    `json:"phone"`
	Fax                 string    `json:"fax"`
	DeleteFlg           bool      `json:"delete_flg,omitempty"`
	Deleted             time.Time `json:"deleted,omitempty"`
	IsLocked            bool      `json:"is_locked,omitempty"`
	SnakeName           string    `json:"snake_name,omitempty"`
	AirportDistance     float64   `json:"airport_distance,omitempty"`
	AirportDistanceUnit string    `json:"airport_distance_unit,omitempty"`
	Imported            time.Time `json:"imported,omitempty"`
	CreateOperatorID    int64     `json:"create_operator_id,omitempty"`
	UpdateOperatorID    int64     `json:"update_operator_id,omitempty"`
	common.Times        `gorm:"embedded"`
}

// HtTmPropertyLangs 施設の言語別情報テーブル
type HtTmPropertyLangs struct {
	TmPropertyLangID      int64  `gorm:"primaryKey;autoIncrement:true" json:"tm_property_lang_id,omitempty"`
	StateProvinceName     string `json:"state_province_name,omitempty"`
	PropertyID            int64  `json:"property_id,omitempty"`
	LangCd                string `json:"lang_cd,omitempty"`
	Name                  string `json:"name,omitempty"`
	PostalCode            string `json:"postal_code,omitempty"`
	City                  string `json:"city,omitempty"`
	Line1                 string `json:"line_1,omitempty" gorm:"column:line_1"`
	Line2                 string `json:"line_2,omitempty" gorm:"column:line_2"`
	Line3                 string `json:"line_3,omitempty" gorm:"column:line_3"`
	CheckinBegin          string `json:"checkin_begin,omitempty"`
	CheckinEnd            string `json:"checkin_end,omitempty"`
	Checkout              string `json:"checkout,omitempty"`
	Instructions          string `json:"instructions,omitempty"`
	SpecialInstructions   string `json:"special_instructions,omitempty"`
	PolicyKnowBeforeYouGo string `json:"policy_know_before_you_go,omitempty"`
	MinAge                int    `json:"min_age"`
	HeroImage             string `json:"hero_image,omitempty"`
	common.Times          `gorm:"embedded"`
}

// HtTmCategorise 各カテゴリー情報テーブル
type HtTmCategorise struct {
	Name string `json:"name"`
}

// InitFacilityOutput 施設一覧の出力
type InitFacilityOutput struct {
	PropertyID        int64  `json:"property_id,omitempty"`
	WholesalerID      int64  `json:"wholesaler_id,omitempty"`
	Name              string `json:"name,omitempty"`
	StateProvinceName string `json:"state_province_name"`
	City              string `json:"city"`
	Line1             string `json:"line_1" gorm:"column:line_1"`
	Line2             string `json:"line_2" gorm:"column:line_2"`
	Line3             string `json:"line_3" gorm:"column:line_3"`
	DispPriority      bool   `json:"disp_priority"`
}

// UpdateDispPriorityInput サイト更新フラグ更新の入力
type UpdateDispPriorityInput struct {
	PropertyID   int64 `json:"property_id"`
	WholesalerID int64 `json:"wholesaler_id"`
	DispPriority bool  `json:"disp_priority"`
}

// BaseInfoInput 施設基本情報の取得入力
type BaseInfoInput struct {
	PropertyID int64 `json:"property_id" param:"propertyId" validate:"required"`
}

// BaseInfoOutput 施設基本情報の取得出力
type BaseInfoOutput struct {
	PropertyID        int64  `json:"property_id"`
	Name              string `json:"name"`
	PostalCode        string `json:"postal_code"`
	StateProvinceName string `json:"state_province_name"`
	City              string `json:"city"`
	Line1             string `json:"line_1" gorm:"column:line_1"`
	Line2             string `json:"line_2" gorm:"column:line_2"`
	Line3             string `json:"line_3" gorm:"column:line_3"`
	Phone             string `json:"phone"`
	Fax               string `json:"fax"`
	CategoryName      string `json:"category_name"`
	ConnectID         string `json:"connect_id"`
	ConnectPassword   string `json:"connect_password"`
}

// DetailInput 施設詳細情報の取得入力
type DetailInput struct {
	PropertyID int64 `json:"property_id" param:"propertyId" validate:"required"`
}

// DetailOutput 施設詳細情報の取得出力
type DetailOutput struct {
	PropertyID                   int64     `json:"property_id"`
	CheckinBegin                 string    `json:"checkin_begin"`
	CheckinEnd                   string    `json:"checkin_end"`
	Checkout                     string    `json:"checkout"`
	Instructions                 string    `json:"instructions"`
	SpecialInstructions          string    `json:"special_instructions"`
	PolicyKnowBeforeYouGo        string    `json:"policy_know_before_you_go"`
	FeeMandatory                 string    `json:"fee_mandatory"`
	FeeOptional                  string    `json:"fee_optional"`
	DescriptionAmenity           string    `json:"description_amenity"`
	DescriptionAttractions       string    `json:"description_attractions"`
	DescriptionBusinessAmenities string    `json:"description_business_amenities"`
	DescriptionDining            string    `json:"description_dining"`
	DescriptionLocation          string    `json:"description_location"`
	DescriptionHeadline          string    `json:"description_headline"`
	DescriptionRooms             string    `json:"description_rooms"`
	DispPriority                 bool      `json:"disp_priority"`
	Amenities                    []Amenity `gorm:"-" json:"amenities"`
}

// SaveBaseInfoInput 施設基本情報更新の入力
type SaveBaseInfoInput struct {
	PropertyID        int64  `json:"property_id" validate:"required"`
	Name              string `json:"name" validate:"required"`
	PostalCode        string `json:"postal_code" validate:"required"`
	StateProvinceName string `json:"state_province_name" validate:"required"`
	City              string `json:"city" validate:"required"`
	Line1             string `json:"line_1"`
	Line2             string `json:"line_2"`
	Line3             string `json:"line_3"`
	Phone             string `json:"phone"`
	Fax               string `json:"fax"`
	ConnectID         string `json:"connect_id"`
	ConnectPassword   string `json:"connect_password"`
}

// SaveDetailInput 施設詳細情報更新の入力
type SaveDetailInput struct {
	PropertyID                   int64     `json:"property_id" validate:"required"`
	CheckinBegin                 string    `json:"checkin_begin"`
	CheckinEnd                   string    `json:"checkin_end"`
	Checkout                     string    `json:"checkout"`
	Instructions                 string    `json:"instructions"`
	SpecialInstructions          string    `json:"special_instructions"`
	PolicyKnowBeforeYouGo        string    `json:"policy_know_before_you_go"`
	FeeMandatory                 string    `json:"fee_mandatory"`
	FeeOptional                  string    `json:"fee_optional"`
	DescriptionAmenity           string    `json:"description_amenity"`
	DescriptionAttractions       string    `json:"description_attractions"`
	DescriptionBusinessAmenities string    `json:"description_business_amenities"`
	DescriptionDining            string    `json:"description_dining"`
	DescriptionLocation          string    `json:"description_location"`
	DescriptionHeadline          string    `json:"description_headline"`
	DescriptionRooms             string    `json:"description_rooms"`
	Amenities                    []Amenity `json:"amenities"`
}

// Amenity 施設アメニティ
type Amenity struct {
	PropertyAmenityID string `json:"property_amenity_id"`
	AmenityName       string `json:"amenity_name"`
}

// IFacilityUsecase 施設関連のusecaseのインターフェース
type IFacilityUsecase interface {
	FetchAll(hmUser account.HtTmHotelManager) ([]InitFacilityOutput, error)
	UpdateDispPriority(request *UpdateDispPriorityInput) error
	FetchBaseInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, request *BaseInfoInput) (*BaseInfoOutput, error)
	SaveBaseInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, request *SaveBaseInfoInput) error
	FetchDetail(request *BaseInfoInput) (*DetailOutput, error)
	SaveDetail(request *SaveDetailInput) error
	FetchAllAmenities() ([]Amenity, error)
}

// IParentUsecase 親アカウントに特化したユースケース（特定の施設やホールセラーに依存しない）
type IParentUsecase interface {
	FetchAll(hmUser account.HtTmHotelManager) ([]InitFacilityOutput, error)
}

// IFacilityRepository 施設関連のrepositoryのインターフェース
type IFacilityRepository interface {
	common.Repository
	FetchAllClientCompanies(hotelManagerID int64) ([]HtTmProperties, error)
	FetchAllFacilitiesByPropertyID(propertyIDs []int64) ([]InitFacilityOutput, error)
	FetchProperty(propertyID int64) (*HtTmProperties, error)
	FetchCategory(propertyID int64) (*HtTmCategorise, error)
	UpdateProperty(property *HtTmProperties) error
	UpsertPropertyLangsBase(property *HtTmPropertyLangs) error
	UpsertPropertyLangsDetail(upsertData *HtTmPropertyLangs) error
}
