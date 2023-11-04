package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"gorm.io/gorm"
)

// imageTlRepository
type imageTemaRepository struct {
	db *gorm.DB
}

// NewImageTlRepository
func NewImageTemaRepository(db *gorm.DB) image.IImageTemaRepository {
	return &imageTemaRepository{
		db: db,
	}
}

// CreatePlanOwnImagesTema
func (i *imageTemaRepository) CreatePlanOwnImagesTema(images []image.HtTmPlanOwnImagesTemas) error {
	return i.db.Create(images).Error
}

// FetchImagesByPlanID Get multiple images associated with multiple plan_ids
func (i *imageTemaRepository) FetchImagesByPlanID(planIDList []int64) ([]image.PlanImagesOutput, error) {
	result := []image.PlanImagesOutput{}
	err := i.db.
		Select("image.image_tema_id as image_id, bind.plan_id, image.url as href, image.title as caption, bind.order").
		Table("ht_tm_image_temas as image").
		Joins("INNER JOIN ht_tm_plan_own_images_temas AS bind ON image.image_tema_id = bind.plan_image_tema_id").
		Where("bind.plan_id IN ?", planIDList).
		Find(&result).Error
	return result, err
}

// FetchImagesByRoomTypeID Get multiple images associated with multiple room_type_ids
func (i *imageTemaRepository) FetchRoomImagesByRoomTypeID(roomTypeIDList []int64) ([]image.RoomImagesOutput, error) {
	result := []image.RoomImagesOutput{}
	err := i.db.
		Select("image.image_tema_id as image_id, bind.room_type_id, image.url as href, image.title as caption, bind.order").
		Table("ht_tm_image_temas as image").
		Joins("INNER JOIN ht_tm_room_own_images_temas AS bind ON image.image_tema_id = bind.room_image_tema_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Find(&result).Error
	return result, err
}
