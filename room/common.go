package room

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
)

// RoomTypeTable 部屋テーブル
type RoomTypeTable struct {
	RoomTypeID              int64     `gorm:"primaryKey;autoIncrement:true" json:"room_type_id"`
	PropertyID              int64     `json:"property_id" validate:"required"`
	RoomTypeCode            string    `json:"room_type_code" validate:"required,max=200"`
	Name                    string    `json:"name" validate:"required,max=35"`
	RoomKindID              int64     `json:"room_kind_id"`
	RoomDesc                string    `json:"room_desc" validate:"max=500"`
	StockSettingStart       time.Time `gorm:"type:time" json:"stock_setting_start" validate:"required_without=IsSettingStockYearRound"`
	StockSettingEnd         time.Time `gorm:"type:time" json:"stock_setting_end" validate:"required_without=IsSettingStockYearRound"`
	IsSettingStockYearRound bool      `json:"is_setting_stock_year_round"`
	RoomCount               int16     `json:"room_count" validate:"required"`
	OcuMin                  int       `json:"ocu_min" validate:"required"`
	OcuMax                  int       `json:"ocu_max" validate:"required"`
	IsSmoking               bool      `json:"is_smoking"`
	IsStopSales             bool      `json:"is_stop_sales"`
	IsDelete                bool      `json:"is_delete,omitempty"`
	common.Times            `gorm:"embedded"`
}

// RoomKindTable 部屋種別マスタテーブル
type HtTmRoomKind struct {
	RoomKindID   int64  `gorm:"primaryKey;autoIncrement:true" json:"room_kind_id"`
	KindName     string `json:"kind_name" validate:"required"`
	IsDelete     bool   `json:"is_delete,omitempty"`
	common.Times `gorm:"embedded"`
}

// DeleteInput 部屋削除の入力
type DeleteInput struct {
	RoomTypeID int64 `json:"room_type_id" validate:"required"`
}

// AllAmenitiesOutput アメニティ一覧の出力
type AllAmenitiesOutput struct {
	AmenityID int64  `json:"amenity_id"`
	Name      string `json:"name"`
}

// SaveInput 部屋作成・更新の入力
type SaveInput struct {
	RoomTypeTable
	AmenityIDList []int64                 `json:"amenity_id_list"`
	Images        []image.RoomImagesInput `json:"images"`
}

// ListInput 一覧の入力
type ListInput struct {
	PropertyID int64 `json:"property_id" param:"propertyId" validate:"required"`
	common.Paging
}

// ListOutput 一覧の出力
type ListOutput struct {
	RoomTypeTable
	Href         string   `json:"href"`
	ImageLength  int      `json:"image_length"`
	AmenityNames []string `json:"amenity_names,omitempty"`
}

// DetailInput 詳細の入力
type DetailInput struct {
	PropertyID int64 `json:"property_id" param:"propertyId" validate:"required"`
	RoomTypeID int64 `json:"room_type_id" param:"roomTypeId" validate:"required"`
}

// DetailOutput 詳細の出力
type DetailOutput struct {
	RoomTypeTable
	AmenityIDList []int64                  `json:"amenity_id_list"`
	Images        []image.RoomImagesOutput `json:"images"`
}

// StopSalesInput 売止更新の入力
type StopSalesInput struct {
	RoomTypeID  int64 `json:"room_type_id" validate:"required"`
	IsStopSales bool  `json:"is_stop_sales"`
}

type RoomData struct {
	RoomTypeID              int64                     `gorm:"primaryKey;autoIncrement:true" json:"room_type_id"`
	PropertyID              int64                     `json:"property_id" validate:"required"`
	RoomTypeCode            string                    `json:"room_type_code" validate:"required"`
	Name                    string                    `json:"name" validate:"required"`
	RoomKindID              int64                     `json:"room_kind_id"`
	RoomDesc                string                    `json:"room_desc"`
	StockSettingStart       time.Time                 `gorm:"type:time" json:"stock_setting_start" validate:"required_without=IsSettingStockYearRound"`
	StockSettingEnd         time.Time                 `gorm:"type:time" json:"stock_setting_end" validate:"required_without=IsSettingStockYearRound"`
	IsSettingStockYearRound bool                      `json:"is_setting_stock_year_round"`
	RoomCount               int16                     `json:"room_count"`
	OcuMin                  int                       `json:"ocu_min" validate:"required"`
	OcuMax                  int                       `json:"ocu_max" validate:"required"`
	IsSmoking               bool                      `json:"is_smoking"`
	IsStopSales             bool                      `json:"is_stop_sales"`
	IsDelete                bool                      `json:"is_delete,omitempty"`
	AmenityIDList           []int                     `json:"amenity_id_list" validate:"required"`
	Images                  []Image                   `json:"images" validate:"required"`
	Stocks                  map[string]SaveStockInput `json:"stocks" validate:"required"`
	common.Times            `gorm:"embedded"`
}

// SaveStockInput Create/update inventory
type SaveStockInput struct {
	Stock       int16 `json:"stock"`
	IsStopSales bool  `json:"is_stop_sales"`
}

type Image struct {
	ImageID int    `json:"image_id"`
	Href    string `json:"href"`
	Order   int    `json:"order"`
	Caption string `json:"caption"`
}

// IRoomUsecase 部屋関連のusecaseのインターフェース
type IRoomUsecase interface {
	FetchList(request *ListInput) ([]ListOutput, error)
	FetchAllAmenities() ([]AllAmenitiesOutput, error)
	Create(reuqest *SaveInput) error
	CreateOrUpdateBulk(request []RoomData) error
	FetchDetail(request *DetailInput) (*DetailOutput, error)
	Update(request *SaveInput) error
	Delete(roomTypeID int64) error
	UpdateStopSales(request *StopSalesInput) error
}

// IRoomCommonUsecase 部屋関連の共通usecaseのインターフェース
type IRoomCommonUsecase interface {
	FetchAllRoomKinds() ([]HtTmRoomKind, error)
}

// IRoomCommonRepository 部屋関連の共通repositoryのインターフェース
type IRoomCommonRepository interface {
	FetchAllRoomKinds() ([]HtTmRoomKind, error)
}
