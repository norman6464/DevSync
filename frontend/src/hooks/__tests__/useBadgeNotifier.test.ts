import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useBadgeNotifier } from '../useBadgeNotifier';
import { notifyBadgeEarned } from '../../api/badges';
import type { BadgeResult } from '../../types/badge';

const mockShowToast = vi.fn();

vi.mock('../../contexts/ToastContext', () => ({
  useToast: () => ({ showToast: mockShowToast }),
}));

vi.mock('../../api/badges', () => ({
  notifyBadgeEarned: vi.fn().mockResolvedValue({}),
}));

const makeBadge = (overrides: Partial<BadgeResult>): BadgeResult => ({
  id: 'b0',
  name: 'badges.firstPost',
  description: '',
  category: 'post',
  earned: false,
  ...overrides,
});

const badge1 = makeBadge({ id: 'b1', name: 'badges.firstPost', earned: true });
const badge2 = makeBadge({ id: 'b2', name: 'badges.tenPosts', earned: true });
const badge3 = makeBadge({ id: 'b3', name: 'badges.hundredPosts', earned: false });

// localStorage mock
const store: Record<string, string> = {};
const mockGetItem = vi.fn((key: string) => store[key] ?? null);
const mockSetItem = vi.fn((key: string, value: string) => { store[key] = value; });
const mockRemoveItem = vi.fn((key: string) => { delete store[key]; });

Object.defineProperty(globalThis, 'localStorage', {
  value: {
    getItem: mockGetItem,
    setItem: mockSetItem,
    removeItem: mockRemoveItem,
  },
  writable: true,
});

describe('useBadgeNotifier', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    Object.keys(store).forEach((k) => delete store[k]);
  });

  it('空バッジ配列では何もしないこと', () => {
    renderHook(() => useBadgeNotifier([]));

    expect(mockShowToast).not.toHaveBeenCalled();
    expect(mockSetItem).not.toHaveBeenCalled();
  });

  it('初回ロード時はToast表示せずlocalStorageに保存すること', () => {
    renderHook(() => useBadgeNotifier([badge1, badge2, badge3]));

    expect(mockShowToast).not.toHaveBeenCalled();
    expect(mockSetItem).toHaveBeenCalledWith('devsync_earned_badges', JSON.stringify(['b1', 'b2']));
  });

  it('新規バッジ検出時にToast表示とAPI通知が行われること', () => {
    store['devsync_earned_badges'] = JSON.stringify(['b1']);

    renderHook(() => useBadgeNotifier([badge1, badge2, badge3]));

    expect(mockShowToast).toHaveBeenCalledTimes(1);
    expect(notifyBadgeEarned).toHaveBeenCalledWith('b2');
    expect(mockSetItem).toHaveBeenCalledWith('devsync_earned_badges', JSON.stringify(['b1', 'b2']));
  });

  it('既存バッジのみの場合はToast表示しないこと', () => {
    store['devsync_earned_badges'] = JSON.stringify(['b1', 'b2']);

    renderHook(() => useBadgeNotifier([badge1, badge2]));

    expect(mockShowToast).not.toHaveBeenCalled();
  });

  it('localStorageの不正JSONでエラーハンドリングすること', () => {
    store['devsync_earned_badges'] = 'invalid-json';

    renderHook(() => useBadgeNotifier([badge1]));

    expect(mockRemoveItem).toHaveBeenCalledWith('devsync_earned_badges');
    expect(mockShowToast).not.toHaveBeenCalled();
    expect(mockSetItem).toHaveBeenCalledWith('devsync_earned_badges', JSON.stringify(['b1']));
  });

  it('localStorageに不正な型が保存されている場合は無視すること', () => {
    store['devsync_earned_badges'] = JSON.stringify({ not: 'array' });

    renderHook(() => useBadgeNotifier([badge1]));

    // 不正な型は空配列として扱われ、初回ロードとして処理
    expect(mockShowToast).not.toHaveBeenCalled();
    expect(mockSetItem).toHaveBeenCalledWith('devsync_earned_badges', JSON.stringify(['b1']));
  });
});
