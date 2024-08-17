package client

import (
	"grpcapp/proto/pkg"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


func RegisterUserClient() *pkg.UserServiceClient {
	userClient, err := grpc.NewClient("0.0.0.0:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := pkg.NewUserServiceClient(userClient)
	return &c

}
