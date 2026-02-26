import client from './client';
import type { Note } from './notes';

export interface NoteLink {
  id: number;
  source_note_id: number;
  target_note_id: number;
  source_note?: Note;
  target_note?: Note;
  created_at: string;
}

export interface CreateNoteLinkRequest {
  target_note_id: number;
}

export const createNoteLink = (sourceNoteId: number, data: CreateNoteLinkRequest) =>
  client.post<{ message: string }>(`/notes/${sourceNoteId}/links`, data);

export const getNoteLinks = (noteId: number) =>
  client.get<NoteLink[]>(`/notes/${noteId}/links`);

export const getNoteBacklinks = (noteId: number) =>
  client.get<NoteLink[]>(`/notes/${noteId}/backlinks`);

export const deleteNoteLink = (sourceNoteId: number, targetNoteId: number) =>
  client.delete(`/notes/${sourceNoteId}/links/${targetNoteId}`);

export const getNoteLinkStats = (noteId: number) =>
  client.get(`/notes/${noteId}/link-stats`);
