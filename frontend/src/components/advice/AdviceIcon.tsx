import {
  Flame, Map, Target, Lightbulb, TrendingUp,
  Award, BookOpen, HelpCircle,
} from 'lucide-react';

const iconMap: Record<string, React.ElementType> = {
  streak_recovery: Flame,
  stalled_roadmap: Map,
  goal_overdue: Target,
  tech_suggestion: Lightbulb,
  goal_suggestion: Target,
  difficulty_up: TrendingUp,
  praise: Award,
  general: BookOpen,
};

const colorMap: Record<string, string> = {
  streak_recovery: 'text-red-400',
  stalled_roadmap: 'text-yellow-400',
  goal_overdue: 'text-orange-400',
  tech_suggestion: 'text-blue-400',
  goal_suggestion: 'text-green-400',
  difficulty_up: 'text-purple-400',
  praise: 'text-yellow-300',
  general: 'text-gray-400',
};

interface AdviceIconProps {
  type: string;
  size?: number;
  className?: string;
}

export default function AdviceIcon({ type, size = 20, className = '' }: AdviceIconProps) {
  const Icon = iconMap[type] || HelpCircle;
  const color = colorMap[type] || 'text-gray-400';

  return <Icon size={size} className={`${color} ${className}`} />;
}
