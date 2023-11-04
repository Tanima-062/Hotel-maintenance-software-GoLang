package notification

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtTmPropertyNotifications お知らせ情報テーブル
type HtTmPropertyNotifications struct {
	NotificationID int64  `gorm:"primaryKey;autoIncrement:true" json:"notification_id,omitempty" validate:"isdefault"`
	NoticeTitle    string `json:"notice_title" validate:"max=128,required"`
	NoticeBody     string `json:"notice_body" validate:"required"`
	Status         int    `json:"status" validate:"oneof=0 1,number"`
	IsLatest       bool   `json:"is_latest"`
	common.Times   `gorm:"embedded" validate:"-"`
}

// DetailInput 詳細取得の入力
type DetailInput struct {
	NotificationID int64 `param:"id" validate:"required"`
}

// ListOutput 一覧の出力
type ListOutput struct {
	NotificationID int64  `json:"notification_id,omitempty"`
	NoticeTitle    string `json:"notice_title"`
	Status         int    `json:"status"`
	IsLatest       bool   `json:"is_latest"`
	common.Times   `gorm:"embedded"`
}

// INotificationUsecase お知らせ関連のusecaseのインターフェース
type INotificationUsecase interface {
	FetchList(paging *common.Paging) ([]ListOutput, error)
	FetchDetail(request *DetailInput) (HtTmPropertyNotifications, error)
	Create(notifications []HtTmPropertyNotifications) error
}

// INotificationRepository お知らせ関連のrepositoryのインターフェース
type INotificationRepository interface {
	// FetchAll limit, offsetに従って複数件取得
	FetchAll(paging *common.Paging) ([]ListOutput, error)
	// FetchOne 主キーで1件取得
	FetchOne(request *DetailInput) (HtTmPropertyNotifications, error)
	// BatchInsert お知らせ複数件作成
	BatchInsert(notifications []HtTmPropertyNotifications) error
}
