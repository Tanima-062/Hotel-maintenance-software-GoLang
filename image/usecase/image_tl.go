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

// imageTlUsecase TL画像関連usecase
type imageTlUsecase struct {
	ARepository   account.IAccountRepository
	ITlRepository image.IImageTlRepository
	ImageStorage  image.IImageStorage
}

// NewImageTlUsecase インスタンス生成
func NewImageTlUsecase(db *gorm.DB) image.IImageUsecase {
	return &imageTlUsecase{
		ARepository:   aInfra.NewAccountRepository(db),
		ITlRepository: iInfra.NewImageTlRepository(db),
		ImageStorage:  iInfra.NewImageStorage(),
	}
}

// FetchAll 画像一覧
func (i *imageTlUsecase) FetchAll(request *image.ListInput) ([]image.ImagesOutput, error) {
	response := []image.ImagesOutput{}
	tlImages, err := i.ITlRepository.FetchImagesByPropertyID(request.PropertyID)
	if err != nil {
		return response, nil
	}
	for _, v := range tlImages {
		response = append(response, image.ImagesOutput{
			ImageID:    v.ImageTlID,
			PropertyID: v.PropertyID,
			Method:     v.Method,
			Href:       v.Href,
			CategoryCd: v.CategoryCd,
			IsMain:     v.HeroImage,
			Caption:    v.Caption,
			SortNum:    v.SortNum,
			GcsInfo:    v.GcsInfo,
		})
	}
	return response, nil
}

// Update 画像情報更新
func (i *imageTlUsecase) Update(request *image.UpdateInput) error {
	// 更新前の対象データ取得
	imageData, err := i.ITlRepository.FetchImageByImageTlID(request.ImageID)
	if err != nil {
		return err
	}

	// トランザクション生成
	tx, txErr := i.ITlRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	imageTxRepo := iInfra.NewImageTlRepository(tx)

	// 更新
	if err := imageTxRepo.UpdateImageTl(request); err != nil {
		i.ITlRepository.TxRollback(tx)
		return err
	}

	// メイン画像から外す
	if !request.IsMain && imageData.HeroImage {
		// メイン画像取得
		mainImages, fetchImagesErr := imageTxRepo.FetchMainImages(imageData.PropertyID)
		if fetchImagesErr != nil {
			i.ITlRepository.TxRollback(tx)
			return fetchImagesErr
		}
		input := []image.UpdateSortNumInput{}
		for index, mainImage := range mainImages {
			sortNum := index + 1
			// sortNumとmainImage.SortNumが一致しないデータだけ更新（ソート順の並び替え）
			if sortNum != mainImage.SortNum {
				input = append(input, image.UpdateSortNumInput{
					ImageID: mainImage.ImageTlID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := imageTxRepo.UpdateImageSort(&input); err != nil {
				i.ITlRepository.TxRollback(tx)
				return err
			}
		}
	}
	// コミットとロールバック
	if err := i.ITlRepository.TxCommit(tx); err != nil {
		i.ITlRepository.TxRollback(tx)
		return err
	}
	return nil
}

// UpdateIsMain メイン画像フラグ更新
func (i *imageTlUsecase) UpdateIsMain(request *image.UpdateIsMainInput) error {
	// ソート順更新
	if err := i.ITlRepository.UpdateIsMain(request); err != nil {
		return err
	}
	// 画像情報取得
	imageData, fetchErr := i.ITlRepository.FetchImageByImageTlID(request.ImageID)
	if fetchErr != nil {
		return fetchErr
	}
	// メイン画像から外す
	if !request.IsMain {
		// メイン画像取得
		mainImages, fetchImagesErr := i.ITlRepository.FetchMainImages(imageData.PropertyID)
		if fetchImagesErr != nil {
			return fetchImagesErr
		}

		input := []image.UpdateSortNumInput{}
		for index, mainImage := range mainImages {
			// ソート順設定
			sortNum := index + 1
			// 取得した画像と一致しなければsortNumの値を設定
			if mainImage.SortNum != sortNum {
				input = append(input, image.UpdateSortNumInput{
					ImageID: mainImage.ImageTlID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := i.ITlRepository.TrUpdateImageSort(&input); err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateSortNum メイン画像の並び順更新
func (i *imageTlUsecase) UpdateSortNum(request *[]image.UpdateSortNumInput) error {
	return i.ITlRepository.TrUpdateImageSort(request)
}

// Delete 画像削除
func (i *imageTlUsecase) Delete(imageID int64) error {
	gcsInfo := &image.GcsInfo{}
	imageData, fetchErr := i.ITlRepository.FetchImageByImageTlID(imageID)
	if fetchErr != nil {
		return fetchErr
	}

	if err := json.Unmarshal([]byte(imageData.GcsInfo), gcsInfo); err != nil {
		return err
	}

	// トランザクション生成
	tx, txErr := i.ITlRepository.TxStart()
	if txErr != nil {
		return txErr
	}

	imageTxRepo := iInfra.NewImageTlRepository(tx)

	if err := imageTxRepo.DeleteImage(imageID); err != nil {
		i.ITlRepository.TxRollback(tx)
		return err
	}

	// メイン画像に設定していたら外す
	if imageData.HeroImage {
		// メイン画像取得
		mainImages, fetchImagesErr := imageTxRepo.FetchMainImages(imageData.PropertyID)
		if fetchImagesErr != nil {
			i.ITlRepository.TxRollback(tx)
			return fetchImagesErr
		}
		input := []image.UpdateSortNumInput{}
		for index, mainImage := range mainImages {
			sortNum := index + 1
			// sortNumとmainImage.SortNumが一致しないデータだけ更新（ソート順の並び替え）
			if sortNum != mainImage.SortNum {
				input = append(input, image.UpdateSortNumInput{
					ImageID: mainImage.ImageTlID,
					SortNum: sortNum,
				})
			}
		}
		// 更新するデータがあればソート順更新
		if len(input) > 0 {
			if err := imageTxRepo.UpdateImageSort(&input); err != nil {
				i.ITlRepository.TxRollback(tx)
				return err
			}
		}
	}

	// コミットとロールバック
	if err := i.ITlRepository.TxCommit(tx); err != nil {
		i.ITlRepository.TxRollback(tx)
		return err
	}
	// GCSの画像を削除
	if err := i.ImageStorage.Delete(gcsInfo.Bucket, gcsInfo.Name); err != nil {
		return err
	}

	return nil
}

// Create 画像作成
func (i *imageTlUsecase) Create(request *image.UploadInput, file *multipart.FileHeader, hmUser account.HtTmHotelManager) error {
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

	upData := &image.HtTmImageTls{
		PropertyID: hmUser.PropertyID,
		Method:     "GET",
		Href:       "https://storage.googleapis.com/" + bucketName + "/" + fileAllPath,
		Links:      "",
		CategoryCd: request.CategoryCd,
		HeroImage:  request.IsMain,
		Caption:    request.Caption,
		SortNum:    request.SortNum,
		GcsInfo:    string(gcsInfoBytes),
	}
	upData.CreatedAt = time.Now()
	upData.UpdatedAt = time.Now()
	return i.ITlRepository.CreateImage(upData)
}

// CountMainImages メイン画像の数を数える
func (i *imageTlUsecase) CountMainImages(propertyID int64) int64 {
	return i.ITlRepository.CountMainImagesPerPropertyID(propertyID)
}
