package image

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmImageTls TLの画像テーブル
type HtTmImageTls struct {
	ImageTlID    int64  `gorm:"primaryKey;autoIncrement:true" json:"image_tl_id,omitempty"`
	PropertyID   int64  `json:"property_id,omitempty"`
	Links        string `json:"links,omitempty"`
	Method       string `json:"method,omitempty"`
	Href         string `json:"href,omitempty"`
	CategoryCd   string `json:"category_cd,omitempty"`
	HeroImage    bool   `json:"HeroImage"`
	Caption      string `json:"caption,omitempty"`
	SortNum      int    `json:"sort_num"`
	GcsInfo      string `json:"gcs_info,omitempty"`
	common.Times `gorm:"embedded"`
}

// HtTmRoomOwnImagesTls
type HtTmRoomOwnImagesTls struct {
	RoomOwnImagesID int64 `gorm:"primaryKey;autoIncrement:true" json:"room_own_images_id,omitempty"`
	RoomTypeID      int64 `json:"room_type_id,omitempty"`
	RoomImageTlID   int64 `json:"room_image_tl_id,omitempty"`
	Order           uint8 `json:"order,omitempty"`
	common.Times
}

// HtTmPlanOwnImagesTls
type HtTmPlanOwnImagesTls struct {
	PlanOwnImagesID int64 `gorm:"primaryKey;autoIncrement:true" json:"plan_own_images_id,omitempty"`
	PlanID          int64 `json:"plan_id,omitempty"`
	PlanImageTlID   int64 `json:"plan_image_tl_id,omitempty"`
	Order           uint8 `json:"order,omitempty"`
	common.Times
}

// IImageTlRepository TLの画像関連のrepositoryのインターフェース
type IImageTlRepository interface {
	common.Repository
	// FetchImagesByPropertyID property_idに紐づく複数件の画像を取得
	FetchImagesByPropertyID(propertyID int64) ([]HtTmImageTls, error)
	// FetchImageByImageTlID image_tl_idに紐づく1件の画像を取得
	FetchImageByImageTlID(imageTlID int64) (HtTmImageTls, error)
	// UpdateImageTl image_tl_idに紐づく画像データを変更
	UpdateImageTl(request *UpdateInput) error
	// UpdateIsMain image_tl_idに紐づく画像データのisMainフラグを更新
	UpdateIsMain(request *UpdateIsMainInput) error
	// UpdateImageSort 主キーに紐づくsort_numを更新(トランザクションを使用しない)
	UpdateImageSort(request *[]UpdateSortNumInput) error
	// TrUpdateImageSort transactionを使い、主キーに紐づくsort_numを更新
	TrUpdateImageSort(request *[]UpdateSortNumInput) error
	// DeleteImage image_tl_idに紐づく画像データを削除
	DeleteImage(imageTlID int64) error
	// CreateImage 画像データを作成
	CreateImage(record *HtTmImageTls) error
	// CountMainImagesPerPropertyID property_idに紐づくメイン画像の数をカウント
	CountMainImagesPerPropertyID(propertyID int64) int64
	// FetchMainImages メイン画像を取得
	FetchMainImages(propertyID int64) ([]HtTmImageTls, error)
	// CreatePlanOwnImagesTL
	CreatePlanOwnImagesTl(images []HtTmPlanOwnImagesTls) error
	// ClearRoomImage room_type_id
	ClearRoomImage(roomTypeID int64) error
	// CreateRoomOwnImagesTl
	CreateRoomOwnImagesTl(images []HtTmRoomOwnImagesTls) error
	// ClearPlanImage plan_id
	ClearPlanImage(planID int64) error
	// FetchImagesByRoomTypeID room_type_id
	FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]RoomImagesOutput, error)
	// FetchImagesByPlanID Get multiple images associated with multiple plan_ids
	FetchImagesByPlanID(planIDList []int64) ([]PlanImagesOutput, error)
}
