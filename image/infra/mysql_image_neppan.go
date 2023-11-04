package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"gorm.io/gorm"
)

// imageNeppanRepository ねっぱん画像関連repository
type imageNeppanRepository struct {
	db *gorm.DB
}

// NewImageNeppanRepository インスタンス生成
func NewImageNeppanRepository(db *gorm.DB) image.IImageNeppanRepository {
	return &imageNeppanRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (r *imageNeppanRepository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (r *imageNeppanRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (r *imageNeppanRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchImagesByRoomTypeID room_type_idに紐づく複数件の画像を取得
func (i *imageNeppanRepository) FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]image.RoomImagesOutput, error) {
	result := []image.RoomImagesOutput{}
	err := i.db.
		Select("image.room_image_neppan_id as image_id, bind.room_type_id, image.href, image.caption, bind.order").
		Table("ht_tm_image_neppans as image").
		Joins("INNER JOIN ht_tm_room_own_images_neppans AS bind ON image.room_image_neppan_id = bind.room_image_neppan_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Order("bind.room_type_id, bind.order").
		Find(&result).Error
	return result, err
}

// FetchImagesByPlanID 複数のplan_idに紐づく複数件の画像を取得
func (i *imageNeppanRepository) FetchImagesByPlanID(planIDList []int64) ([]image.PlanImagesOutput, error) {
	result := []image.PlanImagesOutput{}
	err := i.db.
		Select("image.room_image_neppan_id as image_id, bind.plan_id, image.href, image.caption, bind.order").
		Table("ht_tm_image_neppans as image").
		Joins("INNER JOIN ht_tm_plan_own_images_neppans AS bind ON image.room_image_neppan_id = bind.room_image_neppan_id").
		Where("bind.plan_id IN ?", planIDList).
		Order("bind.plan_id, bind.order").
		Find(&result).Error
	return result, err
}

// FetchImagesByPropertyID property_idに紐づく複数件の画像を取得
func (i *imageNeppanRepository) FetchImagesByPropertyID(propertyID int64) ([]image.HtTmImageNeppans, error) {
	result := []image.HtTmImageNeppans{}
	err := i.db.Where("property_id = ?", propertyID).Order("updated_at DESC, created_at DESC").Find(&result).Error
	return result, err
}

// FetchImageByRoomImageNeppanID room_image_neppan_idに紐づく画像を１件取得
func (i *imageNeppanRepository) FetchImageByRoomImageNeppanID(roomImageNeppanID int64) (image.HtTmImageNeppans, error) {
	result := image.HtTmImageNeppans{}
	err := i.db.Where("room_image_neppan_id = ?", roomImageNeppanID).First(&result).Error
	return result, err
}

// FetchRoomImagesByRoomImageNeppanID property_idに紐づく複数件の画像を取得
func (i *imageNeppanRepository) FetchRoomImagesByRoomImageNeppanID(roomImageNeppanID int64) ([]image.HtTmRoomOwnImagesNeppans, error) {
	result := []image.HtTmRoomOwnImagesNeppans{}
	err := i.db.Where("room_image_neppan_id = ?", roomImageNeppanID).Find(&result).Error
	return result, err
}

// FetchRoomImagesByRoomImageDirectID room_type_idに紐づく複数件の画像を取得
func (i *imageNeppanRepository) FetchRoomImagesByRoomTypeID(roomTypeID int64) ([]image.HtTmRoomOwnImagesNeppans, error) {
	result := []image.HtTmRoomOwnImagesNeppans{}
	err := i.db.Where("room_type_id = ?", roomTypeID).Order("'order'").Find(&result).Error
	return result, err
}

// UpdateImageNeppan room_image_neppan_idに紐づく画像データを更新
func (i *imageNeppanRepository) UpdateImageNeppan(request *image.UpdateInput) error {
	return i.db.Model(&image.HtTmImageNeppans{}).
		Where("room_image_neppan_id", request.ImageID).
		Updates(map[string]interface{}{
			"category_cd": request.CategoryCd,
			"caption":     request.Caption,
			"is_main":     request.IsMain,
			"main_order":  request.SortNum,
		}).Error
}

// UpdateIsMain room_image_neppan_idに紐づく画像データのisMainフラグを変更
func (i *imageNeppanRepository) UpdateIsMain(request *image.UpdateIsMainInput) error {
	return i.db.Model(&image.HtTmImageNeppans{}).
		Where("room_image_neppan_id", request.ImageID).
		Updates(map[string]interface{}{"is_main": request.IsMain, "main_order": *request.SortNum}).
		Error
}

// UpdateImageSort 主キーに紐づくsort_numを更新(トランザクションを使用しない)
func (i *imageNeppanRepository) UpdateImageSort(request *[]image.UpdateSortNumInput) error {
	for _, value := range *request {
		if err := i.db.Model(&image.HtTmImageNeppans{}).
			Where("room_image_neppan_id = ?", value.ImageID).
			Updates(map[string]interface{}{"main_order": value.SortNum}).Error; err != nil {
			return err
		}
	}
	return nil
}

// TrUpdateImageSort transactionを使い、主キーに紐づくsort_numを更新
func (i *imageNeppanRepository) TrUpdateImageSort(request *[]image.UpdateSortNumInput) error {
	return i.db.Transaction(func(tx *gorm.DB) error {
		for _, value := range *request {
			if err := tx.Model(&image.HtTmImageNeppans{}).
				Where("room_image_neppan_id", value.ImageID).
				Update("main_order", value.SortNum).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateRoomImageOrder room_own_images_idに紐づく画像データのorderを更新
func (i *imageNeppanRepository) UpdateRoomImageOrder(request *image.HtTmRoomOwnImagesNeppans) error {
	return i.db.Model(&image.HtTmRoomOwnImagesNeppans{}).
		Where("room_own_images_id", request.RoomOwnImagesID).
		Updates(map[string]interface{}{"order": request.Order}).
		Error
}

// DeleteImage room_image_neppan_idに紐づく画像データを削除
func (i *imageNeppanRepository) DeleteImage(roomImageNeppanID int64) error {
	return i.db.Delete(&image.HtTmImageNeppans{}, "room_image_neppan_id = ?", roomImageNeppanID).Error
}

// CreateImage 画像データを作成
func (i *imageNeppanRepository) CreateImage(record *image.HtTmImageNeppans) error {
	return i.db.Create(record).Error
}

// ClearRoomImage room_type_idに紐づく画像を削除
func (i *imageNeppanRepository) ClearRoomImage(roomTypeID int64) error {
	return i.db.Delete(&image.HtTmRoomOwnImagesNeppans{}, "room_type_id = ?", roomTypeID).Error
}

// ClearRoomImageByRoomImageNeppanID room_image_neppan_idに紐づく画像を削除
func (i *imageNeppanRepository) ClearRoomImageByRoomImageNeppanID(roomImageNeppanID int64) error {
	return i.db.Delete(&image.HtTmRoomOwnImagesNeppans{}, "room_image_neppan_id = ?", roomImageNeppanID).Error
}

// CreateRoomOwnImagesNeppan 部屋画像紐付けテーブルにレコード作成
func (i *imageNeppanRepository) CreateRoomOwnImagesNeppan(images []image.HtTmRoomOwnImagesNeppans) error {
	return i.db.Create(images).Error
}

// ClearPlanImage plan_idに紐づく画像を削除
func (i *imageNeppanRepository) ClearPlanImage(planID int64) error {
	return i.db.Delete(&image.HtTmPlanOwnImagesNeppans{}, "plan_id = ?", planID).Error
}

// CreatePlanOwnImagesNeppan プラン画像紐付けテーブルにレコード作成
func (i *imageNeppanRepository) CreatePlanOwnImagesNeppan(images []image.HtTmPlanOwnImagesNeppans) error {
	return i.db.Create(images).Error
}

// CountMainImagesPerPropertyID property_idに紐づくメイン画像の数をカウント
func (i *imageNeppanRepository) CountMainImagesPerPropertyID(propertyID int64) int64 {
	var result int64
	i.db.Model(&image.HtTmImageNeppans{}).
		Where("property_id = ?", propertyID).
		Where("is_main = ?", true).
		Count(&result)
	return result
}

// FetchMainImages メイン画像を取得
func (i *imageNeppanRepository) FetchMainImages(propertyID int64) ([]image.HtTmImageNeppans, error) {
	result := []image.HtTmImageNeppans{}
	err := i.db.Model(&image.HtTmImageNeppans{}).
		Where("property_id = ?", propertyID).
		Where("is_main = ?", true).
		Order("main_order").
		Find(&result).Error
	return result, err
}
