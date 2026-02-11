import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { updateUser } from '../api/users';
import { getGitHubConnectURL, disconnectGitHub, syncGitHub } from '../api/github';
import { connectZenn, disconnectZenn, syncZenn } from '../api/zenn';
import { connectQiita, disconnectQiita, syncQiita } from '../api/qiita';
import { connectAtCoder, disconnectAtCoder } from '../api/atcoder';
import { deleteAccount } from '../api/auth';
import { getEmailPreferences, updateEmailPreferences } from '../api/emailPreferences';
import toast from 'react-hot-toast';
import { useTranslation } from 'react-i18next';

export function useSettings() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { user, setUser, logout } = useAuthStore();

  // Profile
  const [name, setName] = useState(user?.name || '');
  const [bio, setBio] = useState(user?.bio || '');
  const [avatarUrl, setAvatarUrl] = useState(user?.avatar_url || '');
  const [saving, setSaving] = useState(false);

  // Skills
  const [selectedLanguages, setSelectedLanguages] = useState<string[]>(
    user?.skills_languages ? user.skills_languages.split(',').filter(Boolean) : []
  );
  const [selectedFrameworks, setSelectedFrameworks] = useState<string[]>(
    user?.skills_frameworks ? user.skills_frameworks.split(',').filter(Boolean) : []
  );
  const [savingSkills, setSavingSkills] = useState(false);

  // Email
  const [emailWeeklyReport, setEmailWeeklyReport] = useState(user?.email_weekly_report ?? true);
  const [emailLanguage, setEmailLanguage] = useState(user?.email_language || 'ja');
  const [savingEmail, setSavingEmail] = useState(false);

  // Integrations
  const [syncing, setSyncing] = useState(false);
  const [zennUsername, setZennUsername] = useState('');
  const [connectingZenn, setConnectingZenn] = useState(false);
  const [syncingZenn, setSyncingZenn] = useState(false);
  const [qiitaUsername, setQiitaUsername] = useState('');
  const [connectingQiita, setConnectingQiita] = useState(false);
  const [syncingQiita, setSyncingQiita] = useState(false);
  const [atcoderUsername, setAtcoderUsername] = useState('');
  const [connectingAtcoder, setConnectingAtcoder] = useState(false);
  const [paizaRank, setPaizaRank] = useState(user?.paiza_rank || '');
  const [savingPaiza, setSavingPaiza] = useState(false);

  // Account deletion
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletePassword, setDeletePassword] = useState('');
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    const loadEmailPreferences = async () => {
      try {
        const { data } = await getEmailPreferences();
        setEmailWeeklyReport(data.email_weekly_report);
        setEmailLanguage(data.email_language);
      } catch {
        // デフォルト値を使用
      }
    };
    loadEmailPreferences();
  }, []);

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;
    setSaving(true);
    try {
      const { data } = await updateUser(user.id, { name, bio, avatar_url: avatarUrl });
      setUser(data);
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('settings.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleSaveSkills = async () => {
    if (!user) return;
    setSavingSkills(true);
    try {
      const { data } = await updateUser(user.id, {
        skills_languages: selectedLanguages.join(','),
        skills_frameworks: selectedFrameworks.join(','),
      });
      setUser(data);
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('settings.saveFailed'));
    } finally {
      setSavingSkills(false);
    }
  };

  const handleSaveEmailPreferences = async () => {
    setSavingEmail(true);
    try {
      await updateEmailPreferences({
        email_weekly_report: emailWeeklyReport,
        email_language: emailLanguage,
      });
      toast.success(t('settings.emailPreferencesSaved'));
    } catch {
      toast.error(t('settings.saveFailed'));
    } finally {
      setSavingEmail(false);
    }
  };

  const toggleLanguage = (lang: string) => {
    setSelectedLanguages((prev) =>
      prev.includes(lang) ? prev.filter((l) => l !== lang) : [...prev, lang]
    );
  };

  const toggleFramework = (fw: string) => {
    setSelectedFrameworks((prev) =>
      prev.includes(fw) ? prev.filter((f) => f !== fw) : [...prev, fw]
    );
  };

  const handleConnectGitHub = async () => {
    try {
      const { data } = await getGitHubConnectURL();
      window.location.href = data.url;
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleDisconnectGitHub = async () => {
    if (!user) return;
    try {
      await disconnectGitHub();
      setUser({ ...user, github_connected: false, github_username: '' });
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleSyncGitHub = async () => {
    setSyncing(true);
    try {
      await syncGitHub();
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    } finally {
      setSyncing(false);
    }
  };

  const handleConnectZenn = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || !zennUsername.trim()) return;
    setConnectingZenn(true);
    try {
      await connectZenn(zennUsername.trim());
      setUser({ ...user, zenn_username: zennUsername.trim() });
      setZennUsername('');
      toast.success(t('settings.zennConnected'));
    } catch {
      toast.error(t('settings.zennInvalidUsername'));
    } finally {
      setConnectingZenn(false);
    }
  };

  const handleDisconnectZenn = async () => {
    if (!user) return;
    try {
      await disconnectZenn();
      setUser({ ...user, zenn_username: '' });
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleSyncZenn = async () => {
    setSyncingZenn(true);
    try {
      await syncZenn();
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    } finally {
      setSyncingZenn(false);
    }
  };

  const handleConnectQiita = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user || !qiitaUsername.trim()) return;
    setConnectingQiita(true);
    try {
      await connectQiita(qiitaUsername.trim());
      setUser({ ...user, qiita_username: qiitaUsername.trim() });
      setQiitaUsername('');
      toast.success(t('settings.qiitaConnected'));
    } catch {
      toast.error(t('settings.qiitaInvalidUsername'));
    } finally {
      setConnectingQiita(false);
    }
  };

  const handleDisconnectQiita = async () => {
    if (!user) return;
    try {
      await disconnectQiita();
      setUser({ ...user, qiita_username: '' });
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleSyncQiita = async () => {
    setSyncingQiita(true);
    try {
      await syncQiita();
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    } finally {
      setSyncingQiita(false);
    }
  };

  const handleConnectAtCoder = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!atcoderUsername.trim()) return;
    setConnectingAtcoder(true);
    try {
      const { data } = await connectAtCoder(atcoderUsername.trim());
      setUser(data);
      setAtcoderUsername('');
      toast.success(t('settings.atcoderConnected'));
    } catch {
      toast.error(t('settings.atcoderInvalidUsername'));
    } finally {
      setConnectingAtcoder(false);
    }
  };

  const handleDisconnectAtCoder = async () => {
    try {
      const { data } = await disconnectAtCoder();
      setUser(data);
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleSavePaizaRank = async () => {
    if (!user) return;
    setSavingPaiza(true);
    try {
      const { data } = await updateUser(user.id, { paiza_rank: paizaRank });
      setUser(data);
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('settings.saveFailed'));
    } finally {
      setSavingPaiza(false);
    }
  };

  const handleDeleteAccount = async () => {
    setDeleting(true);
    try {
      await deleteAccount(deletePassword || undefined);
      await logout();
      navigate('/login');
      toast.success(t('accountManagement.accountDeleted'));
    } catch {
      toast.error(t('accountManagement.deleteFailed'));
    } finally {
      setDeleting(false);
    }
  };

  return {
    user,
    // Profile
    name, setName, bio, setBio, avatarUrl, setAvatarUrl, saving, handleSaveProfile,
    // Skills
    selectedLanguages, selectedFrameworks, savingSkills, toggleLanguage, toggleFramework, handleSaveSkills,
    // Email
    emailWeeklyReport, setEmailWeeklyReport, emailLanguage, setEmailLanguage, savingEmail, handleSaveEmailPreferences,
    // GitHub
    syncing, handleConnectGitHub, handleDisconnectGitHub, handleSyncGitHub,
    // Zenn
    zennUsername, setZennUsername, connectingZenn, syncingZenn, handleConnectZenn, handleDisconnectZenn, handleSyncZenn,
    // Qiita
    qiitaUsername, setQiitaUsername, connectingQiita, syncingQiita, handleConnectQiita, handleDisconnectQiita, handleSyncQiita,
    // AtCoder
    atcoderUsername, setAtcoderUsername, connectingAtcoder, handleConnectAtCoder, handleDisconnectAtCoder,
    // Paiza
    paizaRank, setPaizaRank, savingPaiza, handleSavePaizaRank,
    // Account deletion
    showDeleteModal, setShowDeleteModal, deletePassword, setDeletePassword, deleting, handleDeleteAccount,
  };
}
