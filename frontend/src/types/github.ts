export interface GitHubContribution {
  date: string;
  count: number;
}

export interface GitHubLanguageStat {
  language: string;
  bytes: number;
  repo_count: number;
}

export interface GitHubRepository {
  id: number;
  github_repo_id: number;
  name: string;
  full_name: string;
  description: string;
  language: string;
  stars: number;
  forks: number;
  is_private: boolean;
}
