package handler

import (
	"net/http"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CancelPolicyHandler キャンセルポリシー関連の振り分け
type CancelPolicyHandler struct {
	FTlUsecase     cancelPolicy.ICancelPolicyUsecase
	FNeppanUsecase cancelPolicy.ICancelPolicyUsecase
	FDirectUsecase cancelPolicy.ICancelPolicyUsecase
	FRaku2Usecase  cancelPolicy.ICancelPolicyUsecase
	AUsecase       account.IAccountUsecase
}

// NewCancelPolicyHandler インスタンス生成
func NewCancelPolicyHandler(db *gorm.DB) *CancelPolicyHandler {
	return &CancelPolicyHandler{
		FTlUsecase:     usecase.NewCancelPolicyTlUsecase(db),
		FNeppanUsecase: usecase.NewCancelPolicyNeppanUsecase(db),
		FDirectUsecase: usecase.NewCancelPolicyDirectUsecase(db),
		FRaku2Usecase:  usecase.NewCancelPolicyRaku2Usecase(db),
		AUsecase:       aUsecase.NewAccountUsecase(db),
	}
}

func (f *CancelPolicyHandler) List(c echo.Context) error {
	hmUser, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &cancelPolicy.ListInput{}
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
		cancelPolicyList, _ := f.FTlUsecase.List(request)
		return c.JSON(http.StatusOK, cancelPolicyList)
	case utils.WholesalerIDNeppan:
		cancelPolicyList, _ := f.FNeppanUsecase.List(request)
		return c.JSON(http.StatusOK, cancelPolicyList)
	case utils.WholesalerIDDirect:
		cancelPolicyList, _ := f.FDirectUsecase.List(request)
		return c.JSON(http.StatusOK, cancelPolicyList)
	case utils.WholesalerIDRaku2:
		cancelPolicyList, _ := f.FRaku2Usecase.List(request)
		return c.JSON(http.StatusOK, cancelPolicyList)
	}
	return echo.ErrInternalServerError
}

// Create はキャンセルポリシーを新規作成します
func (f *CancelPolicyHandler) Create(c echo.Context) error {
	hmUser, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &cancelPolicy.CreateInput{}
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
		if err := f.FTlUsecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrBadRequest
		}
	case utils.WholesalerIDNeppan:
		if err := f.FNeppanUsecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDDirect:
		if err := f.FDirectUsecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDRaku2:
		if err := f.FRaku2Usecase.Create(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	}

	return c.NoContent(http.StatusOK)
}

// Detail キャンセルポリシー詳細情報取得
func (f *CancelPolicyHandler) Detail(c echo.Context) error {
	hmUser, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &cancelPolicy.DetailInput{}
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
		cancelPolicy, err := f.FTlUsecase.Detail(request)
		if err == nil {
			return c.JSON(http.StatusOK, cancelPolicy)
		}
	case utils.WholesalerIDNeppan:
		cancelPolicy, err := f.FNeppanUsecase.Detail(request)
		if err == nil {
			return c.JSON(http.StatusOK, cancelPolicy)
		}
	case utils.WholesalerIDDirect:
		cancelPolicy, err := f.FDirectUsecase.Detail(request)
		if err == nil {
			return c.JSON(http.StatusOK, cancelPolicy)
		}
	case utils.WholesalerIDRaku2:
		cancelPolicy, err := f.FRaku2Usecase.Detail(request)
		if err == nil {
			return c.JSON(http.StatusOK, cancelPolicy)
		}
	}

	return echo.ErrInternalServerError
}

// Save キャンセルポリシーの保存
func (f *CancelPolicyHandler) Save(c echo.Context) error {
	hmUser, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &cancelPolicy.UpdateInput{}
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
		if err := f.FTlUsecase.Save(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDNeppan:
		if err := f.FNeppanUsecase.Save(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDDirect:
		if err := f.FDirectUsecase.Save(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDRaku2:
		if err := f.FRaku2Usecase.Save(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	default:
		return echo.ErrInternalServerError
	}
	return c.NoContent(http.StatusOK)
}

// Save キャンセルポリシーの削除
func (f *CancelPolicyHandler) Delete(c echo.Context) error {
	hmUser, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &cancelPolicy.DeleteInput{}
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
		if err := f.FTlUsecase.Delete(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDNeppan:
		if err := f.FNeppanUsecase.Delete(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDDirect:
		if err := f.FDirectUsecase.Delete(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	case utils.WholesalerIDRaku2:
		if err := f.FRaku2Usecase.Delete(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	default:
		return echo.ErrInternalServerError
	}

	return c.NoContent(http.StatusOK)
}

// キャンセルポリシーに紐付くプラン一覧を取得
func (f *CancelPolicyHandler) FetchPlans(c echo.Context) error {
	hmUser, err := f.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &cancelPolicy.PlanListInput{}
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
		planList, _ := f.FTlUsecase.PlanList(request)
		return c.JSON(http.StatusOK, planList)
	case utils.WholesalerIDNeppan:
		planList, _ := f.FNeppanUsecase.PlanList(request)
		return c.JSON(http.StatusOK, planList)
	case utils.WholesalerIDDirect:
		planList, _ := f.FDirectUsecase.PlanList(request)
		return c.JSON(http.StatusOK, planList)
	case utils.WholesalerIDRaku2:
		planList, _ := f.FRaku2Usecase.PlanList(request)
		return c.JSON(http.StatusOK, planList)
	}
	return echo.ErrInternalServerError
}

// getHmUser トークンからHMアカウント情報を取得
func (f *CancelPolicyHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return f.AUsecase.FetchHMUserByToken(claimParam)
}
