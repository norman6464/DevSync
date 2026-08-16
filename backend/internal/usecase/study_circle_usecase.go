package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// validStudyCircleStatuses は有効なスタディサークルステータスのマップ。
var validStudyCircleStatuses = map[string]bool{
	string(model.StudyCircleStatusActive):    true,
	string(model.StudyCircleStatusCompleted): true,
	string(model.StudyCircleStatusArchived):  true,
}

// validStudyCircleMemberRoles は有効なメンバー役割のマップ。
var validStudyCircleMemberRoles = map[string]bool{
	string(model.StudyCircleRoleOwner):  true,
	string(model.StudyCircleRoleMember): true,
}

// findOwnedStudyCircle はサークルを取得し、userID がオーナーであることを検証して返す。
// 取得できなかった場合は不在・DB 障害のいずれも 404 として扱う（移行前の挙動を維持している）。
func findOwnedStudyCircle(ctx context.Context, circles repository.StudyCircleRepository, circleID, userID uint) (*model.StudyCircle, error) {
	circle, err := circles.FindByID(ctx, circleID)
	if err != nil || circle == nil {
		return nil, domain.ErrNotFound
	}
	if circle.OwnerID != userID {
		return nil, domain.ErrForbidden
	}
	return circle, nil
}

// findOwnedStudyCircleStep はサークルのオーナーであることを確認したうえでステップを取得する。
// 別サークルのステップ ID を指定された場合も 404 を返す。
func findOwnedStudyCircleStep(ctx context.Context, circles repository.StudyCircleRepository, circleID, stepID, userID uint) (*model.StudyCircleStep, error) {
	if _, err := findOwnedStudyCircle(ctx, circles, circleID, userID); err != nil {
		return nil, err
	}
	step, err := circles.FindStepByID(ctx, stepID)
	if err != nil || step == nil {
		return nil, domain.ErrNotFound
	}
	if step.CircleID != circleID {
		return nil, domain.ErrNotFound
	}
	return step, nil
}

// requireStudyCircleMember はユーザーがサークルのメンバーであることを検証する。
// メンバーでなければ 403 を返す。
func requireStudyCircleMember(ctx context.Context, circles repository.StudyCircleRepository, circleID, userID uint) error {
	isMember, err := circles.IsMember(ctx, circleID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return domain.ErrForbidden
	}
	return nil
}

// CreateStudyCircleUseCase はスタディサークルを作成する。
type CreateStudyCircleUseCase struct {
	circles repository.StudyCircleRepository
}

// NewCreateStudyCircleUseCase は CreateStudyCircleUseCase を生成する。
func NewCreateStudyCircleUseCase(circles repository.StudyCircleRepository) *CreateStudyCircleUseCase {
	return &CreateStudyCircleUseCase{circles: circles}
}

// Execute はサークルを作成し、オーナーをメンバーとして自動追加する。
// 続けて memberIDs のユーザーを追加するが、招待の失敗はサークル作成を巻き戻さない。
func (uc *CreateStudyCircleUseCase) Execute(ctx context.Context, circle *model.StudyCircle, memberIDs []uint) error {
	if err := domain.ValidateStringLength(circle.Name, 1, 100, "サークル名"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(circle.Topic, 0, 200, "トピック"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(circle.Description, 0, 1000, "説明"); err != nil {
		return err
	}
	circle.Name = strings.TrimSpace(circle.Name)
	circle.Topic = strings.TrimSpace(circle.Topic)
	circle.Description = strings.TrimSpace(circle.Description)
	if circle.MaxMembers < 3 || circle.MaxMembers > 10 {
		circle.MaxMembers = 5
	}
	circle.Status = model.StudyCircleStatusActive

	if err := uc.circles.CreateWithOwner(ctx, circle); err != nil {
		return err
	}

	for _, memberID := range memberIDs {
		if memberID != circle.OwnerID {
			_ = uc.circles.AddMember(ctx, circle.ID, memberID, model.StudyCircleRoleMember)
		}
	}

	return nil
}

// ListMyStudyCirclesUseCase は自分が参加しているサークル一覧を取得する。
type ListMyStudyCirclesUseCase struct {
	circles repository.StudyCircleRepository
}

// NewListMyStudyCirclesUseCase は ListMyStudyCirclesUseCase を生成する。
func NewListMyStudyCirclesUseCase(circles repository.StudyCircleRepository) *ListMyStudyCirclesUseCase {
	return &ListMyStudyCirclesUseCase{circles: circles}
}

// Execute は参加サークル一覧を総件数付きで返す。
func (uc *ListMyStudyCirclesUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.StudyCircle, int64, error) {
	return uc.circles.FindByUserID(ctx, userID, limit, offset)
}

// ListStudyCirclesByStatusUseCase は参加サークルをステータスで絞り込んで取得する。
type ListStudyCirclesByStatusUseCase struct {
	circles repository.StudyCircleRepository
}

// NewListStudyCirclesByStatusUseCase は ListStudyCirclesByStatusUseCase を生成する。
func NewListStudyCirclesByStatusUseCase(circles repository.StudyCircleRepository) *ListStudyCirclesByStatusUseCase {
	return &ListStudyCirclesByStatusUseCase{circles: circles}
}

// Execute は指定ステータスの参加サークルを返す。未知のステータスは 400。
func (uc *ListStudyCirclesByStatusUseCase) Execute(ctx context.Context, userID uint, status string) ([]model.StudyCircle, error) {
	if !validStudyCircleStatuses[status] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なステータスです", nil)
	}
	return uc.circles.GetByStatus(ctx, userID, status)
}

// GetStudyCircleUseCase はサークル詳細を取得する。
type GetStudyCircleUseCase struct {
	circles repository.StudyCircleRepository
}

// NewGetStudyCircleUseCase は GetStudyCircleUseCase を生成する。
func NewGetStudyCircleUseCase(circles repository.StudyCircleRepository) *GetStudyCircleUseCase {
	return &GetStudyCircleUseCase{circles: circles}
}

// Execute はサークル詳細を返す。メンバーのみアクセスできる。
func (uc *GetStudyCircleUseCase) Execute(ctx context.Context, id, userID uint) (*model.StudyCircle, error) {
	circle, err := uc.circles.FindByID(ctx, id)
	if err != nil || circle == nil {
		return nil, domain.ErrNotFound
	}
	if err := requireStudyCircleMember(ctx, uc.circles, id, userID); err != nil {
		return nil, err
	}
	return circle, nil
}

// UpdateStudyCircleUseCase はサークル情報を更新する。
type UpdateStudyCircleUseCase struct {
	circles repository.StudyCircleRepository
}

// NewUpdateStudyCircleUseCase は UpdateStudyCircleUseCase を生成する。
func NewUpdateStudyCircleUseCase(circles repository.StudyCircleRepository) *UpdateStudyCircleUseCase {
	return &UpdateStudyCircleUseCase{circles: circles}
}

// Execute はサークルの名前・トピック・説明を部分更新する。オーナーのみ。
func (uc *UpdateStudyCircleUseCase) Execute(ctx context.Context, id, userID uint, name, topic, description *string) (*model.StudyCircle, error) {
	circle, err := findOwnedStudyCircle(ctx, uc.circles, id, userID)
	if err != nil {
		return nil, err
	}

	if name != nil {
		if err := domain.ValidateStringLength(*name, 1, 100, "サークル名"); err != nil {
			return nil, err
		}
		circle.Name = strings.TrimSpace(*name)
	}
	if topic != nil {
		if err := domain.ValidateStringLength(*topic, 1, 200, "トピック"); err != nil {
			return nil, err
		}
		circle.Topic = strings.TrimSpace(*topic)
	}
	if description != nil {
		if err := domain.ValidateStringLength(*description, 0, 1000, "説明"); err != nil {
			return nil, err
		}
		circle.Description = strings.TrimSpace(*description)
	}

	if err := uc.circles.Update(ctx, circle); err != nil {
		return nil, err
	}
	return circle, nil
}

// DeleteStudyCircleUseCase はサークルを削除する。
type DeleteStudyCircleUseCase struct {
	circles repository.StudyCircleRepository
}

// NewDeleteStudyCircleUseCase は DeleteStudyCircleUseCase を生成する。
func NewDeleteStudyCircleUseCase(circles repository.StudyCircleRepository) *DeleteStudyCircleUseCase {
	return &DeleteStudyCircleUseCase{circles: circles}
}

// Execute はサークルと関連データを削除する。オーナーのみ。
func (uc *DeleteStudyCircleUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := findOwnedStudyCircle(ctx, uc.circles, id, userID); err != nil {
		return err
	}
	return uc.circles.Delete(ctx, id)
}

// ListStudyCircleMembersUseCase はサークルのメンバー一覧を取得する。
type ListStudyCircleMembersUseCase struct {
	circles repository.StudyCircleRepository
}

// NewListStudyCircleMembersUseCase は ListStudyCircleMembersUseCase を生成する。
func NewListStudyCircleMembersUseCase(circles repository.StudyCircleRepository) *ListStudyCircleMembersUseCase {
	return &ListStudyCircleMembersUseCase{circles: circles}
}

// Execute はメンバー一覧を返す。メンバーのみアクセスできる。
func (uc *ListStudyCircleMembersUseCase) Execute(ctx context.Context, circleID, userID uint) ([]model.StudyCircleMember, error) {
	if err := requireStudyCircleMember(ctx, uc.circles, circleID, userID); err != nil {
		return nil, err
	}
	return uc.circles.GetMembers(ctx, circleID)
}

// AddStudyCircleMemberUseCase はサークルにメンバーを追加する。
type AddStudyCircleMemberUseCase struct {
	circles repository.StudyCircleRepository
}

// NewAddStudyCircleMemberUseCase は AddStudyCircleMemberUseCase を生成する。
func NewAddStudyCircleMemberUseCase(circles repository.StudyCircleRepository) *AddStudyCircleMemberUseCase {
	return &AddStudyCircleMemberUseCase{circles: circles}
}

// Execute はメンバーを追加する。リクエスト者がメンバーであること・対象が未参加であること・
// 人数上限に達していないことを検証する。
func (uc *AddStudyCircleMemberUseCase) Execute(ctx context.Context, circleID, userID, targetUserID uint) error {
	if err := requireStudyCircleMember(ctx, uc.circles, circleID, userID); err != nil {
		return err
	}

	alreadyMember, err := uc.circles.IsMember(ctx, circleID, targetUserID)
	if err != nil {
		return err
	}
	if alreadyMember {
		return domain.ErrBadRequest
	}

	circle, err := uc.circles.FindByID(ctx, circleID)
	if err != nil || circle == nil {
		return domain.ErrNotFound
	}

	added, err := uc.circles.AddMemberWithinLimit(ctx, circleID, targetUserID, model.StudyCircleRoleMember)
	if err != nil {
		return err
	}
	if !added {
		return domain.NewError(domain.ErrCodeBadRequest, "メンバー上限に達しました", nil)
	}
	return nil
}

// UpdateStudyCircleMemberRoleUseCase はメンバーの役割を更新する。
type UpdateStudyCircleMemberRoleUseCase struct {
	circles repository.StudyCircleRepository
}

// NewUpdateStudyCircleMemberRoleUseCase は UpdateStudyCircleMemberRoleUseCase を生成する。
func NewUpdateStudyCircleMemberRoleUseCase(circles repository.StudyCircleRepository) *UpdateStudyCircleMemberRoleUseCase {
	return &UpdateStudyCircleMemberRoleUseCase{circles: circles}
}

// Execute はメンバーの役割を更新する。オーナーのみ操作でき、オーナー自身の役割は変更できない。
func (uc *UpdateStudyCircleMemberRoleUseCase) Execute(ctx context.Context, circleID, userID, targetUserID uint, role string) error {
	if !validStudyCircleMemberRoles[role] {
		return domain.NewError(domain.ErrCodeBadRequest, "無効な役割です", nil)
	}
	if _, err := findOwnedStudyCircle(ctx, uc.circles, circleID, userID); err != nil {
		return err
	}
	if userID == targetUserID {
		return domain.NewError(domain.ErrCodeBadRequest, "オーナー自身の役割は変更できません", nil)
	}
	isMember, err := uc.circles.IsMember(ctx, circleID, targetUserID)
	if err != nil {
		return err
	}
	if !isMember {
		return domain.NewError(domain.ErrCodeNotFound, "指定されたユーザーはメンバーではありません", nil)
	}
	return uc.circles.UpdateMemberRole(ctx, circleID, targetUserID, model.StudyCircleMemberRole(role))
}

// RemoveStudyCircleMemberUseCase はサークルからメンバーを除外する。
type RemoveStudyCircleMemberUseCase struct {
	circles repository.StudyCircleRepository
}

// NewRemoveStudyCircleMemberUseCase は RemoveStudyCircleMemberUseCase を生成する。
func NewRemoveStudyCircleMemberUseCase(circles repository.StudyCircleRepository) *RemoveStudyCircleMemberUseCase {
	return &RemoveStudyCircleMemberUseCase{circles: circles}
}

// Execute はメンバーを除外する。自分自身なら誰でも退出でき、他人の除外はオーナーのみ。
func (uc *RemoveStudyCircleMemberUseCase) Execute(ctx context.Context, circleID, userID, targetUserID uint) error {
	circle, err := uc.circles.FindByID(ctx, circleID)
	if err != nil || circle == nil {
		return domain.ErrNotFound
	}

	if userID != targetUserID && circle.OwnerID != userID {
		return domain.ErrForbidden
	}

	return uc.circles.RemoveMember(ctx, circleID, targetUserID)
}

// CreateStudyCircleStepUseCase はサークルにステップを追加する。
type CreateStudyCircleStepUseCase struct {
	circles repository.StudyCircleRepository
}

// NewCreateStudyCircleStepUseCase は CreateStudyCircleStepUseCase を生成する。
func NewCreateStudyCircleStepUseCase(circles repository.StudyCircleRepository) *CreateStudyCircleStepUseCase {
	return &CreateStudyCircleStepUseCase{circles: circles}
}

// Execute はステップを追加する。オーナーのみ。
func (uc *CreateStudyCircleStepUseCase) Execute(ctx context.Context, circleID, userID uint, step *model.StudyCircleStep) error {
	if err := domain.ValidateStringLength(step.ResourceURL, 0, 2000, "参考URL"); err != nil {
		return err
	}
	if _, err := findOwnedStudyCircle(ctx, uc.circles, circleID, userID); err != nil {
		return err
	}
	step.CircleID = circleID
	return uc.circles.CreateStep(ctx, step)
}

// UpdateStudyCircleStepUseCase はステップを更新する。
type UpdateStudyCircleStepUseCase struct {
	circles repository.StudyCircleRepository
}

// NewUpdateStudyCircleStepUseCase は UpdateStudyCircleStepUseCase を生成する。
func NewUpdateStudyCircleStepUseCase(circles repository.StudyCircleRepository) *UpdateStudyCircleStepUseCase {
	return &UpdateStudyCircleStepUseCase{circles: circles}
}

// Execute はステップのタイトル・説明を部分更新する。オーナーのみ。
func (uc *UpdateStudyCircleStepUseCase) Execute(ctx context.Context, circleID, userID, stepID uint, title, description *string) (*model.StudyCircleStep, error) {
	step, err := findOwnedStudyCircleStep(ctx, uc.circles, circleID, stepID, userID)
	if err != nil {
		return nil, err
	}

	if title != nil {
		if err := domain.ValidateStringLength(*title, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		step.Title = strings.TrimSpace(*title)
	}
	if description != nil {
		if err := domain.ValidateStringLength(*description, 0, 1000, "説明"); err != nil {
			return nil, err
		}
		step.Description = *description
	}

	if err := uc.circles.UpdateStep(ctx, step); err != nil {
		return nil, err
	}
	return step, nil
}

// DeleteStudyCircleStepUseCase はステップを削除する。
type DeleteStudyCircleStepUseCase struct {
	circles repository.StudyCircleRepository
}

// NewDeleteStudyCircleStepUseCase は DeleteStudyCircleStepUseCase を生成する。
func NewDeleteStudyCircleStepUseCase(circles repository.StudyCircleRepository) *DeleteStudyCircleStepUseCase {
	return &DeleteStudyCircleStepUseCase{circles: circles}
}

// Execute はステップを削除する。オーナーのみ。
func (uc *DeleteStudyCircleStepUseCase) Execute(ctx context.Context, circleID, userID, stepID uint) error {
	if _, err := findOwnedStudyCircleStep(ctx, uc.circles, circleID, stepID, userID); err != nil {
		return err
	}
	return uc.circles.DeleteStep(ctx, stepID)
}

// ReorderStudyCircleStepsUseCase はステップの表示順序を変更する。
type ReorderStudyCircleStepsUseCase struct {
	circles repository.StudyCircleRepository
}

// NewReorderStudyCircleStepsUseCase は ReorderStudyCircleStepsUseCase を生成する。
func NewReorderStudyCircleStepsUseCase(circles repository.StudyCircleRepository) *ReorderStudyCircleStepsUseCase {
	return &ReorderStudyCircleStepsUseCase{circles: circles}
}

// Execute はステップの表示順序をまとめて更新する。オーナーのみ。
func (uc *ReorderStudyCircleStepsUseCase) Execute(ctx context.Context, circleID, userID uint, orders []model.StepOrder) error {
	if _, err := findOwnedStudyCircle(ctx, uc.circles, circleID, userID); err != nil {
		return err
	}
	return uc.circles.ReorderSteps(ctx, circleID, orders)
}

// UpdateStudyCircleProgressUseCase は自分のステップ進捗を更新する。
type UpdateStudyCircleProgressUseCase struct {
	circles repository.StudyCircleRepository
}

// NewUpdateStudyCircleProgressUseCase は UpdateStudyCircleProgressUseCase を生成する。
func NewUpdateStudyCircleProgressUseCase(circles repository.StudyCircleRepository) *UpdateStudyCircleProgressUseCase {
	return &UpdateStudyCircleProgressUseCase{circles: circles}
}

// Execute は自分のステップ進捗を更新する。メンバーのみ。完了時のみ完了日時を記録する。
func (uc *UpdateStudyCircleProgressUseCase) Execute(ctx context.Context, circleID, userID, stepID uint, isCompleted bool) error {
	if err := requireStudyCircleMember(ctx, uc.circles, circleID, userID); err != nil {
		return err
	}

	progress := &model.StudyCircleMemberProgress{
		CircleID:    circleID,
		StepID:      stepID,
		UserID:      userID,
		IsCompleted: isCompleted,
	}
	if isCompleted {
		now := time.Now()
		progress.CompletedAt = &now
	}

	return uc.circles.UpsertProgress(ctx, progress)
}

// ListStudyCircleProgressUseCase は全メンバーの進捗を取得する。
type ListStudyCircleProgressUseCase struct {
	circles repository.StudyCircleRepository
}

// NewListStudyCircleProgressUseCase は ListStudyCircleProgressUseCase を生成する。
func NewListStudyCircleProgressUseCase(circles repository.StudyCircleRepository) *ListStudyCircleProgressUseCase {
	return &ListStudyCircleProgressUseCase{circles: circles}
}

// Execute は全メンバーの進捗を返す。メンバーのみアクセスできる。
func (uc *ListStudyCircleProgressUseCase) Execute(ctx context.Context, circleID, userID uint) ([]model.StudyCircleMemberProgress, error) {
	if err := requireStudyCircleMember(ctx, uc.circles, circleID, userID); err != nil {
		return nil, err
	}
	return uc.circles.GetProgress(ctx, circleID)
}

// CreateStudyCircleCheckinUseCase は日次チェックインを作成する。
type CreateStudyCircleCheckinUseCase struct {
	circles repository.StudyCircleRepository
}

// NewCreateStudyCircleCheckinUseCase は CreateStudyCircleCheckinUseCase を生成する。
func NewCreateStudyCircleCheckinUseCase(circles repository.StudyCircleRepository) *CreateStudyCircleCheckinUseCase {
	return &CreateStudyCircleCheckinUseCase{circles: circles}
}

// Execute は日次チェックインを作成する。メンバーのみ・同じ日に 2 回目は 409。
func (uc *CreateStudyCircleCheckinUseCase) Execute(ctx context.Context, circleID, userID uint, content string) (*model.StudyCircleCheckin, error) {
	if err := domain.ValidateStringLength(content, 1, 5000, "チェックイン内容"); err != nil {
		return nil, err
	}
	content = strings.TrimSpace(content)

	if err := requireStudyCircleMember(ctx, uc.circles, circleID, userID); err != nil {
		return nil, err
	}

	done, err := uc.circles.HasCheckedInToday(ctx, circleID, userID)
	if err != nil {
		return nil, err
	}
	if done {
		return nil, domain.NewError(domain.ErrCodeConflict, "本日は既にチェックイン済みです", nil)
	}

	checkin := &model.StudyCircleCheckin{
		CircleID: circleID,
		UserID:   userID,
		Date:     time.Now().Format("2006-01-02"),
		Content:  content,
	}
	if err := uc.circles.CreateCheckin(ctx, checkin); err != nil {
		return nil, err
	}
	return checkin, nil
}

// ListStudyCircleCheckinsUseCase はチェックイン履歴を取得する。
type ListStudyCircleCheckinsUseCase struct {
	circles repository.StudyCircleRepository
}

// NewListStudyCircleCheckinsUseCase は ListStudyCircleCheckinsUseCase を生成する。
func NewListStudyCircleCheckinsUseCase(circles repository.StudyCircleRepository) *ListStudyCircleCheckinsUseCase {
	return &ListStudyCircleCheckinsUseCase{circles: circles}
}

// Execute はチェックイン履歴を新しい順で返す。メンバーのみアクセスできる。
func (uc *ListStudyCircleCheckinsUseCase) Execute(ctx context.Context, circleID, userID uint) ([]model.StudyCircleCheckin, error) {
	if err := requireStudyCircleMember(ctx, uc.circles, circleID, userID); err != nil {
		return nil, err
	}
	return uc.circles.GetCheckins(ctx, circleID)
}

// GetStudyCircleStreakRankingUseCase はストリークランキングを取得する。
type GetStudyCircleStreakRankingUseCase struct {
	circles repository.StudyCircleRepository
}

// NewGetStudyCircleStreakRankingUseCase は GetStudyCircleStreakRankingUseCase を生成する。
func NewGetStudyCircleStreakRankingUseCase(circles repository.StudyCircleRepository) *GetStudyCircleStreakRankingUseCase {
	return &GetStudyCircleStreakRankingUseCase{circles: circles}
}

// Execute はストリークランキングを返す。メンバーのみアクセスできる。
func (uc *GetStudyCircleStreakRankingUseCase) Execute(ctx context.Context, circleID, userID uint) ([]model.CircleMemberStreak, error) {
	if err := requireStudyCircleMember(ctx, uc.circles, circleID, userID); err != nil {
		return nil, err
	}
	return uc.circles.GetStreakRanking(ctx, circleID)
}

// SearchStudyCirclesUseCase はキーワードでサークルを検索する。
type SearchStudyCirclesUseCase struct {
	circles repository.StudyCircleRepository
}

// NewSearchStudyCirclesUseCase は SearchStudyCirclesUseCase を生成する。
func NewSearchStudyCirclesUseCase(circles repository.StudyCircleRepository) *SearchStudyCirclesUseCase {
	return &SearchStudyCirclesUseCase{circles: circles}
}

// Execute は名前・トピック・説明への部分一致でサークルを検索する。
func (uc *SearchStudyCirclesUseCase) Execute(ctx context.Context, query string, limit, offset int) ([]model.StudyCircle, int64, error) {
	return uc.circles.Search(ctx, query, limit, offset)
}

// CountStudyCirclesUseCase は参加しているサークル数を取得する。
type CountStudyCirclesUseCase struct {
	circles repository.StudyCircleRepository
}

// NewCountStudyCirclesUseCase は CountStudyCirclesUseCase を生成する。
func NewCountStudyCirclesUseCase(circles repository.StudyCircleRepository) *CountStudyCirclesUseCase {
	return &CountStudyCirclesUseCase{circles: circles}
}

// Execute は指定ユーザーが参加しているサークル総数を返す。
func (uc *CountStudyCirclesUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.circles.CountByUserID(ctx, userID)
}
