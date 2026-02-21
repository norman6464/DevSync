import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useNotes, useNoteSearch, useArchivedNotes } from '../useNotes';
import { getMyNotes, createNote, deleteNote, toggleFavorite, archiveNote, duplicateNote, searchNotes, getArchivedNotes, unarchiveNote } from '../../api/notes';
import toast from 'react-hot-toast';

vi.mock('../../api/notes', () => ({
  getMyNotes: vi.fn(),
  createNote: vi.fn(),
  updateNote: vi.fn(),
  deleteNote: vi.fn(),
  toggleFavorite: vi.fn(),
  archiveNote: vi.fn(),
  unarchiveNote: vi.fn(),
  getArchivedNotes: vi.fn(),
  searchNotes: vi.fn(),
  duplicateNote: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn(), success: vi.fn() },
}));

const mockNotes = [
  { id: 1, title: 'Go学習メモ', content: '...', is_favorite: false, is_archived: false, created_at: '2026-01-01', updated_at: '2026-01-01' },
  { id: 2, title: 'React Tips', content: '...', is_favorite: true, is_archived: false, created_at: '2026-01-02', updated_at: '2026-01-02' },
  { id: 3, title: 'Docker設定', content: '...', is_favorite: false, is_archived: false, created_at: '2026-01-03', updated_at: '2026-01-03' },
];

describe('useNotes', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getMyNotes).mockResolvedValue({ data: { data: mockNotes, total: 3, page: 1, limit: 20 } } as never);
    vi.stubGlobal('confirm', () => true);
  });

  it('初期状態でノート一覧が取得されること', async () => {
    const { result } = renderHook(() => useNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.notes).toHaveLength(3);
    expect(result.current.total).toBe(3);
  });

  it('ノート作成が成功すること', async () => {
    const newNote = { id: 10, title: '新ノート', content: 'test', is_favorite: false, is_archived: false, created_at: '2026-02-01', updated_at: '2026-02-01' };
    vi.mocked(createNote).mockResolvedValue({ data: newNote } as never);

    const { result } = renderHook(() => useNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createNote({ title: '新ノート', content: 'test' });
    });

    expect(created).toEqual(newNote);
    expect(toast.success).toHaveBeenCalled();
    expect(result.current.notes.some(n => n.id === 10)).toBe(true);
  });

  it('ノート作成失敗時にエラートーストが表示されnullが返ること', async () => {
    vi.mocked(createNote).mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createNote({ title: 'テスト', content: 'test' });
    });

    expect(created).toBeNull();
    expect(toast.error).toHaveBeenCalled();
  });

  it('ノート削除が成功すること', async () => {
    vi.mocked(deleteNote).mockResolvedValue(undefined as never);

    const { result } = renderHook(() => useNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let success: boolean | undefined;
    await act(async () => {
      success = await result.current.deleteNote(1);
    });

    expect(success).toBe(true);
    expect(result.current.notes.find(n => n.id === 1)).toBeUndefined();
    expect(toast.success).toHaveBeenCalled();
  });

  it('お気に入りトグルが動作すること', async () => {
    vi.mocked(toggleFavorite).mockResolvedValue(undefined as never);

    const { result } = renderHook(() => useNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.toggleFavorite(1);
    });

    expect(toggleFavorite).toHaveBeenCalledWith(1);
    expect(result.current.notes.find(n => n.id === 1)?.is_favorite).toBe(true);
  });

  it('アーカイブが成功し一覧から消えること', async () => {
    vi.mocked(archiveNote).mockResolvedValue(undefined as never);

    const { result } = renderHook(() => useNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.archiveNote(1);
    });

    expect(archiveNote).toHaveBeenCalledWith(1);
    expect(result.current.notes.find(n => n.id === 1)).toBeUndefined();
  });

  it('複製が成功し一覧に追加されること', async () => {
    const dup = { id: 20, title: 'Go学習メモ (コピー)', content: '...', is_favorite: false, is_archived: false, created_at: '2026-02-01', updated_at: '2026-02-01' };
    vi.mocked(duplicateNote).mockResolvedValue({ data: dup } as never);

    const { result } = renderHook(() => useNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let duplicated: unknown;
    await act(async () => {
      duplicated = await result.current.duplicateNote(1);
    });

    expect(duplicated).toEqual(dup);
    expect(result.current.notes.some(n => n.id === 20)).toBe(true);
  });

  it('favoriteNotesがお気に入りのみ返すこと', async () => {
    const { result } = renderHook(() => useNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.favoriteNotes).toHaveLength(1);
    expect(result.current.favoriteNotes[0].id).toBe(2);
  });
});

describe('useNoteSearch', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('検索クエリでノートが取得されること', async () => {
    vi.mocked(searchNotes).mockResolvedValue({ data: { data: [mockNotes[0]], total: 1, page: 1, limit: 20 } } as never);

    const { result } = renderHook(() => useNoteSearch('Go'));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.notes).toHaveLength(1);
    expect(searchNotes).toHaveBeenCalledWith('Go', 1, 20);
  });

  it('空クエリの場合APIを呼ばないこと', async () => {
    const { result } = renderHook(() => useNoteSearch(''));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.notes).toEqual([]);
    expect(searchNotes).not.toHaveBeenCalled();
  });
});

describe('useArchivedNotes', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.stubGlobal('confirm', () => true);
  });

  it('アーカイブ済みノートが取得されること', async () => {
    const archivedNotes = [{ id: 5, title: 'Archived', content: '...', is_favorite: false, is_archived: true, created_at: '2026-01-01', updated_at: '2026-01-01' }];
    vi.mocked(getArchivedNotes).mockResolvedValue({ data: { data: archivedNotes, total: 1, page: 1, limit: 20 } } as never);

    const { result } = renderHook(() => useArchivedNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.notes).toHaveLength(1);
    expect(result.current.total).toBe(1);
  });

  it('アーカイブ解除で一覧から消えること', async () => {
    const archivedNotes = [{ id: 5, title: 'Archived', content: '...', is_favorite: false, is_archived: true, created_at: '2026-01-01', updated_at: '2026-01-01' }];
    vi.mocked(getArchivedNotes).mockResolvedValue({ data: { data: archivedNotes, total: 1, page: 1, limit: 20 } } as never);
    vi.mocked(unarchiveNote).mockResolvedValue(undefined as never);

    const { result } = renderHook(() => useArchivedNotes());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.unarchiveNote(5);
    });

    expect(result.current.notes.find(n => n.id === 5)).toBeUndefined();
  });
});
