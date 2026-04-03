package service

import (
	"LiveChat/internal/model"
	"LiveChat/util"
	"context"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey = "secret"
)

type service struct {
	model.Repository
	timeout time.Duration
}

type MyJWTClaims struct {
	ID       string `json:"id"`
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewService(repository model.Repository) model.Service {
	return &service{
		Repository: repository,
		timeout:    10 * time.Second,
	}
}

func (s *service) CreateUser(ctx context.Context, req *model.CreateUserReq) (*model.CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &model.CreateUserRes{
		ID:       strconv.Itoa(int(r.ID)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

func (s *service) Login(ctx context.Context, req *model.LoginUserReq) (*model.LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &model.LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, u.Password)
	if err != nil {
		return &model.LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		ID:       strconv.Itoa(int(u.ID)),
		UserID:   u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &model.LoginUserRes{}, err
	}

	return &model.LoginUserRes{
		AccessToken: ss,
		Username:    u.Username,
		ID:          strconv.Itoa(int(u.ID)),
	}, nil
}
