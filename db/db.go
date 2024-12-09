package db

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"not null"`
	Surname string `gorm:"not null"`
	Method  string `gorm:"not null"`
	Token   string `gorm:"unique"`
}

var db *gorm.DB

func Init() {
	var err error
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)
	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	db.AutoMigrate(&User{})
}

func CreateUserSession(session User) (User, error) {
	err := db.Create(&session).Error
	return session, err
}

func GetUserSession(token string) (User, error) {
	var user User
	err := db.Where("token = ?", token).First(&user).Error
	return user, err
}

func DeleteUserSession(token string) error {
	err := db.Delete("token = ?", token).Error
	return err
}
