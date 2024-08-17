package server

import (
	"context"
	"fmt"
	models "grpcapp/model"
	"grpcapp/proto/pkg"
	"grpcapp/utils"
	"grpcapp/validators"
	"log"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

var errDB = status.Errorf(codes.Internal, "database error")
var errInvalidCred = status.Errorf(codes.PermissionDenied, "invalid credentials")
var errUserNotFound = status.Errorf(codes.NotFound, "user not found")
var errServer = status.Errorf(codes.Internal, "server error occured")
var errInvalidUserId = status.Errorf(codes.Unauthenticated, "invalid user id")

// var errUnknown = fmt.Errorf("errUnknown error occured")

type AuthServer struct {
	db        *gorm.DB
	Validator *validator.Validate
	Config    *utils.Config
}

func NewAuthServer(db *gorm.DB) *AuthServer {
	config, _ := utils.LoadConfig(".")
	vv := validator.New()
	err := vv.RegisterValidation("IsEmail", validators.IsEmail, true)

	if err != nil {
		fmt.Printf("error registering IsEmail validator %v", err.Error())
	}
	return &AuthServer{db: db, Validator: vv, Config: &config}
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,min=3,max=50,IsEmail"`
	Password string `json:"password" validate:"required,min=8,max=40"`
	Username string `json:"username" validate:"required,min=3,max=12"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,min=3,max=50,IsEmail"`
	Password string `json:"password" validate:"required,min=8,max=40"`
}

func (authService *AuthServer) Login(ctx context.Context, req *pkg.LoginRequest) (*pkg.LoginResponse, error) {
	var loginData LoginRequest = LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := authService.Validator.Struct(loginData); err != nil {
		return nil, status.Error(codes.InvalidArgument, validators.ValidatorErrorFormater(err).Error())

	}
	var user models.User
	user.Email = loginData.Email
	fmt.Println(loginData.Email, loginData.Password)
	if err := authService.db.Unscoped().Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		return nil, errUserNotFound
	}

	if err := user.ComparePassword(loginData.Password); err != nil {

		return nil, errInvalidCred
	}

	tokenString, err := user.CreateJWT()

	if err != nil {
		return nil, errServer
	}

	return &pkg.LoginResponse{
		AuthToken: tokenString,
	}, nil
}

func (authService *AuthServer) SignUp(ctx context.Context, req *pkg.SignUpRequest) (*pkg.SignUpResponse, error) {
	signUpRequest := SignUpRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Username: req.GetUsername(),
	}
	var user models.User
	if err := authService.Validator.Struct(signUpRequest); err != nil {
		return nil, validators.ValidatorErrorFormater(err)
	}
	var err error

	fmt.Println(signUpRequest.Username, signUpRequest.Email)
	if err = authService.db.Unscoped().Where("username = ? OR email = ?", signUpRequest.Username, signUpRequest.Email).First(&user).Error; err == nil {
		if signUpRequest.Email == user.Email {
			return nil, fmt.Errorf("email already exists")
		} else {
			return nil, fmt.Errorf("username already exists")
		}
	}
	if err != gorm.ErrRecordNotFound {
		return nil, errDB
	}

	user.Email = signUpRequest.Email
	user.Password = signUpRequest.Password
	user.Username = signUpRequest.Username

	if err := authService.db.Create(&user).Error; err != nil {
		return nil, errDB
	}
	return &pkg.SignUpResponse{
		Name:     user.Name,
		Username: user.Username,
	}, nil
}

func RunAuthServer(db *gorm.DB) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	lis, err := net.Listen("tcp", "0.0.0.0:2000")

	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	authService := NewAuthServer(db)
	pkg.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)
	fmt.Println("server listening")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to server:", err)
	}
}
