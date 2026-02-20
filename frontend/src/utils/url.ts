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
