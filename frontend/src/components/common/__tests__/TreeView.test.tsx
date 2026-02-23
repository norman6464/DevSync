import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import TreeView from '../TreeView';

const nodes = [
  {
    id: '1',
    label: 'src',
    children: [
      { id: '1-1', label: 'components' },
      { id: '1-2', label: 'pages', children: [{ id: '1-2-1', label: 'Home.tsx' }] },
    ],
  },
  { id: '2', label: 'README.md' },
];

describe('TreeView', () => {
  it('ルートノードが表示される', () => {
    render(<TreeView nodes={nodes} />);
    expect(screen.getByText('src')).toBeInTheDocument();
    expect(screen.getByText('README.md')).toBeInTheDocument();
  });

  it('子ノードが初期状態で非表示', () => {
    render(<TreeView nodes={nodes} />);
    expect(screen.queryByText('components')).not.toBeInTheDocument();
  });

  it('クリックで子ノードが展開される', async () => {
    const user = userEvent.setup();
    render(<TreeView nodes={nodes} />);
    await user.click(screen.getByText('src'));
    expect(screen.getByText('components')).toBeInTheDocument();
    expect(screen.getByText('pages')).toBeInTheDocument();
  });

  it('再クリックで子ノードが折りたたまれる', async () => {
    const user = userEvent.setup();
    render(<TreeView nodes={nodes} />);
    await user.click(screen.getByText('src'));
    await user.click(screen.getByText('src'));
    expect(screen.queryByText('components')).not.toBeInTheDocument();
  });

  it('ネストされた子ノードが展開される', async () => {
    const user = userEvent.setup();
    render(<TreeView nodes={nodes} />);
    await user.click(screen.getByText('src'));
    await user.click(screen.getByText('pages'));
    expect(screen.getByText('Home.tsx')).toBeInTheDocument();
  });

  it('展開アイコンが子ありノードに表示される', () => {
    const { container } = render(<TreeView nodes={nodes} />);
    const chevrons = container.querySelectorAll('.lucide-chevron-right');
    expect(chevrons.length).toBeGreaterThan(0);
  });

  it('リーフノードには展開アイコンなし', () => {
    const { container } = render(<TreeView nodes={[{ id: '1', label: 'file.txt' }]} />);
    expect(container.querySelector('.lucide-chevron-right')).not.toBeInTheDocument();
  });

  it('onSelectコールバックが呼ばれる', async () => {
    const onSelect = vi.fn();
    const user = userEvent.setup();
    render(<TreeView nodes={nodes} onSelect={onSelect} />);
    await user.click(screen.getByText('README.md'));
    expect(onSelect).toHaveBeenCalledWith('2');
  });

  it('defaultExpandedIdsで初期展開される', () => {
    render(<TreeView nodes={nodes} defaultExpandedIds={['1']} />);
    expect(screen.getByText('components')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<TreeView nodes={nodes} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
