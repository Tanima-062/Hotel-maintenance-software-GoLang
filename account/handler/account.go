package handler

import (
	"net/http"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// AccountHandler アカウント関連の振り分け
type AccountHandler struct {
	AUsecase       account.IAccountUsecase
	ANeppanUsecase account.IAccountNeppanUsecase
	ATemaUsecase   account.IAccountTemaUsecase
	ARaku2Usecase  account.IAccountRaku2Usecase
}

// NewAccountHandler インスタンス生成
func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{
		AUsecase:       usecase.NewAccountUsecase(db),
		ANeppanUsecase: usecase.NewAccountNeppanUsecase(db),
		ATemaUsecase:   usecase.NewAccountTemaUsecase(db),
		ARaku2Usecase:  usecase.NewAccountRaku2Usecase(db),
	}
}

// Login ログイン
func (a *AccountHandler) Login(c echo.Context) error {
	request := &account.LoginInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	token, err := a.AUsecase.Login(request)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrUnauthorized
	}

	return c.JSON(http.StatusOK, account.TokenOutput{APIToken: token})
}

// Logout ログアウト
func (a *AccountHandler) Logout(c echo.Context) error {
	claimParam, ClaimParamErr := utils.GetHmUser(c)
	if ClaimParamErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	a.AUsecase.Logout(claimParam)
	return c.NoContent(http.StatusOK)
}

// CheckToken トークン確認
// セキュリティ上、無効な問い合わせはすべてUnauthorizedで返す
func (a *AccountHandler) CheckToken(c echo.Context) error {
	claimParam, ClaimParamErr := utils.GetHmUser(c)
	if ClaimParamErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	token, err := a.AUsecase.CheckToken(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	return c.JSON(http.StatusOK, account.TokenOutput{APIToken: token})
}

// AccountDetail ログイン中のHMアカウント情報を取得
func (a *AccountHandler) AccountDetail(c echo.Context) error {
	claimParam, ClaimParamErr := utils.GetHmUser(c)
	if ClaimParamErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	fetchedHmUser, err := a.AUsecase.FetchDetail(claimParam)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusOK, fetchedHmUser)
}

// ChangePassword パスワード変更処理
func (a *AccountHandler) ChangePassword(c echo.Context) error {
	request := &account.ChangePasswordInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)
	if err := a.AUsecase.ChangePassword(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrUnauthorized
	}
	return c.NoContent(http.StatusOK)
}

// CheckConnect ホールセラー接続用のユーザがあるかどうか
func (a *AccountHandler) CheckConnect(c echo.Context) error {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	_, fErr := a.AUsecase.FetchHMUserByToken(claimParam)
	if fErr != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}
	request := &account.CheckConnectInput{}
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
	case utils.WholesalerIDTema:
		return c.JSON(http.StatusOK, a.ATemaUsecase.FetchConnectUser(request))
	case utils.WholesalerIDNeppan:
		return c.JSON(http.StatusOK, a.ANeppanUsecase.FetchConnectUser(request))
	case utils.WholesalerIDRaku2:
		return c.JSON(http.StatusOK, a.ARaku2Usecase.FetchConnectUser(request))
	}
	return echo.ErrInternalServerError
}

// IsParentAccount 親アカウントかどうか
func (a *AccountHandler) IsParentAccount(c echo.Context) error {
	request := &account.IsParentAccountInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}

	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	return c.JSON(http.StatusOK, a.AUsecase.IsParentAccount(request.HotelManagerID))
}
