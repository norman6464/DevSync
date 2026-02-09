import { useTranslation } from 'react-i18next';
import { MessageSquare, Clock } from 'lucide-react';
import type { AIConversation } from '../../api/advice';

interface ConversationListProps {
  conversations: AIConversation[];
  loading: boolean;
  onSelect: (id: number) => void;
  activeId?: number | null;
}

export default function ConversationList({
  conversations,
  loading,
  onSelect,
  activeId,
}: ConversationListProps) {
  const { t } = useTranslation();

  if (loading) {
    return (
      <div className="text-center text-gray-500 py-8">
        <div className="animate-pulse">{t('common.loading')}...</div>
      </div>
    );
  }

  if (conversations.length === 0) {
    return (
      <div className="text-center text-gray-500 py-8">
        <MessageSquare size={32} className="mx-auto mb-2 opacity-40" />
        <p className="text-sm">{t('advice.noConversations')}</p>
      </div>
    );
  }

  return (
    <div className="space-y-2">
      <h4 className="text-gray-400 text-xs font-medium uppercase tracking-wider mb-3">
        {t('advice.conversations')}
      </h4>
      {conversations.map((conv) => (
        <button
          key={conv.id}
          onClick={() => onSelect(conv.id)}
          className={`w-full text-left p-3 rounded-lg transition-colors ${
            activeId === conv.id
              ? 'bg-blue-600/20 border border-blue-500/30'
              : 'bg-gray-800/50 hover:bg-gray-700/50 border border-transparent'
          }`}
        >
          <p className="text-white text-sm truncate">{conv.title}</p>
          <div className="flex items-center gap-1 mt-1">
            <Clock size={12} className="text-gray-500" />
            <span className="text-xs text-gray-500">
              {new Date(conv.updated_at).toLocaleDateString()}
            </span>
            {conv.messages && (
              <span className="text-xs text-gray-500 ml-2">
                {conv.messages.length} {t('advice.messagesCount', { count: conv.messages.length })}
              </span>
            )}
          </div>
        </button>
      ))}
    </div>
  );
}
