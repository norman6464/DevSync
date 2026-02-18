import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { LearningResource, CreateResourceRequest, ResourceCategory, ResourceDifficulty } from '../../types/resource';
import { buttonSecondaryClass, inputClass, labelClass, textareaClass } from '../../constants/styles';

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
  const [tagInput, setTagInput] = useState('');
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

  const addTag = () => {
    if (tagInput.trim() && !tags.includes(tagInput.trim())) {
      setTags([...tags, tagInput.trim()]);
      setTagInput('');
    }
  };

  const removeTag = (tag: string) => {
    setTags(tags.filter(t => t !== tag));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
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
          className={textareaClass}
          placeholder={t('resources.descriptionPlaceholder')}
        />
      </div>

      {/* Tags */}
      <div>
        <label htmlFor="resource-tags" className={labelClass}>
          {t('resources.tags')}
        </label>
        <div className="flex gap-2">
          <input
            id="resource-tags"
            type="text"
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            onKeyPress={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
            className={`${inputClass} flex-1`}
            placeholder={t('resources.tagsPlaceholder')}
          />
          <button
            type="button"
            onClick={addTag}
            className="px-4 py-2 bg-gray-600 hover:bg-gray-500 text-white rounded-lg transition-colors"
          >
            {t('common.add')}
          </button>
        </div>
        {tags.length > 0 && (
          <div className="flex flex-wrap gap-2 mt-2">
            {tags.map((tag, idx) => (
              <span
                key={idx}
                className="inline-flex items-center gap-1 px-2 py-1 bg-gray-700 text-gray-300 text-sm rounded"
              >
                #{tag}
                <button
                  type="button"
                  onClick={() => removeTag(tag)}
                  className="text-gray-400 hover:text-white"
                >
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12" />
                  </svg>
                </button>
              </span>
            ))}
          </div>
        )}
      </div>

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
