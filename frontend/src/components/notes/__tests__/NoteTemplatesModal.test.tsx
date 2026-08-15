import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import NoteTemplatesModal from '../NoteTemplatesModal';

export interface NoteTemplate {
  id: string;
  title: string;
  content: string;
  category: string;
  tags: string;
}

const mockOnSelect = vi.fn();
const mockOnClose = vi.fn();

const defaultProps = {
  isOpen: true,
  onSelect: mockOnSelect,
  onClose: mockOnClose,
};

describe('NoteTemplatesModal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('モーダルが開いている時にタイトルが表示される', () => {
    render(<NoteTemplatesModal {...defaultProps} />);
    expect(screen.getByText('テンプレートからノートを作成')).toBeInTheDocument();
  });

  it('モーダルが閉じている時は何も表示されない', () => {
    render(<NoteTemplatesModal {...defaultProps} isOpen={false} />);
    expect(screen.queryByText('テンプレートからノートを作成')).not.toBeInTheDocument();
  });

  it('カテゴリフィルターが表示される', () => {
    render(<NoteTemplatesModal {...defaultProps} />);
    // カテゴリ名はカード内のラベルにも出るため、accessible name が
    // カテゴリ名と完全一致するフィルターボタンに絞って確認する
    expect(screen.getByRole('button', { name: 'すべて' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '学習ノート' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'プロジェクト' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '復習' })).toBeInTheDocument();
  });

  it('テンプレートカードが表示される', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    // いくつかの一般的なテンプレートタイトルを確認
    expect(screen.getByText('コーディング学習メモ')).toBeInTheDocument();
    expect(screen.getByText('読書ノート')).toBeInTheDocument();
  });

  it('テンプレートの説明が表示される', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    // カードの説明欄にはテンプレートのタグ文字列が表示される
    expect(screen.getByText('#学習,#コーディング')).toBeInTheDocument();
    expect(screen.getByText('#読書,#書籍')).toBeInTheDocument();
  });

  it('テンプレート選択時にonSelectが呼ばれる', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    const templateButton = screen.getByText('コーディング学習メモ').closest('button');
    expect(templateButton).toBeInTheDocument();

    if (templateButton) {
      fireEvent.click(templateButton);
      expect(mockOnSelect).toHaveBeenCalledTimes(1);
      expect(mockOnSelect).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'コーディング学習メモ',
          category: '学習ノート',
        })
      );
    }
  });

  it('閉じるボタンクリックでonCloseが呼ばれる', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    const closeButton = screen.getByLabelText('閉じる');
    fireEvent.click(closeButton);
    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('キャンセルボタンクリックでonCloseが呼ばれる', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    const cancelButton = screen.getByText('キャンセル');
    fireEvent.click(cancelButton);
    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('カテゴリフィルター選択で表示が絞り込まれる', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    // 「学習ノート」フィルターをクリック
    const learningFilter = screen.getByRole('button', { name: '学習ノート' });
    fireEvent.click(learningFilter);

    // 学習ノートカテゴリのテンプレートが表示される
    expect(screen.getByText('コーディング学習メモ')).toBeInTheDocument();

    // 他のカテゴリのテンプレートは表示されない
    expect(screen.queryByText('プロジェクト計画')).not.toBeInTheDocument();
  });

  it('各テンプレートカードにアイコンが表示される', () => {
    const { container } = render(<NoteTemplatesModal {...defaultProps} />);

    // lucide-reactのアイコンが存在することを確認
    const icons = container.querySelectorAll('svg');
    expect(icons.length).toBeGreaterThan(0);
  });

  it('テンプレートが複数存在する', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    // 複数のテンプレートボタンが存在することを確認
    const templateCards = screen.getAllByRole('button').filter(
      (btn) => !btn.textContent?.includes('キャンセル') && !btn.getAttribute('aria-label')?.includes('閉じる')
    );

    expect(templateCards.length).toBeGreaterThan(8); // 少なくとも8つ以上のテンプレート
  });

  it('テンプレートにタグが含まれている', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    // タグ表示を確認
    const tagElements = screen.getAllByText(/#/);
    expect(tagElements.length).toBeGreaterThan(0);
  });

  it('フィルター選択時にスタイルが変更される', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    const learningFilter = screen.getByRole('button', { name: '学習ノート' });

    // クリック前のスタイルを確認
    expect(learningFilter).not.toHaveClass('text-blue-400');

    // クリック
    fireEvent.click(learningFilter);

    // クリック後のスタイルを確認
    expect(learningFilter).toHaveClass('text-blue-400');
  });

  it('各テンプレートにカテゴリラベルが表示される', () => {
    render(<NoteTemplatesModal {...defaultProps} />);

    // カテゴリラベルの確認
    const categoryLabels = screen.getAllByText(/学習ノート|プロジェクト|復習/);
    expect(categoryLabels.length).toBeGreaterThan(0);
  });
});
