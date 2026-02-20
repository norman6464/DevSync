import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ProfileCompletenessCard from '../ProfileCompletenessCard';

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>);

describe('ProfileCompletenessCard', () => {
  it('完成度パーセンテージを表示する', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={60} missingFields={[]} />);
    expect(screen.getByText('60%')).toBeInTheDocument();
  });

  it('タイトルを表示する', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={50} missingFields={[]} />);
    expect(screen.getByText('プロフィール完成度')).toBeInTheDocument();
  });

  it('プログレスバーが表示される', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={75} missingFields={[]} />);
    const progressbar = screen.getByRole('progressbar');
    expect(progressbar).toHaveAttribute('aria-valuenow', '75');
  });

  it('不足フィールドのリンクを表示する', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={50} missingFields={['avatar', 'bio']} />);
    expect(screen.getByText(/アバター画像を設定/)).toBeInTheDocument();
    expect(screen.getByText(/自己紹介を追加/)).toBeInTheDocument();
  });

  it('不足フィールドのリンクが/settingsに向く', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={50} missingFields={['avatar']} />);
    const link = screen.getByText(/アバター画像を設定/).closest('a');
    expect(link).toHaveAttribute('href', '/settings');
  });

  it('100%の場合はnullを返す', () => {
    const { container } = renderWithRouter(<ProfileCompletenessCard percentage={100} missingFields={[]} />);
    expect(container.firstChild).toBeNull();
  });
});
