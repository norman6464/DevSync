import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import CircleSearchCard from '../CircleSearchCard';
import type { StudyCircle } from '../../../types/studyCircle';

const makeCircle = (overrides: Partial<StudyCircle> & { member_count?: number } = {}): StudyCircle =>
  ({
    id: 1,
    name: 'React 学習会',
    topic: 'React 入門',
    description: '毎週集まって進めます',
    status: 'active',
    member_count: 3,
    max_members: 5,
    created_at: '2026-01-01T00:00:00Z',
    ...overrides,
  }) as StudyCircle;

const renderCard = (circle: StudyCircle) =>
  render(
    <MemoryRouter>
      <CircleSearchCard circle={circle} />
    </MemoryRouter>,
  );

describe('サークル検索カード', () => {
  // ステータスは `studyCircle.<status>` を引く。複数形の名前空間を見ると
  // どのステータスも解決できず、生キーが画面に出る。
  it.each([
    ['active', 'アクティブ'],
    ['completed', '完了'],
    ['archived', 'アーカイブ'],
  ])('ステータス %s を日本語で表示する', (status, label) => {
    renderCard(makeCircle({ status: status as StudyCircle['status'] }));

    expect(screen.getByText(label)).toBeInTheDocument();
    expect(document.body.textContent).not.toContain('studyCircle');
  });

  it('メンバー数の見出しを日本語で表示する', () => {
    renderCard(makeCircle());

    expect(screen.getByText('メンバー')).toBeInTheDocument();
  });

  it('定員に達していれば満員と表示する', () => {
    renderCard(makeCircle({ member_count: 5, max_members: 5 }));

    expect(screen.getByText('満員')).toBeInTheDocument();
  });
});
