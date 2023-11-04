package plan

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"github.com/Adventureinc/hotel-hm-api/src/price"
)

// PlanTable プランテーブル
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
	IsPackage                bool      `json:"is_package"`
	IsNoCancel               bool      `json:"is_no_cancel"`
	IsStopSales              bool      `json:"is_stop_sales"`
	IsDelete                 bool      `json:"is_delete,omitempty"`
	common.Times             `gorm:"embedded"`
}

// ListInput 一覧の入力
type ListInput struct {
	PropertyID int64 `json:"property_id" param:"propertyId" validate:"required"`
	common.Paging
}

// ListOutput 一覧の出力
type ListOutput struct {
	RoomTypeID      int64          `json:"room_type_id,omitempty"`
	RoomName        string         `json:"room_name,omitempty" validate:"required,max=35"`
	RoomIsStopSales bool           `json:"room_is_stop_sales"`
	RoomImageHref   string         `json:"room_image_href"`
	Plans           []DetailOutput `json:"plans"`
}

// DetailInput 詳細の入力
type DetailInput struct {
	PropertyID int64 `json:"property_id" param:"propertyId" validate:"required"`
	PlanID     int64 `json:"plan_id" param:"planId" validate:"required"`
}

// DetailOutput 詳細の出力
type DetailOutput struct {
	PlanTable
	RoomName           string                   `json:"room_name,omitempty"`
	ActiveRooms        []int64                  `json:"active_rooms"`
	Images             []image.PlanImagesOutput `json:"images"`
	ChildRates         []price.ChildRateTable   `json:"child_rates,omitempty"`
	PlanCancelPolicyID *uint64                  `json:"plan_cancel_policy_id"`
	CheckinStart       string                   `json:"checkin_start"`
	CheckinEnd         string                   `json:"checkin_end"`
	Checkout           string                   `json:"checkout"`
}

// SaveInput 作成・更新の入力
type SaveInput struct {
	PlanTable
	SelectedRooms      []int64                 `json:"selected_rooms"`
	ChildRates         []price.ChildRateTable  `json:"child_rates"`
	Images             []image.PlanImagesInput `json:"images"`
	PlanCancelPolicyId *uint64                 `json:"plan_cancel_policy_id"`
	CheckinStart       string                  `json:"checkin_start"`
	CheckinEnd         string                  `json:"checkin_end"`
	Checkout           string                  `json:"checkout"`
}

// DeleteInput 削除の入力
type DeleteInput struct {
	PlanID int64 `json:"plan_id" param:"planId" validate:"required"`
}

// CheckDuplicatePlanCode plan_codeの重複チェック用
type CheckDuplicatePlanCode struct {
	PlanCode   string
	RoomTypeID int64
}

// StopSalesInput 売止更新
type StopSalesInput struct {
	PlanID      int64 `json:"plan_id" validate:"required"`
	IsStopSales bool  `json:"is_stop_sales"`
}

type HtTmPlanCheckInOuts struct {
	PlanCheckInOutId string `gorm:"primaryKey;autoIncrement:true"`
	WholesalerID     int
	PropertyID       int64
	PlanID           int64
	CheckInBegin     string
	CheckInEnd       string
	CheckOut         string
}

type CheckInOutInfo struct {
	WholesalerID int
	PropertyID   int64
	PlanID       int64
	CheckInBegin string
	CheckInEnd   string
	CheckOut     string
}

// IPlanUsecase プラン関連のusecaseのインターフェース
type IPlanUsecase interface {
	FetchList(request *ListInput) ([]ListOutput, error)
	Detail(request *DetailInput) (*DetailOutput, error)
	Create(request *SaveInput) error
	Update(request *SaveInput) error
	Delete(planID int64) error
	UpdateStopSales(request *StopSalesInput) error
}

// ICommonPlanRepository プラン共通処理repositoryのインターフェース
type ICommonPlanRepository interface {
	FetchCheckInOut(propertyId int64, planId int64) (*HtTmPlanCheckInOuts, error)
	UpsertCheckInOut(checkInOutInfo CheckInOutInfo) error
	DeleteCheckInOut(wholesalerId int, planId int64) error
}
