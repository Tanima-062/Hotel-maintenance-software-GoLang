package handler

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/log"
	"github.com/Adventureinc/hotel-hm-api/src/common/log/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"github.com/Adventureinc/hotel-hm-api/src/room/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// RoomHandler 部屋関連の振り分け
type RoomHandler struct {
	RCommonUsecase  room.IRoomCommonUsecase
	RNeppanUsecase  room.IRoomUsecase
	RTlUsecase      room.IRoomBulkUsecase
	RDirectUsecase  room.IRoomUsecase
	RRaku2Usecase   room.IRoomUsecase
	AUsecase        account.IAccountUsecase
	RLogRepository  log.ILogRepository
	RTlRepository   room.IRoomTlRepository
	RTemaUsecase    room.IRoomTemaUseCase
	RTemaRepository room.IRoomTemaRepository
}

var requestDataTl []room.RoomData
var requestDataTema []room.RoomDataTema

// NewRoomHandler インスタンス生成
func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{
		RCommonUsecase:  usecase.NewRoomCommonUsecase(db),
		RNeppanUsecase:  usecase.NewRoomNeppanUsecase(db),
		RTlUsecase:      usecase.NewRoomTlUsecase(db),
		RDirectUsecase:  usecase.NewRoomDirectUsecase(db),
		RRaku2Usecase:   usecase.NewRoomRaku2Usecase(db),
		AUsecase:        aUsecase.NewAccountUsecase(db),
		RLogRepository:  infra.NewLogRepository(db),
		RTlRepository:   rInfra.NewRoomTlRepository(db),
		RTemaUsecase:    usecase.NewRoomTemaUseCase(db),
		RTemaRepository: rInfra.NewRoomTemaRepository(db),
	}
}

// List 一覧
func (r *RoomHandler) List(c echo.Context) error {
	hmUser, err := r.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &room.ListInput{}
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
		list, _ := r.RNeppanUsecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDTl:
		list, _ := r.RTlUsecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDDirect:
		list, _ := r.RDirectUsecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDRaku2:
		list, _ := r.RRaku2Usecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDTema:
		list, _ := r.RTemaUsecase.FetchList(request)
		return c.JSON(http.StatusOK, list)
	default:
		return echo.ErrInternalServerError
	}
}

// FetchAllAmenities アメニティ取得
func (r *RoomHandler) FetchAllAmenities(c echo.Context) error {
	hmUser, err := r.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	switch hmUser.WholesalerID {
	case utils.WholesalerIDNeppan:
		list, _ := r.RNeppanUsecase.FetchAllAmenities()
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDTl:
		list, _ := r.RTlUsecase.FetchAllAmenities()
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDDirect:
		list, _ := r.RDirectUsecase.FetchAllAmenities()
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDRaku2:
		list, _ := r.RRaku2Usecase.FetchAllAmenities()
		return c.JSON(http.StatusOK, list)
	case utils.WholesalerIDTema:
		list, _ := r.RTemaUsecase.FetchAllAmenities()
		return c.JSON(http.StatusOK, list)
	default:
		return echo.ErrInternalServerError
	}
}

// FetchAllRoomKinds 部屋種別一覧取得
func (r *RoomHandler) FetchAllRoomKinds(c echo.Context) error {
	hmUser, err := r.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	if hmUser.WholesalerID == utils.WholesalerIDTl {
		roomKinds, err := r.RTlRepository.FetchAllRoomKindTls()
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, roomKinds)
	}

	roomKinds, err := r.RCommonUsecase.FetchAllRoomKinds()
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, roomKinds)

}

// Detail 詳細
func (r *RoomHandler) Detail(c echo.Context) error {
	hmUser, err := r.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &room.DetailInput{}
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
		detail, err := r.RNeppanUsecase.FetchDetail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)
	case utils.WholesalerIDTl:
		detail, err := r.RTlUsecase.FetchDetail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)
	case utils.WholesalerIDTema:
		detail, err := r.RTemaUsecase.FetchDetail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)
	case utils.WholesalerIDDirect:
		detail, err := r.RDirectUsecase.FetchDetail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)

	case utils.WholesalerIDRaku2:
		detail, err := r.RRaku2Usecase.FetchDetail(request)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, detail)

	default:
		return echo.ErrInternalServerError
	}
}

// Create 作成
func (r *RoomHandler) Create(c echo.Context) error {
	hmUser, err := r.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &room.SaveInput{}
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
		if err := r.RNeppanUsecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	case utils.WholesalerIDDirect:
		if err := r.RDirectUsecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	case utils.WholesalerIDRaku2:
		if err := r.RRaku2Usecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}

	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// CreateOrUpdateBulk Bulk data for Room and Stock
func (r *RoomHandler) CreateOrUpdateBulk(c echo.Context) error {
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

	go r.processBulkRequest(wholesalerId, c)

	return c.JSON(http.StatusOK, map[string]string{"message": "Request accepted successfully!"})
}

func (r *RoomHandler) processBulkRequest(wholesalerId int, c echo.Context) {

	ProcessStartTime := time.Now()
	hostUrl := c.Request().Host
	ActivityLogID, _ := r.RLogRepository.StoreBulkActivityLog(utils.LogServiceRoom, utils.LogTypeMaster, hostUrl, time.Now())

	switch wholesalerId {
	case utils.WholesalerIDTl:
		if err := r.RTlUsecase.CreateOrUpdateBulk(requestDataTl); err != nil {
			c.Echo().Logger.Error(err)
			// Log the error and set the error message
			errorMessage := err.Error()
			err := r.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
			if err != nil {
				return
			}
		} else {
			// Task succeeded, set IsSuccess to true
			err := r.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, true, "")
			if err != nil {
				return
			}
		}
	case utils.WholesalerIDTema:
		if err := r.RTemaUsecase.CreateOrUpdateBulk(requestDataTema); err != nil {
			c.Echo().Logger.Error(err)
			// Log the error and set the error message
			errorMessage := err.Error()
			err := r.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
			if err != nil {
				return
			}
		} else {
			// Task succeeded, set IsSuccess to true
			err := r.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, true, "")
			if err != nil {
				return
			}
		}
	default:
		c.Echo().Logger.Error("Invalid wholesalerId")
		// Log the error for invalid wholesalerID
		errorMessage := "Invalid wholesalerID"
		err := r.RLogRepository.UpdateBulkActivityLog(ActivityLogID, ProcessStartTime, false, errorMessage)
		if err != nil {
			return
		}
	}
}

// Update 更新
func (r *RoomHandler) Update(c echo.Context) error {
	hmUser, err := r.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &room.SaveInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	if request.RoomTypeID == 0 {
		c.Echo().Logger.Error("room_type_id がありません。")
		return echo.ErrBadRequest
	}

	switch hmUser.WholesalerID {
	case utils.WholesalerIDNeppan:
		if err := r.RNeppanUsecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := r.RDirectUsecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDRaku2:
		if err := r.RRaku2Usecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// Delete 部屋削除
func (r *RoomHandler) Delete(c echo.Context) error {
	hmUser, err := r.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &room.DeleteInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	if request.RoomTypeID == 0 {
		c.Echo().Logger.Error("room_type_id がありません。")
		return echo.ErrBadRequest
	}

	switch hmUser.WholesalerID {
	case utils.WholesalerIDNeppan:
		if err := r.RNeppanUsecase.Delete(request.RoomTypeID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := r.RDirectUsecase.Delete(request.RoomTypeID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDRaku2:
		if err := r.RRaku2Usecase.Delete(request.RoomTypeID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// UpdateStopSales 部屋と紐づくプラン・在庫の売止更新
func (r *RoomHandler) UpdateStopSales(c echo.Context) error {
	hmUser, err := r.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &room.StopSalesInput{}
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
		if err := r.RNeppanUsecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := r.RDirectUsecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDRaku2:
		if err := r.RRaku2Usecase.UpdateStopSales(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// getHmUser トークンからHMアカウント情報を取得
func (r *RoomHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return r.AUsecase.FetchHMUserByToken(claimParam)
}
