package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	planInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	priceInfra "github.com/Adventureinc/hotel-hm-api/src/price/infra"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	sInfra "github.com/Adventureinc/hotel-hm-api/src/stock/infra"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"time"
)

// StockTemaUsecase
type StockTemaUsecase struct {
	STemaRepository     stock.IStockTemaRepository
	RTemaRepository     room.IRoomTemaRepository
	PlanTemaRepository  plan.IPlanTemaRepository
	PriceTemaRepository price.IPriceTemaRepository
}

// NewStockTemaUsecase creates a new StockTemaUsecase
func NewStockTemaUsecase(db *gorm.DB) stock.IStockTemaUsecase {
	return &StockTemaUsecase{
		STemaRepository:     sInfra.NewStockTemaRepository(db),
		RTemaRepository:     rInfra.NewRoomTemaRepository(db),
		PlanTemaRepository:  planInfra.NewPlanTemaRepository(db),
		PriceTemaRepository: priceInfra.NewPriceTemaRepository(db),
	}
}

// FetchCalendar returns the calendar for the specified stock account
func (s *StockTemaUsecase) FetchCalendar(hmUser account.HtTmHotelManager, request stock.CalendarInput) (*[]stock.CalendarOutputTema, error) {
	response := []stock.CalendarOutputTema{}
	startDate := request.BaseDate
	t, _ := time.Parse("2006-01-02", request.BaseDate)
	endDate := t.AddDate(0, 0, 14).Format("2006-01-02")

	// Get plan and room in parallel
	roomCh := make(chan []room.HtTmRoomTypeTemas)
	planCh := make(chan []price.HtTmPlanTemas)

	go s.fetchRooms(roomCh, hmUser.PropertyID)
	go s.fetchPlans(planCh, hmUser.PropertyID)
	rooms, plans := <-roomCh, <-planCh

	// Create roomTypeID list and planID list for stock and price data acquisition
	roomTypeCodeList := []int64{}
	for _, value := range rooms {
		roomTypeCode, _ := strconv.Atoi(value.RoomTypeTema.RoomTypeCode)
		roomTypeCodeList = append(roomTypeCodeList, int64(roomTypeCode))
	}

	planCodeList := []int64{}
	for _, value := range plans {
		planCodeList = append(planCodeList, value.TemaPlanTable.PackagePlanCode)
	}

	// Get inventory, price and number of sales per plan
	stockCh := make(chan []stock.HtTmStockTemas)
	priceCh := make(chan []price.HtTmPriceTemas)
	bookingCh := make(chan []stock.BookingCount)
	go s.fetchStocks(stockCh, roomTypeCodeList, startDate, endDate)
	go s.fetchPrices(priceCh, planCodeList, startDate, endDate)
	go s.fetchBookings(bookingCh, planCodeList, startDate, endDate)
	stocks := <-stockCh
	prices := <-priceCh
	bookings := <-bookingCh

	for _, roomData := range rooms {
		// Set room information
		calendarRecord := &stock.CalendarOutputTema{}
		calendarRecord.RoomTypeID = roomData.RoomTypeTema.RoomTypeID
		calendarRecord.RoomName = roomData.RoomTypeTema.Name
		calendarRecord.IsStopSales = roomData.RoomTypeTema.IsStopSales
		calendarRecord.OcuMin = roomData.RoomTypeTema.OcuMin
		calendarRecord.OcuMax = roomData.RoomTypeTema.OcuMax

		// Set stock information
		stockList := []stock.CalendarStockTema{}
		roomTypeCode, _ := strconv.Atoi(roomData.RoomTypeTema.RoomTypeCode)
		for _, stockData := range stocks {
			if stockData.StockTableTema.RoomTypeCode == int64(roomTypeCode) {
				stockList = append(stockList, stock.CalendarStockTema{
					RoomTypeCode: stockData.StockTableTema.RoomTypeCode,
					AriDate:      stockData.StockTableTema.AriDate.Format("2006-01-02"),
					Stock:        stockData.StockTableTema.Stock,
					IsStopSales:  stockData.StockTableTema.Disable,
				})
			}
		}
		calendarRecord.Stocks = stockList

		// Set a plan linked to the room
		for _, planData := range plans {
			if roomData.RoomTypeTema.RoomTypeID != planData.TemaPlanTable.RoomTypeID {
				continue
			}
			planRecord := &stock.CalendarPlanTema{
				PlanID:   planData.TemaPlanTable.PlanID,
				PlanName: planData.TemaPlanTable.PlanName,
				Disable:  planData.TemaPlanTable.Available,
			}
			stockAndPrices := map[string]stock.CalendarPriceTema{}

			// Set the price data linked to the plan
			for _, priceData := range prices {
				if priceData.PriceTemaTable.PackagePlanCode != planData.TemaPlanTable.PackagePlanCode {
					continue
				}
				baseDate := priceData.PriceDate.Format("2006-01-02")
				p := stockAndPrices[baseDate]

				priceType := priceData.PriceTemaTable.TemaPriceType
				pt := reflect.ValueOf(&priceType).Elem().Type()
				if roomTypeCode == priceData.PriceTemaTable.RoomTypeCode {
					for i := 0; i < 6; i++ {
						field := pt.Field(i)
						rv := reflect.ValueOf(&priceType)
						value := reflect.Indirect(rv).FieldByName(field.Name)
						p.Prices = append(p.Prices, stock.Price{
							Type:  "0" + strconv.Itoa(i+1),
							Price: value.Int(),
						})
					}
				}
				stockAndPrices[baseDate] = p
			}

			// Set sales data linked to the plan
			for _, bookingData := range bookings {
				if bookingData.PlanID != planData.TemaPlanTable.PlanID {
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

// UpdateBulkTema process for stock
func (s *StockTemaUsecase) UpdateBulkTema(request []stock.StockDataTema) error {
	// transaction generation
	tx, txErr := s.STemaRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	//Bulk data insert from request
	for _, requestData := range request {
		roomType, _ := s.STemaRepository.FetchRoomTypeIdByRoomTypeCode(requestData.PropertyID, requestData.RoomTypeCode)
		//if no room type data not found
		if (roomType != room.HtTmRoomTypeTemas{}) {
			// Update `RoomTypeTemas`
			if err := s.STemaRepository.UpdateRoomBulkTema(&roomType); err != nil {
				log.Error(err)
				continue
			}

			for ariDate, stockData := range requestData.Stocks {
				// fetch stock detail
				bookingData, _ := s.STemaRepository.FetchBookingCountByRoomTypeId(roomType.RoomTypeCode, ariDate)
				//check if booking data found
				if (bookingData != stock.StockTableTema{}) {
					// update stock detail
					if err := s.STemaRepository.UpdateStocksBulk(roomType.RoomTypeCode, ariDate, int64(stockData.Stock), stockData.Disable); err != nil {
						s.STemaRepository.TxRollback(tx)
						return err
					}
				} else {
					var stockInputData []stock.HtTmStockTemas
					parsedUseDate, _ := time.Parse("2006-01-02", ariDate)
					roomTypeCode, _ := strconv.ParseInt(roomType.RoomTypeCode, 10, 64)
					stockInputData = append(stockInputData, stock.HtTmStockTemas{
						StockTableTema: stock.StockTableTema{
							PropertyID:   roomType.PropertyID,
							RoomTypeCode: roomTypeCode,
							AriDate:      parsedUseDate,
							Stock:        stockData.Stock,
							Disable:      stockData.Disable,
						},
					})
					// create new stock details
					if err := s.STemaRepository.CreateStocks(stockInputData); err != nil {
						s.STemaRepository.TxRollback(tx)
						return err
					}
				}
			}
		}
	}
	// commit and rollback
	if err := s.STemaRepository.TxCommit(tx); err != nil {
		s.STemaRepository.TxRollback(tx)
		return err
	}
	return nil
}

// fetchRooms fetch room information
func (s *StockTemaUsecase) fetchRooms(ch chan<- []room.HtTmRoomTypeTemas, propertyID int64) {
	rooms, roomErr := s.RTemaRepository.FetchRoomsByPropertyID(room.ListInput{PropertyID: propertyID})
	if roomErr != nil {
		ch <- []room.HtTmRoomTypeTemas{}
	}
	ch <- rooms
}

// fetchPlans fetch plan information
func (s *StockTemaUsecase) fetchPlans(ch chan<- []price.HtTmPlanTemas, propertyID int64) {
	plans, planErr := s.PlanTemaRepository.FetchAllByPropertyID(plan.ListInput{PropertyID: propertyID})
	if planErr != nil {
		ch <- []price.HtTmPlanTemas{}
	}
	ch <- plans
}

// fetchStocks fetch stock information
func (s *StockTemaUsecase) fetchStocks(ch chan<- []stock.HtTmStockTemas, roomTypeCodeList []int64, startDate string, endDate string) {
	stocks, stockErr := s.STemaRepository.FetchAllByRoomTypeCodeList(roomTypeCodeList, startDate, endDate)

	if stockErr != nil {
		ch <- []stock.HtTmStockTemas{}
	}
	ch <- stocks
}

// fetchPrices fetch price information
func (s *StockTemaUsecase) fetchPrices(ch chan<- []price.HtTmPriceTemas, planCodeList []int64, startDate string, endDate string) {
	prices, priceErr := s.PriceTemaRepository.FetchAllByPlanCodeList(planCodeList, startDate, endDate)
	if priceErr != nil {
		ch <- []price.HtTmPriceTemas{}
	}
	ch <- prices
}

// fetchBookings fetch booking information
func (s *StockTemaUsecase) fetchBookings(ch chan<- []stock.BookingCount, planIDList []int64, startDate string, endDate string) {
	bookings, bookingErr := s.STemaRepository.FetchAllBookingsByPlanIDList(planIDList, startDate, endDate)
	if bookingErr != nil {
		ch <- []stock.BookingCount{}
	}
	ch <- bookings
}
