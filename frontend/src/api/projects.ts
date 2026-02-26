import client from './client';
import type { Project, CreateProjectRequest, UpdateProjectRequest } from '../types/project';

export const createProject = async (data: CreateProjectRequest): Promise<Project> => {
  const res = await client.post('/projects', data);
  return res.data;
};

export const getProjects = async (limit = 20, offset = 0): Promise<{ projects: Project[]; total: number }> => {
  const res = await client.get('/projects', { params: { limit, offset } });
  return res.data;
};

export const getProjectById = async (id: number): Promise<Project> => {
  const res = await client.get(`/projects/${id}`);
  return res.data;
};

export const getProjectsByUserId = async (userId: number): Promise<Project[]> => {
  const res = await client.get(`/projects/user/${userId}`);
  return res.data;
};

export const getFeaturedProjects = async (userId: number): Promise<Project[]> => {
  const res = await client.get(`/projects/user/${userId}/featured`);
  return res.data;
};

export const updateProject = async (id: number, data: UpdateProjectRequest): Promise<Project> => {
  const res = await client.put(`/projects/${id}`, data);
  return res.data;
};

export const deleteProject = async (id: number): Promise<void> => {
  await client.delete(`/projects/${id}`);
};

// --- マイルストーン ---

export interface ProjectMilestone {
  id: number;
  project_id: number;
  title: string;
  description: string;
  status: 'not_started' | 'in_progress' | 'completed';
  due_date: string | null;
  completed_at: string | null;
  created_at: string;
  updated_at: string;
}

export interface CreateMilestoneRequest {
  title: string;
  description?: string;
  due_date?: string;
}

export interface UpdateMilestoneRequest {
  title?: string;
  description?: string;
  due_date?: string;
  status?: 'not_started' | 'in_progress' | 'completed';
}

export const createMilestone = async (projectId: number, data: CreateMilestoneRequest): Promise<void> => {
  await client.post(`/projects/${projectId}/milestones`, data);
};

export const getMilestones = async (projectId: number): Promise<ProjectMilestone[]> => {
  const res = await client.get(`/projects/${projectId}/milestones`);
  return res.data.milestones;
};

export const updateMilestone = async (milestoneId: number, data: UpdateMilestoneRequest): Promise<ProjectMilestone> => {
  const res = await client.put(`/projects/milestones/${milestoneId}`, data);
  return res.data;
};

export const deleteMilestone = async (milestoneId: number): Promise<void> => {
  await client.delete(`/projects/milestones/${milestoneId}`);
};
