package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/room"
	"gorm.io/gorm"
)

// roomRepository 部屋関連repository
type roomCommonRepository struct {
	db *gorm.DB
}

// NewRoomCommonRepository インスタンス生成
func NewRoomCommonRepository(db *gorm.DB) room.IRoomCommonRepository {
	return &roomCommonRepository{
		db: db,
	}
}

// FetchAllRoomKinds 部屋の種類マスタを全件取得
func (r *roomCommonRepository) FetchAllRoomKinds() ([]room.HtTmRoomKind, error) {
	result := []room.HtTmRoomKind{}
	err := r.db.
		Model(&room.HtTmRoomKind{}).
		Where("is_delete = 0").
		Find(&result).Error
	return result, err
}
