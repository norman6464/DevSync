import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { Lightbulb, ChevronRight, MessageSquare } from 'lucide-react';
import { useAdvice } from '../../hooks/useAdvice';
import AdviceIcon from '../advice/AdviceIcon';

export default function AIAdviceWidget() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { advices, llmAvailable, loading } = useAdvice();

  // 上位3件の未読アドバイスを表示
  const topAdvices = advices.filter((a) => !a.is_read).slice(0, 3);

  if (loading) {
    return (
      <div className="bg-gray-800/50 rounded-xl border border-gray-700 p-4">
        <div className="animate-pulse space-y-3">
          <div className="h-4 bg-gray-700 rounded w-1/2" />
          <div className="h-3 bg-gray-700 rounded w-full" />
          <div className="h-3 bg-gray-700 rounded w-3/4" />
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-800/50 rounded-xl border border-gray-700 p-4">
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-2">
          <Lightbulb size={18} className="text-yellow-400" />
          <h3 className="text-white font-medium text-sm">{t('advice.title')}</h3>
        </div>
        <button
          onClick={() => navigate('/advice')}
          className="text-xs text-blue-400 hover:text-blue-300 flex items-center gap-1"
        >
          {t('advice.viewAll')}
          <ChevronRight size={14} />
        </button>
      </div>

      {topAdvices.length === 0 ? (
        <p className="text-gray-500 text-xs">{t('advice.noAdvice')}</p>
      ) : (
        <div className="space-y-2">
          {topAdvices.map((advice) => {
            const params = advice.params ? JSON.parse(advice.params) : {};
            return (
              <div
                key={advice.id}
                className="flex items-start gap-2 p-2 rounded-lg hover:bg-gray-700/30 cursor-pointer transition-colors"
                onClick={() => {
                  if (advice.action_url) navigate(advice.action_url);
                }}
              >
                <AdviceIcon type={advice.type} size={16} className="mt-0.5" />
                <div className="min-w-0">
                  <p className="text-white text-xs font-medium truncate">
                    {t(advice.title_key, params)}
                  </p>
                  <p className="text-gray-500 text-xs truncate">
                    {t(advice.message_key, params)}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      )}

      {llmAvailable && (
        <button
          onClick={() => navigate('/advice')}
          className="w-full mt-3 flex items-center justify-center gap-2 py-2 bg-blue-600/20 hover:bg-blue-600/30 text-blue-400 rounded-lg text-xs transition-colors border border-blue-600/20"
        >
          <MessageSquare size={14} />
          {t('advice.askAI')}
        </button>
      )}
    </div>
  );
}
