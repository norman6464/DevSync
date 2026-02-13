package repository

import (
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StudyCircleRepository はスタディサークルのデータアクセスを提供するリポジトリ実装。
type StudyCircleRepository struct {
	db *gorm.DB
}

// NewStudyCircleRepository は新しいStudyCircleRepositoryインスタンスを生成する。
func NewStudyCircleRepository(db *gorm.DB) *StudyCircleRepository {
	return &StudyCircleRepository{db: db}
}

// Create はサークルをDBに保存する。
func (r *StudyCircleRepository) Create(circle *model.StudyCircle) error {
	return r.db.Create(circle).Error
}

// FindByID はIDでサークルを取得する。Steps, Members, Owner をプリロードする。
func (r *StudyCircleRepository) FindByID(id uint) (*model.StudyCircle, error) {
	var circle model.StudyCircle
	err := r.db.
		Preload("Owner").
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("Members").
		Preload("Members.User").
		First(&circle, id).Error
	if err != nil {
		return nil, err
	}
	return &circle, nil
}

// FindByUserID はユーザーが参加しているサークル一覧を返す。
func (r *StudyCircleRepository) FindByUserID(userID uint) ([]model.StudyCircle, error) {
	var circles []model.StudyCircle
	err := r.db.
		Preload("Owner").
		Preload("Members").
		Preload("Members.User").
		Joins("JOIN study_circle_members ON study_circle_members.circle_id = study_circles.id").
		Where("study_circle_members.user_id = ?", userID).
		Order("study_circles.updated_at DESC").
		Find(&circles).Error
	return circles, err
}

// Update はサークル情報を更新する。
func (r *StudyCircleRepository) Update(circle *model.StudyCircle) error {
	return r.db.Save(circle).Error
}

// Delete はサークルと関連データをトランザクション内で削除する。
func (r *StudyCircleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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
func (r *StudyCircleRepository) AddMember(circleID, userID uint, role model.StudyCircleMemberRole) error {
	member := model.StudyCircleMember{
		CircleID: circleID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
	return r.db.Create(&member).Error
}

// RemoveMember はメンバーを削除する。
func (r *StudyCircleRepository) RemoveMember(circleID, userID uint) error {
	return r.db.Where("circle_id = ? AND user_id = ?", circleID, userID).Delete(&model.StudyCircleMember{}).Error
}

// GetMembers はメンバー一覧を返す。
func (r *StudyCircleRepository) GetMembers(circleID uint) ([]model.StudyCircleMember, error) {
	var members []model.StudyCircleMember
	err := r.db.Preload("User").Where("circle_id = ?", circleID).Find(&members).Error
	return members, err
}

// IsMember はユーザーがメンバーかどうかを返す。
func (r *StudyCircleRepository) IsMember(circleID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.StudyCircleMember{}).
		Where("circle_id = ? AND user_id = ?", circleID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetMemberCount はメンバー数を返す。
func (r *StudyCircleRepository) GetMemberCount(circleID uint) (int, error) {
	var count int64
	err := r.db.Model(&model.StudyCircleMember{}).Where("circle_id = ?", circleID).Count(&count).Error
	return int(count), err
}

// CreateStep はステップを追加する。
func (r *StudyCircleRepository) CreateStep(step *model.StudyCircleStep) error {
	return r.db.Create(step).Error
}

// UpdateStep はステップを更新する。
func (r *StudyCircleRepository) UpdateStep(step *model.StudyCircleStep) error {
	return r.db.Save(step).Error
}

// DeleteStep はステップを削除する。
func (r *StudyCircleRepository) DeleteStep(stepID uint) error {
	return r.db.Delete(&model.StudyCircleStep{}, stepID).Error
}

// FindStepByID はステップをIDで取得する。
func (r *StudyCircleRepository) FindStepByID(stepID uint) (*model.StudyCircleStep, error) {
	var step model.StudyCircleStep
	err := r.db.First(&step, stepID).Error
	if err != nil {
		return nil, err
	}
	return &step, nil
}

// ReorderSteps はステップの表示順序を更新する。
func (r *StudyCircleRepository) ReorderSteps(circleID uint, stepOrders []StepOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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
func (r *StudyCircleRepository) UpsertProgress(progress *model.StudyCircleMemberProgress) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "circle_id"}, {Name: "step_id"}, {Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_completed", "completed_at"}),
	}).Create(progress).Error
}

// GetProgress はサークル全メンバーの進捗を返す。
func (r *StudyCircleRepository) GetProgress(circleID uint) ([]model.StudyCircleMemberProgress, error) {
	var progress []model.StudyCircleMemberProgress
	err := r.db.Preload("User").Where("circle_id = ?", circleID).Find(&progress).Error
	return progress, err
}

// CreateCheckin はチェックインをDBに保存する。
func (r *StudyCircleRepository) CreateCheckin(checkin *model.StudyCircleCheckin) error {
	return r.db.Create(checkin).Error
}

// GetCheckins はチェックイン一覧を新しい順で返す。
func (r *StudyCircleRepository) GetCheckins(circleID uint) ([]model.StudyCircleCheckin, error) {
	var checkins []model.StudyCircleCheckin
	err := r.db.Preload("User").
		Where("circle_id = ?", circleID).
		Order("created_at DESC").
		Find(&checkins).Error
	return checkins, err
}

// HasCheckedInToday は今日すでにチェックイン済みかを返す。
func (r *StudyCircleRepository) HasCheckedInToday(circleID, userID uint) (bool, error) {
	today := time.Now().Format("2006-01-02")
	var count int64
	err := r.db.Model(&model.StudyCircleCheckin{}).
		Where("circle_id = ? AND user_id = ? AND date = ?", circleID, userID, today).
		Count(&count).Error
	return count > 0, err
}

// GetStreakRanking はサークル内メンバーのストリークランキングを返す。
func (r *StudyCircleRepository) GetStreakRanking(circleID uint) ([]model.CircleMemberStreak, error) {
	// メンバー一覧を取得
	var members []model.StudyCircleMember
	if err := r.db.Preload("User").Where("circle_id = ?", circleID).Find(&members).Error; err != nil {
		return nil, err
	}

	var results []model.CircleMemberStreak
	for _, member := range members {
		// メンバーのチェックイン日付を取得（降順）
		var dates []string
		r.db.Model(&model.StudyCircleCheckin{}).
			Where("circle_id = ? AND user_id = ?", circleID, member.UserID).
			Order("date DESC").
			Pluck("date", &dates)

		streak := calculateCheckinStreak(dates)
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
			CurrentStreak: streak,
			TotalCheckins: len(dates),
		})
	}

	// ストリーク降順でソート
	sort.Slice(results, func(i, j int) bool {
		return results[i].CurrentStreak > results[j].CurrentStreak
	})

	return results, nil
}

// calculateCheckinStreak はチェックイン日付リスト（降順）から連続日数を計算する。
func calculateCheckinStreak(dates []string) int {
	if len(dates) == 0 {
		return 0
	}

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 最新のチェックインが今日または昨日でなければストリーク0
	if dates[0] != today && dates[0] != yesterday {
		return 0
	}

	streak := 1
	for i := 1; i < len(dates); i++ {
		prev, _ := time.Parse("2006-01-02", dates[i-1])
		curr, _ := time.Parse("2006-01-02", dates[i])
		diff := prev.Sub(curr).Hours() / 24
		if diff == 1 {
			streak++
		} else {
			break
		}
	}

	return streak
}

// Search はキーワードでスタディサークルを検索する（名前、トピック、説明に部分一致）。
func (r *StudyCircleRepository) Search(query string, limit, offset int) (interface{}, int64, error) {
	var circles []model.StudyCircle
	var total int64

	searchPattern := "%" + query + "%"
	db := r.db.Preload("Owner").Preload("Members").Preload("Members.User").
		Where("name LIKE ? OR topic LIKE ? OR description LIKE ?", searchPattern, searchPattern, searchPattern).
		Order("created_at DESC")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&circles).Error
	return circles, total, err
}
