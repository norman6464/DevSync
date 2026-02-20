import { useState, useCallback, useMemo } from 'react';
import { useLearningLogs, useLearningLogCalendar } from './useLearningLogs';
import { useAuthStore } from '../store/authStore';
import type { LearningLog, LogCategory } from '../types/learningLog';

export function useLearningLogForm() {
  const user = useAuthStore((s) => s.user);
  const {
    logs, loading, saving,
    createLog, updateLog, deleteLog, toggleFavorite,
  } = useLearningLogs();
  const { calendarData, refetchCalendar } = useLearningLogCalendar(user?.id);

  // UI状態
  const [view, setView] = useState<'list' | 'calendar'>('list');
  const [showForm, setShowForm] = useState(false);
  const [editingLog, setEditingLog] = useState<LearningLog | null>(null);
  const [filterDate, setFilterDate] = useState<string | null>(null);
  const [filterCategory, setFilterCategory] = useState<'all' | LogCategory>('all');
  const [showFavoritesOnly, setShowFavoritesOnly] = useState(false);
  const [sortBy, setSortBy] = useState<'latest' | 'oldest' | 'duration_desc' | 'duration_asc'>('latest');

  // フォーム状態
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [category, setCategory] = useState<LogCategory>('coding');
  const [duration, setDuration] = useState('');

  const resetForm = useCallback(() => {
    setTitle('');
    setContent('');
    setCategory('coding');
    setDuration('');
    setEditingLog(null);
    setShowForm(false);
  }, []);

  const handleSubmit = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;

    const data = {
      title,
      content,
      category,
      duration: duration ? parseInt(duration) : 0,
    };

    if (editingLog) {
      const result = await updateLog(editingLog.id, data);
      if (result) {
        resetForm();
        refetchCalendar();
      }
    } else {
      const result = await createLog(data);
      if (result) {
        resetForm();
        refetchCalendar();
      }
    }
  }, [title, content, category, duration, editingLog, updateLog, createLog, resetForm, refetchCalendar]);

  const handleEdit = useCallback((log: LearningLog) => {
    setEditingLog(log);
    setTitle(log.title);
    setContent(log.content);
    setCategory(log.category);
    setDuration(log.duration ? String(log.duration) : '');
    setShowForm(true);
  }, []);

  const handleDelete = useCallback(async (logId: number) => {
    const ok = await deleteLog(logId);
    if (ok) refetchCalendar();
  }, [deleteLog, refetchCalendar]);

  const handleDateClick = useCallback((date: string) => {
    setFilterDate(date);
    setView('list');
  }, []);

  const clearFilterDate = useCallback(() => {
    setFilterDate(null);
  }, []);

  const filteredLogs = useMemo(() => {
    const filtered = logs.filter((log) => {
      if (filterDate && log.created_at.split('T')[0] !== filterDate) return false;
      if (filterCategory !== 'all' && log.category !== filterCategory) return false;
      if (showFavoritesOnly && !log.is_favorite) return false;
      return true;
    });
    return [...filtered].sort((a, b) => {
      switch (sortBy) {
        case 'oldest':
          return new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
        case 'duration_desc':
          return (b.duration || 0) - (a.duration || 0);
        case 'duration_asc':
          return (a.duration || 0) - (b.duration || 0);
        default:
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
      }
    });
  }, [logs, filterDate, filterCategory, showFavoritesOnly, sortBy]);

  return {
    // データ
    logs, filteredLogs, calendarData, loading, saving,
    // UI状態
    view, setView,
    showForm, setShowForm,
    editingLog,
    filterDate, clearFilterDate,
    filterCategory, setFilterCategory,
    showFavoritesOnly, setShowFavoritesOnly,
    sortBy, setSortBy,
    // フォーム状態
    title, setTitle,
    content, setContent,
    category, setCategory,
    duration, setDuration,
    // アクション
    resetForm,
    handleSubmit,
    handleEdit,
    handleDelete,
    handleDateClick,
    toggleFavorite,
  };
}
