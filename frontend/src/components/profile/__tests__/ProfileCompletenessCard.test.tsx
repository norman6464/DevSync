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

  it('75%以上でTrophyマイルストーンバッジが表示される', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={75} missingFields={[]} />);
    expect(screen.getByText('もうすぐ完成！')).toBeInTheDocument();
  });

  it('50-74%でTrendingUpマイルストーンバッジが表示される', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={50} missingFields={[]} />);
    expect(screen.getByText('半分完成！')).toBeInTheDocument();
  });

  it('25-49%でZapマイルストーンバッジが表示される', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={25} missingFields={[]} />);
    expect(screen.getByText('順調です！')).toBeInTheDocument();
  });

  it('25%未満ではマイルストーンバッジが表示されない', () => {
    renderWithRouter(<ProfileCompletenessCard percentage={20} missingFields={[]} />);
    expect(screen.queryByText('もうすぐ完成！')).not.toBeInTheDocument();
    expect(screen.queryByText('半分完成！')).not.toBeInTheDocument();
    expect(screen.queryByText('順調です！')).not.toBeInTheDocument();
  });
});
