package handler

import (
	"net/http"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/notification"
	"github.com/Adventureinc/hotel-hm-api/src/notification/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// NotificationHandler お知らせ関連の振り分け
type NotificationHandler struct {
	NUsecase notification.INotificationUsecase
}

// NewNotificationHandler インスタンス生成
func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	return &NotificationHandler{
		NUsecase: usecase.NewNotificationUsecase(db),
	}
}

// List 一覧取得
func (n *NotificationHandler) List(c echo.Context) error {
	request := &common.Paging{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	list, _ := n.NUsecase.FetchList(request)
	return c.JSON(http.StatusOK, list)
}

// Detail 詳細
func (n *NotificationHandler) Detail(c echo.Context) error {
	request := &notification.DetailInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	detail, err := n.NUsecase.FetchDetail(request)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, detail)
}

// Create 作成
func (n *NotificationHandler) Create(c echo.Context) error {
	request := &[]notification.HtTmPropertyNotifications{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	for _, r := range *request {
		if err := c.Validate(r); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrBadRequest
		}
	}

	if err := n.NUsecase.Create(*request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}
