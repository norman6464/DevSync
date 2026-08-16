import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import {
  BookOpen,
  FileText,
  Target,
  PenSquare,
  Link as LinkIcon,
  HelpCircle,
  Users,
  FolderGit2,
  Zap
} from 'lucide-react';
import { panelClass } from '../../constants/styles';

interface QuickAction {
  id: string;
  label: string;
  description: string;
  icon: React.ElementType;
  href: string;
  color: string;
}

const QUICK_ACTIONS: QuickAction[] = [
  {
    id: 'learning-log',
    label: '学習ログを追加',
    description: '今日の学習記録を残そう',
    icon: BookOpen,
    href: '/learning-logs',
    color: 'text-blue-400',
  },
  {
    id: 'note',
    label: 'ノートを作成',
    description: '新しい学習ノートを書こう',
    icon: FileText,
    href: '/notes',
    color: 'text-green-400',
  },
  {
    id: 'goal',
    label: '目標を追加',
    description: '新しい学習目標を設定しよう',
    icon: Target,
    href: '/goals',
    color: 'text-orange-400',
  },
  {
    id: 'post',
    label: '投稿を作成',
    description: '学習成果をシェアしよう',
    icon: PenSquare,
    href: '/',
    color: 'text-purple-400',
  },
  {
    id: 'resource',
    label: 'リソースを追加',
    description: '役立つリソースを保存しよう',
    icon: LinkIcon,
    href: '/resources',
    color: 'text-yellow-400',
  },
  {
    id: 'qa',
    label: '問題を質問',
    description: 'Q&Aで質問してみよう',
    icon: HelpCircle,
    href: '/qa',
    color: 'text-red-400',
  },
  {
    id: 'study-circle',
    label: '勉強会を探す',
    description: '学習仲間を見つけよう',
    icon: Users,
    href: '/study-circles',
    color: 'text-indigo-400',
  },
  {
    id: 'project',
    label: 'プロジェクトを開始',
    description: '新しいプロジェクトを始めよう',
    icon: FolderGit2,
    href: '/projects',
    color: 'text-pink-400',
  },
];

export default function QuickActionsWidget() {
  // t は未使用だが、i18n の購読（言語切替時の再レンダー・サスペンド挙動）を維持するため呼び出しは残す
  useTranslation();

  return (
    <div className={panelClass}>
      <div className="flex items-center gap-2 mb-4">
        <Zap className="w-4 h-4 text-yellow-400" />
        <h3 className="text-sm font-medium text-white">クイックアクション</h3>
      </div>

      <div className="grid grid-cols-2 gap-2">
        {QUICK_ACTIONS.map((action) => {
          const Icon = action.icon;

          return (
            <Link
              key={action.id}
              to={action.href}
              className="bg-gray-800/50 border border-gray-700 rounded-lg p-3 hover:border-blue-400/50 hover:bg-gray-800 transition-all group"
            >
              <div className="flex items-start gap-2 mb-1">
                <div className={`${action.color} flex-shrink-0`}>
                  <Icon className="w-5 h-5" />
                </div>
                <div className="flex-1 min-w-0">
                  <div className="text-sm font-medium text-white group-hover:text-blue-400 transition-colors">
                    {action.label}
                  </div>
                </div>
              </div>
              <p className="text-xs text-gray-500 line-clamp-1 ml-7">
                {action.description}
              </p>
            </Link>
          );
        })}
      </div>
    </div>
  );
}
