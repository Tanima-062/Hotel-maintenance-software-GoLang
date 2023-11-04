package handler

import (
	"net/http"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/booking"
	"github.com/Adventureinc/hotel-hm-api/src/booking/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// BookingHandler 予約関連の振り分け
type BookingHandler struct {
	BUsecase booking.IBookingUsecase
	AUsecase account.IAccountUsecase
}

// NewBookingHandler インスタンス生成
func NewBookingHandler(hotelDB *gorm.DB) *BookingHandler {
	return &BookingHandler{
		BUsecase: usecase.NewBookingUsecase(hotelDB),
		AUsecase: aUsecase.NewAccountUsecase(hotelDB),
	}
}

// Search 予約検索
func (b *BookingHandler) Search(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	hmUser, err := b.AUsecase.FetchHMUserByToken(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &booking.SearchInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)
	bookings, _ := b.BUsecase.SearchBookings(&hmUser, claimParam, *request)
	return c.JSON(http.StatusOK, bookings)
}

// Download 予約一覧CSVダウンロードで詳細情報をリストで取得
func (b *BookingHandler) Download(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	hmUser, err := b.AUsecase.FetchHMUserByToken(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &booking.DownloadInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)
	detail, err := b.BUsecase.BookingDownloads(&hmUser, claimParam, *request)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, detail)
}

// Detail 予約詳細
func (b *BookingHandler) Detail(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	hmUser, err := b.AUsecase.FetchHMUserByToken(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &booking.DetailInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)
	detail, err := b.BUsecase.DetailBooking(&hmUser, claimParam, *request)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, detail)
}

// Cancel 予約キャンセル
func (b *BookingHandler) Cancel(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	_, hmErr := b.AUsecase.FetchHMUserByToken(claimParam)
	if hmErr != nil {
		c.Echo().Logger.Error(hmErr)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &booking.CancelInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	success, err := b.BUsecase.CancelBooking(*request)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}
	if success == false {
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// NoShow 予約のNoShow(無断不泊)
func (b *BookingHandler) NoShow(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	_, hmErr := b.AUsecase.FetchHMUserByToken(claimParam)
	if hmErr != nil {
		c.Echo().Logger.Error(hmErr)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &booking.NoShowInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	if err := b.BUsecase.UpdateNoShow(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// getHmUser トークンからHMアカウント情報を取得
func (b *BookingHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return b.AUsecase.FetchHMUserByToken(claimParam)
}
