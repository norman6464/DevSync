import client from './client';
import type { LearningResource, CreateResourceRequest, UpdateResourceRequest, ResourceWithStatus } from '../types/resource';

export const createResource = async (data: CreateResourceRequest): Promise<LearningResource> => {
  const res = await client.post('/resources', data);
  return res.data;
};

export const getPublicResources = async (
  limit = 20,
  offset = 0,
  category?: string,
  difficulty?: string
): Promise<{ resources: LearningResource[]; total: number }> => {
  const res = await client.get('/resources', { params: { limit, offset, category, difficulty } });
  return res.data;
};

export const searchResources = async (
  query: string,
  limit = 20,
  offset = 0
): Promise<{ resources: LearningResource[]; total: number }> => {
  const res = await client.get('/resources/search', { params: { q: query, limit, offset } });
  return res.data;
};

export const getResourceById = async (id: number): Promise<ResourceWithStatus> => {
  const res = await client.get(`/resources/${id}`);
  return {
    ...res.data.resource,
    has_liked: res.data.has_liked,
    has_saved: res.data.has_saved,
  };
};

export const getResourcesByUserId = async (userId: number): Promise<LearningResource[]> => {
  const res = await client.get(`/resources/user/${userId}`);
  return res.data;
};

export const getSavedResources = async (
  limit = 20,
  offset = 0
): Promise<{ resources: LearningResource[]; total: number }> => {
  const res = await client.get('/resources/saved', { params: { limit, offset } });
  return res.data;
};

export const updateResource = async (id: number, data: UpdateResourceRequest): Promise<LearningResource> => {
  const res = await client.put(`/resources/${id}`, data);
  return res.data;
};

export const deleteResource = async (id: number): Promise<void> => {
  await client.delete(`/resources/${id}`);
};

export const likeResource = async (id: number): Promise<void> => {
  await client.post(`/resources/${id}/like`);
};

export const unlikeResource = async (id: number): Promise<void> => {
  await client.delete(`/resources/${id}/like`);
};

export const saveResource = async (id: number): Promise<void> => {
  await client.post(`/resources/${id}/save`);
};

export const unsaveResource = async (id: number): Promise<void> => {
  await client.delete(`/resources/${id}/save`);
};

// --- リソース進捗 ---

export interface UpsertResourceProgressRequest {
  resource_id: number;
  status: 'not_started' | 'in_progress' | 'completed';
  completion_percent: number;
  note?: string;
}

export interface ResourceProgress {
  id: number;
  user_id: number;
  resource_id: number;
  resource?: LearningResource;
  status: 'not_started' | 'in_progress' | 'completed';
  completion_percent: number;
  note: string;
  started_at: string | null;
  completed_at: string | null;
  created_at: string;
  updated_at: string;
}

export const upsertResourceProgress = async (data: UpsertResourceProgressRequest): Promise<ResourceProgress> => {
  const res = await client.put('/resources/progress', data);
  return res.data.progress;
};

export const getMyResourceProgress = async (
  status?: string,
  limit = 20,
  offset = 0
): Promise<{ progresses: ResourceProgress[]; total: number }> => {
  const res = await client.get('/resources/progress', { params: { status, limit, offset } });
  return res.data;
};

export const getResourceProgress = async (resourceId: number): Promise<ResourceProgress> => {
  const res = await client.get(`/resources/${resourceId}/progress`);
  return res.data.progress;
};
