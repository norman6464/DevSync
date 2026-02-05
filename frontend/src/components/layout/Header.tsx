import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';

export default function Header() {
  const { user, logout } = useAuthStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="bg-gray-900 text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
        <Link to="/" className="text-xl font-bold tracking-tight">
          DevSync
        </Link>

        <nav className="flex items-center gap-6">
          <Link to="/" className="hover:text-blue-400 transition-colors">
            Home
          </Link>
          <Link to="/search" className="hover:text-blue-400 transition-colors">
            Search
          </Link>
          <Link to="/rankings" className="hover:text-blue-400 transition-colors">
            Rankings
          </Link>
          <Link to="/chat" className="hover:text-blue-400 transition-colors">
            Chat
          </Link>
        </nav>

        <div className="flex items-center gap-4">
          {user && (
            <Link
              to={`/profile/${user.id}`}
              className="flex items-center gap-2 hover:text-blue-400 transition-colors"
            >
              <div className="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center text-sm font-medium">
                {user.name.charAt(0).toUpperCase()}
              </div>
              <span className="hidden sm:inline">{user.name}</span>
            </Link>
          )}
          <Link to="/settings" className="hover:text-blue-400 transition-colors text-sm">
            Settings
          </Link>
          <button
            onClick={handleLogout}
            className="text-sm text-gray-400 hover:text-white transition-colors"
          >
            Logout
          </button>
        </div>
      </div>
    </header>
  );
}
