import { useTranslation } from 'react-i18next';

interface ShareButtonsProps {
  onDownload: () => void;
  onCopyLink: () => void;
  onShareTwitter: () => void;
  onShareLinkedIn: () => void;
  generating: boolean;
}

export default function ShareButtons({
  onDownload,
  onCopyLink,
  onShareTwitter,
  onShareLinkedIn,
  generating,
}: ShareButtonsProps) {
  const { t } = useTranslation();

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-3" role="group" aria-label={t('sharing.title')}>
      {/* Download */}
      <button
        onClick={onDownload}
        disabled={generating}
        aria-label={t('sharing.download')}
        className="flex flex-col items-center gap-2 p-4 bg-gray-800 hover:bg-gray-700 rounded-md transition-colors disabled:opacity-50"
      >
        <div className="w-10 h-10 flex items-center justify-center bg-green-500/20 text-green-400 rounded-full">
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
        </div>
        <span className="text-sm">
          {generating ? t('common.loading') : t('sharing.download')}
        </span>
      </button>

      {/* Copy Link */}
      <button
        onClick={onCopyLink}
        aria-label={t('sharing.copyLink')}
        className="flex flex-col items-center gap-2 p-4 bg-gray-800 hover:bg-gray-700 rounded-md transition-colors"
      >
        <div className="w-10 h-10 flex items-center justify-center bg-blue-500/20 text-blue-400 rounded-full">
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
          </svg>
        </div>
        <span className="text-sm">{t('sharing.copyLink')}</span>
      </button>

      {/* Twitter */}
      <button
        onClick={onShareTwitter}
        aria-label={t('sharing.shareOnTwitter')}
        className="flex flex-col items-center gap-2 p-4 bg-gray-800 hover:bg-gray-700 rounded-md transition-colors"
      >
        <div className="w-10 h-10 flex items-center justify-center bg-sky-500/20 text-sky-400 rounded-full">
          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
          </svg>
        </div>
        <span className="text-sm">Twitter</span>
      </button>

      {/* LinkedIn */}
      <button
        onClick={onShareLinkedIn}
        aria-label={t('sharing.shareOnLinkedIn')}
        className="flex flex-col items-center gap-2 p-4 bg-gray-800 hover:bg-gray-700 rounded-md transition-colors"
      >
        <div className="w-10 h-10 flex items-center justify-center bg-blue-600/20 text-blue-500 rounded-full">
          <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z" />
          </svg>
        </div>
        <span className="text-sm">LinkedIn</span>
      </button>
    </div>
  );
}
