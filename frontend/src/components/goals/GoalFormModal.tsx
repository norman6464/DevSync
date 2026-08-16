import { useTranslation } from 'react-i18next';
import { type GoalCategory } from '../../api/goals';
import { Modal } from '../common';
import { CATEGORIES } from './goalCategories';
import { inputClass, buttonSecondaryClass, labelClass } from '../../constants/styles';
import { MAX_LENGTH } from '../../utils/formValidation';
import { CharCount } from '../common';

interface GoalFormModalProps {
  isOpen: boolean;
  isEditing: boolean;
  saving: boolean;
  title: string;
  setTitle: (v: string) => void;
  description: string;
  setDescription: (v: string) => void;
  category: GoalCategory;
  setCategory: (v: GoalCategory) => void;
  targetDate: string;
  setTargetDate: (v: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  onCancel: () => void;
}

export default function GoalFormModal({
  isOpen, isEditing, saving,
  title, setTitle, description, setDescription,
  category, setCategory, targetDate, setTargetDate,
  onSubmit, onCancel,
}: GoalFormModalProps) {
  const { t } = useTranslation();

  return (
    <Modal
      isOpen={isOpen}
      onClose={onCancel}
      title={isEditing ? t('goals.editGoal') : t('goals.addGoal')}
      maxWidth="max-w-md"
    >
      <form onSubmit={onSubmit} className="space-y-4">
        <div>
          <label htmlFor="goal-title" className={labelClass}>
            {t('goals.goalTitle')}
          </label>
          <input
            id="goal-title"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder={t('goals.titlePlaceholder')}
            maxLength={MAX_LENGTH.goalTitle}
            className={inputClass}
            required
          />
        </div>
        <div>
          <label htmlFor="goal-description" className={labelClass}>
            {t('goals.description')}
          </label>
          <textarea
            id="goal-description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder={t('goals.descriptionPlaceholder')}
            rows={3}
            maxLength={MAX_LENGTH.goalDescription}
            className={`${inputClass} resize-none`}
          />
          <CharCount value={description} max={MAX_LENGTH.goalDescription} />
        </div>
        <div>
          <label htmlFor="goal-category" className={labelClass}>
            {t('goals.category')}
          </label>
          <select
            id="goal-category"
            value={category}
            onChange={(e) => setCategory(e.target.value as GoalCategory)}
            className={inputClass}
          >
            {CATEGORIES.map((cat) => (
              <option key={cat.value} value={cat.value}>
                {cat.icon} {t(cat.label)}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="goal-target-date" className={labelClass}>
            {t('goals.targetDate')}
          </label>
          <input
            id="goal-target-date"
            type="date"
            value={targetDate}
            onChange={(e) => setTargetDate(e.target.value)}
            className={inputClass}
          />
        </div>
        <div className="flex gap-3 justify-end pt-2">
          <button
            type="button"
            onClick={onCancel}
            className={`${buttonSecondaryClass} text-sm font-medium`}
          >
            {t('common.cancel')}
          </button>
          <button
            type="submit"
            disabled={saving || !title.trim()}
            className={`${buttonSecondaryClass} disabled:opacity-50 text-sm font-medium`}
          >
            {saving ? t('common.loading') : isEditing ? t('common.save') : t('goals.create')}
          </button>
        </div>
      </form>
    </Modal>
  );
}
