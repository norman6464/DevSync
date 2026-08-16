import { Monitor, Rocket, Target, FolderOpen, FileText, type LucideIcon } from 'lucide-react';
import { type GoalCategory } from '../../api/goals';

export const CATEGORIES: { value: GoalCategory; label: string; icon: string; Icon: LucideIcon }[] = [
  { value: 'language', label: 'goals.categoryLanguage', icon: '💻', Icon: Monitor },
  { value: 'framework', label: 'goals.categoryFramework', icon: '🚀', Icon: Rocket },
  { value: 'skill', label: 'goals.categorySkill', icon: '🎯', Icon: Target },
  { value: 'project', label: 'goals.categoryProject', icon: '📁', Icon: FolderOpen },
  { value: 'other', label: 'goals.categoryOther', icon: '📝', Icon: FileText },
];
