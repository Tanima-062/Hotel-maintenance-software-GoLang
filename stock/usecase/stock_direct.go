package usecase

import (
	"strconv"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	planInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	priceInfra "github.com/Adventureinc/hotel-hm-api/src/price/infra"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	sInfra "github.com/Adventureinc/hotel-hm-api/src/stock/infra"
	"gorm.io/gorm"
)

// stockDirectUsecase 直仕入れ在庫関連usecase
type stockDirectUsecase struct {
	SDirectRepository     stock.IStockDirectRepository
	RDirectRepository     room.IRoomDirectRepository
	PlanDirectRepository  plan.IPlanDirectRepository
	PriceDirectRepository price.IPriceDirectRepository
}

func (s *stockDirectUsecase) UpdateBulk(request []stock.StockData) error {
	//TODO implement me
	panic("implement me")
}

// NewStockDirectUsecase インスタンス生成
func NewStockDirectUsecase(db *gorm.DB) stock.IStockUsecase {
	return &stockDirectUsecase{
		SDirectRepository:     sInfra.NewStockDirectRepository(db),
		RDirectRepository:     rInfra.NewRoomDirectRepository(db),
		PlanDirectRepository:  planInfra.NewPlanDirectRepository(db),
		PriceDirectRepository: priceInfra.NewPriceDirectRepository(db),
	}
}

// FetchCalendar 在庫料金カレンダー情報取得
func (s *stockDirectUsecase) FetchCalendar(hmUser account.HtTmHotelManager, request stock.CalendarInput) (*[]stock.CalendarOutput, error) {
	response := []stock.CalendarOutput{}
	startDate := request.BaseDate
	t, _ := time.Parse("2006-01-02", request.BaseDate)
	endDate := t.AddDate(0, 0, 14).Format("2006-01-02")

	// planとroomを並行で取得
	roomCh := make(chan []room.HtTmRoomTypeDirects)
	planCh := make(chan []plan.HtTmPlanDirects)
	go s.fetchRooms(roomCh, hmUser.PropertyID)
	go s.fetchPlans(planCh, hmUser.PropertyID)
	rooms, plans := <-roomCh, <-planCh

	// stockとpriceデータ取得用に、roomTypeID一覧とplanID一覧を作成
	roomTypeIDList := []int64{}
	for _, value := range rooms {
		roomTypeIDList = append(roomTypeIDList, value.RoomTypeID)
	}
	planIDList := []int64{}
	for _, value := range plans {
		planIDList = append(planIDList, value.PlanID)
	}

	// 在庫と料金とプランごとの販売数を取得
	stockCh := make(chan []stock.HtTmStockDirects)
	priceCh := make(chan []price.HtTmPriceDirects)
	bookingCh := make(chan []stock.BookingCount)
	go s.fetchStocks(stockCh, roomTypeIDList, startDate, endDate)
	go s.fetchPrices(priceCh, planIDList, startDate, endDate)
	go s.fetchBookings(bookingCh, planIDList, startDate, endDate)
	stocks, prices, bookings := <-stockCh, <-priceCh, <-bookingCh

	for _, roomData := range rooms {
		// 部屋情報をセット
		calendarRecord := &stock.CalendarOutput{}
		calendarRecord.RoomTypeID = roomData.RoomTypeID
		calendarRecord.RoomName = roomData.Name
		calendarRecord.IsStopSales = roomData.IsStopSales
		calendarRecord.StockSettingStart = roomData.StockSettingStart.Format("2006-01-02")
		calendarRecord.StockSettingEnd = roomData.StockSettingEnd.Format("2006-01-02")
		calendarRecord.OcuMin = roomData.OcuMin
		calendarRecord.OcuMax = roomData.OcuMax
		calendarRecord.IsSettingStockYearRound = roomData.IsSettingStockYearRound

		// 在庫情報をセット
		stockList := []stock.CalendarStock{}
		for _, stockData := range stocks {
			if stockData.RoomTypeID == roomData.RoomTypeID {
				stockList = append(stockList, stock.CalendarStock{
					RoomTypeID:   stockData.RoomTypeID,
					UseDate:      stockData.UseDate.Format("2006-01-02"),
					RoomCount:    stockData.RoomCount,
					BookingCount: stockData.BookingCount,
					Stock:        stockData.Stock,
					IsStopSales:  stockData.IsStopSales,
				})
			}
		}
		calendarRecord.Stocks = stockList

		// 部屋に紐づくプランをセット
		for _, planData := range plans {
			if roomData.RoomTypeID != planData.RoomTypeID {
				continue
			}
			planRecord := &stock.CalendarPlan{
				PlanID:      planData.PlanID,
				PlanName:    planData.Name,
				IsStopSales: planData.IsStopSales,
			}
			stockAndPrices := map[string]stock.CalendarPrice{}

			// プランに紐づく料金データをセット
			for _, priceData := range prices {
				if priceData.PlanID != planData.PlanID {
					continue
				}
				baseDate := priceData.UseDate.Format("2006-01-02")
				numberOfPeople, _ := strconv.Atoi(priceData.RateTypeCode)
				p := stockAndPrices[baseDate]
				p.Prices = append(p.Prices, price.Price{
					Type:  priceData.RateTypeCode,
					Price: priceData.PriceInTax / numberOfPeople,
				})
				stockAndPrices[baseDate] = p
			}

			// プランに紐づく販売データをセット
			for _, bookingData := range bookings {
				if bookingData.PlanID != planData.PlanID {
					continue
				}
				t, _ := time.Parse(time.RFC3339, bookingData.UseDate)
				useDate := t.Format("2006-01-02")
				b := stockAndPrices[useDate]
				b.BookingCount = bookingData.BookingCount
				stockAndPrices[useDate] = b
			}
			planRecord.StockAndPrices = stockAndPrices
			calendarRecord.Plans = append(calendarRecord.Plans, *planRecord)
		}
		response = append(response, *calendarRecord)
	}
	return &response, nil
}

// UpdateStopSales 在庫の売止更新
func (s *stockDirectUsecase) UpdateStopSales(request *stock.StopSalesInput) error {
	// トランザクション生成
	tx, txErr := s.SDirectRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	stockTxRepo := sInfra.NewStockDirectRepository(tx)
	if err := stockTxRepo.UpdateStopSalesByRoomTypeIDList(request.RoomTypeIDs, request.UseDate, request.IsStopSales); err != nil {
		s.SDirectRepository.TxRollback(tx)
		return err
	}

	// コミットとロールバック
	if err := s.SDirectRepository.TxCommit(tx); err != nil {
		s.SDirectRepository.TxRollback(tx)
		return err
	}
	return nil
}

// FetchAll 在庫情報取得
func (s *stockDirectUsecase) FetchAll(request *stock.ListInput) (*[]stock.ListOutput, error) {
	response := []stock.ListOutput{}
	startDate := request.BaseDate
	t, _ := time.Parse("2006-01-02", request.BaseDate)
	endDate := t.AddDate(0, 0, 14).Format("2006-01-02")

	rooms, rErr := s.RDirectRepository.FetchRoomsByPropertyID(room.ListInput{PropertyID: request.PropertyID})
	if rErr != nil {
		return &response, rErr
	}
	roomTypeIDList := []int64{}
	for _, v := range rooms {
		roomTypeIDList = append(roomTypeIDList, v.RoomTypeID)
	}

	stocks, sErr := s.SDirectRepository.FetchAllByRoomTypeIDList(roomTypeIDList, startDate, endDate)

	if sErr != nil {
		return &response, sErr
	}
	for _, roomData := range rooms {
		tempListOutput := stock.ListOutput{
			RoomTypeID:              roomData.RoomTypeID,
			Name:                    roomData.Name,
			StockSettingStart:       roomData.StockSettingStart.Format("2006-01-02"),
			StockSettingEnd:         roomData.StockSettingEnd.Format("2006-01-02"),
			IsSettingStockYearRound: roomData.IsSettingStockYearRound,
			RoomCount:               roomData.RoomCount,
		}
		tempStocks := map[string]stock.ListStockOutput{}
		for _, stockData := range stocks {
			if stockData.RoomTypeID != roomData.RoomTypeID {
				continue
			}

			tempStocks[stockData.UseDate.Format("2006-01-02")] = stock.ListStockOutput{
				RoomCount:    stockData.RoomCount,
				BookingCount: stockData.BookingCount,
			}
		}
		tempListOutput.Stocks = tempStocks
		response = append(response, tempListOutput)
	}
	return &response, nil
}

// Save 在庫作成・更新
func (s *stockDirectUsecase) Save(request *[]stock.SaveInput) error {
	inputData := []stock.HtTmStockDirects{}
	// トランザクション生成
	tx, txErr := s.SDirectRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	stockTxRepo := sInfra.NewStockDirectRepository(tx)
	roomTxRepo := rInfra.NewRoomDirectRepository(tx)

	roomTypeIdList := []int64{}
	for _, roomData := range *request {
		roomTypeIdList = append(roomTypeIdList, roomData.RoomTypeID)
	}

	// roomTypeIdListで部屋情報取得
	rooms, _ := roomTxRepo.FetchRoomListByRoomTypeID(roomTypeIdList)

	// roomTypeIdListで在庫取得
	existStocks, _ := stockTxRepo.FetchStocksByRoomTypeIDList(roomTypeIdList)

	for _, roomData := range *request {
		var isStopSales bool = false
		for _, room := range rooms {
			if roomData.RoomTypeID == room.RoomTypeID {
				isStopSales = room.IsStopSales
				break
			}
		}
		for useDate, stockData := range roomData.Stocks {
			var bokkingCount int16 = 0
			parsedUseDate, _ := time.Parse("2006-01-02", useDate)
			for _, existStock := range existStocks {
				if roomData.RoomTypeID == existStock.RoomTypeID && parsedUseDate == existStock.UseDate {
					bokkingCount = existStock.BookingCount
					break
				}
			}
			inputData = append(inputData, stock.HtTmStockDirects{
				StockTable: stock.StockTable{
					RoomTypeID:  roomData.RoomTypeID,
					UseDate:     parsedUseDate,
					RoomCount:   stockData.RoomCount,
					Stock:       stockData.RoomCount - bokkingCount,
					IsStopSales: isStopSales,
				},
			})
		}
	}
	if err := stockTxRepo.UpsertStocks(inputData); err != nil {
		s.SDirectRepository.TxRollback(tx)
		return err
	}

	// コミットとロールバック
	if err := s.SDirectRepository.TxCommit(tx); err != nil {
		s.SDirectRepository.TxRollback(tx)
		return err
	}
	return nil
}

func (s *stockDirectUsecase) fetchRooms(ch chan<- []room.HtTmRoomTypeDirects, propertyID int64) {
	rooms, roomErr := s.RDirectRepository.FetchRoomsByPropertyID(room.ListInput{PropertyID: propertyID})
	if roomErr != nil {
		ch <- []room.HtTmRoomTypeDirects{}
	}
	ch <- rooms
}

func (s *stockDirectUsecase) fetchPlans(ch chan<- []plan.HtTmPlanDirects, propertyID int64) {
	plans, planErr := s.PlanDirectRepository.FetchAllByPropertyID(plan.ListInput{PropertyID: propertyID})
	if planErr != nil {
		ch <- []plan.HtTmPlanDirects{}
	}
	ch <- plans
}

func (s *stockDirectUsecase) fetchStocks(ch chan<- []stock.HtTmStockDirects, roomTypeIDList []int64, startDate string, endDate string) {
	stocks, stockErr := s.SDirectRepository.FetchAllByRoomTypeIDList(roomTypeIDList, startDate, endDate)
	if stockErr != nil {
		ch <- []stock.HtTmStockDirects{}
	}
	ch <- stocks
}

func (s *stockDirectUsecase) fetchPrices(ch chan<- []price.HtTmPriceDirects, planIDList []int64, startDate string, endDate string) {
	prices, priceErr := s.PriceDirectRepository.FetchAllByPlanIDList(planIDList, startDate, endDate)
	if priceErr != nil {
		ch <- []price.HtTmPriceDirects{}
	}
	ch <- prices
}

func (s *stockDirectUsecase) fetchBookings(ch chan<- []stock.BookingCount, planIDList []int64, startDate string, endDate string) {
	bookings, bookingErr := s.SDirectRepository.FetchAllBookingsByPlanIDList(planIDList, startDate, endDate)
	if bookingErr != nil {
		ch <- []stock.BookingCount{}
	}
	ch <- bookings
}
