import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import ShareButtons from '../ShareButtons';

const defaultProps = {
  onDownload: vi.fn(),
  onCopyLink: vi.fn(),
  onShareTwitter: vi.fn(),
  onShareLinkedIn: vi.fn(),
  generating: false,
};

describe('ShareButtons', () => {
  it('role="group"のコンテナを表示する', () => {
    render(<ShareButtons {...defaultProps} />);
    expect(screen.getByRole('group', { name: 'プロフィールを共有' })).toBeInTheDocument();
  });

  it('4つのシェアボタンを表示する', () => {
    render(<ShareButtons {...defaultProps} />);
    expect(screen.getAllByRole('button')).toHaveLength(4);
  });

  it('ダウンロードボタンクリックでonDownloadが呼ばれる', () => {
    const onDownload = vi.fn();
    render(<ShareButtons {...defaultProps} onDownload={onDownload} />);
    fireEvent.click(screen.getByLabelText('ダウンロード'));
    expect(onDownload).toHaveBeenCalledOnce();
  });

  it('リンクコピーボタンクリックでonCopyLinkが呼ばれる', () => {
    const onCopyLink = vi.fn();
    render(<ShareButtons {...defaultProps} onCopyLink={onCopyLink} />);
    fireEvent.click(screen.getByLabelText('リンクをコピー'));
    expect(onCopyLink).toHaveBeenCalledOnce();
  });

  it('TwitterボタンクリックでonShareTwitterが呼ばれる', () => {
    const onShareTwitter = vi.fn();
    render(<ShareButtons {...defaultProps} onShareTwitter={onShareTwitter} />);
    fireEvent.click(screen.getByLabelText('Twitterで共有'));
    expect(onShareTwitter).toHaveBeenCalledOnce();
  });

  it('LinkedInボタンクリックでonShareLinkedInが呼ばれる', () => {
    const onShareLinkedIn = vi.fn();
    render(<ShareButtons {...defaultProps} onShareLinkedIn={onShareLinkedIn} />);
    fireEvent.click(screen.getByLabelText('LinkedInで共有'));
    expect(onShareLinkedIn).toHaveBeenCalledOnce();
  });

  it('generating=trueでダウンロードボタンが無効になる', () => {
    render(<ShareButtons {...defaultProps} generating={true} />);
    expect(screen.getByLabelText('ダウンロード')).toBeDisabled();
  });

  it('generating=trueでダウンロードボタンのテキストがローディングに変わる', () => {
    render(<ShareButtons {...defaultProps} generating={true} />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('generating=falseでダウンロードボタンのテキストが通常表示', () => {
    render(<ShareButtons {...defaultProps} generating={false} />);
    const downloadButton = screen.getByLabelText('ダウンロード');
    expect(downloadButton).toBeEnabled();
    expect(downloadButton).toHaveTextContent('ダウンロード');
  });

  it('TwitterとLinkedInのラベルテキストを表示する', () => {
    render(<ShareButtons {...defaultProps} />);
    expect(screen.getByText('Twitter')).toBeInTheDocument();
    expect(screen.getByText('LinkedIn')).toBeInTheDocument();
  });
});
