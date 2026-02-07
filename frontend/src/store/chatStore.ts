import { create } from 'zustand';
import type { Message, Conversation } from '../types/message';
import type { ChatRoom, GroupMessage } from '../types/chat';

interface WSGroupMessage {
  type: 'group_message';
  sender_id: number;
  room_id: number;
  content: string;
  sender_name: string;
}

interface ChatState {
  socket: WebSocket | null;
  conversations: Conversation[];
  activeMessages: Message[];
  connected: boolean;
  activeTab: 'dm' | 'group';
  chatRooms: ChatRoom[];
  activeRoomId: number | null;
  groupMessages: GroupMessage[];
  connect: (token: string) => void;
  disconnect: () => void;
  setConversations: (conversations: Conversation[]) => void;
  setActiveMessages: (messages: Message[]) => void;
  addMessage: (message: Message) => void;
  setActiveTab: (tab: 'dm' | 'group') => void;
  setChatRooms: (rooms: ChatRoom[]) => void;
  setActiveRoomId: (id: number | null) => void;
  setGroupMessages: (messages: GroupMessage[]) => void;
  addGroupMessage: (message: GroupMessage) => void;
}

export const useChatStore = create<ChatState>((set, get) => ({
  socket: null,
  conversations: [],
  activeMessages: [],
  connected: false,
  activeTab: 'dm',
  chatRooms: [],
  activeRoomId: null,
  groupMessages: [],

  connect: (token) => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${window.location.host}/ws?token=${token}`);

    ws.onopen = () => set({ connected: true });
    ws.onclose = () => {
      set({ connected: false, socket: null });
      setTimeout(() => {
        const state = get();
        if (!state.connected && token) {
          state.connect(token);
        }
      }, 3000);
    };
    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      if (data.type === 'group_message') {
        const wsMsg = data as WSGroupMessage;
        const state = get();
        if (state.activeRoomId === wsMsg.room_id) {
          const groupMsg: GroupMessage = {
            id: Date.now(),
            chat_room_id: wsMsg.room_id,
            sender_id: wsMsg.sender_id,
            sender: { id: wsMsg.sender_id, name: wsMsg.sender_name } as GroupMessage['sender'],
            content: wsMsg.content,
            created_at: new Date().toISOString(),
          };
          set((s) => ({ groupMessages: [...s.groupMessages, groupMsg] }));
        }
      } else {
        const message = data as Message;
        set((state) => ({
          activeMessages: [...state.activeMessages, message],
        }));
      }
    };

    set({ socket: ws });
  },

  disconnect: () => {
    const { socket } = get();
    if (socket) {
      socket.close();
    }
    set({ socket: null, connected: false });
  },

  setConversations: (conversations) => set({ conversations }),
  setActiveMessages: (messages) => set({ activeMessages: messages }),
  addMessage: (message) =>
    set((state) => ({ activeMessages: [...state.activeMessages, message] })),
  setActiveTab: (tab) => set({ activeTab: tab }),
  setChatRooms: (rooms) => set({ chatRooms: rooms }),
  setActiveRoomId: (id) => set({ activeRoomId: id }),
  setGroupMessages: (messages) => set({ groupMessages: messages }),
  addGroupMessage: (message) =>
    set((state) => ({ groupMessages: [...state.groupMessages, message] })),
}));
