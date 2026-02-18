import { Link } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeSanitize from 'rehype-sanitize';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { useTranslation } from 'react-i18next';
import { likePost, unlikePost, bookmarkPost, unbookmarkPost } from '../../api/posts';
import type { Post } from '../../types/post';
import Avatar from '../common/Avatar';
import { format } from 'date-fns';
import { useState } from 'react';
import { Code2 } from 'lucide-react';
import { cardDarkClass } from '../../constants/styles';

interface PostCardProps {
  post: Post;
  isOwner?: boolean;
  onEdit?: () => void;
  onDelete?: () => void;
  onUpdate?: () => void;
}

export default function PostCard({ post, isOwner = false, onEdit, onDelete, onUpdate }: PostCardProps) {
  const { t } = useTranslation();
  const [liked, setLiked] = useState(post.liked || false);
  const [likeCount, setLikeCount] = useState(post.like_count);
  const [bookmarked, setBookmarked] = useState(post.bookmarked || false);

  const handleLike = async () => {
    try {
      if (liked) {
        await unlikePost(post.id);
        setLiked(false);
        setLikeCount((c) => c - 1);
      } else {
        await likePost(post.id);
        setLiked(true);
        setLikeCount((c) => c + 1);
      }
      onUpdate?.();
    } catch (e) {
      console.warn('Failed to toggle like:', e);
    }
  };

  const handleBookmark = async () => {
    try {
      if (bookmarked) {
        await unbookmarkPost(post.id);
        setBookmarked(false);
      } else {
        await bookmarkPost(post.id);
        setBookmarked(true);
      }
      onUpdate?.();
    } catch (e) {
      console.warn('Failed to toggle bookmark:', e);
    }
  };

  // Parse image URLs
  let imageUrls: string[] = [];
  try {
    if (post.image_urls) {
      imageUrls = JSON.parse(post.image_urls);
    }
  } catch {
    // ignore parse error
  }

  return (
    <div className={cardDarkClass}>
      <div className="flex items-center gap-3 mb-3">
        <Link to={`/profile/${post.user_id}`}>
          <Avatar name={post.user?.name || 'U'} avatarUrl={post.user?.avatar_url} size="sm" />
        </Link>
        <div className="min-w-0">
          <Link to={`/profile/${post.user_id}`} className="font-medium text-sm hover:text-blue-400 transition-colors">
            {post.user?.name}
          </Link>
          <p className="text-xs text-gray-500">
            {format(new Date(post.created_at), 'MMM d, yyyy')}
          </p>
        </div>
      </div>

      <div className="flex items-start justify-between gap-2 mb-2">
        <Link to={`/posts/${post.id}`} className="flex-1 group">
          <h3 className="text-base font-semibold group-hover:text-blue-400 transition-colors">{post.title}</h3>
        </Link>

        {isOwner && (
          <div className="flex gap-1 flex-shrink-0">
            <button
              onClick={onEdit}
              className="p-1.5 text-gray-400 hover:text-white transition-colors"
              title={t('common.edit')}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
              </svg>
            </button>
            <button
              onClick={onDelete}
              className="p-1.5 text-gray-400 hover:text-red-400 transition-colors"
              title={t('common.delete')}
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
              </svg>
            </button>
          </div>
        )}
      </div>

      <Link to={`/posts/${post.id}`} className="block group">
        <div className="text-gray-400 text-sm leading-relaxed prose prose-sm prose-invert max-w-none line-clamp-4">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            rehypePlugins={[rehypeSanitize]}
            components={{
              img: () => null,
              p: ({ children }) => <p className="mb-2">{children}</p>,
              a: ({ children }) => (
                <span className="text-blue-400 hover:underline">{children}</span>
              ),
              code({ className, children, ...props }) {
                const match = /language-(\w+)/.exec(className || '');
                const inline = !match;
                return !inline ? (
                  <SyntaxHighlighter
                    style={vscDarkPlus}
                    language={match[1]}
                    PreTag="div"
                    customStyle={{ borderRadius: '0.5rem', fontSize: '0.75rem', maxHeight: '200px' }}
                  >
                    {String(children).replace(/\n$/, '')}
                  </SyntaxHighlighter>
                ) : (
                  <code className={className} {...props}>
                    {children}
                  </code>
                );
              },
            }}
          >
            {post.content}
          </ReactMarkdown>
        </div>
      </Link>

      {/* Image preview */}
      {imageUrls.length > 0 && (
        <div className="mt-3 flex gap-2 overflow-hidden">
          {imageUrls.slice(0, 4).map((url, i) => (
            <div key={i} className="relative">
              <img
                src={url}
                alt=""
                className="w-20 h-20 object-cover rounded-lg border border-gray-700"
              />
              {i === 3 && imageUrls.length > 4 && (
                <div className="absolute inset-0 bg-black/60 rounded-lg flex items-center justify-center text-white text-sm font-medium">
                  +{imageUrls.length - 4}
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Code Snippets Preview */}
      {post.code_snippets && post.code_snippets.length > 0 && (
        <Link to={`/posts/${post.id}`} className="mt-3 block">
          <div className="border border-gray-700 rounded-lg overflow-hidden bg-gray-800/50">
            {post.code_snippets.slice(0, 2).map((snippet) => (
              <div key={snippet.id} className="border-b border-gray-700 last:border-b-0">
                <div className="flex items-center gap-2 px-3 py-1.5 bg-gray-800">
                  <Code2 className="w-3.5 h-3.5 text-gray-400" />
                  <span className="text-xs text-blue-400 font-mono">{snippet.language}</span>
                  {snippet.file_name && (
                    <span className="text-xs text-gray-500 font-mono">{snippet.file_name}</span>
                  )}
                  {snippet.comment_count > 0 && (
                    <span className="ml-auto text-[10px] text-gray-500">
                      {t('post.snippetComments', { count: snippet.comment_count })}
                    </span>
                  )}
                </div>
                <SyntaxHighlighter
                  language={snippet.language}
                  style={vscDarkPlus}
                  showLineNumbers
                  customStyle={{
                    margin: 0,
                    borderRadius: 0,
                    fontSize: '0.75rem',
                    maxHeight: '120px',
                    overflow: 'hidden',
                  }}
                >
                  {snippet.code.split('\n').slice(0, 6).join('\n')}
                </SyntaxHighlighter>
              </div>
            ))}
            {post.code_snippets.length > 2 && (
              <div className="px-3 py-1.5 text-xs text-gray-500 text-center bg-gray-800/80">
                +{post.code_snippets.length - 2} {t('post.moreSnippets')}
              </div>
            )}
          </div>
        </Link>
      )}

      <div className="flex items-center gap-4 mt-4 pt-3 border-t border-gray-800">
        <button
          onClick={handleLike}
          className={`flex items-center gap-1.5 text-sm transition-colors ${
            liked ? 'text-red-400' : 'text-gray-500 hover:text-red-400'
          }`}
        >
          <svg className="w-4 h-4" fill={liked ? 'currentColor' : 'none'} stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
          </svg>
          {likeCount}
        </button>
        <Link
          to={`/posts/${post.id}`}
          className="flex items-center gap-1.5 text-sm text-gray-500 hover:text-blue-400 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
          </svg>
          {post.comment_count}
        </Link>
        <button
          onClick={handleBookmark}
          className={`flex items-center gap-1.5 text-sm transition-colors ml-auto ${
            bookmarked ? 'text-yellow-400' : 'text-gray-500 hover:text-yellow-400'
          }`}
          title={bookmarked ? t('post.unbookmark') : t('post.bookmark')}
        >
          <svg className="w-4 h-4" fill={bookmarked ? 'currentColor' : 'none'} stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0 1 11.186 0Z" />
          </svg>
        </button>
      </div>
    </div>
  );
}
