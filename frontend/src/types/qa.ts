import type { User } from './user';

export interface Question {
  id: number;
  user_id: number;
  user?: User;
  title: string;
  body: string;
  tags: string; // JSON array of tags
  vote_count: number;
  answer_count: number;
  is_solved: boolean;
  created_at: string;
  updated_at: string;
}

export interface Answer {
  id: number;
  user_id: number;
  user?: User;
  question_id: number;
  body: string;
  vote_count: number;
  is_best: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateQuestionRequest {
  title: string;
  body: string;
  tags?: string;
}

export interface UpdateQuestionRequest extends Partial<CreateQuestionRequest> {}

export interface CreateAnswerRequest {
  body: string;
}

export interface VoteRequest {
  value: 1 | -1;
}
