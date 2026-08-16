/// <reference types="node" />
import { describe, it, expect } from 'vitest';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const SRC_DIR = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const LOCALES_DIR = path.join(SRC_DIR, 'i18n/locales');

type Leaf = string | string[];
type Tree = { [key: string]: Leaf | Tree };

/**
 * ネストした翻訳ツリーを `a.b.c` 形式に平坦化する。
 * 曜日ラベルのように配列で持つ値は、それ自体が 1 つの翻訳なので葉として扱う。
 */
function flatten(tree: Tree, prefix = ''): Record<string, Leaf> {
  const out: Record<string, Leaf> = {};
  for (const [key, value] of Object.entries(tree)) {
    const full = `${prefix}${key}`;
    if (typeof value === 'string' || Array.isArray(value)) out[full] = value;
    else Object.assign(out, flatten(value, `${full}.`));
  }
  return out;
}

function loadLocales(): Record<string, Record<string, Leaf>> {
  const out: Record<string, Record<string, Leaf>> = {};
  for (const file of fs.readdirSync(LOCALES_DIR)) {
    if (!file.endsWith('.json')) continue;
    const raw = fs.readFileSync(path.join(LOCALES_DIR, file), 'utf-8');
    out[path.basename(file, '.json')] = flatten(JSON.parse(raw) as Tree);
  }
  return out;
}

function sourceFiles(dir: string): string[] {
  const out: string[] = [];
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (entry.name === '__tests__' || entry.name === 'test') continue;
      out.push(...sourceFiles(full));
    } else if (/\.tsx?$/.test(entry.name)) {
      out.push(full);
    }
  }
  return out;
}

function record(found: Map<string, string[]>, key: string, file: string): void {
  const rel = path.relative(SRC_DIR, file);
  const sites = found.get(key);
  if (sites) sites.push(rel);
  else found.set(key, [rel]);
}

/** ソース中の `t('some.key')` を集める。 */
function usedKeys(): Map<string, string[]> {
  const pattern = /\bt\(\s*(['"])([a-zA-Z0-9_.-]+)\1/g;
  const found = new Map<string, string[]>();
  for (const file of sourceFiles(SRC_DIR)) {
    const text = fs.readFileSync(file, 'utf-8');
    for (const match of text.matchAll(pattern)) {
      // `t('prefix.' + variable)` は動的キーなので静的には解決できない。
      if (text.slice(match.index + match[0].length).trimStart().startsWith('+')) continue;
      record(found, match[2], file);
    }
  }
  return found;
}

/**
 * `` t(`resources.difficulty.${x}`) `` のような動的キーから、差し込みの手前までの
 * 静的な接頭辞を集める。キーは 1 つのセグメントに解決されるので、
 * 検証対象は「接頭辞の直下にある葉」に限る。
 */
function usedKeyPrefixes(): Map<string, string[]> {
  const pattern = /\bt\(\s*`([^`]*)`/g;
  const found = new Map<string, string[]>();
  for (const file of sourceFiles(SRC_DIR)) {
    const text = fs.readFileSync(file, 'utf-8');
    for (const match of text.matchAll(pattern)) {
      const placeholder = match[1].indexOf('${');
      if (placeholder <= 0) continue;
      record(found, match[1].slice(0, placeholder), file);
    }
  }
  return found;
}

const locales = loadLocales();
const localeNames = Object.keys(locales).sort();

describe('翻訳キー', () => {
  it('locale ファイルが 10 言語ぶん読み込める', () => {
    expect(localeNames).toEqual(['de', 'en', 'es', 'fr', 'ja', 'ko', 'pt', 'ru', 'zh-CN', 'zh-TW']);
  });

  // 未定義キーに対して i18next はキー文字列自体を返すため、
  // 抜けがあると画面に `dashboard.timeline` のような開発者向け文字列が出る。
  it('コードが参照するキーがすべての言語に定義されている', () => {
    const missing: string[] = [];
    for (const [key, sites] of usedKeys()) {
      const absent = localeNames.filter((name) => !(key in locales[name]));
      if (absent.length > 0) {
        missing.push(`${key} — 未定義: ${absent.join(', ')}（使用箇所: ${sites.join(', ')}）`);
      }
    }
    expect(missing).toEqual([]);
  });

  // 動的キーは値が実データ（ステータス・カテゴリ等）で決まるため、どの変種が来ても
  // 引けるように「接頭辞の直下にある葉」が全言語でそろっている必要がある。
  it('動的キーの変種がすべての言語にそろっている', () => {
    const problems: string[] = [];
    for (const [prefix, sites] of usedKeyPrefixes()) {
      const perLocale = new Map(
        localeNames.map((name) => [
          name,
          new Set(
            Object.keys(locales[name]).filter(
              (key) => key.startsWith(prefix) && !key.slice(prefix.length).includes('.'),
            ),
          ),
        ]),
      );
      const union = new Set(localeNames.flatMap((name) => [...perLocale.get(name)!]));
      if (union.size === 0) {
        problems.push(`${prefix}* — どの言語にも該当キーが無い（使用箇所: ${sites.join(', ')}）`);
        continue;
      }
      for (const name of localeNames) {
        const absent = [...union].filter((key) => !perLocale.get(name)!.has(key)).sort();
        if (absent.length > 0) problems.push(`${prefix}* — ${name} に不足: ${absent.join(', ')}`);
      }
    }
    expect(problems).toEqual([]);
  });

  it('コードが参照するキーの値が空文字になっていない', () => {
    const empty: string[] = [];
    for (const key of usedKeys().keys()) {
      for (const name of localeNames) {
        const value = locales[name][key];
        if (typeof value === 'string' && value.trim() === '') empty.push(`${name}: ${key}`);
      }
    }
    expect(empty).toEqual([]);
  });

  // {{count}} のような差し込みが言語によって欠けると、その言語でだけ数値が消える。
  it('差し込み変数がすべての言語で一致する', () => {
    const placeholders = (value: Leaf | undefined) =>
      typeof value === 'string'
        ? [...value.matchAll(/\{\{\s*([a-zA-Z0-9_]+)\s*\}\}/g)].map((m) => m[1]).sort()
        : [];

    const mismatched: string[] = [];
    for (const key of usedKeys().keys()) {
      const expected = placeholders(locales.ja[key]);
      for (const name of localeNames) {
        const actual = placeholders(locales[name][key]);
        if (actual.join(',') !== expected.join(',')) {
          mismatched.push(`${key} — ja: [${expected}] / ${name}: [${actual}]`);
        }
      }
    }
    expect(mismatched).toEqual([]);
  });
});
