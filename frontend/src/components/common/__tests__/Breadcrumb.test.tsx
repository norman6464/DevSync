import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import Breadcrumb from '../Breadcrumb';

const renderWithRouter = (ui: React.ReactElement) => {
  return render(<BrowserRouter>{ui}</BrowserRouter>);
};

describe('Breadcrumb', () => {
  const items = [
    { label: 'ホーム', href: '/' },
    { label: 'プロジェクト', href: '/projects' },
    { label: 'プロジェクト詳細' },
  ];

  it('ブレッドクラムが表示される', () => {
    renderWithRouter(<Breadcrumb items={items} />);

    expect(screen.getByText('ホーム')).toBeInTheDocument();
    expect(screen.getByText('プロジェクト')).toBeInTheDocument();
    expect(screen.getByText('プロジェクト詳細')).toBeInTheDocument();
  });

  it('リンク付きの項目が表示される', () => {
    renderWithRouter(<Breadcrumb items={items} />);

    const homeLink = screen.getByText('ホーム');
    expect(homeLink.closest('a')).toHaveAttribute('href', '/');
  });

  it('最後の項目はリンクではない', () => {
    renderWithRouter(<Breadcrumb items={items} />);

    const lastItem = screen.getByText('プロジェクト詳細');
    expect(lastItem.closest('a')).toBeNull();
  });

  it('セパレーターが表示される', () => {
    const { container } = renderWithRouter(<Breadcrumb items={items} />);

    const separators = container.querySelectorAll('svg');
    // 3項目の場合、セパレーターは2つ
    expect(separators.length).toBe(2);
  });

  it('リンク項目にホバーエフェクトがある', () => {
    renderWithRouter(<Breadcrumb items={items} />);

    const homeLink = screen.getByText('ホーム');
    expect(homeLink).toHaveClass('hover:text-white');
  });

  it('最後の項目は灰色で表示される', () => {
    renderWithRouter(<Breadcrumb items={items} />);

    const lastItem = screen.getByText('プロジェクト詳細');
    expect(lastItem).toHaveClass('text-gray-400');
  });

  it('リンク項目は青色で表示される', () => {
    renderWithRouter(<Breadcrumb items={items} />);

    const homeLink = screen.getByText('ホーム');
    expect(homeLink).toHaveClass('text-blue-400');
  });

  it('1つだけの項目でも表示される', () => {
    renderWithRouter(<Breadcrumb items={[{ label: 'ホーム' }]} />);

    expect(screen.getByText('ホーム')).toBeInTheDocument();
  });

  it('navタグで囲まれている', () => {
    const { container } = renderWithRouter(<Breadcrumb items={items} />);

    const nav = container.querySelector('nav');
    expect(nav).toBeInTheDocument();
  });

  it('aria-labelが設定されている', () => {
    const { container } = renderWithRouter(<Breadcrumb items={items} />);

    const nav = container.querySelector('nav');
    expect(nav).toHaveAttribute('aria-label', 'パンくずリスト');
  });
});
