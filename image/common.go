package image

import (
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// PlanOwnImagesTable 画像とプランを紐付けるテーブル
type PlanOwnImagesTable struct {
	PlanOwnImagesID int64 `gorm:"primaryKey;autoIncrement:true" json:"plan_own_images_id,omitempty"`
	PlanID          int64 `json:"plan_id,omitempty"`
	ImageID         int64 `json:"image_id,omitempty"`
	Order           uint8 `json:"plan_order"`
	common.Times    `gorm:"embedded"`
}

// ImagesOutput 画像一覧の出力
type ImagesOutput struct {
	ImageID      int64  `json:"image_id"`
	PropertyID   int64  `json:"property_id"`
	Method       string `json:"method,omitempty"`
	Href         string `json:"href"`
	CategoryCd   string `json:"category_cd"`
	IsMain       bool   `json:"is_main"`
	Caption      string `json:"caption"`
	SortNum      int    `json:"sort_num"`
	Order        uint8  `json:"order"`
	GcsInfo      string `json:"gcs_info,omitempty"`
	common.Times `gorm:"embedded"`
}

// RoomImagesOutput 部屋画像の出力
type RoomImagesOutput struct {
	ImageID    int64  `json:"image_id"`
	RoomTypeID int64  `json:"room_type_id"`
	Href       string `json:"href"`
	Caption    string `json:"caption"`
	Order      int    `json:"order"`
}

// PlanImagesOutput プラン画像の出力
type PlanImagesOutput struct {
	ImageID int64  `json:"image_id"`
	PlanID  int64  `json:"plan_id"`
	Href    string `json:"href"`
	Caption string `json:"caption"`
	Order   uint8  `json:"order"`
}

// UploadInput 画像アップロードの入力
type UploadInput struct {
	CategoryCd  string `json:"category_cd"`
	IsMain      bool   `json:"is_main"`
	Caption     string `json:"caption"`
	SortNum     int    `json:"sort_num"`
	ContentType string `json:"content_type"`
}

// UpdateInput 画像情報更新の入力
type UpdateInput struct {
	ImageID    int64  `json:"image_id" validate:"required"`
	CategoryCd string `json:"category_cd"`
	IsMain     bool   `json:"is_main"`
	Caption    string `json:"caption"`
	SortNum    int    `json:"sort_num"`
}

// GcsInfo 画像テーブルのカラムに有るgcsのメタデータ。必要なものだけ記載
type GcsInfo struct {
	Name   string
	Bucket string
}

// ListInput 画像一覧の入力
type ListInput struct {
	PropertyID int64 `json:"property_id" param:"propertyId" validate:"required"`
}

// UpdateIsMainInput メイン画像フラグ更新の入力
type UpdateIsMainInput struct {
	ImageID int64 `json:"image_id" validate:"required"`
	IsMain  bool  `json:"is_main"`
	SortNum *int  `json:"sort_num" validate:"required"`
}

// UpdateSortNumInput メイン画像の並び順更新の入力
type UpdateSortNumInput struct {
	ImageID int64 `json:"image_id" validate:"required"`
	SortNum int   `json:"sort_num" validate:"required"`
}

// DeleteInput 画像削除の入力
type DeleteInput struct {
	ImageID int64 `json:"image_id" param:"imageID" validate:"required"`
}

// RoomImagesInput 部屋作成時の画像データ
type RoomImagesInput struct {
	ImageID    int64 `json:"image_id"`
	RoomTypeID int64 `json:"room_type_id"`
	Order      uint8 `json:"order"`
}

// PlanImagesInput プラン作成時の画像データ
type PlanImagesInput struct {
	ImageID int64 `json:"image_id"`
	PlanID  int64 `json:"plan_id"`
	Order   uint8 `json:"order"`
}

// MainImagesCountInput メイン画像カウントの入力
type MainImagesCountInput struct {
	PropertyID   int64 `json:"property_id" param:"propertyId" validate:"required"`
	WholesalerID int64 `json:"wholesaler_id" query:"wholesaler_id" validate:"required"`
}

// IImageUsecase 画像関連のusecaseのインターフェース
type IImageUsecase interface {
	FetchAll(request *ListInput) ([]ImagesOutput, error)
	Update(request *UpdateInput) error
	UpdateIsMain(request *UpdateIsMainInput) error
	UpdateSortNum(request *[]UpdateSortNumInput) error
	Delete(imageID int64) error
	Create(request *UploadInput, file *multipart.FileHeader, hmUser account.HtTmHotelManager) error
	CountMainImages(propertyID int64) int64
}

// IImageStorage 画像関連のstorageのインターフェース
type IImageStorage interface {
	Delete(bucketName string, objectPath string) error
	Create(bucketName string, filename string, file *multipart.FileHeader) (*storage.ObjectAttrs, error)
}
