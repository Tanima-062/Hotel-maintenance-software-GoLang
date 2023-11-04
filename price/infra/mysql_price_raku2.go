package infra

import (
	"fmt"
	"strings"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/price"
	"gorm.io/gorm"
)

// priceRaku2Repository らく通料金関連repository
type priceRaku2Repository struct {
	db *gorm.DB
}

// NewPriceRaku2Repository インスタンス生成
func NewPriceRaku2Repository(db *gorm.DB) price.IPriceRaku2Repository {
	return &priceRaku2Repository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (p *priceRaku2Repository) TxStart() (*gorm.DB, error) {
	tx := p.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (p *priceRaku2Repository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (p *priceRaku2Repository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchChildRates プランに紐づく料金設定を複数件取得
func (p *priceRaku2Repository) FetchChildRates(planID int64) ([]price.HtTmChildRateRaku2s, error) {
	result := []price.HtTmChildRateRaku2s{}
	err := p.db.
		Model(&price.HtTmChildRateRaku2s{}).
		Where("plan_id = ?", planID).
		Find(&result).Error
	return result, err
}

// FetchAllByPlanIDList 期間内の複数のプランIDに紐づく料金を複数件取得
func (p *priceRaku2Repository) FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]price.HtTmPriceRaku2s, error) {
	result := []price.HtTmPriceRaku2s{}
	err := p.db.
		Table("ht_tm_price_raku2s").
		Where("plan_id IN ?", planIDList).
		Where("use_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return result, err
}

// FetchPricesByPlanID 本日以降の料金を複数件取得
func (p *priceRaku2Repository) FetchPricesByPlanID(planID int64) ([]price.HtTmPriceRaku2s, error) {
	result := []price.HtTmPriceRaku2s{}
	err := p.db.
		Table("ht_tm_price_raku2s").
		Where("plan_id = ? ", planID).
		Where("use_date >= ?", time.Now().Format("2006-01-02")).
		Order("use_date ASC").
		Find(&result).Error
	return result, err
}

// UpdateChildPrices 子供料金のみ更新
func (p *priceRaku2Repository) UpdateChildPrices(inputData []price.HtTmPriceRaku2s) error {
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
			// inputDataはht_tm_price_raku2sから取得した既存のデータなので、DUPLICATE KEY UPDATEで更新のみ行われる(新規でINSERTされることはない)
			updateQuery := fmt.Sprintf(`INSERT INTO ht_tm_price_raku2s (
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
