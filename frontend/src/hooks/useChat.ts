import { useEffect, useMemo, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { useChatStore } from '../store/chatStore';
import { getConversations, getMessages, sendMessage as sendMessageApi } from '../api/messages';
import { getFollowing } from '../api/users';
import { getChatRooms, getChatRoomMessages, sendGroupMessage } from '../api/chatRooms';
import type { Conversation } from '../types/message';
import type { User } from '../types/user';

export function useChat() {
  const { userId } = useParams<{ userId: string }>();
  const currentUser = useAuthStore((s) => s.user);
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
  const [newMessage, setNewMessage] = useState('');
  const [showCreateRoom, setShowCreateRoom] = useState(false);
  const [showRoomSettings, setShowRoomSettings] = useState(false);

  useEffect(() => {
    if (!socket) {
      connect();
    }
  }, [socket, connect]);

  const loadChatRooms = () => {
    getChatRooms()
      .then(({ data }) => setChatRooms(data || []))
      .catch((e) => console.warn('Failed to load chat rooms:', e));
  };

  useEffect(() => {
    if (!currentUser) return;

    getConversations()
      .then(({ data }) => setConversations(data || []))
      .catch((e) => console.warn('Failed to load conversations:', e));

    getFollowing(currentUser.id)
      .then(({ data }) => setFollowingUsers(data || []))
      .catch((e) => console.warn('Failed to load following users:', e));

    loadChatRooms();
  }, [currentUser]);

  useEffect(() => {
    if (selectedUserId) {
      getMessages(selectedUserId)
        .then(({ data }) => setActiveMessages(data || []))
        .catch(() => setActiveMessages([]));

    }
  }, [selectedUserId, setActiveMessages]);

  // 選択中ユーザーは selectedUserId と一覧から一意に決まる派生値
  const selectedUser = useMemo(() => {
    if (!selectedUserId) return null;
    const convUser = conversations.find((c) => c.user.id === selectedUserId)?.user;
    const followUser = followingUsers.find((u) => u.id === selectedUserId);
    return convUser || followUser || null;
  }, [selectedUserId, conversations, followingUsers]);

  useEffect(() => {
    if (activeRoomId) {
      getChatRoomMessages(activeRoomId)
        .then(({ data }) => setGroupMessages(data || []))
        .catch(() => setGroupMessages([]));
    }
  }, [activeRoomId, setGroupMessages]);

  const followingWithoutConversation = followingUsers.filter(
    (user) => !conversations.some((conv) => conv.user && conv.user.id === user.id)
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

  const handleSelectRoom = (roomId: number) => {
    setActiveRoomId(roomId);
    setSelectedUserId(null);
  };

  const handleSelectUser = (id: number) => {
    setSelectedUserId(id);
    setActiveRoomId(null);
  };

  return {
    currentUser,
    // Tab
    activeTab, setActiveTab,
    // DM
    conversations, followingUsers, followingWithoutConversation,
    selectedUserId, selectedUser, handleSelectUser,
    activeMessages,
    // Group
    chatRooms, setChatRooms,
    activeRoomId, setActiveRoomId,
    groupMessages, selectedRoom, handleSelectRoom,
    // Message input
    newMessage, setNewMessage, handleSend,
    // Modals
    showCreateRoom, setShowCreateRoom,
    showRoomSettings, setShowRoomSettings,
    loadChatRooms,
  };
}
