package impl

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"

	"study4cash/DB/models"
	"study4cash/auth"
	usersv1 "study4cash/gen/users/v1"
	"study4cash/gen/users/v1/usersv1connect"
)

type UserServer struct {
	db *gorm.DB
}

func (s *UserServer) Login(ctx context.Context, req *usersv1.LoginRequest) (*usersv1.LoginResponse, error) {
	users, err := gorm.G[models.User](s.db).Where("email = ?", req.Email).Find(ctx)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("User not found"))
	}
	user := users[0]
	match, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil {
		return nil, err
	}
	if !match {
		err := connect.NewError(connect.CodeInvalidArgument, errors.New("Invalid password"))
		return nil, err
	}
	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &usersv1.LoginResponse{
		Token:  *token,
		UserId: uint64(user.ID),
	}, nil
}

func (s *UserServer) Register(ctx context.Context, req *usersv1.RegisterRequest) (*usersv1.RegisterResponse, error) {
	user := req.User
	users, err := gorm.G[models.User](s.db).Where("email = ?", user.Email).Find(ctx)
	if err != nil {
		return nil, err
	}
	if len(users) != 0 {
		err := connect.NewError(connect.CodeAlreadyExists, errors.New("User already exists"))
		return nil, err
	}

	password, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}
	userRecord := models.User{
		Email:    user.Email,
		Password: password,
		Name:     user.Name,
		Surname:  user.Surname,
		Active:   true,
	}
	err = gorm.G[models.User](s.db).Create(ctx, &userRecord)
	if err != nil {
		return nil, err
	}

	token, err := auth.GenerateJWT(userRecord.ID)
	if err != nil {
		return nil, err
	}

	return &usersv1.RegisterResponse{
		Token:  *token,
		UserId: uint64(userRecord.ID),
	}, nil
}

func NewUsersServer(db *gorm.DB) (string, http.Handler) {
	user := &UserServer{
		db: db,
	}
	path, handler := usersv1connect.NewUserServiceHandler(
		user,
		connect.WithInterceptors(validate.NewInterceptor()),
	)
	return path, handler
}
