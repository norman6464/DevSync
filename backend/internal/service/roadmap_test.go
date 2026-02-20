package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestRoadmapService はRoadmapServiceのテスト用インスタンスを生成するヘルパー。
func newTestRoadmapService() (*RoadmapService, *MockRoadmapRepository) {
	repo := new(MockRoadmapRepository)
	svc := NewRoadmapService(repo)
	return svc, repo
}

// ============================================================
// GetByID（可視性チェック）
// ============================================================

func TestRoadmapGetByID_PublicRoadmap(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "Public Roadmap", UserID: 1, IsPublic: true}
	roadmap.ID = 1

	repo.On("FindByID", uint(1)).Return(roadmap, nil)

	result, err := svc.GetByID(1, 999) // 他人がアクセス
	assert.NoError(t, err)
	assert.Equal(t, "Public Roadmap", result.Title)
	repo.AssertExpectations(t)
}

func TestRoadmapGetByID_PrivateOwner(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "Private", UserID: 1, IsPublic: false}
	roadmap.ID = 1

	repo.On("FindByID", uint(1)).Return(roadmap, nil)

	result, err := svc.GetByID(1, 1) // 所有者がアクセス
	assert.NoError(t, err)
	assert.Equal(t, "Private", result.Title)
	repo.AssertExpectations(t)
}

func TestRoadmapGetByID_PrivateForbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "Private", UserID: 1, IsPublic: false}
	roadmap.ID = 1

	repo.On("FindByID", uint(1)).Return(roadmap, nil)

	result, err := svc.GetByID(1, 999) // 他人がアクセス
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// ロードマップコピーテスト
// ============================================================

func TestCopyRoadmap_PublicSuccess(t *testing.T) {
	svc, repo := newTestRoadmapService()

	original := &model.Roadmap{Title: "Original", UserID: 1, IsPublic: true}
	original.ID = 1

	copied := &model.Roadmap{Title: "Original", UserID: 2}
	copied.ID = 2

	repo.On("FindByID", uint(1)).Return(original, nil)
	repo.On("CopyRoadmap", uint(1), uint(2)).Return(copied, nil)

	result, err := svc.CopyRoadmap(1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(2), result.UserID)
	repo.AssertExpectations(t)
}

func TestCopyRoadmap_PrivateForbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	original := &model.Roadmap{Title: "Private", UserID: 1, IsPublic: false}
	original.ID = 1

	repo.On("FindByID", uint(1)).Return(original, nil)

	result, err := svc.CopyRoadmap(1, 999) // 他人が非公開をコピー
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestCopyRoadmap_PrivateOwnerCanCopy(t *testing.T) {
	svc, repo := newTestRoadmapService()

	original := &model.Roadmap{Title: "Private", UserID: 1, IsPublic: false}
	original.ID = 1

	copied := &model.Roadmap{Title: "Private", UserID: 1}
	copied.ID = 2

	repo.On("FindByID", uint(1)).Return(original, nil)
	repo.On("CopyRoadmap", uint(1), uint(1)).Return(copied, nil)

	result, err := svc.CopyRoadmap(1, 1) // 所有者は自分のをコピー可能
	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// ロードマップ更新テスト
// ============================================================

func TestRoadmapUpdate_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{Title: "Old", UserID: 1, Status: model.RoadmapStatusActive}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.Roadmap{Title: "New Title"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.Roadmap{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdate_CompletedStatus(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{Title: "Roadmap", UserID: 1, Status: model.RoadmapStatusActive}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.Roadmap{Status: model.RoadmapStatusCompleted}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, model.RoadmapStatusCompleted, result.Status)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

// ============================================================
// UpdateStep（ステップ所属チェック）
// ============================================================

func TestRoadmapUpdateStep_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10, Title: "Old Step"}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	updates := &model.RoadmapStep{Title: "New Step"}
	result, err := svc.UpdateStep(10, 5, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Step", result.Title)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStep_StepBelongsToDifferentRoadmap(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	// ステップは別のロードマップに所属
	step := &model.RoadmapStep{RoadmapID: 20, Title: "Step"}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)

	updates := &model.RoadmapStep{Title: "New"}
	result, err := svc.UpdateStep(10, 5, 1, updates)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStep_NotOwner(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	repo.On("FindByID", uint(10)).Return(roadmap, nil)

	updates := &model.RoadmapStep{Title: "New"}
	result, err := svc.UpdateStep(10, 5, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStep_RepoError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10, Title: "Old"}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(errors.New("db error"))

	_, err := svc.UpdateStep(10, 5, 1, &model.RoadmapStep{Title: "New"})
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// ステップ削除テスト
// ============================================================

func TestRoadmapDeleteStep_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("DeleteStep", uint(5)).Return(nil)

	err := svc.DeleteStep(10, 5, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapDeleteStep_StepBelongsToDifferentRoadmap(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 20}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)

	err := svc.DeleteStep(10, 5, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
	repo.AssertExpectations(t)
}

func TestRoadmapDeleteStep_RepoError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("DeleteStep", uint(5)).Return(errors.New("db error"))

	err := svc.DeleteStep(10, 5, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// ロードマップ削除テスト
// ============================================================

func TestRoadmapDelete_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapDelete_Forbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// ステップ完了状態テスト
// ============================================================

func TestRoadmapUpdateStepCompletion_Complete(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10, IsCompleted: false}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	result, err := svc.UpdateStepCompletion(10, 5, 1, true)
	assert.NoError(t, err)
	assert.True(t, result.IsCompleted)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStepCompletion_Uncomplete(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10, IsCompleted: true}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	result, err := svc.UpdateStepCompletion(10, 5, 1, false)
	assert.NoError(t, err)
	assert.False(t, result.IsCompleted)
	assert.Nil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

// ============================================================
// 存在しないロードマップテスト
// ============================================================

func TestRoadmapGetByID_NotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// テンプレート機能テスト
// ============================================================

func TestGetTemplates_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	templates := []model.Roadmap{
		{Title: "Webフロントエンド", IsTemplate: true, IsPublic: true, Category: "skill"},
		{Title: "バックエンド（Go）", IsTemplate: true, IsPublic: true, Category: "language"},
	}

	repo.On("GetTemplates").Return(templates, nil)

	result, err := svc.GetTemplates()
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.True(t, result[0].IsTemplate)
	assert.True(t, result[1].IsTemplate)
	repo.AssertExpectations(t)
}

func TestGetTemplates_Empty(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("GetTemplates").Return([]model.Roadmap{}, nil)

	result, err := svc.GetTemplates()
	assert.NoError(t, err)
	assert.Len(t, result, 0)
	repo.AssertExpectations(t)
}

func TestCreateFromTemplate_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	template := &model.Roadmap{Title: "Webフロントエンド", UserID: 0, IsTemplate: true, IsPublic: true}
	template.ID = 100

	copied := &model.Roadmap{Title: "Webフロントエンド", UserID: 5, IsTemplate: false}
	copied.ID = 200

	repo.On("FindByID", uint(100)).Return(template, nil)
	repo.On("CopyRoadmap", uint(100), uint(5)).Return(copied, nil)

	result, err := svc.CreateFromTemplate(100, 5)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(5), result.UserID)
	assert.Equal(t, "Webフロントエンド", result.Title)
	repo.AssertExpectations(t)
}

func TestCreateFromTemplate_NotTemplate(t *testing.T) {
	svc, repo := newTestRoadmapService()

	// IsTemplateがfalseの通常ロードマップ
	regular := &model.Roadmap{Title: "My Roadmap", UserID: 1, IsTemplate: false, IsPublic: true}
	regular.ID = 50

	repo.On("FindByID", uint(50)).Return(regular, nil)

	result, err := svc.CreateFromTemplate(50, 5)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestCreateFromTemplate_TemplateNotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.CreateFromTemplate(999, 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestSeedTemplates_CreatesTemplates(t *testing.T) {
	svc, repo := newTestRoadmapService()

	// 既存テンプレートなし
	repo.On("GetTemplates").Return([]model.Roadmap{}, nil)
	// 各テンプレート作成のモック
	repo.On("Create", mock.AnythingOfType("*model.Roadmap")).Return(nil)
	repo.On("CreateStep", mock.AnythingOfType("*model.RoadmapStep")).Return(nil)

	systemUserID := uint(1)
	err := svc.SeedTemplates(systemUserID)
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	// Createが5回以上呼ばれている（5テンプレート）
	createCalls := 0
	for _, call := range repo.Calls {
		if call.Method == "Create" {
			createCalls++
		}
	}
	assert.GreaterOrEqual(t, createCalls, 5)
}

func TestSeedTemplates_SetsUserID(t *testing.T) {
	svc, repo := newTestRoadmapService()

	// 既存テンプレートなし
	repo.On("GetTemplates").Return([]model.Roadmap{}, nil)
	repo.On("Create", mock.AnythingOfType("*model.Roadmap")).Return(nil)
	repo.On("CreateStep", mock.AnythingOfType("*model.RoadmapStep")).Return(nil)

	systemUserID := uint(42)
	err := svc.SeedTemplates(systemUserID)
	assert.NoError(t, err)

	// Createに渡されたRoadmapのUserIDが正しく設定されていることを確認
	for _, call := range repo.Calls {
		if call.Method == "Create" {
			roadmap := call.Arguments.Get(0).(*model.Roadmap)
			assert.Equal(t, systemUserID, roadmap.UserID, "テンプレートのUserIDがシステムユーザーIDと一致すること")
		}
	}
}

func TestSeedTemplates_SkipsIfAlreadyExist(t *testing.T) {
	svc, repo := newTestRoadmapService()

	// 既にテンプレートが存在する
	existing := []model.Roadmap{
		{Title: "Existing Template", IsTemplate: true},
	}
	repo.On("GetTemplates").Return(existing, nil)

	err := svc.SeedTemplates(uint(1))
	assert.NoError(t, err)
	repo.AssertExpectations(t)

	// Createは呼ばれない
	repo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestSeedTemplates_GetTemplatesError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("GetTemplates").Return([]model.Roadmap(nil), errors.New("db error"))

	err := svc.SeedTemplates(uint(1))
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	repo.AssertExpectations(t)
}

func TestSeedTemplates_CreateError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("GetTemplates").Return([]model.Roadmap{}, nil)
	repo.On("Create", mock.AnythingOfType("*model.Roadmap")).Return(errors.New("create error"))

	err := svc.SeedTemplates(uint(1))
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestSeedTemplates_CreateStepError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("GetTemplates").Return([]model.Roadmap{}, nil)
	repo.On("Create", mock.AnythingOfType("*model.Roadmap")).Return(nil)
	repo.On("CreateStep", mock.AnythingOfType("*model.RoadmapStep")).Return(errors.New("create step error"))

	err := svc.SeedTemplates(uint(1))
	assert.Error(t, err)
	assert.Equal(t, "create step error", err.Error())
	repo.AssertExpectations(t)
}

// ============================================================
// Create テスト
// ============================================================

func TestRoadmapCreate_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "New Roadmap", UserID: 1}
	repo.On("Create", roadmap).Return(nil)

	err := svc.Create(roadmap)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapCreate_WhitespaceTitle(t *testing.T) {
	svc, _ := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "   \t\n  ", UserID: 1}
	err := svc.Create(roadmap)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タイトルは必須です")
}

func TestRoadmapCreate_EmptyTitle(t *testing.T) {
	svc, _ := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "", UserID: 1}
	err := svc.Create(roadmap)
	assert.Error(t, err)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestRoadmapGetByUserID_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmaps := []model.Roadmap{
		{Title: "Roadmap 1", UserID: 1},
		{Title: "Roadmap 2", UserID: 1},
	}
	repo.On("GetByUserID", uint(1), 20, 0).Return(roadmaps, int64(2), nil)

	result, total, err := svc.GetByUserID(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestRoadmapGetByUserID_Empty(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("GetByUserID", uint(99), 20, 0).Return([]model.Roadmap{}, int64(0), nil)

	result, total, err := svc.GetByUserID(99, 20, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

func TestRoadmapGetByUserID_Page2(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("GetByUserID", uint(1), 10, 10).Return([]model.Roadmap{}, int64(15), nil)

	result, total, err := svc.GetByUserID(1, 10, 10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(15), total)
}

// ============================================================
// GetPublicRoadmaps テスト
// ============================================================

func TestRoadmapGetPublicRoadmaps_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmaps := []model.Roadmap{
		{Title: "Public 1", IsPublic: true},
		{Title: "Public 2", IsPublic: true},
	}
	repo.On("GetPublicRoadmaps", 10, 0).Return(roadmaps, int64(2), nil)

	result, total, err := svc.GetPublicRoadmaps(10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

// ============================================================
// GetStats テスト
// ============================================================

func TestRoadmapGetStats_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	stats := &model.RoadmapStats{
		TotalRoadmaps:     3,
		CompletedRoadmaps: 1,
	}
	repo.On("GetStats", uint(1)).Return(stats, nil)

	result, err := svc.GetStats(1)
	assert.NoError(t, err)
	assert.Equal(t, 3, result.TotalRoadmaps)
	assert.Equal(t, 1, result.CompletedRoadmaps)
	repo.AssertExpectations(t)
}

// ============================================================
// UpdateVisibility テスト
// ============================================================

func TestRoadmapUpdateVisibility_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{UserID: 1, IsPublic: false}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.UpdateVisibility(1, 1, true)
	assert.NoError(t, err)
	assert.True(t, result.IsPublic)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateVisibility_Forbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.UpdateVisibility(1, 999, true)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// CreateStep テスト
// ============================================================

func TestRoadmapCreateStep_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{Title: "New Step"}

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("CreateStep", step).Return(nil)

	err := svc.CreateStep(10, 1, step)
	assert.NoError(t, err)
	assert.Equal(t, uint(10), step.RoadmapID)
	repo.AssertExpectations(t)
}

func TestRoadmapCreateStep_Forbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	repo.On("FindByID", uint(10)).Return(roadmap, nil)

	step := &model.RoadmapStep{Title: "New Step"}
	err := svc.CreateStep(10, 999, step)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// ReorderSteps テスト
// ============================================================

func TestRoadmapReorderSteps_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	orders := []model.StepOrder{
		{StepID: 1, OrderIndex: 0},
		{StepID: 2, OrderIndex: 1},
	}

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("ReorderSteps", uint(10), orders).Return(nil)

	err := svc.ReorderSteps(10, 1, orders)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapReorderSteps_Forbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	repo.On("FindByID", uint(10)).Return(roadmap, nil)

	orders := []model.StepOrder{{StepID: 1, OrderIndex: 0}}
	err := svc.ReorderSteps(10, 999, orders)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdate_NotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	updates := &model.Roadmap{Title: "New"}
	result, err := svc.Update(99, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoadmapUpdate_RepoError(t *testing.T) {
	svc, repo := newTestRoadmapService()
	existing := &model.Roadmap{Title: "Old", UserID: 1, Status: model.RoadmapStatusActive}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))
	updates := &model.Roadmap{Title: "New"}
	result, err := svc.Update(1, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoadmapUpdateVisibility_NotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	result, err := svc.UpdateVisibility(99, 1, true)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoadmapUpdateVisibility_RepoError(t *testing.T) {
	svc, repo := newTestRoadmapService()
	existing := &model.Roadmap{UserID: 1, IsPublic: false}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(errors.New("db error"))
	result, err := svc.UpdateVisibility(1, 1, true)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoadmapUpdateStepCompletion_Forbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10
	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	result, err := svc.UpdateStepCompletion(10, 5, 999, true)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
}

func TestRoadmapUpdateStepCompletion_StepNotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10
	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(99)).Return(nil, errors.New("not found"))
	result, err := svc.UpdateStepCompletion(10, 99, 1, true)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoadmapUpdateStepCompletion_WrongRoadmap(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10
	step := &model.RoadmapStep{RoadmapID: 20}
	step.ID = 5
	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	result, err := svc.UpdateStepCompletion(10, 5, 1, true)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, result)
}

func TestRoadmapUpdateStepCompletion_RepoError(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10
	step := &model.RoadmapStep{RoadmapID: 10, IsCompleted: false}
	step.ID = 5
	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(errors.New("db error"))
	result, err := svc.UpdateStepCompletion(10, 5, 1, true)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRoadmapDeleteStep_NotOwner(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10
	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	err := svc.DeleteStep(10, 5, 999)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestRoadmapDeleteStep_StepNotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10
	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.DeleteStep(10, 99, 1)
	assert.Error(t, err)
}

func TestRoadmapDelete_NotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.Delete(99, 1)
	assert.Error(t, err)
}

func TestRoadmapCreateStep_NotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	step := &model.RoadmapStep{Title: "New Step"}
	err := svc.CreateStep(99, 1, step)
	assert.Error(t, err)
}

func TestRoadmapUpdate_WithAllFields(t *testing.T) {
	svc, repo := newTestRoadmapService()
	existing := &model.Roadmap{Title: "Old", UserID: 1, Status: model.RoadmapStatusActive}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)
	updates := &model.Roadmap{Title: "New", Description: "Desc", Category: model.RoadmapCategorySkill}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New", result.Title)
	assert.Equal(t, "Desc", result.Description)
	assert.Equal(t, model.RoadmapCategorySkill, result.Category)
}

// ============================================================
// UpdateStep 追加テスト
// ============================================================

func TestUpdateStep_StepBelongsToDifferentRoadmap(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 1

	// ステップが別ロードマップ（ID=2）に属している
	step := &model.RoadmapStep{Title: "Step", RoadmapID: 2}
	step.ID = 10

	repo.On("FindByID", uint(1)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(10)).Return(step, nil)

	updates := &model.RoadmapStep{Title: "Updated"}
	_, err := svc.UpdateStep(1, 10, 1, updates)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateStep_UpdateStepRepoError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 1

	step := &model.RoadmapStep{Title: "Step", RoadmapID: 1}
	step.ID = 10

	repo.On("FindByID", uint(1)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(10)).Return(step, nil)
	repo.On("UpdateStep", step).Return(errors.New("db error"))

	updates := &model.RoadmapStep{Title: "Updated"}
	_, err := svc.UpdateStep(1, 10, 1, updates)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestUpdateStep_FindStepError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 1

	repo.On("FindByID", uint(1)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(10)).Return(nil, errors.New("step not found"))

	updates := &model.RoadmapStep{Title: "Updated"}
	_, err := svc.UpdateStep(1, 10, 1, updates)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// ReorderSteps 追加テスト
// ============================================================

func TestRoadmapReorderSteps_NotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.ReorderSteps(99, 1, []model.StepOrder{})
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// CopyRoadmap 追加テスト
// ============================================================

func TestCopyRoadmap_NotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	_, err := svc.CopyRoadmap(99, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapDeleteStep_FindStepByIDError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.DeleteStep(10, 99, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapDeleteStep_FindByIDError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.DeleteStep(99, 5, 1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStep_AllFields(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10, Title: "Old", Description: "Old Desc", ResourceURL: "https://old.example.com"}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	updates := &model.RoadmapStep{
		Title:       "New Step",
		Description: "New Description",
		ResourceURL: "https://new.example.com",
	}
	result, err := svc.UpdateStep(10, 5, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Step", result.Title)
	assert.Equal(t, "New Description", result.Description)
	assert.Equal(t, "https://new.example.com", result.ResourceURL)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStep_FindByIDError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	updates := &model.RoadmapStep{Title: "New"}
	result, err := svc.UpdateStep(99, 5, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStepCompletion_FindByIDError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.UpdateStepCompletion(99, 5, 1, true)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByStatus テスト
// ============================================================

func TestRoadmapGetByStatus_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmaps := []model.Roadmap{
		{Title: "Active Roadmap", UserID: 1, Status: model.RoadmapStatusActive},
	}
	repo.On("GetByStatus", uint(1), "active").Return(roadmaps, nil)

	result, err := svc.GetByStatus(1, "active")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Active Roadmap", result[0].Title)
	repo.AssertExpectations(t)
}

func TestRoadmapGetByStatus_EmptyResult(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("GetByStatus", uint(1), "completed").Return([]model.Roadmap{}, nil)

	result, err := svc.GetByStatus(1, "completed")
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestRoadmapGetByStatus_InvalidStatus(t *testing.T) {
	svc, _ := newTestRoadmapService()

	result, err := svc.GetByStatus(1, "invalid")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "無効なステータス")
}

func TestRoadmapGetByStatus_RepoError(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("GetByStatus", uint(1), "active").Return([]model.Roadmap{}, errors.New("db error"))

	result, err := svc.GetByStatus(1, "active")
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 空白バイパス脆弱性テスト（Roadmap）
// ============================================================

func TestRoadmapUpdate_WhitespaceTitle(t *testing.T) {
	svc, repo := newTestRoadmapService()
	existing := &model.Roadmap{Title: "Original Title", Description: "Desc", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.Roadmap{Title: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Title", result.Title)
}

func TestRoadmapUpdate_WhitespaceDescription(t *testing.T) {
	svc, repo := newTestRoadmapService()
	existing := &model.Roadmap{Title: "Title", Description: "Original Desc", UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.Roadmap{Description: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Desc", result.Description)
}

func TestRoadmapUpdate_WhitespaceCategory(t *testing.T) {
	svc, repo := newTestRoadmapService()
	existing := &model.Roadmap{Title: "Title", Category: model.RoadmapCategoryLanguage, UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.Roadmap{Category: "   "})
	assert.NoError(t, err)
	assert.Equal(t, model.RoadmapCategoryLanguage, result.Category)
}

func TestRoadmapUpdate_WhitespaceStatus(t *testing.T) {
	svc, repo := newTestRoadmapService()
	existing := &model.Roadmap{Title: "Title", Status: model.RoadmapStatusActive, UserID: 1}
	existing.ID = 1
	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.Update(1, 1, &model.Roadmap{Status: "   "})
	assert.NoError(t, err)
	assert.Equal(t, model.RoadmapStatusActive, result.Status)
}

// ============================================================
// 空白バイパス脆弱性テスト（RoadmapStep）
// ============================================================

func TestRoadmapUpdateStep_WhitespaceTitle(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 1
	step := &model.RoadmapStep{RoadmapID: 1, Title: "Original Step"}
	step.ID = 10
	repo.On("FindByID", uint(1)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(10)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	result, err := svc.UpdateStep(1, 10, 1, &model.RoadmapStep{Title: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Step", result.Title)
}

func TestRoadmapUpdateStep_WhitespaceDescription(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 1
	step := &model.RoadmapStep{RoadmapID: 1, Description: "Original Desc"}
	step.ID = 10
	repo.On("FindByID", uint(1)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(10)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	result, err := svc.UpdateStep(1, 10, 1, &model.RoadmapStep{Description: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "Original Desc", result.Description)
}

func TestRoadmapUpdateStep_WhitespaceResourceURL(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 1
	step := &model.RoadmapStep{RoadmapID: 1, ResourceURL: "https://example.com"}
	step.ID = 10
	repo.On("FindByID", uint(1)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(10)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	result, err := svc.UpdateStep(1, 10, 1, &model.RoadmapStep{ResourceURL: "   "})
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", result.ResourceURL)
}

func TestRoadmapUpdateStep_TrimsPaddedTitle(t *testing.T) {
	svc, repo := newTestRoadmapService()
	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 1
	step := &model.RoadmapStep{RoadmapID: 1, Title: "Original"}
	step.ID = 10
	repo.On("FindByID", uint(1)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(10)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	result, err := svc.UpdateStep(1, 10, 1, &model.RoadmapStep{Title: "  New Title  "})
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
}
