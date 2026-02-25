import { useTranslation } from 'react-i18next';
import { sectionContainerClass, buttonDangerOutlineClass } from '../../../constants/styles';

const GitHubIcon = ({ className }: { className?: string }) => (
  <svg className={className} fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
);

interface Props {
  connected: boolean;
  username?: string;
  syncing: boolean;
  onConnect: () => void;
  onDisconnect: () => void;
  onSync: () => void;
}

export default function GitHubIntegrationCard({ connected, username, syncing, onConnect, onDisconnect, onSync }: Props) {
  const { t } = useTranslation();

  return (
    <div className={sectionContainerClass}>
      <div className="px-6 py-4 border-b border-gray-800">
        <h2 className="text-base font-semibold">{t('settings.github')}</h2>
      </div>
      <div className="p-6">
        {connected ? (
          <div className="space-y-4">
            <div className="flex items-center gap-3">
              <GitHubIcon className="w-8 h-8 text-gray-300" />
              <div>
                <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
                <p className="text-sm text-gray-400">@{username}</p>
              </div>
            </div>
            <div className="flex gap-2">
              <button
                onClick={onSync}
                disabled={syncing}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium border border-gray-700 transition-colors"
              >
                {syncing ? t('settings.syncing') : t('settings.sync')}
              </button>
              <button onClick={onDisconnect} className={buttonDangerOutlineClass}>
                {t('settings.disconnect')}
              </button>
            </div>
          </div>
        ) : (
          <div className="text-center py-4">
            <GitHubIcon className="w-12 h-12 text-gray-600 mx-auto mb-3" />
            <p className="text-gray-400 text-sm mb-4">{t('settings.githubDescription')}</p>
            <button
              onClick={onConnect}
              className="px-5 py-2.5 bg-white hover:bg-gray-100 text-gray-900 rounded-lg font-semibold text-sm transition-colors inline-flex items-center gap-2"
            >
              <GitHubIcon className="w-5 h-5" />
              {t('settings.connect')} GitHub
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
