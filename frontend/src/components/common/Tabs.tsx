import { useState, ReactNode } from 'react';

interface Tab {
  id: string;
  label: string;
  content: ReactNode;
}

interface TabsProps {
  tabs: Tab[];
  defaultActiveId?: string;
  onChange?: (id: string) => void;
}

export default function Tabs({ tabs, defaultActiveId, onChange }: TabsProps) {
  const [activeId, setActiveId] = useState(defaultActiveId || tabs[0]?.id);

  const handleTabClick = (id: string) => {
    setActiveId(id);
    onChange?.(id);
  };

  const activeTab = tabs.find((tab) => tab.id === activeId);

  return (
    <div className="w-full">
      {/* タブヘッダー */}
      <div className="flex border-b border-gray-800">
        {tabs.map((tab) => {
          const isActive = tab.id === activeId;

          return (
            <button
              key={tab.id}
              onClick={() => handleTabClick(tab.id)}
              className={`px-4 py-3 text-sm font-medium transition-colors border-b-2 ${
                isActive
                  ? 'text-blue-400 border-blue-400'
                  : 'text-gray-400 hover:text-white border-transparent'
              }`}
            >
              {tab.label}
            </button>
          );
        })}
      </div>

      {/* タブコンテンツ */}
      <div className="py-4">{activeTab?.content}</div>
    </div>
  );
}
