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
var RequestTemaStockData = stock.StockDataTema{
	PropertyID:   1208011,
	RoomTypeCode: "g1208010",
	Stocks: map[string]stock.UpdateStockTemaInput{
		"2023-07-01": {
			Stock:   12,
			Disable: true,
		},
		"2023-07-02": {
			Stock:   11,
			Disable: false,
		},
	},
}

// request body data
var StockTemaUpdateRequestData = []stock.StockDataTema{
	RequestTemaStockData,
}

// MockStockTemaHandler mock implementation
type MockStockTemaHandler struct {
	mock.Mock
}

func (m *MockStockTemaHandler) UpdateBulkTema([]stock.StockDataTema) error {
	//TODO implement me
	return nil
}

// mocked implementation of Update method
func (m *MockStockTemaHandler) Update(req []stock.StockData) error {
	args := m.Called(StockTemaUpdateRequestData)
	return args.Error(0)
}

// mocked implementation of Update method
func (m *MockStockTemaHandler) Validate(req []stock.StockData) error {
	args := m.Called(StockTemaUpdateRequestData)
	return args.Error(0)
}

func (m *MockStockTemaHandler) FetchAll(request *stock.ListInput) (*[]stock.ListOutput, error) {
	return &[]stock.ListOutput{}, nil
}

func (m *MockStockTemaHandler) FetchCalendar(hmUser account.HtTmHotelManager, request stock.CalendarInput) (*[]stock.CalendarOutputTema, error) {
	return &[]stock.CalendarOutputTema{}, nil
}

func (m *MockStockTemaHandler) Save(request *[]stock.SaveInput) error {
	return nil
}

func (m *MockStockTemaHandler) UpdateBulk(request []stock.StockData) error {
	args := m.Called(StockTemaUpdateRequestData)
	return args.Error(0)
}

func (m *MockStockTemaHandler) UpdateStopSales(request *stock.StopSalesInput) error {
	return nil
}

func (m *MockStockTemaHandler) StoreBulkActivityLog(ServiceName string, Type string, host string, start time.Time) (int64, error) {
	return int64(0), nil
}

func (m *MockStockTemaHandler) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, success bool, errorMessage string) error {
	return nil
}

// TxStart mock
func (m *MockStockTemaHandler) TxStart() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// TxCommit mock
func (m MockStockTemaHandler) TxCommit(tx *gorm.DB) error {
	return nil
}

// TxRollback mock
func (m MockStockTemaHandler) TxRollback(tx *gorm.DB) {
	return
}

// TestStockHandlerUpdateResponseSuccess
func TestStockTemaHandlerUpdateResponseSuccess(t *testing.T) {
	// new mock use case
	mockUseCase := new(MockStockTemaHandler)

	// StockHandler instance with the mock use case
	handler := &handler.StockHandler{
		STemaUsecase:   mockUseCase,
		RLogRepository: mockUseCase,
	}

	// new Echo request context for testing
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/stock/update", nil)
	req.Header.Set("Wholesaler-Id", "4")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// expectations setup on the mock use case
	mockUseCase.On("StoreBulkActivityLog", "Test", "1", "localhost:1323", time.Now()).Return()
	mockUseCase.On("UpdateBulk", StockTemaUpdateRequestData).Return(nil)

	err := handler.UpdateBulk(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestStockHandlerUpdateRequestBindingFailed
func TestStockTemaHandlerUpdateRequestBindingFailed(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/room/update", strings.NewReader("Invalid data"))
	req.Header.Set("Wholesaler-Id", "4")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//bind anything to fail
	c.Bind(stock.StockDataTema{
		PropertyID: -1,
	})
	mockUseCase := new(MockStockTemaHandler)

	handler := &handler.StockHandler{
		STemaUsecase:   mockUseCase,
		RLogRepository: mockUseCase,
	}
	err := handler.UpdateBulk(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}
