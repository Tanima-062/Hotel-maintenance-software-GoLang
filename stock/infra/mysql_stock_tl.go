package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	"gorm.io/gorm"
)

// stockTLRepository Tl inventory related repository
type stockTlRepository struct {
	db *gorm.DB
}

// NewStockTLRepository instantiation
func NewStockTlRepository(db *gorm.DB) stock.IStockTlRepository {
	return &stockTlRepository{
		db: db,
	}
}

// TxStart transaction start
func (s *stockTlRepository) TxStart() (*gorm.DB, error) {
	tx := s.db.Begin()
	return tx, tx.Error
}

// TxCommit transaction commit
func (s *stockTlRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback transaction rollback
func (s *stockTlRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchAllByRoomTypeIDList Acquire multiple items of inventory linked to room_type_id
func (s *stockTlRepository) FetchAllByRoomTypeIDList(roomTypeIDList []int64, startDate string, endDate string) ([]stock.HtTmStockTls, error) {
	result := []stock.HtTmStockTls{}
	err := s.db.
		Table("ht_tm_stock_tls").
		Where("room_type_id IN ?", roomTypeIDList).
		Where("use_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return result, err
}

// FetchAllBookingsByPlanIDList Get multiple sales numbers linked to plan_id
func (s *stockTlRepository) FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]stock.BookingCount, error) {
	result := []stock.BookingCount{}
	err := s.db.
		Select("count(a.`cm_application_id`) as booking_count, b.plan_id, b.use_date").
		Table("ht_th_booking_prices as b").
		Joins("INNER JOIN ht_th_applications as a ON a.cm_application_id = b.cm_application_id").
		Where("a.cancel_flg = 0").
		Where("a.wholesaler_id = ?", utils.WholesalerIDTl).
		Where("b.plan_id IN ?", planIDList).
		Where("b.use_date BETWEEN ? AND ?", startDate, endDate).
		Group("b.plan_id, b.use_date").
		Find(&result).Error
	return result, err
}

// UpdateStopSales Updating the sale stop linked to room_type_id
func (s *stockTlRepository) UpdateStopSales(roomTypeID int64, useDate string, isStopSales bool) error {
	query := s.db.Model(&stock.HtTmStockTls{}).
		Where("room_type_id = ?", roomTypeID)
	if useDate != "" {
		query = query.Where("use_date = ?", useDate)
	}
	return query.Updates(map[string]interface{}{
		"is_stop_sales": isStopSales,
		"updated_at":    time.Now(),
	}).Error
}

// UpdateStopSalesByRoomTypeIDList Updating sales stop linked to room_type_id (multiple)
func (s *stockTlRepository) UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, useDate string, isStopSales bool) error {
	query := s.db.Model(&stock.HtTmStockTls{}).
		Where("room_type_id IN ?", roomTypeIDList)
	if useDate != "" {
		query = query.Where("use_date = ?", useDate)
	}
	return query.Updates(map[string]interface{}{
		"is_stop_sales": isStopSales,
		"updated_at":    time.Now(),
	}).Error
}

// UpsertStocks Create/update inventory
func (s *stockTlRepository) UpsertStocks(inputData []stock.HtTmStockTls) error {

	for _, v := range inputData {
		assignData := map[string]interface{}{
			"room_type_id":  v.RoomTypeID,
			"room_count":    v.RoomCount,
			"booking_count": v.BookingCount,
			"use_date":      v.UseDate,
			"updated_at":    v.UpdatedAt,
		}

		if err := s.db.Model(&stock.HtTmStockTls{}).
			Where("room_type_id = ?", v.RoomTypeID).
			Where("use_date = ?", v.UseDate.Format("2006-01-02")).
			Assign(assignData).
			FirstOrCreate(&stock.HtTmStockTls{}).
			Error; err != nil {
			return err
		}
	}
	return nil
}

// CreateStocks creates a new stock
func (s *stockTlRepository) CreateStocks(inputData []stock.HtTmStockTls) error {
	for _, data := range inputData {
		if err := s.db.Create(&data).Error; err != nil {
			return err
		}
	}
	return nil
}

func (b *stockTlRepository) UpdateStocksBulk(roomTypeID int64, useDate string, stockCount int64, bookingCount int64, isStopSale bool) error {
	return b.db.Model(&stock.HtTmStockTls{}).
		Where("room_type_id = ?", roomTypeID).
		Where("use_date = ?", useDate).
		Updates(map[string]interface{}{
			"room_count":    stockCount + bookingCount,
			"stock":         stockCount,
			"is_stop_sales": isStopSale,
			"updated_at":    time.Now(),
		}).Error
}

func (b *stockTlRepository) FetchBookingCountByRoomTypeId(roomTypeID int64, useDate string) (stock.StockTable, error) {
	result := stock.StockTable{}
	err := b.db.
		Select("a.room_type_id, a.use_date, a.booking_count, a.stock, a.room_count").
		Table("ht_tm_stock_tls AS a").
		Where("a.room_type_id = ?", roomTypeID).
		Where("a.use_date = ?", useDate).
		Find(&result).Error
	return result, err
}
