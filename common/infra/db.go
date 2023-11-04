package infra

import (
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBCon DB読み込み処理
func DBCon() (*gorm.DB, error) {

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	protocol := fmt.Sprintf("tcp(%s:%s)", host, port)

	connectInfo := fmt.Sprintf("%s:%s@%s/%s?parseTime=True&loc=%s", username, password, protocol, database, url.PathEscape("Asia/Tokyo"))
	DB, err := gorm.Open(mysql.Open(connectInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DB, nil
}

// CommonDBCon CommonDB読み込み処理
func CommonDBCon() (*gorm.DB, error) {

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := "common" // hard cording
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	protocol := fmt.Sprintf("tcp(%s:%s)", host, port)

	connectInfo := fmt.Sprintf("%s:%s@%s/%s?parseTime=True&loc=%s", username, password, protocol, database, url.PathEscape("Asia/Tokyo"))
	return gorm.Open(mysql.Open(connectInfo), &gorm.Config{})
}
