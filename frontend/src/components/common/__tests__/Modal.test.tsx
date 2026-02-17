import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import Modal from '../Modal';

describe('Modal', () => {
  const onClose = vi.fn();

  afterEach(() => {
    onClose.mockClear();
  });

  it('isOpen=trueの場合にモーダルが表示される', () => {
    render(
      <Modal isOpen={true} onClose={onClose} title="テストモーダル">
        <p>コンテンツ</p>
      </Modal>
    );
    expect(screen.getByText('テストモーダル')).toBeInTheDocument();
    expect(screen.getByText('コンテンツ')).toBeInTheDocument();
  });

  it('isOpen=falseの場合にモーダルが表示されない', () => {
    render(
      <Modal isOpen={false} onClose={onClose} title="テストモーダル">
        <p>コンテンツ</p>
      </Modal>
    );
    expect(screen.queryByText('テストモーダル')).toBeNull();
  });

  it('背景クリックでonCloseが呼ばれる', () => {
    render(
      <Modal isOpen={true} onClose={onClose} title="テスト">
        <p>コンテンツ</p>
      </Modal>
    );
    // 背景（overlay）をクリック
    const overlay = screen.getByTestId('modal-overlay');
    fireEvent.click(overlay);
    expect(onClose).toHaveBeenCalledOnce();
  });

  it('モーダル内部のクリックでは閉じない', () => {
    render(
      <Modal isOpen={true} onClose={onClose} title="テスト">
        <p>コンテンツ</p>
      </Modal>
    );
    fireEvent.click(screen.getByText('コンテンツ'));
    expect(onClose).not.toHaveBeenCalled();
  });

  it('Escapeキーでoncloseが呼ばれる', () => {
    render(
      <Modal isOpen={true} onClose={onClose} title="テスト">
        <p>コンテンツ</p>
      </Modal>
    );
    fireEvent.keyDown(document, { key: 'Escape' });
    expect(onClose).toHaveBeenCalledOnce();
  });

  it('titleが未指定の場合はヘッダーが表示されない', () => {
    render(
      <Modal isOpen={true} onClose={onClose}>
        <p>コンテンツのみ</p>
      </Modal>
    );
    expect(screen.queryByRole('heading')).toBeNull();
    expect(screen.getByText('コンテンツのみ')).toBeInTheDocument();
  });

  it('maxWidthクラスが適用される', () => {
    render(
      <Modal isOpen={true} onClose={onClose} title="テスト" maxWidth="max-w-lg">
        <p>コンテンツ</p>
      </Modal>
    );
    const dialog = screen.getByRole('dialog');
    expect(dialog.className).toContain('max-w-lg');
  });
});
