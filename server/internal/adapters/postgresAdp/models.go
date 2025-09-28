package postgresAdp

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Username      string `gorm:"primaryKey"`
	Password      string
	Email         string
	Rating        int
	Highestrating int
}

type GameModel struct {
	gorm.Model
	Gameid    string `gorm:"primaryKey"`
	White     string
	Black     string
	Winner    string
	Wonby     string
	Moves     []string
	Date      time.Time
	CreatedAt time.Time
}
