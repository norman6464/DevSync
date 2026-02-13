import { useTranslation } from 'react-i18next';
import { Users, FileText, Users2 } from 'lucide-react';

export type SearchTab = 'users' | 'posts' | 'circles';

interface SearchTabsProps {
  activeTab: SearchTab;
  onTabChange: (tab: SearchTab) => void;
  counts: {
    users: number;
    posts: number;
    circles: number;
  };
}

export default function SearchTabs({ activeTab, onTabChange, counts }: SearchTabsProps) {
  const { t } = useTranslation();

  const tabs: Array<{ id: SearchTab; label: string; icon: typeof Users; count: number }> = [
    { id: 'users', label: t('search.tabs.users'), icon: Users, count: counts.users },
    { id: 'posts', label: t('search.tabs.posts'), icon: FileText, count: counts.posts },
    { id: 'circles', label: t('search.tabs.circles'), icon: Users2, count: counts.circles },
  ];

  return (
    <div className="border-b border-gray-800">
      <div className="flex gap-1">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => onTabChange(tab.id)}
              className={`
                flex items-center gap-2 px-4 py-3 text-sm font-medium transition-colors
                border-b-2 ${
                  isActive
                    ? 'border-blue-500 text-blue-400'
                    : 'border-transparent text-gray-400 hover:text-gray-300'
                }
              `}
            >
              <Icon className="w-4 h-4" />
              <span>{tab.label}</span>
              <span
                className={`
                  px-2 py-0.5 rounded-full text-xs
                  ${isActive ? 'bg-blue-500/20 text-blue-400' : 'bg-gray-800 text-gray-500'}
                `}
              >
                {tab.count}
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
