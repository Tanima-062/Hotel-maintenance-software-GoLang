package infra

import (
	"bufio"
	"context"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/Adventureinc/hotel-hm-api/src/settlement"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const tempWritingFile = "tempWritingFile.pdf"

// settlementStorage 請求関連storage
type settlementStorage struct{}

// NewSettlementStorage インスタンス生成
func NewSettlementStorage() settlement.ISettlementStorage {
	return &settlementStorage{}
}

// Get 請求書をストレージから取得
func (i *settlementStorage) Get(bucketName string, objectPath string) (string, error) {
	ctx := context.Background()
	// client 生成
	cfg, cfgErr := google.JWTConfigFromJSON([]byte(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS_JSON")), storage.ScopeReadWrite)
	if cfgErr != nil {
		return "", cfgErr
	}
	client, sErr := storage.NewClient(ctx, option.WithTokenSource(cfg.TokenSource(ctx)))
	if sErr != nil {
		return "", sErr
	}

	rc, oErr := client.Bucket(bucketName).Object(objectPath).NewReader(ctx)
	defer rc.Close()
	if oErr != nil {
		return "", oErr
	}

	// GCSオブジェクトを書き込むファイルの作成
	f, cErr := os.Create(tempWritingFile)
	if cErr != nil {
		return "", cErr
	}
	// 書き込み
	tee := io.TeeReader(rc, f)
	reader := bufio.NewReader(tee)

	for {
		line, readErr := reader.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return "", readErr
		}
		allLinesProcessed := readErr == io.EOF && len(line) == 0
		if allLinesProcessed {
			break
		}
	}
	return tempWritingFile, nil
}
