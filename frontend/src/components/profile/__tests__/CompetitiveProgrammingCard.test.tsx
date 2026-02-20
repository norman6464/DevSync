import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import CompetitiveProgrammingCard from '../CompetitiveProgrammingCard';
import type { AtCoderRatingInfo } from '../../../api/atcoder';

const mockAtcoder: AtCoderRatingInfo = {
  username: 'testuser',
  rating: 1200,
  color: 'green',
  rank: '5 Kyu',
};

describe('CompetitiveProgrammingCard', () => {
  it('タイトルを表示する', () => {
    render(<CompetitiveProgrammingCard atcoderRating={mockAtcoder} atcoderUsername="testuser" />);
    expect(screen.getByText('競技プログラミング')).toBeInTheDocument();
  });

  it('AtCoderレーティングを表示する', () => {
    render(<CompetitiveProgrammingCard atcoderRating={mockAtcoder} atcoderUsername="testuser" />);
    expect(screen.getByText('1200')).toBeInTheDocument();
    expect(screen.getByText('(5 Kyu)')).toBeInTheDocument();
  });

  it('AtCoderリンクが正しいURLを持つ', () => {
    render(<CompetitiveProgrammingCard atcoderRating={mockAtcoder} atcoderUsername="testuser" />);
    const link = screen.getByText('1200').closest('a');
    expect(link).toHaveAttribute('href', 'https://atcoder.jp/users/testuser');
  });

  it('AtCoderリンクがnoopener noreferrerを持つ', () => {
    render(<CompetitiveProgrammingCard atcoderRating={mockAtcoder} atcoderUsername="testuser" />);
    const link = screen.getByText('1200').closest('a');
    expect(link).toHaveAttribute('rel', 'noopener noreferrer');
  });

  it('paizaランクを表示する', () => {
    render(<CompetitiveProgrammingCard atcoderRating={null} paizaRank="B" />);
    expect(screen.getByText('paiza ランクB')).toBeInTheDocument();
  });

  it('AtCoderとpaiza両方表示できる', () => {
    render(<CompetitiveProgrammingCard atcoderRating={mockAtcoder} atcoderUsername="testuser" paizaRank="A" />);
    expect(screen.getByText('1200')).toBeInTheDocument();
    expect(screen.getByText('paiza ランクA')).toBeInTheDocument();
  });

  it('両方nullの場合はnullを返す', () => {
    const { container } = render(<CompetitiveProgrammingCard atcoderRating={null} />);
    expect(container.firstChild).toBeNull();
  });

  it('AtCoderの色スタイルが適用される', () => {
    render(<CompetitiveProgrammingCard atcoderRating={mockAtcoder} atcoderUsername="testuser" />);
    const ratingEl = screen.getByText('1200');
    expect(ratingEl).toHaveStyle({ color: '#008000' });
  });
});
