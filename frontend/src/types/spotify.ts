export interface SpotifyCurrentlyPlaying {
  is_playing: boolean;
  track_name: string;
  artist_name: string;
  album_name: string;
  album_image: string;
  track_url: string;
  progress_ms: number;
  duration_ms: number;
}

export interface SpotifyRecentTrack {
  track_name: string;
  artist_name: string;
  album_name: string;
  album_image: string;
  track_url: string;
  played_at: string;
}
