import { useTranslation } from 'react-i18next';
import type { User } from '../../types/user';
import GitHubIntegrationCard from './integrations/GitHubIntegrationCard';
import SpotifyIntegrationCard from './integrations/SpotifyIntegrationCard';
import UsernameIntegrationCard from './integrations/UsernameIntegrationCard';
import PaizaRankCard from './integrations/PaizaRankCard';

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

export default function IntegrationSection(props: Props) {
  const { t } = useTranslation();
  const { user } = props;

  return (
    <>
      <GitHubIntegrationCard
        connected={user.github_connected}
        username={user.github_username}
        syncing={props.syncing}
        onConnect={props.onConnectGitHub}
        onDisconnect={props.onDisconnectGitHub}
        onSync={props.onSyncGitHub}
      />

      <SpotifyIntegrationCard
        connected={user.spotify_connected}
        onConnect={props.onConnectSpotify}
        onDisconnect={props.onDisconnectSpotify}
      />

      <UsernameIntegrationCard
        title={t('settings.zenn')}
        description={t('settings.zennDescription')}
        badgeLetter="Z"
        badgeBgClass="bg-blue-500"
        badgeTextClass="text-white"
        badgeBgDisabledClass="bg-blue-500/20"
        badgeTextDisabledClass="text-blue-400"
        connectBtnClass="bg-blue-600 hover:bg-blue-500"
        connected={!!user.zenn_username}
        connectedUsername={user.zenn_username}
        username={props.zennUsername}
        setUsername={props.setZennUsername}
        usernameLabel={t('settings.zennUsername')}
        usernamePlaceholder={t('settings.zennUsernamePlaceholder')}
        connecting={props.connectingZenn}
        onConnect={props.onConnectZenn}
        onDisconnect={props.onDisconnectZenn}
        syncing={props.syncingZenn}
        onSync={props.onSyncZenn}
      />

      <UsernameIntegrationCard
        title={t('settings.qiita')}
        description={t('settings.qiitaDescription')}
        badgeLetter="Q"
        badgeBgClass="bg-green-500"
        badgeTextClass="text-white"
        badgeBgDisabledClass="bg-green-500/20"
        badgeTextDisabledClass="text-green-400"
        connectBtnClass="bg-gray-700 hover:bg-gray-600"
        connected={!!user.qiita_username}
        connectedUsername={user.qiita_username}
        username={props.qiitaUsername}
        setUsername={props.setQiitaUsername}
        usernameLabel={t('settings.qiitaUsername')}
        usernamePlaceholder={t('settings.qiitaUsernamePlaceholder')}
        connecting={props.connectingQiita}
        onConnect={props.onConnectQiita}
        onDisconnect={props.onDisconnectQiita}
        syncing={props.syncingQiita}
        onSync={props.onSyncQiita}
      />

      <UsernameIntegrationCard
        title={t('settings.atcoder')}
        description={t('settings.atcoderDescription')}
        badgeLetter="A"
        badgeBgClass="bg-gray-700"
        badgeTextClass="text-white"
        badgeBgDisabledClass="bg-gray-700/50"
        badgeTextDisabledClass="text-gray-400"
        connectBtnClass="bg-cyan-600 hover:bg-cyan-500"
        connected={!!user.atcoder_username}
        connectedUsername={user.atcoder_username}
        username={props.atcoderUsername}
        setUsername={props.setAtcoderUsername}
        usernameLabel={t('settings.atcoderUsername')}
        usernamePlaceholder={t('settings.atcoderUsernamePlaceholder')}
        connecting={props.connectingAtcoder}
        onConnect={props.onConnectAtCoder}
        onDisconnect={props.onDisconnectAtCoder}
      />

      <PaizaRankCard
        paizaRank={props.paizaRank}
        setPaizaRank={props.setPaizaRank}
        saving={props.savingPaiza}
        onSave={props.onSavePaizaRank}
      />
    </>
  );
}
