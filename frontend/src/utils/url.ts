/**
 * URLがHTTPまたはHTTPSプロトコルかを検証する。
 * javascript:, data:, vbscript: 等の危険なプロトコルを拒否する。
 */
export function isHttpUrl(url: string): boolean {
  if (!url) return false;
  try {
    const parsed = new URL(url);
    return parsed.protocol === 'http:' || parsed.protocol === 'https:';
  } catch {
    return false;
  }
}

/**
 * URLをサニタイズして返す。
 * http/httpsでないURLの場合はundefinedを返す。
 */
export function sanitizeUrl(url: string | undefined): string | undefined {
  if (!url) return undefined;
  return isHttpUrl(url) ? url : undefined;
}

/**
 * URLフィールド配列を検証し、不正なURLがあればエラーメッセージを返す。
 * 全て有効ならnullを返す。
 */
export function findInvalidUrlField(
  fields: { value: string; label: string }[],
): string | null {
  for (const field of fields) {
    if (field.value && !isHttpUrl(field.value)) {
      return field.label;
    }
  }
  return null;
}
