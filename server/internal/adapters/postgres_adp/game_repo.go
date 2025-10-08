package postgres_adp

import (
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
	return nil
}

func (r *PostgresGameRepo) Review(gameId string) (domain.GameReview, error) {
	return domain.GameReview{}, nil
}

func (r *PostgresGameRepo) CreateNewGame(blackName, whiteName string) error {
	return nil
}
