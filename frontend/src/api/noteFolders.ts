import client from './client';
import type { NoteFolder, CreateNoteFolderRequest, UpdateNoteFolderRequest } from '../types/noteFolder';

export const createNoteFolder = (data: CreateNoteFolderRequest) =>
  client.post<NoteFolder>('/note-folders', data);

export const getNoteFolderById = (id: number) =>
  client.get<NoteFolder>(`/note-folders/${id}`);

export const getMyNoteFolders = () =>
  client.get<NoteFolder[]>('/note-folders');

export const getRootNoteFolders = () =>
  client.get<NoteFolder[]>('/note-folders/root');

export const getNoteFolderChildren = (id: number) =>
  client.get<NoteFolder[]>(`/note-folders/${id}/children`);

export const updateNoteFolder = (id: number, data: UpdateNoteFolderRequest) =>
  client.put<NoteFolder>(`/note-folders/${id}`, data);

export const deleteNoteFolder = (id: number) =>
  client.delete(`/note-folders/${id}`);
