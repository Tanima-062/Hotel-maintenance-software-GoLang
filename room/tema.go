package room

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"gorm.io/gorm"
	"time"
)

// HtTmRoomTypeTemas Room type master table
type HtTmRoomTypeTemas struct {
	RoomTypeTema `gorm:"embedded"`
}

// HtTmRoomTemas HtTmRoomTypeTeams てまの部屋テーブル
type HtTmRoomTemas struct {
	RoomTemaID      int64          `gorm:"primaryKey;autoIncrement:true" json:"room_tema_id"`
	PropertyID      int64          `json:"property_id"`
	RoomTypeCode    string         `json:"room_type_code"`
	RoomNameJa      string         `json:"room_name_ja"`
	RoomDescJa      string         `json:"room_desc_ja"`
	RoomNameEn      string         `json:"room_name_en"`
	RoomDescEn      string         `json:"room_desc_en"`
	RoomSize        int            `json:"room_size"`
	RoomSizeUnit    int            `json:"room_size_unit"`
	MinPax          int            `json:"min_pax"`
	MaxPax          int            `json:"max_pax"`
	RoomCategory    int            `json:"room_category"`
	RoomType        int            `json:"room_type"`
	Grade           int            `json:"grade"`
	RoomView        int            `json:"room_view"`
	Equipment       string         `json:"equipment"`
	Amenity         string         `json:"amenity"`
	Feature         string         `json:"feature"`
	Reason          string         `json:"reason"`
	AcceptCondition string         `json:"accept_condition"`
	Available       int            `json:"avaliable"`
	PictureJson     string         `json:"picture_json"`
	IsStopSales     bool           `json:"is_stop_sales"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
	common.Times
}

// RoomDataTema Room type bulk data
type RoomDataTema struct {
	RoomTypeID              int64                     `gorm:"primaryKey;autoIncrement:true" json:"room_type_id"`
	PropertyID              int64                     `json:"property_id" validate:"required"`
	RoomTypeCode            string                    `json:"room_type_code" validate:"required"`
	Name                    string                    `json:"name" validate:"required"`
	RoomKindID              int64                     `json:"room_kind_id" validate:"required"`
	RoomDesc                string                    `json:"room_desc"`
	StockSettingStart       time.Time                 `gorm:"type:time" json:"stock_setting_start"`
	StockSettingEnd         time.Time                 `gorm:"type:time" json:"stock_setting_end"`
	IsSettingStockYearRound bool                      `json:"is_setting_stock_year_round"`
	RoomCount               int16                     `json:"room_count"`
	OcuMin                  int                       `json:"ocu_min" validate:"required"`
	OcuMax                  int                       `json:"ocu_max" validate:"required"`
	IsSmoking               bool                      `json:"is_smoking"`
	IsStopSales             bool                      `json:"is_stop_sales"`
	IsDelete                bool                      `json:"is_delete"`
	AmenityIDList           []int                     `json:"amenity_id_list" validate:"required"`
	Images                  []Image                   `json:"images" validate:"required"`
	Stocks                  map[string]StockInputTema `json:"stocks" validate:"required"`
	common.Times            `gorm:"embedded"`
}

// StockInputTema stock bulk data
type StockInputTema struct {
	Stock       int16 `json:"stock"`
	IsStopSales bool  `json:"is_stop_sales"`
}

// HtTmRoomAmenityTemas amenity master table
type HtTmRoomAmenityTemas struct {
	TemaRoomAmenityID   int64  `gorm:"primaryKey;autoIncrement:true" json:"tema_room_amenity_id,omitempty"`
	TemaRoomAmenityName string `json:"tema_room_amenity_name"`
	LangCd              string `json:"lang_cd"`
	common.Times        `gorm:"embedded"`
}

// RoomAmenitiesTema room amenity master table
type RoomAmenitiesTema struct {
	RoomTypeID          int64  `json:"room_type_id,omitempty"`
	TemaRoomAmenityName string `json:"tema_room_amenity_name"`
	TemaRoomAmenityID   int64  `json:"tema_room_amenity_id"`
}

// HtTmRoomUseAmenityTemas use amenity table
type HtTmRoomUseAmenityTemas struct {
	RoomTypeID        int64 `json:"room_type_id"`
	TemaRoomAmenityID int64 `json:"tema_room_amenity_id"`
	common.Times      `gorm:"embedded"`
}

// HtTmRoomOwnImagesTemas use image table
type HtTmRoomOwnImagesTemas struct {
	RoomOwnImagesID int64 `gorm:"primaryKey;autoIncrement:true" json:"room_own_images_id,omitempty"`
	RoomTypeID      int64 `json:"room_type_id,omitempty"`
	RoomImageTemaID int64 `json:"room_image_tema_id,omitempty"`
	Order           uint8 `json:"order,omitempty"`
	common.Times
}

// RoomTypeTema master table
type RoomTypeTema struct {
	RoomTypeID              int64     `gorm:"primaryKey;autoIncrement:true" json:"room_type_id"`
	PropertyID              int64     `json:"property_id" validate:"required"`
	RoomTypeCode            string    `json:"room_type_code" validate:"required,max=200"`
	Name                    string    `json:"name" validate:"required,max=35"`
	RoomKindID              int64     `json:"room_kind_id"`
	RoomDesc                string    `json:"room_desc" validate:"max=500"`
	StockSettingStart       time.Time `gorm:"type:time" json:"stock_setting_start"`
	StockSettingEnd         time.Time `gorm:"type:time" json:"stock_setting_end"`
	IsSettingStockYearRound bool      `json:"is_setting_stock_year_round"`
	RoomCount               int16     `json:"room_count"`
	OcuMin                  int       `json:"ocu_min" validate:"required"`
	OcuMax                  int       `json:"ocu_max" validate:"required"`
	IsSmoking               bool      `json:"is_smoking"`
	IsStopSales             bool      `json:"is_stop_sales"`
	IsDelete                bool      `json:"is_delete"`
	common.Times            `gorm:"embedded"`
}

// TemaDetailOutput represents amenity and image data
type TemaDetailOutput struct {
	RoomTypeTema
	AmenityIDList []int64          `json:"amenity_id_list"`
	Images        []RoomImagesTema `json:"images"`
}

// ListOutputTema
type ListOutputTema struct {
	RoomTypeTema
	Href         string           `json:"href" gorm:"column:url"`
	ImageLength  int              `json:"image_length"`
	AmenityNames []string         `json:"amenity_names,omitempty"`
	AmenityIDs   []int64          `json:"amenity_ids,omitempty"`
	Images       []RoomImagesTema `json:"images"`
}

// RoomImagesTema
type RoomImagesTema struct {
	ImageID    int64  `json:"image_id"`
	RoomTypeID int64  `json:"room_type_id"`
	Url        string `json:"href" gorm:"column:url"`
	Title      string `json:"caption"`
	Order      int    `json:"order"`
}

// IRoomTemaRepository てま部屋関連のrepositoryのインターフェース
type IRoomTemaRepository interface {
	common.Repository
	// FetchOne 部屋を１件取得
	FetchOne(roomTypeCode int, propertyID int64) (*HtTmRoomTemas, error)
	// FetchList 部屋を複数件取得
	FetchListWithPropertyId(roomTypeCodeList []int, propertyID int64) ([]HtTmRoomTemas, error)
	// FetchRoomTypeIDByRoomTypeCode get room data by type code
	FetchRoomTypeIDByRoomTypeCode(propertyID int64, roomTypeCode string) (HtTmRoomTypeTemas, error)
	// CreateRoomBulkTema insert room data
	CreateRoomBulkTema(roomTable *HtTmRoomTypeTemas) error
	// UpdateRoomBulkTema update room bulk data
	UpdateRoomBulkTema(roomTable *HtTmRoomTypeTemas) error
	// MatchesRoomTypeIDAndPropertyID match room type and property
	MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool
	// FetchRoomByRoomTypeID Fetch room type  by id
	FetchRoomByRoomTypeID(roomTypeID int64) (*HtTmRoomTypeTemas, error)
	//ClearRoomToAmenities clear use amenities
	ClearRoomToAmenities(roomTypeID int64) error
	// CreateRoomToAmenities create new amenities
	CreateRoomToAmenities(roomTypeID int64, tlRoomAmenityID int64) error
	// ClearRoomImage form use table
	ClearRoomImage(roomTypeID int64) error
	// CreateRoomOwnImages map data in own image table
	CreateRoomOwnImages(images []HtTmRoomOwnImagesTemas) error
	//FetchRoomsByPropertyID fetch room properties
	FetchRoomsByPropertyID(req ListInput) ([]HtTmRoomTypeTemas, error)
	// FetchAllAmenities all list of entities
	FetchAllAmenities() ([]HtTmRoomAmenityTemas, error)
	// FetchAmenitiesByRoomTypeID fetch amenities by room typw
	FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]RoomAmenitiesTema, error)
	// FetchImagesByRoomTypeID fetch image by room type
	FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]RoomImagesTema, error)
}

// IRoomTemaUseCase Tema room related usecase interface
type IRoomTemaUseCase interface {
	// CreateOrUpdateBulk room type bulk data insert
	CreateOrUpdateBulk(request []RoomDataTema) error
	// FetchAllAmenities room type bulk data insert
	FetchAllAmenities() ([]AllAmenitiesOutput, error)
	// FetchDetail FetchDetails to fetch room details
	FetchDetail(request *DetailInput) (*TemaDetailOutput, error)
	// FetchList to fetch room list
	FetchList(request *ListInput) ([]ListOutputTema, error)
}
