import { useTranslation } from 'react-i18next';
import { ExternalLink, Play, Sparkles } from 'lucide-react';
import type { YouTubeVideo } from '../../types/youtube';
import { cardClass } from '../../constants/styles';
import { sanitizeUrl } from '../../utils/url';
import { formatDate } from '../../utils/timeFormat';

interface YouTubeVideoCardProps {
  video: YouTubeVideo;
}

export default function YouTubeVideoCard({ video }: YouTubeVideoCardProps) {
  const { t } = useTranslation();
  const videoURL = `https://www.youtube.com/watch?v=${encodeURIComponent(video.video_id)}`;
  const channelURL = `https://www.youtube.com/channel/${encodeURIComponent(video.channel_id)}`;
  const isNew = Date.now() - new Date(video.published_at).getTime() < 7 * 24 * 60 * 60 * 1000;

  return (
    <div className={cardClass}>
      <a href={videoURL} target="_blank" rel="noopener noreferrer" className="block relative group">
        <img
          src={sanitizeUrl(video.thumbnail_url) ?? ''}
          alt={video.title}
          referrerPolicy="no-referrer"
          className="w-full aspect-video object-cover"
          loading="lazy"
        />
        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
          <Play className="w-12 h-12 text-white" fill="white" />
        </div>
        {isNew && (
          <span className="absolute top-2 left-2 inline-flex items-center gap-0.5 px-1.5 py-0.5 bg-red-500/90 text-white text-xs rounded font-medium">
            <Sparkles className="w-3 h-3" />
            {t('youtube.newVideo')}
          </span>
        )}
      </a>

      <div className="p-4">
        <h3 className="text-sm font-semibold text-white line-clamp-2">
          <a href={videoURL} target="_blank" rel="noopener noreferrer"
             className="hover:text-red-400 transition-colors">
            {video.title}
          </a>
        </h3>

        <a href={channelURL} target="_blank" rel="noopener noreferrer"
           className="text-xs text-gray-400 hover:text-gray-300 mt-1 block">
          {video.channel_title}
        </a>

        {video.description && (
          <p className="text-xs text-gray-500 mt-2 line-clamp-2">{video.description}</p>
        )}

        <div className="flex items-center justify-between mt-3 pt-2 border-t border-gray-700">
          <span className="text-xs text-gray-500">
            {formatDate(video.published_at)}
          </span>
          <a href={videoURL} target="_blank" rel="noopener noreferrer"
             className="text-gray-400 hover:text-red-400 transition-colors"
             aria-label={t('youtube.watchOnYoutube')}>
            <ExternalLink className="w-4 h-4" />
          </a>
        </div>
      </div>
    </div>
  );
}
