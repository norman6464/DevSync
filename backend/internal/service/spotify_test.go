package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newSpotifyTestService() (*SpotifyService, *MockUserRepository, *MockSpotifyRepository) {
	userRepo := new(MockUserRepository)
	spotifyRepo := new(MockSpotifyRepository)
	cfg := &config.Config{
		SpotifyClientID:     "test-client-id",
		SpotifyClientSecret: "test-client-secret",
		SpotifyRedirectURL:  "http://localhost:5173/spotify/callback",
	}
	svc := NewSpotifyService(cfg, userRepo, spotifyRepo)
	return svc, userRepo, spotifyRepo
}

func TestSpotifyGetOAuthURL(t *testing.T) {
	svc, _, _ := newSpotifyTestService()

	url := svc.GetOAuthURL("test-state")

	assert.Contains(t, url, "https://accounts.spotify.com/authorize")
	assert.Contains(t, url, "client_id=test-client-id")
	assert.Contains(t, url, "response_type=code")
	assert.Contains(t, url, "state=test-state")
	assert.Contains(t, url, "user-read-currently-playing")
	assert.Contains(t, url, "user-read-recently-played")
}

func TestSpotifyDisconnect_Success(t *testing.T) {
	svc, userRepo, spotifyRepo := newSpotifyTestService()

	user := &model.User{SpotifyConnected: true, SpotifyToken: "token", SpotifyRefreshToken: "refresh"}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("Update", mock.MatchedBy(func(u *model.User) bool {
		return !u.SpotifyConnected && u.SpotifyToken == "" && u.SpotifyRefreshToken == ""
	})).Return(nil)
	spotifyRepo.On("DeleteUserData", uint(1)).Return(nil)

	err := svc.DisconnectSpotify(1)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
	spotifyRepo.AssertExpectations(t)
}

func TestSpotifyDisconnect_UserNotFound(t *testing.T) {
	svc, userRepo, _ := newSpotifyTestService()

	userRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.DisconnectSpotify(999)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestSpotifyDisconnect_UpdateError(t *testing.T) {
	svc, userRepo, _ := newSpotifyTestService()

	user := &model.User{SpotifyConnected: true}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)
	userRepo.On("Update", mock.Anything).Return(errors.New("db error"))

	err := svc.DisconnectSpotify(1)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestSpotifyGetCurrentlyPlaying_UserNotFound(t *testing.T) {
	svc, userRepo, _ := newSpotifyTestService()

	userRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetCurrentlyPlaying(999)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}

func TestSpotifyGetCurrentlyPlaying_NotConnected(t *testing.T) {
	svc, userRepo, _ := newSpotifyTestService()

	user := &model.User{SpotifyConnected: false, SpotifyRefreshToken: ""}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetCurrentlyPlaying(1)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "連携されていません")
}

func TestSpotifyGetRecentlyPlayed_NotConnected(t *testing.T) {
	svc, userRepo, _ := newSpotifyTestService()

	user := &model.User{SpotifyConnected: false, SpotifyRefreshToken: ""}
	user.ID = 1

	userRepo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetRecentlyPlayed(1)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "連携されていません")
}

func TestSpotifyGetRecentlyPlayed_UserNotFound(t *testing.T) {
	svc, userRepo, _ := newSpotifyTestService()

	userRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetRecentlyPlayed(999)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
}
