import { useTranslation } from 'react-i18next';
import { CheckCircle } from 'lucide-react';

interface OnboardingCompleteStepProps {
  saving: boolean;
  onComplete: () => void;
}

export default function OnboardingCompleteStep({
  saving,
  onComplete,
}: OnboardingCompleteStepProps) {
  const { t } = useTranslation();

  return (
    <div className="p-8 text-center">
      <div className="w-16 h-16 bg-green-500/20 rounded-full flex items-center justify-center mx-auto mb-4">
        <CheckCircle className="w-8 h-8 text-green-400" />
      </div>
      <h2 className="text-xl font-semibold text-white mb-2">{t('onboarding.completeTitle')}</h2>
      <p className="text-sm text-gray-400 mb-6">{t('onboarding.completeDescription')}</p>
      <button
        onClick={onComplete}
        disabled={saving}
        className="px-6 py-2.5 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg font-medium text-sm transition-colors"
      >
        {saving ? t('common.loading') : t('onboarding.goToDashboard')}
      </button>
    </div>
  );
}
