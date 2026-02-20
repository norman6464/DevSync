import { useTranslation } from 'react-i18next';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { LANGUAGES, FRAMEWORKS } from '../../constants/skills';

interface OnboardingSkillsStepProps {
  selectedLanguages: string[];
  selectedFrameworks: string[];
  toggleLanguage: (lang: string) => void;
  toggleFramework: (fw: string) => void;
  saving: boolean;
  onSave: () => void;
  onBack: () => void;
  onSkip: () => void;
}

export default function OnboardingSkillsStep({
  selectedLanguages,
  selectedFrameworks,
  toggleLanguage,
  toggleFramework,
  saving,
  onSave,
  onBack,
  onSkip,
}: OnboardingSkillsStepProps) {
  const { t } = useTranslation();

  return (
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
                src={`https://skillicons.dev/icons?${new URLSearchParams({ i: selectedLanguages.join(','), theme: 'dark' })}`}
                alt="Selected languages"
                referrerPolicy="no-referrer"
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
                src={`https://skillicons.dev/icons?${new URLSearchParams({ i: selectedFrameworks.join(','), theme: 'dark' })}`}
                alt="Selected frameworks"
                referrerPolicy="no-referrer"
                className="h-12"
              />
            </div>
          )}
        </div>
      </div>
      <div className="px-6 py-4 border-t border-gray-800 flex justify-between">
        <button
          onClick={onBack}
          className="px-4 py-2 text-gray-400 hover:text-white text-sm transition-colors inline-flex items-center gap-1"
        >
          <ChevronLeft className="w-4 h-4" />
          {t('onboarding.back')}
        </button>
        <div className="flex gap-3">
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
    </div>
  );
}
