package image

import "github.com/Adventureinc/hotel-hm-api/src/common"

// HtTmPlanOwnImagesTemas
type HtTmPlanOwnImagesTemas struct {
	PlanOwnImagesID int64 `gorm:"primaryKey;autoIncrement:true" json:"plan_own_images_id,omitempty"`
	PlanID          int64 `json:"plan_id,omitempty"`
	PlanImageTemaID int64 `json:"plan_image_tema_id,omitempty"`
	Order           uint8 `json:"order,omitempty"`
	common.Times
}

type IImageTemaRepository interface {
	// CreatePlanOwnImagesTema
	CreatePlanOwnImagesTema(images []HtTmPlanOwnImagesTemas) error
	FetchRoomImagesByRoomTypeID(roomTypeIDList []int64) ([]RoomImagesOutput, error)
	FetchImagesByPlanID(planIDList []int64) ([]PlanImagesOutput, error)
}
