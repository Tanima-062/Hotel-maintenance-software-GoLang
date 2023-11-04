package usecase

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	cpInfra "github.com/Adventureinc/hotel-hm-api/src/cancelPolicy/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	iInfra "github.com/Adventureinc/hotel-hm-api/src/image/infra"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	pInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"github.com/Adventureinc/hotel-hm-api/src/price"
	priceInfra "github.com/Adventureinc/hotel-hm-api/src/price/infra"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"gorm.io/gorm"
)

// planNeppanUsecase ねっぱんプラン関連usecase
type planNeppanUsecase struct {
	PNeppanRepository             plan.IPlanNeppanRepository
	RNeppanRepository             room.IRoomNeppanRepository
	INeppanRepository             image.IImageNeppanRepository
	ICommonCancelPolicyRepository cancelPolicy.ICancelPolicyCommonRepository
	ICommonPlanRepository         plan.ICommonPlanRepository
}

// NewPlanNeppanUsecase インスタンス生成
func NewPlanNeppanUsecase(db *gorm.DB) plan.IPlanUsecase {
	return &planNeppanUsecase{
		PNeppanRepository:             pInfra.NewPlanNeppanRepository(db),
		RNeppanRepository:             rInfra.NewRoomNeppanRepository(db),
		INeppanRepository:             iInfra.NewImageNeppanRepository(db),
		ICommonCancelPolicyRepository: cpInfra.NewCommonCancelPolicyRepository(db),
		ICommonPlanRepository:         pInfra.NewPlanCommonRepository(db),
	}
}

// FetchList プラン一覧取得
func (p *planNeppanUsecase) FetchList(request *plan.ListInput) ([]plan.ListOutput, error) {
	response := []plan.ListOutput{}
	roomCh := make(chan []room.HtTmRoomTypeNeppans)
	planCh := make(chan []plan.HtTmPlanNeppans)
	go p.fetchRooms(roomCh, request.PropertyID)
	go p.fetchPlans(planCh, request.PropertyID)
	rooms, plans := <-roomCh, <-planCh

	roomImageCh := make(chan []image.RoomImagesOutput)
	planImageCh := make(chan []image.PlanImagesOutput)
	go p.fetchRoomImages(roomImageCh, rooms)
	go p.fetchPlanImages(planImageCh, plans)
	roomImages, planImages := <-roomImageCh, <-planImageCh

	for _, roomData := range rooms {
		record := &plan.ListOutput{}
		// 1つの部屋に複数のプランを紐付かせる
		for _, planData := range plans {
			var temp plan.DetailOutput
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

		// プランが一つもなかったら、その部屋情報を返さない
		if len(record.Plans) == 0 {
			continue
		}
		record.RoomTypeID = roomData.RoomTypeTable.RoomTypeID
		record.RoomName = roomData.RoomTypeTable.Name
		record.RoomIsStopSales = roomData.RoomTypeTable.IsStopSales

		// roomTypeIdが一致する画像を設定
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

// Detail プラン詳細
func (p *planNeppanUsecase) Detail(request *plan.DetailInput) (*plan.DetailOutput, error) {
	response := &plan.DetailOutput{}

	if p.PNeppanRepository.MatchesPlanIDAndPropertyID(request.PlanID, request.PropertyID) == false {
		return response, fmt.Errorf("Error: %s", "この施設ではこのプランを閲覧できません。")
	}

	planCh := make(chan plan.HtTmPlanNeppans)
	childRatesCh := make(chan []plan.HtTmChildRateNeppans)
	planImageCh := make(chan []image.PlanImagesOutput)
	cancelPolicyCh := make(chan *cancelPolicy.HtThPlanCancelPolicyRelations)
	checkInOutCh := make(chan *plan.HtTmPlanCheckInOuts)

	go p.fetchPlan(planCh, request.PlanID)
	go p.fetchChildRates(childRatesCh, request.PlanID)
	go p.fetchPlanImages(planImageCh, []plan.HtTmPlanNeppans{{PlanTable: plan.PlanTable{PlanID: request.PlanID}}})
	go p.fetchPlanCancelPolicy(cancelPolicyCh, request.PropertyID, request.PlanID)
	go p.fetchCheckInOut(checkInOutCh, request.PropertyID, request.PlanID)

	planDetail, childRates, images, assignedCancelPolicy, checkInOut := <-planCh, <-childRatesCh, <-planImageCh, <-cancelPolicyCh, <-checkInOutCh

	roomData, rErr := p.RNeppanRepository.FetchRoomByRoomTypeID(planDetail.RoomTypeID)
	if rErr != nil {
		return response, rErr
	}

	activePlanTables, err := p.PNeppanRepository.FetchActiveByPlanGroupID(planDetail.PlanGroupID)
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
	if assignedCancelPolicy != nil {
		response.PlanCancelPolicyID = &assignedCancelPolicy.PlanCancelPolicyID
	}
	if checkInOut != nil {
		response.CheckinStart = checkInOut.CheckInBegin
		response.CheckinEnd = checkInOut.CheckInEnd
		response.Checkout = checkInOut.CheckOut
	}

	return response, nil
}

// createPlan プラン作成の共通処理
func (p *planNeppanUsecase) createPlan(request *plan.SaveInput, tx *gorm.DB) error {
	// plan_codeのチェック
	// 同じ部屋に紐づく同じプランコードがあった場合は重複エラーとする
	planCodeList := []plan.CheckDuplicatePlanCode{}
	for _, roomTypeID := range request.SelectedRooms {
		planCodeList = append(planCodeList, plan.CheckDuplicatePlanCode{RoomTypeID: roomTypeID, PlanCode: request.PlanCode})
	}

	duplicate := p.PNeppanRepository.CheckPlanCode(request.PropertyID, planCodeList)
	if duplicate > 0 {
		return fmt.Errorf("DuplicateError")
	}

	planTxRepo := pInfra.NewPlanNeppanRepository(tx)
	roomTxRepo := rInfra.NewRoomNeppanRepository(tx)

	// room_type_idの数だけプランを新規作成
	planTables := []plan.HtTmPlanNeppans{}
	for _, roomTypeID := range request.SelectedRooms {
		// 部屋IDの確認と、売止フラグを取得するため
		roomData, rErr := roomTxRepo.FetchRoomByRoomTypeID(roomTypeID)
		if rErr != nil {
			return rErr
		}
		tempPlan := request.PlanTable
		tempPlan.LangCd = "ja-JP"
		tempPlan.RoomTypeID = roomTypeID
		tempPlan.IsStopSales = roomData.IsStopSales
		planTables = append(planTables, plan.HtTmPlanNeppans{
			PlanTable: tempPlan,
		})
	}
	if request.PlanGroupID == 0  {
		// プランの新規作成
		if err := planTxRepo.CreatePlansNeppan(planTables); err != nil {
			return err
		}
	} else {
		// 既存のプラングループへのプラン追加
		if err := planTxRepo.MakePlansNeppan(planTables); err != nil {
			return err
		}
	}

	// 子供料金設定の登録
	childRateTables := []plan.HtTmChildRateNeppans{}
	for _, child := range request.ChildRates {
		for _, createdPlan := range planTables {
			fromAge, toAge := p.calculateAgeFromChildRateType(child.ChildRateType)
			childRateTables = append(childRateTables, plan.HtTmChildRateNeppans{ChildRateTable: price.ChildRateTable{
				PlanID:        createdPlan.PlanID,
				ChildRateType: child.ChildRateType,
				FromAge:       fromAge,
				ToAge:         toAge,
				Receive:       child.Receive,
				RateCategory:  child.RateCategory,
				Rate:          child.Rate,
				CalcCategory:  child.CalcCategory,
			}})
		}
	}
	if err := planTxRepo.CreateChildRateNeppan(childRateTables); err != nil {
		return err
	}

	// プランに画像を紐付ける
	imageTxRepo := iInfra.NewImageNeppanRepository(tx)
	for _, imageData := range request.Images {
		record := []image.HtTmPlanOwnImagesNeppans{}
		for _, createdPlan := range planTables {
			record = append(record, image.HtTmPlanOwnImagesNeppans{
				RoomImageNeppanID: imageData.ImageID,
				PlanID:            createdPlan.PlanID,
				Order:             imageData.Order,
			})
		}
		if err := imageTxRepo.CreatePlanOwnImagesNeppan(record); err != nil {
			return err
		}
	}

	// キャンセルポリシーを紐づける
	cpTxRepo := cpInfra.NewCommonCancelPolicyRepository(tx)
	// 紐付けする場合のみrequest.PlanCancelPolicyIdに値が入る。紐付けない場合はnilが入る。
	if request.PlanCancelPolicyId != nil {
		for _, createdPlan := range planTables {
			if err := cpTxRepo.UpsertPlanCancelPolicyRelation(utils.WholesalerIDNeppan, createdPlan.PropertyID, createdPlan.PlanID, *request.PlanCancelPolicyId); err != nil {
				return err
			}
		}
	}

	// プランのチェックイン/アウト時間を保存する
	cioTxRepo := pInfra.NewPlanCommonRepository(tx)
	for _, createdPlan := range planTables {
		info := plan.CheckInOutInfo{
			WholesalerID: utils.WholesalerIDNeppan,
			PropertyID:   request.PropertyID,
			PlanID:       createdPlan.PlanID,
			CheckInBegin: request.CheckinStart,
			CheckInEnd:   request.CheckinEnd,
			CheckOut:     request.Checkout,
		}

		if err := cioTxRepo.UpsertCheckInOut(info); err != nil {
			return err
		}
	}
	return nil
}

// Create プラン作成
func (p *planNeppanUsecase) Create(request *plan.SaveInput) error {
	// トランザクション生成
	tx, txErr := p.PNeppanRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	// プラン作成
	if err := p.createPlan(request, tx); err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}

	// コミットとロールバック
	if err := p.PNeppanRepository.TxCommit(tx); err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}
	return nil
}

// Update プラン更新
func (p *planNeppanUsecase) Update(request *plan.SaveInput) error {
	// トランザクション生成
	tx, txErr := p.PNeppanRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	planTxRepo := pInfra.NewPlanNeppanRepository(tx)
	priceTxRepo := priceInfra.NewPriceNeppanRepository(tx)

	// 全件のプラン取得
	allPlanTables, err := planTxRepo.FetchAllByPlanGroupID(request.PlanGroupID)
	if err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}

	// プランの追加リストと更新リストを作成
	var updatePlanIDs []int64
	var additionalRoomTypeIDs []int64
	for _, selectedRoomTypeID := range request.SelectedRooms {
		addFlag := true
		for _, planTable := range allPlanTables {
			if planTable.RoomTypeID == selectedRoomTypeID {
				updatePlanIDs = append(updatePlanIDs, planTable.PlanID)
				addFlag = false
				break
			}
		}
		if addFlag == true {
			additionalRoomTypeIDs = append(additionalRoomTypeIDs, selectedRoomTypeID)
		}
	}

	// プランの追加
	if len(additionalRoomTypeIDs) > 0 {
		additionalRequest := *request
		additionalRequest.SelectedRooms = additionalRoomTypeIDs
		additionalRequest.PlanID = 0
		if err := p.createPlan(&additionalRequest, tx); err != nil{
			p.PNeppanRepository.TxRollback(tx)
			return err
		}
	}

	// プランの更新 
	if len(updatePlanIDs) > 0 {
		planTable := &plan.HtTmPlanNeppans{PlanTable: request.PlanTable}
		planTable.IsDelete = false
		if err := planTxRepo.UpdatePlanNeppan(planTable, updatePlanIDs); err != nil {
			p.PNeppanRepository.TxRollback(tx)
			return err
		}

		// 子供料金設定の更新
		for _, updatePlanID := range updatePlanIDs {
			for _, child := range request.ChildRates {
				childRateTable := &plan.HtTmChildRateNeppans{ChildRateTable: price.ChildRateTable{
					ChildRateType: child.ChildRateType,
					PlanID:        updatePlanID,
					Receive:       child.Receive,
					RateCategory:  child.RateCategory,
					Rate:          child.Rate,
					CalcCategory:  child.CalcCategory,
				}}
				if err := planTxRepo.UpdateChildRateNeppan(childRateTable); err != nil {
					p.PNeppanRepository.TxRollback(tx)
					return err
				}
			}
		}

		// 子供料金設定を料金に反映
		childRates, childRateErr := priceTxRepo.FetchChildRates(updatePlanIDs[0])
		if childRateErr != nil {
			p.PNeppanRepository.TxRollback(tx)
			return childRateErr
		}
		for _, updatePlanID := range updatePlanIDs {
			prices, _ := priceTxRepo.FetchPricesByPlanID(updatePlanID)
			if len(prices) > 0 {
				var inputData []price.HtTmPriceNeppans
				for _, priceData := range prices {
					// 2人以上の大人料金を参照する場合は複数人分の大人料金から子供料金の割引計算をしてしまうので、大人料金1人分から子供料金を計算をする
					// 参照する料金の人数
					numberOfPeople, _ := strconv.Atoi(priceData.RateTypeCode)
					// 人数分で割った大人料金
					priceInTax := priceData.PriceInTax / numberOfPeople
					// 人数分で割った大人料金から子供料金を計算する
					childPrices := p.settingChildPrices(childRates, *planTable, price.Price{Type: priceData.RateTypeCode, Price: priceInTax})
					inputData = append(inputData, price.HtTmPriceNeppans{
						PriceTable: price.PriceTable{
							PriceID:          priceData.PriceID,
							PlanID:           updatePlanID,
							UseDate:          priceData.UseDate,
							RateTypeCode:     priceData.RateTypeCode,
							Price:            priceData.Price,
							PriceInTax:       priceData.PriceInTax,
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
				if updateErr := priceTxRepo.UpdateChildPrices(inputData); updateErr != nil {
					p.PNeppanRepository.TxRollback(tx)
					return updateErr
				}
			} // end of if
		} // end of for

		// 画像を一度削除して、部屋と画像を再度紐付ける
		imageTxRepo := iInfra.NewImageNeppanRepository(tx)
		for _, updatePlanID := range updatePlanIDs {
			if err := imageTxRepo.ClearPlanImage(updatePlanID); err != nil {
				p.PNeppanRepository.TxRollback(tx)
				return err
			}
			for _, imageData := range request.Images {
				record := []image.HtTmPlanOwnImagesNeppans{}
				record = append(record, image.HtTmPlanOwnImagesNeppans{
					RoomImageNeppanID: imageData.ImageID,
					PlanID:            updatePlanID,
					Order:             imageData.Order,
				})
				if err := imageTxRepo.CreatePlanOwnImagesNeppan(record); err != nil {
					p.PNeppanRepository.TxRollback(tx)
					return err
				}
			}
		}
		
		// キャンセルポリシーを紐づける
		cpTxRepo := cpInfra.NewCommonCancelPolicyRepository(tx)
		for _, updatePlanID := range updatePlanIDs {
			// 紐付けする場合のみrequest.PlanCancelPolicyIdに値が入る。紐付けない場合はnilが入る。
			if request.PlanCancelPolicyId != nil {
				if err := cpTxRepo.UpsertPlanCancelPolicyRelation(utils.WholesalerIDNeppan, planTable.PropertyID, updatePlanID, *request.PlanCancelPolicyId); err != nil {
					p.PNeppanRepository.TxRollback(tx)
					return err
				}
			} else {
				if err := cpTxRepo.DeletePlanCancelPolicyRelation(utils.WholesalerIDNeppan, updatePlanID); err != nil {
					p.PNeppanRepository.TxRollback(tx)
					return err
				}
			}
		}
		
		// プランのチェックイン/アウト時間を保存する
		cioTxRepo := pInfra.NewPlanCommonRepository(tx)
		for _, updatePlanID := range updatePlanIDs {
			info := plan.CheckInOutInfo{
				WholesalerID: utils.WholesalerIDNeppan,
				PropertyID:   request.PropertyID,
				PlanID:       updatePlanID,
				CheckInBegin: request.CheckinStart,
				CheckInEnd:   request.CheckinEnd,
				CheckOut:     request.Checkout,
			}
			if err := cioTxRepo.UpsertCheckInOut(info); err != nil {
				p.PNeppanRepository.TxRollback(tx)
				return err
			}
		}
	}

	// プランの削除リストを作成
	var deletePlanIDs []int64
	for _, planTable := range allPlanTables {
		var deleteFlag = true
		for _, selectedRoomTypeID := range request.SelectedRooms {
			if selectedRoomTypeID == planTable.RoomTypeID {
				deleteFlag = false
				break
			}
		}
		if deleteFlag == true {
			deletePlanIDs = append(deletePlanIDs, planTable.PlanID)
		}
	}
	// プランの削除
	if len(deletePlanIDs) > 0 {
		if err := planTxRepo.DeletePlanNeppan(deletePlanIDs); err != nil {
			p.PNeppanRepository.TxRollback(tx)
			return err
		}
	}

	// コミットとロールバック
	if err := p.PNeppanRepository.TxCommit(tx); err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}
	return nil
}

// Delete プラン削除
func (p *planNeppanUsecase) Delete(planID int64) error {
	// トランザクション生成
	tx, txErr := p.PNeppanRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	planTxRepo := pInfra.NewPlanNeppanRepository(tx)

	// アクティブなプランを全件取得
	planGroupID, err := planTxRepo.FetchPlanGroupID(planID)
	if err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}
	activePlanTables, err := planTxRepo.FetchActiveByPlanGroupID(planGroupID)
	if err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}

	// 削除するプランIDのリスト作成
	var deletePlanIDs []int64
	for _, planTable := range activePlanTables {
		deletePlanIDs = append(deletePlanIDs, planTable.PlanID)
	}

	if err := planTxRepo.DeletePlanNeppan(deletePlanIDs); err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}

	cpTxRepo := cpInfra.NewCommonCancelPolicyRepository(tx)
	if err := cpTxRepo.DeletePlanCancelPolicyRelation(utils.WholesalerIDNeppan, planID); err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}

	cioTxRepo := pInfra.NewPlanCommonRepository(tx)
	if err := cioTxRepo.DeleteCheckInOut(utils.WholesalerIDNeppan, planID); err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}

	// コミットとロールバック
	if err := p.PNeppanRepository.TxCommit(tx); err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}

	return nil
}

// UpdateStopSales プラン売止
func (p *planNeppanUsecase) UpdateStopSales(request *plan.StopSalesInput) error {

	// トランザクション生成
	tx, txErr := p.PNeppanRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	planTxRepo := pInfra.NewPlanNeppanRepository(tx)
	planTxRepo.UpdateStopSales([]int64{request.PlanID}, request.IsStopSales)

	// 部屋に紐づくプランが一つでも販売中に戻されたら、部屋の売止も解除する
	if request.IsStopSales == false {
		planData, pErr := planTxRepo.FetchOne(request.PlanID)
		if pErr != nil {
			p.PNeppanRepository.TxRollback(tx)
			return pErr
		}
		roomTxRepo := rInfra.NewRoomNeppanRepository(tx)
		if err := roomTxRepo.UpdateStopSales(planData.RoomTypeID, request.IsStopSales); err != nil {
			p.PNeppanRepository.TxRollback(tx)
			return err
		}
	}

	// コミットとロールバック
	if err := p.PNeppanRepository.TxCommit(tx); err != nil {
		p.PNeppanRepository.TxRollback(tx)
		return err
	}
	return nil
}

func (p *planNeppanUsecase) fetchRooms(ch chan<- []room.HtTmRoomTypeNeppans, propertyID int64) {
	rooms, roomErr := p.RNeppanRepository.FetchRoomsByPropertyID(room.ListInput{PropertyID: propertyID})
	if roomErr != nil {
		ch <- []room.HtTmRoomTypeNeppans{}
	}
	ch <- rooms
}

func (p *planNeppanUsecase) fetchPlans(ch chan<- []plan.HtTmPlanNeppans, propertyID int64) {
	plans, planErr := p.PNeppanRepository.FetchAllByPropertyID(plan.ListInput{PropertyID: propertyID})
	if planErr != nil {
		ch <- []plan.HtTmPlanNeppans{}
	}
	ch <- plans
}

func (p *planNeppanUsecase) fetchRoomImages(ch chan<- []image.RoomImagesOutput, rooms []room.HtTmRoomTypeNeppans) {
	var roomIDList []int64
	for _, roomData := range rooms {
		roomIDList = append(roomIDList, roomData.RoomTypeID)
	}
	images, imageErr := p.INeppanRepository.FetchImagesByRoomTypeID(roomIDList)
	if imageErr != nil {
		ch <- []image.RoomImagesOutput{}
	}
	ch <- images
}

func (p *planNeppanUsecase) fetchPlanImages(ch chan<- []image.PlanImagesOutput, plans []plan.HtTmPlanNeppans) {
	var planIDList []int64
	for _, planData := range plans {
		planIDList = append(planIDList, planData.PlanID)
	}
	images, imageErr := p.INeppanRepository.FetchImagesByPlanID(planIDList)
	if imageErr != nil {
		ch <- []image.PlanImagesOutput{}
	}
	ch <- images
}

func (p *planNeppanUsecase) fetchPlan(ch chan<- plan.HtTmPlanNeppans, planID int64) {
	planData, planErr := p.PNeppanRepository.FetchOne(planID)
	if planErr != nil {
		ch <- plan.HtTmPlanNeppans{}
	}
	ch <- planData
}

func (p *planNeppanUsecase) fetchChildRates(ch chan<- []plan.HtTmChildRateNeppans, planID int64) {
	childRates, planErr := p.PNeppanRepository.FetchChildRates(planID)
	if planErr != nil {
		ch <- []plan.HtTmChildRateNeppans{}
	}
	ch <- childRates
}

func (p *planNeppanUsecase) fetchPlanCancelPolicy(ch chan<- *cancelPolicy.HtThPlanCancelPolicyRelations, propertyID int64, planID int64) {
	assignedPolicy, err := p.ICommonCancelPolicyRepository.FindAssignedPlanCancelPolicy(propertyID, planID)
	if err != nil {
		ch <- nil
	}

	ch <- assignedPolicy
}

func (p *planNeppanUsecase) fetchCheckInOut(ch chan<- *plan.HtTmPlanCheckInOuts, propertyID int64, planID int64) {
	checkInOut, err := p.ICommonPlanRepository.FetchCheckInOut(propertyID, planID)
	if err != nil {
		ch <- nil
	}

	ch <- checkInOut
}

func (p *planNeppanUsecase) calculateAgeFromChildRateType(childRateType int8) (int8, int8) {
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

// calcChildRate　子供料金の料金単位
func (p *planNeppanUsecase) calcChildRate(rateCategory int8, price int, rate int) int {
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

// settingChildPrices 子供料金の計算
func (p *planNeppanUsecase) settingChildPrices(childRates []price.HtTmChildRateNeppans, planData plan.HtTmPlanNeppans, priceData price.Price) []int {
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
			childPrice1InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice1 = int(float64(childPrice1InTax) / 11 * 10)
			childPrice1InTax = childPrice1InTax * numberOfPeople
			childPrice1 = childPrice1 * numberOfPeople
		}
		if childRate.ChildRateType == 2 {
			childPrice2InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice2 = int(float64(childPrice2InTax) / 11 * 10)
			childPrice2InTax = childPrice2InTax * numberOfPeople
			childPrice2 = childPrice2 * numberOfPeople
		}
		if childRate.ChildRateType == 3 {
			childPrice3InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice3 = int(float64(childPrice3InTax) / 11 * 10)
			childPrice3InTax = childPrice3InTax * numberOfPeople
			childPrice3 = childPrice3 * numberOfPeople
		}
		if childRate.ChildRateType == 4 {
			childPrice4InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice4 = int(float64(childPrice4InTax) / 11 * 10)
			childPrice4InTax = childPrice4InTax * numberOfPeople
			childPrice4 = childPrice4 * numberOfPeople
		}
		if childRate.ChildRateType == 5 {
			childPrice5InTax = p.calcChildRate(childRate.RateCategory, priceData.Price, childRate.Rate)
			childPrice5 = int(float64(childPrice5InTax) / 11 * 10)
			childPrice5InTax = childPrice5InTax * numberOfPeople
			childPrice5 = childPrice5 * numberOfPeople
		}
		if childRate.ChildRateType == 6 {
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