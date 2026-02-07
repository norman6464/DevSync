import { getUser, getFollowers, getFollowing } from '../api/users';
import { getUserPosts } from '../api/posts';
import { getContributions, getLanguages, getRepos } from '../api/github';
import { getZennArticles, getZennStats, type ZennArticle, type ZennStats } from '../api/zenn';
import { getQiitaArticles, getQiitaStats, type QiitaArticle, type QiitaStats } from '../api/qiita';
import { getUserGoals, getGoalStats, type LearningGoal, type LearningGoalStats } from '../api/goals';
import type { User } from '../types/user';
import type { Post } from '../types/post';
import type { GitHubContribution, GitHubLanguageStat, GitHubRepository } from '../types/github';
import { useAsyncData } from './useAsyncData';

interface ProfileData {
  user: User;
  posts: Post[];
  contributions: GitHubContribution[];
  languages: GitHubLanguageStat[];
  repos: GitHubRepository[];
  zennArticles: ZennArticle[];
  zennStats: ZennStats | null;
  qiitaArticles: QiitaArticle[];
  qiitaStats: QiitaStats | null;
  goals: LearningGoal[];
  goalStats: LearningGoalStats | null;
  followerCount: number;
  followingCount: number;
}

export function useProfile(id: string | undefined) {
  const userId = id ? parseInt(id) : 0;

  const { data, loading, refetch } = useAsyncData(
    async (): Promise<ProfileData> => {
      const [userRes, postsRes, followersRes, followingRes] = await Promise.all([
        getUser(userId),
        getUserPosts(userId),
        getFollowers(userId),
        getFollowing(userId),
      ]);

      const userData = userRes.data;
      let contributions: GitHubContribution[] = [];
      let languages: GitHubLanguageStat[] = [];
      let repos: GitHubRepository[] = [];

      if (userData.github_connected) {
        const [contribRes, langRes, reposRes] = await Promise.all([
          getContributions(userId),
          getLanguages(userId),
          getRepos(userId),
        ]);
        contributions = contribRes.data || [];
        languages = langRes.data || [];
        repos = reposRes.data || [];
      }

      let zennArticles: ZennArticle[] = [];
      let zennStats: ZennStats | null = null;
      if (userData.zenn_username) {
        const [articlesRes, statsRes] = await Promise.all([
          getZennArticles(userId),
          getZennStats(userId),
        ]);
        zennArticles = articlesRes.data || [];
        zennStats = statsRes.data;
      }

      let qiitaArticles: QiitaArticle[] = [];
      let qiitaStats: QiitaStats | null = null;
      if (userData.qiita_username) {
        const [articlesRes, statsRes] = await Promise.all([
          getQiitaArticles(userId),
          getQiitaStats(userId),
        ]);
        qiitaArticles = articlesRes.data || [];
        qiitaStats = statsRes.data;
      }

      const [goalsRes, goalStatsRes] = await Promise.all([
        getUserGoals(userId),
        getGoalStats(userId),
      ]);

      return {
        user: userData,
        posts: postsRes.data || [],
        contributions,
        languages,
        repos,
        zennArticles,
        zennStats,
        qiitaArticles,
        qiitaStats,
        goals: goalsRes.data || [],
        goalStats: goalStatsRes.data,
        followerCount: (followersRes.data || []).length,
        followingCount: (followingRes.data || []).length,
      };
    },
    { deps: [userId], enabled: !!userId }
  );

  return {
    user: data?.user ?? null,
    posts: data?.posts ?? [],
    contributions: data?.contributions ?? [],
    languages: data?.languages ?? [],
    repos: data?.repos ?? [],
    zennArticles: data?.zennArticles ?? [],
    zennStats: data?.zennStats ?? null,
    qiitaArticles: data?.qiitaArticles ?? [],
    qiitaStats: data?.qiitaStats ?? null,
    goals: data?.goals ?? [],
    goalStats: data?.goalStats ?? null,
    followerCount: data?.followerCount ?? 0,
    followingCount: data?.followingCount ?? 0,
    loading,
    refetch,
  };
}
