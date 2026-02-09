import client from './client';

export interface AIAdvice {
  id: number;
  user_id: number;
  type: string;
  priority: number;
  title_key: string;
  message_key: string;
  params: string;
  action_url: string;
  is_read: boolean;
  expires_at: string | null;
  created_at: string;
}

export interface AIMessage {
  id: number;
  conversation_id: number;
  role: 'user' | 'assistant' | 'system';
  content: string;
  tokens_used: number;
  created_at: string;
}

export interface AIConversation {
  id: number;
  user_id: number;
  title: string;
  created_at: string;
  updated_at: string;
  messages?: AIMessage[];
}

export interface AdviceResponse {
  advices: AIAdvice[];
  llm_available: boolean;
  daily_chat_remaining: number;
}

export const getAdvice = () =>
  client.get<AdviceResponse>('/advice');

export const markAdviceRead = (id: number) =>
  client.put(`/advice/${id}/read`);

export const chatWithAI = (message: string, conversationId?: number) =>
  client.post<AIConversation>('/advice/chat', {
    message,
    conversation_id: conversationId || 0,
  });

export const getConversations = (limit = 20, offset = 0) =>
  client.get<AIConversation[]>('/advice/conversations', {
    params: { limit, offset },
  });

export const getConversation = (id: number) =>
  client.get<AIConversation>(`/advice/conversations/${id}`);

export const deleteConversation = (id: number) =>
  client.delete(`/advice/conversations/${id}`);
