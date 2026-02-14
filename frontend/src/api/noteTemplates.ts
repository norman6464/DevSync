import client from './client';

export interface NoteTemplate {
  id: number;
  user_id: number;
  name: string;
  description: string;
  default_title: string;
  content_template: string;
  default_tags: string;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateNoteTemplateRequest {
  name: string;
  description?: string;
  default_title?: string;
  content_template: string;
  default_tags?: string;
  is_default?: boolean;
}

export interface UpdateNoteTemplateRequest {
  name?: string;
  description?: string;
  default_title?: string;
  content_template?: string;
  default_tags?: string;
  is_default?: boolean;
}

export const createNoteTemplate = (data: CreateNoteTemplateRequest) =>
  client.post<NoteTemplate>('/note-templates', data);

export const updateNoteTemplate = (id: number, data: UpdateNoteTemplateRequest) =>
  client.put<NoteTemplate>(`/note-templates/${id}`, data);

export const deleteNoteTemplate = (id: number) =>
  client.delete(`/note-templates/${id}`);

export const getNoteTemplate = (id: number) =>
  client.get<NoteTemplate>(`/note-templates/${id}`);

export const getMyNoteTemplates = () =>
  client.get<NoteTemplate[]>('/note-templates');

export const getDefaultNoteTemplate = () =>
  client.get<NoteTemplate>('/note-templates/default');

export const useNoteTemplate = (id: number) =>
  client.post<{ id: number; title: string; content: string; tags: string }>(`/note-templates/${id}/use`);
