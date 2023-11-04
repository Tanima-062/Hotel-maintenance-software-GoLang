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
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	sInfra "github.com/Adventureinc/hotel-hm-api/src/stock/infra"
	"gorm.io/gorm"
)

// roomDirectUsecase 直仕入れ部屋関連usecase
type roomDirectUsecase struct {
	RDirectRepository room.IRoomDirectRepository
	IDirectRepository image.IImageDirectRepository
}

func (r *roomDirectUsecase) CreateOrUpdateBulk(request []room.RoomData) error {
	//TODO implement me
	panic("implement me")
}

// NewRoomDirectUsecase インスタンス生成
func NewRoomDirectUsecase(db *gorm.DB) room.IRoomUsecase {
	return &roomDirectUsecase{
		RDirectRepository: rInfra.NewRoomDirectRepository(db),
		IDirectRepository: iInfra.NewImageDirectRepository(db),
	}
}

// FetchList 一覧取得
func (r *roomDirectUsecase) FetchList(request *room.ListInput) ([]room.ListOutput, error) {
	response := []room.ListOutput{}
	rooms, roomErr := r.RDirectRepository.FetchRoomsByPropertyID(*request)
	if roomErr != nil {
		return response, roomErr
	}
	roomTypeIDList := []int64{}
	for _, value := range rooms {
		roomTypeIDList = append(roomTypeIDList, value.RoomTypeID)
	}

	amenityCh := make(chan []room.RoomAmenitiesDirect)
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
				roomResponse.AmenityNames = append(roomResponse.AmenityNames, amenityData.DirectRoomAmenityName)
			}
		}
		response = append(response, roomResponse)
	}

	return response, nil
}

// FetchAllAmenities アメニティ一覧取得
func (r *roomDirectUsecase) FetchAllAmenities() ([]room.AllAmenitiesOutput, error) {
	response := []room.AllAmenitiesOutput{}
	amenities, amenitiesErr := r.RDirectRepository.FetchAllAmenities()
	if amenitiesErr != nil {
		return response, amenitiesErr
	}
	for _, v := range amenities {
		response = append(response, room.AllAmenitiesOutput{
			AmenityID: v.DirectRoomAmenityID,
			Name:      v.DirectRoomAmenityName,
		})
	}
	return response, nil
}

// FetchDetail 部屋詳細取得
func (r *roomDirectUsecase) FetchDetail(request *room.DetailInput) (*room.DetailOutput, error) {
	response := &room.DetailOutput{}
	if r.RDirectRepository.MatchesRoomTypeIDAndPropertyID(request.RoomTypeID, request.PropertyID) == false {
		return response, fmt.Errorf("Error: %s", "この施設ではこの部屋を閲覧できません。")
	}

	roomDetail, roomErr := r.RDirectRepository.FetchRoomByRoomTypeID(request.RoomTypeID)
	if roomErr != nil {
		return response, roomErr
	}

	amenityCh := make(chan []room.RoomAmenitiesDirect)
	imageCh := make(chan []image.RoomImagesOutput)
	go r.fetchAmenitiesByRoomTypeID(amenityCh, []int64{request.RoomTypeID})
	go r.fetchImagesByRoomTypeID(imageCh, []int64{request.RoomTypeID})
	amenities, images := <-amenityCh, <-imageCh

	response.RoomTypeTable = roomDetail.RoomTypeTable

	for _, amenityData := range amenities {
		response.AmenityIDList = append(response.AmenityIDList, amenityData.DirectRoomAmenityID)
	}
	response.Images = images

	return response, nil
}

// Create 部屋作成
func (r *roomDirectUsecase) Create(request *room.SaveInput) error {
	// 部屋コードの重複チェック
	duplicate := r.RDirectRepository.CountRoomTypeCode(request.PropertyID, request.RoomTypeCode)
	if duplicate > 0 {
		return fmt.Errorf("DuplicateError")
	}

	// トランザクション生成
	tx, txErr := r.RDirectRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	roomTxRepo := rInfra.NewRoomDirectRepository(tx)

	// 部屋作成
	roomTable := &room.HtTmRoomTypeDirects{
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
	if err := roomTxRepo.CreateRoomDirect(roomTable); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}

	// 一度アメニティを全件削除してから登録し直す
	if err := roomTxRepo.ClearRoomToAmenities(roomTable.RoomTypeID); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}
	for _, amenityID := range request.AmenityIDList {
		if err := roomTxRepo.CreateRoomToAmenities(roomTable.RoomTypeID, amenityID); err != nil {
			r.RDirectRepository.TxRollback(tx)
			return err
		}
	}

	// 部屋と画像を紐付ける
	if len(request.Images) > 0 {
		iRepo := iInfra.NewImageDirectRepository(tx)
		for _, imageData := range request.Images {
			record := []image.HtTmRoomOwnImagesDirects{}

			record = append(record, image.HtTmRoomOwnImagesDirects{
				RoomImageDirectID: imageData.ImageID,
				RoomTypeID:        roomTable.RoomTypeID,
				Order:             imageData.Order,
			})

			if err := iRepo.CreateRoomOwnImagesDirect(record); err != nil {
				r.RDirectRepository.TxRollback(tx)
				return err
			}
		}
	}

	// 部屋を登録した際、提供数を初期値として在庫に初期データを投入する
	insertStock := []stock.HtTmStockDirects{}
	now := time.Now()
	start := request.StockSettingStart
	end := request.StockSettingEnd
	if request.IsSettingStockYearRound == true {
		// 通年設定の場合は現在から１年分の期間を設定
		start = now
		end = now.AddDate(1, 0, 0)
	} else {
		// 在庫設定期間開始日付が現在より前の場合、現在から設定開始
		if request.StockSettingStart.Before(now) {
			start = now
		} else {
			start = request.StockSettingStart
		}
		// 終了日の在庫も作成するように+1日する
		end = request.StockSettingEnd.AddDate(0, 0, 1)
	}

	stockTxRepo := sInfra.NewStockDirectRepository(tx)
	if start.Before(end) {
		for date := start; date.Before(end); date = date.AddDate(0, 0, 1) {
			insertStock = append(insertStock, stock.HtTmStockDirects{
				StockTable: stock.StockTable{
					RoomTypeID: roomTable.RoomTypeID,
					RoomCount:  request.RoomCount,
					Stock:      request.RoomCount,
					UseDate:    date,
				}})
		}
		if err := stockTxRepo.CreateStocks(insertStock); err != nil {
			r.RDirectRepository.TxRollback(tx)
			return err
		}
	}

	// コミットとロールバック
	if err := r.RDirectRepository.TxCommit(tx); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}
	return nil
}

// Update 部屋更新
func (r *roomDirectUsecase) Update(request *room.SaveInput) error {
	// トランザクション生成
	tx, txErr := r.RDirectRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	roomTxRepo := rInfra.NewRoomDirectRepository(tx)
	// 部屋更新
	roomTable := &room.HtTmRoomTypeDirects{
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
	if err := roomTxRepo.UpdateRoomDirect(roomTable); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}

	// 一度アメニティを全件削除してから登録し直す
	if err := roomTxRepo.ClearRoomToAmenities(roomTable.RoomTypeID); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}
	for _, amenityID := range request.AmenityIDList {
		if err := roomTxRepo.CreateRoomToAmenities(roomTable.RoomTypeID, amenityID); err != nil {
			r.RDirectRepository.TxRollback(tx)
			return err
		}
	}

	// 画像を一度削除して、部屋と画像を再度紐付ける
	iRepo := iInfra.NewImageDirectRepository(tx)
	if err := iRepo.ClearRoomImage(roomTable.RoomTypeID); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}

	for _, imageData := range request.Images {
		record := []image.HtTmRoomOwnImagesDirects{}

		record = append(record, image.HtTmRoomOwnImagesDirects{
			RoomImageDirectID: imageData.ImageID,
			RoomTypeID:        roomTable.RoomTypeID,
			Order:             imageData.Order,
		})

		if err := iRepo.CreateRoomOwnImagesDirect(record); err != nil {
			r.RDirectRepository.TxRollback(tx)
			return err
		}
	}

	stockTxRepo := sInfra.NewStockDirectRepository(tx)

	// 部屋を登録した際、提供数を初期値として在庫に初期データを投入する
	// 既存データ更新と新規データ作成の振り分けはstockTxRepo.UpsertStocks()内で行うので、ここでは新規データ作成の準備のみする
	now := time.Now()
	start := request.StockSettingStart
	end := request.StockSettingEnd

	if request.IsSettingStockYearRound == true {
		// 通年設定の場合は現在から１年分の期間を設定
		start = now
		end = now.AddDate(1, 0, 0)
	} else {
		// 在庫設定期間開始日付が現在より前の場合、現在から設定開始
		if request.StockSettingStart.Before(now) {
			start = now
		} else {
			start = request.StockSettingStart
		}
		// 終了日の在庫も作成するように+1日する
		end = request.StockSettingEnd.AddDate(0, 0, 1)
	}

	var inputData []stock.HtTmStockDirects
	if start.Before(end) {
		for date := start; date.Before(end); date = date.AddDate(0, 0, 1) {
			inputData = append(inputData, stock.HtTmStockDirects{
				StockTable: stock.StockTable{
					RoomTypeID:  request.RoomTypeID,
					UseDate:     date,
					RoomCount:   request.RoomCount,
					Stock:       request.RoomCount,
					IsStopSales: request.IsStopSales,
					Times: common.Times{
						UpdatedAt: time.Now(),
					},
				},
			})
		}
		// 在庫データ作成・更新
		if updateErr := stockTxRepo.UpsertStocks(inputData); updateErr != nil {
			r.RDirectRepository.TxRollback(tx)
			return updateErr
		}
	}

	// コミットとロールバック
	if err := r.RDirectRepository.TxCommit(tx); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}
	return nil
}

// Delete 部屋削除
func (r *roomDirectUsecase) Delete(roomTypeID int64) error {
	// トランザクション生成
	tx, txErr := r.RDirectRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	roomTxRepo := rInfra.NewRoomDirectRepository(tx)
	// 部屋論理削除
	if err := roomTxRepo.DeleteRoomDirect(roomTypeID); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}
	// 部屋に紐づくプランを論理削除
	planTxRepo := pInfra.NewPlanDirectRepository(tx)
	if err := planTxRepo.DeletePlanByRoomTypeID(roomTypeID); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}
	// コミットとロールバック
	if err := r.RDirectRepository.TxCommit(tx); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}
	return nil
}

// UpdateStopSales 部屋の売止更新
func (r *roomDirectUsecase) UpdateStopSales(request *room.StopSalesInput) error {
	// トランザクション生成
	tx, txErr := r.RDirectRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	// 部屋の売止
	roomTxRepo := rInfra.NewRoomDirectRepository(tx)
	if err := roomTxRepo.UpdateStopSales(request.RoomTypeID, request.IsStopSales); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}

	// 部屋に紐づくプランの売止
	planTxRepo := pInfra.NewPlanDirectRepository(tx)
	planList, pErr := planTxRepo.FetchAllByRoomTypeID(request.RoomTypeID)
	if pErr != nil {
		r.RDirectRepository.TxRollback(tx)
		return pErr
	}
	planIDList := []int64{}
	for _, v := range planList {
		planIDList = append(planIDList, v.PlanID)
	}
	if err := planTxRepo.UpdateStopSales(planIDList, request.IsStopSales); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}

	// コミットとロールバック
	if err := r.RDirectRepository.TxCommit(tx); err != nil {
		r.RDirectRepository.TxRollback(tx)
		return err
	}
	return nil
}

func (r *roomDirectUsecase) fetchAmenitiesByRoomTypeID(ch chan<- []room.RoomAmenitiesDirect, roomTypeIDList []int64) {
	res, err := r.RDirectRepository.FetchAmenitiesByRoomTypeID(roomTypeIDList)
	if err != nil {
		ch <- []room.RoomAmenitiesDirect{}
	}
	ch <- res
}

func (r *roomDirectUsecase) fetchImagesByRoomTypeID(ch chan<- []image.RoomImagesOutput, roomTypeIDList []int64) {
	res, err := r.IDirectRepository.FetchImagesByRoomTypeID(roomTypeIDList)
	if err != nil {
		ch <- []image.RoomImagesOutput{}
	}
	ch <- res
}
