import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { X, LogOut, Trash2 } from 'lucide-react';
import toast from 'react-hot-toast';
import {
  getChatRoomMembers, updateChatRoom, deleteChatRoom,
  addChatRoomMember, removeChatRoomMember,
} from '../../api/chatRooms';
import { inputClass, labelClass, textareaClass, buttonPrimaryClass } from '../../constants/styles';
import { useConfirm } from '../../hooks';
import type { User } from '../../types/user';
import type { ChatRoom, ChatRoomMember } from '../../types/chat';
import ConfirmDialog from '../common/ConfirmDialog';
import ChatRoomMemberSection from './ChatRoomMemberSection';

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
  const [loading, setLoading] = useState(false);
  const { confirm, dialogProps } = useConfirm();

  const isOwner = room.owner_id === currentUserId;

  useEffect(() => {
    getChatRoomMembers(room.id)
      .then(({ data }) => setMembers(data || []))
      .catch(() => { toast.error(t('chat.loadMembersFailed')); });
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
      toast.error(t('chat.updateRoomFailed'));
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
      toast.error(t('chat.addMemberFailed'));
    }
  };

  const handleRemoveMember = async (userId: number) => {
    try {
      await removeChatRoomMember(room.id, userId);
      setMembers((prev) => prev.filter((m) => m.user_id !== userId));
    } catch {
      toast.error(t('chat.removeMemberFailed'));
    }
  };

  const handleLeave = async () => {
    const confirmed = await confirm({
      title: t('chat.leaveGroup'),
      message: t('chat.confirmLeave'),
      variant: 'warning',
      confirmText: t('chat.leaveGroup'),
    });
    if (!confirmed) return;
    try {
      await removeChatRoomMember(room.id, currentUserId);
      onLeft();
    } catch {
      toast.error(t('chat.leaveFailed'));
    }
  };

  const handleDelete = async () => {
    const confirmed = await confirm({
      title: t('common.delete'),
      message: t('chat.confirmDelete'),
      variant: 'danger',
      confirmText: t('common.delete'),
    });
    if (!confirmed) return;
    try {
      await deleteChatRoom(room.id);
      onDeleted();
    } catch {
      toast.error(t('chat.chatDeleteFailed'));
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-800 rounded-md p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-white">{t('chat.roomSettings')}</h2>
          <button onClick={onClose} className="p-1 text-gray-400 hover:text-white transition-colors" aria-label={t('common.close')}>
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Edit Room Info (Owner only) */}
        {isOwner && (
          <div className="space-y-3 mb-6">
            <div>
              <label className={labelClass}>
                {t('chat.groupName')}
              </label>
              <input
                type="text"
                value={name}
                onChange={(e) => setName(e.target.value)}
                maxLength={100}
                className={inputClass}
              />
            </div>
            <div>
              <label className={labelClass}>
                {t('chat.groupDescription')}
              </label>
              <textarea
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                rows={2}
                maxLength={500}
                className={textareaClass}
              />
            </div>
            <button
              onClick={handleUpdate}
              disabled={!name.trim() || loading}
              className={`${buttonPrimaryClass} w-full text-sm disabled:opacity-50`}
            >
              {t('chat.editRoom')}
            </button>
          </div>
        )}

        {/* Members */}
        <ChatRoomMemberSection
          members={members}
          availableUsers={availableUsers}
          isOwner={isOwner}
          currentUserId={currentUserId}
          ownerUserId={room.owner_id}
          onAddMember={handleAddMember}
          onRemoveMember={handleRemoveMember}
        />

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
        <ConfirmDialog {...dialogProps} />
      </div>
    </div>
  );
}
