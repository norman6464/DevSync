import { describe, it, expect } from 'vitest';
import { isHttpUrl, sanitizeUrl, findInvalidUrlField } from '../url';

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

describe('findInvalidUrlField', () => {
  it('全て有効なURLの場合nullを返す', () => {
    expect(findInvalidUrlField([
      { value: 'https://example.com', label: 'Demo URL' },
      { value: 'https://github.com/repo', label: 'GitHub URL' },
    ])).toBeNull();
  });

  it('不正なURLがある場合そのラベルを返す', () => {
    expect(findInvalidUrlField([
      { value: 'https://example.com', label: 'Demo URL' },
      { value: 'not-a-url', label: 'GitHub URL' },
    ])).toBe('GitHub URL');
  });

  it('空の値はスキップする', () => {
    expect(findInvalidUrlField([
      { value: '', label: 'Demo URL' },
      { value: 'https://example.com', label: 'GitHub URL' },
    ])).toBeNull();
  });

  it('空配列の場合nullを返す', () => {
    expect(findInvalidUrlField([])).toBeNull();
  });

  it('最初の不正URLのラベルを返す', () => {
    expect(findInvalidUrlField([
      { value: 'invalid1', label: 'Field A' },
      { value: 'invalid2', label: 'Field B' },
    ])).toBe('Field A');
  });
});
