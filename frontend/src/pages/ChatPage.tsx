import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { MessageSquare, Users, Plus, Settings, Send } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { useChatStore } from '../store/chatStore';
import { getConversations, getMessages, sendMessage as sendMessageApi } from '../api/messages';
import { getFollowing } from '../api/users';
import { getChatRooms, getChatRoomMessages, sendGroupMessage } from '../api/chatRooms';
import type { Conversation, Message } from '../types/message';
import type { User } from '../types/user';
import type { ChatRoom } from '../types/chat';
import Avatar from '../components/common/Avatar';
import CreateRoomModal from '../components/chat/CreateRoomModal';
import RoomSettingsModal from '../components/chat/RoomSettingsModal';
import { format } from 'date-fns';

export default function ChatPage() {
  const { t } = useTranslation();
  const { userId } = useParams<{ userId: string }>();
  const currentUser = useAuthStore((s) => s.user);
  const token = useAuthStore((s) => s.token);
  const {
    socket, connect, activeMessages, setActiveMessages,
    activeTab, setActiveTab,
    chatRooms, setChatRooms,
    activeRoomId, setActiveRoomId,
    groupMessages, setGroupMessages, addGroupMessage,
  } = useChatStore();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [followingUsers, setFollowingUsers] = useState<User[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<number | null>(
    userId ? parseInt(userId) : null
  );
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [newMessage, setNewMessage] = useState('');
  const [showCreateRoom, setShowCreateRoom] = useState(false);
  const [showRoomSettings, setShowRoomSettings] = useState(false);

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

    // Load chat rooms
    loadChatRooms();
  }, [currentUser]);

  const loadChatRooms = () => {
    getChatRooms()
      .then(({ data }) => setChatRooms(data || []))
      .catch(() => {});
  };

  useEffect(() => {
    if (selectedUserId) {
      getMessages(selectedUserId)
        .then(({ data }) => setActiveMessages(data || []))
        .catch(() => setActiveMessages([]));

      const convUser = conversations.find((c) => c.user.id === selectedUserId)?.user;
      const followUser = followingUsers.find((u) => u.id === selectedUserId);
      setSelectedUser(convUser || followUser || null);
    }
  }, [selectedUserId, setActiveMessages, conversations, followingUsers]);

  useEffect(() => {
    if (activeRoomId) {
      getChatRoomMessages(activeRoomId)
        .then(({ data }) => setGroupMessages(data || []))
        .catch(() => setGroupMessages([]));
    }
  }, [activeRoomId, setGroupMessages]);

  const followingWithoutConversation = followingUsers.filter(
    (user) => !conversations.some((conv) => conv.user.id === user.id)
  );

  const selectedRoom = chatRooms.find((r) => r.id === activeRoomId);

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim()) return;

    if (activeTab === 'group' && activeRoomId) {
      try {
        const { data } = await sendGroupMessage(activeRoomId, newMessage);
        addGroupMessage(data);
        if (socket && socket.readyState === WebSocket.OPEN) {
          socket.send(JSON.stringify({
            type: 'group_message',
            room_id: activeRoomId,
            content: newMessage,
          }));
        }
        setNewMessage('');
      } catch {
        // handle error
      }
    } else if (selectedUserId) {
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
    }
  };

  const handleSelectRoom = (room: ChatRoom) => {
    setActiveRoomId(room.id);
    setSelectedUserId(null);
  };

  const handleSelectUser = (id: number) => {
    setSelectedUserId(id);
    setActiveRoomId(null);
  };

  return (
    <div className="flex h-[calc(100vh-7rem)] bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
      {/* Sidebar */}
      <div className="w-80 border-r border-gray-800 flex flex-col">
        {/* Tab Switcher */}
        <div className="flex border-b border-gray-800">
          <button
            onClick={() => setActiveTab('dm')}
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
            onClick={() => setActiveTab('group')}
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
              {/* Existing Conversations */}
              {conversations.length > 0 && (
                <div>
                  <div className="px-5 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider">
                    {t('chat.recentChats')}
                  </div>
                  {conversations.map((conv) => (
                    <button
                      key={conv.user.id}
                      onClick={() => handleSelectUser(conv.user.id)}
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

              {/* Following Users without Conversation */}
              {followingWithoutConversation.length > 0 && (
                <div>
                  <div className="px-5 py-2 text-xs font-medium text-gray-500 uppercase tracking-wider border-t border-gray-800 mt-2 pt-3">
                    {t('chat.following')}
                  </div>
                  {followingWithoutConversation.map((user) => (
                    <button
                      key={user.id}
                      onClick={() => handleSelectUser(user.id)}
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
              {/* Create Group Button */}
              <div className="p-3">
                <button
                  onClick={() => setShowCreateRoom(true)}
                  className="w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-gray-800 hover:bg-gray-700 text-white rounded-lg text-sm font-medium transition-colors"
                >
                  <Plus className="w-4 h-4" />
                  {t('chat.createGroup')}
                </button>
              </div>

              {/* Group List */}
              {chatRooms.length > 0 ? (
                chatRooms.map((room) => (
                  <button
                    key={room.id}
                    onClick={() => handleSelectRoom(room)}
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

      {/* Chat Area */}
      <div className="flex-1 flex flex-col">
        {activeTab === 'group' && activeRoomId && selectedRoom ? (
          <>
            {/* Group Chat Header */}
            <div className="px-6 py-3 border-b border-gray-800 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-gray-700 rounded-full flex items-center justify-center">
                  <Users className="w-4 h-4 text-gray-300" />
                </div>
                <div>
                  <div className="font-medium text-sm">{selectedRoom.name}</div>
                  {selectedRoom.description && (
                    <div className="text-xs text-gray-500">{selectedRoom.description}</div>
                  )}
                </div>
              </div>
              <button
                onClick={() => setShowRoomSettings(true)}
                className="p-2 text-gray-400 hover:text-white transition-colors rounded-md"
              >
                <Settings className="w-5 h-5" />
              </button>
            </div>

            {/* Group Messages */}
            <div className="flex-1 overflow-y-auto p-6 space-y-4">
              {groupMessages.map((msg) => {
                const isOwn = msg.sender_id === currentUser?.id;
                return (
                  <div
                    key={msg.id}
                    className={`flex ${isOwn ? 'justify-end' : 'justify-start'}`}
                  >
                    <div className={`max-w-sm ${isOwn ? '' : 'flex items-start gap-2'}`}>
                      {!isOwn && (
                        <Avatar
                          name={msg.sender?.name || ''}
                          avatarUrl={msg.sender?.avatar_url}
                          size="xs"
                        />
                      )}
                      <div>
                        {!isOwn && (
                          <p className="text-xs text-gray-400 mb-1">{msg.sender?.name}</p>
                        )}
                        <div
                          className={`px-4 py-2.5 rounded-2xl ${
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
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Group Input */}
            <form onSubmit={handleSend} className="p-4 border-t border-gray-800 flex gap-3">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder={t('chat.groupMessagePlaceholder')}
                className="flex-1 px-4 py-2.5 bg-gray-800/50 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow"
              />
              <button
                type="submit"
                disabled={!newMessage.trim()}
                className="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 disabled:hover:bg-gray-700 text-white rounded-lg font-medium text-sm transition-colors flex items-center gap-2"
              >
                <Send className="w-4 h-4" />
                {t('chat.send')}
              </button>
            </form>
          </>
        ) : selectedUserId ? (
          <>
            {/* DM Chat Header */}
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

            {/* DM Messages */}
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

            {/* DM Input */}
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
      {showCreateRoom && (
        <CreateRoomModal
          followingUsers={followingUsers}
          onClose={() => setShowCreateRoom(false)}
          onCreated={(room) => {
            setChatRooms([room, ...chatRooms]);
            setActiveRoomId(room.id);
            setShowCreateRoom(false);
          }}
        />
      )}

      {showRoomSettings && selectedRoom && (
        <RoomSettingsModal
          room={selectedRoom}
          currentUserId={currentUser?.id || 0}
          followingUsers={followingUsers}
          onClose={() => setShowRoomSettings(false)}
          onUpdated={loadChatRooms}
          onDeleted={() => {
            setActiveRoomId(null);
            setShowRoomSettings(false);
            loadChatRooms();
          }}
          onLeft={() => {
            setActiveRoomId(null);
            setShowRoomSettings(false);
            loadChatRooms();
          }}
        />
      )}
    </div>
  );
}
