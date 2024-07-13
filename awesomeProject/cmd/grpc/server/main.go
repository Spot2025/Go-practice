package main

import (
	"awesomeProject/accounts/models"
	"awesomeProject/proto"
	"context"
	"database/sql"
	"fmt"
	"net"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedBankServer
	db *sql.DB
}

func (s *server) GetAccountFromStorage(Name string) (models.Account, error) {
	row := s.db.QueryRow("SELECT name, amount FROM accounts WHERE name=$1", Name)
	response := models.Account{}
	err := row.Scan(&response.Name, &response.Amount)
	
	return response, err
}

func (s *server) InsertAccountToStorage(account models.Account) error {
	_, err := s.db.Exec("INSERT INTO accounts(name, amount) VALUES($1, $2)", account.Name, account.Amount)
	return err
}

func (s *server) DeleteAccountFromStorage(Name string) error {
	_, err := s.db.Exec("DELETE FROM accounts WHERE name=$1", Name)
	return err
}

func (s *server) ChangeAccountAmountInStorage(account models.Account) error {
	_, err := s.db.Exec("UPDATE accounts SET amount=$1 WHERE name=$2", account.Amount, account.Name)
	return err
}

func (s *server) ChangeAccountNameInStorage(Name, NewName string) error {
	_, err := s.db.Exec("UPDATE accounts SET name=$1 WHERE name=$2", NewName, Name)
	return err
}

func (s *server) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty name")
	}

	if _, err := s.GetAccountFromStorage(req.Name); err == nil {
		return nil, status.Error(codes.AlreadyExists, "account already exists")
	}

	err := s.InsertAccountToStorage(models.Account{
		Name:   req.Name,
		Amount: int(req.Amount),
	})

	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *server) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.GetAccountResponse, error) {
	account, err := s.GetAccountFromStorage(req.Name)

	if err != nil {
		return nil, status.Error(codes.NotFound, "account not found")
	}

	response := proto.GetAccountResponse{
		Name: account.Name,
		Amount: int32(account.Amount),
	}

	return &response, nil
}

func (s *server) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.Empty, error) {
	_, err := s.GetAccountFromStorage(req.Name)

	if err != nil {
		return nil, status.Error(codes.NotFound, "account not found")
	}

	err = s.DeleteAccountFromStorage(req.Name)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *server) PatchAccount(ctx context.Context, req *proto.PatchAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty name")
	}

	

	if _, err := s.GetAccountFromStorage(req.Name); err != nil {
		return nil, status.Error(codes.NotFound, "account not found")
	}

	err := s.ChangeAccountAmountInStorage(models.Account{
		Name:   req.Name,
		Amount: int(req.Amount),
	})

	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *server) ChangeAccount(ctx context.Context, req *proto.ChangeAccountRequest) (*proto.Empty, error) {
	if len(req.Name) == 0 || len(req.NewName) == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty name")
	}

	if _, err := s.GetAccountFromStorage(req.Name); err != nil {
		return nil, status.Error(codes.NotFound, "account not found")
	} 
	if _, err := s.GetAccountFromStorage(req.NewName); err == nil {
		return nil, status.Error(codes.AlreadyExists, "account with this name already exists")
	}

	err := s.ChangeAccountNameInStorage(req.Name, req.NewName)
	if err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func main() {
	connectionString := "host=0.0.0.0 port=5432 dbname=postgres user=postgres password=password"
	
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 1323))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	proto.RegisterBankServer(s, &server{db: db})
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}