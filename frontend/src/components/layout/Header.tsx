import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../../store/authStore';
import Avatar from '../common/Avatar';
import ThemeToggle from '../common/ThemeToggle';
import LanguageSelector from '../common/LanguageSelector';
import NotificationDropdown from '../notifications/NotificationDropdown';

export default function Header() {
  const { t } = useTranslation();
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const isActive = (path: string) =>
    location.pathname === path ? 'text-white bg-gray-800' : 'text-gray-400 hover:text-white';

  return (
    <header className="bg-gray-900/80 backdrop-blur-sm border-b border-gray-800 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 h-16 flex items-center gap-4">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-2 shrink-0">
          <svg className="w-8 h-8 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M16 18l2-2-2-2" />
            <path d="M8 6L6 8l2 2" />
            <path d="M14.5 4l-5 16" />
          </svg>
          <span className="text-lg font-bold text-white hidden sm:block">DevSync</span>
        </Link>

        {/* Nav */}
        <nav className="flex items-center gap-1 ml-2">
          <Link to="/" className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${isActive('/')}`}>
            {t('nav.dashboard')}
          </Link>
          <Link to="/search" className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${isActive('/search')}`}>
            {t('nav.explore')}
          </Link>
          <Link to="/rankings" className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${isActive('/rankings')}`}>
            {t('nav.rankings')}
          </Link>
          <Link to="/chat" className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${isActive('/chat')}`}>
            {t('nav.chat')}
          </Link>
          <Link to="/goals" className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${isActive('/goals')}`}>
            {t('nav.goals')}
          </Link>
          <Link to="/reports" className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors ${isActive('/reports')}`}>
            {t('nav.reports')}
          </Link>
        </nav>

        {/* Right side */}
        <div className="flex items-center gap-2 ml-auto">
          <LanguageSelector />
          <ThemeToggle />
          <NotificationDropdown />

          <Link
            to="/settings"
            className="p-2 text-gray-400 hover:text-white transition-colors rounded-md"
            title={t('nav.settings')}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a7 7 0 0 1 0 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a7 7 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a7 7 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.992a7 7 0 0 1 0-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124q.108-.066.22-.128c.332-.183.582-.495.644-.869l.214-1.281Z" />
              <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
            </svg>
          </Link>

          {user && (
            <Link
              to={`/profile/${user.id}`}
              className="flex items-center gap-2 ml-1"
            >
              <Avatar name={user.name} avatarUrl={user.avatar_url} size="sm" />
            </Link>
          )}

          <button
            onClick={handleLogout}
            className="p-2 text-gray-400 hover:text-white transition-colors rounded-md"
            title={t('nav.signOut')}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3 0 3-3m0 0-3-3m3 3H9" />
            </svg>
          </button>
        </div>
      </div>
    </header>
  );
}
