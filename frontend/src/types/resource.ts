import type { User } from './user';

export type ResourceCategory = 'book' | 'video' | 'article' | 'course' | 'tutorial' | 'podcast' | 'tool' | 'other';
export type ResourceDifficulty = 'beginner' | 'intermediate' | 'advanced';

export interface LearningResource {
  id: number;
  user_id: number;
  user?: User;
  title: string;
  description: string;
  url: string;
  category: ResourceCategory;
  difficulty: ResourceDifficulty;
  tags: string; // JSON array of tags
  image_url: string;
  is_public: boolean;
  like_count: number;
  save_count: number;
  created_at: string;
  updated_at: string;
}

export interface CreateResourceRequest {
  title: string;
  description?: string;
  url?: string;
  category: ResourceCategory;
  difficulty?: ResourceDifficulty;
  tags?: string;
  image_url?: string;
  is_public?: boolean;
}

export interface UpdateResourceRequest extends Partial<CreateResourceRequest> {}

export interface ResourceWithStatus extends LearningResource {
  has_liked?: boolean;
  has_saved?: boolean;
}
