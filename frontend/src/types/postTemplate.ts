export interface PostTemplate {
  id: number;
  user_id: number;
  name: string;
  title_template: string;
  content_template: string;
  created_at: string;
  updated_at: string;
}

export interface CreatePostTemplateRequest {
  name: string;
  title_template?: string;
  content_template: string;
}

export interface UpdatePostTemplateRequest {
  name?: string;
  title_template?: string;
  content_template?: string;
}

export interface PostTemplateListResponse {
  templates: PostTemplate[];
  total: number;
  limit: number;
  offset: number;
}
