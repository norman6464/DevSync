import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import OnboardingProfileStep from '../OnboardingProfileStep';

const defaultProps = {
  name: '',
  setName: vi.fn(),
  bio: '',
  setBio: vi.fn(),
  saving: false,
  onSave: vi.fn(),
  onSkip: vi.fn(),
};

describe('OnboardingProfileStep', () => {
  it('タイトルと説明が表示される', () => {
    render(<OnboardingProfileStep {...defaultProps} />);
    expect(screen.getByText('DevSyncへようこそ！')).toBeInTheDocument();
    expect(screen.getByText('まずはプロフィールを設定しましょう。')).toBeInTheDocument();
  });

  it('名前と自己紹介の入力欄が表示される', () => {
    render(<OnboardingProfileStep {...defaultProps} />);
    expect(screen.getByPlaceholderText('表示名を入力')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('自己紹介を入力')).toBeInTheDocument();
  });

  it('名前の変更がsetNameを呼び出す', () => {
    const setName = vi.fn();
    render(<OnboardingProfileStep {...defaultProps} setName={setName} />);
    fireEvent.change(screen.getByPlaceholderText('表示名を入力'), { target: { value: 'テスト太郎' } });
    expect(setName).toHaveBeenCalledWith('テスト太郎');
  });

  it('自己紹介の変更がsetBioを呼び出す', () => {
    const setBio = vi.fn();
    render(<OnboardingProfileStep {...defaultProps} setBio={setBio} />);
    fireEvent.change(screen.getByPlaceholderText('自己紹介を入力'), { target: { value: 'エンジニアです' } });
    expect(setBio).toHaveBeenCalledWith('エンジニアです');
  });

  it('次へボタンがonSaveを呼び出す', () => {
    const onSave = vi.fn();
    render(<OnboardingProfileStep {...defaultProps} onSave={onSave} />);
    fireEvent.click(screen.getByText('次へ'));
    expect(onSave).toHaveBeenCalledTimes(1);
  });

  it('スキップボタンがonSkipを呼び出す', () => {
    const onSkip = vi.fn();
    render(<OnboardingProfileStep {...defaultProps} onSkip={onSkip} />);
    fireEvent.click(screen.getByText('スキップ'));
    expect(onSkip).toHaveBeenCalledTimes(1);
  });

  it('saving中はボタンが無効化される', () => {
    render(<OnboardingProfileStep {...defaultProps} saving={true} />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('初期値が入力欄に表示される', () => {
    render(<OnboardingProfileStep {...defaultProps} name="初期名" bio="初期自己紹介" />);
    expect(screen.getByDisplayValue('初期名')).toBeInTheDocument();
    expect(screen.getByDisplayValue('初期自己紹介')).toBeInTheDocument();
  });
});
