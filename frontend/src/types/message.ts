import type { User } from './user';

export interface Message {
  id: number;
  sender_id: number;
  sender: User;
  receiver_id: number;
  receiver: User;
  content: string;
  read: boolean;
  created_at: string;
}

export interface Conversation {
  user: User;
  last_message: Message;
  unread_count: number;
}
