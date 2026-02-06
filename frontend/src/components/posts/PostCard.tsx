import { Link } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { likePost, unlikePost } from '../../api/posts';
import type { Post } from '../../types/post';
import Avatar from '../common/Avatar';
import { format } from 'date-fns';
import { useState } from 'react';

interface PostCardProps {
  post: Post;
  onUpdate?: () => void;
}

export default function PostCard({ post, onUpdate }: PostCardProps) {
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
