package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// UserServiceInterface はUserServiceが実装すべきインターフェース。
type UserServiceInterface interface {
	GetAll(query string) ([]model.User, error)
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
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

// GetAll はユーザー一覧を返す。クエリパラメータqで検索可能。
func (h *UserHandler) GetAll(c *gin.Context) {
	q := c.Query("q")
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
		respondNotFound(c, "user not found")
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

	user, err := h.service.GetByUsername(username)
	if err != nil {
		respondNotFound(c, "user not found")
		return
	}
	respondOK(c, user)
}

// Update はユーザープロフィールを更新する。
// 本人のみ更新可能（userIDとパスパラメータのIDが一致する必要がある）。
func (h *UserHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	userID := c.GetUint("userID")
	if userID != id {
		respondForbidden(c, "cannot update other user's profile")
		return
	}

	existing, err := h.service.GetByID(id)
	if err != nil {
		respondNotFound(c, "user not found")
		return
	}

	input := bindJSON[dto.UpdateUserRequest](c)
	if input == nil {
		return
	}

	if input.Name != "" {
		existing.Name = input.Name
	}
	existing.Bio = input.Bio
	existing.AvatarURL = input.AvatarURL
	if input.SkillsLanguages != nil {
		existing.SkillsLanguages = *input.SkillsLanguages
	}
	if input.SkillsFrameworks != nil {
		existing.SkillsFrameworks = *input.SkillsFrameworks
	}
	if input.AtCoderUsername != nil {
		existing.AtCoderUsername = *input.AtCoderUsername
	}
	if input.PaizaRank != nil {
		existing.PaizaRank = *input.PaizaRank
	}
	if input.OnboardingCompleted != nil {
		existing.OnboardingCompleted = *input.OnboardingCompleted
	}

	if err := h.service.Update(existing); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, existing)
}
