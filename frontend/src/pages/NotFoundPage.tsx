import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Home, Search } from 'lucide-react';

export default function NotFoundPage() {
  const { t } = useTranslation();

  return (
    <div className="min-h-screen bg-gray-950 flex items-center justify-center px-4">
      <div className="text-center max-w-md">
        {/* 404 Icon */}
        <div className="mb-8">
          <div className="inline-flex items-center justify-center w-32 h-32 rounded-full bg-gray-900 border-2 border-gray-800">
            <span className="text-6xl font-bold text-gray-400">404</span>
          </div>
        </div>

        {/* Title */}
        <h1 className="text-3xl font-bold text-white mb-4">
          {t('errors.pageNotFound', 'ページが見つかりません')}
        </h1>

        {/* Description */}
        <p className="text-gray-400 mb-8">
          {t('errors.pageNotFoundDesc', 'お探しのページは存在しないか、移動または削除された可能性があります。')}
        </p>

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-3 justify-center">
          <Link
            to="/"
            className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium transition-colors"
          >
            <Home className="w-5 h-5" />
            {t('common.backToHome', 'ホームに戻る')}
          </Link>
          <Link
            to="/search"
            className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-gray-800 hover:bg-gray-700 text-white rounded-lg font-medium transition-colors"
          >
            <Search className="w-5 h-5" />
            {t('common.search', '検索')}
          </Link>
        </div>

        {/* Helpful Links */}
        <div className="mt-12 pt-8 border-t border-gray-800">
          <p className="text-sm text-gray-500 mb-4">
            {t('errors.helpfulLinks', 'よく見られるページ')}
          </p>
          <div className="flex flex-wrap gap-3 justify-center text-sm">
            <Link to="/dashboard" className="text-blue-400 hover:text-blue-300">
              {t('nav.dashboard', 'ダッシュボード')}
            </Link>
            <span className="text-gray-700">•</span>
            <Link to="/goals" className="text-blue-400 hover:text-blue-300">
              {t('nav.goals', '目標')}
            </Link>
            <span className="text-gray-700">•</span>
            <Link to="/qa" className="text-blue-400 hover:text-blue-300">
              {t('nav.qa', 'Q&A')}
            </Link>
            <span className="text-gray-700">•</span>
            <Link to="/resources" className="text-blue-400 hover:text-blue-300">
              {t('nav.resources', 'リソース')}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
