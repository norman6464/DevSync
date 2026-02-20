import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import StarRating from '../StarRating';

describe('StarRating', () => {
  describe('読み取り専用モード', () => {
    it('5つの星SVGを表示する', () => {
      const { container } = render(<StarRating rating={3} />);
      expect(container.querySelectorAll('svg')).toHaveLength(5);
    });

    it('ratingに応じて星にyellow-400クラスを適用する', () => {
      const { container } = render(<StarRating rating={3} />);
      const svgs = container.querySelectorAll('svg');
      expect(svgs[0].className.baseVal).toContain('text-yellow-400');
      expect(svgs[2].className.baseVal).toContain('text-yellow-400');
      expect(svgs[3].className.baseVal).toContain('text-gray-600');
      expect(svgs[4].className.baseVal).toContain('text-gray-600');
    });

    it('ボタンが表示されない（クリック不可）', () => {
      render(<StarRating rating={3} />);
      expect(screen.queryAllByRole('button')).toHaveLength(0);
    });

    it('デフォルトサイズはsmクラス', () => {
      const { container } = render(<StarRating rating={1} />);
      expect(container.querySelector('svg')?.className.baseVal).toContain('w-4 h-4');
    });

    it('size="lg"でlgクラスを適用する', () => {
      const { container } = render(<StarRating rating={1} size="lg" />);
      expect(container.querySelector('svg')?.className.baseVal).toContain('w-8 h-8');
    });
  });

  describe('インタラクティブモード', () => {
    it('5つのボタンを表示する', () => {
      const onChange = vi.fn();
      render(<StarRating rating={2} onChange={onChange} />);
      expect(screen.getAllByRole('button')).toHaveLength(5);
    });

    it('星クリックでonChangeが呼ばれる', () => {
      const onChange = vi.fn();
      render(<StarRating rating={2} onChange={onChange} />);
      fireEvent.click(screen.getAllByRole('button')[3]);
      expect(onChange).toHaveBeenCalledWith(4);
    });

    it('1番目の星クリックでonChange(1)が呼ばれる', () => {
      const onChange = vi.fn();
      render(<StarRating rating={0} onChange={onChange} />);
      fireEvent.click(screen.getAllByRole('button')[0]);
      expect(onChange).toHaveBeenCalledWith(1);
    });

    it('ボタンのtype属性がbuttonである', () => {
      const onChange = vi.fn();
      render(<StarRating rating={2} onChange={onChange} />);
      screen.getAllByRole('button').forEach((btn) => {
        expect(btn).toHaveAttribute('type', 'button');
      });
    });

    it('size="md"でmdクラスを適用する', () => {
      const onChange = vi.fn();
      const { container } = render(<StarRating rating={1} onChange={onChange} size="md" />);
      expect(container.querySelector('svg')?.className.baseVal).toContain('w-6 h-6');
    });
  });
});
