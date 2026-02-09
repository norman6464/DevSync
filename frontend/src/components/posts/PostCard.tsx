import { Link } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { useTranslation } from 'react-i18next';
import { likePost, unlikePost } from '../../api/posts';
import type { Post } from '../../types/post';
import Avatar from '../common/Avatar';
import { format } from 'date-fns';
import { useState } from 'react';
import { Code2 } from 'lucide-react';

interface PostCardProps {
  post: Post;
  onUpdate?: () => void;
}

export default function PostCard({ post, onUpdate }: PostCardProps) {
  const { t } = useTranslation();
  const [liked, setLiked] = useState(post.liked || false);
  const [likeCount, setLikeCount] = useState(post.like_count);

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
    } catch {
      // handle error
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
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-5 hover:border-gray-700 transition-colors">
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

      <Link to={`/posts/${post.id}`} className="block group">
        <h3 className="text-base font-semibold mb-1.5 group-hover:text-blue-400 transition-colors">{post.title}</h3>
        <div className="text-gray-400 text-sm leading-relaxed prose prose-sm prose-invert max-w-none line-clamp-4">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
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
      </div>
    </div>
  );
}
