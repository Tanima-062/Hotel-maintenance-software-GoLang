package handler

import (
	"net/http"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/settlement"
	"github.com/Adventureinc/hotel-hm-api/src/settlement/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SettlementHandler 請求関連の振り分け
type SettlementHandler struct {
	SUsecase settlement.ISettlementUsecase
	AUsecase account.IAccountUsecase
}

// NewSettlementHandler インスタンス生成
func NewSettlementHandler(db *gorm.DB) *SettlementHandler {
	return &SettlementHandler{
		SUsecase: usecase.NewSettlementUsecase(db),
		AUsecase: aUsecase.NewAccountUsecase(db),
	}
}

// List 請求書一覧
func (s *SettlementHandler) List(c echo.Context) error {
	_, err := s.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &settlement.ListInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	list, _ := s.SUsecase.FetchAll(*request)
	return c.JSON(http.StatusOK, list)
}

// Approve 請求書の承認
func (s *SettlementHandler) Approve(c echo.Context) error {
	_, err := s.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &settlement.UpdateInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	if err := s.SUsecase.Approve(*request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// Download 請求書ダウンロード
func (s *SettlementHandler) Download(c echo.Context) error {
	_, err := s.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &settlement.DownloadInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	tempFileName, downloadFileName, dErr := s.SUsecase.Download(request)
	if dErr != nil {
		c.Echo().Logger.Error(dErr)
		return echo.ErrInternalServerError
	}
	return c.Attachment(tempFileName, downloadFileName)
}

// FetchInfo 精算情報取得
func (s *SettlementHandler) FetchInfo(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	_, hmErr := s.AUsecase.FetchHMUserByToken(claimParam)
	if hmErr != nil {
		c.Echo().Logger.Error(hmErr)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &settlement.InfoInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	info, _ := s.SUsecase.FetchInfo(request, claimParam)
	return c.JSON(http.StatusOK, info)
}

// SaveInfo 精算情報更新
func (s *SettlementHandler) SaveInfo(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	_, hmErr := s.AUsecase.FetchHMUserByToken(claimParam)
	if hmErr != nil {
		c.Echo().Logger.Error(hmErr)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &settlement.SaveInfoInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	if err := s.SUsecase.SaveInfo(request, claimParam); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// getHmUser トークンからHMアカウント情報を取得
func (s *SettlementHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return s.AUsecase.FetchHMUserByToken(claimParam)
}
