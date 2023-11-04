package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"gorm.io/gorm"
)

// roomRaku2Repository らく通部屋関連repository
type roomRaku2Repository struct {
	db *gorm.DB
}

// NewRoomRaku2Repository インスタンス生成
func NewRoomRaku2Repository(db *gorm.DB) room.IRoomRaku2Repository {
	return &roomRaku2Repository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (r *roomRaku2Repository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (r *roomRaku2Repository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (r *roomRaku2Repository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchRoomsByPropertyID propertyIDに紐づく部屋複数件取得
func (r *roomRaku2Repository) FetchRoomsByPropertyID(req room.ListInput) ([]room.HtTmRoomTypeRaku2s, error) {
	result := []room.HtTmRoomTypeRaku2s{}
	query := r.db.
		Table("ht_tm_room_type_raku2s AS room").
		Where("property_id = ? AND is_delete = 0", req.PropertyID)
	if req.Paging.Limit > 0 {
		query = query.Limit(req.Paging.Limit).Offset(req.Paging.Offset)
	}
	err := query.Find(&result).Error

	return result, err
}

// FetchRoomByRoomTypeID roomTypeIDに紐づく部屋を1件取得
func (r *roomRaku2Repository) FetchRoomByRoomTypeID(roomTypeID int64) (*room.HtTmRoomTypeRaku2s, error) {
	result := &room.HtTmRoomTypeRaku2s{}
	err := r.db.
		Table("ht_tm_room_type_raku2s AS room").
		Where("room_type_id = ? AND is_delete = 0", roomTypeID).
		First(&result).Error
	return result, err
}

// FetchRoomListByRoomTypeID roomTypeIDに紐づく部屋を複数件取得
func (r *roomRaku2Repository) FetchRoomListByRoomTypeID(roomTypeIDList []int64) ([]room.HtTmRoomTypeRaku2s, error) {
	result := []room.HtTmRoomTypeRaku2s{}
	err := r.db.
		Table("ht_tm_room_type_raku2s AS room").
		Where("room_type_id IN ? AND is_delete = 0", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// MatchesRoomTypeIDAndPropertyID propertyIDとroomTypeIDが紐付いているか
func (r *roomRaku2Repository) MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool {
	var result int64
	r.db.Model(&room.HtTmRoomTypeRaku2s{}).
		Where("room_type_id = ?", roomTypeID).
		Where("property_id = ?", propertyID).
		Count(&result)
	return result > 0
}

// FetchAmenitiesByRoomTypeID 部屋に紐づくアメニティを複数件取得
func (r *roomRaku2Repository) FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomAmenitiesRaku2, error) {
	result := []room.RoomAmenitiesRaku2{}
	err := r.db.
		Select("bind.room_type_id, bind.raku2_room_amenity_id, amenity.raku2_room_amenity_name").
		Table("ht_tm_room_use_amenity_raku2s AS bind").
		Joins("INNER JOIN ht_tm_room_amenity_raku2s AS amenity ON bind.raku2_room_amenity_id = amenity.raku2_room_amenity_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// FetchAllAmenities 部屋のアメニティを複数件取得
func (r *roomRaku2Repository) FetchAllAmenities() ([]room.HtTmRoomAmenityRaku2s, error) {
	result := []room.HtTmRoomAmenityRaku2s{}
	err := r.db.Find(&result).Error
	return result, err
}

// CountRoomTypeCode 部屋コードの重複件数
func (r *roomRaku2Repository) CountRoomTypeCode(propertyID int64, roomTypeCode string) int64 {
	var result int64
	r.db.Model(&room.HtTmRoomTypeRaku2s{}).
		Where("property_id  = ?", propertyID).
		Where("room_type_code = ?", roomTypeCode).
		Where("is_delete = ?", 0).
		Count(&result)
	return result
}

// CreateRoomRaku2 部屋作成
func (r *roomRaku2Repository) CreateRoomRaku2(roomTable *room.HtTmRoomTypeRaku2s) error {
	return r.db.Create(&roomTable).Error
}

// UpdateRoomRaku2 部屋更新
func (r *roomRaku2Repository) UpdateRoomRaku2(roomTable *room.HtTmRoomTypeRaku2s) error {
	return r.db.Model(&room.HtTmRoomTypeRaku2s{}).
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

// DeleteRoomRaku2 部屋を論理削除
func (r *roomRaku2Repository) DeleteRoomRaku2(roomTypeID int64) error {
	return r.db.Model(&room.HtTmRoomTypeRaku2s{}).
		Where("room_type_id = ?", roomTypeID).Updates(map[string]interface{}{"is_delete": 1}).Error
}

// ClearRoomToAmenities 部屋に紐づくアメニティを削除
func (r *roomRaku2Repository) ClearRoomToAmenities(roomTypeID int64) error {
	return r.db.Delete(&room.HtTmRoomUseAmenityRaku2s{}, "room_type_id = ?", roomTypeID).Error
}

// CreateRoomToAmenities 部屋に紐づくアメニティを作成
func (r *roomRaku2Repository) CreateRoomToAmenities(roomTypeID int64, raku2RoomAmenityID int64) error {
	return r.db.Create(&room.HtTmRoomUseAmenityRaku2s{
		RoomTypeID:         roomTypeID,
		Raku2RoomAmenityID: raku2RoomAmenityID,
		Times:              common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
	}).Error
}

// UpdateStopSales 部屋の売止更新
func (r *roomRaku2Repository) UpdateStopSales(roomTypeID int64, isStopSales bool) error {
	return r.db.Model(&room.HtTmRoomTypeRaku2s{}).
		Where("room_type_id = ?", roomTypeID).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}

// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく部屋の売止の更新
func (r *roomRaku2Repository) UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, isStopSales bool) error {
	return r.db.Model(&room.HtTmRoomTypeRaku2s{}).
		Where("room_type_id IN ?", roomTypeIDList).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}
