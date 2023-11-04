package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// cancelPolicyDirectRepository 直仕入れキャンセルポリシー関連repository
type cancelPolicyDirectRepository struct {
	db *gorm.DB
}

// NewCancelPolicyDirectRepository インスタンス生成
func NewCancelPolicyDirectRepository(db *gorm.DB) cancelPolicy.ICancelPolicyDirectRepository {
	return &cancelPolicyDirectRepository{
		db: db,
	}
}

// Update キャンセルポリシー更新
func (c *cancelPolicyDirectRepository) Update(propertyID int64, cancelPolicy string) error {
	return c.db.Model(&facility.HtTmPropertyDirects{}).
		Where("property_id = ?", propertyID).
		Updates(map[string]interface{}{
			"cancel_penalty_json": cancelPolicy,
			"updated_at":          time.Now(),
		}).Error
}
