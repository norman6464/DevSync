import { useTranslation } from 'react-i18next';
import { ExternalLink, Play } from 'lucide-react';
import type { YouTubeVideo } from '../../types/youtube';
import { cardClass } from '../../constants/styles';

interface YouTubeVideoCardProps {
  video: YouTubeVideo;
}

export default function YouTubeVideoCard({ video }: YouTubeVideoCardProps) {
  const { t } = useTranslation();
  const videoURL = `https://www.youtube.com/watch?v=${video.video_id}`;
  const channelURL = `https://www.youtube.com/channel/${video.channel_id}`;

  return (
    <div className={cardClass}>
      <a href={videoURL} target="_blank" rel="noopener noreferrer" className="block relative group">
        <img
          src={video.thumbnail_url}
          alt={video.title}
          referrerPolicy="no-referrer"
          className="w-full aspect-video object-cover"
          loading="lazy"
        />
        <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
          <Play className="w-12 h-12 text-white" fill="white" />
        </div>
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
            {new Date(video.published_at).toLocaleDateString()}
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
