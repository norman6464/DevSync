import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useQuestions } from '../useQuestions';
import { getQuestions, createQuestion, deleteQuestion } from '../../api/qa';
import type { Question } from '../../types/qa';
import toast from 'react-hot-toast';

vi.mock('../../api/qa', () => ({
  getQuestions: vi.fn(),
  searchQuestions: vi.fn(),
  createQuestion: vi.fn(),
  updateQuestion: vi.fn(),
  deleteQuestion: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn(), success: vi.fn() },
}));

const mockQuestions = [
  { id: 1, title: 'Goの並行処理', body: '...', tags: '["Go"]', vote_count: 3, answer_count: 2, is_solved: true, created_at: '2026-01-01', updated_at: '2026-01-01' },
  { id: 2, title: 'React Hooks', body: '...', tags: '["React"]', vote_count: 1, answer_count: 0, is_solved: false, created_at: '2026-01-02', updated_at: '2026-01-02' },
  { id: 3, title: 'Docker入門', body: '...', tags: '["Docker"]', vote_count: 5, answer_count: 3, is_solved: false, created_at: '2026-01-03', updated_at: '2026-01-03' },
  { id: 4, title: 'TypeScript型', body: '...', tags: '["TypeScript"]', vote_count: 0, answer_count: 1, is_solved: true, created_at: '2026-01-04', updated_at: '2026-01-04' },
] as Question[];

describe('useQuestions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getQuestions).mockResolvedValue({ questions: mockQuestions, total: 4 });
  });

  it('初期状態で質問一覧が取得されること', async () => {
    const { result } = renderHook(() => useQuestions());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.questions).toHaveLength(4);
    expect(result.current.total).toBe(4);
    expect(result.current.sort).toBe('newest');
    expect(result.current.solvedFilter).toBe('all');
  });

  it('solvedFilterで解決済みのみフィルタリングされること', async () => {
    const { result } = renderHook(() => useQuestions());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setSolvedFilter('solved');
    });

    expect(result.current.questions).toHaveLength(2);
    expect(result.current.questions.every(q => q.is_solved)).toBe(true);
  });

  it('solvedFilterで未解決のみフィルタリングされること', async () => {
    const { result } = renderHook(() => useQuestions());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setSolvedFilter('unsolved');
    });

    expect(result.current.questions).toHaveLength(2);
    expect(result.current.questions.every(q => !q.is_solved)).toBe(true);
  });

  it('質問作成が成功すること', async () => {
    const newQuestion = { id: 10, title: '新質問', body: 'test', tags: '[]', vote_count: 0, answer_count: 0, is_solved: false, created_at: '2026-02-01', updated_at: '2026-02-01' } as Question;
    vi.mocked(createQuestion).mockResolvedValue(newQuestion);

    const { result } = renderHook(() => useQuestions());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createQuestion({ title: '新質問', body: 'test', tags: '[]' });
    });

    expect(created).toEqual(newQuestion);
    expect(toast.success).toHaveBeenCalled();
    expect(result.current.questions.some(q => q.id === 10)).toBe(true);
  });

  it('質問作成失敗時にエラートーストが表示されnullが返ること', async () => {
    vi.mocked(createQuestion).mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useQuestions());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createQuestion({ title: 'テスト', body: 'test', tags: '[]' });
    });

    expect(created).toBeNull();
    expect(toast.error).toHaveBeenCalled();
  });

  it('質問削除が成功すること', async () => {
    vi.mocked(deleteQuestion).mockResolvedValue(undefined);

    const { result } = renderHook(() => useQuestions());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let success: boolean | undefined;
    await act(async () => {
      success = await result.current.deleteQuestion(mockQuestions[0] as never);
    });

    expect(success).toBe(true);
    expect(result.current.questions.find(q => q.id === 1)).toBeUndefined();
    expect(toast.success).toHaveBeenCalled();
  });

  it('ソート変更でページが0にリセットされること', async () => {
    const { result } = renderHook(() => useQuestions());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setPage(2);
    });
    expect(result.current.page).toBe(2);

    act(() => {
      result.current.setSort('votes');
    });
    expect(result.current.page).toBe(0);
    expect(result.current.sort).toBe('votes');
  });

  it('solvedFilter変更でページが0にリセットされること', async () => {
    const { result } = renderHook(() => useQuestions());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    act(() => {
      result.current.setPage(3);
    });
    expect(result.current.page).toBe(3);

    act(() => {
      result.current.setSolvedFilter('solved');
    });
    expect(result.current.page).toBe(0);
  });
});
