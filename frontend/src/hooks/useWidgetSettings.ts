import { useState, useCallback } from 'react';
import { useAsyncData } from './useAsyncData';
import { getWidgetSettings, updateWidgetSettings } from '../api/widgetSettings';
import type { WidgetConfig } from '../types/widgetSettings';
import { useToast } from './useToast';
import { useTranslation } from 'react-i18next';

const DEFAULT_WIDGETS: WidgetConfig[] = [
  { key: 'userProfile', visible: true, order: 0 },
  { key: 'level', visible: true, order: 1 },
  { key: 'streak', visible: true, order: 2 },
  { key: 'dailyChallenge', visible: true, order: 3 },
  { key: 'weeklyChallenge', visible: true, order: 4 },
  { key: 'studyCircle', visible: true, order: 5 },
  { key: 'quickEntry', visible: true, order: 6 },
  { key: 'quickActions', visible: true, order: 7 },
  { key: 'recommendedUsers', visible: true, order: 8 },
  { key: 'trending', visible: true, order: 9 },
  { key: 'aiAdvice', visible: true, order: 10 },
  { key: 'goalsProgress', visible: true, order: 11 },
  { key: 'recentNotifications', visible: true, order: 12 },
  { key: 'quickStats', visible: true, order: 13 },
];

export function useWidgetSettings() {
  const { t } = useTranslation();
  const toast = useToast();
  const [editing, setEditing] = useState(false);
  const [localWidgets, setLocalWidgets] = useState<WidgetConfig[]>([]);
  const [saving, setSaving] = useState(false);

  const { data: widgets, loading, refetch } = useAsyncData(
    async () => {
      const res = await getWidgetSettings();
      if (res.data?.settings) {
        try {
          return JSON.parse(res.data.settings) as WidgetConfig[];
        } catch {
          return DEFAULT_WIDGETS;
        }
      }
      return DEFAULT_WIDGETS;
    },
    { initialData: DEFAULT_WIDGETS }
  );

  const startEditing = useCallback(() => {
    setLocalWidgets([...widgets].sort((a, b) => a.order - b.order));
    setEditing(true);
  }, [widgets]);

  const cancelEditing = useCallback(() => {
    setEditing(false);
    setLocalWidgets([]);
  }, []);

  const toggleVisibility = useCallback((key: string) => {
    setLocalWidgets((prev) =>
      prev.map((w) => (w.key === key ? { ...w, visible: !w.visible } : w))
    );
  }, []);

  const moveUp = useCallback((key: string) => {
    setLocalWidgets((prev) => {
      const sorted = [...prev].sort((a, b) => a.order - b.order);
      const idx = sorted.findIndex((w) => w.key === key);
      if (idx <= 0) return prev;
      const newOrder = sorted.map((w, i) => ({ ...w, order: i }));
      const temp = newOrder[idx].order;
      newOrder[idx].order = newOrder[idx - 1].order;
      newOrder[idx - 1].order = temp;
      return newOrder.sort((a, b) => a.order - b.order);
    });
  }, []);

  const moveDown = useCallback((key: string) => {
    setLocalWidgets((prev) => {
      const sorted = [...prev].sort((a, b) => a.order - b.order);
      const idx = sorted.findIndex((w) => w.key === key);
      if (idx < 0 || idx >= sorted.length - 1) return prev;
      const newOrder = sorted.map((w, i) => ({ ...w, order: i }));
      const temp = newOrder[idx].order;
      newOrder[idx].order = newOrder[idx + 1].order;
      newOrder[idx + 1].order = temp;
      return newOrder.sort((a, b) => a.order - b.order);
    });
  }, []);

  const saveSettings = useCallback(async () => {
    setSaving(true);
    try {
      await updateWidgetSettings(localWidgets);
      await refetch();
      setEditing(false);
      toast.success(t('widgetSettings.saveSuccess'));
    } catch {
      toast.error(t('widgetSettings.saveFailed'));
    } finally {
      setSaving(false);
    }
  }, [localWidgets, refetch, toast, t]);

  const sortedWidgets = [...widgets].sort((a, b) => a.order - b.order);
  const sortedLocalWidgets = [...localWidgets].sort((a, b) => a.order - b.order);

  return {
    widgets: sortedWidgets,
    loading,
    editing,
    localWidgets: sortedLocalWidgets,
    saving,
    startEditing,
    cancelEditing,
    toggleVisibility,
    moveUp,
    moveDown,
    saveSettings,
  };
}
