package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common/infra"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
	"gorm.io/gorm"
)

const (
	// http://www.unicodemap.org/range/62/Hiragana/
	hiraganaLo = 0x3041 // ぁ

	// http://www.unicodemap.org/range/63/Katakana/
	katakanaLo = 0x30a1 // ァ

	codeDiff = katakanaLo - hiraganaLo
)

func RequestLog(c echo.Context, request interface{}) {
	request_log, _ := json.Marshal(request)
	c.Echo().Logger.Info("request >>> " + string(request_log))
}

// GenerateToken token発行スクリプト
func GenerateToken(hotelManagerID int64) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = os.Getenv("JWT_CLAIMS_NAME")
	claims["str"] = os.Getenv("JWT_CLAIMS_STR")
	claims["hotelManagerID"] = hotelManagerID
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

	// Generate encoded token and send it as response
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// GetExtensionFromContentType ContentTypeから拡張子を取得
// 必要に応じて追記してください
func GetExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpeg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}

// RemoveDoubleQuotation ダブルクオーテーション削除
func RemoveDoubleQuotation(str string) string {
	return strings.Replace(str, "\"", "", -1)
}

// ConvertNewlineCodeToBrTag 改行コードをbrタグに変換する処理
func ConvertNewlineCodeToBrTag(str string) string {
	return strings.NewReplacer(
		"\r\n", "<br>",
		"\r", "<br>",
		"\n", "<br>",
	).Replace(str)
}

// ConvertBrTagToNewlineCode brタグを改行コードに変換する処理
func ConvertBrTagToNewlineCode(str string) string {
	return strings.Replace(str, "<br>", "\n", -1)
}

// GetHmUser トークンからhotelmanagerIDとAPITokenを抜き出す処理
func GetHmUser(c echo.Context) (*account.ClaimParam, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	if claims["name"].(string) != os.Getenv("JWT_CLAIMS_NAME") || claims["str"].(string) != os.Getenv("JWT_CLAIMS_STR") {
		return &account.ClaimParam{}, fmt.Errorf("Error: %s", "tokenが不正です")
	}
	return &account.ClaimParam{HotelManagerID: int64(claims["hotelManagerID"].(float64)), APIToken: user.Raw}, nil
}

// PublicHoliday 休日用の構造体
type PublicHoliday struct {
	Name string
	Date time.Time
}

// GetHoliday 祝日を取得する処理
func GetHoliday() ([]PublicHoliday, error) {
	var holiday []PublicHoliday
	commonDB, DBErr := infra.CommonDBCon()
	if DBErr != nil {
		return holiday, DBErr
	}
	pool, cErr := commonDB.DB()
	if cErr != nil {
		pool.Close()
		return holiday, cErr
	}
	err := commonDB.Table("cm_tm_public_holiday").Where("date >= ?", time.Now().Format("2006-01-02")).Find(&holiday).Error
	pool.Close()
	return holiday, err
}

// GetBookingStatus 予約ステータスを返却する処理
func GetBookingStatus(cancelFlg bool, noShowFlg bool, arrival string, depature string) uint8 {
	if cancelFlg == true && noShowFlg == false {
		return ReserveStatusCancel
	}
	if cancelFlg == true && noShowFlg == true {
		return ReserveStatusNoShow
	}
	now := time.Now()
	arrivalTime, _ := time.Parse(time.RFC3339, arrival)
	departureTime, _ := time.Parse(time.RFC3339, depature)
	departureTime = departureTime.Add(24 * time.Hour)
	if now.Before(arrivalTime) && now.Before(departureTime) {
		return ReserveStatusReserved
	}
	if now.After(arrivalTime) && now.Before(departureTime) {
		return ReserveStatusStaying
	}
	if now.After(departureTime) {
		return ReserveStatusStayed
	}
	return 0
}

// BankInfo はcm_tm_bank_1, cm_tm_bank_2のスキーマ
type BankInfo struct {
	BankId       int32
	BranchBankId uint32
	BankNameKana string
	BankName     string
	DivisionId   uint8
	LineId       uint8
}

// ActiveBank はcm_tt_active_bankのスキーマ
type ActiveBank struct {
	ActiveBankId uint32
}

// GetAvailableBanksByKeyword キーワードに部分一致する銀行名を取得する
func GetAvailableBanksByKeyword(keyword string, limit int) ([]BankInfo, error) {
	var bankInfo []BankInfo
	commonDB, DBErr := infra.CommonDBCon()
	if DBErr != nil {
		return bankInfo, DBErr
	}
	db, cErr := commonDB.DB()
	if cErr != nil {
		return bankInfo, cErr
	}
	defer db.Close()

	tblNo, bIdErr := getActiveBankTableId(commonDB)
	if bIdErr != nil {
		return bankInfo, bIdErr
	}

	// 全角キーワード
	widenKw := fmt.Sprintf("%%%s%%", width.Widen.String(keyword))
	// 半角キーワード
	narrowKw := fmt.Sprintf("%%%s%%", width.Narrow.String(norm.NFD.String(HiraganaToKatakana(keyword))))
	err := commonDB.
		Table(fmt.Sprintf("cm_tm_bank_%d", tblNo)).
		Where("(bank_name like ? or bank_name_kana like ?) and branch_bank_id=0", widenKw, narrowKw).
		Order("bank_id").
		Limit(limit).
		Find(&bankInfo).Error

	if err != nil {
		return nil, err
	}

	return bankInfo, err
}

// GetBranchBanksByBankId キーワードに部分一致する銀行名を取得する
func GetBranchBanksByKeyword(bankId string, keyword string, limit int) ([]BankInfo, error) {
	// TBD 共通化する //
	var bankInfo []BankInfo
	commonDB, DBErr := infra.CommonDBCon()
	if DBErr != nil {
		return bankInfo, DBErr
	}
	db, cErr := commonDB.DB()
	if cErr != nil {
		return bankInfo, cErr
	}
	defer db.Close()

	tblNo, bIdErr := getActiveBankTableId(commonDB)
	if bIdErr != nil {
		return bankInfo, bIdErr
	}
	// 全角キーワード
	widenKw := fmt.Sprintf("%%%s%%", width.Widen.String(keyword))
	// 半角キーワード
	narrowKw := fmt.Sprintf("%%%s%%", width.Narrow.String(norm.NFD.String(HiraganaToKatakana(keyword))))

	// TBD 共通化する //
	err := commonDB.
		Table(fmt.Sprintf("cm_tm_bank_%d", tblNo)).
		Where("bank_id = ? and (bank_name like ? or bank_name_kana like ?)", bankId, widenKw, narrowKw).
		Limit(limit).
		Find(&bankInfo).Error

	if err != nil {
		return nil, err
	}

	return bankInfo, err
}

func getActiveBankTableId(commonDB *gorm.DB) (uint32, error) {
	var b []ActiveBank
	err := commonDB.Table("cm_tt_active_bank").Select([]string{"active_bank_id"}).Where("active_flg=?", 1).Find(&b).Error
	if err != nil || len(b) == 0 {
		return 0, nil
	}
	return b[0].ActiveBankId, err
}

// HiraganaToKatakana はひらがなをカタカナに変換する。
func HiraganaToKatakana(str string) string {
	src := []rune(str)
	dst := make([]rune, len(src))
	for i, r := range src {
		switch {
		case unicode.In(r, unicode.Hiragana):
			dst[i] = r + codeDiff
		default:
			dst[i] = r
		}
	}
	return string(dst)
}
