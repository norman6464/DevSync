import { useTranslation } from 'react-i18next';
import { Trophy, CheckCircle2 } from 'lucide-react';
import { useWeeklyChallenge } from '../../hooks/useWeeklyChallenge';
import { panelClass } from '../../constants/styles';

export default function WeeklyChallengeWidget() {
  const { t } = useTranslation();
  const { challenge, loading } = useWeeklyChallenge();

  if (loading || !challenge) return null;

  const progress = Math.min(
    Math.round((challenge.current_value / challenge.target_value) * 100),
    100
  );

  return (
    <div className={panelClass}>
      <div className="flex items-center justify-between mb-3">
        <h3 className="flex items-center gap-2 text-sm font-medium text-white">
          <Trophy className="w-4 h-4 text-yellow-400" />
          {t('weeklyChallenge.title')}
        </h3>
        {challenge.is_completed && (
          <span className="flex items-center gap-1 text-xs text-green-400">
            <CheckCircle2 className="w-3.5 h-3.5" />
            {t('weeklyChallenge.completed')}
          </span>
        )}
      </div>

      <div className={`rounded-lg p-3 ${
        challenge.is_completed
          ? 'bg-green-900/20 border border-green-800/30'
          : 'bg-gray-800/50'
      }`}>
        <p className={`text-sm mb-3 ${
          challenge.is_completed ? 'text-green-300' : 'text-gray-300'
        }`}>
          {t(`weeklyChallenge.challenges.${challenge.description}`)}
        </p>

        {/* Progress bar */}
        <div className="mb-2">
          <div className="flex justify-between text-xs text-gray-500 mb-1">
            <span>{challenge.current_value} / {challenge.target_value}</span>
            <span>{progress}%</span>
          </div>
          <div className="w-full h-2 bg-gray-700 rounded-full overflow-hidden">
            <div
              className={`h-full rounded-full transition-all duration-500 ${
                challenge.is_completed
                  ? 'bg-green-500'
                  : progress >= 75
                    ? 'bg-yellow-500'
                    : 'bg-blue-500'
              }`}
              style={{ width: `${progress}%` }}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
