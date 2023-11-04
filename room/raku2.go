package room

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmRoomTypeRaku2s らく通の部屋テーブル
type HtTmRoomTypeRaku2s struct {
	RoomTypeTable `gorm:"embedded"`
}

// HtTmRoomUseAmenityRaku2s らく通の部屋と紐づくアメニティテーブル
type HtTmRoomUseAmenityRaku2s struct {
	RoomTypeID         int64 `json:"room_type_id"`
	Raku2RoomAmenityID int64 `json:"raku2_room_amenity_id"`
	common.Times       `gorm:"embedded"`
}

// HtTmRoomAmenityRaku2s らく通のアメニティテーブル
type HtTmRoomAmenityRaku2s struct {
	Raku2RoomAmenityID   int64  `gorm:"primaryKey;autoIncrement:true" json:"raku2_room_amenity_id,omitempty"`
	Raku2RoomAmenityName string `json:"raku2_room_amenity_name"`
	LangCd               string `json:"lang_cd"`
	common.Times         `gorm:"embedded"`
}

// RoomAmenitiesRaku2 アメニティDB取得結果
type RoomAmenitiesRaku2 struct {
	RoomTypeID           int64  `json:"room_type_id,omitempty"`
	Raku2RoomAmenityName string `json:"raku2_room_amenity_name"`
	Raku2RoomAmenityID   int64  `json:"raku2_room_amenity_id"`
}

// IRoomRaku2Repository らく通部屋関連のrepositoryのインターフェース
type IRoomRaku2Repository interface {
	common.Repository
	// FetchRoomsByPropertyID propertyIDに紐づく部屋複数件取得
	FetchRoomsByPropertyID(req ListInput) ([]HtTmRoomTypeRaku2s, error)
	// FetchRoomByRoomTypeID roomTypeIDに紐づく部屋を1件取得
	FetchRoomByRoomTypeID(roomTypeID int64) (*HtTmRoomTypeRaku2s, error)
	// FetchRoomListByRoomTypeID roomTypeIDに紐づく部屋を複数件取得
	FetchRoomListByRoomTypeID(roomTypeIDList []int64) ([]HtTmRoomTypeRaku2s, error)
	// MatchesRoomTypeIDAndPropertyID propertyIDとroomTypeIDが紐付いているか
	MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool
	// FetchAmenitiesByRoomTypeID 部屋に紐づくアメニティを複数件取得
	FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]RoomAmenitiesRaku2, error)
	// FetchAllAmenities 部屋のアメニティを複数件取得
	FetchAllAmenities() ([]HtTmRoomAmenityRaku2s, error)
	// CountRoomTypeCode 部屋コードの重複件数
	CountRoomTypeCode(propertyID int64, roomTypeCode string) int64
	// CreateRoomRaku2 部屋作成
	CreateRoomRaku2(roomTable *HtTmRoomTypeRaku2s) error
	// UpdateRoomRaku2 部屋更新
	UpdateRoomRaku2(roomTable *HtTmRoomTypeRaku2s) error
	// DeleteRoomRaku2 部屋を論理削除
	DeleteRoomRaku2(roomTypeID int64) error
	// ClearRoomToAmenities 部屋に紐づくアメニティを削除
	ClearRoomToAmenities(roomTypeID int64) error
	// CreateRoomToAmenities 部屋に紐づくアメニティを作成
	CreateRoomToAmenities(roomTypeID int64, raku2RoomAmenityID int64) error
	// UpdateStopSales 部屋の売止更新
	UpdateStopSales(roomTypeID int64, isStopSales bool) error
	// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく部屋の売止の更新
	UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, isStopSales bool) error
}
