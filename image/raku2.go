package image

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmImageRaku2s らく通の画像テーブル
type HtTmImageRaku2s struct {
	RoomImageRaku2ID int64  `gorm:"primaryKey;autoIncrement:true" json:"room_image_raku2_id,omitempty"`
	PropertyID       int64  `json:"property_id,omitempty"`
	Method           string `json:"method,omitempty"`
	Href             string `json:"href,omitempty"`
	CategoryCd       string `json:"category_cd,omitempty"`
	IsMain           bool   `json:"is_main,omitempty"`
	Caption          string `json:"caption,omitempty"`
	MainOrder        int    `json:"main_order"`
	GcsInfo          string `json:"gcs_info,omitempty"`
	common.Times     `gorm:"embedded"`
}

// HtTmRoomOwnImagesRaku2s らく通の部屋画像紐付けテーブル
type HtTmRoomOwnImagesRaku2s struct {
	RoomOwnImagesID  int64 `gorm:"primaryKey;autoIncrement:true" json:"room_own_images_id,omitempty"`
	RoomTypeID       int64 `json:"room_type_id,omitempty"`
	RoomImageRaku2ID int64 `json:"room_image_raku2_id,omitempty"`
	Order            uint8 `json:"order,omitempty"`
	common.Times
}

// HtTmPlanOwnImagesRaku2s らく通のプラン画像紐付けテーブル
type HtTmPlanOwnImagesRaku2s struct {
	PlanOwnImagesID  int64 `gorm:"primaryKey;autoIncrement:true" json:"plan_own_images_id,omitempty"`
	PlanID           int64 `json:"plan_id,omitempty"`
	RoomImageRaku2ID int64 `json:"room_image_raku2_id,omitempty"`
	Order            uint8 `json:"order,omitempty"`
	common.Times
}

// IImageRaku2Repository らく通の画像関連のrepositoryのインターフェース
type IImageRaku2Repository interface {
	common.Repository
	// FetchImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
	FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]RoomImagesOutput, error)
	// FetchImagesByPlanID 複数のplan_idに紐づく複数件の画像を取得
	FetchImagesByPlanID(planIDList []int64) ([]PlanImagesOutput, error)
	// FetchImagesByPropertyID property_idに紐づく複数件の画像を取得
	FetchImagesByPropertyID(propertyID int64) ([]HtTmImageRaku2s, error)
	// FetchImageByRoomImageRaku2ID room_image_raku2_idに紐づく画像を１件取得
	FetchImageByRoomImageRaku2ID(roomImageRaku2ID int64) (HtTmImageRaku2s, error)
	// FetchRoomImagesByRoomImageRaku2ID room_image_raku2_idに紐づく複数件の画像を取得
	FetchRoomImagesByRoomImageRaku2ID(roomImageRaku2ID int64) ([]HtTmRoomOwnImagesRaku2s, error)
	// FetchRoomImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
	FetchRoomImagesByRoomTypeID(roomTypeID int64) ([]HtTmRoomOwnImagesRaku2s, error)
	// UpdateImageRaku2 room_image_neppan_idに紐づく画像データを変更
	UpdateImageRaku2(request *UpdateInput) error
	// UpdateIsMain room_image_raku2_idに紐づく画像データのisMainフラグを変更
	UpdateIsMain(request *UpdateIsMainInput) error
	// UpdateImageSort 主キーに紐づくsort_numを更新(トランザクションを使用しない)
	UpdateImageSort(request *[]UpdateSortNumInput) error
	// TrUpdateImageSort transactionを使い、主キーに紐づくsort_numを更新
	TrUpdateImageSort(request *[]UpdateSortNumInput) error
	// UpdateRoomImageOrder room_own_images_idに紐づく画像データのorderを更新
	UpdateRoomImageOrder(request *HtTmRoomOwnImagesRaku2s) error
	// DeleteImage room_image_raku2_idに紐づく画像データを削除
	DeleteImage(roomImageRaku2ID int64) error
	// CreateImage 画像データを作成
	CreateImage(record *HtTmImageRaku2s) error
	// ClearRoomImage room_type_idに紐づく画像を削除
	ClearRoomImage(roomTypeID int64) error
	// ClearRoomImageByRoomImageRaku2ID room_image_raku2_idに紐づく画像を削除
	ClearRoomImageByRoomImageRaku2ID(roomImageRaku2ID int64) error
	// CreateRoomOwnImagesRaku2 部屋画像紐付けテーブルにレコード作成
	CreateRoomOwnImagesRaku2(images []HtTmRoomOwnImagesRaku2s) error
	// ClearPlanImage plan_idに紐づく画像を削除
	ClearPlanImage(planID int64) error
	// CreatePlanOwnImagesRaku2 プラン画像紐付けテーブルにレコード作成
	CreatePlanOwnImagesRaku2(images []HtTmPlanOwnImagesRaku2s) error
	// CountMainImagesPerPropertyID property_idに紐づくメイン画像の数をカウント
	CountMainImagesPerPropertyID(propertyID int64) int64
	// FetchMainImages メイン画像を取得
	FetchMainImages(propertyID int64) ([]HtTmImageRaku2s, error)
}
