package services

import (
	"context"
	db "project-sqlc/internal/db/models"
)

type IRoleService interface {
	GetRole(ctx context.Context, id int32) (db.Role, error)
	GetUserRoles(ctx context.Context, userId int32) ([]db.GetUserRolesRow, error)
}

type RoleService struct {
	db *db.Queries
}

func NewRoleService(db *db.Queries) IRoleService {
	return &RoleService{db: db}
}

func (s *RoleService) GetRole(ctx context.Context, id int32) (db.Role, error) {
	role, err := s.db.GetRole(ctx, int32(id))
	if err != nil {
		return db.Role{}, err
	}
	return role, nil
}

func (s *RoleService) GetUserRoles(ctx context.Context, userId int32) ([]db.GetUserRolesRow, error) {
	roles, err := s.db.GetUserRoles(ctx, int32(userId))
	if err != nil {
		return []db.GetUserRolesRow{}, err
	}
	return roles, nil
}
