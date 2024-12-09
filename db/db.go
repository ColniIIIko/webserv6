package db

import (
	"os"

	"gorm.io/driver/postgres"
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
		connectionString = os.Getenv("DB_CONN")
	)

	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
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
