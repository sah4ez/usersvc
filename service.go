package usersvc

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

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]User
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]User{},
	}
}

func (s *inmemService) PostUser(ctx context.Context, p User) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[p.ID]; ok {
		return ErrAlreadyExists
	}
	s.m[p.ID] = p
	return nil
}

func (s *inmemService) GetUser(ctx context.Context, id string) (User, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return p, nil
}

func (s *inmemService) PatchUser(ctx context.Context, id string, p User) error {
	if p.ID != "" && id != p.ID {
		return ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, ok := s.m[id]
	if !ok {
		return ErrNotFound
	}

	s.m[id] = p
	return nil
}

func (s *inmemService) GetUsers(ctx context.Context) ([]User, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	users := make([]User, len(s.m))
	i := 0
	for _, v := range s.m {
		users[i] = v
		i++
	}
	return users, nil
}
