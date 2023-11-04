package image

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmImageNeppans ねっぱんの画像テーブル
type HtTmImageNeppans struct {
	RoomImageNeppanID int64  `gorm:"primaryKey;autoIncrement:true" json:"room_image_neppan_id,omitempty"`
	PropertyID        int64  `json:"property_id,omitempty"`
	Method            string `json:"method,omitempty"`
	Href              string `json:"href,omitempty"`
	CategoryCd        string `json:"category_cd,omitempty"`
	IsMain            bool   `json:"is_main,omitempty"`
	Caption           string `json:"caption,omitempty"`
	MainOrder         int    `json:"main_order"`
	GcsInfo           string `json:"gcs_info,omitempty"`
	common.Times      `gorm:"embedded"`
}

// HtTmRoomOwnImagesNeppans ねっぱんの部屋画像紐付けテーブル
type HtTmRoomOwnImagesNeppans struct {
	RoomOwnImagesID   int64 `gorm:"primaryKey;autoIncrement:true" json:"room_own_images_id,omitempty"`
	RoomTypeID        int64 `json:"room_type_id,omitempty"`
	RoomImageNeppanID int64 `json:"room_image_neppan_id,omitempty"`
	Order             uint8 `json:"order,omitempty"`
	common.Times
}

// HtTmPlanOwnImagesNeppans ねっぱんのプラン画像紐付けテーブル
type HtTmPlanOwnImagesNeppans struct {
	PlanOwnImagesID   int64 `gorm:"primaryKey;autoIncrement:true" json:"plan_own_images_id,omitempty"`
	PlanID            int64 `json:"plan_id,omitempty"`
	RoomImageNeppanID int64 `json:"room_image_neppan_id,omitempty"`
	Order             uint8 `json:"order,omitempty"`
	common.Times
}

// IImageNeppanRepository ねっぱんの画像関連のrepositoryのインターフェース
type IImageNeppanRepository interface {
	common.Repository
	// FetchImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
	FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]RoomImagesOutput, error)
	// FetchImagesByPlanID 複数のplan_idに紐づく複数件の画像を取得
	FetchImagesByPlanID(planIDList []int64) ([]PlanImagesOutput, error)
	// FetchImagesByPropertyID property_idに紐づく複数件の画像を取得
	FetchImagesByPropertyID(propertyID int64) ([]HtTmImageNeppans, error)
	// FetchImageByRoomImageNeppanID room_image_neppan_idに紐づく画像を１件取得
	FetchImageByRoomImageNeppanID(roomImageNeppanID int64) (HtTmImageNeppans, error)
	// FetchRoomImagesByRoomImageNeppanID property_idに紐づく複数件の画像を取得
	FetchRoomImagesByRoomImageNeppanID(roomImageNeppanID int64) ([]HtTmRoomOwnImagesNeppans, error)
	// FetchRoomImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
	FetchRoomImagesByRoomTypeID(roomTypeID int64) ([]HtTmRoomOwnImagesNeppans, error)
	// UpdateImageNeppan room_image_neppan_idに紐づく画像データを変更
	UpdateImageNeppan(request *UpdateInput) error
	// UpdateIsMain room_image_neppan_idに紐づく画像データのisMainフラグを変更
	UpdateIsMain(request *UpdateIsMainInput) error
	// UpdateImageSort 主キーに紐づくsort_numを更新(トランザクションを使用しない)
	UpdateImageSort(request *[]UpdateSortNumInput) error
	// TrUpdateImageSort transactionを使い、主キーに紐づくsort_numを更新
	TrUpdateImageSort(request *[]UpdateSortNumInput) error
	// UpdateRoomImageOrder room_own_images_idに紐づく画像データのorderを更新
	UpdateRoomImageOrder(request *HtTmRoomOwnImagesNeppans) error
	// DeleteImage room_image_neppan_idに紐づく画像データを削除
	DeleteImage(roomImageNeppanID int64) error
	// CreateImage 画像データを作成
	CreateImage(record *HtTmImageNeppans) error
	// ClearRoomImage room_type_idに紐づく画像を削除
	ClearRoomImage(roomTypeID int64) error
	// ClearRoomImageByRoomImageNeppanID room_image_neppan_idに紐づく画像を削除
	ClearRoomImageByRoomImageNeppanID(roomImageNeppanID int64) error
	// CreateRoomOwnImagesNeppan 部屋画像紐付けテーブルにレコード作成
	CreateRoomOwnImagesNeppan(images []HtTmRoomOwnImagesNeppans) error
	// ClearPlanImage plan_idに紐づく画像を削除
	ClearPlanImage(planID int64) error
	// CreatePlanOwnImagesNeppan プラン画像紐付けテーブルにレコード作成
	CreatePlanOwnImagesNeppan(images []HtTmPlanOwnImagesNeppans) error
	// CountMainImagesPerPropertyID property_idに紐づくメイン画像の数をカウント
	CountMainImagesPerPropertyID(propertyID int64) int64
	// FetchMainImages メイン画像を取得
	FetchMainImages(propertyID int64) ([]HtTmImageNeppans, error)
}
