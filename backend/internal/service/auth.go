package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

func (s *AuthService) Register(input RegisterInput) (*AuthResponse, error) {
	existing, _ := s.userRepo.FindByEmail(input.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) Login(input LoginInput) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: *user}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return uint(userID), nil
}

func (s *AuthService) GenerateLoginState() (string, error) {
	claims := jwt.MapClaims{
		"purpose": "github_login",
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateLoginState(state string) error {
	token, err := jwt.Parse(state, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("invalid state")
	}
	purpose, _ := claims["purpose"].(string)
	if purpose != "github_login" {
		return errors.New("invalid state purpose")
	}
	return nil
}

func (s *AuthService) GitHubLogin(ghUser *GitHubUserInfo, accessToken string) (*AuthResponse, error) {
	// 1. Try to find by GitHub ID
	user, err := s.userRepo.FindByGitHubID(ghUser.ID)
	if err == nil && user != nil {
		user.GitHubToken = accessToken
		user.GitHubUsername = ghUser.Login
		if ghUser.AvatarURL != "" {
			user.AvatarURL = ghUser.AvatarURL
		}
		s.userRepo.Update(user)

		token, err := s.generateToken(user.ID)
		if err != nil {
			return nil, err
		}
		return &AuthResponse{Token: token, User: *user}, nil
	}

	// 2. Try to find by email and link
	if ghUser.Email != "" {
		user, err = s.userRepo.FindByEmail(ghUser.Email)
		if err == nil && user != nil {
			user.GitHubID = ghUser.ID
			user.GitHubToken = accessToken
			user.GitHubUsername = ghUser.Login
			user.GitHubConnected = true
			if user.AvatarURL == "" {
				user.AvatarURL = ghUser.AvatarURL
			}
			s.userRepo.Update(user)

			token, err := s.generateToken(user.ID)
			if err != nil {
				return nil, err
			}
			return &AuthResponse{Token: token, User: *user}, nil
		}
	}

	// 3. Create new user
	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}
	email := ghUser.Email
	if email == "" {
		email = ghUser.Login + "@github.local"
	}

	newUser := &model.User{
		Name:            name,
		Email:           email,
		GitHubID:        ghUser.ID,
		GitHubUsername:  ghUser.Login,
		GitHubToken:     accessToken,
		GitHubConnected: true,
		AvatarURL:       ghUser.AvatarURL,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	token, err := s.generateToken(newUser.ID)
	if err != nil {
		return nil, err
	}
	return &AuthResponse{Token: token, User: *newUser}, nil
}

func (s *AuthService) GenerateOAuthState(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"purpose": "oauth_state",
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateOAuthState(state string) (uint, error) {
	token, err := jwt.Parse(state, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid state")
	}

	purpose, _ := claims["purpose"].(string)
	if purpose != "oauth_state" {
		return 0, errors.New("invalid state purpose")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid state claims")
	}

	return uint(userID), nil
}

func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
