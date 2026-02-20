import { useTranslation } from 'react-i18next';
import { MessageSquare, Users, Plus } from 'lucide-react';
import type { Conversation } from '../../types/message';
import type { ChatRoom } from '../../types/chat';
import type { User } from '../../types/user';
import Avatar from '../common/Avatar';

interface ChatSidebarProps {
  activeTab: 'dm' | 'group';
  onTabChange: (tab: 'dm' | 'group') => void;
  conversations: Conversation[];
  followingWithoutConversation: User[];
  selectedUserId: number | null;
  onSelectUser: (userId: number) => void;
  chatRooms: ChatRoom[];
  activeRoomId: number | null;
  onSelectRoom: (roomId: number) => void;
  onCreateRoom: () => void;
}

export default function ChatSidebar({
  activeTab,
  onTabChange,
  conversations,
  followingWithoutConversation,
  selectedUserId,
  onSelectUser,
  chatRooms,
  activeRoomId,
  onSelectRoom,
  onCreateRoom,
}: ChatSidebarProps) {
  const { t } = useTranslation();

  return (
    <div className="w-80 border-r border-gray-800 flex flex-col">
      {/* Tab Switcher */}
      <div className="flex border-b border-gray-800">
        <button
          onClick={() => onTabChange('dm')}
          className={`flex-1 flex items-center justify-center gap-2 px-4 py-3 text-sm font-medium transition-colors ${
            activeTab === 'dm'
              ? 'text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-white'
          }`}
        >
          <MessageSquare className="w-4 h-4" />
          {t('chat.dmTab')}
        </button>
        <button
          onClick={() => onTabChange('group')}
          className={`flex-1 flex items-center justify-center gap-2 px-4 py-3 text-sm font-medium transition-colors ${
            activeTab === 'group'
              ? 'text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-white'
          }`}
        >
          <Users className="w-4 h-4" />
          {t('chat.groupTab')}
        </button>
      </div>

      <div className="flex-1 overflow-y-auto">
        {activeTab === 'dm' ? (
          <>
            {conversations.length > 0 && (
              <div>
                <div className="px-5 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider">
                  {t('chat.recentChats')}
                </div>
                {conversations.filter((conv) => conv.user).map((conv) => (
                  <button
                    key={conv.user.id}
                    onClick={() => onSelectUser(conv.user.id)}
                    className={`w-full flex items-center gap-3 px-5 py-3 transition-colors text-left border-l-2 ${
                      selectedUserId === conv.user.id && activeTab === 'dm'
                        ? 'bg-gray-800/70 border-l-blue-500'
                        : 'border-l-transparent hover:bg-gray-800/40'
                    }`}
                  >
                    <Avatar name={conv.user.name} avatarUrl={conv.user.avatar_url} size="sm" />
                    <div className="flex-1 min-w-0">
                      <div className="font-medium text-sm">{conv.user.name}</div>
                      {conv.last_message && (
                        <div className="text-xs text-gray-500 truncate mt-0.5">{conv.last_message.content}</div>
                      )}
                    </div>
                    {conv.unread_count > 0 && (
                      <span className="bg-blue-600 text-white text-xs font-medium rounded-full min-w-[1.25rem] h-5 flex items-center justify-center px-1.5">
                        {conv.unread_count}
                      </span>
                    )}
                  </button>
                ))}
              </div>
            )}

            {followingWithoutConversation.length > 0 && (
              <div>
                <div className="px-5 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider border-t border-gray-800 mt-2 pt-3">
                  {t('chat.following')}
                </div>
                {followingWithoutConversation.map((user) => (
                  <button
                    key={user.id}
                    onClick={() => onSelectUser(user.id)}
                    className={`w-full flex items-center gap-3 px-5 py-3 transition-colors text-left border-l-2 ${
                      selectedUserId === user.id && activeTab === 'dm'
                        ? 'bg-gray-800/70 border-l-green-500'
                        : 'border-l-transparent hover:bg-gray-800/40'
                    }`}
                  >
                    <Avatar name={user.name} avatarUrl={user.avatar_url} size="sm" />
                    <div className="flex-1 min-w-0">
                      <div className="font-medium text-sm">{user.name}</div>
                      <div className="text-xs text-gray-500 mt-0.5">{t('chat.startNewChat')}</div>
                    </div>
                  </button>
                ))}
              </div>
            )}

            {conversations.length === 0 && followingWithoutConversation.length === 0 && (
              <div className="p-6 text-center text-gray-500 text-sm">
                {t('chat.noConversations')}
              </div>
            )}
          </>
        ) : (
          <>
            <div className="p-3">
              <button
                onClick={onCreateRoom}
                className="w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-white rounded-lg text-sm font-medium transition-colors"
              >
                <Plus className="w-4 h-4" />
                {t('chat.createGroup')}
              </button>
            </div>

            {chatRooms.length > 0 ? (
              chatRooms.map((room) => (
                <button
                  key={room.id}
                  onClick={() => onSelectRoom(room.id)}
                  className={`w-full flex items-center gap-3 px-5 py-3 transition-colors text-left border-l-2 ${
                    activeRoomId === room.id
                      ? 'bg-gray-800/70 border-l-blue-500'
                      : 'border-l-transparent hover:bg-gray-800/40'
                  }`}
                >
                  <div className="w-8 h-8 bg-gray-700 rounded-full flex items-center justify-center">
                    <Users className="w-4 h-4 text-gray-300" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-sm">{room.name}</div>
                    {room.description && (
                      <div className="text-xs text-gray-500 truncate mt-0.5">{room.description}</div>
                    )}
                  </div>
                </button>
              ))
            ) : (
              <div className="p-6 text-center text-gray-500 text-sm">
                {t('chat.noRooms')}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
