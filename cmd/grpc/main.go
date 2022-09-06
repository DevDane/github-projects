package main

import (
	"fmt"
	"net"
	"os"

	"devdane.com/repos"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	godotenv.Load(".env")
	gs := grpc.NewServer()
	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	reflection.Register(gs)

	svc := repos.NewService()
	svc.RegisterService(gs)

	if err := gs.Serve(lis); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
