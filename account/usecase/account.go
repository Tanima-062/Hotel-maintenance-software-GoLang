package usecase

import (
	"fmt"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"gorm.io/gorm"
)

type accountUsecase struct {
	ARepository account.IAccountRepository
}

// NewAccountUsecase インスタンス生成
func NewAccountUsecase(db *gorm.DB) account.IAccountUsecase {
	return &accountUsecase{
		ARepository: infra.NewAccountRepository(db),
	}
}

// Login ログイン時のID＆PWチェックとトークン発行
func (a *accountUsecase) Login(LoginInput *account.LoginInput) (string, error) {
	usernameEnc, eErr := utils.Encrypt(LoginInput.Username)
	if eErr != nil {
		return "", eErr
	}
	passwordEnc, eErr := utils.Encrypt(LoginInput.Password)
	if eErr != nil {
		return "", eErr
	}

	hmUser := a.ARepository.FetchHMUserByLoginInfo(&account.HtTmHotelManager{
		UsernameEnc: usernameEnc,
		PasswordEnc: passwordEnc,
	})
	if hmUser.HotelManagerID == 0 {
		return "", fmt.Errorf("Error: %s", "ユーザーIDもしくはパスワードが正しくありません。")
	}

	token, tokenErr := utils.GenerateToken(hmUser.HotelManagerID)
	if tokenErr != nil {
		return "", tokenErr
	}

	claimParam := &account.ClaimParam{
		HotelManagerID: hmUser.HotelManagerID,
		APIToken:       token,
	}
	hmUser.LoginedAt = time.Now()
	if DBErr := a.ARepository.SaveLoginInfo(&hmUser, claimParam, token); DBErr != nil {
		return "", DBErr
	}
	return token, nil
}

// Logout ログアウト処理。トークンを削除するだけ
func (a *accountUsecase) Logout(claimParam *account.ClaimParam) error {
	return a.ARepository.DeleteAPIToken(claimParam)
}

// CheckToken 正しいトークンか確認し、トークン再度生成する
func (a *accountUsecase) CheckToken(claimParam *account.ClaimParam) (string, error) {
	fetchedHmUser, fetchErr := a.ARepository.FetchHMUserByToken(claimParam)
	if fetchErr != nil {
		return "", fetchErr
	}

	newToken, tokenErr := utils.GenerateToken(fetchedHmUser.HotelManagerID)
	if tokenErr != nil {
		return "", tokenErr
	}

	if DBErr := a.ARepository.SaveLoginInfo(&fetchedHmUser, claimParam, newToken); DBErr != nil {
		return "", DBErr
	}
	return newToken, nil
}

// FetchDetail HMアカウント情報の取得
func (a *accountUsecase) FetchDetail(claimParam *account.ClaimParam) (*account.HtTmHotelManager, error) {
	fetchedHmUser, fetchErr := a.ARepository.FetchHMUserByToken(claimParam)
	if fetchErr != nil {
		return &account.HtTmHotelManager{}, fetchErr
	}
	firstName, dErr := utils.Decrypt(fetchedHmUser.FirstNameEnc)
	if dErr != nil {
		return &account.HtTmHotelManager{}, dErr
	}
	lastName, dErr := utils.Decrypt(fetchedHmUser.LastNameEnc)
	if dErr != nil {
		return &account.HtTmHotelManager{}, dErr
	}
	email, dErr := utils.Decrypt(fetchedHmUser.EmailEnc)
	if dErr != nil {
		return &account.HtTmHotelManager{}, dErr
	}
	username, dErr := utils.Decrypt(fetchedHmUser.UsernameEnc)
	if dErr != nil {
		return &account.HtTmHotelManager{}, dErr
	}
	fetchedHmUser.FirstNameEnc = firstName
	fetchedHmUser.LastNameEnc = lastName
	fetchedHmUser.EmailEnc = email
	fetchedHmUser.UsernameEnc = username
	fetchedHmUser.PasswordEnc = ""
	return &fetchedHmUser, nil
}

// ChangePassword パスワード変更
func (a *accountUsecase) ChangePassword(request *account.ChangePasswordInput) error {
	usernameEnc, eErr := utils.Encrypt(request.Username)
	if eErr != nil {
		return eErr
	}
	passwordEnc, eErr := utils.Encrypt(request.Password)
	if eErr != nil {
		return eErr
	}
	hmUser := a.ARepository.FetchHMUserByLoginInfo(&account.HtTmHotelManager{
		UsernameEnc: usernameEnc,
		PasswordEnc: passwordEnc,
	})
	if hmUser.HotelManagerID == 0 {
		return fmt.Errorf("Error: %s", "ユーザーIDもしくはパスワードが正しくありません。")
	}

	newPasswordEnc, eErr := utils.Encrypt(request.NewPassword)
	if eErr != nil {
		return eErr
	}
	return a.ARepository.UpdatePassword(hmUser.HotelManagerID, newPasswordEnc)
}

// FetchHMUser トークンに紐づくHMアカウントを取得
func (a *accountUsecase) FetchHMUserByToken(claimParam *account.ClaimParam) (account.HtTmHotelManager, error) {
	return a.ARepository.FetchHMUserByToken(claimParam)
}

// isParentAccount 親アカウントかどうか
func (a *accountUsecase) IsParentAccount(hotelManagerID int64) bool {
	hmUser, err := a.ARepository.FetchOne(hotelManagerID)
	if err != nil {
		return false
	}
	return hmUser.PropertyID == 0
}
