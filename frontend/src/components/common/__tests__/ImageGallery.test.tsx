import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ImageGallery from '../ImageGallery';

const images = [
  { src: '/img1.jpg', alt: '画像1' },
  { src: '/img2.jpg', alt: '画像2' },
  { src: '/img3.jpg', alt: '画像3' },
  { src: '/img4.jpg', alt: '画像4' },
];

describe('ImageGallery', () => {
  it('すべての画像が表示される', () => {
    render(<ImageGallery images={images} />);
    const imgs = screen.getAllByRole('img');
    expect(imgs.length).toBe(4);
  });

  it('alt属性が正しく設定される', () => {
    render(<ImageGallery images={images} />);
    expect(screen.getByAlt('画像1')).toBeInTheDocument();
    expect(screen.getByAlt('画像2')).toBeInTheDocument();
  });

  it('画像クリックでライトボックスが表示される', async () => {
    const user = userEvent.setup();
    render(<ImageGallery images={images} />);
    await user.click(screen.getByAlt('画像1'));
    expect(screen.getByTestId('lightbox')).toBeInTheDocument();
  });

  it('ライトボックスで閉じるボタンが動作する', async () => {
    const user = userEvent.setup();
    render(<ImageGallery images={images} />);
    await user.click(screen.getByAlt('画像1'));
    await user.click(screen.getByLabelText('閉じる'));
    expect(screen.queryByTestId('lightbox')).not.toBeInTheDocument();
  });

  it('ライトボックスで次へナビゲーション', async () => {
    const user = userEvent.setup();
    render(<ImageGallery images={images} />);
    await user.click(screen.getByAlt('画像1'));
    await user.click(screen.getByLabelText('次へ'));
    const lightboxImg = screen.getByTestId('lightbox').querySelector('img');
    expect(lightboxImg).toHaveAttribute('alt', '画像2');
  });

  it('ライトボックスで前へナビゲーション', async () => {
    const user = userEvent.setup();
    render(<ImageGallery images={images} />);
    await user.click(screen.getByAlt('画像2'));
    await user.click(screen.getByLabelText('前へ'));
    const lightboxImg = screen.getByTestId('lightbox').querySelector('img');
    expect(lightboxImg).toHaveAttribute('alt', '画像1');
  });

  it('カラム数が適用される', () => {
    const { container } = render(<ImageGallery images={images} columns={3} />);
    const grid = container.querySelector('[style]');
    expect(grid).toHaveStyle({ gridTemplateColumns: 'repeat(3, minmax(0, 1fr))' });
  });

  it('デフォルトカラム数は3', () => {
    const { container } = render(<ImageGallery images={images} />);
    const grid = container.querySelector('[style]');
    expect(grid).toHaveStyle({ gridTemplateColumns: 'repeat(3, minmax(0, 1fr))' });
  });

  it('空画像でメッセージが表示される', () => {
    render(<ImageGallery images={[]} emptyMessage="画像がありません" />);
    expect(screen.getByText('画像がありません')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<ImageGallery images={images} className="custom-class" />);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
