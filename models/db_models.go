package models

import (
	"genosha/db"
	"github.com/jinzhu/gorm"
	"log"
)

type (
	User struct {
		gorm.Model
		Guid         string `gorm:"type:varchar(1024);not null;primary_key;unique"`
		Name         string `gorm:"not null"`
		Email        string `gorm:"not null;unique"`
		LowEmail     string `gorm:"not null;unique"`
		PasswordHash string `gorm:"type:varchar(1024);not null"`
		Version      string `gorm:"not null;default:''"`
		Token        string `gorm:"not null;default:''"`
		TokenCount   int    `gorm:"not null;default:0"`
	}
	Role struct {
		gorm.Model
		Name string `gorm:"type:varchar(32);unique;not null"`
	}
	UserToRole struct {
		gorm.Model
		UserGuid string `gorm:"type:varchar(32);not null"`
		Role     string `gorm:"type:varchar(32);not null"`
	}
)

func init() {
	log.Println("coming to models")
	db.Pg.DB().SetMaxIdleConns(5)
	db.Pg.DB().SetMaxOpenConns(15)
	db.Pg.AutoMigrate(&User{})
	db.Pg.AutoMigrate(&Role{})
	db.Pg.AutoMigrate(&UserToRole{})
	db.Pg.Model(&UserToRole{}).AddForeignKey("role", "roles(name)", "CASCADE", "CASCADE")
}
