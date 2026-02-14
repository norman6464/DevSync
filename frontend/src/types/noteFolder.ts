export interface NoteFolder {
  id: number;
  user_id: number;
  parent_id?: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface CreateNoteFolderRequest {
  name: string;
  parent_id?: number;
}

export interface UpdateNoteFolderRequest {
  name?: string;
  parent_id?: number;
}
