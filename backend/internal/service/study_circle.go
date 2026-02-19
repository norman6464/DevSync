package service

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// StudyCircleService はスタディサークルのビジネスロジックを提供する。
type StudyCircleService struct {
	repo repository.StudyCircleRepositoryInterface
}

// NewStudyCircleService は新しいStudyCircleServiceインスタンスを生成する。
func NewStudyCircleService(repo repository.StudyCircleRepositoryInterface) *StudyCircleService {
	return &StudyCircleService{repo: repo}
}

func (s *StudyCircleService) findAndCheckOwnership(id, userID uint) (*model.StudyCircle, error) {
	circle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	if circle.OwnerID != userID {
		return nil, ErrForbidden
	}
	return circle, nil
}

func (s *StudyCircleService) findAndCheckStepOwnership(circleID, stepID, userID uint) (*model.StudyCircleStep, error) {
	if _, err := s.findAndCheckOwnership(circleID, userID); err != nil {
		return nil, err
	}
	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return nil, ErrNotFound
	}
	if step.CircleID != circleID {
		return nil, ErrNotFound
	}
	return step, nil
}

// Create はサークルを作成し、オーナーをメンバーとして自動追加する。
func (s *StudyCircleService) Create(circle *model.StudyCircle, memberIDs []uint) error {
	if circle.MaxMembers < 3 || circle.MaxMembers > 10 {
		circle.MaxMembers = 5
	}
	circle.Status = model.StudyCircleStatusActive

	if err := s.repo.Create(circle); err != nil {
		return err
	}

	if err := s.repo.AddMember(circle.ID, circle.OwnerID, model.StudyCircleRoleOwner); err != nil {
		return err
	}

	for _, memberID := range memberIDs {
		if memberID != circle.OwnerID {
			_ = s.repo.AddMember(circle.ID, memberID, model.StudyCircleRoleMember)
		}
	}

	return nil
}

// GetMyCircles はユーザーが参加しているサークル一覧を返す。
func (s *StudyCircleService) GetMyCircles(userID uint) ([]model.StudyCircle, error) {
	return s.repo.FindByUserID(userID)
}

// GetByID はサークル詳細を返す。メンバーのみアクセス可能。
func (s *StudyCircleService) GetByID(id, userID uint) (*model.StudyCircle, error) {
	circle, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}

	isMember, err := s.repo.IsMember(id, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}

	return circle, nil
}

// Update はサークル情報を更新する。オーナーのみ。
func (s *StudyCircleService) Update(id, userID uint, name, topic, description *string) (*model.StudyCircle, error) {
	circle, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if name != nil {
		circle.Name = *name
	}
	if topic != nil {
		circle.Topic = *topic
	}
	if description != nil {
		circle.Description = *description
	}

	if err := s.repo.Update(circle); err != nil {
		return nil, err
	}
	return circle, nil
}

// Delete はサークルを削除する。オーナーのみ。
func (s *StudyCircleService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// GetMembers はサークルメンバー一覧を返す。メンバーのみ。
func (s *StudyCircleService) GetMembers(circleID, userID uint) ([]model.StudyCircleMember, error) {
	isMember, err := s.repo.IsMember(circleID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}
	return s.repo.GetMembers(circleID)
}

// AddMember はメンバーを追加する。リクエスト者がメンバーであること＋上限チェック。
func (s *StudyCircleService) AddMember(circleID, userID, targetUserID uint) error {
	isMember, err := s.repo.IsMember(circleID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrForbidden
	}

	circle, err := s.repo.FindByID(circleID)
	if err != nil {
		return ErrNotFound
	}

	count, err := s.repo.GetMemberCount(circleID)
	if err != nil {
		return err
	}
	if count >= circle.MaxMembers {
		return domain.NewError(domain.ErrCodeBadRequest, "メンバー上限に達しました", nil)
	}

	return s.repo.AddMember(circleID, targetUserID, model.StudyCircleRoleMember)
}

// RemoveMember はメンバーを除外する。オーナーは誰でも除外可、メンバーは自分のみ退出可。
func (s *StudyCircleService) RemoveMember(circleID, userID, targetUserID uint) error {
	circle, err := s.repo.FindByID(circleID)
	if err != nil {
		return ErrNotFound
	}

	if userID == targetUserID {
		return s.repo.RemoveMember(circleID, targetUserID)
	}

	if circle.OwnerID != userID {
		return ErrForbidden
	}

	return s.repo.RemoveMember(circleID, targetUserID)
}

// CreateStep はステップを追加する。オーナーのみ。
func (s *StudyCircleService) CreateStep(circleID, userID uint, step *model.StudyCircleStep) error {
	if _, err := s.findAndCheckOwnership(circleID, userID); err != nil {
		return err
	}

	step.CircleID = circleID
	return s.repo.CreateStep(step)
}

// UpdateStep はステップを更新する。オーナーのみ。
func (s *StudyCircleService) UpdateStep(circleID, userID, stepID uint, title, description *string) (*model.StudyCircleStep, error) {
	step, err := s.findAndCheckStepOwnership(circleID, stepID, userID)
	if err != nil {
		return nil, err
	}

	if title != nil {
		step.Title = *title
	}
	if description != nil {
		step.Description = *description
	}

	if err := s.repo.UpdateStep(step); err != nil {
		return nil, err
	}
	return step, nil
}

// DeleteStep はステップを削除する。オーナーのみ。
func (s *StudyCircleService) DeleteStep(circleID, userID, stepID uint) error {
	if _, err := s.findAndCheckStepOwnership(circleID, stepID, userID); err != nil {
		return err
	}
	return s.repo.DeleteStep(stepID)
}

// ReorderSteps はステップの順序を変更する。オーナーのみ。
func (s *StudyCircleService) ReorderSteps(circleID, userID uint, orders []model.StepOrder) error {
	if _, err := s.findAndCheckOwnership(circleID, userID); err != nil {
		return err
	}
	return s.repo.ReorderSteps(circleID, orders)
}

// UpdateProgress は自分のステップ進捗を更新する。メンバーのみ。
func (s *StudyCircleService) UpdateProgress(circleID, userID, stepID uint, isCompleted bool) error {
	isMember, err := s.repo.IsMember(circleID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return ErrForbidden
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

	return s.repo.UpsertProgress(progress)
}

// GetProgress は全メンバーの進捗を返す。メンバーのみ。
func (s *StudyCircleService) GetProgress(circleID, userID uint) ([]model.StudyCircleMemberProgress, error) {
	isMember, err := s.repo.IsMember(circleID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}
	return s.repo.GetProgress(circleID)
}

// CreateCheckin は日次チェックインを作成する。1日1回制限。
func (s *StudyCircleService) CreateCheckin(circleID, userID uint, content string) (*model.StudyCircleCheckin, error) {
	isMember, err := s.repo.IsMember(circleID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}

	done, err := s.repo.HasCheckedInToday(circleID, userID)
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
	if err := s.repo.CreateCheckin(checkin); err != nil {
		return nil, err
	}
	return checkin, nil
}

// GetCheckins はチェックイン履歴を返す。メンバーのみ。
func (s *StudyCircleService) GetCheckins(circleID, userID uint) ([]model.StudyCircleCheckin, error) {
	isMember, err := s.repo.IsMember(circleID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}
	return s.repo.GetCheckins(circleID)
}

// GetStreakRanking はストリークランキングを返す。メンバーのみ。
func (s *StudyCircleService) GetStreakRanking(circleID, userID uint) ([]model.CircleMemberStreak, error) {
	isMember, err := s.repo.IsMember(circleID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, ErrForbidden
	}
	return s.repo.GetStreakRanking(circleID)
}

// SearchCircles はキーワードでスタディサークルを検索する。
func (s *StudyCircleService) SearchCircles(query string, limit, offset int) (interface{}, int64, error) {
	return s.repo.Search(query, limit, offset)
}
