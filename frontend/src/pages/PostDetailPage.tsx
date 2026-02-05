import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getPost, getComments, createComment } from '../api/posts';
import type { Post, Comment } from '../types/post';
import PostCard from '../components/posts/PostCard';
import Avatar from '../components/common/Avatar';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { format } from 'date-fns';

export default function PostDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [post, setPost] = useState<Post | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [newComment, setNewComment] = useState('');
  const [loading, setLoading] = useState(true);

  const fetchData = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const [postRes, commentsRes] = await Promise.all([
        getPost(parseInt(id)),
        getComments(parseInt(id)),
      ]);
      setPost(postRes.data);
      setComments(commentsRes.data || []);
    } catch {
      // handle error
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [id]);

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim() || !id) return;
    try {
      await createComment(parseInt(id), newComment);
      setNewComment('');
      fetchData();
    } catch {
      // handle error
    }
  };

  if (loading) return <LoadingSpinner />;
  if (!post) return <div className="text-center text-gray-400 py-12">Post not found</div>;

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <PostCard post={post} onUpdate={fetchData} />

      <div className="bg-gray-900 rounded-lg p-6">
        <h3 className="font-semibold mb-4">Comments ({comments.length})</h3>

        <form onSubmit={handleSubmitComment} className="flex gap-2 mb-6">
          <input
            type="text"
            value={newComment}
            onChange={(e) => setNewComment(e.target.value)}
            placeholder="Write a comment..."
            className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:border-blue-500"
          />
          <button
            type="submit"
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm transition-colors"
          >
            Comment
          </button>
        </form>

        <div className="space-y-4">
          {comments.map((comment) => (
            <div key={comment.id} className="flex gap-3">
              <Link to={`/profile/${comment.user_id}`}>
                <Avatar name={comment.user?.name || 'U'} avatarUrl={comment.user?.avatar_url} size="sm" />
              </Link>
              <div>
                <div className="flex items-center gap-2">
                  <Link to={`/profile/${comment.user_id}`} className="font-medium text-sm hover:text-blue-400">
                    {comment.user?.name}
                  </Link>
                  <span className="text-xs text-gray-500">
                    {format(new Date(comment.created_at), 'MMM d, yyyy HH:mm')}
                  </span>
                </div>
                <p className="text-sm text-gray-300 mt-1">{comment.content}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
