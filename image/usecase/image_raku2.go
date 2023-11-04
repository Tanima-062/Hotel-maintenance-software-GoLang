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

// imageRaku2Usecase らく通画像関連usecase
type imageRaku2Usecase struct {
	ARepository      account.IAccountRepository
	IRaku2Repository image.IImageRaku2Repository
	ImageStorage     image.IImageStorage
}

// NewImageRaku2Usecase インスタンス生成
func NewImageRaku2Usecase(db *gorm.DB) image.IImageUsecase {
	return &imageRaku2Usecase{
		ARepository:      aInfra.NewAccountRepository(db),
		IRaku2Repository: iInfra.NewImageRaku2Repository(db),
		ImageStorage:     iInfra.NewImageStorage(),
	}
}

// FetchAll 画像一覧
func (i *imageRaku2Usecase) FetchAll(request *image.ListInput) ([]image.ImagesOutput, error) {
	response := []image.ImagesOutput{}

	raku2Images, err := i.IRaku2Repository.FetchImagesByPropertyID(request.PropertyID)
	if err != nil {
		return response, nil
	}
	for _, v := range raku2Images {
		response = append(response, image.ImagesOutput{
			ImageID:    v.RoomImageRaku2ID,
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
func (i *imageRaku2Usecase) Update(request *image.UpdateInput) error {
	// 更新前の対象データ取得
	imageData, err := i.IRaku2Repository.FetchImageByRoomImageRaku2ID(request.ImageID)
	if err != nil {
		return err
	}

	// トランザクション生成
	tx, txErr := i.IRaku2Repository.TxStart()
	if txErr != nil {
		return txErr
	}

	imageTxRepo := iInfra.NewImageRaku2Repository(tx)

	// 更新
	if err := imageTxRepo.UpdateImageRaku2(request); err != nil {
		i.IRaku2Repository.TxRollback(tx)
		return err
	}

	// メイン画像から外す
	if !request.IsMain && imageData.IsMain {
		// メイン画像取得
		mainImages, fetchImagesErr := imageTxRepo.FetchMainImages(imageData.PropertyID)
		if fetchImagesErr != nil {
			i.IRaku2Repository.TxRollback(tx)
			return fetchImagesErr
		}
		input := []image.UpdateSortNumInput{}
		for index, mainImage := range mainImages {
			sortNum := index + 1
			// sortNumとMainOrderが一致しないデータだけ更新（ソート順の並び替え）
			if sortNum != mainImage.MainOrder {
				input = append(input, image.UpdateSortNumInput{
					ImageID: mainImage.RoomImageRaku2ID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := imageTxRepo.UpdateImageSort(&input); err != nil {
				i.IRaku2Repository.TxRollback(tx)
				return err
			}
		}
	}

	// category_cdを変更する場合のみ削除処理実行
	if request.CategoryCd != imageData.CategoryCd {
		// 対象データ取得
		roomImages, err := imageTxRepo.FetchRoomImagesByRoomImageRaku2ID(request.ImageID)
		if err != nil {
			i.IRaku2Repository.TxRollback(tx)
			return err
		}
		if len(roomImages) > 0 {
			// 対象データ削除
			if err := imageTxRepo.ClearRoomImageByRoomImageRaku2ID(request.ImageID); err != nil {
				i.IRaku2Repository.TxRollback(tx)
				return err
			}
			for _, roomImage := range roomImages {
				// RoomTypeIDで紐づいてるデータ取得
				roomImagesByRoomTypeId, _ := imageTxRepo.FetchRoomImagesByRoomTypeID(roomImage.RoomTypeID)
				for index, value := range roomImagesByRoomTypeId {
					orderNum := uint8(index) + 1
					// orderを並び替えて更新（orderNumとvalue.Orderが一致しないものだけ更新）
					if orderNum != value.Order {
						input := &image.HtTmRoomOwnImagesRaku2s{
							RoomOwnImagesID: value.RoomOwnImagesID,
							Order:           orderNum,
						}
						if err := imageTxRepo.UpdateRoomImageOrder(input); err != nil {
							i.IRaku2Repository.TxRollback(tx)
							return err
						}
					}
				}
			}
		}
	}
	// コミットとロールバック
	if err := i.IRaku2Repository.TxCommit(tx); err != nil {
		i.IRaku2Repository.TxRollback(tx)
		return err
	}
	return nil
}

// UpdateIsMain メイン画像フラグ更新
func (i *imageRaku2Usecase) UpdateIsMain(request *image.UpdateIsMainInput) error {
	// ソート順更新
	if err := i.IRaku2Repository.UpdateIsMain(request); err != nil {
		return err
	}
	// 画像情報取得
	imageData, fetchErr := i.IRaku2Repository.FetchImageByRoomImageRaku2ID(request.ImageID)
	if fetchErr != nil {
		return fetchErr
	}
	// メイン画像から外す
	if !request.IsMain {
		// メイン画像取得
		mainImages, fetchImagesErr := i.IRaku2Repository.FetchMainImages(imageData.PropertyID)
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
					ImageID: mainImage.RoomImageRaku2ID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := i.IRaku2Repository.TrUpdateImageSort(&input); err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateSortNum メイン画像の並び順更新
func (i *imageRaku2Usecase) UpdateSortNum(request *[]image.UpdateSortNumInput) error {
	return i.IRaku2Repository.TrUpdateImageSort(request)
}

// Delete 画像削除
func (i *imageRaku2Usecase) Delete(imageID int64) error {
	imageData, fetchErr := i.IRaku2Repository.FetchImageByRoomImageRaku2ID(imageID)
	if fetchErr != nil {
		return fetchErr
	}
	gcsInfo := &image.GcsInfo{}
	if err := json.Unmarshal([]byte(imageData.GcsInfo), gcsInfo); err != nil {
		return err
	}

	// トランザクション生成
	tx, txErr := i.IRaku2Repository.TxStart()
	if txErr != nil {
		return txErr
	}

	imageTxRepo := iInfra.NewImageRaku2Repository(tx)

	if err := imageTxRepo.DeleteImage(imageID); err != nil {
		i.IRaku2Repository.TxRollback(tx)
		return err
	}

	// メイン画像に設定していたら外す
	if imageData.IsMain {
		// メイン画像取得
		mainImages, fetchImagesErr := imageTxRepo.FetchMainImages(imageData.PropertyID)
		if fetchImagesErr != nil {
			i.IRaku2Repository.TxRollback(tx)
			return fetchImagesErr
		}
		input := []image.UpdateSortNumInput{}
		for index, mainImage := range mainImages {
			sortNum := index + 1
			// sortNumとMainOrderが一致しないデータだけ更新（ソート順の並び替え）
			if sortNum != mainImage.MainOrder {
				input = append(input, image.UpdateSortNumInput{
					ImageID: mainImage.RoomImageRaku2ID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := imageTxRepo.UpdateImageSort(&input); err != nil {
				i.IRaku2Repository.TxRollback(tx)
				return err
			}
		}
	}

	// 対象データ取得
	roomImages, err := imageTxRepo.FetchRoomImagesByRoomImageRaku2ID(imageID)
	if err != nil {
		i.IRaku2Repository.TxRollback(tx)
		return err
	}
	if len(roomImages) > 0 {
		// 対象データ削除
		if err := imageTxRepo.ClearRoomImageByRoomImageRaku2ID(imageID); err != nil {
			i.IRaku2Repository.TxRollback(tx)
			return err
		}
		for _, roomImage := range roomImages {
			// RoomTypeIDで紐づいてるデータ取得
			roomImagesByRoomTypeId, _ := imageTxRepo.FetchRoomImagesByRoomTypeID(roomImage.RoomTypeID)
			for index, value := range roomImagesByRoomTypeId {
				orderNum := uint8(index) + 1
				// orderを並び替えて更新（orderNumとvalue.Orderが一致しないものだけ更新）
				if orderNum != value.Order {
					input := &image.HtTmRoomOwnImagesRaku2s{
						RoomOwnImagesID: value.RoomOwnImagesID,
						Order:           orderNum,
					}
					if err := imageTxRepo.UpdateRoomImageOrder(input); err != nil {
						i.IRaku2Repository.TxRollback(tx)
						return err
					}
				}
			}
		}
	}

	// コミットとロールバック
	if err := i.IRaku2Repository.TxCommit(tx); err != nil {
		i.IRaku2Repository.TxRollback(tx)
		return err
	}
	// GCSの画像を削除
	if err := i.ImageStorage.Delete(gcsInfo.Bucket, gcsInfo.Name); err != nil {
		return err
	}

	return nil
}

// Create 画像作成
func (i *imageRaku2Usecase) Create(request *image.UploadInput, file *multipart.FileHeader, hmUser account.HtTmHotelManager) error {
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

	upData := &image.HtTmImageRaku2s{
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
	return i.IRaku2Repository.CreateImage(upData)
}

// CountMainImages メイン画像の数を数える
func (i *imageRaku2Usecase) CountMainImages(propertyID int64) int64 {
	return i.IRaku2Repository.CountMainImagesPerPropertyID(propertyID)
}
