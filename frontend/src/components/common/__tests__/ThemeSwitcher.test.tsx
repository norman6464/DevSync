import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ThemeSwitcher from '../ThemeSwitcher';

describe('ThemeSwitcher', () => {
  it('3つのテーマオプションが表示される', () => {
    render(<ThemeSwitcher value="dark" onChange={() => {}} />);
    expect(screen.getByLabelText('ライト')).toBeInTheDocument();
    expect(screen.getByLabelText('ダーク')).toBeInTheDocument();
    expect(screen.getByLabelText('システム')).toBeInTheDocument();
  });

  it('現在のテーマがハイライトされる', () => {
    render(<ThemeSwitcher value="dark" onChange={() => {}} />);
    const darkBtn = screen.getByLabelText('ダーク');
    expect(darkBtn.className).toContain('bg-blue-600');
  });

  it('クリックでonChangeが呼ばれる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<ThemeSwitcher value="dark" onChange={onChange} />);
    await user.click(screen.getByLabelText('ライト'));
    expect(onChange).toHaveBeenCalledWith('light');
  });

  it('システムモードを選択できる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<ThemeSwitcher value="dark" onChange={onChange} />);
    await user.click(screen.getByLabelText('システム'));
    expect(onChange).toHaveBeenCalledWith('system');
  });

  it('ラベルが表示される', () => {
    render(<ThemeSwitcher value="dark" onChange={() => {}} label="テーマ" />);
    expect(screen.getByText('テーマ')).toBeInTheDocument();
  });

  it('disabled状態で全ボタンが無効', () => {
    render(<ThemeSwitcher value="dark" onChange={() => {}} disabled />);
    const buttons = screen.getAllByRole('button');
    buttons.forEach((btn) => expect(btn).toBeDisabled());
  });

  it('各ボタンにアイコンが表示される', () => {
    const { container } = render(<ThemeSwitcher value="dark" onChange={() => {}} />);
    const svgs = container.querySelectorAll('svg');
    expect(svgs.length).toBe(3);
  });

  it('非アクティブボタンはハイライトされない', () => {
    render(<ThemeSwitcher value="dark" onChange={() => {}} />);
    const lightBtn = screen.getByLabelText('ライト');
    expect(lightBtn.className).not.toContain('bg-blue-600');
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<ThemeSwitcher value="dark" onChange={() => {}} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
