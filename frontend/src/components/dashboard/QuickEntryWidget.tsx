import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { BookOpen, Code, GraduationCap, Users, Layers, Clock, Check } from 'lucide-react';
import { useQuickEntry } from '../../hooks/useQuickEntry';
import type { LogCategory } from '../../types/learningLog';
import { panelClass } from '../../constants/styles';

const CATEGORIES: { id: LogCategory; icon: React.ElementType; color: string }[] = [
  { id: 'coding', icon: Code, color: 'text-blue-400 border-blue-500/50 bg-blue-500/10' },
  { id: 'reading', icon: BookOpen, color: 'text-green-400 border-green-500/50 bg-green-500/10' },
  { id: 'course', icon: GraduationCap, color: 'text-purple-400 border-purple-500/50 bg-purple-500/10' },
  { id: 'meetup', icon: Users, color: 'text-orange-400 border-orange-500/50 bg-orange-500/10' },
  { id: 'other', icon: Layers, color: 'text-gray-400 border-gray-500/50 bg-gray-500/10' },
];

const DURATION_PRESETS = [15, 30, 60, 90, 120];

export default function QuickEntryWidget() {
  const { t } = useTranslation();
  const { recentCategories, submitting, submit } = useQuickEntry();
  const [selectedCategory, setSelectedCategory] = useState<LogCategory | null>(null);
  const [duration, setDuration] = useState<number>(30);
  const [note, setNote] = useState('');
  const [showSuccess, setShowSuccess] = useState(false);

  const sortedCategories = [...CATEGORIES].sort((a, b) => {
    const aIdx = recentCategories.indexOf(a.id);
    const bIdx = recentCategories.indexOf(b.id);
    if (aIdx === -1 && bIdx === -1) return 0;
    if (aIdx === -1) return 1;
    if (bIdx === -1) return -1;
    return aIdx - bIdx;
  });

  const handleSubmit = async () => {
    if (!selectedCategory) return;
    const ok = await submit(selectedCategory, duration, note);
    if (ok) {
      setShowSuccess(true);
      setSelectedCategory(null);
      setDuration(30);
      setNote('');
      setTimeout(() => setShowSuccess(false), 2000);
    }
  };

  return (
    <div className={panelClass}>
      <div className="flex items-center gap-2 mb-4">
        <Clock className="w-4 h-4 text-blue-400" />
        <h3 className="text-sm font-medium text-white">{t('quickEntry.title')}</h3>
      </div>

      {showSuccess ? (
        <div className="flex flex-col items-center justify-center py-6 animate-pulse">
          <div className="w-12 h-12 rounded-full bg-green-500/20 flex items-center justify-center mb-2">
            <Check className="w-6 h-6 text-green-400" />
          </div>
          <p className="text-sm text-green-400">{t('quickEntry.success')}</p>
        </div>
      ) : (
        <div className="space-y-3">
          {/* Category Selection */}
          <div>
            <p className="text-xs text-gray-500 mb-2">{t('quickEntry.selectCategory')}</p>
            <div className="flex flex-wrap gap-1.5">
              {sortedCategories.map(({ id, icon: Icon, color }) => (
                <button
                  key={id}
                  onClick={() => setSelectedCategory(id)}
                  className={`flex items-center gap-1 px-2.5 py-1.5 text-xs rounded-lg border transition-all ${
                    selectedCategory === id
                      ? color
                      : 'border-gray-700 text-gray-400 hover:border-gray-600'
                  }`}
                >
                  <Icon className="w-3.5 h-3.5" />
                  {t(`quickEntry.categories.${id}`)}
                </button>
              ))}
            </div>
          </div>

          {/* Duration */}
          {selectedCategory && (
            <>
              <div>
                <p className="text-xs text-gray-500 mb-2">{t('quickEntry.duration')}</p>
                <div className="flex flex-wrap gap-1.5">
                  {DURATION_PRESETS.map((d) => (
                    <button
                      key={d}
                      onClick={() => setDuration(d)}
                      className={`px-2.5 py-1.5 text-xs rounded-lg border transition-all ${
                        duration === d
                          ? 'border-blue-500/50 bg-blue-500/10 text-blue-400'
                          : 'border-gray-700 text-gray-400 hover:border-gray-600'
                      }`}
                    >
                      {d}{t('quickEntry.minutes')}
                    </button>
                  ))}
                </div>
              </div>

              {/* Note */}
              <div>
                <input
                  type="text"
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                  placeholder={t('quickEntry.notePlaceholder')}
                  className="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-200 placeholder-gray-500 focus:border-blue-500 focus:outline-none"
                />
              </div>

              {/* Submit */}
              <button
                onClick={handleSubmit}
                disabled={submitting}
                className="w-full py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-500 disabled:bg-gray-700 disabled:text-gray-500 rounded-lg transition-colors"
              >
                {submitting ? t('quickEntry.submitting') : t('quickEntry.submit')}
              </button>
            </>
          )}
        </div>
      )}
    </div>
  );
}
