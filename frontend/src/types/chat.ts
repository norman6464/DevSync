import type { User } from './user';

export interface ChatRoom {
  id: number;
  name: string;
  description: string;
  owner_id: number;
  owner?: User;
  created_at: string;
  updated_at: string;
}

export interface ChatRoomMember {
  id: number;
  chat_room_id: number;
  user_id: number;
  user?: User;
  joined_at: string;
}

export interface GroupMessage {
  id: number;
  chat_room_id: number;
  sender_id: number;
  sender?: User;
  content: string;
  created_at: string;
}
