package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/notification"
	"gorm.io/gorm"
)

// notificationRepository お知らせ関連repository
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository インスタンス生成
func NewNotificationRepository(db *gorm.DB) notification.INotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

// FetchAll limit, offsetに従って複数件取得
func (a *notificationRepository) FetchAll(paging *common.Paging) ([]notification.ListOutput, error) {
	result := []notification.ListOutput{}
	query := a.db.Model(&notification.HtTmPropertyNotifications{})
	if paging.Limit > 0 {
		query = query.Limit(paging.Limit).Offset(paging.Offset)
	}
	err := query.Find(&result).Error
	return result, err
}

// FetchOne 主キーで1件取得
func (a *notificationRepository) FetchOne(request *notification.DetailInput) (notification.HtTmPropertyNotifications, error) {
	result := notification.HtTmPropertyNotifications{}
	err := a.db.First(&result, request.NotificationID).Error
	return result, err
}

// BatchInsert お知らせ複数件作成
func (a *notificationRepository) BatchInsert(notifications []notification.HtTmPropertyNotifications) error {
	return a.db.Create(notifications).Error
}
