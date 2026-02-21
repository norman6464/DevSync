// ノートのソートオプション定義
export const NOTES_SORT_OPTIONS = [
  { value: 'latest' as const, label: 'notes.sortLatest' },
  { value: 'oldest' as const, label: 'notes.sortOldest' },
  { value: 'updated' as const, label: 'notes.sortUpdated' },
  { value: 'favorites_first' as const, label: 'notes.sortFavorites' },
] as const;

export type NoteSortValue = typeof NOTES_SORT_OPTIONS[number]['value'];
