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
