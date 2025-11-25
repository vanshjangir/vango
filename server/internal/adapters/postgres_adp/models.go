package postgres_adp

import (
	"time"

	"github.com/vanshjangir/vango/server/internal/domain"
)

type UserModel struct {
	Username      string `gorm:"primaryKey"`
	Password      string
	Email         string
	Rating        int
	Highestrating int
}

type GameModel struct {
	Gameid    int `gorm:"primaryKey;autoIncrement"`
	White     string
	Black     string
	Winner    int
	Wonby     string
	Moves     domain.StringArray
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

func (GameModel) TableName() string {
	return "games"
}
