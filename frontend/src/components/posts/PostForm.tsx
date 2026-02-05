import { useState } from 'react';
import toast from 'react-hot-toast';

interface PostFormProps {
  onSubmit: (title: string, content: string) => Promise<void>;
}

export default function PostForm({ onSubmit }: PostFormProps) {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !content.trim()) return;
    setLoading(true);
    try {
      await onSubmit(title, content);
      setTitle('');
      setContent('');
      toast.success('Post created');
    } catch {
      toast.error('Failed to create post');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="bg-gray-900 rounded-lg p-6 space-y-3">
      <input
        type="text"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Title"
        className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:border-blue-500"
      />
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="What did you learn today?"
        rows={3}
        className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:border-blue-500"
      />
      <button
        type="submit"
        disabled={loading}
        className="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded-md font-medium text-sm transition-colors"
      >
        {loading ? 'Posting...' : 'Post'}
      </button>
    </form>
  );
}
