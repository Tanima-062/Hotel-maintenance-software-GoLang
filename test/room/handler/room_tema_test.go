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

const WholesalerTemaId = "4"

var RequestTemaRoomData = roomBulk.RoomDataTema{
	RoomTypeID:   31,
	PropertyID:   1208011,
	Name:         "Test002",
	RoomTypeCode: "g1208010",
	RoomKindID:   4,
	RoomDesc:     "01",
	OcuMin:       3,
	OcuMax:       4,
	IsStopSales:  false,
	AmenityIDList: []int{
		1, 2,
	},
	Images: []roomBulk.Image{
		{
			ImageID: 273,
			Href:    "http://placeholder.com/300X300",
			Order:   1,
			Caption: "Test",
		},
	},
	Stocks: map[string]roomBulk.StockInputTema{
		"2023-07-01": roomBulk.StockInputTema{
			Stock:       12,
			IsStopSales: true,
		},
	},
}

// request body data
var roomTemaCreateOrUpdateRequestDataArray = []roomBulk.RoomDataTema{
	RequestTemaRoomData,
}

// MockRoomTemaBulkUseCase mock implementation
type MockRoomTemaBulkUseCase struct {
	mock.Mock
}

func (m *MockRoomTemaBulkUseCase) CreateOrUpdateBulk(request []roomBulk.RoomDataTema) error {
	return nil
}

func (m *MockRoomTemaBulkUseCase) Create(room *roomBulk.SaveInput) error {
	return nil
}

func (m *MockRoomTemaBulkUseCase) Delete(roomTypeID int64) error {
	return nil
}

func (m *MockRoomTemaBulkUseCase) FetchList(request *roomBulk.ListInput) ([]roomBulk.ListOutputTema, error) {
	return []roomBulk.ListOutputTema{}, nil
}

func (m *MockRoomTemaBulkUseCase) FetchAllAmenities() ([]roomBulk.AllAmenitiesOutput, error) {
	return []roomBulk.AllAmenitiesOutput{}, nil
}

func (m *MockRoomTemaBulkUseCase) FetchDetail(request *roomBulk.DetailInput) (*roomBulk.TemaDetailOutput, error) {
	return nil, nil
}

func (m *MockRoomTemaBulkUseCase) Update(request *roomBulk.SaveInput) error {
	return nil
}

func (m *MockRoomTemaBulkUseCase) UpdateStopSales(request *roomBulk.StopSalesInput) error {
	return nil
}

func (m *MockRoomTemaBulkUseCase) StoreBulkActivityLog(ServiceName string, Type string, host string, start time.Time) (int64, error) {
	return int64(0), nil
}

func (m *MockRoomTemaBulkUseCase) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, success bool, errorMessage string) error {
	return nil
}

// TxStart mock
func (m *MockRoomTemaBulkUseCase) TxStart() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// TxCommit mock
func (m MockRoomTemaBulkUseCase) TxCommit(tx *gorm.DB) error {
	return nil
}

// TxRollback mock
func (m MockRoomTemaBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

// TestRoomBulkHandlerCreateOrUpdateResponseSuccess
func TestTemaRoomBulkHandlerCreateOrUpdateResponseSuccess(t *testing.T) {
	mockUseCase := new(MockRoomTemaBulkUseCase)

	// Create a new RoomBulkHandler instance with the mock use case
	handler := &roomHandler.RoomHandler{
		RTemaUsecase:   mockUseCase,
		RLogRepository: mockUseCase,
	}

	// Echo request context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/room", nil)
	req.Header.Set("Wholesaler-Id", WholesalerTemaId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Bind(roomTemaCreateOrUpdateRequestDataArray)
	mockUseCase.On("StoreBulkActivityLog", "Test", "2", "localhost:1323", time.Now()).Return()
	mockUseCase.On("CreateOrUpdateBulk", roomTemaCreateOrUpdateRequestDataArray).Return(nil)
	err := handler.CreateOrUpdateBulk(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestRoomBulkHandlerCreateOrUpdateRequestBindingFailed
func TestTemaRoomBulkHandlerCreateOrUpdateRequestBindingFailed(t *testing.T) {
	mockUseCase := new(MockRoomTemaBulkUseCase)
	handler := &roomHandler.RoomHandler{
		RTemaUsecase:   mockUseCase,
		RLogRepository: mockUseCase,
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/room", strings.NewReader("Invalid data"))
	req.Header.Set("Wholesaler-Id", WholesalerTemaId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	emptyArray := []roomBulk.RoomDataTema{}
	//bind anything to fail
	c.Bind(emptyArray)

	err := handler.CreateOrUpdateBulk(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}
