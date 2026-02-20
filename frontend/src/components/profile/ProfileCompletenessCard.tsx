import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

interface ProfileCompletenessCardProps {
  percentage: number;
  missingFields: string[];
}

export default function ProfileCompletenessCard({ percentage, missingFields }: ProfileCompletenessCardProps) {
  const { t } = useTranslation();

  if (percentage >= 100) return null;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-5">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-semibold flex items-center gap-2">
          <svg className="w-4 h-4 text-yellow-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24" aria-hidden="true">
            <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75 11.25 15 15 9.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
          </svg>
          {t('profile.completeness')}
        </h3>
        <span className="text-sm font-bold text-yellow-400">{percentage}%</span>
      </div>
      <div className="h-2 bg-gray-800 rounded-full overflow-hidden mb-3" role="progressbar" aria-valuenow={percentage} aria-valuemin={0} aria-valuemax={100} aria-label={`${t('profile.completeness')}: ${percentage}%`}>
        <div
          className="h-full bg-gradient-to-r from-yellow-500 to-green-500 rounded-full transition-all duration-500"
          style={{ width: `${percentage}%` }}
        />
      </div>
      <div className="flex flex-wrap gap-2">
        {missingFields.map((field) => (
          <Link
            key={field}
            to="/settings"
            className="px-2.5 py-1 bg-gray-800 hover:bg-gray-700 text-gray-400 hover:text-white text-xs rounded-lg transition-colors"
          >
            + {t(`profile.missing.${field}`)}
          </Link>
        ))}
      </div>
    </div>
  );
}
