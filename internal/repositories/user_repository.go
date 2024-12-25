package repositories

import (
	"context"
	db "project-sqlc/internal/db/models"
	"project-sqlc/internal/dto"
)

type IUserRepository interface {
	GetUser(ctx context.Context, id int32) (db.User, error)
	GetUsers(ctx context.Context, baseDto dto.GetUsersDto) ([]db.User, error)
	CreateUser(ctx context.Context, user db.User) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
}

type userRepository struct {
	db *db.Queries
}

func NewUserRepository(db *db.Queries) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUser(ctx context.Context, id int32) (db.User, error) {
	user, err := r.db.GetUserById(ctx, int32(id))
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUsers(ctx context.Context, baseDto dto.GetUsersDto) ([]db.User, error) {
	users, err := r.db.GetUsers(ctx, db.GetUsersParams{
		Name:   "%" + baseDto.Name + "%",
		Limit:  baseDto.GetTake(),
		Offset: baseDto.GetSkip(),
	})
	if err != nil {
		return []db.User{}, err
	}
	return users, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user db.User) (db.User, error) {
	id, err := r.db.CreateUser(ctx, db.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return db.User{}, err
	}
	return r.GetUser(ctx, int32(id))
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.db.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}
	return user, nil
}
