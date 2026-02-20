import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import OnboardingSkillsStep from '../OnboardingSkillsStep';

const defaultProps = {
  selectedLanguages: [] as string[],
  selectedFrameworks: [] as string[],
  toggleLanguage: vi.fn(),
  toggleFramework: vi.fn(),
  saving: false,
  onSave: vi.fn(),
  onBack: vi.fn(),
  onSkip: vi.fn(),
};

describe('OnboardingSkillsStep', () => {
  it('タイトルと説明が表示される', () => {
    render(<OnboardingSkillsStep {...defaultProps} />);
    expect(screen.getByText('スキルを選択')).toBeInTheDocument();
    expect(screen.getByText('使用している言語やフレームワークを選択してください。')).toBeInTheDocument();
  });

  it('言語ボタンが表示される', () => {
    render(<OnboardingSkillsStep {...defaultProps} />);
    expect(screen.getByText('javascript')).toBeInTheDocument();
    expect(screen.getByText('python')).toBeInTheDocument();
    expect(screen.getByText('go')).toBeInTheDocument();
  });

  it('言語クリックでtoggleLanguageが呼ばれる', () => {
    const toggleLanguage = vi.fn();
    render(<OnboardingSkillsStep {...defaultProps} toggleLanguage={toggleLanguage} />);
    fireEvent.click(screen.getByText('python'));
    expect(toggleLanguage).toHaveBeenCalledWith('python');
  });

  it('フレームワーククリックでtoggleFrameworkが呼ばれる', () => {
    const toggleFramework = vi.fn();
    render(<OnboardingSkillsStep {...defaultProps} toggleFramework={toggleFramework} />);
    fireEvent.click(screen.getByText('react'));
    expect(toggleFramework).toHaveBeenCalledWith('react');
  });

  it('選択済み言語のプレビューが表示される', () => {
    render(<OnboardingSkillsStep {...defaultProps} selectedLanguages={['javascript', 'python']} />);
    const imgs = screen.getAllByRole('img');
    const langImg = imgs.find(img => img.getAttribute('alt') === 'Selected languages');
    expect(langImg).toBeInTheDocument();
  });

  it('戻るボタンがonBackを呼び出す', () => {
    const onBack = vi.fn();
    render(<OnboardingSkillsStep {...defaultProps} onBack={onBack} />);
    fireEvent.click(screen.getByText('戻る'));
    expect(onBack).toHaveBeenCalledTimes(1);
  });

  it('スキップボタンがonSkipを呼び出す', () => {
    const onSkip = vi.fn();
    render(<OnboardingSkillsStep {...defaultProps} onSkip={onSkip} />);
    fireEvent.click(screen.getByText('スキップ'));
    expect(onSkip).toHaveBeenCalledTimes(1);
  });

  it('次へボタンがonSaveを呼び出す', () => {
    const onSave = vi.fn();
    render(<OnboardingSkillsStep {...defaultProps} onSave={onSave} />);
    fireEvent.click(screen.getByText('次へ'));
    expect(onSave).toHaveBeenCalledTimes(1);
  });
});
