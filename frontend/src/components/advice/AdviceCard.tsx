import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { Check } from 'lucide-react';
import AdviceIcon from './AdviceIcon';
import type { AIAdvice } from '../../api/advice';
import { parseJsonObject } from '../../utils/json';
import { linkSmallClass } from '../../constants/styles';

interface AdviceCardProps {
  advice: AIAdvice;
  onMarkRead: (id: number) => void;
}

export default function AdviceCard({ advice, onMarkRead }: AdviceCardProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();

  const params = parseJsonObject(advice.params);

  const priorityBorder: Record<number, string> = {
    1: 'border-l-red-500',
    2: 'border-l-orange-500',
    3: 'border-l-blue-500',
    4: 'border-l-green-500',
    5: 'border-l-gray-500',
  };

  return (
    <div
      className={`bg-gray-800/50 rounded-lg border-l-4 ${priorityBorder[advice.priority] || 'border-l-gray-500'} p-4 ${
        advice.is_read ? 'opacity-60' : ''
      }`}
    >
      <div className="flex items-start gap-3">
        <AdviceIcon type={advice.type} size={24} className="mt-0.5 flex-shrink-0" />
        <div className="flex-1 min-w-0">
          <h4 className="text-white font-medium text-sm">
            {t(advice.title_key, params)}
          </h4>
          <p className="text-gray-400 text-xs mt-1">
            {t(advice.message_key, params)}
          </p>
          <div className="flex items-center gap-2 mt-3">
            {advice.action_url?.startsWith('/') && (
              <button
                onClick={() => navigate(advice.action_url)}
                className={linkSmallClass}
              >
                {t('advice.viewAll')} &rarr;
              </button>
            )}
            {!advice.is_read && (
              <button
                onClick={() => onMarkRead(advice.id)}
                className="ml-auto text-xs text-gray-500 hover:text-gray-300 transition-colors flex items-center gap-1"
              >
                <Check size={12} />
                {t('advice.markRead')}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
