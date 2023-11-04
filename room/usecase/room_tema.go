package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	iInfra "github.com/Adventureinc/hotel-hm-api/src/image/infra"
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"gorm.io/gorm"
	"time"
)

type RoomTemaUseCase struct {
	RTemaRepository room.IRoomTemaRepository
	ITemaRepository image.IImageTemaRepository
}

func NewRoomTemaUseCase(db *gorm.DB) room.IRoomTemaUseCase {
	return &RoomTemaUseCase{
		RTemaRepository: infra.NewRoomTemaRepository(db),
		ITemaRepository: iInfra.NewImageTemaRepository(db),
	}
}
func (r *RoomTemaUseCase) FetchList(request *room.ListInput) ([]room.ListOutputTema, error) {
	response := []room.ListOutputTema{}
	rooms, roomErr := r.RTemaRepository.FetchRoomsByPropertyID(*request)
	if roomErr != nil {
		return response, roomErr
	}

	roomTypeIDList := []int64{}
	for _, value := range rooms {
		roomTypeIDList = append(roomTypeIDList, value.RoomTypeID)
	}

	amenityCh := make(chan []room.RoomAmenitiesTema)
	imageCh := make(chan []room.RoomImagesTema)
	go r.fetchAmenitiesByRoomTypeID(amenityCh, roomTypeIDList)
	go r.fetchImagesByRoomTypeID(imageCh, roomTypeIDList)
	amenities, images := <-amenityCh, <-imageCh

	for _, v := range rooms {
		roomResponse := room.ListOutputTema{RoomTypeTema: v.RoomTypeTema}
		for _, imageData := range images {
			if imageData.RoomTypeID != v.RoomTypeID {
				continue
			}
			roomResponse.ImageLength++
			if imageData.RoomTypeID == v.RoomTypeID && len(roomResponse.Images) == 0 {
				roomResponse.Href = imageData.Url
				roomResponse.Images = append(roomResponse.Images, imageData)
			}
		}
		for _, amenityData := range amenities {
			if amenityData.RoomTypeID == v.RoomTypeID {
				roomResponse.AmenityNames = append(roomResponse.AmenityNames, amenityData.TemaRoomAmenityName)
				roomResponse.AmenityIDs = append(roomResponse.AmenityIDs, amenityData.TemaRoomAmenityID)
			}
		}
		response = append(response, roomResponse)
	}

	return response, nil
}

func (r *RoomTemaUseCase) CreateOrUpdateBulk(request []room.RoomDataTema) error {
	// transaction generation
	tx, txErr := r.RTemaRepository.TxStart()
	if txErr != nil {
		return txErr
	}
	//Bulk data insert from request
	for _, data := range request {
		roomTable := &room.HtTmRoomTypeTemas{
			RoomTypeTema: room.RoomTypeTema{
				PropertyID:   data.PropertyID,
				RoomTypeCode: data.RoomTypeCode,
				Name:         data.Name,
				RoomKindID:   data.RoomKindID,
				RoomDesc:     data.RoomDesc,
				OcuMin:       data.OcuMin,
				OcuMax:       data.OcuMax,
				IsStopSales:  data.IsStopSales,
				Times:        common.Times{UpdatedAt: time.Now(), CreatedAt: time.Now()},
			},
		}

		// Room code duplication check
		roomType, _ := r.RTemaRepository.FetchRoomTypeIDByRoomTypeCode(data.PropertyID, data.RoomTypeCode)
		// Checking if the room type is already in the database
		if (roomType != room.HtTmRoomTypeTemas{}) {
			roomTable.RoomTypeID = roomType.RoomTypeID
			// Update `HtTmRoomTypeTemas`
			if err := r.RTemaRepository.UpdateRoomBulkTema(roomTable); err != nil {
				r.RTemaRepository.TxRollback(tx)
				return err
			}
		} else {
			// Insert into `HtTmRoomTypeTemas`
			if err := r.RTemaRepository.CreateRoomBulkTema(roomTable); err != nil {
				r.RTemaRepository.TxRollback(tx)
				return err
			}
		}

		// Delete all amenities and then register again
		if err := r.RTemaRepository.ClearRoomToAmenities(roomTable.RoomTypeID); err != nil {
			r.RTemaRepository.TxRollback(tx)
			return err
		}

		// Insert amenities
		for _, amenityID := range data.AmenityIDList {
			if err := r.RTemaRepository.CreateRoomToAmenities(roomTable.RoomTypeID, int64(amenityID)); err != nil {
				r.RTemaRepository.TxRollback(tx)
				return err
			}
		}

		// Delete the image once and associate the room and the image again
		if err := r.RTemaRepository.ClearRoomImage(roomTable.RoomTypeID); err != nil {
			r.RTemaRepository.TxRollback(tx)
			return err
		}

		for _, imageData := range data.Images {
			var record []room.HtTmRoomOwnImagesTemas
			record = append(record, room.HtTmRoomOwnImagesTemas{
				RoomImageTemaID: int64(imageData.ImageID),
				RoomTypeID:      roomTable.RoomTypeID,
				Order:           uint8(imageData.Order),
			})

			if err := r.RTemaRepository.CreateRoomOwnImages(record); err != nil {
				r.RTemaRepository.TxRollback(tx)
				return err
			}
		}
	}

	// commit and rollback
	if err := r.RTemaRepository.TxCommit(tx); err != nil {
		r.RTemaRepository.TxRollback(tx)
		return err
	}
	return nil
}

func (r *RoomTemaUseCase) FetchAllAmenities() ([]room.AllAmenitiesOutput, error) {
	response := []room.AllAmenitiesOutput{}
	amenities, amenitiesErr := r.RTemaRepository.FetchAllAmenities()
	if amenitiesErr != nil {
		return response, amenitiesErr
	}
	for _, v := range amenities {
		response = append(response, room.AllAmenitiesOutput{
			AmenityID: v.TemaRoomAmenityID,
			Name:      v.TemaRoomAmenityName,
		})
	}
	return response, nil
}

// FetchDetail Get room details
func (r *RoomTemaUseCase) FetchDetail(request *room.DetailInput) (*room.TemaDetailOutput, error) {
	response := &room.TemaDetailOutput{}
	//fetch a single room details
	roomDetail, roomErr := r.RTemaRepository.FetchRoomByRoomTypeID(request.RoomTypeID)
	if roomErr != nil {
		return response, roomErr
	}

	amenityCh := make(chan []room.RoomAmenitiesTema)
	imageCh := make(chan []room.RoomImagesTema)
	go r.fetchAmenitiesByRoomTypeID(amenityCh, []int64{request.RoomTypeID})
	go r.fetchImagesByRoomTypeID(imageCh, []int64{request.RoomTypeID})
	amenities, images := <-amenityCh, <-imageCh

	response.RoomTypeTema = roomDetail.RoomTypeTema
	for _, amenityData := range amenities {
		response.AmenityIDList = append(response.AmenityIDList, amenityData.TemaRoomAmenityID)
	}
	response.Images = images

	return response, nil
}

func (r *RoomTemaUseCase) fetchAmenitiesByRoomTypeID(ch chan<- []room.RoomAmenitiesTema, roomTypeIDList []int64) {
	res, err := r.RTemaRepository.FetchAmenitiesByRoomTypeID(roomTypeIDList)
	if err != nil {
		ch <- []room.RoomAmenitiesTema{}
	}
	ch <- res
}

func (r *RoomTemaUseCase) fetchImagesByRoomTypeID(ch chan<- []room.RoomImagesTema, roomTypeIDList []int64) {
	res, err := r.RTemaRepository.FetchImagesByRoomTypeID(roomTypeIDList)
	if err != nil {
		ch <- []room.RoomImagesTema{}
	}
	ch <- res
}
