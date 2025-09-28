package postgresAdp

import (
	"github.com/vanshjangir/rapidgo/server/internal/domain"
	"gorm.io/gorm"
)

type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) FindByUsername(username string) (domain.User, error) {
	var user UserModel
	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return domain.User{}, err
	}
	return domain.User{
        Username: user.Username,
        Email: user.Email,
        Password: user.Password,
        Rating: user.Rating,
        Highestrating: user.Highestrating,
    }, nil
}

func (r *PostgresUserRepo) FindByEmail(email string) (domain.User, error) {
	var user UserModel
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return domain.User{}, err
	}
	return domain.User{
        Username: user.Username,
        Email: user.Email,
        Password: user.Password,
        Rating: user.Rating,
        Highestrating: user.Highestrating,
    }, nil
}

func (r *PostgresUserRepo) InsertUser(user domain.User) error {
    tx := r.db.Begin()
    err := tx.Create(&UserModel{
        Username: user.Username,
        Email: user.Email,
        Password: user.Password,
        Rating: user.Rating,
        Highestrating: user.Highestrating,
    }).Error
    if err != nil {
        tx.Rollback()
        return err
    }
    tx.Commit()
    return nil
}

func (r *PostgresUserRepo) ChangeUsername(oldName, newName string) error {
	tx := r.db.Begin()
	if err := tx.Model(&UserModel{}).Where("username = ?", oldName).
		Update("username", newName).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&GameModel{}).Where("white = ?", oldName).
		Update("white", newName).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&GameModel{}).Where("black = ?", oldName).
		Update("black", newName).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
