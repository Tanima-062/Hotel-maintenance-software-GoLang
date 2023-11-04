package booking

import (
	"encoding/xml"
)

// TL 予約情報取得APUのRoot
type XmlEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    XmlBody
}
type XmlBody struct {
	XMLName           xml.Name `xml:"Body"`
	XmlOTA_HotelResRS XmlOTA_HotelResRS
}
type XmlOTA_HotelResRS struct {
	XMLName           xml.Name `xml:"OTA_HotelResRS"`
	HotelReservations XmlHotelReservations
}

type XmlHotelReservations struct {
	XMLName          xml.Name `xml:"HotelReservations"`
	HotelReservation XmlHotelReservation
}

type XmlHotelReservation struct {
	XMLName   xml.Name `xml:"HotelReservation"`
	RoomStays XmlRoomStays
}

type XmlRoomStays struct {
	XMLName  xml.Name `xml:"RoomStays"`
	RoomStay XmlRoomStay
}

type XmlRoomStay struct {
	XMLName   xml.Name       `xml:"RoomStay"`
	RoomTypes []XmlRoomTypes `xml:"RoomTypes"`
	RatePlans XmlRatePlans   `xml:"RatePlans"`
}

type XmlRoomTypes struct {
	XMLName  xml.Name    `xml:"RoomTypes"`
	RoomType XmlRoomType `xml:"RoomType"`
}
type XmlRoomType struct {
	XMLName         xml.Name `xml:"RoomType"`
	RoomTypeCode    string   `xml:"RoomTypeCode,attr"`
	RoomDescription XmlRoomDescription
}
type XmlRoomDescription struct {
	XMLName xml.Name `xml:"RoomDescription"`
	Name    string   `xml:"Name,attr"`
}

type XmlRatePlans struct {
	XMLName  xml.Name `xml:"RatePlans"`
	RatePlan XmlRatePlan
}
type XmlRatePlan struct {
	XMLName      xml.Name `xml:"RatePlan"`
	RatePlanName string   `xml:"RatePlanName,attr"`
}

// IBookingTlAPI 予約関連のTL APIのインターフェース
type IBookingTlAPI interface {
	RetrieveBooking(url string, body string) (XmlRoomStay, error)
}
