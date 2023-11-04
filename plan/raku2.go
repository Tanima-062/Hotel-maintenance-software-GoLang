package plan

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/price"
)

// HtTmPlanRaku2s らく通のプランテーブル
type HtTmPlanRaku2s struct {
	PlanTable `gorm:"embedded"`
}

// HtTmPlanGroupIDRaku2s らく通のプラングループIDの採番テーブル
type HtTmPlanGroupIDRaku2s struct {
	PlanGroupID int64 `json:"plan_group_id"`
}

// HtTmChildRateRaku2s らく通の子供料金設定テーブル
type HtTmChildRateRaku2s struct {
	price.ChildRateTable `gorm:"embedded"`
}

// IPlanRaku2Repository らく通プラン関連のrepositoryのインターフェース
type IPlanRaku2Repository interface {
	common.Repository
	// FetchAllByPropertyID 削除されていないproperty_idに紐づくプラン複数件取得
	FetchAllByPropertyID(req ListInput) ([]HtTmPlanRaku2s, error)
	// FetchAllByRoomTypeID 削除されていないroom_type_idに紐づくプラン複数件取得
	FetchAllByRoomTypeID(roomTypeID int64) ([]HtTmPlanRaku2s, error)
	// FetchAllByPlanGroupID plan_group_idに紐づくプラン複数件取得
	FetchAllByPlanGroupID(planGroupID int64) ([]HtTmPlanRaku2s, error)
	// FetchActiveByPlanGroupID plan_group_idに紐づくアクティブなプラン複数件取得
	FetchActiveByPlanGroupID(planGroupID int64) ([]HtTmPlanRaku2s, error)
	// FetchAllByCancelPolicyID cancel_policy_idに紐づく削除されていないプランを複数件取得
	FetchAllByCancelPolicyID(cancelPolicyID uint64) ([]HtTmPlanRaku2s, error)
	// FetchOne plan_idに紐づく削除されていないプランを1件取得
	FetchOne(planID int64) (HtTmPlanRaku2s, error)
	// FetchList plan_idに紐づく削除されていないプランを複数件取得
	FetchList(planIDList []int64) ([]HtTmPlanRaku2s, error)
	// MatchesPlanIDAndPropertyID propertyIDとplanIDが紐付いているか
	MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool
	// FetchChildRates plan_idに紐づく子供料金設定を複数件取得
	FetchChildRates(planID int64) ([]HtTmChildRateRaku2s, error)
	// DeletePlanByRoomTypeID room_type_idに紐づくプランを論理削除
	DeletePlanByRoomTypeID(roomTypeID int64) error
	// CheckPlanCode room_type_idとplan_codeの組み合わせで合致するものの件数を取得
	CheckPlanCode(propertyID int64, planCodeList []CheckDuplicatePlanCode) int64
	// MakePlansRaku2 プランを複数件作成
	MakePlansRaku2(planTables []HtTmPlanRaku2s) error
	// CreatePlansRaku2 プランを複数件新規作成
	CreatePlansRaku2(planTables []HtTmPlanRaku2s) error
	// CreateChildRateRaku2 子供料金を複数件新規作成
	CreateChildRateRaku2(childRates []HtTmChildRateRaku2s) error
	// UpdatePlanRaku2 プラン情報を更新
	UpdatePlanRaku2(planTable *HtTmPlanRaku2s, planIDs []int64) error
	// UpdateChildRateRaku2 子供料金設定を更新
	UpdateChildRateRaku2(child *HtTmChildRateRaku2s) error
	// FetchPlanGroupID plan_idに基づくplan_group_idの取得
	FetchPlanGroupID(planID int64) (int64, error)
	// DeletePlanRaku2 プランを一斉に論理削除
	DeletePlanRaku2(planIDs []int64) error
	// UpdateStopSales プランの売止を更新
	UpdateStopSales(planIDList []int64, isStopSales bool) error
	// UpdateStopSalesByRoomID 部屋に紐づくプランの売止フラグを複数件更新
	UpdateStopSalesByRoomID(roomTypeID int64, isStopSales bool) error
	// UpdateStopSalesByRoomIDList room_type_id(複数)に紐づくプランの売止フラグを複数件更新
	UpdateStopSalesByRoomIDList(roomTypeIDList []int64, isStopSales bool) error
}
