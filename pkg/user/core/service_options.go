package core

import (
	"math/rand"

	"github.com/alanchchen/go-project-skeleton/pkg/user/api"
)

type Option func(*UserService)

func WithBuiltInUsers(name string) Option {
	return func(s *UserService) {
		s.users = append(s.users, &api.User{
			Name: name,
			Id:   rand.Int63(),
		})
	}
}
