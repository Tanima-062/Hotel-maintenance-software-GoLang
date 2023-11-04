package common

import (
	"time"

	"gorm.io/gorm"
)

// Times テーブル共通情報（作成日と更新日）
type Times struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Paging ページング処理用
type Paging struct {
	Limit  int `json:"limit" query:"limit"`   // 件数
	Offset int `json:"offset" query:"offset"` // ページ番号
}

// Repository リポジトリ共通
type Repository interface {
	// TxStart トランザクションスタート
	TxStart() (*gorm.DB, error)
	// TxCommit トランザクションコミット
	TxCommit(tx *gorm.DB) error
	// TxRollback トランザクションロールバック
	TxRollback(tx *gorm.DB)
}
