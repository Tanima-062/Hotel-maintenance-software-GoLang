package settlement

import (
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common"
)

const (
	settlement string = "settlement"
	details    string = "details"
)

// HtTmHotelSettlements 請求書テーブル
type HtTmHotelSettlements struct {
	ID             int64     `json:"id"`
	HotelManagerID int64     `json:"hotel_mamager_id"`
	FileType       string    `json:"file_type"`
	SourcePath     string    `json:"source_type"`
	DateDivision   time.Time `gorm:"type:time" json:"date_division"`
	DateOfIssue    time.Time `gorm:"type:time" json:"date_of_issue"`
	FixedDate      time.Time `gorm:"type:time" json:"fixed_date"`
	ApproveFlg     bool      `json:"approve_flg"`
	Generation     int64     `json:"generation"`
	OperatorID     int64     `json:"operator_id"`
	ExpiredAt      time.Time `gorm:"type:time" json:"expired_at"`
	common.Times   `gorm:"embedded"`
}

// HtTmSettlementAccounts 精算情報テーブル
type HtTmSettlementAccounts struct {
	ID                    int64  `json:"id"`
	HotelManagerID        int64  `json:"hotel_mamager_id"`
	PropertyID            int64  `json:"property_id"`
	Addressee             string `json:"addressee"`
	BankName              string `json:"bank_name"`
	BankNameRuby          string `json:"bank_name_ruby"`
	BankCode              string `json:"bank_code"`
	BankBranch            string `json:"bank_branch"`
	BankBranchRuby        string `json:"bank_branch_ruby"`
	BankBranchCode        string `json:"bank_branch_code"`
	BankAccountType       string `json:"bank_account_type"`
	BankAccountNumber     string `json:"bank_account_number"`
	BankAccountHolder     string `json:"bank_account_holder"`
	ClosingDate           int16  `json:"closing_date" gorm:"default:99"`
	PaymentDateMonthLater int16  `json:"payment_date_month_later"`
	Memo                  string `json:"memo"`
	IsUserUpdate          bool   `json:"is_user_update"`
	common.Times          `gorm:"embedded"`
}

// HtThHotelManagerSettlementNotifications 通知先テーブル
type HtThHotelManagerSettlementNotifications struct {
	ID                   int64  `json:"id"`
	HotelManagerID       int64  `json:"hotel_mamager_id"`
	EmailEnc             string `json:"email_enc"`
	NotifyDate           string `json:"notify_date"`
	NotificationSchedule string `json:"notification_schedule"`
	Memo                 string `json:"memo"`
	common.Times         `gorm:"embedded"`
}

// ListInput 精算書一覧の入力
type ListInput struct {
	PropertyID   int64 `json:"property_id" param:"propertyId" valiadte:"required"`
	WholesalerID int64 `json:"wholesaler_id" query:"wholesaler_id" validate:"required"`
}

// ListOutput 精算書一覧の出力
type ListOutput struct {
	SettlementID int64  `json:"settlement_id"`
	TargetDate   string `json:"target_date"`
	IssueDate    string `json:"issue_date"`
	FixedDate    string `json:"fixed_date"`
	Status       uint8  `json:"status"`
	IsApprove    bool   `json:"is_approve"`
}

// UpdateInput 承認更新の入力
type UpdateInput struct {
	SettlementID int64 `json:"settlement_id" valiadte:"required"`
	IsApprove    bool  `json:"is_approve"`
}

// DownloadInput ダウンロードの入力
type DownloadInput struct {
	SettlementID int64 `json:"settlement_id" param:"settlementId" valiadte:"required"`
}

// InfoInput 精算情報取得の入力
type InfoInput struct {
	PropertyID int64 `json:"property_id" param:"propertyId" valiadte:"required"`
}

// InfoOutput 精算情報取得の出力
type InfoOutput struct {
	AccountID         int64    `json:"account_id"`
	Addressee         string   `json:"addressee"`
	BankName          string   `json:"bank_name"`
	BankNameRuby      string   `json:"bank_name_ruby"`
	BankCode          string   `json:"bank_code"`
	BankBranch        string   `json:"bank_branch"`
	BankBranchRuby    string   `json:"bank_branch_ruby"`
	BankBranchCode    string   `json:"bank_branch_code"`
	BankAccountType   string   `json:"bank_account_type"`
	BankAccountNumber string   `json:"bank_account_number"`
	BankAccountHolder string   `json:"bank_account_holder"`
	Emails            []string `json:"emails"`
}

// SaveInfoInput 精算情報の作成・更新の入力
type SaveInfoInput struct {
	AccountID         int64    `json:"account_id"`
	PropertyID        int64    `json:"property_id" validate:"required"`
	WholesalerID      int64    `json:"wholesaler_id" validate:"required"`
	Addressee         string   `json:"addressee" validate:"required"`
	BankName          string   `json:"bank_name" validate:"required"`
	BankNameRuby      string   `json:"bank_name_ruby" validate:"required"`
	BankCode          string   `json:"bank_code" validate:"required"`
	BankBranch        string   `json:"bank_branch" validate:"required"`
	BankBranchRuby    string   `json:"bank_branch_ruby" validate:"required"`
	BankBranchCode    string   `json:"bank_branch_code" validate:"required"`
	BankAccountType   string   `json:"bank_account_type" validate:"required"`
	BankAccountNumber string   `json:"bank_account_number" validate:"required"`
	BankAccountHolder string   `json:"bank_account_holder" validate:"required"`
	Emails            []string `json:"emails"`
}

// ISettlementUsecase 請求関連のusecaseのインターフェース
type ISettlementUsecase interface {
	FetchAll(req ListInput) ([]ListOutput, error)
	Approve(req UpdateInput) error
	Download(req *DownloadInput) (string, string, error)
	FetchInfo(req *InfoInput, claimParam *account.ClaimParam) (*InfoOutput, error)
	SaveInfo(req *SaveInfoInput, claimParam *account.ClaimParam) error
}

// ISettlementRepository 請求関連のrepositoryのインターフェース
type ISettlementRepository interface {
	common.Repository
	// FetchAll 請求書一覧を複数件取得
	FetchAll(hotelManagerID int64) (*[]HtTmHotelSettlements, error)
	// UpdateApproveFlg 請求書の承認フラグを更新
	UpdateApproveFlg(settlementID int64, approveFlg bool) error
	// FetchOne 精算情報を1件取得
	FetchOne(settlementID int64) (*HtTmHotelSettlements, error)
	// FetchAccount 口座情報を１件取得
	FetchAccount(propertyID int64) (*HtTmSettlementAccounts, error)
	// UpsertAccount 口座情報を更新・作成
	UpsertAccount(upsertData *HtTmSettlementAccounts) error
	// FetchMails 通知先を複数件取得
	FetchMails(HotelManagerID int64) (*[]HtThHotelManagerSettlementNotifications, error)
	// CreateNotificationMails 通知先を複数件作成
	CreateNotificationMails(insertData *[]HtThHotelManagerSettlementNotifications) error
	// ClearNotificationMails 通知先を複数件削除
	ClearNotificationMails(HotelManagerID int64, existData []string) error
}

// ISettlementStorage 請求関連のrstorageのインターフェース
type ISettlementStorage interface {
	Get(bucketName string, objectPath string) (string, error)
}
