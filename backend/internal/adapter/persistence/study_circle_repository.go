package persistence

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// studyCircleRepository は [repository.StudyCircleRepository] の GORM 実装。
type studyCircleRepository struct {
	db *gorm.DB
}

// NewStudyCircleRepository は StudyCircleRepository の GORM 実装を返す。
func NewStudyCircleRepository(db *gorm.DB) repository.StudyCircleRepository {
	return &studyCircleRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.StudyCircleRepository = (*studyCircleRepository)(nil)

// Create はサークルをDBに保存する。
func (r *studyCircleRepository) Create(ctx context.Context, circle *model.StudyCircle) error {
	return r.db.WithContext(ctx).Create(circle).Error
}

// FindByID はIDでサークルを取得する。Owner, Steps, Members をプリロードする。
// 不在の場合は (nil, nil) を返す。
func (r *studyCircleRepository) FindByID(ctx context.Context, id uint) (*model.StudyCircle, error) {
	var circle model.StudyCircle
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("Members").
		Preload("Members.User").
		First(&circle, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &circle, nil
}

// GetByStatus はユーザーが参加しているサークルをステータスでフィルタリングして返す。
func (r *studyCircleRepository) GetByStatus(ctx context.Context, userID uint, status string) ([]model.StudyCircle, error) {
	var circles []model.StudyCircle
	err := r.db.WithContext(ctx).
		Preload("Owner").
		Preload("Members").
		Preload("Members.User").
		Joins("JOIN study_circle_members ON study_circle_members.circle_id = study_circles.id").
		Where("study_circle_members.user_id = ? AND study_circles.status = ?", userID, status).
		Order("study_circles.updated_at DESC").
		Find(&circles).Error
	return circles, err
}

// FindByUserID はユーザーが参加しているサークル一覧をページネーション付きで返す。
func (r *studyCircleRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.StudyCircle, int64, error) {
	scope := r.db.WithContext(ctx).
		Joins("JOIN study_circle_members ON study_circle_members.circle_id = study_circles.id").
		Where("study_circle_members.user_id = ?", userID)

	var total int64
	if err := scope.Session(&gorm.Session{}).Model(&model.StudyCircle{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var circles []model.StudyCircle
	err := scope.Session(&gorm.Session{}).
		Preload("Owner").
		Preload("Members").
		Preload("Members.User").
		Order("study_circles.updated_at DESC").
		Limit(limit).Offset(offset).
		Find(&circles).Error
	return circles, total, err
}

// Update はサークル情報を更新する。
func (r *studyCircleRepository) Update(ctx context.Context, circle *model.StudyCircle) error {
	return r.db.WithContext(ctx).Save(circle).Error
}

// Delete はサークルと関連データをトランザクション内で削除する。
func (r *studyCircleRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("circle_id = ?", id).Delete(&model.StudyCircleCheckin{}).Error; err != nil {
			return err
		}
		if err := tx.Where("circle_id = ?", id).Delete(&model.StudyCircleMemberProgress{}).Error; err != nil {
			return err
		}
		if err := tx.Where("circle_id = ?", id).Delete(&model.StudyCircleStep{}).Error; err != nil {
			return err
		}
		if err := tx.Where("circle_id = ?", id).Delete(&model.StudyCircleMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.StudyCircle{}, id).Error
	})
}

// AddMember はメンバーを追加する。
func (r *studyCircleRepository) AddMember(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error {
	member := model.StudyCircleMember{
		CircleID: circleID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
	return r.db.WithContext(ctx).Create(&member).Error
}

// RemoveMember はメンバーを削除する。
func (r *studyCircleRepository) RemoveMember(ctx context.Context, circleID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("circle_id = ? AND user_id = ?", circleID, userID).
		Delete(&model.StudyCircleMember{}).Error
}

// UpdateMemberRole はメンバーの役割を更新する。
func (r *studyCircleRepository) UpdateMemberRole(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error {
	return r.db.WithContext(ctx).Model(&model.StudyCircleMember{}).
		Where("circle_id = ? AND user_id = ?", circleID, userID).
		Update("role", role).Error
}

// GetMembers はメンバー一覧を返す。
func (r *studyCircleRepository) GetMembers(ctx context.Context, circleID uint) ([]model.StudyCircleMember, error) {
	var members []model.StudyCircleMember
	err := r.db.WithContext(ctx).Preload("User").Where("circle_id = ?", circleID).Find(&members).Error
	return members, err
}

// IsMember はユーザーがメンバーかどうかを返す。
func (r *studyCircleRepository) IsMember(ctx context.Context, circleID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.StudyCircleMember{}).
		Where("circle_id = ? AND user_id = ?", circleID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetMemberCount はメンバー数を返す。
func (r *studyCircleRepository) GetMemberCount(ctx context.Context, circleID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.StudyCircleMember{}).
		Where("circle_id = ?", circleID).Count(&count).Error
	return int(count), err
}

// CountByUserID は指定ユーザーが参加しているスタディサークル総数を返す。
func (r *studyCircleRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.StudyCircleMember{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CreateStep はステップを追加する。
func (r *studyCircleRepository) CreateStep(ctx context.Context, step *model.StudyCircleStep) error {
	return r.db.WithContext(ctx).Create(step).Error
}

// UpdateStep はステップを更新する。
func (r *studyCircleRepository) UpdateStep(ctx context.Context, step *model.StudyCircleStep) error {
	return r.db.WithContext(ctx).Save(step).Error
}

// DeleteStep はステップを削除する。
func (r *studyCircleRepository) DeleteStep(ctx context.Context, stepID uint) error {
	return r.db.WithContext(ctx).Delete(&model.StudyCircleStep{}, stepID).Error
}

// FindStepByID はステップをIDで取得する。不在の場合は (nil, nil) を返す。
func (r *studyCircleRepository) FindStepByID(ctx context.Context, stepID uint) (*model.StudyCircleStep, error) {
	var step model.StudyCircleStep
	err := r.db.WithContext(ctx).First(&step, stepID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &step, nil
}

// ReorderSteps はステップの表示順序を更新する。
func (r *studyCircleRepository) ReorderSteps(ctx context.Context, circleID uint, stepOrders []model.StepOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, o := range stepOrders {
			if err := tx.Model(&model.StudyCircleStep{}).
				Where("id = ? AND circle_id = ?", o.StepID, circleID).
				Update("order_index", o.OrderIndex).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpsertProgress はメンバーのステップ進捗を更新（なければ作成）する。
func (r *studyCircleRepository) UpsertProgress(ctx context.Context, progress *model.StudyCircleMemberProgress) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "circle_id"}, {Name: "step_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_completed", "completed_at"}),
	}).Create(progress).Error
}

// GetProgress はサークル全メンバーの進捗を返す。
func (r *studyCircleRepository) GetProgress(ctx context.Context, circleID uint) ([]model.StudyCircleMemberProgress, error) {
	var progress []model.StudyCircleMemberProgress
	err := r.db.WithContext(ctx).Preload("User").Where("circle_id = ?", circleID).Find(&progress).Error
	return progress, err
}

// CreateCheckin はチェックインをDBに保存する。
func (r *studyCircleRepository) CreateCheckin(ctx context.Context, checkin *model.StudyCircleCheckin) error {
	return r.db.WithContext(ctx).Create(checkin).Error
}

// GetCheckins はチェックイン一覧を新しい順で返す。
func (r *studyCircleRepository) GetCheckins(ctx context.Context, circleID uint) ([]model.StudyCircleCheckin, error) {
	var checkins []model.StudyCircleCheckin
	err := r.db.WithContext(ctx).Preload("User").
		Where("circle_id = ?", circleID).
		Order("created_at DESC").
		Find(&checkins).Error
	return checkins, err
}

// HasCheckedInToday は今日すでにチェックイン済みかを返す。
func (r *studyCircleRepository) HasCheckedInToday(ctx context.Context, circleID, userID uint) (bool, error) {
	today := time.Now().Format("2006-01-02")
	var count int64
	err := r.db.WithContext(ctx).Model(&model.StudyCircleCheckin{}).
		Where("circle_id = ? AND user_id = ? AND date = ?", circleID, userID, today).
		Count(&count).Error
	return count > 0, err
}

// GetStreakRanking はサークル内メンバーのストリークランキングを返す（連続日数の降順）。
func (r *studyCircleRepository) GetStreakRanking(ctx context.Context, circleID uint) ([]model.CircleMemberStreak, error) {
	db := r.db.WithContext(ctx)

	var members []model.StudyCircleMember
	if err := db.Preload("User").Where("circle_id = ?", circleID).Find(&members).Error; err != nil {
		return nil, err
	}

	var results []model.CircleMemberStreak
	for _, member := range members {
		var dates []string
		db.Model(&model.StudyCircleCheckin{}).
			Where("circle_id = ? AND user_id = ?", circleID, member.UserID).
			Order("date DESC").
			Pluck("date", &dates)

		userName := ""
		avatarURL := ""
		if member.User != nil {
			userName = member.User.Name
			avatarURL = member.User.AvatarURL
		}

		results = append(results, model.CircleMemberStreak{
			UserID:        member.UserID,
			UserName:      userName,
			AvatarURL:     avatarURL,
			CurrentStreak: calculateCheckinStreak(dates),
			TotalCheckins: len(dates),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CurrentStreak > results[j].CurrentStreak
	})

	return results, nil
}

// calculateCheckinStreak はチェックイン日付リスト（降順）から連続日数を計算する。
// 最新のチェックインが今日でも昨日でもなければ連続は途切れているとみなして 0 を返す。
func calculateCheckinStreak(dates []string) int {
	if len(dates) == 0 {
		return 0
	}

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if dates[0] != today && dates[0] != yesterday {
		return 0
	}

	streak := 1
	for i := 1; i < len(dates); i++ {
		prev, _ := time.Parse("2006-01-02", dates[i-1])
		curr, _ := time.Parse("2006-01-02", dates[i])
		if prev.Sub(curr).Hours()/24 != 1 {
			break
		}
		streak++
	}

	return streak
}

// Search はキーワードでスタディサークルを検索する（名前・トピック・説明に部分一致）。
func (r *studyCircleRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.StudyCircle, int64, error) {
	pattern := escapeLikePattern(query)
	scope := r.db.WithContext(ctx).Model(&model.StudyCircle{}).
		Where("name LIKE ? OR topic LIKE ? OR description LIKE ?", pattern, pattern, pattern)

	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var circles []model.StudyCircle
	err := scope.Session(&gorm.Session{}).
		Preload("Owner").Preload("Members").Preload("Members.User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&circles).Error
	return circles, total, err
}
