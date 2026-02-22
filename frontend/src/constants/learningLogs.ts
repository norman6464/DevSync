import { Code, BookOpen, GraduationCap, Users, FileText, type LucideIcon } from 'lucide-react';
import type { LogCategory } from '../types/learningLog';

// 学習ログのカテゴリ定義
export const LOG_CATEGORIES: { value: LogCategory; label: string; Icon: LucideIcon }[] = [
  { value: 'coding', label: 'learningLogs.categoryCoding', Icon: Code },
  { value: 'reading', label: 'learningLogs.categoryReading', Icon: BookOpen },
  { value: 'course', label: 'learningLogs.categoryCourse', Icon: GraduationCap },
  { value: 'meetup', label: 'learningLogs.categoryMeetup', Icon: Users },
  { value: 'other', label: 'learningLogs.categoryOther', Icon: FileText },
];

// カテゴリ情報を取得する
export const getCategoryInfo = (cat: LogCategory) =>
  LOG_CATEGORIES.find((c) => c.value === cat) || LOG_CATEGORIES[4];

// カテゴリの色クラスを取得する
export const getCategoryColor = (cat: LogCategory) => {
  switch (cat) {
    case 'coding':
      return 'text-blue-400 bg-blue-400/10';
    case 'reading':
      return 'text-green-400 bg-green-400/10';
    case 'course':
      return 'text-purple-400 bg-purple-400/10';
    case 'meetup':
      return 'text-orange-400 bg-orange-400/10';
    default:
      return 'text-gray-400 bg-gray-400/10';
  }
};
