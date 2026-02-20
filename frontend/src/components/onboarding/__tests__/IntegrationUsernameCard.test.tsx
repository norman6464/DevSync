import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import IntegrationUsernameCard from '../IntegrationUsernameCard';

const defaultProps = {
  icon: <div data-testid="icon">Z</div>,
  serviceName: 'Zenn',
  description: 'テスト説明',
  username: '',
  onUsernameChange: vi.fn(),
  placeholder: 'ユーザー名を入力',
  connecting: false,
  onConnect: vi.fn(),
  buttonClassName: 'bg-blue-600 text-white rounded-lg font-medium text-sm',
};

describe('IntegrationUsernameCard', () => {
  it('サービス名と説明を表示する', () => {
    render(<IntegrationUsernameCard {...defaultProps} />);
    expect(screen.getByText('Zenn')).toBeInTheDocument();
    expect(screen.getByText('テスト説明')).toBeInTheDocument();
  });

  it('アイコンを表示する', () => {
    render(<IntegrationUsernameCard {...defaultProps} />);
    expect(screen.getByTestId('icon')).toBeInTheDocument();
  });

  it('未接続時に入力フィールドと接続ボタンを表示する', () => {
    render(<IntegrationUsernameCard {...defaultProps} />);
    expect(screen.getByPlaceholderText('ユーザー名を入力')).toBeInTheDocument();
    expect(screen.getByText('連携する')).toBeInTheDocument();
  });

  it('接続済み時にユーザー名を表示する', () => {
    render(<IntegrationUsernameCard {...defaultProps} connectedUsername="testuser" />);
    expect(screen.getByText(/連携済み - @testuser/)).toBeInTheDocument();
    expect(screen.queryByPlaceholderText('ユーザー名を入力')).not.toBeInTheDocument();
  });

  it('入力変更でonUsernameChangeが呼ばれる', () => {
    const onUsernameChange = vi.fn();
    render(<IntegrationUsernameCard {...defaultProps} onUsernameChange={onUsernameChange} />);
    fireEvent.change(screen.getByPlaceholderText('ユーザー名を入力'), { target: { value: 'user1' } });
    expect(onUsernameChange).toHaveBeenCalledOnce();
  });

  it('接続ボタンクリックでonConnectが呼ばれる', () => {
    const onConnect = vi.fn();
    render(<IntegrationUsernameCard {...defaultProps} username="testuser" onConnect={onConnect} />);
    fireEvent.click(screen.getByText('連携する'));
    expect(onConnect).toHaveBeenCalledOnce();
  });

  it('username空の場合、接続ボタンがdisabledになる', () => {
    render(<IntegrationUsernameCard {...defaultProps} username="" />);
    expect(screen.getByText('連携する')).toBeDisabled();
  });

  it('connecting=trueの場合、ボタンがdisabledでローディング表示', () => {
    render(<IntegrationUsernameCard {...defaultProps} username="test" connecting={true} />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
    expect(screen.getByText('読み込み中...')).toBeDisabled();
  });

  it('入力フィールドのmaxLengthが50である', () => {
    render(<IntegrationUsernameCard {...defaultProps} />);
    expect(screen.getByPlaceholderText('ユーザー名を入力')).toHaveAttribute('maxLength', '50');
  });

  it('接続済み時にCheckCircleアイコン（SVG）が表示される', () => {
    const { container } = render(<IntegrationUsernameCard {...defaultProps} connectedUsername="user1" />);
    expect(container.querySelector('svg')).toBeInTheDocument();
  });
});
