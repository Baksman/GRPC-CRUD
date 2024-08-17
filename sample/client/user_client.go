package main

// import (
// 	"context"
// 	"fmt"
// 	"grpcapp/userproto/pkg"
// 	"io"
// 	"log"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// var users = []*pkg.CreateuserRequest{
// 	{
// 		Name:     "Ibrahim",
// 		Username: "Baksman",
// 	},
// 	{
// 		Name:     "Yussuf",
// 		Username: "Yusure",
// 	},
// 	{
// 		Name:     "Muda",
// 		Username: "Mzinida",
// 	},
// }

// func main() {
// 	userClient, err := grpc.NewClient("0.0.0.0:2000", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer userClient.Close()

// 	c := pkg.NewUserServiceClient(userClient)

// 	userResponse, err := c.CreateUser(context.Background(), &pkg.CreateuserRequest{
// 		Name:     "Ibrahim",
// 		Username: "Baksman",
// 	})

// 	if err != nil {
// 		log.Fatal("error getting response: ", err)
// 	}

// 	fmt.Printf("%v %v", userResponse.Name, userResponse.Username)

// 	stream, err := c.CreateUserStream(context.Background())
// 	waitc := make(chan struct{})
// 	go func() {
// 		for {
// 			in, err := stream.Recv()
// 			if err == io.EOF {
// 				waitc <- struct{}{}
// 				close(waitc)
// 				return
// 			}
// 			if err != nil {
// 				log.Fatalf("Failed to receive a user : %v", err)
// 			}
// 			log.Printf("username is %s,name %s ", in.Username, in.Name)
// 		}
// 	}()

// 	for _, note := range users {
// 		if err := stream.Send(note); err != nil {
// 			log.Fatalf("Failed to send a note: %v", err)
// 		}
// 	}
// 	stream.CloseSend()
// 	<-waitc
// }
