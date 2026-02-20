import { useTranslation } from 'react-i18next';
import { FileText } from 'lucide-react';
import { type ZennArticle, type ZennStats } from '../../api/zenn';
import { type QiitaArticle, type QiitaStats } from '../../api/qiita';

interface ProfileArticlesSectionProps {
  zennUsername?: string;
  zennArticles: ZennArticle[];
  zennStats: ZennStats | null;
  qiitaUsername?: string;
  qiitaArticles: QiitaArticle[];
  qiitaStats: QiitaStats | null;
}

export default function ProfileArticlesSection({
  zennUsername,
  zennArticles,
  zennStats,
  qiitaUsername,
  qiitaArticles,
  qiitaStats,
}: ProfileArticlesSectionProps) {
  const { t } = useTranslation();

  return (
    <>
      {/* Zenn Articles */}
      {zennUsername && zennArticles.length > 0 && (
        <div>
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide flex items-center gap-2">
              <span className="w-5 h-5 bg-blue-500 rounded text-white text-xs flex items-center justify-center font-bold">Z</span>
              {t('profile.zennArticles')}
              {zennStats && <span className="text-xs text-gray-500 font-normal ml-2">{zennStats.total_articles} {t('profile.articles')} · {zennStats.total_likes} {t('post.like')}s</span>}
            </h2>
            <a href={`https://zenn.dev/${encodeURIComponent(zennUsername)}`} target="_blank" rel="noopener noreferrer" className="text-xs text-gray-500 hover:text-blue-400 transition-colors flex items-center gap-1">{t('profile.viewAllOnZenn')}<svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" /></svg></a>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {zennArticles.slice(0, 6).map((article) => (
              <a key={article.id} href={`https://zenn.dev/${encodeURIComponent(zennUsername)}/articles/${encodeURIComponent(article.slug)}`} target="_blank" rel="noopener noreferrer" className="bg-gray-900 border border-gray-800 rounded-md p-4 hover:border-gray-600 transition-colors group">
                <div className="flex items-start gap-3">
                  <span className="text-2xl">{article.emoji || '📝'}</span>
                  <div className="min-w-0 flex-1">
                    <div className="font-medium text-sm text-blue-400 group-hover:text-blue-300 line-clamp-2">{article.title}</div>
                    <div className="flex items-center gap-3 mt-2">
                      <span className="flex items-center gap-1 text-xs text-gray-400"><svg className="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/></svg>{article.liked_count}</span>
                      <span className="flex items-center gap-1 text-xs text-gray-400"><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" /></svg>{article.comments_count}</span>
                      <span className="px-2 py-0.5 bg-gray-800 text-gray-400 text-xs rounded">{article.article_type === 'tech' ? 'Tech' : 'Idea'}</span>
                    </div>
                  </div>
                </div>
              </a>
            ))}
          </div>
        </div>
      )}

      {/* Qiita Articles */}
      {qiitaUsername && qiitaArticles.length > 0 && (
        <div>
          <div className="flex items-center justify-between mb-3">
            <h2 className="text-sm font-semibold text-gray-300 uppercase tracking-wide flex items-center gap-2">
              <span className="w-5 h-5 bg-green-500 rounded text-white text-xs flex items-center justify-center font-bold">Q</span>
              {t('profile.qiitaArticles')}
              {qiitaStats && <span className="text-xs text-gray-500 font-normal ml-2">{qiitaStats.total_articles} {t('profile.articles')} · {qiitaStats.total_likes} {t('post.like')}s</span>}
            </h2>
            <a href={`https://qiita.com/${encodeURIComponent(qiitaUsername)}`} target="_blank" rel="noopener noreferrer" className="text-xs text-gray-500 hover:text-green-400 transition-colors flex items-center gap-1">{t('profile.viewAllOnQiita')}<svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" /></svg></a>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {qiitaArticles.slice(0, 6).map((article) => (
              <a key={article.id} href={article.url} target="_blank" rel="noopener noreferrer" className="bg-gray-900 border border-gray-800 rounded-md p-4 hover:border-gray-600 transition-colors group">
                <div className="flex items-start gap-3">
                  <FileText className="w-6 h-6 text-green-400 flex-shrink-0" aria-hidden="true" />
                  <div className="min-w-0 flex-1">
                    <div className="font-medium text-sm text-green-400 group-hover:text-green-300 line-clamp-2">{article.title}</div>
                    <div className="flex items-center gap-3 mt-2">
                      <span className="flex items-center gap-1 text-xs text-gray-400"><svg className="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/></svg>{article.likes_count}</span>
                      <span className="flex items-center gap-1 text-xs text-gray-400"><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" d="M7.5 8.25h9m-9 3H12m-9.75 1.51c0 1.6 1.123 2.994 2.707 3.227 1.129.166 2.27.293 3.423.379.35.026.67.21.865.501L12 21l2.755-4.133a1.14 1.14 0 0 1 .865-.501 48.172 48.172 0 0 0 3.423-.379c1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" /></svg>{article.comments_count}</span>
                      {article.tags && <span className="px-2 py-0.5 bg-gray-800 text-gray-400 text-xs rounded truncate max-w-[100px]">{article.tags.split(',')[0]}</span>}
                    </div>
                  </div>
                </div>
              </a>
            ))}
          </div>
        </div>
      )}
    </>
  );
}
