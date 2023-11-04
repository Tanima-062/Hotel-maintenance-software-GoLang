package usecase_test

import (
	"errors"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	planUsecase "github.com/Adventureinc/hotel-hm-api/src/plan/usecase"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var request = []price.TemaPlanData{

	{
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
			IsPublishedYearRound:    false,
			IsAccommodatedYearRound: false,
			TaxCategory:             true,
			MinStayCategory:         1,
			MaxStayCategory:         1,
			MealConditionBreakfast:  true,
			MealConditionDinner:     true,
			MealConditionLunch:      true,
		},
		RoomTypeCode: "1",
		ChildRates: []price.ChildRateTable{
			{
				ChildRateID:   1,
				ChildRateType: 1,
				PlanID:        1,
				FromAge:       1,
				ToAge:         1,
				Receive:       true,
				RateCategory:  1,
				Rate:          1,
				CalcCategory:  true,
			},
		},

		Images: []image.PlanImagesInput{
			{
				ImageID: 1,
				PlanID:  1,
				Order:   1,
			},
		},
	},
	{
		TemaPlanTable: price.TemaPlanTable{
			PlanID:     12,
			PropertyID: 22,
		},
		RoomTypeCode: "2",
		ChildRates: []price.ChildRateTable{
			{
				ChildRateID: 2,
			},
		},
		Images: []image.PlanImagesInput{},
	},
	{
		TemaPlanTable: price.TemaPlanTable{
			PlanID:     1,
			PropertyID: 2,
		},
		RoomTypeCode: "3",
		ChildRates: []price.ChildRateTable{
			{
				ChildRateID: 2,
			},
		},
		Images: []image.PlanImagesInput{},
	},
	{
		TemaPlanTable: price.TemaPlanTable{
			PlanID:     2,
			PropertyID: 3,
		},
		RoomTypeCode: "4",
		ChildRates:   []price.ChildRateTable{},
		Images:       []image.PlanImagesInput{},
	},
	{
		TemaPlanTable: price.TemaPlanTable{
			PlanID:     0,
			PropertyID: 4,
		},
		RoomTypeCode: "5",
		ChildRates: []price.ChildRateTable{
			{
				ChildRateID: 2,
			},
		},
		Images: []image.PlanImagesInput{},
	},
	{
		TemaPlanTable: price.TemaPlanTable{
			PlanID:     -1,
			PropertyID: 5,
		},
		RoomTypeCode: "6",
		ChildRates: []price.ChildRateTable{
			{
				ChildRateID: 2,
			},
		},
		Images: []image.PlanImagesInput{
			{
				ImageID: 236,
				PlanID:  501,
				Order:   1,
			},
		},
	},
}
var flag = 0

// MockRoomBulkUseCase mock implementation
type MockPlanTemaBulkUseCase struct {
	mock.Mock
}

func (m *MockPlanTemaBulkUseCase) CreatePlanOwnImagesTema(images []image.HtTmPlanOwnImagesTemas) error {
	if images[0].PlanImageTemaID == 236 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) FetchRoomImagesByRoomTypeID(roomTypeIDList []int64) ([]image.RoomImagesOutput, error) {
	var result []image.RoomImagesOutput
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchImagesByPlanID(planIDList []int64) ([]image.PlanImagesOutput, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockPlanTemaBulkUseCase) FetchOnePlan(propertyID int64, packagePlanCode int) (*plan.HtTmPlanTemas, error) {
	result := &plan.HtTmPlanTemas{}
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchOne(roomTypeCode int, propertyID int64) (*room.HtTmRoomTemas, error) {
	result := &room.HtTmRoomTemas{}
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchOneWithPlanId(planID int64) (price.HtTmPlanTemas, error) {
	result := price.HtTmPlanTemas{}
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchList(propertyID int64, packagePlanCodeList []int) ([]plan.HtTmPlanTemas, error) {
	var result []plan.HtTmPlanTemas
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchAllByPropertyID(req plan.ListInput) ([]price.HtTmPlanTemas, error) {
	var result []price.HtTmPlanTemas
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) GetPlanIfPlanCodeExist(propertyID int64, planCode int64, roomTypeID int64) (price.HtTmPlanTemas, error) {

	if propertyID == 3 {
		return price.HtTmPlanTemas{
			TemaPlanTable: price.TemaPlanTable{
				PlanID:          2,
				PackagePlanCode: planCode,
				PropertyID:      propertyID,
				RoomTypeID:      roomTypeID,
			},
		}, nil
	}
	if propertyID == 5 {
		return price.HtTmPlanTemas{
			TemaPlanTable: price.TemaPlanTable{
				PlanID:          -1,
				PackagePlanCode: planCode,
				PropertyID:      propertyID,
				RoomTypeID:      roomTypeID,
			},
		}, nil
	}
	if propertyID == 4 {
		return price.HtTmPlanTemas{}, errors.New("new err")
	}
	return price.HtTmPlanTemas{
		TemaPlanTable: price.TemaPlanTable{
			PlanID:          0,
			PackagePlanCode: planCode,
			PropertyID:      propertyID,
			RoomTypeID:      roomTypeID,
		},
	}, nil
}

func (m *MockPlanTemaBulkUseCase) UpdatePlanBulkTema(planTable price.HtTmPlanTemas, planID int64) error {
	if planID == 2 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) GetNextPlanID() (price.HtTmPlanTemas, error) {
	if flag == 3 {
		flag = 0
		return price.HtTmPlanTemas{}, errors.New("new err")
	}

	return price.HtTmPlanTemas{
		TemaPlanTable: price.TemaPlanTable{
			PlanID: 500,
		},
	}, nil
}

func (m *MockPlanTemaBulkUseCase) CreatePlanBulkTema(planTable price.HtTmPlanTemas) error {
	if planTable.PlanID == 501 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) ClearChildRateTema(planID int64) error {
	if planID == 0 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) ClearImageTema(planID int64) error {
	if planID == 0 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) CreateChildRateTema(childRates []price.HtTmChildRateTemas) error {

	if childRates[0].ChildRateID == 0 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) DeletePlanTema(planCode int64, roomTypeIDs []int64) error {
	if planCode == 0 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) MatchesPlanIDAndPropertyID(planID int64, propertyID int64) bool {
	//TODO implement me
	return false
}

func (m *MockPlanTemaBulkUseCase) FetchActiveByPlanCode(planCode string) ([]price.HtTmPlanTemas, error) {
	//TODO implement me
	return []price.HtTmPlanTemas{}, nil
}

func (m *MockPlanTemaBulkUseCase) FetchChildRates(planID int64) ([]price.HtTmChildRateTemas, error) {
	var result []price.HtTmChildRateTemas
	return result, nil
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
	if flag == 1 || flag == 2 {
		flag = 0
		return errors.New("new err")
	}

	return nil
}

// TxRollback mock
func (m *MockPlanTemaBulkUseCase) TxRollback(tx *gorm.DB) {
	return
}

func (m *MockPlanTemaBulkUseCase) FetchListWithPropertyId(roomTypeCodeList []int, propertyID int64) ([]room.HtTmRoomTemas, error) {
	var result []room.HtTmRoomTemas
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchRoomTypeIdByRoomTypeCode(propertyID int64, roomTypeCode string) (room.HtTmRoomTypeTemas, error) {
	if propertyID == 2 {
		return room.HtTmRoomTypeTemas{}, errors.New("new err")
	}
	return room.HtTmRoomTypeTemas{
		RoomTypeTema: room.RoomTypeTema{
			RoomTypeID:   31,
			RoomTypeCode: roomTypeCode,
			PropertyID:   propertyID,
		},
	}, nil

}

func (m *MockPlanTemaBulkUseCase) CreateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error {

	if roomTable.PropertyID == 1208013 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) UpdateRoomBulkTema(roomTable *room.HtTmRoomTypeTemas) error {
	if roomTable.RoomTypeID == 37 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) MatchesRoomTypeIDAndPropertyID(roomTypeID int64, propertyID int64) bool {
	//TODO implement me
	return false
}

func (m *MockPlanTemaBulkUseCase) FetchRoomByRoomTypeID(roomTypeID int64) (*room.HtTmRoomTypeTemas, error) {
	result := &room.HtTmRoomTypeTemas{}
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) ClearRoomToAmenities(roomTypeID int64) error {
	if roomTypeID == 0 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) CreateRoomToAmenities(roomTypeID int64, tlRoomAmenityID int64) error {

	if roomTypeID == 34 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) ClearRoomImage(roomTypeID int64) error {

	if roomTypeID == 35 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) CreateRoomOwnImages(images []room.HtTmRoomOwnImagesTemas) error {
	if images[0].RoomTypeID == 36 {
		return errors.New("new err")
	}
	return nil
}

func (m *MockPlanTemaBulkUseCase) FetchRoomsByPropertyID(req room.ListInput) ([]room.HtTmRoomTypeTemas, error) {
	var result []room.HtTmRoomTypeTemas
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchAllAmenities() ([]room.HtTmRoomAmenityTemas, error) {
	var result []room.HtTmRoomAmenityTemas
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchAmenitiesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomAmenitiesTema, error) {
	var result []room.RoomAmenitiesTema
	return result, nil
}

func (m *MockPlanTemaBulkUseCase) FetchImagesByRoomTypeID(roomTypeIDList []int64) ([]room.RoomImagesTema, error) {
	var result []room.RoomImagesTema
	return result, nil
}
func TemaPlanBulkCreateDataProcess(reqData []price.TemaPlanData) (string, error) {
	mockUseCase := new(MockPlanTemaBulkUseCase)
	//mock required repositories with instance of tema use case
	useCases := &planUsecase.PlanTemaUsecase{
		PTemaRepository: mockUseCase,
		RTemaRepository: mockUseCase,
		ITemaRepository: mockUseCase,
	}
	return useCases.CreateBulk(reqData)
}

// request body data
func TestTemaPlanBulkCreateDataProcess(t *testing.T) {
	//using loop to run all the test data to successfully run all the success and error cases
	for _, data := range request {
		var reqArray = []price.TemaPlanData{
			data,
		}
		if data.RoomTypeCode == "2" {
			flag = 1
		}
		if data.RoomTypeCode == "3" {
			flag = 2
		}
		if data.TemaPlanTable.PlanID == 0 {
			flag = 3
		}
		if data.TemaPlanTable.PlanID == -1 {
			flag = 4
		}
		res, err := TemaPlanBulkCreateDataProcess(reqArray)
		if res != "" {
			assert.NotNilf(t, res, res)
			assert.Equal(t, res, res)
		} else {
			assert.Error(t, err)
		}

	}
}
