package infra

import (
	"strings"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/booking"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"gorm.io/gorm"
)

// bookingRepository 予約関連repository
type bookingRepository struct {
	hotelDB *gorm.DB
}

// NewBookingRepository インスタンス生成
func NewBookingRepository(hotelDB *gorm.DB) booking.IBookingRepository {
	return &bookingRepository{
		hotelDB: hotelDB,
	}
}

// FetchBookings 検索条件に基づいて予約一覧を複数件取得（hotelリポジトリ参照）
func (b *bookingRepository) FetchBookings(req booking.SearchInput) ([]booking.SearchDBOutput, error) {
	targetWholesalers := []int{utils.WholesalerIDTl, utils.WholesalerIDTema, utils.WholesalerIDNeppan, utils.WholesalerIDDirect, utils.WholesalerIDRaku2}
	result := []booking.SearchDBOutput{}
	query := b.hotelDB.Debug().
		Select("applications.cm_application_id",
			"applications.total_pay_in_tax",
			"applications.cancel_flg",
			"applications.noshow_flg",
			"cm_th_application.application_cd",
			"tours_application.id as tour_id",
			"applications.arrival",
			"applications.departure",
			"applications.given_name_enc",
			"applications.family_name_enc",
			"applications.phone_enc",
			"cm_th_application.payment_limit_dt",
			"IFNULL((SELECT count(1) FROM skyticket.cm_th_payment WHERE cm_th_payment.cm_application_id = applications.cm_application_id) ,0) AS payment_count").
		Table("ht_th_applications AS applications").
		Joins("LEFT JOIN ht_tm_property_langs AS langs ON applications.property_id = langs.property_id").
		Joins("LEFT JOIN skyticket.cm_th_application ON applications.cm_application_id = cm_th_application.cm_application_id").
		Joins("LEFT JOIN skyticket.cm_th_organized_tours_application AS tours_application ON applications.cm_application_id = tours_application.cm_application_id").
		Where("applications.property_id = ?", req.PropertyID).
		Where("applications.wholesaler_id IN ?", targetWholesalers).
		Where("langs.lang_cd = ?", "ja-JP")
	// 申し込み番号
	var ApplicationIDs []int64
	// 申込番号のリストが空でないかチェック
	if len(req.ApplicationIDs) != 0 {
		// 0以外をスライスApplicationIDsに詰め直す
		for _, v := range req.ApplicationIDs {
			if v != 0 {
				ApplicationIDs = append(ApplicationIDs, v)
			}
		}
		// ApplicationIDsが空でければ、WHERE句追加
		if len(ApplicationIDs) != 0 {
			query = query.Where(b.hotelDB.Where("applications.cm_application_id IN ?", ApplicationIDs).
				Or("applications.itinerary_id IN ?", ApplicationIDs))
		}
	}
	// 申込日　開始日
	if req.ApplicationStart != "" {
		query = query.Where("applications.created_at >= ?", req.ApplicationStart+" 00:00:00")
	}
	// 申込日　終了日
	if req.ApplicationEnd != "" {
		query = query.Where("applications.created_at <= ?", req.ApplicationEnd+" 23:59:59")
	}
	// checkout 開始日
	if req.CheckoutStart != "" {
		query = query.Where("applications.departure >= ?", req.CheckoutStart)
	}
	// checkout 終了日
	if req.CheckoutEnd != "" {
		t, _ := time.Parse("2006-01-02", req.CheckoutEnd)
		query = query.Where(
			b.hotelDB.Where("applications.departure <= ? AND applications.cancel_flg = ?", req.CheckoutEnd, 0).
				Or("applications.departure <= ? AND applications.cancel_flg = ?", t.AddDate(0, 0, 1).Format("2006-01-02"), 1))
	}
	// checkin 開始日
	if req.CheckinStart != "" {
		query = query.Where("applications.arrival >= ?", req.CheckinStart)
	}
	// checkin 終了日
	if req.CheckinEnd != "" {
		query = query.Where("applications.arrival <= ?", req.CheckinEnd)
	}
	// family_name_enc 予約者性
	if len(req.FamilyNameEncList) != 0 {
		query = query.Where("applications.family_name_enc IN ?", req.FamilyNameEncList)
	}
	// given_name_enc 予約者名
	if len(req.GivenNameEncList) != 0 {
		query = query.Where("applications.given_name_enc IN ?", req.GivenNameEncList)
	}
	// phone_enc 電話番号
	if req.PhoneEnc != "" {
		query = query.Where("applications.phone_enc = ?", req.PhoneEnc)
	}

	switch req.Status {
	case utils.ReserveStatusReserved:
		query = query.Where("applications.cancel_flg = ?", 0).
			Where("applications.departure > ?", time.Now().Format("2006-01-02")).
			Where("applications.arrival > ?", time.Now().Format("2006-01-02"))
	case utils.ReserveStatusCancel:
		query = query.Where("applications.cancel_flg = ?", 1).
			Where("applications.noshow_flg = ?", 0)
	case utils.ReserveStatusNoShow:
		query = query.Where("applications.cancel_flg = ?", 1).
			Where("applications.noshow_flg = ?", 1)
	case utils.ReserveStatusStaying:
		query = query.Where("applications.cancel_flg = ?", 0).
			Where("applications.departure >= ?", time.Now().Format("2006-01-02")).
			Where("applications.arrival <= ?", time.Now().Format("2006-01-02"))
	case utils.ReserveStatusStayed:
		query = query.Where("applications.cancel_flg = ?", 0).
			Where("applications.departure < ?", time.Now().Format("2006-01-02"))
	}
	err := query.Group(strings.Join([]string{
		"applications.ht_th_application_id",
		"applications.created_at",
		"langs.name",
		"applications.total_pay_in_tax",
		"applications.total_pay_ex_tax",
		"applications.cancel_flg",
		"applications.cancel_fee",
		"applications.noshow_flg",
		"applications.noshow_fee",
		"applications.canceled_dt",
		"cm_th_application.payment_limit_dt", /*支払期限日時*/
		"cm_th_application.total_price",      /*GoTo割引適用後代金*/
		"applications.itinerary_id",
		"applications.affiliate_reference_id",
		"applications.customer_ip",
		"applications.cm_application_id",
		"applications.wholesaler_id",
		"applications.arrival",
		"applications.departure",
		"applications.given_name_enc",
		"applications.family_name_enc",
		"applications.email_enc",
		"applications.phone_enc",
		"applications.memo",
		"applications.is_misuse"}, ",")).
		Find(&result).Error
	return result, err
}

// FetchBookingDownloadData 予約詳細情報を複数取得（hotelリポジトリ参照）
func (b *bookingRepository) FetchBookingDownloadData(req booking.DownloadInput) (*[]booking.BookingDownloadDBOutput, error) {
	result := []booking.BookingDownloadDBOutput{}
	err := b.hotelDB.
		Select("applications.cm_application_id",
			"applications.ht_th_application_id",
			"applications.wholesaler_id",
			"applications.property_id",
			"applications.itinerary_id",
			"applications.created_at",
			"applications.total_pay_in_tax",
			"applications.cancel_flg",
			"applications.cancel_fee",
			"applications.canceled_dt",
			"applications.noshow_flg",
			"applications.noshow_fee",
			"cm_th_application.application_cd",
			"tours_application.id as tour_id",
			"applications.arrival",
			"applications.departure",
			"applications.stays",
			"applications.room_num",
			"applications.given_name_enc",
			"applications.family_name_enc",
			"applications.phone_enc",
			"applications.email_enc",
			"cm_th_application.payment_limit_dt",
			"IFNULL((SELECT count(1) FROM skyticket.cm_th_payment WHERE cm_th_payment.cm_application_id = applications.cm_application_id) ,0) AS payment_count").
		Table("ht_th_applications AS applications").
		Joins("LEFT JOIN ht_tm_property_langs AS langs ON applications.property_id = langs.property_id").
		Joins("LEFT JOIN skyticket.cm_th_application ON applications.cm_application_id = cm_th_application.cm_application_id").
		Joins("LEFT JOIN skyticket.cm_th_organized_tours_application AS tours_application ON applications.cm_application_id = tours_application.cm_application_id").
		Where("applications.cm_application_id IN ?", req.CmApplicationIDs).
		Where("applications.property_id = ?", req.PropertyID).
		Where("langs.lang_cd = ?", "ja-JP").
		Group("applications.cm_application_id").
		Find(&result).Error
	return &result, err
}

// FetchDetailApplicationData 予約詳細情報を一件取得（hotelリポジトリ参照）
func (b *bookingRepository) FetchDetailApplicationData(req booking.DetailInput) (*booking.DetailApplicationDBOutput, error) {
	result := &booking.DetailApplicationDBOutput{}
	err := b.hotelDB.
		Select("applications.cm_application_id",
			"applications.ht_th_application_id",
			"applications.wholesaler_id",
			"applications.property_id",
			"applications.itinerary_id",
			"applications.created_at",
			"applications.total_pay_in_tax",
			"applications.cancel_flg",
			"applications.cancel_fee",
			"applications.canceled_dt",
			"applications.noshow_flg",
			"applications.noshow_fee",
			"cm_th_application.application_cd",
			"tours_application.id as tour_id",
			"applications.arrival",
			"applications.departure",
			"applications.stays",
			"applications.room_num",
			"applications.given_name_enc",
			"applications.family_name_enc",
			"applications.phone_enc",
			"applications.email_enc",
			"cm_th_application.payment_limit_dt",
			"IFNULL((SELECT count(1) FROM skyticket.cm_th_payment WHERE cm_th_payment.cm_application_id = applications.cm_application_id) ,0) AS payment_count").
		Table("ht_th_applications AS applications").
		Joins("LEFT JOIN ht_tm_property_langs AS langs ON applications.property_id = langs.property_id").
		Joins("LEFT JOIN skyticket.cm_th_application ON applications.cm_application_id = cm_th_application.cm_application_id").
		Joins("LEFT JOIN skyticket.cm_th_organized_tours_application AS tours_application ON applications.cm_application_id = tours_application.cm_application_id").
		Where("applications.cm_application_id = ?", req.CmApplicationID).
		Where("applications.property_id = ?", req.PropertyID).
		Where("langs.lang_cd = ?", "ja-JP").
		First(result).Error
	return result, err
}

// FetchBookingRoomsByApplicationID ht_th_application_idに基づく部屋の予約情報を複数件取得
func (b *bookingRepository) FetchBookingRoomsByApplicationID(HtThApplicationID int64) ([]booking.HtThBookingRooms, error) {
	result := []booking.HtThBookingRooms{}
	err := b.hotelDB.
		Model(&booking.HtThBookingRooms{}).
		Where("ht_th_application_id = ?", HtThApplicationID).
		Find(&result).Error
	return result, err
}

// FetchBookingRoomListByApplicationID ht_th_application_idに基づく部屋の予約情報を複数件取得
func (b *bookingRepository) FetchBookingRoomListByApplicationID(HtThApplicationIDs []int64) ([]booking.HtThBookingRooms, error) {
	result := []booking.HtThBookingRooms{}
	err := b.hotelDB.
		Model(&booking.HtThBookingRooms{}).
		Where("ht_th_application_id IN ?", HtThApplicationIDs).
		Find(&result).Error
	return result, err
}

// FetchRoomListTlsByItineraryID itinerary_idに基づく部屋・プラン情報を複数件取得(TLのみ)
func (b *bookingRepository) FetchRoomListTlsByItineraryID(ItineraryIDs []string) ([]booking.HtTmItineraryTls, error) {
	result := []booking.HtTmItineraryTls{}
	err := b.hotelDB.
		Model(&booking.HtTmItineraryTls{}).
		Where("itinerary_id IN ?", ItineraryIDs).
		Find(&result).Error
	return result, err
}

// FetchNoShowData 予約IDに基づくNoShowのデータを１件取得
func (b *bookingRepository) FetchNoShowData(CmApplicationID int64) (booking.HtThApplications, error) {
	result := booking.HtThApplications{}
	err := b.hotelDB.
		Model(&booking.HtThApplications{}).
		Where("cm_application_id = ?", CmApplicationID).
		Where("cancel_flg = 1").
		Where("arrival <= ?", time.Now().Format("2006-01-02")).
		First(&result).Error
	return result, err
}

// UpdateNoShow ht_th_application_idに基づくデータのNoShowフラグを更新
func (b *bookingRepository) UpdateNoShow(HtThApplicationID int64, noShowFlg bool, noShowFee float32) error {
	return b.hotelDB.
		Model(&booking.HtThApplications{}).
		Where("ht_th_application_id = ?", HtThApplicationID).
		Updates(map[string]interface{}{
			"noshow_fee": noShowFee,
			"noshow_flg": noShowFlg,
			"updated_at": time.Now(),
		}).Error
}

// FetchFlashSaleData 予約IDに基づくセールデータを取得
func (b *bookingRepository) FetchFlashSaleData(CmApplicationIDs []int64) ([]booking.CmThFlashSale, error) {
	result := []booking.CmThFlashSale{}
	err := b.hotelDB.
		Table("skyticket.cm_th_flash_sale").
		Where("cm_application_id in ?", CmApplicationIDs).
		Find(&result).Error
	return result, err
}

// FetchBookingPriceData 予約IDに基づく予約料金データを取得
func (b *bookingRepository) FetchBookingPriceData(CmApplicationIDs []int64) ([]booking.HtThBookingPrices, error) {
	result := []booking.HtThBookingPrices{}
	err := b.hotelDB.
		Model(&booking.HtThBookingPrices{}).
		Where("cm_application_id in ?", CmApplicationIDs).
		Find(&result).Error
	return result, err
}
