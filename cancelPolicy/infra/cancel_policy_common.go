package infra

import (
	"encoding/json"

	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	cp "github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommonCancelPolicyRepository struct {
	db *gorm.DB
}

func NewCommonCancelPolicyRepository(db *gorm.DB) cancelPolicy.ICancelPolicyCommonRepository {
	return &CommonCancelPolicyRepository{db: db}
}

// CreatePlanCancelPolicy はプランごとのキャンセルポリシーを作成します。
func (r *CommonCancelPolicyRepository) CreatePlanCancelPolicy(policyName string, wholesalerID int, propertyID int64, cancelPolicy string) error {
	row := cp.HtTmPlanCancelPolicies{
		PolicyName:        policyName,
		PropertyID:        propertyID,
		WholesalerID:      wholesalerID,
		CancelPenaltyJSON: cancelPolicy,
	}

	return r.db.Model(&cp.HtTmPlanCancelPolicies{}).Create(&row).Error
}

// FindPlanCancelPolicy はプランごとのキャンセルポリシーを返却します。
func (r *CommonCancelPolicyRepository) FindPlanCancelPolicy(policyId uint64) (cancelPolicy.CancelPolicyJSONWithName, error) {
	var row cp.HtTmPlanCancelPolicies
	var ret cp.CancelPolicyJSONWithName

	err := r.db.Model(&cp.HtTmPlanCancelPolicies{}).Where("plan_cancel_policy_id = ?", policyId).Find(&row).Error
	if err != nil {
		return ret, err
	}

	ret.CancelPolicyName = &row.PolicyName
	jsonData := []byte(row.CancelPenaltyJSON)
	if jErr := json.Unmarshal(jsonData, &ret.CancelPolicyJSON); jErr != nil {
		return ret, jErr
	}

	return ret, err
}

func (r *CommonCancelPolicyRepository) PlanCancelPolicyList(propertyID int64) ([]cancelPolicy.HtTmPlanCancelPolicies, error) {
	var ret []cp.HtTmPlanCancelPolicies
	err := r.db.Model(&cp.HtTmPlanCancelPolicies{}).Where("property_id = ?", propertyID).Find(&ret).Error

	return ret, err
}

func (r *CommonCancelPolicyRepository) UpdatePlanCancelPolicy(policyId uint64, policyName string, cancelPolicy string) error {
	return r.db.Model(&cp.HtTmPlanCancelPolicies{}).
		Where("plan_cancel_policy_id = ?", policyId).
		Updates(map[string]interface{}{
			"policy_name":         policyName,
			"cancel_penalty_json": cancelPolicy,
			//			"updated_at":          time.Now(),
		}).Error
}

func (r *CommonCancelPolicyRepository) DeletePlanCancelPolicy(policyId uint64) error {
	return r.db.Model(&cancelPolicy.HtTmPlanCancelPolicies{}).
		Delete("plan_cancel_policy_id = ?", policyId).Error
}

func (r *CommonCancelPolicyRepository) FindAssignedPlanCancelPolicy(propertyId int64, planId int64) (*cancelPolicy.HtThPlanCancelPolicyRelations, error) {
	var row cp.HtThPlanCancelPolicyRelations
	err := r.db.Model(&cancelPolicy.HtThPlanCancelPolicyRelations{}).
		Where("property_id", propertyId).
		Where("plan_id", planId).
		First(&row).Error

	return &row, err
}

func (r *CommonCancelPolicyRepository) UpsertPlanCancelPolicyRelation(wholesalerID int, propertyId int64, planId int64, plancancelPolicyId uint64) error {
	record := &cancelPolicy.HtThPlanCancelPolicyRelations{
		WholesalerID:       wholesalerID,
		PropertyID:         propertyId,
		PlanID:             planId,
		PlanCancelPolicyID: plancancelPolicyId,
	}

	return r.db.Model(&cancelPolicy.HtThPlanCancelPolicyRelations{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "wholesaler_id"}, {Name: "property_id"}, {Name: "plan_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"plan_cancel_policy_id": plancancelPolicyId}),
		}).Create(&record).Error
}

func (r *CommonCancelPolicyRepository) DeletePlanCancelPolicyRelation(wholesalerID int, planId int64) error {
	return r.db.Model(&cancelPolicy.HtThPlanCancelPolicyRelations{}).
		Where("wholesaler_id = ? AND plan_id = ?", wholesalerID, planId).
		Delete(&cancelPolicy.HtThPlanCancelPolicyRelations{}).Error
}
