import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { X, Sparkles, Clock } from 'lucide-react';
import type { GoalCategory } from '../../api/goals';
import { GOAL_TEMPLATES, GOAL_CATEGORIES, type GoalTemplate, getTemplatesByCategory } from '../../constants/goalTemplates';
import { buttonPrimaryClass, buttonSecondaryClass } from '../../constants/styles';

interface GoalTemplatesModalProps {
  isOpen: boolean;
  onSelect: (template: GoalTemplate) => void;
  onClose: () => void;
}

export default function GoalTemplatesModal({ isOpen, onSelect, onClose }: GoalTemplatesModalProps) {
  const { t } = useTranslation();
  const [selectedCategory, setSelectedCategory] = useState<GoalCategory | 'all'>('all');

  const filteredTemplates = getTemplatesByCategory(selectedCategory);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-gray-900 border border-gray-800 rounded-lg w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-800">
          <div className="flex items-center gap-2">
            <Sparkles className="w-5 h-5 text-yellow-400" />
            <h2 className="text-lg font-semibold">テンプレートから目標を作成</h2>
          </div>
          <button
            onClick={onClose}
            aria-label="閉じる"
            className="p-1 hover:bg-gray-800 rounded transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Category Filters */}
        <div className="p-4 border-b border-gray-800">
          <div className="flex gap-2 flex-wrap">
            <button
              onClick={() => setSelectedCategory('all')}
              className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                selectedCategory === 'all'
                  ? 'bg-blue-500/20 text-blue-400 border border-blue-400/30'
                  : 'bg-gray-800/50 text-gray-400 hover:text-white border border-gray-700'
              }`}
            >
              すべて
            </button>
            {GOAL_CATEGORIES.map(({ value, label, Icon }) => (
              <button
                key={value}
                onClick={() => setSelectedCategory(value)}
                className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                  selectedCategory === value
                    ? 'bg-blue-500/20 text-blue-400 border border-blue-400/30'
                    : 'bg-gray-800/50 text-gray-400 hover:text-white border border-gray-700'
                }`}
              >
                <Icon className="w-3.5 h-3.5" />
                {label}
              </button>
            ))}
          </div>
        </div>

        {/* Templates Grid */}
        <div className="flex-1 overflow-y-auto p-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
            {filteredTemplates.map((template) => {
              const categoryInfo = GOAL_CATEGORIES.find((c) => c.value === template.category);
              const Icon = categoryInfo?.Icon;

              return (
                <button
                  key={template.id}
                  onClick={() => onSelect(template)}
                  className="bg-gray-800/50 border border-gray-700 rounded-lg p-4 text-left hover:border-blue-400/50 hover:bg-gray-800 transition-all group"
                >
                  <div className="flex items-start gap-3">
                    {Icon && (
                      <div className="w-10 h-10 bg-gradient-to-br from-blue-500/20 to-purple-500/20 rounded-lg flex items-center justify-center flex-shrink-0 group-hover:from-blue-500/30 group-hover:to-purple-500/30 transition-colors">
                        <Icon className="w-5 h-5 text-blue-400" />
                      </div>
                    )}
                    <div className="flex-1 min-w-0">
                      <h3 className="font-medium text-white mb-1 group-hover:text-blue-400 transition-colors">
                        {template.title}
                      </h3>
                      <p className="text-xs text-gray-400 mb-2 line-clamp-2">
                        {template.description}
                      </p>
                      <div className="flex items-center gap-3">
                        <span className="inline-flex items-center gap-1 px-2 py-0.5 text-xs bg-gray-700 text-gray-300 rounded">
                          {categoryInfo?.label}
                        </span>
                        <span className="inline-flex items-center gap-1 text-xs text-gray-500">
                          <Clock className="w-3 h-3" />
                          {template.estimatedDays}日
                        </span>
                      </div>
                    </div>
                  </div>
                </button>
              );
            })}
          </div>

          {filteredTemplates.length === 0 && (
            <div className="text-center py-12 text-gray-500">
              このカテゴリにテンプレートはありません
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="p-4 border-t border-gray-800 flex justify-end">
          <button onClick={onClose} className={buttonSecondaryClass}>
            キャンセル
          </button>
        </div>
      </div>
    </div>
  );
}
