package handler

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/log"
	"github.com/Adventureinc/hotel-hm-api/src/common/log/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"github.com/Adventureinc/hotel-hm-api/src/price/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// PriceHandler 料金関連の振り分け
type PriceHandler struct {
	PDirectUsecase price.IPriceUsecase
	PTlUsecase     price.IPriceBulkTlUsecase
	PTemaUsecase   price.IPriceBulkTemaUsecase
	AUsecase       account.IAccountUsecase
	RLogRepository log.ILogRepository
}

var requestTl []price.PriceData
var requestTema []price.PriceTemaData

// NewPriceHandler インスタンス生成
func NewPriceHandler(db *gorm.DB) *PriceHandler {
	return &PriceHandler{
		PDirectUsecase: usecase.NewPriceDirectUsecase(db),
		PTlUsecase:     usecase.NewPriceTlUsecase(db),
		PTemaUsecase:   usecase.NewPriceTemaUsecase(db),
		AUsecase:       aUsecase.NewAccountUsecase(db),
		RLogRepository: infra.NewLogRepository(db),
	}
}

// Detail 料金取得
func (p *PriceHandler) Detail(c echo.Context) error {
	hmUser, err := p.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &price.DetailInput{}
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
		detail, err := p.PDirectUsecase.FetchDetail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)
	}
	return echo.ErrInternalServerError
}

// Save 料金の作成・更新
func (p *PriceHandler) Save(c echo.Context) error {
	hmUser, err := p.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &[]price.SaveInput{}
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
		if err := p.PDirectUsecase.Save(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.NoContent(http.StatusOK)
	}
	return echo.ErrInternalServerError
}

// Update price from bulk data
func (p *PriceHandler) UpdateBulk(c echo.Context) error {
	wholesalerId, _ := strconv.Atoi(c.Request().Header.Get("Wholesaler-Id"))

	switch wholesalerId {
	case utils.WholesalerIDTl:
		if err := c.Bind(&requestTl); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrBadRequest
		}

		// validate request data
		errorMessages := utils.Validate(c, requestTl)
		if len(errorMessages) > 0 {
			errorMessageShow := utils.ErrorMessageShow{
				Message: "Unprocessable entity",
				Errors:  errorMessages,
			}
			return c.JSON(http.StatusUnprocessableEntity, errorMessageShow)
		}
		// log request payload
		utils.RequestLog(c, requestTl)
	case utils.WholesalerIDTema:
		if err := c.Bind(&requestTema); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrBadRequest
		}
		// validate request data
		errorMessages := utils.Validate(c, requestTema)
		if len(errorMessages) > 0 {
			errorMessageShow := utils.ErrorMessageShow{
				Message: "Unprocessable entity",
				Errors:  errorMessages,
			}
			return c.JSON(http.StatusUnprocessableEntity, errorMessageShow)
		}
		// log request payload
		utils.RequestLog(c, requestTema)
	}

	go p.processBulkRequest(wholesalerId, c)

	return c.JSON(http.StatusOK, map[string]string{"message": "Request accepted successfully!"})
}

func (p *PriceHandler) processBulkRequest(wholesalerId int, c echo.Context) {

	ProcessStartTime := time.Now()
	hostUrl := c.Request().Host
	ActivityLogID, _ := p.RLogRepository.StoreBulkActivityLog(utils.LogServicePrice, utils.LogTypeDifferential, hostUrl, time.Now())

	switch wholesalerId {
	case utils.WholesalerIDTl:
		if msg, err := p.PTlUsecase.Update(requestTl); err != nil {
			c.Echo().Logger.Error(err)
			c.Echo().Logger.Error(msg)
			// Log the error and set the error message
			errorMessage := err.Error()
			err := p.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
			if err != nil {
				return
			}
		} else {
			// Task succeeded, set IsSuccess to true
			err := p.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, true, "")
			if err != nil {
				return
			}
		}
	case utils.WholesalerIDTema:
		if msg, err := p.PTemaUsecase.Update(requestTema); err != nil {
			c.Echo().Logger.Error(err)
			c.Echo().Logger.Error(msg)
			// Log the error and set the error message
			errorMessage := err.Error()
			err := p.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
			if err != nil {
				return
			}
		} else {
			// Task succeeded, set IsSuccess to true
			err := p.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, true, "")
			if err != nil {
				return
			}
		}
	default:
		c.Echo().Logger.Error("Invalid wholesalerId")
		// Log the error for invalid wholesalerID
		errorMessage := "Invalid wholesalerID"
		err := p.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
		if err != nil {
			return
		}
	}
}

// getHmUser トークンからHMアカウント情報を取得
func (p *PriceHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return p.AUsecase.FetchHMUserByToken(claimParam)
}
