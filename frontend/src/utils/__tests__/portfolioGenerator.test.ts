import { describe, it, expect, vi } from 'vitest';
import { generatePortfolioHTML, downloadPortfolio } from '../portfolioGenerator';
import type { PortfolioData } from '../portfolioGenerator';
import type { User } from '../../types/user';
import type { LearningGoal } from '../../api/goals';

const baseUser: User = {
  id: 1,
  username: 'testuser',
  name: 'Test User',
  email: 'test@example.com',
  avatar_url: 'https://example.com/avatar.png',
  bio: 'A developer',
  github_id: 123,
  github_username: 'testuser',
  github_connected: true,
  spotify_connected: false,
  zenn_username: '',
  qiita_username: '',
  atcoder_username: '',
  paiza_rank: '',
  skills_languages: 'TypeScript,Go',
  skills_frameworks: 'React,Gin',
  onboarding_completed: true,
  email_weekly_report: false,
  email_language: 'ja',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

const baseData: PortfolioData = {
  user: baseUser,
  languages: [
    { language: 'TypeScript', bytes: 50000, repo_count: 5 },
    { language: 'Go', bytes: 30000, repo_count: 3 },
  ],
  repos: [
    { id: 1, github_repo_id: 100, name: 'my-project', full_name: 'testuser/my-project', description: 'A cool project', language: 'TypeScript', stars: 10, forks: 2, is_private: false },
  ],
  goals: [
    {
      id: 1,
      user_id: 1,
      title: 'Learn Rust',
      description: '',
      category: 'language',
      status: 'active',
      progress: 60,
      target_hours: 0,
      target_date: '2026-06-01',
      is_public: true,
      created_at: '2026-01-01',
      updated_at: '2026-01-01',
      completed_at: null,
    } satisfies LearningGoal,
  ],
  totalContributions: 500,
  followerCount: 20,
  followingCount: 15,
};

describe('portfolioGenerator', () => {
  describe('generatePortfolioHTML', () => {
    it('ユーザー名をHTMLに含める', () => {
      const html = generatePortfolioHTML(baseData, 'minimal');
      expect(html).toContain('Test User');
    });

    it('XSSペイロードがエスケープされる', () => {
      const data: PortfolioData = {
        ...baseData,
        user: { ...baseUser, name: '<script>alert("xss")</script>', bio: '<img onerror="hack()">' },
      };
      const html = generatePortfolioHTML(data, 'minimal');
      expect(html).not.toContain('<script>');
      expect(html).toContain('&lt;script&gt;');
      expect(html).not.toContain('<img onerror');
      expect(html).toContain('&lt;img onerror');
    });

    it('javascript: URLをサニタイズする', () => {
      const data: PortfolioData = {
        ...baseData,
        user: { ...baseUser, avatar_url: 'javascript:alert(1)' },
      };
      const html = generatePortfolioHTML(data, 'minimal');
      expect(html).not.toContain('javascript:');
    });

    it('正常なhttps URLはそのまま出力する', () => {
      const html = generatePortfolioHTML(baseData, 'minimal');
      expect(html).toContain('https://example.com/avatar.png');
    });

    it('minimalテーマのスタイルを含む', () => {
      const html = generatePortfolioHTML(baseData, 'minimal');
      expect(html).toContain('--bg: #ffffff');
    });

    it('modernテーマのスタイルを含む', () => {
      const html = generatePortfolioHTML(baseData, 'modern');
      expect(html).toContain('--bg: #0f172a');
    });

    it('gradientテーマのスタイルを含む', () => {
      const html = generatePortfolioHTML(baseData, 'gradient');
      expect(html).toContain('backdrop-filter');
    });

    it('言語セクションを表示する', () => {
      const html = generatePortfolioHTML(baseData, 'minimal');
      expect(html).toContain('TypeScript');
      expect(html).toContain('Go');
    });

    it('リポジトリ情報を表示する', () => {
      const html = generatePortfolioHTML(baseData, 'minimal');
      expect(html).toContain('my-project');
      expect(html).toContain('A cool project');
    });

    it('学習目標のプログレスを表示する', () => {
      const html = generatePortfolioHTML(baseData, 'minimal');
      expect(html).toContain('Learn Rust');
      expect(html).toContain('width: 60%');
    });

    it('スキルタグを表示する', () => {
      const html = generatePortfolioHTML(baseData, 'minimal');
      expect(html).toContain('TypeScript');
      expect(html).toContain('React');
      expect(html).toContain('Gin');
    });

    it('統計情報を表示する', () => {
      const html = generatePortfolioHTML(baseData, 'minimal');
      expect(html).toContain('500');
      expect(html).toContain('20');
      expect(html).toContain('15');
    });

    it('プログレス値を0-100の範囲にクランプする', () => {
      const data: PortfolioData = {
        ...baseData,
        goals: [
          { id: 1, user_id: 1, title: 'Over 100', description: '', category: 'language', status: 'active', progress: 150, target_date: '', created_at: '', updated_at: '' } as LearningGoal,
        ],
      };
      const html = generatePortfolioHTML(data, 'minimal');
      expect(html).toContain('width: 100%');
      expect(html).not.toContain('width: 150%');
    });
  });

  describe('downloadPortfolio', () => {
    it('危険なファイル名文字をサニタイズする', () => {
      const createObjectURL = vi.fn(() => 'blob:test');
      const revokeObjectURL = vi.fn();
      vi.stubGlobal('URL', { createObjectURL, revokeObjectURL });

      const mockElement = { href: '', download: '', click: vi.fn() };
      vi.spyOn(document, 'createElement').mockReturnValue(mockElement as unknown as HTMLAnchorElement);
      vi.spyOn(document.body, 'appendChild').mockImplementation(() => mockElement as unknown as HTMLAnchorElement);
      vi.spyOn(document.body, 'removeChild').mockImplementation(() => mockElement as unknown as HTMLAnchorElement);

      downloadPortfolio('<html></html>', 'my/file:name?.html');

      expect(mockElement.download).toBe('myfilename.html');
      expect(mockElement.click).toHaveBeenCalled();
      expect(revokeObjectURL).toHaveBeenCalled();

      vi.restoreAllMocks();
    });
  });
});
