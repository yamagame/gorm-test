package main

import (
	"example/yamagame/gorm-test/model"
	"fmt"
	"os"
	"strconv"

	"github.com/oklog/ulid/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DB_USER = "root"
const DB_PASS = "pass"
const DB_HOST = "localhost"
const DB_PORT = "3316"
const DB_NAME = "gorm_test_db"
const DB_LOCA = "UTC"

type Result struct {
	ULID       string
	Name       string
	City       string
	SchoolName string
}

func DB() *gorm.DB {
	dsn := DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?charset=utf8mb4&parseTime=True&loc=" + DB_LOCA
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Addr{})
	db.AutoMigrate(&model.School{})
}

func create(db *gorm.DB, ulid, name, city, schoolname string) {
	db.Create(&model.User{Name: name, ULID: ulid})
	db.Create(&model.Addr{City: city, ULID: ulid})
	db.Create(&model.School{Name: schoolname, ULID: ulid})
}

func list(db *gorm.DB) []string {
	var ulids []string
	var users []*model.User
	db.Find(&users)
	for _, user := range users {
		ulids = append(ulids, user.ULID)
	}
	return ulids
}

func join1(db *gorm.DB, ulids []string) []*Result {
	users := []*Result{}
	for _, ulid := range ulids {
		var user model.User
		db.Select("ULID", "Name").Where("ul_id = ?", ulid).Take(&user)
		var addr model.Addr
		db.Select("ULID", "City").Where("ul_id = ?", ulid).Take(&addr)
		var school model.School
		db.Select("ULID", "Name").Where("ul_id = ?", ulid).Take(&school)
		users = append(users, &Result{
			Name:       user.Name,
			City:       addr.City,
			SchoolName: school.Name,
		})
	}
	return users
}

func join2(db *gorm.DB, ulids []string) []*Result {
	users := []*Result{}
	for _, ulid := range ulids {
		var user Result
		db.Table("users").
			Select("users.ul_id as ulid", "users.name", "addrs.city", "schools.name as school_name").
			Where("users.ul_id = ?", ulid).
			Joins("left join addrs on addrs.ul_id = users.ul_id").
			Joins("left join schools on schools.ul_id = users.ul_id").
			Take(&user)
		users = append(users, &user)
	}
	return users
}

func main() {
	db := DB()

	cmd := ""
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	option := ""
	if len(os.Args) > 2 {
		option = os.Args[2]
	}

	switch cmd {
	case "migrate":
		migrate(db)
	case "create":
		if option == "" {
			option = "100"
		}
		count, _ := strconv.Atoi(option)
		for i := 0; i < count; i++ {
			id := ulid.Make()
			idx := fmt.Sprintf("-%d", i)
			create(db, id.String(), "name"+idx, "city"+idx, "school"+idx)
		}
	case "list":
		ulids := list(db)
		if option != "" {
			for _, ulid := range ulids {
				fmt.Println(ulid)
			}
		}
	case "join1":
		ulids := list(db)
		users := join1(db, ulids)
		if option != "" {
			for _, user := range users {
				fmt.Println(user)
			}
		}
	case "join2":
		ulids := list(db)
		users := join2(db, ulids)
		if option != "" {
			for _, user := range users {
				fmt.Println(user)
			}
		}
	}
}
