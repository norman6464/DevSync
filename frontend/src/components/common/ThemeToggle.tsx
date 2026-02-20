import { useTranslation } from 'react-i18next';
import { Sun, Moon, Monitor } from 'lucide-react';
import { useThemeStore } from '../../store/themeStore';

export default function ThemeToggle() {
  const { t } = useTranslation();
  const { theme, setTheme } = useThemeStore();

  const themes = [
    {
      value: 'light' as const,
      icon: <Sun className="w-4 h-4" />,
      label: t('settings.light'),
    },
    {
      value: 'dark' as const,
      icon: <Moon className="w-4 h-4" />,
      label: t('settings.dark'),
    },
    {
      value: 'system' as const,
      icon: <Monitor className="w-4 h-4" />,
      label: t('settings.system'),
    },
  ];

  return (
    <div className="flex items-center gap-1 p-1 bg-gray-800 rounded-lg">
      {themes.map((item) => (
        <button
          key={item.value}
          onClick={() => setTheme(item.value)}
          className={`p-2 rounded-md transition-all ${
            theme === item.value
              ? 'bg-gray-700 text-white shadow-sm'
              : 'text-gray-400 hover:text-white'
          }`}
          title={item.label}
          aria-label={t('settings.setTheme', { theme: item.label })}
        >
          {item.icon}
        </button>
      ))}
    </div>
  );
}
