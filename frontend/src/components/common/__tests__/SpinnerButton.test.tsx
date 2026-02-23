import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SpinnerButton from '../SpinnerButton';

describe('SpinnerButton', () => {
  it('ラベルが表示される', () => {
    render(<SpinnerButton onClick={() => {}}>送信</SpinnerButton>);
    expect(screen.getByText('送信')).toBeInTheDocument();
  });

  it('クリックでonClickが呼ばれる', async () => {
    const onClick = vi.fn();
    const user = userEvent.setup();
    render(<SpinnerButton onClick={onClick}>送信</SpinnerButton>);
    await user.click(screen.getByRole('button'));
    expect(onClick).toHaveBeenCalled();
  });

  it('ローディング中にスピナーが表示される', () => {
    render(<SpinnerButton onClick={() => {}} loading>送信</SpinnerButton>);
    expect(screen.getByTestId('spinner')).toBeInTheDocument();
  });

  it('ローディング中はdisabled', () => {
    render(<SpinnerButton onClick={() => {}} loading>送信</SpinnerButton>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('ローディング中にラベルテキストが表示される', () => {
    render(<SpinnerButton onClick={() => {}} loading loadingText="送信中...">送信</SpinnerButton>);
    expect(screen.getByText('送信中...')).toBeInTheDocument();
  });

  it('primaryバリアントのスタイルが適用される', () => {
    const { container } = render(<SpinnerButton onClick={() => {}} variant="primary">送信</SpinnerButton>);
    expect(container.querySelector('.bg-blue-600')).toBeInTheDocument();
  });

  it('secondaryバリアントのスタイルが適用される', () => {
    const { container } = render(<SpinnerButton onClick={() => {}} variant="secondary">送信</SpinnerButton>);
    expect(container.querySelector('.bg-gray-700')).toBeInTheDocument();
  });

  it('dangerバリアントのスタイルが適用される', () => {
    const { container } = render(<SpinnerButton onClick={() => {}} variant="danger">送信</SpinnerButton>);
    expect(container.querySelector('.bg-red-600')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<SpinnerButton onClick={() => {}} size="sm">送信</SpinnerButton>);
    expect(container.querySelector('.px-3')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<SpinnerButton onClick={() => {}} size="lg">送信</SpinnerButton>);
    expect(container.querySelector('.px-6')).toBeInTheDocument();
  });

  it('disabled状態が適用される', () => {
    render(<SpinnerButton onClick={() => {}} disabled>送信</SpinnerButton>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<SpinnerButton onClick={() => {}} className="custom-class">送信</SpinnerButton>);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
