import { describe, it, expect } from 'vitest';
import { isHttpUrl, sanitizeUrl } from '../url';

describe('isHttpUrl', () => {
  it('https URLを許可する', () => {
    expect(isHttpUrl('https://example.com')).toBe(true);
    expect(isHttpUrl('https://example.com/path?q=1')).toBe(true);
  });

  it('http URLを許可する', () => {
    expect(isHttpUrl('http://example.com')).toBe(true);
  });

  it('javascript: URLを拒否する', () => {
    expect(isHttpUrl('javascript:alert(1)')).toBe(false);
  });

  it('data: URLを拒否する', () => {
    expect(isHttpUrl('data:text/html,<script>alert(1)</script>')).toBe(false);
  });

  it('vbscript: URLを拒否する', () => {
    expect(isHttpUrl('vbscript:msgbox')).toBe(false);
  });

  it('空文字列を拒否する', () => {
    expect(isHttpUrl('')).toBe(false);
  });

  it('不正なURLを拒否する', () => {
    expect(isHttpUrl('not-a-url')).toBe(false);
  });

  it('ftp: URLを拒否する', () => {
    expect(isHttpUrl('ftp://example.com')).toBe(false);
  });
});

describe('sanitizeUrl', () => {
  it('有効なhttps URLをそのまま返す', () => {
    expect(sanitizeUrl('https://example.com/avatar.png')).toBe('https://example.com/avatar.png');
  });

  it('javascript: URLに対してundefinedを返す', () => {
    expect(sanitizeUrl('javascript:alert(1)')).toBeUndefined();
  });

  it('undefinedに対してundefinedを返す', () => {
    expect(sanitizeUrl(undefined)).toBeUndefined();
  });

  it('空文字列に対してundefinedを返す', () => {
    expect(sanitizeUrl('')).toBeUndefined();
  });
});
