import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../../store/authStore';
import Avatar from '../common/Avatar';

export default function UserProfileWidget() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  if (!user) return null;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <div className="flex items-center gap-3 mb-3">
        <Avatar name={user.name} avatarUrl={user.avatar_url} size="md" />
        <div className="min-w-0">
          <Link to={`/profile/${user.username}`} className="font-medium text-sm hover:text-blue-400 block truncate">
            {user.name}
          </Link>
          {user.github_username && (
            <p className="text-xs text-gray-500 truncate">@{user.github_username}</p>
          )}
        </div>
      </div>
      <div className="space-y-0.5">
        <Link to={`/profile/${user.username}`} className="flex items-center gap-2 text-sm text-gray-400 hover:text-white hover:bg-gray-800 py-1.5 px-2 rounded-md transition-colors">
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.9 17.9 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" /></svg>
          {t('dashboard.yourProfile')}
        </Link>
        <Link to="/settings" className="flex items-center gap-2 text-sm text-gray-400 hover:text-white hover:bg-gray-800 py-1.5 px-2 rounded-md transition-colors">
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24" aria-hidden="true"><path strokeLinecap="round" strokeLinejoin="round" d="M10.343 3.94c.09-.542.56-.94 1.11-.94h1.093c.55 0 1.02.398 1.11.94l.149.894c.07.424.384.764.78.93l.164.076c.374.174.798.203 1.167.067l.838-.34a1.114 1.114 0 0 1 1.365.486l.547.948c.271.47.163 1.07-.258 1.414l-.69.577a1.18 1.18 0 0 0-.378 1.001c.006.1.006.2 0 .3a1.18 1.18 0 0 0 .378 1.001l.69.577c.421.345.529.944.258 1.414l-.547.948a1.114 1.114 0 0 1-1.365.486l-.838-.34a1.18 1.18 0 0 0-1.167.067l-.164.076c-.396.166-.71.506-.78.93l-.149.894c-.09.542-.56.94-1.11.94h-1.093c-.55 0-1.02-.398-1.11-.94l-.149-.894a1.18 1.18 0 0 0-.78-.93l-.164-.076a1.18 1.18 0 0 0-1.167-.067l-.838.34a1.114 1.114 0 0 1-1.365-.486l-.547-.948a1.114 1.114 0 0 1 .258-1.414l.69-.577a1.18 1.18 0 0 0 .378-1.001 2 2 0 0 1 0-.3 1.18 1.18 0 0 0-.378-1.001l-.69-.577a1.114 1.114 0 0 1-.258-1.414l.547-.948a1.114 1.114 0 0 1 1.365-.486l.838.34c.369.136.793.107 1.167-.067l.164-.076c.396-.166.71-.506.78-.93l.149-.894Z" /><path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" /></svg>
          {t('nav.settings')}
        </Link>
      </div>
      {!user.github_connected && (
        <Link
          to="/settings"
          className="mt-3 flex items-center gap-2 text-sm text-amber-400 hover:text-amber-300 py-2 px-3 bg-amber-400/10 border border-amber-400/20 rounded-lg transition-colors"
        >
          <svg className="w-4 h-4 shrink-0" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
          {t('dashboard.connectGitHub')}
        </Link>
      )}
    </div>
  );
}
