import client from './client';

export interface Note {
  id: number;
  user_id: number;
  folder_id?: number | null;
  title: string;
  content: string;
  tags: string;
  is_favorite: boolean;
  is_archived: boolean;
  created_at: string;
  updated_at: string;
}

export interface NoteFolder {
  id: number;
  user_id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface CreateNoteRequest {
  title: string;
  content: string;
  tags?: string;
  folder_id?: number | null;
}

export interface UpdateNoteRequest {
  title?: string;
  content?: string;
  tags?: string;
  folder_id?: number | null;
}

export interface NotePaginatedResponse {
  data: Note[];
  total: number;
  page: number;
  limit: number;
}

export const createNote = (data: CreateNoteRequest) =>
  client.post<Note>('/notes', data);

export const updateNote = (id: number, data: UpdateNoteRequest) =>
  client.put<Note>(`/notes/${id}`, data);

export const deleteNote = (id: number) =>
  client.delete(`/notes/${id}`);

export const getNote = (id: number) =>
  client.get<Note>(`/notes/${id}`);

export const getMyNotes = (page = 1, limit = 20) =>
  client.get<NotePaginatedResponse>(`/notes?page=${page}&limit=${limit}`);

export const getNotesByFolder = (folderId: number) =>
  client.get<Note[]>(`/notes/folder/${folderId}`);

export const searchNotes = (query: string, page = 1, limit = 20) =>
  client.get<NotePaginatedResponse>(`/notes/search?q=${encodeURIComponent(query)}&page=${page}&limit=${limit}`);

export const toggleFavorite = (id: number) =>
  client.put<{ message: string }>(`/notes/${id}/favorite`);

export const archiveNote = (id: number) =>
  client.put<{ message: string }>(`/notes/${id}/archive`);

export const unarchiveNote = (id: number) =>
  client.put<{ message: string }>(`/notes/${id}/unarchive`);

export const getArchivedNotes = (page = 1, limit = 20) =>
  client.get<NotePaginatedResponse>(`/notes/archived?page=${page}&limit=${limit}`);

export const duplicateNote = (id: number) =>
  client.post<Note>(`/notes/${id}/duplicate`);

export const exportNoteMarkdown = (id: number) =>
  client.get<Blob>(`/notes/${id}/export`, { responseType: 'blob' });
