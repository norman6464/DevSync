export interface User {
  id: number;
  name: string;
  email: string;
  avatar_url: string;
  bio: string;
  github_id: number;
  github_username: string;
  github_connected: boolean;
  zenn_username: string;
  qiita_username: string;
  skills_languages: string;
  skills_frameworks: string;
  onboarding_completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}
