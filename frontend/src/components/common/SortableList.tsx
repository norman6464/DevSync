import { ChevronUp, ChevronDown, GripVertical, X } from 'lucide-react';

interface SortableItem {
  id: string;
  label: string;
}

interface SortableListProps {
  items: SortableItem[];
  onReorder: (items: SortableItem[]) => void;
  onRemove?: (id: string) => void;
  className?: string;
}

export default function SortableList({
  items,
  onReorder,
  onRemove,
  className = '',
}: SortableListProps) {
  const moveItem = (index: number, direction: 'up' | 'down') => {
    const newItems = [...items];
    const targetIndex = direction === 'up' ? index - 1 : index + 1;
    [newItems[index], newItems[targetIndex]] = [newItems[targetIndex], newItems[index]];
    onReorder(newItems);
  };

  return (
    <div className={`space-y-1 ${className}`.trim()}>
      {items.map((item, index) => (
        <div
          key={item.id}
          className="flex items-center gap-2 px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg"
        >
          <GripVertical className="w-4 h-4 text-gray-500 flex-shrink-0" />
          <span className="flex-1 text-sm text-gray-200">{item.label}</span>
          <div className="flex items-center gap-1">
            <button
              type="button"
              data-testid="move-up"
              disabled={index === 0}
              onClick={() => moveItem(index, 'up')}
              className="p-1 text-gray-400 hover:text-white disabled:text-gray-700 disabled:cursor-not-allowed"
            >
              <ChevronUp className="w-4 h-4" />
            </button>
            <button
              type="button"
              data-testid="move-down"
              disabled={index === items.length - 1}
              onClick={() => moveItem(index, 'down')}
              className="p-1 text-gray-400 hover:text-white disabled:text-gray-700 disabled:cursor-not-allowed"
            >
              <ChevronDown className="w-4 h-4" />
            </button>
            {onRemove && (
              <button
                type="button"
                data-testid="remove-item"
                onClick={() => onRemove(item.id)}
                className="p-1 text-gray-400 hover:text-red-400"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
