package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/plan"
	"github.com/Adventureinc/hotel-hm-api/src/price"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	planInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	priceInfra "github.com/Adventureinc/hotel-hm-api/src/price/infra"
	"github.com/labstack/gommon/log"

	"math"
	"strconv"

	"gorm.io/gorm"
)

// priceTlUsecase Tl price related usecase
type priceTlUsecase struct {
	PriceTlRepository price.IPriceTlRepository
	PlanTlRepository  plan.IPlanTlRepository
}

// NewPriceTlUsecase instantiation
func NewPriceTlUsecase(db *gorm.DB) price.IPriceBulkTlUsecase {
	return &priceTlUsecase{
		PriceTlRepository: priceInfra.NewPriceTlRepository(db),
		PlanTlRepository:  planInfra.NewPlanTlRepository(db),
	}
}

// Update price data from bulk request
func (p *priceTlUsecase) Update(request []price.PriceData) (string, error) {
	// transaction generation
	tx, txErr := p.PriceTlRepository.TxStart()
	if txErr != nil {
		return "", txErr
	}

	for _, requestData := range request {
		planResult, err := p.PlanTlRepository.GetPlanByPropertyIDAndPlanCodeAndRoomTypeCode(requestData.PropertyID, requestData.PlanCode, requestData.RoomTypeCode)
		if err != nil {
			log.Error(err)
			continue
		}

		// loop over plan data
		for _, planResultData := range planResult {
			planResultData.PlanTable.PublishingStartDate = requestData.PublishingStartDate
			planResultData.PlanTable.PublishingEndDate = requestData.PublishingEndDate
			planResultData.PlanTable.IsPublishedYearRound = requestData.IsPublishedYearRound
			if err := p.PlanTlRepository.UpdatePlanBulkTl(planResultData, planResultData.PlanTable.PlanID); err != nil {
				log.Error(err)
			}
			childRateTables, _ := p.PlanTlRepository.FetchChildRates(planResultData.PlanTable.PlanID)
			for useDate, priceData := range requestData.Prices {
				for i := range priceData {
					priceDataF := p.GetPriceData(planResultData.PlanTable, childRateTables, priceData[i], useDate)
					// check if price exists by planID and rateTypeCode and useDate and rateTypeCode
					isFound, _ := p.PriceTlRepository.CheckIfPriceExistsByPlanIDAndRateTypeCodeAndUseDate(planResultData.PlanTable.PlanID, priceData[i].Type, useDate)
					if isFound {
						// update price and rate_type_code for existing useDate
						if err := p.PriceTlRepository.UpdatePrice(
							planResultData.PlanTable.PlanID,
							useDate,
							priceData[i].Type,
							priceData[i].Price,
							priceData[i].IsStopSales,
							priceDataF,
						); err != nil {
							log.Error(err)
							continue
						}
					} else {
						// create price and rate_type_code
						if err := p.PriceTlRepository.CreatePrice(priceDataF); err != nil {
							log.Error(err)
						}
					}
				}
			}
		}
	}

	// commit and rollback
	if err := p.PriceTlRepository.TxCommit(tx); err != nil {
		p.PlanTlRepository.TxRollback(tx)
		return "Something went wrong!", err
	}

	return "Price has been updated", nil
}

// Save price data
func (p *priceTlUsecase) GetPriceData(
	request price.PlanTable,
	childRateTables []price.HtTmChildRateTls,
	priceData price.Price,
	date string,
) price.HtTmPriceTls {
	childPrices := p.settingChildPrices(childRateTables, request.TaxCategory, priceData.Price, priceData.Type)
	AdultPrice, AdultPriceInTax := p.settingTax(request.TaxCategory, priceData.Price)
	numberOfPeople, _ := strconv.Atoi(priceData.Type)
	AdultPrice = AdultPrice * numberOfPeople
	AdultPriceInTax = AdultPriceInTax * numberOfPeople
	parsedDate, _ := time.Parse("2006-01-02", date)

	inputData := price.HtTmPriceTls{
		PriceTable: price.PriceTable{
			PlanID:           request.PlanID,
			UseDate:          parsedDate,
			RateTypeCode:     priceData.Type,
			Price:            AdultPrice,
			PriceInTax:       AdultPriceInTax,
			ChildPrice1:      childPrices[0],
			ChildPrice1InTax: childPrices[1],
			ChildPrice2:      childPrices[2],
			ChildPrice2InTax: childPrices[3],
			ChildPrice3:      childPrices[4],
			ChildPrice3InTax: childPrices[5],
			ChildPrice4:      childPrices[6],
			ChildPrice4InTax: childPrices[7],
			ChildPrice5:      childPrices[8],
			ChildPrice5InTax: childPrices[9],
			ChildPrice6:      childPrices[10],
			ChildPrice6InTax: childPrices[11],
			Times: common.Times{
				UpdatedAt: time.Now(),
			},
		},
	}

	return inputData
}

func (p *priceTlUsecase) settingTax(taxCategory bool, price int) (int, int) {
	if taxCategory == true {
		priceInTax := float64(price) * 1.1
		return price, int(priceInTax)
	}
	priceNotInTax := float64(price) / 11 * 10
	return int(priceNotInTax), price
}

func (p *priceTlUsecase) calcChildRate(rateCategory int8, price int, rate int) int {
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

func (p *priceTlUsecase) settingChildPrices(childRates []price.HtTmChildRateTls, taxCategory bool, priceData int, priceType string) []int {
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
	numberOfPeople, _ := strconv.Atoi(priceType)

	for _, childRate := range childRates {
		// Calculate the unit of charge (yen, discounted yen, %), number of people, and tax according to the type of
		// child (elementary school student, etc.)
		if childRate.ChildRateType == utils.ChildRateTypeA {
			temp := p.calcChildRate(childRate.RateCategory, priceData, childRate.Rate)
			childPrice1, childPrice1InTax = p.settingTax(taxCategory, temp)
			childPrice1InTax = childPrice1InTax * numberOfPeople
			childPrice1 = childPrice1 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeB {
			temp := p.calcChildRate(childRate.RateCategory, priceData, childRate.Rate)
			childPrice2, childPrice2InTax = p.settingTax(taxCategory, temp)
			childPrice2InTax = childPrice2InTax * numberOfPeople
			childPrice2 = childPrice2 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeC {
			temp := p.calcChildRate(childRate.RateCategory, priceData, childRate.Rate)
			childPrice3, childPrice3InTax = p.settingTax(taxCategory, temp)
			childPrice3InTax = childPrice3InTax * numberOfPeople
			childPrice3 = childPrice3 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeD {
			temp := p.calcChildRate(childRate.RateCategory, priceData, childRate.Rate)
			childPrice4, childPrice4InTax = p.settingTax(taxCategory, temp)
			childPrice4InTax = childPrice4InTax * numberOfPeople
			childPrice4 = childPrice4 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeE {
			temp := p.calcChildRate(childRate.RateCategory, priceData, childRate.Rate)
			childPrice5, childPrice5InTax = p.settingTax(taxCategory, temp)
			childPrice5InTax = childPrice5InTax * numberOfPeople
			childPrice5 = childPrice5 * numberOfPeople
		}
		if childRate.ChildRateType == utils.ChildRateTypeF {
			temp := p.calcChildRate(childRate.RateCategory, priceData, childRate.Rate)
			childPrice6, childPrice6InTax = p.settingTax(taxCategory, temp)
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
