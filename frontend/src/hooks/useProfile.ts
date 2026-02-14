import { getUser, getUserByUsername, getFollowers, getFollowing } from '../api/users';
import { getUserPosts } from '../api/posts';
import { getContributions, getLanguages, getRepos } from '../api/github';
import { getZennArticles, getZennStats, type ZennArticle, type ZennStats } from '../api/zenn';
import { getQiitaArticles, getQiitaStats, type QiitaArticle, type QiitaStats } from '../api/qiita';
import { getAtCoderRating, type AtCoderRatingInfo } from '../api/atcoder';
import { getUserGoals, getGoalStats, type LearningGoal, type LearningGoalStats } from '../api/goals';
import { getUserBadges } from '../api/badges';
import { getStreakInfo } from '../api/learningLogs';
import type { User } from '../types/user';
import type { Post } from '../types/post';
import type { GitHubContribution, GitHubLanguageStat, GitHubRepository } from '../types/github';
import type { BadgeResult } from '../types/badge';
import type { StreakInfo } from '../types/learningLog';
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
  badges: BadgeResult[];
  atcoderRating: AtCoderRatingInfo | null;
  streakInfo: StreakInfo | null;
}

export function useProfile(usernameOrId: string | undefined) {
  // usernameOrIdが数値の場合はID、そうでない場合はusername
  const isId = usernameOrId && /^\d+$/.test(usernameOrId);
  const userId = isId ? parseInt(usernameOrId) : 0;
  const username = !isId ? usernameOrId : '';

  const { data, loading, refetch } = useAsyncData(
    async (): Promise<ProfileData> => {
      // usernameまたはIDでユーザーを取得
      const userRes = username
        ? await getUserByUsername(username)
        : await getUser(userId);

      const userData = userRes.data;
      const actualUserId = userData.id;

      const [postsRes, followersRes, followingRes] = await Promise.all([
        getUserPosts(actualUserId),
        getFollowers(actualUserId),
        getFollowing(actualUserId),
      ]);
      let contributions: GitHubContribution[] = [];
      let languages: GitHubLanguageStat[] = [];
      let repos: GitHubRepository[] = [];

      if (userData.github_connected) {
        const [contribRes, langRes, reposRes] = await Promise.all([
          getContributions(actualUserId),
          getLanguages(actualUserId),
          getRepos(actualUserId),
        ]);
        contributions = contribRes.data || [];
        languages = langRes.data || [];
        repos = reposRes.data || [];
      }

      let zennArticles: ZennArticle[] = [];
      let zennStats: ZennStats | null = null;
      if (userData.zenn_username) {
        const [articlesRes, statsRes] = await Promise.all([
          getZennArticles(actualUserId),
          getZennStats(actualUserId),
        ]);
        zennArticles = articlesRes.data || [];
        zennStats = statsRes.data;
      }

      let qiitaArticles: QiitaArticle[] = [];
      let qiitaStats: QiitaStats | null = null;
      if (userData.qiita_username) {
        const [articlesRes, statsRes] = await Promise.all([
          getQiitaArticles(actualUserId),
          getQiitaStats(actualUserId),
        ]);
        qiitaArticles = articlesRes.data || [];
        qiitaStats = statsRes.data;
      }

      let atcoderRating: AtCoderRatingInfo | null = null;
      if (userData.atcoder_username) {
        try {
          const atcoderRes = await getAtCoderRating(userData.atcoder_username);
          atcoderRating = atcoderRes.data;
        } catch {
          // AtCoder API取得失敗時は無視
        }
      }

      const [goalsRes, goalStatsRes, badgesRes, streakRes] = await Promise.all([
        getUserGoals(actualUserId),
        getGoalStats(actualUserId),
        getUserBadges(actualUserId),
        getStreakInfo(actualUserId),
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
        atcoderRating,
        goals: goalsRes.data || [],
        goalStats: goalStatsRes.data,
        followerCount: (followersRes.data || []).length,
        followingCount: (followingRes.data || []).length,
        badges: badgesRes.data?.badges || [],
        streakInfo: streakRes.data || null,
      };
    },
    { deps: [usernameOrId], enabled: !!usernameOrId }
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
    badges: data?.badges ?? [],
    atcoderRating: data?.atcoderRating ?? null,
    streakInfo: data?.streakInfo ?? null,
    loading,
    refetch,
  };
}
