package handler

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/log"
	"github.com/Adventureinc/hotel-hm-api/src/common/log/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	"github.com/Adventureinc/hotel-hm-api/src/stock/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// StockHandler 在庫関連の振り分け
type StockHandler struct {
	SNeppanUsecase stock.IStockUsecase
	STlUsecase     stock.IStockUsecase
	STemaUsecase   stock.IStockTemaUsecase
	SDirectUsecase stock.IStockUsecase
	SRaku2Usecase  stock.IStockUsecase
	AUsecase       account.IAccountUsecase
	RLogRepository log.ILogRepository
}

var requestStockTl []stock.StockData
var requestStockTema []stock.StockDataTema

// NewStockHandler インスタンス生成
func NewStockHandler(db *gorm.DB) *StockHandler {
	return &StockHandler{
		SNeppanUsecase: usecase.NewStockNeppanUsecase(db),
		STlUsecase:     usecase.NewStockTlUsecase(db),
		STemaUsecase:   usecase.NewStockTemaUsecase(db),
		SDirectUsecase: usecase.NewStockDirectUsecase(db),
		SRaku2Usecase:  usecase.NewStockRaku2Usecase(db),
		AUsecase:       aUsecase.NewAccountUsecase(db),
		RLogRepository: infra.NewLogRepository(db),
	}
}

// Calendar 在庫料金カレンダー
func (s *StockHandler) Calendar(c echo.Context) error {
	hmUser, err := s.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &stock.CalendarInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	switch hmUser.WholesalerID {
	case utils.WholesalerIDNeppan:
		cal, _ := s.SNeppanUsecase.FetchCalendar(hmUser, *request)
		return c.JSON(http.StatusOK, cal)
	case utils.WholesalerIDTl:
		cal, _ := s.STlUsecase.FetchCalendar(hmUser, *request)
		return c.JSON(http.StatusOK, cal)
	case utils.WholesalerIDDirect:
		cal, _ := s.SDirectUsecase.FetchCalendar(hmUser, *request)
		return c.JSON(http.StatusOK, cal)
	case utils.WholesalerIDRaku2:
		cal, _ := s.SRaku2Usecase.FetchCalendar(hmUser, *request)
		return c.JSON(http.StatusOK, cal)
	case utils.WholesalerIDTema:
		cal, _ := s.STemaUsecase.FetchCalendar(hmUser, *request)
		return c.JSON(http.StatusOK, cal)
	}
	return echo.ErrInternalServerError
}

// UpdateStopSales 在庫の売止
func (s *StockHandler) UpdateStopSales(c echo.Context) error {
	hmUser, err := s.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &stock.StopSalesInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	switch hmUser.WholesalerID {
	case utils.WholesalerIDNeppan:
		if err := s.SNeppanUsecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.NoContent(http.StatusOK)
	case utils.WholesalerIDDirect:
		if err := s.SDirectUsecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.NoContent(http.StatusOK)
	case utils.WholesalerIDRaku2:
		if err := s.SRaku2Usecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.NoContent(http.StatusOK)
	}
	return echo.ErrInternalServerError
}

// FetchAll 在庫の一覧
func (s *StockHandler) FetchAll(c echo.Context) error {
	hmUser, err := s.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &stock.ListInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	switch hmUser.WholesalerID {
	case utils.WholesalerIDDirect:
		detail, _ := s.SDirectUsecase.FetchAll(request)
		return c.JSON(http.StatusOK, detail)
	}
	return echo.ErrInternalServerError
}

// Save 在庫の保存
func (s *StockHandler) Save(c echo.Context) error {
	hmUser, err := s.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &[]stock.SaveInput{}
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
	utils.RequestLog(c, request)

	switch hmUser.WholesalerID {
	case utils.WholesalerIDDirect:
		if err := s.SDirectUsecase.Save(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.NoContent(http.StatusOK)
	}
	return echo.ErrInternalServerError
}

// UpdateBulk updates the bulk request with stock data
func (s *StockHandler) UpdateBulk(c echo.Context) error {

	wholesalerId, _ := strconv.Atoi(c.Request().Header.Get("Wholesaler-Id"))
	switch wholesalerId {
	case utils.WholesalerIDTl:
		if err := c.Bind(&requestStockTl); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrBadRequest
		}

		// validate request data
		errorMessages := utils.Validate(c, requestStockTl)
		if len(errorMessages) > 0 {
			errorMessageShow := utils.ErrorMessageShow{
				Message: "Unprocessable entity",
				Errors:  errorMessages,
			}
			return c.JSON(http.StatusUnprocessableEntity, errorMessageShow)
		}
		utils.RequestLog(c, requestStockTl)

	case utils.WholesalerIDTema:
		if err := c.Bind(&requestStockTema); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrBadRequest
		}

		// validate request data
		errorMessages := utils.Validate(c, requestStockTema)
		if len(errorMessages) > 0 {
			errorMessageShow := utils.ErrorMessageShow{
				Message: "Unprocessable entity",
				Errors:  errorMessages,
			}
			return c.JSON(http.StatusUnprocessableEntity, errorMessageShow)
		}
		utils.RequestLog(c, requestStockTema)

	}

	go s.processBulkRequest(wholesalerId, c)

	return c.JSON(http.StatusOK, map[string]string{"message": " Request accepted successfully!"})
}

// processBulkRequest sends a request to the specified service
func (s *StockHandler) processBulkRequest(wholesalerId int, c echo.Context) {

	ProcessStartTime := time.Now()
	hostUrl := c.Request().Host
	ActivityLogID, _ := s.RLogRepository.StoreBulkActivityLog(utils.LogServiceStock, utils.LogTypeDifferential, hostUrl, time.Now())

	switch wholesalerId {
	case utils.WholesalerIDTl:
		if err := s.STlUsecase.UpdateBulk(requestStockTl); err != nil {
			c.Echo().Logger.Error(err)
			// Log the error and set the error message
			errorMessage := err.Error()
			err := s.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
			if err != nil {
				return
			}
		} else {
			// Task succeeded, set IsSuccess to true
			err := s.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, true, "")
			if err != nil {
				return
			}
		}
	case utils.WholesalerIDTema:
		if err := s.STemaUsecase.UpdateBulkTema(requestStockTema); err != nil {
			c.Echo().Logger.Error(err)
			// Log the error and set the error message
			errorMessage := err.Error()
			err := s.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
			if err != nil {
				return
			}
		} else {
			// Task succeeded, set IsSuccess to true
			err := s.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, true, "")
			if err != nil {
				return
			}
		}
	default:
		c.Echo().Logger.Error("Invalid wholesalerId")
		// Log the error for invalid wholesalerID
		errorMessage := "Invalid wholesalerID"
		err := s.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
		if err != nil {
			return
		}
	}
}

// getHmUser トークンからHMアカウント情報を取得
func (s *StockHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return s.AUsecase.FetchHMUserByToken(claimParam)
}
