import { useState, KeyboardEvent } from 'react';
import { X } from 'lucide-react';

interface TagInputProps {
  value: string[];
  onChange: (value: string[]) => void;
  label?: string;
  placeholder?: string;
  maxTags?: number;
  disabled?: boolean;
  className?: string;
}

export default function TagInput({
  value,
  onChange,
  label,
  placeholder = 'タグを入力...',
  maxTags,
  disabled = false,
  className = '',
}: TagInputProps) {
  const [input, setInput] = useState('');

  const isMaxReached = maxTags != null && value.length >= maxTags;

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      const tag = input.trim();
      if (!tag || value.includes(tag)) return;
      onChange([...value, tag]);
      setInput('');
    }
  };

  const removeTag = (index: number) => {
    onChange(value.filter((_, i) => i !== index));
  };

  return (
    <div className={`${className}`.trim()}>
      {label && <label className="block text-sm text-gray-400 mb-1">{label}</label>}
      <div className="flex flex-wrap gap-2 p-2 bg-gray-800 border border-gray-700 rounded-lg min-h-[42px]">
        {value.map((tag, i) => (
          <span
            key={tag}
            className="flex items-center gap-1 px-2 py-1 bg-gray-700 rounded text-sm text-gray-200"
          >
            {tag}
            {!disabled && (
              <button
                type="button"
                aria-label="削除"
                onClick={() => removeTag(i)}
                className="text-gray-400 hover:text-white"
              >
                <X className="w-3 h-3" />
              </button>
            )}
          </span>
        ))}
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={value.length === 0 ? placeholder : ''}
          disabled={disabled || isMaxReached}
          className="flex-1 min-w-[80px] bg-transparent text-sm text-gray-200 placeholder-gray-500 focus:outline-none disabled:opacity-50"
        />
      </div>
    </div>
  );
}
