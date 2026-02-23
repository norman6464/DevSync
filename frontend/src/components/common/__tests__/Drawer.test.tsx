import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Drawer from '../Drawer';

describe('Drawer', () => {
  it('開いた状態でコンテンツが表示される', () => {
    render(<Drawer open onClose={() => {}}>内容</Drawer>);
    expect(screen.getByText('内容')).toBeInTheDocument();
  });

  it('閉じた状態で非表示', () => {
    render(<Drawer open={false} onClose={() => {}}>内容</Drawer>);
    expect(screen.queryByText('内容')).not.toBeInTheDocument();
  });

  it('タイトルが表示される', () => {
    render(<Drawer open onClose={() => {}} title="設定">内容</Drawer>);
    expect(screen.getByText('設定')).toBeInTheDocument();
  });

  it('閉じるボタンで onClose が呼ばれる', async () => {
    const onClose = vi.fn();
    const user = userEvent.setup();
    render(<Drawer open onClose={onClose} title="設定">内容</Drawer>);
    await user.click(screen.getByRole('button'));
    expect(onClose).toHaveBeenCalled();
  });

  it('オーバーレイクリックで閉じる', async () => {
    const onClose = vi.fn();
    const user = userEvent.setup();
    render(<Drawer open onClose={onClose}>内容</Drawer>);
    await user.click(screen.getByTestId('drawer-overlay'));
    expect(onClose).toHaveBeenCalled();
  });

  it('右からスライドイン（デフォルト）', () => {
    const { container } = render(<Drawer open onClose={() => {}}>内容</Drawer>);
    expect(container.querySelector('.right-0')).toBeInTheDocument();
  });

  it('左からスライドイン', () => {
    const { container } = render(<Drawer open onClose={() => {}} position="left">内容</Drawer>);
    expect(container.querySelector('.left-0')).toBeInTheDocument();
  });

  it('下からスライドイン', () => {
    const { container } = render(<Drawer open onClose={() => {}} position="bottom">内容</Drawer>);
    expect(container.querySelector('.bottom-0')).toBeInTheDocument();
  });

  it('カスタム幅が適用される', () => {
    render(<Drawer open onClose={() => {}} width="400px">内容</Drawer>);
    const panel = screen.getByText('内容').closest('[style]');
    expect(panel).toHaveStyle({ width: '400px' });
  });

  it('カスタムクラス名が適用される', () => {
    render(<Drawer open onClose={() => {}} className="custom-class">内容</Drawer>);
    expect(screen.getByText('内容').closest('.custom-class')).toBeInTheDocument();
  });

  it('子要素が正しく描画される', () => {
    render(
      <Drawer open onClose={() => {}}>
        <p>パラグラフ1</p>
        <p>パラグラフ2</p>
      </Drawer>
    );
    expect(screen.getByText('パラグラフ1')).toBeInTheDocument();
    expect(screen.getByText('パラグラフ2')).toBeInTheDocument();
  });
});
