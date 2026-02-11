import { useTranslation } from 'react-i18next';
import { useSettings } from '../hooks/useSettings';
import ProfileSection from './settings/ProfileSection';
import SkillsSection from './settings/SkillsSection';
import IntegrationSection from './settings/IntegrationSection';
import AccountSection from './settings/AccountSection';

export default function SettingsPage() {
  const { t } = useTranslation();
  const s = useSettings();

  if (!s.user) return null;

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">{t('settings.title')}</h1>

      <ProfileSection
        name={s.name} setName={s.setName}
        bio={s.bio} setBio={s.setBio}
        avatarUrl={s.avatarUrl} setAvatarUrl={s.setAvatarUrl}
        saving={s.saving} onSubmit={s.handleSaveProfile}
      />

      <SkillsSection
        selectedLanguages={s.selectedLanguages}
        selectedFrameworks={s.selectedFrameworks}
        savingSkills={s.savingSkills}
        toggleLanguage={s.toggleLanguage}
        toggleFramework={s.toggleFramework}
        onSave={s.handleSaveSkills}
      />

      <AccountSection
        user={s.user}
        emailWeeklyReport={s.emailWeeklyReport} setEmailWeeklyReport={s.setEmailWeeklyReport}
        emailLanguage={s.emailLanguage} setEmailLanguage={s.setEmailLanguage}
        savingEmail={s.savingEmail} onSaveEmailPreferences={s.handleSaveEmailPreferences}
        showDeleteModal={s.showDeleteModal} setShowDeleteModal={s.setShowDeleteModal}
        deletePassword={s.deletePassword} setDeletePassword={s.setDeletePassword}
        deleting={s.deleting} onDeleteAccount={s.handleDeleteAccount}
      />

      <IntegrationSection
        user={s.user}
        syncing={s.syncing}
        onConnectGitHub={s.handleConnectGitHub} onDisconnectGitHub={s.handleDisconnectGitHub} onSyncGitHub={s.handleSyncGitHub}
        zennUsername={s.zennUsername} setZennUsername={s.setZennUsername}
        connectingZenn={s.connectingZenn} syncingZenn={s.syncingZenn}
        onConnectZenn={s.handleConnectZenn} onDisconnectZenn={s.handleDisconnectZenn} onSyncZenn={s.handleSyncZenn}
        qiitaUsername={s.qiitaUsername} setQiitaUsername={s.setQiitaUsername}
        connectingQiita={s.connectingQiita} syncingQiita={s.syncingQiita}
        onConnectQiita={s.handleConnectQiita} onDisconnectQiita={s.handleDisconnectQiita} onSyncQiita={s.handleSyncQiita}
        atcoderUsername={s.atcoderUsername} setAtcoderUsername={s.setAtcoderUsername}
        connectingAtcoder={s.connectingAtcoder}
        onConnectAtCoder={s.handleConnectAtCoder} onDisconnectAtCoder={s.handleDisconnectAtCoder}
        paizaRank={s.paizaRank} setPaizaRank={s.setPaizaRank}
        savingPaiza={s.savingPaiza} onSavePaizaRank={s.handleSavePaizaRank}
      />
    </div>
  );
}
