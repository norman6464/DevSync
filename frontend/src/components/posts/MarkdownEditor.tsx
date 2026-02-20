import { useState, useRef, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeSanitize from 'rehype-sanitize';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { uploadImage } from '../../api/posts';
import { sanitizeUrl } from '../../utils/url';
import MarkdownToolbar from './MarkdownToolbar';

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
  placeholder,
  minHeight = '200px',
  onImagesChange,
}: MarkdownEditorProps) {
  const { t } = useTranslation();
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

  return (
    <div className="border border-gray-700 rounded-lg overflow-hidden bg-gray-800">
      {/* Tabs */}
      <div className="flex border-b border-gray-700" role="tablist">
        <button
          type="button"
          role="tab"
          aria-selected={activeTab === 'write'}
          aria-controls="editor-write-panel"
          onClick={() => setActiveTab('write')}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === 'write'
              ? 'bg-gray-700 text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-white hover:bg-gray-750'
          }`}
        >
          {t('editor.write')}
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={activeTab === 'preview'}
          aria-controls="editor-preview-panel"
          onClick={() => setActiveTab('preview')}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            activeTab === 'preview'
              ? 'bg-gray-700 text-white border-b-2 border-blue-500'
              : 'text-gray-400 hover:text-white hover:bg-gray-750'
          }`}
        >
          {t('editor.preview')}
        </button>

        {/* Toolbar */}
        {activeTab === 'write' && (
          <MarkdownToolbar onAction={handleToolbarAction} />
        )}
      </div>

      {/* Content */}
      <div style={{ minHeight }}>
        {activeTab === 'write' ? (
          <div className="relative" id="editor-write-panel" role="tabpanel">
            <textarea
              ref={textareaRef}
              value={value}
              onChange={(e) => onChange(e.target.value)}
              onPaste={handlePaste}
              onDrop={handleDrop}
              onDragOver={handleDragOver}
              placeholder={placeholder || t('editor.placeholder')}
              maxLength={10000}
              className="w-full p-4 bg-transparent text-white resize-none focus:outline-none font-mono text-sm"
              style={{ minHeight }}
            />
            {uploading && (
              <div className="absolute inset-0 bg-gray-900/80 flex items-center justify-center" aria-live="polite">
                <div className="flex items-center gap-2 text-blue-400">
                  <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24" aria-hidden="true">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    />
                  </svg>
                  <span className="text-sm">{t('editor.uploading')}</span>
                </div>
              </div>
            )}
          </div>
        ) : (
          <div
            className="p-4 prose prose-invert prose-sm max-w-none"
            id="editor-preview-panel"
            role="tabpanel"
            style={{ minHeight }}
          >
            {value ? (
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                rehypePlugins={[rehypeSanitize]}
                components={{
                  code({ className, children, ...props }) {
                    const match = /language-(\w+)/.exec(className || '');
                    const inline = !match;
                    return !inline ? (
                      <SyntaxHighlighter
                        style={vscDarkPlus}
                        language={match[1]}
                        PreTag="div"
                        customStyle={{ borderRadius: '0.5rem', fontSize: '0.875rem' }}
                      >
                        {String(children).replace(/\n$/, '')}
                      </SyntaxHighlighter>
                    ) : (
                      <code className={className} {...props}>
                        {children}
                      </code>
                    );
                  },
                }}
              >{value}</ReactMarkdown>
            ) : (
              <p className="text-gray-500 italic">{t('editor.nothingToPreview')}</p>
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
          <div className="text-xs text-gray-500 mb-2">{t('editor.attachedImages')}</div>
          <div className="flex flex-wrap gap-2">
            {uploadedImages.map((url, i) => (
              <div key={i} className="relative group">
                <img
                  src={sanitizeUrl(url) ?? ''}
                  alt={t('editor.uploadedImage', { number: i + 1 })}
                  className="w-16 h-16 object-cover rounded border border-gray-600"
                />
                <button
                  type="button"
                  aria-label={t('editor.removeImage', { number: i + 1 })}
                  onClick={() => {
                    const newImages = uploadedImages.filter((_, j) => j !== i);
                    setUploadedImages(newImages);
                    onImagesChange?.(newImages);
                  }}
                  className="absolute -top-1 -right-1 w-5 h-5 bg-red-500 rounded-full text-white text-xs opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                >
                  <span aria-hidden="true">x</span>
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Hint */}
      <div className="border-t border-gray-700 px-4 py-2 text-xs text-gray-500 flex items-center gap-4">
        <span>{t('editor.markdownSupported')}</span>
        <span>{t('editor.pasteOrDragImages')}</span>
        <span className="ml-auto">{value.length.toLocaleString()} / 10,000</span>
      </div>
    </div>
  );
}
