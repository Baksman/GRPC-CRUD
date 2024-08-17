package client

import (
	"grpcapp/proto/pkg"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterAuthClient() *pkg.AuthServiceClient {
	authClient, err := grpc.NewClient("0.0.0.0:2000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := pkg.NewAuthServiceClient(authClient)
	return &c
}
