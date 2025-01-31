package users

import "gorm.io/gorm"

type User struct {
	gorm.Model
	SteamID   string `gorm:"unique:not null"`
	Username  string
	AvatarURL string
	SteamLVL  int
}
