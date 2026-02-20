import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import ReportCharts from '../ReportCharts';
import { type ActivityReport } from '../../../api/reports';

const makeReport = (overrides: Partial<ActivityReport> = {}): ActivityReport => ({
  period: 'weekly',
  start_date: '2025-01-01',
  end_date: '2025-01-07',
  user_id: 1,
  total_contributions: 50,
  posts_created: 10,
  comments_created: 5,
  likes_received: 20,
  goals_completed: 3,
  goals_progress: 75,
  new_followers: 2,
  messages_exchanged: 15,
  daily_contributions: [
    { date: '2025-01-01', contributions: 10, posts: 2, comments: 1 },
    { date: '2025-01-02', contributions: 5, posts: 1, comments: 0 },
  ],
  top_languages: [
    { language: 'TypeScript', bytes: 50000, repos: 3 },
    { language: 'Go', bytes: 30000, repos: 2 },
  ],
  ...overrides,
});

const renderCharts = (report = makeReport(), maxContribution = 10) =>
  render(<ReportCharts report={report} maxContribution={maxContribution} />);

describe('ReportCharts', () => {
  it('日別アクティビティのタイトルが表示される', () => {
    renderCharts();
    expect(screen.getByText('日別アクティビティ')).toBeInTheDocument();
  });

  it('よく使う言語のタイトルが表示される', () => {
    renderCharts();
    expect(screen.getByText('よく使う言語')).toBeInTheDocument();
  });

  it('学習目標の進捗のタイトルが表示される', () => {
    renderCharts();
    expect(screen.getByText('学習目標の進捗')).toBeInTheDocument();
  });

  it('目標進捗パーセンテージが表示される', () => {
    renderCharts();
    expect(screen.getByText('75%')).toBeInTheDocument();
  });

  it('達成目標数が表示される', () => {
    renderCharts();
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  it('言語名が表示される', () => {
    renderCharts();
    expect(screen.getByText('TypeScript')).toBeInTheDocument();
    expect(screen.getByText('Go')).toBeInTheDocument();
  });

  it('言語バイト数がフォーマットされて表示される', () => {
    renderCharts();
    expect(screen.getByText(/48\.8 KB/)).toBeInTheDocument();
    expect(screen.getByText(/29\.3 KB/)).toBeInTheDocument();
  });

  it('daily_contributionsが空の場合はチャートが非表示', () => {
    renderCharts(makeReport({ daily_contributions: [] }));
    expect(screen.queryByText('日別アクティビティ')).not.toBeInTheDocument();
  });

  it('top_languagesが空の場合は言語セクションが非表示', () => {
    renderCharts(makeReport({ top_languages: [] }));
    expect(screen.queryByText('よく使う言語')).not.toBeInTheDocument();
  });

  it('goals_progressが0%の場合も表示される', () => {
    renderCharts(makeReport({ goals_progress: 0 }));
    expect(screen.getByText('0%')).toBeInTheDocument();
  });

  it('平均進捗ラベルが表示される', () => {
    renderCharts();
    expect(screen.getByText('アクティブな目標の平均進捗')).toBeInTheDocument();
  });
});
