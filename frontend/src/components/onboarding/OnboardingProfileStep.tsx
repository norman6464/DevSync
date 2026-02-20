import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { ChevronRight } from 'lucide-react';
import { inputClass, labelClass } from '../../constants/styles';

interface OnboardingProfileStepProps {
  name: string;
  setName: (name: string) => void;
  bio: string;
  setBio: (bio: string) => void;
  saving: boolean;
  onSave: () => void;
  onSkip: () => void;
}

export default function OnboardingProfileStep({
  name,
  setName,
  bio,
  setBio,
  saving,
  onSave,
  onSkip,
}: OnboardingProfileStepProps) {
  const { t } = useTranslation();
  const handleNameChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setName(e.target.value), [setName]);
  const handleBioChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => setBio(e.target.value), [setBio]);

  return (
    <div>
      <div className="px-6 py-5 border-b border-gray-800">
        <h2 className="text-xl font-semibold text-white">{t('onboarding.welcomeTitle')}</h2>
        <p className="text-sm text-gray-400 mt-1">{t('onboarding.welcomeDescription')}</p>
      </div>
      <div className="p-6 space-y-4">
        <div>
          <label className={labelClass}>{t('settings.name')}</label>
          <input
            type="text"
            autoComplete="name"
            value={name}
            onChange={handleNameChange}
            placeholder={t('onboarding.namePlaceholder')}
            maxLength={200}
            className={inputClass}
          />
        </div>
        <div>
          <label className={labelClass}>{t('settings.bio')}</label>
          <textarea
            value={bio}
            onChange={handleBioChange}
            rows={3}
            maxLength={500}
            placeholder={t('onboarding.bioPlaceholder')}
            className={`${inputClass} resize-none`}
          />
        </div>
      </div>
      <div className="px-6 py-4 border-t border-gray-800 flex justify-end gap-3">
        <button
          onClick={onSkip}
          className="px-4 py-2 text-gray-400 hover:text-white text-sm transition-colors"
        >
          {t('onboarding.skip')}
        </button>
        <button
          onClick={onSave}
          disabled={saving}
          className="px-5 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors inline-flex items-center gap-1"
        >
          {saving ? t('common.loading') : t('onboarding.next')}
          <ChevronRight className="w-4 h-4" />
        </button>
      </div>
    </div>
  );
}
