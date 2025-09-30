package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Bilal-Cplusoft/sun_ready/internal/models"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repo.UserRepo
	jwtSecret string
}

func NewAuthService(userRepo *repo.UserRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

type Claims struct {
	UserID    int `json:"user_id"`
	CompanyID int `json:"company_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(ctx context.Context, email, password, firstName, lastName string, companyID int) (*models.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedStr := string(hashedPassword)
	user := &models.User{
		Email:     email,
		Password:  &hashedStr,
		FirstName: &firstName,
		LastName:  &lastName,
		CompanyID: companyID,
		Type:      int16(models.UserTypeSales),
		Disabled:  false,
		IsManager: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if user.Disabled {
		return "", nil, errors.New("user account is disabled")
	}

	if user.Password == nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.GenerateToken(user.ID, user.CompanyID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *AuthService) GenerateToken(userID, companyID int) (string, error) {
	claims := &Claims{
		UserID:    userID,
		CompanyID: companyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
