import { useTranslation } from 'react-i18next';
import { Sparkles } from 'lucide-react';
import { buttonSecondaryClass } from '../../constants/styles';
import { LANGUAGES, FRAMEWORKS } from '../../constants/skills';

interface Props {
  selectedLanguages: string[];
  selectedFrameworks: string[];
  savingSkills: boolean;
  toggleLanguage: (lang: string) => void;
  toggleFramework: (fw: string) => void;
  onSave: () => void;
}

export default function SkillsSection({
  selectedLanguages, selectedFrameworks, savingSkills,
  toggleLanguage, toggleFramework, onSave,
}: Props) {
  const { t } = useTranslation();

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-800">
        <h2 className="text-base font-semibold flex items-center gap-2">
          <Sparkles className="w-5 h-5 text-yellow-400" /> {t('settings.skills')}
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
          onClick={onSave}
          disabled={savingSkills}
          className={`${buttonSecondaryClass} px-5 font-medium text-sm disabled:opacity-50`}
        >
          {savingSkills ? t('common.loading') : t('common.save')}
        </button>
      </div>
    </div>
  );
}
