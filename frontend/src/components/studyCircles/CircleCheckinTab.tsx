import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { StudyCircleCheckin } from '../../types/studyCircle';
import Avatar from '../common/Avatar';

interface CircleCheckinTabProps {
  checkins: StudyCircleCheckin[];
  onCheckin: (content: string) => Promise<unknown>;
}

export default function CircleCheckinTab({ checkins, onCheckin }: CircleCheckinTabProps) {
  const { t } = useTranslation();
  const [checkinContent, setCheckinContent] = useState('');

  const handleCheckin = async () => {
    if (!checkinContent.trim()) return;
    const result = await onCheckin(checkinContent);
    if (result) {
      setCheckinContent('');
    }
  };

  return (
    <div className="space-y-4">
      {/* Checkin Form */}
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
        <h3 className="text-sm font-medium text-white mb-2">{t('studyCircle.checkin.title')}</h3>
        <div className="flex gap-2">
          <input
            type="text"
            value={checkinContent}
            onChange={(e) => setCheckinContent(e.target.value)}
            placeholder={t('studyCircle.checkin.placeholder')}
            className="flex-1 bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-purple-500"
            onKeyDown={(e) => e.key === 'Enter' && handleCheckin()}
          />
          <button
            onClick={handleCheckin}
            disabled={!checkinContent.trim()}
            className="px-4 py-2 bg-purple-600 hover:bg-purple-500 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors shrink-0"
          >
            {t('studyCircle.checkin.submit')}
          </button>
        </div>
      </div>

      {/* Checkin History */}
      <div className="bg-gray-900 border border-gray-800 rounded-md p-4">
        <h3 className="text-sm font-medium text-white mb-3">{t('studyCircle.checkin.history')}</h3>
        {checkins.length === 0 ? (
          <p className="text-xs text-gray-500 text-center py-4">{t('studyCircle.checkin.noCheckins')}</p>
        ) : (
          <div className="space-y-2">
            {checkins.map((ci) => (
              <div key={ci.id} className="flex items-start gap-2.5 p-2 rounded-lg hover:bg-gray-800/50">
                <Avatar name={ci.user?.name || ''} avatarUrl={ci.user?.avatar_url} size="xs" />
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-medium text-white">{ci.user?.name}</span>
                    <span className="text-[10px] text-gray-600">{ci.date}</span>
                  </div>
                  <p className="text-xs text-gray-400 mt-0.5">{ci.content}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
