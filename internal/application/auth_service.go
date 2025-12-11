package application

import (
	"context"
	"errors"
	"time"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecret 密钥，生产环境应该从环境变量读取
var JWTSecret = []byte("your_super_secret_key_change_this_in_prod")

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
}

// Login 验证并返回 Token
func (s *AuthService) LoginJwt(ctx context.Context, username, password string) (string, error) {
	ctx, err := s.LoginCheck(ctx, username, password)
	if err != nil {
		return "", err
	}
	userId := ctx.Value("userID").(uint)
	role := ctx.Value("role").(string)
	// 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时过期
	})

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) LoginCheck(ctx context.Context, username, password string) (context.Context, error) {
	// 1. 查询用户
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return ctx, errors.New("invalid username or password") // 模糊报错，防止枚举攻击
	}

	// 2. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return ctx, errors.New("invalid username or password")
	}
	ctx = context.WithValue(ctx, "userID", user.ID)
	ctx = context.WithValue(ctx, "role", user.Role)
	return ctx, nil
}

// Register 注册新用户 (用于初始化管理员)
func (s *AuthService) Register(ctx context.Context, username, password string, role string) error {
	// 1. 哈希密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPwd),
		Role:         role,
		IsActive:     true,
	}

	return s.userRepo.Save(ctx, user)
}
