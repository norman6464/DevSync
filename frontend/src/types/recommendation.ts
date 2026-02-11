/** おすすめユーザー情報 */
export interface RecommendedUser {
  user: {
    id: number;
    name: string;
    avatar_url: string;
    github_username: string;
    skills_languages: string;
    skills_frameworks: string;
  };
  common_skills: string[];
  match_score: number;
}

/** トレンド投稿 */
export interface TrendingPost {
  id: number;
  user_id: number;
  user: {
    id: number;
    name: string;
    avatar_url: string;
  };
  title: string;
  content: string;
  like_count: number;
  comment_count: number;
  created_at: string;
}

/** トレンド学習リソース */
export interface TrendingResource {
  id: number;
  user_id: number;
  user: {
    id: number;
    name: string;
    avatar_url: string;
  };
  title: string;
  description: string;
  url: string;
  category: string;
  difficulty: string;
  like_count: number;
  save_count: number;
  created_at: string;
}
