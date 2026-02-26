import { useTranslation } from 'react-i18next';
import type { WidgetConfig } from '../../types/widgetSettings';

interface WidgetSettingsPanelProps {
  widgets: WidgetConfig[];
  saving: boolean;
  onToggleVisibility: (key: string) => void;
  onMoveUp: (key: string) => void;
  onMoveDown: (key: string) => void;
  onSave: () => void;
  onCancel: () => void;
}

export default function WidgetSettingsPanel({
  widgets,
  saving,
  onToggleVisibility,
  onMoveUp,
  onMoveDown,
  onSave,
  onCancel,
}: WidgetSettingsPanelProps) {
  const { t } = useTranslation();

  return (
    <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-4">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-semibold text-white">{t('widgetSettings.title')}</h3>
        <div className="flex gap-2">
          <button
            onClick={onCancel}
            className="px-3 py-1 text-xs text-gray-400 hover:text-white border border-gray-600 rounded transition-colors"
          >
            {t('common.cancel')}
          </button>
          <button
            onClick={onSave}
            disabled={saving}
            className="px-3 py-1 text-xs bg-blue-600 hover:bg-blue-500 text-white rounded transition-colors disabled:opacity-50"
          >
            {saving ? t('common.saving') : t('common.save')}
          </button>
        </div>
      </div>
      <ul className="space-y-1">
        {widgets.map((widget, idx) => (
          <li
            key={widget.key}
            className="flex items-center gap-2 px-2 py-1.5 rounded hover:bg-gray-700/50 transition-colors"
          >
            <button
              onClick={() => onToggleVisibility(widget.key)}
              className={`w-5 h-5 flex items-center justify-center rounded border transition-colors ${
                widget.visible
                  ? 'bg-blue-600 border-blue-500 text-white'
                  : 'border-gray-600 text-gray-500'
              }`}
              aria-label={widget.visible ? t('widgetSettings.hide') : t('widgetSettings.show')}
            >
              {widget.visible && (
                <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="3" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                </svg>
              )}
            </button>
            <span className={`flex-1 text-sm ${widget.visible ? 'text-white' : 'text-gray-500'}`}>
              {t(`widgetSettings.widgets.${widget.key}`)}
            </span>
            <div className="flex gap-0.5">
              <button
                onClick={() => onMoveUp(widget.key)}
                disabled={idx === 0}
                className="p-1 text-gray-400 hover:text-white disabled:text-gray-700 transition-colors"
                aria-label={t('widgetSettings.moveUp')}
              >
                <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M5 15l7-7 7 7" />
                </svg>
              </button>
              <button
                onClick={() => onMoveDown(widget.key)}
                disabled={idx === widgets.length - 1}
                className="p-1 text-gray-400 hover:text-white disabled:text-gray-700 transition-colors"
                aria-label={t('widgetSettings.moveDown')}
              >
                <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
