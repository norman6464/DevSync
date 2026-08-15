import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Plus, Check, X } from 'lucide-react';
import toast from 'react-hot-toast';
import { useBookmarkCollections } from '../../hooks/useBookmarkCollections';
import { addPostToCollection, removePostFromCollection } from '../../api/bookmarkCollections';

interface BookmarkCollectionPickerProps {
  postId: number;
  onClose: () => void;
}

export default function BookmarkCollectionPicker({ postId, onClose }: BookmarkCollectionPickerProps) {
  const { t } = useTranslation();
  const { collections, loading, create, refetch } = useBookmarkCollections();
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newName, setNewName] = useState('');
  const [addingTo, setAddingTo] = useState<number | null>(null);

  const handleAddToCollection = async (collectionId: number) => {
    setAddingTo(collectionId);
    try {
      await addPostToCollection(collectionId, postId);
      toast.success(t('bookmarkCollections.addedToCollection'));
      await refetch();
    } catch {
      toast.error(t('bookmarkCollections.alreadyInCollection'));
    } finally {
      setAddingTo(null);
    }
  };

  const handleRemoveFromCollection = async (collectionId: number) => {
    try {
      await removePostFromCollection(collectionId, postId);
      toast.success(t('bookmarkCollections.removedFromCollection'));
      await refetch();
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  const handleCreate = async () => {
    if (!newName.trim()) return;
    try {
      const col = await create(newName.trim());
      await addPostToCollection(col.id, postId);
      toast.success(t('bookmarkCollections.createdAndAdded'));
      setNewName('');
      setShowCreateForm(false);
    } catch {
      toast.error(t('errors.somethingWrong'));
    }
  };

  return (
    <div className="absolute bottom-10 right-0 z-20 w-64 bg-gray-800 border border-gray-700 rounded-lg shadow-xl p-3">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm font-medium text-gray-200">
          {t('bookmarkCollections.saveToCollection')}
        </span>
        <button onClick={onClose} className="text-gray-500 hover:text-gray-300">
          <X className="w-4 h-4" />
        </button>
      </div>

      {loading ? (
        <div className="text-xs text-gray-500 py-2">{t('common.loading')}</div>
      ) : (
        <div className="max-h-40 overflow-y-auto space-y-1">
          {collections.map((col) => (
            <div
              key={col.id}
              className="flex items-center justify-between px-2 py-1.5 rounded hover:bg-gray-700/50 group"
            >
              <button
                onClick={() => handleAddToCollection(col.id)}
                disabled={addingTo === col.id}
                className="flex items-center gap-2 text-sm text-gray-300 flex-1 text-left"
              >
                <span
                  className="w-2.5 h-2.5 rounded-full flex-shrink-0"
                  style={{ backgroundColor: col.color || '#3b82f6' }}
                />
                <span className="truncate">{col.name}</span>
              </button>
              <button
                onClick={() => handleRemoveFromCollection(col.id)}
                className="text-gray-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-opacity"
                aria-label={t('bookmarkCollections.removeFromCollection')}
              >
                <X className="w-3.5 h-3.5" />
              </button>
            </div>
          ))}
          {collections.length === 0 && (
            <p className="text-xs text-gray-500 py-1">{t('bookmarkCollections.noCollections')}</p>
          )}
        </div>
      )}

      {showCreateForm ? (
        <div className="mt-2 flex gap-1.5">
          <input
            type="text"
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
            placeholder={t('bookmarkCollections.collectionName')}
            className="flex-1 bg-gray-900 border border-gray-700 rounded px-2 py-1 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-blue-500"
            autoFocus
          />
          <button
            onClick={handleCreate}
            className="p-1 text-green-400 hover:text-green-300"
          >
            <Check className="w-4 h-4" />
          </button>
        </div>
      ) : (
        <button
          onClick={() => setShowCreateForm(true)}
          className="mt-2 flex items-center gap-1.5 text-sm text-blue-400 hover:text-blue-300 w-full"
        >
          <Plus className="w-4 h-4" />
          {t('bookmarkCollections.createNew')}
        </button>
      )}
    </div>
  );
}
