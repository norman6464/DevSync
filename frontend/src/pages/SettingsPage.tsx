import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../store/authStore';
import { updateUser } from '../api/users';
import { getGitHubConnectURL, disconnectGitHub, syncGitHub } from '../api/github';
import { connectZenn, disconnectZenn, syncZenn } from '../api/zenn';
import { deleteAccount } from '../api/auth';
import toast from 'react-hot-toast';

// skillicons.dev supported icons
const LANGUAGES = [
  'java', 'typescript', 'javascript', 'python', 'go', 'rust', 'cpp', 'c', 'cs',
  'php', 'ruby', 'swift', 'kotlin', 'scala', 'elixir', 'haskell', 'lua', 'perl',
  'r', 'dart', 'html', 'css', 'sass', 'bash', 'powershell', 'sql', 'graphql'
];

const FRAMEWORKS = [
  'react', 'vue', 'angular', 'svelte', 'nextjs', 'nuxtjs', 'gatsby', 'astro',
  'spring', 'django', 'flask', 'fastapi', 'rails', 'laravel', 'express', 'nestjs',
  'gin', 'fiber', 'actix', 'rocket', 'tailwind', 'bootstrap', 'materialui',
  'prisma', 'graphql', 'apollo', 'redux', 'nodejs', 'deno', 'bun'
];

export default function SettingsPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { user, setUser, logout } = useAuthStore();
  const [name, setName] = useState(user?.name || '');
  const [bio, setBio] = useState(user?.bio || '');
  const [avatarUrl, setAvatarUrl] = useState(user?.avatar_url || '');
  const [selectedLanguages, setSelectedLanguages] = useState<string[]>(
    user?.skills_languages ? user.skills_languages.split(',').filter(Boolean) : []
  );
  const [selectedFrameworks, setSelectedFrameworks] = useState<string[]>(
    user?.skills_frameworks ? user.skills_frameworks.split(',').filter(Boolean) : []
  );
  const [saving, setSaving] = useState(false);
  const [savingSkills, setSavingSkills] = useState(false);
  const [syncing, setSyncing] = useState(false);
  const [zennUsername, setZennUsername] = useState('');
  const [connectingZenn, setConnectingZenn] = useState(false);
  const [syncingZenn, setSyncingZenn] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [deletePassword, setDeletePassword] = useState('');
  const [deleting, setDeleting] = useState(false);

  if (!user) return null;

  const handleSaveProfile = async (e: React.FormEvent) => {
    e.preventDefault();
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

  const handleDeleteAccount = async () => {
    setDeleting(true);
    try {
      await deleteAccount(deletePassword || undefined);
      logout();
      navigate('/login');
      toast.success(t('accountManagement.accountDeleted'));
    } catch {
      toast.error(t('accountManagement.deleteFailed'));
    } finally {
      setDeleting(false);
    }
  };

  const handleConnectZenn = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!zennUsername.trim()) return;
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

  const inputClass = "w-full px-3 py-2 bg-gray-800/50 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow";

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">{t('settings.title')}</h1>

      {/* Profile Section */}
      <form onSubmit={handleSaveProfile} className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.profile')}</h2>
        </div>
        <div className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1.5">{t('settings.name')}</label>
            <input type="text" value={name} onChange={(e) => setName(e.target.value)} className={inputClass} />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1.5">{t('settings.bio')}</label>
            <textarea value={bio} onChange={(e) => setBio(e.target.value)} rows={3} placeholder="Tell us about yourself" className={`${inputClass} resize-none`} />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1.5">{t('settings.avatar')}</label>
            <input type="text" value={avatarUrl} onChange={(e) => setAvatarUrl(e.target.value)} placeholder="https://..." className={inputClass} />
          </div>
        </div>
        <div className="px-6 py-4 border-t border-gray-800 flex justify-end">
          <button
            type="submit"
            disabled={saving}
            className="px-5 py-2 bg-green-600 hover:bg-green-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors"
          >
            {saving ? t('common.loading') : t('common.save')}
          </button>
        </div>
      </form>

      {/* Skills Section */}
      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold flex items-center gap-2">
            <span>✨</span> {t('settings.skills')}
          </h2>
          <p className="text-xs text-gray-500 mt-1">{t('settings.selectLanguages')}</p>
        </div>
        <div className="p-6 space-y-6">
          {/* Languages */}
          <div>
            <h3 className="text-sm font-medium text-gray-300 mb-3 flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M17.25 6.75 22.5 12l-5.25 5.25m-10.5 0L1.5 12l5.25-5.25m7.5-3-4.5 16.5" />
              </svg>
              {t('profile.languages')}
            </h3>
            <div className="flex flex-wrap gap-2">
              {LANGUAGES.map((lang) => (
                <button
                  key={lang}
                  type="button"
                  onClick={() => toggleLanguage(lang)}
                  className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all border ${
                    selectedLanguages.includes(lang)
                      ? 'bg-blue-600/20 text-blue-300 border-blue-500/50'
                      : 'bg-gray-800 text-gray-400 border-gray-700 hover:border-gray-500'
                  }`}
                >
                  {lang}
                </button>
              ))}
            </div>
            {selectedLanguages.length > 0 && (
              <div className="mt-3 p-3 bg-gray-800/50 rounded-lg">
                <p className="text-xs text-gray-500 mb-2">{t('settings.preview')}:</p>
                <img
                  src={`https://skillicons.dev/icons?i=${selectedLanguages.join(',')}&theme=dark`}
                  alt="Selected languages"
                  className="h-12"
                />
              </div>
            )}
          </div>

          {/* Frameworks */}
          <div>
            <h3 className="text-sm font-medium text-gray-300 mb-3 flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="M15.59 14.37a6 6 0 0 1-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 0 0 6.16-12.12A14.98 14.98 0 0 0 9.631 8.41m5.96 5.96a14.926 14.926 0 0 1-5.841 2.58m-.119-8.54a6 6 0 0 0-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 0 0-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 0 1-2.448-2.448 14.9 14.9 0 0 1 .06-.312m-2.24 2.39a4.493 4.493 0 0 0-1.757 4.306 4.493 4.493 0 0 0 4.306-1.758M16.5 9a1.5 1.5 0 1 1-3 0 1.5 1.5 0 0 1 3 0Z" />
              </svg>
              {t('profile.frameworks')}
            </h3>
            <div className="flex flex-wrap gap-2">
              {FRAMEWORKS.map((fw) => (
                <button
                  key={fw}
                  type="button"
                  onClick={() => toggleFramework(fw)}
                  className={`px-3 py-1.5 rounded-lg text-xs font-medium transition-all border ${
                    selectedFrameworks.includes(fw)
                      ? 'bg-purple-600/20 text-purple-300 border-purple-500/50'
                      : 'bg-gray-800 text-gray-400 border-gray-700 hover:border-gray-500'
                  }`}
                >
                  {fw}
                </button>
              ))}
            </div>
            {selectedFrameworks.length > 0 && (
              <div className="mt-3 p-3 bg-gray-800/50 rounded-lg">
                <p className="text-xs text-gray-500 mb-2">{t('settings.preview')}:</p>
                <img
                  src={`https://skillicons.dev/icons?i=${selectedFrameworks.join(',')}&theme=dark`}
                  alt="Selected frameworks"
                  className="h-12"
                />
              </div>
            )}
          </div>
        </div>
        <div className="px-6 py-4 border-t border-gray-800 flex justify-end">
          <button
            type="button"
            onClick={handleSaveSkills}
            disabled={savingSkills}
            className="px-5 py-2 bg-green-600 hover:bg-green-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors"
          >
            {savingSkills ? t('common.loading') : t('common.save')}
          </button>
        </div>
      </div>

      {/* GitHub Integration */}
      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-800">
          <h2 className="text-base font-semibold">{t('settings.github')}</h2>
        </div>
        <div className="p-6">
          {user.github_connected ? (
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <svg className="w-8 h-8 text-gray-300" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                <div>
                  <p className="text-sm font-medium text-green-400">{t('settings.connected')}</p>
                  <p className="text-sm text-gray-400">@{user.github_username}</p>
                </div>
              </div>
              <div className="flex gap-2">
                <button
                  onClick={handleSyncGitHub}
                  disabled={syncing}
                  className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium border border-gray-700 transition-colors"
                >
                  {syncing ? t('settings.syncing') : t('settings.sync')}
                </button>
                <button
                  onClick={handleDisconnectGitHub}
                  className="px-4 py-2 text-red-400 hover:text-red-300 border border-red-400/30 hover:border-red-400/50 rounded-lg text-sm font-medium transition-colors"
                >
                  {t('settings.disconnect')}
                </button>
              </div>
            </div>
          ) : (
            <div className="text-center py-4">
              <svg className="w-12 h-12 text-gray-600 mx-auto mb-3" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
              <p className="text-gray-400 text-sm mb-4">Connect your GitHub account to sync contributions, repos, and languages.</p>
              <button
                onClick={handleConnectGitHub}
                className="px-5 py-2.5 bg-white hover:bg-gray-100 text-gray-900 rounded-lg font-semibold text-sm transition-colors inline-flex items-center gap-2"
              >
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                {t('settings.connect')} GitHub
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Zenn Integration */}
      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
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
                  onClick={handleSyncZenn}
                  disabled={syncingZenn}
                  className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium border border-gray-700 transition-colors"
                >
                  {syncingZenn ? t('settings.syncing') : t('settings.sync')}
                </button>
                <button
                  onClick={handleDisconnectZenn}
                  className="px-4 py-2 text-red-400 hover:text-red-300 border border-red-400/30 hover:border-red-400/50 rounded-lg text-sm font-medium transition-colors"
                >
                  {t('settings.disconnect')}
                </button>
              </div>
            </div>
          ) : (
            <form onSubmit={handleConnectZenn} className="space-y-4">
              <div className="text-center py-2">
                <div className="w-12 h-12 bg-blue-500/20 rounded-lg flex items-center justify-center text-blue-400 font-bold text-xl mx-auto mb-3">Z</div>
                <p className="text-gray-400 text-sm mb-4">{t('settings.zennDescription')}</p>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-300 mb-1.5">{t('settings.zennUsername')}</label>
                <input
                  type="text"
                  value={zennUsername}
                  onChange={(e) => setZennUsername(e.target.value)}
                  placeholder="your-zenn-username"
                  className={inputClass}
                />
              </div>
              <button
                type="submit"
                disabled={connectingZenn || !zennUsername.trim()}
                className="w-full px-5 py-2.5 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg font-semibold text-sm transition-colors"
              >
                {connectingZenn ? t('common.loading') : t('settings.connect')}
              </button>
            </form>
          )}
        </div>
      </div>

      {/* Danger Zone - Account Deletion */}
      <div className="bg-gray-900 border border-red-500/30 rounded-xl overflow-hidden">
        <div className="px-6 py-4 border-b border-red-500/30 bg-red-500/5">
          <h2 className="text-base font-semibold text-red-400">{t('accountManagement.dangerZone')}</h2>
        </div>
        <div className="p-6">
          <p className="text-gray-400 text-sm mb-4">
            {t('accountManagement.deleteWarning')}
          </p>
          <button
            onClick={() => setShowDeleteModal(true)}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm font-medium transition-colors"
          >
            {t('accountManagement.deleteAccount')}
          </button>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      {showDeleteModal && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 px-4">
          <div className="bg-gray-900 border border-gray-700 rounded-xl max-w-md w-full p-6">
            <h3 className="text-lg font-semibold text-white mb-2">
              {t('accountManagement.confirmDelete')}
            </h3>
            <p className="text-gray-400 text-sm mb-4">
              {t('accountManagement.deleteConfirmText')}
            </p>

            {/* Show password field for non-GitHub-only users */}
            {!user.github_connected && (
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-300 mb-1.5">
                  {t('auth.password')}
                </label>
                <input
                  type="password"
                  value={deletePassword}
                  onChange={(e) => setDeletePassword(e.target.value)}
                  className={inputClass}
                  placeholder="••••••••"
                />
              </div>
            )}

            <div className="flex gap-3 justify-end">
              <button
                onClick={() => {
                  setShowDeleteModal(false);
                  setDeletePassword('');
                }}
                className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-sm font-medium transition-colors"
              >
                {t('common.cancel')}
              </button>
              <button
                onClick={handleDeleteAccount}
                disabled={deleting}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors"
              >
                {deleting ? t('common.loading') : t('accountManagement.deleteAccount')}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
