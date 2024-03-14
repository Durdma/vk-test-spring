package service

import (
	"context"
	"vk-test-spring/internal/repository"
)

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) CreateUser(ctx context.Context, input UserInput) error {
	return nil
}

func (s *UsersService) DeleteUser(ctx context.Context, userId string) error {
	return nil
}

func (s *UsersService) ChangeRole(ctx context.Context, userId string, role string) error {
	return nil
}

func (s *UsersService) GetUserIdRole(username string, password string) (string, string, error) {
	userId, role, err := s.repo.GetUserIdRole(username, password)
	if userId == "" || role == "" || err != nil {
		return "", "", err
	}

	return userId, role, nil
}
