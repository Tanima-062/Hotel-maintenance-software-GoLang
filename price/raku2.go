package price

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmPriceRaku2s らく通の料金テーブル
type HtTmPriceRaku2s struct {
	PriceTable `gorm:"embedded"`
}

// HtTmChildRateRaku2s らく通の子供料金設定テーブル
type HtTmChildRateRaku2s struct {
	ChildRateTable `gorm:"embedded"`
}

// IPriceRaku2Repository らく通料金関連のrepositoryのインターフェース
type IPriceRaku2Repository interface {
	common.Repository
	// FetchChildRates プランに紐づく料金設定を複数件取得
	FetchChildRates(planID int64) ([]HtTmChildRateRaku2s, error)
	// FetchAllByPlanIDList 期間内の複数のプランIDに紐づく料金を複数件取得
	FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]HtTmPriceRaku2s, error)
	// FetchPricesByPlanID 本日以降の料金を複数件取得
	FetchPricesByPlanID(planID int64) ([]HtTmPriceRaku2s, error)
	// UpdateChildPrices 子供料金のみ更新
	UpdateChildPrices(inputData []HtTmPriceRaku2s) error
}
