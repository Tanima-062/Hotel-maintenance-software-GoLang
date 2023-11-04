package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"gorm.io/gorm"
	"time"
)

// planTemaRepository てまプラン関連repository
type planTemaRepository struct {
	db *gorm.DB
}

// NewPlanTemaRepository インスタンス生成
func NewPlanTemaRepository(db *gorm.DB) plan.IPlanTemaRepository {
	return &planTemaRepository{
		db: db,
	}
}

// TxStart transaction start
func (p *planTemaRepository) TxStart() (*gorm.DB, error) {
	tx := p.db.Begin()
	return tx, tx.Error
}

// TxCommit transaction commit
func (p *planTemaRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback transaction rollback
func (p *planTemaRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchOne プランを1件取得
func (p *planTemaRepository) FetchOnePlan(propertyID int64, packagePlanCode int) (*plan.HtTmPlanTemas, error) {
	result := &plan.HtTmPlanTemas{}
	err := p.db.
		Table("ht_tm_plan_temas").
		Where("property_id = ?", propertyID).
		Where("package_plan_code = ?", packagePlanCode).
		First(&result).Error
	return result, err
}

// FetchList プランを複数件取得
func (p *planTemaRepository) FetchList(propertyID int64, packagePlanCodeList []int) ([]plan.HtTmPlanTemas, error) {
	result := []plan.HtTmPlanTemas{}
	err := p.db.
		Table("ht_tm_plan_temas").
		Where("property_id = ?", propertyID).
		Where("package_plan_code IN ?", packagePlanCodeList).
		Find(&result).Error
	return result, err
}

// GetPlanIfPlanCodeExist Is plan exist
func (p *planTemaRepository) GetPlanIfPlanCodeExist(propertyID int64, planCode int64, roomTypeID int64) (price.HtTmPlanTemas, error) {
	result := price.HtTmPlanTemas{}
	err := p.db.Table("ht_tm_plan_temas").
		Select("plan_tema_id, room_type_id").
		Where("property_id  = ? And package_plan_code = ? And room_type_id = ? And deleted_at IS NULL", propertyID, planCode, roomTypeID).
		First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (p *planTemaRepository) UpdatePlanBulkTema(planTable price.HtTmPlanTemas, planID int64) error {
	return p.db.Model(&price.HtTmPlanTemas{}).
		Where("plan_tema_id = ?", planID).
		Updates(map[string]interface{}{
			"plan_tema_id":                planTable.PlanID,
			"plan_name":                   planTable.PlanName,
			"property_id":                 planTable.PropertyID,
			"package_plan_code":           planTable.PackagePlanCode,
			"lang_cd":                     planTable.LangCd,
			"desc":                        planTable.Desc,
			"room_type_id":                planTable.RoomTypeID,
			"available":                   planTable.Available,
			"rate_type":                   planTable.RateType,
			"tax_category":                planTable.TaxCategory,
			"stay_limit_min":              planTable.StayLimitMin,
			"stay_limit_max":              planTable.StayLimitMax,
			"adv_bk_create_start_enabled": planTable.AdvBKCreateStartEnabled,
			"adv_bk_create_end_enabled":   planTable.AdvBKCreateEndEnabled,
			"adv_bk_create_start_d":       planTable.AdvBkCreateStartD,
			"adv_bk_create_end_d":         planTable.AdvBkCreateEndD,
			"adv_bk_create_start_h":       planTable.AdvBkCreateStartH,
			"adv_bk_create_end_h":         planTable.AdvBkCreateEndH,
			"adv_bk_create_start_m":       planTable.AdvBkCreateStartM,
			"adv_bk_create_end_m":         planTable.AdvBkCreateEndM,
			"adv_bk_modify_end_enabled":   planTable.AdvBkModifyEndEnabled,
			"adv_bk_modify_end_d":         planTable.AdvBkModifyEndD,
			"adv_bk_modify_end_m":         planTable.AdvBkModifyEndM,
			"adv_bk_modify_end_h":         planTable.AdvBkModifyEndH,
			"adv_bk_cancel_end_enabled":   planTable.AdvBkCancelEndEnabled,
			"adv_bk_cancel_end_h":         planTable.AdvBkCancelEndH,
			"adv_bk_cancel_end_m":         planTable.AdvBkCancelEndM,
			"adv_bk_cancel_end_d":         planTable.AdvBkCancelEndD,
			"is_accommodated_year_round":  planTable.IsAccommodatedYearRound,
			"is_published_year_round":     planTable.IsPublishedYearRound,
			"min_stay_category":           planTable.MinStayCategory,
			"max_stay_category":           planTable.MaxStayCategory,
			"meal_condition_breakfast":    planTable.MealConditionBreakfast,
			"meal_condition_dinner":       planTable.MealConditionDinner,
			"meal_condition_lunch":        planTable.MealConditionLunch,
			"children_acceptable1":        planTable.ChildrenAcceptable1,
			"children_acceptable2":        planTable.ChildrenAcceptable2,
			"children_acceptable3":        planTable.ChildrenAcceptable3,
			"children_acceptable4":        planTable.ChildrenAcceptable4,
			"children_acceptable5":        planTable.ChildrenAcceptable5,
			"service_charge_type":         planTable.ServiceChargeType,
			"service_charge_value":        planTable.ServiceChargeValue,
			"optional_items":              planTable.OptionalItems,
			"cancelpolicy_json":           planTable.CancelpolicyJson,
			"children_json":               planTable.ChildrenJson,
			"picture_json":                planTable.PictureJson,
			"checkin_time_start_h":        planTable.CheckinTimeStartH,
			"checkin_time_end_h":          planTable.CheckinTimeEndM,
			"checkin_time_start_m":        planTable.CheckinTimeStartM,
			"checkin_time_end_m":          planTable.CheckinTimeEndM,
			"listing_period_start_h":      planTable.ListingPeriodStartH,
			"listing_period_end_h":        planTable.ListingPeriodEndH,
			"listing_period_start_m":      planTable.ListingPeriodStartM,
			"listing_period_end_m":        planTable.ListingPeriodEndM,
			"reserve_period_start":        planTable.ReservePeriodStart,
			"reserve_period_end":          planTable.ReservePeriodEnd,
			"listing_period_start":        planTable.ListingPeriodStart,
			"listing_period_end":          planTable.ListingPeriodEnd,
			"payment":                     planTable.Payment,
			"plan_type":                   planTable.PlanType,
			"is_room_charge":              planTable.IsRoomCharge,
			"updated_at":                  time.Now(),
		}).Error
}

func (p *planTemaRepository) GetNextPlanID() (price.HtTmPlanTemas, error) {
	result := price.HtTmPlanTemas{}
	err := p.db.Select("plan_tema_id").Last(&result).Error
	return result, err
}

// CreatePlansTema Create plan
func (p *planTemaRepository) CreatePlanBulkTema(planTable price.HtTmPlanTemas) error {
	return p.db.Create(&planTable).Error
}

func (p *planTemaRepository) ClearChildRateTema(planID int64) error {
	return p.db.Delete(&plan.HtTmChildRateTemas{}, "plan_id = ?", planID).Error
}

func (p *planTemaRepository) ClearImageTema(planID int64) error {
	return p.db.Delete(&image.HtTmPlanOwnImagesTemas{}, "plan_id = ?", planID).Error
}

// CreateChildRateTema Create multiple child fares
func (p *planTemaRepository) CreateChildRateTema(childRates []price.HtTmChildRateTemas) error {
	return p.db.Create(&childRates).Error
}

func (p *planTemaRepository) DeletePlanTema(planCode int64, roomTypeIDs []int64) error {
	return p.db.Model(&price.HtTmPlanTemas{}).
		Where("package_plan_code = ? And room_type_id NOT IN ?", planCode, roomTypeIDs).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
		}).Error
}

// FetchAllByPropertyID Acquire multiple plans linked to property_id that has not been deleted
func (p *planTemaRepository) FetchAllByPropertyID(req plan.ListInput) ([]price.HtTmPlanTemas, error) {
	result := []price.HtTmPlanTemas{}
	query := p.db.
		Select("plan_tema_id, room_type_id, package_plan_code, plan_name, available, property_id").
		Table("ht_tm_plan_temas").
		Where("property_id = ?", req.PropertyID)
	if req.Paging.Limit > 0 {
		query = query.Limit(req.Paging.Limit).Offset(req.Paging.Offset)
	}
	err := query.Find(&result).Error
	return result, err
}

// MatchesPlanIDAndPropertyID Are propertyID and planID linked?
func (p *planTemaRepository) MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool {
	var result int64
	p.db.Model(&price.HtTmPlanTemas{}).
		Where("plan_tema_id = ?", planID).
		Where("property_id = ?", propertyID).
		Count(&result)
	return result > 0
}

// FetchActiveByPlanCode Get multiple active plans linked to plan_group_id
func (p *planTemaRepository) FetchActiveByPlanCode(planCode string) ([]price.HtTmPlanTemas, error) {
	result := []price.HtTmPlanTemas{}
	err := p.db.
		Table("ht_tm_plan_temas").
		Where("plan_code = ? AND is_delete = 0", planCode).
		Find(&result).Error
	return result, err
}

// FetchOne Get one undeleted plan associated with planID
func (p *planTemaRepository) FetchOneWithPlanID(planID int64) (price.HtTmPlanTemas, error) {
	result := price.HtTmPlanTemas{}
	err := p.db.
		Table("ht_tm_plan_temas").
		Where("plan_tema_id = ?", planID).
		First(&result).Error
	return result, err
}

// FetchChildRates Get multiple child price settings linked to plan_id
func (p *planTemaRepository) FetchChildRates(planID int64) ([]price.HtTmChildRateTemas, error) {
	result := []price.HtTmChildRateTemas{}
	err := p.db.
		Table("ht_tm_child_rate_temas").
		Where("plan_id = ?", planID).
		Find(&result).Error
	return result, err
}
