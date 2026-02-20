import type { User } from './user';

export type ReviewStatus = 'not_started' | 'reading' | 'completed';

export interface BookReview {
  id: number;
  user_id: number;
  user?: User;
  title: string;
  author: string;
  isbn: string;
  rating: number; // 1-5
  review: string;
  image_url: string;
  status: ReviewStatus;
  is_archived: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateBookReviewRequest {
  title: string;
  author?: string;
  isbn?: string;
  rating: number;
  review?: string;
  image_url?: string;
}

export interface UpdateBookReviewRequest extends Partial<CreateBookReviewRequest> {}
