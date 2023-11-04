package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	fInfra "github.com/Adventureinc/hotel-hm-api/src/facility/infra"
	"github.com/Adventureinc/hotel-hm-api/src/plan"
	pInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"gorm.io/gorm"
)

// cancelPolicyRaku2Usecase らく通のキャンセルポリシー関連usecase
type cancelPolicyRaku2Usecase struct {
	CRaku2Repository  cancelPolicy.ICancelPolicyRaku2Repository
	FRaku2Repository  facility.IFacilityRaku2Repository
	PRaku2Repository  plan.IPlanRaku2Repository
	CCommonRepository cancelPolicy.ICancelPolicyCommonRepository
}

// NewCancelPolicyRaku2Usecase インスタンス生成
func NewCancelPolicyRaku2Usecase(db *gorm.DB) cancelPolicy.ICancelPolicyUsecase {
	return &cancelPolicyRaku2Usecase{
		CRaku2Repository:  infra.NewCancelPolicyRaku2Repository(db),
		FRaku2Repository:  fInfra.NewFacilityRaku2Repository(db),
		PRaku2Repository:  pInfra.NewPlanRaku2Repository(db),
		CCommonRepository: infra.NewCommonCancelPolicyRepository(db),
	}
}

// Detail キャンセルポリシー詳細
func (c *cancelPolicyRaku2Usecase) Detail(req *cancelPolicy.DetailInput) (*cancelPolicy.CancelPolicyJSONWithName, error) {
	response := &cancelPolicy.CancelPolicyJSONWithName{}
	if req.PropertyID != nil {
		facilityData, err := c.FRaku2Repository.FirstOrCreate(*req.PropertyID)
		if err != nil {
			return response, err
		}

		response.CancelPolicyName = nil
		// 施設のデフォルトキャンセルポリシーが未登録の場合、未登録であることを示す意味で当日キャンセルとノーショーはuint8最大値の255を返却
		if facilityData.CancelPenaltyJSON == "" {
			response.CancelPolicyJSON.Settings.CaseOfCancellationToday.Rate = 255
			response.CancelPolicyJSON.Settings.CaseOfNoShow.Rate = 255
			return response, nil
		}
		// マッチする条件を指定
		repRate := regexp.MustCompile("{\"Rate\":\"(\\d{1,})\"}")
		// 文字列{"Rate":"[一回以上繰り返しの数値]"}にマッチする場合、数値を囲っている""を削除して書き換える
		if isMatchRate := repRate.MatchString(facilityData.CancelPenaltyJSON); isMatchRate {
			facilityData.CancelPenaltyJSON = repRate.ReplaceAllString(facilityData.CancelPenaltyJSON, "{\"Rate\":$1}")
		}
		// マッチする条件を指定
		repDeposit := regexp.MustCompile("\"Deposit\":\"(\\d{1,})\"")
		// 文字列"Deposit":"[一回以上繰り返しの数値]"にマッチする場合、数値を囲っている""を削除して書き換える
		if isMatchDeposit := repDeposit.MatchString(facilityData.CancelPenaltyJSON); isMatchDeposit {
			facilityData.CancelPenaltyJSON = repDeposit.ReplaceAllString(facilityData.CancelPenaltyJSON, "\"Deposit\":$1")
		}
		jsonData := []byte(facilityData.CancelPenaltyJSON)
		if jErr := json.Unmarshal(jsonData, &response.CancelPolicyJSON); jErr != nil {
			return response, jErr
		}
		return response, nil

	} else if req.PlanCancelPolicyID != nil {
		ret, err := c.CCommonRepository.FindPlanCancelPolicy(*req.PlanCancelPolicyID)
		if err != nil {
			return response, err
		}

		response.CancelPolicyName = ret.CancelPolicyName
		response.CancelPolicyJSON = ret.CancelPolicyJSON

	} else {
		return nil, errors.New("Invalid Argument")
	}

	return response, nil
}

// Save キャンセルポリシー保存
func (c *cancelPolicyRaku2Usecase) Save(req *cancelPolicy.UpdateInput) error {
	cancelPolicy := req.CancelPolicyJSON
	jsonData, jErr := json.Marshal(cancelPolicy)
	if jErr != nil {
		return jErr
	}

	var err error
	if req.PlanCancelPolicyID != nil {
		err = c.CCommonRepository.UpdatePlanCancelPolicy(*req.PlanCancelPolicyID, req.PolicyName, string(jsonData))
	} else if req.PropertyID != nil {
		err = c.CRaku2Repository.Update(*req.PropertyID, string(jsonData))
	} else {
		err = fmt.Errorf("[cancelPolicyDirectUsecase] invalid argument")
	}

	return err
}

// Create キャンセルポリシー新規作成
func (c *cancelPolicyRaku2Usecase) Create(req *cancelPolicy.CreateInput) error {
	cancelPolicy := req.CancelPolicyJSON
	jsonData, jErr := json.Marshal(cancelPolicy)
	if jErr != nil {
		return jErr
	}

	return c.CCommonRepository.CreatePlanCancelPolicy(req.CancelPolicyName, utils.WholesalerIDRaku2, req.PropertyID, string(jsonData))
}

// List キャンセルポリシー一覧返却
func (c *cancelPolicyRaku2Usecase) List(req *cancelPolicy.ListInput) ([]cancelPolicy.CancelPolicyInfo, error) {
	list, err := c.CCommonRepository.PlanCancelPolicyList(req.PropertyID)
	if err != nil {
		return []cancelPolicy.CancelPolicyInfo{}, err
	}

	var ret []cancelPolicy.CancelPolicyInfo
	for _, v := range list {
		ret = append(ret, cancelPolicy.CancelPolicyInfo{PolicyID: v.PlanCancelPolicyID, PolicyName: v.PolicyName})
	}

	return ret, nil
}

// Delete プランごとのキャンセルポリシー削除
func (c *cancelPolicyRaku2Usecase) Delete(req *cancelPolicy.DeleteInput) error {
	return c.CCommonRepository.DeletePlanCancelPolicy(req.PlanCancelPolicyID)
}

// PlanList プラン一覧返却
func (c *cancelPolicyRaku2Usecase) PlanList(req *cancelPolicy.PlanListInput) ([]cancelPolicy.PlanInfo, error) {
	list, err := c.PRaku2Repository.FetchAllByCancelPolicyID(req.PlanCancelPolicyID)
	if err != nil {
		return []cancelPolicy.PlanInfo{}, err
	}

	var ret []cancelPolicy.PlanInfo
	for _, v := range list {
		ret = append(ret, cancelPolicy.PlanInfo{PlanID: v.PlanID, PlanName: v.Name})
	}

	return ret, nil
}
