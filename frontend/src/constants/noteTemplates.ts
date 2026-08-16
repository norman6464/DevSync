import { BookOpen, Lightbulb, FolderOpen, Calendar, type LucideIcon } from 'lucide-react';

export type NoteCategory = '学習ノート' | 'プロジェクト' | '復習' | 'その他';

export interface NoteTemplate {
  id: string;
  title: string;
  content: string;
  category: NoteCategory;
  tags: string;
}

// カテゴリ情報
export const NOTE_CATEGORIES: { value: NoteCategory; label: string; Icon: LucideIcon }[] = [
  { value: '学習ノート', label: '学習ノート', Icon: BookOpen },
  { value: 'プロジェクト', label: 'プロジェクト', Icon: FolderOpen },
  { value: '復習', label: '復習', Icon: Calendar },
  { value: 'その他', label: 'その他', Icon: Lightbulb },
];

// ノートテンプレート
export const NOTE_TEMPLATES: NoteTemplate[] = [
  // 学習ノート
  {
    id: 'coding-memo',
    title: 'コーディング学習メモ',
    content: `# コーディング学習メモ

## 学習日時
${new Date().toLocaleDateString('ja-JP')}

## 学習内容
-

## 学んだこと
-

## つまずいたポイント
-

## コードスニペット
\`\`\`
// ここにコードを記述
\`\`\`

## 次回やること
- `,
    category: '学習ノート',
    tags: '#学習,#コーディング',
  },
  {
    id: 'book-reading',
    title: '読書ノート',
    content: `# 読書ノート

## 書籍情報
- **タイトル**:
- **著者**:
- **出版社**:
- **読了日**: ${new Date().toLocaleDateString('ja-JP')}

## 要約
-

## 重要なポイント
-

## 印象に残った引用
>

## 実践したいこと
-

## 評価
⭐️⭐️⭐️⭐️⭐️ (5段階)`,
    category: '学習ノート',
    tags: '#読書,#書籍',
  },
  {
    id: 'online-course',
    title: 'オンラインコース学習',
    content: `# オンラインコース学習

## コース情報
- **コース名**:
- **プラットフォーム**: (Udemy / Coursera / その他)
- **進捗**: [ ] 10% [ ] 30% [ ] 50% [ ] 80% [x] 100%

## 今日の学習内容
- セクション:
- レッスン:

## 学んだこと
1.
2.
3.

## 演習・課題
-

## メモ
- `,
    category: '学習ノート',
    tags: '#オンラインコース,#学習',
  },
  {
    id: 'tech-article',
    title: '技術記事まとめ',
    content: `# 技術記事まとめ

## 記事情報
- **タイトル**:
- **URL**:
- **投稿日**:

## 要点
-

## 技術スタック
-

## 実装のポイント
-

## 参考になったコード
\`\`\`
// ここにコードを記述
\`\`\`

## 所感
- `,
    category: '学習ノート',
    tags: '#技術記事,#まとめ',
  },

  // プロジェクト
  {
    id: 'project-plan',
    title: 'プロジェクト計画',
    content: `# プロジェクト計画

## プロジェクト名


## 概要


## 目的・ゴール
-

## 技術スタック
- **フロントエンド**:
- **バックエンド**:
- **データベース**:
- **インフラ**:

## スケジュール
- [ ] 要件定義 (〜)
- [ ] 設計 (〜)
- [ ] 実装 (〜)
- [ ] テスト (〜)
- [ ] リリース (〜)

## タスク
- [ ]
- [ ]
- [ ]

## リスク・課題
- `,
    category: 'プロジェクト',
    tags: '#プロジェクト,#計画',
  },
  {
    id: 'tech-investigation',
    title: '技術調査',
    content: `# 技術調査

## 調査テーマ


## 背景・目的


## 調査内容
### 選択肢1:
- メリット:
- デメリット:
- 使用例:

### 選択肢2:
- メリット:
- デメリット:
- 使用例:

## 比較表
| 項目 | 選択肢1 | 選択肢2 |
|------|--------|--------|
| パフォーマンス |  |  |
| 学習コスト |  |  |
| コミュニティ |  |  |

## 結論


## 参考リンク
- `,
    category: 'プロジェクト',
    tags: '#技術調査,#比較検討',
  },
  {
    id: 'troubleshooting',
    title: 'トラブルシューティング',
    content: `# トラブルシューティング

## 発生日時
${new Date().toLocaleString('ja-JP')}

## 問題の概要


## エラーメッセージ
\`\`\`
// エラー内容
\`\`\`

## 環境
- OS:
- 言語/フレームワークバージョン:
- その他:

## 試したこと
1.
2.
3.

## 解決方法


## 根本原因


## 再発防止策
- `,
    category: 'プロジェクト',
    tags: '#トラブルシューティング,#エラー',
  },
  {
    id: 'release-note',
    title: 'リリースノート',
    content: `# リリースノート

## バージョン
v0.0.0

## リリース日
${new Date().toLocaleDateString('ja-JP')}

## 新機能
-

## 改善
-

## バグ修正
-

## 破壊的変更
-

## 非推奨
-

## 移行ガイド


## 既知の問題
- `,
    category: 'プロジェクト',
    tags: '#リリース,#バージョン',
  },

  // 復習
  {
    id: 'daily-review',
    title: 'デイリーレビュー',
    content: `# デイリーレビュー

## 日付
${new Date().toLocaleDateString('ja-JP')}

## 今日やったこと
-

## 学んだこと
-

## よかったこと
-

## 改善点
-

## 明日やること
- `,
    category: '復習',
    tags: '#振り返り,#デイリー',
  },
  {
    id: 'weekly-review',
    title: 'ウィークリーレビュー',
    content: `# ウィークリーレビュー

## 期間
〜

## 今週の成果
-

## 達成したゴール
- [x]
- [ ]

## 学んだこと
-

## 課題・改善点
-

## 来週の目標
- `,
    category: '復習',
    tags: '#振り返り,#ウィークリー',
  },
  {
    id: 'monthly-review',
    title: '月次振り返り',
    content: `# 月次振り返り

## 対象月
${new Date().getFullYear()}年${new Date().getMonth() + 1}月

## 今月の成果
### 技術面
-

### プロジェクト
-

### 学習
-

## 統計
- 学習時間: 〜時間
- コミット数: 〜回
- 完了タスク: 〜個

## よかったこと
-

## 反省点
-

## 来月の目標
1.
2.
3. `,
    category: '復習',
    tags: '#振り返り,#月次',
  },

  // その他
  {
    id: 'idea-memo',
    title: 'アイデアメモ',
    content: `# アイデアメモ

## アイデア名


## 概要


## 背景・きっかけ


## 詳細
-

## 期待される効果
-

## 実現に必要なこと
-

## 参考
- `,
    category: 'その他',
    tags: '#アイデア,#メモ',
  },
];

// カテゴリでテンプレートをフィルター
export const getTemplatesByCategory = (category: NoteCategory | 'all') => {
  if (category === 'all') return NOTE_TEMPLATES;
  return NOTE_TEMPLATES.filter((t) => t.category === category);
};
