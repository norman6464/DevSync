import client from './client';
import type { UserDashboardStats } from '../types/dashboard';

export const getUserDashboardStats = (userId: number) =>
  client.get<UserDashboardStats>(`/users/${userId}/dashboard-stats`);
