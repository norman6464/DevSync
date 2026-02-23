import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import InfiniteScroll from '../InfiniteScroll';

const mockIntersectionObserver = vi.fn();

beforeEach(() => {
  mockIntersectionObserver.mockReturnValue({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
  });
  window.IntersectionObserver = mockIntersectionObserver;
});

describe('InfiniteScroll', () => {
  it('子要素が表示される', () => {
    render(
      <InfiniteScroll onLoadMore={() => {}} hasMore>
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(screen.getByText('コンテンツ')).toBeInTheDocument();
  });

  it('ローディング中にスピナーが表示される', () => {
    const { container } = render(
      <InfiniteScroll onLoadMore={() => {}} hasMore loading>
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(container.querySelector('.animate-spin')).toBeInTheDocument();
  });

  it('ローディング中でなければスピナー非表示', () => {
    const { container } = render(
      <InfiniteScroll onLoadMore={() => {}} hasMore>
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(container.querySelector('.animate-spin')).not.toBeInTheDocument();
  });

  it('データなし時にメッセージ表示', () => {
    render(
      <InfiniteScroll onLoadMore={() => {}} hasMore={false} endMessage="すべて読み込みました">
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(screen.getByText('すべて読み込みました')).toBeInTheDocument();
  });

  it('hasMoreがtrueの場合はendMessage非表示', () => {
    render(
      <InfiniteScroll onLoadMore={() => {}} hasMore endMessage="すべて読み込みました">
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(screen.queryByText('すべて読み込みました')).not.toBeInTheDocument();
  });

  it('IntersectionObserverが作成される', () => {
    render(
      <InfiniteScroll onLoadMore={() => {}} hasMore>
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(mockIntersectionObserver).toHaveBeenCalled();
  });

  it('センチネル要素が存在する', () => {
    const { container } = render(
      <InfiniteScroll onLoadMore={() => {}} hasMore>
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(container.querySelector('[data-testid="sentinel"]')).toBeInTheDocument();
  });

  it('hasMoreがfalseの場合センチネル非表示', () => {
    const { container } = render(
      <InfiniteScroll onLoadMore={() => {}} hasMore={false}>
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(container.querySelector('[data-testid="sentinel"]')).not.toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(
      <InfiniteScroll onLoadMore={() => {}} hasMore className="custom-class">
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('カスタムローディングテキストが表示される', () => {
    render(
      <InfiniteScroll onLoadMore={() => {}} hasMore loading loadingText="読み込み中...">
        <div>コンテンツ</div>
      </InfiniteScroll>
    );

    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });
});
