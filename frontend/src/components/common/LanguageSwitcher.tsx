import { useState } from 'react';
import { Globe, ChevronDown } from 'lucide-react';

interface Language {
  code: string;
  label: string;
  flag?: string;
}

interface LanguageSwitcherProps {
  languages: Language[];
  current: string;
  onChange: (code: string) => void;
  disabled?: boolean;
  className?: string;
}

export default function LanguageSwitcher({
  languages,
  current,
  onChange,
  disabled = false,
  className = '',
}: LanguageSwitcherProps) {
  const [open, setOpen] = useState(false);
  const currentLang = languages.find((l) => l.code === current);

  const handleSelect = (code: string) => {
    onChange(code);
    setOpen(false);
  };

  return (
    <div className={`relative inline-block ${className}`.trim()}>
      <button
        type="button"
        onClick={() => setOpen(!open)}
        disabled={disabled}
        className="flex items-center gap-2 px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-sm text-gray-200 hover:border-gray-600 disabled:opacity-50"
      >
        <Globe className="w-4 h-4 text-gray-400" />
        {currentLang?.flag && <span>{currentLang.flag}</span>}
        <span>{currentLang?.label}</span>
        <ChevronDown className="w-3 h-3 text-gray-400" />
      </button>
      {open && (
        <div className="absolute top-full left-0 mt-1 w-full bg-gray-800 border border-gray-700 rounded-lg shadow-xl overflow-hidden z-50">
          {languages.map((lang) => (
            <div
              key={lang.code}
              role="option"
              aria-selected={lang.code === current}
              onClick={() => handleSelect(lang.code)}
              className={`flex items-center gap-2 px-3 py-2 text-sm cursor-pointer hover:bg-gray-700/50 ${
                lang.code === current ? 'bg-gray-700 text-white' : 'text-gray-300'
              }`}
            >
              {lang.flag && <span>{lang.flag}</span>}
              <span>{lang.label}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
