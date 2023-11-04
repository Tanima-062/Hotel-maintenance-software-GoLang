package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"gorm.io/gorm"
	"time"
)

// priceTLRepository
type priceTlRepository struct {
	db *gorm.DB
}

// NewPriceTLRepository
func NewPriceTlRepository(db *gorm.DB) price.IPriceTlRepository {
	return &priceTlRepository{
		db: db,
	}
}

// TxStart transaction start
func (p *priceTlRepository) TxStart() (*gorm.DB, error) {
	tx := p.db.Begin()
	return tx, tx.Error
}

// TxCommit transaction commit
func (p *priceTlRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback transaction rollback
func (p *priceTlRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchChildRates Get multiple price settings linked to a plan
func (p *priceTlRepository) FetchChildRates(planID int64) ([]price.HtTmChildRateTls, error) {
	result := []price.HtTmChildRateTls{}
	err := p.db.
		Model(&price.HtTmChildRateTls{}).
		Where("plan_id = ?", planID).
		Find(&result).Error
	return result, err
}

// FetchAllByPlanIDList Get multiple charges associated with multiple plan IDs within the period
func (p *priceTlRepository) FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]price.HtTmPriceTls, error) {
	result := []price.HtTmPriceTls{}
	err := p.db.
		Table("ht_tm_price_tls").
		Where("plan_id IN ?", planIDList).
		Where("use_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return result, err
}

// FetchPricesByPlanID Get multiple charges from today onwards
func (p *priceTlRepository) FetchPricesByPlanID(planID int64) ([]price.HtTmPriceTls, error) {
	result := []price.HtTmPriceTls{}
	err := p.db.
		Table("ht_tm_price_tls").
		Where("plan_id = ? ", planID).
		Where("use_date >= ?", time.Now().Format("2006-01-02")).
		Order("use_date ASC").
		Find(&result).Error
	return result, err
}

// FetchChildRatesByPlanID Get multiple price settings linked to plan
func (p *priceTlRepository) FetchChildRatesByPlanID(planID int64) ([]price.HtTmChildRateTls, error) {
	result := []price.HtTmChildRateTls{}
	err := p.db.
		Model(&price.HtTmChildRateTls{}).
		Where("plan_id = ?", planID).
		Find(&result).Error
	return result, err
}

// Get plan id by plan code
func (p *priceTlRepository) GetPlanIdByPlanCode(planID int64) ([]price.HtTmPriceTls, error) {
	result := []price.HtTmPriceTls{}
	err := p.db.
		Table("ht_tm_price_tls").
		Where("plan_id = ? ", planID).
		Where("use_date >= ?", time.Now().Format("2006-01-02")).
		Order("use_date ASC").
		Find(&result).Error
	return result, err
}

func (p *priceTlRepository) UpdatePrice(planID int64, useDate string, rateTypeCode string, price int, isStopSales bool, priceData price.HtTmPriceTls) error {
	err := p.db.Table("ht_tm_price_tls").
		Where("plan_id = ?", planID).
		Where("use_date = ?", useDate).
		Where("rate_type_code = ?", rateTypeCode).
		Updates(map[string]interface{}{
			"price":               price,
			"price_in_tax":        priceData.PriceTable.PriceInTax,
			"child_price1":        priceData.PriceTable.ChildPrice1,
			"child_price1_in_tax": priceData.PriceTable.ChildPrice1InTax,
			"child_price2":        priceData.PriceTable.ChildPrice2,
			"child_price2_in_tax": priceData.PriceTable.ChildPrice2InTax,
			"child_price3":        priceData.PriceTable.ChildPrice3,
			"child_price3_in_tax": priceData.PriceTable.ChildPrice3InTax,
			"child_price4":        priceData.PriceTable.ChildPrice4,
			"child_price4_in_tax": priceData.PriceTable.ChildPrice4InTax,
			"child_price5":        priceData.PriceTable.ChildPrice5,
			"child_price5_in_tax": priceData.PriceTable.ChildPrice5InTax,
			"child_price6":        priceData.PriceTable.ChildPrice6,
			"child_price6_in_tax": priceData.PriceTable.ChildPrice6InTax,
			"regular_price":       priceData.PriceTable.RegularPrice,
			"is_stop_sales":       isStopSales,
			"updated_at":          time.Now(),
		}).Error
	return err
}

func (p *priceTlRepository) CheckIfPriceExistsByPlanIDAndRateTypeCodeAndUseDate(planID int64, rateTypeCode string, useDate string) (bool, error) {
	result := price.HtTmPriceTls{}
	err := p.db.
		Table("ht_tm_price_tls").
		Where("plan_id = ? ", planID).
		Where("rate_type_code = ? ", rateTypeCode).
		Where("use_date = ?", useDate).
		First(&result).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// get price by property_id and plan code
func (p *priceTlRepository) GetPriceByPropertyIDAndPlanCode(propertyID int64, planCode string) ([]price.HtTmPriceTls, error) {
	planResult := []price.HtTmPlanTls{}
	var planIds []int64
	err := p.db.
		Select("plan_id").
		Table("ht_tm_plan_tls").
		Where("property_id = ? AND plan_code = ? AND is_delete = 0", propertyID, planCode).
		Find(&planResult).Error

	if err != nil {
		return nil, err
	} else {
		for _, planData := range planResult {
			planIds = append(planIds, planData.PlanID)
		}
	}
	result := []price.HtTmPriceTls{}
	errP := p.db.
		Table("ht_tm_price_tls").
		Where("plan_id IN ?", planIds).
		Find(&result).Error
	if errP != nil {
		return nil, err
	}
	return result, errP
}

// get price by plan id and rate_type_code
func (p *priceTlRepository) GetPriceByPlanIDAndRateTypeCode(planID int64, rateTypeCode string) (price.HtTmPriceTls, error) {
	result := price.HtTmPriceTls{}
	errP := p.db.
		Table("ht_tm_price_tls").
		Where("plan_id = ?", planID).
		Where("rate_type_code = ?", rateTypeCode).
		First(&result).Error

	return result, errP
}

// create price
func (p *priceTlRepository) CreatePrice(priceTable price.HtTmPriceTls) error {
	return p.db.Create(&priceTable).Error
}
