import { Link } from 'react-router-dom';
import { likePost, unlikePost } from '../../api/posts';
import { useAuthStore } from '../../store/authStore';
import type { Post } from '../../types/post';
import Avatar from '../common/Avatar';
import { format } from 'date-fns';
import { useState } from 'react';

interface PostCardProps {
  post: Post;
  onUpdate?: () => void;
}

export default function PostCard({ post, onUpdate }: PostCardProps) {
  const currentUser = useAuthStore((s) => s.user);
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

  return (
    <div className="bg-gray-900 rounded-lg p-6">
      <div className="flex items-center gap-3 mb-3">
        <Link to={`/profile/${post.user_id}`}>
          <Avatar name={post.user?.name || 'U'} avatarUrl={post.user?.avatar_url} size="sm" />
        </Link>
        <div>
          <Link to={`/profile/${post.user_id}`} className="font-medium hover:text-blue-400 text-sm">
            {post.user?.name}
          </Link>
          <p className="text-xs text-gray-500">
            {format(new Date(post.created_at), 'MMM d, yyyy HH:mm')}
          </p>
        </div>
      </div>

      <Link to={`/posts/${post.id}`}>
        <h3 className="text-lg font-semibold mb-2 hover:text-blue-400">{post.title}</h3>
      </Link>
      <p className="text-gray-300 text-sm whitespace-pre-wrap">{post.content}</p>

      <div className="flex items-center gap-4 mt-4 pt-3 border-t border-gray-800">
        <button
          onClick={handleLike}
          className={`flex items-center gap-1 text-sm transition-colors ${
            liked ? 'text-red-400' : 'text-gray-400 hover:text-red-400'
          }`}
        >
          {liked ? '\u2665' : '\u2661'} {likeCount}
        </button>
        <Link
          to={`/posts/${post.id}`}
          className="text-sm text-gray-400 hover:text-blue-400 transition-colors"
        >
          Comments ({post.comment_count})
        </Link>
      </div>
    </div>
  );
}
