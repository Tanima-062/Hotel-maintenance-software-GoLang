package room

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
)

// HtTmRoomTypeTls
type HtTmRoomTypeTls struct {
	RoomTypeTable `gorm:"embedded"`
}

// HtTmRoomUseAmenityTls
type HtTmRoomUseAmenityTls struct {
	RoomTypeID       int64 `json:"room_type_id"`
	TlsRoomAmenityID int64 `json:"tls_room_amenity_id"`
	common.Times     `gorm:"embedded"`
}

// HtTmRoomAmenityTls
type HtTmRoomAmenityTls struct {
	TlsRoomAmenityID   int64  `gorm:"primaryKey;autoIncrement:true" json:"tls_room_amenity_id,omitempty"`
	TlsRoomAmenityName string `json:"tls_room_amenity_name"`
	LangCd             string `json:"lang_cd"`
	common.Times       `gorm:"embedded"`
}

// RoomAmenitiesTl
type RoomAmenitiesTl struct {
	RoomTypeID         int64  `json:"room_type_id,omitempty"`
	TlsRoomAmenityName string `json:"tls_room_amenity_name"`
	TlsRoomAmenityID   int64  `json:"tls_room_amenity_id"`
}

// RoomKindTable 部屋種別マスタテーブル
type HtTmRoomKindsTls struct {
	RoomKindID   int64  `gorm:"primaryKey;autoIncrement:true" json:"room_kind_id"`
	KindName     string `json:"kind_name" validate:"required"`
	IsDelete     bool   `json:"is_delete,omitempty"`
	common.Times `gorm:"embedded"`
}

type ListOutputTl struct {
	RoomTypeTable
	Href         string                   `json:"href"`
	ImageLength  int                      `json:"image_length"`
	AmenityNames []string                 `json:"amenity_names,omitempty"`
	AmenityIDs   []int64                  `json:"amenity_ids,omitempty"`
	Images       []image.RoomImagesOutput `json:"images"`
}

type IRoomBulkUsecase interface {
	FetchList(request *ListInput) ([]ListOutputTl, error)
	CreateOrUpdateBulk(request []RoomData) error
	FetchDetail(request *DetailInput) (*DetailOutput, error)
	FetchAllAmenities() ([]AllAmenitiesOutput, error)
}

// IRoomTLRepository
type IRoomTlRepository interface {
	common.Repository
	// FetchRoomsByPropertyID
	FetchRoomsByPropertyID(req ListInput) ([]HtTmRoomTypeTls, error)
	// FetchRoomByRoomTypeID
	FetchRoomByRoomTypeID(roomTypeID int64) (*HtTmRoomTypeTls, error)
	// FetchRoomListByRoomTypeID
	FetchRoomListByRoomTypeID(roomTypeIDList []int64) ([]HtTmRoomTypeTls, error)
	// MatchesRoomTypeIDAndPropertyID
	MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool
	// FetchAmenitiesByRoomTypeID
	FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]RoomAmenitiesTl, error)
	// FetchAllAmenities
	FetchAllAmenities() ([]HtTmRoomAmenityTls, error)
	// CountRoomTypeCode
	CountRoomTypeCode(propertyID int64, roomTypeCode string) int64
	// FetchRoomTypeIdByRoomTypeCode to get room type data
	FetchRoomTypeIdByRoomTypeCode(propertyID int64, roomTypeCode string) (HtTmRoomTypeTls, error)
	// CreateRoomTl
	CreateRoomTl(roomTable *HtTmRoomTypeTls) error
	// UpdateRoomTl
	UpdateRoomTl(roomTable *HtTmRoomTypeTls) error
	// UpdateRoomBulkTl
	UpdateRoomBulkTl(roomTable *HtTmRoomTypeTls) error
	// CreateRoomBulkTl
	CreateRoomBulkTl(roomTable *HtTmRoomTypeTls) error
	// DeleteRoomTl
	DeleteRoomTl(roomTypeID int64) error
	// ClearRoomToAmenities
	ClearRoomToAmenities(roomTypeID int64) error
	// CreateRoomToAmenities
	CreateRoomToAmenities(roomTypeID int64, tlRoomAmenityID int64) error
	// UpdateStopSales
	UpdateStopSales(roomTypeID int64, isStopSales bool) error
	// UpdateStopSalesByRoomTypeIDList
	UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, isStopSales bool) error
	FetchAllRoomKindTls() ([]HtTmRoomKindsTls, error)
}
