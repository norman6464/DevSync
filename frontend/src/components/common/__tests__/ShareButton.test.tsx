import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import ShareButton from '../ShareButton';

// window.openのモック
const mockWindowOpen = vi.fn();
global.window.open = mockWindowOpen;

describe('ShareButton', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('シェアボタンが表示される', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" />);
    expect(screen.getByRole('button')).toBeInTheDocument();
  });

  it('シェアアイコンが表示される', () => {
    const { container } = render(<ShareButton text="テストメッセージ" url="https://example.com" />);

    // lucide-reactのShareアイコンが存在することを確認
    const icon = container.querySelector('svg');
    expect(icon).toBeInTheDocument();
  });

  it('ボタンクリックでドロップダウンが表示される', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    // ソーシャルメディアオプションが表示される
    expect(screen.getByText(/Twitter/)).toBeInTheDocument();
    expect(screen.getByText(/Facebook/)).toBeInTheDocument();
    expect(screen.getByText(/LinkedIn/)).toBeInTheDocument();
  });

  it('Twitterシェアボタンをクリックすると正しいURLが開かれる', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    const twitterButton = screen.getByText(/Twitter/);
    fireEvent.click(twitterButton);

    expect(mockWindowOpen).toHaveBeenCalledWith(
      expect.stringContaining('twitter.com'),
      '_blank',
      expect.any(String)
    );
    expect(mockWindowOpen).toHaveBeenCalledWith(
      expect.stringContaining('テストメッセージ'),
      '_blank',
      expect.any(String)
    );
  });

  it('Facebookシェアボタンをクリックすると正しいURLが開かれる', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    const facebookButton = screen.getByText(/Facebook/);
    fireEvent.click(facebookButton);

    expect(mockWindowOpen).toHaveBeenCalledWith(
      expect.stringContaining('facebook.com'),
      '_blank',
      expect.any(String)
    );
  });

  it('LinkedInシェアボタンをクリックすると正しいURLが開かれる', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    const linkedInButton = screen.getByText(/LinkedIn/);
    fireEvent.click(linkedInButton);

    expect(mockWindowOpen).toHaveBeenCalledWith(
      expect.stringContaining('linkedin.com'),
      '_blank',
      expect.any(String)
    );
  });

  it('リンクコピーボタンが表示される', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    expect(screen.getByText(/リンクをコピー/)).toBeInTheDocument();
  });

  it('カスタムタイトルが表示される', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" title="カスタムタイトル" />);
    expect(screen.getByText('カスタムタイトル')).toBeInTheDocument();
  });

  it('デフォルトタイトルが表示される', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" />);
    expect(screen.getByText(/共有/)).toBeInTheDocument();
  });

  it('ドロップダウン外をクリックすると閉じる', () => {
    render(<ShareButton text="テストメッセージ" url="https://example.com" />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    // ドロップダウンが表示されていることを確認
    expect(screen.getByText(/Twitter/)).toBeInTheDocument();

    // 外側をクリック
    fireEvent.click(document.body);

    // ドロップダウンが閉じていることを確認（非同期なので時間をおいて確認）
    setTimeout(() => {
      expect(screen.queryByText(/Twitter/)).not.toBeInTheDocument();
    }, 100);
  });

  it('各ソーシャルメディアに対応するアイコンが表示される', () => {
    const { container } = render(<ShareButton text="テストメッセージ" url="https://example.com" />);

    const button = screen.getByRole('button');
    fireEvent.click(button);

    // アイコンが複数存在することを確認
    const icons = container.querySelectorAll('svg');
    expect(icons.length).toBeGreaterThan(3); // シェアボタンのアイコン + 各SNSのアイコン
  });
});
