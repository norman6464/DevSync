import { useState, useRef, useEffect, ReactNode } from 'react';

interface DropdownItem {
  id: string;
  label?: string;
  disabled?: boolean;
}

interface DropdownDivider {
  id: 'divider';
}

interface DropdownProps {
  trigger: ReactNode;
  items: (DropdownItem | DropdownDivider)[];
  onSelect: (id: string) => void;
  className?: string;
}

export default function Dropdown({
  trigger,
  items,
  onSelect,
  className = '',
}: DropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (ref.current && !ref.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleSelect = (item: DropdownItem | DropdownDivider) => {
    if (item.id === 'divider') return;
    if ('disabled' in item && item.disabled) return;
    onSelect(item.id);
    setIsOpen(false);
  };

  return (
    <div ref={ref} className={`relative inline-block ${className}`.trim()}>
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="inline-flex items-center"
      >
        {trigger}
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-48 bg-gray-800 border border-gray-700 rounded-lg shadow-lg z-50 py-1">
          {items.map((item, index) => {
            if (item.id === 'divider') {
              return <div key={`divider-${index}`} className="border-t border-gray-700 my-1" />;
            }

            const menuItem = item as DropdownItem;
            return (
              <button
                key={menuItem.id}
                type="button"
                onClick={() => handleSelect(menuItem)}
                disabled={menuItem.disabled}
                className={`w-full text-left px-4 py-2 text-sm transition-colors ${
                  menuItem.disabled
                    ? 'text-gray-600 cursor-not-allowed'
                    : 'text-gray-200 hover:bg-gray-700'
                }`}
              >
                {menuItem.label}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}
