package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteOwnerOf はノートの所有者 ID を返す。
func noteOwnerOf(n *model.Note) uint { return n.UserID }

// requireOwnedNote は指定ノートを取得し、userID が所有者であることを検証する。
//
// 不在のときも DB 障害のときも同じ 404 を返す。これは移行前の挙動をそのまま引き継いだもので、
// 「不在」と「DB 障害」を区別する改善は移行済みスライスをまとめて扱う別チケットの範囲。
func requireOwnedNote(
	ctx context.Context,
	notes repository.NoteReader,
	noteID, userID uint,
	notFoundMessage string,
) (*model.Note, error) {
	note, err := notes.FindByID(ctx, noteID)
	if err != nil || note == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, notFoundMessage, err)
	}
	if note.UserID != userID {
		return nil, domain.ErrForbidden
	}
	return note, nil
}

// CreateNoteLinkUseCase はノート間リンクを作成する。
type CreateNoteLinkUseCase struct {
	links repository.NoteLinkRepository
	notes repository.NoteReader
}

// NewCreateNoteLinkUseCase は CreateNoteLinkUseCase を生成する。
func NewCreateNoteLinkUseCase(links repository.NoteLinkRepository, notes repository.NoteReader) *CreateNoteLinkUseCase {
	return &CreateNoteLinkUseCase{links: links, notes: notes}
}

// Execute はリンク元の所有権とリンク先の存在・重複を検査したうえでリンクを作成する。
func (uc *CreateNoteLinkUseCase) Execute(ctx context.Context, sourceNoteID, targetNoteID, userID uint) error {
	if sourceNoteID == targetNoteID {
		return domain.NewError(domain.ErrCodeValidation, "同じノートへのリンクは作成できません", nil)
	}

	if _, err := requireOwnedNote(ctx, uc.notes, sourceNoteID, userID, "ソースノートが見つかりません"); err != nil {
		return err
	}

	// リンク先は所有者でなくてよい。存在だけを確認する。
	target, err := uc.notes.FindByID(ctx, targetNoteID)
	if err != nil || target == nil {
		return domain.NewError(domain.ErrCodeNotFound, "リンク先のノートが見つかりません", err)
	}

	exists, err := uc.links.Exists(ctx, sourceNoteID, targetNoteID)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewError(domain.ErrCodeValidation, "このリンクは既に存在します", nil)
	}

	return uc.links.Create(ctx, &model.NoteLink{
		SourceNoteID: sourceNoteID,
		TargetNoteID: targetNoteID,
	})
}

// ListNoteLinksUseCase は指定ノートからのリンク一覧を取得する。
type ListNoteLinksUseCase struct {
	links repository.NoteLinkRepository
}

// NewListNoteLinksUseCase は ListNoteLinksUseCase を生成する。
func NewListNoteLinksUseCase(links repository.NoteLinkRepository) *ListNoteLinksUseCase {
	return &ListNoteLinksUseCase{links: links}
}

// Execute はリンク一覧を返す。所有権は検証しない（移行前の挙動を維持している）。
func (uc *ListNoteLinksUseCase) Execute(ctx context.Context, sourceNoteID uint) ([]model.NoteLink, error) {
	return uc.links.FindBySourceNoteID(ctx, sourceNoteID)
}

// ListNoteBacklinksUseCase は指定ノートへのリンク一覧（バックリンク）を取得する。
type ListNoteBacklinksUseCase struct {
	links repository.NoteLinkRepository
}

// NewListNoteBacklinksUseCase は ListNoteBacklinksUseCase を生成する。
func NewListNoteBacklinksUseCase(links repository.NoteLinkRepository) *ListNoteBacklinksUseCase {
	return &ListNoteBacklinksUseCase{links: links}
}

// Execute はバックリンク一覧を返す。所有権は検証しない（移行前の挙動を維持している）。
func (uc *ListNoteBacklinksUseCase) Execute(ctx context.Context, targetNoteID uint) ([]model.NoteLink, error) {
	return uc.links.FindByTargetNoteID(ctx, targetNoteID)
}

// DeleteNoteLinkUseCase はノート間リンクを削除する。
type DeleteNoteLinkUseCase struct {
	links repository.NoteLinkRepository
	notes repository.NoteReader
}

// NewDeleteNoteLinkUseCase は DeleteNoteLinkUseCase を生成する。
func NewDeleteNoteLinkUseCase(links repository.NoteLinkRepository, notes repository.NoteReader) *DeleteNoteLinkUseCase {
	return &DeleteNoteLinkUseCase{links: links, notes: notes}
}

// Execute はリンク元の所有権を検証したうえでリンクを削除する。
func (uc *DeleteNoteLinkUseCase) Execute(ctx context.Context, sourceNoteID, targetNoteID, userID uint) error {
	if _, err := requireOwnedNote(ctx, uc.notes, sourceNoteID, userID, "ソースノートが見つかりません"); err != nil {
		return err
	}
	return uc.links.Delete(ctx, sourceNoteID, targetNoteID)
}

// GetNoteLinkStatsUseCase はノートのリンク統計を取得する。
type GetNoteLinkStatsUseCase struct {
	links repository.NoteLinkRepository
	notes repository.NoteReader
}

// NewGetNoteLinkStatsUseCase は GetNoteLinkStatsUseCase を生成する。
func NewGetNoteLinkStatsUseCase(links repository.NoteLinkRepository, notes repository.NoteReader) *GetNoteLinkStatsUseCase {
	return &GetNoteLinkStatsUseCase{links: links, notes: notes}
}

// Execute はノートの所有者に対してフォワードリンク数とバックリンク数を返す。
//
// リンクの作成・削除と違い、ノートが不在のときは 404 ではなく 500 になる。
// 移行前も所有権 helper が取得エラーをそのまま返しており、その挙動を維持している。
func (uc *GetNoteLinkStatsUseCase) Execute(ctx context.Context, noteID, userID uint) (*model.NoteLinkStats, error) {
	if _, err := ensureOwner(ctx, uc.notes.FindByID, noteID, userID, noteOwnerOf); err != nil {
		return nil, err
	}

	forwardCount, err := uc.links.CountBySourceNoteID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	backlinkCount, err := uc.links.CountByTargetNoteID(ctx, noteID)
	if err != nil {
		return nil, err
	}

	return &model.NoteLinkStats{
		NoteID:           noteID,
		ForwardLinkCount: forwardCount,
		BacklinkCount:    backlinkCount,
	}, nil
}
