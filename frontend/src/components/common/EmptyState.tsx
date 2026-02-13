import type { LucideIcon } from 'lucide-react';

interface EmptyStateProps {
  icon: LucideIcon;
  title?: string;
  message: string;
  actionLabel?: string;
  onAction?: () => void;
}

/**
 * EmptyState - データがない状態を表示する統一コンポーネント
 *
 * @param icon - 表示するLucideアイコン
 * @param title - 見出し（オプション）
 * @param message - メインメッセージ
 * @param actionLabel - CTAボタンのラベル（オプション）
 * @param onAction - CTAボタンのクリックハンドラー（オプション）
 */
export default function EmptyState({
  icon: Icon,
  title,
  message,
  actionLabel,
  onAction,
}: EmptyStateProps) {
  return (
    <div className="text-center py-12 px-4">
      {/* Icon Container */}
      <div className="inline-flex items-center justify-center w-16 h-16 mb-4 bg-gray-800 rounded-full">
        <Icon className="w-8 h-8 text-gray-500" />
      </div>

      {/* Title */}
      {title && (
        <h3 className="text-lg font-semibold text-white mb-2">
          {title}
        </h3>
      )}

      {/* Message */}
      <p className="text-gray-400 text-sm mb-6 max-w-md mx-auto">
        {message}
      </p>

      {/* Action Button */}
      {actionLabel && onAction && (
        <button
          onClick={onAction}
          className="px-6 py-2.5 bg-gray-700 hover:bg-gray-600 text-white rounded-lg font-medium text-sm transition-colors"
        >
          {actionLabel}
        </button>
      )}
    </div>
  );
}
