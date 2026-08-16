import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import LearningStatsCards from '../LearningStatsCards';
import LearningStreakCard from '../LearningStreakCard';
import i18n from '../../../i18n';
import * as learningLogsApi from '../../../api/learningLogs';
import type { LearningLog } from '../../../types/learningLog';

vi.mock('../../../api/learningLogs');
vi.mock('../../../store/authStore', () => ({
  useAuthStore: vi.fn((selector: (s: { user: { id: 1 } }) => unknown) => selector({ user: { id: 1 } })),
}));
vi.mock('../../../hooks/useStreakFreeze', () => ({
  useStreakFreeze: () => ({ freezeStatus: null, useFreeze: vi.fn() }),
}));

const log: LearningLog = {
  id: 1,
  user_id: 1,
  title: 'Go学習',
  content: '...',
  category: 'coding',
  duration: 60,
  source: 'manual',
  is_favorite: false,
  created_at: '2026-01-01',
  updated_at: '2026-01-01',
};

/** 画面テキストに日本語（ひらがな・カタカナ・漢字）が含まれるか。 */
const hasJapanese = (text: string) => /[぀-ヿ一-鿿]/.test(text);

describe('学習ログカードの言語切替', () => {
  beforeEach(async () => {
    vi.clearAllMocks();
    vi.mocked(learningLogsApi.getMyLogs).mockResolvedValue({
      data: [log],
    } as Awaited<ReturnType<typeof learningLogsApi.getMyLogs>>);
    vi.mocked(learningLogsApi.getStreakInfo).mockResolvedValue({
      data: { current_streak: 3, longest_streak: 5, total_days: 10 },
    } as Awaited<ReturnType<typeof learningLogsApi.getStreakInfo>>);
    vi.mocked(learningLogsApi.getCalendarData).mockResolvedValue({
      data: [],
    } as unknown as Awaited<ReturnType<typeof learningLogsApi.getCalendarData>>);
    await i18n.changeLanguage('en');
  });

  afterEach(async () => {
    await i18n.changeLanguage('ja');
  });

  // 部分的な i18n 適用のまま日本語ハードコードが残ると、en でも日本語が表示される。
  it('en で統計カードに日本語が表示されない', async () => {
    const { container } = render(<LearningStatsCards />);

    await screen.findByText('Learning stats');
    expect(hasJapanese(container.textContent ?? '')).toBe(false);
  });

  it('en でストリークカードに日本語が表示されない', async () => {
    const { container } = render(<LearningStreakCard />);

    await screen.findByText('Learning Streak');
    expect(hasJapanese(container.textContent ?? '')).toBe(false);
  });
});
