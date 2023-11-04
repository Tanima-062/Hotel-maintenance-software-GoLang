package usecase

import (
	"fmt"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	iInfra "github.com/Adventureinc/hotel-hm-api/src/image/infra"
	pInfra "github.com/Adventureinc/hotel-hm-api/src/plan/infra"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"gorm.io/gorm"
)

// roomRaku2Usecase らく通部屋関連usecase
type roomRaku2Usecase struct {
	RRaku2Repository room.IRoomRaku2Repository
	IRaku2Repository image.IImageRaku2Repository
}

func (r *roomRaku2Usecase) CreateOrUpdateBulk(request []room.RoomData) error {
	//TODO implement me
	panic("implement me")
}

// NewRoomRaku2Usecase インスタンス生成
func NewRoomRaku2Usecase(db *gorm.DB) room.IRoomUsecase {
	return &roomRaku2Usecase{
		RRaku2Repository: rInfra.NewRoomRaku2Repository(db),
		IRaku2Repository: iInfra.NewImageRaku2Repository(db),
	}
}

// FetchList 一覧取得
func (r *roomRaku2Usecase) FetchList(request *room.ListInput) ([]room.ListOutput, error) {
	response := []room.ListOutput{}
	rooms, roomErr := r.RRaku2Repository.FetchRoomsByPropertyID(*request)
	if roomErr != nil {
		return response, roomErr
	}

	roomTypeIDList := []int64{}
	for _, value := range rooms {
		roomTypeIDList = append(roomTypeIDList, value.RoomTypeID)
	}

	amenityCh := make(chan []room.RoomAmenitiesRaku2)
	imageCh := make(chan []image.RoomImagesOutput)
	go r.fetchAmenitiesByRoomTypeID(amenityCh, roomTypeIDList)
	go r.fetchImagesByRoomTypeID(imageCh, roomTypeIDList)
	amenities, images := <-amenityCh, <-imageCh

	for _, v := range rooms {
		roomResponse := room.ListOutput{RoomTypeTable: v.RoomTypeTable}
		for _, imageData := range images {
			if imageData.RoomTypeID != v.RoomTypeID {
				continue
			}
			roomResponse.ImageLength++
			if imageData.Order == 1 {
				roomResponse.Href = imageData.Href
			}
		}
		for _, amenityData := range amenities {
			if amenityData.RoomTypeID == v.RoomTypeID {
				roomResponse.AmenityNames = append(roomResponse.AmenityNames, amenityData.Raku2RoomAmenityName)
			}
		}
		response = append(response, roomResponse)
	}

	return response, nil
}

// FetchAllAmenities アメニティ一覧取得
func (r *roomRaku2Usecase) FetchAllAmenities() ([]room.AllAmenitiesOutput, error) {
	response := []room.AllAmenitiesOutput{}
	amenities, amenitiesErr := r.RRaku2Repository.FetchAllAmenities()
	if amenitiesErr != nil {
		return response, amenitiesErr
	}
	for _, v := range amenities {
		response = append(response, room.AllAmenitiesOutput{
			AmenityID: v.Raku2RoomAmenityID,
			Name:      v.Raku2RoomAmenityName,
		})
	}
	return response, nil
}

// FetchDetail 部屋詳細取得
func (r *roomRaku2Usecase) FetchDetail(request *room.DetailInput) (*room.DetailOutput, error) {
	response := &room.DetailOutput{}
	if r.RRaku2Repository.MatchesRoomTypeIDAndPropertyID(request.RoomTypeID, request.PropertyID) == false {
		return response, fmt.Errorf("Error: %s", "この施設ではこの部屋を閲覧できません。")
	}

	roomDetail, roomErr := r.RRaku2Repository.FetchRoomByRoomTypeID(request.RoomTypeID)
	if roomErr != nil {
		return response, roomErr
	}

	amenityCh := make(chan []room.RoomAmenitiesRaku2)
	imageCh := make(chan []image.RoomImagesOutput)
	go r.fetchAmenitiesByRoomTypeID(amenityCh, []int64{request.RoomTypeID})
	go r.fetchImagesByRoomTypeID(imageCh, []int64{request.RoomTypeID})
	amenities, images := <-amenityCh, <-imageCh

	response.RoomTypeTable = roomDetail.RoomTypeTable

	for _, amenityData := range amenities {
		response.AmenityIDList = append(response.AmenityIDList, amenityData.Raku2RoomAmenityID)
	}
	response.Images = images

	return response, nil
}

// Create 部屋作成
func (r *roomRaku2Usecase) Create(request *room.SaveInput) error {
	// 部屋コードの重複チェック
	duplicate := r.RRaku2Repository.CountRoomTypeCode(request.PropertyID, request.RoomTypeCode)
	if duplicate > 0 {
		return fmt.Errorf("DuplicateError")
	}

	// トランザクション生成
	tx, txErr := r.RRaku2Repository.TxStart()
	if txErr != nil {
		return txErr
	}
	roomTxRepo := rInfra.NewRoomRaku2Repository(tx)

	// 部屋作成
	roomTable := &room.HtTmRoomTypeRaku2s{
		RoomTypeTable: room.RoomTypeTable{
			PropertyID:              request.PropertyID,
			RoomTypeCode:            request.RoomTypeCode,
			Name:                    request.Name,
			RoomKindID:              request.RoomKindID,
			RoomDesc:                request.RoomDesc,
			StockSettingStart:       request.StockSettingStart,
			StockSettingEnd:         request.StockSettingEnd,
			IsSettingStockYearRound: request.IsSettingStockYearRound,
			RoomCount:               request.RoomCount,
			OcuMin:                  request.OcuMin,
			OcuMax:                  request.OcuMax,
			IsSmoking:               request.IsSmoking,
			IsStopSales:             request.IsStopSales,
			IsDelete:                request.IsDelete,
			Times:                   common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
		},
	}
	if err := roomTxRepo.CreateRoomRaku2(roomTable); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}

	// 一度アメニティを全件削除してから登録し直す
	if err := roomTxRepo.ClearRoomToAmenities(roomTable.RoomTypeID); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}
	for _, amenityID := range request.AmenityIDList {
		if err := roomTxRepo.CreateRoomToAmenities(roomTable.RoomTypeID, amenityID); err != nil {
			r.RRaku2Repository.TxRollback(tx)
			return err
		}
	}

	// 部屋と画像を紐付ける
	iRepo := iInfra.NewImageRaku2Repository(tx)
	for _, imageData := range request.Images {
		record := []image.HtTmRoomOwnImagesRaku2s{}

		record = append(record, image.HtTmRoomOwnImagesRaku2s{
			RoomImageRaku2ID: imageData.ImageID,
			RoomTypeID:       roomTable.RoomTypeID,
			Order:            imageData.Order,
		})

		if err := iRepo.CreateRoomOwnImagesRaku2(record); err != nil {
			r.RRaku2Repository.TxRollback(tx)
			return err
		}
	}

	// コミットとロールバック
	if err := r.RRaku2Repository.TxCommit(tx); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}

	return nil
}

// Update 部屋更新
func (r *roomRaku2Usecase) Update(request *room.SaveInput) error {
	// トランザクション生成
	tx, txErr := r.RRaku2Repository.TxStart()
	if txErr != nil {
		return txErr
	}
	roomTxRepo := rInfra.NewRoomRaku2Repository(tx)
	// 部屋更新
	roomTable := &room.HtTmRoomTypeRaku2s{
		RoomTypeTable: room.RoomTypeTable{
			RoomTypeID:              request.RoomTypeID,
			PropertyID:              request.PropertyID,
			RoomTypeCode:            request.RoomTypeCode,
			Name:                    request.Name,
			RoomKindID:              request.RoomKindID,
			RoomDesc:                request.RoomDesc,
			StockSettingStart:       request.StockSettingStart,
			StockSettingEnd:         request.StockSettingEnd,
			IsSettingStockYearRound: request.IsSettingStockYearRound,
			RoomCount:               request.RoomCount,
			OcuMin:                  request.OcuMin,
			OcuMax:                  request.OcuMax,
			IsSmoking:               request.IsSmoking,
			IsStopSales:             request.IsStopSales,
			IsDelete:                request.IsDelete,
			Times:                   common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
		},
	}
	if err := roomTxRepo.UpdateRoomRaku2(roomTable); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}

	// 一度アメニティを全件削除してから登録し直す
	if err := roomTxRepo.ClearRoomToAmenities(roomTable.RoomTypeID); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}
	for _, amenityID := range request.AmenityIDList {
		if err := roomTxRepo.CreateRoomToAmenities(roomTable.RoomTypeID, amenityID); err != nil {
			r.RRaku2Repository.TxRollback(tx)
			return err
		}
	}

	// 画像を一度削除して、部屋と画像を再度紐付ける
	iRepo := iInfra.NewImageRaku2Repository(tx)
	if err := iRepo.ClearRoomImage(roomTable.RoomTypeID); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}

	for _, imageData := range request.Images {
		record := []image.HtTmRoomOwnImagesRaku2s{}

		record = append(record, image.HtTmRoomOwnImagesRaku2s{
			RoomImageRaku2ID: imageData.ImageID,
			RoomTypeID:       roomTable.RoomTypeID,
			Order:            imageData.Order,
		})

		if err := iRepo.CreateRoomOwnImagesRaku2(record); err != nil {
			r.RRaku2Repository.TxRollback(tx)
			return err
		}
	}

	// コミットとロールバック
	if err := r.RRaku2Repository.TxCommit(tx); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}
	return nil
}

// Delete 部屋削除
func (r *roomRaku2Usecase) Delete(roomTypeID int64) error {
	// トランザクション生成
	tx, txErr := r.RRaku2Repository.TxStart()
	if txErr != nil {
		return txErr
	}
	roomTxRepo := rInfra.NewRoomRaku2Repository(tx)
	// 部屋論理削除
	if err := roomTxRepo.DeleteRoomRaku2(roomTypeID); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}
	// 部屋に紐づくプランを論理削除
	planTxRepo := pInfra.NewPlanRaku2Repository(tx)
	if err := planTxRepo.DeletePlanByRoomTypeID(roomTypeID); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}
	// コミットとロールバック
	if err := r.RRaku2Repository.TxCommit(tx); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}
	return nil
}

// UpdateStopSales 部屋と紐づくプラン・在庫の売止更新
func (r *roomRaku2Usecase) UpdateStopSales(request *room.StopSalesInput) error {
	// トランザクション生成
	tx, txErr := r.RRaku2Repository.TxStart()
	if txErr != nil {
		return txErr
	}
	// 部屋の売止
	roomTxRepo := rInfra.NewRoomRaku2Repository(tx)
	if err := roomTxRepo.UpdateStopSales(request.RoomTypeID, request.IsStopSales); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}

	// 部屋に紐づくプランの売止
	planTxRepo := pInfra.NewPlanRaku2Repository(tx)
	planList, pErr := planTxRepo.FetchAllByRoomTypeID(request.RoomTypeID)
	if pErr != nil {
		r.RRaku2Repository.TxRollback(tx)
		return pErr
	}
	planIDList := []int64{}
	for _, v := range planList {
		planIDList = append(planIDList, v.PlanID)
	}
	if err := planTxRepo.UpdateStopSales(planIDList, request.IsStopSales); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}

	// コミットとロールバック
	if err := r.RRaku2Repository.TxCommit(tx); err != nil {
		r.RRaku2Repository.TxRollback(tx)
		return err
	}
	return nil
}

func (r *roomRaku2Usecase) fetchAmenitiesByRoomTypeID(ch chan<- []room.RoomAmenitiesRaku2, roomTypeIDList []int64) {
	res, err := r.RRaku2Repository.FetchAmenitiesByRoomTypeID(roomTypeIDList)
	if err != nil {
		ch <- []room.RoomAmenitiesRaku2{}
	}
	ch <- res
}

func (r *roomRaku2Usecase) fetchImagesByRoomTypeID(ch chan<- []image.RoomImagesOutput, roomTypeIDList []int64) {
	res, err := r.IRaku2Repository.FetchImagesByRoomTypeID(roomTypeIDList)
	if err != nil {
		ch <- []image.RoomImagesOutput{}
	}
	ch <- res
}
