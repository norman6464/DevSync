import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../store/authStore';
import { useChatStore } from '../store/chatStore';
import { getConversations, getMessages, sendMessage as sendMessageApi } from '../api/messages';
import { getFollowing } from '../api/users';
import type { Conversation, Message } from '../types/message';
import type { User } from '../types/user';
import Avatar from '../components/common/Avatar';
import { format } from 'date-fns';

export default function ChatPage() {
  const { t } = useTranslation();
  const { userId } = useParams<{ userId: string }>();
  const currentUser = useAuthStore((s) => s.user);
  const token = useAuthStore((s) => s.token);
  const { socket, connect, activeMessages, setActiveMessages } = useChatStore();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [followingUsers, setFollowingUsers] = useState<User[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<number | null>(
    userId ? parseInt(userId) : null
  );
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [newMessage, setNewMessage] = useState('');

  useEffect(() => {
    if (token && !socket) {
      connect(token);
    }
  }, [token, socket, connect]);

  useEffect(() => {
    // Load conversations
    getConversations()
      .then(({ data }) => setConversations(data || []))
      .catch(() => {});

    // Load following users
    if (currentUser) {
      getFollowing(currentUser.id)
        .then(({ data }) => setFollowingUsers(data || []))
        .catch(() => {});
    }
  }, [currentUser]);

  useEffect(() => {
    if (selectedUserId) {
      getMessages(selectedUserId)
        .then(({ data }) => setActiveMessages(data || []))
        .catch(() => setActiveMessages([]));

      // Find user info from conversations or following
      const convUser = conversations.find((c) => c.user.id === selectedUserId)?.user;
      const followUser = followingUsers.find((u) => u.id === selectedUserId);
      setSelectedUser(convUser || followUser || null);
    }
  }, [selectedUserId, setActiveMessages, conversations, followingUsers]);

  // Filter following users that don't have existing conversations
  const followingWithoutConversation = followingUsers.filter(
    (user) => !conversations.some((conv) => conv.user.id === user.id)
  );

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim() || !selectedUserId) return;

    try {
      const { data } = await sendMessageApi(selectedUserId, newMessage);
      setActiveMessages([...activeMessages, data]);
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(JSON.stringify({
          type: 'message',
          receiver_id: selectedUserId,
          content: newMessage,
        }));
      }
      setNewMessage('');
    } catch {
      // handle error
    }
  };

  return (
    <div className="flex h-[calc(100vh-7rem)] bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
      {/* Sidebar */}
      <div className="w-80 border-r border-gray-800 flex flex-col">
        <div className="px-5 py-4 border-b border-gray-800">
          <h2 className="font-semibold text-sm flex items-center gap-2">
            <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M20.25 8.511c.884.284 1.5 1.128 1.5 2.097v4.286c0 1.136-.847 2.1-1.98 2.193-.34.027-.68.052-1.02.072v3.091l-3-3c-1.354 0-2.694-.055-4.02-.163a2.115 2.115 0 0 1-.825-.242m9.345-8.334a2.126 2.126 0 0 0-.476-.095 48.64 48.64 0 0 0-8.048 0c-1.131.094-1.976 1.057-1.976 2.192v4.286c0 .837.46 1.58 1.155 1.951m9.345-8.334V6.637c0-1.621-1.152-3.026-2.76-3.235A48.455 48.455 0 0 0 11.25 3c-2.115 0-4.198.137-6.24.402-1.608.209-2.76 1.614-2.76 3.235v6.226c0 1.621 1.152 3.026 2.76 3.235.577.075 1.157.14 1.74.194V21l4.155-4.155" />
            </svg>
            {t('chat.title')}
          </h2>
        </div>
        <div className="flex-1 overflow-y-auto">
          {/* Existing Conversations */}
          {conversations.length > 0 && (
            <div>
              <div className="px-5 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider">
                {t('chat.recentChats')}
              </div>
              {conversations.map((conv) => (
                <button
                  key={conv.user.id}
                  onClick={() => setSelectedUserId(conv.user.id)}
                  className={`w-full flex items-center gap-3 px-5 py-3 transition-colors text-left border-l-2 ${
                    selectedUserId === conv.user.id
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

          {/* Following Users without Conversation */}
          {followingWithoutConversation.length > 0 && (
            <div>
              <div className="px-5 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider border-t border-gray-800 mt-2 pt-3">
                {t('chat.following')}
              </div>
              {followingWithoutConversation.map((user) => (
                <button
                  key={user.id}
                  onClick={() => setSelectedUserId(user.id)}
                  className={`w-full flex items-center gap-3 px-5 py-3 transition-colors text-left border-l-2 ${
                    selectedUserId === user.id
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

          {/* Empty State */}
          {conversations.length === 0 && followingWithoutConversation.length === 0 && (
            <div className="p-6 text-center text-gray-500 text-sm">
              {t('chat.noConversations')}
            </div>
          )}
        </div>
      </div>

      {/* Chat Area */}
      <div className="flex-1 flex flex-col">
        {selectedUserId ? (
          <>
            {/* Chat Header */}
            <div className="px-6 py-3 border-b border-gray-800 flex items-center gap-3">
              {selectedUser && (
                <>
                  <Avatar
                    name={selectedUser.name}
                    avatarUrl={selectedUser.avatar_url}
                    size="sm"
                  />
                  <div>
                    <div className="font-medium text-sm">{selectedUser.name}</div>
                  </div>
                </>
              )}
            </div>

            {/* Messages */}
            <div className="flex-1 overflow-y-auto p-6 space-y-4">
              {activeMessages.map((msg: Message) => {
                const isOwn = msg.sender_id === currentUser?.id;
                return (
                  <div
                    key={msg.id}
                    className={`flex ${isOwn ? 'justify-end' : 'justify-start'}`}
                  >
                    <div
                      className={`max-w-sm px-4 py-2.5 rounded-2xl ${
                        isOwn
                          ? 'bg-blue-600 text-white rounded-br-md'
                          : 'bg-gray-800 text-gray-100 rounded-bl-md'
                      }`}
                    >
                      <p className="text-sm leading-relaxed">{msg.content}</p>
                      <p className={`text-xs mt-1 ${isOwn ? 'text-blue-200' : 'text-gray-500'}`}>
                        {format(new Date(msg.created_at), 'HH:mm')}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Input */}
            <form onSubmit={handleSend} className="p-4 border-t border-gray-800 flex gap-3">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder={t('chat.typeMessage')}
                className="flex-1 px-4 py-2.5 bg-gray-800/50 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow"
              />
              <button
                type="submit"
                disabled={!newMessage.trim()}
                className="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 disabled:hover:bg-gray-700 text-white rounded-lg font-medium text-sm transition-colors flex items-center gap-2"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5" />
                </svg>
                {t('chat.send')}
              </button>
            </form>
          </>
        ) : (
          <div className="flex-1 flex flex-col items-center justify-center text-gray-500">
            <svg className="w-16 h-16 text-gray-700 mb-4" fill="none" stroke="currentColor" strokeWidth="1" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M20.25 8.511c.884.284 1.5 1.128 1.5 2.097v4.286c0 1.136-.847 2.1-1.98 2.193-.34.027-.68.052-1.02.072v3.091l-3-3c-1.354 0-2.694-.055-4.02-.163a2.115 2.115 0 0 1-.825-.242m9.345-8.334a2.126 2.126 0 0 0-.476-.095 48.64 48.64 0 0 0-8.048 0c-1.131.094-1.976 1.057-1.976 2.192v4.286c0 .837.46 1.58 1.155 1.951m9.345-8.334V6.637c0-1.621-1.152-3.026-2.76-3.235A48.455 48.455 0 0 0 11.25 3c-2.115 0-4.198.137-6.24.402-1.608.209-2.76 1.614-2.76 3.235v6.226c0 1.621 1.152 3.026 2.76 3.235.577.075 1.157.14 1.74.194V21l4.155-4.155" />
            </svg>
            <p className="text-sm">{t('chat.startConversation')}</p>
          </div>
        )}
      </div>
    </div>
  );
}
