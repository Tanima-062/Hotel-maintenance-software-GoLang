package usecase

import (
	"github.com/Adventureinc/hotel-hm-api/src/room"
	rInfra "github.com/Adventureinc/hotel-hm-api/src/room/infra"
	"gorm.io/gorm"
)

// roomCommonUsecase 直仕入れ部屋関連usecase
type roomCommonUsecase struct {
	RComonRepository room.IRoomCommonRepository
}

// NewRoomDirectUsecase インスタンス生成
func NewRoomCommonUsecase(db *gorm.DB) room.IRoomCommonUsecase {
	return &roomCommonUsecase{
		RComonRepository: rInfra.NewRoomCommonRepository(db),
	}
}

// FetchAllRoomKinds 部屋種別一覧取得
func (r *roomCommonUsecase) FetchAllRoomKinds() ([]room.HtTmRoomKind, error) {
	roomKinds, roomKindsErr := r.RComonRepository.FetchAllRoomKinds()
	if roomKindsErr != nil {
		return []room.HtTmRoomKind{}, roomKindsErr
	}
	return roomKinds, nil
}
