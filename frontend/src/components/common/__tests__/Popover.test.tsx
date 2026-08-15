import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Popover from '../Popover';

describe('Popover', () => {
  it('トリガーが表示される', () => {
    render(
      <Popover trigger={<button>開く</button>}>
        <p>ポップオーバー内容</p>
      </Popover>
    );
    expect(screen.getByText('開く')).toBeInTheDocument();
  });

  it('初期状態でコンテンツが非表示', () => {
    render(
      <Popover trigger={<button>開く</button>}>
        <p>ポップオーバー内容</p>
      </Popover>
    );
    expect(screen.queryByText('ポップオーバー内容')).not.toBeInTheDocument();
  });

  it('クリックでコンテンツが表示される', async () => {
    const user = userEvent.setup();
    render(
      <Popover trigger={<button>開く</button>}>
        <p>ポップオーバー内容</p>
      </Popover>
    );
    await user.click(screen.getByText('開く'));
    expect(screen.getByText('ポップオーバー内容')).toBeInTheDocument();
  });

  it('再クリックで閉じる', async () => {
    const user = userEvent.setup();
    render(
      <Popover trigger={<button>開く</button>}>
        <p>ポップオーバー内容</p>
      </Popover>
    );
    await user.click(screen.getByText('開く'));
    await user.click(screen.getByText('開く'));
    expect(screen.queryByText('ポップオーバー内容')).not.toBeInTheDocument();
  });

  it('上方向に配置される', async () => {
    const user = userEvent.setup();
    const { container } = render(
      <Popover trigger={<button>開く</button>} position="top">
        <p>内容</p>
      </Popover>
    );
    await user.click(screen.getByText('開く'));
    expect(container.querySelector('.bottom-full')).toBeInTheDocument();
  });

  it('下方向に配置される（デフォルト）', async () => {
    const user = userEvent.setup();
    const { container } = render(
      <Popover trigger={<button>開く</button>}>
        <p>内容</p>
      </Popover>
    );
    await user.click(screen.getByText('開く'));
    expect(container.querySelector('.top-full')).toBeInTheDocument();
  });

  it('左方向に配置される', async () => {
    const user = userEvent.setup();
    const { container } = render(
      <Popover trigger={<button>開く</button>} position="left">
        <p>内容</p>
      </Popover>
    );
    await user.click(screen.getByText('開く'));
    expect(container.querySelector('.right-full')).toBeInTheDocument();
  });

  it('右方向に配置される', async () => {
    const user = userEvent.setup();
    const { container } = render(
      <Popover trigger={<button>開く</button>} position="right">
        <p>内容</p>
      </Popover>
    );
    await user.click(screen.getByText('開く'));
    expect(container.querySelector('.left-full')).toBeInTheDocument();
  });

  it('タイトルが表示される', async () => {
    const user = userEvent.setup();
    render(
      <Popover trigger={<button>開く</button>} title="詳細情報">
        <p>内容</p>
      </Popover>
    );
    await user.click(screen.getByText('開く'));
    expect(screen.getByText('詳細情報')).toBeInTheDocument();
  });

  it('カスタムクラス名が適用される', async () => {
    const user = userEvent.setup();
    render(
      <Popover trigger={<button>開く</button>} className="custom-class">
        <p>内容</p>
      </Popover>
    );
    await user.click(screen.getByText('開く'));
    expect(document.querySelector('.custom-class')).toBeInTheDocument();
  });
});
