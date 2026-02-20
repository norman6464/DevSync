import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import RoadmapTemplatesSection from '../RoadmapTemplatesSection';
import type { Roadmap } from '../../../api/roadmaps';

const makeTemplate = (overrides: Partial<Roadmap> = {}): Roadmap => ({
  id: 1,
  user_id: 1,
  title: 'テストテンプレート',
  description: 'テンプレートの説明',
  category: 'language',
  is_public: true,
  is_template: true,
  step_count: 3,
  completed_step_count: 0,
  progress: 0,
  status: 'active',
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
  completed_at: null,
  steps: [
    { id: 1, roadmap_id: 1, title: 'ステップ1', description: '説明1', order_index: 0, is_completed: false, resource_url: '', created_at: '', updated_at: '' },
    { id: 2, roadmap_id: 1, title: 'ステップ2', description: '', order_index: 1, is_completed: false, resource_url: '', created_at: '', updated_at: '' },
  ],
  ...overrides,
});

const defaultProps = {
  templates: [makeTemplate()],
  showTemplates: false,
  expandedTemplate: null,
  creating: false,
  toggleTemplates: vi.fn(),
  toggleExpandedTemplate: vi.fn(),
  handleUseTemplate: vi.fn(),
};

const renderSection = (props: Partial<typeof defaultProps> = {}) =>
  render(
    <MemoryRouter>
      <RoadmapTemplatesSection {...defaultProps} {...props} />
    </MemoryRouter>
  );

describe('RoadmapTemplatesSection', () => {
  it('テンプレートが空の場合何も表示しない', () => {
    const { container } = renderSection({ templates: [] });
    expect(container.innerHTML).toBe('');
  });

  it('テンプレートヘッダーが表示される', () => {
    renderSection();
    expect(screen.getByText('おすすめテンプレート')).toBeInTheDocument();
  });

  it('テンプレート数バッジが表示される', () => {
    renderSection({ templates: [makeTemplate(), makeTemplate({ id: 2, title: 'テンプレート2' })] });
    expect(screen.getByText(/2/)).toBeInTheDocument();
  });

  it('トグルボタンクリックでtoggleTemplatesが呼ばれる', () => {
    const toggleTemplates = vi.fn();
    renderSection({ toggleTemplates });
    fireEvent.click(screen.getByText('おすすめテンプレート'));
    expect(toggleTemplates).toHaveBeenCalledOnce();
  });

  it('showTemplates=trueでテンプレート一覧が表示される', () => {
    renderSection({ showTemplates: true });
    expect(screen.getByText('テストテンプレート')).toBeInTheDocument();
    expect(screen.getByText('テンプレートの説明')).toBeInTheDocument();
  });

  it('showTemplates=falseでテンプレート一覧が非表示', () => {
    renderSection({ showTemplates: false });
    expect(screen.queryByText('テストテンプレート')).not.toBeInTheDocument();
  });

  it('テンプレート利用ボタンクリックでhandleUseTemplateが呼ばれる', () => {
    const handleUseTemplate = vi.fn();
    renderSection({ showTemplates: true, handleUseTemplate });
    fireEvent.click(screen.getByText('このテンプレートを使う'));
    expect(handleUseTemplate).toHaveBeenCalledWith(1);
  });

  it('creating=trueで利用ボタンが無効になる', () => {
    renderSection({ showTemplates: true, creating: true });
    const button = screen.getByText('読み込み中...');
    expect(button).toBeDisabled();
  });

  it('展開時にステップが表示される', () => {
    renderSection({ showTemplates: true, expandedTemplate: 1 });
    expect(screen.getByText('ステップ1')).toBeInTheDocument();
    expect(screen.getByText('ステップ2')).toBeInTheDocument();
  });
});
