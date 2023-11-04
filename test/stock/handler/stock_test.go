package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	"github.com/Adventureinc/hotel-hm-api/src/stock/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// stock request test data set
var RequestStockData = stock.StockData{
	PropertyID:   1208010,
	RoomTypeCode: "g1208010",
	Stocks: map[string]stock.UpdateStockInput{
		"2023-07-01": stock.UpdateStockInput{
			Stock:       12,
			IsStopSales: true,
		},
		"2023-07-02": stock.UpdateStockInput{
			Stock:       11,
			IsStopSales: false,
		},
	},
}

// request body data
var StockUpdateRequestData = []stock.StockData{
	RequestStockData,
}

// MockStockHandler mock implementation
type MockStockHandler struct {
	mock.Mock
}

// mocked implementation of Update method
func (m *MockStockHandler) Update(req []stock.StockData) error {
	args := m.Called(StockUpdateRequestData)
	return args.Error(0)
}

// mocked implementation of Update method
func (m *MockStockHandler) Validate(req []stock.StockData) error {
	args := m.Called(StockUpdateRequestData)
	return args.Error(0)
}

func (m *MockStockHandler) FetchAll(request *stock.ListInput) (*[]stock.ListOutput, error) {
	return &[]stock.ListOutput{}, nil
}

func (m *MockStockHandler) FetchCalendar(hmUser account.HtTmHotelManager, request stock.CalendarInput) (*[]stock.CalendarOutput, error) {
	return &[]stock.CalendarOutput{}, nil
}

func (m *MockStockHandler) Save(request *[]stock.SaveInput) error {
	return nil
}

func (m *MockStockHandler) UpdateBulk(request []stock.StockData) error {
	args := m.Called(StockUpdateRequestData)
	return args.Error(0)
}

func (m *MockStockHandler) UpdateStopSales(request *stock.StopSalesInput) error {
	return nil
}

func (m *MockStockHandler) StoreBulkActivityLog(ServiceName string, Type string, host string, start time.Time) (int64, error) {
	return int64(0), nil
}

func (m *MockStockHandler) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, success bool, errorMessage string) error {
	return nil
}

// TxStart mock
func (m *MockStockHandler) TxStart() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// TxCommit mock
func (m MockStockHandler) TxCommit(tx *gorm.DB) error {
	return nil
}

// TxRollback mock
func (m MockStockHandler) TxRollback(tx *gorm.DB) {
	return
}

// TestStockHandlerUpdateResponseSuccess
func TestStockHandlerUpdateResponseSuccess(t *testing.T) {
	// new mock use case
	mockUseCase := new(MockStockHandler)

	// StockHandler instance with the mock use case
	handler := &handler.StockHandler{
		STlUsecase:     mockUseCase,
		RLogRepository: mockUseCase,
	}

	// new Echo request context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/stock/update", nil)
	req.Header.Set("Wholesaler-Id", "3")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// expectations setup on the mock use case
	mockUseCase.On("StoreBulkActivityLog", "Test", "1", "localhost:1323", time.Now()).Return()
	mockUseCase.On("UpdateBulk", StockUpdateRequestData).Return(nil)

	err := handler.UpdateBulk(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestStockHandlerUpdateRequestBindingFailed
func TestStockHandlerUpdateRequestBindingFailed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/room/update", strings.NewReader("Invalid data"))
	req.Header.Set("Wholesaler-Id", "3")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//bind anything to fail
	c.Bind(mock.Anything)

	handler := &handler.StockHandler{}
	err := handler.UpdateBulk(c)

	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}
