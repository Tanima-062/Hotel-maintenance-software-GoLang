package usecase

import (
	"strings"
	"unicode"
	"errors"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aInfra "github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	nInfra "github.com/Adventureinc/hotel-hm-api/src/facility/infra"
	tInfra "github.com/Adventureinc/hotel-hm-api/src/facility/infra"
	"gorm.io/gorm"
)

// facilityTlUsecase TLの施設関連usecase
type facilityTlUsecase struct {
	ARepository   account.IAccountRepository
	FRepository   facility.IFacilityRepository
	FTlRepository facility.IFacilityTlRepository
}

// NewFacilityTlUsecase インスタンス生成
func NewFacilityTlUsecase(db *gorm.DB) facility.IFacilityUsecase {
	return &facilityTlUsecase{
		ARepository:   aInfra.NewAccountRepository(db),
		FRepository:   tInfra.NewFacilityRepository(db),
		FTlRepository: tInfra.NewFacilityTlRepository(db),
	}
}

// FetchAll アカウントに紐づく施設情報をすべて取得
func (f *facilityTlUsecase) FetchAll(hmUser account.HtTmHotelManager) ([]facility.InitFacilityOutput, error) {
	// 子施設の場合は当該施設の情報のみ返却する。
	if hmUser.PropertyID != facility.ParentPropertyId {
		return f.FTlRepository.FetchAllFacilities([]int64{hmUser.PropertyID}, hmUser.WholesalerID)
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
	return f.FTlRepository.FetchAllFacilities(propertyIDs, hmUser.WholesalerID)
}

// UpdateDispPriority サイト公開フラグを更新
func (f *facilityTlUsecase) UpdateDispPriority(request *facility.UpdateDispPriorityInput) error {
	return f.FTlRepository.UpdateDispPriority(request.PropertyID, request.DispPriority)
}

// FetchBaseInfo 施設基本情報を取得
func (f *facilityTlUsecase) FetchBaseInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, request *facility.BaseInfoInput) (*facility.BaseInfoOutput, error) {
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
		ConnectID:         "",
		ConnectPassword:   "",
	}, nil
}

// FetchDetail 施設詳細情報を取得
func (f *facilityTlUsecase) FetchDetail(request *facility.BaseInfoInput) (*facility.DetailOutput, error) {
	response := &facility.DetailOutput{}
	property, pErr := f.FTlRepository.FetchPropertyDetail(request.PropertyID)
	if pErr != nil {
		return &facility.DetailOutput{}, pErr
	}
	amenities, _ := f.FTlRepository.FetchAmenities(request.PropertyID)
	response.CheckinBegin = f.convertUpperTimeToLower(property.CheckinBegin)
	response.CheckinEnd = f.convertUpperTimeToLower(property.CheckinEnd)
	response.Checkout = f.convertUpperTimeToLower(property.Checkout)
	response.Instructions = utils.ConvertBrTagToNewlineCode(property.Instructions)
	response.SpecialInstructions = utils.ConvertBrTagToNewlineCode(property.SpecialInstructions)
	response.PolicyKnowBeforeYouGo = utils.ConvertBrTagToNewlineCode(property.PolicyKnowBeforeYouGo)
	response.FeeMandatory = utils.ConvertBrTagToNewlineCode(property.FeeMandatory)
	response.DescriptionAmenity = utils.ConvertBrTagToNewlineCode(property.DescriptionAmenity)
	response.DescriptionAttractions = utils.ConvertBrTagToNewlineCode(property.DescriptionAttractions)
	response.DescriptionBusinessAmenities = utils.ConvertBrTagToNewlineCode(property.DescriptionBusinessAmenities)
	response.DescriptionDining = utils.ConvertBrTagToNewlineCode(property.DescriptionDining)
	response.DescriptionLocation = utils.ConvertBrTagToNewlineCode(property.DescriptionLocation)
	response.DescriptionHeadline = utils.ConvertBrTagToNewlineCode(property.DescriptionHeadline)
	response.DescriptionRooms = utils.ConvertBrTagToNewlineCode(property.DescriptionRooms)
	response.FeeOptional = utils.ConvertBrTagToNewlineCode(property.FeeOptional)

	for _, amenityData := range *amenities {
		response.Amenities = append(response.Amenities, facility.Amenity{
			PropertyAmenityID: amenityData.PropertyAmenityID,
			AmenityName:       amenityData.AmenityName,
		})
	}
	return response, nil
}

// SaveBaseInfo 施設基本情報を更新
func (f *facilityTlUsecase) SaveBaseInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, request *facility.SaveBaseInfoInput) error {
	// トランザクション生成
	tx, txErr := f.FRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	txFacilityRepo := nInfra.NewFacilityRepository(tx)

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

	// コミットとロールバック
	if err := f.FRepository.TxCommit(tx); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}
	return nil
}

// SaveDetail 施設詳細情報を更新
func (f *facilityTlUsecase) SaveDetail(request *facility.SaveDetailInput) error {
	// トランザクション生成
	tx, txErr := f.FRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	txFacilityRepo := nInfra.NewFacilityTlRepository(tx)

	// 施設情報更新
	if err := txFacilityRepo.UpsertPropertyTl(&facility.HtTmPropertyTls{
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
	var insertAmenities []facility.HtTmPropertyTlsUseAmenity
	for _, amenityData := range request.Amenities {
		insertAmenities = append(insertAmenities, facility.HtTmPropertyTlsUseAmenity{
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
func (f *facilityTlUsecase) FetchAllAmenities() ([]facility.Amenity, error) {
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
func (f *facilityTlUsecase) convertUpperTimeToLower(time string) string {
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
