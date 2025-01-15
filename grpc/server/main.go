package main

import (
	"context"
	"fmt"
	pb "grpc/proto"
	"log"
	"net"
	"sort"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersServer struct {
	pb.UnimplementedUsersServer
	users sync.Map
}

func (s *UsersServer) AddUser(ctx context.Context, in *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	var resp pb.AddUserResponse

	if _, ok := s.users.Load(in.User.Email); ok {
		resp.Error = fmt.Sprintf("User with email %s exists", in.User.Email)
	} else {
		s.users.Store(in.User.Email, in.User)
	}
	return &resp, nil
}

func (s *UsersServer) ListUsers(ctx context.Context, in *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var list []string

	s.users.Range(func(key, _ interface{}) bool {
		list = append(list, key.(string))
		return true
	})

	sort.Strings(list)

	offset := int(in.Offset)
	end := int(in.Offset + in.Limit)
	if end > len(list) {
		end = len(list)
	}
	if offset >= end {
		offset = 0
		end = 0
	}
	resp := pb.ListUsersResponse{
		Count:  int32(len(list)),
		Emails: list[offset:end],
	}
	return &resp, nil
}

func (s *UsersServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var resp pb.GetUserResponse
	if user, ok := s.users.Load(in.Email); ok {
		resp.User = user.(*pb.User)
	} else {
		return nil, status.Errorf(codes.NotFound, `Пользователь с email %s не найден`, in.Email)
	}
	return &resp, nil
}

func (s *UsersServer) DelUser(ctx context.Context, in *pb.DelUserRequest) (*pb.DelUserResponse, error) {
	var resp pb.DelUserResponse

	if _, ok := s.users.LoadAndDelete(in.Email); !ok {
		resp.Error = fmt.Sprintf("User with email %s not found", in.Email)
	}
	return &resp, nil
}

func main() {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	pb.RegisterUsersServer(s, &UsersServer{})

	fmt.Println("Servar GRPC start working")

	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}

}
