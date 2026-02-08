import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { X, UserPlus, UserMinus, LogOut, Trash2 } from 'lucide-react';
import {
  getChatRoomMembers, updateChatRoom, deleteChatRoom,
  addChatRoomMember, removeChatRoomMember,
} from '../../api/chatRooms';
import type { User } from '../../types/user';
import type { ChatRoom, ChatRoomMember } from '../../types/chat';
import Avatar from '../common/Avatar';

interface Props {
  room: ChatRoom;
  currentUserId: number;
  followingUsers: User[];
  onClose: () => void;
  onUpdated: () => void;
  onDeleted: () => void;
  onLeft: () => void;
}

export default function RoomSettingsModal({
  room, currentUserId, followingUsers,
  onClose, onUpdated, onDeleted, onLeft,
}: Props) {
  const { t } = useTranslation();
  const [members, setMembers] = useState<ChatRoomMember[]>([]);
  const [name, setName] = useState(room.name);
  const [description, setDescription] = useState(room.description);
  const [showAddMember, setShowAddMember] = useState(false);
  const [loading, setLoading] = useState(false);

  const isOwner = room.owner_id === currentUserId;

  useEffect(() => {
    getChatRoomMembers(room.id)
      .then(({ data }) => setMembers(data || []))
      .catch(() => {});
  }, [room.id]);

  const memberUserIds = members.map((m) => m.user_id);
  const availableUsers = followingUsers.filter((u) => !memberUserIds.includes(u.id));

  const handleUpdate = async () => {
    if (!name.trim()) return;
    setLoading(true);
    try {
      await updateChatRoom(room.id, { name: name.trim(), description: description.trim() });
      onUpdated();
    } catch {
      // handle error
    } finally {
      setLoading(false);
    }
  };

  const handleAddMember = async (userId: number) => {
    try {
      await addChatRoomMember(room.id, userId);
      const { data } = await getChatRoomMembers(room.id);
      setMembers(data || []);
    } catch {
      // handle error
    }
  };

  const handleRemoveMember = async (userId: number) => {
    try {
      await removeChatRoomMember(room.id, userId);
      setMembers((prev) => prev.filter((m) => m.user_id !== userId));
    } catch {
      // handle error
    }
  };

  const handleLeave = async () => {
    if (!confirm(t('chat.confirmLeave'))) return;
    try {
      await removeChatRoomMember(room.id, currentUserId);
      onLeft();
    } catch {
      // handle error
    }
  };

  const handleDelete = async () => {
    if (!confirm(t('chat.confirmDelete'))) return;
    try {
      await deleteChatRoom(room.id);
      onDeleted();
    } catch {
      // handle error
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-800 rounded-xl p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-white">{t('chat.roomSettings')}</h2>
          <button onClick={onClose} className="p-1 text-gray-400 hover:text-white transition-colors">
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Edit Room Info (Owner only) */}
        {isOwner && (
          <div className="space-y-3 mb-6">
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-1">
                {t('chat.groupName')}
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-1">
                {t('chat.groupDescription')}
              </label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={2}
                className="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
              />
            </div>
            <button
              onClick={handleUpdate}
              disabled={!name.trim() || loading}
              className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors"
            >
              {t('chat.editRoom')}
            </button>
          </div>
        )}

        {/* Members */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-sm font-medium text-gray-300">
              {t('chat.members')} ({members.length})
            </h3>
            <button
              onClick={() => setShowAddMember(!showAddMember)}
              className="flex items-center gap-1 text-xs text-blue-400 hover:text-blue-300"
            >
              <UserPlus className="w-3.5 h-3.5" />
              {t('chat.addMember')}
            </button>
          </div>

          {/* Add Member List */}
          {showAddMember && availableUsers.length > 0 && (
            <div className="mb-3 p-2 bg-gray-700/50 rounded-lg space-y-1">
              {availableUsers.map((user) => (
                <button
                  key={user.id}
                  onClick={() => handleAddMember(user.id)}
                  className="w-full flex items-center gap-2 px-2 py-1.5 rounded hover:bg-gray-600 transition-colors"
                >
                  <Avatar name={user.name} avatarUrl={user.avatar_url} size="xs" />
                  <span className="text-sm text-white">{user.name}</span>
                  <UserPlus className="w-3.5 h-3.5 text-green-400 ml-auto" />
                </button>
              ))}
            </div>
          )}

          {/* Member List */}
          <div className="space-y-1">
            {members.map((member) => (
              <div
                key={member.id}
                className="flex items-center gap-3 px-2 py-2 rounded-lg"
              >
                <Avatar
                  name={member.user?.name || ''}
                  avatarUrl={member.user?.avatar_url}
                  size="sm"
                />
                <div className="flex-1 min-w-0">
                  <div className="text-sm text-white">{member.user?.name}</div>
                  {member.user_id === room.owner_id && (
                    <span className="text-xs text-yellow-400">{t('chat.owner')}</span>
                  )}
                </div>
                {isOwner && member.user_id !== currentUserId && (
                  <button
                    onClick={() => handleRemoveMember(member.user_id)}
                    className="p-1 text-gray-500 hover:text-red-400 transition-colors"
                    title={t('chat.removeMember')}
                  >
                    <UserMinus className="w-4 h-4" />
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Actions */}
        <div className="space-y-2 border-t border-gray-700 pt-4">
          {!isOwner && (
            <button
              onClick={handleLeave}
              className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-yellow-400 rounded-lg text-sm transition-colors"
            >
              <LogOut className="w-4 h-4" />
              {t('chat.leaveGroup')}
            </button>
          )}
          {isOwner && (
            <button
              onClick={handleDelete}
              className="w-full flex items-center justify-center gap-2 px-4 py-2 bg-red-600/20 hover:bg-red-600/30 text-red-400 rounded-lg text-sm transition-colors"
            >
              <Trash2 className="w-4 h-4" />
              {t('chat.deleteGroup')}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
