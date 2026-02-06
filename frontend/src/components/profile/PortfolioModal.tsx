import { useState, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import {
  generatePortfolioHTML,
  downloadPortfolio,
  type PortfolioTheme,
  type PortfolioData,
} from '../../utils/portfolioGenerator';

interface PortfolioModalProps {
  isOpen: boolean;
  onClose: () => void;
  data: PortfolioData;
}

const themes: { id: PortfolioTheme; name: string; description: string }[] = [
  { id: 'minimal', name: 'Minimal', description: 'Clean and simple white design' },
  { id: 'modern', name: 'Modern', description: 'Dark theme with cards' },
  { id: 'gradient', name: 'Gradient', description: 'Colorful glassmorphism' },
];

export default function PortfolioModal({
  isOpen,
  onClose,
  data,
}: PortfolioModalProps) {
  const { t } = useTranslation();
  const [selectedTheme, setSelectedTheme] = useState<PortfolioTheme>('modern');
  const [previewMode, setPreviewMode] = useState(false);

  const generatedHTML = useMemo(() => {
    return generatePortfolioHTML(data, selectedTheme);
  }, [data, selectedTheme]);

  if (!isOpen) return null;

  const handleDownload = () => {
    const filename = `portfolio-${data.user.github_username || data.user.name || 'developer'}.html`;
    downloadPortfolio(generatedHTML, filename);
    toast.success(t('portfolio.downloaded'));
  };

  const handleCopyHTML = async () => {
    try {
      await navigator.clipboard.writeText(generatedHTML);
      toast.success(t('portfolio.htmlCopied'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleOpenPreview = () => {
    const blob = new Blob([generatedHTML], { type: 'text/html' });
    const url = URL.createObjectURL(blob);
    window.open(url, '_blank');
    setTimeout(() => URL.revokeObjectURL(url), 10000);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/70 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative bg-gray-900 border border-gray-800 rounded-xl shadow-2xl max-w-4xl w-full mx-4 max-h-[90vh] overflow-hidden flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-800">
          <h2 className="text-lg font-semibold">{t('portfolio.title')}</h2>
          <button
            onClick={onClose}
            className="p-2 text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition-colors"
          >
            <svg
              className="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-6 space-y-6">
          {/* Theme Selection */}
          <div>
            <h3 className="text-sm font-semibold text-gray-300 mb-3">
              {t('portfolio.selectTheme')}
            </h3>
            <div className="grid grid-cols-3 gap-3">
              {themes.map((theme) => (
                <button
                  key={theme.id}
                  onClick={() => setSelectedTheme(theme.id)}
                  className={`p-4 rounded-xl border transition-all text-left ${
                    selectedTheme === theme.id
                      ? 'border-purple-500 bg-purple-500/10'
                      : 'border-gray-700 hover:border-gray-600 bg-gray-800/50'
                  }`}
                >
                  <div className="flex items-center gap-2 mb-2">
                    {theme.id === 'minimal' && (
                      <div className="w-8 h-8 rounded-lg bg-white flex items-center justify-center">
                        <div className="w-4 h-4 rounded bg-gray-300"></div>
                      </div>
                    )}
                    {theme.id === 'modern' && (
                      <div className="w-8 h-8 rounded-lg bg-gray-800 border border-gray-700 flex items-center justify-center">
                        <div className="w-4 h-4 rounded bg-blue-500"></div>
                      </div>
                    )}
                    {theme.id === 'gradient' && (
                      <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-900 to-blue-900 flex items-center justify-center">
                        <div className="w-4 h-4 rounded bg-purple-400/50 backdrop-blur"></div>
                      </div>
                    )}
                    <span className="font-medium text-sm">{theme.name}</span>
                  </div>
                  <p className="text-xs text-gray-500">{theme.description}</p>
                </button>
              ))}
            </div>
          </div>

          {/* Preview */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-semibold text-gray-300">
                {t('portfolio.preview')}
              </h3>
              <div className="flex gap-2">
                <button
                  onClick={() => setPreviewMode(!previewMode)}
                  className="text-xs px-3 py-1.5 bg-gray-800 hover:bg-gray-700 rounded-lg transition-colors"
                >
                  {previewMode ? t('portfolio.hideCode') : t('portfolio.showCode')}
                </button>
                <button
                  onClick={handleOpenPreview}
                  className="text-xs px-3 py-1.5 bg-gray-800 hover:bg-gray-700 rounded-lg transition-colors flex items-center gap-1.5"
                >
                  <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                  </svg>
                  {t('portfolio.openInNewTab')}
                </button>
              </div>
            </div>

            {previewMode ? (
              <div className="bg-gray-950 rounded-xl p-4 overflow-auto max-h-80">
                <pre className="text-xs text-gray-400 font-mono whitespace-pre-wrap break-all">
                  {generatedHTML.slice(0, 3000)}
                  {generatedHTML.length > 3000 && '...'}
                </pre>
              </div>
            ) : (
              <div className="bg-gray-950 rounded-xl overflow-hidden border border-gray-800">
                <iframe
                  srcDoc={generatedHTML}
                  title="Portfolio Preview"
                  className="w-full h-80 bg-white"
                  sandbox="allow-same-origin"
                />
              </div>
            )}
          </div>

          {/* Info */}
          <div className="bg-gray-800/50 rounded-xl p-4">
            <h4 className="text-sm font-medium text-gray-300 mb-2 flex items-center gap-2">
              <svg className="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="m11.25 11.25.041-.02a.75.75 0 0 1 1.063.852l-.708 2.836a.75.75 0 0 0 1.063.853l.041-.021M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9-3.75h.008v.008H12V8.25Z" />
              </svg>
              {t('portfolio.included')}
            </h4>
            <ul className="text-xs text-gray-500 space-y-1">
              <li>• {t('portfolio.includesProfile')}</li>
              <li>• {t('portfolio.includesStats')}</li>
              <li>• {t('portfolio.includesSkills')}</li>
              <li>• {t('portfolio.includesRepos')}</li>
              <li>• {t('portfolio.includesGoals')}</li>
            </ul>
          </div>
        </div>

        {/* Footer */}
        <div className="p-4 border-t border-gray-800 flex gap-3">
          <button
            onClick={handleCopyHTML}
            className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-gray-800 hover:bg-gray-700 rounded-lg transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.666 3.888A2.25 2.25 0 0 0 13.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 0 1-.75.75H9a.75.75 0 0 1-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 0 1-2.25 2.25H6.75A2.25 2.25 0 0 1 4.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 0 1 1.927-.184" />
            </svg>
            {t('portfolio.copyHTML')}
          </button>
          <button
            onClick={handleDownload}
            className="flex-1 flex items-center justify-center gap-2 px-4 py-2.5 bg-purple-600 hover:bg-purple-500 text-white rounded-lg transition-colors"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 16.5v2.25A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75V16.5M16.5 12 12 16.5m0 0L7.5 12m4.5 4.5V3" />
            </svg>
            {t('portfolio.download')}
          </button>
        </div>
      </div>
    </div>
  );
}
