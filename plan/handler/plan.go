package handler

import (
	"github.com/Adventureinc/hotel-hm-api/src/common/log"
	"github.com/Adventureinc/hotel-hm-api/src/common/log/infra"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"net/http"
	"strconv"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"github.com/Adventureinc/hotel-hm-api/src/plan/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// PlanHandler プラン関連の振り分け
type PlanHandler struct {
	PNeppanUsecase   plan.IPlanUsecase
	PBulkTlUsecase   plan.IPlanBulkUsecase
	PBulkTemaUsecase plan.IPlanBulkTemaUsecase
	PDirectUsecase   plan.IPlanUsecase
	PRaku2Usecase    plan.IPlanUsecase
	AUsecase         account.IAccountUsecase
	RLogRepository   log.ILogRepository
}

var requestDataTl []price.PlanData
var requestDataTema []price.TemaPlanData

// NewPlanHandler インスタンス生成
func NewPlanHandler(db *gorm.DB) *PlanHandler {
	return &PlanHandler{
		PNeppanUsecase:   usecase.NewPlanNeppanUsecase(db),
		PBulkTlUsecase:   usecase.NewPlanTlUsecase(db),
		PBulkTemaUsecase: usecase.NewPlanTemaUsecase(db),
		PDirectUsecase:   usecase.NewPlanDirectUsecase(db),
		PRaku2Usecase:    usecase.NewPlanRaku2Usecase(db),
		AUsecase:         aUsecase.NewAccountUsecase(db),
		RLogRepository:   infra.NewLogRepository(db),
	}
}

// List 一覧取得
func (p *PlanHandler) List(c echo.Context) error {
	hmUser, err := p.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &plan.ListInput{}
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
		list, _ := p.PNeppanUsecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDTl:
		list, _ := p.PBulkTlUsecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDDirect:
		list, _ := p.PDirectUsecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDRaku2:
		list, _ := p.PRaku2Usecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDTema:

		list, _ := p.PBulkTemaUsecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	}
	return echo.ErrInternalServerError
}

// Detail 詳細取得
func (p *PlanHandler) Detail(c echo.Context) error {
	hmUser, err := p.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &plan.DetailInput{}
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
		detail, err := p.PNeppanUsecase.Detail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)

	case utils.WholesalerIDTl:
		detail, err := p.PBulkTlUsecase.Detail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)
	case utils.WholesalerIDTema:
		detail, err := p.PBulkTemaUsecase.Detail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)
	case utils.WholesalerIDDirect:
		detail, err := p.PDirectUsecase.Detail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)
	case utils.WholesalerIDRaku2:
		detail, err := p.PRaku2Usecase.Detail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)
	}
	return echo.ErrInternalServerError
}

// Create 新規作成
func (p *PlanHandler) Create(c echo.Context) error {
	hmUser, err := p.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &plan.SaveInput{}
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
		if err := p.PNeppanUsecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	case utils.WholesalerIDDirect:
		if err := p.PDirectUsecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	case utils.WholesalerIDRaku2:
		if err := p.PRaku2Usecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// Update 更新
func (p *PlanHandler) Update(c echo.Context) error {
	hmUser, err := p.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &plan.SaveInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	if request.PlanID == 0 {
		c.Echo().Logger.Error("plan_id がありません。")
		return echo.ErrBadRequest
	}

	switch hmUser.WholesalerID {
	case utils.WholesalerIDNeppan:
		if err := p.PNeppanUsecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	case utils.WholesalerIDDirect:
		if err := p.PDirectUsecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	case utils.WholesalerIDRaku2:
		if err := p.PRaku2Usecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// Delete 削除
func (p *PlanHandler) Delete(c echo.Context) error {
	hmUser, err := p.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &plan.DeleteInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	if request.PlanID == 0 {
		c.Echo().Logger.Error("room_type_id がありません。")
		return echo.ErrBadRequest
	}

	switch hmUser.WholesalerID {
	case utils.WholesalerIDNeppan:
		if err := p.PNeppanUsecase.Delete(request.PlanID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := p.PDirectUsecase.Delete(request.PlanID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDRaku2:
		if err := p.PRaku2Usecase.Delete(request.PlanID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// UpdateStopSales プランの売止更新
func (p *PlanHandler) UpdateStopSales(c echo.Context) error {
	hmUser, err := p.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &plan.StopSalesInput{}
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
		if err := p.PNeppanUsecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := p.PDirectUsecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDRaku2:
		if err := p.PRaku2Usecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// getHmUser トークンからHMアカウント情報を取得
func (p *PlanHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return p.AUsecase.FetchHMUserByToken(claimParam)
}

// CreateOrUpdateBulk create or update plan data from bulk request
func (p *PlanHandler) CreateOrUpdateBulk(c echo.Context) error {
	wholesalerId, _ := strconv.Atoi(c.Request().Header.Get("Wholesaler-Id"))

	switch wholesalerId {
	case utils.WholesalerIDTl:
		if err := c.Bind(&requestDataTl); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrBadRequest
		}

		// validate request data
		errorMessages := utils.Validate(c, requestDataTl)
		if len(errorMessages) > 0 {
			errorMessageShow := utils.ErrorMessageShow{
				Message: "Unprocessable entity",
				Errors:  errorMessages,
			}
			return c.JSON(http.StatusUnprocessableEntity, errorMessageShow)
		}
		// log request payload
		utils.RequestLog(c, requestDataTl)
	case utils.WholesalerIDTema:
		if err := c.Bind(&requestDataTema); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrBadRequest
		}

		// validate request data
		errorMessages := utils.Validate(c, requestDataTema)
		if len(errorMessages) > 0 {
			errorMessageShow := utils.ErrorMessageShow{
				Message: "Unprocessable entity",
				Errors:  errorMessages,
			}
			return c.JSON(http.StatusUnprocessableEntity, errorMessageShow)
		}
		// log request payload
		utils.RequestLog(c, requestDataTema)

	}
	go p.processBulkRequest(wholesalerId, c)

	return c.JSON(http.StatusOK, map[string]string{"message": "Request accepted successfully!"})
}

func (p *PlanHandler) processBulkRequest(wholesalerId int, c echo.Context) {

	ProcessStartTime := time.Now()
	hostUrl := c.Request().Host
	ActivityLogID, _ := p.RLogRepository.StoreBulkActivityLog(utils.LogServicePlan, utils.LogTypeMaster, hostUrl, time.Now())

	switch wholesalerId {
	case utils.WholesalerIDTl:
		if msg, err := p.PBulkTlUsecase.CreateBulk(requestDataTl); err != nil {
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
		if msg, err := p.PBulkTemaUsecase.CreateBulk(requestDataTema); err != nil {
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
