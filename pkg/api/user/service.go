package user

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"sync"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
)

//go:generate mockgen -source=api.pb.go -destination=mock/api.go -package=mock

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewService(opts ...Option) *UserService {
	s := &UserService{
		mu:     new(sync.RWMutex),
		logger: log.New(ioutil.Discard, "", 0),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type UserService struct {
	users  []*User
	mu     *sync.RWMutex
	logger *log.Logger
}

func (s *UserService) Bind(server *grpc.Server) {
	RegisterServiceServer(server, s)
}

func (s *UserService) AddUser(ctx context.Context, req *AddUserRequest) (*Users, error) {
	if req.Name == "" {
		s.logger.Println("the given name is empty")
		return nil, errors.New("please tell me your name")
	}

	user := &User{
		Name: req.Name,
		Id:   rand.Int63(),
	}

	s.mu.Lock()
	s.users = append(s.users, user)
	s.mu.Unlock()

	return &Users{
		Users: []*User{
			user,
		},
	}, nil
}

func (s *UserService) FindUserById(ctx context.Context, req *FindUserByIdRequest) (*Users, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Id == req.Id {
			return &Users{
				Users: []*User{
					u,
				},
			}, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *UserService) FindUserByName(ctx context.Context, req *FindUserByNameRequest) (*Users, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, u := range s.users {
		if u.Name == req.Name {
			return &Users{
				Users: []*User{
					u,
				},
			}, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *UserService) ListUsers(ctx context.Context, _ *empty.Empty) (*Users, error) {
	return &Users{
		Users: s.users,
	}, nil
}

// ----------------------------------------------------------------------------

func randomID(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return string(b)
}
