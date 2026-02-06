import type { User } from './user';
import type { GitHubRepository } from './github';

export interface Project {
  id: number;
  user_id: number;
  user?: User;
  title: string;
  description: string;
  tech_stack: string; // JSON array of technologies
  demo_url: string;
  github_url: string;
  image_url: string;
  role: string;
  start_date: string | null;
  end_date: string | null;
  featured: boolean;
  github_repo_id: number | null;
  github_repo?: GitHubRepository;
  created_at: string;
  updated_at: string;
}

export interface CreateProjectRequest {
  title: string;
  description?: string;
  tech_stack?: string;
  demo_url?: string;
  github_url?: string;
  image_url?: string;
  role?: string;
  start_date?: string;
  end_date?: string;
  featured?: boolean;
  github_repo_id?: number;
}

export interface UpdateProjectRequest extends Partial<CreateProjectRequest> {}
