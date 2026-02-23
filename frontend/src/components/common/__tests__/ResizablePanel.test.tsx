import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import ResizablePanel from '../ResizablePanel';

describe('ResizablePanel', () => {
  it('コンテンツが表示される', () => {
    render(<ResizablePanel>パネル内容</ResizablePanel>);
    expect(screen.getByText('パネル内容')).toBeInTheDocument();
  });

  it('初期幅が適用される', () => {
    render(<ResizablePanel defaultSize={300}>内容</ResizablePanel>);
    const panel = screen.getByText('内容').closest('[style]');
    expect(panel).toHaveStyle({ width: '300px' });
  });

  it('リサイズハンドルが表示される', () => {
    render(<ResizablePanel>内容</ResizablePanel>);
    expect(screen.getByTestId('resize-handle')).toBeInTheDocument();
  });

  it('最小幅が設定される', () => {
    render(<ResizablePanel minSize={100} defaultSize={100}>内容</ResizablePanel>);
    const panel = screen.getByText('内容').closest('[style]');
    expect(panel).toHaveStyle({ width: '100px' });
  });

  it('水平方向がデフォルト', () => {
    render(<ResizablePanel>内容</ResizablePanel>);
    const handle = screen.getByTestId('resize-handle');
    expect(handle.className).toContain('cursor-col-resize');
  });

  it('垂直方向に設定可能', () => {
    render(<ResizablePanel direction="vertical" defaultSize={200}>内容</ResizablePanel>);
    const handle = screen.getByTestId('resize-handle');
    expect(handle.className).toContain('cursor-row-resize');
  });

  it('垂直方向で高さが適用される', () => {
    render(<ResizablePanel direction="vertical" defaultSize={200}>内容</ResizablePanel>);
    const panel = screen.getByText('内容').closest('[style]');
    expect(panel).toHaveStyle({ height: '200px' });
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(<ResizablePanel className="custom-class">内容</ResizablePanel>);
    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('子要素が正しく描画される', () => {
    render(
      <ResizablePanel>
        <div>子要素1</div>
        <div>子要素2</div>
      </ResizablePanel>
    );
    expect(screen.getByText('子要素1')).toBeInTheDocument();
    expect(screen.getByText('子要素2')).toBeInTheDocument();
  });

  it('デフォルトサイズが250px', () => {
    render(<ResizablePanel>内容</ResizablePanel>);
    const panel = screen.getByText('内容').closest('[style]');
    expect(panel).toHaveStyle({ width: '250px' });
  });
});
