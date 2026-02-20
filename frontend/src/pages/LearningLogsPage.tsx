import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Code, BookOpen, GraduationCap, Users, FileText, Calendar, List, Download, Clock, Star, type LucideIcon } from 'lucide-react';
import { useLearningLogForm } from '../hooks/useLearningLogForm';
import { useWeeklyDuration, useStreak } from '../hooks';
import { useAuthStore } from '../store/authStore';
import { exportLogsCSV, type ExportPeriod } from '../api/learningLogs';
import type { LogCategory } from '../types/learningLog';
import LogCalendar from '../components/learning-logs/LogCalendar';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { Modal } from '../components/common';
import { inputClass, buttonSecondaryClass, labelClass } from '../constants/styles';

const CATEGORIES: { value: LogCategory; label: string; Icon: LucideIcon }[] = [
  { value: 'coding', label: 'learningLogs.categoryCoding', Icon: Code },
  { value: 'reading', label: 'learningLogs.categoryReading', Icon: BookOpen },
  { value: 'course', label: 'learningLogs.categoryCourse', Icon: GraduationCap },
  { value: 'meetup', label: 'learningLogs.categoryMeetup', Icon: Users },
  { value: 'other', label: 'learningLogs.categoryOther', Icon: FileText },
];

const getCategoryInfo = (cat: LogCategory) =>
  CATEGORIES.find((c) => c.value === cat) || CATEGORIES[4];

const getCategoryColor = (cat: LogCategory) => {
  switch (cat) {
    case 'coding': return 'text-blue-400 bg-blue-400/10';
    case 'reading': return 'text-green-400 bg-green-400/10';
    case 'course': return 'text-purple-400 bg-purple-400/10';
    case 'meetup': return 'text-orange-400 bg-orange-400/10';
    default: return 'text-gray-400 bg-gray-400/10';
  }
};

export default function LearningLogsPage() {
  const { t } = useTranslation();
  const user = useAuthStore((s) => s.user);
  const { weeklyDuration } = useWeeklyDuration(user?.id);
  const { streakInfo } = useStreak(user?.id);
  const [exporting, setExporting] = useState(false);
  const {
    filteredLogs, calendarData, loading, saving,
    view, setView,
    showForm, setShowForm,
    editingLog,
    filterDate, clearFilterDate,
    filterCategory, setFilterCategory,
    showFavoritesOnly, setShowFavoritesOnly,
    title, setTitle,
    content, setContent,
    category, setCategory,
    duration, setDuration,
    resetForm, handleSubmit, handleEdit, handleDelete, handleDateClick, toggleFavorite,
  } = useLearningLogForm();

  const handleExport = useCallback(async (period: ExportPeriod) => {
    setExporting(true);
    try {
      const res = await exportLogsCSV(period);
      const url = URL.createObjectURL(res.data);
      const a = document.createElement('a');
      a.href = url;
      a.download = `learning-logs-${period}-${new Date().toISOString().slice(0, 10)}.csv`;
      a.click();
      URL.revokeObjectURL(url);
    } finally {
      setExporting(false);
    }
  }, []);

  const handleShowForm = useCallback(() => setShowForm(true), [setShowForm]);
  const handleViewList = useCallback(() => { setView('list'); clearFilterDate(); }, [setView, clearFilterDate]);
  const handleViewCalendar = useCallback(() => { setView('calendar'); clearFilterDate(); }, [setView, clearFilterDate]);
  const handleFilterAll = useCallback(() => setFilterCategory('all'), [setFilterCategory]);
  const handleToggleFavoritesFilter = useCallback(() => setShowFavoritesOnly(prev => !prev), [setShowFavoritesOnly]);
  const handleTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setTitle(e.target.value), [setTitle]);
  const handleContentChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => setContent(e.target.value), [setContent]);
  const handleCategoryChange = useCallback((e: React.ChangeEvent<HTMLSelectElement>) => setCategory(e.target.value as LogCategory), [setCategory]);
  const handleDurationChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setDuration(e.target.value), [setDuration]);

  if (loading) return <div className="py-12"><LoadingSpinner /></div>;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">{t('learningLogs.title')}</h1>
          <p className="text-sm text-gray-400 mt-1">{t('learningLogs.subtitle')}</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="relative group">
            <button
              disabled={exporting}
              className={`${buttonSecondaryClass} font-medium text-sm flex items-center gap-1.5`}
            >
              <Download className="w-4 h-4" />
              {exporting ? t('common.loading') : t('learningLogs.exportCSV')}
            </button>
            <div className="absolute right-0 top-full mt-1 w-40 bg-gray-800 border border-gray-700 rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-10">
              {(['7', '30', '90', 'all'] as ExportPeriod[]).map((p) => (
                <button
                  key={p}
                  onClick={() => handleExport(p)}
                  className="w-full text-left px-3 py-2 text-sm text-gray-300 hover:text-white hover:bg-gray-700 first:rounded-t-lg last:rounded-b-lg"
                >
                  {p === 'all' ? t('learningLogs.exportAll') : t('learningLogs.exportDays', { days: p })}
                </button>
              ))}
            </div>
          </div>
          <button
            onClick={handleShowForm}
            className={`${buttonSecondaryClass} font-medium text-sm`}
          >
            {t('learningLogs.addLog')}
          </button>
        </div>
      </div>

      {/* Weekly Summary */}
      <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4 flex items-center gap-3">
          <Clock className="w-8 h-8 text-blue-400" />
          <div>
            <p className="text-xs text-gray-400">{t('learningLogs.weeklyDuration')}</p>
            <p className="text-lg font-bold text-white">
              {weeklyDuration >= 60
                ? t('learningLogs.hoursMinutes', { hours: Math.floor(weeklyDuration / 60), minutes: weeklyDuration % 60 })
                : t('learningLogs.durationMinutes', { minutes: weeklyDuration })}
            </p>
          </div>
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4 flex items-center gap-3">
          <Calendar className="w-8 h-8 text-orange-400" />
          <div>
            <p className="text-xs text-gray-400">{t('learningLogs.currentStreak')}</p>
            <p className="text-lg font-bold text-white">{streakInfo?.current_streak ?? 0}{t('learningLogs.days')}</p>
          </div>
        </div>
        <div className="bg-gray-900 border border-gray-800 rounded-md p-4 flex items-center gap-3 col-span-2 md:col-span-1">
          <FileText className="w-8 h-8 text-green-400" />
          <div>
            <p className="text-xs text-gray-400">{t('learningLogs.logCount')}</p>
            <p className="text-lg font-bold text-white">{filteredLogs.length}</p>
          </div>
        </div>
      </div>

      {/* View Toggle */}
      <div className="flex items-center gap-2">
        <button
          onClick={handleViewList}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
            view === 'list' ? 'bg-purple-500/20 text-purple-400' : 'text-gray-400 hover:text-white'
          }`}
        >
          <List className="w-4 h-4" />
          {t('learningLogs.list')}
        </button>
        <button
          onClick={handleViewCalendar}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
            view === 'calendar' ? 'bg-purple-500/20 text-purple-400' : 'text-gray-400 hover:text-white'
          }`}
        >
          <Calendar className="w-4 h-4" />
          {t('learningLogs.calendar')}
        </button>
        <button
          onClick={handleToggleFavoritesFilter}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
            showFavoritesOnly ? 'bg-yellow-500/20 text-yellow-400' : 'text-gray-400 hover:text-white'
          }`}
        >
          <Star className={`w-4 h-4 ${showFavoritesOnly ? 'fill-yellow-400' : ''}`} />
          {t('learningLogs.favorites')}
        </button>
        {filterDate && (
          <div className="flex items-center gap-2 ml-2">
            <span className="text-sm text-purple-400">
              {t('learningLogs.logsOnDate', { date: filterDate })}
            </span>
            <button
              onClick={clearFilterDate}
              className="text-xs text-gray-500 hover:text-white"
            >
              &times;
            </button>
          </div>
        )}
      </div>

      {/* Category Filter */}
      <div className="flex flex-col gap-2">
        <span className="text-sm text-gray-400">{t('learningLogs.category')}</span>
        <div className="flex flex-wrap gap-2">
          <button
            onClick={handleFilterAll}
            className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
              filterCategory === 'all'
                ? 'bg-purple-500/20 text-purple-400'
                : 'bg-gray-800/50 text-gray-400 hover:text-white'
            }`}
          >
            {t('learningLogs.filterAll')}
          </button>
          {CATEGORIES.map(({ value, label, Icon }) => (
            <button
              key={value}
              onClick={() => setFilterCategory(value)}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                filterCategory === value
                  ? 'bg-purple-500/20 text-purple-400'
                  : 'bg-gray-800/50 text-gray-400 hover:text-white'
              }`}
            >
              <Icon className="w-4 h-4" />
              {t(label)}
            </button>
          ))}
        </div>
      </div>

      {/* Calendar View */}
      {view === 'calendar' && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
          <LogCalendar entries={calendarData} onDateClick={handleDateClick} />
        </div>
      )}

      {/* Create/Edit Form Modal */}
      <Modal
        isOpen={showForm}
        onClose={resetForm}
        title={editingLog ? t('learningLogs.editLog') : t('learningLogs.addLog')}
        maxWidth="max-w-md"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="log-title" className={labelClass}>
              {t('learningLogs.logTitle')}
            </label>
            <input
              id="log-title"
              type="text"
              value={title}
              onChange={handleTitleChange}
              placeholder={t('learningLogs.titlePlaceholder')}
              className={inputClass}
              required
            />
          </div>
          <div>
            <label htmlFor="log-content" className={labelClass}>
              {t('learningLogs.content')}
            </label>
            <textarea
              id="log-content"
              value={content}
              onChange={handleContentChange}
              placeholder={t('learningLogs.contentPlaceholder')}
              rows={4}
              className={`${inputClass} resize-none`}
              required
            />
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label htmlFor="log-category" className={labelClass}>
                {t('learningLogs.category')}
              </label>
              <select
                id="log-category"
                value={category}
                onChange={handleCategoryChange}
                className={inputClass}
              >
                {CATEGORIES.map((cat) => (
                  <option key={cat.value} value={cat.value}>
                    {t(cat.label)}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label htmlFor="log-duration" className={labelClass}>
                {t('learningLogs.duration')}
              </label>
              <input
                id="log-duration"
                type="number"
                value={duration}
                onChange={handleDurationChange}
                placeholder={t('learningLogs.durationPlaceholder')}
                min="0"
                className={inputClass}
              />
            </div>
          </div>
          <div className="flex gap-3 justify-end pt-2">
            <button
              type="button"
              onClick={resetForm}
              className={`${buttonSecondaryClass} text-sm font-medium`}
            >
              {t('common.cancel')}
            </button>
            <button
              type="submit"
              disabled={saving || !title.trim() || !content.trim()}
              className="px-4 py-2 bg-purple-600 hover:bg-purple-500 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors"
            >
              {saving ? t('common.loading') : editingLog ? t('common.save') : t('learningLogs.addLog')}
            </button>
          </div>
        </form>
      </Modal>

      {/* Logs List */}
      {view === 'list' && (
        <>
          {filteredLogs.length === 0 ? (
            <div className="bg-gray-900 border border-gray-800 rounded-md p-12 text-center">
              <p className="text-gray-400 text-sm mb-4">
                {filterDate ? t('learningLogs.noLogsOnDate') : t('learningLogs.noLogs')}
              </p>
              {!filterDate && (
                <button
                  onClick={handleShowForm}
                  className={`${buttonSecondaryClass} font-medium text-sm`}
                >
                  {t('learningLogs.addLog')}
                </button>
              )}
            </div>
          ) : (
            <div className="space-y-3">
              {filteredLogs.map((log) => {
                const catInfo = getCategoryInfo(log.category);
                const CatIcon = catInfo.Icon;
                return (
                  <div key={log.id} className="bg-gray-900 border border-gray-800 rounded-md p-4">
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex items-start gap-3 min-w-0 flex-1">
                        <CatIcon className="w-5 h-5 text-purple-400 flex-shrink-0 mt-0.5" />
                        <div className="min-w-0 flex-1">
                          <div className="flex items-center gap-2 flex-wrap">
                            <h3 className="font-medium">{log.title}</h3>
                            <span className={`px-2 py-0.5 text-xs rounded-full ${getCategoryColor(log.category)}`}>
                              {t(catInfo.label)}
                            </span>
                          </div>
                          <p className="text-sm text-gray-400 mt-1 whitespace-pre-wrap">{log.content}</p>
                          <div className="flex items-center gap-4 mt-2 text-xs text-gray-500">
                            <span>{new Date(log.created_at).toLocaleDateString()}</span>
                            {log.duration > 0 && (
                              <span>{t('learningLogs.durationMinutes', { minutes: log.duration })}</span>
                            )}
                          </div>
                        </div>
                      </div>

                      <div className="flex items-center gap-1">
                        <button
                          onClick={() => toggleFavorite(log.id)}
                          className={`p-2 transition-colors ${log.is_favorite ? 'text-yellow-400' : 'text-gray-400 hover:text-yellow-400'}`}
                          title={t('learningLogs.toggleFavorite')}
                        >
                          <Star className={`w-4 h-4 ${log.is_favorite ? 'fill-yellow-400' : ''}`} />
                        </button>
                        <button
                          onClick={() => handleEdit(log)}
                          className="p-2 text-gray-400 hover:text-blue-400 transition-colors"
                          title={t('common.edit')}
                        >
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
                          </svg>
                        </button>
                        <button
                          onClick={() => handleDelete(log.id)}
                          className="p-2 text-gray-400 hover:text-red-400 transition-colors"
                          title={t('common.delete')}
                        >
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                          </svg>
                        </button>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </>
      )}
    </div>
  );
}
