package stock

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmStockNeppans ねっぱんの在庫テーブル
type HtTmStockNeppans struct {
	StockTable `gorm:"embedded"`
}

// IStockNeppanRepository ねっぱん在庫関連のrepositoryのインターフェース
type IStockNeppanRepository interface {
	common.Repository
	// FetchAllByRoomTypeIDList room_type_idに紐づく在庫を複数件取得
	FetchAllByRoomTypeIDList(roomTypeIDList []int64, startDate string, endDate string) ([]HtTmStockNeppans, error)
	// FetchAllBookingsByPlanIDList plan_idに紐づく販売数を複数件取得
	FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]BookingCount, error)
	// UpdateStopSales room_type_idに紐づく売止の更新
	UpdateStopSales(roomTypeID int64, useDate string, isStopSales bool) error
	// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく売止の更新
	UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, useDate string, isStopSales bool) error
	// UpsertStocks 在庫の作成・更新
	UpsertStocks(inputData []HtTmStockNeppans) error
}
