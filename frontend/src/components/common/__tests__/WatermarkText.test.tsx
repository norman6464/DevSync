import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import WatermarkText from '../WatermarkText';

describe('WatermarkText', () => {
  it('子要素が表示される', () => {
    render(<WatermarkText text="DRAFT">コンテンツ</WatermarkText>);
    expect(screen.getByText('コンテンツ')).toBeInTheDocument();
  });

  it('ウォーターマークテキストが表示される', () => {
    render(<WatermarkText text="DRAFT">コンテンツ</WatermarkText>);
    expect(screen.getByText('DRAFT')).toBeInTheDocument();
  });

  it('ウォーターマークが半透明', () => {
    render(<WatermarkText text="DRAFT">コンテンツ</WatermarkText>);
    const watermark = screen.getByText('DRAFT');
    expect(watermark.className).toContain('opacity-');
  });

  it('カスタム透明度が適用される', () => {
    render(<WatermarkText text="DRAFT" opacity={20}>コンテンツ</WatermarkText>);
    const watermark = screen.getByText('DRAFT');
    expect(watermark.style.opacity).toBe('0.2');
  });

  it('カスタム角度が適用される', () => {
    render(<WatermarkText text="DRAFT" rotate={-30}>コンテンツ</WatermarkText>);
    const watermark = screen.getByText('DRAFT');
    expect(watermark.style.transform).toContain('rotate(-30deg)');
  });

  it('デフォルト角度は-45度', () => {
    render(<WatermarkText text="DRAFT">コンテンツ</WatermarkText>);
    const watermark = screen.getByText('DRAFT');
    expect(watermark.style.transform).toContain('rotate(-45deg)');
  });

  it('フォントサイズが設定される', () => {
    render(<WatermarkText text="DRAFT" fontSize="4rem">コンテンツ</WatermarkText>);
    const watermark = screen.getByText('DRAFT');
    expect(watermark.style.fontSize).toBe('4rem');
  });

  it('relative配置で子要素を包む', () => {
    const { container } = render(<WatermarkText text="DRAFT">コンテンツ</WatermarkText>);
    expect(container.querySelector('.relative')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<WatermarkText text="DRAFT" className="custom-class">コンテンツ</WatermarkText>);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('ウォーターマークがpointer-events-none', () => {
    render(<WatermarkText text="DRAFT">コンテンツ</WatermarkText>);
    const watermark = screen.getByText('DRAFT');
    expect(watermark.className).toContain('pointer-events-none');
  });
});
