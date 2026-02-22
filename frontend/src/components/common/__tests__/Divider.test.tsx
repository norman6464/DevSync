import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import Divider from '../Divider';

describe('Divider', () => {
  it('区切り線が表示される', () => {
    const { container } = render(<Divider />);

    const divider = container.querySelector('.border-gray-800');
    expect(divider).toBeInTheDocument();
  });

  it('水平方向の区切り線が表示される', () => {
    const { container } = render(<Divider />);

    const divider = container.querySelector('.border-t');
    expect(divider).toBeInTheDocument();
  });

  it('垂直方向の区切り線が表示される', () => {
    const { container } = render(<Divider orientation="vertical" />);

    const divider = container.querySelector('.border-l');
    expect(divider).toBeInTheDocument();
  });

  it('テキスト付き区切り線が表示される', () => {
    render(<Divider>テキスト</Divider>);

    expect(screen.getByText('テキスト')).toBeInTheDocument();
  });

  it('テキストが中央に表示される', () => {
    render(<Divider textAlign="center">中央</Divider>);

    const text = screen.getByText('中央');
    expect(text.parentElement).toHaveClass('justify-center');
  });

  it('テキストが左に表示される', () => {
    render(<Divider textAlign="left">左</Divider>);

    const text = screen.getByText('左');
    expect(text.parentElement).toHaveClass('justify-start');
  });

  it('テキストが右に表示される', () => {
    render(<Divider textAlign="right">右</Divider>);

    const text = screen.getByText('右');
    expect(text.parentElement).toHaveClass('justify-end');
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<Divider className="custom-class" />);

    const divider = container.querySelector('.custom-class');
    expect(divider).toBeInTheDocument();
  });

  it('デフォルトは水平方向', () => {
    const { container } = render(<Divider />);

    const divider = container.querySelector('.border-t');
    expect(divider).toBeInTheDocument();
  });

  it('テキストがない場合はシンプルな線のみ', () => {
    const { container } = render(<Divider />);

    const divider = container.querySelector('.border-t');
    expect(divider).toBeInTheDocument();
    expect(container.querySelector('.flex')).not.toBeInTheDocument();
  });
});
