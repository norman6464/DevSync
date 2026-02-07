import { Outlet } from 'react-router-dom';
import Header from './Header';
import ScrollToTop from '../common/ScrollToTop';
import ScrollToTopOnNavigate from '../common/ScrollToTopOnNavigate';

export default function Layout() {
  return (
    <div className="min-h-screen bg-gray-950 text-gray-100">
      <ScrollToTopOnNavigate />
      <Header />
      <main className="max-w-7xl mx-auto px-4 py-6">
        <Outlet />
      </main>
      <ScrollToTop />
    </div>
  );
}
