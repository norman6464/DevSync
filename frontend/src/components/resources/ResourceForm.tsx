import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { LearningResource, CreateResourceRequest, ResourceCategory, ResourceDifficulty } from '../../types/resource';
import { buttonSecondaryClass, inputClass, labelClass, textareaClass } from '../../constants/styles';
import { findInvalidUrlField } from '../../utils/url';
import TagInput from '../common/TagInput';
import toast from 'react-hot-toast';

interface ResourceFormProps {
  resource?: LearningResource;
  onSubmit: (data: CreateResourceRequest) => Promise<void>;
  onCancel: () => void;
  loading?: boolean;
}

const categories: ResourceCategory[] = ['book', 'video', 'article', 'course', 'tutorial', 'podcast', 'tool', 'other'];
const difficulties: ResourceDifficulty[] = ['beginner', 'intermediate', 'advanced'];

export default function ResourceForm({ resource, onSubmit, onCancel, loading }: ResourceFormProps) {
  const { t } = useTranslation();
  const [title, setTitle] = useState(resource?.title || '');
  const [description, setDescription] = useState(resource?.description || '');
  const [url, setUrl] = useState(resource?.url || '');
  const [category, setCategory] = useState<ResourceCategory>(resource?.category || 'article');
  const [difficulty, setDifficulty] = useState<ResourceDifficulty | ''>(resource?.difficulty || '');
  const [tags, setTags] = useState<string[]>(() => {
    if (resource?.tags) {
      try {
        return JSON.parse(resource.tags);
      } catch {
        return [];
      }
    }
    return [];
  });
  const [imageUrl, setImageUrl] = useState(resource?.image_url || '');
  const [isPublic, setIsPublic] = useState(resource?.is_public ?? true);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const invalidField = findInvalidUrlField([
      { value: url, label: t('resources.url') },
      { value: imageUrl, label: t('resources.imageUrl') },
    ]);
    if (invalidField) {
      toast.error(t('common.invalidUrl', { field: invalidField }));
      return;
    }
    await onSubmit({
      title,
      description,
      url,
      category,
      difficulty: difficulty || undefined,
      tags: JSON.stringify(tags),
      image_url: imageUrl,
      is_public: isPublic,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Title */}
      <div>
        <label htmlFor="resource-title" className={labelClass}>
          {t('resources.title')} *
        </label>
        <input
          id="resource-title"
          type="text"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          maxLength={300}
          className={inputClass}
          placeholder={t('resources.titlePlaceholder')}
        />
      </div>

      {/* URL */}
      <div>
        <label htmlFor="resource-url" className={labelClass}>
          {t('resources.url')}
        </label>
        <input
          id="resource-url"
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          maxLength={2048}
          className={inputClass}
          placeholder="https://..."
        />
      </div>

      {/* Category & Difficulty */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label htmlFor="resource-category" className={labelClass}>
            {t('resources.category')} *
          </label>
          <select
            id="resource-category"
            value={category}
            onChange={(e) => setCategory(e.target.value as ResourceCategory)}
            className={inputClass}
          >
            {categories.map(cat => (
              <option key={cat} value={cat}>
                {t(`resources.categories.${cat}`)}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="resource-difficulty" className={labelClass}>
            {t('resources.difficultyLabel')}
          </label>
          <select
            id="resource-difficulty"
            value={difficulty}
            onChange={(e) => setDifficulty(e.target.value as ResourceDifficulty | '')}
            className={inputClass}
          >
            <option value="">{t('resources.selectDifficulty')}</option>
            {difficulties.map(diff => (
              <option key={diff} value={diff}>
                {t(`resources.difficulty.${diff}`)}
              </option>
            ))}
          </select>
        </div>
      </div>

      {/* Description */}
      <div>
        <label htmlFor="resource-description" className={labelClass}>
          {t('resources.description')}
        </label>
        <textarea
          id="resource-description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          rows={3}
          maxLength={1000}
          className={textareaClass}
          placeholder={t('resources.descriptionPlaceholder')}
        />
      </div>

      {/* Tags */}
      <TagInput
        tags={tags}
        onChange={setTags}
        label={t('resources.tags')}
        id="resource-tags"
        placeholder={t('resources.tagsPlaceholder')}
        maxLength={50}
        prefix="#"
      />

      {/* Image URL */}
      <div>
        <label htmlFor="resource-image-url" className={labelClass}>
          {t('resources.imageUrl')}
        </label>
        <input
          id="resource-image-url"
          type="url"
          value={imageUrl}
          onChange={(e) => setImageUrl(e.target.value)}
          maxLength={2048}
          className={inputClass}
          placeholder="https://..."
        />
      </div>

      {/* Public */}
      <div className="flex items-center gap-2">
        <input
          type="checkbox"
          id="isPublic"
          checked={isPublic}
          onChange={(e) => setIsPublic(e.target.checked)}
          className="w-4 h-4 rounded bg-gray-700 border-gray-600 text-green-500 focus:ring-gray-500"
        />
        <label htmlFor="isPublic" className="text-sm text-gray-300">
          {t('resources.makePublic')}
        </label>
      </div>

      {/* Buttons */}
      <div className="flex gap-3 pt-4">
        <button
          type="button"
          onClick={onCancel}
          className={`flex-1 ${buttonSecondaryClass}`}
        >
          {t('common.cancel')}
        </button>
        <button
          type="submit"
          disabled={loading || !title.trim()}
          className={`flex-1 ${buttonSecondaryClass} disabled:bg-gray-600 disabled:cursor-not-allowed`}
        >
          {loading ? t('common.saving') : t('common.save')}
        </button>
      </div>
    </form>
  );
}
