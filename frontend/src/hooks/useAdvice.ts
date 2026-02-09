import { useState, useCallback } from 'react';
import { useAsyncData } from './useAsyncData';
import {
  getAdvice,
  markAdviceRead,
  chatWithAI,
  getConversations,
  getConversation,
} from '../api/advice';
import type { AIAdvice, AIConversation, AIMessage } from '../api/advice';

export function useAdvice() {
  const { data, loading, refetch } = useAsyncData(
    async () => {
      const res = await getAdvice();
      return res.data;
    },
    {
      initialData: {
        advices: [] as AIAdvice[],
        llm_available: false,
        daily_chat_remaining: 0,
      },
    }
  );

  const markRead = useCallback(
    async (id: number) => {
      try {
        await markAdviceRead(id);
        await refetch();
      } catch {
        // エラーは無視
      }
    },
    [refetch]
  );

  return {
    advices: data.advices,
    llmAvailable: data.llm_available,
    dailyChatRemaining: data.daily_chat_remaining,
    loading,
    refetch,
    markRead,
  };
}

export function useAIChat() {
  const [messages, setMessages] = useState<AIMessage[]>([]);
  const [sending, setSending] = useState(false);
  const [conversationId, setConversationId] = useState<number | null>(null);

  const sendMessage = useCallback(
    async (message: string) => {
      if (!message.trim() || sending) return null;
      setSending(true);

      // ユーザーメッセージを即座に表示
      const userMsg: AIMessage = {
        id: Date.now(),
        conversation_id: conversationId || 0,
        role: 'user',
        content: message,
        tokens_used: 0,
        created_at: new Date().toISOString(),
      };
      setMessages((prev) => [...prev, userMsg]);

      try {
        const res = await chatWithAI(message, conversationId || undefined);
        const conv = res.data;
        setConversationId(conv.id);
        if (conv.messages) {
          setMessages(conv.messages.filter((m) => m.role !== 'system'));
        }
        return conv;
      } catch {
        // エラー時はユーザーメッセージを残す
        return null;
      } finally {
        setSending(false);
      }
    },
    [conversationId, sending]
  );

  const loadConversation = useCallback(async (id: number) => {
    try {
      const res = await getConversation(id);
      const conv = res.data;
      setConversationId(conv.id);
      if (conv.messages) {
        setMessages(conv.messages.filter((m) => m.role !== 'system'));
      }
      return conv;
    } catch {
      return null;
    }
  }, []);

  const startNewChat = useCallback(() => {
    setConversationId(null);
    setMessages([]);
  }, []);

  return {
    messages,
    sending,
    conversationId,
    sendMessage,
    loadConversation,
    startNewChat,
  };
}

export function useConversations() {
  const { data: conversations, loading, refetch } = useAsyncData(
    async () => {
      const res = await getConversations();
      return res.data;
    },
    { initialData: [] as AIConversation[] }
  );

  return { conversations, loading, refetch };
}
