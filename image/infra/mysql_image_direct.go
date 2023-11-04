package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"gorm.io/gorm"
)

// imageDirectRepository 直仕入れ画像関連repository
type imageDirectRepository struct {
	db *gorm.DB
}

// NewImageDirectRepository インスタンス生成
func NewImageDirectRepository(db *gorm.DB) image.IImageDirectRepository {
	return &imageDirectRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (r *imageDirectRepository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (r *imageDirectRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (r *imageDirectRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
func (i *imageDirectRepository) FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]image.RoomImagesOutput, error) {
	result := []image.RoomImagesOutput{}
	err := i.db.
		Select("image.room_image_direct_id as image_id, bind.room_type_id, image.href, image.caption, bind.order").
		Table("ht_tm_image_directs as image").
		Joins("INNER JOIN ht_tm_room_own_images_directs AS bind ON image.room_image_direct_id = bind.room_image_direct_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Order("bind.room_type_id, bind.order").
		Find(&result).Error
	return result, err
}

// FetchImagesByPlanID 複数のplan_idに紐づく複数件の画像を取得
func (i *imageDirectRepository) FetchImagesByPlanID(planIDList []int64) ([]image.PlanImagesOutput, error) {
	result := []image.PlanImagesOutput{}
	err := i.db.
		Select("image.room_image_direct_id as image_id, bind.plan_id, image.href, image.caption, bind.order").
		Table("ht_tm_image_directs as image").
		Joins("INNER JOIN ht_tm_plan_own_images_directs AS bind ON image.room_image_direct_id = bind.room_image_direct_id").
		Where("bind.plan_id IN ?", planIDList).
		Order("bind.plan_id, bind.order").
		Find(&result).Error
	return result, err
}

// FetchImagesByPropertyID property_idに紐づく複数件の画像を取得
func (i *imageDirectRepository) FetchImagesByPropertyID(propertyID int64) ([]image.HtTmImageDirects, error) {
	result := []image.HtTmImageDirects{}
	err := i.db.Where("property_id = ?", propertyID).Order("updated_at DESC, created_at DESC").Find(&result).Error
	return result, err
}

// FetchImageByRoomImageDirectID room_image_direct_idに紐づく画像を１件取得
func (i *imageDirectRepository) FetchImageByRoomImageDirectID(roomImageDirectID int64) (image.HtTmImageDirects, error) {
	result := image.HtTmImageDirects{}
	err := i.db.Where("room_image_direct_id = ?", roomImageDirectID).First(&result).Error
	return result, err
}

// FetchRoomImagesByRoomImageDirectID room_image_direct_idに紐づく複数件の画像を取得
func (i *imageDirectRepository) FetchRoomImagesByRoomImageDirectID(roomImageDirectID int64) ([]image.HtTmRoomOwnImagesDirects, error) {
	result := []image.HtTmRoomOwnImagesDirects{}
	err := i.db.Where("room_image_direct_id = ?", roomImageDirectID).Find(&result).Error
	return result, err
}

// FetchRoomImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
func (i *imageDirectRepository) FetchRoomImagesByRoomTypeID(roomTypeID int64) ([]image.HtTmRoomOwnImagesDirects, error) {
	result := []image.HtTmRoomOwnImagesDirects{}
	err := i.db.Where("room_type_id = ?", roomTypeID).Order("'order'").Find(&result).Error
	return result, err
}

// UpdateImageNeppan room_image_neppan_idに紐づく画像データを更新
func (i *imageDirectRepository) UpdateImageDirect(request *image.UpdateInput) error {
	return i.db.Model(&image.HtTmImageDirects{}).
		Where("room_image_direct_id", request.ImageID).
		Updates(map[string]interface{}{
			"category_cd": request.CategoryCd,
			"caption":     request.Caption,
			"is_main":     request.IsMain,
			"main_order":  request.SortNum,
		}).Error
}

// UpdateIsMain room_image_direct_idに紐づく画像データのisMainフラグを更新
func (i *imageDirectRepository) UpdateIsMain(request *image.UpdateIsMainInput) error {
	return i.db.Model(&image.HtTmImageDirects{}).
		Where("room_image_direct_id", request.ImageID).
		Updates(map[string]interface{}{"is_main": request.IsMain, "main_order": *request.SortNum}).
		Error
}

// UpdateImageSort 主キーに紐づくsort_numを更新(トランザクションを使用しない)
func (i *imageDirectRepository) UpdateImageSort(request *[]image.UpdateSortNumInput) error {
	for _, value := range *request {
		if err := i.db.Model(&image.HtTmImageDirects{}).
			Where("room_image_direct_id = ?", value.ImageID).
			Updates(map[string]interface{}{"main_order": value.SortNum}).Error; err != nil {
			return err
		}
	}
	return nil
}

// TrUpdateImageSort transactionを使い、主キーに紐づくsort_numを更新
func (i *imageDirectRepository) TrUpdateImageSort(request *[]image.UpdateSortNumInput) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		for _, value := range *request {
			if err := tx.Model(&image.HtTmImageDirects{}).
				Where("room_image_direct_id", value.ImageID).
				Update("main_order", value.SortNum).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateRoomImageOrder room_own_images_idに紐づく画像データのorderを更新
func (i *imageDirectRepository) UpdateRoomImageOrder(request *image.HtTmRoomOwnImagesDirects) error {
	return i.db.Model(&image.HtTmRoomOwnImagesDirects{}).
		Where("room_own_images_id", request.RoomOwnImagesID).
		Updates(map[string]interface{}{"order": request.Order}).
		Error
}

// DeleteImage room_image_direct_idに紐づく画像データを削除
func (i *imageDirectRepository) DeleteImage(roomImageDirectID int64) error {
	return i.db.Delete(&image.HtTmImageDirects{}, "room_image_direct_id = ?", roomImageDirectID).Error
}

// CreateImage 画像データを作成
func (i *imageDirectRepository) CreateImage(record *image.HtTmImageDirects) error {
	return i.db.Create(record).Error
}

// ClearRoomImage room_type_idに紐づく画像を削除
func (i *imageDirectRepository) ClearRoomImage(roomTypeID int64) error {
	return i.db.Delete(&image.HtTmRoomOwnImagesDirects{}, "room_type_id = ?", roomTypeID).Error
}

// ClearRoomImageByRoomImageDirectID room_image_direct_idに紐づく画像を削除
func (i *imageDirectRepository) ClearRoomImageByRoomImageDirectID(roomImageDirectID int64) error {
	return i.db.Delete(&image.HtTmRoomOwnImagesDirects{}, "room_image_direct_id = ?", roomImageDirectID).Error
}

// CreateRoomOwnImagesDirect 部屋画像紐付けテーブルにレコード作成
func (i *imageDirectRepository) CreateRoomOwnImagesDirect(images []image.HtTmRoomOwnImagesDirects) error {
	return i.db.Create(images).Error
}

// ClearPlanImage plan_idに紐づく画像を削除
func (i *imageDirectRepository) ClearPlanImage(planID int64) error {
	return i.db.Delete(&image.HtTmPlanOwnImagesDirects{}, "plan_id = ?", planID).Error
}

// CreatePlanOwnImagesDirect プラン画像紐付けテーブルにレコード作成
func (i *imageDirectRepository) CreatePlanOwnImagesDirect(images []image.HtTmPlanOwnImagesDirects) error {
	return i.db.Create(images).Error
}

// CountMainImagesPerPropertyID property_idに紐づくメイン画像の数をカウント
func (i *imageDirectRepository) CountMainImagesPerPropertyID(propertyID int64) int64 {
	var result int64
	i.db.Model(&image.HtTmImageDirects{}).
		Where("property_id = ?", propertyID).
		Where("is_main = ?", true).
		Count(&result)
	return result
}

// FetchMainImages メイン画像を取得
func (i *imageDirectRepository) FetchMainImages(propertyID int64) ([]image.HtTmImageDirects, error) {
	result := []image.HtTmImageDirects{}
	err := i.db.Model(&image.HtTmImageDirects{}).
		Where("property_id = ?", propertyID).
		Where("is_main = ?", true).
		Order("main_order").
		Find(&result).Error
	return result, err
}
