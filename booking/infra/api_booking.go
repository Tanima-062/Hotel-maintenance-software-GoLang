package infra

import (
	"fmt"
	"os"

	"github.com/Adventureinc/hotel-hm-api/src/booking"
	"github.com/Adventureinc/hotel-hm-api/src/common/infra"
)

// bookingAPI TL 予約APIクライアント
type bookingAPI struct {
	client *infra.APIClient
}

// cancelResponse キャンセル時の返却値
type cancelResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewBookingAPI インスタンス生成
func NewBookingAPI() booking.IBookingAPI {
	return &bookingAPI{
		client: new(infra.APIClient),
	}
}

// CancelBooking 予約キャンセルAPI（現状、adminのキャンセル処理を実行するだけ）
func (a *bookingAPI) CancelBooking(cmApplicationID int64, cancelFee int64, noShow uint8) (bool, error) {
	response := &cancelResponse{}

	a.client.URL = fmt.Sprintf("%s/%s/%d/%d/%s?noshow=%d",
		os.Getenv("HOTEL_ADMIN_API_PREFIX"),
		"cancel_application_from_hm",
		cmApplicationID,
		cancelFee,
		os.Getenv("HM_API_KEY"),
		noShow)
	a.client.Response = response
	if err := a.client.Get(); err != nil {
		return false, err
	}

	return response.Status == 200, nil
}
