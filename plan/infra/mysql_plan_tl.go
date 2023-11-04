package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"gorm.io/gorm"
)

// planTLRepository Tl plan related repository
type planTlRepository struct {
	db *gorm.DB
}

// NewPlanTLRepository instantiation
func NewPlanTlRepository(db *gorm.DB) plan.IPlanTlRepository {
	return &planTlRepository{
		db: db,
	}
}

// TxStart transaction start
func (p *planTlRepository) TxStart() (*gorm.DB, error) {
	tx := p.db.Begin()
	return tx, tx.Error
}

// TxCommit transaction commit
func (p *planTlRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback transaction rollback
func (p *planTlRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchAllByPropertyID Acquire multiple plans linked to property_id that has not been deleted
func (p *planTlRepository) FetchAllByPropertyID(req plan.ListInput) ([]price.HtTmPlanTls, error) {
	result := []price.HtTmPlanTls{}
	query := p.db.
		Table("ht_tm_plan_tls").
		Where("property_id = ? AND is_delete = 0", req.PropertyID)
	if req.Paging.Limit > 0 {
		query = query.Limit(req.Paging.Limit).Offset(req.Paging.Offset)
	}
	err := query.Find(&result).Error
	return result, err
}

// FetchActiveByPlanCode Get multiple active plans linked to plan_group_id
func (p *planTlRepository) FetchActiveByPlanCode(planCode string) ([]price.HtTmPlanTls, error) {
	result := []price.HtTmPlanTls{}
	err := p.db.
		Table("ht_tm_plan_tls").
		Where("plan_code = ? AND is_delete = 0", planCode).
		Find(&result).Error
	return result, err
}

// FetchOne Get one undeleted plan associated with plan_id
func (p *planTlRepository) FetchOne(planID int64) (price.HtTmPlanTls, error) {
	result := price.HtTmPlanTls{}
	err := p.db.
		Table("ht_tm_plan_tls").
		Where("plan_id = ? AND is_delete = 0", planID).
		First(&result).Error
	return result, err
}

// FetchList Get multiple undeleted plans associated with plan_id
func (p *planTlRepository) FetchList(planIDList []int64) ([]price.HtTmPlanTls, error) {
	result := []price.HtTmPlanTls{}
	err := p.db.
		Table("ht_tm_plan_tls").
		Where("plan_id IN ? AND is_delete = 0", planIDList).
		Find(&result).Error
	return result, err
}

// MatchesPlanIDAndPropertyID Are propertyID and planID linked?
func (p *planTlRepository) MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool {
	var result int64
	p.db.Model(&price.HtTmPlanTls{}).
		Where("plan_id = ?", planID).
		Where("property_id = ?", propertyID).
		Count(&result)
	return result > 0
}

// FetchChildRates Get multiple child price settings linked to plan_id
func (p *planTlRepository) FetchChildRates(planID int64) ([]price.HtTmChildRateTls, error) {
	result := []price.HtTmChildRateTls{}
	err := p.db.
		Table("ht_tm_child_rate_tls").
		Where("plan_id = ?", planID).
		Find(&result).Error
	return result, err
}

// CreateChildRateTL Create multiple child fares
func (p *planTlRepository) CreateChildRateTl(childRates []price.HtTmChildRateTls) error {
	return p.db.Create(&childRates).Error
}

// GetPlanIfPlanCodeExist Is plan exist
func (p *planTlRepository) GetPlanIfPlanCodeExist(propertyID int64, planCode string, roomTypeID int64) (price.HtTmPlanTls, error) {
	result := price.HtTmPlanTls{}
	err := p.db.Table("ht_tm_plan_tls").
		Select("plan_id, room_type_id").
		Where("property_id  = ? And plan_code = ? And room_type_id = ? And is_delete=0", propertyID, planCode, roomTypeID).
		First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

// CreatePlansTl Create plan
func (p *planTlRepository) CreatePlanBulkTl(planTable price.HtTmPlanTls) error {
	return p.db.Create(&planTable).Error
}

// FetchPlan Get plan with plan_id
func (p *planTlRepository) FetchPlan(planID int64) (price.HtTmPlanTls, error) {
	var result price.HtTmPlanTls
	err := p.db.
		Table("ht_tm_plan_tls").
		Where("plan_id = ? AND is_delete = 0", planID).
		First(&result).Error
	return result, err
}

func (p *planTlRepository) UpdatePlanBulkTl(planTable price.HtTmPlanTls, planID int64) error {
	return p.db.Model(&price.HtTmPlanTls{}).
		Where("plan_id = ?", planID).
		Updates(map[string]interface{}{
			"name":                       planTable.Name,
			"description":                planTable.Description,
			"tax_category":               planTable.TaxCategory,
			"charge_category":            planTable.ChargeCategory,
			"accommodation_period_start": planTable.AccommodationPeriodStart,
			"accommodation_period_end":   planTable.AccommodationPeriodEnd,
			"is_accommodated_year_round": planTable.IsAccommodatedYearRound,
			"publishing_start_date":      planTable.PublishingStartDate,
			"publishing_end_date":        planTable.PublishingEndDate,
			"is_published_year_round":    planTable.IsPublishedYearRound,
			"reserve_accept_date":        planTable.ReserveAcceptDate,
			"reserve_accept_time":        planTable.ReserveAcceptTime,
			"reserve_deadline_date":      planTable.ReserveDeadlineDate,
			"reserve_deadline_time":      planTable.ReserveDeadlineTime,
			"min_stay_category":          planTable.MinStayCategory,
			"min_stay_num":               planTable.MinStayNum,
			"max_stay_category":          planTable.MaxStayCategory,
			"max_stay_num":               planTable.MaxStayNum,
			"meal_condition_breakfast":   planTable.MealConditionBreakfast,
			"meal_condition_dinner":      planTable.MealConditionDinner,
			"meal_condition_lunch":       planTable.MealConditionLunch,
			"is_no_cancel":               planTable.IsNoCancel,
			"is_delete":                  planTable.IsDelete,
			"cancel_policy":              planTable.CancelPolicy,
			"is_stop_sales":              planTable.IsStopSales,
			"updated_at":                 time.Now(),
		}).Error
}

func (p *planTlRepository) ClearChildRateTl(planID int64) error {
	return p.db.Delete(&plan.HtTmChildRateTls{}, "plan_id = ?", planID).Error
}

func (p *planTlRepository) ClearImageTl(planID int64) error {
	return p.db.Delete(&image.HtTmPlanOwnImagesTls{}, "plan_id = ?", planID).Error
}

func (p *planTlRepository) DeletePlanTl(planCode string, roomTypeIDs []int64) error {
	return p.db.Model(&price.HtTmPlanTls{}).
		Where("plan_code = ? And room_type_id NOT IN ?", planCode, roomTypeIDs).
		Updates(map[string]interface{}{
			"is_delete": 1,
		}).Error
}

func (p *planTlRepository) GetNextPlanID() (price.HtTmPlanTls, error) {
	result := price.HtTmPlanTls{}
	err := p.db.Select("plan_id").Last(&result).Error
	return result, err
}

func (p *planTlRepository) GetPlanByPropertyIDAndPlanCodeAndRoomTypeCode(propertyID int64, planCode string, roomTypeCode string) ([]price.HtTmPlanTls, error) {
	result := []price.HtTmPlanTls{}
	err := p.db.
		Table("ht_tm_plan_tls AS plan").
		Joins("INNER JOIN ht_tm_room_type_tls AS room ON plan.room_type_id = room.room_type_id").
		Where("plan.property_id = ? AND plan.plan_code = ? AND room.room_type_code = ?", propertyID, planCode, roomTypeCode).
		Find(&result).Error
	return result, err
}
