package repositories

import (
	"context"
	db "project-sqlc/internal/db"
	models "project-sqlc/internal/db/models"
	"project-sqlc/internal/dto"

	"github.com/jackc/pgx/v5"
)

type IUserRepository interface {
	GetUser(ctx context.Context, id int32) (models.User, error)
	GetUsers(ctx context.Context, baseDto dto.GetUsersDto) ([]models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type userRepository struct {
	db      *db.Database
	queries *models.Queries
}

func NewUserRepository(db *db.Database) IUserRepository {
	return &userRepository{db: db, queries: db.Query}
}

func (r *userRepository) GetUser(ctx context.Context, id int32) (models.User, error) {
	user, err := r.queries.GetUserById(ctx, int32(id))
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUsers(ctx context.Context, baseDto dto.GetUsersDto) ([]models.User, error) {
	users, err := r.queries.GetUsers(ctx, models.GetUsersParams{
		Name:   "%" + baseDto.Name + "%",
		Limit:  baseDto.GetTake(),
		Offset: baseDto.GetSkip(),
	})
	if err != nil {
		return []models.User{}, err
	}
	return users, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	id, err := r.queries.CreateUser(ctx, models.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return models.User{}, err
	}
	return r.GetUser(ctx, int32(id))
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	// user, err := r.queries.GetUserByEmail(ctx, email)
	rows, err := r.db.Pool.Query(ctx, "SELECT * FROM users WHERE email = $1", "lhquan1999@gmail.com")
	if err != nil {
		return models.User{}, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[models.User])
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
