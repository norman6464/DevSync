import { describe, it, expect } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import Tooltip from '../Tooltip';

describe('Tooltip', () => {
  it('ツールチップが表示される', async () => {
    render(
      <Tooltip content="ヘルプテキスト">
        <button>ホバーしてください</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバーしてください');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      expect(screen.getByText('ヘルプテキスト')).toBeInTheDocument();
    });
  });

  it('ホバーを離れるとツールチップが非表示になる', async () => {
    render(
      <Tooltip content="ヘルプテキスト">
        <button>ホバーしてください</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバーしてください');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      expect(screen.getByText('ヘルプテキスト')).toBeInTheDocument();
    });

    fireEvent.mouseLeave(trigger);

    await waitFor(() => {
      expect(screen.queryByText('ヘルプテキスト')).not.toBeInTheDocument();
    });
  });

  it('子要素が表示される', () => {
    render(
      <Tooltip content="ヘルプ">
        <button>ボタン</button>
      </Tooltip>
    );

    expect(screen.getByText('ボタン')).toBeInTheDocument();
  });

  it('上方向のツールチップが表示される', async () => {
    const { container } = render(
      <Tooltip content="上" position="top">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = container.querySelector('.bottom-full');
      expect(tooltip).toBeInTheDocument();
    });
  });

  it('下方向のツールチップが表示される', async () => {
    const { container } = render(
      <Tooltip content="下" position="bottom">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = container.querySelector('.top-full');
      expect(tooltip).toBeInTheDocument();
    });
  });

  it('左方向のツールチップが表示される', async () => {
    const { container } = render(
      <Tooltip content="左" position="left">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = container.querySelector('.right-full');
      expect(tooltip).toBeInTheDocument();
    });
  });

  it('右方向のツールチップが表示される', async () => {
    const { container } = render(
      <Tooltip content="右" position="right">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = container.querySelector('.left-full');
      expect(tooltip).toBeInTheDocument();
    });
  });

  it('ツールチップに背景色がある', async () => {
    const { container } = render(
      <Tooltip content="テキスト">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = screen.getByText('テキスト');
      expect(tooltip).toHaveClass('bg-gray-900');
    });
  });

  it('ツールチップに白いテキストがある', async () => {
    const { container } = render(
      <Tooltip content="テキスト">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = screen.getByText('テキスト');
      expect(tooltip).toHaveClass('text-white');
    });
  });

  it('ツールチップに角丸がある', async () => {
    const { container } = render(
      <Tooltip content="テキスト">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = screen.getByText('テキスト');
      expect(tooltip).toHaveClass('rounded');
    });
  });

  it('ツールチップにパディングがある', async () => {
    const { container } = render(
      <Tooltip content="テキスト">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = screen.getByText('テキスト');
      expect(tooltip).toHaveClass('px-2', 'py-1');
    });
  });

  it('デフォルト位置は上', async () => {
    const { container } = render(
      <Tooltip content="デフォルト">
        <button>ホバー</button>
      </Tooltip>
    );

    const trigger = screen.getByText('ホバー');
    fireEvent.mouseEnter(trigger);

    await waitFor(() => {
      const tooltip = container.querySelector('.bottom-full');
      expect(tooltip).toBeInTheDocument();
    });
  });

  it('複数のツールチップが独立して動作する', async () => {
    render(
      <>
        <Tooltip content="ツールチップ1">
          <button>ボタン1</button>
        </Tooltip>
        <Tooltip content="ツールチップ2">
          <button>ボタン2</button>
        </Tooltip>
      </>
    );

    const button1 = screen.getByText('ボタン1');
    fireEvent.mouseEnter(button1);

    await waitFor(() => {
      expect(screen.getByText('ツールチップ1')).toBeInTheDocument();
      expect(screen.queryByText('ツールチップ2')).not.toBeInTheDocument();
    });
  });
});
