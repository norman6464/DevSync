import { useState, useEffect, useRef, useCallback } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import {
  ChevronDown,
  Target,
  BarChart3,
  Map,
  FolderKanban,
  BookOpen,
  BookMarked,
  Library,
  Settings,
  Lightbulb,
  LineChart,
  Users,
  Youtube,
  Menu,
  X,
  LogOut,
} from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import Avatar from '../common/Avatar';
import ThemeToggle from '../common/ThemeToggle';
import LanguageSelector from '../common/LanguageSelector';
import NotificationDropdown from '../notifications/NotificationDropdown';
import { useClickOutside } from '../../hooks/useClickOutside';

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
  { path: '/learning-logs', key: 'nav.learningLogs', icon: BookMarked },
  { path: '/notes', key: 'notes.title', icon: BookOpen },
  { path: '/analytics', key: 'nav.analytics', icon: LineChart },
  { path: '/reports', key: 'nav.reports', icon: BarChart3 },
  { path: '/roadmaps', key: 'nav.roadmaps', icon: Map },
  { path: '/projects', key: 'nav.projects', icon: FolderKanban },
  { path: '/resources', key: 'nav.resources', icon: Library },
  { path: '/book-reviews', key: 'nav.bookReviews', icon: Library },
  { path: '/study-circles', key: 'nav.studyCircles', icon: Users },
  { path: '/advice', key: 'nav.advice', icon: Lightbulb },
  { path: '/youtube', key: 'nav.youtube', icon: Youtube },
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

  const handleLogout = async () => {
    await logout();
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

  const closeMobile = useCallback(() => setMobileOpen(false), []);
  const closeMore = useCallback(() => setMoreOpen(false), []);
  useClickOutside(menuRef, mobileOpen, closeMobile);
  useClickOutside(moreRef, moreOpen, closeMore);

  return (
    <header className="bg-gray-900/80 backdrop-blur-sm border-b border-gray-800 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 h-16 flex items-center gap-3">
        {/* Logo */}
        <Link to="/" className="flex items-center gap-2 shrink-0">
          <svg className="w-8 h-8 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
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
            <X className="w-6 h-6" aria-hidden="true" />
          ) : (
            <Menu className="w-6 h-6" aria-hidden="true" />
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
            <div
              className="absolute top-full right-0 mt-1 w-48 bg-gray-900 border border-gray-700 rounded-lg shadow-sm py-1 z-50"
              role="menu"
            >
              {moreItems.map(({ path, key, icon: Icon }) => (
                <Link
                  key={path}
                  to={path}
                  role="menuitem"
                  className={`flex items-center gap-2.5 px-4 py-2.5 text-sm transition-colors ${
                    location.pathname === path
                      ? 'text-white bg-gray-800'
                      : 'text-gray-300 hover:text-white hover:bg-gray-800/50'
                  }`}
                  onClick={() => setMoreOpen(false)}
                >
                  <Icon className="w-4 h-4" aria-hidden="true" />
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
              to={`/profile/${user.username}`}
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
            <LogOut className="w-5 h-5" aria-hidden="true" />
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
