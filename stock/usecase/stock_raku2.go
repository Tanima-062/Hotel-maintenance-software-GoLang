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

// stockRaku2Usecase らく通在庫関連usecase
type stockRaku2Usecase struct {
	SRaku2Repository     stock.IStockRaku2Repository
	RRaku2Repository     room.IRoomRaku2Repository
	PlanRaku2Repository  plan.IPlanRaku2Repository
	PriceRaku2Repository price.IPriceRaku2Repository
}

func (s *stockRaku2Usecase) UpdateBulk(request []stock.StockData) error {
	//TODO implement me
	panic("implement me")
}

// NewStockRaku2Usecase インスタンス生成
func NewStockRaku2Usecase(db *gorm.DB) stock.IStockUsecase {
	return &stockRaku2Usecase{
		SRaku2Repository:     sInfra.NewStockRaku2Repository(db),
		RRaku2Repository:     rInfra.NewRoomRaku2Repository(db),
		PlanRaku2Repository:  planInfra.NewPlanRaku2Repository(db),
		PriceRaku2Repository: priceInfra.NewPriceRaku2Repository(db),
	}
}

// FetchCalendar 在庫料金カレンダー情報取得
func (s *stockRaku2Usecase) FetchCalendar(hmUser account.HtTmHotelManager, request stock.CalendarInput) (*[]stock.CalendarOutput, error) {
	response := []stock.CalendarOutput{}
	startDate := request.BaseDate
	t, _ := time.Parse("2006-01-02", request.BaseDate)
	endDate := t.AddDate(0, 0, 14).Format("2006-01-02")

	// planとroomを並行で取得
	roomCh := make(chan []room.HtTmRoomTypeRaku2s)
	planCh := make(chan []plan.HtTmPlanRaku2s)
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
	stockCh := make(chan []stock.HtTmStockRaku2s)
	priceCh := make(chan []price.HtTmPriceRaku2s)
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
				p := stockAndPrices[baseDate]
				numberOfPeople, _ := strconv.Atoi(priceData.RateTypeCode)
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
func (s *stockRaku2Usecase) UpdateStopSales(request *stock.StopSalesInput) error {
	// トランザクション生成
	tx, txErr := s.SRaku2Repository.TxStart()
	if txErr != nil {
		return txErr
	}
	stockTxRepo := sInfra.NewStockRaku2Repository(tx)
	if err := stockTxRepo.UpdateStopSalesByRoomTypeIDList(request.RoomTypeIDs, request.UseDate, request.IsStopSales); err != nil {
		s.SRaku2Repository.TxRollback(tx)
		return err
	}

	// コミットとロールバック
	if err := s.SRaku2Repository.TxCommit(tx); err != nil {
		s.SRaku2Repository.TxRollback(tx)
		return err
	}
	return nil
}

// FetchAll 在庫情報取得
func (s *stockRaku2Usecase) FetchAll(request *stock.ListInput) (*[]stock.ListOutput, error) {
	response := []stock.ListOutput{}
	startDate := request.BaseDate
	t, _ := time.Parse("2006-01-02", request.BaseDate)
	endDate := t.AddDate(0, 0, 14).Format("2006-01-02")

	rooms, rErr := s.RRaku2Repository.FetchRoomsByPropertyID(room.ListInput{PropertyID: request.PropertyID})
	if rErr != nil {
		return &response, rErr
	}
	roomTypeIDList := []int64{}
	for _, v := range rooms {
		roomTypeIDList = append(roomTypeIDList, v.RoomTypeID)
	}

	stocks, sErr := s.SRaku2Repository.FetchAllByRoomTypeIDList(roomTypeIDList, startDate, endDate)

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
func (s *stockRaku2Usecase) Save(request *[]stock.SaveInput) error {
	inputData := []stock.HtTmStockRaku2s{}
	// トランザクション生成
	tx, txErr := s.SRaku2Repository.TxStart()
	if txErr != nil {
		return txErr
	}

	stockTxRepo := sInfra.NewStockRaku2Repository(tx)
	roomTxRepo := rInfra.NewRoomRaku2Repository(tx)

	for _, roomData := range *request {
		fetchedRoomData, rErr := roomTxRepo.FetchRoomByRoomTypeID(roomData.RoomTypeID)
		if rErr != nil {
			s.SRaku2Repository.TxRollback(tx)
			return rErr
		}
		for useDate, stockData := range roomData.Stocks {
			parsedUseDate, _ := time.Parse("2006-01-02", useDate)
			inputData = append(inputData, stock.HtTmStockRaku2s{
				StockTable: stock.StockTable{
					RoomTypeID:  roomData.RoomTypeID,
					UseDate:     parsedUseDate,
					RoomCount:   stockData.RoomCount,
					IsStopSales: fetchedRoomData.IsStopSales,
				},
			})
		}
	}
	if err := stockTxRepo.UpsertStocks(inputData); err != nil {
		s.SRaku2Repository.TxRollback(tx)
		return err
	}

	// コミットとロールバック
	if err := s.SRaku2Repository.TxCommit(tx); err != nil {
		s.SRaku2Repository.TxRollback(tx)
		return err
	}
	return nil
}

func (s *stockRaku2Usecase) fetchRooms(ch chan<- []room.HtTmRoomTypeRaku2s, propertyID int64) {
	rooms, roomErr := s.RRaku2Repository.FetchRoomsByPropertyID(room.ListInput{PropertyID: propertyID})
	if roomErr != nil {
		ch <- []room.HtTmRoomTypeRaku2s{}
	}
	ch <- rooms
}

func (s *stockRaku2Usecase) fetchPlans(ch chan<- []plan.HtTmPlanRaku2s, propertyID int64) {
	plans, planErr := s.PlanRaku2Repository.FetchAllByPropertyID(plan.ListInput{PropertyID: propertyID})
	if planErr != nil {
		ch <- []plan.HtTmPlanRaku2s{}
	}
	ch <- plans
}

func (s *stockRaku2Usecase) fetchStocks(ch chan<- []stock.HtTmStockRaku2s, roomTypeIDList []int64, startDate string, endDate string) {
	stocks, stockErr := s.SRaku2Repository.FetchAllByRoomTypeIDList(roomTypeIDList, startDate, endDate)
	if stockErr != nil {
		ch <- []stock.HtTmStockRaku2s{}
	}
	ch <- stocks
}

func (s *stockRaku2Usecase) fetchPrices(ch chan<- []price.HtTmPriceRaku2s, planIDList []int64, startDate string, endDate string) {
	prices, priceErr := s.PriceRaku2Repository.FetchAllByPlanIDList(planIDList, startDate, endDate)
	if priceErr != nil {
		ch <- []price.HtTmPriceRaku2s{}
	}
	ch <- prices
}

func (s *stockRaku2Usecase) fetchBookings(ch chan<- []stock.BookingCount, planIDList []int64, startDate string, endDate string) {
	bookings, bookingErr := s.SRaku2Repository.FetchAllBookingsByPlanIDList(planIDList, startDate, endDate)
	if bookingErr != nil {
		ch <- []stock.BookingCount{}
	}
	ch <- bookings
}
