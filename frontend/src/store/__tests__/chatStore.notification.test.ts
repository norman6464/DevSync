import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useChatStore } from '../chatStore';

/**
 * WebSocket のふりをする最小のスタブ。
 * connect() が組み立てた onmessage をテストから直接呼べるようにする。
 */
class FakeWebSocket {
  static last: FakeWebSocket | null = null;
  onopen: (() => void) | null = null;
  onclose: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  close = vi.fn();

  constructor() {
    FakeWebSocket.last = this;
  }
}

describe('chatStore の通知メッセージ', () => {
  beforeEach(() => {
    vi.stubGlobal('WebSocket', FakeWebSocket);
    useChatStore.setState({
      socket: null,
      connected: false,
      notificationSignal: 0,
      groupMessages: [],
      activeMessages: [],
      activeRoomId: null,
    });
    useChatStore.getState().connect();
  });

  const receive = (payload: unknown) => {
    FakeWebSocket.last?.onmessage?.({ data: JSON.stringify(payload) });
  };

  it('通知を受け取るとシグナルが進む', () => {
    receive({ type: 'notification', id: 1, notification_type: 'follow', actor_id: 2 });

    expect(useChatStore.getState().notificationSignal).toBe(1);
  });

  it('通知が続けて届いた回数だけシグナルが進む', () => {
    receive({ type: 'notification', id: 1, notification_type: 'follow', actor_id: 2 });
    receive({ type: 'notification', id: 2, notification_type: 'like', actor_id: 3 });

    expect(useChatStore.getState().notificationSignal).toBe(2);
  });

  it('通知はメッセージ一覧に混ざらない', () => {
    receive({ type: 'notification', id: 1, notification_type: 'follow', actor_id: 2 });

    expect(useChatStore.getState().activeMessages).toHaveLength(0);
    expect(useChatStore.getState().groupMessages).toHaveLength(0);
  });

  it('グループメッセージではシグナルが進まない', () => {
    useChatStore.setState({ activeRoomId: 5 });
    receive({
      type: 'group_message',
      sender_id: 2,
      room_id: 5,
      content: 'hello',
      sender_name: 'bob',
    });

    expect(useChatStore.getState().notificationSignal).toBe(0);
    expect(useChatStore.getState().groupMessages).toHaveLength(1);
  });

  it('壊れた JSON では何も起きない', () => {
    FakeWebSocket.last?.onmessage?.({ data: 'not json' });

    expect(useChatStore.getState().notificationSignal).toBe(0);
  });
});
