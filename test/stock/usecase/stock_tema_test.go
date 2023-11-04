package usecase_test

import (
	"errors"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	stockUseCase "github.com/Adventureinc/hotel-hm-api/src/stock/usecase"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var flag = 0

// MockRoomStockBulkUseCase mock implementation
type MockStockTemaBulkUseCase struct {
	mock.Mock
}

func (m MockStockTemaBulkUseCase) TxStart() (*gorm.DB, error) {
	db, _, err := sqlmock.New()
	gorm, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if flag == 2 {
		return nil, errors.New("new err")
	}
	return gorm.Debug(), err
}

// TxCommit mock
func (m MockStockTemaBulkUseCase) TxCommit(tx *gorm.DB) error {

	if flag == 1 {
		return errors.New("new err")
	}
	return nil
}

func (m MockStockTemaBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

func (m MockStockTemaBulkUseCase) FetchRoomTypeIdByRoomTypeCode(propertyID int64, roomTypeCode string) (room.HtTmRoomTypeTemas, error) {
	return room.HtTmRoomTypeTemas{
		RoomTypeTema: room.RoomTypeTema{
			RoomTypeID:   1,
			RoomTypeCode: roomTypeCode,
			PropertyID:   propertyID,
		},
	}, nil
}

func (m MockStockTemaBulkUseCase) UpdateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error {
	if roomTable.RoomTypeCode == "1208101" {
		return errors.New("new err")
	}
	return nil
}

func (m MockStockTemaBulkUseCase) FetchBookingCountByRoomTypeId(roomTypeID string, useDate string) (stock.StockTableTema, error) {

	if roomTypeID == "1208100" || useDate == "2023-07-02" {
		return stock.StockTableTema{}, errors.New("new err")
	}
	if roomTypeID == "1208103" || useDate == "2023-09-02" {
		return stock.StockTableTema{}, errors.New("new err")
	}
	return stock.StockTableTema{
		StockID: 1,
		Stock:   20,
		Disable: true,
	}, nil
}

func (m *MockStockTemaBulkUseCase) UpdateStocksBulk(propertyCode string, ariDate string, stock int64, disable bool) error {
	if propertyCode == "1208102" {
		return errors.New("new err")
	}
	return nil
}

func (m *MockStockTemaBulkUseCase) CreateStocks(inputData []stock.HtTmStockTemas) error {

	var lan = len(inputData)
	if inputData[lan-1].RoomTypeCode == 1208103 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockStockTemaBulkUseCase) FetchAllByRoomTypeCodeList(roomTypeCodeList []int64, startDate string, endDate string) ([]stock.HtTmStockTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockStockTemaBulkUseCase) FetchAllBookingsByPlanIDList(planIDList []int64, startDate string, endDate string) ([]stock.BookingCount, error) {
	//TODO implement me
	panic("implement me")
}

var request = []stock.StockDataTema{
	{
		PropertyID:   1208918,
		RoomTypeCode: "1208010",
		Stocks: map[string]stock.UpdateStockTemaInput{
			"2023-07-01": stock.UpdateStockTemaInput{
				Stock:   12,
				Disable: true,
			},
			"2023-07-02": stock.UpdateStockTemaInput{
				Stock:   15,
				Disable: false,
			},
		},
	},
	{
		PropertyID:   1208100,
		RoomTypeCode: "1208100",
		Stocks: map[string]stock.UpdateStockTemaInput{
			"2023-08-01": stock.UpdateStockTemaInput{
				Stock:   12,
				Disable: true,
			},
			"2023-08-02": stock.UpdateStockTemaInput{
				Stock:   15,
				Disable: false,
			},
		},
	},
	{
		PropertyID:   1208101,
		RoomTypeCode: "1208101",
		Stocks: map[string]stock.UpdateStockTemaInput{
			"2023-09-01": stock.UpdateStockTemaInput{
				Stock:   12,
				Disable: true,
			},
			"2023-09-02": stock.UpdateStockTemaInput{
				Stock:   15,
				Disable: false,
			},
		},
	},
	{
		PropertyID:   1208102,
		RoomTypeCode: "1208102",
		Stocks: map[string]stock.UpdateStockTemaInput{
			"2023-09-01": stock.UpdateStockTemaInput{
				Stock:   12,
				Disable: true,
			},
			"2023-09-02": stock.UpdateStockTemaInput{
				Stock:   15,
				Disable: false,
			},
		},
	},
	{
		PropertyID:   1208103,
		RoomTypeCode: "1208103",
		Stocks: map[string]stock.UpdateStockTemaInput{
			"2023-09-01": stock.UpdateStockTemaInput{
				Stock:   12,
				Disable: true,
			},
			"2023-09-02": stock.UpdateStockTemaInput{
				Stock:   15,
				Disable: false,
			},
		},
	},

	{
		PropertyID:   1208104,
		RoomTypeCode: "1208104",
		Stocks: map[string]stock.UpdateStockTemaInput{
			"2023-09-01": stock.UpdateStockTemaInput{
				Stock:   12,
				Disable: true,
			},
			"2023-09-02": stock.UpdateStockTemaInput{
				Stock:   15,
				Disable: false,
			},
		},
	},
	{
		PropertyID:   1208105,
		RoomTypeCode: "1208105",
		Stocks: map[string]stock.UpdateStockTemaInput{
			"2023-09-01": stock.UpdateStockTemaInput{
				Stock:   12,
				Disable: true,
			},
			"2023-09-02": stock.UpdateStockTemaInput{
				Stock:   15,
				Disable: false,
			},
		},
	},
}

func TemaStockUsecaseUpdateBulkTema(request []stock.StockDataTema) error {
	mockUseCase := new(MockStockTemaBulkUseCase)
	useCases := &stockUseCase.StockTemaUsecase{
		STemaRepository: mockUseCase,
	}

	ret := useCases.UpdateBulkTema(request)
	return ret
}

func TestTemaStockUsecaseUpdateBulkTema(t *testing.T) {

	for _, data := range request {
		var reqArray = []stock.StockDataTema{
			data,
		}
		if data.RoomTypeCode == "1208104" {
			flag = 1
		}
		if data.RoomTypeCode == "1208105" {
			flag = 2
		}
		res := TemaStockUsecaseUpdateBulkTema(reqArray)
		if res == nil {
			assert.Nil(t, res)
			assert.Equal(t, res, nil)
		} else {
			assert.Error(t, res)
			assert.Equal(t, res, errors.New("new err"))
		}

	}
}
