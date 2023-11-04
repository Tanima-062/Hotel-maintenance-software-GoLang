package room

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmRoomTypeDirects 直仕入れの部屋テーブル
type HtTmRoomTypeDirects struct {
	RoomTypeTable `gorm:"embedded"`
}

// HtTmRoomUseAmenityDirects 直仕入れの部屋と紐づくアメニティテーブル
type HtTmRoomUseAmenityDirects struct {
	RoomTypeID          int64 `json:"room_type_id"`
	DirectRoomAmenityID int64 `json:"direct_room_amenity_id"`
	common.Times        `gorm:"embedded"`
}

// HtTmRoomAmenityDirects 直仕入れのアメニティテーブル
type HtTmRoomAmenityDirects struct {
	DirectRoomAmenityID   int64  `gorm:"primaryKey;autoIncrement:true" json:"direct_room_amenity_id,omitempty"`
	DirectRoomAmenityName string `json:"direct_room_amenity_name"`
	LangCd                string `json:"lang_cd"`
	common.Times          `gorm:"embedded"`
}

// RoomAmenitiesDirect アメニティDB取得結果
type RoomAmenitiesDirect struct {
	RoomTypeID            int64  `json:"room_type_id,omitempty"`
	DirectRoomAmenityName string `json:"direct_room_amenity_name"`
	DirectRoomAmenityID   int64  `json:"direct_room_amenity_id"`
}

// IRoomDirectRepository 直仕入れ部屋関連のrepositoryのインターフェース
type IRoomDirectRepository interface {
	common.Repository
	// FetchRoomsByPropertyID propertyIDに紐づく部屋複数件取得
	FetchRoomsByPropertyID(req ListInput) ([]HtTmRoomTypeDirects, error)
	// FetchRoomByRoomTypeID roomTypeIDに紐づく部屋を1件取得
	FetchRoomByRoomTypeID(roomTypeID int64) (*HtTmRoomTypeDirects, error)
	// FetchRoomListByRoomTypeID roomTypeIDに紐づく部屋を複数件取得
	FetchRoomListByRoomTypeID(roomTypeIDList []int64) ([]HtTmRoomTypeDirects, error)
	// MatchesRoomTypeIDAndPropertyID propertyIDとroomTypeIDが紐付いているか
	MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool
	// FetchAmenitiesByRoomTypeID 部屋に紐づくアメニティを複数件取得
	FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]RoomAmenitiesDirect, error)
	// FetchAllAmenities 部屋のアメニティを複数件取得
	FetchAllAmenities() ([]HtTmRoomAmenityDirects, error)
	// CountRoomTypeCode 部屋コードの重複件数
	CountRoomTypeCode(propertyID int64, roomTypeCode string) int64
	// CreateRoomDirect 部屋作成
	CreateRoomDirect(roomTable *HtTmRoomTypeDirects) error
	// UpdateRoomDirect 部屋更新
	UpdateRoomDirect(roomTable *HtTmRoomTypeDirects) error
	// DeleteRoomDirect 部屋を論理削除
	DeleteRoomDirect(roomTypeID int64) error
	// ClearRoomToAmenities 部屋に紐づくアメニティを削除
	ClearRoomToAmenities(roomTypeID int64) error
	// CreateRoomToAmenities 部屋に紐づくアメニティを作成
	CreateRoomToAmenities(roomTypeID int64, directRoomAmenityID int64) error
	// UpdateStopSales 部屋の売止更新
	UpdateStopSales(roomTypeID int64, isStopSales bool) error
	// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく部屋の売止の更新
	UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, isStopSales bool) error
}
