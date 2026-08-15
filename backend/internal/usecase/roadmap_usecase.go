package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// errRoadmapNotFound は port が「不在」を表す nil を返したときに返すエラー。
// DomainError ではないため handler では 500 になり、リポジトリの生エラーが
// そのまま返っていた移行前の挙動と一致する。
var errRoadmapNotFound = errors.New("ロードマップが見つかりません")

// validRoadmapStatuses は有効なロードマップステータスのマップ。
var validRoadmapStatuses = map[string]bool{
	string(model.RoadmapStatusActive):    true,
	string(model.RoadmapStatusCompleted): true,
}

// ownerOfRoadmap は所有権判定に使う所有者 ID の取り出し。
func ownerOfRoadmap(r *model.Roadmap) uint { return r.UserID }

// findOwnedRoadmapStep はロードマップの所有者であることを確認したうえでステップを取得する。
// ステップが別のロードマップに属している場合は 400 を返す。
func findOwnedRoadmapStep(ctx context.Context, roadmaps repository.RoadmapRepository, roadmapID, stepID, userID uint) (*model.RoadmapStep, error) {
	if _, err := ensureOwner(ctx, roadmaps.FindByID, roadmapID, userID, ownerOfRoadmap); err != nil {
		return nil, err
	}
	step, err := roadmaps.FindStepByID(ctx, stepID)
	if err != nil {
		return nil, err
	}
	if step == nil {
		return nil, errRoadmapNotFound
	}
	if step.RoadmapID != roadmapID {
		return nil, domain.ErrBadRequest
	}
	return step, nil
}

// CreateRoadmapUseCase はロードマップを作成する。
type CreateRoadmapUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewCreateRoadmapUseCase は CreateRoadmapUseCase を生成する。
func NewCreateRoadmapUseCase(roadmaps repository.RoadmapRepository) *CreateRoadmapUseCase {
	return &CreateRoadmapUseCase{roadmaps: roadmaps}
}

// Execute はタイトルと説明を検証したうえでロードマップを作成する。
func (uc *CreateRoadmapUseCase) Execute(ctx context.Context, roadmap *model.Roadmap) error {
	if err := domain.ValidateStringLength(roadmap.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(roadmap.Description, 0, 1000, "説明"); err != nil {
		return err
	}
	roadmap.Title = strings.TrimSpace(roadmap.Title)
	roadmap.Description = strings.TrimSpace(roadmap.Description)
	return uc.roadmaps.Create(ctx, roadmap)
}

// GetRoadmapUseCase はロードマップを 1 件取得する。
type GetRoadmapUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewGetRoadmapUseCase は GetRoadmapUseCase を生成する。
func NewGetRoadmapUseCase(roadmaps repository.RoadmapRepository) *GetRoadmapUseCase {
	return &GetRoadmapUseCase{roadmaps: roadmaps}
}

// Execute はロードマップを返す。非公開のものは所有者しか取得できない。
func (uc *GetRoadmapUseCase) Execute(ctx context.Context, id, userID uint) (*model.Roadmap, error) {
	roadmap, err := uc.roadmaps.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if roadmap == nil {
		return nil, errRoadmapNotFound
	}
	if roadmap.UserID != userID && !roadmap.IsPublic {
		return nil, domain.ErrForbidden
	}
	return roadmap, nil
}

// ListRoadmapsByUserUseCase は指定ユーザーのロードマップ一覧を取得する。
type ListRoadmapsByUserUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewListRoadmapsByUserUseCase は ListRoadmapsByUserUseCase を生成する。
func NewListRoadmapsByUserUseCase(roadmaps repository.RoadmapRepository) *ListRoadmapsByUserUseCase {
	return &ListRoadmapsByUserUseCase{roadmaps: roadmaps}
}

// Execute は指定ユーザーのロードマップを新しい順で返す。
func (uc *ListRoadmapsByUserUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	return uc.roadmaps.GetByUserID(ctx, userID, limit, offset)
}

// ListRoadmapsByStatusUseCase はロードマップをステータスで絞り込んで取得する。
type ListRoadmapsByStatusUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewListRoadmapsByStatusUseCase は ListRoadmapsByStatusUseCase を生成する。
func NewListRoadmapsByStatusUseCase(roadmaps repository.RoadmapRepository) *ListRoadmapsByStatusUseCase {
	return &ListRoadmapsByStatusUseCase{roadmaps: roadmaps}
}

// Execute は指定ステータスのロードマップを返す。未知のステータスは 400。
func (uc *ListRoadmapsByStatusUseCase) Execute(ctx context.Context, userID uint, status string) ([]model.Roadmap, error) {
	if !validRoadmapStatuses[status] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なステータスです", nil)
	}
	return uc.roadmaps.GetByStatus(ctx, userID, status)
}

// ListPublicRoadmapsUseCase は公開ロードマップ一覧を取得する。
type ListPublicRoadmapsUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewListPublicRoadmapsUseCase は ListPublicRoadmapsUseCase を生成する。
func NewListPublicRoadmapsUseCase(roadmaps repository.RoadmapRepository) *ListPublicRoadmapsUseCase {
	return &ListPublicRoadmapsUseCase{roadmaps: roadmaps}
}

// Execute は公開ロードマップを新しい順で返す。
func (uc *ListPublicRoadmapsUseCase) Execute(ctx context.Context, limit, offset int) ([]model.Roadmap, int64, error) {
	return uc.roadmaps.GetPublicRoadmaps(ctx, limit, offset)
}

// UpdateRoadmapUseCase はロードマップを更新する。
type UpdateRoadmapUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewUpdateRoadmapUseCase は UpdateRoadmapUseCase を生成する。
func NewUpdateRoadmapUseCase(roadmaps repository.RoadmapRepository) *UpdateRoadmapUseCase {
	return &UpdateRoadmapUseCase{roadmaps: roadmaps}
}

// Execute はタイトル・説明・カテゴリ・ステータスを部分更新する。所有者のみ。
// 前後の空白を除いて空になった項目は「変更なし」として扱う。
func (uc *UpdateRoadmapUseCase) Execute(ctx context.Context, id, userID uint, updates *model.Roadmap) (*model.Roadmap, error) {
	roadmap, err := ensureOwner(ctx, uc.roadmaps.FindByID, id, userID, ownerOfRoadmap)
	if err != nil {
		return nil, err
	}

	if title := strings.TrimSpace(updates.Title); title != "" {
		if err := domain.ValidateStringLength(title, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		roadmap.Title = title
	}
	if desc := strings.TrimSpace(updates.Description); desc != "" {
		if err := domain.ValidateStringLength(desc, 1, 1000, "説明"); err != nil {
			return nil, err
		}
		roadmap.Description = desc
	}
	if strings.TrimSpace(string(updates.Category)) != "" {
		roadmap.Category = updates.Category
	}
	if strings.TrimSpace(string(updates.Status)) != "" {
		roadmap.Status = updates.Status
		if roadmap.Status == model.RoadmapStatusCompleted && roadmap.CompletedAt == nil {
			now := time.Now()
			roadmap.CompletedAt = &now
		}
	}

	if err := uc.roadmaps.Update(ctx, roadmap); err != nil {
		return nil, err
	}
	return roadmap, nil
}

// UpdateRoadmapVisibilityUseCase はロードマップの公開/非公開を切り替える。
type UpdateRoadmapVisibilityUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewUpdateRoadmapVisibilityUseCase は UpdateRoadmapVisibilityUseCase を生成する。
func NewUpdateRoadmapVisibilityUseCase(roadmaps repository.RoadmapRepository) *UpdateRoadmapVisibilityUseCase {
	return &UpdateRoadmapVisibilityUseCase{roadmaps: roadmaps}
}

// Execute は公開状態を更新する。所有者のみ。
func (uc *UpdateRoadmapVisibilityUseCase) Execute(ctx context.Context, id, userID uint, isPublic bool) (*model.Roadmap, error) {
	roadmap, err := ensureOwner(ctx, uc.roadmaps.FindByID, id, userID, ownerOfRoadmap)
	if err != nil {
		return nil, err
	}

	roadmap.IsPublic = isPublic
	if err := uc.roadmaps.Update(ctx, roadmap); err != nil {
		return nil, err
	}
	return roadmap, nil
}

// DeleteRoadmapUseCase はロードマップを削除する。
type DeleteRoadmapUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewDeleteRoadmapUseCase は DeleteRoadmapUseCase を生成する。
func NewDeleteRoadmapUseCase(roadmaps repository.RoadmapRepository) *DeleteRoadmapUseCase {
	return &DeleteRoadmapUseCase{roadmaps: roadmaps}
}

// Execute はロードマップを削除する。所有者のみ。
func (uc *DeleteRoadmapUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.roadmaps.FindByID, id, userID, ownerOfRoadmap); err != nil {
		return err
	}
	return uc.roadmaps.Delete(ctx, id)
}

// CopyRoadmapUseCase は公開ロードマップを自分用に複製する。
type CopyRoadmapUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewCopyRoadmapUseCase は CopyRoadmapUseCase を生成する。
func NewCopyRoadmapUseCase(roadmaps repository.RoadmapRepository) *CopyRoadmapUseCase {
	return &CopyRoadmapUseCase{roadmaps: roadmaps}
}

// Execute はロードマップを複製する。非公開のものは所有者しか複製できない。
func (uc *CopyRoadmapUseCase) Execute(ctx context.Context, roadmapID, userID uint) (*model.Roadmap, error) {
	original, err := uc.roadmaps.FindByID(ctx, roadmapID)
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, errRoadmapNotFound
	}
	if !original.IsPublic && original.UserID != userID {
		return nil, domain.ErrForbidden
	}

	return copyRoadmapOrNotFound(ctx, uc.roadmaps, roadmapID, userID)
}

// copyRoadmapOrNotFound は複製を実行し、複製元が存在確認との間に消えていた場合は
// 不在として扱う（nil を成功として返さない）。
func copyRoadmapOrNotFound(ctx context.Context, roadmaps repository.RoadmapRepository, originalID, userID uint) (*model.Roadmap, error) {
	copied, err := roadmaps.CopyRoadmap(ctx, originalID, userID)
	if err != nil {
		return nil, err
	}
	if copied == nil {
		return nil, errRoadmapNotFound
	}
	return copied, nil
}

// ListRoadmapTemplatesUseCase はテンプレート一覧を取得する。
type ListRoadmapTemplatesUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewListRoadmapTemplatesUseCase は ListRoadmapTemplatesUseCase を生成する。
func NewListRoadmapTemplatesUseCase(roadmaps repository.RoadmapRepository) *ListRoadmapTemplatesUseCase {
	return &ListRoadmapTemplatesUseCase{roadmaps: roadmaps}
}

// Execute はテンプレートのロードマップをステップ付きで返す。
func (uc *ListRoadmapTemplatesUseCase) Execute(ctx context.Context) ([]model.Roadmap, error) {
	return uc.roadmaps.GetTemplates(ctx)
}

// CreateRoadmapFromTemplateUseCase はテンプレートから自分用のロードマップを作る。
type CreateRoadmapFromTemplateUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewCreateRoadmapFromTemplateUseCase は CreateRoadmapFromTemplateUseCase を生成する。
func NewCreateRoadmapFromTemplateUseCase(roadmaps repository.RoadmapRepository) *CreateRoadmapFromTemplateUseCase {
	return &CreateRoadmapFromTemplateUseCase{roadmaps: roadmaps}
}

// Execute はテンプレートを複製する。テンプレートでないロードマップを指定した場合は 400。
func (uc *CreateRoadmapFromTemplateUseCase) Execute(ctx context.Context, templateID, userID uint) (*model.Roadmap, error) {
	template, err := uc.roadmaps.FindByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, errRoadmapNotFound
	}
	if !template.IsTemplate {
		return nil, domain.ErrBadRequest
	}
	return copyRoadmapOrNotFound(ctx, uc.roadmaps, templateID, userID)
}

// CreateRoadmapStepUseCase はロードマップにステップを追加する。
type CreateRoadmapStepUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewCreateRoadmapStepUseCase は CreateRoadmapStepUseCase を生成する。
func NewCreateRoadmapStepUseCase(roadmaps repository.RoadmapRepository) *CreateRoadmapStepUseCase {
	return &CreateRoadmapStepUseCase{roadmaps: roadmaps}
}

// Execute はステップを追加する。所有者のみ。
func (uc *CreateRoadmapStepUseCase) Execute(ctx context.Context, roadmapID, userID uint, step *model.RoadmapStep) error {
	if _, err := ensureOwner(ctx, uc.roadmaps.FindByID, roadmapID, userID, ownerOfRoadmap); err != nil {
		return err
	}
	step.RoadmapID = roadmapID
	return uc.roadmaps.CreateStep(ctx, step)
}

// UpdateRoadmapStepUseCase はステップを更新する。
type UpdateRoadmapStepUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewUpdateRoadmapStepUseCase は UpdateRoadmapStepUseCase を生成する。
func NewUpdateRoadmapStepUseCase(roadmaps repository.RoadmapRepository) *UpdateRoadmapStepUseCase {
	return &UpdateRoadmapStepUseCase{roadmaps: roadmaps}
}

// Execute はステップのタイトル・説明・リソース URL を部分更新する。所有者のみ。
func (uc *UpdateRoadmapStepUseCase) Execute(ctx context.Context, roadmapID, stepID, userID uint, updates *model.RoadmapStep) (*model.RoadmapStep, error) {
	step, err := findOwnedRoadmapStep(ctx, uc.roadmaps, roadmapID, stepID, userID)
	if err != nil {
		return nil, err
	}

	if title := strings.TrimSpace(updates.Title); title != "" {
		if err := domain.ValidateStringLength(title, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		step.Title = title
	}
	if desc := strings.TrimSpace(updates.Description); desc != "" {
		if err := domain.ValidateStringLength(desc, 1, 1000, "説明"); err != nil {
			return nil, err
		}
		step.Description = desc
	}
	if url := strings.TrimSpace(updates.ResourceURL); url != "" {
		if err := domain.ValidateStringLength(url, 1, 500, "リソースURL"); err != nil {
			return nil, err
		}
		step.ResourceURL = url
	}

	if err := uc.roadmaps.UpdateStep(ctx, step); err != nil {
		return nil, err
	}
	return step, nil
}

// UpdateRoadmapStepCompletionUseCase はステップの完了状態を切り替える。
type UpdateRoadmapStepCompletionUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewUpdateRoadmapStepCompletionUseCase は UpdateRoadmapStepCompletionUseCase を生成する。
func NewUpdateRoadmapStepCompletionUseCase(roadmaps repository.RoadmapRepository) *UpdateRoadmapStepCompletionUseCase {
	return &UpdateRoadmapStepCompletionUseCase{roadmaps: roadmaps}
}

// Execute は完了状態を更新する。完了にしたときだけ完了日時を記録し、
// 未完了に戻したときは完了日時を消す。所有者のみ。
func (uc *UpdateRoadmapStepCompletionUseCase) Execute(ctx context.Context, roadmapID, stepID, userID uint, isCompleted bool) (*model.RoadmapStep, error) {
	step, err := findOwnedRoadmapStep(ctx, uc.roadmaps, roadmapID, stepID, userID)
	if err != nil {
		return nil, err
	}

	step.IsCompleted = isCompleted
	switch {
	case isCompleted && step.CompletedAt == nil:
		now := time.Now()
		step.CompletedAt = &now
	case !isCompleted:
		step.CompletedAt = nil
	}

	if err := uc.roadmaps.UpdateStep(ctx, step); err != nil {
		return nil, err
	}
	return step, nil
}

// BatchCompleteRoadmapStepsUseCase は複数ステップをまとめて完了にする。
type BatchCompleteRoadmapStepsUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewBatchCompleteRoadmapStepsUseCase は BatchCompleteRoadmapStepsUseCase を生成する。
func NewBatchCompleteRoadmapStepsUseCase(roadmaps repository.RoadmapRepository) *BatchCompleteRoadmapStepsUseCase {
	return &BatchCompleteRoadmapStepsUseCase{roadmaps: roadmaps}
}

// Execute は指定ステップをまとめて完了にする。所有者のみ。
// 別のロードマップのステップが混ざっていたらその時点で中断する（それまでの更新は残る）。
func (uc *BatchCompleteRoadmapStepsUseCase) Execute(ctx context.Context, roadmapID, userID uint, stepIDs []uint) (*model.Roadmap, error) {
	if _, err := ensureOwner(ctx, uc.roadmaps.FindByID, roadmapID, userID, ownerOfRoadmap); err != nil {
		return nil, err
	}

	now := time.Now()
	for _, stepID := range stepIDs {
		step, err := uc.roadmaps.FindStepByID(ctx, stepID)
		if err != nil {
			return nil, err
		}
		if step == nil {
			return nil, errRoadmapNotFound
		}
		if step.RoadmapID != roadmapID {
			return nil, domain.ErrBadRequest
		}
		if step.IsCompleted {
			continue
		}
		step.IsCompleted = true
		step.CompletedAt = &now
		if err := uc.roadmaps.UpdateStep(ctx, step); err != nil {
			return nil, err
		}
	}

	roadmap, err := uc.roadmaps.FindByID(ctx, roadmapID)
	if err != nil {
		return nil, err
	}
	if roadmap == nil {
		return nil, errRoadmapNotFound
	}
	return roadmap, nil
}

// DeleteRoadmapStepUseCase はステップを削除する。
type DeleteRoadmapStepUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewDeleteRoadmapStepUseCase は DeleteRoadmapStepUseCase を生成する。
func NewDeleteRoadmapStepUseCase(roadmaps repository.RoadmapRepository) *DeleteRoadmapStepUseCase {
	return &DeleteRoadmapStepUseCase{roadmaps: roadmaps}
}

// Execute はステップを削除する。所有者のみ。
func (uc *DeleteRoadmapStepUseCase) Execute(ctx context.Context, roadmapID, stepID, userID uint) error {
	if _, err := findOwnedRoadmapStep(ctx, uc.roadmaps, roadmapID, stepID, userID); err != nil {
		return err
	}
	return uc.roadmaps.DeleteStep(ctx, stepID)
}

// ReorderRoadmapStepsUseCase はステップの表示順序を変更する。
type ReorderRoadmapStepsUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewReorderRoadmapStepsUseCase は ReorderRoadmapStepsUseCase を生成する。
func NewReorderRoadmapStepsUseCase(roadmaps repository.RoadmapRepository) *ReorderRoadmapStepsUseCase {
	return &ReorderRoadmapStepsUseCase{roadmaps: roadmaps}
}

// Execute はステップの表示順序をまとめて更新する。所有者のみ。
func (uc *ReorderRoadmapStepsUseCase) Execute(ctx context.Context, roadmapID, userID uint, orders []model.StepOrder) error {
	if _, err := ensureOwner(ctx, uc.roadmaps.FindByID, roadmapID, userID, ownerOfRoadmap); err != nil {
		return err
	}
	return uc.roadmaps.ReorderSteps(ctx, roadmapID, orders)
}

// CountRoadmapsUseCase はロードマップ数を取得する。
type CountRoadmapsUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewCountRoadmapsUseCase は CountRoadmapsUseCase を生成する。
func NewCountRoadmapsUseCase(roadmaps repository.RoadmapRepository) *CountRoadmapsUseCase {
	return &CountRoadmapsUseCase{roadmaps: roadmaps}
}

// Execute は指定ユーザーのロードマップ総数を返す。
func (uc *CountRoadmapsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.roadmaps.CountByUserID(ctx, userID)
}

// SeedRoadmapTemplatesUseCase はプリセットのテンプレートを初期登録する。
type SeedRoadmapTemplatesUseCase struct {
	roadmaps repository.RoadmapRepository
}

// NewSeedRoadmapTemplatesUseCase は SeedRoadmapTemplatesUseCase を生成する。
func NewSeedRoadmapTemplatesUseCase(roadmaps repository.RoadmapRepository) *SeedRoadmapTemplatesUseCase {
	return &SeedRoadmapTemplatesUseCase{roadmaps: roadmaps}
}

// Execute はテンプレートを初期登録する。既にテンプレートがあれば何もしない。
// userID には外部キー制約を満たすためのシステムユーザー ID を渡す。
func (uc *SeedRoadmapTemplatesUseCase) Execute(ctx context.Context, userID uint) error {
	existing, err := uc.roadmaps.GetTemplates(ctx)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	for _, tmpl := range presetRoadmapTemplates {
		roadmap := &model.Roadmap{
			UserID:      userID,
			Title:       tmpl.Title,
			Description: tmpl.Description,
			Category:    tmpl.Category,
			IsPublic:    true,
			IsTemplate:  true,
			StepCount:   len(tmpl.Steps),
			Status:      model.RoadmapStatusActive,
		}
		if err := uc.roadmaps.Create(ctx, roadmap); err != nil {
			return err
		}
		for i, stepDef := range tmpl.Steps {
			step := &model.RoadmapStep{
				RoadmapID:   roadmap.ID,
				Title:       stepDef.Title,
				Description: stepDef.Description,
				OrderIndex:  i,
				ResourceURL: stepDef.ResourceURL,
			}
			if err := uc.roadmaps.CreateStep(ctx, step); err != nil {
				return err
			}
		}
	}
	return nil
}
