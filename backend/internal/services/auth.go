package services

import (
	"errors"
	"orchestrix/internal/config"
	"orchestrix/internal/middleware"
	"orchestrix/internal/models"
	"orchestrix/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg      *config.Config
	userRepo *repository.UserRepository
	orgRepo  *repository.OrganizationRepository
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		cfg:      cfg,
		userRepo: repository.NewUserRepository(),
		orgRepo:  repository.NewOrganizationRepository(),
	}
}

func (s *AuthService) Register(orgName, email, password, name, role string) (*models.User, string, error) {
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, "", errors.New("email already registered")
	}

	org, err := s.orgRepo.Create(orgName)
	if err != nil {
		return nil, "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user, err := s.userRepo.Create(org.ID, email, string(hashedPassword), name, role)
	if err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := &middleware.Claims{
		UserID:         user.ID,
		Email:          user.Email,
		OrganizationID: user.OrganizationID,
		Role:           user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) GetUserByID(id int64) (*models.User, error) {
	return s.userRepo.GetByID(id)
}
