package plan

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"gorm.io/gorm"
)

// HtTmPlanTemas てまのプランテーブル
type HtTmPlanTemas struct {
	PlanTemaID          int            `json:"plan_tema_ids"`
	PropertyID          int64          `json:"property_id"`
	PackagePlanCode     int            `json:"package_plan_code"`
	PlanName            string         `json:"plan_name"`
	Desc                string         `json:"desc"`
	LangCd              string         `json:"lang_cd"`
	PlanType            int            `json:"plan_type"`
	Payment             int            `json:"payment"`
	ListingPeriodStart  string         `json:"listing_period_start"`
	ListingPeriodEnd    string         `json:"listing_period_end"`
	IsRoomCharge        int            `json:"is_room_charge"`
	Available           int            `json:"available"`
	RateType            int            `json:"rate_type"`
	ListingPeriodStartH int            `json:"listing_period_start_h"`
	ListingPeriodStartM int            `json:"listing_period_start_m"`
	ListingPeriodEndH   int            `json:"listing_period_end_h"`
	ListingPeriodEndM   int            `json:"listing_period_end_m"`
	ReservePeriodStart  string         `json:"reserve_period_start"`
	ReservePeriodEnd    string         `json:"reserve_period_end"`
	CheckinTimeStartH   int            `json:"checkin_time_start_h"`
	CheckinTimeStartM   int            `json:"checkin_time_start_m"`
	CheckinTimeEndH     int            `json:"checkin_time_end_h"`
	CheckinTimeEndM     int            `json:"checkin_time_end_m"`
	CheckoutTimeEndH    int            `json:"checkout_time_end_h"`
	CheckoutTimeEndM    int            `json:"checkout_time_end_m"`
	StayLimitMin        int            `json:"stay_limit_min"`
	StayLimitMax        int            `json:"stay_limit_max"`
	AdvBkCreateStartD   int            `json:"adv_bk_create_start_d"`
	AdvBkCreateStartH   int            `json:"adv_bk_create_start_h"`
	AdvBkCreateStartM   int            `json:"adv_bk_create_start_m"`
	AdvBkCreateEndD     int            `json:"adv_bk_create_end_d"`
	AdvBkCreateEndH     int            `json:"adv_bk_create_end_h"`
	AdvBkCreateEndM     int            `json:"adv_bk_create_end_m"`
	AdvBkModifyEndD     int            `json:"adv_bk_modify_end_d"`
	AdvBkModifyEndH     int            `json:"adv_bk_modify_end_h"`
	AdvBkModifyEndM     int            `json:"adv_bk_modify_end_m"`
	AdvBkCancelEndD     int            `json:"adv_bk_cancel_end_d"`
	AdvBkCancelEndH     int            `json:"adv_bk_cancel_end_h"`
	AdvBkCancelEndM     int            `json:"adv_bk_cancel_end_m"`
	ServiceChargeType   int            `json:"service_charge_type"`
	ServiceChargeValue  int            `json:"service_charge_value"`
	OptionalItems       string         `json:"optional_items"`
	CancelpolicyJson    string         `json:"cancelpolicy_json"`
	ChildrenJson        string         `json:"children_json"`
	ChildrenAcceptable1 int            `json:"children_acceptable1"`
	ChildrenAcceptable2 int            `json:"children_acceptable2"`
	ChildrenAcceptable3 int            `json:"children_acceptable3"`
	ChildrenAcceptable4 int            `json:"children_acceptable4"`
	ChildrenAcceptable5 int            `json:"children_acceptable5"`
	PictureJson         string         `json:"picture_json"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at"`
	common.Times
}
type TemaBulkListOutput struct {
	RoomTypeID      int64                  `json:"room_type_id,omitempty"`
	RoomName        string                 `json:"room_name,omitempty" validate:"required,max=35"`
	RoomIsStopSales bool                   `json:"room_is_stop_sales"`
	RoomImageHref   string                 `json:"room_image_href"`
	Plans           []TemaBulkDetailOutput `json:"plans"`
}
type PlanPolicy struct {
	CheckinStart string `json:"checkin_start,omitempty"`
	CheckinEnd   string `json:"checkin_end,omitempty" `
	Checkout     string `json:"checkout"`
}
type TemaBulkDetailOutput struct {
	price.TemaPlanTable
	Name       string                   `json:"room_name,omitempty"`
	Images     []image.PlanImagesOutput `json:"images"`
	ChildRates []price.ChildRateTable   `json:"child_rates,omitempty"`
}

// HtTmChildRateTemas Tema child price setting table
type HtTmChildRateTemas struct {
	price.ChildRateTable `gorm:"embedded"`
}

type IPlanBulkTemaUsecase interface {
	FetchList(request *ListInput) ([]TemaBulkListOutput, error)
	CreateBulk(request []price.TemaPlanData) (string, error)
	Detail(request *DetailInput) (*TemaBulkDetailOutput, error)
}

// IPlanTemaRepository てまプラン関連のrepositoryのインターフェース
type IPlanTemaRepository interface {
	common.Repository
	// FetchOne プランを1件取得
	FetchOnePlan(propertyID int64, packagePlanCode int) (*HtTmPlanTemas, error)
	// FetchOne With PlanID
	FetchOneWithPlanID(planID int64) (price.HtTmPlanTemas, error)
	// FetchList プランを複数件取得
	FetchList(propertyID int64, packagePlanCodeList []int) ([]HtTmPlanTemas, error)
	//fetchAllByPropertyID
	FetchAllByPropertyID(req ListInput) ([]price.HtTmPlanTemas, error)
	// GetPlanIfPlanCodeExist
	GetPlanIfPlanCodeExist(propertyID int64, planCode int64, roomTypeID int64) (price.HtTmPlanTemas, error)
	// UpdatePlanTema
	UpdatePlanBulkTema(planTable price.HtTmPlanTemas, planID int64) error
	//Get Next Plan Id
	GetNextPlanID() (price.HtTmPlanTemas, error)
	// CreatePlansTema Create plan
	CreatePlanBulkTema(planTable price.HtTmPlanTemas) error
	// ClearChildRateTema
	ClearChildRateTema(planID int64) error
	// ClearChildRateTema
	ClearImageTema(planID int64) error
	// CreateChildRateTema Create multiple child fares
	CreateChildRateTema(childRates []price.HtTmChildRateTemas) error
	// DeletePlanTema
	DeletePlanTema(planCode int64, roomTypeIDs []int64) error
	// MatchesPlanIDAndPropertyID Are propertyID and planID linked?
	MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool
	// FetchActiveByPlanGroupID Get multiple active plans linked to plan_group_id
	FetchActiveByPlanCode(planCode string) ([]price.HtTmPlanTemas, error)
	// FetchChildRates Get multiple child price settings linked to plan_id
	FetchChildRates(planID int64) ([]price.HtTmChildRateTemas, error)
}
