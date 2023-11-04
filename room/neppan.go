package room

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmRoomTypeNeppans ねっぱんの部屋テーブル
type HtTmRoomTypeNeppans struct {
	RoomTypeTable `gorm:"embedded"`
}

// HtTmRoomUseAmenityNeppans ねっぱんの部屋と紐づくアメニティテーブル
type HtTmRoomUseAmenityNeppans struct {
	RoomTypeID          int64 `json:"room_type_id"`
	NeppanRoomAmenityID int64 `json:"neppan_room_amenity_id"`
	common.Times        `gorm:"embedded"`
}

// HtTmRoomAmenityNeppans ねっぱんのアメニティテーブル
type HtTmRoomAmenityNeppans struct {
	NeppanRoomAmenityID   int64  `gorm:"primaryKey;autoIncrement:true" json:"neppan_room_amenity_id,omitempty"`
	NeppanRoomAmenityName string `json:"neppan_room_amenity_name"`
	LangCd                string `json:"lang_cd"`
	common.Times          `gorm:"embedded"`
}

// RoomAmenitiesNeppan アメニティDB取得結果
type RoomAmenitiesNeppan struct {
	RoomTypeID            int64  `json:"room_type_id,omitempty"`
	NeppanRoomAmenityName string `json:"neppan_room_amenity_name"`
	NeppanRoomAmenityID   int64  `json:"neppan_room_amenity_id"`
}

// IRoomNeppanRepository ねっぱん部屋関連のrepositoryのインターフェース
type IRoomNeppanRepository interface {
	common.Repository
	// FetchRoomsByPropertyID propertyIDに紐づく部屋複数件取得
	FetchRoomsByPropertyID(req ListInput) ([]HtTmRoomTypeNeppans, error)
	// FetchRoomByRoomTypeID roomTypeIDに紐づく部屋を1件取得
	FetchRoomByRoomTypeID(roomTypeID int64) (*HtTmRoomTypeNeppans, error)
	// FetchRoomListByRoomTypeID roomTypeIDに紐づく部屋を複数件取得
	FetchRoomListByRoomTypeID(roomTypeIDList []int64) ([]HtTmRoomTypeNeppans, error)
	// MatchesRoomTypeIDAndPropertyID propertyIDとroomTypeIDが紐付いているか
	MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool
	// FetchAmenitiesByRoomTypeID 部屋に紐づくアメニティを複数件取得
	FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]RoomAmenitiesNeppan, error)
	// FetchAllAmenities 部屋のアメニティを複数件取得
	FetchAllAmenities() ([]HtTmRoomAmenityNeppans, error)
	// CountRoomTypeCode 部屋コードの重複件数
	CountRoomTypeCode(propertyID int64, roomTypeCode string) int64
	// CreateRoomNeppan 部屋作成
	CreateRoomNeppan(roomTable *HtTmRoomTypeNeppans) error
	// UpdateRoomNeppan 部屋更新
	UpdateRoomNeppan(roomTable *HtTmRoomTypeNeppans) error
	// DeleteRoomNeppan 部屋を論理削除
	DeleteRoomNeppan(roomTypeID int64) error
	// ClearRoomToAmenities 部屋に紐づくアメニティを削除
	ClearRoomToAmenities(roomTypeID int64) error
	// CreateRoomToAmenities 部屋に紐づくアメニティを作成
	CreateRoomToAmenities(roomTypeID int64, neppanRoomAmenityID int64) error
	// UpdateStopSales 部屋の売止更新
	UpdateStopSales(roomTypeID int64, isStopSales bool) error
	// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく部屋の売止の更新
	UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, isStopSales bool) error
}
