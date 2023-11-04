package infra

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	"github.com/Adventureinc/hotel-hm-api/src/image"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// imageStorage 画像関連storage
type imageStorage struct{}

// NewImageStorage インスタンス生成
func NewImageStorage() image.IImageStorage {
	return &imageStorage{}
}

// Delete GCSの画像を削除
func (i *imageStorage) Delete(bucketName string, objectPath string) error {
	ctx := context.Background()
	// client 生成
	cfg, cfgErr := google.JWTConfigFromJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")), storage.ScopeReadWrite)
	if cfgErr != nil {
		return cfgErr
	}

	client, sErr := storage.NewClient(ctx, option.WithTokenSource(cfg.TokenSource(ctx)))
	if sErr != nil {
		return sErr
	}

	if err := client.Bucket(bucketName).Object(objectPath).Delete(ctx); err != nil {
		return err
	}
	return nil
}

// Create GCSに画像をアップロード
func (i *imageStorage) Create(bucketName string, filename string, file *multipart.FileHeader) (*storage.ObjectAttrs, error) {
	attrs := &storage.ObjectAttrs{}
	ctx := context.Background()
	// client 生成
	cfg, cfgErr := google.JWTConfigFromJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")), storage.ScopeReadWrite)
	if cfgErr != nil {
		return attrs, cfgErr
	}

	client, sErr := storage.NewClient(ctx, option.WithTokenSource(cfg.TokenSource(ctx)))
	if sErr != nil {
		return attrs, sErr
	}

	// file準備
	src, err := file.Open()
	if err != nil {
		return attrs, err
	}
	defer src.Close()

	// 書き込み
	sw := client.Bucket(bucketName).Object(filename).NewWriter(ctx)
	if _, ioErr := io.Copy(sw, src); ioErr != nil {
		return attrs, ioErr
	}
	if swErr := sw.Close(); swErr != nil {
		return attrs, swErr
	}

	return client.Bucket(bucketName).Object(filename).Attrs(ctx)
}
