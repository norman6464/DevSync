import { useTranslation } from 'react-i18next';
import { sectionContainerClass, buttonDangerOutlineClass } from '../../../constants/styles';

const SpotifyIcon = ({ className }: { className?: string }) => (
  <svg className={className} viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z" />
  </svg>
);

interface Props {
  connected: boolean;
  onConnect: () => void;
  onDisconnect: () => void;
}

export default function SpotifyIntegrationCard({ connected, onConnect, onDisconnect }: Props) {
  const { t } = useTranslation();

  return (
    <div className={sectionContainerClass}>
      <div className="px-6 py-4 border-b border-gray-800">
        <h2 className="text-base font-semibold">{t('settings.spotify')}</h2>
      </div>
      <div className="p-6">
        {connected ? (
          <div className="space-y-4">
            <div className="flex items-center gap-3">
              <SpotifyIcon className="w-8 h-8 text-green-500" />
              <div>
                <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
              </div>
            </div>
            <div className="flex gap-2">
              <button onClick={onDisconnect} className={buttonDangerOutlineClass}>
                {t('settings.disconnect')}
              </button>
            </div>
          </div>
        ) : (
          <div className="text-center py-4">
            <SpotifyIcon className="w-12 h-12 text-green-500/40 mx-auto mb-3" />
            <p className="text-gray-400 text-sm mb-4">{t('settings.spotifyDescription')}</p>
            <button
              onClick={onConnect}
              className="px-5 py-2.5 bg-green-600 hover:bg-green-500 text-white rounded-lg font-semibold text-sm transition-colors inline-flex items-center gap-2"
            >
              <SpotifyIcon className="w-5 h-5" />
              {t('settings.connect')} Spotify
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
