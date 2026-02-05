import type { GitHubLanguageStat } from '../../types/github';

interface LanguageChartProps {
  languages: GitHubLanguageStat[];
}

const COLORS = [
  '#3b82f6', '#ef4444', '#22c55e', '#f59e0b', '#8b5cf6',
  '#ec4899', '#06b6d4', '#f97316', '#14b8a6', '#6366f1',
];

export default function LanguageChart({ languages }: LanguageChartProps) {
  const totalBytes = languages.reduce((sum, l) => sum + l.bytes, 0);
  if (totalBytes === 0) return null;

  const sorted = [...languages].sort((a, b) => b.bytes - a.bytes).slice(0, 10);

  return (
    <div className="space-y-3">
      <div className="flex h-4 rounded-full overflow-hidden">
        {sorted.map((lang, i) => (
          <div
            key={lang.language}
            style={{
              width: `${(lang.bytes / totalBytes) * 100}%`,
              backgroundColor: COLORS[i % COLORS.length],
            }}
            title={`${lang.language}: ${((lang.bytes / totalBytes) * 100).toFixed(1)}%`}
          />
        ))}
      </div>
      <div className="flex flex-wrap gap-3">
        {sorted.map((lang, i) => (
          <div key={lang.language} className="flex items-center gap-1.5 text-sm">
            <div
              className="w-3 h-3 rounded-full"
              style={{ backgroundColor: COLORS[i % COLORS.length] }}
            />
            <span className="text-gray-300">{lang.language}</span>
            <span className="text-gray-500">
              {((lang.bytes / totalBytes) * 100).toFixed(1)}%
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
