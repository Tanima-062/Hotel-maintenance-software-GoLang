package handler

import (
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"github.com/Adventureinc/hotel-hm-api/src/price/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const WholesalerTemaId = "4"

// Test request body data
var RequestPriceData = price.PriceTemaData{
	PropertyID:           1208010,
	Disable:              true,
	RoomTypeCode:         "2",
	PackagePlanCode:      2,
	IsRoomCharge:         0,
	ListingPeriodStart:   "2023-10-22",
	ListingPeriodEnd:     "2023-10-28",
	IsPublishedYearRound: false,
	MinPax:               0,
	MaxPax:               5,
	PriceList: map[string][]price.PriceTema{
		"2023-07-01": []price.PriceTema{
			{
				Price: 1000,
			},
			{
				Price: 2000,
			},
		},
		"2023-07-02": []price.PriceTema{
			{
				Price: 1000,
			},
			{
				Price: 2000,
			},
		},
	},
}
var priceTemaCreateRequestData = []price.PriceTemaData{
	RequestPriceData,
}

type MockTemaplanBulkUseCase struct {
	mock.Mock
}

func (m MockTemaplanBulkUseCase) TxStart() (*gorm.DB, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockTemaplanBulkUseCase) TxCommit(tx *gorm.DB) error {
	//TODO implement me
	panic("implement me")
}

func (m MockTemaplanBulkUseCase) TxRollback(tx *gorm.DB) {
	//TODO implement me
	panic("implement me")
}

func (m MockTemaplanBulkUseCase) StoreBulkActivityLog(ServiceName string, Type string, Host string, start time.Time) (int64, error) {
	//TODO implement me
	return int64(0), nil
}

func (m MockTemaplanBulkUseCase) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, status bool, errorMessage string) error {
	//TODO implement me
	return nil
}

func (m MockTemaplanBulkUseCase) Update(request []price.PriceTemaData) (string, error) {
	//TODO implement me
	return "Price has been updated", nil
}

// TestPlanBulkHandlerCreateResponseSuccess
func TestTemaPriceBulkHandlerCreateResponseSuccess(t *testing.T) {
	// Create a new Echo request context for testing
	e := echo.New()

	// Create a new mock use case
	mockUseCase := new(MockTemaplanBulkUseCase)

	// Create a new PlanBulkHandler instance with the mock use case
	handler := &handler.PriceHandler{
		PTemaUsecase:   mockUseCase,
		RLogRepository: mockUseCase,
	}

	req := httptest.NewRequest(http.MethodPost, "/bulk/price", nil)
	req.Header.Set("Wholesaler-Id", WholesalerTemaId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Bind(priceTemaCreateRequestData)
	mockUseCase.On("RLogRepository", "Test", "1", "localhost:1323", time.Now()).Return()
	mockUseCase.On("UpdateBulk", priceTemaCreateRequestData).Return(nil)
	err := handler.UpdateBulk(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestPlanBulkHandlerCreateBindingFailed
func TestTemaPriceBulkHandlerCreateBindingFailed(t *testing.T) {
	mockUseCase := new(MockTemaplanBulkUseCase)
	handler := &handler.PriceHandler{
		PTemaUsecase:   mockUseCase,
		RLogRepository: mockUseCase,
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/price", strings.NewReader("Invalid data"))
	req.Header.Set("Wholesaler-Id", WholesalerTemaId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//bind anything to fail
	c.Bind(mock.Anything)

	err := handler.UpdateBulk(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}
