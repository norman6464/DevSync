import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Alert from '../Alert';

describe('Alert', () => {
  it('メッセージが表示される', () => {
    render(<Alert variant="info" message="情報メッセージ" />);

    expect(screen.getByText('情報メッセージ')).toBeInTheDocument();
  });

  it('タイトルが表示される', () => {
    render(<Alert variant="info" title="タイトル" message="メッセージ" />);

    expect(screen.getByText('タイトル')).toBeInTheDocument();
  });

  it('infoバリアントのスタイルが適用される', () => {
    const { container } = render(<Alert variant="info" message="情報" />);

    expect(container.querySelector('.border-blue-500')).toBeInTheDocument();
  });

  it('successバリアントのスタイルが適用される', () => {
    const { container } = render(<Alert variant="success" message="成功" />);

    expect(container.querySelector('.border-green-500')).toBeInTheDocument();
  });

  it('warningバリアントのスタイルが適用される', () => {
    const { container } = render(<Alert variant="warning" message="警告" />);

    expect(container.querySelector('.border-yellow-500')).toBeInTheDocument();
  });

  it('errorバリアントのスタイルが適用される', () => {
    const { container } = render(<Alert variant="error" message="エラー" />);

    expect(container.querySelector('.border-red-500')).toBeInTheDocument();
  });

  it('アイコンが表示される', () => {
    const { container } = render(<Alert variant="info" message="情報" />);

    expect(container.querySelector('.lucide-info')).toBeInTheDocument();
  });

  it('閉じるボタンが表示される', () => {
    render(<Alert variant="info" message="情報" onClose={() => {}} />);

    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('閉じるボタンクリックでコールバックが呼ばれる', async () => {
    const onClose = vi.fn();
    const user = userEvent.setup();
    render(<Alert variant="info" message="情報" onClose={onClose} />);

    await user.click(screen.getByRole('button'));

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('閉じるボタンがない場合は表示されない', () => {
    render(<Alert variant="info" message="情報" />);

    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Alert variant="info" message="情報" className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('role=alertが設定されている', () => {
    render(<Alert variant="error" message="エラー" />);

    expect(screen.getByRole('alert')).toBeInTheDocument();
  });
});
