import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import type { StudyCircle } from '../../types/studyCircle';

interface CircleSearchCardProps {
  circle: StudyCircle;
}

export default function CircleSearchCard({ circle }: CircleSearchCardProps) {
  const { t } = useTranslation();

  return (
    <Link
      to={`/study-circles/${circle.id}`}
      className="block bg-gray-900 border border-gray-800 rounded-md p-5 hover:border-gray-700 transition-colors"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          <h3 className="font-semibold text-white mb-1">{circle.name}</h3>
          <p className="text-sm text-blue-400 mb-2">{circle.topic}</p>
          {circle.description && <p className="text-sm text-gray-400 line-clamp-2">{circle.description}</p>}
        </div>
        <div className="text-right flex-shrink-0">
          <div className="text-sm text-gray-400">
            {circle.member_count || 0} / {circle.max_members || '\u221E'}
          </div>
          <div className="text-xs text-gray-500 mt-1">{t('studyCircles.members')}</div>
        </div>
      </div>
    </Link>
  );
}
