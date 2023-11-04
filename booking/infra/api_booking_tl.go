package infra

import (
	"os"

	"github.com/Adventureinc/hotel-hm-api/src/booking"
	"github.com/Adventureinc/hotel-hm-api/src/common/infra"
)

// bookingTlAPI TL予約APIクライアント
type bookingTlAPI struct {
	client *infra.APIClient
}

var TlAPIPrefix = os.Getenv("TL_API_PREFIX")

// NewBookingTlAPI インスタンス生成
func NewBookingTlAPI() booking.IBookingTlAPI {
	return &bookingTlAPI{
		client: new(infra.APIClient),
	}
}

// RetrieveBooking 予約情報を取得するAPI
func (a *bookingTlAPI) RetrieveBooking(url string, body string) (booking.XmlRoomStay, error) {
	response := &booking.XmlEnvelope{}

	a.client.URL = url
	a.client.Data = body
	a.client.Response = response

	if err := a.client.PostXml([]byte(body)); err != nil {
		return booking.XmlRoomStay{}, err
	}

	return response.Body.XmlOTA_HotelResRS.HotelReservations.HotelReservation.RoomStays.RoomStay, nil
}
