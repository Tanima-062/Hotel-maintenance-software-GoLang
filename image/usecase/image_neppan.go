package usecase

import (
	"encoding/json"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/account"
	aInfra "github.com/Adventureinc/hotel-hm-api/src/account/infra"
	"github.com/Adventureinc/hotel-hm-api/src/common/utils"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	iInfra "github.com/Adventureinc/hotel-hm-api/src/image/infra"
	"gorm.io/gorm"
)

// imageNeppanUsecase ねっぱん画像関連usecase
type imageNeppanUsecase struct {
	ARepository       account.IAccountRepository
	INeppanRepository image.IImageNeppanRepository
	ImageStorage      image.IImageStorage
}

// NewImageNeppanUsecase インスタンス生成
func NewImageNeppanUsecase(db *gorm.DB) image.IImageUsecase {
	return &imageNeppanUsecase{
		ARepository:       aInfra.NewAccountRepository(db),
		INeppanRepository: iInfra.NewImageNeppanRepository(db),
		ImageStorage:      iInfra.NewImageStorage(),
	}
}

// FetchAll 画像一覧
func (i *imageNeppanUsecase) FetchAll(request *image.ListInput) ([]image.ImagesOutput, error) {
	response := []image.ImagesOutput{}

	neppanImages, err := i.INeppanRepository.FetchImagesByPropertyID(request.PropertyID)
	if err != nil {
		return response, nil
	}
	for _, v := range neppanImages {
		response = append(response, image.ImagesOutput{
			ImageID:    v.RoomImageNeppanID,
			PropertyID: v.PropertyID,
			Method:     v.Method,
			Href:       v.Href,
			CategoryCd: v.CategoryCd,
			IsMain:     v.IsMain,
			Caption:    v.Caption,
			SortNum:    v.MainOrder,
			GcsInfo:    v.GcsInfo,
		})
	}
	return response, nil
}

// Update 画像情報更新
func (i *imageNeppanUsecase) Update(request *image.UpdateInput) error {
	// 更新前の対象データ取得
	imageData, err := i.INeppanRepository.FetchImageByRoomImageNeppanID(request.ImageID)
	if err != nil {
		return err
	}

	// トランザクション生成
	tx, txErr := i.INeppanRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	imageTxRepo := iInfra.NewImageNeppanRepository(tx)

	// 更新
	if err := imageTxRepo.UpdateImageNeppan(request); err != nil {
		i.INeppanRepository.TxRollback(tx)
		return err
	}

	// メイン画像から外す
	if !request.IsMain && imageData.IsMain {
		// メイン画像取得
		mainImages, fetchImagesErr := imageTxRepo.FetchMainImages(imageData.PropertyID)
		if fetchImagesErr != nil {
			i.INeppanRepository.TxRollback(tx)
			return fetchImagesErr
		}
		input := []image.UpdateSortNumInput{}
		for index, mainImage := range mainImages {
			sortNum := index + 1
			// sortNumとMainOrderが一致しないデータだけ更新（ソート順の並び替え）
			if sortNum != mainImage.MainOrder {
				input = append(input, image.UpdateSortNumInput{
					ImageID: mainImage.RoomImageNeppanID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := imageTxRepo.UpdateImageSort(&input); err != nil {
				i.INeppanRepository.TxRollback(tx)
				return err
			}
		}
	}

	// category_cdを変更する場合のみ削除処理実行
	if request.CategoryCd != imageData.CategoryCd {
		// 対象データ取得
		roomImages, err := imageTxRepo.FetchRoomImagesByRoomImageNeppanID(request.ImageID)
		if err != nil {
			i.INeppanRepository.TxRollback(tx)
			return err
		}
		if len(roomImages) > 0 {
			// 対象データ削除
			if err := imageTxRepo.ClearRoomImageByRoomImageNeppanID(request.ImageID); err != nil {
				i.INeppanRepository.TxRollback(tx)
				return err
			}
			for _, roomImage := range roomImages {
				// RoomTypeIDで紐づいてるデータ取得
				roomImagesByRoomTypeId, _ := imageTxRepo.FetchRoomImagesByRoomTypeID(roomImage.RoomTypeID)
				for index, value := range roomImagesByRoomTypeId {
					orderNum := uint8(index) + 1
					// orderを並び替えて更新（orderNumとvalue.Orderが一致しないものだけ更新）
					if orderNum != value.Order {
						input := &image.HtTmRoomOwnImagesNeppans{
							RoomOwnImagesID: value.RoomOwnImagesID,
							Order:           orderNum,
						}
						if err := imageTxRepo.UpdateRoomImageOrder(input); err != nil {
							i.INeppanRepository.TxRollback(tx)
							return err
						}
					}
				}
			}
		}
	}
	// コミットとロールバック
	if err := i.INeppanRepository.TxCommit(tx); err != nil {
		i.INeppanRepository.TxRollback(tx)
		return err
	}
	return nil
}

// UpdateIsMain メイン画像フラグ更新
func (i *imageNeppanUsecase) UpdateIsMain(request *image.UpdateIsMainInput) error {
	// ソート順更新
	if err := i.INeppanRepository.UpdateIsMain(request); err != nil {
		return err
	}
	// 画像情報取得
	imageData, fetchErr := i.INeppanRepository.FetchImageByRoomImageNeppanID(request.ImageID)
	if fetchErr != nil {
		return fetchErr
	}
	// メイン画像から外す
	if !request.IsMain {
		// メイン画像取得
		mainImages, fetchImagesErr := i.INeppanRepository.FetchMainImages(imageData.PropertyID)
		if fetchImagesErr != nil {
			return fetchImagesErr
		}

		input := []image.UpdateSortNumInput{}
		for index, mainImage := range mainImages {
			// ソート順設定
			sortNum := index + 1
			// 取得した画像と一致しなければsortNumの値を設定
			if mainImage.MainOrder != sortNum {
				input = append(input, image.UpdateSortNumInput{
					ImageID: mainImage.RoomImageNeppanID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := i.INeppanRepository.TrUpdateImageSort(&input); err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateSortNum メイン画像の並び順更新
func (i *imageNeppanUsecase) UpdateSortNum(request *[]image.UpdateSortNumInput) error {
	return i.INeppanRepository.TrUpdateImageSort(request)
}

// Delete 画像削除
func (i *imageNeppanUsecase) Delete(imageID int64) error {
	imageData, fetchErr := i.INeppanRepository.FetchImageByRoomImageNeppanID(imageID)
	if fetchErr != nil {
		return fetchErr
	}
	gcsInfo := &image.GcsInfo{}
	if err := json.Unmarshal([]byte(imageData.GcsInfo), gcsInfo); err != nil {
		return err
	}

	// トランザクション生成
	tx, txErr := i.INeppanRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	imageTxRepo := iInfra.NewImageNeppanRepository(tx)

	if err := imageTxRepo.DeleteImage(imageID); err != nil {
		i.INeppanRepository.TxRollback(tx)
		return err
	}

	// メイン画像に設定していたら外す
	if imageData.IsMain {
		// メイン画像取得
		mainImages, fetchImagesErr := imageTxRepo.FetchMainImages(imageData.PropertyID)
		if fetchImagesErr != nil {
			i.INeppanRepository.TxRollback(tx)
			return fetchImagesErr
		}
		input := []image.UpdateSortNumInput{}
		for index, mainImage := range mainImages {
			sortNum := index + 1
			// sortNumとMainOrderが一致しないデータだけ更新（ソート順の並び替え）
			if sortNum != mainImage.MainOrder {
				input = append(input, image.UpdateSortNumInput{
					ImageID: mainImage.RoomImageNeppanID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := imageTxRepo.UpdateImageSort(&input); err != nil {
				i.INeppanRepository.TxRollback(tx)
				return err
			}
		}
	}

	// 対象データ取得
	roomImages, err := imageTxRepo.FetchRoomImagesByRoomImageNeppanID(imageID)
	if err != nil {
		i.INeppanRepository.TxRollback(tx)
		return err
	}
	if len(roomImages) > 0 {
		// 対象データ削除
		if err := imageTxRepo.ClearRoomImageByRoomImageNeppanID(imageID); err != nil {
			i.INeppanRepository.TxRollback(tx)
			return err
		}
		for _, roomImage := range roomImages {
			// RoomTypeIDで紐づいてるデータ取得
			roomImagesByRoomTypeId, _ := imageTxRepo.FetchRoomImagesByRoomTypeID(roomImage.RoomTypeID)
			for index, value := range roomImagesByRoomTypeId {
				orderNum := uint8(index) + 1
				// orderを並び替えて更新（orderNumとvalue.Orderが一致しないものだけ更新）
				if orderNum != value.Order {
					input := &image.HtTmRoomOwnImagesNeppans{
						RoomOwnImagesID: value.RoomOwnImagesID,
						Order:           orderNum,
					}
					if err := imageTxRepo.UpdateRoomImageOrder(input); err != nil {
						i.INeppanRepository.TxRollback(tx)
						return err
					}
				}
			}
		}
	}

	// コミットとロールバック
	if err := i.INeppanRepository.TxCommit(tx); err != nil {
		i.INeppanRepository.TxRollback(tx)
		return err
	}
	// GCSの画像を削除
	if err := i.ImageStorage.Delete(gcsInfo.Bucket, gcsInfo.Name); err != nil {
		return err
	}

	return nil
}

// Create 画像作成
func (i *imageNeppanUsecase) Create(request *image.UploadInput, file *multipart.FileHeader, hmUser account.HtTmHotelManager) error {
	// upload情報設定する
	bucketName := os.Getenv("GCS_PROP_IMG_BUCKET_NAME")
	formatedTime := time.Now().Format("20060102150405") // 時間フォーマット
	filename := strconv.FormatInt(hmUser.PropertyID, 10) + "_" + formatedTime + utils.GetExtensionFromContentType(request.ContentType)
	fileAllPath := utils.ImgBasePath + "/" + strconv.FormatInt(hmUser.WholesalerID, 10) + "/" + filename

	// upload
	rawGcsInfo, upErr := i.ImageStorage.Create(bucketName, fileAllPath, file)
	if upErr != nil {
		return upErr
	}

	gcsInfoBytes, marshalErr := json.Marshal(map[string]string{"name": rawGcsInfo.Name, "bucket": rawGcsInfo.Bucket, "contentType": rawGcsInfo.ContentType})
	if marshalErr != nil {
		return marshalErr
	}

	upData := &image.HtTmImageNeppans{
		PropertyID: hmUser.PropertyID,
		Method:     "GET",
		Href:       "https://storage.googleapis.com/" + bucketName + "/" + fileAllPath,
		CategoryCd: request.CategoryCd,
		IsMain:     request.IsMain,
		Caption:    request.Caption,
		MainOrder:  request.SortNum,
		GcsInfo:    string(gcsInfoBytes),
	}
	upData.CreatedAt = time.Now()
	upData.UpdatedAt = time.Now()
	return i.INeppanRepository.CreateImage(upData)
}

// CountMainImages メイン画像の数を数える
func (i *imageNeppanUsecase) CountMainImages(propertyID int64) int64 {
	return i.INeppanRepository.CountMainImagesPerPropertyID(propertyID)
}
