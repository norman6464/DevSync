import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import CircleCheckinTab from '../CircleCheckinTab';
import type { StudyCircleCheckin } from '../../../types/studyCircle';

const mockCheckins: StudyCircleCheckin[] = [
  {
    id: 1, circle_id: 1, user_id: 1,
    user: { id: 1, name: 'Alice', avatar_url: '' },
    date: '2026-02-18', content: 'Reactの基礎を学んだ',
    created_at: '2026-02-18T10:00:00Z',
  },
  {
    id: 2, circle_id: 1, user_id: 2,
    user: { id: 2, name: 'Bob', avatar_url: '' },
    date: '2026-02-18', content: 'TypeScriptの型を復習した',
    created_at: '2026-02-18T11:00:00Z',
  },
];

describe('CircleCheckinTab', () => {
  it('チェックインフォームが表示される', () => {
    const onCheckin = vi.fn();
    render(<CircleCheckinTab checkins={[]} onCheckin={onCheckin} />);
    expect(screen.getByText('日次チェックイン')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('今日やったことを一言で...')).toBeInTheDocument();
  });

  it('チェックイン履歴が表示される', () => {
    const onCheckin = vi.fn();
    render(<CircleCheckinTab checkins={mockCheckins} onCheckin={onCheckin} />);
    expect(screen.getByText('Reactの基礎を学んだ')).toBeInTheDocument();
    expect(screen.getByText('TypeScriptの型を復習した')).toBeInTheDocument();
  });

  it('ユーザー名が表示される', () => {
    const onCheckin = vi.fn();
    render(<CircleCheckinTab checkins={mockCheckins} onCheckin={onCheckin} />);
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Bob')).toBeInTheDocument();
  });

  it('送信ボタンが空入力で無効化される', () => {
    const onCheckin = vi.fn();
    render(<CircleCheckinTab checkins={[]} onCheckin={onCheckin} />);
    const submitButton = screen.getByText('チェックイン');
    expect(submitButton).toBeDisabled();
  });

  it('テキスト入力後に送信ボタンが有効になる', () => {
    const onCheckin = vi.fn();
    render(<CircleCheckinTab checkins={[]} onCheckin={onCheckin} />);
    const input = screen.getByPlaceholderText('今日やったことを一言で...');
    fireEvent.change(input, { target: { value: 'テスト内容' } });
    const submitButton = screen.getByText('チェックイン');
    expect(submitButton).not.toBeDisabled();
  });

  it('チェックインが空の場合は空状態メッセージが表示される', () => {
    const onCheckin = vi.fn();
    render(<CircleCheckinTab checkins={[]} onCheckin={onCheckin} />);
    expect(screen.getByText('チェックインがまだありません')).toBeInTheDocument();
  });
});
