import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Users, ChevronRight } from 'lucide-react';
import { useStudyCircles } from '../../hooks';
import Avatar from '../common/Avatar';

export default function StudyCircleWidget() {
  const { t } = useTranslation();
  const { circles, loading } = useStudyCircles();

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
      <div className="flex items-center justify-between mb-3">
        <h3 className="flex items-center gap-2 text-sm font-medium text-white">
          <Users className="w-4 h-4 text-purple-400" />
          {t('studyCircle.widget.title')}
        </h3>
        <Link to="/study-circles" className="text-xs text-gray-400 hover:text-purple-400 transition-colors">
          {t('studyCircle.widget.viewAll')}
        </Link>
      </div>

      {loading ? (
        <div className="space-y-2">
          {[1, 2].map((i) => (
            <div key={i} className="h-14 bg-gray-800 rounded-lg animate-pulse" />
          ))}
        </div>
      ) : circles.length === 0 ? (
        <div className="text-center py-4">
          <p className="text-xs text-gray-500 mb-2">{t('studyCircle.widget.noCircles')}</p>
          <Link
            to="/study-circles"
            className="text-xs text-purple-400 hover:text-purple-300 transition-colors"
          >
            {t('studyCircle.widget.viewAll')}
          </Link>
        </div>
      ) : (
        <div className="space-y-2">
          {circles.slice(0, 2).map((circle) => (
            <Link
              key={circle.id}
              to={`/study-circles/${circle.id}`}
              className="block p-2.5 rounded-lg hover:bg-gray-800/50 transition-colors"
            >
              <div className="flex items-center justify-between mb-1">
                <span className="text-xs font-medium text-white truncate">{circle.name}</span>
                <ChevronRight className="w-3 h-3 text-gray-600 shrink-0" />
              </div>
              <div className="flex items-center justify-between">
                <span className="text-[10px] text-purple-400 truncate">{circle.topic}</span>
                <div className="flex items-center gap-0.5 shrink-0 ml-2">
                  {circle.members?.slice(0, 3).map((m) => (
                    <Avatar key={m.id} name={m.user?.name || ''} avatarUrl={m.user?.avatar_url} size="xs" />
                  ))}
                </div>
              </div>
            </Link>
          ))}
          {circles.length > 2 && (
            <Link
              to="/study-circles"
              className="flex items-center justify-center gap-1 text-xs text-gray-400 hover:text-purple-400 pt-1 transition-colors"
            >
              {t('studyCircle.widget.viewAll')}
              <ChevronRight className="w-3 h-3" />
            </Link>
          )}
        </div>
      )}
    </div>
  );
}
