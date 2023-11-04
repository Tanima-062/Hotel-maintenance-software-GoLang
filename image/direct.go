package image

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmImageDirects 直仕入れの画像テーブル
type HtTmImageDirects struct {
	RoomImageDirectID int64  `gorm:"primaryKey;autoIncrement:true" json:"room_image_direct_id,omitempty"`
	PropertyID        int64  `json:"property_id,omitempty"`
	Method            string `json:"method,omitempty"`
	Href              string `json:"href,omitempty"`
	CategoryCd        string `json:"category_cd,omitempty"`
	IsMain            bool   `json:"is_main"`
	Caption           string `json:"caption,omitempty"`
	MainOrder         int    `json:"main_order"`
	GcsInfo           string `json:"gcs_info,omitempty"`
	common.Times      `gorm:"embedded"`
}

// HtTmRoomOwnImagesDirects 直仕入れの部屋画像紐付けテーブル
type HtTmRoomOwnImagesDirects struct {
	RoomOwnImagesID   int64 `gorm:"primaryKey;autoIncrement:true" json:"room_own_images_id,omitempty"`
	RoomTypeID        int64 `json:"room_type_id,omitempty"`
	RoomImageDirectID int64 `json:"room_image_direct_id,omitempty"`
	Order             uint8 `json:"order"`
	common.Times
}

// HtTmPlanOwnImagesDirects 直仕入れのプラン画像紐付けテーブル
type HtTmPlanOwnImagesDirects struct {
	PlanOwnImagesID   int64 `gorm:"primaryKey;autoIncrement:true" json:"plan_own_images_id,omitempty"`
	PlanID            int64 `json:"plan_id,omitempty"`
	RoomImageDirectID int64 `json:"room_image_direct_id,omitempty"`
	Order             uint8 `json:"order"`
	common.Times
}

// IImageDirectRepository 直仕入れの画像関連のrepositoryのインターフェース
type IImageDirectRepository interface {
	common.Repository
	// FetchImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
	FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]RoomImagesOutput, error)
	// FetchImagesByPlanID 複数のplan_idに紐づく複数件の画像を取得
	FetchImagesByPlanID(planIDList []int64) ([]PlanImagesOutput, error)
	// FetchImagesByPropertyID property_idに紐づく複数件の画像を取得
	FetchImagesByPropertyID(propertyID int64) ([]HtTmImageDirects, error)
	// FetchImageByRoomImageDirectID room_image_direct_idに紐づく画像を１件取得
	FetchImageByRoomImageDirectID(roomImageDirectID int64) (HtTmImageDirects, error)
	// FetchRoomImagesByRoomImageDirectID room_image_direct_idに紐づく画像を複数件取得
	FetchRoomImagesByRoomImageDirectID(roomImageDirectID int64) ([]HtTmRoomOwnImagesDirects, error)
	// FetchRoomImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
	FetchRoomImagesByRoomTypeID(roomTypeID int64) ([]HtTmRoomOwnImagesDirects, error)
	// UpdateImageDirect room_image_direct_idに紐づく画像データを更新
	UpdateImageDirect(request *UpdateInput) error
	// UpdateIsMain room_image_direct_idに紐づく画像データのisMainフラグを更新
	UpdateIsMain(request *UpdateIsMainInput) error
	// UpdateImageSort 主キーに紐づくsort_numを更新(トランザクションを使用しない)
	UpdateImageSort(request *[]UpdateSortNumInput) error
	// TrUpdateImageSort transactionを使い、主キーに紐づくsort_numを更新
	TrUpdateImageSort(request *[]UpdateSortNumInput) error
	// UpdateRoomImageOrder room_own_images_idに紐づく画像データのorderを更新
	UpdateRoomImageOrder(request *HtTmRoomOwnImagesDirects) error
	// DeleteImage room_image_direct_idに紐づく画像データを削除
	DeleteImage(roomImageDirectID int64) error
	// CreateImage 画像データを作成
	CreateImage(record *HtTmImageDirects) error
	// ClearRoomImage room_type_idに紐づく画像を削除
	ClearRoomImage(roomTypeID int64) error
	// ClearRoomImageByRoomImageDirectID room_image_direct_idに紐づく画像を削除
	ClearRoomImageByRoomImageDirectID(roomImageDirectID int64) error
	// CreateRoomOwnImagesDirect 部屋画像紐付けテーブルにレコード作成
	CreateRoomOwnImagesDirect(images []HtTmRoomOwnImagesDirects) error
	// ClearPlanImage plan_idに紐づく画像を削除
	ClearPlanImage(planID int64) error
	// CreatePlanOwnImagesDirect プラン画像紐付けテーブルにレコード作成
	CreatePlanOwnImagesDirect(images []HtTmPlanOwnImagesDirects) error
	// CountMainImagesPerPropertyID property_idに紐づくメイン画像の数をカウント
	CountMainImagesPerPropertyID(propertyID int64) int64
	// FetchMainImages メイン画像を取得
	FetchMainImages(propertyID int64) ([]HtTmImageDirects, error)
}
