package usecase

import (
	"fmt"
	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	cpInfra "github.com/Adventureinc/hotel-hm-api/src/cancelPolicy/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	iInfra "github.com/Adventureinc/hotel-hm-api/src/image/infra"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	pInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"math"
	"strconv"
)

// planTlUsecase Tl plan related usecase
type planTlUsecase struct {
	PTlRepository                 plan.IPlanTlRepository
	RTlRepository                 room.IRoomTlRepository
	ITlRepository                 image.IImageTlRepository
	ICommonCancelPolicyRepository cancelPolicy.ICancelPolicyCommonRepository
	ICommonPlanRepository         plan.ICommonPlanRepository
}

// NewPlanTLUsecase Instantiation
func NewPlanTlUsecase(db *gorm.DB) plan.IPlanBulkUsecase {
	return &planTlUsecase{
		PTlRepository:                 pInfra.NewPlanTlRepository(db),
		RTlRepository:                 rInfra.NewRoomTlRepository(db),
		ITlRepository:                 iInfra.NewImageTlRepository(db),
		ICommonCancelPolicyRepository: cpInfra.NewCommonCancelPolicyRepository(db),
		ICommonPlanRepository:         pInfra.NewPlanCommonRepository(db),
	}
}

// FetchList Get plan list
func (p *planTlUsecase) FetchList(request *plan.ListInput) ([]plan.BulkListOutput, error) {
	response := []plan.BulkListOutput{}
	roomCh := make(chan []room.HtTmRoomTypeTls)
	planCh := make(chan []price.HtTmPlanTls)
	go p.fetchRooms(roomCh, request.PropertyID)
	go p.fetchPlans(planCh, request.PropertyID)
	rooms, plans := <-roomCh, <-planCh

	roomImageCh := make(chan []image.RoomImagesOutput)
	planImageCh := make(chan []image.PlanImagesOutput)
	go p.fetchRoomImages(roomImageCh, rooms)
	go p.fetchPlanImages(planImageCh, plans)
	roomImages, planImages := <-roomImageCh, <-planImageCh

	for _, roomData := range rooms {
		record := &plan.BulkListOutput{}
		// Link multiple plans to one room
		for _, planData := range plans {
			var temp plan.BulkDetailOutput
			if roomData.RoomTypeID == planData.RoomTypeID {
				temp.PlanTable = planData.PlanTable
			}
			for _, planImage := range planImages {
				if planImage.PlanID == planData.PlanID && planImage.Order == 1 {
					temp.Images = append(temp.Images, planImage)
					break
				}
			}
			if temp.PlanTable.PlanID != 0 {
				record.Plans = append(record.Plans, temp)
			}
		}

		// If there is no plan, do not return the room information
		if len(record.Plans) == 0 {
			continue
		}
		record.RoomTypeID = roomData.RoomTypeTable.RoomTypeID
		record.RoomName = roomData.RoomTypeTable.Name
		record.RoomIsStopSales = roomData.RoomTypeTable.IsStopSales

		// roomTypeId set matching images
		for _, roomImage := range roomImages {
			if roomImage.RoomTypeID == roomData.RoomTypeID && roomImage.Order == 1 {
				record.RoomImageHref = roomImage.Href
				break
			}
		}
		response = append(response, *record)
	}
	return response, nil
}

// Detail Plan details
func (p *planTlUsecase) Detail(request *plan.DetailInput) (*plan.BulkDetailOutput, error) {
	response := &plan.BulkDetailOutput{}

	if p.PTlRepository.MatchesPlanIDAndPropertyID(request.PlanID, request.PropertyID) == false {
		return response, fmt.Errorf("Error: %s", "This plan cannot be viewed at this property.")
	}

	planCh := make(chan price.HtTmPlanTls)
	childRatesCh := make(chan []price.HtTmChildRateTls)
	planImageCh := make(chan []image.PlanImagesOutput)
	checkInOutCh := make(chan *plan.HtTmPlanCheckInOuts)

	go p.fetchPlan(planCh, request.PlanID)
	go p.fetchChildRates(childRatesCh, request.PlanID)
	go p.fetchPlanImages(planImageCh, []price.HtTmPlanTls{{PlanTable: price.PlanTable{PlanID: request.PlanID}}})

	planDetail, childRates, images, checkInOut := <-planCh, <-childRatesCh, <-planImageCh, <-checkInOutCh

	roomData, rErr := p.RTlRepository.FetchRoomByRoomTypeID(planDetail.RoomTypeID)
	if rErr != nil {
		return response, rErr
	}

	activePlanTables, err := p.PTlRepository.FetchActiveByPlanCode(planDetail.PlanTable.PlanCode)
	if err != nil {
		return response, err
	}
	var activeRooms []int64
	for _, planTable := range activePlanTables {
		activeRooms = append(activeRooms, planTable.RoomTypeID)
	}

	response.PlanTable = planDetail.PlanTable
	response.RoomName = roomData.Name
	response.ActiveRooms = activeRooms
	for _, childRate := range childRates {
		response.ChildRates = append(response.ChildRates, childRate.ChildRateTable)
	}
	response.Images = images
	if checkInOut != nil {
		response.CheckinStart = checkInOut.CheckInBegin
		response.CheckinEnd = checkInOut.CheckInEnd
		response.Checkout = checkInOut.CheckOut
	}

	return response, nil
}

func (p *planTlUsecase) fetchRooms(ch chan<- []room.HtTmRoomTypeTls, propertyID int64) {
	rooms, roomErr := p.RTlRepository.FetchRoomsByPropertyID(room.ListInput{PropertyID: propertyID})
	if roomErr != nil {
		ch <- []room.HtTmRoomTypeTls{}
	}
	ch <- rooms
}

func (p *planTlUsecase) fetchPlans(ch chan<- []price.HtTmPlanTls, propertyID int64) {
	plans, planErr := p.PTlRepository.FetchAllByPropertyID(plan.ListInput{PropertyID: propertyID})
	if planErr != nil {
		ch <- []price.HtTmPlanTls{}
	}
	ch <- plans
}

func (p *planTlUsecase) fetchRoomImages(ch chan<- []image.RoomImagesOutput, rooms []room.HtTmRoomTypeTls) {
	var roomIDList []int64
	for _, roomData := range rooms {
		roomIDList = append(roomIDList, roomData.RoomTypeID)
	}
	images, imageErr := p.ITlRepository.FetchImagesByRoomTypeID(roomIDList)
	if imageErr != nil {
		ch <- []image.RoomImagesOutput{}
	}
	ch <- images
}

func (p *planTlUsecase) fetchPlanImages(ch chan<- []image.PlanImagesOutput, plans []price.HtTmPlanTls) {
	var planIDList []int64
	for _, planData := range plans {
		planIDList = append(planIDList, planData.PlanID)
	}
	images, imageErr := p.ITlRepository.FetchImagesByPlanID(planIDList)
	if imageErr != nil {
		ch <- []image.PlanImagesOutput{}
	}
	ch <- images
}

func (p *planTlUsecase) fetchPlan(ch chan<- price.HtTmPlanTls, planID int64) {
	planData, planErr := p.PTlRepository.FetchOne(planID)
	if planErr != nil {
		ch <- price.HtTmPlanTls{}
	}
	ch <- planData
}

func (p *planTlUsecase) fetchChildRates(ch chan<- []price.HtTmChildRateTls, planID int64) {
	childRates, planErr := p.PTlRepository.FetchChildRates(planID)
	if planErr != nil {
		ch <- []price.HtTmChildRateTls{}
	}
	ch <- childRates
}

func (p *planTlUsecase) calculateAgeFromChildRateType(childRateType int8) (int8, int8) {
	switch childRateType {
	case utils.ChildRateTypeA:
		return 9, 11
	case utils.ChildRateTypeB:
		return 6, 8
	case utils.ChildRateTypeC:
		return 0, 5
	case utils.ChildRateTypeD:
		return 0, 5
	case utils.ChildRateTypeE:
		return 0, 5
	case utils.ChildRateTypeF:
		return 0, 5
	}
	return 0, 0
}

// calcChildRate Fee Unit for Child Fare
func (p *planTlUsecase) calcChildRate(rateCategory int8, price int, rate int) int {
	switch rateCategory {
	case 0: // rate
		return int(float64(price) * (float64(rate) / 100))
	case 1: // fixed amount
		return rate
	case 2: // Yen discount
		return int(math.Max(float64(price-rate), 0))
	}
	return 0
}

// settingChildPrices Calculating Child Fares
func (p *planTlUsecase) settingChildPrices(childRates []price.HtTmChildRateTls, priceData price.Price) []int {
	childPrice1 := 0
	childPrice1InTax := 0
	childPrice2 := 0
	childPrice2InTax := 0
	childPrice3 := 0
	childPrice3InTax := 0
	childPrice4 := 0
	childPrice4InTax := 0
	childPrice5 := 0
	childPrice5InTax := 0
	childPrice6 := 0
	childPrice6InTax := 0

	// Number of people
	numberOfPeople, _ := strconv.Atoi(priceData.Type)

	for _, childRate := range childRates {
		// Calculate the unit of charge (yen, discounted yen, %), number of people, and tax according to the type of child (elementary school student, etc.)
		if childRate.ChildRateType == utils.ChildRateTypeA {
			childPrice1InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice1 = int(float64(childPrice1InTax) / 11 * 10)
			childPrice1InTax = childPrice1InTax * numberOfPeople
			childPrice1 = childPrice1 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeB {
			childPrice2InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice2 = int(float64(childPrice2InTax) / 11 * 10)
			childPrice2InTax = childPrice2InTax * numberOfPeople
			childPrice2 = childPrice2 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeC {
			childPrice3InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice3 = int(float64(childPrice3InTax) / 11 * 10)
			childPrice3InTax = childPrice3InTax * numberOfPeople
			childPrice3 = childPrice3 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeD {
			childPrice4InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice4 = int(float64(childPrice4InTax) / 11 * 10)
			childPrice4InTax = childPrice4InTax * numberOfPeople
			childPrice4 = childPrice4 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeE {
			childPrice5InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice5 = int(float64(childPrice5InTax) / 11 * 10)
			childPrice5InTax = childPrice5InTax * numberOfPeople
			childPrice5 = childPrice5 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeF {
			childPrice6InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice6 = int(float64(childPrice6InTax) / 11 * 10)
			childPrice6InTax = childPrice6InTax * numberOfPeople
			childPrice6 = childPrice6 * numberOfPeople
		}
	}
	return []int{
		childPrice1,
		childPrice1InTax,
		childPrice2,
		childPrice2InTax,
		childPrice3,
		childPrice3InTax,
		childPrice4,
		childPrice4InTax,
		childPrice5,
		childPrice5InTax,
		childPrice6,
		childPrice6InTax,
	}
}

// Bulk Create or Update Plan
func (p *planTlUsecase) CreateBulk(request []price.PlanData) (string, error) {
	// transaction generation
	tx, txErr := p.PTlRepository.TxStart()
	if txErr != nil {
		log.Error(txErr)
		return "Something went wrong", txErr
	}
	existingPlans := make(map[string][]int64)
	for i := range request {
		planTable := price.HtTmPlanTls{
			PlanTable: request[i].PlanTable,
		}
		roomTypeData, rErr := p.RTlRepository.FetchRoomTypeIdByRoomTypeCode(planTable.PlanTable.PropertyID, request[i].RoomTypeCode)
		if rErr != nil {
			log.Error(rErr)
			continue
		}
		roomTypeID := roomTypeData.RoomTypeTable.RoomTypeID
		planTable.PlanTable.LangCd = "ja-JP"
		planTable.PlanTable.RoomTypeID = roomTypeID

		// Get Plan Detail by PropertyID, PlanCode, roomTypeID
		planR, _ := p.PTlRepository.GetPlanIfPlanCodeExist(request[i].PropertyID, request[i].PlanCode, roomTypeID)
		if planR.PlanTable.PlanID > 0 {
			existingPlans[request[i].PlanCode] = append(existingPlans[request[i].PlanCode], roomTypeID)
			// Update plan
			planTable.PlanTable.PlanID = planR.PlanTable.PlanID
			if err := p.PTlRepository.UpdatePlanBulkTl(planTable, planTable.PlanTable.PlanID); err != nil {
				log.Error(err)
				continue
			}
		} else {
			// Create new plan
			planLastData, err := p.PTlRepository.GetNextPlanID()
			if err != nil {
				planTable.PlanTable.PlanID = 1
			} else {
				planTable.PlanTable.PlanID = planLastData.PlanID + 1
			}

			if err := p.PTlRepository.CreatePlanBulkTl(planTable); err != nil {
				log.Error(err)
			}
		}

		// Registering Child Pricing
		// Delete existing child_rates
		if err := p.PTlRepository.ClearChildRateTl(planR.PlanTable.PlanID); err != nil {
			log.Error(err)
		}

		// Then insert
		childRateTables := []price.HtTmChildRateTls{}

		for _, child := range request[i].ChildRates {
			fromAge, toAge := p.calculateAgeFromChildRateType(child.ChildRateType)
			childRateTables = append(childRateTables, price.HtTmChildRateTls{
				ChildRateTable: price.ChildRateTable{
					PlanID:        planTable.PlanTable.PlanID,
					ChildRateType: child.ChildRateType,
					FromAge:       fromAge,
					ToAge:         toAge,
					Receive:       child.Receive,
					RateCategory:  child.RateCategory,
					Rate:          child.Rate,
					CalcCategory:  child.CalcCategory,
				},
			})
		}

		if err := p.PTlRepository.CreateChildRateTl(childRateTables); err != nil {
			log.Error(err)
		}

		// Attach image to plan
		// Delete existing Images
		if err := p.PTlRepository.ClearImageTl(planR.PlanTable.PlanID); err != nil {
			log.Error(err)
		} else {
			// Then insert
			for _, imageData := range request[i].Images {
				var record []image.HtTmPlanOwnImagesTls
				record = append(record, image.HtTmPlanOwnImagesTls{
					PlanImageTlID: imageData.ImageID,
					PlanID:        planTable.PlanTable.PlanID,
					Order:         imageData.Order,
				})
				if err := p.ITlRepository.CreatePlanOwnImagesTl(record); err != nil {
					log.Error(err)
				}
			}
		}

		// Save plan check-in/out times
		info := plan.CheckInOutInfo{
			WholesalerID: utils.WholesalerIDTl,
			PropertyID:   request[i].PlanTable.PropertyID,
			PlanID:       planTable.PlanTable.PlanID,
			CheckInBegin: request[i].CheckinStart,
			CheckInEnd:   request[i].CheckinEnd,
			CheckOut:     request[i].Checkout,
		}

		if err := p.ICommonPlanRepository.UpsertCheckInOut(info); err != nil {
			log.Error(err)
		}

	}

	for key, value := range existingPlans {
		if err := p.PTlRepository.DeletePlanTl(key, value); err != nil {
			log.Error(err)
		}
	}

	// commit and rollback
	if err := p.PTlRepository.TxCommit(tx); err != nil {
		p.PTlRepository.TxRollback(tx)
		log.Error(err)
		return "Something went wrong!", err
	}

	return "Bulk plan data processing complete successfully.", nil
}
