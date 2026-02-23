import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import MasonryGrid from '../MasonryGrid';

describe('MasonryGrid', () => {
  it('子要素が表示される', () => {
    render(
      <MasonryGrid columns={3}>
        <div>アイテム1</div>
        <div>アイテム2</div>
        <div>アイテム3</div>
      </MasonryGrid>
    );
    expect(screen.getByText('アイテム1')).toBeInTheDocument();
    expect(screen.getByText('アイテム2')).toBeInTheDocument();
    expect(screen.getByText('アイテム3')).toBeInTheDocument();
  });

  it('指定カラム数のグリッドが生成される', () => {
    const { container } = render(
      <MasonryGrid columns={3}>
        <div>1</div><div>2</div><div>3</div>
      </MasonryGrid>
    );
    const cols = container.querySelectorAll('[data-testid="masonry-column"]');
    expect(cols.length).toBe(3);
  });

  it('2カラムのグリッドが生成される', () => {
    const { container } = render(
      <MasonryGrid columns={2}>
        <div>1</div><div>2</div>
      </MasonryGrid>
    );
    const cols = container.querySelectorAll('[data-testid="masonry-column"]');
    expect(cols.length).toBe(2);
  });

  it('アイテムがカラムに分配される', () => {
    render(
      <MasonryGrid columns={2}>
        <div>A</div><div>B</div><div>C</div><div>D</div>
      </MasonryGrid>
    );
    expect(screen.getByText('A')).toBeInTheDocument();
    expect(screen.getByText('B')).toBeInTheDocument();
    expect(screen.getByText('C')).toBeInTheDocument();
    expect(screen.getByText('D')).toBeInTheDocument();
  });

  it('ギャップが適用される', () => {
    const { container } = render(
      <MasonryGrid columns={2} gap={16}>
        <div>1</div><div>2</div>
      </MasonryGrid>
    );
    const grid = container.firstChild;
    expect(grid).toHaveStyle({ gap: '16px' });
  });

  it('デフォルトギャップは8px', () => {
    const { container } = render(
      <MasonryGrid columns={2}>
        <div>1</div><div>2</div>
      </MasonryGrid>
    );
    const grid = container.firstChild;
    expect(grid).toHaveStyle({ gap: '8px' });
  });

  it('空の場合何も表示されない', () => {
    const { container } = render(<MasonryGrid columns={3}>{[]}</MasonryGrid>);
    const cols = container.querySelectorAll('[data-testid="masonry-column"]');
    expect(cols.length).toBe(3);
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(
      <MasonryGrid columns={2} className="custom-class">
        <div>1</div>
      </MasonryGrid>
    );
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('デフォルトカラム数は3', () => {
    const { container } = render(
      <MasonryGrid>
        <div>1</div><div>2</div><div>3</div>
      </MasonryGrid>
    );
    const cols = container.querySelectorAll('[data-testid="masonry-column"]');
    expect(cols.length).toBe(3);
  });
});
