import client from './client';
import type { BookReview, CreateBookReviewRequest, UpdateBookReviewRequest } from '../types/bookReview';

export const createBookReview = async (data: CreateBookReviewRequest): Promise<BookReview> => {
  const res = await client.post('/book-reviews', data);
  return res.data;
};

export const getBookReviews = async (limit = 20, offset = 0): Promise<{ reviews: BookReview[]; total: number }> => {
  const res = await client.get('/book-reviews', { params: { limit, offset } });
  return res.data;
};

export const getBookReviewById = async (id: number): Promise<BookReview> => {
  const res = await client.get(`/book-reviews/${id}`);
  return res.data;
};

export const getBookReviewsByUserId = async (userId: number): Promise<BookReview[]> => {
  const res = await client.get(`/book-reviews/user/${userId}`);
  return res.data;
};

export const updateBookReview = async (id: number, data: UpdateBookReviewRequest): Promise<BookReview> => {
  const res = await client.put(`/book-reviews/${id}`, data);
  return res.data;
};

export const deleteBookReview = async (id: number): Promise<void> => {
  await client.delete(`/book-reviews/${id}`);
};
