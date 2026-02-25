import { useTranslation } from 'react-i18next';
import { inputClass, sectionContainerClass, labelClass, buttonDangerOutlineClass } from '../../../constants/styles';

interface Props {
  title: string;
  description: string;
  // Icon badge
  badgeLetter: string;
  badgeBgClass: string;
  badgeTextClass: string;
  badgeBgDisabledClass: string;
  badgeTextDisabledClass: string;
  // Connect button
  connectBtnClass: string;
  // Connected state
  connected: boolean;
  connectedUsername?: string;
  // Username form
  username: string;
  setUsername: (v: string) => void;
  usernameLabel: string;
  usernamePlaceholder: string;
  // Actions
  connecting: boolean;
  onConnect: (e: React.FormEvent) => void;
  onDisconnect: () => void;
  // Optional sync
  syncing?: boolean;
  onSync?: () => void;
}

export default function UsernameIntegrationCard(props: Props) {
  const { t } = useTranslation();

  return (
    <div className={sectionContainerClass}>
      <div className="px-6 py-4 border-b border-gray-800">
        <h2 className="text-base font-semibold">{props.title}</h2>
      </div>
      <div className="p-6">
        {props.connected ? (
          <div className="space-y-4">
            <div className="flex items-center gap-3">
              <div className={`w-8 h-8 ${props.badgeBgClass} rounded-lg flex items-center justify-center ${props.badgeTextClass} font-bold text-sm`}>
                {props.badgeLetter}
              </div>
              <div>
                <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
                {props.connectedUsername && (
                  <p className="text-sm text-gray-400">@{props.connectedUsername}</p>
                )}
              </div>
            </div>
            <div className="flex gap-2">
              {props.onSync && (
                <button
                  onClick={props.onSync}
                  disabled={props.syncing}
                  className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium border border-gray-700 transition-colors"
                >
                  {props.syncing ? t('settings.syncing') : t('settings.sync')}
                </button>
              )}
              <button
                onClick={props.onDisconnect}
                className={buttonDangerOutlineClass}
              >
                {t('settings.disconnect')}
              </button>
            </div>
          </div>
        ) : (
          <form onSubmit={props.onConnect} className="space-y-4">
            <div className="text-center py-2">
              <div className={`w-12 h-12 ${props.badgeBgDisabledClass} rounded-lg flex items-center justify-center ${props.badgeTextDisabledClass} font-bold text-xl mx-auto mb-3`}>
                {props.badgeLetter}
              </div>
              <p className="text-gray-400 text-sm mb-4">{props.description}</p>
            </div>
            <div>
              <label htmlFor={`integration-${props.badgeLetter.toLowerCase()}-username`} className={labelClass}>
                {props.usernameLabel}
              </label>
              <input
                id={`integration-${props.badgeLetter.toLowerCase()}-username`}
                type="text"
                value={props.username}
                onChange={(e) => props.setUsername(e.target.value)}
                placeholder={props.usernamePlaceholder}
                maxLength={50}
                className={inputClass}
              />
            </div>
            <button
              type="submit"
              disabled={props.connecting || !props.username.trim()}
              className={`w-full px-5 py-2.5 ${props.connectBtnClass} disabled:opacity-50 text-white rounded-lg font-semibold text-sm transition-colors`}
            >
              {props.connecting ? t('common.loading') : t('settings.connect')}
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
