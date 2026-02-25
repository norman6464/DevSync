// フォームバリデーション定数・ユーティリティ

/** 各フォームのフィールドごとの最大文字数定義 */
export const MAX_LENGTH = {
  // 投稿
  postTitle: 300,
  postContent: 5000,
  // ノート
  noteTitle: 200,
  noteContent: 10000,
  noteTags: 500,
  // 学習ログ
  logTitle: 200,
  logContent: 5000,
  // 学習目標
  goalTitle: 200,
  goalDescription: 2000,
  // Q&A
  questionTitle: 500,
  questionBody: 5000,
  questionTags: 500,
  answerBody: 5000,
  // 書籍レビュー
  bookTitle: 300,
  bookAuthor: 200,
  bookIsbn: 20,
  bookReview: 2000,
  // リソース
  resourceTitle: 300,
  resourceDescription: 1000,
  // プロジェクト
  projectTitle: 200,
  projectDescription: 2000,
  projectRole: 100,
  // 共通
  url: 2048,
  message: 1000,
} as const;

/** 値がトリム後に空でないことを検証する */
export function isRequired(value: string): boolean {
  return value.trim().length > 0;
}

/** 複数フィールドがすべて必須入力済みかを検証する */
export function areAllRequired(...values: string[]): boolean {
  return values.every(isRequired);
}
