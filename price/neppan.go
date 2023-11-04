package price

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmPriceNeppans ねっぱんの料金テーブル
type HtTmPriceNeppans struct {
	PriceTable `gorm:"embedded"`
}

// HtTmChildRateNeppans ねっぱんの子供料金設定テーブル
type HtTmChildRateNeppans struct {
	ChildRateTable `gorm:"embedded"`
}

// IPriceNeppanRepository ねっぱん料金関連のrepositoryのインターフェース
type IPriceNeppanRepository interface {
	common.Repository
	// FetchChildRates プランに紐づく料金設定を複数件取得
	FetchChildRates(planID int64) ([]HtTmChildRateNeppans, error)
	// FetchAllByPlanIDList 期間内の複数のプランIDに紐づく料金を複数件取得
	FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]HtTmPriceNeppans, error)
	// FetchPricesByPlanID 本日以降の料金を複数件取得
	FetchPricesByPlanID(planID int64) ([]HtTmPriceNeppans, error)
	// UpdateChildPrices 子供料金のみ更新
	UpdateChildPrices(inputData []HtTmPriceNeppans) error
}
