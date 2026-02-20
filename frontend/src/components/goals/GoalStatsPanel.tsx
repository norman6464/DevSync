import { useTranslation } from 'react-i18next';

interface GoalStatsPanelProps {
  total: number;
  active: number;
  completed: number;
  paused: number;
  overdue: number;
}

export default function GoalStatsPanel({ total, active, completed, paused, overdue }: GoalStatsPanelProps) {
  const { t } = useTranslation();

  const stats = [
    { value: total, label: t('goals.totalGoals'), color: '' },
    { value: active, label: t('goals.activeGoals'), color: 'text-blue-400' },
    { value: completed, label: t('goals.completedGoals'), color: 'text-green-400' },
    { value: paused, label: t('goals.pausedGoals'), color: 'text-yellow-400' },
    { value: overdue, label: t('goals.overdueGoals'), color: 'text-red-400' },
  ];

  return (
    <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
      {stats.map((stat) => (
        <div key={stat.label} className="bg-gray-900 border border-gray-800 rounded-md p-4">
          <p className={`text-2xl font-bold ${stat.color}`}>{stat.value}</p>
          <p className="text-sm text-gray-400">{stat.label}</p>
        </div>
      ))}
    </div>
  );
}
