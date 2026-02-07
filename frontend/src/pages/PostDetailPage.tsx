import { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { usePostDetail } from '../hooks';
import PostCard from '../components/posts/PostCard';
import Avatar from '../components/common/Avatar';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { format } from 'date-fns';

export default function PostDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { post, comments, loading, submitting, submitComment, refetch } = usePostDetail(id);
  const [newComment, setNewComment] = useState('');

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim()) return;
    const success = await submitComment(newComment);
    if (success) setNewComment('');
  };

  if (loading) return <div className="py-12"><LoadingSpinner /></div>;
  if (!post) return <div className="text-center text-gray-400 py-12">Post not found</div>;

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <PostCard post={post} onUpdate={refetch} />

      {/* Comments Section */}
      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <div className="px-6 py-4 border-b border-gray-800 flex items-center gap-2">
          <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
          </svg>
          <h3 className="text-sm font-semibold">Comments ({comments.length})</h3>
        </div>

        {/* Comment Form */}
        <div className="px-6 py-4 border-b border-gray-800">
          <form onSubmit={handleSubmitComment} className="flex gap-3">
            <input
              type="text"
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              placeholder="Write a comment..."
              className="flex-1 px-4 py-2.5 bg-gray-800/50 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow"
            />
            <button
              type="submit"
              disabled={submitting || !newComment.trim()}
              className="px-5 py-2.5 bg-green-600 hover:bg-green-500 disabled:opacity-40 disabled:hover:bg-green-600 text-white rounded-lg font-medium text-sm transition-colors"
            >
              {submitting ? 'Posting...' : 'Comment'}
            </button>
          </form>
        </div>

        {/* Comments List */}
        {comments.length === 0 ? (
          <div className="px-6 py-8 text-center text-gray-500 text-sm">
            No comments yet. Be the first to comment.
          </div>
        ) : (
          <div className="divide-y divide-gray-800/50">
            {comments.map((comment) => (
              <div key={comment.id} className="px-6 py-4 flex gap-3">
                <Link to={`/profile/${comment.user_id}`} className="flex-shrink-0">
                  <Avatar name={comment.user?.name || 'U'} avatarUrl={comment.user?.avatar_url} size="sm" />
                </Link>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2 flex-wrap">
                    <Link
                      to={`/profile/${comment.user_id}`}
                      className="font-medium text-sm hover:text-blue-400 transition-colors"
                    >
                      {comment.user?.name}
                    </Link>
                    <span className="text-xs text-gray-600">
                      {format(new Date(comment.created_at), 'MMM d, yyyy · HH:mm')}
                    </span>
                  </div>
                  <p className="text-sm text-gray-300 mt-1 leading-relaxed">{comment.content}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
