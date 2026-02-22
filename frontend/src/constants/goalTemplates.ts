import { Code, Layers, Wrench, FolderGit2, FileText, type LucideIcon } from 'lucide-react';
import type { GoalCategory } from '../api/goals';

export interface GoalTemplate {
  id: string;
  title: string;
  description: string;
  category: GoalCategory;
  estimatedDays: number;
}

// カテゴリ情報
export const GOAL_CATEGORIES: { value: GoalCategory; label: string; Icon: LucideIcon }[] = [
  { value: 'language', label: '言語', Icon: Code },
  { value: 'framework', label: 'フレームワーク', Icon: Layers },
  { value: 'skill', label: 'スキル', Icon: Wrench },
  { value: 'project', label: 'プロジェクト', Icon: FolderGit2 },
  { value: 'other', label: 'その他', Icon: FileText },
];

// ゴールテンプレート
export const GOAL_TEMPLATES: GoalTemplate[] = [
  // 言語
  {
    id: 'typescript-basics',
    title: 'TypeScript基礎マスター',
    description: 'TypeScriptの型システム、インターフェース、ジェネリクスなどの基礎を学習',
    category: 'language',
    estimatedDays: 30,
  },
  {
    id: 'go-backend',
    title: 'Go言語バックエンド開発',
    description: 'Goの基礎からGin、GORMを使ったREST API開発まで習得',
    category: 'language',
    estimatedDays: 45,
  },
  {
    id: 'python-data-science',
    title: 'Pythonデータサイエンス',
    description: 'Pandas、NumPy、Matplotlibを使ったデータ分析の基礎を学習',
    category: 'language',
    estimatedDays: 60,
  },
  {
    id: 'rust-systems',
    title: 'Rustシステムプログラミング',
    description: 'Rustの所有権、借用、ライフタイムなどの基礎概念を理解',
    category: 'language',
    estimatedDays: 90,
  },

  // フレームワーク
  {
    id: 'react-advanced',
    title: 'React実践',
    description: 'React Hooks、状態管理、パフォーマンス最適化などの実践的なスキルを習得',
    category: 'framework',
    estimatedDays: 45,
  },
  {
    id: 'nextjs-fullstack',
    title: 'Next.jsフルスタック開発',
    description: 'Next.js 14のApp Router、Server Components、APIルートを使った開発',
    category: 'framework',
    estimatedDays: 60,
  },
  {
    id: 'vue-composition',
    title: 'Vue 3 Composition API',
    description: 'Vue 3のComposition API、Pinia、Vue Routerを使った開発',
    category: 'framework',
    estimatedDays: 45,
  },
  {
    id: 'spring-boot',
    title: 'Spring Boot開発',
    description: 'Spring BootでのREST API開発、JPA、Securityの実装',
    category: 'framework',
    estimatedDays: 60,
  },

  // スキル
  {
    id: 'docker-kubernetes',
    title: 'Docker & Kubernetes',
    description: 'コンテナ技術の基礎からKubernetesを使った本番運用まで',
    category: 'skill',
    estimatedDays: 60,
  },
  {
    id: 'aws-certified',
    title: 'AWS認定資格取得',
    description: 'AWS Certified Solutions Architect - Associate取得を目指す',
    category: 'skill',
    estimatedDays: 90,
  },
  {
    id: 'algorithm-leetcode',
    title: 'アルゴリズム・データ構造',
    description: 'LeetCodeを使った競技プログラミング・アルゴリズム学習',
    category: 'skill',
    estimatedDays: 120,
  },
  {
    id: 'ci-cd-github-actions',
    title: 'CI/CD構築（GitHub Actions）',
    description: 'GitHub Actionsを使った自動テスト、デプロイパイプラインの構築',
    category: 'skill',
    estimatedDays: 30,
  },
  {
    id: 'graphql-apollo',
    title: 'GraphQL & Apollo',
    description: 'GraphQLスキーマ設計、Apollo Client/Serverの実装',
    category: 'skill',
    estimatedDays: 45,
  },

  // プロジェクト
  {
    id: 'portfolio-website',
    title: 'ポートフォリオサイト作成',
    description: '自分のスキルやプロジェクトを紹介するポートフォリオサイトを構築',
    category: 'project',
    estimatedDays: 30,
  },
  {
    id: 'todo-fullstack',
    title: 'フルスタックTodoアプリ',
    description: '認証、CRUD、リアルタイム同期を持つTodoアプリを開発',
    category: 'project',
    estimatedDays: 45,
  },
  {
    id: 'blog-cms',
    title: 'ブログCMS開発',
    description: 'マークダウン対応、タグ機能、検索機能を持つブログシステム',
    category: 'project',
    estimatedDays: 60,
  },
  {
    id: 'realtime-chat',
    title: 'リアルタイムチャットアプリ',
    description: 'WebSocketを使ったリアルタイムメッセージング機能を実装',
    category: 'project',
    estimatedDays: 45,
  },

  // その他
  {
    id: 'tech-book-reading',
    title: '技術書読破',
    description: '「Clean Code」「リファクタリング」などの名著を読んで実践',
    category: 'other',
    estimatedDays: 60,
  },
  {
    id: 'oss-contribution',
    title: 'OSSコントリビューション',
    description: 'オープンソースプロジェクトへの貢献を通じた実践的学習',
    category: 'other',
    estimatedDays: 90,
  },
];

// カテゴリ情報を取得
export const getCategoryInfo = (category: GoalCategory) =>
  GOAL_CATEGORIES.find((c) => c.value === category) || GOAL_CATEGORIES[4];

// カテゴリでテンプレートをフィルター
export const getTemplatesByCategory = (category: GoalCategory | 'all') => {
  if (category === 'all') return GOAL_TEMPLATES;
  return GOAL_TEMPLATES.filter((t) => t.category === category);
};
