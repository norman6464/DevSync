import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import OnboardingCompleteStep from '../OnboardingCompleteStep';

describe('OnboardingCompleteStep', () => {
  it('完了タイトルと説明が表示される', () => {
    render(<OnboardingCompleteStep saving={false} onComplete={vi.fn()} />);
    expect(screen.getByText('セットアップ完了！')).toBeInTheDocument();
    expect(screen.getByText('準備が整いました。ダッシュボードで活動を始めましょう。')).toBeInTheDocument();
  });

  it('ダッシュボードへボタンが表示される', () => {
    render(<OnboardingCompleteStep saving={false} onComplete={vi.fn()} />);
    expect(screen.getByText('ダッシュボードへ')).toBeInTheDocument();
  });

  it('ボタンクリックでonCompleteが呼ばれる', () => {
    const onComplete = vi.fn();
    render(<OnboardingCompleteStep saving={false} onComplete={onComplete} />);
    fireEvent.click(screen.getByText('ダッシュボードへ'));
    expect(onComplete).toHaveBeenCalledTimes(1);
  });

  it('saving中はローディング表示になる', () => {
    render(<OnboardingCompleteStep saving={true} onComplete={vi.fn()} />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
    expect(screen.queryByText('ダッシュボードへ')).not.toBeInTheDocument();
  });

  it('saving中はボタンが無効化される', () => {
    render(<OnboardingCompleteStep saving={true} onComplete={vi.fn()} />);
    const button = screen.getByRole('button');
    expect(button).toBeDisabled();
  });
});
