package usecase

import (
	"fmt"
	"github.com/Adventureinc/hotel-hm-api/src/stock"
	sInfra "github.com/Adventureinc/hotel-hm-api/src/stock/infra"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	iInfra "github.com/Adventureinc/hotel-hm-api/src/image/infra"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"gorm.io/gorm"
)

// roomTlUsecase TL room related usecase
type roomTlUsecase struct {
	RTlRepository room.IRoomTlRepository
	ITlRepository image.IImageTlRepository
	STlRepository stock.IStockTlRepository
}

// NewRoomTlUsecase instantiation
func NewRoomTlUsecase(db *gorm.DB) room.IRoomBulkUsecase {
	return &roomTlUsecase{
		RTlRepository: rInfra.NewRoomTlRepository(db),
		ITlRepository: iInfra.NewImageTlRepository(db),
		STlRepository: sInfra.NewStockTlRepository(db),
	}
}

// FetchList List acquisition
func (r *roomTlUsecase) FetchList(request *room.ListInput) ([]room.ListOutputTl, error) {
	response := []room.ListOutputTl{}
	rooms, roomErr := r.RTlRepository.FetchRoomsByPropertyID(*request)
	if roomErr != nil {
		return response, roomErr
	}

	roomTypeIDList := []int64{}
	for _, value := range rooms {
		roomTypeIDList = append(roomTypeIDList, value.RoomTypeID)
	}

	amenityCh := make(chan []room.RoomAmenitiesTl)
	imageCh := make(chan []image.RoomImagesOutput)
	go r.fetchAmenitiesByRoomTypeID(amenityCh, roomTypeIDList)
	go r.fetchImagesByRoomTypeID(imageCh, roomTypeIDList)
	amenities, images := <-amenityCh, <-imageCh

	for _, v := range rooms {
		roomResponse := room.ListOutputTl{RoomTypeTable: v.RoomTypeTable}
		for _, imageData := range images {
			if imageData.RoomTypeID != v.RoomTypeID {
				continue
			}
			roomResponse.ImageLength++
			if imageData.RoomTypeID == v.RoomTypeID && len(roomResponse.Images) == 0 {
				roomResponse.Href = imageData.Href
				roomResponse.Images = append(roomResponse.Images, imageData)
			}
		}
		for _, amenityData := range amenities {
			if amenityData.RoomTypeID == v.RoomTypeID {
				roomResponse.AmenityNames = append(roomResponse.AmenityNames, amenityData.TlsRoomAmenityName)
				roomResponse.AmenityIDs = append(roomResponse.AmenityIDs, amenityData.TlsRoomAmenityID)
			}
		}
		response = append(response, roomResponse)
	}

	return response, nil
}

// FetchAllAmenities Get amenity list
func (r *roomTlUsecase) FetchAllAmenities() ([]room.AllAmenitiesOutput, error) {
	response := []room.AllAmenitiesOutput{}
	amenities, amenitiesErr := r.RTlRepository.FetchAllAmenities()
	if amenitiesErr != nil {
		return response, amenitiesErr
	}
	for _, v := range amenities {
		response = append(response, room.AllAmenitiesOutput{
			AmenityID: v.TlsRoomAmenityID,
			Name:      v.TlsRoomAmenityName,
		})
	}
	return response, nil
}

// FetchDetail Get room details
func (r *roomTlUsecase) FetchDetail(request *room.DetailInput) (*room.DetailOutput, error) {
	response := &room.DetailOutput{}
	if r.RTlRepository.MatchesRoomTypeIDAndPropertyID(request.RoomTypeID, request.PropertyID) == false {
		return response, fmt.Errorf("Error: %s", "この施設ではこの部屋を閲覧できません。")
	}

	roomDetail, roomErr := r.RTlRepository.FetchRoomByRoomTypeID(request.RoomTypeID)
	if roomErr != nil {
		return response, roomErr
	}

	amenityCh := make(chan []room.RoomAmenitiesTl)
	imageCh := make(chan []image.RoomImagesOutput)
	go r.fetchAmenitiesByRoomTypeID(amenityCh, []int64{request.RoomTypeID})
	go r.fetchImagesByRoomTypeID(imageCh, []int64{request.RoomTypeID})
	amenities, images := <-amenityCh, <-imageCh

	response.RoomTypeTable = roomDetail.RoomTypeTable

	for _, amenityData := range amenities {
		response.AmenityIDList = append(response.AmenityIDList, amenityData.TlsRoomAmenityID)
	}
	response.Images = images

	return response, nil
}

// Create room creation
func (r *roomTlUsecase) Create(request *room.SaveInput) error {
	// Room code duplication check
	duplicate := r.RTlRepository.CountRoomTypeCode(request.PropertyID, request.RoomTypeCode)
	if duplicate > 0 {
		return fmt.Errorf("DuplicateError")
	}

	// transaction generation
	tx, txErr := r.RTlRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	roomTxRepo := rInfra.NewRoomTlRepository(tx)

	// room creation
	roomTable := &room.HtTmRoomTypeTls{
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

	if err := roomTxRepo.CreateRoomTl(roomTable); err != nil {
		r.RTlRepository.TxRollback(tx) //done
		return err
	}

	// Delete all amenities and then register again
	if err := roomTxRepo.ClearRoomToAmenities(roomTable.RoomTypeID); err != nil {
		r.RTlRepository.TxRollback(tx)
		return err
	}
	for _, amenityID := range request.AmenityIDList {
		if err := roomTxRepo.CreateRoomToAmenities(roomTable.RoomTypeID, amenityID); err != nil {
			r.RTlRepository.TxRollback(tx)
			return err
		}
	}

	// Associating a room with an image
	iRepo := iInfra.NewImageTlRepository(tx)
	for _, imageData := range request.Images {
		record := []image.HtTmRoomOwnImagesTls{}

		record = append(record, image.HtTmRoomOwnImagesTls{
			RoomImageTlID: imageData.ImageID,
			RoomTypeID:    roomTable.RoomTypeID,
			Order:         imageData.Order,
		})

		if err := iRepo.CreateRoomOwnImagesTl(record); err != nil {
			r.RTlRepository.TxRollback(tx) //done
			return err
		}
	}

	// commit and rollback
	if err := r.RTlRepository.TxCommit(tx); err != nil {
		r.RTlRepository.TxRollback(tx)
		return err
	}

	return nil
}

// Update room update
func (r *roomTlUsecase) Update(request *room.SaveInput) error {
	// transaction generation
	tx, txErr := r.RTlRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	roomTxRepo := rInfra.NewRoomTlRepository(tx)
	// room update
	roomTable := &room.HtTmRoomTypeTls{
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
	if err := roomTxRepo.UpdateRoomTl(roomTable); err != nil {
		r.RTlRepository.TxRollback(tx)
		return err
	}

	// Delete all amenities and then register again
	if err := roomTxRepo.ClearRoomToAmenities(roomTable.RoomTypeID); err != nil {
		r.RTlRepository.TxRollback(tx)
		return err
	}
	for _, amenityID := range request.AmenityIDList {
		if err := roomTxRepo.CreateRoomToAmenities(roomTable.RoomTypeID, amenityID); err != nil {
			r.RTlRepository.TxRollback(tx)
			return err
		}
	}

	// Delete the image once and associate the room and the image again
	iRepo := iInfra.NewImageTlRepository(tx)
	if err := iRepo.ClearRoomImage(roomTable.RoomTypeID); err != nil {
		r.RTlRepository.TxRollback(tx)
		return err
	}

	for _, imageData := range request.Images {
		record := []image.HtTmRoomOwnImagesTls{}

		record = append(record, image.HtTmRoomOwnImagesTls{
			RoomImageTlID: imageData.ImageID,
			RoomTypeID:    roomTable.RoomTypeID,
			Order:         imageData.Order,
		})

		if err := iRepo.CreateRoomOwnImagesTl(record); err != nil {
			r.RTlRepository.TxRollback(tx)
			return err
		}
	}

	// commit and rollback
	if err := r.RTlRepository.TxCommit(tx); err != nil {
		r.RTlRepository.TxRollback(tx)
		return err
	}
	return nil
}

// Delete delete room
func (r *roomTlUsecase) Delete(roomTypeID int64) error {
	return nil
}

func (r *roomTlUsecase) fetchAmenitiesByRoomTypeID(ch chan<- []room.RoomAmenitiesTl, roomTypeIDList []int64) {
	res, err := r.RTlRepository.FetchAmenitiesByRoomTypeID(roomTypeIDList)
	if err != nil {
		ch <- []room.RoomAmenitiesTl{}
	}
	ch <- res
}

func (r *roomTlUsecase) fetchImagesByRoomTypeID(ch chan<- []image.RoomImagesOutput, roomTypeIDList []int64) {
	res, err := r.ITlRepository.FetchImagesByRoomTypeID(roomTypeIDList)
	if err != nil {
		ch <- []image.RoomImagesOutput{}
	}
	ch <- res
}

// CreateOrUpdateBulk creates or updates
func (r *roomTlUsecase) CreateOrUpdateBulk(request []room.RoomData) error {
	// transaction generation
	tx, txErr := r.RTlRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	//Bulk data insert from request
	for _, data := range request {
		roomTable := &room.HtTmRoomTypeTls{
			RoomTypeTable: room.RoomTypeTable{
				PropertyID:              data.PropertyID,
				RoomTypeCode:            data.RoomTypeCode,
				Name:                    data.Name,
				RoomKindID:              data.RoomKindID,
				RoomDesc:                data.RoomDesc,
				StockSettingStart:       data.StockSettingStart,
				StockSettingEnd:         data.StockSettingEnd,
				IsSettingStockYearRound: data.IsSettingStockYearRound,
				RoomCount:               data.RoomCount,
				OcuMin:                  data.OcuMin,
				OcuMax:                  data.OcuMax,
				IsSmoking:               data.IsSmoking,
				IsStopSales:             data.IsStopSales,
				IsDelete:                data.IsDelete,
				Times:                   common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
			},
		}

		// Room code duplication check
		roomType, _ := r.RTlRepository.FetchRoomTypeIdByRoomTypeCode(data.PropertyID, data.RoomTypeCode)

		// Checking if the room type is already in the database
		if (roomType != room.HtTmRoomTypeTls{}) {
			roomTable.RoomTypeID = roomType.RoomTypeID
			// Update `RoomTypeTls`
			if err := r.RTlRepository.UpdateRoomBulkTl(roomTable); err != nil {
				r.RTlRepository.TxRollback(tx)
				return err
			}
		} else {
			// Insert into `RoomTypeTls`
			if err := r.RTlRepository.CreateRoomBulkTl(roomTable); err != nil {
				r.RTlRepository.TxRollback(tx)
				return err
			}
		}

		// Delete all amenities and then register again
		if err := r.RTlRepository.ClearRoomToAmenities(roomTable.RoomTypeID); err != nil {
			r.RTlRepository.TxRollback(tx)
			return err
		}
		// Insert amenities
		for _, amenityID := range data.AmenityIDList {
			if err := r.RTlRepository.CreateRoomToAmenities(roomTable.RoomTypeID, int64(amenityID)); err != nil {
				r.RTlRepository.TxRollback(tx)
				return err
			}
		}

		// Delete the image once and associate the room and the image again
		if err := r.ITlRepository.ClearRoomImage(roomTable.RoomTypeID); err != nil {
			r.ITlRepository.TxRollback(tx)
			return err
		}
		for _, imageData := range data.Images {
			var record []image.HtTmRoomOwnImagesTls
			record = append(record, image.HtTmRoomOwnImagesTls{
				RoomImageTlID: int64(imageData.ImageID),
				RoomTypeID:    roomTable.RoomTypeID,
				Order:         uint8(imageData.Order),
			})

			if err := r.ITlRepository.CreateRoomOwnImagesTl(record); err != nil {
				r.ITlRepository.TxRollback(tx)
				return err
			}
		}

	}

	// commit and rollback
	if err := r.RTlRepository.TxCommit(tx); err != nil {
		r.RTlRepository.TxRollback(tx)
		return err
	}

	return nil
}
