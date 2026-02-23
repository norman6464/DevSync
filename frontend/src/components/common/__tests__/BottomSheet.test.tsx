import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import BottomSheet from '../BottomSheet';

describe('BottomSheet', () => {
  it('開いた状態でコンテンツが表示される', () => {
    render(<BottomSheet open onClose={() => {}}>内容</BottomSheet>);
    expect(screen.getByText('内容')).toBeInTheDocument();
  });

  it('閉じた状態で非表示', () => {
    render(<BottomSheet open={false} onClose={() => {}}>内容</BottomSheet>);
    expect(screen.queryByText('内容')).not.toBeInTheDocument();
  });

  it('タイトルが表示される', () => {
    render(<BottomSheet open onClose={() => {}} title="メニュー">内容</BottomSheet>);
    expect(screen.getByText('メニュー')).toBeInTheDocument();
  });

  it('ドラッグハンドルが表示される', () => {
    render(<BottomSheet open onClose={() => {}}>内容</BottomSheet>);
    expect(screen.getByTestId('drag-handle')).toBeInTheDocument();
  });

  it('オーバーレイクリックでonCloseが呼ばれる', async () => {
    const onClose = vi.fn();
    const user = userEvent.setup();
    render(<BottomSheet open onClose={onClose}>内容</BottomSheet>);
    await user.click(screen.getByTestId('bottomsheet-overlay'));
    expect(onClose).toHaveBeenCalled();
  });

  it('カスタム高さが適用される', () => {
    render(<BottomSheet open onClose={() => {}} maxHeight="80vh">内容</BottomSheet>);
    const sheet = screen.getByText('内容').closest('[style]');
    expect(sheet).toHaveStyle({ maxHeight: '80vh' });
  });

  it('デフォルト高さは50vh', () => {
    render(<BottomSheet open onClose={() => {}}>内容</BottomSheet>);
    const sheet = screen.getByText('内容').closest('[style]');
    expect(sheet).toHaveStyle({ maxHeight: '50vh' });
  });

  it('子要素が正しく描画される', () => {
    render(
      <BottomSheet open onClose={() => {}}>
        <p>要素1</p>
        <p>要素2</p>
      </BottomSheet>
    );
    expect(screen.getByText('要素1')).toBeInTheDocument();
    expect(screen.getByText('要素2')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    render(<BottomSheet open onClose={() => {}} className="custom-class">内容</BottomSheet>);
    expect(screen.getByText('内容').closest('.custom-class')).toBeInTheDocument();
  });

  it('閉じるボタンが表示される', async () => {
    const onClose = vi.fn();
    const user = userEvent.setup();
    render(<BottomSheet open onClose={onClose} title="タイトル">内容</BottomSheet>);
    await user.click(screen.getByLabelText('閉じる'));
    expect(onClose).toHaveBeenCalled();
  });
});
