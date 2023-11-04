package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// planCommonRepository プラン共通処理のrepository
type planCommonRepository struct {
	db *gorm.DB
}

// NewPlanDirectRepository インスタンス生成
func NewPlanCommonRepository(db *gorm.DB) plan.ICommonPlanRepository {
	return &planCommonRepository{
		db: db,
	}
}

// プランのチェックイン/アウト日を取得します。
func (c *planCommonRepository) FetchCheckInOut(propertyId int64, planId int64) (*plan.HtTmPlanCheckInOuts, error) {
	var row plan.HtTmPlanCheckInOuts
	err := c.db.Model(&plan.HtTmPlanCheckInOuts{}).Where("property_id = ?", propertyId).Where("plan_id = ?", planId).First(&row).Error

	return &row, err
}

// プランのチェックイン/アウト日を登録/更新します。
func (c *planCommonRepository) UpsertCheckInOut(checkInOutInfo plan.CheckInOutInfo) error {
	record := &plan.HtTmPlanCheckInOuts{
		WholesalerID: checkInOutInfo.WholesalerID,
		PropertyID:   checkInOutInfo.PropertyID,
		PlanID:       checkInOutInfo.PlanID,
		CheckInBegin: checkInOutInfo.CheckInBegin,
		CheckInEnd:   checkInOutInfo.CheckInEnd,
		CheckOut:     checkInOutInfo.CheckOut,
	}

	updateColumns := map[string]interface{}{
		"check_in_begin": checkInOutInfo.CheckInBegin,
		"check_in_end":   checkInOutInfo.CheckInEnd,
		"check_out":      checkInOutInfo.CheckOut,
	}

	return c.db.Model(&plan.HtTmPlanCheckInOuts{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "wholesaler_id"}, {Name: "property_id"}, {Name: "plan_id"}},
			DoUpdates: clause.Assignments(updateColumns),
		}).Create(&record).Error
}

func (c *planCommonRepository) DeleteCheckInOut(wholesalerId int, planId int64) error {
	return c.db.
		Where("wholesaler_id = ? AND plan_id = ?", wholesalerId, planId).
		Delete(&plan.HtTmPlanCheckInOuts{}).Error
}
