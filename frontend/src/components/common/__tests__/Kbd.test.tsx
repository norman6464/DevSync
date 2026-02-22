import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Kbd from '../Kbd';

describe('Kbd', () => {
  it('単一キーが表示される', () => {
    render(<Kbd keys={['Enter']} />);

    expect(screen.getByText('Enter')).toBeInTheDocument();
  });

  it('複数キーが+で区切られる', () => {
    render(<Kbd keys={['Ctrl', 'S']} />);

    expect(screen.getByText('Ctrl')).toBeInTheDocument();
    expect(screen.getByText('S')).toBeInTheDocument();
    expect(screen.getByText('+')).toBeInTheDocument();
  });

  it('kbd要素が使用される', () => {
    const { container } = render(<Kbd keys={['A']} />);

    expect(container.querySelector('kbd')).toBeInTheDocument();
  });

  it('smサイズが適用される', () => {
    const { container } = render(<Kbd keys={['A']} size="sm" />);

    expect(container.querySelector('.text-xs')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<Kbd keys={['A']} size="lg" />);

    expect(container.querySelector('.text-base')).toBeInTheDocument();
  });

  it('キーキャップ風のスタイルが適用される', () => {
    const { container } = render(<Kbd keys={['A']} />);

    const kbd = container.querySelector('kbd');
    expect(kbd).toHaveClass('rounded');
  });

  it('3つのキーの組み合わせが表示される', () => {
    render(<Kbd keys={['Ctrl', 'Shift', 'P']} />);

    expect(screen.getByText('Ctrl')).toBeInTheDocument();
    expect(screen.getByText('Shift')).toBeInTheDocument();
    expect(screen.getByText('P')).toBeInTheDocument();
    expect(screen.getAllByText('+')).toHaveLength(2);
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Kbd keys={['A']} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('特殊キーが正しく表示される', () => {
    render(<Kbd keys={['⌘', '⇧', 'K']} />);

    expect(screen.getByText('⌘')).toBeInTheDocument();
    expect(screen.getByText('⇧')).toBeInTheDocument();
    expect(screen.getByText('K')).toBeInTheDocument();
  });
});
