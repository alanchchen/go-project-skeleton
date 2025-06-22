package core

import (
	"context"
	"errors"
	"math/rand"
	"sync"

	empty "github.com/golang/protobuf/ptypes/empty"

	"github.com/alanchchen/go-project-skeleton/pkg/user/api"
)

func NewService(opts ...Option) *UserService {
	s := &UserService{
		mu: new(sync.RWMutex),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type UserService struct {
	api.UnimplementedUserServiceServer

	users []*api.User
	mu    *sync.RWMutex
}

func (s *UserService) AddUser(ctx context.Context, req *api.AddUserRequest) (*api.Users, error) {
	if req.Name == "" {
		return nil, errors.New("please tell me your name")
	}

	user := &api.User{
		Name: req.Name,
		Id:   rand.Int63(),
	}

	s.mu.Lock()
	s.users = append(s.users, user)
	s.mu.Unlock()

	return &api.Users{
		Users: []*api.User{
			user,
		},
	}, nil
}

func (s *UserService) FindUserById(ctx context.Context, req *api.FindUserByIdRequest) (*api.Users, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Id == req.Id {
			return &api.Users{
				Users: []*api.User{
					u,
				},
			}, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *UserService) FindUserByName(ctx context.Context, req *api.FindUserByNameRequest) (*api.Users, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Name == req.Name {
			return &api.Users{
				Users: []*api.User{
					u,
				},
			}, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *UserService) ListUsers(ctx context.Context, _ *empty.Empty) (*api.Users, error) {
	return &api.Users{
		Users: s.users,
	}, nil
}
