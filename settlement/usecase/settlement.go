package usecase

import (
	"os"
	"strings"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aInfra "github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"github.com/Adventureinc/hotel-hm-api/src/settlement"
	sInfra "github.com/Adventureinc/hotel-hm-api/src/settlement/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"gorm.io/gorm"
)

// settlementUsecase 請求関連usecase
type settlementUsecase struct {
	SRepository       settlement.ISettlementRepository
	SettlementStorage settlement.ISettlementStorage
	ARepository       account.IAccountRepository
}

// NewSettlementUsecase インスタンス生成
func NewSettlementUsecase(db *gorm.DB) settlement.ISettlementUsecase {
	return &settlementUsecase{
		SRepository:       sInfra.NewSettlementRepository(db),
		SettlementStorage: sInfra.NewSettlementStorage(),
		ARepository:       aInfra.NewAccountRepository(db),
	}
}

// 精算書一覧
func (s *settlementUsecase) FetchAll(req settlement.ListInput) ([]settlement.ListOutput, error) {
	res := []settlement.ListOutput{}
	hmUser, hErr := s.ARepository.FetchHMUserByPropertyID(req.PropertyID, req.WholesalerID)
	if hErr != nil {
		return res, hErr
	}
	settlements, sErr := s.SRepository.FetchAll(hmUser.HotelManagerID)
	if sErr != nil {
		return res, sErr
	}
	for _, settlementData := range *settlements {
		status := uint8(0) // 未承認
		now := time.Now()
		if now.After(settlementData.FixedDate) {
			status = 2 // 承認確定
		} else if settlementData.ApproveFlg {
			status = 1 // 承認済み
		}
		res = append(res, settlement.ListOutput{
			SettlementID: settlementData.ID,
			TargetDate:   settlementData.DateDivision.Format("2006-01-02"),
			IssueDate:    settlementData.DateOfIssue.Format("2006-01-02"),
			FixedDate:    settlementData.FixedDate.Format("2006-01-02"),
			IsApprove:    settlementData.ApproveFlg,
			Status:       status,
		})
	}
	return res, nil
}

// Approve 承認処理
func (s *settlementUsecase) Approve(req settlement.UpdateInput) error {
	return s.SRepository.UpdateApproveFlg(req.SettlementID, req.IsApprove)
}

// Download 精算書ダウンロード
func (s *settlementUsecase) Download(req *settlement.DownloadInput) (string, string, error) {
	settlementData, err := s.SRepository.FetchOne(req.SettlementID)
	if err != nil {
		return "", "", err
	}

	tempFileName, sErr := s.SettlementStorage.Get(os.Getenv("GCS_PROP_SETTLEMENT_BUCKET_NAME"), settlementData.SourcePath)
	if sErr != nil {
		return "", "", err
	}

	// パスの最後から、hotel_manager_id（XX-2020年...のXX-のところ）を除いた文字列をファイル名とする
	t := strings.Split(settlementData.SourcePath, "/")
	downloadFileName := strings.Split(t[len(t)-1], "-")
	return tempFileName, downloadFileName[len(downloadFileName)-1], nil
}

// FetchAccounts 精算情報取得
func (s *settlementUsecase) FetchInfo(req *settlement.InfoInput, claimParam *account.ClaimParam) (*settlement.InfoOutput, error) {
	response := &settlement.InfoOutput{}
	settlementAccount, sErr := s.SRepository.FetchAccount(req.PropertyID)
	if sErr != nil {
		return response, sErr
	}
	mails, mErr := s.SRepository.FetchMails(settlementAccount.HotelManagerID)
	if mErr != nil {
		return response, mErr
	}

	addressee, dErr := utils.Decrypt(settlementAccount.Addressee)
	if dErr != nil {
		return response, dErr
	}
	bankName, dErr := utils.Decrypt(settlementAccount.BankName)
	if dErr != nil {
		return response, dErr
	}
	bankNameRuby, dErr := utils.Decrypt(settlementAccount.BankNameRuby)
	if dErr != nil {
		return response, dErr
	}
	bankCode, dErr := utils.Decrypt(settlementAccount.BankCode)
	if dErr != nil {
		return response, dErr
	}
	bankBranch, dErr := utils.Decrypt(settlementAccount.BankBranch)
	if dErr != nil {
		return response, dErr
	}
	bankBranchRuby, dErr := utils.Decrypt(settlementAccount.BankBranchRuby)
	if dErr != nil {
		return response, dErr
	}
	bankBranchCode, dErr := utils.Decrypt(settlementAccount.BankBranchCode)
	if dErr != nil {
		return response, dErr
	}
	bankAccountType, dErr := utils.Decrypt(settlementAccount.BankAccountType)
	if dErr != nil {
		return response, dErr
	}
	bankAccountNumber, dErr := utils.Decrypt(settlementAccount.BankAccountNumber)
	if dErr != nil {
		return response, dErr
	}
	bankAccountHolder, dErr := utils.Decrypt(settlementAccount.BankAccountHolder)
	if dErr != nil {
		return response, dErr
	}

	var emails []string
	for _, v := range *mails {
		emailDec, dErr := utils.Decrypt(v.EmailEnc)
		if dErr != nil {
			return response, dErr
		}
		emails = append(emails, emailDec)
	}

	response.AccountID = settlementAccount.ID
	response.Addressee = addressee
	response.BankName = bankName
	response.BankNameRuby = bankNameRuby
	response.BankCode = bankCode
	response.BankBranch = bankBranch
	response.BankBranchRuby = bankBranchRuby
	response.BankBranchCode = bankBranchCode
	response.BankAccountType = bankAccountType
	response.BankAccountNumber = bankAccountNumber
	response.BankAccountHolder = bankAccountHolder
	response.Emails = emails
	return response, nil
}

// SaveInfo 精算情報を作成更新
func (s *settlementUsecase) SaveInfo(req *settlement.SaveInfoInput, claimParam *account.ClaimParam) error {
	hmUser, hErr := s.ARepository.FetchHMUserByPropertyID(req.PropertyID, req.WholesalerID)
	if hErr != nil {
		return hErr
	}
	fetchedMails, mErr := s.SRepository.FetchMails(hmUser.HotelManagerID)
	if mErr != nil {
		return mErr
	}

	// 各データの暗号化
	addresseeEnc, eErr := utils.Encrypt(req.Addressee)
	if eErr != nil {
		return eErr
	}
	bankNameEnc, eErr := utils.Encrypt(req.BankName)
	if eErr != nil {
		return eErr
	}
	bankNameRubyEnc, eErr := utils.Encrypt(req.BankNameRuby)
	if eErr != nil {
		return eErr
	}
	bankCodeEnc, eErr := utils.Encrypt(req.BankCode)
	if eErr != nil {
		return eErr
	}
	bankBranchEnc, eErr := utils.Encrypt(req.BankBranch)
	if eErr != nil {
		return eErr
	}
	bankBranchRubyEnc, eErr := utils.Encrypt(req.BankBranchRuby)
	if eErr != nil {
		return eErr
	}
	bankBranchCodeEnc, eErr := utils.Encrypt(req.BankBranchCode)
	if eErr != nil {
		return eErr
	}
	bankAccountTypeEnc, eErr := utils.Encrypt(req.BankAccountType)
	if eErr != nil {
		return eErr
	}
	bankAccountNumberEnc, eErr := utils.Encrypt(req.BankAccountNumber)
	if eErr != nil {
		return eErr
	}
	bankAccountHolderEnc, eErr := utils.Encrypt(req.BankAccountHolder)
	if eErr != nil {
		return eErr
	}

	var emailsEnc []string
	for _, v := range req.Emails {
		emailEnc, eErr := utils.Encrypt(v)
		if eErr != nil {
			return eErr
		}
		emailsEnc = append(emailsEnc, emailEnc)
	}

	upsertData := &settlement.HtTmSettlementAccounts{
		ID:                req.AccountID,
		PropertyID:        req.PropertyID,
		HotelManagerID:    hmUser.HotelManagerID,
		Addressee:         addresseeEnc,
		BankName:          bankNameEnc,
		BankNameRuby:      bankNameRubyEnc,
		BankCode:          bankCodeEnc,
		BankBranch:        bankBranchEnc,
		BankBranchRuby:    bankBranchRubyEnc,
		BankBranchCode:    bankBranchCodeEnc,
		BankAccountType:   bankAccountTypeEnc,
		BankAccountNumber: bankAccountNumberEnc,
		BankAccountHolder: bankAccountHolderEnc,
		IsUserUpdate:      true,
	}

	tx, txErr := s.SRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	txSettlementRepo := sInfra.NewSettlementRepository(tx)
	if err := txSettlementRepo.UpsertAccount(upsertData); err != nil {
		s.SRepository.TxRollback(tx)
		return err
	}
	insertData := []settlement.HtThHotelManagerSettlementNotifications{}
	existMailData := []string{}
	for _, mail := range emailsEnc {
		duplicate := false
		existMailData = append(existMailData, mail)
		// すでに登録しているメールアドレスと一致している場合、除外する
		for _, fetchedMail := range *fetchedMails {
			if fetchedMail.EmailEnc == "" {
				continue
			}
			if fetchedMail.EmailEnc == mail {
				duplicate = true
			}
		}
		if duplicate == true {
			duplicate = false
			continue
		}
		// insertするデータ
		insertData = append(insertData, settlement.HtThHotelManagerSettlementNotifications{
			HotelManagerID: hmUser.HotelManagerID,
			EmailEnc:       mail,
		})
	}
	if err := txSettlementRepo.ClearNotificationMails(hmUser.HotelManagerID, existMailData); err != nil {
		s.SRepository.TxRollback(tx)
		return err
	}

	if len(insertData) > 0 {
		if err := txSettlementRepo.CreateNotificationMails(&insertData); err != nil {
			s.SRepository.TxRollback(tx)
			return err
		}
	}

	// コミットとロールバック
	if err := s.SRepository.TxCommit(tx); err != nil {
		s.SRepository.TxRollback(tx)
		return err
	}
	return nil
}
