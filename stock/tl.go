package stock

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmStockTLs Tl inventory table
type HtTmStockTls struct {
	StockTable `gorm:"embedded"`
}

// IStockTLRepository Tl inventory-related repository interface
type IStockTlRepository interface {
	common.Repository
	// FetchAllByRoomTypeIDList Acquire multiple items of inventory linked to room_type_id
	FetchAllByRoomTypeIDList(roomTypeIDList []int64, startDate string, endDate string) ([]HtTmStockTls, error)
	// FetchAllBookingsByPlanIDList Get multiple sales numbers linked to plan_id
	FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]BookingCount, error)
	// UpdateStopSales Updating the sale stop linked to room_type_id
	UpdateStopSales(roomTypeID int64, useDate string, isStopSales bool) error
	// UpdateStopSalesByRoomTypeIDList Updating sales stop linked to room_type_id (multiple)
	UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, useDate string, isStopSales bool) error
	// UpsertStocks Create/update inventory
	UpsertStocks(inputData []HtTmStockTls) error
	// CreateStocks Create inventory
	CreateStocks(inputData []HtTmStockTls) error
	// FetchBookingCountByRoomTypeId
	FetchBookingCountByRoomTypeId(roomTypeID int64, useDate string) (StockTable, error)
	//UpdateStocks Update multiple inventory
	UpdateStocksBulk(propertyID int64, useDate string, stock int64, bookingCount int64, isStopSale bool) error
}
