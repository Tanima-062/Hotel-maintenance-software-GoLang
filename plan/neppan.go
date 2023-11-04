package plan

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/price"
)

// HtTmPlanNeppans ねっぱんのプランテーブル
type HtTmPlanNeppans struct {
	PlanTable `gorm:"embedded"`
}

// HtTmPlanGroupIDNeppans ねっぱんのプラングループIDの採番テーブル
type HtTmPlanGroupIDNeppans struct {
	PlanGroupID int64 `json:"plan_group_id"`
}

// HtTmChildRateNeppans ねっぱんの子供料金設定テーブル
type HtTmChildRateNeppans struct {
	price.ChildRateTable `gorm:"embedded"`
}

// IPlanNeppanRepository ねっぱんプラン関連のrepositoryのインターフェース
type IPlanNeppanRepository interface {
	common.Repository
	// FetchAllByPropertyID 削除されていないproperty_idに紐づくプラン複数件取得
	FetchAllByPropertyID(req ListInput) ([]HtTmPlanNeppans, error)
	// FetchAllByRoomTypeID 削除されていないroom_type_idに紐づくプラン複数件取得
	FetchAllByRoomTypeID(roomTypeID int64) ([]HtTmPlanNeppans, error)
	// FetchAllByPlanGroupID plan_group_idに紐づくプラン複数件取得
	FetchAllByPlanGroupID(planGroupID int64) ([]HtTmPlanNeppans, error)
	// FetchActiveByPlanGroupID plan_group_idに紐づくアクティブなプラン複数件取得
	FetchActiveByPlanGroupID(planGroupID int64) ([]HtTmPlanNeppans, error)
	// FetchAllByCancelPolicyID cancel_policy_idに紐づく削除されていないプランを複数件取得
	FetchAllByCancelPolicyID(cancelPolicyID uint64) ([]HtTmPlanNeppans, error)
	// FetchOne plan_idに紐づく削除されていないプランを1件取得
	FetchOne(planID int64) (HtTmPlanNeppans, error)
	// FetchList plan_idに紐づく削除されていないプランを複数件取得
	FetchList(planIDList []int64) ([]HtTmPlanNeppans, error)
	// MatchesPlanIDAndPropertyID propertyIDとplanIDが紐付いているか
	MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool
	// FetchChildRates plan_idに紐づく子供料金設定を複数件取得
	FetchChildRates(planID int64) ([]HtTmChildRateNeppans, error)
	// DeletePlanByRoomTypeID room_type_idに紐づくプランを論理削除
	DeletePlanByRoomTypeID(roomTypeID int64) error
	// CheckPlanCode room_type_idとplan_codeの組み合わせで合致するものの件数を取得
	CheckPlanCode(propertyID int64, planCodeList []CheckDuplicatePlanCode) int64
	// MakePlansNeppan プランを複数件作成
	MakePlansNeppan(planTables []HtTmPlanNeppans) error
	// CreatePlansNeppan プランを複数件新規作成
	CreatePlansNeppan(planTables []HtTmPlanNeppans) error
	// CreateChildRateNeppan 子供料金を複数件新規作成
	CreateChildRateNeppan(childRates []HtTmChildRateNeppans) error
	// UpdatePlanNeppan プラン情報を更新
	UpdatePlanNeppan(planTable *HtTmPlanNeppans, planIDs []int64) error
	// UpdateChildRateNeppan 子供料金設定を更新
	UpdateChildRateNeppan(child *HtTmChildRateNeppans) error
	// FetchPlanGroupID plan_idに基づくplan_group_idの取得
	FetchPlanGroupID(planIDs int64) (int64, error)
	// DeletePlanNeppan プランを一斉に論理削除
	DeletePlanNeppan(planID []int64) error
	// UpdateStopSales プランの売止を更新
	UpdateStopSales(planIDList []int64, isStopSales bool) error
	// UpdateStopSalesByRoomID 部屋に紐づくプランの売止フラグを複数件更新
	UpdateStopSalesByRoomID(roomTypeID int64, isStopSales bool) error
	// UpdateStopSalesByRoomIDList room_type_id(複数)に紐づくプランの売止フラグを複数件更新
	UpdateStopSalesByRoomIDList(roomTypeIDList []int64, isStopSales bool) error
}
