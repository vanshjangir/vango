package postgres_adp

import (
	"fmt"

	"github.com/vanshjangir/rapidgo/server/internal/domain"
	"gorm.io/gorm"
)

type PostgresGameRepo struct {
	db *gorm.DB
}

func NewPostgresGameRepo(db *gorm.DB) *PostgresGameRepo {
	return &PostgresGameRepo{db: db}
}

func (r *PostgresGameRepo) SaveGame(g *domain.Game) error {
	tx := r.db.Begin()
	
	err := tx.Model(&GameModel{}).
		Where("gameid = ?", g.Id).
		Updates(map[string]any{
			"winner": g.Winner,
			"wonby": g.WonBy,
			"moves": g.History,
		}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	
	tx.Commit()
	return nil
}

func (r *PostgresGameRepo) Review(gameId int) (domain.GameReview, error) {
	var review GameModel
	if err := r.db.Take(&review, "gameid = ?", gameId).Error; err != nil {
		return domain.GameReview{}, err
	}
	
	return domain.GameReview{
		Id:        review.Gameid,
		BlackName: review.Black,
		WhiteName: review.White,
		Winner:    review.Winner,
		Moves:     review.Moves,
	}, nil
}

func (r *PostgresGameRepo) CreateNewGame(
	blackName, whiteName string,
) (int, error) {
    tx := r.db.Begin()
	game := GameModel{White: whiteName, Black: blackName}
    result := tx.Create(&game)
    if result.Error != nil {
        tx.Rollback()
		return -1, fmt.Errorf("CreateNewGame: %v", result.Error)
    }
    tx.Commit()
    return  game.Gameid, nil
}
