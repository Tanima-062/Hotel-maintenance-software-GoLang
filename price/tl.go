package price

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
)

type PlanTable struct {
	PlanID                   int64     `gorm:"primaryKey;autoIncrement:true" json:"plan_id,omitempty"`
	PlanGroupID              int64     `json:"plan_group_id"`
	RoomTypeID               int64     `json:"room_type_id,omitempty"`
	PropertyID               int64     `json:"property_id,omitempty" validate:"required"`
	PlanCode                 string    `json:"plan_code,omitempty" validate:"required"`
	LangCd                   string    `json:"lang_cd,omitempty"`
	Name                     string    `json:"name" validate:"required"`
	Description              string    `json:"description"`
	ChargeCategory           int8      `json:"charge_category" gorm:"default:1"`
	TaxCategory              bool      `json:"tax_category"`
	AccommodationPeriodStart time.Time `gorm:"type:time" json:"accommodation_period_start"`
	AccommodationPeriodEnd   time.Time `gorm:"type:time" json:"accommodation_period_end"`
	IsAccommodatedYearRound  bool      `json:"is_accommodated_year_round"`
	PublishingStartDate      time.Time `gorm:"type:time" json:"publishing_start_date"`
	PublishingEndDate        time.Time `gorm:"type:time" json:"publishing_end_date"`
	IsPublishedYearRound     bool      `json:"is_published_year_round"`
	ReserveAcceptDate        int16     `json:"reserve_accept_date"`
	ReserveAcceptTime        string    `json:"reserve_accept_time,omitempty"`
	ReserveDeadlineDate      int16     `json:"reserve_deadline_date"`
	ReserveDeadlineTime      string    `json:"reserve_deadline_time,omitempty"`
	MinStayCategory          bool      `json:"min_stay_category"`
	MinStayNum               int8      `json:"min_stay_num"`
	MaxStayCategory          bool      `json:"max_stay_category"`
	MaxStayNum               int8      `json:"max_stay_num"`
	MealConditionBreakfast   bool      `json:"meal_condition_breakfast"`
	MealConditionDinner      bool      `json:"meal_condition_dinner"`
	MealConditionLunch       bool      `json:"meal_condition_lunch"`
	IsNoCancel               bool      `json:"is_no_cancel"`
	IsStopSales              bool      `json:"is_stop_sales"`
	CancelPolicy             string    `json:"cancel_policy"`
	IsDelete                 bool      `json:"is_delete,omitempty"`
	common.Times             `gorm:"embedded"`
}

type PlanData struct {
	PlanTable
	RoomTypeCode       string                  `json:"room_type_code,omitempty" validate:"required"`
	SelectedRooms      []int64                 `json:"selected_rooms"`
	ChildRates         []ChildRateTable        `json:"child_rates"`
	Images             []image.PlanImagesInput `json:"images"`
	PlanCancelPolicyId *uint64                 `json:"plan_cancel_policy_id"`
	CheckinStart       string                  `json:"checkin_start"`
	CheckinEnd         string                  `json:"checkin_end"`
	Checkout           string                  `json:"checkout"`
	Prices             map[string][]Price      `json:"prices"`
}

type PriceData struct {
	PropertyID           int64              `json:"property_id"    validate:"required"`
	PlanCode             string             `json:"plan_code"      validate:"required"`
	RoomTypeCode         string             `json:"room_type_code" validate:"required"`
	PublishingStartDate  time.Time          `gorm:"type:time"  json:"publishing_start_date"`
	PublishingEndDate    time.Time          `gorm:"type:time"  json:"publishing_end_date"`
	IsPublishedYearRound bool               `json:"is_published_year_round"`
	Prices               map[string][]Price `json:"prices"`
}

// HtTmPlanTLs Tl plan table
type HtTmPlanTls struct {
	PlanTable `gorm:"embedded"`
}

// HtTmPriceTLs
type HtTmPriceTls struct {
	PriceTable `gorm:"embedded"`
}

// HtTmChildRateTLs
type HtTmChildRateTls struct {
	ChildRateTable `gorm:"embedded"`
}

// IPriceBulkUsecase
type IPriceBulkTlUsecase interface {
	GetPriceData(request PlanTable, childRateTables []HtTmChildRateTls, priceData Price, date string) HtTmPriceTls
	Update(request []PriceData) (string, error)
}

// IPriceTLRepository
type IPriceTlRepository interface {
	common.Repository
	// FetchChildRates Get multiple price settings linked to a plan
	FetchChildRates(planID int64) ([]HtTmChildRateTls, error)
	// FetchAllByPlanIDList Get multiple charges associated with multiple plan IDs within the period
	FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]HtTmPriceTls, error)
	// FetchPricesByPlanID Get multiple charges from today onwards
	FetchPricesByPlanID(planID int64) ([]HtTmPriceTls, error)

	// FetchChildRatesByPlanID Get multiple price settings linked to plan
	FetchChildRatesByPlanID(planID int64) ([]HtTmChildRateTls, error)
	// update prices
	UpdatePrice(planID int64, useDate string, rateTypeCode string, price int, isStopSales bool, priceData HtTmPriceTls) error
	//CheckIfPriceExistsByPlanIDRateTypeCodeAndUseDate
	CheckIfPriceExistsByPlanIDAndRateTypeCodeAndUseDate(planID int64, rateTypeCode string, useDate string) (bool, error)
	// GetPriceByPropertyIDAndPlanCode GetPrice
	GetPriceByPropertyIDAndPlanCode(propertyID int64, planCode string) ([]HtTmPriceTls, error)
	// create price
	CreatePrice(priceTable HtTmPriceTls) error
	// get price by plan_id and rate_type_code
	GetPriceByPlanIDAndRateTypeCode(planID int64, rateTypeCode string) (HtTmPriceTls, error)
}
