import client from './client';
import type { Note } from './notes';

export interface NoteVersion {
  id: number;
  note_id: number;
  version_number: number;
  title: string;
  content: string;
  tags: string;
  created_at: string;
}

export interface NoteVersionListResponse {
  versions: NoteVersion[];
  total: number;
  limit: number;
  offset: number;
}

export const getNoteVersions = (noteId: number, limit = 20, offset = 0) =>
  client.get<NoteVersionListResponse>(`/notes/${noteId}/versions`, { params: { limit, offset } });

export const getNoteVersion = (noteId: number, versionId: number) =>
  client.get<NoteVersion>(`/notes/${noteId}/versions/${versionId}`);

export const restoreNoteVersion = (noteId: number, versionId: number) =>
  client.post<Note>(`/notes/${noteId}/versions/${versionId}/restore`);
