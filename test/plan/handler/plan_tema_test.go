package handler_test

import (
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"github.com/Adventureinc/hotel-hm-api/src/plan/handler"
	"github.com/Adventureinc/hotel-hm-api/src/price"
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
var RequestTemaPlanData = price.TemaPlanData{
	TemaPlanTable: price.TemaPlanTable{
		PlanID:                  1,
		PlanGroupID:             1,
		PropertyID:              1,
		PackagePlanCode:         1,
		PlanName:                "1",
		Desc:                    "1",
		LangCd:                  "1",
		PlanType:                1,
		Payment:                 1,
		ListingPeriodStart:      "1",
		ListingPeriodEnd:        "1",
		IsRoomCharge:            1,
		Available:               true,
		RateType:                1,
		ListingPeriodStartH:     1,
		ListingPeriodStartM:     1,
		ListingPeriodEndH:       1,
		ListingPeriodEndM:       1,
		ReservePeriodStart:      "1",
		ReservePeriodEnd:        "1",
		CheckinTimeStartH:       1,
		CheckinTimeStartM:       1,
		CheckinTimeEndH:         1,
		CheckinTimeEndM:         1,
		CheckoutTimeEndH:        1,
		CheckoutTimeEndM:        1,
		StayLimitMin:            1,
		StayLimitMax:            1,
		AdvBKCreateStartEnabled: 1,
		AdvBKCreateEndEnabled:   1,
		AdvBkCreateStartD:       1,
		ReserveAcceptTime:       "1",
		AdvBkCreateStartH:       1,
		AdvBkCreateStartM:       1,
		ReserveDeadlineTime:     "1",
		AdvBkCreateEndD:         1,
		AdvBkCreateEndH:         1,
		AdvBkCreateEndM:         1,
		AdvBkModifyEndEnabled:   1,
		AdvBkModifyEndD:         1,
		AdvBkModifyEndH:         1,
		AdvBkModifyEndM:         1,
		AdvBkCancelEndEnabled:   1,
		AdvBkCancelEndD:         1,
		AdvBkCancelEndH:         1,
		AdvBkCancelEndM:         1,
		ServiceChargeType:       1,
		ServiceChargeValue:      1,
		OptionalItems:           "1",
		CancelpolicyJson:        "1",
		ChildrenJson:            "1",
		ChildrenAcceptable1:     1,
		ChildrenAcceptable2:     1,
		ChildrenAcceptable3:     1,
		ChildrenAcceptable4:     1,
		ChildrenAcceptable5:     1,
		PictureJson:             "1",
		RoomTypeID:              1,
		IsPublishedYearRound:    true,
		IsAccommodatedYearRound: true,
		TaxCategory:             true,
		MinStayCategory:         1,
		MaxStayCategory:         1,
		MealConditionBreakfast:  false,
		MealConditionDinner:     false,
		MealConditionLunch:      false,
	},
	RoomTypeCode: "1",
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
}

var planTemaCreateRequestData = []price.TemaPlanData{
	RequestTemaPlanData,
}

// MockTemaPlanBulkUseCase mock implementation
type MockTemaPlanBulkUseCase struct {
	mock.Mock
}

// mocked implementation of Create method
func (m *MockTemaPlanBulkUseCase) FetchList(req *plan.ListInput) ([]plan.TemaBulkListOutput, error) {
	args := m.Called(req)
	return []plan.TemaBulkListOutput{}, args.Error(0)
}

// mocked implementation of Create method
func (m *MockTemaPlanBulkUseCase) CreateBulk(req []price.TemaPlanData) (string, error) {
	args := m.Called(planTemaCreateRequestData)
	return "any message", args.Error(0)
}

func (m *MockTemaPlanBulkUseCase) Detail(request *plan.DetailInput) (*plan.TemaBulkDetailOutput, error) {
	return &plan.TemaBulkDetailOutput{}, nil
}

func (m *MockTemaPlanBulkUseCase) StoreBulkActivityLog(ServiceName string, Type string, host string, start time.Time) (int64, error) {
	return int64(0), nil
}

func (m *MockTemaPlanBulkUseCase) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, success bool, errorMessage string) error {
	return nil
}

// TxStart mock
func (m *MockTemaPlanBulkUseCase) TxStart() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

// TxCommit mock
func (m MockTemaPlanBulkUseCase) TxCommit(tx *gorm.DB) error {
	return nil
}

// TxRollback mock
func (m MockTemaPlanBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

// TestPlanBulkHandlerCreateResponseSuccess
func TestTemaPlanBulkHandlerCreateResponseSuccess(t *testing.T) {
	// Create a new Echo request context for testing
	e := echo.New()

	// Create a new mock use case
	mockUseCase := new(MockTemaPlanBulkUseCase)

	// Create a new PlanBulkHandler instance with the mock use case
	handler := &handler.PlanHandler{
		PBulkTemaUsecase: mockUseCase,
		RLogRepository:   mockUseCase,
	}

	req := httptest.NewRequest(http.MethodPost, "/bulk/plan", nil)
	req.Header.Set("Wholesaler-Id", WholesalerTemaId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Bind(planTemaCreateRequestData)
	mockUseCase.On("StoreBulkActivityLog", "Test", "1", "localhost:1323", time.Now()).Return()
	mockUseCase.On("CreateBulk", planTemaCreateRequestData).Return(nil)
	err := handler.CreateOrUpdateBulk(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// TestPlanBulkHandlerCreateBindingFailed
func TestTemaPlanBulkHandlerCreateBindingFailed(t *testing.T) {
	handler := &handler.PlanHandler{}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/bulk/plan", strings.NewReader("Invalid data"))
	req.Header.Set("Wholesaler-Id", WholesalerTemaId)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//bind anything to fail
	c.Bind(mock.Anything)

	err := handler.CreateOrUpdateBulk(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}
