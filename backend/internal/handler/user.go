package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// UserServiceInterface はUserServiceが実装すべきインターフェース。
type UserServiceInterface interface {
	GetAll(query string) ([]model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	UpdateProfile(id, userID uint, input *service.UpdateProfileInput) (*model.User, error)
	GetProfileCompleteness(userID uint) (*service.ProfileCompleteness, error)
}

// UserHandler はユーザー関連のHTTPハンドラ。
// ユーザー検索・詳細取得・プロフィール更新を処理する。
type UserHandler struct {
	service UserServiceInterface
}

// NewUserHandler は新しいUserHandlerインスタンスを生成する。
func NewUserHandler(s UserServiceInterface) *UserHandler {
	return &UserHandler{service: s}
}

// GetAll はユーザー一覧を返す。クエリパラメータqで検索可能（最大100文字）。
func (h *UserHandler) GetAll(c *gin.Context) {
	q, ok := parseOptionalSearchQuery(c, maxUserSearchQueryLen)
	if !ok {
		return
	}
	users, err := h.service.GetAll(q)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, users)
}

// GetByID は指定IDのユーザー情報を返す。
func (h *UserHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	user, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, user)
}

// GetByUsername は指定ユーザー名のユーザー情報を返す。
func (h *UserHandler) GetByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		respondBadRequest(c, "username required")
		return
	}
	if len([]rune(username)) > maxUsernameLen {
		respondBadRequest(c, "ユーザー名は50文字以下である必要があります")
		return
	}

	user, err := h.service.GetByUsername(username)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, user)
}

// Update はユーザープロフィールを更新する。
// 所有権チェック・フィールドマッピング・バリデーションはService層で実施する。
func (h *UserHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateUserRequest](c)
	if input == nil {
		return
	}

	profileInput := &service.UpdateProfileInput{
		Name:                input.Name,
		Bio:                 input.Bio,
		AvatarURL:           input.AvatarURL,
		SkillsLanguages:     input.SkillsLanguages,
		SkillsFrameworks:    input.SkillsFrameworks,
		AtCoderUsername:     input.AtCoderUsername,
		PaizaRank:           input.PaizaRank,
		OnboardingCompleted: input.OnboardingCompleted,
	}

	user, err := h.service.UpdateProfile(id, userID, profileInput)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, user)
}

// GetProfileCompleteness は認証ユーザーのプロフィール完成度を返す。
func (h *UserHandler) GetProfileCompleteness(c *gin.Context) {
	userID := c.GetUint("userID")
	result, err := h.service.GetProfileCompleteness(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, result)
}
