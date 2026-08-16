import { useState, useEffect, useRef, useMemo } from 'react';
import { Search } from 'lucide-react';

interface Command {
  id: string;
  label: string;
  category: string;
  icon?: string;
}

interface CommandPaletteProps {
  open: boolean;
  commands: Command[];
  onSelect: (command: Command) => void;
  onClose: () => void;
  className?: string;
}

export default function CommandPalette({
  open,
  commands,
  onSelect,
  onClose,
  className = '',
}: CommandPaletteProps) {
  const [query, setQuery] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  // 開いた瞬間に検索欄をリセットする。effect で setState すると余計な再レンダーが
  // 入るため、公式の「render 中に前回値と比較して調整する」パターンで行う。
  const [prevOpen, setPrevOpen] = useState(open);
  if (prevOpen !== open) {
    setPrevOpen(open);
    if (open) setQuery('');
  }

  useEffect(() => {
    if (open) {
      setTimeout(() => inputRef.current?.focus(), 0);
    }
  }, [open]);

  const filtered = useMemo(
    () => commands.filter((cmd) => cmd.label.toLowerCase().includes(query.toLowerCase())),
    [commands, query]
  );

  const grouped = useMemo(() => {
    const map = new Map<string, Command[]>();
    filtered.forEach((cmd) => {
      const group = map.get(cmd.category) || [];
      group.push(cmd);
      map.set(cmd.category, group);
    });
    return map;
  }, [filtered]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-[20vh]">
      <div data-testid="palette-overlay" className="absolute inset-0 bg-black/50" onClick={onClose} />
      <div className={`relative w-full max-w-lg bg-gray-900 border border-gray-700 rounded-xl shadow-2xl overflow-hidden ${className}`.trim()}>
        <div className="flex items-center gap-2 px-4 py-3 border-b border-gray-800">
          <Search className="w-4 h-4 text-gray-500" />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="コマンドを検索..."
            className="flex-1 bg-transparent text-sm text-gray-200 placeholder-gray-500 focus:outline-none"
            autoFocus
          />
        </div>
        <div className="max-h-[300px] overflow-y-auto p-2">
          {filtered.length === 0 ? (
            <p className="text-center text-sm text-gray-500 py-4">コマンドが見つかりません</p>
          ) : (
            Array.from(grouped.entries()).map(([category, cmds]) => (
              <div key={category}>
                <p className="text-xs text-gray-500 px-2 py-1 font-medium">{category}</p>
                {cmds.map((cmd) => (
                  <button
                    key={cmd.id}
                    type="button"
                    onClick={() => onSelect(cmd)}
                    className="w-full flex items-center gap-2 px-3 py-2 text-sm text-gray-300 hover:bg-gray-800 rounded-lg transition-colors"
                  >
                    {cmd.icon && <span>{cmd.icon}</span>}
                    <span>{cmd.label}</span>
                  </button>
                ))}
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
