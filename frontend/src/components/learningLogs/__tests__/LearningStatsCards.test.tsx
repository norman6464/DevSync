import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import LearningStatsCards from '../LearningStatsCards';
import * as learningLogsApi from '../../../api/learningLogs';
import type { LearningLog } from '../../../types/learningLog';

vi.mock('../../../api/learningLogs');
vi.mock('../../../store/authStore', () => ({
  useAuthStore: vi.fn(() => ({ user: { id: 1, name: 'Test User' } })),
}));

const mockLearningLogs: LearningLog[] = [
  {
    id: 1,
    user_id: 1,
    title: 'React学習',
    content: 'Hooks を学んだ',
    category: 'coding',
    duration: 120,
    source: 'manual',
    is_favorite: false,
    created_at: '2026-02-22T10:00:00Z',
    updated_at: '2026-02-22T10:00:00Z',
  },
  {
    id: 2,
    user_id: 1,
    title: '技術書読書',
    content: 'Clean Code を読んだ',
    category: 'reading',
    duration: 60,
    source: 'manual',
    is_favorite: true,
    created_at: '2026-02-21T14:00:00Z',
    updated_at: '2026-02-21T14:00:00Z',
  },
  {
    id: 3,
    user_id: 1,
    title: 'Udemy受講',
    content: 'TypeScript コースを受講',
    category: 'course',
    duration: 90,
    source: 'manual',
    is_favorite: false,
    created_at: '2026-02-20T09:00:00Z',
    updated_at: '2026-02-20T09:00:00Z',
  },
];

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('LearningStatsCards', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(learningLogsApi.getMyLogs).mockResolvedValue({
      data: mockLearningLogs,
    } as Awaited<ReturnType<typeof learningLogsApi.getMyLogs>>);
  });

  it('タイトルが表示される', async () => {
    renderWithRouter(<LearningStatsCards />);
    // 見出し「学習統計」はロード完了後にのみ描画される
    expect(await screen.findByText('学習統計')).toBeInTheDocument();
  });

  it('総学習時間が表示される', async () => {
    renderWithRouter(<LearningStatsCards />);

    expect(await screen.findByText('総学習時間')).toBeInTheDocument();
    // 120 + 60 + 90 = 270分 = 4.5時間（カードは時間と分の両方を表示する）
    expect(screen.getByText('4.5')).toBeInTheDocument();
    expect(screen.getByText('270分')).toBeInTheDocument();
  });

  it('学習ログ数が表示される', async () => {
    renderWithRouter(<LearningStatsCards />);

    await waitFor(() => {
      expect(screen.getByText('3')).toBeInTheDocument();
      expect(screen.getByText(/学習ログ数/)).toBeInTheDocument();
    });
  });

  it('平均学習時間が表示される', async () => {
    renderWithRouter(<LearningStatsCards />);

    expect(await screen.findByText('平均学習時間')).toBeInTheDocument();
    // (120 + 60 + 90) / 3 = 90分 = 1.5時間（カードは時間と分/回の両方を表示する）
    expect(screen.getByText('1.5')).toBeInTheDocument();
    expect(screen.getByText('90分/回')).toBeInTheDocument();
  });

  it('カテゴリ別統計が表示される', async () => {
    renderWithRouter(<LearningStatsCards />);

    expect(await screen.findByText('カテゴリ別学習時間')).toBeInTheDocument();
    // 3 カテゴリすべてのラベルが表示される
    expect(screen.getByText('コーディング')).toBeInTheDocument();
    expect(screen.getByText('読書')).toBeInTheDocument();
    expect(screen.getByText('コース')).toBeInTheDocument();
  });

  it('各カテゴリの学習時間が表示される', async () => {
    renderWithRouter(<LearningStatsCards />);

    await waitFor(() => {
      // coding: 120分, reading: 60分, course: 90分
      const timeElements = screen.getAllByText(/分|時間/);
      expect(timeElements.length).toBeGreaterThan(0);
    });
  });

  it('ローディング状態が表示される', () => {
    vi.mocked(learningLogsApi.getMyLogs).mockImplementation(
      () => new Promise(() => {}) as ReturnType<typeof learningLogsApi.getMyLogs>
    );

    const { container } = renderWithRouter(<LearningStatsCards />);

    // ローディングスケルトンまたはスピナーの確認
    const loadingElements = container.querySelectorAll('.animate-pulse');
    expect(loadingElements.length).toBeGreaterThan(0);
  });

  it('学習ログがゼロの場合の表示', async () => {
    vi.mocked(learningLogsApi.getMyLogs).mockResolvedValue({
      data: [],
    } as unknown as Awaited<ReturnType<typeof learningLogsApi.getMyLogs>>);

    renderWithRouter(<LearningStatsCards />);

    expect(
      await screen.findByText('まだ学習ログがありません')
    ).toBeInTheDocument();
    // 総学習時間・学習ログ数・平均学習時間の 3 カードがいずれも 0 を表示する
    expect(screen.getAllByText('0')).toHaveLength(3);
  });

  it('カードが複数表示される', async () => {
    const { container } = renderWithRouter(<LearningStatsCards />);

    await waitFor(() => {
      // 統計カードが存在することを確認
      const cards = container.querySelectorAll('.bg-gray-900, .bg-gray-800');
      expect(cards.length).toBeGreaterThan(2);
    });
  });

  it('アイコンが表示される', async () => {
    const { container } = renderWithRouter(<LearningStatsCards />);

    await waitFor(() => {
      // lucide-reactのアイコンが存在することを確認
      const icons = container.querySelectorAll('svg');
      expect(icons.length).toBeGreaterThan(0);
    });
  });

  it('グリッドレイアウトで表示される', async () => {
    const { container } = renderWithRouter(<LearningStatsCards />);

    await waitFor(() => {
      const grid = container.querySelector('.grid');
      expect(grid).toBeInTheDocument();
    });
  });

  it('カテゴリ別の割合が表示される', async () => {
    renderWithRouter(<LearningStatsCards />);

    await waitFor(() => {
      // パーセンテージ表示を確認
      const percentages = screen.getAllByText(/%/);
      expect(percentages.length).toBeGreaterThan(0);
    });
  });
});
