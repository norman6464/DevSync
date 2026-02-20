import { useState, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { Bell, CheckCheck } from 'lucide-react';
import { useNotifications } from '../hooks';
import NotificationFilters from '../components/notifications/NotificationFilters';
import NotificationItem from '../components/notifications/NotificationItem';
import EmptyState from '../components/common/EmptyState';
import LoadingSpinner from '../components/common/LoadingSpinner';
import { buttonSecondaryClass } from '../constants/styles';

export default function NotificationsPage() {
  const { t } = useTranslation();
  const [showUnreadOnly, setShowUnreadOnly] = useState(false);
  const {
    notifications, unreadCount, total, loading,
    page, setPage, limit,
    filterType, setFilterType,
    markAsRead, markAllAsRead, deleteNotification,
  } = useNotifications();

  const filteredNotifications = showUnreadOnly
    ? notifications.filter((n) => !n.read)
    : notifications;

  const handleToggleUnreadOnly = useCallback(() => setShowUnreadOnly((prev) => !prev), []);
  const handlePreviousPage = useCallback(() => setPage(page - 1), [page, setPage]);
  const handleNextPage = useCallback(() => setPage(page + 1), [page, setPage]);

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-white">{t('notifications.pageTitle')}</h1>
          <p className="text-gray-400 text-sm mt-1">{t('notifications.pageSubtitle')}</p>
        </div>
        {unreadCount > 0 && (
          <button
            onClick={markAllAsRead}
            className={`flex items-center gap-2 ${buttonSecondaryClass}`}
          >
            <CheckCheck className="w-5 h-5" aria-hidden="true" />
            {t('notifications.markAllRead')}
          </button>
        )}
      </div>

      <NotificationFilters
        filterType={filterType}
        setFilterType={setFilterType}
        showUnreadOnly={showUnreadOnly}
        onToggleUnreadOnly={handleToggleUnreadOnly}
      />

      {loading ? (
        <div className="flex justify-center items-center min-h-[400px]">
          <LoadingSpinner />
        </div>
      ) : filteredNotifications.length === 0 ? (
        <EmptyState
          icon={Bell}
          message={t('notifications.empty')}
        />
      ) : (
        <>
          <ul className="space-y-2" role="list" aria-label={t('notifications.listLabel')}>
            {filteredNotifications.map((notification) => (
              <NotificationItem
                key={notification.id}
                notification={notification}
                onMarkAsRead={markAsRead}
                onDelete={deleteNotification}
              />
            ))}
          </ul>

          {total > limit && (
            <div className="flex justify-center gap-2 mt-8">
              <button
                onClick={handlePreviousPage}
                disabled={page <= 1}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {t('common.previous')}
              </button>
              <span className="px-4 py-2 text-gray-400">
                {page} / {totalPages}
              </span>
              <button
                onClick={handleNextPage}
                disabled={page >= totalPages}
                className="px-4 py-2 bg-gray-800 hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed text-white rounded-lg transition-colors"
              >
                {t('common.next')}
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
