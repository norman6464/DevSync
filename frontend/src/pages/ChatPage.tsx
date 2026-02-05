import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { useChatStore } from '../store/chatStore';
import { getConversations, getMessages, sendMessage as sendMessageApi } from '../api/messages';
import type { Conversation, Message } from '../types/message';
import Avatar from '../components/common/Avatar';
import { format } from 'date-fns';

export default function ChatPage() {
  const { userId } = useParams<{ userId: string }>();
  const currentUser = useAuthStore((s) => s.user);
  const token = useAuthStore((s) => s.token);
  const { socket, connect, activeMessages, setActiveMessages } = useChatStore();
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<number | null>(
    userId ? parseInt(userId) : null
  );
  const [newMessage, setNewMessage] = useState('');

  useEffect(() => {
    if (token && !socket) {
      connect(token);
    }
  }, [token, socket, connect]);

  useEffect(() => {
    getConversations()
      .then(({ data }) => setConversations(data || []))
      .catch(() => {});
  }, []);

  useEffect(() => {
    if (selectedUserId) {
      getMessages(selectedUserId)
        .then(({ data }) => setActiveMessages(data || []))
        .catch(() => setActiveMessages([]));
    }
  }, [selectedUserId, setActiveMessages]);

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
    <div className="flex h-[calc(100vh-7rem)] bg-gray-900 rounded-lg overflow-hidden">
      <div className="w-80 border-r border-gray-800 overflow-y-auto">
        <div className="p-4 border-b border-gray-800">
          <h2 className="font-semibold">Messages</h2>
        </div>
        {conversations.map((conv) => (
          <button
            key={conv.user.id}
            onClick={() => setSelectedUserId(conv.user.id)}
            className={`w-full flex items-center gap-3 p-4 hover:bg-gray-800 transition-colors text-left ${
              selectedUserId === conv.user.id ? 'bg-gray-800' : ''
            }`}
          >
            <Avatar name={conv.user.name} avatarUrl={conv.user.avatar_url} size="sm" />
            <div className="flex-1 min-w-0">
              <div className="font-medium text-sm">{conv.user.name}</div>
              {conv.last_message && (
                <div className="text-xs text-gray-400 truncate">{conv.last_message.content}</div>
              )}
            </div>
            {conv.unread_count > 0 && (
              <span className="bg-blue-600 text-xs rounded-full px-2 py-0.5">
                {conv.unread_count}
              </span>
            )}
          </button>
        ))}
      </div>

      <div className="flex-1 flex flex-col">
        {selectedUserId ? (
          <>
            <div className="flex-1 overflow-y-auto p-4 space-y-3">
              {activeMessages.map((msg: Message) => (
                <div
                  key={msg.id}
                  className={`flex ${msg.sender_id === currentUser?.id ? 'justify-end' : 'justify-start'}`}
                >
                  <div
                    className={`max-w-xs px-4 py-2 rounded-lg ${
                      msg.sender_id === currentUser?.id
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-800 text-gray-100'
                    }`}
                  >
                    <p className="text-sm">{msg.content}</p>
                    <p className="text-xs opacity-60 mt-1">
                      {format(new Date(msg.created_at), 'HH:mm')}
                    </p>
                  </div>
                </div>
              ))}
            </div>
            <form onSubmit={handleSend} className="p-4 border-t border-gray-800 flex gap-2">
              <input
                type="text"
                value={newMessage}
                onChange={(e) => setNewMessage(e.target.value)}
                placeholder="Type a message..."
                className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:border-blue-500"
              />
              <button
                type="submit"
                className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors"
              >
                Send
              </button>
            </form>
          </>
        ) : (
          <div className="flex-1 flex items-center justify-center text-gray-400">
            Select a conversation
          </div>
        )}
      </div>
    </div>
  );
}
