import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import LanguageSwitcher from '../LanguageSwitcher';

const languages = [
  { code: 'ja', label: '日本語', flag: '🇯🇵' },
  { code: 'en', label: 'English', flag: '🇺🇸' },
  { code: 'ko', label: '한국어', flag: '🇰🇷' },
];

describe('LanguageSwitcher', () => {
  it('現在の言語が表示される', () => {
    render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} />);
    expect(screen.getByText('日本語')).toBeInTheDocument();
  });

  it('クリックでドロップダウンが開く', async () => {
    const user = userEvent.setup();
    render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} />);
    await user.click(screen.getByText('日本語'));
    expect(screen.getByText('English')).toBeInTheDocument();
    expect(screen.getByText('한국어')).toBeInTheDocument();
  });

  it('言語選択でonChangeが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<LanguageSwitcher languages={languages} current="ja" onChange={onChange} />);
    await user.click(screen.getByText('日本語'));
    await user.click(screen.getByText('English'));
    expect(onChange).toHaveBeenCalledWith('en');
  });

  it('選択後にドロップダウンが閉じる', async () => {
    const user = userEvent.setup();
    render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} />);
    await user.click(screen.getByText('日本語'));
    await user.click(screen.getByText('English'));
    expect(screen.queryByText('한국어')).not.toBeInTheDocument();
  });

  it('フラグが表示される', () => {
    render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} />);
    expect(screen.getByText('🇯🇵')).toBeInTheDocument();
  });

  it('ドロップダウンでフラグが表示される', async () => {
    const user = userEvent.setup();
    render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} />);
    await user.click(screen.getByText('日本語'));
    expect(screen.getByText('🇺🇸')).toBeInTheDocument();
  });

  it('現在の言語がハイライトされる', async () => {
    const user = userEvent.setup();
    render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} />);
    await user.click(screen.getByText('日本語'));
    const items = screen.getAllByRole('option');
    expect(items[0].className).toContain('bg-gray-700');
  });

  it('disabled状態でクリックが無効', () => {
    render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} disabled />);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('アイコンが表示される', () => {
    const { container } = render(<LanguageSwitcher languages={languages} current="ja" onChange={() => {}} />);
    expect(container.querySelector('.lucide-globe')).toBeInTheDocument();
  });
});
