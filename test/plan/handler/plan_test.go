package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/image"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"github.com/Adventureinc/hotel-hm-api/src/plan/handler"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

const WholesalerTlId = "3"

// Test request body data
var RequestPlanData = price.PlanData{
	PlanTable: price.PlanTable{
		PlanID:                   123,
		PlanGroupID:              3,
		RoomTypeID:               2,
		PropertyID:               1208010,
		PlanCode:                 "G1212",
		LangCd:                   "ja-JP",
		Name:                     "Plan01",
		Description:              "N/A",
		ChargeCategory:           6,
		TaxCategory:              false,
		AccommodationPeriodStart: time.Now(),
		AccommodationPeriodEnd:   time.Now(),
		IsAccommodatedYearRound:  false,
		PublishingStartDate:      time.Now(),
		PublishingEndDate:        time.Now(),
		IsPublishedYearRound:     false,
		ReserveAcceptDate:        5,
		ReserveAcceptTime:        "Test accept time",
		ReserveDeadlineDate:      3,
		ReserveDeadlineTime:      "Test deadline time",
		MinStayCategory:          false,
		MinStayNum:               2,
		MaxStayCategory:          false,
		MaxStayNum:               4,
		MealConditionBreakfast:   false,
		MealConditionDinner:      false,
		MealConditionLunch:       false,
		IsNoCancel:               false,
		IsStopSales:              false,
		IsDelete:                 false,
	},
	SelectedRooms: []int64{
		1, 2,
	},
	Images: []image.PlanImagesInput{
		{
			ImageID: 273,
			PlanID:  1,
			Order:   1,
		},
	},
	ChildRates: []price.ChildRateTable{
		{
			ChildRateID:   1,
			ChildRateType: 1,
			PlanID:        237,
			FromAge:       5,
			ToAge:         18,
			Receive:       false,
			RateCategory:  4,
			Rate:          4,
			CalcCategory:  false,
		},
		{
			ChildRateID:   2,
			ChildRateType: 1,
			PlanID:        237,
			FromAge:       5,
			ToAge:         18,
			Receive:       false,
			RateCategory:  5,
			Rate:          8,
			CalcCategory:  false,
		},
	},
	Prices: map[string][]price.Price{
		"2023-06-01": {
			price.Price{
				Type:  "test",
				Price: 1222,
			},
		},
		"2023-06-02": {
			price.Price{
				Type:  "test",
				Price: 3333,
			},
		},
	},
}

var planCreateRequestData = []price.PlanData{
	RequestPlanData,
}

// MockplanBulkUseCase mock implementation
type MockplanBulkUseCase struct {
	mock.Mock
}

// mocked implementation of Create method
func (m *MockplanBulkUseCase) FetchList(req *plan.ListInput) ([]plan.BulkListOutput, error) {
	args := m.Called(req)
	return []plan.BulkListOutput{}, args.Error(0)
}

// mocked implementation of Create method
func (m *MockplanBulkUseCase) CreateBulk(req []price.PlanData) (string, error) {
	args := m.Called(planCreateRequestData)
	return "any message", args.Error(0)
}

func (m *MockplanBulkUseCase) Detail(request *plan.DetailInput) (*plan.BulkDetailOutput, error) {
	return &plan.BulkDetailOutput{}, nil
}

func (m *MockplanBulkUseCase) StoreBulkActivityLog(ServiceName string, Type string, host string, start time.Time) (int64, error) {
	return int64(0), nil
}

func (m *MockplanBulkUseCase) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, success bool, errorMessage string) error {
	return nil
}

// TxStart mock
func (m *MockplanBulkUseCase) TxStart() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// TxCommit mock
func (m MockplanBulkUseCase) TxCommit(tx *gorm.DB) error {
	return nil
}

// TxRollback mock
func (m MockplanBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

// TestPlanBulkHandlerCreateResponseSuccess
func TestPlanBulkHandlerCreateResponseSuccess(t *testing.T) {
	// Create a new Echo request context for testing
	e := echo.New()

	// Create a new mock use case
	mockUseCase := new(MockplanBulkUseCase)

	// Create a new PlanBulkHandler instance with the mock use case
	handler := &handler.PlanHandler{
		PBulkTlUsecase: mockUseCase,
		RLogRepository: mockUseCase,
	}

	req := httptest.NewRequest(http.MethodPost, "/bulk/plan", nil)
	req.Header.Set("Wholesaler-Id", WholesalerTlId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Bind(planCreateRequestData)
	mockUseCase.On("StoreBulkActivityLog", "Test", "1", "localhost:1323", time.Now()).Return()
	mockUseCase.On("CreateBulk", planCreateRequestData).Return(nil)
	err := handler.CreateOrUpdateBulk(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestPlanBulkHandlerCreateBindingFailed
func TestPlanBulkHandlerCreateBindingFailed(t *testing.T) {
	handler := &handler.PlanHandler{}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/plan", strings.NewReader("Invalid data"))
	req.Header.Set("Wholesaler-Id", WholesalerTlId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//bind anything to fail
	c.Bind(mock.Anything)

	err := handler.CreateOrUpdateBulk(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}
