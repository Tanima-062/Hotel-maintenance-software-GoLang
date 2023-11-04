package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	planInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	priceInfra "github.com/Adventureinc/hotel-hm-api/src/price/infra"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"time"
)

// priceTemaUsecase Tema price related usecase
type PriceTemaUsecase struct {
	PriceTemaRepository price.IPriceTemaRepository
	PlanTemaRepository  plan.IPlanTemaRepository
}

// NewPriceTemaUsecase instantiation
func NewPriceTemaUsecase(db *gorm.DB) price.IPriceBulkTemaUsecase {
	return &PriceTemaUsecase{
		PriceTemaRepository: priceInfra.NewPriceTemaRepository(db),
		PlanTemaRepository:  planInfra.NewPlanTemaRepository(db),
	}
}

// Update price data from bulk request
func (p *PriceTemaUsecase) Update(request []price.PriceTemaData) (string, error) {
	// transaction generation
	tx, txErr := p.PriceTemaRepository.TxStart()
	if txErr != nil {
		return "", txErr
	}

	for _, requestData := range request {
		for priceDate, priceData := range requestData.PriceList {
			roomTypeCode, _ := strconv.Atoi(requestData.RoomTypeCode)
			// delete existing price
			err := p.PriceTemaRepository.DeletePriceTema(requestData.PropertyID, requestData.PackagePlanCode, roomTypeCode, priceDate)
			if err != nil {
				log.Error(err)
			}
			date, _ := time.Parse("2006-01-02", priceDate)
			priceTable := price.HtTmPriceTemas{
				PriceTemaTable: price.PriceTemaTable{
					PriceDate:       date,
					PropertyID:      requestData.PropertyID,
					RoomTypeCode:    roomTypeCode,
					PackagePlanCode: requestData.PackagePlanCode,
					Disable:         requestData.Disable,
					MinPax:          requestData.MinPax,
					MaxPax:          requestData.MaxPax,
				},
			}
			priceType := price.TemaPriceType{}
			pt := reflect.ValueOf(&priceType).Elem().Type()
			for i := range priceData {
				field := pt.Field(i)
				rv := reflect.ValueOf(&priceType)
				fieldName := reflect.Indirect(rv).FieldByName(field.Name)
				fieldName.SetInt(priceData[i].Price)
				priceTable.TemaPriceType = priceType
			}
			if err := p.PriceTemaRepository.CreatePrice(priceTable); err != nil {
				log.Error(err)
			}
		}
	}
	var retErr error = nil
	errMsg := ""
	// commit and rollback
	if err := p.PriceTemaRepository.TxCommit(tx); err != nil {
		p.PlanTemaRepository.TxRollback(tx)
		retErr = err
		errMsg = "Something went wrong!"
	}
	errMsg = "Price has been updated"
	return errMsg, retErr
}
