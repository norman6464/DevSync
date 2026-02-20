import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Download, Search } from 'lucide-react';
import { useLearningLogForm } from '../hooks/useLearningLogForm';
import { useWeeklyDuration, useStreak } from '../hooks';
import { useAuthStore } from '../store/authStore';
import { exportLogsCSV, type ExportPeriod } from '../api/learningLogs';
import LogCalendar from '../components/learning-logs/LogCalendar';
import LogCard from '../components/learning-logs/LogCard';
import LogFiltersBar from '../components/learning-logs/LogFiltersBar';
import LogFormModal from '../components/learning-logs/LogFormModal';
import WeeklySummaryCard from '../components/learning-logs/WeeklySummaryCard';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { inputClass, buttonSecondaryClass } from '../constants/styles';


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
  const handleToggleFavoritesFilter = useCallback(() => setShowFavoritesOnly(prev => !prev), [setShowFavoritesOnly]);
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

      <LogFiltersBar
        view={view}
        filterCategory={filterCategory}
        showFavoritesOnly={showFavoritesOnly}
        sortBy={sortBy}
        filterDate={filterDate}
        onViewList={handleViewList}
        onViewCalendar={handleViewCalendar}
        onToggleFavorites={handleToggleFavoritesFilter}
        onFilterCategory={setFilterCategory}
        onSortBy={setSortBy}
        onClearFilterDate={clearFilterDate}
      />

      {/* Calendar View */}
      {view === 'calendar' && (
        <div className="bg-gray-900 border border-gray-800 rounded-md p-6">
          <LogCalendar entries={calendarData} onDateClick={handleDateClick} />
        </div>
      )}

      <LogFormModal
        isOpen={showForm}
        editingLog={editingLog}
        title={title}
        setTitle={setTitle}
        content={content}
        setContent={setContent}
        category={category}
        setCategory={setCategory}
        duration={duration}
        setDuration={setDuration}
        saving={saving}
        onSubmit={handleSubmit}
        onClose={resetForm}
      />

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
