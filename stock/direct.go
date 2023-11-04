package stock

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmStockDirects 直仕入れの在庫テーブル
type HtTmStockDirects struct {
	StockTable `gorm:"embedded"`
}

// IStockDirectRepository 直仕入れ在庫関連のrepositoryのインターフェース
type IStockDirectRepository interface {
	common.Repository
	// FetchAllByRoomTypeIDList room_type_idに紐づく在庫を複数件取得
	FetchAllByRoomTypeIDList(roomTypeIDList []int64, startDate string, endDate string) ([]HtTmStockDirects, error)
	// FetchAllBookingsByPlanIDList plan_idに紐づく販売数を複数件取得
	FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]BookingCount, error)
	// FetchStocksByRoomTypeIDList room_type_idに紐づく本日以降の在庫を複数件取得
	FetchStocksByRoomTypeIDList(roomTypeIDList []int64) ([]HtTmStockDirects, error)
	// UpdateStopSales room_type_idに紐づく売止の更新
	UpdateStopSales(roomTypeID int64, useDate string, isStopSales bool) error
	// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく売止の更新
	UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, useDate string, isStopSales bool) error
	// UpsertStocks 在庫の作成・更新
	UpsertStocks(inputData []HtTmStockDirects) error
	// FetchStock 日付とroom_type_idに紐づく在庫を１件取得
	FetchStock(roomTypeID int64, useDate string) (*HtTmStockDirects, error)
	// CreateStocks 在庫を複数件作成
	CreateStocks(inputData []HtTmStockDirects) error
}
