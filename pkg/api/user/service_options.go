package user

import "log"

type Option func(*UserService)

func UseLogger(logger *log.Logger) Option {
	return func(s *UserService) {
		s.logger = logger
	}
}
