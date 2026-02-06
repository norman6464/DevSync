import { useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import type { PortfolioTheme } from '../../utils/portfolioGenerator';
import { generatePortfolioHTML, downloadPortfolio } from '../../utils/portfolioGenerator';
import type { User } from '../../types/user';
import type { GitHubLanguageStat, GitHubRepository } from '../../types/github';
import type { LearningGoal } from '../../api/goals';

interface PortfolioModalProps {
  isOpen: boolean;
  onClose: () => void;
  user: User;
  languages: GitHubLanguageStat[];
  repos: GitHubRepository[];
  goals: LearningGoal[];
  totalContributions: number;
  followerCount: number;
  followingCount: number;
}

const themes: { id: PortfolioTheme; label: string; color: string }[] = [
  { id: 'minimal', label: 'Minimal', color: 'bg-white border-2 border-gray-300' },
  { id: 'modern', label: 'Modern', color: 'bg-slate-900' },
  { id: 'gradient', label: 'Gradient', color: 'bg-gradient-to-br from-indigo-500 to-purple-600' },
];

export default function PortfolioModal({
  isOpen,
  onClose,
  user,
  languages,
  repos,
  goals,
  totalContributions,
  followerCount,
  followingCount
}: PortfolioModalProps) {
  const { t } = useTranslation();
  const [theme, setTheme] = useState<PortfolioTheme>('modern');
  const [showCode, setShowCode] = useState(false);

  const data = { user, languages, repos, goals, totalContributions, followerCount, followingCount };
  const html = useMemo(() => generatePortfolioHTML(data, theme), [data, theme]);

  const handleDownload = () => {
    const filename = `${user.name.toLowerCase().replace(/\s+/g, '-')}-portfolio.html`;
    downloadPortfolio(html, filename);
    toast.success(t('sharing.downloaded'));
  };

  const handleCopyHTML = async () => {
    await navigator.clipboard.writeText(html);
    toast.success(t('portfolio.copied'));
  };

  const handleOpenInNewTab = () => {
    const blob = new Blob([html], { type: 'text/html' });
    const url = URL.createObjectURL(blob);
    window.open(url, '_blank');
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-800 rounded-xl w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-700">
          <h2 className="text-xl font-semibold text-white">{t('portfolio.title')}</h2>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-white transition-colors"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Theme Selection */}
        <div className="p-4 border-b border-gray-700">
          <label className="block text-sm font-medium text-gray-300 mb-2">
            {t('portfolio.theme')}
          </label>
          <div className="flex gap-3">
            {themes.map((t) => (
              <button
                key={t.id}
                onClick={() => setTheme(t.id)}
                className={`flex flex-col items-center gap-2 p-3 rounded-lg transition-all ${
                  theme === t.id ? 'ring-2 ring-green-500' : 'hover:bg-gray-700'
                }`}
              >
                <div className={`w-16 h-12 rounded-md ${t.color}`} />
                <span className="text-sm text-gray-300">{t.label}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Preview */}
        <div className="flex-1 overflow-auto p-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-300">{t('portfolio.preview')}</span>
            <button
              onClick={() => setShowCode(!showCode)}
              className="text-sm text-gray-400 hover:text-white transition-colors"
            >
              {showCode ? t('portfolio.hideCode') : t('portfolio.showCode')}
            </button>
          </div>

          {showCode ? (
            <pre className="bg-gray-900 p-4 rounded-lg overflow-auto text-xs text-gray-300 max-h-96">
              {html}
            </pre>
          ) : (
            <div className="bg-white rounded-lg overflow-hidden">
              <iframe
                srcDoc={html}
                className="w-full h-96 border-0"
                title="Portfolio Preview"
              />
            </div>
          )}
        </div>

        {/* Actions */}
        <div className="p-4 border-t border-gray-700 flex flex-wrap gap-3">
          <button
            onClick={handleOpenInNewTab}
            className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
            </svg>
            {t('portfolio.openInNewTab')}
          </button>
          <button
            onClick={handleDownload}
            className="flex items-center gap-2 px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-lg transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
            </svg>
            {t('portfolio.download')}
          </button>
          <button
            onClick={handleCopyHTML}
            className="flex items-center gap-2 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 17.25v3.375c0 .621-.504 1.125-1.125 1.125h-9.75a1.125 1.125 0 0 1-1.125-1.125V7.875c0-.621.504-1.125 1.125-1.125H6.75a9.06 9.06 0 0 1 1.5.124m7.5 10.376h3.375c.621 0 1.125-.504 1.125-1.125V11.25c0-4.46-3.243-8.161-7.5-8.876a9.06 9.06 0 0 0-1.5-.124H9.375c-.621 0-1.125.504-1.125 1.125v3.5m7.5 10.375H9.375a1.125 1.125 0 0 1-1.125-1.125v-9.25m12 6.625v-1.875a3.375 3.375 0 0 0-3.375-3.375h-1.5a1.125 1.125 0 0 1-1.125-1.125v-1.5a3.375 3.375 0 0 0-3.375-3.375H9.75" />
            </svg>
            {t('portfolio.copy')}
          </button>
        </div>
      </div>
    </div>
  );
}
