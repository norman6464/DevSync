import { useTranslation } from 'react-i18next';
import { MessageSquare, Users, Plus, Settings, Send } from 'lucide-react';
import { useChat } from '../hooks/useChat';
import type { Message } from '../types/message';
import Avatar from '../components/common/Avatar';
import CreateRoomModal from '../components/chat/CreateRoomModal';
import RoomSettingsModal from '../components/chat/RoomSettingsModal';
import { format } from 'date-fns';

export default function ChatPage() {
  const { t } = useTranslation();
  const c = useChat();

  return (
    <div className="flex h-[calc(100vh-7rem)] bg-gray-900 border border-gray-800 rounded-md overflow-hidden">
      {/* Sidebar */}
      <div className="w-80 border-r border-gray-800 flex flex-col">
        {/* Tab Switcher */}
        <div className="flex border-b border-gray-800">
          <button
            onClick={() => c.setActiveTab('dm')}
            className={`flex-1 flex items-center justify-center gap-2 px-4 py-3 text-sm font-medium transition-colors ${
              c.activeTab === 'dm'
                ? 'text-white border-b-2 border-blue-500'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            <MessageSquare className="w-4 h-4" />
            {t('chat.dmTab')}
          </button>
          <button
            onClick={() => c.setActiveTab('group')}
            className={`flex-1 flex items-center justify-center gap-2 px-4 py-3 text-sm font-medium transition-colors ${
              c.activeTab === 'group'
                ? 'text-white border-b-2 border-blue-500'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            <Users className="w-4 h-4" />
            {t('chat.groupTab')}
          </button>
        </div>

        <div className="flex-1 overflow-y-auto">
          {c.activeTab === 'dm' ? (
            <>
              {c.conversations.length > 0 && (
                <div>
                  <div className="px-5 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    {t('chat.recentChats')}
                  </div>
                  {c.conversations.map((conv) => (
                    <button
                      key={conv.user.id}
                      onClick={() => c.handleSelectUser(conv.user.id)}
                      className={`w-full flex items-center gap-3 px-5 py-3 transition-colors text-left border-l-2 ${
                        c.selectedUserId === conv.user.id && c.activeTab === 'dm'
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

              {c.followingWithoutConversation.length > 0 && (
                <div>
                  <div className="px-5 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider border-t border-gray-800 mt-2 pt-3">
                    {t('chat.following')}
                  </div>
                  {c.followingWithoutConversation.map((user) => (
                    <button
                      key={user.id}
                      onClick={() => c.handleSelectUser(user.id)}
                      className={`w-full flex items-center gap-3 px-5 py-3 transition-colors text-left border-l-2 ${
                        c.selectedUserId === user.id && c.activeTab === 'dm'
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

              {c.conversations.length === 0 && c.followingWithoutConversation.length === 0 && (
                <div className="p-6 text-center text-gray-500 text-sm">
                  {t('chat.noConversations')}
                </div>
              )}
            </>
          ) : (
            <>
              <div className="p-3">
                <button
                  onClick={() => c.setShowCreateRoom(true)}
                  className="w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-white rounded-lg text-sm font-medium transition-colors"
                >
                  <Plus className="w-4 h-4" />
                  {t('chat.createGroup')}
                </button>
              </div>

              {c.chatRooms.length > 0 ? (
                c.chatRooms.map((room) => (
                  <button
                    key={room.id}
                    onClick={() => c.handleSelectRoom(room.id)}
                    className={`w-full flex items-center gap-3 px-5 py-3 transition-colors text-left border-l-2 ${
                      c.activeRoomId === room.id
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

      {/* Chat Area */}
      <div className="flex-1 flex flex-col">
        {c.activeTab === 'group' && c.activeRoomId && c.selectedRoom ? (
          <>
            <div className="px-6 py-3 border-b border-gray-800 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-gray-700 rounded-full flex items-center justify-center">
                  <Users className="w-4 h-4 text-gray-300" />
                </div>
                <div>
                  <div className="font-medium text-sm">{c.selectedRoom.name}</div>
                  {c.selectedRoom.description && (
                    <div className="text-xs text-gray-500">{c.selectedRoom.description}</div>
                  )}
                </div>
              </div>
              <button
                onClick={() => c.setShowRoomSettings(true)}
                className="p-2 text-gray-400 hover:text-white transition-colors rounded-md"
              >
                <Settings className="w-5 h-5" />
              </button>
            </div>

            <div className="flex-1 overflow-y-auto p-6 space-y-3">
              {c.groupMessages.map((msg) => {
                const isOwn = msg.sender_id === c.currentUser?.id;
                return (
                  <div
                    key={msg.id}
                    className={`flex items-end gap-2 ${isOwn ? 'justify-end' : 'justify-start'}`}
                  >
                    {!isOwn && (
                      <Avatar
                        name={msg.sender?.name || ''}
                        avatarUrl={msg.sender?.avatar_url}
                        size="xs"
                      />
                    )}
                    {isOwn && (
                      <span className="chat-bubble-time-own text-xs mb-0.5">
                        {format(new Date(msg.created_at), 'HH:mm')}
                      </span>
                    )}
                    <div>
                      {!isOwn && (
                        <p className="text-xs text-gray-400 mb-1">{msg.sender?.name}</p>
                      )}
                      <div
                        className={`px-4 py-2.5 rounded-md ${
                          isOwn
                            ? 'chat-bubble-own rounded-br-md'
                            : 'chat-bubble-other rounded-bl-md'
                        }`}
                      >
                        <p className="text-sm leading-relaxed">{msg.content}</p>
                      </div>
                    </div>
                    {!isOwn && (
                      <span className="chat-bubble-time-other text-xs mb-0.5">
                        {format(new Date(msg.created_at), 'HH:mm')}
                      </span>
                    )}
                  </div>
                );
              })}
            </div>

            <form onSubmit={c.handleSend} className="p-4 border-t border-gray-800 flex gap-3">
              <input
                type="text"
                value={c.newMessage}
                onChange={(e) => c.setNewMessage(e.target.value)}
                placeholder={t('chat.groupMessagePlaceholder')}
                className="flex-1 px-4 py-2.5 bg-gray-800/50 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow"
              />
              <button
                type="submit"
                disabled={!c.newMessage.trim()}
                className="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 disabled:hover:bg-gray-700 text-white rounded-lg font-medium text-sm transition-colors flex items-center gap-2"
              >
                <Send className="w-4 h-4" />
                {t('chat.send')}
              </button>
            </form>
          </>
        ) : c.selectedUserId ? (
          <>
            <div className="px-6 py-3 border-b border-gray-800 flex items-center gap-3">
              {c.selectedUser && (
                <>
                  <Avatar
                    name={c.selectedUser.name}
                    avatarUrl={c.selectedUser.avatar_url}
                    size="sm"
                  />
                  <div>
                    <div className="font-medium text-sm">{c.selectedUser.name}</div>
                  </div>
                </>
              )}
            </div>

            <div className="flex-1 overflow-y-auto p-6 space-y-3">
              {c.activeMessages.map((msg: Message) => {
                const isOwn = msg.sender_id === c.currentUser?.id;
                return (
                  <div
                    key={msg.id}
                    className={`flex items-end gap-2 ${isOwn ? 'justify-end' : 'justify-start'}`}
                  >
                    {isOwn && (
                      <span className="chat-bubble-time-own text-xs mb-0.5">
                        {format(new Date(msg.created_at), 'HH:mm')}
                      </span>
                    )}
                    <div
                      className={`max-w-sm px-4 py-2.5 rounded-md ${
                        isOwn
                          ? 'chat-bubble-own rounded-br-md'
                          : 'chat-bubble-other rounded-bl-md'
                      }`}
                    >
                      <p className="text-sm leading-relaxed">{msg.content}</p>
                    </div>
                    {!isOwn && (
                      <span className="chat-bubble-time-other text-xs mb-0.5">
                        {format(new Date(msg.created_at), 'HH:mm')}
                      </span>
                    )}
                  </div>
                );
              })}
            </div>

            <form onSubmit={c.handleSend} className="p-4 border-t border-gray-800 flex gap-3">
              <input
                type="text"
                value={c.newMessage}
                onChange={(e) => c.setNewMessage(e.target.value)}
                placeholder={t('chat.typeMessage')}
                className="flex-1 px-4 py-2.5 bg-gray-800/50 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow"
              />
              <button
                type="submit"
                disabled={!c.newMessage.trim()}
                className="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 disabled:hover:bg-gray-700 text-white rounded-lg font-medium text-sm transition-colors flex items-center gap-2"
              >
                <Send className="w-4 h-4" />
                {t('chat.send')}
              </button>
            </form>
          </>
        ) : (
          <div className="flex-1 flex flex-col items-center justify-center text-gray-500">
            <MessageSquare className="w-16 h-16 text-gray-700 mb-4" />
            <p className="text-sm">{t('chat.startConversation')}</p>
          </div>
        )}
      </div>

      {/* Modals */}
      {c.showCreateRoom && (
        <CreateRoomModal
          followingUsers={c.followingUsers}
          onClose={() => c.setShowCreateRoom(false)}
          onCreated={(room) => {
            c.setChatRooms([room, ...c.chatRooms]);
            c.setActiveRoomId(room.id);
            c.setShowCreateRoom(false);
          }}
        />
      )}

      {c.showRoomSettings && c.selectedRoom && (
        <RoomSettingsModal
          room={c.selectedRoom}
          currentUserId={c.currentUser?.id || 0}
          followingUsers={c.followingUsers}
          onClose={() => c.setShowRoomSettings(false)}
          onUpdated={c.loadChatRooms}
          onDeleted={() => {
            c.setActiveRoomId(null);
            c.setShowRoomSettings(false);
            c.loadChatRooms();
          }}
          onLeft={() => {
            c.setActiveRoomId(null);
            c.setShowRoomSettings(false);
            c.loadChatRooms();
          }}
        />
      )}
    </div>
  );
}
