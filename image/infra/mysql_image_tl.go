package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"gorm.io/gorm"
)

// imageTlRepository TL画像関連repository
type imageTlRepository struct {
	db *gorm.DB
}

// NewImageTlRepository インスタンス生成
func NewImageTlRepository(db *gorm.DB) image.IImageTlRepository {
	return &imageTlRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (r *imageTlRepository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (r *imageTlRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (r *imageTlRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchImagesByPropertyID property_idに紐づく複数件の画像を取得
func (i *imageTlRepository) FetchImagesByPropertyID(propertyID int64) ([]image.HtTmImageTls, error) {
	result := []image.HtTmImageTls{}
	err := i.db.Where("property_id = ?", propertyID).Order("updated_at DESC, created_at DESC").Find(&result).Error
	return result, err
}

// FetchImageByImageTlID image_tl_idに紐づく1件の画像を取得
func (i *imageTlRepository) FetchImageByImageTlID(imageTlID int64) (image.HtTmImageTls, error) {
	result := image.HtTmImageTls{}
	err := i.db.Where("image_tl_id = ?", imageTlID).First(&result).Error
	return result, err
}

// UpdateImageTl image_tl_idに紐づく画像データを更新
func (i *imageTlRepository) UpdateImageTl(request *image.UpdateInput) error {
	return i.db.Model(&image.HtTmImageTls{}).
		Where("image_tl_id", request.ImageID).
		Updates(map[string]interface{}{
			"category_cd": request.CategoryCd,
			"caption":     request.Caption,
			"hero_image":  request.IsMain,
			"sort_num":    request.SortNum,
		}).Error
}

// UpdateIsMain image_tl_idに紐づく画像データのisMainフラグを更新
func (i *imageTlRepository) UpdateIsMain(request *image.UpdateIsMainInput) error {
	return i.db.Model(&image.HtTmImageTls{ImageTlID: request.ImageID}).
		Select("hero_image", "sort_num").
		Updates(map[string]interface{}{"hero_image": request.IsMain, "sort_num": *request.SortNum}).
		Error
}

// UpdateImageSort 主キーに紐づくsort_numを更新(トランザクションを使用しない)
func (i *imageTlRepository) UpdateImageSort(request *[]image.UpdateSortNumInput) error {
	for _, value := range *request {
		if err := i.db.Model(&image.HtTmImageTls{}).
			Where("image_tl_id = ?", value.ImageID).
			Updates(map[string]interface{}{"sort_num": value.SortNum}).Error; err != nil {
			return err
		}
	}
	return nil
}

// TrUpdateImageSort transactionを使い、主キーに紐づくsort_numを更新
func (i *imageTlRepository) TrUpdateImageSort(request *[]image.UpdateSortNumInput) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		for _, value := range *request {
			if err := tx.Model(&image.HtTmImageTls{ImageTlID: value.ImageID}).
				Update("sort_num", value.SortNum).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteImage image_tl_idに紐づく画像データを削除
func (i *imageTlRepository) DeleteImage(imageTlID int64) error {
	return i.db.Delete(&image.HtTmImageTls{}, "image_tl_id = ?", imageTlID).Error
}

// CreateImage 画像データを作成
func (i *imageTlRepository) CreateImage(record *image.HtTmImageTls) error {
	return i.db.Create(record).Error
}

// CountMainImagesPerPropertyID property_idに紐づくメイン画像の数をカウント
func (i *imageTlRepository) CountMainImagesPerPropertyID(propertyID int64) int64 {
	var result int64
	i.db.Model(&image.HtTmImageTls{}).
		Where("property_id = ?", propertyID).
		Where("hero_image = ?", true).
		Count(&result)
	return result
}

// FetchMainImages メイン画像を取得
func (i *imageTlRepository) FetchMainImages(propertyID int64) ([]image.HtTmImageTls, error) {
	result := []image.HtTmImageTls{}
	err := i.db.Model(&image.HtTmImageTls{}).
		Where("property_id = ?", propertyID).
		Where("hero_image = ?", true).
		Order("sort_num").
		Find(&result).Error
	return result, err
}

// CreatePlanOwnImagesTL
func (i *imageTlRepository) CreatePlanOwnImagesTl(images []image.HtTmPlanOwnImagesTls) error {
	return i.db.Create(images).Error
}

// ClearRoomImage room_type_id
func (i *imageTlRepository) ClearRoomImage(roomTypeID int64) error {
	return i.db.Delete(&image.HtTmRoomOwnImagesTls{}, "room_type_id = ?", roomTypeID).Error
}

// CreateRoomOwnImagesTl
func (i *imageTlRepository) CreateRoomOwnImagesTl(images []image.HtTmRoomOwnImagesTls) error {
	return i.db.Create(images).Error
}

// ClearPlanImage plan_id
func (i *imageTlRepository) ClearPlanImage(planID int64) error {
	return i.db.Delete(&image.HtTmPlanOwnImagesTls{}, "plan_id = ?", planID).Error
}

// FetchImagesByRoomTypeID Get multiple images associated with multiple room_type_ids
func (i *imageTlRepository) FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]image.RoomImagesOutput, error) {
	result := []image.RoomImagesOutput{}
	err := i.db.
		Select("image.image_tl_id as image_id, bind.room_type_id, image.href, image.caption, bind.order").
		Table("ht_tm_image_tls as image").
		Joins("INNER JOIN ht_tm_room_own_images_tls AS bind ON image.image_tl_id = bind.room_image_tl_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// FetchImagesByPlanID Get multiple images associated with multiple plan_ids
func (i *imageTlRepository) FetchImagesByPlanID(planIDList []int64) ([]image.PlanImagesOutput, error) {
	result := []image.PlanImagesOutput{}
	err := i.db.
		Select("image.image_tl_id as image_id, bind.plan_id, image.href, image.caption, bind.order").
		Table("ht_tm_image_tls as image").
		Joins("INNER JOIN ht_tm_plan_own_images_tls AS bind ON image.image_tl_id = bind.plan_image_tl_id").
		Where("bind.plan_id IN ?", planIDList).
		Find(&result).Error
	return result, err
}
