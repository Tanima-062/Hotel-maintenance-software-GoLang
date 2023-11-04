package usecase

import (
	planInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	priceInfra "github.com/Adventureinc/hotel-hm-api/src/price/infra"
	"gorm.io/gorm"
	"strconv"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	sInfra "github.com/Adventureinc/hotel-hm-api/src/stock/infra"
	"github.com/labstack/gommon/log"
)

// stockTlUsecase Usecase related to Tl purchase inventory
type stockTlUsecase struct {
	STlRepository     stock.IStockTlRepository
	RTlRepository     room.IRoomTlRepository
	PlanTlRepository  plan.IPlanTlRepository
	PriceTlRepository price.IPriceTlRepository
}

// NewStockTLUsecase instantiation
func NewStockTlUsecase(db *gorm.DB) *stockTlUsecase {
	return &stockTlUsecase{
		STlRepository:     sInfra.NewStockTlRepository(db),
		RTlRepository:     rInfra.NewRoomTlRepository(db),
		PlanTlRepository:  planInfra.NewPlanTlRepository(db),
		PriceTlRepository: priceInfra.NewPriceTlRepository(db),
	}
}

func (s *stockTlUsecase) UpdateBulk(request []stock.StockData) error {
	// transaction generation
	tx, txErr := s.STlRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	//Bulk data insert from request
	for _, requestData := range request {
		roomType, _ := s.RTlRepository.FetchRoomTypeIdByRoomTypeCode(requestData.PropertyID, requestData.RoomTypeCode)
		//if no room type data not found
		if (roomType != room.HtTmRoomTypeTls{}) {
			roomType.StockSettingStart = requestData.StockSettingStart
			roomType.StockSettingEnd = requestData.StockSettingEnd
			roomType.IsSettingStockYearRound = requestData.IsSettingStockYearRound
			// Update `RoomTypeTls`
			if err := s.RTlRepository.UpdateRoomBulkTl(&roomType); err != nil {
				log.Error(err)
				continue
			}

			for useDate, stockData := range requestData.Stocks {
				// fetch stock detail
				bookingData, _ := s.STlRepository.FetchBookingCountByRoomTypeId(roomType.RoomTypeID, useDate)
				//check if booking data found
				if (bookingData != stock.StockTable{}) {
					// update stock detail
					if err := s.STlRepository.UpdateStocksBulk(roomType.RoomTypeID, useDate, int64(stockData.Stock), int64(bookingData.BookingCount), stockData.IsStopSales); err != nil {
						s.STlRepository.TxRollback(tx)
						return err
					}
				} else {
					var stockInputData []stock.HtTmStockTls
					parsedUseDate, _ := time.Parse("2006-01-02", useDate)
					stockInputData = append(stockInputData, stock.HtTmStockTls{
						StockTable: stock.StockTable{
							RoomTypeID:  roomType.RoomTypeID,
							UseDate:     parsedUseDate,
							RoomCount:   stockData.Stock + bookingData.BookingCount,
							Stock:       stockData.Stock,
							IsStopSales: stockData.IsStopSales,
						},
					})
					if err := s.STlRepository.CreateStocks(stockInputData); err != nil {
						s.STlRepository.TxRollback(tx)
						return err
					}
				}
			}
		}
	}
	// commit and rollback
	if err := s.STlRepository.TxCommit(tx); err != nil {
		s.STlRepository.TxRollback(tx)
		return err
	}
	return nil
}

// FetchCalendar Get inventory price calendar information
func (s *stockTlUsecase) FetchCalendar(hmUser account.HtTmHotelManager, request stock.CalendarInput) (*[]stock.CalendarOutput, error) {
	response := []stock.CalendarOutput{}
	startDate := request.BaseDate
	t, _ := time.Parse("2006-01-02", request.BaseDate)
	endDate := t.AddDate(0, 0, 14).Format("2006-01-02")

	// Get plan and room in parallel
	roomCh := make(chan []room.HtTmRoomTypeTls)
	planCh := make(chan []price.HtTmPlanTls)
	go s.fetchRooms(roomCh, hmUser.PropertyID)
	go s.fetchPlans(planCh, hmUser.PropertyID)
	rooms, plans := <-roomCh, <-planCh

	// Create roomTypeID list and planID list for stock and price data acquisition
	roomTypeIDList := []int64{}
	for _, value := range rooms {
		roomTypeIDList = append(roomTypeIDList, value.RoomTypeID)
	}
	planIDList := []int64{}
	for _, value := range plans {
		planIDList = append(planIDList, value.PlanID)
	}

	// Get inventory, price and number of sales per plan
	stockCh := make(chan []stock.HtTmStockTls)
	priceCh := make(chan []price.HtTmPriceTls)
	bookingCh := make(chan []stock.BookingCount)
	go s.fetchStocks(stockCh, roomTypeIDList, startDate, endDate)
	go s.fetchPrices(priceCh, planIDList, startDate, endDate)
	go s.fetchBookings(bookingCh, planIDList, startDate, endDate)
	stocks, prices, bookings := <-stockCh, <-priceCh, <-bookingCh

	for _, roomData := range rooms {
		// Set room information
		calendarRecord := &stock.CalendarOutput{}
		calendarRecord.RoomTypeID = roomData.RoomTypeID
		calendarRecord.RoomName = roomData.Name
		calendarRecord.IsStopSales = roomData.IsStopSales
		calendarRecord.StockSettingStart = roomData.StockSettingStart.Format("2006-01-02")
		calendarRecord.StockSettingEnd = roomData.StockSettingEnd.Format("2006-01-02")
		calendarRecord.OcuMin = roomData.OcuMin
		calendarRecord.OcuMax = roomData.OcuMax
		calendarRecord.IsSettingStockYearRound = roomData.IsSettingStockYearRound

		// Set stock information
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

		// Set a plan linked to the room
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

			// Set the price data linked to the plan
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

			// Set sales data linked to the plan
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

// UpdateStopSales Inventory discontinued update
func (s *stockTlUsecase) UpdateStopSales(request *stock.StopSalesInput) error {
	// transaction generation
	tx, txErr := s.STlRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	stockTxRepo := sInfra.NewStockTlRepository(tx)
	if err := stockTxRepo.UpdateStopSalesByRoomTypeIDList(request.RoomTypeIDs, request.UseDate, request.IsStopSales); err != nil {
		s.STlRepository.TxRollback(tx)
		return err
	}

	// commit and rollback
	if err := s.STlRepository.TxCommit(tx); err != nil {
		s.STlRepository.TxRollback(tx)
		return err
	}
	return nil
}

// FetchAll Get inventory information
func (s *stockTlUsecase) FetchAll(request *stock.ListInput) (*[]stock.ListOutput, error) {
	response := []stock.ListOutput{}
	startDate := request.BaseDate
	t, _ := time.Parse("2006-01-02", request.BaseDate)
	endDate := t.AddDate(0, 0, 14).Format("2006-01-02")

	rooms, rErr := s.RTlRepository.FetchRoomsByPropertyID(room.ListInput{PropertyID: request.PropertyID})
	if rErr != nil {
		return &response, rErr
	}
	roomTypeIDList := []int64{}
	for _, v := range rooms {
		roomTypeIDList = append(roomTypeIDList, v.RoomTypeID)
	}

	stocks, sErr := s.STlRepository.FetchAllByRoomTypeIDList(roomTypeIDList, startDate, endDate)

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

// Save Inventory creation/update
func (s *stockTlUsecase) Save(request *[]stock.SaveInput) error {
	inputData := []stock.HtTmStockTls{}
	// transaction generation
	tx, txErr := s.STlRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	stockTxRepo := sInfra.NewStockTlRepository(tx)
	roomTxRepo := rInfra.NewRoomTlRepository(tx)

	for _, roomData := range *request {
		fetchedRoomData, rErr := roomTxRepo.FetchRoomByRoomTypeID(roomData.RoomTypeID)
		if rErr != nil {
			s.STlRepository.TxRollback(tx)
			return rErr
		}
		for useDate, stockData := range roomData.Stocks {
			parsedUseDate, _ := time.Parse("2006-01-02", useDate)
			inputData = append(inputData, stock.HtTmStockTls{
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
		s.STlRepository.TxRollback(tx)
		return err
	}

	// commit and rollback
	if err := s.STlRepository.TxCommit(tx); err != nil {
		s.STlRepository.TxRollback(tx)
		return err
	}
	return nil
}

func (s *stockTlUsecase) fetchRooms(ch chan<- []room.HtTmRoomTypeTls, propertyID int64) {
	rooms, roomErr := s.RTlRepository.FetchRoomsByPropertyID(room.ListInput{PropertyID: propertyID})
	if roomErr != nil {
		ch <- []room.HtTmRoomTypeTls{}
	}
	ch <- rooms
}

func (s *stockTlUsecase) fetchPlans(ch chan<- []price.HtTmPlanTls, propertyID int64) {
	plans, planErr := s.PlanTlRepository.FetchAllByPropertyID(plan.ListInput{PropertyID: propertyID})
	if planErr != nil {
		ch <- []price.HtTmPlanTls{}
	}
	ch <- plans
}

func (s *stockTlUsecase) fetchStocks(ch chan<- []stock.HtTmStockTls, roomTypeIDList []int64, startDate string, endDate string) {
	stocks, stockErr := s.STlRepository.FetchAllByRoomTypeIDList(roomTypeIDList, startDate, endDate)
	if stockErr != nil {
		ch <- []stock.HtTmStockTls{}
	}
	ch <- stocks
}

func (s *stockTlUsecase) fetchPrices(ch chan<- []price.HtTmPriceTls, planIDList []int64, startDate string, endDate string) {
	prices, priceErr := s.PriceTlRepository.FetchAllByPlanIDList(planIDList, startDate, endDate)
	if priceErr != nil {
		ch <- []price.HtTmPriceTls{}
	}
	ch <- prices
}

func (s *stockTlUsecase) fetchBookings(ch chan<- []stock.BookingCount, planIDList []int64, startDate string, endDate string) {
	bookings, bookingErr := s.STlRepository.FetchAllBookingsByPlanIDList(planIDList, startDate, endDate)
	if bookingErr != nil {
		ch <- []stock.BookingCount{}
	}
	ch <- bookings
}
