package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/settlement"
	"gorm.io/gorm"
)

// settlementRepository 請求関連repository
type settlementRepository struct {
	db *gorm.DB
}

// NewSettlementRepository インスタンス生成
func NewSettlementRepository(db *gorm.DB) settlement.ISettlementRepository {
	return &settlementRepository{
		db: db,
	}
}

// TxStart トランザクションスタート
func (s *settlementRepository) TxStart() (*gorm.DB, error) {
	tx := s.db.Debug().Begin()
	return tx, tx.Error
}

// TxCommit トランザクションコミット
func (s *settlementRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// TxRollback トランザクション ロールバック
func (s *settlementRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

// FetchAll 請求書一覧を複数件取得
func (s *settlementRepository) FetchAll(hotelManagerID int64) (*[]settlement.HtTmHotelSettlements, error) {
	res := &[]settlement.HtTmHotelSettlements{}
	err := s.db.
		Model(&settlement.HtTmHotelSettlements{}).
		Where("hotel_manager_id = ?", hotelManagerID).
		Where("expired_at IS NULL").
		Group("date_division").
		Order("date_division DESC, created_at DESC").
		Find(res).Error
	return res, err
}

// UpdateApproveFlg 請求書の承認フラグを更新
func (s *settlementRepository) UpdateApproveFlg(settlementID int64, approveFlg bool) error {
	return s.db.
		Model(&settlement.HtTmHotelSettlements{}).
		Where("id = ?", settlementID).
		Update("approve_flg", approveFlg).
		Error
}

// FetchOne 精算情報を1件取得
func (s *settlementRepository) FetchOne(settlementID int64) (*settlement.HtTmHotelSettlements, error) {
	res := &settlement.HtTmHotelSettlements{}
	err := s.db.
		Model(&settlement.HtTmHotelSettlements{}).
		Where("id = ?", settlementID).
		First(res).Error
	return res, err
}

// FetchAccount 口座情報を１件取得
func (s *settlementRepository) FetchAccount(propertyID int64) (*settlement.HtTmSettlementAccounts, error) {
	res := &settlement.HtTmSettlementAccounts{}
	err := s.db.
		Select("ht_tm_settlement_accounts.id",
			"ht_tm_settlement_accounts.hotel_manager_id",
			"ht_tm_settlement_accounts.property_id",
			"ht_tm_settlement_accounts.addressee",
			"ht_tm_settlement_accounts.bank_name",
			"ht_tm_settlement_accounts.bank_name_ruby",
			"ht_tm_settlement_accounts.bank_code",
			"ht_tm_settlement_accounts.bank_branch",
			"ht_tm_settlement_accounts.bank_branch_ruby",
			"ht_tm_settlement_accounts.bank_branch_code",
			"ht_tm_settlement_accounts.bank_account_type",
			"ht_tm_settlement_accounts.bank_account_number",
			"ht_tm_settlement_accounts.bank_account_holder",
			"ht_tm_settlement_accounts.closing_date",
			"ht_tm_settlement_accounts.payment_date_month_later",
			"ht_tm_settlement_accounts.memo",
			"ht_tm_settlement_accounts.is_user_update",
			"ht_tm_settlement_accounts.created_at",
			"ht_tm_settlement_accounts.updated_at").
		Table("ht_tm_settlement_accounts").
		Joins("INNER JOIN ht_tm_hotel_managers ON ht_tm_settlement_accounts.hotel_manager_id = ht_tm_hotel_managers.hotel_manager_id").
		Where("ht_tm_settlement_accounts.property_id = ?", propertyID).
		Where("ht_tm_settlement_accounts.hotel_manager_id <> 0"). // 現状、hotel_manager_idが0の情報が登録されている為、暫定的に0以外を条件文に追加
		Where("ht_tm_hotel_managers.del_flg = 0").
		First(res).Error
	return res, err
}

// UpsertAccount 口座情報を更新・作成
func (s *settlementRepository) UpsertAccount(upsertData *settlement.HtTmSettlementAccounts) error {
	assignData := map[string]interface{}{
		"hotel_manager_id":    upsertData.HotelManagerID,
		"property_id":         upsertData.PropertyID,
		"addressee":           upsertData.Addressee,
		"bank_name":           upsertData.BankName,
		"bank_name_ruby":      upsertData.BankNameRuby,
		"bank_code":           upsertData.BankCode,
		"bank_branch":         upsertData.BankBranch,
		"bank_branch_ruby":    upsertData.BankBranchRuby,
		"bank_branch_code":    upsertData.BankBranchCode,
		"bank_account_type":   upsertData.BankAccountType,
		"bank_account_number": upsertData.BankAccountNumber,
		"bank_account_holder": upsertData.BankAccountHolder,
		"is_user_update":      upsertData.IsUserUpdate,
	}
	return s.db.Model(&settlement.HtTmSettlementAccounts{}).
		Where("id = ?", upsertData.ID).
		Assign(assignData).
		FirstOrCreate(&settlement.HtTmSettlementAccounts{}).
		Error
}

// FetchMails 通知先を複数件取得
func (s *settlementRepository) FetchMails(HotelManagerID int64) (*[]settlement.HtThHotelManagerSettlementNotifications, error) {
	res := &[]settlement.HtThHotelManagerSettlementNotifications{}
	err := s.db.Debug().
		Model(&settlement.HtThHotelManagerSettlementNotifications{}).
		Where("hotel_manager_id = ?", HotelManagerID).
		Find(res).Error
	return res, err
}

// CreateNotificationMails 通知先を複数件作成
func (s *settlementRepository) CreateNotificationMails(insertData *[]settlement.HtThHotelManagerSettlementNotifications) error {
	return s.db.Model(&settlement.HtThHotelManagerSettlementNotifications{}).Create(insertData).Error
}

// ClearNotificationMails 通知先を複数件削除
func (s *settlementRepository) ClearNotificationMails(HotelManagerID int64, existData []string) error {
	query := s.db.Where("hotel_manager_id = ?", HotelManagerID)
	if len(existData) > 0 {
		query = query.Where("email_enc NOT IN ?", existData)
	}
	return query.Delete(&settlement.HtThHotelManagerSettlementNotifications{}).Error
}
