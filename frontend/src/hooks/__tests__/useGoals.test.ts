import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import type { AxiosResponse } from 'axios';
import { useGoals } from '../useGoals';
import { getMyGoals, createGoal, updateGoal, deleteGoal, duplicateGoal, type LearningGoal, type GoalCategory } from '../../api/goals';
import toast from 'react-hot-toast';

vi.mock('../../api/goals', () => ({
  getMyGoals: vi.fn(),
  createGoal: vi.fn(),
  updateGoal: vi.fn(),
  deleteGoal: vi.fn(),
  duplicateGoal: vi.fn(),
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn(), success: vi.fn() },
}));

const mockGoals = [
  { id: 1, title: 'Go習得', status: 'active', category: 'programming', progress: 50 },
  { id: 2, title: 'AWS資格', status: 'completed', category: 'certification', progress: 100 },
  { id: 3, title: 'React復習', status: 'paused', category: 'programming', progress: 30 },
  { id: 4, title: 'Docker学習', status: 'active', category: 'devops', progress: 20 },
];

describe('useGoals', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getMyGoals).mockResolvedValue({ data: mockGoals } as unknown as AxiosResponse<LearningGoal[]>);
  });

  it('目標がステータス別に正しく分類されること', async () => {
    const { result } = renderHook(() => useGoals());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.goals).toHaveLength(4);
    expect(result.current.activeGoals).toHaveLength(2);
    expect(result.current.completedGoals).toHaveLength(1);
    expect(result.current.pausedGoals).toHaveLength(1);
  });

  it('目標作成が成功すること', async () => {
    const newGoal = { id: 10, title: '新目標', status: 'active', category: 'programming', progress: 0 };
    vi.mocked(createGoal).mockResolvedValue({ data: newGoal } as unknown as AxiosResponse<LearningGoal>);

    const { result } = renderHook(() => useGoals());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createGoal({ title: '新目標', description: '', category: 'programming' as unknown as GoalCategory });
    });

    expect(created).toEqual(newGoal);
    expect(toast.success).toHaveBeenCalled();
  });

  it('目標更新でprogress=100の場合に完了トーストが表示されること', async () => {
    const updated = { ...mockGoals[0], progress: 100, status: 'completed' };
    vi.mocked(updateGoal).mockResolvedValue({ data: updated } as unknown as AxiosResponse<LearningGoal>);

    const { result } = renderHook(() => useGoals());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.updateGoal(1, { progress: 100 });
    });

    expect(toast.success).toHaveBeenCalled();
  });

  it('目標削除が成功すること', async () => {
    vi.mocked(deleteGoal).mockResolvedValue(undefined as unknown as AxiosResponse);

    const { result } = renderHook(() => useGoals());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let success: boolean | undefined;
    await act(async () => {
      success = await result.current.deleteGoal(1);
    });

    expect(success).toBe(true);
    expect(result.current.goals.find(g => g.id === 1)).toBeUndefined();
  });

  it('目標複製が成功すること', async () => {
    const duplicated = { id: 20, title: 'Go習得 (コピー)', status: 'active', category: 'programming', progress: 0 };
    vi.mocked(duplicateGoal).mockResolvedValue({ data: duplicated } as unknown as AxiosResponse<LearningGoal>);

    const { result } = renderHook(() => useGoals());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.duplicateGoal(1);
    });

    expect(created).toEqual(duplicated);
    expect(result.current.goals.some(g => g.id === 20)).toBe(true);
  });

  it('作成失敗時にエラートーストが表示されnullが返ること', async () => {
    vi.mocked(createGoal).mockRejectedValue(new Error('fail'));

    const { result } = renderHook(() => useGoals());

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let created: unknown;
    await act(async () => {
      created = await result.current.createGoal({ title: 'テスト', description: '', category: 'programming' as unknown as GoalCategory });
    });

    expect(created).toBeNull();
    expect(toast.error).toHaveBeenCalled();
  });
});
