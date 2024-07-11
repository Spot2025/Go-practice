package main

import (
	"awesomeProject/proto"
	"context"
	"flag"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Command struct {
	Port   	int
	Host   	string
	Cmd    	string
	Name   	string
	Amount 	int
	NewName string
}

func main() {
	portVal := flag.Int("port", 8080, "server port")
	hostVal := flag.String("host", "0.0.0.0", "server host")
	cmdVal := flag.String("cmd", "", "command to execute")
	nameVal := flag.String("name", "", "name of account")
	amountVal := flag.Int("amount", 0, "amount of account")
	newNameVal := flag.String("new_name", "", "new name of account")

	flag.Parse()

	cmd := Command{
		Port:    *portVal,
		Host:    *hostVal,
		Cmd:     *cmdVal,
		Name:    *nameVal,
		Amount:  *amountVal,
		NewName: *newNameVal,
	}

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", cmd.Host, cmd.Port),
								grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = conn.Close()
	}()

	c := proto.NewBankClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	switch cmd.Cmd {
	case "create":
		_, err := c.CreateAccount(ctx, &proto.CreateAccountRequest{
			Name: cmd.Name,
			Amount: int32(cmd.Amount),
		})
		if err != nil {
			panic(err)
		}
	case "get":
		resp, err := c.GetAccount(ctx, &proto.GetAccountRequest{
			Name: cmd.Name,
		})
		if err != nil {
			panic(err)
		}

		fmt.Printf("response account name: %s and amount: %d\n", resp.Name, resp.Amount)
	case "delete":
		_, err := c.DeleteAccount(ctx, &proto.DeleteAccountRequest{
			Name: cmd.Name,
		})
		if err != nil {
			panic(err)
		}
	case "patch_amount":
		_, err := c.PatchAccount(ctx, &proto.PatchAccountRequest{
			Name: cmd.Name,
			Amount: int32(cmd.Amount),
		})
		if err != nil {
			panic(err)
		}
	case "patch_name":
		_, err := c.ChangeAccount(ctx, &proto.ChangeAccountRequest{
			Name: cmd.Name,
			NewName: cmd.NewName,
		})
		if err != nil {
			panic(err)
		}
	default:
		fmt.Printf("unknown command %s", cmd.Cmd)
	}
}	