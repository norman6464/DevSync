import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { I18nextProvider } from 'react-i18next';
import i18n from '../../../i18n';
import PageLoader from '../PageLoader';

describe('PageLoader', () => {
  const renderWithI18n = (component: React.ReactElement) => {
    return render(
      <I18nextProvider i18n={i18n}>
        {component}
      </I18nextProvider>
    );
  };

  it('デフォルトでレンダリングされる', () => {
    renderWithI18n(<PageLoader />);
    expect(screen.getByText('読み込み中...')).toBeInTheDocument();
  });

  it('カスタムメッセージが表示される', () => {
    renderWithI18n(<PageLoader message="データを取得中..." />);
    expect(screen.getByText('データを取得中...')).toBeInTheDocument();
  });

  it('メッセージなしでレンダリングされる', () => {
    renderWithI18n(<PageLoader showMessage={false} />);
    expect(screen.queryByText('読み込み中...')).not.toBeInTheDocument();
  });

  it('フルハイト表示される', () => {
    const { container } = renderWithI18n(<PageLoader fullHeight />);
    const wrapper = container.firstChild;
    expect(wrapper).toHaveClass('min-h-screen');
  });

  it('デフォルトではフルハイトではない', () => {
    const { container } = renderWithI18n(<PageLoader />);
    const wrapper = container.firstChild;
    expect(wrapper).not.toHaveClass('min-h-screen');
    expect(wrapper).toHaveClass('py-12');
  });
});
