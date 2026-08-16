import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import LearningStreakCard from '../LearningStreakCard';
import * as learningLogsApi from '../../../api/learningLogs';
import type { StreakInfo, CalendarEntry } from '../../../types/learningLog';

vi.mock('../../../api/learningLogs');
vi.mock('../../../store/authStore', () => ({
  useAuthStore: vi.fn(() => ({ user: { id: 1, name: 'Test User' } })),
}));

const mockStreakInfo: StreakInfo = {
  current_streak: 7,
  longest_streak: 15,
  total_days: 45,
  last_log_date: '2026-02-22T00:00:00Z',
};

const mockCalendarData: CalendarEntry[] = [
  { date: '2026-02-22', count: 3 },
  { date: '2026-02-21', count: 2 },
  { date: '2026-02-20', count: 1 },
  { date: '2026-02-19', count: 4 },
  { date: '2026-02-18', count: 2 },
  { date: '2026-02-17', count: 1 },
  { date: '2026-02-16', count: 3 },
];

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('LearningStreakCard', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(learningLogsApi.getStreakInfo).mockResolvedValue({
      data: mockStreakInfo,
    } as any);
    vi.mocked(learningLogsApi.getCalendarData).mockResolvedValue({
      data: mockCalendarData,
    } as any);
  });

  it('タイトルが表示される', async () => {
    renderWithRouter(<LearningStreakCard />);
    // ローディング完了後に見出しが描画される
    expect(
      await screen.findByRole('heading', { name: '学習ストリーク' })
    ).toBeInTheDocument();
  });

  it('現在のストリークが表示される', async () => {
    renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      expect(screen.getByText('7')).toBeInTheDocument();
      expect(screen.getByText(/日連続/)).toBeInTheDocument();
    });
  });

  it('最長ストリークが表示される', async () => {
    renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      expect(screen.getByText('15')).toBeInTheDocument();
      expect(screen.getByText(/最長記録/)).toBeInTheDocument();
    });
  });

  it('合計学習日数が表示される', async () => {
    renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      expect(screen.getByText('45')).toBeInTheDocument();
      expect(screen.getByText(/合計学習日数/)).toBeInTheDocument();
    });
  });

  it('ストリークがゼロの場合の表示', async () => {
    vi.mocked(learningLogsApi.getStreakInfo).mockResolvedValue({
      data: {
        current_streak: 0,
        longest_streak: 0,
        total_days: 0,
        last_log_date: '',
      },
    } as any);

    renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      expect(screen.getByText(/今日から学習を始めよう/)).toBeInTheDocument();
    });
  });

  it('カレンダーグリッドが表示される', async () => {
    const { container } = renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      // カレンダーグリッドの存在を確認
      const calendar = container.querySelector('[data-testid="streak-calendar"]');
      expect(calendar).toBeInTheDocument();
    });
  });

  it('学習日にはアクティブなマーカーが表示される', async () => {
    const { container } = renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      // 学習日のマーカー（例: bg-blue-500クラス）が存在することを確認
      const activeMarkers = container.querySelectorAll('.bg-blue-500, .bg-green-500, .bg-orange-500');
      expect(activeMarkers.length).toBeGreaterThan(0);
    });
  });

  it('ストリークレベルが表示される', async () => {
    renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      // ストリークレベル（例: "初心者", "中級者", "上級者"など）
      const levelText = screen.getByText(/レベル|初心者|中級者|上級者|マスター/);
      expect(levelText).toBeInTheDocument();
    });
  });

  it('炎のアイコンが表示される', async () => {
    const { container } = renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      // Flameアイコンのsvgが存在することを確認
      const icons = container.querySelectorAll('svg');
      expect(icons.length).toBeGreaterThan(0);
    });
  });

  it('ローディング状態が表示される', () => {
    vi.mocked(learningLogsApi.getStreakInfo).mockImplementation(
      () => new Promise(() => {}) as any
    );
    vi.mocked(learningLogsApi.getCalendarData).mockImplementation(
      () => new Promise(() => {}) as any
    );

    const { container } = renderWithRouter(<LearningStreakCard />);

    // ローディングスケルトンまたはスピナーの確認
    const loadingElements = container.querySelectorAll('.animate-pulse');
    expect(loadingElements.length).toBeGreaterThan(0);
  });

  it('高いストリークには特別な表示がある', async () => {
    vi.mocked(learningLogsApi.getStreakInfo).mockResolvedValue({
      data: {
        current_streak: 30,
        longest_streak: 30,
        total_days: 100,
        last_log_date: '2026-02-22T00:00:00Z',
      },
    } as any);

    renderWithRouter(<LearningStreakCard />);

    // 30日以上のストリークには特別なレベルバッジとマイルストーンメッセージが表示される
    expect(await screen.findByText('上級者')).toBeInTheDocument();
    expect(screen.getByText('30日達成！習慣化できています')).toBeInTheDocument();
    // 「30」は現在のストリークと最長記録の 2 箇所に表示される
    expect(screen.getAllByText('30')).toHaveLength(2);
  });

  it('直近30日のカレンダーが表示される', async () => {
    const { container } = renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      const calendar = container.querySelector('[data-testid="streak-calendar"]');
      expect(calendar).toBeInTheDocument();

      // 30日分のセルが存在することを確認
      const cells = container.querySelectorAll('[data-date]');
      expect(cells.length).toBeGreaterThanOrEqual(28); // 最低4週間分
    });
  });

  it('マイルストーンメッセージが表示される', async () => {
    renderWithRouter(<LearningStreakCard />);

    await waitFor(() => {
      // マイルストーンメッセージ（例: "7日達成！", "次は10日を目指そう"など）
      const messages = screen.queryAllByText(/達成|目指そう|素晴らしい|継続/);
      expect(messages.length).toBeGreaterThan(0);
    });
  });
});
