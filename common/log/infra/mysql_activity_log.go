package infra

import (
	"github.com/Adventureinc/hotel-hm-api/src/common/log"
	"gorm.io/gorm"
	"time"
)

type logRepository struct {
	db *gorm.DB
}

func (l *logRepository) TxStart() (*gorm.DB, error) {
	tx := l.db.Begin()
	return tx, tx.Error
}

func (l *logRepository) TxCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (l *logRepository) TxRollback(tx *gorm.DB) {
	tx.Rollback()
}

func NewLogRepository(db *gorm.DB) log.ILogRepository {
	return &logRepository{
		db: db,
	}
}

func (l *logRepository) StoreBulkActivityLog(ServiceName string, Type string, Host string, start time.Time) (int64, error) {
	newLog := &log.HtThHmBulkActivityLog{
		ServiceName:    ServiceName,
		Type:           Type,
		ProcessStartAt: start,
		HostUrl:        Host,
		CreatedAt:      time.Now(),
	}

	err := l.db.Create(newLog).Error
	if err != nil {
		return 0, err
	}

	return newLog.ActivityLogID, nil
}

func (l *logRepository) UpdateBulkActivityLog(ActivityLogID int64, ProcessStartTime time.Time, status bool, errorMessage string) error {

	end := time.Now()
	duration := end.Sub(ProcessStartTime)
	seconds := int(duration.Seconds())

	return l.db.Model(&log.HtThHmBulkActivityLog{}).
		Where("hm_bulk_activity_log_id = ?", ActivityLogID).
		Updates(map[string]interface{}{
			"process_end_at": time.Now(),
			"duration":       seconds,
			"is_success":     status,
			"error_message":  errorMessage,
			"updated_at":     time.Now(),
		}).Error
}
