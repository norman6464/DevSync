import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Calendar, List, Download, Star, ArrowDownWideNarrow, Search } from 'lucide-react';
import { useLearningLogForm } from '../hooks/useLearningLogForm';
import { useWeeklyDuration, useStreak } from '../hooks';
import { useAuthStore } from '../store/authStore';
import { exportLogsCSV, type ExportPeriod } from '../api/learningLogs';
import type { LogCategory } from '../types/learningLog';
import LogCalendar from '../components/learning-logs/LogCalendar';
import LogCard, { CATEGORIES } from '../components/learning-logs/LogCard';
import WeeklySummaryCard from '../components/learning-logs/WeeklySummaryCard';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { Modal } from '../components/common';
import { inputClass, buttonSecondaryClass, labelClass } from '../constants/styles';


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
    sortBy, setSortBy,
    searchQuery, setSearchQuery,
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
  const handleSearchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value), [setSearchQuery]);

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
      <WeeklySummaryCard
        weeklyDuration={weeklyDuration}
        streakInfo={streakInfo}
        logCount={filteredLogs.length}
      />

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
        <input
          type="text"
          value={searchQuery}
          onChange={handleSearchChange}
          placeholder={t('learningLogs.searchPlaceholder')}
          className={`${inputClass} pl-9`}
        />
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

      {/* Sort */}
      <div className="flex flex-col gap-2">
        <span className="text-sm text-gray-400 flex items-center gap-1.5">
          <ArrowDownWideNarrow className="w-4 h-4" />
          {t('learningLogs.sort')}
        </span>
        <div className="flex flex-wrap gap-2">
          {([
            { value: 'latest', label: 'learningLogs.sortLatest' },
            { value: 'oldest', label: 'learningLogs.sortOldest' },
            { value: 'duration_desc', label: 'learningLogs.sortDurationDesc' },
            { value: 'duration_asc', label: 'learningLogs.sortDurationAsc' },
          ] as const).map((opt) => (
            <button
              key={opt.value}
              onClick={() => setSortBy(opt.value)}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                sortBy === opt.value
                  ? 'bg-purple-500/20 text-purple-400'
                  : 'bg-gray-800/50 text-gray-400 hover:text-white'
              }`}
            >
              {t(opt.label)}
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
              maxLength={200}
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
              maxLength={5000}
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
              {filteredLogs.map((log) => (
                <LogCard
                  key={log.id}
                  log={log}
                  onEdit={handleEdit}
                  onDelete={handleDelete}
                  onToggleFavorite={toggleFavorite}
                />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
}
