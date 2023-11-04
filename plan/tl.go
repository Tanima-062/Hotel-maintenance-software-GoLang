package plan

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"github.com/Adventureinc/hotel-hm-api/src/price"
)

// HtTmPlanGroupIDTls Tl plan group ID numbering table
type HtTmPlanGroupIDTls struct {
	PlanGroupID int64 `json:"plan_group_id"`
}

// HtTmChildRateTls Tl child price setting table
type HtTmChildRateTls struct {
	price.ChildRateTable `gorm:"embedded"`
}

type BulkListOutput struct {
	RoomTypeID      int64              `json:"room_type_id,omitempty"`
	RoomName        string             `json:"room_name,omitempty" validate:"required,max=35"`
	RoomIsStopSales bool               `json:"room_is_stop_sales"`
	RoomImageHref   string             `json:"room_image_href"`
	Plans           []BulkDetailOutput `json:"plans"`
}

type BulkDetailOutput struct {
	price.PlanTable
	RoomName     string                   `json:"room_name,omitempty"`
	ActiveRooms  []int64                  `json:"active_rooms"`
	Images       []image.PlanImagesOutput `json:"images"`
	ChildRates   []price.ChildRateTable   `json:"child_rates,omitempty"`
	CheckinStart string                   `json:"checkin_start"`
	CheckinEnd   string                   `json:"checkin_end"`
	Checkout     string                   `json:"checkout"`
}

type IPlanBulkUsecase interface {
	FetchList(request *ListInput) ([]BulkListOutput, error)
	CreateBulk(request []price.PlanData) (string, error)
	Detail(request *DetailInput) (*BulkDetailOutput, error)
}

// IPlanTlRepository Tl plan related repository interface
type IPlanTlRepository interface {
	common.Repository
	// FetchAllByPropertyID Acquire multiple plans linked to property_id that has not been deleted
	FetchAllByPropertyID(req ListInput) ([]price.HtTmPlanTls, error)
	// FetchActiveByPlanGroupID Get multiple active plans linked to plan_group_id
	FetchActiveByPlanCode(planCode string) ([]price.HtTmPlanTls, error)
	// FetchOne Get one undeleted plan associated with plan_id
	FetchOne(planID int64) (price.HtTmPlanTls, error)
	// FetchList Get multiple undeleted plans associated with plan_id
	FetchList(planIDList []int64) ([]price.HtTmPlanTls, error)
	// MatchesPlanIDAndPropertyID Are propertyID and planID linked?
	MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool
	// FetchChildRates Get multiple child price settings linked to plan_id
	FetchChildRates(planID int64) ([]price.HtTmChildRateTls, error)
	// CreateChildRateTL Create multiple child fares
	CreateChildRateTl(childRates []price.HtTmChildRateTls) error

	// FetchPlan with plan_id
	FetchPlan(planID int64) (price.HtTmPlanTls, error)
	// CreatePlansTL Create plan with updated plan group id
	CreatePlanBulkTl(planTable price.HtTmPlanTls) error
	// GetPlanIfPlanCodeExist
	GetPlanIfPlanCodeExist(propertyID int64, planCode string, roomTypeID int64) (price.HtTmPlanTls, error)
	// UpdatePlanTl
	UpdatePlanBulkTl(planTable price.HtTmPlanTls, planID int64) error
	// ClearChildRateTl
	ClearChildRateTl(planID int64) error
	// ClearChildRateTl
	ClearImageTl(planID int64) error
	// Get Next Plan Id
	GetNextPlanID() (price.HtTmPlanTls, error)
	// Get plan by propertyID and PlanCode
	GetPlanByPropertyIDAndPlanCodeAndRoomTypeCode(propertyID int64, planCode string, roomTypeCode string) ([]price.HtTmPlanTls, error)
	// DeletePlanTl
	DeletePlanTl(planCode string, roomTypeIDs []int64) error
}
