package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"gorm.io/gorm"
)

// roomNeppanRepository ねっぱん部屋関連repository
type roomNeppanRepository struct {
	db *gorm.DB
}

// NewRoomNeppanRepository インスタンス生成
func NewRoomNeppanRepository(db *gorm.DB) room.IRoomNeppanRepository {
	return &roomNeppanRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (r *roomNeppanRepository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (r *roomNeppanRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (r *roomNeppanRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchRoomsByPropertyID propertyIDに紐づく部屋複数件取得
func (r *roomNeppanRepository) FetchRoomsByPropertyID(req room.ListInput) ([]room.HtTmRoomTypeNeppans, error) {
	result := []room.HtTmRoomTypeNeppans{}
	query := r.db.
		Table("ht_tm_room_type_neppans AS room").
		Where("property_id = ? AND is_delete = 0", req.PropertyID)
	if req.Paging.Limit > 0 {
		query = query.Limit(req.Paging.Limit).Offset(req.Paging.Offset)
	}
	err := query.Find(&result).Error

	return result, err
}

// FetchRoomByRoomTypeID roomTypeIDに紐づく部屋を1件取得
func (r *roomNeppanRepository) FetchRoomByRoomTypeID(roomTypeID int64) (*room.HtTmRoomTypeNeppans, error) {
	result := &room.HtTmRoomTypeNeppans{}
	err := r.db.
		Table("ht_tm_room_type_neppans AS room").
		Where("room_type_id = ? AND is_delete = 0", roomTypeID).
		First(&result).Error
	return result, err
}

// FetchRoomListByRoomTypeID roomTypeIDに紐づく部屋を複数件取得
func (r *roomNeppanRepository) FetchRoomListByRoomTypeID(roomTypeIDList []int64) ([]room.HtTmRoomTypeNeppans, error) {
	result := []room.HtTmRoomTypeNeppans{}
	err := r.db.
		Table("ht_tm_room_type_neppans AS room").
		Where("room_type_id IN ? AND is_delete = 0", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// MatchesRoomTypeIDAndPropertyID propertyIDとroomTypeIDが紐付いているか
func (r *roomNeppanRepository) MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool {
	var result int64
	r.db.Model(&room.HtTmRoomTypeNeppans{}).
		Where("room_type_id = ?", roomTypeID).
		Where("property_id = ?", propertyID).
		Count(&result)
	return result > 0
}

// FetchAmenitiesByRoomTypeID 部屋に紐づくアメニティを複数件取得
func (r *roomNeppanRepository) FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomAmenitiesNeppan, error) {
	result := []room.RoomAmenitiesNeppan{}
	err := r.db.
		Select("bind.room_type_id, bind.neppan_room_amenity_id, amenity.neppan_room_amenity_name").
		Table("ht_tm_room_use_amenity_neppans AS bind").
		Joins("INNER JOIN ht_tm_room_amenity_neppans AS amenity ON bind.neppan_room_amenity_id = amenity.neppan_room_amenity_id").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// FetchAllAmenities 部屋のアメニティを複数件取得
func (r *roomNeppanRepository) FetchAllAmenities() ([]room.HtTmRoomAmenityNeppans, error) {
	result := []room.HtTmRoomAmenityNeppans{}
	err := r.db.Find(&result).Error
	return result, err
}

// CountRoomTypeCode 部屋コードの重複件数
func (r *roomNeppanRepository) CountRoomTypeCode(propertyID int64, roomTypeCode string) int64 {
	var result int64
	r.db.Model(&room.HtTmRoomTypeNeppans{}).
		Where("property_id  = ?", propertyID).
		Where("room_type_code = ?", roomTypeCode).
		Where("is_delete = ?", 0).
		Count(&result)
	return result
}

// CreateRoomNeppan 部屋作成
func (r *roomNeppanRepository) CreateRoomNeppan(roomTable *room.HtTmRoomTypeNeppans) error {
	return r.db.Create(&roomTable).Error
}

// UpdateRoomNeppan 部屋更新
func (r *roomNeppanRepository) UpdateRoomNeppan(roomTable *room.HtTmRoomTypeNeppans) error {
	return r.db.Model(&room.HtTmRoomTypeNeppans{}).
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

// DeleteRoomNeppan 部屋を論理削除
func (r *roomNeppanRepository) DeleteRoomNeppan(roomTypeID int64) error {
	return r.db.Model(&room.HtTmRoomTypeNeppans{}).
		Where("room_type_id = ?", roomTypeID).Updates(map[string]interface{}{"is_delete": 1}).Error
}

// ClearRoomToAmenities 部屋に紐づくアメニティを削除
func (r *roomNeppanRepository) ClearRoomToAmenities(roomTypeID int64) error {
	return r.db.Delete(&room.HtTmRoomUseAmenityNeppans{}, "room_type_id = ?", roomTypeID).Error
}

// CreateRoomToAmenities 部屋に紐づくアメニティを作成
func (r *roomNeppanRepository) CreateRoomToAmenities(roomTypeID int64, neppanRoomAmenityID int64) error {
	return r.db.Create(&room.HtTmRoomUseAmenityNeppans{
		RoomTypeID:          roomTypeID,
		NeppanRoomAmenityID: neppanRoomAmenityID,
		Times:               common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
	}).Error
}

// UpdateStopSales 部屋の売止更新
func (r *roomNeppanRepository) UpdateStopSales(roomTypeID int64, isStopSales bool) error {
	return r.db.Model(&room.HtTmRoomTypeNeppans{}).
		Where("room_type_id = ?", roomTypeID).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}

// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく部屋の売止の更新
func (r *roomNeppanRepository) UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, isStopSales bool) error {
	return r.db.Model(&room.HtTmRoomTypeNeppans{}).
		Where("room_type_id IN ?", roomTypeIDList).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}
