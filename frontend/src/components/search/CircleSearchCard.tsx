import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Users } from 'lucide-react';
import type { StudyCircle } from '../../types/studyCircle';
import { badgeBaseClass } from '../../constants/styles';

const statusColor: Record<string, string> = {
  active: 'bg-green-400/10 text-green-400',
  completed: 'bg-blue-400/10 text-blue-400',
  archived: 'bg-gray-400/10 text-gray-400',
};

interface CircleSearchCardProps {
  circle: StudyCircle & { member_count?: number };
}

export default function CircleSearchCard({ circle }: CircleSearchCardProps) {
  const { t } = useTranslation();
  const isFull = circle.max_members > 0 && (circle.member_count || 0) >= circle.max_members;

  return (
    <Link
      to={`/study-circles/${circle.id}`}
      className="block bg-gray-900 border border-gray-800 rounded-md p-5 hover:border-gray-700 transition-colors"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="flex-1">
          <div className="flex items-center gap-2 flex-wrap mb-1">
            <h3 className="font-semibold text-white">{circle.name}</h3>
            {circle.status && (
              <span className={`${badgeBaseClass} ${statusColor[circle.status] || statusColor.active}`}>
                {t(`studyCircle.${circle.status}`)}
              </span>
            )}
            {isFull && (
              <span className="inline-flex items-center gap-0.5 px-1.5 py-0.5 text-xs rounded bg-red-400/10 text-red-400 font-medium">
                <Users className="w-3 h-3" />
                {t('studyCircle.full')}
              </span>
            )}
          </div>
          <p className="text-sm text-blue-400 mb-2">{circle.topic}</p>
          {circle.description && <p className="text-sm text-gray-400 line-clamp-2">{circle.description}</p>}
        </div>
        <div className="text-right flex-shrink-0">
          <div className="text-sm text-gray-400">
            {circle.member_count || 0} / {circle.max_members || '\u221E'}
          </div>
          <div className="text-xs text-gray-500 mt-1">{t('studyCircle.members')}</div>
        </div>
      </div>
    </Link>
  );
}
