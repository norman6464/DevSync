import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import ChatRoomMemberSection from '../ChatRoomMemberSection';
import type { ChatRoomMember } from '../../../types/chat';
import type { User } from '../../../types/user';

const mockMembers: ChatRoomMember[] = [
  { id: 1, chat_room_id: 1, user_id: 10, user: { id: 10, name: 'Alice' } as User, joined_at: '2026-01-01' },
  { id: 2, chat_room_id: 1, user_id: 20, user: { id: 20, name: 'Bob' } as User, joined_at: '2026-01-02' },
];

const mockAvailableUsers: User[] = [
  { id: 30, name: 'Charlie' } as User,
];

const defaultProps = {
  members: mockMembers,
  availableUsers: mockAvailableUsers,
  isOwner: true,
  currentUserId: 10,
  ownerUserId: 10,
  onAddMember: vi.fn(),
  onRemoveMember: vi.fn(),
};

describe('ChatRoomMemberSection', () => {
  it('メンバー数を表示する', () => {
    render(<ChatRoomMemberSection {...defaultProps} />);
    expect(screen.getByText(/メンバー \(2\)/)).toBeInTheDocument();
  });

  it('メンバー名を表示する', () => {
    render(<ChatRoomMemberSection {...defaultProps} />);
    expect(screen.getByText('Alice')).toBeInTheDocument();
    expect(screen.getByText('Bob')).toBeInTheDocument();
  });

  it('オーナーラベルを表示する', () => {
    render(<ChatRoomMemberSection {...defaultProps} />);
    expect(screen.getByText('オーナー')).toBeInTheDocument();
  });

  it('メンバー追加ボタンが表示される', () => {
    render(<ChatRoomMemberSection {...defaultProps} />);
    expect(screen.getByText('メンバー追加')).toBeInTheDocument();
  });

  it('メンバー追加ボタンクリックで追加候補リストが表示される', () => {
    render(<ChatRoomMemberSection {...defaultProps} />);
    fireEvent.click(screen.getByText('メンバー追加'));
    expect(screen.getByText('Charlie')).toBeInTheDocument();
  });

  it('追加候補クリックでonAddMemberが呼ばれる', () => {
    const onAddMember = vi.fn();
    render(<ChatRoomMemberSection {...defaultProps} onAddMember={onAddMember} />);
    fireEvent.click(screen.getByText('メンバー追加'));
    fireEvent.click(screen.getByText('Charlie'));
    expect(onAddMember).toHaveBeenCalledWith(30);
  });

  it('オーナーは他メンバーの削除ボタンが表示される', () => {
    render(<ChatRoomMemberSection {...defaultProps} />);
    const removeButtons = screen.getAllByTitle('メンバー削除');
    expect(removeButtons).toHaveLength(1);
  });

  it('削除ボタンクリックでonRemoveMemberが呼ばれる', () => {
    const onRemoveMember = vi.fn();
    render(<ChatRoomMemberSection {...defaultProps} onRemoveMember={onRemoveMember} />);
    fireEvent.click(screen.getByTitle('メンバー削除'));
    expect(onRemoveMember).toHaveBeenCalledWith(20);
  });

  it('非オーナーの場合は削除ボタンが表示されない', () => {
    render(<ChatRoomMemberSection {...defaultProps} isOwner={false} />);
    expect(screen.queryByTitle('メンバー削除')).not.toBeInTheDocument();
  });

  it('追加候補が空の場合はリストが表示されない', () => {
    render(<ChatRoomMemberSection {...defaultProps} availableUsers={[]} />);
    fireEvent.click(screen.getByText('メンバー追加'));
    expect(screen.queryByText('Charlie')).not.toBeInTheDocument();
  });
});
