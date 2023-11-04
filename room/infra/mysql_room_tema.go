package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"gorm.io/gorm"
)

// roomTemaRepository てま部屋関連repository
type roomTemaRepository struct {
	db *gorm.DB
}

// NewRoomTemaRepository インスタンス生成
func NewRoomTemaRepository(db *gorm.DB) room.IRoomTemaRepository {
	return &roomTemaRepository{
		db: db,
	}
}

// TxStart
func (r *roomTemaRepository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit transaction commit
func (r *roomTemaRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback
func (r *roomTemaRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchOne 部屋を１件取得
func (r *roomTemaRepository) FetchOne(roomTypeCode int, propertyID int64) (*room.HtTmRoomTemas, error) {
	result := &room.HtTmRoomTemas{}
	err := r.db.
		Model(&room.HtTmRoomTemas{}).
		Where("room_type_code = ?", roomTypeCode).
		Where("property_id = ?", propertyID).
		First(&result).Error
	return result, err
}

// FetchAllAmenities get room amentities
func (r *roomTemaRepository) FetchAllAmenities() ([]room.HtTmRoomAmenityTemas, error) {
	result := []room.HtTmRoomAmenityTemas{}
	err := r.db.
		Where("lang_cd = ?", "ja-JP").Find(&result).Error
	return result, err
}

// FetchList 部屋を複数件取得
func (r *roomTemaRepository) FetchListWithPropertyId(roomTypeCodeList []int, propertyID int64) ([]room.HtTmRoomTemas, error) {
	result := []room.HtTmRoomTemas{}
	err := r.db.
		Table("ht_tm_room_temas").
		Where("room_type_code IN ?", roomTypeCodeList).
		Where("property_id = ?", propertyID).
		Find(&result).Error
	return result, err
}

// FetchRoomTypeIdByRoomTypeCode Fetch RoomTypeID
func (r *roomTemaRepository) FetchRoomTypeIDByRoomTypeCode(propertyID int64, roomTypeCode string) (room.HtTmRoomTypeTemas, error) {
	result := room.HtTmRoomTypeTemas{}
	err := r.db.
		Table("ht_tm_room_type_temas AS a").
		Where("a.property_id = ?", propertyID).
		Where("a.room_type_code = ?", roomTypeCode).
		Find(&result).Error
	return result, err
}

// FetchRoomsByPropertyID get all room type by property id
func (r *roomTemaRepository) FetchRoomsByPropertyID(req room.ListInput) ([]room.HtTmRoomTypeTemas, error) {
	result := []room.HtTmRoomTypeTemas{}
	query := r.db.
		Table("ht_tm_room_type_temas AS room").
		Where("property_id = ?", req.PropertyID)
	if req.Paging.Limit > 0 {
		query = query.Limit(req.Paging.Limit).Offset(req.Paging.Offset)
	}
	err := query.Find(&result).Error

	return result, err
}

// CreateRoomBulkTema insert a new room into the room type database
func (r *roomTemaRepository) CreateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error {
	return r.db.Create(&roomTable).Error
}

// ClearRoomToAmenities clear room amenities
func (r *roomTemaRepository) ClearRoomToAmenities(roomTypeID int64) error {
	return r.db.Delete(&room.HtTmRoomUseAmenityTemas{}, "room_type_id = ?", roomTypeID).Error
}

// MatchesRoomTypeIDAndPropertyID
func (r *roomTemaRepository) MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool {
	var result int64
	r.db.Model(&room.HtTmRoomTypeTemas{}).
		Where("room_type_id = ?", roomTypeID).
		Where("property_id = ?", propertyID).
		Count(&result)
	return result > 0
}

// FetchRoomByRoomTypeID fetch room type
func (r *roomTemaRepository) FetchRoomByRoomTypeID(roomTypeID int64) (*room.HtTmRoomTypeTemas, error) {
	result := &room.HtTmRoomTypeTemas{}
	err := r.db.
		Table("ht_tm_room_type_temas AS room").
		Where("room_type_id = ? AND is_delete = 0", roomTypeID).
		First(&result).Error
	return result, err
}

// CreateRoomToAmenities insert new amenities
func (r *roomTemaRepository) CreateRoomToAmenities(roomTypeID int64, TemaRoomAmenityID int64) error {
	return r.db.Create(&room.HtTmRoomUseAmenityTemas{
		RoomTypeID:        roomTypeID,
		TemaRoomAmenityID: TemaRoomAmenityID,
		Times:             common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
	}).Error
}

// ClearRoomImage deletes room image
func (r *roomTemaRepository) ClearRoomImage(roomTypeID int64) error {
	return r.db.Delete(&room.HtTmRoomOwnImagesTemas{}, "room_type_id = ?", roomTypeID).Error
}

// CreateRoomOwnImages insert new room own images
func (r *roomTemaRepository) CreateRoomOwnImages(images []room.HtTmRoomOwnImagesTemas) error {
	return r.db.Create(images).Error
}

// UpdateRoomBulkTema update room type
func (r *roomTemaRepository) UpdateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error {
	return r.db.Model(&room.HtTmRoomTypeTemas{}).
		Where("room_type_id = ?", roomTable.RoomTypeID).
		Updates(map[string]interface{}{
			"name":          roomTable.Name,
			"room_kind_id":  roomTable.RoomKindID,
			"room_desc":     roomTable.RoomDesc,
			"ocu_min":       roomTable.OcuMin,
			"ocu_max":       roomTable.OcuMax,
			"is_stop_sales": roomTable.IsStopSales,
			"updated_at":    time.Now(),
		}).Error
}

// FetchAmenitiesByRoomTypeID get room type by room type id
func (r *roomTemaRepository) FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomAmenitiesTema, error) {
	result := []room.RoomAmenitiesTema{}
	err := r.db.
		Select("bind.room_type_id, bind.tema_room_amenity_id, amenity.tema_room_amenity_name").
		Table("ht_tm_room_use_amenity_temas AS bind").
		Joins("INNER JOIN ht_tm_room_amenity_temas AS amenity ON bind.tema_room_amenity_id = amenity.tema_room_amenity_id").
		Where("amenity.lang_cd = ?", "ja-JP").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// FetchImagesByRoomTypeID get multiple images
func (r *roomTemaRepository) FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomImagesTema, error) {
	result := []room.RoomImagesTema{}
	err := r.db.
		Select("image.image_tema_id as image_id, bind.room_type_id, image.url, image.title, bind.order").
		Table("ht_tm_image_temas as image").
		Joins("INNER JOIN ht_tm_room_own_images_temas AS bind ON image.image_tema_id = bind.room_image_tema_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Find(&result).Error
	return result, err
}
