package stock

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"time"
)

// HtTmStockTemas stock master table
type HtTmStockTemas struct {
	StockTableTema `gorm:"embedded"`
}

// StockDataTema stock data
type StockDataTema struct {
	PropertyID   int64                           `json:"property_id" validate:"required"`
	RoomTypeCode string                          `json:"room_type_code" validate:"required"`
	Stocks       map[string]UpdateStockTemaInput `json:"stocks" validate:"required"`
}

type UpdateStockTemaInput struct {
	Stock   int16 `json:"stock"`
	Disable bool  `json:"is_stop_sales"`
}

// StockTableTema stock master table data
type StockTableTema struct {
	StockID      int64     `gorm:"primaryKey;autoIncrement:true" json:"ari_tema_id,omitempty"`
	PropertyID   int64     `json:"property_id"`
	RoomTypeCode int64     `json:"room_type_code,omitempty"`
	AriDate      time.Time `gorm:"type:time" json:"use_date" gorm:"column:ari_date"`
	Stock        int16     `json:"stock"`
	Disable      bool      `json:"is_stop_sales" gorm:"column:disable"`
	common.Times `gorm:"embedded"`
}

// CalendarOutputTema calendar data structure
type CalendarOutputTema struct {
	RoomTypeID              int64               `json:"room_type_id"`
	RoomName                string              `json:"room_name"`
	IsStopSales             bool                `json:"is_stop_sales" gorm:"column:disable"`
	StockSettingStart       string              `json:"stock_setting_start"`
	StockSettingEnd         string              `json:"stock_setting_end"`
	OcuMin                  int                 `json:"ocu_min"`
	OcuMax                  int                 `json:"ocu_max"`
	IsSettingStockYearRound bool                `json:"is_setting_stock_year_round"`
	Stocks                  []CalendarStockTema `json:"stocks"`
	Plans                   []CalendarPlanTema  `json:"plans"`
}

// CalendarStockTema stock data structure
type CalendarStockTema struct {
	RoomTypeCode int64  `json:"room_type_code,omitempty"`
	AriDate      string `json:"use_date" gorm:"column:ari_date"`
	Stock        int16  `json:"stock"`
	IsStopSales  bool   `json:"is_stop_sales" gorm:"column:disable"`
}

// CalendarPlanTema stock and price data structure
type CalendarPlanTema struct {
	PlanID         int64                        `json:"plan_id"`
	PlanName       string                       `json:"plan_name"`
	Disable        bool                         `json:"is_stop_sales" gorm:"column:disable"`
	StockAndPrices map[string]CalendarPriceTema `json:"stock_and_prices"`
}

// CalendarPriceTema calendar price data structure
type CalendarPriceTema struct {
	BookingCount int     `json:"booking_count"`
	Prices       []Price `json:"prices"`
}

// Price 料金データ
type Price struct {
	Type    string `json:"type"`
	Price   int64  `json:"price"`
	Disable bool   `json:"is_stop_sales" gorm:"column:disable"`
}

type IStockTemaRepository interface {
	common.Repository
	// FetchRoomTypeIdByRoomTypeCode fetch room type by room type code
	FetchRoomTypeIdByRoomTypeCode(propertyID int64, roomTypeCode string) (room.HtTmRoomTypeTemas, error)
	// UpdateRoomBulkTema update room bulk data
	UpdateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error
	// FetchBookingCountByRoomTypeId fetch stock detail
	FetchBookingCountByRoomTypeId(roomTypeID string, useDate string) (StockTableTema, error)
	// UpdateStocksBulk update stock detail
	UpdateStocksBulk(propertyCode string, ariDate string, stock int64, disable bool) error
	//CreateStocks create stock detail
	CreateStocks(inputData []HtTmStockTemas) error
	// FetchAllByRoomTypeCodeList FetchAllByRoomTypeIDList fetch all room types
	FetchAllByRoomTypeCodeList(roomTypeCodeList []int64, startDate string, endDate string) ([]HtTmStockTemas, error)
	// FetchAllBookingsByPlanIDList fetch all bookings by plan
	FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]BookingCount, error)
}

type IStockTemaUsecase interface {
	// UpdateBulkTema update stock data
	UpdateBulkTema(request []StockDataTema) error
	// FetchCalendar fetch calender data
	FetchCalendar(hmUser account.HtTmHotelManager, request CalendarInput) (*[]CalendarOutputTema, error)
}
