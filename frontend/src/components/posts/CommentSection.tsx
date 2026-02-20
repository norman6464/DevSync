import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { sectionContainerClass, messageInputClass } from '../../constants/styles';
import type { Comment } from '../../types/post';
import CommentItem from './CommentItem';
import MentionInput from '../common/MentionInput';

interface CommentSectionProps {
  comments: Comment[];
  submitting: boolean;
  onSubmitComment: (content: string, parentId?: number) => Promise<boolean>;
}

export default function CommentSection({ comments, submitting, onSubmitComment }: CommentSectionProps) {
  const { t } = useTranslation();
  const [newComment, setNewComment] = useState('');
  const [replyingTo, setReplyingTo] = useState<number | null>(null);
  const [replyContent, setReplyContent] = useState('');

  const handleSubmitComment = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim()) return;
    const success = await onSubmitComment(newComment);
    if (success) setNewComment('');
  }, [newComment, onSubmitComment]);

  const handleSubmitReply = useCallback(async (e: React.FormEvent, parentId: number) => {
    e.preventDefault();
    if (!replyContent.trim()) return;
    const success = await onSubmitComment(replyContent, parentId);
    if (success) {
      setReplyContent('');
      setReplyingTo(null);
    }
  }, [replyContent, onSubmitComment]);

  const handleReply = useCallback((commentId: number) => {
    setReplyingTo((prev) => (prev === commentId ? null : commentId));
    setReplyContent('');
  }, []);

  return (
    <div className={sectionContainerClass}>
      <div className="px-6 py-4 border-b border-gray-800 flex items-center gap-2">
        <svg aria-hidden="true" className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
        </svg>
        <h3 className="text-sm font-semibold">{t('post.comments')} ({comments.length})</h3>
      </div>

      {/* Comment Form */}
      <div className="px-6 py-4 border-b border-gray-800">
        <form onSubmit={handleSubmitComment} className="flex gap-3">
          <MentionInput
            value={newComment}
            onChange={setNewComment}
            placeholder={t('post.writeComment')}
            className={messageInputClass}
            disabled={submitting}
          />
          <button
            type="submit"
            disabled={submitting || !newComment.trim()}
            className="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 disabled:hover:bg-gray-700 text-white rounded-lg font-medium text-sm transition-colors"
          >
            {submitting ? t('post.posting') : t('post.comment')}
          </button>
        </form>
      </div>

      {/* Comments List */}
      {comments.length === 0 ? (
        <div className="px-6 py-8 text-center text-gray-500 text-sm">
          {t('post.noComments')}
        </div>
      ) : (
        <div className="divide-y divide-gray-800/50">
          {comments.map((comment) => (
            <div key={comment.id} className="px-6 py-4">
              <CommentItem
                comment={comment}
                onReply={handleReply}
                isReplying={replyingTo === comment.id}
              />

              {/* Reply Form */}
              {replyingTo === comment.id && (
                <div className="ml-11 mt-3">
                  <form onSubmit={(e) => handleSubmitReply(e, comment.id)} className="flex gap-2">
                    <MentionInput
                      value={replyContent}
                      onChange={setReplyContent}
                      placeholder={t('post.writeReply')}
                      className={messageInputClass}
                      disabled={submitting}
                    />
                    <button
                      type="submit"
                      disabled={submitting || !replyContent.trim()}
                      className="px-4 py-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 disabled:hover:bg-gray-700 text-white rounded-lg font-medium text-xs transition-colors"
                    >
                      {t('post.reply')}
                    </button>
                  </form>
                </div>
              )}

              {/* Replies */}
              {comment.replies && comment.replies.length > 0 && (
                <div className="ml-11 mt-3 space-y-3 border-l-2 border-gray-800 pl-4">
                  {comment.replies.map((reply) => (
                    <CommentItem key={reply.id} comment={reply} size="xs" />
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
