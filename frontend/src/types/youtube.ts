export interface YouTubeVideo {
  id: number;
  video_id: string;
  title: string;
  description: string;
  channel_id: string;
  channel_title: string;
  thumbnail_url: string;
  published_at: string;
}

export interface YouTubeSearchResponse {
  videos: YouTubeVideo[];
  query: string;
  cached: boolean;
  total: number;
}

export interface YouTubeRecommendResponse {
  videos: YouTubeVideo[];
  skills: string[];
  available: boolean;
}
