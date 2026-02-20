import { Link } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeSanitize from 'rehype-sanitize';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { useTranslation } from 'react-i18next';
import { Code2 } from 'lucide-react';
import type { Post } from '../../types/post';
import { sanitizeUrl } from '../../utils/url';

interface PostCardContentProps {
  post: Post;
  imageUrls: string[];
}

export default function PostCardContent({ post, imageUrls }: PostCardContentProps) {
  const { t } = useTranslation();

  return (
    <>
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
                src={sanitizeUrl(url) ?? ''}
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
    </>
  );
}
