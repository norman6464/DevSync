import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import ChatSidebar from '../ChatSidebar';

const makeConversation = (overrides = {}) => ({
  user: { id: 1, name: 'Alice', avatar_url: '' },
  last_message: { content: '最新メッセージ' },
  unread_count: 0,
  ...overrides,
});

const makeRoom = (overrides = {}) => ({
  id: 1,
  name: 'テストルーム',
  description: 'テスト用グループ',
  owner_id: 1,
  created_at: '',
  updated_at: '',
  ...overrides,
});

const makeUser = (overrides = {}) => ({
  id: 10,
  name: 'Bob',
  avatar_url: '',
  ...overrides,
});

const defaultProps = {
  activeTab: 'dm' as const,
  onTabChange: vi.fn(),
  conversations: [],
  followingWithoutConversation: [],
  selectedUserId: null,
  onSelectUser: vi.fn(),
  chatRooms: [],
  activeRoomId: null,
  onSelectRoom: vi.fn(),
  onCreateRoom: vi.fn(),
};

const renderSidebar = (props = {}) =>
  render(<ChatSidebar {...defaultProps} {...props} />);

describe('ChatSidebar', () => {
  it('DMタブとグループタブが表示される', () => {
    renderSidebar();
    expect(screen.getByText('DM')).toBeInTheDocument();
    expect(screen.getByText('グループ')).toBeInTheDocument();
  });

  it('グループタブクリックでonTabChangeが呼ばれる', () => {
    const onTabChange = vi.fn();
    renderSidebar({ onTabChange });
    fireEvent.click(screen.getByText('グループ'));
    expect(onTabChange).toHaveBeenCalledWith('group');
  });

  it('DMタブで会話一覧が表示される', () => {
    renderSidebar({
      conversations: [makeConversation()],
    });
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('最新メッセージ')).toBeInTheDocument();
  });

  it('未読数が表示される', () => {
    renderSidebar({
      conversations: [makeConversation({ unread_count: 3 })],
    });
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  it('会話クリックでonSelectUserが呼ばれる', () => {
    const onSelectUser = vi.fn();
    renderSidebar({
      conversations: [makeConversation()],
      onSelectUser,
    });
    fireEvent.click(screen.getByText('Alice'));
    expect(onSelectUser).toHaveBeenCalledWith(1);
  });

  it('フォロー中ユーザーが表示される', () => {
    renderSidebar({
      followingWithoutConversation: [makeUser()],
    });
    expect(screen.getByText('Bob')).toBeInTheDocument();
    expect(screen.getByText('新しいチャットを開始')).toBeInTheDocument();
  });

  it('会話もフォローも空の場合に空メッセージが表示される', () => {
    renderSidebar();
    expect(screen.getByText('まだ会話がありません')).toBeInTheDocument();
  });

  it('グループタブでルーム一覧が表示される', () => {
    renderSidebar({
      activeTab: 'group',
      chatRooms: [makeRoom()],
    });
    expect(screen.getByText('テストルーム')).toBeInTheDocument();
    expect(screen.getByText('テスト用グループ')).toBeInTheDocument();
  });

  it('グループ作成ボタンクリックでonCreateRoomが呼ばれる', () => {
    const onCreateRoom = vi.fn();
    renderSidebar({ activeTab: 'group', onCreateRoom });
    fireEvent.click(screen.getByText('グループ作成'));
    expect(onCreateRoom).toHaveBeenCalled();
  });

  it('ルームクリックでonSelectRoomが呼ばれる', () => {
    const onSelectRoom = vi.fn();
    renderSidebar({
      activeTab: 'group',
      chatRooms: [makeRoom()],
      onSelectRoom,
    });
    fireEvent.click(screen.getByText('テストルーム'));
    expect(onSelectRoom).toHaveBeenCalledWith(1);
  });

  it('グループタブでルームが空の場合に空メッセージが表示される', () => {
    renderSidebar({ activeTab: 'group' });
    expect(screen.getByText('グループがありません')).toBeInTheDocument();
  });
});
