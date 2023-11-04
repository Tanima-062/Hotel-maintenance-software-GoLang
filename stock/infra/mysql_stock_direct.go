package infra

import (
	"fmt"
	"strings"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	"gorm.io/gorm"
)

const (
	// 一度にまとめて投げるとプレースホルダ長すぎるエラーが起こるのでSQLを分割する。
	// Error 1390: Prepared statement contains too many placeholders
	placeholderLimit int = 1000
)

// stockDirectRepository 直仕入れ在庫関連repository
type stockDirectRepository struct {
	db *gorm.DB
}

// NewStockDirectRepository インスタンス生成
func NewStockDirectRepository(db *gorm.DB) stock.IStockDirectRepository {
	return &stockDirectRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (s *stockDirectRepository) TxStart() (*gorm.DB, error) {
	tx := s.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (s *stockDirectRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (s *stockDirectRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchAllByRoomTypeIDList room_type_idに紐づく在庫を複数件取得
func (s *stockDirectRepository) FetchAllByRoomTypeIDList(roomTypeIDList []int64, startDate string, endDate string) ([]stock.HtTmStockDirects, error) {
	result := []stock.HtTmStockDirects{}
	err := s.db.
		Table("ht_tm_stock_directs").
		Where("room_type_id IN ?", roomTypeIDList).
		Where("use_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error

	return result, err
}

// FetchAllBookingsByPlanIDList plan_idに紐づく販売数を複数件取得
func (s *stockDirectRepository) FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]stock.BookingCount, error) {
	result := []stock.BookingCount{}
	err := s.db.
		Select("count(a.`cm_application_id`) as booking_count, b.plan_id, b.use_date").
		Table("ht_th_booking_prices as b").
		Joins("INNER JOIN ht_th_applications as a ON a.cm_application_id = b.cm_application_id").
		Where("a.cancel_flg = 0").
		Where("a.wholesaler_id = ?", utils.WholesalerIDDirect).
		Where("b.plan_id IN ?", planIDList).
		Where("b.use_date BETWEEN ? AND ?", startDate, endDate).
		Group("b.plan_id, b.use_date").
		Find(&result).Error
	return result, err
}

// FetchStocksByRoomTypeIDList room_type_idに紐づく本日以降の在庫を複数件取得
func (s *stockDirectRepository) FetchStocksByRoomTypeIDList(roomTypeIDList []int64) ([]stock.HtTmStockDirects, error) {
	result := []stock.HtTmStockDirects{}
	err := s.db.
		Table("ht_tm_stock_directs").
		Where("room_type_id IN ? ", roomTypeIDList).
		Where("use_date >= ?", time.Now().Format("2006-01-02")).
		Order("use_date ASC").
		Find(&result).Error
	return result, err
}

// UpdateStopSales room_type_idに紐づく売止の更新
func (s *stockDirectRepository) UpdateStopSales(roomTypeID int64, useDate string, isStopSales bool) error {
	query := s.db.Model(&stock.HtTmStockDirects{}).
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
func (s *stockDirectRepository) UpdateStopSalesByRoomTypeIDList(roomTypeIDList []int64, useDate string, isStopSales bool) error {
	query := s.db.Model(&stock.HtTmStockDirects{}).
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
func (s *stockDirectRepository) UpsertStocks(inputData []stock.HtTmStockDirects) error {

	// Update対象の列を抽出する
	q := s.db.Model(&stock.HtTmStockDirects{})
	for _, v := range inputData {
		q.Or("room_type_id = ? AND use_date = ?", v.RoomTypeID, v.UseDate.Format("2006-01-02"))
	}

	var updateTargetRow []stock.HtTmStockDirects
	q.Find(&updateTargetRow)

	updateValues := [][]interface{}{}
	insertValues := [][]interface{}{}

	// 新規データを更新/新規に分類する
	for _, v := range inputData {
		// Update Query
		var uq []interface{} = nil
		for _, ur := range updateTargetRow {
			if v.RoomTypeID == ur.RoomTypeID && v.UseDate.Format("2006-01-02") == ur.UseDate.Format("2006-01-02") {
				// 提供数-予約数＝在庫数なのでここで計算する
				stock := v.RoomCount - ur.BookingCount
				uq = append(uq, ur.StockID)
				uq = append(uq, ur.RoomTypeID)
				uq = append(uq, ur.UseDate.Format("2006-01-02"))
				uq = append(uq, v.RoomCount)
				uq = append(uq, ur.BookingCount)
				uq = append(uq, stock)
				uq = append(uq, ur.IsStopSales)
				uq = append(uq, time.Now()) // updated_at
				break
			}
		}

		if uq != nil {
			updateValues = append(updateValues, uq)
			continue
		}

		// Insert Query
		var iv []interface{} = nil
		iv = append(iv, v.RoomTypeID)
		iv = append(iv, v.UseDate.Format("2006-01-02"))
		iv = append(iv, v.RoomCount)
		iv = append(iv, v.BookingCount)
		iv = append(iv, v.Stock)
		iv = append(iv, v.IsStopSales)
		iv = append(iv, time.Now()) // updated_at
		iv = append(iv, time.Now()) // created_at

		insertValues = append(insertValues, iv)
	}

	tmpPlaceHolder := []string{}
	tmpValues := []interface{}{}
	for i, v := range updateValues {
		tmpPlaceHolder = append(tmpPlaceHolder, "(?,?,?,?,?,?,?,?)")
		tmpValues = append(tmpValues, v...)
		// クエリの上限もしくは最後に到達したら
		if i%(placeholderLimit+1) == 0 || i == len(updateValues)-1 {
			// Update
			updateQuery := fmt.Sprintf(`INSERT INTO ht_tm_stock_directs (
										stock_id,
										room_type_id,
										use_date,
										room_count,
										booking_count,
										stock,
										is_stop_sales,
										updated_at)
									  VALUES %s ON DUPLICATE KEY UPDATE
										room_count=VALUES(room_count),
										stock=VALUES(stock),
										is_stop_sales=VALUES(is_stop_sales),
										updated_at=VALUES(updated_at)`, strings.Join(tmpPlaceHolder, ","))
			if err := s.db.Exec(updateQuery, tmpValues...).Error; err != nil {
				return err
			}

			tmpPlaceHolder = []string{}
			tmpValues = []interface{}{}
		}
	}

	for i, v := range insertValues {
		tmpPlaceHolder = append(tmpPlaceHolder, "(?,?,?,?,?,?,?,?)")
		tmpValues = append(tmpValues, v...)
		// クエリの上限もしくは最後に到達したら
		if i%(placeholderLimit+1) == 0 || i == len(insertValues)-1 {
			// Insert
			insertQuery := fmt.Sprintf(`INSERT INTO ht_tm_stock_directs (
										room_type_id,
										use_date,
										room_count,
										booking_count,
										stock,
										is_stop_sales,
										created_at,
										updated_at
									  ) VALUES %s `, strings.Join(tmpPlaceHolder, ","))

			if err := s.db.Exec(insertQuery, tmpValues...).Error; err != nil {
				return err
			}

			tmpPlaceHolder = []string{}
			tmpValues = []interface{}{}
		}
	}

	return nil
}

// FetchStock 日付とroom_type_idに紐づく在庫を１件取得
func (s *stockDirectRepository) FetchStock(roomTypeID int64, useDate string) (*stock.HtTmStockDirects, error) {
	result := &stock.HtTmStockDirects{}
	err := s.db.
		Table("ht_tm_stock_directs").
		Where("room_type_id = ?", roomTypeID).
		Where("use_date = ?", useDate).
		First(result).Error
	return result, err
}

// CreateStocks 在庫を複数件作成
func (s *stockDirectRepository) CreateStocks(inputData []stock.HtTmStockDirects) error {
	return s.db.Create(&inputData).Error
}
