import { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Send, Loader2, Plus, MessageSquare } from 'lucide-react';
import type { AIMessage } from '../../api/advice';

interface AIChatPanelProps {
  messages: AIMessage[];
  sending: boolean;
  dailyChatRemaining: number;
  onSend: (message: string) => void;
  onNewChat: () => void;
}

export default function AIChatPanel({
  messages,
  sending,
  dailyChatRemaining,
  onSend,
  onNewChat,
}: AIChatPanelProps) {
  const { t } = useTranslation();
  const [input, setInput] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() || sending || dailyChatRemaining <= 0) return;
    onSend(input.trim());
    setInput('');
  };

  return (
    <div className="bg-gray-800/50 rounded-xl border border-gray-700 flex flex-col h-[600px]">
      {/* ヘッダー */}
      <div className="flex items-center justify-between p-4 border-b border-gray-700">
        <div className="flex items-center gap-2">
          <MessageSquare size={18} className="text-blue-400" />
          <h3 className="text-white font-medium">{t('advice.chatTitle')}</h3>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-xs text-gray-400">
            {t('advice.remainingToday', { count: dailyChatRemaining })}
          </span>
          <button
            onClick={onNewChat}
            className="text-xs text-blue-400 hover:text-blue-300 flex items-center gap-1"
          >
            <Plus size={14} />
            {t('advice.newChat')}
          </button>
        </div>
      </div>

      {/* メッセージ一覧 */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.length === 0 && (
          <div className="text-center text-gray-500 mt-20">
            <MessageSquare size={40} className="mx-auto mb-3 opacity-40" />
            <p className="text-sm">{t('advice.chatPlaceholder')}</p>
          </div>
        )}
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
          >
            <div
              className={`max-w-[80%] rounded-lg px-4 py-2 text-sm ${
                msg.role === 'user'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-700 text-gray-200'
              }`}
            >
              <p className="whitespace-pre-wrap">{msg.content}</p>
            </div>
          </div>
        ))}
        {sending && (
          <div className="flex justify-start">
            <div className="bg-gray-700 rounded-lg px-4 py-2">
              <Loader2 size={16} className="animate-spin text-gray-400" />
            </div>
          </div>
        )}
        <div ref={messagesEndRef} />
      </div>

      {/* 入力フォーム */}
      <form onSubmit={handleSubmit} className="p-4 border-t border-gray-700">
        {dailyChatRemaining <= 0 ? (
          <p className="text-center text-yellow-400 text-sm">
            {t('advice.rateLimitReached')}
          </p>
        ) : (
          <div className="flex gap-2">
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder={t('advice.chatPlaceholder')}
              disabled={sending}
              className="flex-1 bg-gray-700 border border-gray-600 rounded-lg px-4 py-2 text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            />
            <button
              type="submit"
              disabled={sending || !input.trim()}
              className="bg-blue-600 hover:bg-blue-700 text-white rounded-lg px-4 py-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {sending ? (
                <Loader2 size={16} className="animate-spin" />
              ) : (
                <Send size={16} />
              )}
            </button>
          </div>
        )}
      </form>
    </div>
  );
}
