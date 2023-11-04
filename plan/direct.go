package plan

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/price"
)

// HtTmPlanDirects 直仕入れのプランテーブル
type HtTmPlanDirects struct {
	PlanTable `gorm:"embedded"`
}

// HtTmPlanGroupIDDirects 直仕入れのプラングループIDの採番テーブル
type HtTmPlanGroupIDDirects struct {
	PlanGroupID int64 `json:"plan_group_id"`
}

// HtTmChildRateDirects 直仕入れの子供料金設定テーブル
type HtTmChildRateDirects struct {
	price.ChildRateTable `gorm:"embedded"`
}

// IPlanDirectRepository 直仕入れプラン関連のrepositoryのインターフェース
type IPlanDirectRepository interface {
	common.Repository
	// FetchAllByPropertyID 削除されていないproperty_idに紐づくプラン複数件取得
	FetchAllByPropertyID(req ListInput) ([]HtTmPlanDirects, error)
	// FetchAllByRoomTypeID 削除されていないroom_type_idに紐づくプラン複数件取得
	FetchAllByRoomTypeID(roomTypeID int64) ([]HtTmPlanDirects, error)
	// FetchAllByPlanGroupID plan_group_idに紐づくプラン複数件取得
	FetchAllByPlanGroupID(planGroupID int64) ([]HtTmPlanDirects, error)
	// FetchActiveByPlanGroupID plan_group_idに紐づくアクティブなプラン複数件取得
	FetchActiveByPlanGroupID(planGroupID int64) ([]HtTmPlanDirects, error)
	// FetchAllByCancelPolicyID cancel_policy_idに紐づく削除されていないプランを複数件取得
	FetchAllByCancelPolicyID(cancelPolicyID uint64) ([]HtTmPlanDirects, error)
	// FetchOne plan_idに紐づく削除されていないプランを1件取得
	FetchOne(planID int64) (HtTmPlanDirects, error)
	// FetchList plan_idに紐づく削除されていないプランを複数件取得
	FetchList(planIDList []int64) ([]HtTmPlanDirects, error)
	// MatchesPlanIDAndPropertyID propertyIDとplanIDが紐付いているか
	MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool
	// FetchChildRates plan_idに紐づく子供料金設定を複数件取得
	FetchChildRates(planID int64) ([]HtTmChildRateDirects, error)
	// DeletePlanByRoomTypeID room_type_idに紐づくプランを論理削除
	DeletePlanByRoomTypeID(roomTypeID int64) error
	// CheckPlanCode room_type_idとplan_codeの組み合わせで合致するものの件数を取得
	CheckPlanCode(propertyID int64, planCodeList []CheckDuplicatePlanCode) int64
	// MakePlansDirect プランを複数件作成
	MakePlansDirect(planTables []HtTmPlanDirects) error
	// CreatePlansDirect プランを複数件新規作成
	CreatePlansDirect(planTables []HtTmPlanDirects) error
	// CreateChildRateDirect 子供料金を複数件新規作成
	CreateChildRateDirect(childRates []HtTmChildRateDirects) error
	// UpdatePlanDirect プラン情報を更新
	UpdatePlanDirect(planTable *HtTmPlanDirects, planIDs []int64) error
	// UpdateChildRateDirect 子供料金設定を更新
	UpdateChildRateDirect(child *HtTmChildRateDirects) error
	// FetchPlanGroupID plan_idに基づくplan_group_idの取得
	FetchPlanGroupID(planID int64) (int64, error)
	// DeletePlanDirect プランを一斉に論理削除
	DeletePlanDirect(planIDs []int64) error
	// UpdateStopSales プランの売止を更新
	UpdateStopSales(planIDList []int64, isStopSales bool) error
	// UpdateStopSalesByRoomID 部屋に紐づくプランの売止フラグを複数件更新
	UpdateStopSalesByRoomID(roomTypeID int64, isStopSales bool) error
	// UpdateStopSalesByRoomIDList room_type_id(複数)に紐づくプランの売止フラグを複数件更新
	UpdateStopSalesByRoomIDList(roomTypeIDList []int64, isStopSales bool) error
}
