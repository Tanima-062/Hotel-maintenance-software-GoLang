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
	"strconv"
)

// PlanTemaUsecase Tl plan related usecase
type PlanTemaUsecase struct {
	PTemaRepository               plan.IPlanTemaRepository
	RTemaRepository               room.IRoomTemaRepository
	ITemaRepository               image.IImageTemaRepository
	ICommonCancelPolicyRepository cancelPolicy.ICancelPolicyCommonRepository
	ICommonPlanRepository         plan.ICommonPlanRepository
}

// NewplanTemaUsecase Instantiation
func NewPlanTemaUsecase(db *gorm.DB) plan.IPlanBulkTemaUsecase {
	return &PlanTemaUsecase{
		PTemaRepository:               pInfra.NewPlanTemaRepository(db),
		RTemaRepository:               rInfra.NewRoomTemaRepository(db),
		ITemaRepository:               iInfra.NewImageTemaRepository(db),
		ICommonCancelPolicyRepository: cpInfra.NewCommonCancelPolicyRepository(db),
		ICommonPlanRepository:         pInfra.NewPlanCommonRepository(db),
	}
}

// Bulk Create or Update Plan
func (p *PlanTemaUsecase) CreateBulk(request []price.TemaPlanData) (string, error) {
	// transaction generation
	tx, txErr := p.PTemaRepository.TxStart()
	if txErr != nil {
		log.Error(txErr)
		return "Something went wrong", txErr
	}
	existingPlans := make(map[int64][]int64)
	for i := range request {
		planTable := price.HtTmPlanTemas{
			TemaPlanTable: request[i].TemaPlanTable,
		}
		// Fetch RoomTypeID by RoomTypeCode
		roomTypeData, rErr := p.RTemaRepository.FetchRoomTypeIDByRoomTypeCode(planTable.TemaPlanTable.PropertyID, request[i].RoomTypeCode)
		if rErr != nil {
			log.Error(rErr)
			continue
		}
		roomTypeID := roomTypeData.RoomTypeTema.RoomTypeID
		planTable.TemaPlanTable.LangCd = "ja-JP"
		planTable.TemaPlanTable.RoomTypeID = roomTypeID

		// Get Plan Detail by PropertyID, PlanCode, roomTypeID
		planR, _ := p.PTemaRepository.GetPlanIfPlanCodeExist(request[i].PropertyID, request[i].TemaPlanTable.PackagePlanCode, roomTypeID)
		if planR.TemaPlanTable.PlanID > 0 {
			existingPlans[request[i].PackagePlanCode] = append(existingPlans[request[i].PackagePlanCode], roomTypeID)
			planTable.TemaPlanTable.PlanID = planR.TemaPlanTable.PlanID
			// Update plan
			if err := p.PTemaRepository.UpdatePlanBulkTema(planTable, planTable.TemaPlanTable.PlanID); err != nil {
				log.Error(err)
				continue
			}
		} else {
			// Create new plan
			planLastData, err := p.PTemaRepository.GetNextPlanID()
			if err != nil {
				planTable.TemaPlanTable.PlanID = 1
			} else {
				planTable.TemaPlanTable.PlanID = planLastData.PlanID + 1
			}
			// Create plan
			if err := p.PTemaRepository.CreatePlanBulkTema(planTable); err != nil {
				log.Error(err)
			}
		}

		// Registering Child Pricing
		// Delete existing child_rates
		if err := p.PTemaRepository.ClearChildRateTema(planR.TemaPlanTable.PlanID); err != nil {
			log.Error(err)
		}

		// Then insert
		childRateTables := []price.HtTmChildRateTemas{}

		for _, child := range request[i].ChildRates {
			fromAge, toAge := p.calculateAgeFromChildRateType(child.ChildRateType)
			childRateTables = append(childRateTables, price.HtTmChildRateTemas{
				ChildRateTable: price.ChildRateTable{
					PlanID:        planTable.TemaPlanTable.PlanID,
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

		if err := p.PTemaRepository.CreateChildRateTema(childRateTables); err != nil {
			log.Error(err)
		}

		// Attach image to plan
		// Delete existing Images
		if err := p.PTemaRepository.ClearImageTema(planR.TemaPlanTable.PlanID); err != nil {
			log.Error(err)
		} else {
			// Then insert
			for _, imageData := range request[i].Images {

				var record []image.HtTmPlanOwnImagesTemas
				record = append(record, image.HtTmPlanOwnImagesTemas{
					PlanImageTemaID: imageData.ImageID,
					PlanID:          planTable.TemaPlanTable.PlanID,
					Order:           imageData.Order,
				})
				if err := p.ITemaRepository.CreatePlanOwnImagesTema(record); err != nil {
					log.Error(err)
				}
			}
		}
	}

	for key, value := range existingPlans {
		// delete plan
		if err := p.PTemaRepository.DeletePlanTema(key, value); err != nil {
			log.Error(err)
		}
	}

	// commit and rollback
	if err := p.PTemaRepository.TxCommit(tx); err != nil {
		p.PTemaRepository.TxRollback(tx)
		log.Error(err)
		return "Something went wrong!", err
	}

	return "Bulk plan data processing complete successfully.", nil
}

// calculateAgeFromChildRateType
func (p *PlanTemaUsecase) calculateAgeFromChildRateType(childRateType int8) (int8, int8) {
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

// Detail Plan details
func (p *PlanTemaUsecase) Detail(request *plan.DetailInput) (*plan.TemaBulkDetailOutput, error) {
	response := &plan.TemaBulkDetailOutput{}

	if p.PTemaRepository.MatchesPlanIDAndPropertyID(request.PlanID, request.PropertyID) == false {
		return response, fmt.Errorf("Error: %s", "This plan cannot be viewed at this property.")
	}

	planCh := make(chan price.HtTmPlanTemas)
	childRatesCh := make(chan []price.HtTmChildRateTemas)
	planImageCh := make(chan []image.PlanImagesOutput)

	go p.fetchPlan(planCh, request.PlanID)
	go p.fetchChildRates(childRatesCh, request.PlanID)
	go p.fetchPlanImages(planImageCh, []price.HtTmPlanTemas{{TemaPlanTable: price.TemaPlanTable{PlanID: request.PlanID}}})

	planDetail, childRates, images := <-planCh, <-childRatesCh, <-planImageCh

	// Fetch room by roomTypeID
	roomData, rErr := p.RTemaRepository.FetchRoomByRoomTypeID(planDetail.RoomTypeID)
	if rErr != nil {
		return response, rErr
	}

	response.TemaPlanTable = planDetail.TemaPlanTable

	checkinStart := strconv.Itoa(planDetail.TemaPlanTable.CheckinTimeStartH) + ":" + strconv.Itoa(planDetail.TemaPlanTable.CheckinTimeStartM)
	checkinEnd := strconv.Itoa(planDetail.TemaPlanTable.CheckinTimeEndH) + ":" + strconv.Itoa(planDetail.TemaPlanTable.CheckinTimeEndM)
	checkout := strconv.Itoa(planDetail.TemaPlanTable.CheckoutTimeEndH) + ":" + strconv.Itoa(planDetail.TemaPlanTable.CheckoutTimeEndM)

	response.TemaPlanTable.CheckinStart = checkinStart
	response.TemaPlanTable.CheckinEnd = checkinEnd
	response.TemaPlanTable.Checkout = checkout

	response.Name = roomData.Name
	for _, childRate := range childRates {
		response.ChildRates = append(response.ChildRates, childRate.ChildRateTable)
	}
	response.Images = images

	return response, nil
}

// FetchList Get plan list
func (p *PlanTemaUsecase) FetchList(request *plan.ListInput) ([]plan.TemaBulkListOutput, error) {
	response := []plan.TemaBulkListOutput{}
	roomCh := make(chan []room.HtTmRoomTypeTemas)
	planCh := make(chan []price.HtTmPlanTemas)
	go p.fetchRooms(roomCh, request.PropertyID)
	go p.fetchPlans(planCh, request.PropertyID)
	rooms, plans := <-roomCh, <-planCh

	roomImageCh := make(chan []image.RoomImagesOutput)
	planImageCh := make(chan []image.PlanImagesOutput)
	go p.fetchRoomImages(roomImageCh, rooms)
	go p.fetchPlanImages(planImageCh, plans)
	roomImages, planImages := <-roomImageCh, <-planImageCh

	for _, roomData := range rooms {
		record := &plan.TemaBulkListOutput{}
		// Link multiple plans to one room
		for _, planData := range plans {
			var temp plan.TemaBulkDetailOutput

			if roomData.RoomTypeTema.RoomTypeID == planData.TemaPlanTable.RoomTypeID && planData.TemaPlanTable.PlanID != 0 {
				temp.TemaPlanTable = planData.TemaPlanTable

				for _, planImage := range planImages {
					if planImage.PlanID == planData.TemaPlanTable.PlanID && planImage.Order == 1 {
						temp.Images = append(temp.Images, planImage)
						break
					}
				}

				record.Plans = append(record.Plans, temp)
			}
		}

		// If there is no plan, do not return the room information
		if len(record.Plans) == 0 {
			continue
		}
		record.RoomTypeID = roomData.RoomTypeTema.RoomTypeID
		record.RoomName = roomData.RoomTypeTema.Name
		record.RoomIsStopSales = roomData.RoomTypeTema.IsStopSales

		// roomTypeId set matching images
		for _, roomImage := range roomImages {
			if roomImage.RoomTypeID == roomData.RoomTypeTema.RoomTypeID && roomImage.Order == 1 {
				record.RoomImageHref = roomImage.Href
				break
			}
		}
		response = append(response, *record)

	}
	return response, nil
}
func (p *PlanTemaUsecase) fetchRooms(ch chan<- []room.HtTmRoomTypeTemas, propertyID int64) {
	// Fetch all rooms by propertyID
	rooms, roomErr := p.RTemaRepository.FetchRoomsByPropertyID(room.ListInput{PropertyID: propertyID})
	if roomErr != nil {
		ch <- []room.HtTmRoomTypeTemas{}
	}
	ch <- rooms
}
func (p *PlanTemaUsecase) fetchPlans(ch chan<- []price.HtTmPlanTemas, propertyID int64) {
	// Fetch all plans by propertyID
	plans, planErr := p.PTemaRepository.FetchAllByPropertyID(plan.ListInput{PropertyID: propertyID})
	if planErr != nil {
		ch <- []price.HtTmPlanTemas{}
	}
	ch <- plans
}
func (p *PlanTemaUsecase) fetchRoomImages(ch chan<- []image.RoomImagesOutput, rooms []room.HtTmRoomTypeTemas) {
	var roomIDList []int64
	for _, roomData := range rooms {
		roomIDList = append(roomIDList, roomData.RoomTypeID)
	}
	// Fetch room images by roomTypeID
	images, imageErr := p.ITemaRepository.FetchRoomImagesByRoomTypeID(roomIDList)
	if imageErr != nil {
		ch <- []image.RoomImagesOutput{}
	}
	ch <- images
}
func (p *PlanTemaUsecase) fetchPlanImages(ch chan<- []image.PlanImagesOutput, plans []price.HtTmPlanTemas) {
	var planIDList []int64
	for _, planData := range plans {
		planIDList = append(planIDList, planData.TemaPlanTable.PlanID)
	}
	// Fetch images by planID
	images, imageErr := p.ITemaRepository.FetchImagesByPlanID(planIDList)
	if imageErr != nil {
		ch <- []image.PlanImagesOutput{}
	}
	ch <- images
}
func (p *PlanTemaUsecase) fetchPlan(ch chan<- price.HtTmPlanTemas, planID int64) {
	// Fetch one plan with planID
	planData, planErr := p.PTemaRepository.FetchOneWithPlanID(planID)
	if planErr != nil {
		ch <- price.HtTmPlanTemas{}
	}
	ch <- planData
}

func (p *PlanTemaUsecase) fetchChildRates(ch chan<- []price.HtTmChildRateTemas, planID int64) {
	// Fetch child rates by planID
	childRates, planErr := p.PTemaRepository.FetchChildRates(planID)
	if planErr != nil {
		ch <- []price.HtTmChildRateTemas{}
	}
	ch <- childRates
}
