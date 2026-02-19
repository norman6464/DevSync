import client from './client';
import type { YouTubeSearchResponse, YouTubeRecommendResponse } from '../types/youtube';

export const searchYouTubeVideos = async (
  query: string,
  language = 'ja'
): Promise<YouTubeSearchResponse> => {
  const res = await client.get('/youtube/search', { params: { q: query, lang: language } });
  return res.data;
};

export const getYouTubeRecommendations = async (): Promise<YouTubeRecommendResponse> => {
  const res = await client.get('/youtube/recommend');
  return res.data;
};

export const getYouTubeStatus = async (): Promise<{ available: boolean }> => {
  const res = await client.get('/youtube/status');
  return res.data;
};
