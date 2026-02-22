import { useState, ReactNode } from 'react';
import { ChevronDown } from 'lucide-react';

interface AccordionItem {
  id: string;
  title: string;
  content: ReactNode;
}

interface AccordionProps {
  items: AccordionItem[];
  single?: boolean;
  defaultOpenIds?: string[];
  onChange?: (openIds: string[]) => void;
  className?: string;
}

export default function Accordion({
  items,
  single = false,
  defaultOpenIds = [],
  onChange,
  className = '',
}: AccordionProps) {
  const [openIds, setOpenIds] = useState<string[]>(defaultOpenIds);

  const toggle = (id: string) => {
    const newOpenIds = openIds.includes(id)
      ? openIds.filter((openId) => openId !== id)
      : single
        ? [id]
        : [...openIds, id];

    setOpenIds(newOpenIds);
    onChange?.(newOpenIds);
  };

  return (
    <div className={`divide-y divide-gray-800 border border-gray-800 rounded-lg ${className}`.trim()}>
      {items.map((item) => {
        const isOpen = openIds.includes(item.id);

        return (
          <div key={item.id}>
            <button
              onClick={() => toggle(item.id)}
              className="flex items-center justify-between w-full px-4 py-3 text-left text-gray-200 hover:bg-gray-800/50 transition-colors"
              aria-expanded={isOpen}
            >
              <span>{item.title}</span>
              <ChevronDown
                className={`w-4 h-4 text-gray-400 transition-transform ${isOpen ? 'rotate-180' : ''}`}
              />
            </button>
            {isOpen && (
              <div className="px-4 py-3 text-gray-300">
                {item.content}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
