import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import FileUpload from '../FileUpload';

describe('FileUpload', () => {
  it('アップロードエリアが表示される', () => {
    render(<FileUpload onUpload={() => {}} />);

    expect(screen.getByText('ファイルをドラッグ&ドロップ')).toBeInTheDocument();
  });

  it('アップロードアイコンが表示される', () => {
    const { container } = render(<FileUpload onUpload={() => {}} />);

    expect(container.querySelector('svg')).toBeInTheDocument();
  });

  it('クリックでファイル選択ができる', () => {
    const { container } = render(<FileUpload onUpload={() => {}} />);

    const input = container.querySelector('input[type="file"]');
    expect(input).toBeInTheDocument();
  });

  it('ファイル選択時にコールバックが呼ばれる', async () => {
    const onUpload = vi.fn();
    const { container } = render(<FileUpload onUpload={onUpload} />);

    const input = container.querySelector('input[type="file"]') as HTMLInputElement;
    const file = new File(['test'], 'test.txt', { type: 'text/plain' });

    fireEvent.change(input, { target: { files: [file] } });

    expect(onUpload).toHaveBeenCalledWith([file]);
  });

  it('複数ファイルが選択できる', () => {
    const { container } = render(<FileUpload onUpload={() => {}} multiple />);

    const input = container.querySelector('input[type="file"]');
    expect(input).toHaveAttribute('multiple');
  });

  it('ファイルタイプが制限される', () => {
    const { container } = render(<FileUpload onUpload={() => {}} accept=".jpg,.png" />);

    const input = container.querySelector('input[type="file"]');
    expect(input).toHaveAttribute('accept', '.jpg,.png');
  });

  it('許可されたファイル形式が表示される', () => {
    render(<FileUpload onUpload={() => {}} accept=".jpg,.png" />);

    expect(screen.getByText('.jpg,.png')).toBeInTheDocument();
  });

  it('最大ファイルサイズが表示される', () => {
    render(<FileUpload onUpload={() => {}} maxSizeMB={5} />);

    expect(screen.getByText('最大 5MB')).toBeInTheDocument();
  });

  it('ドラッグオーバーでハイライトされる', () => {
    const { container } = render(<FileUpload onUpload={() => {}} />);

    const dropZone = container.querySelector('[data-testid="drop-zone"]')!;
    fireEvent.dragOver(dropZone);

    expect(dropZone).toHaveClass('border-blue-500');
  });

  it('ドラッグリーブでハイライトが解除される', () => {
    const { container } = render(<FileUpload onUpload={() => {}} />);

    const dropZone = container.querySelector('[data-testid="drop-zone"]')!;
    fireEvent.dragOver(dropZone);
    fireEvent.dragLeave(dropZone);

    expect(dropZone).not.toHaveClass('border-blue-500');
  });

  it('無効状態が適用される', () => {
    const { container } = render(<FileUpload onUpload={() => {}} disabled />);

    expect(container.querySelector('.opacity-50')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<FileUpload onUpload={() => {}} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
