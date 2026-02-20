import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import toast from 'react-hot-toast';
import { Plus } from 'lucide-react';
import { inputClass } from '../../constants/styles';
import { useConfirm } from '../../hooks';
import ConfirmDialog from '../common/ConfirmDialog';
import MarkdownEditor from './MarkdownEditor';
import CodeSnippetInput, { type SnippetInputData } from './CodeSnippetInput';

interface PostFormProps {
  post?: { id: number; title: string; content: string; image_urls?: string };
  onSubmit: (
    title: string,
    content: string,
    imageUrls?: string,
    codeSnippets?: { language: string; file_name?: string; code: string }[],
    isDraft?: boolean
  ) => Promise<void>;
  onCancel?: () => void;
  loading?: boolean;
}

export default function PostForm({ post, onSubmit, onCancel, loading: externalLoading }: PostFormProps) {
  const { t } = useTranslation();
  const [title, setTitle] = useState(post?.title || '');
  const [content, setContent] = useState(post?.content || '');
  const [imageUrls, setImageUrls] = useState<string[]>(
    post?.image_urls ? JSON.parse(post.image_urls) : []
  );
  const [snippets, setSnippets] = useState<SnippetInputData[]>([]);
  const [loading, setLoading] = useState(false);
  const [expanded, setExpanded] = useState(!!post);
  const { confirm, dialogProps } = useConfirm();

  const addSnippet = () => {
    setSnippets([...snippets, { language: '', file_name: '', code: '' }]);
  };

  const updateSnippet = (index: number, value: SnippetInputData) => {
    const updated = [...snippets];
    updated[index] = value;
    setSnippets(updated);
  };

  const removeSnippet = (index: number) => {
    setSnippets(snippets.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent, isDraft = false) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;
    setLoading(true);
    try {
      const imageUrlsJson = imageUrls.length > 0 ? JSON.stringify(imageUrls) : undefined;
      const validSnippets = snippets
        .filter((s) => s.language && s.code.trim())
        .map((s) => ({
          language: s.language,
          file_name: s.file_name || undefined,
          code: s.code,
        }));
      await onSubmit(title, content, imageUrlsJson, validSnippets.length > 0 ? validSnippets : undefined, isDraft);
      if (!post) {
        setTitle('');
        setContent('');
        setImageUrls([]);
        setSnippets([]);
        setExpanded(false);
      }
      toast.success(post ? t('post.postUpdated') : (isDraft ? t('post.draftSaved') : t('post.postCreated')));
    } catch {
      toast.error(post ? t('post.updateFailed') : (isDraft ? t('post.draftFailed') : t('post.postFailed')));
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-gray-900 border border-gray-800 rounded-md p-4">
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        onFocus={() => setExpanded(true)}
        placeholder={t('post.createTitle')}
        maxLength={300}
        className={inputClass}
      />
      <p className="text-xs text-gray-500 text-right mt-1">{title.length}/300</p>
      {expanded && (
        <>
          <div className="mt-3">
            <MarkdownEditor
              value={content}
              onChange={setContent}
              placeholder={t('post.createContent')}
              minHeight="150px"
              onImagesChange={setImageUrls}
            />
          </div>

          {/* Code Snippets Section */}
          {snippets.length > 0 && (
            <div className="mt-3 space-y-3" role="group" aria-label={t('post.addSnippet')}>
              {snippets.map((snippet, i) => (
                <CodeSnippetInput
                  key={i}
                  index={i}
                  value={snippet}
                  onChange={(val) => updateSnippet(i, val)}
                  onRemove={() => removeSnippet(i)}
                />
              ))}
            </div>
          )}

          <div className="flex justify-between items-center mt-3">
            <div className="flex items-center gap-3">
              <div className="text-xs text-gray-500">
                {t('post.markdownSupported')}
              </div>
              <button
                type="button"
                onClick={addSnippet}
                className="flex items-center gap-1 px-2.5 py-1.5 text-xs text-blue-400 hover:text-blue-300 hover:bg-blue-500/10 rounded-lg transition-colors"
              >
                <Plus className="w-3.5 h-3.5" aria-hidden="true" />
                {t('post.addSnippet')}
              </button>
            </div>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={async () => {
                  if (title.trim() || content.trim()) {
                    const confirmed = await confirm({
                      title: t('common.confirm'),
                      message: t('post.confirmDiscard'),
                      variant: 'warning',
                      confirmText: t('common.discard'),
                    });
                    if (!confirmed) return;
                  }
                  if (onCancel) {
                    onCancel();
                  } else {
                    setExpanded(false);
                    setTitle('');
                    setContent('');
                    setImageUrls([]);
                    setSnippets([]);
                  }
                }}
                className="px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors rounded-lg"
              >
                {t('common.cancel')}
              </button>
              {!post && (
                <button
                  type="button"
                  onClick={(e) => handleSubmit(e, true)}
                  disabled={(externalLoading || loading) || !title.trim() || !content.trim()}
                  className="px-4 py-2 text-sm text-gray-400 hover:text-white disabled:opacity-40 transition-colors rounded-lg border border-gray-700 hover:border-gray-600"
                >
                  {t('post.saveDraft')}
                </button>
              )}
              <button
                type="submit"
                disabled={(externalLoading || loading) || !title.trim() || !content.trim()}
                className="px-5 py-2 bg-gray-700 hover:bg-gray-600 disabled:opacity-40 disabled:hover:bg-gray-700 text-white rounded-lg font-medium text-sm transition-colors"
              >
                {(externalLoading || loading) ? (post ? t('post.updating') : t('post.posting')) : (post ? t('post.update') : t('post.post'))}
              </button>
            </div>
          </div>
        </>
      )}
      <ConfirmDialog {...dialogProps} />
    </form>
  );
}
