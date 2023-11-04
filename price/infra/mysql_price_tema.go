package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"gorm.io/gorm"
	"time"
)

// priceTemaRepository
type priceTemaRepository struct {
	db *gorm.DB
}

// NewpriceTemaRepository
func NewPriceTemaRepository(db *gorm.DB) price.IPriceTemaRepository {
	return &priceTemaRepository{
		db: db,
	}
}

// TxStart transaction start
func (p *priceTemaRepository) TxStart() (*gorm.DB, error) {
	tx := p.db.Begin()
	return tx, tx.Error
}

// TxCommit transaction commit
func (p *priceTemaRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback transaction rollback
func (p *priceTemaRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchAllByPlanCodeList Get multiple charges associated with multiple plan IDs within the period
func (p *priceTemaRepository) FetchAllByPlanCodeList(planCodeList []int64, startDate string, endDate string) ([]price.HtTmPriceTemas, error) {
	result := []price.HtTmPriceTemas{}
	err := p.db.
		Table("ht_tm_price_temas").
		Where("package_plan_code IN ?", planCodeList).
		Where("price_date BETWEEN ? AND ?", startDate, endDate).
		Find(&result).Error
	return result, err
}

// FetchPricesByPlanID Get multiple charges from today onwards
func (p *priceTemaRepository) FetchPricesByPlanID(planID int64) ([]price.HtTmPriceTemas, error) {
	result := []price.HtTmPriceTemas{}
	err := p.db.
		Table("ht_tm_price_temas").
		Where("plan_tema_id = ? ", planID).
		Where("price_date >= ?", time.Now().Format("2006-01-02")).
		Order("price_date ASC").
		Find(&result).Error
	return result, err
}

func (p *priceTemaRepository) DeletePriceTema(propertyID int64, packagePlanCode int64, roomTypeCode int, priceDate string) error {
	return p.db.Delete(&price.HtTmPriceTemas{}, "property_id = ? And package_plan_code = ? And room_type_code = ? And price_date = ?", propertyID, packagePlanCode, roomTypeCode, priceDate).Error
}

// create price
func (p *priceTemaRepository) CreatePrice(priceTable price.HtTmPriceTemas) error {
	return p.db.Create(&priceTable).Error
}

//func (p *priceTemaRepository) FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]price.HtTmPriceTemas, error) {
//	result := []price.HtTmPriceTemas{}
//	err := p.db.
//		Table("ht_tm_price_temas").
//		Where("package_plan_code IN ?", planIDList).
//		Where("price_date BETWEEN ? AND ?", startDate, endDate).
//		Find(&result).Error
//	return result, err
//}
