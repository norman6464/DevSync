import { useState } from 'react';
import toast from 'react-hot-toast';
import MarkdownEditor from './MarkdownEditor';

interface PostFormProps {
  onSubmit: (title: string, content: string, imageUrls?: string) => Promise<void>;
}

export default function PostForm({ onSubmit }: PostFormProps) {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [imageUrls, setImageUrls] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [expanded, setExpanded] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;
    setLoading(true);
    try {
      const imageUrlsJson = imageUrls.length > 0 ? JSON.stringify(imageUrls) : undefined;
      await onSubmit(title, content, imageUrlsJson);
      setTitle('');
      setContent('');
      setImageUrls([]);
      setExpanded(false);
      toast.success('Post created');
    } catch {
      toast.error('Failed to create post');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-gray-900 border border-gray-800 rounded-xl p-4">
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        onFocus={() => setExpanded(true)}
        placeholder="What did you learn today?"
        className="w-full px-3 py-2.5 bg-gray-800/50 border border-gray-700 rounded-lg text-white text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-shadow"
      />
      {expanded && (
        <>
          <div className="mt-3">
            <MarkdownEditor
              value={content}
              onChange={setContent}
              placeholder="Share your thoughts... (Markdown supported)"
              minHeight="150px"
              onImagesChange={setImageUrls}
            />
          </div>
          <div className="flex justify-between items-center mt-3">
            <div className="text-xs text-gray-500">
              Supports **bold**, *italic*, `code`, and more
            </div>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={() => {
                  setExpanded(false);
                  setTitle('');
                  setContent('');
                  setImageUrls([]);
                }}
                className="px-4 py-2 text-sm text-gray-400 hover:text-white transition-colors rounded-lg"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={loading || !title.trim() || !content.trim()}
                className="px-5 py-2 bg-green-600 hover:bg-green-500 disabled:opacity-40 disabled:hover:bg-green-600 text-white rounded-lg font-medium text-sm transition-colors"
              >
                {loading ? 'Posting...' : 'Post'}
              </button>
            </div>
          </div>
        </>
      )}
    </form>
  );
}
