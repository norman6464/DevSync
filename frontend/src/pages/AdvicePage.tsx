import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Lightbulb, MessageSquare } from 'lucide-react';
import { useAdvice, useAIChat, useConversations } from '../hooks/useAdvice';
import AdviceCard from '../components/advice/AdviceCard';
import AIChatPanel from '../components/advice/AIChatPanel';
import ConversationList from '../components/advice/ConversationList';

export default function AdvicePage() {
  const { t } = useTranslation();
  const { advices, llmAvailable, dailyChatRemaining, loading, markRead } = useAdvice();
  const { messages, sending, conversationId, sendMessage, loadConversation, startNewChat, deleteCurrentConversation } = useAIChat();
  const { conversations, loading: convsLoading, refetch: refetchConvs, removeConversation } = useConversations();
  const [activeTab, setActiveTab] = useState<'advice' | 'chat'>('advice');

  const handleSend = useCallback(async (message: string) => {
    await sendMessage(message);
    refetchConvs();
  }, [sendMessage, refetchConvs]);

  const handleSelectConversation = useCallback((id: number) => {
    loadConversation(id);
    setActiveTab('chat');
  }, [loadConversation]);

  const handleDeleteConversation = useCallback(async (id: number) => {
    const success = await removeConversation(id);
    if (success && conversationId === id) {
      startNewChat();
    }
  }, [removeConversation, conversationId, startNewChat]);

  const handleDeleteCurrentConversation = useCallback(async () => {
    const success = await deleteCurrentConversation();
    if (success) {
      refetchConvs();
    }
  }, [deleteCurrentConversation, refetchConvs]);

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      {/* ヘッダー */}
      <div className="flex items-center gap-3 mb-6">
        <Lightbulb size={28} className="text-yellow-400" />
        <div>
          <h1 className="text-2xl font-bold text-white">{t('advice.title')}</h1>
          <p className="text-gray-400 text-sm">{t('advice.subtitle')}</p>
        </div>
      </div>

      {/* モバイルタブ切替 */}
      {llmAvailable && (
        <div className="flex gap-2 mb-6 lg:hidden">
          <button
            onClick={() => setActiveTab('advice')}
            className={`flex-1 py-2 px-4 rounded-lg text-sm font-medium transition-colors ${
              activeTab === 'advice'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-800 text-gray-400 hover:text-white'
            }`}
          >
            <Lightbulb size={16} className="inline mr-1" />
            {t('advice.title')}
          </button>
          <button
            onClick={() => setActiveTab('chat')}
            className={`flex-1 py-2 px-4 rounded-lg text-sm font-medium transition-colors ${
              activeTab === 'chat'
                ? 'bg-blue-600 text-white'
                : 'bg-gray-800 text-gray-400 hover:text-white'
            }`}
          >
            <MessageSquare size={16} className="inline mr-1" />
            {t('advice.askAI')}
          </button>
        </div>
      )}

      <div className="flex gap-6">
        {/* 左カラム: アドバイス一覧 */}
        <div
          className={`flex-1 ${
            llmAvailable && activeTab === 'chat' ? 'hidden lg:block' : ''
          }`}
        >
          {loading ? (
            <div className="space-y-3">
              {[1, 2, 3].map((i) => (
                <div key={i} className="bg-gray-800/50 rounded-lg p-4 animate-pulse">
                  <div className="h-4 bg-gray-800 rounded w-1/2 mb-2" />
                  <div className="h-3 bg-gray-800 rounded w-full" />
                </div>
              ))}
            </div>
          ) : advices.length === 0 ? (
            <div className="text-center py-16 text-gray-500">
              <Lightbulb size={48} className="mx-auto mb-4 opacity-40" />
              <p>{t('advice.noAdvice')}</p>
            </div>
          ) : (
            <div className="space-y-3">
              {advices.map((advice) => (
                <AdviceCard
                  key={advice.id || `${advice.type}-${advice.title_key}`}
                  advice={advice}
                  onMarkRead={markRead}
                />
              ))}
            </div>
          )}
        </div>

        {/* 右カラム: AIチャット + 会話履歴 */}
        {llmAvailable && (
          <div
            className={`w-full lg:w-[450px] flex-shrink-0 space-y-4 ${
              activeTab === 'advice' ? 'hidden lg:block' : ''
            }`}
          >
            <AIChatPanel
              messages={messages}
              sending={sending}
              dailyChatRemaining={dailyChatRemaining}
              conversationId={conversationId}
              onSend={handleSend}
              onNewChat={startNewChat}
              onDelete={handleDeleteCurrentConversation}
            />
            <ConversationList
              conversations={conversations}
              loading={convsLoading}
              onSelect={handleSelectConversation}
              onDelete={handleDeleteConversation}
              activeId={conversationId}
            />
          </div>
        )}
      </div>
    </div>
  );
}
