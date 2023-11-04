package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	roomBulk "github.com/Adventureinc/hotel-hm-api/src/room"
	roomHandler "github.com/Adventureinc/hotel-hm-api/src/room/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const WholesalerTlId = "3"

var RequestRoomData = roomBulk.RoomData{
	RoomTypeID:              121,
	PropertyID:              1208010,
	Name:                    "Test001",
	RoomTypeCode:            "g1208010",
	RoomKindID:              2,
	RoomDesc:                "01",
	StockSettingStart:       time.Now(),
	StockSettingEnd:         time.Now(),
	IsSettingStockYearRound: true,
	RoomCount:               12,
	OcuMin:                  3,
	OcuMax:                  4,
	IsSmoking:               false,
	IsStopSales:             false,
	IsDelete:                false,
	AmenityIDList: []int{
		1, 2,
	},
	Images: []roomBulk.Image{
		{
			ImageID: 273,
			Href:    "hhtp://placeholder.com/300X300",
			Order:   1,
			Caption: "Test",
		},
	},
	Stocks: map[string]roomBulk.SaveStockInput{
		"2023-07-01": roomBulk.SaveStockInput{
			Stock:       12,
			IsStopSales: true,
		},
	},
}

// request body data
var roomCreateOrUpdateRequestDataArray = []roomBulk.RoomData{
	RequestRoomData,
}

// MockRoomBulkUseCase mock implementation
type MockRoomBulkUseCase struct {
	mock.Mock
}

func (m *MockRoomBulkUseCase) Create(room *roomBulk.SaveInput) error {
	return nil
}

func (m *MockRoomBulkUseCase) CreateOrUpdateBulk(request []roomBulk.RoomData) error {
	args := m.Called(roomCreateOrUpdateRequestDataArray)
	return args.Error(0)
}

func (m *MockRoomBulkUseCase) Delete(roomTypeID int64) error {
	return nil
}

func (m *MockRoomBulkUseCase) FetchList(request *roomBulk.ListInput) ([]roomBulk.ListOutputTl, error) {
	return []roomBulk.ListOutputTl{}, nil
}

func (m *MockRoomBulkUseCase) FetchAllAmenities() ([]roomBulk.AllAmenitiesOutput, error) {
	return []roomBulk.AllAmenitiesOutput{}, nil
}

func (m *MockRoomBulkUseCase) FetchDetail(request *roomBulk.DetailInput) (*roomBulk.DetailOutput, error) {
	return nil, nil
}

func (m *MockRoomBulkUseCase) Update(request *roomBulk.SaveInput) error {
	return nil
}

func (m *MockRoomBulkUseCase) UpdateStopSales(request *roomBulk.StopSalesInput) error {
	return nil
}

func (m *MockRoomBulkUseCase) StoreBulkActivityLog(ServiceName string, Type string, host string, start time.Time) (int64, error) {
	return int64(0), nil
}

func (m *MockRoomBulkUseCase) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, success bool, errorMessage string) error {
	return nil
}

// TxStart mock
func (m *MockRoomBulkUseCase) TxStart() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// TxCommit mock
func (m MockRoomBulkUseCase) TxCommit(tx *gorm.DB) error {
	return nil
}

// TxRollback mock
func (m MockRoomBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

// TestRoomBulkHandlerCreateOrUpdateResponseSuccess
func TestRoomBulkHandlerCreateOrUpdateResponseSuccess(t *testing.T) {
	mockUseCase := new(MockRoomBulkUseCase)

	// Create a new RoomBulkHandler instance with the mock use case
	handler := &roomHandler.RoomHandler{
		RTlUsecase:     mockUseCase,
		RLogRepository: mockUseCase,
	}

	// Echo request context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/room", nil)
	req.Header.Set("Wholesaler-Id", WholesalerTlId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Bind(roomCreateOrUpdateRequestDataArray)
	mockUseCase.On("StoreBulkActivityLog", "Test", "1", "localhost:1323", time.Now()).Return()
	mockUseCase.On("CreateOrUpdateBulk", roomCreateOrUpdateRequestDataArray).Return(nil)
	err := handler.CreateOrUpdateBulk(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestRoomBulkHandlerCreateOrUpdateRequestBindingFailed
func TestRoomBulkHandlerCreateOrUpdateRequestBindingFailed(t *testing.T) {
	handler := &roomHandler.RoomHandler{}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/room", strings.NewReader("Invalid data"))
	req.Header.Set("Wholesaler-Id", WholesalerTlId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//bind anything to fail
	c.Bind(mock.Anything)

	err := handler.CreateOrUpdateBulk(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}
