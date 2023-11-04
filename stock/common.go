package stock

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/price"
)

// StockTable 在庫テーブル
type StockTable struct {
	StockID      int64     `gorm:"primaryKey;autoIncrement:true" json:"price_id,omitempty"`
	RoomTypeID   int64     `json:"room_type_id,omitempty"`
	UseDate      time.Time `gorm:"type:time" json:"use_date"`
	RoomCount    int16     `json:"room_count"`
	BookingCount int16     `json:"booking_count"`
	Stock        int16     `json:"stock"`
	IsStopSales  bool      `json:"is_stop_sales"`
	common.Times `gorm:"embedded"`
}

// CalendarInput 在庫料金カレンダーの入力
type CalendarInput struct {
	BaseDate string `json:"base_date" param:"baseDate"`
}

// CalendarOutput 在庫料金カレンダーの出力
type CalendarOutput struct {
	RoomTypeID              int64           `json:"room_type_id"`
	RoomName                string          `json:"room_name"`
	IsStopSales             bool            `json:"is_stop_sales"`
	StockSettingStart       string          `json:"stock_setting_start"`
	StockSettingEnd         string          `json:"stock_setting_end"`
	OcuMin                  int             `json:"ocu_min"`
	OcuMax                  int             `json:"ocu_max"`
	IsSettingStockYearRound bool            `json:"is_setting_stock_year_round"`
	Stocks                  []CalendarStock `json:"stocks"`
	Plans                   []CalendarPlan  `json:"plans"`
}

// CalendarPlan カレンダーのプラン情報
type CalendarPlan struct {
	PlanID         int64                    `json:"plan_id"`
	PlanName       string                   `json:"plan_name"`
	IsStopSales    bool                     `json:"is_stop_sales"`
	StockAndPrices map[string]CalendarPrice `json:"stock_and_prices"`
}

// CalendarPrice カレンダーの料金情報
type CalendarPrice struct {
	BookingCount int           `json:"booking_count"`
	Prices       []price.Price `json:"prices"`
}

// CalendarStock カレンダーの在庫情報
type CalendarStock struct {
	RoomTypeID   int64  `json:"room_type_id,omitempty"`
	UseDate      string `json:"use_date"`
	RoomCount    int16  `json:"room_count"`
	BookingCount int16  `json:"booking_count"`
	Stock        int16  `json:"stock"`
	IsStopSales  bool   `json:"is_stop_sales"`
}

// BookingCount カレンダーの在庫販売数
type BookingCount struct {
	BookingCount int    `json:"booking_count"`
	PlanID       int64  `json:"plan_id"`
	UseDate      string `json:"use_date"`
}

// StopSalesInput 売止更新
type StopSalesInput struct {
	RoomTypeIDs []int64 `json:"room_type_ids" validate:"required"`
	UseDate     string  `json:"use_date" validate:"required"`
	IsStopSales bool    `json:"is_stop_sales"`
}

// ListInput 在庫一覧の入力
type ListInput struct {
	PropertyID int64  `json:"property_id" param:"propertyId" validate:"required"`
	BaseDate   string `json:"base_date" param:"baseDate" validate:"required"`
}

// ListOutput 在庫一覧の出力
type ListOutput struct {
	RoomTypeID              int64                      `json:"room_type_id"`
	Name                    string                     `json:"name"`
	StockSettingStart       string                     `json:"stock_setting_start"`
	StockSettingEnd         string                     `json:"stock_setting_end"`
	IsSettingStockYearRound bool                       `json:"is_setting_stock_year_round"`
	RoomCount               int16                      `json:"room_count"`
	Stocks                  map[string]ListStockOutput `json:"stocks"`
}

// ListStockOutput 在庫一覧の在庫情報
type ListStockOutput struct {
	RoomCount    int16 `json:"room_count"`
	BookingCount int16 `json:"booking_count"`
}

// SaveInput 在庫作成・更新の入力
type SaveInput struct {
	RoomTypeID int64                     `json:"room_type_id" validate:"required"`
	Stocks     map[string]SaveStockInput `json:"stocks"`
}

// SaveStockInput 作成・更新する在庫数
type SaveStockInput struct {
	RoomCount int16 `json:"room_count"`
}

// StockData update stock
type StockData struct {
	PropertyID              int64                       `json:"property_id" validate:"required"`
	RoomTypeCode            string                      `json:"room_type_code" validate:"required"`
	StockSettingStart       time.Time                   `gorm:"type:time" json:"stock_setting_start,omitempty"`
	StockSettingEnd         time.Time                   `gorm:"type:time" json:"stock_setting_end,omitempty"`
	IsSettingStockYearRound bool                        `json:"is_setting_stock_year_round,omitempty"`
	Stocks                  map[string]UpdateStockInput `json:"stocks" validate:"required"`
}

// UpdateStockInput update inventory
type UpdateStockInput struct {
	Stock       int16 `json:"stock"`
	IsStopSales bool  `json:"is_stop_sales"`
}

// IStockUsecase 在庫関連のusecaseのインターフェース
type IStockUsecase interface {
	FetchCalendar(hmUser account.HtTmHotelManager, request CalendarInput) (*[]CalendarOutput, error)
	UpdateStopSales(request *StopSalesInput) error
	FetchAll(request *ListInput) (*[]ListOutput, error)
	Save(request *[]SaveInput) error
	UpdateBulk(request []StockData) error
}
