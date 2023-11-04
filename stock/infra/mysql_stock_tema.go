package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	"gorm.io/gorm"
	"time"
)

// stockTemaRepository tema inventory related repository
type stockTemaRepository struct {
	db *gorm.DB
}

// NewStockTemaRepository creates a new stock tema inventory repository
func NewStockTemaRepository(db *gorm.DB) stock.IStockTemaRepository {
	return &stockTemaRepository{
		db: db,
	}
}

// TxStart Tx inventory related repository
func (s *stockTemaRepository) TxStart() (*gorm.DB, error) {
	tx := s.db.Begin()
	return tx, tx.Error
}

// TxCommit Commit
func (s *stockTemaRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback Tx inventory related repository
func (s *stockTemaRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchRoomTypeIdByRoomTypeCode fetch room type by code
func (s *stockTemaRepository) FetchRoomTypeIdByRoomTypeCode(propertyID int64, roomTypeCode string) (room.HtTmRoomTypeTemas, error) {
	result := room.HtTmRoomTypeTemas{}
	err := s.db.
		Table("ht_tm_room_type_temas AS a").
		Where("a.property_id = ?", propertyID).
		Where("a.room_type_code = ?", roomTypeCode).
		Where("a.is_delete = ?", 0).
		Find(&result).Error
	return result, err
}

// UpdateRoomBulkTema update rooom bulk data
func (s *stockTemaRepository) UpdateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error {
	return s.db.Model(&room.HtTmRoomTypeTemas{}).
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

// FetchBookingCountByRoomTypeId returns the number of books in the room type table
func (s *stockTemaRepository) FetchBookingCountByRoomTypeId(roomTypeCode string, ariDate string) (stock.StockTableTema, error) {
	result := stock.StockTableTema{}
	err := s.db.
		Select("a.ari_tema_id, a.ari_date, a.stock, a.disable, a.room_type_code").
		Table("ht_tm_stock_temas AS a").
		Where("a.room_type_code = ?", roomTypeCode).
		Where("a.ari_date = ?", ariDate).
		Find(&result).Error
	return result, err
}

// UpdateStocksBulk  updates the stock table with the specified room type and quantity in the database
func (s *stockTemaRepository) UpdateStocksBulk(roomTypeCode string, ariDate string, stockCount int64, disable bool) error {
	return s.db.Model(&stock.HtTmStockTemas{}).
		Where("room_type_code = ?", roomTypeCode).
		Where("ari_date = ?", ariDate).
		Updates(map[string]interface{}{
			"stock":      stockCount,
			"disable":    disable,
			"updated_at": time.Now(),
		}).Error
}

// CreateStocks updates the stock table with the specified room type and quantity in the database
func (s *stockTemaRepository) CreateStocks(inputData []stock.HtTmStockTemas) error {
	for _, data := range inputData {
		if err := s.db.Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}

// FetchAllStocksByRoomTypeIDList FetchAllByRoomTypeIDList Acquire multiple items of inventory linked to room_type_id
func (s *stockTemaRepository) FetchAllStocksByRoomTypeIDList(roomTypeIDList []int64, startDate string, endDate string) ([]stock.StockTableTema, error) {
	result := []stock.StockTableTema{}
	err := s.db.
		Select("a.*, b.room_type_id").
		Table("ht_tm_stock_temas as a").
		Joins("INNER JOIN ht_tm_room_type_temas as b ON a.room_type_code = b.room_type_code").
		Where("b.room_type_id IN ?", roomTypeIDList).
		Where("a.ari_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return result, err
}

// FetchAllBookingsByPlanIDList Get multiple sales numbers linked to plan_id
func (s *stockTemaRepository) FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]stock.BookingCount, error) {
	result := []stock.BookingCount{}
	err := s.db.
		Select("count(a.`cm_application_id`) as booking_count, b.plan_id, b.use_date").
		Table("ht_th_booking_prices as b").
		Joins("INNER JOIN ht_th_applications as a ON a.cm_application_id = b.cm_application_id").
		Where("a.cancel_flg = 0").
		Where("a.wholesaler_id = ?", utils.WholesalerIDTema).
		Where("b.plan_id IN ?", planIDList).
		Where("b.use_date BETWEEN ? AND ?", startDate, endDate).
		Group("b.plan_id, b.use_date").
		Find(&result).Error
	return result, err
}

// FetchAllByRoomTypeCodeList get multiple room code
func (s *stockTemaRepository) FetchAllByRoomTypeCodeList(roomTypeCodeList []int64, startDate string, endDate string) ([]stock.HtTmStockTemas, error) {
	result := []stock.HtTmStockTemas{}
	err := s.db.
		Table("ht_tm_stock_temas").
		Where("room_type_code IN ?", roomTypeCodeList).
		Where("ari_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return result, err
}
