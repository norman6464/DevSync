import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ColorPicker from '../ColorPicker';

describe('ColorPicker', () => {
  it('プリセットカラーが表示される', () => {
    const { container } = render(<ColorPicker value="#ef4444" onChange={() => {}} />);

    const swatches = container.querySelectorAll('[data-testid="color-swatch"]');
    expect(swatches.length).toBeGreaterThan(0);
  });

  it('選択中の色がハイライトされる', () => {
    const { container } = render(<ColorPicker value="#ef4444" onChange={() => {}} />);

    const selected = container.querySelector('.ring-2');
    expect(selected).toBeInTheDocument();
  });

  it('色をクリックで選択できる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    const { container } = render(<ColorPicker value="#ef4444" onChange={onChange} />);

    const swatches = container.querySelectorAll('[data-testid="color-swatch"]');
    await user.click(swatches[1]);

    expect(onChange).toHaveBeenCalled();
  });

  it('カスタムカラー入力が表示される', () => {
    render(<ColorPicker value="#ef4444" onChange={() => {}} showInput />);

    expect(screen.getByDisplayValue('#ef4444')).toBeInTheDocument();
  });

  it('カスタムカラーを入力できる', async () => {
    const onChange = vi.fn();
    const user = userEvent.setup();
    render(<ColorPicker value="#ef4444" onChange={onChange} showInput />);

    const input = screen.getByDisplayValue('#ef4444');
    await user.clear(input);
    await user.type(input, '#000000');

    expect(onChange).toHaveBeenCalled();
  });

  it('ラベルが表示される', () => {
    render(<ColorPicker value="#ef4444" onChange={() => {}} label="テーマカラー" />);

    expect(screen.getByText('テーマカラー')).toBeInTheDocument();
  });

  it('カスタムカラーパレットが使用される', () => {
    const colors = ['#ff0000', '#00ff00', '#0000ff'];
    const { container } = render(<ColorPicker value="#ff0000" onChange={() => {}} colors={colors} />);

    const swatches = container.querySelectorAll('[data-testid="color-swatch"]');
    expect(swatches.length).toBe(3);
  });

  it('smサイズが適用される', () => {
    const { container } = render(<ColorPicker value="#ef4444" onChange={() => {}} size="sm" />);

    expect(container.querySelector('.w-6')).toBeInTheDocument();
  });

  it('lgサイズが適用される', () => {
    const { container } = render(<ColorPicker value="#ef4444" onChange={() => {}} size="lg" />);

    expect(container.querySelector('.w-10')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<ColorPicker value="#ef4444" onChange={() => {}} className="custom-class" />);

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });
});
