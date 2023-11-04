package usecase

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy"
	"github.com/Adventureinc/hotel-hm-api/src/cancelPolicy/infra"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	fInfra "github.com/Adventureinc/hotel-hm-api/src/facility/infra"
	"gorm.io/gorm"
)

// cancelPolicyTlUsecase TLのキャンセルポリシー関連usecase
type cancelPolicyTlUsecase struct {
	CTlRepository     cancelPolicy.ICancelPolicyTlRepository
	FTlRepository     facility.IFacilityTlRepository
	CCommonRepository cancelPolicy.ICancelPolicyCommonRepository
}

// NewCancelPolicyTlUsecase インスタンス生成
func NewCancelPolicyTlUsecase(db *gorm.DB) cancelPolicy.ICancelPolicyUsecase {
	return &cancelPolicyTlUsecase{
		CTlRepository:     infra.NewCancelPolicyTlRepository(db),
		FTlRepository:     fInfra.NewFacilityTlRepository(db),
		CCommonRepository: infra.NewCommonCancelPolicyRepository(db),
	}
}

// Detail キャンセルポリシー詳細
func (c *cancelPolicyTlUsecase) Detail(req *cancelPolicy.DetailInput) (*cancelPolicy.CancelPolicyJSONWithName, error) {
	response := &cancelPolicy.CancelPolicyJSONWithName{}
	facilityData, err := c.FTlRepository.FetchPropertyByPropertyID(*req.PropertyID)
	if err != nil {
		return response, err
	}
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
}

// Save キャンセルポリシー保存
func (c *cancelPolicyTlUsecase) Save(req *cancelPolicy.UpdateInput) error {
	cancelPolicy := req.CancelPolicyJSON
	jsonData, jErr := json.Marshal(cancelPolicy)
	if jErr != nil {
		return jErr
	}
	return c.CTlRepository.UpsertCancelPolicyTl(*req.PropertyID, string(jsonData))
}

// Create キャンセルポリシー作成 (プランごとのキャンセルポリシーは未対応のため)
func (c *cancelPolicyTlUsecase) Create(req *cancelPolicy.CreateInput) error {
	return fmt.Errorf("TLリンカーンはプランごとのキャンセルポリシーに対応していません")
}

// List キャンセルポリシー一覧返却
func (c *cancelPolicyTlUsecase) List(req *cancelPolicy.ListInput) ([]cancelPolicy.CancelPolicyInfo, error) {
	return []cancelPolicy.CancelPolicyInfo{}, fmt.Errorf("TLリンカーンはプランごとのキャンセルポリシーに対応していません")
}

// Delete プランごとのキャンセルポリシー削除
func (c *cancelPolicyTlUsecase) Delete(req *cancelPolicy.DeleteInput) error {
	return fmt.Errorf("TLリンカーンはプランごとのキャンセルポリシーに対応していません")
}

// PlanList プラン一覧返却
func (c *cancelPolicyTlUsecase) PlanList(req *cancelPolicy.PlanListInput) ([]cancelPolicy.PlanInfo, error) {
	return []cancelPolicy.PlanInfo{}, fmt.Errorf("TLリンカーンはプランごとのキャンセルポリシーに対応していません")
}
