package usecase

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/notification"
	"github.com/Adventureinc/hotel-hm-api/src/notification/infra"
	"gorm.io/gorm"
)

// notificationUsecase お知らせ関連usecase
type notificationUsecase struct {
	NRepository notification.INotificationRepository
}

// NewNotificationUsecase インスタンス生成
func NewNotificationUsecase(db *gorm.DB) notification.INotificationUsecase {
	return &notificationUsecase{
		NRepository: infra.NewNotificationRepository(db),
	}
}

// 全件取得
func (n *notificationUsecase) FetchList(paging *common.Paging) ([]notification.ListOutput, error) {
	return n.NRepository.FetchAll(paging)
}

// 詳細
func (n *notificationUsecase) FetchDetail(request *notification.DetailInput) (notification.HtTmPropertyNotifications, error) {
	return n.NRepository.FetchOne(request)
}

// 新規作成
func (n *notificationUsecase) Create(notifications []notification.HtTmPropertyNotifications) error {
	for index := range notifications {
		notifications[index].CreatedAt = time.Now()
		notifications[index].UpdatedAt = time.Now()
	}
	return n.NRepository.BatchInsert(notifications)
}
