import client from './client';
import type { SpotifyCurrentlyPlaying, SpotifyRecentTrack } from '../types/spotify';

export const getSpotifyConnectURL = () =>
  client.get<{ url: string }>('/spotify/connect');

export const spotifyCallback = (code: string, state: string) =>
  client.get('/spotify/callback', { params: { code, state } });

export const getCurrentlyPlaying = (userId: number) =>
  client.get<SpotifyCurrentlyPlaying | null>(`/spotify/currently-playing/${userId}`);

export const getRecentlyPlayed = (userId: number) =>
  client.get<SpotifyRecentTrack[]>(`/spotify/recently-played/${userId}`);

export const disconnectSpotify = () =>
  client.delete('/spotify/disconnect');
