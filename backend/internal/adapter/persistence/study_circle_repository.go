package persistence

import (
	"context"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// studyCircleRepository は [repository.StudyCircleRepository] の sqlc(pgx) 実装。
// CreateWithOwner・Delete・AddMemberWithinLimit・ReorderSteps は複数文を
// 1トランザクションで行うため、Queries だけでなくトランザクションを開始できる
// *pgxpool.Pool を直接保持する。
type studyCircleRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewStudyCircleRepository は StudyCircleRepository の sqlc(pgx) 実装を返す。
func NewStudyCircleRepository(pool *pgxpool.Pool) repository.StudyCircleRepository {
	return &studyCircleRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.StudyCircleRepository = (*studyCircleRepository)(nil)

func toModelStudyCircle(row sqlcgen.StudyCircle) model.StudyCircle {
	return model.StudyCircle{
		ID:          uint(row.ID),
		Name:        row.Name,
		Topic:       row.Topic,
		Description: fromStringPtr(row.Description),
		OwnerID:     uint(row.OwnerID),
		MaxMembers:  int(row.MaxMembers),
		Status:      model.StudyCircleStatus(fromStringPtr(row.Status)),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

func toModelStudyCircleMember(row sqlcgen.StudyCircleMember) model.StudyCircleMember {
	return model.StudyCircleMember{
		ID:       uint(row.ID),
		CircleID: uint(row.CircleID),
		UserID:   uint(row.UserID),
		Role:     model.StudyCircleMemberRole(fromStringPtr(row.Role)),
		JoinedAt: timeValue(fromTimestamptz(row.JoinedAt)),
	}
}

func toModelStudyCircleStep(row sqlcgen.StudyCircleStep) model.StudyCircleStep {
	return model.StudyCircleStep{
		ID:          uint(row.ID),
		CircleID:    uint(row.CircleID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		OrderIndex:  int(fromInt64PtrValue(row.OrderIndex)),
		ResourceURL: fromStringPtr(row.ResourceUrl),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

func toModelStudyCircleMemberProgress(row sqlcgen.StudyCircleMemberProgress) model.StudyCircleMemberProgress {
	return model.StudyCircleMemberProgress{
		ID:          uint(row.ID),
		CircleID:    uint(row.CircleID),
		StepID:      uint(row.StepID),
		UserID:      uint(row.UserID),
		IsCompleted: row.IsCompleted,
		CompletedAt: fromTimestamptz(row.CompletedAt),
	}
}

func toModelStudyCircleCheckin(row sqlcgen.StudyCircleCheckin) model.StudyCircleCheckin {
	return model.StudyCircleCheckin{
		ID:        uint(row.ID),
		CircleID:  uint(row.CircleID),
		UserID:    uint(row.UserID),
		Date:      row.Date,
		Content:   row.Content,
		CreatedAt: timeValue(fromTimestamptz(row.CreatedAt)),
	}
}

// attachStudyCircleMembers は複数メンバー行（User Preload込み）を model.StudyCircle へ付与する。
func attachStudyCircleMembers(ctx context.Context, q *sqlcgen.Queries, circle *model.StudyCircle) error {
	memberRows, err := q.ListStudyCircleMembersWithUserByCircle(ctx, int64(circle.ID))
	if err != nil {
		return err
	}
	members := make([]model.StudyCircleMember, len(memberRows))
	for i, row := range memberRows {
		members[i] = toModelStudyCircleMember(row.StudyCircleMember)
		user := toModelUser(row.User)
		members[i].User = &user
	}
	circle.Members = members
	return nil
}

// CreateWithOwner はサークル行の作成とオーナーのメンバー登録を 1 トランザクションで保存する。
// サークル行だけ残るとメンバー条件の各操作にオーナー本人すら入れなくなるため、
// オーナーのメンバー登録までを 1 トランザクションで確定する。
// Statusは移行前のGORMの `gorm:"default:'active'"` に相当し、未指定（ゼロ値）なら active を補う。
func (r *studyCircleRepository) CreateWithOwner(ctx context.Context, circle *model.StudyCircle) error {
	status := circle.Status
	if status == "" {
		status = model.StudyCircleStatusActive
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	row, err := q.CreateStudyCircle(ctx, sqlcgen.CreateStudyCircleParams{
		Name:        circle.Name,
		Topic:       circle.Topic,
		Description: &circle.Description,
		OwnerID:     int64(circle.OwnerID),
		MaxMembers:  int64(circle.MaxMembers),
		Status:      (*string)(&status),
	})
	if err != nil {
		return err
	}
	*circle = toModelStudyCircle(row)

	roleOwner := string(model.StudyCircleRoleOwner)
	if err := q.CreateStudyCircleMember(ctx, sqlcgen.CreateStudyCircleMemberParams{
		CircleID: int64(circle.ID),
		UserID:   int64(circle.OwnerID),
		Role:     &roleOwner,
		JoinedAt: toTimestamptzNotNull(time.Now()),
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// FindByID はIDでサークルを取得する。Owner, Steps, Members をプリロードする。
// 不在の場合は (nil, nil) を返す。
func (r *studyCircleRepository) FindByID(ctx context.Context, id uint) (*model.StudyCircle, error) {
	row, err := r.q.GetStudyCircleWithOwnerByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	circle := toModelStudyCircle(row.StudyCircle)
	owner := toModelUser(row.User)
	circle.Owner = &owner

	stepRows, err := r.q.ListStudyCircleStepsByCircle(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	circle.Steps = make([]model.StudyCircleStep, len(stepRows))
	for i, row := range stepRows {
		circle.Steps[i] = toModelStudyCircleStep(row)
	}

	if err := attachStudyCircleMembers(ctx, r.q, &circle); err != nil {
		return nil, err
	}

	return &circle, nil
}

// GetByStatus はユーザーが参加しているサークルをステータスでフィルタリングして返す。
func (r *studyCircleRepository) GetByStatus(ctx context.Context, userID uint, status string) ([]model.StudyCircle, error) {
	rows, err := r.q.ListStudyCirclesByUserAndStatus(ctx, sqlcgen.ListStudyCirclesByUserAndStatusParams{
		UserID: int64(userID),
		Status: &status,
	})
	if err != nil {
		return nil, err
	}

	circles := make([]model.StudyCircle, len(rows))
	for i, row := range rows {
		circles[i] = toModelStudyCircle(row.StudyCircle)
		owner := toModelUser(row.User)
		circles[i].Owner = &owner
		if err := attachStudyCircleMembers(ctx, r.q, &circles[i]); err != nil {
			return nil, err
		}
	}
	return circles, nil
}

// FindByUserID はユーザーが参加しているサークル一覧をページネーション付きで返す。
func (r *studyCircleRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.StudyCircle, int64, error) {
	total, err := r.q.CountStudyCirclesByUserMembership(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListStudyCirclesByUser(ctx, sqlcgen.ListStudyCirclesByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	circles := make([]model.StudyCircle, len(rows))
	for i, row := range rows {
		circles[i] = toModelStudyCircle(row.StudyCircle)
		owner := toModelUser(row.User)
		circles[i].Owner = &owner
		if err := attachStudyCircleMembers(ctx, r.q, &circles[i]); err != nil {
			return nil, 0, err
		}
	}
	return circles, total, nil
}

// Update はサークル情報を更新する（GORMのSave＝全カラム上書きに相当）。
func (r *studyCircleRepository) Update(ctx context.Context, circle *model.StudyCircle) error {
	row, err := r.q.UpdateStudyCircle(ctx, sqlcgen.UpdateStudyCircleParams{
		ID:          int64(circle.ID),
		Name:        circle.Name,
		Topic:       circle.Topic,
		Description: &circle.Description,
		MaxMembers:  int64(circle.MaxMembers),
		Status:      (*string)(&circle.Status),
	})
	if err != nil {
		return err
	}
	*circle = toModelStudyCircle(row)
	return nil
}

// Delete はサークルと関連データをトランザクション内で削除する。
func (r *studyCircleRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	circleID := int64(id)
	if err := q.DeleteStudyCircleCheckinsByCircle(ctx, circleID); err != nil {
		return err
	}
	if err := q.DeleteStudyCircleMemberProgressByCircle(ctx, circleID); err != nil {
		return err
	}
	if err := q.DeleteStudyCircleStepsByCircle(ctx, circleID); err != nil {
		return err
	}
	if err := q.DeleteStudyCircleMembersByCircle(ctx, circleID); err != nil {
		return err
	}
	if err := q.DeleteStudyCircle(ctx, circleID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AddMember はメンバーを追加する。
func (r *studyCircleRepository) AddMember(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error {
	roleStr := string(role)
	return r.q.CreateStudyCircleMember(ctx, sqlcgen.CreateStudyCircleMemberParams{
		CircleID: int64(circleID),
		UserID:   int64(userID),
		Role:     &roleStr,
		JoinedAt: toTimestamptzNotNull(time.Now()),
	})
}

// AddMemberWithinLimit はサークル行をロックして現在人数を数え、上限未満のときだけ
// 同一トランザクションでメンバーを追加する。上限到達時は (false, nil) を返す。
// 「数える → 追加する」を同時実行しても上限を超えないよう、
// サークル行の行ロックで直列化してから人数を確定する。
func (r *studyCircleRepository) AddMemberWithinLimit(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	maxMembers, err := q.LockStudyCircleForMemberLimit(ctx, int64(circleID))
	if err != nil {
		return false, err
	}
	count, err := q.CountStudyCircleMembersByCircle(ctx, int64(circleID))
	if err != nil {
		return false, err
	}
	if count >= maxMembers {
		return false, tx.Commit(ctx)
	}

	roleStr := string(role)
	if err := q.CreateStudyCircleMember(ctx, sqlcgen.CreateStudyCircleMemberParams{
		CircleID: int64(circleID),
		UserID:   int64(userID),
		Role:     &roleStr,
		JoinedAt: toTimestamptzNotNull(time.Now()),
	}); err != nil {
		return false, err
	}
	return true, tx.Commit(ctx)
}

// RemoveMember はメンバーを削除する。
func (r *studyCircleRepository) RemoveMember(ctx context.Context, circleID, userID uint) error {
	return r.q.DeleteStudyCircleMember(ctx, sqlcgen.DeleteStudyCircleMemberParams{
		CircleID: int64(circleID),
		UserID:   int64(userID),
	})
}

// UpdateMemberRole はメンバーの役割を更新する。
func (r *studyCircleRepository) UpdateMemberRole(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error {
	roleStr := string(role)
	return r.q.UpdateStudyCircleMemberRole(ctx, sqlcgen.UpdateStudyCircleMemberRoleParams{
		CircleID: int64(circleID),
		UserID:   int64(userID),
		Role:     &roleStr,
	})
}

// GetMembers はメンバー一覧を返す。
func (r *studyCircleRepository) GetMembers(ctx context.Context, circleID uint) ([]model.StudyCircleMember, error) {
	rows, err := r.q.ListStudyCircleMembersWithUserByCircle(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}
	members := make([]model.StudyCircleMember, len(rows))
	for i, row := range rows {
		members[i] = toModelStudyCircleMember(row.StudyCircleMember)
		user := toModelUser(row.User)
		members[i].User = &user
	}
	return members, nil
}

// IsMember はユーザーがメンバーかどうかを返す。
func (r *studyCircleRepository) IsMember(ctx context.Context, circleID, userID uint) (bool, error) {
	count, err := r.q.CountStudyCircleMembership(ctx, sqlcgen.CountStudyCircleMembershipParams{
		CircleID: int64(circleID),
		UserID:   int64(userID),
	})
	return count > 0, err
}

// CountByUserID は指定ユーザーが参加しているスタディサークル総数を返す。
func (r *studyCircleRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountStudyCircleMembershipsByUser(ctx, int64(userID))
}

// CreateStep はステップを追加する。
func (r *studyCircleRepository) CreateStep(ctx context.Context, step *model.StudyCircleStep) error {
	row, err := r.q.CreateStudyCircleStep(ctx, sqlcgen.CreateStudyCircleStepParams{
		CircleID:    int64(step.CircleID),
		Title:       step.Title,
		Description: &step.Description,
		OrderIndex:  toInt64Ptr(step.OrderIndex),
		ResourceUrl: &step.ResourceURL,
	})
	if err != nil {
		return err
	}
	*step = toModelStudyCircleStep(row)
	return nil
}

// UpdateStep はステップを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *studyCircleRepository) UpdateStep(ctx context.Context, step *model.StudyCircleStep) error {
	row, err := r.q.UpdateStudyCircleStep(ctx, sqlcgen.UpdateStudyCircleStepParams{
		ID:          int64(step.ID),
		Title:       step.Title,
		Description: &step.Description,
		OrderIndex:  toInt64Ptr(step.OrderIndex),
		ResourceUrl: &step.ResourceURL,
	})
	if err != nil {
		return err
	}
	*step = toModelStudyCircleStep(row)
	return nil
}

// DeleteStep はステップを削除する。
func (r *studyCircleRepository) DeleteStep(ctx context.Context, stepID uint) error {
	return r.q.DeleteStudyCircleStep(ctx, int64(stepID))
}

// FindStepByID はステップをIDで取得する。不在の場合は (nil, nil) を返す。
func (r *studyCircleRepository) FindStepByID(ctx context.Context, stepID uint) (*model.StudyCircleStep, error) {
	row, err := r.q.GetStudyCircleStepByID(ctx, int64(stepID))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	step := toModelStudyCircleStep(row)
	return &step, nil
}

// ReorderSteps はステップの表示順序を更新する。
func (r *studyCircleRepository) ReorderSteps(ctx context.Context, circleID uint, stepOrders []model.StepOrder) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for _, o := range stepOrders {
		if err := q.ReorderStudyCircleStep(ctx, sqlcgen.ReorderStudyCircleStepParams{
			ID:         int64(o.StepID),
			CircleID:   int64(circleID),
			OrderIndex: toInt64Ptr(o.OrderIndex),
		}); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// UpsertProgress はメンバーのステップ進捗を更新（なければ作成）する。
func (r *studyCircleRepository) UpsertProgress(ctx context.Context, progress *model.StudyCircleMemberProgress) error {
	row, err := r.q.UpsertStudyCircleMemberProgress(ctx, sqlcgen.UpsertStudyCircleMemberProgressParams{
		CircleID:    int64(progress.CircleID),
		StepID:      int64(progress.StepID),
		UserID:      int64(progress.UserID),
		IsCompleted: progress.IsCompleted,
		CompletedAt: toTimestamptz(progress.CompletedAt),
	})
	if err != nil {
		return err
	}
	*progress = toModelStudyCircleMemberProgress(row)
	return nil
}

// GetProgress はサークル全メンバーの進捗を返す。
func (r *studyCircleRepository) GetProgress(ctx context.Context, circleID uint) ([]model.StudyCircleMemberProgress, error) {
	rows, err := r.q.ListStudyCircleMemberProgressWithUser(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}
	progress := make([]model.StudyCircleMemberProgress, len(rows))
	for i, row := range rows {
		progress[i] = toModelStudyCircleMemberProgress(row.StudyCircleMemberProgress)
		user := toModelUser(row.User)
		progress[i].User = &user
	}
	return progress, nil
}

// CreateCheckin はチェックインをDBに保存する。
func (r *studyCircleRepository) CreateCheckin(ctx context.Context, checkin *model.StudyCircleCheckin) error {
	row, err := r.q.CreateStudyCircleCheckin(ctx, sqlcgen.CreateStudyCircleCheckinParams{
		CircleID: int64(checkin.CircleID),
		UserID:   int64(checkin.UserID),
		Date:     checkin.Date,
		Content:  checkin.Content,
	})
	if err != nil {
		return err
	}
	*checkin = toModelStudyCircleCheckin(row)
	return nil
}

// GetCheckins はチェックイン一覧を新しい順で返す。
func (r *studyCircleRepository) GetCheckins(ctx context.Context, circleID uint) ([]model.StudyCircleCheckin, error) {
	rows, err := r.q.ListStudyCircleCheckinsWithUser(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}
	checkins := make([]model.StudyCircleCheckin, len(rows))
	for i, row := range rows {
		checkins[i] = toModelStudyCircleCheckin(row.StudyCircleCheckin)
		user := toModelUser(row.User)
		checkins[i].User = &user
	}
	return checkins, nil
}

// HasCheckedInToday は今日すでにチェックイン済みかを返す。
func (r *studyCircleRepository) HasCheckedInToday(ctx context.Context, circleID, userID uint) (bool, error) {
	today := time.Now().Format("2006-01-02")
	count, err := r.q.CountStudyCircleCheckinsToday(ctx, sqlcgen.CountStudyCircleCheckinsTodayParams{
		CircleID: int64(circleID),
		UserID:   int64(userID),
		Date:     today,
	})
	return count > 0, err
}

// GetStreakRanking はサークル内メンバーのストリークランキングを返す（連続日数の降順）。
func (r *studyCircleRepository) GetStreakRanking(ctx context.Context, circleID uint) ([]model.CircleMemberStreak, error) {
	memberRows, err := r.q.ListStudyCircleMembersWithUserByCircle(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}

	// メンバーごとに 1 クエリずつ発行すると N+1 になるため、サークル分のチェックインを
	// 1 回で引いてメモリ上で user_id ごとにまとめる。date の降順は calculateCheckinStreak の前提。
	checkinRows, err := r.q.ListStudyCircleCheckinDatesByCircle(ctx, int64(circleID))
	if err != nil {
		return nil, err
	}

	datesByUser := make(map[uint][]string, len(memberRows))
	for _, checkin := range checkinRows {
		userID := uint(checkin.UserID)
		datesByUser[userID] = append(datesByUser[userID], checkin.Date)
	}

	var results []model.CircleMemberStreak
	for _, row := range memberRows {
		member := toModelStudyCircleMember(row.StudyCircleMember)
		user := toModelUser(row.User)
		dates := datesByUser[member.UserID]

		results = append(results, model.CircleMemberStreak{
			UserID:        member.UserID,
			UserName:      user.Name,
			AvatarURL:     user.AvatarURL,
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

	total, err := r.q.CountSearchStudyCircles(ctx, pattern)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.SearchStudyCircles(ctx, sqlcgen.SearchStudyCirclesParams{
		Name:   pattern,
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	circles := make([]model.StudyCircle, len(rows))
	for i, row := range rows {
		circles[i] = toModelStudyCircle(row.StudyCircle)
		owner := toModelUser(row.User)
		circles[i].Owner = &owner
		if err := attachStudyCircleMembers(ctx, r.q, &circles[i]); err != nil {
			return nil, 0, err
		}
	}
	return circles, total, nil
}
