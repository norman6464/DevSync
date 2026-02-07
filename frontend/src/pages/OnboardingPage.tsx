import { useState } from 'react';
import { Navigate, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { User, Code, Link, CheckCircle, ChevronRight, ChevronLeft } from 'lucide-react';
import { useAuthStore } from '../store/authStore';
import { updateUser } from '../api/users';
import { getGitHubConnectURL } from '../api/github';
import { connectZenn } from '../api/zenn';
import { connectQiita } from '../api/qiita';
import toast from 'react-hot-toast';

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

const STEPS = [
  { id: 1, icon: User },
  { id: 2, icon: Code },
  { id: 3, icon: Link },
  { id: 4, icon: CheckCircle },
];

export default function OnboardingPage() {
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
  const [connectingZenn, setConnectingZenn] = useState(false);
  const [connectingQiita, setConnectingQiita] = useState(false);

  if (!user) return null;
  if (user.onboarding_completed) return <Navigate to="/" replace />;

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

  const handleConnectQiita = async () => {
    if (!qiitaUsername.trim()) return;
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

  const handleComplete = async () => {
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

  const inputClass = "w-full px-3 py-2 bg-gray-800/50 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow";

  return (
    <div className="min-h-screen bg-gray-950 flex flex-col items-center justify-center px-4 py-12">
      <div className="w-full max-w-xl">
        {/* Logo */}
        <div className="text-center mb-8">
          <svg className="w-12 h-12 mx-auto mb-4 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M16 18l2-2-2-2" />
            <path d="M8 6L6 8l2 2" />
            <path d="M14.5 4l-5 16" />
          </svg>
        </div>

        {/* Step Indicator */}
        <div className="flex items-center justify-center gap-2 mb-8">
          {STEPS.map((s, i) => {
            const StepIcon = s.icon;
            return (
              <div key={s.id} className="flex items-center">
                <div
                  className={`w-10 h-10 rounded-full flex items-center justify-center border-2 transition-colors ${
                    step === s.id
                      ? 'border-blue-500 bg-blue-500/20 text-blue-400'
                      : step > s.id
                      ? 'border-green-500 bg-green-500/20 text-green-400'
                      : 'border-gray-700 bg-gray-800 text-gray-500'
                  }`}
                >
                  <StepIcon className="w-5 h-5" />
                </div>
                {i < STEPS.length - 1 && (
                  <div
                    className={`w-12 h-0.5 mx-1 transition-colors ${
                      step > s.id ? 'bg-green-500' : 'bg-gray-700'
                    }`}
                  />
                )}
              </div>
            );
          })}
        </div>

        {/* Step Content */}
        <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
          {/* Step 1: Profile */}
          {step === 1 && (
            <div>
              <div className="px-6 py-5 border-b border-gray-800">
                <h2 className="text-xl font-semibold text-white">{t('onboarding.welcomeTitle')}</h2>
                <p className="text-sm text-gray-400 mt-1">{t('onboarding.welcomeDescription')}</p>
              </div>
              <div className="p-6 space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-1.5">
                    {t('settings.name')}
                  </label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    placeholder={t('onboarding.namePlaceholder')}
                    className={inputClass}
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-300 mb-1.5">
                    {t('settings.bio')}
                  </label>
                  <textarea
                    value={bio}
                    onChange={(e) => setBio(e.target.value)}
                    rows={3}
                    placeholder={t('onboarding.bioPlaceholder')}
                    className={`${inputClass} resize-none`}
                  />
                </div>
              </div>
              <div className="px-6 py-4 border-t border-gray-800 flex justify-end gap-3">
                <button
                  onClick={() => setStep(2)}
                  className="px-4 py-2 text-gray-400 hover:text-white text-sm transition-colors"
                >
                  {t('onboarding.skip')}
                </button>
                <button
                  onClick={handleSaveProfile}
                  disabled={saving}
                  className="px-5 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors inline-flex items-center gap-1"
                >
                  {saving ? t('common.loading') : t('onboarding.next')}
                  <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          )}

          {/* Step 2: Skills */}
          {step === 2 && (
            <div>
              <div className="px-6 py-5 border-b border-gray-800">
                <h2 className="text-xl font-semibold text-white">{t('onboarding.skillsTitle')}</h2>
                <p className="text-sm text-gray-400 mt-1">{t('onboarding.skillsDescription')}</p>
              </div>
              <div className="p-6 space-y-6">
                {/* Languages */}
                <div>
                  <h3 className="text-sm font-medium text-gray-300 mb-3">{t('profile.languages')}</h3>
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
                  <h3 className="text-sm font-medium text-gray-300 mb-3">{t('profile.frameworks')}</h3>
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
              <div className="px-6 py-4 border-t border-gray-800 flex justify-between">
                <button
                  onClick={() => setStep(1)}
                  className="px-4 py-2 text-gray-400 hover:text-white text-sm transition-colors inline-flex items-center gap-1"
                >
                  <ChevronLeft className="w-4 h-4" />
                  {t('onboarding.back')}
                </button>
                <div className="flex gap-3">
                  <button
                    onClick={() => setStep(3)}
                    className="px-4 py-2 text-gray-400 hover:text-white text-sm transition-colors"
                  >
                    {t('onboarding.skip')}
                  </button>
                  <button
                    onClick={handleSaveSkills}
                    disabled={saving}
                    className="px-5 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors inline-flex items-center gap-1"
                  >
                    {saving ? t('common.loading') : t('onboarding.next')}
                    <ChevronRight className="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          )}

          {/* Step 3: Integrations */}
          {step === 3 && (
            <div>
              <div className="px-6 py-5 border-b border-gray-800">
                <h2 className="text-xl font-semibold text-white">{t('onboarding.integrationsTitle')}</h2>
                <p className="text-sm text-gray-400 mt-1">{t('onboarding.integrationsDescription')}</p>
              </div>
              <div className="p-6 space-y-4">
                {/* GitHub */}
                <div className="p-4 bg-gray-800/50 rounded-lg border border-gray-700">
                  <div className="flex items-center gap-3 mb-3">
                    <svg className="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                    <div>
                      <h3 className="text-sm font-medium text-white">GitHub</h3>
                      <p className="text-xs text-gray-400">{t('onboarding.githubDescription')}</p>
                    </div>
                  </div>
                  {user.github_connected ? (
                    <div className="flex items-center gap-2 text-green-400 text-sm">
                      <CheckCircle className="w-4 h-4" />
                      <span>{t('settings.connected')} - @{user.github_username}</span>
                    </div>
                  ) : (
                    <button
                      onClick={handleConnectGitHub}
                      className="w-full py-2 bg-white hover:bg-gray-100 text-gray-900 rounded-lg font-medium text-sm transition-colors inline-flex items-center justify-center gap-2"
                    >
                      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
                      {t('settings.connect')} GitHub
                    </button>
                  )}
                </div>

                {/* Zenn */}
                <div className="p-4 bg-gray-800/50 rounded-lg border border-gray-700">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="w-6 h-6 bg-blue-500 rounded flex items-center justify-center text-white font-bold text-xs">Z</div>
                    <div>
                      <h3 className="text-sm font-medium text-white">Zenn</h3>
                      <p className="text-xs text-gray-400">{t('onboarding.zennDescription')}</p>
                    </div>
                  </div>
                  {user.zenn_username ? (
                    <div className="flex items-center gap-2 text-green-400 text-sm">
                      <CheckCircle className="w-4 h-4" />
                      <span>{t('settings.connected')} - @{user.zenn_username}</span>
                    </div>
                  ) : (
                    <div className="flex gap-2">
                      <input
                        type="text"
                        value={zennUsername}
                        onChange={(e) => setZennUsername(e.target.value)}
                        placeholder={t('settings.zennUsername')}
                        className={`${inputClass} flex-1`}
                      />
                      <button
                        onClick={handleConnectZenn}
                        disabled={connectingZenn || !zennUsername.trim()}
                        className="px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors whitespace-nowrap"
                      >
                        {connectingZenn ? t('common.loading') : t('settings.connect')}
                      </button>
                    </div>
                  )}
                </div>

                {/* Qiita */}
                <div className="p-4 bg-gray-800/50 rounded-lg border border-gray-700">
                  <div className="flex items-center gap-3 mb-3">
                    <div className="w-6 h-6 bg-green-500 rounded flex items-center justify-center text-white font-bold text-xs">Q</div>
                    <div>
                      <h3 className="text-sm font-medium text-white">Qiita</h3>
                      <p className="text-xs text-gray-400">{t('onboarding.qiitaDescription')}</p>
                    </div>
                  </div>
                  {user.qiita_username ? (
                    <div className="flex items-center gap-2 text-green-400 text-sm">
                      <CheckCircle className="w-4 h-4" />
                      <span>{t('settings.connected')} - @{user.qiita_username}</span>
                    </div>
                  ) : (
                    <div className="flex gap-2">
                      <input
                        type="text"
                        value={qiitaUsername}
                        onChange={(e) => setQiitaUsername(e.target.value)}
                        placeholder={t('settings.qiitaUsername')}
                        className={`${inputClass} flex-1`}
                      />
                      <button
                        onClick={handleConnectQiita}
                        disabled={connectingQiita || !qiitaUsername.trim()}
                        className="px-4 py-2 bg-green-600 hover:bg-green-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors whitespace-nowrap"
                      >
                        {connectingQiita ? t('common.loading') : t('settings.connect')}
                      </button>
                    </div>
                  )}
                </div>
              </div>
              <div className="px-6 py-4 border-t border-gray-800 flex justify-between">
                <button
                  onClick={() => setStep(2)}
                  className="px-4 py-2 text-gray-400 hover:text-white text-sm transition-colors inline-flex items-center gap-1"
                >
                  <ChevronLeft className="w-4 h-4" />
                  {t('onboarding.back')}
                </button>
                <button
                  onClick={() => setStep(4)}
                  className="px-5 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg font-medium text-sm transition-colors inline-flex items-center gap-1"
                >
                  {t('onboarding.next')}
                  <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            </div>
          )}

          {/* Step 4: Complete */}
          {step === 4 && (
            <div className="p-8 text-center">
              <div className="w-16 h-16 bg-green-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
                <CheckCircle className="w-8 h-8 text-green-400" />
              </div>
              <h2 className="text-xl font-semibold text-white mb-2">{t('onboarding.completeTitle')}</h2>
              <p className="text-sm text-gray-400 mb-6">{t('onboarding.completeDescription')}</p>
              <button
                onClick={handleComplete}
                disabled={saving}
                className="px-6 py-2.5 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors"
              >
                {saving ? t('common.loading') : t('onboarding.goToDashboard')}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
