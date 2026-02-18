import { useState, useCallback } from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { useTranslation } from 'react-i18next';
import { Copy, Check, MessageSquarePlus } from 'lucide-react';
import { useSnippetComments } from '../../hooks';
import type { CodeSnippet } from '../../types/post';
import Avatar from '../common/Avatar';
import { format } from 'date-fns';

interface CodeSnippetViewerProps {
  snippet: CodeSnippet;
  showComments?: boolean;
}

export default function CodeSnippetViewer({ snippet, showComments = true }: CodeSnippetViewerProps) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);
  const [commentLine, setCommentLine] = useState<number | null>(null);
  const [commentText, setCommentText] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const { comments, addComment, removeComment } = useSnippetComments(
    showComments ? snippet.id : 0
  );

  const handleCopy = useCallback(() => {
    navigator.clipboard.writeText(snippet.code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }, [snippet.code]);

  const handleAddComment = async () => {
    if (!commentText.trim() || commentLine === null) return;
    setSubmitting(true);
    const ok = await addComment(commentLine, commentText);
    if (ok) {
      setCommentText('');
      setCommentLine(null);
    }
    setSubmitting(false);
  };

  // Group comments by line number
  const commentsByLine = comments.reduce<Record<number, typeof comments>>((acc, c) => {
    if (!acc[c.line_number]) acc[c.line_number] = [];
    acc[c.line_number].push(c);
    return acc;
  }, {});

  const codeLines = snippet.code.split('\n');

  return (
    <div className="border border-gray-700 rounded-lg overflow-hidden bg-gray-900">
      {/* Header */}
      <div className="flex items-center justify-between px-4 py-2 border-b border-gray-700 bg-gray-800">
        <div className="flex items-center gap-2">
          <span className="px-2 py-0.5 bg-blue-500/20 text-blue-400 text-xs font-mono rounded">
            {snippet.language}
          </span>
          {snippet.file_name && (
            <span className="text-sm text-gray-300 font-mono">{snippet.file_name}</span>
          )}
        </div>
        <button
          onClick={handleCopy}
          className="flex items-center gap-1.5 px-2.5 py-1 text-xs text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-colors"
        >
          {copied ? (
            <>
              <Check className="w-3.5 h-3.5 text-green-400" aria-hidden="true" />
              <span className="text-green-400">{t('post.codeCopied')}</span>
            </>
          ) : (
            <>
              <Copy className="w-3.5 h-3.5" aria-hidden="true" />
              {t('post.copyCode')}
            </>
          )}
        </button>
      </div>

      {/* Code with inline comments */}
      {showComments ? (
        <div className="overflow-x-auto">
          <table className="w-full">
            <tbody>
              {codeLines.map((line, i) => {
                const lineNum = i + 1;
                const lineComments = commentsByLine[lineNum] || [];
                return (
                  <tr key={i} className="group">
                    <td className="w-0">
                      <div className="flex items-center">
                        {/* Add comment button */}
                        <button
                          onClick={() => setCommentLine(commentLine === lineNum ? null : lineNum)}
                          className="w-5 h-full flex items-center justify-center text-transparent group-hover:text-blue-400 hover:bg-gray-800 transition-colors"
                          aria-label={t('post.addLineComment')}
                        >
                          <MessageSquarePlus className="w-3 h-3" aria-hidden="true" />
                        </button>
                        {/* Line number */}
                        <span className="px-2 text-xs text-gray-600 select-none font-mono text-right min-w-[2.5rem] inline-block">
                          {lineNum}
                        </span>
                      </div>
                    </td>
                    <td className="px-4">
                      <pre className="text-sm font-mono text-gray-200 whitespace-pre">
                        <SyntaxHighlighter
                          language={snippet.language}
                          style={vscDarkPlus}
                          customStyle={{
                            background: 'transparent',
                            padding: 0,
                            margin: 0,
                            fontSize: '0.875rem',
                          }}
                          PreTag="span"
                          wrapLines={false}
                        >
                          {line || ' '}
                        </SyntaxHighlighter>
                      </pre>

                      {/* Inline comment form */}
                      {commentLine === lineNum && (
                        <div className="my-2 p-3 bg-gray-800 border border-gray-700 rounded-lg">
                          <textarea
                            value={commentText}
                            onChange={(e) => setCommentText(e.target.value)}
                            placeholder={t('post.lineCommentPlaceholder')}
                            className="w-full px-3 py-2 bg-gray-900 border border-gray-600 rounded text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500 resize-none"
                            rows={2}
                            autoFocus
                          />
                          <div className="flex justify-end gap-2 mt-2">
                            <button
                              type="button"
                              onClick={() => { setCommentLine(null); setCommentText(''); }}
                              className="px-3 py-1.5 text-xs text-gray-400 hover:text-white transition-colors"
                            >
                              {t('common.cancel')}
                            </button>
                            <button
                              type="button"
                              onClick={handleAddComment}
                              disabled={submitting || !commentText.trim()}
                              className="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 disabled:opacity-40 text-white text-xs rounded transition-colors"
                            >
                              {t('post.addLineComment')}
                            </button>
                          </div>
                        </div>
                      )}

                      {/* Existing comments for this line */}
                      {lineComments.length > 0 && (
                        <div className="my-1 space-y-1">
                          {lineComments.map((comment) => (
                            <div
                              key={comment.id}
                              className="flex items-start gap-2 p-2 bg-gray-800/60 border-l-2 border-blue-500/40 rounded-r"
                            >
                              <Avatar
                                name={comment.user?.name || 'U'}
                                avatarUrl={comment.user?.avatar_url}
                                size="xs"
                              />
                              <div className="flex-1 min-w-0">
                                <div className="flex items-center gap-2">
                                  <span className="text-xs font-medium text-gray-300">
                                    {comment.user?.name}
                                  </span>
                                  <span className="text-[10px] text-gray-600">
                                    {format(new Date(comment.created_at), 'MMM d, HH:mm')}
                                  </span>
                                </div>
                                <p className="text-xs text-gray-400 mt-0.5">{comment.content}</p>
                              </div>
                              <button
                                onClick={() => removeComment(comment.id)}
                                className="text-xs text-gray-600 hover:text-red-400 transition-colors"
                                aria-label={t('post.deleteComment')}
                              >
                                ×
                              </button>
                            </div>
                          ))}
                        </div>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      ) : (
        /* Simple highlight without comments */
        <SyntaxHighlighter
          language={snippet.language}
          style={vscDarkPlus}
          showLineNumbers
          customStyle={{
            margin: 0,
            borderRadius: 0,
            fontSize: '0.875rem',
          }}
        >
          {snippet.code}
        </SyntaxHighlighter>
      )}
    </div>
  );
}
