package usecase

import (
	"testing"
	"time"

	"github.com/Adventureinc/hotel-hm-api/src/room"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestFetchListNeppan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := []string{"room_type_id", "property_id", "room_type_code", "name", "room_desc", "stock_setting_start", "stock_setting_end", "is_setting_stock_year_round", "room_count", "ocu_min", "ocu_max", "is_smoking", "is_stop_sales", "is_delete", "created_at", "updated_at"}
	mock.ExpectQuery("SELECT (.+) FROM ht_tm_room_type_directs AS room WHERE property_id = (.+) AND is_delete = 0").
		WithArgs(1208039).
		WillReturnRows(
			sqlmock.NewRows(columns).
				AddRow(89, 1208039, "R001", "Room1", "", stringToTime("2001-01-01 00:00:00"), stringToTime("2001-01-01 00:00:00"), true, 1, 1, 3, false, false, false, stringToTime("2021-04-23 17:40:36"), stringToTime("2021-05-07 17:28:03")).
				AddRow(90, 1208039, "3456", "Royal", "", stringToTime("2001-01-01 00:00:00"), stringToTime("2001-01-01 00:00:00"), true, 1, 5, 6, false, false, false, stringToTime("2021-04-23 17:49:47"), stringToTime("2021-04-23 17:49:47")).
				AddRow(100, 1208039, "Room2", "Room2", "", stringToTime("2001-01-01 00:00:00"), stringToTime("2001-01-01 00:00:00"), true, 1, 1, 1, false, false, false, stringToTime("2021-05-11 12:18:54"), stringToTime("2021-05-11 12:18:54")),
		)

	columns = []string{"image_id", "room_type_id", "href", "caption", "order"}
	mock.ExpectQuery("SELECT image.room_image_direct_id as image_id, bind.room_type_id, image.href, image.caption, bind.order FROM ht_tm_image_directs as image INNER JOIN ht_tm_room_own_images_directs AS bind ON image.room_image_direct_id = bind.room_image_direct_id WHERE bind.room_type_id IN (.+) ORDER BY bind.room_type_id, bind.order").
		WithArgs(89, 90, 100).
		WillReturnRows(
			sqlmock.NewRows(columns).AddRow("85", "89", "https://storage.googleapis.com/dev-hotel/hotel/property-images/7/1208039_20210507172736.png", "", "1"),
		)

	columns = []string{"image_id", "room_type_id", "href", "caption", "order"}
	mock.ExpectQuery("SELECT bind.room_type_id, bind.direct_room_amenity_id, amenity.direct_room_amenity_name FROM ht_tm_room_use_amenity_directs AS bind INNER JOIN ht_tm_room_amenity_directs AS amenity ON bind.direct_room_amenity_id = amenity.direct_room_amenity_id WHERE bind.room_type_id IN (.+)").
		WithArgs(89, 90, 100).
		WillReturnRows(
			sqlmock.NewRows(columns),
		)

	gorm, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("initializing err %s", err)
	}
	uc := NewRoomDirectUsecase(gorm.Debug())

	request := &room.ListInput{PropertyID: 1208039}
	l, err := uc.FetchList(request)
	if err != nil {
		t.Fatalf("err is %s", err)
	}

	if len(l) != 3 {
		t.Fatalf("RoomDirectUsecase.FetchListの個数が正しくありません。。Expected: %d, Actual: %d", 1, len(l))
	}
}

func stringToTime(str string) time.Time {
	var layout = "2006-01-02 15:04:05"
	t, _ := time.Parse(layout, str)
	return t
}
