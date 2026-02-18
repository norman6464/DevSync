import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../store/authStore';
import { updateUser } from '../api/users';
import { getGitHubConnectURL } from '../api/github';
import { connectZenn } from '../api/zenn';
import { connectQiita } from '../api/qiita';
import { connectAtCoder } from '../api/atcoder';
import toast from 'react-hot-toast';

export function useOnboarding() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { user, setUser } = useAuthStore();
  const [step, setStep] = useState(1);
  const [saving, setSaving] = useState(false);

  // Step 1: Profile
  const [name, setName] = useState(user?.name || '');
  const [bio, setBio] = useState(user?.bio || '');

  // Step 2: Skills
  const [selectedLanguages, setSelectedLanguages] = useState<string[]>(
    user?.skills_languages ? user.skills_languages.split(',').filter(Boolean) : []
  );
  const [selectedFrameworks, setSelectedFrameworks] = useState<string[]>(
    user?.skills_frameworks ? user.skills_frameworks.split(',').filter(Boolean) : []
  );

  // Step 3: Integrations
  const [zennUsername, setZennUsername] = useState('');
  const [qiitaUsername, setQiitaUsername] = useState('');
  const [atcoderUsername, setAtcoderUsername] = useState('');
  const [connectingZenn, setConnectingZenn] = useState(false);
  const [connectingQiita, setConnectingQiita] = useState(false);
  const [connectingAtcoder, setConnectingAtcoder] = useState(false);
  const [paizaRank, setPaizaRank] = useState('');

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

  const handleSaveProfile = async () => {
    if (!user) return;
    setSaving(true);
    try {
      const { data } = await updateUser(user.id, { name, bio });
      setUser(data);
      setStep(2);
    } catch {
      toast.error(t('settings.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleSaveSkills = async () => {
    if (!user) return;
    setSaving(true);
    try {
      const { data } = await updateUser(user.id, {
        skills_languages: selectedLanguages.join(','),
        skills_frameworks: selectedFrameworks.join(','),
      });
      setUser(data);
      setStep(3);
    } catch {
      toast.error(t('settings.saveFailed'));
    } finally {
      setSaving(false);
    }
  };

  const handleConnectGitHub = async () => {
    try {
      localStorage.setItem('onboarding_redirect', 'true');
      const { data } = await getGitHubConnectURL();
      window.location.href = data.url;
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleConnectZenn = async () => {
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

  const handleConnectQiita = async () => {
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

  const handleConnectAtCoder = async () => {
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

  const handleSavePaizaRank = async () => {
    if (!user || !paizaRank) return;
    try {
      const { data } = await updateUser(user.id, { paiza_rank: paizaRank });
      setUser(data);
      toast.success(t('settings.saved'));
    } catch {
      toast.error(t('settings.saveFailed'));
    }
  };

  const handleComplete = async () => {
    if (!user) return;
    setSaving(true);
    try {
      const { data } = await updateUser(user.id, { onboarding_completed: true });
      setUser(data);
      navigate('/');
    } catch {
      toast.error(t('errors.somethingWrong'));
    } finally {
      setSaving(false);
    }
  };

  return {
    user,
    step,
    setStep,
    saving,
    // Step 1
    name,
    setName,
    bio,
    setBio,
    handleSaveProfile,
    // Step 2
    selectedLanguages,
    selectedFrameworks,
    toggleLanguage,
    toggleFramework,
    handleSaveSkills,
    // Step 3
    zennUsername,
    setZennUsername,
    qiitaUsername,
    setQiitaUsername,
    atcoderUsername,
    setAtcoderUsername,
    connectingZenn,
    connectingQiita,
    connectingAtcoder,
    paizaRank,
    setPaizaRank,
    handleConnectGitHub,
    handleConnectZenn,
    handleConnectQiita,
    handleConnectAtCoder,
    handleSavePaizaRank,
    handleComplete,
  };
}
