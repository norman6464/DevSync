import { useState, useEffect, useRef } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  ChevronDown,
  Target,
  BarChart3,
  Map,
  FolderKanban,
  BookOpen,
  Library,
  Settings,
} from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import Avatar from '../common/Avatar';
import ThemeToggle from '../common/ThemeToggle';
import LanguageSelector from '../common/LanguageSelector';
import NotificationDropdown from '../notifications/NotificationDropdown';

/** 常に表示する主要ナビ（5個以内に抑える） */
const navItems = [
  { path: '/', key: 'nav.dashboard' },
  { path: '/search', key: 'nav.explore' },
  { path: '/rankings', key: 'nav.rankings' },
  { path: '/chat', key: 'nav.chat' },
  { path: '/qa', key: 'nav.qa' },
] as const;

/** 「その他」ドロップダウンに格納する項目 */
const moreItems = [
  { path: '/goals', key: 'nav.goals', icon: Target },
  { path: '/reports', key: 'nav.reports', icon: BarChart3 },
  { path: '/roadmaps', key: 'nav.roadmaps', icon: Map },
  { path: '/projects', key: 'nav.projects', icon: FolderKanban },
  { path: '/resources', key: 'nav.resources', icon: Library },
  { path: '/book-reviews', key: 'nav.bookReviews', icon: BookOpen },
  { path: '/settings', key: 'nav.settings', icon: Settings },
] as const;

export default function Header() {
  const { t } = useTranslation();
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();
  const location = useLocation();
  const [mobileOpen, setMobileOpen] = useState(false);
  const [moreOpen, setMoreOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);
  const moreRef = useRef<HTMLDivElement>(null);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const isActive = (path: string) =>
    location.pathname === path
      ? 'text-white bg-gray-800'
      : 'text-white/70 hover:text-white';

  const isMoreActive = moreItems.some((item) => location.pathname === item.path);

  // Close menus on route change
  useEffect(() => {
    setMobileOpen(false);
    setMoreOpen(false);
  }, [location.pathname]);

  // Close mobile menu on outside click
  useEffect(() => {
    if (!mobileOpen) return;
    const handleClick = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setMobileOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, [mobileOpen]);

  // Close more dropdown on outside click
  useEffect(() => {
    if (!moreOpen) return;
    const handleClick = (e: MouseEvent) => {
      if (moreRef.current && !moreRef.current.contains(e.target as Node)) {
        setMoreOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, [moreOpen]);

  return (
    <header className="bg-gray-900/80 backdrop-blur-sm border-b border-gray-800 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 h-16 flex items-center gap-3">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-2 shrink-0">
          <svg className="w-8 h-8 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M16 18l2-2-2-2" />
            <path d="M8 6L6 8l2 2" />
            <path d="M14.5 4l-5 16" />
          </svg>
          <span className="text-lg font-bold text-white hidden sm:block">DevSync</span>
        </Link>

        {/* Mobile hamburger button */}
        <button
          className="md:hidden p-2 text-gray-300 hover:text-white transition-colors rounded-md"
          onClick={() => setMobileOpen(!mobileOpen)}
          aria-label={mobileOpen ? t('common.close') : t('common.menu')}
          aria-expanded={mobileOpen}
        >
          {mobileOpen ? (
            <svg className="w-6 h-6" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          ) : (
            <svg className="w-6 h-6" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
            </svg>
          )}
        </button>

        {/* Desktop Nav */}
        <nav className="hidden md:flex items-center gap-1 ml-2 min-w-0">
          {navItems.map(({ path, key }) => (
            <Link
              key={path}
              to={path}
              className={`px-3 py-1.5 rounded-md text-sm font-medium transition-colors whitespace-nowrap ${isActive(path)}`}
            >
              {t(key)}
            </Link>
          ))}
        </nav>

        {/* More dropdown */}
        <div className="hidden md:block relative shrink-0" ref={moreRef}>
          <button
            onClick={() => setMoreOpen(!moreOpen)}
            className={`flex items-center gap-1 px-3 py-1.5 rounded-md text-sm font-medium transition-colors whitespace-nowrap ${
              isMoreActive ? 'text-white bg-gray-800' : 'text-white/70 hover:text-white'
            }`}
            aria-expanded={moreOpen}
          >
            {t('nav.more')}
            <ChevronDown className={`w-3.5 h-3.5 transition-transform ${moreOpen ? 'rotate-180' : ''}`} />
          </button>

          {moreOpen && (
            <div className="absolute top-full right-0 mt-1 w-48 bg-gray-900 border border-gray-700 rounded-lg shadow-xl py-1 z-50">
              {moreItems.map(({ path, key, icon: Icon }) => (
                <Link
                  key={path}
                  to={path}
                  className={`flex items-center gap-2.5 px-4 py-2.5 text-sm transition-colors ${
                    location.pathname === path
                      ? 'text-white bg-gray-800'
                      : 'text-gray-300 hover:text-white hover:bg-gray-800/50'
                  }`}
                  onClick={() => setMoreOpen(false)}
                >
                  <Icon className="w-4 h-4" />
                  {t(key)}
                </Link>
              ))}
            </div>
          )}
        </div>

        {/* Right side — shrink-0 ensures icons are never compressed */}
        <div className="flex items-center gap-2 ml-auto shrink-0">
          <LanguageSelector />
          <ThemeToggle />
          <NotificationDropdown />

          {user && (
            <Link
              to={`/profile/${user.id}`}
              className="flex items-center gap-2 ml-1"
              aria-label={t('nav.profile')}
            >
              <Avatar name={user.name} avatarUrl={user.avatar_url} size="sm" />
            </Link>
          )}

          <button
            onClick={handleLogout}
            className="p-2 text-gray-300 hover:text-white transition-colors rounded-md"
            title={t('nav.signOut')}
            aria-label={t('nav.signOut')}
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="1.5" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3 0 3-3m0 0-3-3m3 3H9" />
            </svg>
          </button>
        </div>
      </div>

      {/* Mobile Nav Overlay */}
      {mobileOpen && (
        <div ref={menuRef} className="md:hidden border-t border-gray-800 bg-gray-900/95 backdrop-blur-sm">
          <nav className="max-w-7xl mx-auto px-4 py-3 flex flex-col gap-1">
            {navItems.map(({ path, key }) => (
              <Link
                key={path}
                to={path}
                className={`px-4 py-2.5 rounded-md text-sm font-medium transition-colors ${isActive(path)}`}
              >
                {t(key)}
              </Link>
            ))}

            {/* More items */}
            <div className="border-t border-gray-800 my-1" />
            {moreItems.map(({ path, key, icon: Icon }) => (
              <Link
                key={path}
                to={path}
                className={`flex items-center gap-2.5 px-4 py-2.5 rounded-md text-sm font-medium transition-colors ${isActive(path)}`}
              >
                <Icon className="w-4 h-4" />
                {t(key)}
              </Link>
            ))}
          </nav>
        </div>
      )}
    </header>
  );
}
