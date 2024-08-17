package server

import (
	"context"
	"errors"
	"fmt"
	models "grpcapp/model"
	"grpcapp/proto/pkg"
	"grpcapp/utils"
	"log"
	"net"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var config utils.Config

type UserServer struct {
	db        *gorm.DB
	Validator *validator.Validate
	Config    *utils.Config
	mutex     sync.RWMutex
}

func NewUserServer(db *gorm.DB) *UserServer {
	var once sync.Once
	once.Do(func() {
		config, _ = utils.LoadConfig(".")
	})

	validator := validator.New()
	return &UserServer{
		Validator: validator,
		db:        db,
		Config:    &config,
	}
}

func (userService *UserServer) GetUserDetails(context context.Context, empty *pkg.Empty) (*pkg.User, error) {
	var user models.User
	rawUserId := context.Value("userId")

	userId, ok := rawUserId.(uint)

	if !ok {
		return nil, errInvalidUserId
	}

	if err := userService.db.First(&user, "id = ?", userId).Error; err != nil {
		if isNotfound := errors.Is(err, gorm.ErrRecordNotFound); isNotfound {
			return nil, errUserNotFound
		}

		return nil, errDB

	}
	return &pkg.User{
		Name:           user.Name,
		HashedPassword: user.Password,
		Email:          user.Email,
	}, nil
}

func (userService *UserServer) DeleteAccount(context context.Context, req *pkg.Empty) (*pkg.Empty, error) {
	var user models.User
	rawUserId := context.Value("userId")

	userId, ok := rawUserId.(uint)

	if !ok {
		return nil, errInvalidUserId
	}

	if err := userService.db.Where("id = ?", userId).Delete(&user).Error; err != nil {
		if isNotfound := errors.Is(err, gorm.ErrRecordNotFound); isNotfound {
			return nil, errUserNotFound
		}
		return nil, errDB
	}

	return &pkg.Empty{}, nil
}

func (userService *UserServer) UpdateUserDetails(context context.Context, req *pkg.UpdateUserDetailsRequest) (*pkg.User, error) {
	var user models.User
	rawUserId := context.Value("userId")

	userId, ok := rawUserId.(uint)

	if !ok {
		return nil, errInvalidUserId
	}

	if err := userService.db.First(&user, "id = ?", userId).Error; err != nil {
		if isNotfound := errors.Is(err, gorm.ErrRecordNotFound); isNotfound {
			return nil, errUserNotFound
		}

		return nil, errDB

	}
	user.Name = req.GetName()

	if err := userService.db.Save(&user); err != nil {
		return nil, errDB
	}
	return &pkg.User{
		Name:           user.Name,
		HashedPassword: user.Password,
		Email:          user.Email,
	}, nil
}

func (userService *UserServer) GetAllUsers(context.Context, *pkg.Empty) (*pkg.UserList, error) {
	users := []models.User{}

	if err := userService.db.Find(&users).Error; err != nil {
		if isNotfound := errors.Is(err, gorm.ErrRecordNotFound); isNotfound {
			return nil, errUserNotFound
		}

		return nil, errDB
	}

	pkgUsers := []*pkg.User{}
	for i := 0; i < len(users); i++ {
		currentUser := users[i]
		currentPbUser := pkg.User{
			Name:           currentUser.Username,
			HashedPassword: currentUser.Password,
			Email:          currentUser.Email,
		}
		userService.mutex.Lock()
		pkgUsers = append(pkgUsers, &currentPbUser)
		userService.mutex.Unlock()
		// pbUsers = append(pbUsers, currentPbUser)
	}
	return &pkg.UserList{
		Users: pkgUsers,
	}, nil
}

func RunUserServer(db *gorm.DB) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lis, err := net.Listen("tcp", "0.0.0.0:3000")

	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			authInterceptor,
		),
	)

	userService := NewUserServer(db)
	pkg.RegisterUserServiceServer(grpcServer, userService)

	reflection.Register(grpcServer)

	fmt.Println("server listening")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

func authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Incoming gRPC request: %s", info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
	}

	tokenString := authHeader[0]
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, fmt.Errorf("invalid token signature")

		}
		return nil, fmt.Errorf("error parsing token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	ctx2 := context.WithValue(ctx, "userId", claims.ID)

	return handler(ctx2, req)
}
