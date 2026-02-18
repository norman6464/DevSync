import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Send, FileText, Check, Clock } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { useAutoSave } from '../../hooks/useAutoSave';
import Avatar from '../common/Avatar';

interface QuickPostFormProps {
  onSubmit: (title: string, content: string, isDraft?: boolean) => Promise<void>;
}

export default function QuickPostForm({ onSubmit }: QuickPostFormProps) {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const [content, setContent] = useState('');
  const [isFocused, setIsFocused] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showRestorePrompt, setShowRestorePrompt] = useState(false);

  // 自動保存
  const { saveStatus, lastSaved, clearSaved, getSaved } = useAutoSave({
    key: 'devsync_draft_post',
    data: content.trim() ? { content } : null,
    delay: 3000,
  });

  // 初回マウント時に保存データがあれば復元を提案
  useEffect(() => {
    const saved = getSaved();
    if (saved && saved.content && !content) {
      setShowRestorePrompt(true);
    }
  }, []);

  const handleRestore = () => {
    const saved = getSaved();
    if (saved && saved.content) {
      setContent(saved.content);
      setIsFocused(true);
    }
    setShowRestorePrompt(false);
  };

  const handleDismissRestore = () => {
    clearSaved();
    setShowRestorePrompt(false);
  };

  const handleSubmit = async (isDraft = false) => {
    if (!content.trim()) return;

    setIsSubmitting(true);
    try {
      // タイトルは本文の最初の50文字から自動生成
      const title = content.trim().substring(0, 50) + (content.length > 50 ? '...' : '');
      await onSubmit(title, content, isDraft);
      setContent('');
      setIsFocused(false);
      // 送信成功後、保存データをクリア
      clearSaved();
    } finally {
      setIsSubmitting(false);
    }
  };

  // 保存状態の表示テキスト
  const getSaveStatusText = () => {
    if (saveStatus === 'saving') return t('post.saving');
    if (saveStatus === 'saved' && lastSaved) {
      const seconds = Math.floor((Date.now() - lastSaved.getTime()) / 1000);
      if (seconds < 60) return t('post.savedJustNow');
      return t('post.savedMinutesAgo', { minutes: Math.floor(seconds / 60) });
    }
    return '';
  };

  return (
    <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-800 rounded-xl p-4 transition-all">
      {/* 復元プロンプト */}
      {showRestorePrompt && (
        <div role="alert" className="mb-3 p-3 bg-blue-900/20 border border-blue-800/30 rounded-lg flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Clock className="w-4 h-4 text-blue-400" aria-hidden="true" />
            <p className="text-sm text-blue-300">{t('post.restorePrompt')}</p>
          </div>
          <div className="flex gap-2">
            <button
              onClick={handleDismissRestore}
              className="px-2 py-1 text-xs text-gray-400 hover:text-white transition-colors"
            >
              {t('common.dismiss')}
            </button>
            <button
              onClick={handleRestore}
              className="px-3 py-1 text-xs bg-blue-600 hover:bg-blue-700 text-white rounded transition-colors"
            >
              {t('post.restore')}
            </button>
          </div>
        </div>
      )}

      <div className="flex gap-3">
        <Avatar name={user?.name || ''} avatarUrl={user?.avatar_url} size="sm" />
        <div className="flex-1">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            onFocus={() => setIsFocused(true)}
            placeholder={t('post.createTitle')}
            className={`w-full bg-gray-800/50 border border-gray-700 rounded-lg px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent resize-none transition-all ${
              isFocused ? 'min-h-[120px]' : 'min-h-[60px]'
            }`}
            disabled={isSubmitting}
          />

          {isFocused && (
            <div className="mt-3 flex items-center justify-between animate-in fade-in duration-200">
              <div className="flex items-center gap-2">
                <p className="text-xs text-gray-500">
                  {t('post.markdownSupported')}
                </p>
                {/* 保存状態の表示 */}
                {saveStatus !== 'idle' && (
                  <div className="flex items-center gap-1 text-xs text-gray-400">
                    {saveStatus === 'saving' ? (
                      <>
                        <Clock className="w-3 h-3 animate-pulse" aria-hidden="true" />
                        <span>{getSaveStatusText()}</span>
                      </>
                    ) : (
                      <>
                        <Check className="w-3 h-3 text-green-500" aria-hidden="true" />
                        <span>{getSaveStatusText()}</span>
                      </>
                    )}
                  </div>
                )}
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => setIsFocused(false)}
                  disabled={isSubmitting}
                  className="px-3 py-1.5 text-sm text-gray-400 hover:text-white transition-colors disabled:opacity-50"
                >
                  {t('common.cancel')}
                </button>
                <button
                  onClick={() => handleSubmit(true)}
                  disabled={!content.trim() || isSubmitting}
                  className="flex items-center gap-1.5 px-3 py-1.5 text-sm bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <FileText className="w-4 h-4" />
                  {t('post.saveDraft')}
                </button>
                <button
                  onClick={() => handleSubmit(false)}
                  disabled={!content.trim() || isSubmitting}
                  className="flex items-center gap-1.5 px-4 py-1.5 text-sm bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <Send className="w-4 h-4" />
                  {isSubmitting ? t('post.posting') : t('post.post')}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
