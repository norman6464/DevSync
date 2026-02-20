import { useState, useMemo } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Sparkles, Rocket, FolderOpen, Pin } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { useProfile, usePostSeries, usePostCollections, useProfileCompleteness, usePinnedPosts } from '../hooks';
import { sectionContainerClass, emptyStateClass } from '../constants/styles';
import { PageLoader } from '../components/common';
import ProfileHeader from '../components/profile/ProfileHeader';
import ContributionCalendar from '../components/profile/ContributionCalendar';
import LanguageChart from '../components/profile/LanguageChart';
import BadgeDisplay from '../components/profile/BadgeDisplay';
import StreakDisplay from '../components/profile/StreakDisplay';
import LevelDisplay from '../components/profile/LevelDisplay';
import PostCard from '../components/posts/PostCard';
import PostSeriesCard from '../components/series/PostSeriesCard';
import ProfileCompletenessCard from '../components/profile/ProfileCompletenessCard';
import CompetitiveProgrammingCard from '../components/profile/CompetitiveProgrammingCard';
import ShareModal from '../components/profile/ShareModal';
import PortfolioModal from '../components/profile/PortfolioModal';
import SpotifyNowPlaying from '../components/profile/SpotifyNowPlaying';
import ProfileGoalsSection from '../components/profile/ProfileGoalsSection';
import ProfileRepositoriesSection from '../components/profile/ProfileRepositoriesSection';
import ProfileArticlesSection from '../components/profile/ProfileArticlesSection';


export default function ProfilePage() {
  const { t } = useTranslation();
  const { username } = useParams<{ username: string }>();
  const currentUser = useAuthStore((s) => s.user);
  const {
    user, posts, contributions, languages, repos,
    zennArticles, zennStats, qiitaArticles, qiitaStats, atcoderRating,
    goals, goalStats, followerCount, followingCount, badges, streakInfo, loading,
  } = useProfile(username);

  const [shareModalOpen, setShareModalOpen] = useState(false);
  const [portfolioModalOpen, setPortfolioModalOpen] = useState(false);
  const { series } = usePostSeries(user?.id);
  const { collections } = usePostCollections(user?.id);
  const { pins } = usePinnedPosts(user?.id);
  const { percentage, missingFields } = useProfileCompleteness();
  const totalContributions = useMemo(() => contributions.reduce((sum, c) => sum + c.count, 0), [contributions]);

  if (loading) return <PageLoader />;
  if (!user) return <div className="text-center text-gray-400 py-12">{t('errors.notFound')}</div>;

  const isOwnProfile = currentUser?.id === user.id;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <ProfileHeader
        user={user}
        isOwnProfile={isOwnProfile}
        followerCount={followerCount}
        followingCount={followingCount}
        onShareClick={() => setShareModalOpen(true)}
        onPortfolioClick={() => setPortfolioModalOpen(true)}
      />

      {/* Profile Completeness */}
      {isOwnProfile && (
        <ProfileCompletenessCard percentage={percentage} missingFields={missingFields} />
      )}

      {/* Skills */}
      {(user.skills_languages || user.skills_frameworks) && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-6 space-y-4">
          {user.skills_languages && (
            <div>
              <h3 className="text-sm font-semibold text-gray-300 mb-3 flex items-center gap-2"><Sparkles className="w-4 h-4 text-yellow-400" aria-hidden="true" /> {t('profile.languages')}</h3>
              <a href="https://skillicons.dev" target="_blank" rel="noopener noreferrer"><img src={`https://skillicons.dev/icons?${new URLSearchParams({ i: user.skills_languages, theme: 'dark' })}`} alt="Languages" referrerPolicy="no-referrer" className="h-12" /></a>
            </div>
          )}
          {user.skills_frameworks && (
            <div>
              <h3 className="text-sm font-semibold text-gray-300 mb-3 flex items-center gap-2"><Rocket className="w-4 h-4 text-blue-400" aria-hidden="true" /> {t('profile.frameworks')}</h3>
              <a href="https://skillicons.dev" target="_blank" rel="noopener noreferrer"><img src={`https://skillicons.dev/icons?${new URLSearchParams({ i: user.skills_frameworks, theme: 'dark' })}`} alt="Frameworks" referrerPolicy="no-referrer" className="h-12" /></a>
            </div>
          )}
        </div>
      )}

      {/* AtCoder Rating & paiza Rank */}
      <CompetitiveProgrammingCard
        atcoderRating={atcoderRating}
        atcoderUsername={user.atcoder_username}
        paizaRank={user.paiza_rank}
      />

      {/* Spotify Now Playing */}
      {user.spotify_connected && <SpotifyNowPlaying userId={user.id} />}

      {/* GitHub Data */}
      {user.github_connected && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className={`lg:col-span-2 ${sectionContainerClass}`}>
            <div className="px-6 py-4 border-b border-gray-800"><h2 className="text-sm font-semibold">{t('profile.contributions')}</h2></div>
            <div className="p-6"><ContributionCalendar contributions={contributions} /></div>
          </div>
          {languages.length > 0 && (
            <div className={sectionContainerClass}>
              <div className="px-6 py-4 border-b border-gray-800"><h2 className="text-sm font-semibold">{t('profile.languages')}</h2></div>
              <div className="p-6"><LanguageChart languages={languages} /></div>
            </div>
          )}
        </div>
      )}

      <LevelDisplay userId={user.id} />

      <StreakDisplay streakInfo={streakInfo} />

      <BadgeDisplay badges={badges} />

      {/* Repositories */}
      {user.github_connected && (
        <ProfileRepositoriesSection repos={repos} githubUsername={user.github_username} />
      )}

      {/* Zenn & Qiita Articles */}
      <ProfileArticlesSection
        zennUsername={user.zenn_username}
        zennArticles={zennArticles}
        zennStats={zennStats}
        qiitaUsername={user.qiita_username}
        qiitaArticles={qiitaArticles}
        qiitaStats={qiitaStats}
      />

      {/* Learning Goals */}
      <ProfileGoalsSection goals={goals} goalStats={goalStats} isOwnProfile={isOwnProfile} />

      {/* Post Series */}
      {series.length > 0 && (
        <div>
          <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">{t('series.title')}</h2>
          <div className="space-y-3">
            {series.map((s) => (
              <PostSeriesCard key={s.id} series={s} isOwner={isOwnProfile} />
            ))}
          </div>
        </div>
      )}

      {/* Post Collections */}
      {collections.length > 0 && (
        <div>
          <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">{t('collections.title')}</h2>
          <div className="grid gap-3 grid-cols-1 sm:grid-cols-2">
            {collections.map((c) => (
              <Link
                key={c.id}
                to={`/collections/${c.id}`}
                className="bg-gray-900 border border-gray-800 rounded-md p-4 hover:border-gray-700 transition-colors"
              >
                <div className="flex items-center gap-2">
                  <FolderOpen className="w-4 h-4 text-purple-400" aria-hidden="true" />
                  <h3 className="text-white font-medium truncate">{c.title}</h3>
                </div>
                {c.description && (
                  <p className="text-gray-400 text-sm mt-1 line-clamp-2">{c.description}</p>
                )}
              </Link>
            ))}
          </div>
        </div>
      )}

      {/* Pinned Posts */}
      {pins.length > 0 && (
        <div>
          <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3 flex items-center gap-2">
            <Pin className="w-4 h-4 text-yellow-400" aria-hidden="true" />
            {t('pinnedPosts.title')}
          </h2>
          <div className="space-y-3">
            {pins.map((pin) => pin.post && (
              <PostCard
                key={pin.id}
                post={pin.post}
                isOwner={currentUser?.id === pin.post.user_id}
              />
            ))}
          </div>
        </div>
      )}

      {/* Posts */}
      <div>
        <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide mb-3">{t('profile.posts')}</h2>
        {posts.length === 0 ? (
          <div className={`${emptyStateClass} text-gray-400 text-sm`}>{t('profile.noPosts')}</div>
        ) : (
          <div className="space-y-3">{posts.map((post) => (
            <PostCard
              key={post.id}
              post={post}
              isOwner={currentUser?.id === post.user_id}
            />
          ))}</div>
        )}
      </div>

      <ShareModal isOpen={shareModalOpen} onClose={() => setShareModalOpen(false)} user={user} followerCount={followerCount} followingCount={followingCount} totalContributions={totalContributions} languages={languages} postCount={posts.length} />
      <PortfolioModal isOpen={portfolioModalOpen} onClose={() => setPortfolioModalOpen(false)} user={user} languages={languages} repos={repos} goals={goals} totalContributions={totalContributions} followerCount={followerCount} followingCount={followingCount} />
    </div>
  );
}
