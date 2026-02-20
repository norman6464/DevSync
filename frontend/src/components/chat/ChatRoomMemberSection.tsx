import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { UserPlus, UserMinus } from 'lucide-react';
import type { User } from '../../types/user';
import type { ChatRoomMember } from '../../types/chat';
import Avatar from '../common/Avatar';
import { linkSmallClass } from '../../constants/styles';

interface ChatRoomMemberSectionProps {
  members: ChatRoomMember[];
  availableUsers: User[];
  isOwner: boolean;
  currentUserId: number;
  ownerUserId: number;
  onAddMember: (userId: number) => void;
  onRemoveMember: (userId: number) => void;
}

export default function ChatRoomMemberSection({
  members,
  availableUsers,
  isOwner,
  currentUserId,
  ownerUserId,
  onAddMember,
  onRemoveMember,
}: ChatRoomMemberSectionProps) {
  const { t } = useTranslation();
  const [showAddMember, setShowAddMember] = useState(false);

  return (
    <div className="mb-6">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-medium text-gray-300">
          {t('chat.members')} ({members.length})
        </h3>
        <button
          onClick={() => setShowAddMember(!showAddMember)}
          className={`${linkSmallClass} flex items-center gap-1`}
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
              onClick={() => onAddMember(user.id)}
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
              {member.user_id === ownerUserId && (
                <span className="text-xs text-yellow-400">{t('chat.owner')}</span>
              )}
            </div>
            {isOwner && member.user_id !== currentUserId && (
              <button
                onClick={() => onRemoveMember(member.user_id)}
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
  );
}
