import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import LevelBadge from '../LevelBadge';

describe('LevelBadge', () => {
  it('レベル値を表示する', () => {
    render(<LevelBadge level={5} />);
    expect(screen.getByText(/5/)).toBeInTheDocument();
  });

  it('レベル1-10でgray系の色を適用する', () => {
    const { container } = render(<LevelBadge level={5} />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-gray-500/20');
  });

  it('レベル11-20でgreen系の色を適用する', () => {
    const { container } = render(<LevelBadge level={15} />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-green-500/20');
  });

  it('レベル21-30でblue系の色を適用する', () => {
    const { container } = render(<LevelBadge level={25} />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-blue-500/20');
  });

  it('レベル31-40でpurple系の色を適用する', () => {
    const { container } = render(<LevelBadge level={35} />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-purple-500/20');
  });

  it('レベル41以上でyellow系の色を適用する', () => {
    const { container } = render(<LevelBadge level={50} />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-yellow-500/20');
  });

  it('デフォルトでsmサイズを適用する', () => {
    const { container } = render(<LevelBadge level={1} />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('text-xs');
  });

  it('size=mdで大きめのサイズを適用する', () => {
    const { container } = render(<LevelBadge level={1} size="md" />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('text-sm');
  });

  it('境界値レベル11でgreen系を適用する', () => {
    const { container } = render(<LevelBadge level={11} />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-green-500/20');
  });

  it('境界値レベル41でyellow系を適用する', () => {
    const { container } = render(<LevelBadge level={41} />);
    const badge = container.querySelector('span');
    expect(badge?.className).toContain('bg-yellow-500/20');
  });
});
