package log

import (
	"github.com/Adventureinc/hotel-hm-api/src/common"
	"time"
)

type HtThHmBulkActivityLog struct {
	ActivityLogID  int64     `gorm:"primaryKey;autoIncrement:true" json:"hm_bulk_activity_log_id"`
	ServiceName    string    `json:"service_name"`
	Type           string    `json:"type"`
	ProcessStartAt time.Time `gorm:"type:time" json:"process_start_at"`
	ProcessEndAt   time.Time `gorm:"type:time" json:"process_end_at"`
	Duration       int64     `json:"duration"`
	IsSuccess      bool      `json:"is_success"`
	ErrorMessage   string    `json:"error_message"`
	HostUrl        string    `json:"host_url"`
	CreatedAt      time.Time `gorm:"type:time" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:time" json:"updated_at"`
}

// ILogRepository represents a repository for logging information
type ILogRepository interface {
	common.Repository
	StoreBulkActivityLog(ServiceName string, Type string, Host string, start time.Time) (int64, error)
	UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, status bool, errorMessage string) error
}
