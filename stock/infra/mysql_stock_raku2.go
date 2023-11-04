package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	"gorm.io/gorm"
)

// stockRaku2Repository らく通在庫関連repository
type stockRaku2Repository struct {
	db *gorm.DB
}

// NewStockRaku2Repository インスタンス生成
func NewStockRaku2Repository(db *gorm.DB) stock.IStockRaku2Repository {
	return &stockRaku2Repository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (s *stockRaku2Repository) TxStart() (*gorm.DB, error) {
	tx := s.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (s *stockRaku2Repository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (s *stockRaku2Repository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchAllByRoomTypeIDList room_type_idに紐づく在庫を複数件取得
func (s *stockRaku2Repository) FetchAllByRoomTypeIDList(roomTypeIDList []int64, startDate string, endDate string) ([]stock.HtTmStockRaku2s, error) {
	result := []stock.HtTmStockRaku2s{}
	err := s.db.
		Table("ht_tm_stock_raku2s").
		Where("room_type_id IN ?", roomTypeIDList).
		Where("use_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return result, err
}

// FetchAllBookingsByPlanIDList plan_idに紐づく販売数を複数件取得
func (s *stockRaku2Repository) FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]stock.BookingCount, error) {
	result := []stock.BookingCount{}
	err := s.db.
		Select("count(a.`cm_application_id`) as booking_count, b.plan_id, b.use_date").
		Table("ht_th_booking_prices as b").
		Joins("INNER JOIN ht_th_applications as a ON a.cm_application_id = b.cm_application_id").
		Where("a.cancel_flg = 0").
		Where("a.wholesaler_id = ?", utils.WholesalerIDRaku2).
		Where("b.plan_id IN ?", planIDList).
		Where("b.use_date BETWEEN ? AND ?", startDate, endDate).
		Group("b.plan_id, b.use_date").
		Find(&result).Error
	return result, err
}

// UpdateStopSales room_type_idに紐づく売止の更新
func (s *stockRaku2Repository) UpdateStopSales(roomTypeID int64, useDate string, isStopSales bool) error {
	query := s.db.Model(&stock.HtTmStockRaku2s{}).
		Where("room_type_id = ?", roomTypeID)
	if useDate != "" {
		query = query.Where("use_date = ?", useDate)
	}
	return query.Updates(map[string]interface{}{
		"is_stop_sales": isStopSales,
		"updated_at":    time.Now(),
	}).Error
}

// UpdateStopSalesByRoomTypeIDList room_type_id(複数)に紐づく売止の更新
func (s *stockRaku2Repository) UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, useDate string, isStopSales bool) error {
	query := s.db.Model(&stock.HtTmStockRaku2s{}).
		Where("room_type_id IN ?", roomTypeIDList)
	if useDate != "" {
		query = query.Where("use_date = ?", useDate)
	}
	return query.Updates(map[string]interface{}{
		"is_stop_sales": isStopSales,
		"updated_at":    time.Now(),
	}).Error
}

// UpsertStocks 在庫の作成・更新
func (s *stockRaku2Repository) UpsertStocks(inputData []stock.HtTmStockRaku2s) error {

	for _, v := range inputData {
		assignData := map[string]interface{}{
			"room_type_id":  v.RoomTypeID,
			"room_count":    v.RoomCount,
			"booking_count": v.BookingCount,
			"use_date":      v.UseDate,
			"updated_at":    v.UpdatedAt,
		}

		if err := s.db.Model(&stock.HtTmStockRaku2s{}).
			Where("room_type_id = ?", v.RoomTypeID).
			Where("use_date = ?", v.UseDate.Format("2006-01-02")).
			Assign(assignData).
			FirstOrCreate(&stock.HtTmStockRaku2s{}).
			Error; err != nil {
			return err
		}
	}
	return nil
}
