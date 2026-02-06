import { useState, useRef, useCallback } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { uploadImage } from '../../api/posts';

interface MarkdownEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  minHeight?: string;
  onImagesChange?: (urls: string[]) => void;
}

export default function MarkdownEditor({
  value,
  onChange,
  placeholder = 'Write your content here... (Markdown supported)',
  minHeight = '200px',
  onImagesChange,
}: MarkdownEditorProps) {
  const [activeTab, setActiveTab] = useState<'write' | 'preview'>('write');
  const [uploading, setUploading] = useState(false);
  const [uploadedImages, setUploadedImages] = useState<string[]>([]);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const insertText = useCallback(
    (before: string, after: string = '') => {
      const textarea = textareaRef.current;
      if (!textarea) return;

      const start = textarea.selectionStart;
      const end = textarea.selectionEnd;
      const selectedText = value.substring(start, end);
      const newText = value.substring(0, start) + before + selectedText + after + value.substring(end);

      onChange(newText);

      // Restore cursor position
      setTimeout(() => {
        textarea.focus();
        const newPos = start + before.length + selectedText.length + after.length;
        textarea.setSelectionRange(newPos, newPos);
      }, 0);
    },
    [value, onChange]
  );

  const handleToolbarAction = (action: string) => {
    switch (action) {
      case 'bold':
        insertText('**', '**');
        break;
      case 'italic':
        insertText('*', '*');
        break;
      case 'strikethrough':
        insertText('~~', '~~');
        break;
      case 'heading':
        insertText('## ');
        break;
      case 'link':
        insertText('[', '](url)');
        break;
      case 'code':
        insertText('`', '`');
        break;
      case 'codeblock':
        insertText('```\n', '\n```');
        break;
      case 'quote':
        insertText('> ');
        break;
      case 'list':
        insertText('- ');
        break;
      case 'orderedlist':
        insertText('1. ');
        break;
      case 'task':
        insertText('- [ ] ');
        break;
      case 'image':
        fileInputRef.current?.click();
        break;
    }
  };

  const handleImageUpload = async (files: FileList | null) => {
    if (!files || files.length === 0) return;

    setUploading(true);
    try {
      const uploadPromises = Array.from(files).map((file) => uploadImage(file));
      const results = await Promise.all(uploadPromises);
      const urls = results.map((r) => r.url);

      // Insert image markdown
      const imageMarkdown = urls.map((url) => `![image](${url})`).join('\n');
      insertText(imageMarkdown + '\n');

      // Update uploaded images list
      const newImages = [...uploadedImages, ...urls];
      setUploadedImages(newImages);
      onImagesChange?.(newImages);
    } catch (error) {
      console.error('Failed to upload image:', error);
    } finally {
      setUploading(false);
    }
  };

  const handlePaste = async (e: React.ClipboardEvent) => {
    const items = e.clipboardData?.items;
    if (!items) return;

    for (const item of items) {
      if (item.type.startsWith('image/')) {
        e.preventDefault();
        const file = item.getAsFile();
        if (file) {
          const dataTransfer = new DataTransfer();
          dataTransfer.items.add(file);
          await handleImageUpload(dataTransfer.files);
        }
        break;
      }
    }
  };

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault();
    const files = e.dataTransfer?.files;
    if (files && files.length > 0) {
      const imageFiles = Array.from(files).filter((f) => f.type.startsWith('image/'));
      if (imageFiles.length > 0) {
        const dataTransfer = new DataTransfer();
        imageFiles.forEach((f) => dataTransfer.items.add(f));
        await handleImageUpload(dataTransfer.files);
      }
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  const toolbarButtons = [
    { action: 'heading', icon: 'H', title: 'Heading' },
    { action: 'bold', icon: 'B', title: 'Bold', className: 'font-bold' },
    { action: 'italic', icon: 'I', title: 'Italic', className: 'italic' },
    { action: 'strikethrough', icon: 'S', title: 'Strikethrough', className: 'line-through' },
    { action: 'divider' },
    { action: 'link', icon: '🔗', title: 'Link' },
    { action: 'image', icon: '🖼️', title: 'Image' },
    { action: 'divider' },
    { action: 'code', icon: '<>', title: 'Inline Code' },
    { action: 'codeblock', icon: '{}', title: 'Code Block' },
    { action: 'divider' },
    { action: 'quote', icon: '"', title: 'Quote' },
    { action: 'list', icon: '•', title: 'Bullet List' },
    { action: 'orderedlist', icon: '1.', title: 'Numbered List' },
    { action: 'task', icon: '☐', title: 'Task List' },
  ];

  return (
    <div className="border border-gray-700 rounded-lg overflow-hidden bg-gray-800">
      {/* Tabs */}
      <div className="flex border-b border-gray-700">
        <button
          type="button"
          onClick={() => setActiveTab('write')}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === 'write'
              ? 'bg-gray-700 text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-white hover:bg-gray-750'
          }`}
        >
          Write
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('preview')}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === 'preview'
              ? 'bg-gray-700 text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-white hover:bg-gray-750'
          }`}
        >
          Preview
        </button>

        {/* Toolbar */}
        {activeTab === 'write' && (
          <div className="flex items-center gap-0.5 ml-auto px-2">
            {toolbarButtons.map((btn, i) =>
              btn.action === 'divider' ? (
                <div key={i} className="w-px h-4 bg-gray-600 mx-1" />
              ) : (
                <button
                  key={btn.action}
                  type="button"
                  onClick={() => handleToolbarAction(btn.action)}
                  className={`p-1.5 text-xs text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-colors ${btn.className || ''}`}
                  title={btn.title}
                >
                  {btn.icon}
                </button>
              )
            )}
          </div>
        )}
      </div>

      {/* Content */}
      <div style={{ minHeight }}>
        {activeTab === 'write' ? (
          <div className="relative">
            <textarea
              ref={textareaRef}
              value={value}
              onChange={(e) => onChange(e.target.value)}
              onPaste={handlePaste}
              onDrop={handleDrop}
              onDragOver={handleDragOver}
              placeholder={placeholder}
              className="w-full p-4 bg-transparent text-white resize-none focus:outline-none font-mono text-sm"
              style={{ minHeight }}
            />
            {uploading && (
              <div className="absolute inset-0 bg-gray-900/80 flex items-center justify-center">
                <div className="flex items-center gap-2 text-blue-400">
                  <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    />
                  </svg>
                  <span className="text-sm">Uploading...</span>
                </div>
              </div>
            )}
          </div>
        ) : (
          <div
            className="p-4 prose prose-invert prose-sm max-w-none"
            style={{ minHeight }}
          >
            {value ? (
              <ReactMarkdown remarkPlugins={[remarkGfm]}>{value}</ReactMarkdown>
            ) : (
              <p className="text-gray-500 italic">Nothing to preview</p>
            )}
          </div>
        )}
      </div>

      {/* Hidden file input */}
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        multiple
        className="hidden"
        onChange={(e) => handleImageUpload(e.target.files)}
      />

      {/* Uploaded images preview */}
      {uploadedImages.length > 0 && (
        <div className="border-t border-gray-700 p-3">
          <div className="text-xs text-gray-500 mb-2">Attached images:</div>
          <div className="flex flex-wrap gap-2">
            {uploadedImages.map((url, i) => (
              <div key={i} className="relative group">
                <img
                  src={url}
                  alt={`Uploaded ${i + 1}`}
                  className="w-16 h-16 object-cover rounded border border-gray-600"
                />
                <button
                  type="button"
                  onClick={() => {
                    const newImages = uploadedImages.filter((_, j) => j !== i);
                    setUploadedImages(newImages);
                    onImagesChange?.(newImages);
                  }}
                  className="absolute -top-1 -right-1 w-5 h-5 bg-red-500 rounded-full text-white text-xs opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                >
                  x
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Hint */}
      <div className="border-t border-gray-700 px-4 py-2 text-xs text-gray-500 flex items-center gap-4">
        <span>Markdown supported</span>
        <span>Paste or drag images to upload</span>
      </div>
    </div>
  );
}
