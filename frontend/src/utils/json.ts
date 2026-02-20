export function parseJsonArray<T = string>(json: string | undefined | null): T[] {
  if (!json) return [];
  try {
    return JSON.parse(json);
  } catch {
    return [];
  }
}
