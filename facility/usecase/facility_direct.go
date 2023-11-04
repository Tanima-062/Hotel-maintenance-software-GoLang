package usecase

import (
	"strconv"
	"strings"
	"unicode"
	"errors"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aInfra "github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	fInfra "github.com/Adventureinc/hotel-hm-api/src/facility/infra"
	nInfra "github.com/Adventureinc/hotel-hm-api/src/facility/infra"
	"gorm.io/gorm"
)

// facilityDirectUsecase 直仕入れの施設関連usecase
type facilityDirectUsecase struct {
	ARepository       account.IAccountRepository
	FRepository       facility.IFacilityRepository
	FDirectRepository facility.IFacilityDirectRepository
}

// NewFacilityDirectUsecase インスタンス生成
func NewFacilityDirectUsecase(db *gorm.DB) facility.IFacilityUsecase {
	return &facilityDirectUsecase{
		ARepository:       aInfra.NewAccountRepository(db),
		FRepository:       fInfra.NewFacilityRepository(db),
		FDirectRepository: fInfra.NewFacilityDirectRepository(db),
	}
}

// FetchAll アカウントに紐づく施設情報をすべて取得
func (f *facilityDirectUsecase) FetchAll(hmUser account.HtTmHotelManager) ([]facility.InitFacilityOutput, error) {
	// 子施設の場合は当該施設の情報のみ返却する。
	if hmUser.PropertyID != facility.ParentPropertyId {
		return f.FDirectRepository.FetchAllFacilities([]int64{hmUser.PropertyID}, hmUser.WholesalerID)
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
	return f.FDirectRepository.FetchAllFacilities(propertyIDs, hmUser.WholesalerID)
}

// UpdateDispPriority サイト公開フラグを更新
func (f *facilityDirectUsecase) UpdateDispPriority(request *facility.UpdateDispPriorityInput) error {
	// ht_tm_property_directsにレコードがないケースがあるため、存在しない場合表示する瞬間に最低限のレコードを作成する
	if _, err := f.FDirectRepository.FirstOrCreate(request.PropertyID); err != nil {
		return err
	}
	return f.FDirectRepository.UpdateDispPriority(request.PropertyID, request.DispPriority)
}

// FetchBaseInfo 施設基本情報を取得
func (f *facilityDirectUsecase) FetchBaseInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, request *facility.BaseInfoInput) (*facility.BaseInfoOutput, error) {
	// ht_tm_property_directsにレコードがないケースがあるため、存在しない場合表示する瞬間に最低限のレコードを作成する
	if _, err := f.FDirectRepository.FirstOrCreate(request.PropertyID); err != nil {
		return &facility.BaseInfoOutput{}, err
	}
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
func (f *facilityDirectUsecase) FetchDetail(request *facility.BaseInfoInput) (*facility.DetailOutput, error) {
	response := &facility.DetailOutput{}
	property, pErr := f.FDirectRepository.FetchPropertyDetail(request.PropertyID)
	if pErr != nil {
		return response, pErr
	}
	amenities, _ := f.FDirectRepository.FetchAmenities(request.PropertyID)
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
			PropertyAmenityID: strconv.FormatInt(amenityData.PropertyAmenityID, 10),
			AmenityName:       amenityData.AmenityName,
		})
	}
	return response, nil
}

// SaveBaseInfo 施設基本情報を更新
func (f *facilityDirectUsecase) SaveBaseInfo(hmUser *account.HtTmHotelManager, claimParam *account.ClaimParam, request *facility.SaveBaseInfoInput) error {
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
func (f *facilityDirectUsecase) SaveDetail(request *facility.SaveDetailInput) error {
	// トランザクション生成
	tx, txErr := f.FRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	txFacilityDirectRepo := fInfra.NewFacilityDirectRepository(tx)

	// 施設情報更新
	if err := txFacilityDirectRepo.UpsertPropertyDirect(&facility.HtTmPropertyDirects{
		PropertyID:                   request.PropertyID,
		FeeMandatory:                 utils.ConvertNewlineCodeToBrTag(request.FeeMandatory),
		FeeOptional:                  utils.ConvertNewlineCodeToBrTag(request.FeeOptional),
		DescriptionAmenity:           utils.ConvertNewlineCodeToBrTag(request.DescriptionAmenity),
		DescriptionAttractions:       utils.ConvertNewlineCodeToBrTag(request.DescriptionAttractions),
		DescriptionBusinessAmenities: utils.ConvertNewlineCodeToBrTag(request.DescriptionBusinessAmenities),
		DescriptionDining:            utils.ConvertNewlineCodeToBrTag(request.DescriptionDining),
		DescriptionLocation:          utils.ConvertNewlineCodeToBrTag(request.DescriptionLocation),
		DescriptionHeadline:          utils.ConvertNewlineCodeToBrTag(request.DescriptionHeadline),
		DescriptionRooms:             utils.ConvertNewlineCodeToBrTag(request.DescriptionRooms),
		CheckinBegin:                 request.CheckinBegin,
		CheckinEnd:                   request.CheckinEnd,
		Checkout:                     request.Checkout,
		Instructions:                 utils.ConvertNewlineCodeToBrTag(request.Instructions),
		SpecialInstructions:          utils.ConvertNewlineCodeToBrTag(request.SpecialInstructions),
		PolicyKnowBeforeYouGo:        utils.ConvertNewlineCodeToBrTag(request.PolicyKnowBeforeYouGo),
	}); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}

	// アメニティの情報をリセット
	if err := txFacilityDirectRepo.ClearPropertyAmenity(request.PropertyID); err != nil {
		f.FRepository.TxRollback(tx)
		return err
	}
	var insertAmenities []facility.HtTmPropertyDirectsUseAmenity
	for _, amenityData := range request.Amenities {
		insertAmenities = append(insertAmenities, facility.HtTmPropertyDirectsUseAmenity{
			PropertyID:        request.PropertyID,
			PropertyAmenityID: amenityData.PropertyAmenityID,
		})
	}
	// アメニティの情報を登録
	if len(insertAmenities) > 0 {
		if err := txFacilityDirectRepo.CreatePropertyAmenity(insertAmenities); err != nil {
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

// FetchAllAmenities 施設のアメニティを取得
func (f *facilityDirectUsecase) FetchAllAmenities() ([]facility.Amenity, error) {
	res := []facility.Amenity{}
	amenities, err := f.FDirectRepository.FetchAllAmenities()
	if err != nil {
		return res, err
	}
	for _, amenityData := range *amenities {
		res = append(res, facility.Amenity{
			PropertyAmenityID: strconv.FormatInt(amenityData.PropertyAmenityID, 10),
			AmenityName:       amenityData.AmenityName,
		})
	}
	return res, nil
}

// 元々フリー入力だった項目で時刻が全角で登録されていてプルダウン表示できないケースが少なくないため、半角に変換する
func (f *facilityDirectUsecase) convertUpperTimeToLower(time string) string {
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

