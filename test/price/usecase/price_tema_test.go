package usecase_test

import (
	"errors"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	priceUsecase "github.com/Adventureinc/hotel-hm-api/src/price/usecase"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var request = []price.PriceTemaData{
	{
		PropertyID:           1,
		Disable:              true,
		RoomTypeCode:         "1",
		PackagePlanCode:      1,
		IsRoomCharge:         1,
		ListingPeriodStart:   "1",
		ListingPeriodEnd:     "1",
		IsPublishedYearRound: true,
		MinPax:               1,
		MaxPax:               1,
		PriceList: map[string][]price.PriceTema{
			"2023-07-01": []price.PriceTema{
				{Price: 1},
			},
		},
	},
	{
		PropertyID: -1,
		PriceList:  map[string][]price.PriceTema{},
	},
	{
		PropertyID:           2,
		Disable:              true,
		RoomTypeCode:         "1",
		PackagePlanCode:      1,
		IsRoomCharge:         1,
		ListingPeriodStart:   "1",
		ListingPeriodEnd:     "1",
		IsPublishedYearRound: true,
		MinPax:               1,
		MaxPax:               1,
		PriceList: map[string][]price.PriceTema{
			"2023-07-01": []price.PriceTema{
				{Price: 1},
			},
		},
	},
	{
		PropertyID:           3,
		Disable:              true,
		RoomTypeCode:         "1",
		PackagePlanCode:      1,
		IsRoomCharge:         1,
		ListingPeriodStart:   "1",
		ListingPeriodEnd:     "1",
		IsPublishedYearRound: true,
		MinPax:               1,
		MaxPax:               1,
		PriceList: map[string][]price.PriceTema{
			"2023-07-01": []price.PriceTema{
				{Price: 3},
			},
		},
	},
	{
		PropertyID: -3,
		PriceList:  map[string][]price.PriceTema{},
	},
}
var flag = 0

// MockRoomBulkUseCase mock implementation
type MockPlanTemaBulkUseCase struct {
	mock.Mock
}

func (m MockPlanTemaBulkUseCase) FetchOnePlan(propertyID int64, packagePlanCode int) (*plan.HtTmPlanTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) FetchOneWithPlanId(planID int64) (price.HtTmPlanTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) FetchList(propertyID int64, packagePlanCodeList []int) ([]plan.HtTmPlanTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) FetchAllByPropertyID(req plan.ListInput) ([]price.HtTmPlanTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) GetPlanIfPlanCodeExist(propertyID int64, planCode int64, roomTypeID int64) (price.HtTmPlanTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) UpdatePlanBulkTema(planTable price.HtTmPlanTemas, planID int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) GetNextPlanID() (price.HtTmPlanTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) CreatePlanBulkTema(planTable price.HtTmPlanTemas) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) ClearChildRateTema(planID int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) ClearImageTema(planID int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) CreateChildRateTema(childRates []price.HtTmChildRateTemas) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) DeletePlanTema(planCode int64, roomTypeIDs []int64) error {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) FetchActiveByPlanCode(planCode string) ([]price.HtTmPlanTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) FetchChildRates(planID int64) ([]price.HtTmChildRateTemas, error) {
	//TODO implement me
	panic("implement me")
}

// TxStart mock
func (m *MockPlanTemaBulkUseCase) TxStart() (*gorm.DB, error) {
	db, _, err := sqlmock.New()
	gorm, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})

	if flag == 1 {
		flag = 0
		return nil, errors.New("new err")
	}
	return gorm.Debug(), err
}

// TxCommit mock
func (m *MockPlanTemaBulkUseCase) TxCommit(tx *gorm.DB) error {
	if flag == 2 {
		flag = 0
		return errors.New("new err")
	}

	return nil
}

// TxRollback mock
func (m *MockPlanTemaBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

func (m MockPlanTemaBulkUseCase) FetchAllByPlanCodeList(planCodeList []int64, startDate string, endDate string) ([]price.HtTmPriceTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) FetchPricesByPlanID(planID int64) ([]price.HtTmPriceTemas, error) {
	//TODO implement me
	panic("implement me")
}

func (m MockPlanTemaBulkUseCase) DeletePriceTema(propertyID int64, packagePlanCode int64, roomTypeCode int, priceDate string) error {

	if propertyID == 2 {
		return errors.New("new err")
	}
	return nil
}

func (m MockPlanTemaBulkUseCase) CreatePrice(priceTable price.HtTmPriceTemas) error {
	if priceTable.PriceTemaTable.PropertyID == 3 {
		return errors.New("new err")
	}
	return nil
}

func (m MockPlanTemaBulkUseCase) FetchAllByPlanIDList(planIDList []int64, startDate string, endDate string) ([]price.HtTmPriceTemas, error) {
	//TODO implement me
	panic("implement me")
}

func TemaPriceUpdateDataProcess(reqData []price.PriceTemaData) (string, error) {
	mockUseCase := new(MockPlanTemaBulkUseCase)
	//mock required repositories with instance of tema use case
	useCases := &priceUsecase.PriceTemaUsecase{
		PriceTemaRepository: mockUseCase,
		PlanTemaRepository:  mockUseCase,
	}
	return useCases.Update(reqData)
}

func TestTemaPriceUpdateDataProcess(t *testing.T) {
	//using loop to run all the test data to successfully run all the success and error cases
	for _, data := range request {
		var reqArray = []price.PriceTemaData{
			data,
		}
		if data.PropertyID == -1 {
			flag = 1
		}
		if data.PropertyID == -3 {
			flag = 2
		}
		res, err := TemaPriceUpdateDataProcess(reqArray)
		if res != "" {
			assert.NotNilf(t, res, res)
			assert.Equal(t, res, res)
		} else {
			assert.Error(t, err)
		}

	}
}
