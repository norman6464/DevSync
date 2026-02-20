import { useState, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { getCurrentlyPlaying, getRecentlyPlayed } from '../../api/spotify';
import type { SpotifyCurrentlyPlaying, SpotifyRecentTrack } from '../../types/spotify';
import { sanitizeUrl } from '../../utils/url';

interface Props {
  userId: number;
}

const SpotifyIcon = ({ className }: { className?: string }) => (
  <svg className={className} viewBox="0 0 24 24" fill="currentColor">
    <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.419 1.56-.299.421-1.02.599-1.559.3z" />
  </svg>
);

const formatTime = (ms: number) => {
  const min = Math.floor(ms / 60000);
  const sec = Math.floor((ms % 60000) / 1000);
  return `${min}:${sec.toString().padStart(2, '0')}`;
};

export default function SpotifyNowPlaying({ userId }: Props) {
  const { t } = useTranslation();
  const [currentTrack, setCurrentTrack] = useState<SpotifyCurrentlyPlaying | null>(null);
  const [recentTracks, setRecentTracks] = useState<SpotifyRecentTrack[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchData = useCallback(async () => {
    try {
      const { data: current } = await getCurrentlyPlaying(userId);
      setCurrentTrack(current);

      if (!current) {
        const { data: recent } = await getRecentlyPlayed(userId);
        setRecentTracks(recent || []);
      }
    } catch {
      // Spotify未連携やエラーの場合は何も表示しない
    } finally {
      setLoading(false);
    }
  }, [userId]);

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 30000);
    return () => clearInterval(interval);
  }, [fetchData]);

  if (loading) return null;
  if (!currentTrack && recentTracks.length === 0) return null;

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
      <h3 className="text-sm font-semibold text-gray-300 mb-4 flex items-center gap-2">
        <SpotifyIcon className="w-5 h-5 text-green-500" />
        {currentTrack ? t('spotify.nowPlaying') : t('spotify.recentlyPlayed')}
      </h3>

      {currentTrack ? (
        <a
          href={sanitizeUrl(currentTrack.track_url) || '#'}
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center gap-4 p-3 bg-gray-800/50 rounded-lg hover:bg-gray-800 transition-colors"
        >
          {currentTrack.album_image && (
            <img
              src={sanitizeUrl(currentTrack.album_image) ?? ''}
              alt={currentTrack.album_name}
              referrerPolicy="no-referrer"
              className="w-16 h-16 rounded-lg object-cover"
            />
          )}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              {currentTrack.is_playing && (
                <span className="flex gap-0.5" aria-label={t('spotify.nowPlaying')}>
                  <span className="w-1 h-3 bg-green-500 rounded-full animate-pulse" />
                  <span className="w-1 h-4 bg-green-500 rounded-full animate-pulse" style={{ animationDelay: '0.2s' }} />
                  <span className="w-1 h-2 bg-green-500 rounded-full animate-pulse" style={{ animationDelay: '0.4s' }} />
                </span>
              )}
              <span className="font-medium text-white truncate">{currentTrack.track_name}</span>
            </div>
            <p className="text-sm text-gray-400 truncate">{currentTrack.artist_name}</p>
            <p className="text-xs text-gray-500 truncate">{currentTrack.album_name}</p>
            <div className="mt-1.5">
              <div className="h-1 bg-gray-700 rounded-full overflow-hidden">
                <div
                  className="h-full bg-green-500 rounded-full transition-all"
                  style={{ width: `${(currentTrack.progress_ms / currentTrack.duration_ms) * 100}%` }}
                />
              </div>
              <div className="flex justify-between text-xs text-gray-500 mt-0.5">
                <span>{formatTime(currentTrack.progress_ms)}</span>
                <span>{formatTime(currentTrack.duration_ms)}</span>
              </div>
            </div>
          </div>
        </a>
      ) : (
        <div className="space-y-2">
          {recentTracks.slice(0, 5).map((track, index) => (
            <a
              key={`${track.track_url}-${index}`}
              href={sanitizeUrl(track.track_url) || '#'}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-800/50 transition-colors"
            >
              {track.album_image && (
                <img
                  src={sanitizeUrl(track.album_image) ?? ''}
                  alt={track.album_name}
                  referrerPolicy="no-referrer"
                  className="w-10 h-10 rounded object-cover"
                />
              )}
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-white truncate">{track.track_name}</p>
                <p className="text-xs text-gray-400 truncate">{track.artist_name}</p>
              </div>
            </a>
          ))}
        </div>
      )}
    </div>
  );
}
