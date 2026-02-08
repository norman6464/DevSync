import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Sparkles, CheckCircle2 } from 'lucide-react';
import { getDailyChallenge, isChallengeCompleted, markChallengeCompleted } from '../../utils/dailyChallenge';

export default function DailyChallengeWidget() {
  const { t } = useTranslation();
  const [completed, setCompleted] = useState(isChallengeCompleted);
  const challengeKey = getDailyChallenge();

  const handleComplete = () => {
    markChallengeCompleted();
    setCompleted(true);
  };

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
      <h3 className="flex items-center gap-2 text-sm font-medium text-white mb-3">
        <Sparkles className="w-4 h-4 text-yellow-400" />
        {t('dailyChallenge.title')}
      </h3>

      <div className={`rounded-lg p-3 ${completed ? 'bg-green-900/20 border border-green-800/30' : 'bg-gray-800/50'}`}>
        <p className={`text-sm mb-3 ${completed ? 'text-green-300 line-through' : 'text-gray-300'}`}>
          {t(`dailyChallenge.challenges.${challengeKey}`)}
        </p>

        {completed ? (
          <div className="flex items-center gap-2 text-green-400">
            <CheckCircle2 className="w-4 h-4" />
            <span className="text-xs font-medium">{t('dailyChallenge.completed')}</span>
          </div>
        ) : (
          <button
            onClick={handleComplete}
            className="w-full py-1.5 px-3 text-xs font-medium text-white bg-yellow-600 hover:bg-yellow-500 rounded-lg transition-colors"
          >
            {t('dailyChallenge.markComplete')}
          </button>
        )}
      </div>
    </div>
  );
}
