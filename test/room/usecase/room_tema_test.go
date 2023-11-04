package usecase_test

import (
	"errors"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	roomUseCase "github.com/Adventureinc/hotel-hm-api/src/room/usecase"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var request = []room.RoomDataTema{
	{
		RoomTypeID:   31,
		PropertyID:   1208011,
		Name:         "Test002",
		RoomTypeCode: "290",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   33,
		PropertyID:   1208011,
		Name:         "Test002",
		RoomTypeCode: "290",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   34,
		PropertyID:   1208011,
		Name:         "Test002",
		RoomTypeCode: "299",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   35,
		PropertyID:   1208011,
		Name:         "Test002",
		RoomTypeCode: "300",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   36,
		PropertyID:   1208011,
		Name:         "Test002",
		RoomTypeCode: "301",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   37,
		PropertyID:   1208011,
		Name:         "Test002",
		RoomTypeCode: "296",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   32,
		PropertyID:   1208012,
		Name:         "Test002",
		RoomTypeCode: "291",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   38,
		PropertyID:   1208013,
		Name:         "Test002",
		RoomTypeCode: "297",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   39,
		PropertyID:   1208016,
		Name:         "Test002",
		RoomTypeCode: "366",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
	{
		RoomTypeID:   40,
		PropertyID:   1208018,
		Name:         "Test002",
		RoomTypeCode: "366",
		RoomKindID:   4,
		RoomDesc:     "01",
		OcuMin:       3,
		OcuMax:       4,
		IsStopSales:  false,
		AmenityIDList: []int{
			1, 2,
		},
		Images: []room.Image{
			{
				ImageID: 273,
				Href:    "http://placeholder.com/300X300",
				Order:   1,
				Caption: "Test",
			},
		},
		Stocks: map[string]room.StockInputTema{
			"2023-07-01": room.StockInputTema{
				Stock:       12,
				IsStopSales: true,
			},
		},
	},
}
var flag = 0

// MockRoomBulkUseCase mock implementation
type MockRoomTemaBulkUseCase struct {
	mock.Mock
}

// TxStart mock
func (m *MockRoomTemaBulkUseCase) TxStart() (*gorm.DB, error) {
	db, _, err := sqlmock.New()
	gorm, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if flag == 2 {
		return nil, errors.New("new err")
	}
	return gorm.Debug(), err
}

// TxCommit mock
func (m *MockRoomTemaBulkUseCase) TxCommit(tx *gorm.DB) error {

	if flag == 1 {
		return errors.New("new err")
	}
	return nil
}

// TxRollback mock
func (m *MockRoomTemaBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

func (m *MockRoomTemaBulkUseCase) FetchOne(roomTypeCode int, propertyID int64) (*room.HtTmRoomTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRoomTemaBulkUseCase) FetchListWithPropertyId(roomTypeCodeList []int, propertyID int64) ([]room.HtTmRoomTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRoomTemaBulkUseCase) FetchRoomTypeIdByRoomTypeCode(propertyID int64, roomTypeCode string) (room.HtTmRoomTypeTemas, error) {
	if propertyID == 1208012 {
		return room.HtTmRoomTypeTemas{
			RoomTypeTema: room.RoomTypeTema{},
		}, errors.New("new err")
	}
	if roomTypeCode == "296" {
		return room.HtTmRoomTypeTemas{
			RoomTypeTema: room.RoomTypeTema{
				RoomTypeID:   37,
				RoomTypeCode: roomTypeCode,
				PropertyID:   propertyID,
			},
		}, nil
	}
	if propertyID == 1208013 {
		return room.HtTmRoomTypeTemas{
			RoomTypeTema: room.RoomTypeTema{},
		}, errors.New("new err")
	}
	if roomTypeCode == "299" {
		return room.HtTmRoomTypeTemas{
			RoomTypeTema: room.RoomTypeTema{
				RoomTypeID:   34,
				RoomTypeCode: roomTypeCode,
				PropertyID:   propertyID,
			},
		}, nil
	}
	if roomTypeCode == "300" {
		return room.HtTmRoomTypeTemas{
			RoomTypeTema: room.RoomTypeTema{
				RoomTypeID:   35,
				RoomTypeCode: roomTypeCode,
				PropertyID:   propertyID,
			},
		}, nil
	}
	if propertyID == 1208016 {
		return room.HtTmRoomTypeTemas{
			RoomTypeTema: room.RoomTypeTema{
				RoomTypeID:   39,
				RoomTypeCode: roomTypeCode,
				PropertyID:   propertyID,
			},
		}, nil
	}
	if roomTypeCode == "301" {
		return room.HtTmRoomTypeTemas{
			RoomTypeTema: room.RoomTypeTema{
				RoomTypeID:   36,
				RoomTypeCode: roomTypeCode,
				PropertyID:   propertyID,
			},
		}, nil
	}
	return room.HtTmRoomTypeTemas{
		RoomTypeTema: room.RoomTypeTema{
			RoomTypeID:   31,
			RoomTypeCode: roomTypeCode,
			PropertyID:   propertyID,
		},
	}, nil

}

func (m *MockRoomTemaBulkUseCase) CreateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error {

	if roomTable.PropertyID == 1208013 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockRoomTemaBulkUseCase) UpdateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error {
	if roomTable.RoomTypeID == 37 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockRoomTemaBulkUseCase) MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MockRoomTemaBulkUseCase) FetchRoomByRoomTypeID(roomTypeID int64) (*room.HtTmRoomTypeTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRoomTemaBulkUseCase) ClearRoomToAmenities(roomTypeID int64) error {
	if roomTypeID == 0 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockRoomTemaBulkUseCase) CreateRoomToAmenities(roomTypeID int64, tlRoomAmenityID int64) error {

	if roomTypeID == 34 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockRoomTemaBulkUseCase) ClearRoomImage(roomTypeID int64) error {

	if roomTypeID == 35 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockRoomTemaBulkUseCase) CreateRoomOwnImages(images []room.HtTmRoomOwnImagesTemas) error {
	if images[0].RoomTypeID == 36 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockRoomTemaBulkUseCase) FetchRoomsByPropertyID(req room.ListInput) ([]room.HtTmRoomTypeTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRoomTemaBulkUseCase) FetchAllAmenities() ([]room.HtTmRoomAmenityTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRoomTemaBulkUseCase) FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomAmenitiesTema, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockRoomTemaBulkUseCase) FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomImagesTema, error) {
	//TODO implement me
	panic("implement me")
}
func TemaRoomBulkCreateOrUpdateDataProcess(reqData []room.RoomDataTema) error {
	mockUseCase := new(MockRoomTemaBulkUseCase)
	//mock required repositories with instance of tema use case
	useCases := &roomUseCase.RoomTemaUseCase{
		RTemaRepository: mockUseCase,
	}
	return useCases.CreateOrUpdateBulk(reqData)
}

// request body data
func TestTemaRoomBulkCreateOrUpdateDataProcess(t *testing.T) {
	//using loop to run all the test data to successfully run all the success and error cases
	for _, data := range request {
		var reqArray = []room.RoomDataTema{
			data,
		}
		if data.RoomTypeID == 39 {
			flag = 1
		}
		if data.RoomTypeID == 40 {
			flag = 2
		}
		res := TemaRoomBulkCreateOrUpdateDataProcess(reqArray)
		if res == nil {
			assert.Nil(t, res)
			assert.Equal(t, res, nil)
		} else {
			assert.Error(t, res)
			assert.Equal(t, res, errors.New("new err"))
		}

	}
}
