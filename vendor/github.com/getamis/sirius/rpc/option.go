// Copyright 2017 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpc

import (
	"crypto/tls"

	"github.com/getamis/sirius/metrics"
	"google.golang.org/grpc"
)

type ServerOption func(*Server)

// APIs to be registered to the RPC server
func APIs(apis ...API) ServerOption {
	return func(s *Server) {
		s.apis = apis
	}
}

// Credentials for the RPC server
func Credentials(credentials *tls.Config) ServerOption {
	return func(s *Server) {
		s.credentials = credentials
	}
}

func Metrics(metrics metrics.ServerMetrics) ServerOption {
	return func(s *Server) {
		s.grpcMetrics = metrics
	}
}

func StreamInterceptors(interceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInterceptors = interceptors
	}
}

func UnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInterceptors = interceptors
	}
}
