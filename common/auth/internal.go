package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/labstack/echo/v4"
)

func Internal(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// check `API_KEY` exist in request header
		apiKey := c.Request().Header.Get(os.Getenv("ADV_INTERNAL_API_KEY_HEADER"))
		if apiKey != os.Getenv("ADV_INTERNAL_API_KEY") {
			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Invalid %s", os.Getenv("ADV_INTERNAL_API_KEY_HEADER")))
		}
		// check `Wholesaler-Id` exist in header
		wholesalerId := c.Request().Header.Get("Wholesaler-Id")
		if c.Request().Header.Get("Wholesaler-Id") == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Unknown Wholesaler-Id")
		} else {
			wid, _ := strconv.Atoi(wholesalerId)
			if !containsInValidWholesalerList(wid) {
				return echo.NewHTTPError(http.StatusBadRequest, "Unknown Wholesaler-Id")
			}
		}

		return next(c)
	}
}

func containsInValidWholesalerList(wholesalerId int) bool {
	validWholesalers := []int{
		utils.WholesalerIDParent,
		utils.WholesalerIDTl,
		utils.WholesalerIDTema,
		utils.WholesalerIDNeppan,
		utils.WholesalerIDDirect,
		utils.WholesalerIDRaku2,
	}

	for i := 0; i < len(validWholesalers); i++ {
		if wholesalerId == validWholesalers[i] {
			return true
		}
	}

	return false
}
