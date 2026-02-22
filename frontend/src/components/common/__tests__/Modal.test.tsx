import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import Modal from '../Modal';

describe('Modal', () => {
  const mockOnClose = vi.fn();

  it('モーダルが表示される', () => {
    render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    expect(screen.getByText('コンテンツ')).toBeInTheDocument();
  });

  it('閉じている時は表示されない', () => {
    render(
      <Modal isOpen={false} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    expect(screen.queryByText('コンテンツ')).not.toBeInTheDocument();
  });

  it('閉じるボタンが表示される', () => {
    const { container } = render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    const closeButton = container.querySelector('[aria-label="閉じる"]');
    expect(closeButton).toBeInTheDocument();
  });

  it('閉じるボタンをクリックするとonCloseが呼ばれる', () => {
    const { container } = render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    const closeButton = container.querySelector('[aria-label="閉じる"]');
    fireEvent.click(closeButton!);

    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('ESCキーでonCloseが呼ばれる', () => {
    render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    fireEvent.keyDown(document, { key: 'Escape' });

    expect(mockOnClose).toHaveBeenCalledTimes(1);
  });

  it('オーバーレイが表示される', () => {
    const { container } = render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    const overlay = container.querySelector('.bg-black\\/50');
    expect(overlay).toBeInTheDocument();
  });

  it('ヘッダーが表示される', () => {
    render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Header>タイトル</Modal.Header>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    expect(screen.getByText('タイトル')).toBeInTheDocument();
  });

  it('ボディが表示される', () => {
    render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>ボディコンテンツ</Modal.Body>
      </Modal>
    );

    expect(screen.getByText('ボディコンテンツ')).toBeInTheDocument();
  });

  it('フッターが表示される', () => {
    render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
        <Modal.Footer>フッター</Modal.Footer>
      </Modal>
    );

    expect(screen.getByText('フッター')).toBeInTheDocument();
  });

  it('モーダルに背景色がある', () => {
    const { container } = render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    const modal = container.querySelector('.bg-gray-900');
    expect(modal).toBeInTheDocument();
  });

  it('モーダルにボーダーがある', () => {
    const { container } = render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    const modal = container.querySelector('.border');
    expect(modal).toBeInTheDocument();
  });

  it('モーダルに角丸がある', () => {
    const { container } = render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    const modal = container.querySelector('.rounded-lg');
    expect(modal).toBeInTheDocument();
  });

  it('ヘッダーにタイトルが表示される', () => {
    render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Header>
          <Modal.Title>モーダルタイトル</Modal.Title>
        </Modal.Header>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    expect(screen.getByText('モーダルタイトル')).toBeInTheDocument();
  });

  it('タイトルが太字で表示される', () => {
    render(
      <Modal isOpen={true} onClose={mockOnClose}>
        <Modal.Header>
          <Modal.Title>タイトル</Modal.Title>
        </Modal.Header>
        <Modal.Body>コンテンツ</Modal.Body>
      </Modal>
    );

    const title = screen.getByText('タイトル');
    expect(title).toHaveClass('font-semibold');
  });

  it('複数のモーダルが独立して動作する', () => {
    const mockOnClose2 = vi.fn();

    render(
      <>
        <Modal isOpen={true} onClose={mockOnClose}>
          <Modal.Body>モーダル1</Modal.Body>
        </Modal>
        <Modal isOpen={false} onClose={mockOnClose2}>
          <Modal.Body>モーダル2</Modal.Body>
        </Modal>
      </>
    );

    expect(screen.getByText('モーダル1')).toBeInTheDocument();
    expect(screen.queryByText('モーダル2')).not.toBeInTheDocument();
  });
});
