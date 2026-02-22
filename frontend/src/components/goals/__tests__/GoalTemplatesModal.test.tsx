import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import GoalTemplatesModal from '../GoalTemplatesModal';
import type { GoalCategory } from '../../../api/goals';

interface GoalTemplate {
  id: string;
  title: string;
  description: string;
  category: GoalCategory;
  estimatedDays: number;
}

const mockOnSelect = vi.fn();
const mockOnClose = vi.fn();

const defaultProps = {
  isOpen: true,
  onSelect: mockOnSelect,
  onClose: mockOnClose,
};

describe('GoalTemplatesModal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('モーダルが開いている時にタイトルが表示される', () => {
    render(<GoalTemplatesModal {...defaultProps} />);
    expect(screen.getByText('テンプレートから目標を作成')).toBeInTheDocument();
  });

  it('モーダルが閉じている時は何も表示されない', () => {
    render(<GoalTemplatesModal {...defaultProps} isOpen={false} />);
    expect(screen.queryByText('テンプレートから目標を作成')).not.toBeInTheDocument();
  });

  it('カテゴリフィルターが表示される', () => {
    render(<GoalTemplatesModal {...defaultProps} />);
    expect(screen.getByText('すべて')).toBeInTheDocument();
    expect(screen.getByText('言語')).toBeInTheDocument();
    expect(screen.getByText('フレームワーク')).toBeInTheDocument();
    expect(screen.getByText('スキル')).toBeInTheDocument();
    expect(screen.getByText('プロジェクト')).toBeInTheDocument();
  });

  it('テンプレートカードが表示される', () => {
    render(<GoalTemplatesModal {...defaultProps} />);

    // いくつかの一般的なテンプレートタイトルを確認
    expect(screen.getByText('TypeScript基礎マスター')).toBeInTheDocument();
    expect(screen.getByText('React実践')).toBeInTheDocument();
  });

  it('テンプレートの説明が表示される', () => {
    render(<GoalTemplatesModal {...defaultProps} />);

    // 説明文の一部を確認
    const descriptions = screen.getAllByText(/を学習/);
    expect(descriptions.length).toBeGreaterThan(0);
  });

  it('推定期間が表示される', () => {
    render(<GoalTemplatesModal {...defaultProps} />);

    // 推定期間表示の確認（例: "30日"）
    const durations = screen.getAllByText(/日$/);
    expect(durations.length).toBeGreaterThan(0);
  });

  it('テンプレート選択時にonSelectが呼ばれる', () => {
    render(<GoalTemplatesModal {...defaultProps} />);

    const templateButton = screen.getByText('TypeScript基礎マスター').closest('button');
    expect(templateButton).toBeInTheDocument();

    if (templateButton) {
      fireEvent.click(templateButton);
      expect(mockOnSelect).toHaveBeenCalledTimes(1);
      expect(mockOnSelect).toHaveBeenCalledWith(
        expect.objectContaining({
          title: 'TypeScript基礎マスター',
          category: 'language',
        })
      );
    }
  });

  it('閉じるボタンクリックでonCloseが呼ばれる', () => {
    render(<GoalTemplatesModal {...defaultProps} />);

    const closeButton = screen.getByLabelText('閉じる');
    fireEvent.click(closeButton);
    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('キャンセルボタンクリックでonCloseが呼ばれる', () => {
    render(<GoalTemplatesModal {...defaultProps} />);

    const cancelButton = screen.getByText('キャンセル');
    fireEvent.click(cancelButton);
    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('カテゴリフィルター選択で表示が絞り込まれる', () => {
    render(<GoalTemplatesModal {...defaultProps} />);

    // 「言語」フィルターをクリック
    const languageFilter = screen.getByText('言語');
    fireEvent.click(languageFilter);

    // 言語カテゴリのテンプレートが表示される
    expect(screen.getByText('TypeScript基礎マスター')).toBeInTheDocument();

    // 他のカテゴリのテンプレートは表示されない（例: React実践はフレームワーク）
    expect(screen.queryByText('React実践')).not.toBeInTheDocument();
  });

  it('各テンプレートカードにカテゴリアイコンが表示される', () => {
    const { container } = render(<GoalTemplatesModal {...defaultProps} />);

    // lucide-reactのアイコンが存在することを確認
    const icons = container.querySelectorAll('svg');
    expect(icons.length).toBeGreaterThan(0);
  });

  it('テンプレートが複数存在する', () => {
    render(<GoalTemplatesModal {...defaultProps} />);

    // 複数のテンプレートボタンが存在することを確認
    const templateCards = screen.getAllByRole('button').filter(
      (btn) => !btn.textContent?.includes('キャンセル') && !btn.getAttribute('aria-label')?.includes('閉じる')
    );

    expect(templateCards.length).toBeGreaterThan(5); // 少なくとも5つ以上のテンプレート
  });
});
