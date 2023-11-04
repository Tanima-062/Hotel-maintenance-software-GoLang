package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// cancelPolicyDirectRepository TLキャンセルポリシー関連repository
type cancelPolicyTlRepository struct {
	db *gorm.DB
}

// NewCancelPolicyTlRepository インスタンス生成
func NewCancelPolicyTlRepository(db *gorm.DB) cancelPolicy.ICancelPolicyTlRepository {
	return &cancelPolicyTlRepository{
		db: db,
	}
}

// Update キャンセルポリシー更新
func (c *cancelPolicyTlRepository) UpsertCancelPolicyTl(propertyID int64, cancelPolicy string) error {
	assignData := map[string]interface{}{
		"property_id":         propertyID,
		"lang_cd":             "ja-JP",
		"cancel_penalty_json": cancelPolicy,
	}

	return c.db.Model(&facility.HtTmPropertyTls{}).
		Where("property_id = ?", propertyID).
		Assign(assignData).
		FirstOrCreate(&facility.HtTmPropertyTls{}).
		Error
}
