package usecase

import (
	"encoding/json"
	"os"
	"time"
	"strings"
	"unicode"
	"errors"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aInfra "github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	nInfra "github.com/Adventureinc/hotel-hm-api/src/facility/infra"
	"gorm.io/gorm"
)

// facilityTemaUsecase てまの施設関連usecase
type facilityTemaUsecase struct {
	ARepository     account.IAccountRepository
	ATemaRepository account.IAccountTemaRepository
	FRepository     facility.IFacilityRepository
	FTemaRepository facility.IFacilityTemaRepository
	FTlRepository   facility.IFacilityTlRepository
}

// NewFacilityTemaUsecase インスタンス生成
func NewFacilityTemaUsecase(db *gorm.DB) facility.IFacilityTemaUsecase {
	return &facilityTemaUsecase{
		ARepository:     aInfra.NewAccountRepository(db),
		ATemaRepository: aInfra.NewAccountTemaRepository(db),
		FRepository:     nInfra.NewFacilityRepository(db),
		FTemaRepository: nInfra.NewFacilityTemaRepository(db),
		FTlRepository:   nInfra.NewFacilityTlRepository(db),
	}
}

// FetchAll アカウントに紐づく施設情報をすべて取得
func (f *facilityTemaUsecase) FetchAll(hmUser account.HtTmHotelManager) ([]facility.InitFacilityOutput, error) {
	// 子施設の場合は当該施設の情報のみ返却する。
	if hmUser.PropertyID != facility.ParentPropertyId {
		return f.FTemaRepository.FetchAllFacilities([]int64{hmUser.PropertyID}, hmUser.WholesalerID)
	}

	// 親施設の場合のみは親施設と紐づく子施設の情報をを返却する。
	properties, _ := f.FRepository.FetchAllClientCompanies(hmUser.HotelManagerID)
	propertyIDs := []int64{}
	if len(properties) == 0 {
		propertyIDs = append(propertyIDs, hmUser.PropertyID)
	} else {
		for _, property := range properties {
			propertyIDs = append(propertyIDs, property.PropertyID)
		}
	}
	return f.FTemaRepository.FetchAllFacilities(propertyIDs, hmUser.WholesalerID)
}

// UpdateDispPriority サイト公開フラグを更新
func (f *facilityTemaUsecase) UpdateDispPriority(request *facility.UpdateDispPriorityInput) error {
	return f.FTemaRepository.UpdateDispPriority(request.PropertyID, request.DispPriority)
}

// FetchBaseInfo 施設基本情報を取得
func (f *facilityTemaUsecase) FetchBaseInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, request *facility.BaseInfoInput) (*facility.BaseInfoOutput, error) {
	property, pErr := f.FRepository.FetchProperty(request.PropertyID)
	if pErr != nil {
		return &facility.BaseInfoOutput{}, pErr
	}
	var categoryName string
	category, err := f.FRepository.FetchCategory(request.PropertyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == false {
			return &facility.BaseInfoOutput{}, err
		}
		categoryName = ""
	} else {
		categoryName = category.Name
	}

	connectUser, cErr := f.ATemaRepository.FetchConnectUser(request.PropertyID)
	var password string
	// 連携アカウントがない場合はエラーにしないで処理せずそのまま通す
	if cErr == nil && connectUser.Username != "" && connectUser.PasswordEnc != "" {
		passwordDec, dErr := utils.Decrypt(connectUser.PasswordEnc)
		if dErr != nil {
			return &facility.BaseInfoOutput{}, dErr
		}
		password = passwordDec
	} else {
		password = connectUser.PasswordEnc
	}

	return &facility.BaseInfoOutput{
		PropertyID:        property.PropertyID,
		Name:              property.Name,
		PostalCode:        property.PostalCode,
		StateProvinceName: property.StateProvinceName,
		City:              property.City,
		Line1:             property.Line1,
		Line2:             property.Line2,
		Line3:             property.Line3,
		Phone:             property.Phone,
		Fax:               property.Fax,
		CategoryName:      categoryName,
		ConnectID:         connectUser.Username,
		ConnectPassword:   password,
	}, nil
}

// FetchDetail 施設詳細情報を取得
func (f *facilityTemaUsecase) FetchDetail(request *facility.BaseInfoInput) (*facility.DetailOutput, error) {
	response := &facility.DetailOutput{}
	property, pErr := f.FTemaRepository.FetchPropertyDetail(request.PropertyID)
	if pErr != nil {
		return response, pErr
	}
	amenities, _ := f.FTemaRepository.FetchAmenities(request.PropertyID)

	response.PropertyID = request.PropertyID
	response.CheckinBegin = f.convertUpperTimeToLower(property.CheckinBegin)
	response.CheckinEnd = f.convertUpperTimeToLower(property.CheckinEnd)
	response.Checkout = f.convertUpperTimeToLower(property.Checkout)
	response.Instructions = utils.ConvertBrTagToNewlineCode(property.Instructions)
	response.SpecialInstructions = utils.ConvertBrTagToNewlineCode(property.SpecialInstructions)
	response.PolicyKnowBeforeYouGo = utils.ConvertBrTagToNewlineCode(property.PolicyKnowBeforeYouGo)
	response.FeeMandatory = utils.ConvertBrTagToNewlineCode(property.FeeMandatory)
	response.FeeOptional = utils.ConvertBrTagToNewlineCode(property.FeeOptional)
	response.DescriptionAmenity = utils.ConvertBrTagToNewlineCode(property.DescriptionAmenity)
	response.DescriptionAttractions = utils.ConvertBrTagToNewlineCode(property.DescriptionAttractions)
	response.DescriptionBusinessAmenities = utils.ConvertBrTagToNewlineCode(property.DescriptionBusinessAmenities)
	response.DescriptionDining = utils.ConvertBrTagToNewlineCode(property.DescriptionDining)
	response.DescriptionLocation = utils.ConvertBrTagToNewlineCode(property.DescriptionLocation)
	response.DescriptionHeadline = utils.ConvertBrTagToNewlineCode(property.DescriptionHeadline)
	response.DescriptionRooms = utils.ConvertBrTagToNewlineCode(property.DescriptionRooms)
	for _, amenityData := range *amenities {
		response.Amenities = append(response.Amenities, facility.Amenity{
			PropertyAmenityID: amenityData.PropertyAmenityID,
			AmenityName:       amenityData.AmenityName,
		})
	}
	return response, nil
}

// IsRegisteredConnect 連動IDが既に使われているかをチェック
func (f *facilityTemaUsecase) IsRegisteredConnect(request *facility.SaveBaseInfoInput) (bool, error) {
	isRegistered := false
	count, err := f.ATemaRepository.FetchCountOtherConnectedID(request.PropertyID, request.ConnectID)

	if err != nil {
		return isRegistered, err
	}
	if count > 0 {
		isRegistered = true
	}

	return isRegistered, nil
}

// SaveBaseInfo 施設基本情報を更新
func (f *facilityTemaUsecase) SaveBaseInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, request *facility.SaveBaseInfoInput) error {
	// トランザクション生成
	tx, txErr := f.FRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	txFacilityRepo := nInfra.NewFacilityRepository(tx)
	txAccountRepo := aInfra.NewAccountTemaRepository(tx)

	// 施設情報更新
	if err := txFacilityRepo.UpdateProperty(&facility.HtTmProperties{
		PropertyID:        request.PropertyID,
		Name:              request.Name,
		PostalCode:        request.PostalCode,
		StateProvinceName: request.StateProvinceName,
		City:              request.City,
		Line1:             request.Line1,
		Line2:             request.Line2,
		Line3:             request.Line3,
		Phone:             request.Phone,
		Fax:               request.Fax,
	}); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}
	if err := txFacilityRepo.UpsertPropertyLangsBase(&facility.HtTmPropertyLangs{
		PropertyID:        request.PropertyID,
		Name:              request.Name,
		PostalCode:        request.PostalCode,
		StateProvinceName: request.StateProvinceName,
		City:              request.City,
		Line1:             request.Line1,
		Line2:             request.Line2,
		Line3:             request.Line3,
	}); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}

	// パスワードの暗号化
	loginPWEnc, eErr := utils.Encrypt(os.Getenv("TEMA_LOGIN_PW"))
	if eErr != nil {
		f.FRepository.TxRollback(tx)
		return eErr
	}
	passwordEnc, eErr := utils.Encrypt(request.ConnectPassword)
	if eErr != nil {
		f.FRepository.TxRollback(tx)
		return eErr
	}

	// urlsに格納するJSONデータ準備
	urlList, _ := json.Marshal(map[string]string{
		"GetBookingResultRQ": os.Getenv("TEMA_GET_BOOKING_RESULT_RQ"),
		"GetRoomListRQ":      os.Getenv("TEMA_GET_ROOM_LIST_RQ"),
		"GetPlanListRQ":      os.Getenv("TEMA_GET_PLAN_LIST_RQ"),
		"GetAriListRQ":       os.Getenv("TEMA_GET_ARI_LIST_RQ"),
		"GetPriceListRQ":     os.Getenv("TEMA_GET_PRICE_LIST_RQ"),
	})
	assignData := &account.HtTmWholesalerApiAccounts{
		PropertyID:  request.PropertyID,
		Name:        request.Name,
		LoginID:     os.Getenv("TEMA_LOGIN_ID"),
		LoginPWEnc:  loginPWEnc,
		Username:    request.ConnectID,
		PasswordEnc: passwordEnc,
		Urls:        string(urlList),
	}

	connectUser, _ := f.ATemaRepository.FetchConnectUser(request.PropertyID)
	if connectUser.WholesalerAccountID == 0 {
		assignData.CreatedAt = time.Now()
	}

	if err := txAccountRepo.UpsertConnectUser(assignData); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}
	// コミットとロールバック
	if err := f.FRepository.TxCommit(tx); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}
	return nil
}

// SaveDetail 施設詳細情報を更新
func (f *facilityTemaUsecase) SaveDetail(request *facility.SaveDetailInput) error {
	// トランザクション生成
	tx, txErr := f.FRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	txFacilityRepo := nInfra.NewFacilityTemaRepository(tx)

	// 施設情報更新
	if err := txFacilityRepo.UpsertPropertyTema(&facility.HtTmPropertyTemas{
		PropertyID:                   request.PropertyID,
		CheckinBegin:                 request.CheckinBegin,
		CheckinEnd:                   request.CheckinEnd,
		Checkout:                     request.Checkout,
		Instructions:                 utils.ConvertNewlineCodeToBrTag(request.Instructions),
		SpecialInstructions:          utils.ConvertNewlineCodeToBrTag(request.SpecialInstructions),
		PolicyKnowBeforeYouGo:        utils.ConvertNewlineCodeToBrTag(request.PolicyKnowBeforeYouGo),
		FeeMandatory:                 utils.ConvertNewlineCodeToBrTag(request.FeeMandatory),
		FeeOptional:                  utils.ConvertNewlineCodeToBrTag(request.FeeOptional),
		DescriptionAmenity:           utils.ConvertNewlineCodeToBrTag(request.DescriptionAmenity),
		DescriptionAttractions:       utils.ConvertNewlineCodeToBrTag(request.DescriptionAttractions),
		DescriptionBusinessAmenities: utils.ConvertNewlineCodeToBrTag(request.DescriptionBusinessAmenities),
		DescriptionDining:            utils.ConvertNewlineCodeToBrTag(request.DescriptionDining),
		DescriptionLocation:          utils.ConvertNewlineCodeToBrTag(request.DescriptionLocation),
		DescriptionHeadline:          utils.ConvertNewlineCodeToBrTag(request.DescriptionHeadline),
		DescriptionRooms:             utils.ConvertNewlineCodeToBrTag(request.DescriptionRooms),
	}); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}

	// アメニティの情報をリセット
	if err := txFacilityRepo.ClearPropertyAmenity(request.PropertyID); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}
	var insertAmenities []facility.HtTmPropertyTemaUseAmenity
	for _, amenityData := range request.Amenities {
		insertAmenities = append(insertAmenities, facility.HtTmPropertyTemaUseAmenity{
			PropertyID:        request.PropertyID,
			PropertyAmenityID: amenityData.PropertyAmenityID,
		})
	}

	// アメニティの情報を登録
	if len(insertAmenities) > 0 {
		if err := txFacilityRepo.CreatePropertyAmenity(insertAmenities); err != nil {
			f.FRepository.TxRollback(tx)
			return err
		}
	}

	// コミットとロールバック
	if err := f.FRepository.TxCommit(tx); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}
	return nil
}

// FetchAllAmenities 施設アメニティを取得
func (f *facilityTemaUsecase) FetchAllAmenities() ([]facility.Amenity, error) {
	res := []facility.Amenity{}
	amenities, err := f.FTlRepository.FetchAllAmenities()
	if err != nil {
		return res, err
	}
	for _, amenityData := range *amenities {
		res = append(res, facility.Amenity{
			PropertyAmenityID: amenityData.PropertyAmenityID,
			AmenityName:       amenityData.AmenityName,
		})
	}
	return res, nil
}

// 元々フリー入力だった項目で時刻が全角で登録されていてプルダウン表示できないケースが少なくないため、半角に変換する
func (f *facilityTemaUsecase) convertUpperTimeToLower(time string) string {
	// 全角コロンを半角コロンに変換
	time = strings.Replace(time, "：", ":", -1)
	// 全角数字を半角数字に変換
	var numConv = unicode.SpecialCase{
		unicode.CaseRange{
			Lo: 0xff10, // '０'
			Hi: 0xff19, // '９'
			Delta: [unicode.MaxCase]rune{
				0,
				0x0030 - 0xff10, // '0' - '０'
				0,
			},
		},
	}
	return strings.ToLowerSpecial(numConv, time)
}
