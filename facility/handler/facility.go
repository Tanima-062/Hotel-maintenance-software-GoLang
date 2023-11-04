package handler

import (
	"net/http"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	"github.com/Adventureinc/hotel-hm-api/src/facility/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// FacilityHandler 施設関連の振り分け
type FacilityHandler struct {
	FParentUsecase facility.IParentUsecase
	FTlUsecase     facility.IFacilityUsecase
	FTemaUsecase   facility.IFacilityTemaUsecase
	FNeppanUsecase facility.IFacilityNeppanUsecase
	FDirectUsecase facility.IFacilityUsecase
	FRaku2Usecase  facility.IFacilityRaku2Usecase
	AUsecase       account.IAccountUsecase
}

// NewFacilityHandler インスタンス生成
func NewFacilityHandler(db *gorm.DB) *FacilityHandler {
	return &FacilityHandler{
		FParentUsecase: usecase.NewFacilityParentUsecase(db),
		FTlUsecase:     usecase.NewFacilityTlUsecase(db),
		FTemaUsecase:   usecase.NewFacilityTemaUsecase(db),
		FNeppanUsecase: usecase.NewFacilityNeppanUsecase(db),
		FDirectUsecase: usecase.NewFacilityDirectUsecase(db),
		FRaku2Usecase:  usecase.NewFacilityRaku2Usecase(db),
		AUsecase:       aUsecase.NewAccountUsecase(db),
	}
}

// FetchAll アカウントに紐づく施設一覧取得
func (f *FacilityHandler) FetchAll(c echo.Context) error {
	hmUser, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	switch hmUser.WholesalerID {
	case utils.WholesalerIDParent:
		facilities, err := f.FParentUsecase.FetchAll(hmUser)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDTema:
		facilities, err := f.FTemaUsecase.FetchAll(hmUser)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDTl:
		facilities, err := f.FTlUsecase.FetchAll(hmUser)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDNeppan:
		facilities, err := f.FNeppanUsecase.FetchAll(hmUser)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDDirect:
		facilities, err := f.FDirectUsecase.FetchAll(hmUser)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDRaku2:
		facilities, err := f.FRaku2Usecase.FetchAll(hmUser)
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	}
	return echo.ErrInternalServerError
}

// UpdateDispPriority 施設のサイト公開フラグを更新
func (f *FacilityHandler) UpdateDispPriority(c echo.Context) error {
	_, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &facility.UpdateDispPriorityInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	switch request.WholesalerID {
	case utils.WholesalerIDTl:
		if err := f.FTlUsecase.UpdateDispPriority(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDTema:
		if err := f.FTemaUsecase.UpdateDispPriority(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDNeppan:
		if err := f.FNeppanUsecase.UpdateDispPriority(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDDirect:
		if err := f.FDirectUsecase.UpdateDispPriority(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDRaku2:
		if err := f.FRaku2Usecase.UpdateDispPriority(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// FetchBaseInfo 施設の基本情報を取得
func (f *FacilityHandler) FetchBaseInfo(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	hmUser, err := f.AUsecase.FetchHMUserByToken(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &facility.BaseInfoInput{}
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
	case utils.WholesalerIDTl:
		res, _ := f.FTlUsecase.FetchBaseInfo(&hmUser, claimParam, request)
		return c.JSON(http.StatusOK, res)
	case utils.WholesalerIDTema:
		res, _ := f.FTemaUsecase.FetchBaseInfo(&hmUser, claimParam, request)
		return c.JSON(http.StatusOK, res)
	case utils.WholesalerIDNeppan:
		res, _ := f.FNeppanUsecase.FetchBaseInfo(&hmUser, claimParam, request)
		return c.JSON(http.StatusOK, res)
	case utils.WholesalerIDDirect:
		res, _ := f.FDirectUsecase.FetchBaseInfo(&hmUser, claimParam, request)
		return c.JSON(http.StatusOK, res)
	case utils.WholesalerIDRaku2:
		res, _ := f.FRaku2Usecase.FetchBaseInfo(&hmUser, claimParam, request)
		return c.JSON(http.StatusOK, res)
	}
	return echo.ErrInternalServerError
}

// FetchDetail 施設の詳細情報を取得
func (f *FacilityHandler) FetchDetail(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	hmUser, err := f.AUsecase.FetchHMUserByToken(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &facility.BaseInfoInput{}
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
	case utils.WholesalerIDTl:
		res, _ := f.FTlUsecase.FetchDetail(request)
		return c.JSON(http.StatusOK, res)
	case utils.WholesalerIDTema:
		res, _ := f.FTemaUsecase.FetchDetail(request)
		return c.JSON(http.StatusOK, res)
	case utils.WholesalerIDNeppan:
		res, _ := f.FNeppanUsecase.FetchDetail(request)
		return c.JSON(http.StatusOK, res)
	case utils.WholesalerIDDirect:
		res, _ := f.FDirectUsecase.FetchDetail(request)
		return c.JSON(http.StatusOK, res)
	case utils.WholesalerIDRaku2:
		res, _ := f.FRaku2Usecase.FetchDetail(request)
		return c.JSON(http.StatusOK, res)
	}
	return echo.ErrInternalServerError
}

// SaveBaseInfo 施設の基本情報を保存
func (f *FacilityHandler) SaveBaseInfo(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	hmUser, err := f.AUsecase.FetchHMUserByToken(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &facility.SaveBaseInfoInput{}
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
	case utils.WholesalerIDTl:
		if err := f.FTlUsecase.SaveBaseInfo(&hmUser, claimParam, request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDTema:
		// 連携ID重複登録チェック
		isRegistered, err := f.FTemaUsecase.IsRegisteredConnect(request);
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		if isRegistered {
			goto isRegisteredError
		}
		if err := f.FTemaUsecase.SaveBaseInfo(&hmUser, claimParam, request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDNeppan:
		// 連携ID重複登録チェック
		isRegistered, err := f.FNeppanUsecase.IsRegisteredConnect(request);
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		if isRegistered {
			goto isRegisteredError
		}
		if err := f.FNeppanUsecase.SaveBaseInfo(&hmUser, claimParam, request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDDirect:
		if err := f.FDirectUsecase.SaveBaseInfo(&hmUser, claimParam, request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDRaku2:
		// 連携ID重複登録チェック
		isRegistered, err := f.FRaku2Usecase.IsRegisteredConnect(request);
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		if isRegistered {
			goto isRegisteredError
		}
		if err := f.FRaku2Usecase.SaveBaseInfo(&hmUser, claimParam, request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	}
	return c.NoContent(http.StatusOK)

	// ラベル 連動ID重複登録エラー
	isRegisteredError:
		return echo.NewHTTPError(http.StatusConflict, "connect id is already registered.")

}

// SaveDetail 施設の詳細情報の保存
func (f *FacilityHandler) SaveDetail(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	hmUser, err := f.AUsecase.FetchHMUserByToken(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &facility.SaveDetailInput{}
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
	case utils.WholesalerIDTl:
		if err := f.FTlUsecase.SaveDetail(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDTema:
		if err := f.FTemaUsecase.SaveDetail(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDNeppan:
		if err := f.FNeppanUsecase.SaveDetail(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDDirect:
		if err := f.FDirectUsecase.SaveDetail(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDRaku2:
		if err := f.FRaku2Usecase.SaveDetail(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	}
	return c.NoContent(http.StatusOK)
}

// FetchAllAmenities 施設のアメニティ一覧を取得
func (f *FacilityHandler) FetchAllAmenities(c echo.Context) error {
	hmUser, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	switch hmUser.WholesalerID {
	case utils.WholesalerIDTema:
		facilities, err := f.FTlUsecase.FetchAllAmenities()
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDTl:
		facilities, err := f.FTlUsecase.FetchAllAmenities()
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDNeppan:
		facilities, err := f.FNeppanUsecase.FetchAllAmenities()
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDDirect:
		facilities, err := f.FDirectUsecase.FetchAllAmenities()
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	case utils.WholesalerIDRaku2:
		facilities, err := f.FRaku2Usecase.FetchAllAmenities()
		if err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
		return c.JSON(http.StatusOK, facilities)
	}
	return echo.ErrInternalServerError
}

// getHmUser トークンからHMアカウント情報を取得
func (f *FacilityHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return f.AUsecase.FetchHMUserByToken(claimParam)
}
