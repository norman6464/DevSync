import i18n from '../i18n';

export function formatDate(dateStr: string | null | undefined): string | null {
  if (!dateStr) return null;
  return new Date(dateStr).toLocaleDateString();
}

export function formatDistanceToNow(dateString: string): string {
  const t = i18n.t.bind(i18n);
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / 60000);
  const diffHours = Math.floor(diffMs / 3600000);
  const diffDays = Math.floor(diffMs / 86400000);

  if (diffMins < 1) return t('time.justNow');
  if (diffMins < 60) return t('time.minutesAgo', { count: diffMins });
  if (diffHours < 24) return t('time.hoursAgo', { count: diffHours });
  return t('time.daysAgo', { count: diffDays });
}

/** dateString から windowMs 以内しか経過していないか（NEW バッジ等の判定用）。 */
export function isWithinLast(dateString: string, windowMs: number): boolean {
  return Date.now() - new Date(dateString).getTime() < windowMs;
}

/** dateString が現在時刻より過去か（期限超過の判定用）。 */
export function isPastDate(dateString: string): boolean {
  return new Date(dateString).getTime() < Date.now();
}

/** 開始日〜終了日（未設定なら現在まで）の日数。開始日が無ければ null。 */
export function calcDurationDays(startDate: string | null | undefined, endDate: string | null | undefined): number | null {
  if (!startDate) return null;
  const start = new Date(startDate).getTime();
  const end = endDate ? new Date(endDate).getTime() : Date.now();
  return Math.max(1, Math.ceil((end - start) / (1000 * 60 * 60 * 24)));
}
