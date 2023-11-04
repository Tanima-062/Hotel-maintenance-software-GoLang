package usecase

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/booking"
	"github.com/Adventureinc/hotel-hm-api/src/booking/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	pInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"gorm.io/gorm"
)

const (
	// TimeFormat 時刻取得のフォーマット
	TimeFormat = "2006-01-02 15:04:05"

	// DateFormat 日付取得のフォーマット
	DateFormat = "2006-01-02"

	// KagoshimaCouponFee クーポン1枚あたりの料金
	KagoshimaCouponFee = 1000

	// SapporoCouponFee クーポン1枚あたりの料金
	SapporoCouponFee = 1000
)

// bookingUsecase 　共通部分の予約関連usecase
type bookingUsecase struct {
	BRepository          booking.IBookingRepository
	RoomDirectRepository room.IRoomDirectRepository
	PlanDirectRepository plan.IPlanDirectRepository
	RoomNeppanRepository room.IRoomNeppanRepository
	PlanNeppanRepository plan.IPlanNeppanRepository
	RoomRaku2Repository  room.IRoomRaku2Repository
	PlanRaku2Repository  plan.IPlanRaku2Repository
	RoomTemaRepository   room.IRoomTemaRepository
	PlanTemaRepository   plan.IPlanTemaRepository
	BAPI                 booking.IBookingAPI
}

// NewBookingUsecase インスタンス生成
func NewBookingUsecase(hotelDB *gorm.DB) booking.IBookingUsecase {
	return &bookingUsecase{
		BRepository:          infra.NewBookingRepository(hotelDB),
		RoomDirectRepository: rInfra.NewRoomDirectRepository(hotelDB),
		PlanDirectRepository: pInfra.NewPlanDirectRepository(hotelDB),
		RoomNeppanRepository: rInfra.NewRoomNeppanRepository(hotelDB),
		PlanNeppanRepository: pInfra.NewPlanNeppanRepository(hotelDB),
		RoomRaku2Repository:  rInfra.NewRoomRaku2Repository(hotelDB),
		PlanRaku2Repository:  pInfra.NewPlanRaku2Repository(hotelDB),
		RoomTemaRepository:   rInfra.NewRoomTemaRepository(hotelDB),
		PlanTemaRepository:   pInfra.NewPlanTemaRepository(hotelDB),
		BAPI:                 infra.NewBookingAPI(),
	}
}

// SearchBookings 予約一覧検索
func (b *bookingUsecase) SearchBookings(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, req booking.SearchInput) ([]booking.SearchOutput, error) {
	res := []booking.SearchOutput{}
	/*
	* 予約者性(FamilyName),予約者名(GivenName),電話番号(Phone)は暗号化した文字列を予約テーブルに登録しているので
	* これらを条件に取得するためには一度暗号化する。
	 */
	//　暗号化する予約者名(GivenNameEnc)のリストを作成
	givenNameList := utils.UpperAndLowerStrList(req.GivenName)
	//　暗号化する予約者性(FamilyNameEnc)のリストを作成
	familyNameList := utils.UpperAndLowerStrList(req.FamilyName)

	//　givenNameEncListを作成
	var givenNameEncList []string
	for _, v := range givenNameList {
		givenNameEnc, eErr := utils.Encrypt(v)
		if eErr != nil {
			return res, eErr
		}
		givenNameEncList = append(givenNameEncList, givenNameEnc)
	}
	req.GivenNameEncList = givenNameEncList

	//　familyNameEncListを作成
	var familyNameEncList []string
	for _, v := range familyNameList {
		familyNameEnc, eErr := utils.Encrypt(v)
		if eErr != nil {
			return res, eErr
		}
		familyNameEncList = append(familyNameEncList, familyNameEnc)
	}
	req.FamilyNameEncList = familyNameEncList

	//　Phoneを暗号化
	phoneEnc, eErr := utils.Encrypt(req.Phone)
	if eErr != nil {
		return res, eErr
	}
	req.PhoneEnc = phoneEnc

	bookings, err := b.BRepository.FetchBookings(req)
	if err != nil {
		return res, err
	}

	var givenNames, familyNames, phones []string
	for _, v := range bookings {
		givenNameDec, dErr := utils.Decrypt(v.GivenNameEnc)
		if dErr != nil {
			return res, dErr
		}
		familyNameDec, dErr := utils.Decrypt(v.FamilyNameEnc)
		if dErr != nil {
			return res, dErr
		}
		phoneDec, dErr := utils.Decrypt(v.PhoneEnc)
		if dErr != nil {
			return res, dErr
		}
		givenNames = append(givenNames, givenNameDec)
		familyNames = append(familyNames, familyNameDec)
		phones = append(phones, phoneDec)
	}

	for index, v := range bookings {
		res = append(res, booking.SearchOutput{
			CmApplicationID: v.CmApplicationID,
			ApplicationCd:   v.ApplicationCd,
			TourID:          v.TourID,
			Checkin:         v.Arrival,
			Checkout:        v.Departure,
			GivenName:       givenNames[index],
			FamilyName:      familyNames[index],
			Phone:           phones[index],
			TotalPayInTax:   v.TotalPayInTax,
			Status:          utils.GetBookingStatus(v.CancelFlg, v.NoshowFlg, v.Arrival, v.Departure),
		})
	}
	return res, nil
}

// DetailBooking 予約詳細情報取得
func (b *bookingUsecase) DetailBooking(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, req booking.DetailInput) (*booking.DetailOutput, error) {
	response := &booking.DetailOutput{}
	appData, err := b.BRepository.FetchDetailApplicationData(req)
	if err != nil || appData.HtThApplicationID == 0 {
		return response, err
	}
	givenName, dErr := utils.Decrypt(appData.GivenNameEnc)
	if dErr != nil {
		return response, dErr
	}
	familyName, dErr := utils.Decrypt(appData.FamilyNameEnc)
	if dErr != nil {
		return response, dErr
	}
	email, dErr := utils.Decrypt(appData.EmailEnc)
	if dErr != nil {
		return response, dErr
	}
	phone, dErr := utils.Decrypt(appData.PhoneEnc)
	if dErr != nil {
		return response, dErr
	}

	// 予約料金情報を取得
	bookingPrices, bpErr := b.BRepository.FetchBookingPriceData([]int64{appData.CmApplicationID})
	if bpErr != nil {
		return response, bpErr
	}

	// 取得した予約料金情報から、日付毎の人数と料金の内訳を抽出
	personPrices := map[string][]booking.PersonPrice{}
	for _, bookingPrice := range bookingPrices {
		currentUseDate := bookingPrice.UseDate.Format(DateFormat)
		pd := personPrices[currentUseDate]
		p := &booking.PersonPrice{}
		p.UseDate = currentUseDate
		p.Person = bookingPrice.Person
		p.Child1Person = bookingPrice.Child1Person
		p.Child2Person = bookingPrice.Child2Person
		p.Child3Person = bookingPrice.Child3Person
		p.Child4Person = bookingPrice.Child4Person
		p.Child5Person = bookingPrice.Child5Person
		p.Child6Person = bookingPrice.Child6Person
		p.PriceInTax = bookingPrice.PriceInTax
		p.ChildPrice1InTax = bookingPrice.ChildPrice1InTax
		p.ChildPrice2InTax = bookingPrice.ChildPrice2InTax
		p.ChildPrice3InTax = bookingPrice.ChildPrice3InTax
		p.ChildPrice4InTax = bookingPrice.ChildPrice4InTax
		p.ChildPrice5InTax = bookingPrice.ChildPrice5InTax
		p.ChildPrice6InTax = bookingPrice.ChildPrice6InTax
		pd = append(pd, *p)
		personPrices[currentUseDate] = pd
	}
	response.PersonPrices = personPrices

	// 部屋ごとの人数内訳が1日分だけ必要なため、取得した予約料金情報から宿泊初日のデータのみを切り出す
	firstUseDate := ""
	firstDateBookingPrices := []booking.HtThBookingPrices{}
	for _, bookingPrice := range bookingPrices {
		currentUseDate := bookingPrice.UseDate.Format(DateFormat)
		if firstUseDate == "" || firstUseDate == currentUseDate {
			firstDateBookingPrices = append(firstDateBookingPrices, bookingPrice)
			firstUseDate = currentUseDate
		}
	}

	bookingRooms, rErr := b.BRepository.FetchBookingRoomsByApplicationID(appData.HtThApplicationID)
	if rErr != nil {
		return response, rErr
	}
	// 客室詳細情報を設定する
	temp := []booking.DetailRoomAndPlan{}
	for i, bookingRoom := range bookingRooms {
		// 各部屋の代表者名
		roomGivenName, dErr := utils.Decrypt(bookingRoom.GivenNameEnc)
		if dErr != nil {
			return response, dErr
		}
		roomFamilyName, dErr := utils.Decrypt(bookingRoom.FamilyNameEnc)
		if dErr != nil {
			return response, dErr
		}
		// 格納する客室情報の構造体
		detailRoomAndPlan := &booking.DetailRoomAndPlan{
			RoomID:         bookingRoom.RoomID,
			RoomName:       bookingRoom.RoomName,
			PlanName:       bookingRoom.PlanName,
			FamilyName:     roomFamilyName,
			GivenName:      roomGivenName,
			NumberOfAdults: bookingRoom.NumberOfAdults,
			NumberOfChilds: bookingRoom.NumberOfChilds,
			ChildAges:      bookingRoom.ChildAges,
		}
		// 部屋ごとの人数内訳をマージ（客室情報と料金情報の部屋を明示的に紐づける情報がないため、「客室情報のi部屋目 ＝ 料金情報のi部屋目」とみなす）
		if len(firstDateBookingPrices) > i {
			detailRoomAndPlan.Person = firstDateBookingPrices[i].Person
			detailRoomAndPlan.Child1Person = firstDateBookingPrices[i].Child1Person
			detailRoomAndPlan.Child2Person = firstDateBookingPrices[i].Child2Person
			detailRoomAndPlan.Child3Person = firstDateBookingPrices[i].Child3Person
			detailRoomAndPlan.Child4Person = firstDateBookingPrices[i].Child4Person
			detailRoomAndPlan.Child5Person = firstDateBookingPrices[i].Child5Person
			detailRoomAndPlan.Child6Person = firstDateBookingPrices[i].Child6Person
		} else {
			// 予約時に料金情報を保持しないホールセラーの場合、客室情報に保持している人数内訳をマージ
			detailRoomAndPlan.Person = bookingRoom.NumberOfAdults
			detailRoomAndPlan.Child1Person = bookingRoom.NumberOfUpperGrades
			detailRoomAndPlan.Child2Person = bookingRoom.NumberOfLowerGrades
			detailRoomAndPlan.Child3Person = bookingRoom.NumberOfInfantMealsWithBedding
			detailRoomAndPlan.Child4Person = bookingRoom.NumberOfInfantMealOnly
			detailRoomAndPlan.Child5Person = bookingRoom.NumberOfInfantBeddingOnly
			detailRoomAndPlan.Child6Person = bookingRoom.NumberOfInfantMealsWithoutBedding
		}

		// 部屋・プラン名が予約情報に保持されている場合、部屋・プラン情報をスライスに追加して次に進む
		if bookingRoom.RoomName != "" && bookingRoom.PlanName != "" {
			temp = append(temp, *detailRoomAndPlan)
			continue
		}

		// 部屋名・プラン名が予約情報に保持されていない場合、予約情報のホールセラーIDに応じて部屋・プラン名を取得する
		switch appData.WholesalerID {
		// TLリンカーン（取得不可なので何もしない）
		case utils.WholesalerIDTl:
			break
		// 手間いらず
		case utils.WholesalerIDTema:
			if bookingRoom.RoomName != "" {
				roomId, parseErr := strconv.Atoi(bookingRoom.RoomID)
				if parseErr == nil {
					roomData, rErr := b.RoomTemaRepository.FetchOne(roomId, appData.PropertyID)
					if rErr == nil {
						detailRoomAndPlan.RoomName = roomData.RoomNameJa
					}
				}
			}
			if bookingRoom.PlanName != "" {
				rateId, parseErr := strconv.Atoi(bookingRoom.RateID)
				if parseErr == nil {
					planData, pErr := b.PlanTemaRepository.FetchOnePlan(appData.PropertyID, rateId)
					if pErr == nil {
						detailRoomAndPlan.PlanName = planData.PlanName
					}
				}
			}
		// ねっぱん
		case utils.WholesalerIDNeppan:
			if bookingRoom.RoomName != "" {
				roomId, parseErr := strconv.ParseInt(bookingRoom.RoomID, 10, 64)
				if parseErr == nil {
					roomData, rErr := b.RoomNeppanRepository.FetchRoomByRoomTypeID(roomId)
					if rErr == nil {
						detailRoomAndPlan.RoomName = roomData.Name
					}
				}
			}
			if bookingRoom.PlanName != "" {
				planId, parseErr := strconv.ParseInt(bookingRoom.RateID, 10, 64)
				if parseErr == nil {
					planData, pErr := b.PlanNeppanRepository.FetchOne(planId)
					if pErr == nil {
						detailRoomAndPlan.PlanName = planData.Name
					}
				}
			}
		// 直仕入
		case utils.WholesalerIDDirect:
			if bookingRoom.RoomName != "" {
				roomId, parseErr := strconv.ParseInt(bookingRoom.RoomID, 10, 64)
				if parseErr == nil {
					roomData, rErr := b.RoomDirectRepository.FetchRoomByRoomTypeID(roomId)
					if rErr == nil {
						detailRoomAndPlan.RoomName = roomData.Name
					}
				}
			}
			if bookingRoom.PlanName != "" {
				planId, parseErr := strconv.ParseInt(bookingRoom.RateID, 10, 64)
				if parseErr == nil {
					planData, pErr := b.PlanDirectRepository.FetchOne(planId)
					if pErr == nil {
						detailRoomAndPlan.PlanName = planData.Name
					}
				}
			}
		// らく通
		case utils.WholesalerIDRaku2:
			if bookingRoom.RoomName != "" {
				roomId, parseErr := strconv.ParseInt(bookingRoom.RoomID, 10, 64)
				if parseErr == nil {
					roomData, rErr := b.RoomRaku2Repository.FetchRoomByRoomTypeID(roomId)
					if rErr == nil {
						detailRoomAndPlan.RoomName = roomData.Name
					}
				}
			}
			if bookingRoom.PlanName != "" {
				planId, parseErr := strconv.ParseInt(bookingRoom.RateID, 10, 64)
				if parseErr == nil {
					planData, pErr := b.PlanRaku2Repository.FetchOne(planId)
					if pErr == nil {
						detailRoomAndPlan.PlanName = planData.Name
					}
				}
			}
		}
		// 部屋・プラン情報をスライスに追加（switch文の途中でエラーになった場合も、取得できた分を追加する）
		temp = append(temp, *detailRoomAndPlan)

	}
	response.RoomsAndPlan = temp

	cancelFeeSuggest, err := SuggestCancelFee(bookingRooms, appData)
	if err != nil {
		return response, err
	}
	response.CancelFeeSuggest = cancelFeeSuggest

	// セール情報を取得
	flashSales, fErr := b.BRepository.FetchFlashSaleData([]int64{appData.CmApplicationID})
	if fErr != nil {
		return response, fErr
	}
	// 画面表示用のセール情報を生成
	flashSaleOutput := []booking.FlashSale{}
	for _, flashSale := range flashSales {
		// 個々のセール情報を生成
		temp := &booking.FlashSale{}
		temp.DiscountCashAmount = flashSale.DiscountCashAmount
		temp.DiscountCouponCount = -1

		switch {
		case strings.HasPrefix(flashSale.SaleType, "KENWARI_IMAKAGO2"):
			// 鹿児島キャンペーン
			temp.SaleName = "今こそ鹿児島"
			temp.DiscountCouponCount = int(flashSale.DiscountCouponAmount) / KagoshimaCouponFee
		case strings.HasPrefix(flashSale.SaleType, "KENWARI_SAPPORO2"):
			// 札幌キャンペーン
			temp.SaleName = "さっぽろ冬割"
			temp.DiscountCouponCount = int(flashSale.DiscountCouponAmount) / SapporoCouponFee
		case strings.HasPrefix(flashSale.SaleType, "KENWARI_HAKODATE2021"):
			// 函館キャンペーン
			temp.SaleName = "はこだて割"
		case strings.HasPrefix(flashSale.SaleType, "GOTO"):
			if strings.HasPrefix(flashSale.SaleType, "GOTO_ZENKOKUSHIEN") {
				// 全国旅行支援
				temp.SaleName = "全国旅行支援"
			} else {
				// GoTo
				temp.SaleName = "GoToキャンペーン"
			}
		case strings.HasPrefix(flashSale.SaleType, "MOTTOTOKYO"):
			// もっとTokyo
			temp.SaleName = "もっとTokyo"
		case strings.HasPrefix(flashSale.SaleType, "KENWARI_MOTTOTOKYO"):
			// もっとTokyo
			temp.SaleName = "もっとTokyo"
		case strings.HasPrefix(flashSale.SaleType, "KENWARI_CHIBATOKU"):
			// 千葉とく旅
			temp.SaleName = "千葉とく割"
		default:
			// 名称不明
			temp.SaleName = "割引セール"
		}
		flashSaleOutput = append(flashSaleOutput, *temp)

		// 割引後金額を取得（適用可能な全てのセールを適用済みの金額が設定されているので、どれか1つを取得すればよい）
		response.SalePrice = flashSale.SalePrice
		// セールが1つでも適用済みで割引対象の場合、割引フラグを立てる
		if !(response.DiscountPaymentFlg) && flashSale.DiscountPaymentFlg {
			response.DiscountPaymentFlg = true
		}
		// 割引額を加算
		response.DiscountCashAmount += flashSale.DiscountCashAmount
	}
	response.FlashSales = flashSaleOutput

	response.CmApplicationID = appData.CmApplicationID
	response.WholesalerID = appData.WholesalerID
	response.ApplicationCd = appData.ApplicationCd
	response.TourID = appData.TourID
	response.CreatedAt = appData.CreatedAt
	response.GivenNameEnc = givenName
	response.FamilyNameEnc = familyName
	response.EmailEnc = email
	response.TotalPayInTax = appData.TotalPayInTax
	response.CancelFee = appData.CancelFee
	response.CancelFlg = appData.CancelFlg
	response.CanceledDt = appData.CanceledDt
	response.NoshowFee = appData.NoshowFee
	response.NoshowFlg = appData.NoshowFlg
	response.Arrival = appData.Arrival
	response.Departure = appData.Departure
	response.Stays = appData.Stays
	response.RoomNum = appData.RoomNum
	response.PhoneEnc = phone
	response.Status = utils.GetBookingStatus(appData.CancelFlg, appData.NoshowFlg, appData.Arrival, appData.Departure)
	return response, nil
}

// BookingDownloads CSVでダウンロードする予約詳細情報一覧取得
func (b *bookingUsecase) BookingDownloads(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, req booking.DownloadInput) ([]booking.BookingDownloadOutput, error) {
	response := []booking.BookingDownloadOutput{}
	// 予約情報詳細一覧取得
	bookingDownloads, err := b.BRepository.FetchBookingDownloadData(req)
	if err != nil {
		return response, err
	}
	var givenNames, familyNames, emails, phones []string
	htThApplicationIds := []int64{}
	cmApplicationIDs := []int64{}
	wholesalerIDList := make(map[int64]int64)
	for _, v := range *bookingDownloads {
		// 復号
		givenNameDec, dErr := utils.Decrypt(v.GivenNameEnc)
		if dErr != nil {
			return response, dErr
		}
		familyNameDec, dErr := utils.Decrypt(v.FamilyNameEnc)
		if dErr != nil {
			return response, dErr
		}
		emailDec, dErr := utils.Decrypt(v.EmailEnc)
		if dErr != nil {
			return response, dErr
		}
		phoneDec, dErr := utils.Decrypt(v.PhoneEnc)
		if dErr != nil {
			return response, dErr
		}
		givenNames = append(givenNames, givenNameDec)
		familyNames = append(familyNames, familyNameDec)
		emails = append(emails, emailDec)
		phones = append(phones, phoneDec)
		// HtThApplicationIDのリスト
		htThApplicationIds = append(htThApplicationIds, v.HtThApplicationID)
		// CmApplicationIDのリスト
		cmApplicationIDs = append(cmApplicationIDs, v.CmApplicationID)
		// キーをHtThApplicationID、値をホールセラーIDにしたmapを作成
		wholesalerIDList[v.HtThApplicationID] = v.WholesalerID
	}
	// htThApplicationIdsを元に部屋情報取得
	bookingRoomListById, rErr := b.BRepository.FetchBookingRoomListByApplicationID(htThApplicationIds)
	if rErr != nil {
		return response, rErr
	}

	temaRoomIDList := []int{}
	temaPlanIDList := []int{}
	neppanRoomIDList := []int64{}
	neppanPlanIDList := []int64{}
	directRoomIDList := []int64{}
	directPlanIDList := []int64{}
	raku2RoomIDList := []int64{}
	raku2PlanIDList := []int64{}
	for htThApplicationID, wholesalerID := range wholesalerIDList {
		for _, bookingRoomById := range bookingRoomListById {
			// HtThApplicationIDが一致したらwholesalerIDごとに部屋ID・プランIDをまとめる
			if bookingRoomById.HtThApplicationID == htThApplicationID {
				switch wholesalerID {
				case utils.WholesalerIDTema:
					roomId, _ := strconv.Atoi(bookingRoomById.RoomID)
					planId, _ := strconv.Atoi(bookingRoomById.RateID)
					// RoomIDのリスト
					temaRoomIDList = append(temaRoomIDList, roomId)
					// PlanIDのリスト
					temaPlanIDList = append(temaPlanIDList, planId)
				case utils.WholesalerIDNeppan:
					roomId, _ := strconv.ParseInt(bookingRoomById.RoomID, 10, 64)
					planId, _ := strconv.ParseInt(bookingRoomById.RateID, 10, 64)
					// RoomIDのリスト
					neppanRoomIDList = append(neppanRoomIDList, roomId)
					// PlanIDのリスト
					neppanPlanIDList = append(neppanPlanIDList, planId)
				case utils.WholesalerIDDirect:
					roomId, _ := strconv.ParseInt(bookingRoomById.RoomID, 10, 64)
					planId, _ := strconv.ParseInt(bookingRoomById.RateID, 10, 64)
					// RoomIDのリスト
					directRoomIDList = append(directRoomIDList, roomId)
					// PlanIDのリスト
					directPlanIDList = append(directPlanIDList, planId)
				case utils.WholesalerIDRaku2:
					roomId, _ := strconv.ParseInt(bookingRoomById.RoomID, 10, 64)
					planId, _ := strconv.ParseInt(bookingRoomById.RateID, 10, 64)
					// RoomIDのリスト
					raku2RoomIDList = append(raku2RoomIDList, roomId)
					// PlanIDのリスト
					raku2PlanIDList = append(raku2PlanIDList, planId)
				}
			}
		}
	}
	temaRoomDatas := []room.HtTmRoomTemas{}
	temaPlanDatas := []plan.HtTmPlanTemas{}
	neppanRoomDatas := []room.HtTmRoomTypeNeppans{}
	neppanPlanDatas := []plan.HtTmPlanNeppans{}
	directRoomDatas := []room.HtTmRoomTypeDirects{}
	directPlanDatas := []plan.HtTmPlanDirects{}
	raku2RoomDatas := []room.HtTmRoomTypeRaku2s{}
	raku2PlanDatas := []plan.HtTmPlanRaku2s{}
	switch true {
	case len(temaRoomIDList) > 0 && len(temaPlanIDList) > 0:
		// 予約詳細の部屋情報取得
		temaRoomDatas, _ = b.RoomTemaRepository.FetchListWithPropertyId(temaRoomIDList, req.PropertyID)
		// プラン情報取得
		temaPlanDatas, _ = b.PlanTemaRepository.FetchList(req.PropertyID, temaPlanIDList)
		fallthrough
	case len(neppanRoomIDList) > 0 && len(neppanPlanIDList) > 0:
		// 予約詳細の部屋情報取得
		neppanRoomDatas, _ = b.RoomNeppanRepository.FetchRoomListByRoomTypeID(neppanRoomIDList)
		// プラン情報取得
		neppanPlanDatas, _ = b.PlanNeppanRepository.FetchList(neppanPlanIDList)
		fallthrough
	case len(directRoomIDList) > 0 && len(directPlanIDList) > 0:
		// 予約詳細の部屋情報取得
		directRoomDatas, _ = b.RoomDirectRepository.FetchRoomListByRoomTypeID(directRoomIDList)
		// プラン情報取得
		directPlanDatas, _ = b.PlanDirectRepository.FetchList(directPlanIDList)
		fallthrough
	case len(raku2RoomIDList) > 0 && len(raku2PlanIDList) > 0:
		// 予約詳細の部屋情報取得
		raku2RoomDatas, _ = b.RoomRaku2Repository.FetchRoomListByRoomTypeID(raku2RoomIDList)
		// プラン情報取得
		raku2PlanDatas, _ = b.PlanRaku2Repository.FetchList(raku2PlanIDList)
	}

	// セール情報を取得
	flashSales, fErr := b.BRepository.FetchFlashSaleData(cmApplicationIDs)
	if fErr != nil {
		return response, fErr
	}

	// 予約詳細情報格納
	for index, bookingDownload := range *bookingDownloads {
		temp := []booking.DetailRoomAndPlan{}
		for _, bookingRoomById := range bookingRoomListById {
			bookingRoomId, _ := strconv.ParseInt(bookingRoomById.RoomID, 10, 64)
			bookingRoomPlanId, _ := strconv.ParseInt(bookingRoomById.RateID, 10, 64)
			if bookingRoomById.HtThApplicationID == bookingDownload.HtThApplicationID {
				// 予約情報に部屋名・プラン名が含まれている場合、その内容を格納して次に進む
				if bookingRoomById.RoomName != "" && bookingRoomById.PlanName != "" {
					temp = append(temp, booking.DetailRoomAndPlan{
						RoomID:         bookingRoomById.RoomID,
						RoomName:       bookingRoomById.RoomName,
						PlanName:       bookingRoomById.PlanName,
						NumberOfAdults: bookingRoomById.NumberOfAdults,
						NumberOfChilds: bookingRoomById.NumberOfChilds,
					})
					continue
				}

				// 予約情報に部屋名・プラン名が含まれていない場合、IDを基に取得した内容をマージして格納する
				var roomName string
				var planName string
				// 予約情報のホールセラーIDを見て部屋・プラン情報を格納
				switch bookingDownload.WholesalerID {
				case utils.WholesalerIDTl:
					roomName = bookingRoomById.RoomName
					planName = bookingRoomById.PlanName
				case utils.WholesalerIDTema:
					temaBookingRoomPlanId, _ := strconv.Atoi(bookingRoomById.RateID)
					for _, roomData := range temaRoomDatas {
						if bookingRoomById.RoomID == roomData.RoomTypeCode {
							roomName = roomData.RoomNameJa
							break
						}
					}
					for _, planData := range temaPlanDatas {
						if temaBookingRoomPlanId == planData.PackagePlanCode {
							planName = planData.PlanName
							break
						}
					}
				case utils.WholesalerIDNeppan:
					for _, roomData := range neppanRoomDatas {
						if bookingRoomId == roomData.RoomTypeID {
							roomName = roomData.Name
							break
						}
					}
					for _, planData := range neppanPlanDatas {
						if bookingRoomPlanId == planData.PlanID {
							planName = planData.Name
							break
						}
					}
				case utils.WholesalerIDDirect:
					for _, roomData := range directRoomDatas {
						if bookingRoomId == roomData.RoomTypeID {
							roomName = roomData.Name
							break
						}
					}
					for _, planData := range directPlanDatas {
						if bookingRoomPlanId == planData.PlanID {
							planName = planData.Name
							break
						}
					}
				case utils.WholesalerIDRaku2:
					for _, roomData := range raku2RoomDatas {
						if bookingRoomId == roomData.RoomTypeID {
							roomName = roomData.Name
							break
						}
					}
					for _, planData := range raku2PlanDatas {
						if bookingRoomPlanId == planData.PlanID {
							planName = planData.Name
							break
						}
					}
				}
				// 部屋・プラン情報を詰める（取得できていない項目があっても、取得できた分を返却する）
				temp = append(temp, booking.DetailRoomAndPlan{
					RoomID:         bookingRoomById.RoomID,
					RoomName:       roomName,
					PlanName:       planName,
					NumberOfAdults: bookingRoomById.NumberOfAdults,
					NumberOfChilds: bookingRoomById.NumberOfChilds,
				})
			}
		}

		var salePrice float32
		var discountCashAmount float32
		var discountPaymentFlg bool
		// セール情報から出力項目を生成
		for _, flashSale := range flashSales {
			if flashSale.CmApplicationID != bookingDownload.CmApplicationID {
				continue
			}
			// 割引後金額を取得（適用可能な全てのセールを適用済みの金額が設定されているので、どれか1つを取得すればよい）
			salePrice = flashSale.SalePrice
			// セールが1つでも適用済みで割引対象の場合、割引フラグを立てる
			if flashSale.DiscountPaymentFlg {
				discountPaymentFlg = true
			}
			// 割引額を加算
			discountCashAmount += flashSale.DiscountCashAmount
		}

		response = append(response, booking.BookingDownloadOutput{
			CmApplicationID:    bookingDownload.CmApplicationID,
			ApplicationCd:      bookingDownload.ApplicationCd,
			TourID:             bookingDownload.TourID,
			CreatedAt:          bookingDownload.CreatedAt,
			GivenNameEnc:       givenNames[index],
			FamilyNameEnc:      familyNames[index],
			EmailEnc:           emails[index],
			TotalPayInTax:      bookingDownload.TotalPayInTax,
			SalePrice:          salePrice,
			CancelFee:          bookingDownload.CancelFee,
			CancelFlg:          bookingDownload.CancelFlg,
			CanceledDt:         bookingDownload.CanceledDt,
			NoshowFee:          bookingDownload.NoshowFee,
			NoshowFlg:          bookingDownload.NoshowFlg,
			Arrival:            bookingDownload.Arrival,
			Departure:          bookingDownload.Departure,
			PhoneEnc:           phones[index],
			DiscountPaymentFlg: discountPaymentFlg,
			DiscountCashAmount: discountCashAmount,
			Status:             utils.GetBookingStatus(bookingDownload.CancelFlg, bookingDownload.NoshowFlg, bookingDownload.Arrival, bookingDownload.Departure),
			RoomsAndPlan:       temp,
		})
	}
	return response, nil
}

// CancelBooking 予約キャンセル
func (b *bookingUsecase) CancelBooking(req booking.CancelInput) (bool, error) {
	return b.BAPI.CancelBooking(req.CmApplicationID, req.CancelFee, req.Noshow)
}

// UpdateNoShow NoShowフラグの更新
func (b *bookingUsecase) UpdateNoShow(req *booking.NoShowInput) error {
	appData, err := b.BRepository.FetchNoShowData(req.CmApplicationID)
	if err != nil {
		return err
	}
	if appData.HtThApplicationID == 0 {
		return fmt.Errorf("Error: %s", "NoShow可能な期間ではありません。")
	}
	now := time.Now()
	t := time.Date(appData.CanceledDt.Year(), appData.CanceledDt.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 2, -1)
	if now.After(t) {
		return fmt.Errorf("Error: %s", "NoShow可能日を過ぎています。")
	}
	var noShowFee float32
	if req.NoshowFlg == true {
		noShowFee = appData.TotalPayInTax
	}
	return b.BRepository.UpdateNoShow(appData.HtThApplicationID, req.NoshowFlg, noShowFee)
}

// SuggestCancelFee キャンセル料の入力補助
func SuggestCancelFee(bookingRooms []booking.HtThBookingRooms, appData *booking.DetailApplicationDBOutput) (float32, error) {
	var res float32

	cancelFlg := appData.CancelFlg
	if cancelFlg == true {
		res = 0
		return res, nil
	}

	// 返金不可の場合は満額を返す
	totalPayInTax := appData.TotalPayInTax
	refundable := bookingRooms[0].Refundable
	if refundable == false {
		res = totalPayInTax
		return res, nil
	}

	// キャンセルペナルティーの取得
	// 部屋が複数の場合も、キャンセルペナルティーは同一
	cancelPenalties := bookingRooms[0].CancelPenalties
	// キャンセルポリシーの取得
	jsonData := []byte(cancelPenalties)
	var cancelPolicy []booking.CancelPolicy
	if err := json.Unmarshal(jsonData, &cancelPolicy); err != nil {
		return res, err
	}
	if len(cancelPolicy) == 0 {
		res = 0
		return res, nil
	}

	t := time.Now()
	timeNow := t.Format(TimeFormat)
	for i := 0; i < len(cancelPolicy); i++ {
		// 該当するキャンセルポリシーがある場合
		if (cancelPolicy[i].Start <= timeNow) && ((timeNow <= cancelPolicy[i].End) || (cancelPolicy[i].End == "")) {
			percent := strings.TrimRight(cancelPolicy[i].Percent, "%")
			rate, err := strconv.Atoi(percent)
			if err != nil {
				return res, err
			}
			res := float32(math.Ceil(float64(totalPayInTax * float32(rate) / 100)))
			return res, nil
		}
	}
	// 該当するキャンセルポリシーが無い場合
	res = 0
	return res, nil
}
