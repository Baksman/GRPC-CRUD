package main

// import (
// 	"context"
// 	"fmt"
// 	"grpcapp/userproto/pkg"
// 	"io"
// 	"log"
// 	"net"
// 	"time"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/metadata"
// 	"google.golang.org/grpc/reflection"
// 	"google.golang.org/grpc/status"
// )

// type UserServer struct {
// }

// func (s *UserServer) CreateUser(ctx context.Context, req *pkg.CreateuserRequest) (*pkg.CreateUserReesponse, error) {
// 	username := req.Username
// 	name := req.Name
// 	return &pkg.CreateUserReesponse{
// 		Name:     name,
// 		Username: username,
// 	}, nil
// }

// // var testUser = &pkg.CreateUserReesponse{
// // 	Name:     "Ibrahim shehu",
// // 	Username: "Baksman99",
// // }

// func (s *UserServer) CreateUserStream(stream grpc.BidiStreamingServer[pkg.CreateuserRequest, pkg.CreateUserReesponse]) error {

// 	for i := 0; ; {
// 		user, err := stream.Recv()
// 		if err == io.EOF {
// 			return nil
// 		}
// 		currentUsername := fmt.Sprintf("%v%v", user.GetUsername(), i)
// 		currentName := fmt.Sprintf("%v%v", user.GetUsername(), i)
// 		middlewareValue, _ := stream.Context().Value("name").(string)
// 		fmt.Println(middlewareValue)
// 		stream.Send(&pkg.CreateUserReesponse{
// 			Name:     currentName,
// 			Username: currentUsername,
// 		})
// 		if err != nil {
// 			return err
// 		}

// 	}
// }
// func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	start := time.Now()
// 	log.Printf("Received request   for %s", info.FullMethod)

// 	md, _ := metadata.FromIncomingContext(ctx)
// 	log.Printf("Metadata: %v", md)

// 	resp, err := handler(ctx, req)

// 	log.Printf("Handled request for %s in %v", info.FullMethod, time.Since(start))
// 	return resp, err
// }

// func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	log.Printf("Incoming gRPC request: %s", info.FullMethod)
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	// Add your custom logic here, e.g., authentication, logging, metrics

// 	if !ok {
// 		return nil, status.Errorf(codes.InvalidArgument, "metadata not found")
// 	}

// 	authTokens := md["authorization"]
// 	if len(authTokens) > 0 {
// 		authHeader := authTokens[0]
// 		fmt.Print(authHeader)
// 	}

// 	return handler(ctx, req)
// }

// func main() {
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)

// 	lis, err := net.Listen("tcp", "0.0.0.0:2000")

// 	if err != nil {
// 		log.Fatalf("failed to listen %v", err)
// 	}
// 	// LoggingStreamInterceptor
// 	// option := grpc.UnaryInterceptor(unaryInterceptor)
// 	// option2 := grpc.StreamInterceptor(StreamLoggingInterceptor)
// 	grpcServer := grpc.NewServer(grpc.StreamInterceptor(
// 		stInterceptor,
// 	))

// 	pkg.RegisterUserServiceServer(grpcServer, &UserServer{})
// 	reflection.Register(grpcServer)
// 	fmt.Println("server listening")
// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalln("Failed to serve:", err)
// 	}
// }

// func (css *customServerStream) Context() context.Context {
// 	return css.ctx
// }

// type customServerStream struct {
// 	grpc.ServerStream
// 	ctx context.Context
// }

// func stInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 	log.Printf("%v", info.FullMethod)
// 	ctx := ss.Context()
// 	newStream := &customServerStream{
// 		ServerStream: ss,
// 		ctx:          context.WithValue(ctx, "name", "Baksmnannn"),
// 	}
// 	return handler(srv, newStream)
// 	// err := handler(srv, ss)
// 	// err := handler(srv, newServer)

// 	// return status.Error(codes.Internal, err.Error())
// }
