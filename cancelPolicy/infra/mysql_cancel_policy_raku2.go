package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// cancelPolicyRaku2Repository らく通キャンセルポリシー関連repository
type cancelPolicyRaku2Repository struct {
	db *gorm.DB
}

// NewCancelPolicyRaku2Repository インスタンス生成
func NewCancelPolicyRaku2Repository(db *gorm.DB) cancelPolicy.ICancelPolicyRaku2Repository {
	return &cancelPolicyRaku2Repository{
		db: db,
	}
}

// Update キャンセルポリシー更新
func (c *cancelPolicyRaku2Repository) Update(propertyID int64, cancelPolicy string) error {
	return c.db.Model(&facility.HtTmPropertyRaku2s{}).
		Where("property_id = ?", propertyID).
		Updates(map[string]interface{}{
			"cancel_penalty_json": cancelPolicy,
			"updated_at":          time.Now(),
		}).Error
}
