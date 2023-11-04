package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"gorm.io/gorm"
)

// imageRaku2Repository らく通画像関連repository
type imageRaku2Repository struct {
	db *gorm.DB
}

// NewImageRaku2Repository インスタンス生成
func NewImageRaku2Repository(db *gorm.DB) image.IImageRaku2Repository {
	return &imageRaku2Repository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (r *imageRaku2Repository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (r *imageRaku2Repository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (r *imageRaku2Repository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
func (i *imageRaku2Repository) FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]image.RoomImagesOutput, error) {
	result := []image.RoomImagesOutput{}
	err := i.db.
		Select("image.room_image_raku2_id as image_id, bind.room_type_id, image.href, image.caption, bind.order").
		Table("ht_tm_image_raku2s as image").
		Joins("INNER JOIN ht_tm_room_own_images_raku2s AS bind ON image.room_image_raku2_id = bind.room_image_raku2_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Order("bind.room_type_id, bind.order").
		Find(&result).Error
	return result, err
}

// FetchImagesByPlanID 複数のplan_idに紐づく複数件の画像を取得
func (i *imageRaku2Repository) FetchImagesByPlanID(planIDList []int64) ([]image.PlanImagesOutput, error) {
	result := []image.PlanImagesOutput{}
	err := i.db.
		Select("image.room_image_raku2_id as image_id, bind.plan_id, image.href, image.caption, bind.order").
		Table("ht_tm_image_raku2s as image").
		Joins("INNER JOIN ht_tm_plan_own_images_raku2s AS bind ON image.room_image_raku2_id = bind.room_image_raku2_id").
		Where("bind.plan_id IN ?", planIDList).
		Order("bind.plan_id, bind.order").
		Find(&result).Error
	return result, err
}

// FetchImagesByPropertyID property_idに紐づく複数件の画像を取得
func (i *imageRaku2Repository) FetchImagesByPropertyID(propertyID int64) ([]image.HtTmImageRaku2s, error) {
	result := []image.HtTmImageRaku2s{}
	err := i.db.Where("property_id = ?", propertyID).Order("updated_at DESC, created_at DESC").Find(&result).Error
	return result, err
}

// FetchImageByRoomImageRaku2ID room_image_raku2_idに紐づく画像を１件取得
func (i *imageRaku2Repository) FetchImageByRoomImageRaku2ID(roomImageRaku2ID int64) (image.HtTmImageRaku2s, error) {
	result := image.HtTmImageRaku2s{}
	err := i.db.Where("room_image_raku2_id = ?", roomImageRaku2ID).First(&result).Error
	return result, err
}

// FetchRoomImagesByRoomImageRaku2ID room_image_raku2_idに紐づく複数件の画像を取得
func (i *imageRaku2Repository) FetchRoomImagesByRoomImageRaku2ID(roomImageRaku2ID int64) ([]image.HtTmRoomOwnImagesRaku2s, error) {
	result := []image.HtTmRoomOwnImagesRaku2s{}
	err := i.db.Where("room_image_raku2_id = ?", roomImageRaku2ID).Find(&result).Error
	return result, err
}

// FetchRoomImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
func (i *imageRaku2Repository) FetchRoomImagesByRoomTypeID(roomTypeID int64) ([]image.HtTmRoomOwnImagesRaku2s, error) {
	result := []image.HtTmRoomOwnImagesRaku2s{}
	err := i.db.Where("room_type_id = ?", roomTypeID).Order("'order'").Find(&result).Error
	return result, err
}

// UpdateImageRaku2 room_image_raku2_idに紐づく画像データを更新
func (i *imageRaku2Repository) UpdateImageRaku2(request *image.UpdateInput) error {
	return i.db.Model(&image.HtTmImageRaku2s{}).
		Where("room_image_raku2_id", request.ImageID).
		Updates(map[string]interface{}{
			"category_cd": request.CategoryCd,
			"caption":     request.Caption,
			"is_main":     request.IsMain,
			"main_order":  request.SortNum,
		}).Error
}

// UpdateIsMain room_image_raku2_idに紐づく画像データのisMainフラグを変更
func (i *imageRaku2Repository) UpdateIsMain(request *image.UpdateIsMainInput) error {
	return i.db.Model(&image.HtTmImageRaku2s{}).
		Where("room_image_raku2_id", request.ImageID).
		Updates(map[string]interface{}{"is_main": request.IsMain, "main_order": *request.SortNum}).
		Error
}

// UpdateImageSort 主キーに紐づくsort_numを更新(トランザクションを使用しない)
func (i *imageRaku2Repository) UpdateImageSort(request *[]image.UpdateSortNumInput) error {
	for _, value := range *request {
		if err := i.db.Model(&image.HtTmImageRaku2s{}).
			Where("room_image_raku2_id = ?", value.ImageID).
			Updates(map[string]interface{}{"main_order": value.SortNum}).Error; err != nil {
			return err
		}
	}
	return nil
}

// TrUpdateImageSort transactionを使い、主キーに紐づくsort_numを更新
func (i *imageRaku2Repository) TrUpdateImageSort(request *[]image.UpdateSortNumInput) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		for _, value := range *request {
			if err := tx.Model(&image.HtTmImageRaku2s{}).
				Where("room_image_raku2_id", value.ImageID).
				Update("main_order", value.SortNum).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateRoomImageOrder room_own_images_idに紐づく画像データのorderを更新
func (i *imageRaku2Repository) UpdateRoomImageOrder(request *image.HtTmRoomOwnImagesRaku2s) error {
	return i.db.Model(&image.HtTmRoomOwnImagesRaku2s{}).
		Where("room_own_images_id", request.RoomOwnImagesID).
		Updates(map[string]interface{}{"order": request.Order}).
		Error
}

// DeleteImage room_image_raku2_idに紐づく画像データを削除
func (i *imageRaku2Repository) DeleteImage(roomImageRaku2ID int64) error {
	return i.db.Delete(&image.HtTmImageRaku2s{}, "room_image_raku2_id = ?", roomImageRaku2ID).Error
}

// CreateImage 画像データを作成
func (i *imageRaku2Repository) CreateImage(record *image.HtTmImageRaku2s) error {
	return i.db.Create(record).Error
}

// ClearRoomImage room_type_idに紐づく画像を削除
func (i *imageRaku2Repository) ClearRoomImage(roomTypeID int64) error {
	return i.db.Delete(&image.HtTmRoomOwnImagesRaku2s{}, "room_type_id = ?", roomTypeID).Error
}

// ClearRoomImageByRoomImageRaku2ID room_image_raku2_idに紐づく画像を削除
func (i *imageRaku2Repository) ClearRoomImageByRoomImageRaku2ID(roomImageRaku2ID int64) error {
	return i.db.Delete(&image.HtTmRoomOwnImagesRaku2s{}, "room_image_raku2_id = ?", roomImageRaku2ID).Error
}

// CreateRoomOwnImagesRaku2 部屋画像紐付けテーブルにレコード作成
func (i *imageRaku2Repository) CreateRoomOwnImagesRaku2(images []image.HtTmRoomOwnImagesRaku2s) error {
	return i.db.Create(images).Error
}

// ClearPlanImage plan_idに紐づく画像を削除
func (i *imageRaku2Repository) ClearPlanImage(planID int64) error {
	return i.db.Delete(&image.HtTmPlanOwnImagesRaku2s{}, "plan_id = ?", planID).Error
}

// CreatePlanOwnImagesRaku2 プラン画像紐付けテーブルにレコード作成
func (i *imageRaku2Repository) CreatePlanOwnImagesRaku2(images []image.HtTmPlanOwnImagesRaku2s) error {
	return i.db.Create(images).Error
}

// CountMainImagesPerPropertyID property_idに紐づくメイン画像の数をカウント
func (i *imageRaku2Repository) CountMainImagesPerPropertyID(propertyID int64) int64 {
	var result int64
	i.db.Model(&image.HtTmImageRaku2s{}).
		Where("property_id = ?", propertyID).
		Where("is_main = ?", true).
		Count(&result)
	return result
}

// FetchMainImages メイン画像を取得
func (i *imageRaku2Repository) FetchMainImages(propertyID int64) ([]image.HtTmImageRaku2s, error) {
	result := []image.HtTmImageRaku2s{}
	err := i.db.Model(&image.HtTmImageRaku2s{}).
		Where("property_id = ?", propertyID).
		Where("is_main = ?", true).
		Order("main_order").
		Find(&result).Error
	return result, err
}
