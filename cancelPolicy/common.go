package cancelPolicy

// PlanCancelPolicyJSON プランごとに設定するキャンセルポリシー
type CancelPolicyJSONWithName struct {
	CancelPolicyName *string          `json:"CancelPolicyName"`
	CancelPolicyJSON CancelPolicyJSON `json:"CancelPolicyJSON"`
}

// CancelPolicyJSON JSONで格納されているキャンセルポリシーデータのRoot
type CancelPolicyJSON struct {
	Settings Settings `json:"Settings"`
}

type CancelPolicyInfo struct {
	PolicyID   uint64 `json:"policy_id"`
	PolicyName string `json:"policy_name"`
}

type PlanInfo struct {
	PlanID   int64  `json:"plan_id"`
	PlanName string `json:"plan_name"`
}

type Settings struct {
	CaseOfCancellationToday CaseOfCancellationToday `json:"CaseOfCancellationToday"`
	CaseOfNoShow            CaseOfNoShow            `json:"CaseOfNoShow"`
	NonRefundable           int                     `json:"NonRefundable"`
	Deposit                 int                     `json:"Deposit"`
	AdditionalCases         []AdditionalCases       `json:"AdditionalCases"`
}
type CaseOfCancellationToday struct {
	Rate uint8 `json:"Rate"`
}
type CaseOfNoShow struct {
	Rate uint8 `json:"Rate"`
}
type AdditionalCases struct {
	AdditionalCase AdditionalCase `json:"AdditionalCase"`
}
type AdditionalCase struct {
	StartDays string `json:"StartDays"`
	EndDays   string `json:"EndDays"`
	Rate      string `json:"Rate"`
}

// CreateInput 新規作成時の入力
type CreateInput struct {
	CancelPolicyName string `json:"CancelPolicyName" validate:"required"`
	PropertyID       int64  `json:"PropertyId" validate:"required"`
	CancelPolicyJSON
}

// ListInput 一覧取得時の入力
type ListInput struct {
	PropertyID int64 `param:"propertyId" validate:"required"`
}

// DetailInput 詳細情報の入力
type DetailInput struct {
	PropertyID         *int64  `query:"property_id"`
	PlanCancelPolicyID *uint64 `query:"plan_cancel_policy_id"`
}

// UpdateInput キャンセルポリシー更新の入力
type UpdateInput struct {
	PropertyID         *int64  `json:"PropertyId"`
	PlanCancelPolicyID *uint64 `json:"PlanCancelPolicyId"`
	PolicyName         string  `json:"CancelPolicyName"`
	CancelPolicyJSON
}

// DeleteInput 詳細情報の入力
type DeleteInput struct {
	PlanCancelPolicyID uint64 `param:"planCancelPolicyId" validate:"required"`
}

// PlanListInput プラン一覧取得時の入力
type PlanListInput struct {
	PlanCancelPolicyID uint64 `query:"plan_cancel_policy_id" validate:"required"`
}

type HtTmPlanCancelPolicies struct {
	PlanCancelPolicyID uint64 `gorm:"primary_key;autoIncrement:true"`
	PolicyName         string
	WholesalerID       int
	PropertyID         int64
	CancelPenaltyJSON  string
}

type HtThPlanCancelPolicyRelations struct {
	CancelPolicyPlanIdRelationId uint64 `gorm:"primary_key;autoIncrement:true"`
	WholesalerID                 int
	PropertyID                   int64
	PlanID                       int64
	PlanCancelPolicyID           uint64
}

// ICancelPolicyUsecase キャンセルポリシー関連のusecaseのインターフェース
type ICancelPolicyUsecase interface {
	Create(req *CreateInput) error
	List(req *ListInput) ([]CancelPolicyInfo, error)
	Detail(req *DetailInput) (*CancelPolicyJSONWithName, error)
	Save(req *UpdateInput) error
	Delete(req *DeleteInput) error
	PlanList(req *PlanListInput) ([]PlanInfo, error)
}

// ICancelPolicyCommonRepository 共通のキャンセルポリシー操作インターフェース
type ICancelPolicyCommonRepository interface {
	CreatePlanCancelPolicy(policyName string, wholesalerID int, propertyID int64, cancelPolicy string) error
	FindPlanCancelPolicy(policyId uint64) (CancelPolicyJSONWithName, error)
	PlanCancelPolicyList(propertyID int64) ([]HtTmPlanCancelPolicies, error)
	UpdatePlanCancelPolicy(policyId uint64, policyName string, cancelPolicy string) error
	DeletePlanCancelPolicy(policyId uint64) error
	FindAssignedPlanCancelPolicy(propertyId int64, planId int64) (*HtThPlanCancelPolicyRelations, error)
	UpsertPlanCancelPolicyRelation(wholesalerID int, propertyId int64, planId int64, policyId uint64) error
	DeletePlanCancelPolicyRelation(wholesalerID int, planId int64) error
}
