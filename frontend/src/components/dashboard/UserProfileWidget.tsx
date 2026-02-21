import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { User, Settings, Github } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import Avatar from '../common/Avatar';
import { panelClass } from '../../constants/styles';

export default function UserProfileWidget() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);

  if (!user) return null;

  return (
    <div className={panelClass}>
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
          <User className="w-4 h-4" aria-hidden="true" />
          {t('dashboard.yourProfile')}
        </Link>
        <Link to="/settings" className="flex items-center gap-2 text-sm text-gray-400 hover:text-white hover:bg-gray-800 py-1.5 px-2 rounded-md transition-colors">
          <Settings className="w-4 h-4" aria-hidden="true" />
          {t('nav.settings')}
        </Link>
      </div>
      {!user.github_connected && (
        <Link
          to="/settings"
          className="mt-3 flex items-center gap-2 text-sm text-amber-400 hover:text-amber-300 py-2 px-3 bg-amber-400/10 border border-amber-400/20 rounded-lg transition-colors"
        >
          <Github className="w-4 h-4 shrink-0" aria-hidden="true" />
          {t('dashboard.connectGitHub')}
        </Link>
      )}
    </div>
  );
}
