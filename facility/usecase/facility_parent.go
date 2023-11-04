package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/account"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/facility"
	fInfra "github.com/Adventureinc/hotel-hm-api/src/facility/infra"
	"gorm.io/gorm"
)

// facilityParentUrsecase 親アカウントから取ってくる処理
type facilityParentUsecase struct {
	FRepository       facility.IFacilityRepository
	FTemaRepository   facility.IFacilityTemaRepository
	FTlRepository     facility.IFacilityTlRepository
	FDirectRepository facility.IFacilityDirectRepository
	FNeppanRepository facility.IFacilityNeppanRepository
}

// NewFacilityTlUsecase インスタンス生成
func NewFacilityParentUsecase(db *gorm.DB) facility.IParentUsecase {
	// リポジトリの初期化
	return &facilityParentUsecase{
		FRepository:       fInfra.NewFacilityRepository(db),
		FTemaRepository:   fInfra.NewFacilityTemaRepository(db),
		FTlRepository:     fInfra.NewFacilityTlRepository(db),
		FDirectRepository: fInfra.NewFacilityDirectRepository(db),
		FNeppanRepository: fInfra.NewFacilityNeppanRepository(db),
	}
}

func (f *facilityParentUsecase) FetchAll(hmUser account.HtTmHotelManager) ([]facility.InitFacilityOutput, error) {
	r, _ := f.FRepository.FetchAllClientCompanies(hmUser.HotelManagerID)

	propertyIDs := []int64{}
	response := []facility.InitFacilityOutput{}
	for _, v := range r {
		propertyIDs = append(propertyIDs, v.PropertyID)
	}
	facilities, err := f.FRepository.FetchAllFacilitiesByPropertyID(propertyIDs)
	if err != nil {
		return response, err
	}

	// ホールセラーIDを見て施設詳細情報取得、DispPriority(公開・非公開フラグ)を格納する
	for _, v := range facilities {
		switch v.WholesalerID {
		case utils.WholesalerIDTema:
			property, pErr := f.FTemaRepository.FetchPropertyDetail(v.PropertyID)
			if pErr != nil {
				continue
			}
			v.DispPriority = property.DispPriority
		case utils.WholesalerIDTl:
			property, pErr := f.FTlRepository.FetchPropertyDetail(v.PropertyID)
			if pErr != nil {
				continue
			}
			v.DispPriority = property.DispPriority
		case utils.WholesalerIDNeppan:
			property, pErr := f.FNeppanRepository.FetchPropertyDetail(v.PropertyID)
			if pErr != nil {
				continue
			}
			v.DispPriority = property.DispPriority
		case utils.WholesalerIDDirect:
			property, pErr := f.FDirectRepository.FetchPropertyDetail(v.PropertyID)
			if pErr != nil {
				continue
			}
			v.DispPriority = property.DispPriority
		}
		response = append(response, v)
	}
	return response, nil
}
