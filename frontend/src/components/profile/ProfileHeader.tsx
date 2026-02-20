import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import type { User } from '../../types/user';
import Avatar from '../common/Avatar';
import FollowButton from './FollowButton';

interface ProfileHeaderProps {
  user: User;
  isOwnProfile: boolean;
  followerCount: number;
  followingCount: number;
  onShareClick: () => void;
  onPortfolioClick: () => void;
}

export default function ProfileHeader({
  user,
  isOwnProfile,
  followerCount,
  followingCount,
  onShareClick,
  onPortfolioClick,
}: ProfileHeaderProps) {
  const { t } = useTranslation();

  const socialLinks = [
    { key: 'github', username: user.github_username, href: `https://github.com/${user.github_username}`, color: 'hover:text-blue-400', badge: <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg> },
    { key: 'zenn', username: user.zenn_username, href: `https://zenn.dev/${user.zenn_username}`, color: 'hover:text-blue-400', badge: <span className="w-4 h-4 bg-blue-500 rounded text-white text-xs flex items-center justify-center font-bold">Z</span> },
    { key: 'qiita', username: user.qiita_username, href: `https://qiita.com/${user.qiita_username}`, color: 'hover:text-green-400', badge: <span className="w-4 h-4 bg-green-500 rounded text-white text-xs flex items-center justify-center font-bold">Q</span> },
    { key: 'atcoder', username: user.atcoder_username, href: `https://atcoder.jp/users/${user.atcoder_username}`, color: 'hover:text-cyan-400', badge: <span className="w-4 h-4 bg-gray-700 rounded text-white text-xs flex items-center justify-center font-bold">A</span> },
  ];

  const externalLinkIcon = (
    <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" /></svg>
  );

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
      <div className="flex items-start gap-5">
        <Avatar name={user.name} avatarUrl={user.avatar_url} size="lg" />
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3 flex-wrap">
            <h1 className="text-2xl font-bold">{user.name}</h1>
            {!isOwnProfile && <FollowButton userId={user.id} />}
            <button
              onClick={onShareClick}
              className="flex items-center gap-1.5 px-3 py-1.5 bg-gray-800 hover:bg-gray-700 text-gray-300 text-sm rounded-lg transition-colors"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                <path strokeLinecap="round" strokeLinejoin="round" d="M7.217 10.907a2.25 2.25 0 1 0 0 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186 9.566-5.314m-9.566 7.5 9.566 5.314m0 0a2.25 2.25 0 1 0 3.935 2.186 2.25 2.25 0 0 0-3.935-2.186Zm0-12.814a2.25 2.25 0 1 0 3.933-2.185 2.25 2.25 0 0 0-3.933 2.185Z" />
              </svg>
              {t('sharing.share')}
            </button>
            {isOwnProfile && (
              <button
                onClick={onPortfolioClick}
                className="flex items-center gap-1.5 px-3 py-1.5 bg-purple-600 hover:bg-purple-500 text-sm rounded-lg transition-colors"
                style={{ color: '#ffffff' }}
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 0 0-3.375-3.375h-1.5A1.125 1.125 0 0 1 13.5 7.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 0 0-9-9Z" />
                </svg>
                {t('portfolio.generate')}
              </button>
            )}
          </div>
          {user.bio && <p className="text-gray-400 mt-1 text-sm">{user.bio}</p>}
          <div className="flex flex-wrap gap-3 mt-2">
            {socialLinks.map(link => link.username && (
              <a key={link.key} href={link.href} target="_blank" rel="noopener noreferrer" className={`inline-flex items-center gap-1.5 text-sm text-gray-500 ${link.color} transition-colors`}>
                {link.badge}
                @{link.username}
                {externalLinkIcon}
              </a>
            ))}
          </div>
          <div className="flex gap-4 mt-3 text-sm">
            <Link to={`/profile/${user.username}/followers`} className="text-gray-400 hover:text-blue-400 transition-colors">
              <strong className="text-white">{followerCount}</strong> {t('profile.followers')}
            </Link>
            <Link to={`/profile/${user.username}/following`} className="text-gray-400 hover:text-blue-400 transition-colors">
              <strong className="text-white">{followingCount}</strong> {t('profile.following')}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
