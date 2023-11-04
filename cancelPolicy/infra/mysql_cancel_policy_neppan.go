package infra

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"gorm.io/gorm"
)

// cancelPolicyDirectRepository ねっぱんキャンセルポリシー関連repository
type cancelPolicyNeppanRepository struct {
	db *gorm.DB
}

// NewCancelPolicyNeppanRepository インスタンス生成
func NewCancelPolicyNeppanRepository(db *gorm.DB) cancelPolicy.ICancelPolicyNeppanRepository {
	return &cancelPolicyNeppanRepository{
		db: db,
	}
}

// Update キャンセルポリシー更新
func (c *cancelPolicyNeppanRepository) Update(propertyID int64, cancelPolicy string) error {
	return c.db.Model(&facility.HtTmPropertyNeppans{}).
		Where("property_id = ?", propertyID).
		Updates(map[string]interface{}{
			"cancel_penalty_json": cancelPolicy,
			"updated_at":          time.Now(),
		}).Error
}
