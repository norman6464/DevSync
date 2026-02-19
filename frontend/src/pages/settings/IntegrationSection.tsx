import { useTranslation } from 'react-i18next';
import { inputClass, sectionContainerClass, labelClass, selectClass, buttonDangerOutlineClass } from '../../constants/styles';
import type { User } from '../../types/user';

interface Props {
  user: User;
  // GitHub
  syncing: boolean;
  onConnectGitHub: () => void;
  onDisconnectGitHub: () => void;
  onSyncGitHub: () => void;
  // Zenn
  zennUsername: string;
  setZennUsername: (v: string) => void;
  connectingZenn: boolean;
  syncingZenn: boolean;
  onConnectZenn: (e: React.FormEvent) => void;
  onDisconnectZenn: () => void;
  onSyncZenn: () => void;
  // Qiita
  qiitaUsername: string;
  setQiitaUsername: (v: string) => void;
  connectingQiita: boolean;
  syncingQiita: boolean;
  onConnectQiita: (e: React.FormEvent) => void;
  onDisconnectQiita: () => void;
  onSyncQiita: () => void;
  // AtCoder
  atcoderUsername: string;
  setAtcoderUsername: (v: string) => void;
  connectingAtcoder: boolean;
  onConnectAtCoder: (e: React.FormEvent) => void;
  onDisconnectAtCoder: () => void;
  // Spotify
  onConnectSpotify: () => void;
  onDisconnectSpotify: () => void;
  // Paiza
  paizaRank: string;
  setPaizaRank: (v: string) => void;
  savingPaiza: boolean;
  onSavePaizaRank: () => void;
}

const SpotifyIcon = ({ className }: { className?: string }) => (
  <svg className={className} viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z" />
  </svg>
);

const GitHubIcon = ({ className }: { className?: string }) => (
  <svg className={className} fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
);

export default function IntegrationSection(props: Props) {
  const { t } = useTranslation();
  const { user } = props;

  return (
    <>
      {/* GitHub Integration */}
      <div className={sectionContainerClass}>
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.github')}</h2>
        </div>
        <div className="p-6">
          {user.github_connected ? (
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <GitHubIcon className="w-8 h-8 text-gray-300" />
                <div>
                  <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
                  <p className="text-sm text-gray-400">@{user.github_username}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={props.onSyncGitHub}
                  disabled={props.syncing}
                  className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium border border-gray-700 transition-colors"
                >
                  {props.syncing ? t('settings.syncing') : t('settings.sync')}
                </button>
                <button
                  onClick={props.onDisconnectGitHub}
                  className={buttonDangerOutlineClass}
                >
                  {t('settings.disconnect')}
                </button>
              </div>
            </div>
          ) : (
            <div className="text-center py-4">
              <GitHubIcon className="w-12 h-12 text-gray-600 mx-auto mb-3" />
              <p className="text-gray-400 text-sm mb-4">{t('settings.githubDescription')}</p>
              <button
                onClick={props.onConnectGitHub}
                className="px-5 py-2.5 bg-white hover:bg-gray-100 text-gray-900 rounded-lg font-semibold text-sm transition-colors inline-flex items-center gap-2"
              >
                <GitHubIcon className="w-5 h-5" />
                {t('settings.connect')} GitHub
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Spotify Integration */}
      <div className={sectionContainerClass}>
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.spotify')}</h2>
        </div>
        <div className="p-6">
          {user.spotify_connected ? (
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <SpotifyIcon className="w-8 h-8 text-green-500" />
                <div>
                  <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={props.onDisconnectSpotify}
                  className={buttonDangerOutlineClass}
                >
                  {t('settings.disconnect')}
                </button>
              </div>
            </div>
          ) : (
            <div className="text-center py-4">
              <SpotifyIcon className="w-12 h-12 text-green-500/40 mx-auto mb-3" />
              <p className="text-gray-400 text-sm mb-4">{t('settings.spotifyDescription')}</p>
              <button
                onClick={props.onConnectSpotify}
                className="px-5 py-2.5 bg-green-600 hover:bg-green-500 text-white rounded-lg font-semibold text-sm transition-colors inline-flex items-center gap-2"
              >
                <SpotifyIcon className="w-5 h-5" />
                {t('settings.connect')} Spotify
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Zenn Integration */}
      <div className={sectionContainerClass}>
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.zenn')}</h2>
        </div>
        <div className="p-6">
          {user.zenn_username ? (
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-blue-500 rounded-lg flex items-center justify-center text-white font-bold text-sm">Z</div>
                <div>
                  <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
                  <p className="text-sm text-gray-400">@{user.zenn_username}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={props.onSyncZenn}
                  disabled={props.syncingZenn}
                  className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium border border-gray-700 transition-colors"
                >
                  {props.syncingZenn ? t('settings.syncing') : t('settings.sync')}
                </button>
                <button
                  onClick={props.onDisconnectZenn}
                  className={buttonDangerOutlineClass}
                >
                  {t('settings.disconnect')}
                </button>
              </div>
            </div>
          ) : (
            <form onSubmit={props.onConnectZenn} className="space-y-4">
              <div className="text-center py-2">
                <div className="w-12 h-12 bg-blue-500/20 rounded-lg flex items-center justify-center text-blue-400 font-bold text-xl mx-auto mb-3">Z</div>
                <p className="text-gray-400 text-sm mb-4">{t('settings.zennDescription')}</p>
              </div>
              <div>
                <label htmlFor="integration-zenn-username" className={labelClass}>{t('settings.zennUsername')}</label>
                <input
                  id="integration-zenn-username"
                  type="text"
                  value={props.zennUsername}
                  onChange={(e) => props.setZennUsername(e.target.value)}
                  placeholder={t('settings.zennUsernamePlaceholder')}
                  className={inputClass}
                />
              </div>
              <button
                type="submit"
                disabled={props.connectingZenn || !props.zennUsername.trim()}
                className="w-full px-5 py-2.5 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg font-semibold text-sm transition-colors"
              >
                {props.connectingZenn ? t('common.loading') : t('settings.connect')}
              </button>
            </form>
          )}
        </div>
      </div>

      {/* Qiita Integration */}
      <div className={sectionContainerClass}>
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.qiita')}</h2>
        </div>
        <div className="p-6">
          {user.qiita_username ? (
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-green-500 rounded-lg flex items-center justify-center text-white font-bold text-sm">Q</div>
                <div>
                  <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
                  <p className="text-sm text-gray-400">@{user.qiita_username}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={props.onSyncQiita}
                  disabled={props.syncingQiita}
                  className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium border border-gray-700 transition-colors"
                >
                  {props.syncingQiita ? t('settings.syncing') : t('settings.sync')}
                </button>
                <button
                  onClick={props.onDisconnectQiita}
                  className={buttonDangerOutlineClass}
                >
                  {t('settings.disconnect')}
                </button>
              </div>
            </div>
          ) : (
            <form onSubmit={props.onConnectQiita} className="space-y-4">
              <div className="text-center py-2">
                <div className="w-12 h-12 bg-green-500/20 rounded-lg flex items-center justify-center text-green-400 font-bold text-xl mx-auto mb-3">Q</div>
                <p className="text-gray-400 text-sm mb-4">{t('settings.qiitaDescription')}</p>
              </div>
              <div>
                <label htmlFor="integration-qiita-username" className={labelClass}>{t('settings.qiitaUsername')}</label>
                <input
                  id="integration-qiita-username"
                  type="text"
                  value={props.qiitaUsername}
                  onChange={(e) => props.setQiitaUsername(e.target.value)}
                  placeholder={t('settings.qiitaUsernamePlaceholder')}
                  className={inputClass}
                />
              </div>
              <button
                type="submit"
                disabled={props.connectingQiita || !props.qiitaUsername.trim()}
                className="w-full px-5 py-2.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-50 text-white rounded-lg font-semibold text-sm transition-colors"
              >
                {props.connectingQiita ? t('common.loading') : t('settings.connect')}
              </button>
            </form>
          )}
        </div>
      </div>

      {/* AtCoder Integration */}
      <div className={sectionContainerClass}>
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.atcoder')}</h2>
        </div>
        <div className="p-6">
          {user.atcoder_username ? (
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 bg-gray-700 rounded-lg flex items-center justify-center text-white font-bold text-sm">A</div>
                <div>
                  <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
                  <p className="text-sm text-gray-400">@{user.atcoder_username}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={props.onDisconnectAtCoder}
                  className={buttonDangerOutlineClass}
                >
                  {t('settings.disconnect')}
                </button>
              </div>
            </div>
          ) : (
            <form onSubmit={props.onConnectAtCoder} className="space-y-4">
              <div className="text-center py-2">
                <div className="w-12 h-12 bg-gray-700/50 rounded-lg flex items-center justify-center text-gray-400 font-bold text-xl mx-auto mb-3">A</div>
                <p className="text-gray-400 text-sm mb-4">{t('settings.atcoderDescription')}</p>
              </div>
              <div>
                <label htmlFor="integration-atcoder-username" className={labelClass}>{t('settings.atcoderUsername')}</label>
                <input
                  id="integration-atcoder-username"
                  type="text"
                  value={props.atcoderUsername}
                  onChange={(e) => props.setAtcoderUsername(e.target.value)}
                  placeholder={t('settings.atcoderUsernamePlaceholder')}
                  className={inputClass}
                />
              </div>
              <button
                type="submit"
                disabled={props.connectingAtcoder || !props.atcoderUsername.trim()}
                className="w-full px-5 py-2.5 bg-cyan-600 hover:bg-cyan-500 disabled:opacity-50 text-white rounded-lg font-semibold text-sm transition-colors"
              >
                {props.connectingAtcoder ? t('common.loading') : t('settings.connect')}
              </button>
            </form>
          )}
        </div>
      </div>

      {/* paiza Rank */}
      <div className={sectionContainerClass}>
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.paiza')}</h2>
        </div>
        <div className="p-6 space-y-4">
          <div className="text-center py-2">
            <div className="w-12 h-12 bg-emerald-700/50 rounded-lg flex items-center justify-center text-emerald-400 font-bold text-xl mx-auto mb-3">P</div>
            <p className="text-gray-400 text-sm mb-4">{t('settings.paizaDescription')}</p>
          </div>
          <div>
            <label htmlFor="integration-paiza-rank" className={labelClass}>{t('settings.paizaRankLabel')}</label>
            <select
              id="integration-paiza-rank"
              value={props.paizaRank}
              onChange={(e) => props.setPaizaRank(e.target.value)}
              className={`${selectClass} w-full`}
            >
              <option value="">{t('settings.paizaSelectRank')}</option>
              <option value="S">S {t('settings.paizaRankS')}</option>
              <option value="A">A {t('settings.paizaRankA')}</option>
              <option value="B">B {t('settings.paizaRankB')}</option>
              <option value="C">C {t('settings.paizaRankC')}</option>
              <option value="D">D {t('settings.paizaRankD')}</option>
              <option value="E">E {t('settings.paizaRankE')}</option>
            </select>
          </div>
          <button
            type="button"
            onClick={props.onSavePaizaRank}
            disabled={props.savingPaiza}
            className="w-full px-5 py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white rounded-lg font-semibold text-sm transition-colors"
          >
            {props.savingPaiza ? t('common.loading') : t('common.save')}
          </button>
        </div>
      </div>
    </>
  );
}
