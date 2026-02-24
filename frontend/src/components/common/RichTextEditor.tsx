import { useRef, useCallback } from 'react';
import { Bold, Italic, Underline, List } from 'lucide-react';

interface RichTextEditorProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
  error?: string;
  placeholder?: string;
  className?: string;
}

export default function RichTextEditor({
  value,
  onChange,
  label,
  error,
  placeholder,
  className = '',
}: RichTextEditorProps) {
  const editorRef = useRef<HTMLDivElement>(null);

  const execCommand = useCallback((command: string) => {
    document.execCommand(command, false);
    if (editorRef.current) {
      onChange(editorRef.current.innerHTML);
    }
  }, [onChange]);

  const handleInput = useCallback(() => {
    if (editorRef.current) {
      onChange(editorRef.current.innerHTML);
    }
  }, [onChange]);

  const showPlaceholder = !value && placeholder;

  return (
    <div className={`${className}`.trim()}>
      {label && <label className="block text-sm text-gray-400 mb-1">{label}</label>}
      <div className="border border-gray-700 rounded-lg overflow-hidden">
        <div data-testid="toolbar" className="flex gap-1 px-2 py-1.5 border-b border-gray-700 bg-gray-800/50">
          <button
            type="button"
            aria-label="太字"
            onClick={() => execCommand('bold')}
            className="p-1.5 text-gray-400 hover:text-white hover:bg-gray-700 rounded"
          >
            <Bold className="w-4 h-4" />
          </button>
          <button
            type="button"
            aria-label="斜体"
            onClick={() => execCommand('italic')}
            className="p-1.5 text-gray-400 hover:text-white hover:bg-gray-700 rounded"
          >
            <Italic className="w-4 h-4" />
          </button>
          <button
            type="button"
            aria-label="下線"
            onClick={() => execCommand('underline')}
            className="p-1.5 text-gray-400 hover:text-white hover:bg-gray-700 rounded"
          >
            <Underline className="w-4 h-4" />
          </button>
          <button
            type="button"
            aria-label="リスト"
            onClick={() => execCommand('insertUnorderedList')}
            className="p-1.5 text-gray-400 hover:text-white hover:bg-gray-700 rounded"
          >
            <List className="w-4 h-4" />
          </button>
        </div>
        <div className="relative">
          {showPlaceholder && (
            <span className="absolute top-3 left-4 text-sm text-gray-500 pointer-events-none">
              {placeholder}
            </span>
          )}
          <div
            ref={editorRef}
            role="textbox"
            contentEditable
            onInput={handleInput}
            dangerouslySetInnerHTML={{ __html: value }}
            className="min-h-[120px] px-4 py-3 bg-gray-800 text-sm text-gray-200 focus:outline-none"
          />
        </div>
      </div>
      {error && <p className="mt-1 text-xs text-red-400">{error}</p>}
    </div>
  );
}
