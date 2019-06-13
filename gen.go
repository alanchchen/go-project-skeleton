package api

//go:generate protoc -I. --go_out=plugins=grpc:. pkg/api/greeter/api.proto
//go:generate protoc -I. --go_out=plugins=grpc:. pkg/api/user/api.proto
