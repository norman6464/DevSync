import client from './client';
import type {
  Question, Answer,
  CreateQuestionRequest, UpdateQuestionRequest,
  CreateAnswerRequest, VoteRequest,
} from '../types/qa';

// Questions
export const createQuestion = async (data: CreateQuestionRequest): Promise<Question> => {
  const res = await client.post('/questions', data);
  return res.data;
};

export const getQuestions = async (
  limit = 20, offset = 0, tag = '', sort = 'newest'
): Promise<{ questions: Question[]; total: number }> => {
  const res = await client.get('/questions', { params: { limit, offset, tag, sort } });
  return res.data;
};

export const searchQuestions = async (
  q: string, limit = 20, offset = 0
): Promise<{ questions: Question[]; total: number }> => {
  const res = await client.get('/questions/search', { params: { q, limit, offset } });
  return res.data;
};

export const getQuestionById = async (id: number): Promise<{ question: Question; user_vote: number }> => {
  const res = await client.get(`/questions/${id}`);
  return res.data;
};

export const getQuestionsByUserId = async (userId: number): Promise<Question[]> => {
  const res = await client.get(`/questions/user/${userId}`);
  return res.data;
};

export const updateQuestion = async (id: number, data: UpdateQuestionRequest): Promise<Question> => {
  const res = await client.put(`/questions/${id}`, data);
  return res.data;
};

export const deleteQuestion = async (id: number): Promise<void> => {
  await client.delete(`/questions/${id}`);
};

export const voteQuestion = async (id: number, data: VoteRequest): Promise<void> => {
  await client.post(`/questions/${id}/vote`, data);
};

export const removeQuestionVote = async (id: number): Promise<void> => {
  await client.delete(`/questions/${id}/vote`);
};

// Answers
export const createAnswer = async (questionId: number, data: CreateAnswerRequest): Promise<Answer> => {
  const res = await client.post(`/questions/${questionId}/answers`, data);
  return res.data;
};

export const updateAnswer = async (questionId: number, answerId: number, data: CreateAnswerRequest): Promise<Answer> => {
  const res = await client.put(`/questions/${questionId}/answers/${answerId}`, data);
  return res.data;
};

export const deleteAnswer = async (questionId: number, answerId: number): Promise<void> => {
  await client.delete(`/questions/${questionId}/answers/${answerId}`);
};

export const setBestAnswer = async (questionId: number, answerId: number): Promise<void> => {
  await client.put(`/questions/${questionId}/answers/${answerId}/best`);
};

export const voteAnswer = async (questionId: number, answerId: number, data: VoteRequest): Promise<void> => {
  await client.post(`/questions/${questionId}/answers/${answerId}/vote`, data);
};

export const removeAnswerVote = async (questionId: number, answerId: number): Promise<void> => {
  await client.delete(`/questions/${questionId}/answers/${answerId}/vote`);
};
