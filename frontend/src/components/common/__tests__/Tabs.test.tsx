import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import Tabs from '../Tabs';

describe('Tabs', () => {
  const tabs = [
    { id: '1', label: 'タブ1', content: 'コンテンツ1' },
    { id: '2', label: 'タブ2', content: 'コンテンツ2' },
    { id: '3', label: 'タブ3', content: 'コンテンツ3' },
  ];

  it('タブが表示される', () => {
    render(<Tabs tabs={tabs} />);

    expect(screen.getByText('タブ1')).toBeInTheDocument();
    expect(screen.getByText('タブ2')).toBeInTheDocument();
    expect(screen.getByText('タブ3')).toBeInTheDocument();
  });

  it('デフォルトで最初のタブのコンテンツが表示される', () => {
    render(<Tabs tabs={tabs} />);

    expect(screen.getByText('コンテンツ1')).toBeInTheDocument();
    expect(screen.queryByText('コンテンツ2')).not.toBeInTheDocument();
  });

  it('タブをクリックするとコンテンツが切り替わる', () => {
    render(<Tabs tabs={tabs} />);

    const tab2 = screen.getByText('タブ2');
    fireEvent.click(tab2);

    expect(screen.getByText('コンテンツ2')).toBeInTheDocument();
    expect(screen.queryByText('コンテンツ1')).not.toBeInTheDocument();
  });

  it('アクティブなタブがハイライトされる', () => {
    render(<Tabs tabs={tabs} />);

    const tab1 = screen.getByText('タブ1');
    expect(tab1).toHaveClass('text-blue-400');
  });

  it('非アクティブなタブは灰色で表示される', () => {
    render(<Tabs tabs={tabs} />);

    const tab2 = screen.getByText('タブ2');
    expect(tab2).toHaveClass('text-gray-400');
  });

  it('アクティブなタブにボーダーがある', () => {
    render(<Tabs tabs={tabs} />);

    const tab1 = screen.getByText('タブ1');
    expect(tab1).toHaveClass('border-b-2', 'border-blue-400');
  });

  it('タブクリック時にonChangeが呼ばれる', () => {
    const mockOnChange = vi.fn();
    render(<Tabs tabs={tabs} onChange={mockOnChange} />);

    const tab2 = screen.getByText('タブ2');
    fireEvent.click(tab2);

    expect(mockOnChange).toHaveBeenCalledWith('2');
  });

  it('defaultActiveIdで初期アクティブタブを指定できる', () => {
    render(<Tabs tabs={tabs} defaultActiveId="2" />);

    expect(screen.getByText('コンテンツ2')).toBeInTheDocument();
    expect(screen.queryByText('コンテンツ1')).not.toBeInTheDocument();
  });

  it('ホバー時にタブの色が変わる', () => {
    render(<Tabs tabs={tabs} />);

    const tab2 = screen.getByText('タブ2');
    expect(tab2).toHaveClass('hover:text-white');
  });

  it('トランジション効果がある', () => {
    render(<Tabs tabs={tabs} />);

    const tab1 = screen.getByText('タブ1');
    expect(tab1).toHaveClass('transition-colors');
  });

  it('複数回クリックしても正しく動作する', () => {
    render(<Tabs tabs={tabs} />);

    const tab2 = screen.getByText('タブ2');
    const tab3 = screen.getByText('タブ3');

    fireEvent.click(tab2);
    expect(screen.getByText('コンテンツ2')).toBeInTheDocument();

    fireEvent.click(tab3);
    expect(screen.getByText('コンテンツ3')).toBeInTheDocument();

    fireEvent.click(tab2);
    expect(screen.getByText('コンテンツ2')).toBeInTheDocument();
  });
});
