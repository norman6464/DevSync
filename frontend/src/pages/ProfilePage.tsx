import { useState, useMemo } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Sparkles, Rocket, FolderOpen, Pin } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { useProfile, usePostSeries, usePostCollections, useProfileCompleteness, usePinnedPosts } from '../hooks';
import { sectionContainerClass } from '../constants/styles';
import { PageLoader } from '../components/common';
import ProfileHeader from '../components/profile/ProfileHeader';
import ContributionCalendar from '../components/profile/ContributionCalendar';
import LanguageChart from '../components/profile/LanguageChart';
import BadgeDisplay from '../components/profile/BadgeDisplay';
import StreakDisplay from '../components/profile/StreakDisplay';
import LevelDisplay from '../components/profile/LevelDisplay';
import PostCard from '../components/posts/PostCard';
import PostSeriesCard from '../components/series/PostSeriesCard';
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
      {isOwnProfile && percentage < 100 && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-5">
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <svg className="w-4 h-4 text-yellow-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
              </svg>
              {t('profile.completeness')}
            </h3>
            <span className="text-sm font-bold text-yellow-400">{percentage}%</span>
          </div>
          <div className="h-2 bg-gray-800 rounded-full overflow-hidden mb-3" role="progressbar" aria-valuenow={percentage} aria-valuemin={0} aria-valuemax={100} aria-label={`${t('profile.completeness')}: ${percentage}%`}>
            <div
              className="h-full bg-gradient-to-r from-yellow-500 to-green-500 rounded-full transition-all duration-500"
              style={{ width: `${percentage}%` }}
            />
          </div>
          <div className="flex flex-wrap gap-2">
            {missingFields.map((field) => (
              <Link
                key={field}
                to="/settings"
                className="px-2.5 py-1 bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white text-xs rounded-lg transition-colors"
              >
                + {t(`profile.missing.${field}`)}
              </Link>
            ))}
          </div>
        </div>
      )}

      {/* Skills */}
      {(user.skills_languages || user.skills_frameworks) && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-6 space-y-4">
          {user.skills_languages && (
            <div>
              <h3 className="text-sm font-semibold text-gray-300 mb-3 flex items-center gap-2"><Sparkles className="w-4 h-4 text-yellow-400" aria-hidden="true" /> {t('profile.languages')}</h3>
              <a href="https://skillicons.dev" target="_blank" rel="noopener noreferrer"><img src={`https://skillicons.dev/icons?${new URLSearchParams({ i: user.skills_languages, theme: 'dark' })}`} alt="Languages" className="h-12" /></a>
            </div>
          )}
          {user.skills_frameworks && (
            <div>
              <h3 className="text-sm font-semibold text-gray-300 mb-3 flex items-center gap-2"><Rocket className="w-4 h-4 text-blue-400" aria-hidden="true" /> {t('profile.frameworks')}</h3>
              <a href="https://skillicons.dev" target="_blank" rel="noopener noreferrer"><img src={`https://skillicons.dev/icons?${new URLSearchParams({ i: user.skills_frameworks, theme: 'dark' })}`} alt="Frameworks" className="h-12" /></a>
            </div>
          )}
        </div>
      )}

      {/* AtCoder Rating & paiza Rank */}
      {(atcoderRating || user.paiza_rank) && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
          <h3 className="text-sm font-semibold text-gray-300 mb-4 flex items-center gap-2">
            <svg className="w-4 h-4 text-cyan-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M18.75 4.236c.982.143 1.954.317 2.916.52A6.003 6.003 0 0 1 16.27 9.728M18.75 4.236V4.5c0 2.108-.966 3.99-2.48 5.228m0 0a6.023 6.023 0 0 1-2.52.587 6.023 6.023 0 0 1-2.52-.587" /></svg>
            {t('profile.competitiveProgramming')}
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
            {atcoderRating && (
              <a href={`https://atcoder.jp/users/${user.atcoder_username}`} target="_blank" rel="noopener noreferrer" className="flex items-center gap-4 p-4 bg-gray-800/50 rounded-lg border border-gray-700 hover:border-gray-600 transition-colors group">
                <div className="w-12 h-12 bg-gray-700 rounded-lg flex items-center justify-center text-white font-bold text-lg">A</div>
                <div>
                  <div className="text-sm text-gray-400">AtCoder</div>
                  <div className="flex items-center gap-2">
                    <span className={`text-xl font-bold atcoder-${atcoderRating.color}`} style={{ color: atcoderRating.color === 'gray' ? '#808080' : atcoderRating.color === 'brown' ? '#804000' : atcoderRating.color === 'green' ? '#008000' : atcoderRating.color === 'cyan' ? '#00C0C0' : atcoderRating.color === 'blue' ? '#0000FF' : atcoderRating.color === 'yellow' ? '#C0C000' : atcoderRating.color === 'orange' ? '#FF8000' : atcoderRating.color === 'red' ? '#FF0000' : '#808080' }}>
                      {atcoderRating.rating}
                    </span>
                    <span className="text-xs text-gray-500">({atcoderRating.rank})</span>
                  </div>
                </div>
              </a>
            )}
            {user.paiza_rank && (
              <div className="flex items-center gap-4 p-4 bg-gray-800/50 rounded-lg border border-gray-700">
                <div className="w-12 h-12 bg-emerald-700 rounded-lg flex items-center justify-center text-white font-bold text-lg">P</div>
                <div>
                  <div className="text-sm text-gray-400">paiza</div>
                  <div className="flex items-center gap-2">
                    <span className="text-xl font-bold text-white">
                      {t('profile.paizaRankLabel', { rank: user.paiza_rank })}
                    </span>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

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
          <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center text-gray-400 text-sm">{t('profile.noPosts')}</div>
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
