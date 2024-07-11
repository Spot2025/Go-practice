package main

import (
	"awesomeProject/accounts/models"
	"awesomeProject/proto"
	"context"
	"fmt"
	"net"
	"sync"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedBankServer
	accounts map[string]*models.Account
	guard    *sync.RWMutex
}

func NewServer() *server {
	return &server{
		accounts: make(map[string]*models.Account),
		guard:    &sync.RWMutex{},
	}
}

func (s *server) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty name")
	}

	s.guard.Lock()

	if _, ok := s.accounts[req.Name]; ok {
		s.guard.Unlock()

		return nil, status.Error(codes.AlreadyExists, "account already exists")
	}

	s.accounts[req.Name] = &models.Account{
		Name:   req.Name,
		Amount: int(req.Amount),
	}

	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	s.guard.Lock()

	account, ok := s.accounts[req.Name]

	s.guard.Unlock()

	if !ok {
		return nil, status.Error(codes.NotFound, "account not found")
	}

	response := proto.GetAccountResponse{
		Name: account.Name,
		Amount: int32(account.Amount),
	}

	return &response, nil
}

func (s *server) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.Empty, error) {
	s.guard.Lock()

	_, ok := s.accounts[req.Name]

	if !ok {
		s.guard.Unlock()

		return nil, status.Error(codes.NotFound, "account not found")
	}

	delete(s.accounts, req.Name)

	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) PatchAccount(ctx context.Context, req *proto.PatchAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty name")
	}

	s.guard.Lock()

	if _, ok := s.accounts[req.Name]; !ok {
		s.guard.Unlock()

		return nil, status.Error(codes.NotFound, "account not found")
	}

	s.accounts[req.Name].Amount = int(req.Amount)

	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func (s *server) ChangeAccount(ctx context.Context, req *proto.ChangeAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 || len(req.NewName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty name")
	}

	s.guard.Lock()

	if _, ok := s.accounts[req.Name]; !ok {
		s.guard.Unlock()

		return nil, status.Error(codes.NotFound, "account not found")
	}

	if _, ok := s.accounts[req.NewName]; ok {
		s.guard.Unlock()

		return nil, status.Error(codes.AlreadyExists, "account with this name already exists")
	}

	s.accounts[req.NewName] = &models.Account{
		Name: req.NewName,
		Amount: s.accounts[req.Name].Amount,
	}

	delete(s.accounts, req.Name)

	s.guard.Unlock()

	return &proto.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 1323))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterBankServer(s, NewServer())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}