package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// planRaku2Repository らく通プラン関連repository
type planRaku2Repository struct {
	db *gorm.DB
}

// NewPlanRaku2Repository インスタンス生成
func NewPlanRaku2Repository(db *gorm.DB) plan.IPlanRaku2Repository {
	return &planRaku2Repository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (p *planRaku2Repository) TxStart() (*gorm.DB, error) {
	tx := p.db.Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (p *planRaku2Repository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (p *planRaku2Repository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchAllByPropertyID 削除されていないproperty_idに紐づくプラン複数件取得
func (p *planRaku2Repository) FetchAllByPropertyID(req plan.ListInput) ([]plan.HtTmPlanRaku2s, error) {
	result := []plan.HtTmPlanRaku2s{}
	query := p.db.
		Table("ht_tm_plan_raku2s").
		Where("property_id = ? AND is_delete = 0", req.PropertyID)
	if req.Paging.Limit > 0 {
		query = query.Limit(req.Paging.Limit).Offset(req.Paging.Offset)
	}
	err := query.Find(&result).Error
	return result, err
}

// FetchAllByRoomTypeID 削除されていないroom_type_idに紐づくプラン複数件取得
func (p *planRaku2Repository) FetchAllByRoomTypeID(roomTypeID int64) ([]plan.HtTmPlanRaku2s, error) {
	result := []plan.HtTmPlanRaku2s{}
	err := p.db.
		Table("ht_tm_plan_raku2s").
		Where("room_type_id = ? AND is_delete = 0", roomTypeID).
		Find(&result).Error
	return result, err
}

// FetchAllByPlanGroupID plan_group_idに紐づくプラン複数件取得
func (p *planRaku2Repository) FetchAllByPlanGroupID(planGroupID int64) ([]plan.HtTmPlanRaku2s, error) {
	result := []plan.HtTmPlanRaku2s{}
	err := p.db.
		Table("ht_tm_plan_raku2s").
		Where("plan_group_id = ?", planGroupID).
		Find(&result).Error
	return result, err
}

// FetchActiveByPlanGroupID plan_group_idに紐づくアクティブなプラン複数件取得
func (p *planRaku2Repository) FetchActiveByPlanGroupID(planGroupID int64) ([]plan.HtTmPlanRaku2s, error) {
	result := []plan.HtTmPlanRaku2s{}
	err := p.db.
		Table("ht_tm_plan_raku2s").
		Where("plan_group_id = ? AND is_delete = 0", planGroupID).
		Find(&result).Error
	return result, err
}

// FetchAllByCancelPolicyID cancel_policy_idに紐づく削除されていないプランを複数件取得
func (p *planRaku2Repository) FetchAllByCancelPolicyID(cancelPolicyID uint64) ([]plan.HtTmPlanRaku2s, error) {
	result := []plan.HtTmPlanRaku2s{}
	err := p.db.
		Table("ht_tm_plan_raku2s").
		Joins("LEFT OUTER JOIN ht_th_plan_cancel_policy_relations ON ht_tm_plan_raku2s.plan_id = ht_th_plan_cancel_policy_relations.plan_id").
		Where("wholesaler_id = ? AND plan_cancel_policy_id = ? AND is_delete = 0", utils.WholesalerIDRaku2, cancelPolicyID).
		Find(&result).Error

	return result, err
}

// FetchOne plan_idに紐づく削除されていないプランを1件取得
func (p *planRaku2Repository) FetchOne(planID int64) (plan.HtTmPlanRaku2s, error) {
	result := plan.HtTmPlanRaku2s{}
	err := p.db.
		Table("ht_tm_plan_raku2s").
		Where("plan_id = ? AND is_delete = 0", planID).
		First(&result).Error
	return result, err
}

// FetchList plan_idに紐づく削除されていないプランを複数件取得
func (p *planRaku2Repository) FetchList(planIDList []int64) ([]plan.HtTmPlanRaku2s, error) {
	result := []plan.HtTmPlanRaku2s{}
	err := p.db.
		Table("ht_tm_plan_raku2s").
		Where("plan_id IN ? AND is_delete = 0", planIDList).
		Find(&result).Error
	return result, err
}

// MatchesPlanIDAndPropertyID propertyIDとplanIDが紐付いているか
func (p *planRaku2Repository) MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool {
	var result int64
	p.db.Model(&plan.HtTmPlanRaku2s{}).
		Where("plan_id = ?", planID).
		Where("property_id = ?", propertyID).
		Count(&result)
	return result > 0
}

// FetchChildRates plan_idに紐づく子供料金設定を複数件取得
func (p *planRaku2Repository) FetchChildRates(planID int64) ([]plan.HtTmChildRateRaku2s, error) {
	result := []plan.HtTmChildRateRaku2s{}
	err := p.db.
		Table("ht_tm_child_rate_raku2s").
		Where("plan_id = ?", planID).
		Find(&result).Error
	return result, err
}

// DeletePlanByRoomTypeID room_type_idに紐づくプランを論理削除
func (p *planRaku2Repository) DeletePlanByRoomTypeID(roomTypeID int64) error {
	return p.db.Model(&plan.HtTmPlanRaku2s{}).
		Where("room_type_id = ?", roomTypeID).
		Update("is_delete", 1).Error
}

// CheckPlanCode room_type_idとplan_codeの組み合わせで合致するものの件数を取得
func (p *planRaku2Repository) CheckPlanCode(propertyID int64, planCodeList []plan.CheckDuplicatePlanCode) int64 {
	var result int64
	q := p.db.Model(&plan.HtTmPlanRaku2s{})
	planCondition := p.db.Model(&plan.HtTmPlanRaku2s{})
	for _, v := range planCodeList {
		planCondition = planCondition.Or(map[string]interface{}{"room_type_id": v.RoomTypeID, "plan_code": v.PlanCode, "is_delete": 0})
	}
	// GroupConditions https://gorm.io/docs/advanced_query.html
	q.Where(planCondition.Where("property_id  = ?", propertyID)).Count(&result)
	return result
}

// selectNextPlanGroupID 次点のプラングループIDの取得
func (p *planRaku2Repository) selectNextPlanGroupID() (plan.HtTmPlanGroupIDRaku2s, error) {
	var result plan.HtTmPlanGroupIDRaku2s
	err := p.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Table("ht_tm_plan_group_id_raku2s").
		Select("plan_group_id").
		First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

// updateNextPlanGroupID 次点のプラングループIDの更新
func (p *planRaku2Repository) updateNextPlanGroupID() error {
	return p.db.
		Exec("UPDATE ht_tm_plan_group_id_raku2s SET plan_group_id = plan_group_id + 1").Error
}

// fetchNextPlanGroupID 新規プラン作成用のプラングループID取得
func (p *planRaku2Repository) fetchNextPlanGroupID() (int64, error) {
	var planGroupID int64
	res, sErr := p.selectNextPlanGroupID()
	if sErr != nil {
		return planGroupID, sErr
	}
	planGroupID = res.PlanGroupID

	// 次の採番用にインクリメントしておく
	if uErr := p.updateNextPlanGroupID(); uErr != nil {
		return planGroupID, uErr
	}

	return planGroupID, nil
}

// MakePlansRaku2 プランを複数件作成
func (p *planRaku2Repository) MakePlansRaku2(planTables []plan.HtTmPlanRaku2s) error {
	return p.db.Create(&planTables).Error
}

// CreatePlansRaku2 プランを複数件新規作成
func (p *planRaku2Repository) CreatePlansRaku2(planTables []plan.HtTmPlanRaku2s) error {
	planGroupID, err := p.fetchNextPlanGroupID()
	if err != nil {
		return err
	}
	for i := 0; i < len(planTables); i++ {
		planTables[i].PlanGroupID = planGroupID
	}
	return p.MakePlansRaku2(planTables)
}

// CreateChildRateRaku2 子供料金を複数件新規作成
func (p *planRaku2Repository) CreateChildRateRaku2(childRates []plan.HtTmChildRateRaku2s) error {
	return p.db.Create(&childRates).Error
}

// UpdatePlanRaku2 プラン情報を更新
func (p *planRaku2Repository) UpdatePlanRaku2(planTable *plan.HtTmPlanRaku2s, planIDs []int64) error {
	return p.db.Model(&plan.HtTmPlanRaku2s{}).
		Where("plan_id IN ?", planIDs).
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
			"is_package":                 planTable.IsPackage,
			"is_no_cancel":               planTable.IsNoCancel,
			"is_delete":                  planTable.IsDelete,
			"updated_at":                 time.Now(),
		}).Error
}

// UpdateChildRateRaku2 子供料金設定を更新
func (p *planRaku2Repository) UpdateChildRateRaku2(child *plan.HtTmChildRateRaku2s) error {
	return p.db.Model(&plan.HtTmChildRateRaku2s{}).
		Where("child_rate_type = ? AND plan_id = ?", child.ChildRateType, child.PlanID).
		Updates(map[string]interface{}{
			"receive":       child.Receive,
			"rate_category": child.RateCategory,
			"rate":          child.Rate,
			"calc_category": child.CalcCategory,
			"updated_at":    time.Now(),
		}).Error
}

// FetchPlanGroupID plan_idに基づくplan_group_idの取得
func (p *planRaku2Repository) FetchPlanGroupID(planID int64) (int64, error) {
	result := plan.HtTmPlanRaku2s{}
	err := p.db.
		Table("ht_tm_plan_raku2s").
		Where("plan_id = ?", planID).
		First(&result).Error
	return result.PlanGroupID, err
}

// DeletePlanRaku2 プランを一斉に論理削除
func (p *planRaku2Repository) DeletePlanRaku2(planIDs []int64) error {
	return p.db.Model(&plan.HtTmPlanRaku2s{}).
		Where("plan_id IN ?", planIDs).
		Updates(map[string]interface{}{
			"is_delete": 1,
		}).Error
}

// UpdateStopSales プランの売止を更新
func (p *planRaku2Repository) UpdateStopSales(planIDList []int64, isStopSales bool) error {
	return p.db.Model(&plan.HtTmPlanRaku2s{}).
		Where("plan_id IN ?", planIDList).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}

// UpdateStopSalesByRoomID 部屋に紐づくプランの売止フラグを複数件更新
func (p *planRaku2Repository) UpdateStopSalesByRoomID(roomTypeID int64, isStopSales bool) error {
	return p.db.Model(&plan.HtTmPlanRaku2s{}).
		Where("room_type_Id = ?", roomTypeID).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}

// UpdateStopSalesByRoomIDList room_type_id(複数)に紐づくプランの売止フラグを複数件更新
func (p *planRaku2Repository) UpdateStopSalesByRoomIDList(roomTypeIDList []int64, isStopSales bool) error {
	return p.db.Model(&plan.HtTmPlanRaku2s{}).
		Where("room_type_Id IN ?", roomTypeIDList).
		Updates(map[string]interface{}{
			"is_stop_sales": isStopSales,
			"updated_at":    time.Now(),
		}).Error
}
