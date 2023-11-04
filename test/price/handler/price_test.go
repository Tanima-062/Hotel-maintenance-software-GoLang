package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/price"
	"github.com/Adventureinc/hotel-hm-api/src/price/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const WholesalerTlId = "3"

// Test request body data
var RequestPriceData = price.PriceData{
	PropertyID: 1208010,
	PlanCode:   "g1208010",
	Prices: map[string][]price.Price{
		"2023-07-01": []price.Price{
			{
				Type:  "Test",
				Price: 1000,
			},
			{
				Type:  "Test01",
				Price: 2000,
			},
		},
		"2023-07-02": []price.Price{
			{
				Type:  "Test",
				Price: 1000,
			},
			{
				Type:  "Test01",
				Price: 2000,
			},
		},
	},
}

// request body data array
var priceUpdateRequestData = []price.PriceData{
	RequestPriceData,
}

// MockPriceBulkUseCase mock implementation
type MockPriceBulkUseCase struct {
	mock.Mock
}

// mocked implementation of save method
func (m *MockPriceBulkUseCase) Save(
	requestData price.PlanTable,
	childRateTables []price.HtTmChildRateTls,
	priceData price.Price,
	date string,
	db *gorm.DB) error {
	return nil
}

// mocked implementation of Update method
func (m *MockPriceBulkUseCase) Update(requestData []price.PriceData) (string, error) {
	args := m.Called(priceUpdateRequestData)
	return "any message", args.Error(0)
}

// mocked implementation of GetPriceData method
func (m *MockPriceBulkUseCase) GetPriceData(request price.PlanTable, childRateTables []price.HtTmChildRateTls, priceData price.Price, date string) price.HtTmPriceTls {
	return price.HtTmPriceTls{}
}

func (m *MockPriceBulkUseCase) StoreBulkActivityLog(ServiceName string, Type string, host string, start time.Time) (int64, error) {
	return int64(0), nil
}

func (m *MockPriceBulkUseCase) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, success bool, errorMessage string) error {
	return nil
}

// TxStart mock
func (m *MockPriceBulkUseCase) TxStart() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// TxCommit mock
func (m MockPriceBulkUseCase) TxCommit(tx *gorm.DB) error {
	return nil
}

// TxRollback mock
func (m MockPriceBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

// TestPriceBulkHandlerUpdateResponseSuccess
func TestPriceBulkHandlerUpdateResponseSuccess(t *testing.T) {
	// mock use case
	mockUseCase := new(MockPriceBulkUseCase)

	// PriceBulkHandler instance with the mock use case
	handler := &handler.PriceHandler{
		PTlUsecase:     mockUseCase,
		RLogRepository: mockUseCase,
	}

	// Echo request context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/plan/price/update", nil)
	req.Header.Set("Wholesaler-Id", WholesalerTlId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Bind(priceUpdateRequestData)
	mockUseCase.On("StoreBulkActivityLog", "Test", "1", "localhost:1323", time.Now()).Return()
	mockUseCase.On("Update", priceUpdateRequestData).Return(nil)
	err := handler.UpdateBulk(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestPriceBulkHandlerUpdateRequestBindingFailed
func TestPriceBulkHandlerUpdateRequestBindingFailed(t *testing.T) {
	handler := &handler.PriceHandler{}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/plan/price/update", strings.NewReader("Invalid data"))
	req.Header.Set("Wholesaler-Id", WholesalerTlId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//bind anything to fail
	c.Bind(mock.Anything)

	err := handler.UpdateBulk(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}
