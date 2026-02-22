import { useState, KeyboardEvent } from 'react';
import { X, Hash } from 'lucide-react';

interface TagSelectorProps {
  value: string;
  onChange: (value: string) => void;
  suggestions?: string[];
  maxTags?: number;
}

export default function TagSelector({
  value,
  onChange,
  suggestions = [],
  maxTags = 10,
}: TagSelectorProps) {
  const [inputValue, setInputValue] = useState('');

  const selectedTags = value
    ? value.split(',').filter((tag) => tag.trim() !== '')
    : [];

  const normalizeTag = (tag: string): string => {
    const trimmed = tag.trim();
    return trimmed.startsWith('#') ? trimmed : `#${trimmed}`;
  };

  const addTags = (tagsToAdd: string[]) => {
    const normalizedNewTags = tagsToAdd
      .map((tag) => tag.trim())
      .filter((tag) => tag !== '')
      .map(normalizeTag);

    const uniqueNewTags = normalizedNewTags.filter(
      (tag) => !selectedTags.includes(tag)
    );

    if (uniqueNewTags.length === 0) return;

    const totalTags = selectedTags.length + uniqueNewTags.length;
    if (totalTags > maxTags) return;

    const newValue = [...selectedTags, ...uniqueNewTags].join(',');
    onChange(newValue);
  };

  const removeTag = (tagToRemove: string) => {
    const newTags = selectedTags.filter((tag) => tag !== tagToRemove);
    onChange(newTags.join(','));
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      const tags = inputValue.split(',').map((t) => t.trim());
      addTags(tags);
      setInputValue('');
    }
  };

  const handleSuggestionClick = (suggestion: string) => {
    addTags([suggestion]);
  };

  return (
    <div className="space-y-3">
      {/* 選択済みタグ */}
      {selectedTags.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {selectedTags.map((tag) => (
            <button
              key={tag}
              onClick={() => removeTag(tag)}
              className="bg-blue-500/20 text-blue-400 px-3 py-1.5 rounded-lg text-sm font-medium flex items-center gap-1.5 hover:bg-red-500/30 transition-colors border border-blue-400/30"
            >
              <Hash className="w-3.5 h-3.5" />
              {tag.replace('#', '')}
              <X className="w-3.5 h-3.5" />
            </button>
          ))}
        </div>
      )}

      {/* タグ入力フィールド */}
      <div className="relative">
        <input
          type="text"
          value={inputValue}
          onChange={(e) => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="タグを入力してEnterキーで追加（カンマ区切りで複数追加可能）"
          className="w-full px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <Hash className="absolute right-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-500" />
      </div>

      {/* 候補タグ */}
      {suggestions.length > 0 && (
        <div className="space-y-2">
          <div className="text-sm text-gray-400">候補タグ:</div>
          <div className="flex flex-wrap gap-2">
            {suggestions.map((suggestion) => {
              const normalizedSuggestion = normalizeTag(suggestion);
              const isSelected = selectedTags.includes(normalizedSuggestion);

              return (
                <button
                  key={suggestion}
                  onClick={() => handleSuggestionClick(suggestion)}
                  disabled={isSelected}
                  className={`px-3 py-1.5 rounded-lg text-sm font-medium flex items-center gap-1.5 transition-colors ${
                    isSelected
                      ? 'bg-blue-500/20 text-blue-400 border border-blue-400/30 cursor-default'
                      : 'bg-gray-800/50 text-gray-300 hover:text-white hover:bg-gray-700 border border-gray-700'
                  }`}
                >
                  <Hash className="w-3.5 h-3.5" />
                  {suggestion}
                </button>
              );
            })}
          </div>
        </div>
      )}

      {/* タグ数制限表示 */}
      {selectedTags.length > 0 && (
        <div className="text-xs text-gray-500">
          {selectedTags.length} / {maxTags} タグ
        </div>
      )}
    </div>
  );
}
