package price

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmPriceDirects 直仕入れの料金テーブル
type HtTmPriceDirects struct {
	PriceTable `gorm:"embedded"`
}

// HtTmChildRateDirects 直仕入れの子供料金設定
type HtTmChildRateDirects struct {
	ChildRateTable `gorm:"embedded"`
}

// IPriceDirectRepository 直仕入れ料金関連のrepositoryのインターフェース
type IPriceDirectRepository interface {
	common.Repository
	// FetchPricesWithInThePeriod 期間内のプランの料金を複数件取得
	FetchPricesWithInThePeriod(planID int64, startDate string, endDate string) ([]HtTmPriceDirects, error)
	// UpsertPrices 料金の作成と更新
	UpsertPrices(inputData []HtTmPriceDirects) error
	// FetchChildRates プランに紐づく料金設定を複数件取得
	FetchChildRates(planID int64) ([]HtTmChildRateDirects, error)
	// FetchChildRatesByPlanIDList 複数プランに紐づく料金設定を複数件取得
	FetchChildRatesByPlanIDList(planIDList []int64) ([]HtTmChildRateDirects, error)
	// FetchAllByPlanIDList 期間内の複数のプランIDに紐づく料金を複数件取得
	FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]HtTmPriceDirects, error)
	// FetchPricesByPlanID 本日以降の料金を複数件取得
	FetchPricesByPlanID(planID int64) ([]HtTmPriceDirects, error)
	// UpdateChildPrices 子供料金のみ更新
	UpdateChildPrices(inputData []HtTmPriceDirects) error
}
