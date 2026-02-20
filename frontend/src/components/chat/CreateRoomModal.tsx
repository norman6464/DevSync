import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { X } from 'lucide-react';
import toast from 'react-hot-toast';
import { createChatRoom } from '../../api/chatRooms';
import type { User } from '../../types/user';
import type { ChatRoom } from '../../types/chat';
import Avatar from '../common/Avatar';
import { inputClass, labelClass, textareaClass, buttonSecondaryClass, buttonPrimaryClass } from '../../constants/styles';

interface Props {
  followingUsers: User[];
  onClose: () => void;
  onCreated: (room: ChatRoom) => void;
}

export default function CreateRoomModal({ followingUsers, onClose, onCreated }: Props) {
  const { t } = useTranslation();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [selectedMembers, setSelectedMembers] = useState<number[]>([]);
  const [loading, setLoading] = useState(false);

  const toggleMember = (userId: number) => {
    setSelectedMembers((prev) =>
      prev.includes(userId) ? prev.filter((id) => id !== userId) : [...prev, userId]
    );
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setLoading(true);
    try {
      const { data } = await createChatRoom({
        name: name.trim(),
        description: description.trim(),
        member_ids: selectedMembers,
      });
      onCreated(data);
    } catch {
      toast.error(t('chat.createFailed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-800 rounded-md p-6 w-full max-w-md max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold text-white">{t('chat.createGroup')}</h2>
          <button onClick={onClose} className="p-1 text-gray-400 hover:text-white transition-colors" aria-label={t('common.close')}>
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="room-name" className={labelClass}>
              {t('chat.groupName')}
            </label>
            <input
              id="room-name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder={t('chat.groupNamePlaceholder')}
              maxLength={100}
              className={inputClass}
              required
            />
          </div>

          <div>
            <label htmlFor="room-description" className={labelClass}>
              {t('chat.groupDescription')}
            </label>
            <textarea
              id="room-description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder={t('chat.groupDescriptionPlaceholder')}
              rows={2}
              maxLength={500}
              className={textareaClass}
            />
          </div>

          {followingUsers.length > 0 && (
            <div>
              <label className="block text-sm font-medium text-gray-300 mb-2">
                {t('chat.selectMembers')}
              </label>
              <div className="space-y-1 max-h-48 overflow-y-auto">
                {followingUsers.map((user) => (
                  <button
                    key={user.id}
                    type="button"
                    onClick={() => toggleMember(user.id)}
                    className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
                      selectedMembers.includes(user.id)
                        ? 'bg-blue-600/20 border border-blue-500/30'
                        : 'hover:bg-gray-700'
                    }`}
                  >
                    <Avatar name={user.name} avatarUrl={user.avatar_url} size="sm" />
                    <span className="text-sm text-white">{user.name}</span>
                    {selectedMembers.includes(user.id) && (
                      <span className="ml-auto text-blue-400 text-xs">&#10003;</span>
                    )}
                  </button>
                ))}
              </div>
            </div>
          )}

          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className={`${buttonSecondaryClass} flex-1 text-sm`}
            >
              {t('common.cancel')}
            </button>
            <button
              type="submit"
              disabled={!name.trim() || loading}
              className={`${buttonPrimaryClass} flex-1 text-sm disabled:opacity-50`}
            >
              {loading ? t('common.saving') : t('chat.createButton')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
