import { useState, useRef, useEffect, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { getUsers } from '../../api/users';
import type { User } from '../../types/user';
import Avatar from './Avatar';

interface MentionInputProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  onSubmit?: () => void;
  disabled?: boolean;
}

export default function MentionInput({
  value,
  onChange,
  placeholder,
  className = '',
  onSubmit,
  disabled,
}: MentionInputProps) {
  const { t } = useTranslation();
  const [suggestions, setSuggestions] = useState<User[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [mentionQuery, setMentionQuery] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);
  const suggestionsRef = useRef<HTMLDivElement>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout>>();

  const findMentionQuery = useCallback((text: string, cursorPos: number) => {
    const beforeCursor = text.slice(0, cursorPos);
    const match = beforeCursor.match(/(?:^|\s)@([a-zA-Z0-9_-]*)$/);
    return match ? match[1] : null;
  }, []);

  const fetchSuggestions = useCallback(async (query: string) => {
    if (query.length < 1) {
      setSuggestions([]);
      setShowSuggestions(false);
      return;
    }
    try {
      const { data } = await getUsers(query);
      const users = (data || []).slice(0, 5);
      setSuggestions(users);
      setShowSuggestions(users.length > 0);
      setSelectedIndex(0);
    } catch {
      setSuggestions([]);
      setShowSuggestions(false);
    }
  }, []);

  const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    onChange(newValue);

    const cursorPos = e.target.selectionStart || 0;
    const query = findMentionQuery(newValue, cursorPos);

    if (query !== null) {
      setMentionQuery(query);
      if (debounceRef.current) clearTimeout(debounceRef.current);
      debounceRef.current = setTimeout(() => fetchSuggestions(query), 200);
    } else {
      setShowSuggestions(false);
      setMentionQuery('');
    }
  }, [onChange, findMentionQuery, fetchSuggestions]);

  const insertMention = useCallback((username: string) => {
    const input = inputRef.current;
    if (!input) return;

    const cursorPos = input.selectionStart || 0;
    const beforeCursor = value.slice(0, cursorPos);
    const afterCursor = value.slice(cursorPos);

    const mentionStart = beforeCursor.lastIndexOf('@');
    const newValue = beforeCursor.slice(0, mentionStart) + `@${username} ` + afterCursor;

    onChange(newValue);
    setShowSuggestions(false);
    setMentionQuery('');

    setTimeout(() => {
      const newCursorPos = mentionStart + username.length + 2;
      input.setSelectionRange(newCursorPos, newCursorPos);
      input.focus();
    }, 0);
  }, [value, onChange]);

  const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
    if (!showSuggestions) {
      if (e.key === 'Enter' && onSubmit) {
        e.preventDefault();
        onSubmit();
      }
      return;
    }

    if (e.key === 'ArrowDown') {
      e.preventDefault();
      setSelectedIndex((i) => Math.min(i + 1, suggestions.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      setSelectedIndex((i) => Math.max(i - 1, 0));
    } else if (e.key === 'Enter' || e.key === 'Tab') {
      e.preventDefault();
      if (suggestions[selectedIndex]) {
        insertMention(suggestions[selectedIndex].username);
      }
    } else if (e.key === 'Escape') {
      setShowSuggestions(false);
    }
  }, [showSuggestions, suggestions, selectedIndex, insertMention, onSubmit]);

  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, []);

  return (
    <div className="relative flex-1">
      <input
        ref={inputRef}
        type="text"
        value={value}
        onChange={handleChange}
        onKeyDown={handleKeyDown}
        placeholder={placeholder}
        disabled={disabled}
        className={className}
      />
      {showSuggestions && (
        <div
          ref={suggestionsRef}
          className="absolute bottom-full left-0 mb-1 w-64 bg-gray-800 border border-gray-700 rounded-lg shadow-xl z-50 overflow-hidden"
        >
          <div className="px-3 py-1.5 text-xs text-gray-500 border-b border-gray-700">
            {t('mention.suggestions')}
          </div>
          {suggestions.map((user, index) => (
            <button
              key={user.id}
              type="button"
              className={`w-full flex items-center gap-2 px-3 py-2 text-left text-sm transition-colors ${
                index === selectedIndex ? 'bg-blue-600/20 text-white' : 'text-gray-300 hover:bg-gray-700'
              }`}
              onMouseDown={(e) => {
                e.preventDefault();
                insertMention(user.username);
              }}
              onMouseEnter={() => setSelectedIndex(index)}
            >
              <Avatar name={user.name} avatarUrl={user.avatar_url} size="xs" />
              <div className="min-w-0">
                <div className="font-medium truncate">{user.name}</div>
                <div className="text-xs text-gray-500">@{user.username}</div>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
