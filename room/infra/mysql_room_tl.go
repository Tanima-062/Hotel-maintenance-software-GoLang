package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"gorm.io/gorm"
)

// roomTRepository
type roomTlRepository struct {
	db *gorm.DB
}

// NewRoomTRepository
func NewRoomTlRepository(db *gorm.DB) room.IRoomTlRepository {
	return &roomTlRepository{
		db: db,
	}
}

// TxStart
func (r *roomTlRepository) TxStart() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

// TxCommit transaction commit
func (r *roomTlRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback
func (r *roomTlRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchRoomsByPropertyID
func (r *roomTlRepository) FetchRoomsByPropertyID(req room.ListInput) ([]room.HtTmRoomTypeTls, error) {
	result := []room.HtTmRoomTypeTls{}
	query := r.db.
		Table("ht_tm_room_type_tls AS room").
		Where("property_id = ? AND is_delete = 0", req.PropertyID)
	if req.Paging.Limit > 0 {
		query = query.Limit(req.Paging.Limit).Offset(req.Paging.Offset)
	}
	err := query.Find(&result).Error

	return result, err
}

func (r *roomTlRepository) FetchRoomTypeIdByRoomTypeCode(propertyID int64, roomTypeCode string) (room.HtTmRoomTypeTls, error) {
	result := room.HtTmRoomTypeTls{}
	err := r.db.
		Table("ht_tm_room_type_tls AS a").
		Where("a.property_id = ?", propertyID).
		Where("a.room_type_code = ?", roomTypeCode).
		Where("a.is_delete = ?", 0).
		Find(&result).Error
	return result, err
}

// FetchRoomByRoomTypeID
func (r *roomTlRepository) FetchRoomByRoomTypeID(roomTypeID int64) (*room.HtTmRoomTypeTls, error) {
	result := &room.HtTmRoomTypeTls{}
	err := r.db.
		Table("ht_tm_room_type_tls AS room").
		Where("room_type_id = ? AND is_delete = 0", roomTypeID).
		First(&result).Error
	return result, err
}

// FetchRoomListByRoomTypeID
func (r *roomTlRepository) FetchRoomListByRoomTypeID(roomTypeIDList []int64) ([]room.HtTmRoomTypeTls, error) {
	result := []room.HtTmRoomTypeTls{}
	err := r.db.
		Table("ht_tm_room_type_tls AS room").
		Where("room_type_id IN ? AND is_delete = 0", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// MatchesRoomTypeIDAndPropertyID
func (r *roomTlRepository) MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool {
	var result int64
	r.db.Model(&room.HtTmRoomTypeTls{}).
		Where("room_type_id = ?", roomTypeID).
		Where("property_id = ?", propertyID).
		Count(&result)
	return result > 0
}

// FetchAmenitiesByRoomTypeID
func (r *roomTlRepository) FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomAmenitiesTl, error) {
	result := []room.RoomAmenitiesTl{}
	err := r.db.
		Select("bind.room_type_id, bind.tls_room_amenity_id, amenity.tls_room_amenity_name").
		Table("ht_tm_room_use_amenity_tls AS bind").
		Joins("INNER JOIN ht_tm_room_amenity_tls AS amenity ON bind.tls_room_amenity_id = amenity.tls_room_amenity_id").
		Where("lang_cd = ?", "ja-JP").
		Where("bind.room_type_id IN ?", roomTypeIDList).
		Find(&result).Error
	return result, err
}

// FetchAllAmenities
func (r *roomTlRepository) FetchAllAmenities() ([]room.HtTmRoomAmenityTls, error) {
	result := []room.HtTmRoomAmenityTls{}
	err := r.db.Where("lang_cd = ?", "ja-JP").Find(&result).Error
	return result, err
}

// CountRoomTypeCode
func (r *roomTlRepository) CountRoomTypeCode(propertyID int64, roomTypeCode string) int64 {
	var result int64
	r.db.Model(&room.HtTmRoomTypeTls{}).
		Where("property_id  = ?", propertyID).
		Where("room_type_code = ?", roomTypeCode).
		Where("is_delete = ?", 0).
		Count(&result)
	return result
}

// CreateRoomTl
func (r *roomTlRepository) CreateRoomTl(roomTable *room.HtTmRoomTypeTls) error {
	return r.db.Create(&roomTable).Error
}

// UpdateRoomTl
func (r *roomTlRepository) UpdateRoomTl(roomTable *room.HtTmRoomTypeTls) error {
	return r.db.Model(&room.HtTmRoomTypeTls{}).
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

// CreateRoomTl
func (r *roomTlRepository) CreateRoomBulkTl(roomTable *room.HtTmRoomTypeTls) error {
	return r.db.Create(&roomTable).Error
}

// UpdateRoomTl
func (r *roomTlRepository) UpdateRoomBulkTl(roomTable *room.HtTmRoomTypeTls) error {
	return r.db.Model(&room.HtTmRoomTypeTls{}).
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
			"is_stop_sales":               roomTable.IsStopSales,
			"updated_at":                  time.Now(),
		}).Error
}

// DeleteRoomTl
func (r *roomTlRepository) DeleteRoomTl(roomTypeID int64) error {
	return r.db.Model(&room.HtTmRoomTypeTls{}).
		Where("room_type_id = ?", roomTypeID).Updates(map[string]interface{}{"is_delete": 1}).Error
}

// ClearRoomToAmenities
func (r *roomTlRepository) ClearRoomToAmenities(roomTypeID int64) error {
	return r.db.Delete(&room.HtTmRoomUseAmenityTls{}, "room_type_id = ?", roomTypeID).Error
}

// CreateRoomToAmenities
func (r *roomTlRepository) CreateRoomToAmenities(roomTypeID int64, TlsRoomAmenityID int64) error {
	return r.db.Create(&room.HtTmRoomUseAmenityTls{
		RoomTypeID:       roomTypeID,
		TlsRoomAmenityID: TlsRoomAmenityID,
		Times:            common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
	}).Error
}

// UpdateStopSales
func (r *roomTlRepository) UpdateStopSales(roomTypeID int64, isStopSales bool) error {
	return r.db.Model(&room.HtTmRoomTypeTls{}).
		Where("room_type_id = ?", roomTypeID).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}

// UpdateStopSalesByRoomTypeIDList
func (r *roomTlRepository) UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, isStopSales bool) error {
	return r.db.Model(&room.HtTmRoomTypeTls{}).
		Where("room_type_id IN ?", roomTypeIDList).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}

func (r *roomTlRepository) FetchAllRoomKindTls() ([]room.HtTmRoomKindsTls, error) {
	result := []room.HtTmRoomKindsTls{}
	err := r.db.
		Model(&room.HtTmRoomKindsTls{}).
		Where("is_delete = 0").
		Find(&result).Error
	return result, err
}
