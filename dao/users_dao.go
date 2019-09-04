package dao

import (
	"genosha/db"
	"genosha/models"
	"github.com/jinzhu/gorm"
)

var localPG = db.Pg



func CreateNewUser(u *models.User) error {
	return localPG.Create(u).Error
}

func UpdateUserByGuid(guid, name string) error {
	var user models.User
	return localPG.Model(&user).Where("guid = ?", guid).Updates(models.User{Name:name}).Error
}

func GetUserByGuid(userGuid string) *models.User {
	var user models.User
	localPG.Where("guid = ?", userGuid).First(&user)
	return &user
}

func GetUserByLowEmail(userLowEmail string) *models.User {
	var user models.User
	localPG.Where("low_email = ?", userLowEmail).First(&user)
	return &user
}

func UpdatePassWDByEmail(userEmail string, PassWD string) error {
	var user models.User
	return localPG.Model(&user).Where("email = ?", userEmail).Update("password_hash", PassWD).Error
}

func UpdateTokenByLowEmail(userLowEmail string, token string) error {
	var user models.User
	return localPG.Model(&user).Where("low_email = ?", userLowEmail).Update("token", token).Error
}

func ResetTokenCountByLowEmail(userLowEmail string) error {
	var user models.User
	return localPG.Model(&user).Where("low_email = ?", userLowEmail).Update("token_count", 0).Error
}

func IncreaseTokenCountByLowEmail(userLowEmail string) error {
	var user models.User
	return localPG.Model(&user).Where("low_email = ?", userLowEmail).UpdateColumn("token_count", gorm.Expr(`token_count + ?`, 1)).Error
}

func UpdateVersionByEmail(userEmail string, Version string) error {
	var user models.User
	return localPG.Model(&user).Where("email = ?", userEmail).Update("version", Version).Error
}

func CreateUserToRole(r *models.UserToRole) error {
	return localPG.Create(r).Error
}

func GetUserRolesByGuid(userGuid string) []string {
	var roles []string
	localPG.Table("user_to_roles").Where("user_guid = ?", userGuid).Order("role").Pluck("role", &roles)
	return roles
}
