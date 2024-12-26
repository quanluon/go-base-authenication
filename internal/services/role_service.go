package services

import (
	"context"
	db "project-sqlc/internal/db"
	models "project-sqlc/internal/db/models"
)

type IRoleService interface {
	GetRole(ctx context.Context, id int32) (models.Role, error)
	GetUserRoles(ctx context.Context, userId int32) ([]models.GetUserRolesRow, error)
}

type RoleService struct {
	db      *db.Database
	queries *models.Queries
}

func NewRoleService(db *db.Database) IRoleService {
	return &RoleService{db: db, queries: db.Query}
}

func (s *RoleService) GetRole(ctx context.Context, id int32) (models.Role, error) {
	role, err := s.queries.GetRole(ctx, int32(id))
	if err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (s *RoleService) GetUserRoles(ctx context.Context, userId int32) ([]models.GetUserRolesRow, error) {
	roles, err := s.queries.GetUserRoles(ctx, int32(userId))
	if err != nil {
		return []models.GetUserRolesRow{}, err
	}
	return roles, nil
}
