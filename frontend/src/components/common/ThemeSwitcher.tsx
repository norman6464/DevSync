import { Sun, Moon, Monitor } from 'lucide-react';

type Theme = 'light' | 'dark' | 'system';

interface ThemeSwitcherProps {
  value: Theme;
  onChange: (theme: Theme) => void;
  label?: string;
  disabled?: boolean;
  className?: string;
}

const themes: { value: Theme; label: string; icon: typeof Sun }[] = [
  { value: 'light', label: 'ライト', icon: Sun },
  { value: 'dark', label: 'ダーク', icon: Moon },
  { value: 'system', label: 'システム', icon: Monitor },
];

export default function ThemeSwitcher({
  value,
  onChange,
  label,
  disabled = false,
  className = '',
}: ThemeSwitcherProps) {
  return (
    <div className={`${className}`.trim()}>
      {label && <p className="text-sm text-gray-400 mb-1">{label}</p>}
      <div className="inline-flex bg-gray-800 rounded-lg p-1">
        {themes.map(({ value: themeValue, label: themeLabel, icon: Icon }) => (
          <button
            key={themeValue}
            type="button"
            aria-label={themeLabel}
            onClick={() => onChange(themeValue)}
            disabled={disabled}
            className={`p-2 rounded-md transition-colors disabled:opacity-50 ${
              value === themeValue
                ? 'bg-blue-600 text-white'
                : 'text-gray-400 hover:text-gray-200'
            }`}
          >
            <Icon className="w-4 h-4" />
          </button>
        ))}
      </div>
    </div>
  );
}
