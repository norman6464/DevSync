import client from './client';

export interface LearningLogTemplate {
  id: number;
  user_id: number;
  name: string;
  default_title: string;
  default_content: string;
  default_category: string;
  default_duration: number;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateLearningLogTemplateInput {
  name: string;
  default_title?: string;
  default_content?: string;
  default_category?: string;
  default_duration?: number;
  is_default?: boolean;
}

export interface UpdateLearningLogTemplateInput {
  name?: string;
  default_title?: string;
  default_content?: string;
  default_category?: string;
  default_duration?: number;
  is_default?: boolean;
}

export const getLearningLogTemplates = () =>
  client.get<LearningLogTemplate[]>('/learning-log-templates');

export const getLearningLogTemplate = (id: number) =>
  client.get<LearningLogTemplate>(`/learning-log-templates/${id}`);

export const getDefaultLearningLogTemplate = () =>
  client.get<LearningLogTemplate>('/learning-log-templates/default');

export const createLearningLogTemplate = (input: CreateLearningLogTemplateInput) =>
  client.post<LearningLogTemplate>('/learning-log-templates', input);

export const updateLearningLogTemplate = (id: number, input: UpdateLearningLogTemplateInput) =>
  client.put<LearningLogTemplate>(`/learning-log-templates/${id}`, input);

export const deleteLearningLogTemplate = (id: number) =>
  client.delete(`/learning-log-templates/${id}`);

export const useLearningLogTemplate = (id: number) =>
  client.post(`/learning-log-templates/${id}/use`);
