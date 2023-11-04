package infra

import (
	"fmt"
	"strings"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/price"
	"gorm.io/gorm"
)

const (
	// 一度にまとめて投げるとプレースホルダ長すぎるエラーが起こるのでSQLを分割する。
	// Error 1390: Prepared statement contains too many placeholders
	placeholderLimit int = 1000
)

// priceDirectRepository 直仕入れ料金関連repository
type priceDirectRepository struct {
	db *gorm.DB
}

// NewPriceDirectRepository インスタンス生成
func NewPriceDirectRepository(db *gorm.DB) price.IPriceDirectRepository {
	return &priceDirectRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (p *priceDirectRepository) TxStart() (*gorm.DB, error) {
	tx := p.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (p *priceDirectRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (p *priceDirectRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchPricesWithInThePeriod 期間内のプランの料金を複数件取得
func (p *priceDirectRepository) FetchPricesWithInThePeriod(planID int64, startDate string, endDate string) ([]price.HtTmPriceDirects, error) {
	result := []price.HtTmPriceDirects{}
	err := p.db.
		Table("ht_tm_price_directs").
		Where("plan_id = ? AND ? >= use_date AND use_date >= ?", planID, endDate, startDate).
		Order("use_date ASC").
		Find(&result).Error
	return result, err
}

// UpsertPrices 料金の作成と更新
func (p *priceDirectRepository) UpsertPrices(inputData []price.HtTmPriceDirects) error {
	// Update対象の列を抽出する
	q := p.db.Model(&price.HtTmPriceDirects{})
	for _, v := range inputData {
		q.Or("plan_id = ? AND rate_type_code >= ? AND use_date = ?", v.PlanID, v.RateTypeCode, v.UseDate.Format("2006-01-02"))
	}
	var updateTargetRow []price.HtTmPriceDirects
	q.Find(&updateTargetRow)

	updateValues := [][]interface{}{}
	insertValues := [][]interface{}{}

	// 新規データを更新/新規に分類する
	for _, v := range inputData {
		// Update Query
		var uq []interface{} = nil
		for _, ur := range updateTargetRow {
			if v.PlanID == ur.PlanID && v.RateTypeCode == ur.RateTypeCode && v.UseDate.Format("2006-01-02") == ur.UseDate.Format("2006-01-02") {
				uq = append(uq, ur.PriceID)
				uq = append(uq, ur.PlanID)
				uq = append(uq, ur.UseDate.Format("2006-01-02"))
				uq = append(uq, ur.RateTypeCode)
				uq = append(uq, v.Price)
				uq = append(uq, v.PriceInTax)
				uq = append(uq, v.ChildPrice1)
				uq = append(uq, v.ChildPrice1InTax)
				uq = append(uq, v.ChildPrice2)
				uq = append(uq, v.ChildPrice2InTax)
				uq = append(uq, v.ChildPrice3)
				uq = append(uq, v.ChildPrice3InTax)
				uq = append(uq, v.ChildPrice4)
				uq = append(uq, v.ChildPrice4InTax)
				uq = append(uq, v.ChildPrice5)
				uq = append(uq, v.ChildPrice5InTax)
				uq = append(uq, v.ChildPrice6)
				uq = append(uq, v.ChildPrice6InTax)
				uq = append(uq, v.RegularPrice)
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
		iv = append(iv, v.PlanID)
		iv = append(iv, v.UseDate.Format("2006-01-02"))
		iv = append(iv, v.RateTypeCode)
		iv = append(iv, v.Price)
		iv = append(iv, v.PriceInTax)
		iv = append(iv, v.ChildPrice1)
		iv = append(iv, v.ChildPrice1InTax)
		iv = append(iv, v.ChildPrice2)
		iv = append(iv, v.ChildPrice2InTax)
		iv = append(iv, v.ChildPrice3)
		iv = append(iv, v.ChildPrice3InTax)
		iv = append(iv, v.ChildPrice4)
		iv = append(iv, v.ChildPrice4InTax)
		iv = append(iv, v.ChildPrice5)
		iv = append(iv, v.ChildPrice5InTax)
		iv = append(iv, v.ChildPrice6)
		iv = append(iv, v.ChildPrice6InTax)
		iv = append(iv, v.RegularPrice)
		iv = append(iv, time.Now()) // updated_at
		iv = append(iv, time.Now()) // created_at

		insertValues = append(insertValues, iv)
	}

	tmpPlaceHolder := []string{}
	tmpValues := []interface{}{}
	for i, v := range updateValues {
		tmpPlaceHolder = append(tmpPlaceHolder, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		tmpValues = append(tmpValues, v...)
		// クエリの上限もしくは最後に到達したら
		if i%(placeholderLimit+1) == 0 || i == len(updateValues)-1 {
			// Update
			updateQuery := fmt.Sprintf(`INSERT INTO ht_tm_price_directs (
                                    price_id, 
                                    plan_id, 
                                    use_date, 
                                    rate_type_code, 
                                    price, price_in_tax, 
                                    child_price1, 
                                    child_price1_in_tax, 
                                    child_price2, 
                                    child_price2_in_tax, 
                                    child_price3, 
                                    child_price3_in_tax, 
                                    child_price4, 
                                    child_price4_in_tax, 
                                    child_price5, 
                                    child_price5_in_tax, 
                                    child_price6, 
                                    child_price6_in_tax, 
                                    regular_price,  
                                    updated_at) 
                                  VALUES %s ON DUPLICATE KEY UPDATE 
                                    price=VALUES(price), 
                                    price_in_tax=VALUES(price_in_tax), 
                                    child_price1=VALUES(child_price1), 
                                    child_price1_in_tax=VALUES(child_price1_in_tax), 
                                    child_price2=VALUES(child_price2), 
                                    child_price2_in_tax=VALUES(child_price2_in_tax), 
                                    child_price3=VALUES(child_price3), 
                                    child_price3_in_tax=VALUES(child_price3_in_tax), 
                                    child_price4=VALUES(child_price4), 
                                    child_price4_in_tax=VALUES(child_price4_in_tax), 
                                    child_price5=VALUES(child_price5), 
                                    child_price5_in_tax=VALUES(child_price5_in_tax), 
                                    child_price6=VALUES(child_price6), 
                                    child_price6_in_tax=VALUES(child_price6_in_tax), 
                                    regular_price=VALUES(regular_price),
                                    updated_at=VALUES(updated_at)`, strings.Join(tmpPlaceHolder, ","))
			if err := p.db.Exec(updateQuery, tmpValues...).Error; err != nil {
				return err
			}

			tmpPlaceHolder = []string{}
			tmpValues = []interface{}{}
		}
	}

	for i, v := range insertValues {
		tmpPlaceHolder = append(tmpPlaceHolder, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		tmpValues = append(tmpValues, v...)
		// クエリの上限もしくは最後に到達したら
		if i%(placeholderLimit+1) == 0 || i == len(insertValues)-1 {
			// Insert
			insertQuery := fmt.Sprintf(`INSERT INTO ht_tm_price_directs (
                                    plan_id, 
                                    use_date, 
                                    rate_type_code, 
                                    price, 
                                    price_in_tax, 
                                    child_price1, 
                                    child_price1_in_tax, 
                                    child_price2, 
                                    child_price2_in_tax, 
                                    child_price3, 
                                    child_price3_in_tax, 
                                    child_price4, 
                                    child_price4_in_tax, 
                                    child_price5, 
                                    child_price5_in_tax, 
                                    child_price6, 
                                    child_price6_in_tax, 
                                    regular_price, 
                                    created_at, 
                                    updated_at
                                  ) VALUES %s `, strings.Join(tmpPlaceHolder, ","))

			if err := p.db.Exec(insertQuery, tmpValues...).Error; err != nil {
				return err
			}

			tmpPlaceHolder = []string{}
			tmpValues = []interface{}{}
		}
	}

	return nil
}

// FetchChildRates プランに紐づく料金設定を複数件取得
func (p *priceDirectRepository) FetchChildRates(planID int64) ([]price.HtTmChildRateDirects, error) {
	result := []price.HtTmChildRateDirects{}
	err := p.db.
		Model(&price.HtTmChildRateDirects{}).
		Where("plan_id = ?", planID).
		Find(&result).Error
	return result, err
}

// FetchChildRatesByPlanIDList 複数プランに紐づく料金設定を複数件取得
func (p *priceDirectRepository) FetchChildRatesByPlanIDList(planIDList []int64) ([]price.HtTmChildRateDirects, error) {
	result := []price.HtTmChildRateDirects{}
	err := p.db.
		Model(&price.HtTmChildRateDirects{}).
		Where("plan_id IN ?", planIDList).
		Find(&result).Error
	return result, err
}

// FetchAllByPlanIDList 期間内の複数のプランIDに紐づく料金を複数件取得
func (p *priceDirectRepository) FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]price.HtTmPriceDirects, error) {
	result := []price.HtTmPriceDirects{}
	err := p.db.
		Table("ht_tm_price_directs").
		Where("plan_id IN ?", planIDList).
		Where("use_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return result, err
}

// FetchPricesByPlanID 本日以降の料金を複数件取得
func (p *priceDirectRepository) FetchPricesByPlanID(planID int64) ([]price.HtTmPriceDirects, error) {
	result := []price.HtTmPriceDirects{}
	err := p.db.
		Table("ht_tm_price_directs").
		Where("plan_id = ? ", planID).
		Where("use_date >= ?", time.Now().Format("2006-01-02")).
		Order("use_date ASC").
		Find(&result).Error
	return result, err
}

// UpdateChildPrices 子供料金のみ更新
func (p *priceDirectRepository) UpdateChildPrices(inputData []price.HtTmPriceDirects) error {

	tmpPlaceHolder := []string{}
	var uv []interface{} = nil
	for i, v := range inputData {
		tmpPlaceHolder = append(tmpPlaceHolder, "(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		uv = append(uv, v.PriceID)
		uv = append(uv, v.PlanID)
		uv = append(uv, v.UseDate)
		uv = append(uv, v.RateTypeCode)
		uv = append(uv, v.Price)
		uv = append(uv, v.PriceInTax)
		uv = append(uv, v.ChildPrice1)
		uv = append(uv, v.ChildPrice1InTax)
		uv = append(uv, v.ChildPrice2)
		uv = append(uv, v.ChildPrice2InTax)
		uv = append(uv, v.ChildPrice3)
		uv = append(uv, v.ChildPrice3InTax)
		uv = append(uv, v.ChildPrice4)
		uv = append(uv, v.ChildPrice4InTax)
		uv = append(uv, v.ChildPrice5)
		uv = append(uv, v.ChildPrice5InTax)
		uv = append(uv, v.ChildPrice6)
		uv = append(uv, v.ChildPrice6InTax)
		uv = append(uv, v.RegularPrice)
		uv = append(uv, time.Now()) // updated_at
		// クエリの上限もしくは最後に到達したら
		if (i+1)%placeholderLimit == 0 || i == len(inputData)-1 {
			// inputDataはht_tm_price_directsから取得した既存のデータなので、DUPLICATE KEY UPDATEで更新のみ行われる(新規でINSERTされることはない)
			updateQuery := fmt.Sprintf(`INSERT INTO ht_tm_price_directs (
			                        price_id,
			                        plan_id,
			                        use_date,
			                        rate_type_code,
			                        price,
									price_in_tax,
			                        child_price1,
			                        child_price1_in_tax,
			                        child_price2,
			                        child_price2_in_tax,
			                        child_price3,
			                        child_price3_in_tax,
			                        child_price4,
			                        child_price4_in_tax,
			                        child_price5,
			                        child_price5_in_tax,
			                        child_price6,
			                        child_price6_in_tax,
			                        regular_price,
			                        updated_at)
			                      VALUES %s ON DUPLICATE KEY UPDATE
			                        child_price1=VALUES(child_price1),
			                        child_price1_in_tax=VALUES(child_price1_in_tax),
			                        child_price2=VALUES(child_price2),
			                        child_price2_in_tax=VALUES(child_price2_in_tax),
			                        child_price3=VALUES(child_price3),
			                        child_price3_in_tax=VALUES(child_price3_in_tax),
			                        child_price4=VALUES(child_price4),
			                        child_price4_in_tax=VALUES(child_price4_in_tax),
			                        child_price5=VALUES(child_price5),
			                        child_price5_in_tax=VALUES(child_price5_in_tax),
			                        child_price6=VALUES(child_price6),
			                        child_price6_in_tax=VALUES(child_price6_in_tax),
			                        updated_at=VALUES(updated_at)`, strings.Join(tmpPlaceHolder, ","))
			if err := p.db.Exec(updateQuery, uv...).Error; err != nil {
				return err
			}

			tmpPlaceHolder = []string{}
			uv = nil
		}
	}
	return nil
}
