package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"gorm.io/gorm"
)

// roomDirectRepository 直仕入れ部屋関連repository
type roomDirectRepository struct {
	db *gorm.DB
}

// NewRoomDirectRepository インスタンス生成
func NewRoomDirectRepository(db *gorm.DB) room.IRoomDirectRepository {
	return &roomDirectRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (r *roomDirectRepository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (r *roomDirectRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (r *roomDirectRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchRoomsByPropertyID propertyIDに紐づく部屋複数件取得
func (r *roomDirectRepository) FetchRoomsByPropertyID(req room.ListInput) ([]room.HtTmRoomTypeDirects, error) {
	result := []room.HtTmRoomTypeDirects{}
	query := r.db.
		Table("ht_tm_room_type_directs AS room").
		Where("property_id = ? AND is_delete = 0", req.PropertyID)
	if req.Paging.Limit > 0 {
		query = query.Limit(req.Paging.Limit).Offset(req.Paging.Offset)
	}
	err := query.Find(&result).Error

	return result, err
}

// FetchRoomByRoomTypeID roomTypeIDに紐づく部屋を1件取得
func (r *roomDirectRepository) FetchRoomByRoomTypeID(roomTypeID int64) (*room.HtTmRoomTypeDirects, error) {
	result := &room.HtTmRoomTypeDirects{}
	err := r.db.
		Table("ht_tm_room_type_directs AS room").
		Where("room_type_id = ? AND is_delete = 0", roomTypeID).
		First(&result).Error
	return result, err
}

// FetchRoomListByRoomTypeID roomTypeIDに紐づく部屋を複数件取得
func (r *roomDirectRepository) FetchRoomListByRoomTypeID(roomTypeIDList []int64) ([]room.HtTmRoomTypeDirects, error) {
	result := []room.HtTmRoomTypeDirects{}
	err := r.db.
		Table("ht_tm_room_type_directs AS room").
		Where("room_type_id IN ? AND is_delete = 0", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// MatchesRoomTypeIDAndPropertyID propertyIDとroomTypeIDが紐付いているか
func (r *roomDirectRepository) MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool {
	var result int64
	r.db.Model(&room.HtTmRoomTypeDirects{}).
		Where("room_type_id = ?", roomTypeID).
		Where("property_id = ?", propertyID).
		Count(&result)
	return result > 0
}

// FetchAmenitiesByRoomTypeID 部屋に紐づくアメニティを複数件取得
func (r *roomDirectRepository) FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomAmenitiesDirect, error) {
	result := []room.RoomAmenitiesDirect{}
	err := r.db.
		Select("bind.room_type_id, bind.direct_room_amenity_id, amenity.direct_room_amenity_name").
		Table("ht_tm_room_use_amenity_directs AS bind").
		Joins("INNER JOIN ht_tm_room_amenity_directs AS amenity ON bind.direct_room_amenity_id = amenity.direct_room_amenity_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// FetchAllAmenities 部屋のアメニティを複数件取得
func (r *roomDirectRepository) FetchAllAmenities() ([]room.HtTmRoomAmenityDirects, error) {
	result := []room.HtTmRoomAmenityDirects{}
	err := r.db.Find(&result).Error
	return result, err
}

// CountRoomTypeCode 部屋コードの重複件数
func (r *roomDirectRepository) CountRoomTypeCode(propertyID int64, roomTypeCode string) int64 {
	var result int64
	r.db.Model(&room.HtTmRoomTypeDirects{}).
		Where("property_id  = ?", propertyID).
		Where("room_type_code = ?", roomTypeCode).
		Where("is_delete = ?", 0).
		Count(&result)
	return result
}

// CreateRoomDirect 部屋作成
func (r *roomDirectRepository) CreateRoomDirect(roomTable *room.HtTmRoomTypeDirects) error {
	return r.db.Create(&roomTable).Error
}

// UpdateRoomDirect 部屋更新
func (r *roomDirectRepository) UpdateRoomDirect(roomTable *room.HtTmRoomTypeDirects) error {
	return r.db.Model(&room.HtTmRoomTypeDirects{}).
		Where("room_type_id = ?", roomTable.RoomTypeID).
		Updates(map[string]interface{}{
			"name":                        roomTable.Name,
			"room_kind_id":                roomTable.RoomKindID,
			"room_desc":                   roomTable.RoomDesc,
			"stock_setting_start":         roomTable.StockSettingStart,
			"stock_setting_end":           roomTable.StockSettingEnd,
			"is_setting_stock_year_round": roomTable.IsSettingStockYearRound,
			"room_count":                  roomTable.RoomCount,
			"ocu_min":                     roomTable.OcuMin,
			"ocu_max":                     roomTable.OcuMax,
			"is_smoking":                  roomTable.IsSmoking,
			"updated_at":                  time.Now(),
		}).Error
}

// DeleteRoomDirect 部屋を論理削除
func (r *roomDirectRepository) DeleteRoomDirect(roomTypeID int64) error {
	return r.db.Model(&room.HtTmRoomTypeDirects{}).
		Where("room_type_id = ?", roomTypeID).Updates(map[string]interface{}{"is_delete": 1}).Error
}

// ClearRoomToAmenities 部屋に紐づくアメニティを削除
func (r *roomDirectRepository) ClearRoomToAmenities(roomTypeID int64) error {
	return r.db.Delete(&room.HtTmRoomUseAmenityDirects{}, "room_type_id = ?", roomTypeID).Error
}

// CreateRoomToAmenities 部屋に紐づくアメニティを作成
func (r *roomDirectRepository) CreateRoomToAmenities(roomTypeID int64, directRoomAmenityID int64) error {
	return r.db.Create(&room.HtTmRoomUseAmenityDirects{
		RoomTypeID:          roomTypeID,
		DirectRoomAmenityID: directRoomAmenityID,
		Times:               common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
	}).Error
}

// UpdateStopSales 部屋の売止更新
func (r *roomDirectRepository) UpdateStopSales(roomTypeID int64, isStopSales bool) error {
	return r.db.Model(&room.HtTmRoomTypeDirects{}).
		Where("room_type_id = ?", roomTypeID).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}

// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく部屋の売止の更新
func (r *roomDirectRepository) UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, isStopSales bool) error {
	return r.db.Model(&room.HtTmRoomTypeDirects{}).
		Where("room_type_id IN ?", roomTypeIDList).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}
