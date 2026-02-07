import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { BookOpen, Video, FileText, GraduationCap, BookMarked, Mic, Wrench, Pin, type LucideIcon } from 'lucide-react';
import type { LearningResource, ResourceCategory, ResourceDifficulty } from '../../types/resource';
import Avatar from '../common/Avatar';

interface ResourceCardProps {
  resource: LearningResource;
  hasLiked?: boolean;
  hasSaved?: boolean;
  onLike?: () => void;
  onUnlike?: () => void;
  onSave?: () => void;
  onUnsave?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  isOwner?: boolean;
  showUser?: boolean;
}

const categoryIcons: Record<ResourceCategory, LucideIcon> = {
  book: BookOpen,
  video: Video,
  article: FileText,
  course: GraduationCap,
  tutorial: BookMarked,
  podcast: Mic,
  tool: Wrench,
  other: Pin,
};

const difficultyColors: Record<ResourceDifficulty, string> = {
  beginner: 'bg-green-500/20 text-green-400',
  intermediate: 'bg-yellow-500/20 text-yellow-400',
  advanced: 'bg-red-500/20 text-red-400',
};

export default function ResourceCard({
  resource,
  hasLiked = false,
  hasSaved = false,
  onLike,
  onUnlike,
  onSave,
  onUnsave,
  onEdit,
  onDelete,
  isOwner = false,
  showUser = true,
}: ResourceCardProps) {
  const { t } = useTranslation();
  const [liked, setLiked] = useState(hasLiked);
  const [saved, setSaved] = useState(hasSaved);
  const [likeCount, setLikeCount] = useState(resource.like_count);
  const [saveCount, setSaveCount] = useState(resource.save_count);

  const tags = resource.tags ? JSON.parse(resource.tags) : [];

  const handleLike = () => {
    if (liked) {
      setLiked(false);
      setLikeCount(prev => Math.max(0, prev - 1));
      onUnlike?.();
    } else {
      setLiked(true);
      setLikeCount(prev => prev + 1);
      onLike?.();
    }
  };

  const handleSave = () => {
    if (saved) {
      setSaved(false);
      setSaveCount(prev => Math.max(0, prev - 1));
      onUnsave?.();
    } else {
      setSaved(true);
      setSaveCount(prev => prev + 1);
      onSave?.();
    }
  };

  return (
    <div className="bg-gray-800 rounded-xl overflow-hidden border border-gray-700 hover:border-gray-600 transition-colors">
      <div className="p-4">
        {/* Header */}
        <div className="flex items-start justify-between gap-2">
          <div className="flex items-center gap-2">
            {(() => { const Icon = categoryIcons[resource.category]; return <Icon className="w-6 h-6 text-gray-400" />; })()}
            <div>
              <span className="text-xs text-gray-400 uppercase">
                {t(`resources.categories.${resource.category}`)}
              </span>
              {resource.difficulty && (
                <span className={`ml-2 px-2 py-0.5 text-xs rounded-full ${difficultyColors[resource.difficulty]}`}>
                  {t(`resources.difficulty.${resource.difficulty}`)}
                </span>
              )}
            </div>
          </div>

          {isOwner && (
            <div className="flex gap-1">
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

        {/* Title & Description */}
        <h3 className="text-lg font-semibold text-white mt-3">
          {resource.url ? (
            <a
              href={resource.url}
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-green-400 transition-colors"
            >
              {resource.title}
            </a>
          ) : (
            resource.title
          )}
        </h3>

        {resource.description && (
          <p className="text-gray-300 text-sm mt-2 line-clamp-2">{resource.description}</p>
        )}

        {/* Tags */}
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-1.5 mt-3">
            {tags.slice(0, 4).map((tag: string, idx: number) => (
              <span
                key={idx}
                className="px-2 py-0.5 bg-gray-700 text-gray-300 text-xs rounded"
              >
                #{tag}
              </span>
            ))}
            {tags.length > 4 && (
              <span className="px-2 py-0.5 text-gray-400 text-xs">
                +{tags.length - 4}
              </span>
            )}
          </div>
        )}

        {/* User & Actions */}
        <div className="flex items-center justify-between mt-4 pt-3 border-t border-gray-700">
          {showUser && resource.user && (
            <Link
              to={`/profile/${resource.user_id}`}
              className="flex items-center gap-2 hover:opacity-80 transition-opacity"
            >
              <Avatar avatarUrl={resource.user.avatar_url} name={resource.user.name} size="sm" />
              <span className="text-sm text-gray-400">{resource.user.name}</span>
            </Link>
          )}

          <div className="flex items-center gap-3 ml-auto">
            {/* Like */}
            <button
              onClick={handleLike}
              className={`flex items-center gap-1 text-sm transition-colors ${
                liked ? 'text-red-400' : 'text-gray-400 hover:text-red-400'
              }`}
            >
              <svg
                className="w-5 h-5"
                fill={liked ? 'currentColor' : 'none'}
                stroke="currentColor"
                strokeWidth="2"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
              </svg>
              {likeCount > 0 && <span>{likeCount}</span>}
            </button>

            {/* Save */}
            <button
              onClick={handleSave}
              className={`flex items-center gap-1 text-sm transition-colors ${
                saved ? 'text-yellow-400' : 'text-gray-400 hover:text-yellow-400'
              }`}
            >
              <svg
                className="w-5 h-5"
                fill={saved ? 'currentColor' : 'none'}
                stroke="currentColor"
                strokeWidth="2"
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M17.593 3.322c1.1.128 1.907 1.077 1.907 2.185V21L12 17.25 4.5 21V5.507c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0 1 11.186 0Z" />
              </svg>
              {saveCount > 0 && <span>{saveCount}</span>}
            </button>

            {/* External Link */}
            {resource.url && (
              <a
                href={resource.url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-gray-400 hover:text-white transition-colors"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 6H5.25A2.25 2.25 0 0 0 3 8.25v10.5A2.25 2.25 0 0 0 5.25 21h10.5A2.25 2.25 0 0 0 18 18.75V10.5m-10.5 6L21 3m0 0h-5.25M21 3v5.25" />
                </svg>
              </a>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
