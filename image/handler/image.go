package handler

import (
	"net/http"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aUsecase "github.com/Adventureinc/hotel-hm-api/src/account/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"github.com/Adventureinc/hotel-hm-api/src/image/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ImageHandler 画像関連の振り分け
type ImageHandler struct {
	ITlUsecase     image.IImageUsecase
	INeppanUsecase image.IImageUsecase
	IDirectUsecase image.IImageUsecase
	IRaku2Usecase  image.IImageUsecase
	AUsecase       account.IAccountUsecase
}

// NewImageHandler インスタンス生成
func NewImageHandler(db *gorm.DB) *ImageHandler {
	return &ImageHandler{
		ITlUsecase:     usecase.NewImageTlUsecase(db),
		INeppanUsecase: usecase.NewImageNeppanUsecase(db),
		IDirectUsecase: usecase.NewImageDirectUsecase(db),
		IRaku2Usecase:  usecase.NewImageRaku2Usecase(db),
		AUsecase:       aUsecase.NewAccountUsecase(db),
	}
}

// FetchAll 画像一覧
func (i *ImageHandler) FetchAll(c echo.Context) error {
	hmUser, err := i.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &image.ListInput{}
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
	case utils.WholesalerIDTema:
		fallthrough
	case utils.WholesalerIDTl:
		images, _ := i.ITlUsecase.FetchAll(request)
		return c.JSON(http.StatusOK, images)

	case utils.WholesalerIDNeppan:
		images, _ := i.INeppanUsecase.FetchAll(request)
		return c.JSON(http.StatusOK, images)

	case utils.WholesalerIDDirect:
		images, _ := i.IDirectUsecase.FetchAll(request)
		return c.JSON(http.StatusOK, images)

	case utils.WholesalerIDRaku2:
		images, _ := i.IRaku2Usecase.FetchAll(request)
		return c.JSON(http.StatusOK, images)
	}

	return echo.ErrInternalServerError
}

// Update 画像更新
func (i *ImageHandler) Update(c echo.Context) error {
	hmUser, err := i.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &image.UpdateInput{}
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
	case utils.WholesalerIDTema:
		fallthrough
	case utils.WholesalerIDTl:
		if err := i.ITlUsecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDNeppan:
		if err := i.INeppanUsecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := i.IDirectUsecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDRaku2:
		if err := i.IRaku2Usecase.Update(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	}

	return c.NoContent(http.StatusOK)
}

// UpdateIsMain メイン画像更新
func (i *ImageHandler) UpdateIsMain(c echo.Context) error {
	hmUser, err := i.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &image.UpdateIsMainInput{}
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
	case utils.WholesalerIDTema:
		fallthrough
	case utils.WholesalerIDTl:
		if err := i.ITlUsecase.UpdateIsMain(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDNeppan:
		if err := i.INeppanUsecase.UpdateIsMain(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := i.IDirectUsecase.UpdateIsMain(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDRaku2:
		if err := i.IRaku2Usecase.UpdateIsMain(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	}

	return c.NoContent(http.StatusOK)
}

// UpdateSortNum メイン画像の並び順更新
func (i *ImageHandler) UpdateSortNum(c echo.Context) error {
	hmUser, err := i.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &[]image.UpdateSortNumInput{}
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
	case utils.WholesalerIDTema:
		fallthrough
	case utils.WholesalerIDTl:
		if err := i.ITlUsecase.UpdateSortNum(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDNeppan:
		if err := i.INeppanUsecase.UpdateSortNum(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := i.IDirectUsecase.UpdateSortNum(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDRaku2:
		if err := i.IRaku2Usecase.UpdateSortNum(request); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	}

	return c.NoContent(http.StatusOK)
}

// Delete 画像の削除
func (i *ImageHandler) Delete(c echo.Context) error {
	hmUser, err := i.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &image.DeleteInput{}
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
	case utils.WholesalerIDTema:
		fallthrough
	case utils.WholesalerIDTl:
		if err := i.ITlUsecase.Delete(request.ImageID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDNeppan:
		if err := i.INeppanUsecase.Delete(request.ImageID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := i.IDirectUsecase.Delete(request.ImageID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDRaku2:
		if err := i.IRaku2Usecase.Delete(request.ImageID); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	}

	return c.NoContent(http.StatusOK)
}

// Create メイン画像の作成
func (i *ImageHandler) Create(c echo.Context) error {
	hmUser, err := i.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	file, fileErr := c.FormFile("imagefile")
	if fileErr != nil {
		c.Echo().Logger.Error(fileErr)
		return echo.ErrBadRequest
	}

	request := &image.UploadInput{}
	if err := c.Bind(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	if err := c.Validate(request); err != nil {
		c.Echo().Logger.Error(err)
		return echo.ErrBadRequest
	}
	utils.RequestLog(c, request)

	request.Caption = utils.RemoveDoubleQuotation(request.Caption)
	request.ContentType = utils.RemoveDoubleQuotation(request.ContentType)
	request.CategoryCd = utils.RemoveDoubleQuotation(request.CategoryCd)

	switch hmUser.WholesalerID {
	case utils.WholesalerIDTema:
		fallthrough
	case utils.WholesalerIDTl:
		if err := i.ITlUsecase.Create(request, file, hmUser); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDNeppan:
		if err := i.INeppanUsecase.Create(request, file, hmUser); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDDirect:
		if err := i.IDirectUsecase.Create(request, file, hmUser); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}

	case utils.WholesalerIDRaku2:
		if err := i.IRaku2Usecase.Create(request, file, hmUser); err != nil {
			c.Echo().Logger.Error(err)
			return echo.ErrInternalServerError
		}
	}

	return c.NoContent(http.StatusOK)
}

// CountMainImages メイン画像設定数取得
func (i *ImageHandler) CountMainImages(c echo.Context) error {
	_, err := i.getHmUser(c)
	if err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired jwt")
	}

	request := &image.MainImagesCountInput{}
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
		fallthrough
	case utils.WholesalerIDTl:
		return c.JSON(http.StatusOK, map[string]int64{"count": i.ITlUsecase.CountMainImages(request.PropertyID)})

	case utils.WholesalerIDNeppan:
		return c.JSON(http.StatusOK, map[string]int64{"count": i.INeppanUsecase.CountMainImages(request.PropertyID)})

	case utils.WholesalerIDDirect:
		return c.JSON(http.StatusOK, map[string]int64{"count": i.IDirectUsecase.CountMainImages(request.PropertyID)})

	case utils.WholesalerIDRaku2:
		return c.JSON(http.StatusOK, map[string]int64{"count": i.IRaku2Usecase.CountMainImages(request.PropertyID)})
	}

	return echo.ErrInternalServerError
}

// getHmUser トークンからHMアカウント情報を取得
func (i *ImageHandler) getHmUser(c echo.Context) (account.HtTmHotelManager, error) {
	claimParam, err := utils.GetHmUser(c)
	if err != nil {
		return account.HtTmHotelManager{}, err
	}
	return i.AUsecase.FetchHMUserByToken(claimParam)
}
