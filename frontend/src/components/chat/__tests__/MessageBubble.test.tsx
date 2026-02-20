import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import MessageBubble from '../MessageBubble';

const baseProps = {
  content: 'テストメッセージ',
  createdAt: '2026-02-19T14:30:00Z',
  isOwn: false,
};

describe('MessageBubble', () => {
  it('メッセージ内容を表示する', () => {
    render(<MessageBubble {...baseProps} />);
    expect(screen.getByText('テストメッセージ')).toBeInTheDocument();
  });

  it('時刻をHH:mm形式で表示する', () => {
    render(<MessageBubble {...baseProps} />);
    expect(screen.getByText('23:30')).toBeInTheDocument();
  });

  it('自分のメッセージはjustify-endクラスを持つ', () => {
    const { container } = render(<MessageBubble {...baseProps} isOwn={true} />);
    expect(container.firstChild).toHaveClass('justify-end');
  });

  it('他人のメッセージはjustify-startクラスを持つ', () => {
    const { container } = render(<MessageBubble {...baseProps} isOwn={false} />);
    expect(container.firstChild).toHaveClass('justify-start');
  });

  it('自分のメッセージにchat-bubble-ownクラスが適用される', () => {
    render(<MessageBubble {...baseProps} isOwn={true} />);
    const bubble = screen.getByText('テストメッセージ').parentElement;
    expect(bubble).toHaveClass('chat-bubble-own');
  });

  it('他人のメッセージにchat-bubble-otherクラスが適用される', () => {
    render(<MessageBubble {...baseProps} isOwn={false} />);
    const bubble = screen.getByText('テストメッセージ').parentElement;
    expect(bubble).toHaveClass('chat-bubble-other');
  });

  it('showSenderInfo=trueで送信者名を表示する', () => {
    render(
      <MessageBubble {...baseProps} showSenderInfo senderName="Alice" />
    );
    expect(screen.getByText('Alice')).toBeInTheDocument();
  });

  it('showSenderInfo=falseでは送信者名を表示しない', () => {
    render(
      <MessageBubble {...baseProps} showSenderInfo={false} senderName="Alice" />
    );
    expect(screen.queryByText('Alice')).not.toBeInTheDocument();
  });

  it('自分のメッセージではshowSenderInfoがあっても送信者名を表示しない', () => {
    render(
      <MessageBubble {...baseProps} isOwn={true} showSenderInfo senderName="Me" />
    );
    expect(screen.queryByText('Me')).not.toBeInTheDocument();
  });

  it('自分のメッセージにchat-bubble-time-ownクラスの時刻を表示する', () => {
    render(<MessageBubble {...baseProps} isOwn={true} />);
    const timeEl = screen.getByText('23:30');
    expect(timeEl).toHaveClass('chat-bubble-time-own');
  });

  it('他人のメッセージにchat-bubble-time-otherクラスの時刻を表示する', () => {
    render(<MessageBubble {...baseProps} isOwn={false} />);
    const timeEl = screen.getByText('23:30');
    expect(timeEl).toHaveClass('chat-bubble-time-other');
  });
});
