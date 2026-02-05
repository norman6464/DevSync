import { create } from 'zustand';
import type { Message, Conversation } from '../types/message';

interface ChatState {
  socket: WebSocket | null;
  conversations: Conversation[];
  activeMessages: Message[];
  connected: boolean;
  connect: (token: string) => void;
  disconnect: () => void;
  setConversations: (conversations: Conversation[]) => void;
  setActiveMessages: (messages: Message[]) => void;
  addMessage: (message: Message) => void;
}

export const useChatStore = create<ChatState>((set, get) => ({
  socket: null,
  conversations: [],
  activeMessages: [],
  connected: false,

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
      const message = JSON.parse(event.data) as Message;
      set((state) => ({
        activeMessages: [...state.activeMessages, message],
      }));
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
}));
