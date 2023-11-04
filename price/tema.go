package price

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"gorm.io/gorm"
	"time"
)

type TemaPlanTable struct {
	PlanID                  int64          `json:"plan_id" gorm:"column:plan_tema_id"`
	PropertyID              int64          `json:"property_id"`
	PackagePlanCode         int64          `json:"plan_code" gorm:"column:package_plan_code"`
	PlanName                string         `json:"name" gorm:"column:plan_name"`
	Desc                    string         `json:"description" gorm:"column:desc"`
	PlanGroupID             int64          `json:"plan_group_id"`
	LangCd                  string         `json:"lang_cd"`
	PlanType                int            `json:"plan_type"`
	Payment                 int            `json:"payment"`
	ListingPeriodStart      string         `json:"publishing_start_date" gorm:"column:listing_period_start"`
	ListingPeriodEnd        string         `json:"publishing_end_date" gorm:"column:listing_period_end"`
	IsRoomCharge            int            `json:"is_room_charge"`
	RateType                int            `json:"charge_category" gorm:"column:rate_type"`
	ListingPeriodStartH     int            `json:"listing_period_start_h"`
	ListingPeriodStartM     int            `json:"listing_period_start_m"`
	ListingPeriodEndH       int            `json:"listing_period_end_h"`
	ListingPeriodEndM       int            `json:"listing_period_end_m"`
	ReservePeriodStart      string         `json:"accommodation_period_start" gorm:"column:reserve_period_start"`
	ReservePeriodEnd        string         `json:"accommodation_period_end" gorm:"column:reserve_period_end"`
	CheckinTimeStartH       int            `json:"checkin_time_start_h"`
	CheckinTimeStartM       int            `json:"checkin_time_start_m"`
	CheckinTimeEndH         int            `json:"checkin_time_end_h"`
	CheckinTimeEndM         int            `json:"checkin_time_end_m"`
	CheckoutTimeEndH        int            `json:"checkout_time_end_h"`
	CheckoutTimeEndM        int            `json:"checkout_time_end_m"`
	CheckinStart            string         `json:"checkin_start" gorm:"-"`
	CheckinEnd              string         `json:"checkin_end" gorm:"-"`
	Checkout                string         `json:"checkout" gorm:"-"`
	StayLimitMin            int            `json:"min_stay_num" gorm:"column:stay_limit_min"`
	StayLimitMax            int            `json:"max_stay_num" gorm:"column:stay_limit_max"`
	AdvBKCreateStartEnabled int            `json:"adv_bk_create_start_enabled"`
	AdvBKCreateEndEnabled   int            `json:"adv_bk_create_end_enabled"`
	AdvBkCreateStartD       int            `json:"adv_bk_create_start_d"`
	ReserveAcceptTime       string         `json:"reserve_accept_time" gorm:"-"`
	AdvBkCreateStartH       int            `json:"adv_bk_create_start_h"`
	AdvBkCreateStartM       int            `json:"adv_bk_create_start_m"`
	ReserveDeadlineTime     string         `json:"reserve_deadline_time" gorm:"-"`
	AdvBkCreateEndD         int            `json:"adv_bk_create_end_d"`
	AdvBkCreateEndH         int            `json:"adv_bk_create_end_h"`
	AdvBkCreateEndM         int            `json:"adv_bk_create_end_m"`
	AdvBkModifyEndEnabled   int            `json:"adv_bk_modify_end_enabled"`
	AdvBkModifyEndD         int            `json:"adv_bk_modify_end_d"`
	AdvBkModifyEndH         int            `json:"adv_bk_modify_end_h"`
	AdvBkModifyEndM         int            `json:"adv_bk_modify_end_m"`
	AdvBkCancelEndEnabled   int            `json:"adv_bk_cancel_end_enabled"`
	AdvBkCancelEndD         int            `json:"adv_bk_cancel_end_d"`
	AdvBkCancelEndH         int            `json:"adv_bk_cancel_end_h"`
	AdvBkCancelEndM         int            `json:"adv_bk_cancel_end_m"`
	ServiceChargeType       int            `json:"service_charge_type"`
	ServiceChargeValue      int            `json:"service_charge_value"`
	OptionalItems           string         `json:"optional_items"`
	CancelpolicyJson        string         `json:"cancelpolicy_json"`
	ChildrenJson            string         `json:"children_json"`
	ChildrenAcceptable1     int            `json:"children_acceptable1"`
	ChildrenAcceptable2     int            `json:"children_acceptable2"`
	ChildrenAcceptable3     int            `json:"children_acceptable3"`
	ChildrenAcceptable4     int            `json:"children_acceptable4"`
	ChildrenAcceptable5     int            `json:"children_acceptable5"`
	PictureJson             string         `json:"picture_json"`
	RoomTypeID              int64          `json:"room_type_id,omitempty"`
	IsPublishedYearRound    bool           `json:"is_published_year_round"`
	IsAccommodatedYearRound bool           `json:"is_accommodated_year_round"`
	TaxCategory             bool           `json:"tax_category"`
	MinStayCategory         int            `json:"min_stay_category"`
	MaxStayCategory         int            `json:"max_stay_category"`
	MealConditionBreakfast  bool           `json:"meal_condition_breakfast"`
	MealConditionDinner     bool           `json:"meal_condition_dinner"`
	MealConditionLunch      bool           `json:"meal_condition_lunch"`
	DeletedAt               gorm.DeletedAt `json:"deleted_at"`
	Available               bool           `json:"available"  gorm:"column:available"`
	common.Times
}

type TemaPlanData struct {
	TemaPlanTable
	RoomTypeCode string                  `json:"room_type_code,omitempty" validate:"required"`
	ChildRates   []ChildRateTable        `json:"child_rates"`
	Images       []image.PlanImagesInput `json:"images"`
}

type HtTmPriceTemas struct {
	PriceTemaTable `gorm:"embedded"`
}
type PriceTemaTable struct {
	PriceTemaID     int64     `gorm:"primaryKey;autoIncrement:true" json:"price_tema_id,omitempty"`
	PropertyID      int64     `json:"property_id,omitempty"`
	RoomTypeCode    int       `json:"room_type_code,omitempty"`
	PackagePlanCode int64     `json:"package_plan_code,omitempty"`
	PriceDate       time.Time `gorm:"type:time" json:"price_date"`
	Disable         bool      `json:"is_stop_sales" gorm:"column:disable"`
	IsRoomCharge    int       `json:"is_room_charge"`
	MinPax          int64     `json:"min_pax"`
	MaxPax          int64     `json:"max_pax"`
	TemaPriceType
	common.Times `gorm:"embedded"`
}

type PriceTemaData struct {
	PropertyID           int64                  `json:"property_id,omitempty"`
	RoomTypeCode         string                 `json:"room_type_code,omitempty"`
	PackagePlanCode      int64                  `json:"package_plan_code,omitempty"`
	IsRoomCharge         int                    `json:"is_room_charge"`
	ListingPeriodStart   string                 `json:"listing_period_start"`
	ListingPeriodEnd     string                 `json:"listing_period_end"`
	IsPublishedYearRound bool                   `json:"is_published_year_round"`
	MinPax               int64                  `json:"min_pax"`
	MaxPax               int64                  `json:"max_pax"`
	Disable              bool                   `json:"is_stop_sales" gorm:"column:disable"`
	PriceList            map[string][]PriceTema `json:"prices"`
}

type PriceTema struct {
	Price int64 `json:"price"`
}

// HtTmPlanTemas Tema plan table
type HtTmPlanTemas struct {
	TemaPlanTable `gorm:"embedded"`
}

// HtTmChildRateTemas
type HtTmChildRateTemas struct {
	ChildRateTable `gorm:"embedded"`
}

// IPriceBulkUsecase
type IPriceBulkTemaUsecase interface {
	Update(request []PriceTemaData) (string, error)
}
type IPriceTemaRepository interface {
	common.Repository

	// FetchAllByPlanIDList Get multiple charges associated with multiple plan IDs within the period
	FetchAllByPlanCodeList(planCodeList []int64, startDate string, endDate string) ([]HtTmPriceTemas, error)
	// FetchPricesByPlanID Get multiple charges from today onwards
	FetchPricesByPlanID(planID int64) ([]HtTmPriceTemas, error)
	// DeletePriceTema
	DeletePriceTema(propertyID int64, packagePlanCode int64, roomTypeCode int, priceDate string) error
	// create price
	CreatePrice(priceTable HtTmPriceTemas) error
}
