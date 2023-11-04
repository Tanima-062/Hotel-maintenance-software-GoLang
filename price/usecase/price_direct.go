package usecase

import (
	"math"
	"strconv"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	planInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	priceInfra "github.com/Adventureinc/hotel-hm-api/src/price/infra"
	"gorm.io/gorm"
)

// priceDirectUsecase 直仕入れ料金関連usecase
type priceDirectUsecase struct {
	PriceDirectRepository price.IPriceDirectRepository
	PlanDirectRepository  plan.IPlanDirectRepository
}

// NewPriceDirectUsecase インスタンス生成
func NewPriceDirectUsecase(db *gorm.DB) price.IPriceUsecase {
	return &priceDirectUsecase{
		PriceDirectRepository: priceInfra.NewPriceDirectRepository(db),
		PlanDirectRepository:  planInfra.NewPlanDirectRepository(db),
	}
}

// FetchDetail 一定期間の料金データを取得
func (p *priceDirectUsecase) FetchDetail(request *price.DetailInput) (price.DetailOutput, error) {
	response := price.DetailOutput{}
	// 日付データは月初月末前後+6日間取得する
	baseStartDate, _ := time.Parse("2006-01-02", request.StartDate)
	baseEndDate, _ := time.Parse("2006-01-02", request.EndDate)
	startDate := time.Date(baseStartDate.Year(), baseStartDate.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, -6).Format("2006-01-02")
	endDate := time.Date(baseEndDate.Year(), baseEndDate.Month(), 1, 0, 0, 0, 0, time.Local).AddDate(0, 1, 6).Format("2006-01-02")
	prices, pErr := p.PriceDirectRepository.FetchPricesWithInThePeriod(request.PlanID, startDate, endDate)
	if pErr != nil {
		return response, pErr
	}

	response.PlanID = request.PlanID
	tempPrices := map[string][]price.Price{}
	for _, priceData := range prices {
		// 人数
		numberOfPeople, _ := strconv.Atoi(priceData.RateTypeCode)
		priceDate := priceData.UseDate.Format("2006-01-02")
		priceInTax := int(float64(priceData.PriceInTax) / float64(numberOfPeople))
		tempPrices[priceDate] = append(tempPrices[priceDate],
			price.Price{Type: priceData.RateTypeCode, Price: priceInTax})
	}
	response.Prices = tempPrices
	return response, nil
}

// Save 料金データを作成・更新
func (p *priceDirectUsecase) Save(request *[]price.SaveInput) error {
	var inputData []price.HtTmPriceDirects
	var planIdList []int64
	for _, planPrices := range *request {
		planIdList = append(planIdList, planPrices.PlanID)
	}
	// 税区分の為、プラン情報取得
	planList, pErr := p.PlanDirectRepository.FetchList(planIdList)
	if pErr != nil {
		return pErr
	}
	// 子供料金設定情報取得
	childRates, cErr := p.PriceDirectRepository.FetchChildRatesByPlanIDList(planIdList)
	if cErr != nil {
		return cErr
	}

	for _, planPrices := range *request {
		planData := plan.HtTmPlanDirects{}
		// PlanIDが一致するプランを格納
		for _, plan := range planList {
			if planPrices.PlanID == plan.PlanID {
				planData = plan
				break
			}
		}
		childRateData := []price.HtTmChildRateDirects{}
		// PlanIDが一致するchildRateのスライス作成
		for _, childRate := range childRates {
			if planPrices.PlanID == childRate.PlanID {
				childRateData = append(childRateData, childRate)
			}
		}
		for date, input := range planPrices.Prices {
			for _, priceData := range input {
				childPrices := p.settingChildPrices(childRateData, planData, priceData)
				AdultPrice, AdultPriceInTax := p.settingTax(planData.TaxCategory, priceData.Price)
				numberOfPeople, _ := strconv.Atoi(priceData.Type)
				AdultPrice = AdultPrice * numberOfPeople
				AdultPriceInTax = AdultPriceInTax * numberOfPeople
				parsedDate, _ := time.Parse("2006-01-02", date)

				inputData = append(inputData, price.HtTmPriceDirects{
					PriceTable: price.PriceTable{
						PlanID:           planPrices.PlanID,
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
				})
			}
		}
	}

	tx, err := p.PriceDirectRepository.TxStart()
	if err != nil {
		return err
	}

	pRepository := priceInfra.NewPriceDirectRepository(tx)

	if err := pRepository.UpsertPrices(inputData); err != nil {
		pRepository.TxRollback(tx)
		return err
	}

	return pRepository.TxCommit(tx)
}

func (p *priceDirectUsecase) settingTax(taxCategory bool, price int) (int, int) {
	if taxCategory == true {
		priceInTax := float64(price) * 1.1
		return price, int(priceInTax)
	}
	priceNotInTax := float64(price) / 11 * 10
	return int(priceNotInTax), price
}

func (p *priceDirectUsecase) calcChildRate(rateCategory int8, price int, rate int) int {
	switch rateCategory {
	case 0: // 率
		return int(float64(price) * (float64(rate) / 100))
	case 1: // 固定金額
		return rate
	case 2: // 円引き
		return int(math.Max(float64(price-rate), 0))
	}
	return 0
}

func (p *priceDirectUsecase) settingChildPrices(childRates []price.HtTmChildRateDirects, planData plan.HtTmPlanDirects, priceData price.Price) []int {
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

	// 人数
	numberOfPeople, _ := strconv.Atoi(priceData.Type)

	for _, childRate := range childRates {
		// 子供の種別（小学生とか）に応じて、料金単位（円、円引き、％）、人数、税金の計算をおこなう
		if childRate.ChildRateType == 1 {
			temp := p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice1, childPrice1InTax = p.settingTax(planData.TaxCategory, temp)
			childPrice1InTax = childPrice1InTax * numberOfPeople
			childPrice1 = childPrice1 * numberOfPeople
		}
		if childRate.ChildRateType == 2 {
			temp := p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice2, childPrice2InTax = p.settingTax(planData.TaxCategory, temp)
			childPrice2InTax = childPrice2InTax * numberOfPeople
			childPrice2 = childPrice2 * numberOfPeople
		}
		if childRate.ChildRateType == 3 {
			temp := p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice3, childPrice3InTax = p.settingTax(planData.TaxCategory, temp)
			childPrice3InTax = childPrice3InTax * numberOfPeople
			childPrice3 = childPrice3 * numberOfPeople
		}
		if childRate.ChildRateType == 4 {
			temp := p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice4, childPrice4InTax = p.settingTax(planData.TaxCategory, temp)
			childPrice4InTax = childPrice4InTax * numberOfPeople
			childPrice4 = childPrice4 * numberOfPeople
		}
		if childRate.ChildRateType == 5 {
			temp := p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice5, childPrice5InTax = p.settingTax(planData.TaxCategory, temp)
			childPrice5InTax = childPrice5InTax * numberOfPeople
			childPrice5 = childPrice5 * numberOfPeople
		}
		if childRate.ChildRateType == 6 {
			temp := p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice6, childPrice6InTax = p.settingTax(planData.TaxCategory, temp)
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
