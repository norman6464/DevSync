import { useAsyncData } from './useAsyncData';
import { getMyGoals, type LearningGoal } from '../api/goals';
import { getNotifications } from '../api/notifications';
import type { Notification } from '../types/notification';

export function useDashboard() {
  const { data: goals, loading: goalsLoading } = useAsyncData(
    async () => {
      const { data } = await getMyGoals();
      return data || [];
    },
    { initialData: [] as LearningGoal[] }
  );

  const { data: notificationsData, loading: notificationsLoading } = useAsyncData(
    async () => {
      const { data } = await getNotifications(1, 5);
      return data;
    },
    { initialData: { notifications: [] as Notification[], total: 0 } }
  );

  const activeGoals = goals.filter((g) => g.status === 'active');
  const completedGoals = goals.filter((g) => g.status === 'completed');
  const avgProgress =
    activeGoals.length > 0
      ? Math.round(activeGoals.reduce((sum, g) => sum + g.progress, 0) / activeGoals.length)
      : 0;

  return {
    goals,
    activeGoals,
    completedGoals,
    avgProgress,
    goalsLoading,
    recentNotifications: notificationsData.notifications,
    notificationTotal: notificationsData.total,
    notificationsLoading,
  };
}
