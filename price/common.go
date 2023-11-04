package price

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"time"
)

// PriceTable 料金テーブル
type PriceTable struct {
	PriceID          int64     `gorm:"primaryKey;autoIncrement:true" json:"price_id,omitempty"`
	PlanID           int64     `json:"plan_id,omitempty" validate:"required"`
	UseDate          time.Time `gorm:"type:time" json:"use_date"`
	RateTypeCode     string    `json:"rate_type_code,omitempty"`
	Price            int       `json:"price" validate:"required"`
	PriceInTax       int       `json:"price_in_tax" validate:"required"`
	ChildPrice1      int       `json:"child_price1" validate:"required" gorm:"column:child_price1"`
	ChildPrice1InTax int       `json:"child_price1_in_tax" validate:"required" gorm:"column:child_price1_in_tax"`
	ChildPrice2      int       `json:"child_price2" validate:"required" gorm:"column:child_price2"`
	ChildPrice2InTax int       `json:"child_price2_in_tax" validate:"required" gorm:"column:child_price2_in_tax"`
	ChildPrice3      int       `json:"child_price3" validate:"required" gorm:"column:child_price3"`
	ChildPrice3InTax int       `json:"child_price3_in_tax" validate:"required" gorm:"column:child_price3_in_tax"`
	ChildPrice4      int       `json:"child_price4" validate:"required" gorm:"column:child_price4"`
	ChildPrice4InTax int       `json:"child_price4_in_tax" validate:"required" gorm:"column:child_price4_in_tax"`
	ChildPrice5      int       `json:"child_price5" validate:"required" gorm:"column:child_price5"`
	ChildPrice5InTax int       `json:"child_price5_in_tax" validate:"required" gorm:"column:child_price5_in_tax"`
	ChildPrice6      int       `json:"child_price6" validate:"required" gorm:"column:child_price6"`
	ChildPrice6InTax int       `json:"child_price6_in_tax" validate:"required" gorm:"column:child_price6_in_tax"`
	RegularPrice     int       `json:"regular_price" validate:"required"`
	common.Times     `gorm:"embedded"`
}

// ChildRateTable 子供料金設定テーブル
// child_rate_type 1 ... 9〜11歳　小学校高学年(ChildA)
//
//	2 ... 6〜8歳　小学校低学年(ChildB)
//	3 ... 0〜5歳　幼児（食事・布団あり）(ChildC)
//	4 ... 0〜5歳　幼児（食事あり・布団なし）(ChildD)
//	5 ... 0〜5歳　幼児（食事なし・布団あり）(ChildE)
//	6 ... 0〜5歳　乳児（食事なし・布団なし）(ChildF)
type ChildRateTable struct {
	ChildRateID   int64 `gorm:"primaryKey;autoIncrement:true" json:"child_rate_id,omitempty"`
	ChildRateType int8  `json:"child_rate_type"`
	PlanID        int64 `json:"plan_id"`
	FromAge       int8  `json:"from_age"`
	ToAge         int8  `json:"to_age"`
	Receive       bool  `json:"receive"`       // 受入
	RateCategory  int8  `json:"rate_category"` // 0: ％, 1: 定額(円), 2: 定額を引く(円引き)
	Rate          int   `json:"rate"`
	CalcCategory  bool  `json:"calc_category"` // 大人人数算出時に数えるかどうか
	common.Times  `gorm:"embedded"`
}

// DetailInput 料金詳細の入力
type DetailInput struct {
	PlanID    int64  `json:"plan_id" param:"planId" validate:"required"`
	StartDate string `json:"start_date" param:"startDate" validate:"required"`
	EndDate   string `json:"end_date" param:"endDate" validate:"required"`
}

// DetailOutput 料金詳細の出力
type DetailOutput struct {
	PlanID int64              `json:"plan_id"`
	Prices map[string][]Price `json:"prices"`
}

// SaveInput データ登録・更新の入力
type SaveInput struct {
	PlanID int64              `json:"plan_id" validate:"required"`
	Prices map[string][]Price `json:"prices"`
}

// Price 料金データ
type Price struct {
	Type        string `json:"type"`
	Price       int    `json:"price"`
	IsStopSales bool   `json:"is_stop_sales"`
}

// IPriceUsecase 料金関連のusecaseのインターフェース
type IPriceUsecase interface {
	FetchDetail(request *DetailInput) (DetailOutput, error)
	Save(request *[]SaveInput) error
}
