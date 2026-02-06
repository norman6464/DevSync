import client from './client';

export interface BookReview {
  id: number;
  user_id: number;
  user?: {
    id: number;
    name: string;
    avatar_url: string;
    github_username: string;
  };
  title: string;
  author: string;
  cover_url: string;
  rating: number;
  content: string;
  like_count: number;
  has_liked?: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateBookReviewRequest {
  title: string;
  author: string;
  cover_url?: string;
  rating: number;
  content: string;
}

export interface UpdateBookReviewRequest {
  title?: string;
  author?: string;
  cover_url?: string;
  rating?: number;
  content?: string;
}

export const createBookReview = (data: CreateBookReviewRequest) =>
  client.post<BookReview>('/book-reviews', data);

export const getBookReviews = (limit = 20, offset = 0) =>
  client.get<BookReview[]>(`/book-reviews?limit=${limit}&offset=${offset}`);

export const getBookReview = (id: number) =>
  client.get<BookReview>(`/book-reviews/${id}`);

export const getUserBookReviews = (userId: number) =>
  client.get<BookReview[]>(`/book-reviews/user/${userId}`);

export const updateBookReview = (id: number, data: UpdateBookReviewRequest) =>
  client.put<BookReview>(`/book-reviews/${id}`, data);

export const deleteBookReview = (id: number) =>
  client.delete(`/book-reviews/${id}`);

export const likeBookReview = (id: number) =>
  client.post(`/book-reviews/${id}/like`);

export const unlikeBookReview = (id: number) =>
  client.delete(`/book-reviews/${id}/like`);
