package service

import (
	"context"
	"errors"
	"sync"
)

type Service interface {
	PostUser(ctx context.Context, p User) error
	GetUser(ctx context.Context, id string) (User, error)
	PatchUser(ctx context.Context, id string, p User) error
	GetUsers(ctx context.Context) ([]User, error)
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

