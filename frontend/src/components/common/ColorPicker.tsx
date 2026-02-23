const DEFAULT_COLORS = [
  '#ef4444', '#f97316', '#eab308', '#22c55e', '#14b8a6',
  '#3b82f6', '#6366f1', '#a855f7', '#ec4899', '#6b7280',
];

const sizeMap = {
  sm: 'w-6 h-6',
  md: 'w-8 h-8',
  lg: 'w-10 h-10',
};

interface ColorPickerProps {
  value: string;
  onChange: (color: string) => void;
  colors?: string[];
  showInput?: boolean;
  label?: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export default function ColorPicker({
  value,
  onChange,
  colors = DEFAULT_COLORS,
  showInput = false,
  label,
  size = 'md',
  className = '',
}: ColorPickerProps) {
  return (
    <div className={`${className}`.trim()}>
      {label && <label className="block text-sm text-gray-400 mb-2">{label}</label>}
      <div className="flex flex-wrap gap-2">
        {colors.map((color) => (
          <button
            key={color}
            type="button"
            data-testid="color-swatch"
            onClick={() => onChange(color)}
            className={`${sizeMap[size]} rounded-full cursor-pointer transition-all ${
              value === color ? 'ring-2 ring-white ring-offset-2 ring-offset-gray-900' : ''
            }`}
            style={{ backgroundColor: color }}
          />
        ))}
      </div>
      {showInput && (
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          className="mt-3 w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded text-sm text-gray-200 font-mono focus:outline-none focus:border-blue-500"
        />
      )}
    </div>
  );
}
