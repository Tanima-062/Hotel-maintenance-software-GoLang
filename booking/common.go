package booking

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

// HtThApplications 予約情報テーブル
type HtThApplications struct {
	HtThApplicationID        int64     `json:"ht_th_application_id"`
	WholesalerID             int64     `json:"wholesaler_id"`
	PropertyID               int64     `json:"property_id"`
	CmApplicationID          int64     `json:"cm_application_id"`
	CustomerIP               string    `json:"customer_ip"`
	CustomerSessID           string    `json:"customer_sess_id"`
	ItineraryID              string    `json:"itinerary_id"`
	AffiliateReferenceID     string    `json:"affiliate_reference_id"`
	AffiliateMetadata        string    `json:"affiliate_metadata"`
	RetrieveLink             string    `json:"retrieve_link"`
	GivenNameEnc             string    `json:"given_name_enc"`
	FamilyNameEnc            string    `json:"family_name_enc"`
	EmailEnc                 string    `json:"email_enc"`
	PhoneEnc                 string    `json:"phone_enc"`
	Line1Enc                 string    `json:"line_1_enc"`
	Line2Enc                 string    `json:"line_2_enc"`
	Line3Enc                 string    `json:"line_3_enc"`
	CityEnc                  string    `json:"city_enc"`
	StateProvinceCodeEnc     string    `json:"state_providence_code_enc"`
	PostalCodeEnc            string    `json:"postal_code_enc"`
	PaymentMethod            string    `json:"payment_method"`
	PaymentType              string    `json:"payment_type"`
	PaymentMeans             string    `json:"payment_means"`
	Arrival                  string    `json:"arrival"`
	Departure                string    `json:"departure"`
	Stays                    int       `json:"stays"`
	RoomNum                  int       `json:"room_num"`
	Hold                     bool      `json:"hold"`
	Holded                   time.Time `gorm:"type:time" json:"holded"`
	Status                   string    `json:"status"`
	CancelPolicyDisp         string    `json:"cancel_policy_disp"`
	CancelPolicyEnDisp       string    `json:"cancel_policy_en_disp"`
	CancelFlg                bool      `json:"cancel_flg"`
	CancelFee                float32   `json:"cancel_fee"`
	NoshowFlg                bool      `json:"noshow_flg"`
	NoshowFee                float32   `json:"noshow_fee"`
	CanceledDt               time.Time `gorm:"type:time" json:"canceled_dt"`
	CancelError              string    `json:"cancel_error"`
	CommisionFee             float32   `json:"commision_fee"`
	CommisionFeeTax          float32   `json:"commision_fee_tax"`
	HandlingCharge           float32   `json:"handling_charge"`
	AdministrativeFee        float32   `json:"administrative_fee"`
	TotalPayExTax            float32   `json:"total_pay_ex_tax"`
	TotalPayInTax            float32   `json:"total_pay_in_tax"`
	TotalStrikethrough       float32   `json:"total_strikethrough"`
	TotalStrikethroughInTax  float32   `json:"total_strikethrough_in_tax"`
	FeesJSON                 string    `json:"fees_json"`
	NightlyJSON              string    `json:"nightly_json"`
	StayJSON                 string    `json:"stay_json"`
	DescriptionsJSON         string    `json:"descriptions_json"`
	DescriptionsEnJSON       string    `json:"descriptions_en_json"`
	OptionalFeesJSON         string    `json:"optional_fees_json"`
	OptionalFeesEnJSON       string    `json:"optional_fees_en_json"`
	PromotionsJSON           string    `json:"promotions_json"`
	CurrencyID               string    `json:"currency_id"`
	LangCd                   string    `json:"lang_cd"`
	UserID                   int64     `json:"user_id"`
	SmsSendFlg               bool      `json:"sms_send_flg"`
	CompletePropertyImg      string    `json:"complete_property_img"`
	CompletePropertyLat      float64   `json:"complete_property_lat"`
	CompletePropertyLng      float64   `json:"complete_property_lng"`
	CompletePropertyArrives  time.Time `gorm:"type:time" json:"complete_property_arrives"`
	CompletePropertyRoomName string    `json:"complete_property_room_name"`
	AirCmApplicationID       int64     `json:"air_cm_application_id"`
	SupplierInfo             string    `json:"supplier_info"`
	IsMisuse                 int       `json:"is_misuse"`
	ContactStatus            int       `json:"contact_status"`
	Memo                     string    `json:"memo"`
	SearchEmailEnc           string    `json:"search_email_enc"`
	SearchGivenNameEnc       string    `json:"search_given_name_enc"`
	SearchFamilyNameEnc      string    `json:"search_family_name_enc"`
	CreateOperatorID         int64     `json:"create_operator_id"`
	UpdateOperatorID         int64     `json:"update_operator_id"`
	common.Times             `gorm:"embedded"`
}

// HtThBookingPrices 予約単価情報テーブル
type HtThBookingPrices struct {
	HtThBookingPriceID int64     `json:"ht_th_booking_price_id"`
	CmApplicationID    int64     `json:"cm_application_id"`
	UseDate            time.Time `json:"use_date"`
	RoomTypeID         int64     `json:"room_type_id"`
	PlanID             int64     `json:"plan_id"`
	Person             int       `json:"person"`
	Child1Person       int       `json:"child_1_person" gorm:"column:child_1_person"`
	Child2Person       int       `json:"child_2_person" gorm:"column:child_2_person"`
	Child3Person       int       `json:"child_3_person" gorm:"column:child_3_person"`
	Child4Person       int       `json:"child_4_person" gorm:"column:child_4_person"`
	Child5Person       int       `json:"child_5_person" gorm:"column:child_5_person"`
	Child6Person       int       `json:"child_6_person" gorm:"column:child_6_person"`
	Price              int       `json:"price"`
	PriceInTax         int       `json:"price_in_tax"`
	ChildPrice1        int       `json:"child_price1" gorm:"column:child_price1"`
	ChildPrice1InTax   int       `json:"child_price1_in_tax" gorm:"column:child_price1_in_tax"`
	ChildPrice2        int       `json:"child_price2" gorm:"column:child_price2"`
	ChildPrice2InTax   int       `json:"child_price2_in_tax" gorm:"column:child_price2_in_tax"`
	ChildPrice3        int       `json:"child_price3" gorm:"column:child_price3"`
	ChildPrice3InTax   int       `json:"child_price3_in_tax" gorm:"column:child_price3_in_tax"`
	ChildPrice4        int       `json:"child_price4" gorm:"column:child_price4"`
	ChildPrice4InTax   int       `json:"child_price4_in_tax" gorm:"column:child_price4_in_tax"`
	ChildPrice5        int       `json:"child_price5" gorm:"column:child_price5"`
	ChildPrice5InTax   int       `json:"child_price5_in_tax" gorm:"column:child_price5_in_tax"`
	ChildPrice6        int       `json:"child_price6" gorm:"column:child_price6"`
	ChildPrice6InTax   int       `json:"child_price6_in_tax" gorm:"column:child_price6_in_tax"`
	common.Times       `gorm:"embedded"`
}

// HtThBookingRooms 予約時の部屋情報テーブル
type HtThBookingRooms struct {
	HtThBookingRoomID                 int64  `json:"ht_th_booking_room_id"`
	HtThApplicationID                 int64  `json:"ht_th_application_id"`
	FamilyNameEnc                     string `json:"family_name_enc"`
	GivenNameEnc                      string `json:"given_name_enc"`
	Refundable                        bool   `json:"refundable"`
	NumberOfChilds                    int    `json:"number_of_childs"`
	NumberOfAdults                    int    `json:"number_of_adults"`
	ChildAges                         string `json:"child_ages"`
	NumberOfUpperGrades               int    `json:"number_of_upper_grades"`
	NumberOfLowerGrades               int    `json:"number_of_lower_grades"`
	NumberOfInfantMealsWithBedding    int    `json:"number_of_infant_meals_with_bedding"`
	NumberOfInfantMealOnly            int    `json:"number_of_infant_meal_only"`
	NumberOfInfantBeddingOnly         int    `json:"number_of_infant_bedding_only"`
	NumberOfInfantMealsWithoutBedding int    `json:"number_of_infant_meals_without_bedding"`
	CancelPenalties                   string `json:"cancel_penalties"`
	RoomID                            string `json:"room_id"`
	RateID                            string `json:"rate_id"`
	RoomName                          string `json:"room_name"`
	PlanName                          string `json:"plan_name"`
}

// HtTmItineraryTls TLの部屋情報・プラン名のテーブル
type HtTmItineraryTls struct {
	ItineraryID  string `json:"itinerary_id"`
	TlPropertyID string `json:"tl_property_id"`
	RoomID       string `json:"room_id"`
	RoomName     string `json:"room_name"`
	PlanName     string `json:"plan_name"`
}

// SearchInput 予約検索取得の入力
type SearchInput struct {
	PropertyID        int64    `json:"property_id" validate:"required"`
	WholesalerID      int64    `json:"wholesaler_id" validate:"required"`
	ApplicationIDs    []int64  `json:"application_ids"`
	ApplicationStart  string   `json:"application_start"`
	ApplicationEnd    string   `json:"application_end"`
	CheckinStart      string   `json:"checkin_start"`
	CheckinEnd        string   `json:"checkin_end"`
	CheckoutStart     string   `json:"checkout_start"`
	CheckoutEnd       string   `json:"checkout_end"`
	FamilyName        string   `json:"family_name"`          /*暗号化前の予約者性*/
	GivenName         string   `json:"given_name"`           /*暗号化前の予約者名*/
	FamilyNameEncList []string `json:"family_name_enc_list"` /*暗号化した予約者性*/
	GivenNameEncList  []string `json:"given_name_enc_list"`  /*暗号化した予約者名*/
	Phone             string   `json:"phone"`                /*暗号化前の電話番号*/
	PhoneEnc          string   `json:"phone_enc"`            /*暗号化した電話番号*/
	Status            uint8    `json:"status"`
}

// SearchDBOutput 予約検索のDB出力
type SearchDBOutput struct {
	CmApplicationID int64   `json:"cm_application_id"`
	CancelFlg       bool    `json:"cancel_flg"`
	NoshowFlg       bool    `json:"noshow_flg"`
	ApplicationCd   string  `json:"application_cd"`
	TourID          int64   `json:"tour_id"`
	Arrival         string  `json:"arrival"`
	Departure       string  `json:"departure"`
	GivenNameEnc    string  `json:"given_name_enc"`
	FamilyNameEnc   string  `json:"family_name_enc"`
	PhoneEnc        string  `json:"phone_enc"`
	TotalPayInTax   float32 `json:"total_pay_in_tax"`
	PaymentLimitDt  string  `json:"payment_limit_dt"`
	PaymentCount    int     `json:"payment_count"`
}

// SearchOutput 予約検索の出力
type SearchOutput struct {
	CmApplicationID int64   `json:"cm_application_id"`
	Status          uint8   `json:"status"`
	ApplicationCd   string  `json:"application_cd"`
	TourID          int64   `json:"tour_id"`
	Checkin         string  `json:"checkin"`
	Checkout        string  `json:"checkout"`
	GivenName       string  `json:"given_name"`
	FamilyName      string  `json:"family_name"`
	Phone           string  `json:"phone"`
	TotalPayInTax   float32 `json:"total_pay_in_tax"`
}

// DownloadInput CSV変換・DLする予約情報取得の入力
type DownloadInput struct {
	CmApplicationIDs []int64 `json:"cm_application_ids" validate:"required"`
	PropertyID       int64   `json:"property_id" validate:"required"`
	WholesalerID     int64   `json:"wholesaler_id" validate:"required"`
}

// BookingDownloadDBOutput CSV変換・DLする予約一覧の詳細情報のDB出力
type BookingDownloadDBOutput struct {
	CmApplicationID             int64     `json:"cm_application_id"`
	HtThApplicationID           int64     `json:"ht_th_application_id"`
	WholesalerID                int64     `json:"wholesaler_id"`
	PropertyID                  int64     `json:"property_id"`
	ItineraryID                 string    `json:"itinerary_id"`
	CancelFlg                   bool      `json:"cancel_flg"`
	CancelFee                   float32   `json:"cancel_fee"`
	CanceledDt                  time.Time `gorm:"type:time" json:"canceled_dt"`
	NoshowFlg                   bool      `json:"noshow_flg"`
	NoshowFee                   float32   `json:"noshow_fee"`
	ApplicationCd               string    `json:"application_cd"`
	TourID                      int64     `json:"tour_id"`
	Arrival                     string    `json:"arrival"`
	Departure                   string    `json:"departure"`
	Stays                       string    `json:"stays"`
	RoomNum                     string    `json:"room_num"`
	GivenNameEnc                string    `json:"given_name_enc"`
	FamilyNameEnc               string    `json:"family_name_enc"`
	PhoneEnc                    string    `json:"phone_enc"`
	EmailEnc                    string    `json:"email_enc"`
	TotalPayInTax               float32   `json:"total_pay_in_tax"`
	SalePrice                   float32   `json:"sale_price"`
	DiscountPaymentFlg          bool      `json:"discount_payment_flg"`
	DiscountCashAmount          float32   `json:"discount_cash_amount"`
	PaymentLimitDt              string    `json:"payment_limit_dt"`
	PaymentCount                int       `json:"payment_count"`
	CreatedAt                   time.Time `json:"created_at"`
}

// BookingDownloadOutput CSV変換・DLする予約一覧の詳細情報の出力
type BookingDownloadOutput struct {
	CmApplicationID             int64               `json:"cm_application_id"`
	ApplicationCd               string              `json:"application_cd"`
	TourID                      int64               `json:"tour_id"`
	CreatedAt                   time.Time           `json:"created_at"`
	GivenNameEnc                string              `json:"given_name_enc"`
	FamilyNameEnc               string              `json:"family_name_enc"`
	EmailEnc                    string              `json:"email_enc"`
	TotalPayInTax               float32             `json:"total_pay_in_tax"`
	SalePrice                   float32             `json:"sale_price"`
	CancelFee                   float32             `json:"cancel_fee"`
	CancelFlg                   bool                `json:"cancel_flg"`
	CanceledDt                  time.Time           `gorm:"type:time" json:"canceled_dt"`
	NoshowFlg                   bool                `json:"noshow_flg"`
	NoshowFee                   float32             `json:"noshow_fee"`
	Arrival                     string              `json:"arrival"`
	Departure                   string              `json:"departure"`
	Stays                       string              `json:"stays"`
	RoomNum                     string              `json:"room_num"`
	PhoneEnc                    string              `json:"phone_enc"`
	DiscountPaymentFlg          bool                `json:"discount_payment_flg"`
	DiscountCashAmount          float32             `json:"discount_cash_amount"`
	Status                      uint8               `json:"status"`
	RoomsAndPlan                []DetailRoomAndPlan `json:"rooms"`
}

// DetailInput 予約詳細情報取得の入力
type DetailInput struct {
	CmApplicationID int64 `json:"cm_application_id" param:"cmApplicationId" validate:"required"`
	PropertyID      int64 `json:"property_id" param:"propertyId" validate:"required"`
	WholesalerID    int64 `json:"wholesaler_id" query:"wholesaler_id" validate:"required"`
}

// DetailApplicationDBOutput 予約詳細情報のDB出力
type DetailApplicationDBOutput struct {
	CmApplicationID             int64     `json:"cm_application_id"`
	HtThApplicationID           int64     `json:"ht_th_application_id"`
	WholesalerID                int64     `json:"wholesaler_id"`
	PropertyID                  int64     `json:"property_id"`
	ItineraryID                 string    `json:"itinerary_id"`
	CancelFlg                   bool      `json:"cancel_flg"`
	CancelFee                   float32   `json:"cancel_fee"`
	CanceledDt                  time.Time `gorm:"type:time" json:"canceled_dt"`
	NoshowFlg                   bool      `json:"noshow_flg"`
	NoshowFee                   float32   `json:"noshow_fee"`
	ApplicationCd               string    `json:"application_cd"`
	TourID                      int64     `json:"tour_id"`
	Arrival                     string    `json:"arrival"`
	Departure                   string    `json:"departure"`
	Stays                       string    `json:"stays"`
	RoomNum                     string    `json:"room_num"`
	GivenNameEnc                string    `json:"given_name_enc"`
	FamilyNameEnc               string    `json:"family_name_enc"`
	PhoneEnc                    string    `json:"phone_enc"`
	EmailEnc                    string    `json:"email_enc"`
	TotalPayInTax               float32   `json:"total_pay_in_tax"`
	PaymentLimitDt              string    `json:"payment_limit_dt"`
	PaymentCount                int       `json:"payment_count"`
	CreatedAt                   time.Time `json:"created_at"`
}

// DetailOutput 予約詳細情報の出力
type DetailOutput struct {
	CmApplicationID             int64                  `json:"cm_application_id"`
	WholesalerID                int64                  `json:"wholesaler_id"`
	ApplicationCd               string                 `json:"application_cd"`
	TourID                      int64                  `json:"tour_id"`
	CreatedAt                   time.Time              `json:"created_at"`
	GivenNameEnc                string                 `json:"given_name_enc"`
	FamilyNameEnc               string                 `json:"family_name_enc"`
	EmailEnc                    string                 `json:"email_enc"`
	TotalPayInTax               float32                `json:"total_pay_in_tax"`
	SalePrice                   float32                `json:"sale_price"`
	CancelFee                   float32                `json:"cancel_fee"`
	CancelFeeSuggest            float32                `json:"cancel_fee_suggest"`
	CancelFlg                   bool                   `json:"cancel_flg"`
	CanceledDt                  time.Time              `gorm:"type:time" json:"canceled_dt"`
	NoshowFlg                   bool                   `json:"noshow_flg"`
	NoshowFee                   float32                `json:"noshow_fee"`
	Arrival                     string                 `json:"arrival"`
	Departure                   string                 `json:"departure"`
	Stays                       string                 `json:"stays"`
	RoomNum                     string                 `json:"room_num"`
	PhoneEnc                    string                 `json:"phone_enc"`
	DiscountPaymentFlg          bool                   `json:"discount_payment_flg"`
	DiscountCashAmount          float32                `json:"discount_cash_amount"`
	Status                      uint8                  `json:"status"`
	RoomsAndPlan                []DetailRoomAndPlan    `json:"rooms"`
	FlashSales                  []FlashSale            `json:"flash_sales"`
	PersonPrices                map[string][]PersonPrice `json:"person_prices"`
}

// CmThFlashSale セール情報のDB出力
type CmThFlashSale struct {
	CmApplicationID             int64      `json:"cm_application_id"`
	SaleType                    string     `json:"sale_type"`
	SalePrice                   float32    `json:"sale_price"`
	DiscountPaymentFlg          bool       `json:"discount_payment_flg"`
	DiscountCashAmount          float32    `json:"discount_cash_amount"`
	DiscountCouponAmount        float32    `json:"discount_coupon_amount"`
}

// FlashSale セール情報の出力
type FlashSale struct {
	SaleName                    string        `json:"sale_name"`
	DiscountCashAmount          float32       `json:"discount_cash_amount"`
	DiscountCouponCount         int           `json:"discount_coupon_count"`
}

// PersonPrice 人数料金情報
type PersonPrice struct {
	UseDate               string        `json:"use_date"`
	Person                int           `json:"person"`
	Child1Person          int           `json:"child_1_person"`
	Child2Person          int           `json:"child_2_person"`
	Child3Person          int           `json:"child_3_person"`
	Child4Person          int           `json:"child_4_person"`
	Child5Person          int           `json:"child_5_person"`
	Child6Person          int           `json:"child_6_person"`
	PriceInTax            int           `json:"price_in_tax"`
	ChildPrice1InTax      int           `json:"child_price1_in_tax"`
	ChildPrice2InTax      int           `json:"child_price2_in_tax"`
	ChildPrice3InTax      int           `json:"child_price3_in_tax"`
	ChildPrice4InTax      int           `json:"child_price4_in_tax"`
	ChildPrice5InTax      int           `json:"child_price5_in_tax"`
	ChildPrice6InTax      int           `json:"child_price6_in_tax"`
}

// DetailRoomAndPlan 予約詳細情報の部屋・プラン
type DetailRoomAndPlan struct {
	RoomID         string `json:"room_id"`
	RoomName       string `json:"room_name"`
	PlanName       string `json:"plan_name"`
	FamilyName     string `json:"family_name"`
	GivenName      string `json:"given_name"`
	NumberOfChilds int    `json:"number_of_childs"`
	NumberOfAdults int    `json:"number_of_adults"`
	ChildAges      string `json:"child_ages"`
	Person         int    `json:"person"`
	Child1Person   int    `json:"child_1_person"`
	Child2Person   int    `json:"child_2_person"`
	Child3Person   int    `json:"child_3_person"`
	Child4Person   int    `json:"child_4_person"`
	Child5Person   int    `json:"child_5_person"`
	Child6Person   int    `json:"child_6_person"`
}

// CancelInput 予約キャンセルの入力
type CancelInput struct {
	CmApplicationID int64 `json:"cm_application_id" validate:"required"`
	CancelFee       int64 `json:"cancel_fee"`
	Noshow          uint8 `json:"noshow"`
}

// CancelPolicy キャンセルポリシー
type CancelPolicy struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	Percent string `json:"percent"`
}

// NoShowInput NoShowの入力
type NoShowInput struct {
	CmApplicationID int64 `json:"cm_application_id" validate:"required"`
	NoshowFlg       bool  `json:"noshow_flg"`
}

// IBookingUsecase 予約関連のusecaseのインターフェース
type IBookingUsecase interface {
	SearchBookings(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, req SearchInput) ([]SearchOutput, error)
	BookingDownloads(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, req DownloadInput) ([]BookingDownloadOutput, error)
	DetailBooking(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, req DetailInput) (*DetailOutput, error)
	CancelBooking(req CancelInput) (bool, error)
	UpdateNoShow(req *NoShowInput) error
}

// IBookingRepository 予約関連のrepositoryのインターフェース
type IBookingRepository interface {
	// FetchBookings 検索条件に基づいて予約一覧を複数件取得（hotelリポジトリ参照）
	FetchBookings(req SearchInput) ([]SearchDBOutput, error)
	// FetchDetailApplicationData 予約詳細情報を一件取得（hotelリポジトリ参照）
	FetchDetailApplicationData(req DetailInput) (*DetailApplicationDBOutput, error)
	// FetchBookingDownloadData 予約詳細情報を複数取得（hotelリポジトリ参照）
	FetchBookingDownloadData(req DownloadInput) (*[]BookingDownloadDBOutput, error)
	// FetchBookingRoomsByApplicationID ht_th_application_idに基づく部屋の予約情報を複数件取得
	FetchBookingRoomsByApplicationID(HtThApplicationID int64) ([]HtThBookingRooms, error)
	// FetchBookingRoomListByApplicationID ht_th_application_idに基づく部屋の予約情報を複数件取得
	FetchBookingRoomListByApplicationID(HtThApplicationIDs []int64) ([]HtThBookingRooms, error)
	// FetchRoomListTlsByItineraryID itinerary_idに基づく部屋・プラン情報を複数件取得
	FetchRoomListTlsByItineraryID(ItineraryIDs []string) ([]HtTmItineraryTls, error)
	// FetchNoShowData 予約IDに基づくNoShowのデータを１件取得
	FetchNoShowData(CmApplicationID int64) (HtThApplications, error)
	// UpdateNoShow ht_th_application_idに基づくデータのNoShowフラグを更新
	UpdateNoShow(HtThApplicationID int64, noShowFlg bool, noShowFee float32) error
	// FetchFlashSaleData 予約IDに基づくセールデータを取得
	FetchFlashSaleData(CmApplicationIDs []int64) ([]CmThFlashSale, error)
	// FetchBookingPriceData 予約IDに基づく予約料金データを取得
	FetchBookingPriceData(CmApplicationIDs []int64) ([]HtThBookingPrices, error)
}

// IBookingAPI 予約関連のAPIのインターフェース
type IBookingAPI interface {
	CancelBooking(cmApplicationID int64, cancelFee int64, noShow uint8) (bool, error)
}
