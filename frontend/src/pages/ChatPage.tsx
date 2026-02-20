import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { MessageSquare, Users, Settings } from 'lucide-react';
import { useChat } from '../hooks/useChat';
import type { ChatRoom } from '../types/chat';
import Avatar from '../components/common/Avatar';
import ChatSidebar from '../components/chat/ChatSidebar';
import CreateRoomModal from '../components/chat/CreateRoomModal';
import RoomSettingsModal from '../components/chat/RoomSettingsModal';
import MessageBubble from '../components/chat/MessageBubble';
import MessageInputForm from '../components/chat/MessageInputForm';

export default function ChatPage() {
  const { t } = useTranslation();
  const c = useChat();

  const handleRoomCreated = useCallback((room: ChatRoom) => {
    c.setChatRooms([room, ...c.chatRooms]);
    c.setActiveRoomId(room.id);
    c.setShowCreateRoom(false);
  }, [c]);

  const handleRoomDeleted = useCallback(() => {
    c.setActiveRoomId(null);
    c.setShowRoomSettings(false);
    c.loadChatRooms();
  }, [c]);

  const handleRoomLeft = useCallback(() => {
    c.setActiveRoomId(null);
    c.setShowRoomSettings(false);
    c.loadChatRooms();
  }, [c]);

  return (
    <div className="flex h-[calc(100vh-7rem)] bg-gray-900 border border-gray-800 rounded-md overflow-hidden">
      {/* Sidebar */}
      <ChatSidebar
        activeTab={c.activeTab}
        onTabChange={c.setActiveTab}
        conversations={c.conversations}
        followingWithoutConversation={c.followingWithoutConversation}
        selectedUserId={c.selectedUserId}
        onSelectUser={c.handleSelectUser}
        chatRooms={c.chatRooms}
        activeRoomId={c.activeRoomId}
        onSelectRoom={c.handleSelectRoom}
        onCreateRoom={() => c.setShowCreateRoom(true)}
      />

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
              {c.groupMessages.map((msg) => (
                <MessageBubble
                  key={msg.id}
                  content={msg.content}
                  createdAt={msg.created_at}
                  isOwn={msg.sender_id === c.currentUser?.id}
                  senderName={msg.sender?.name}
                  senderAvatarUrl={msg.sender?.avatar_url}
                  showSenderInfo
                />
              ))}
            </div>

            <MessageInputForm
              value={c.newMessage}
              onChange={c.setNewMessage}
              onSubmit={c.handleSend}
              placeholder={t('chat.groupMessagePlaceholder')}
            />
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
              {c.activeMessages.map((msg) => (
                <MessageBubble
                  key={msg.id}
                  content={msg.content}
                  createdAt={msg.created_at}
                  isOwn={msg.sender_id === c.currentUser?.id}
                />
              ))}
            </div>

            <MessageInputForm
              value={c.newMessage}
              onChange={c.setNewMessage}
              onSubmit={c.handleSend}
              placeholder={t('chat.typeMessage')}
            />
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
          onCreated={handleRoomCreated}
        />
      )}

      {c.showRoomSettings && c.selectedRoom && (
        <RoomSettingsModal
          room={c.selectedRoom}
          currentUserId={c.currentUser?.id || 0}
          followingUsers={c.followingUsers}
          onClose={() => c.setShowRoomSettings(false)}
          onUpdated={c.loadChatRooms}
          onDeleted={handleRoomDeleted}
          onLeft={handleRoomLeft}
        />
      )}
    </div>
  );
}
