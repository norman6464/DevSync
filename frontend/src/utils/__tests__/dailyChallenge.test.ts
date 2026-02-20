import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import {
  getDailyChallenge,
  isChallengeCompleted,
  markChallengeCompleted,
} from '../dailyChallenge';

describe('getDailyChallenge', () => {
  it('同一日付で常に同じチャレンジを返す', () => {
    const date = new Date(2026, 0, 15);
    const result1 = getDailyChallenge(date);
    const result2 = getDailyChallenge(date);
    expect(result1).toBe(result2);
  });

  it('異なる日付で異なるチャレンジを返す可能性がある', () => {
    const results = new Set<string>();
    for (let d = 1; d <= 30; d++) {
      results.add(getDailyChallenge(new Date(2026, 0, d)));
    }
    expect(results.size).toBeGreaterThan(1);
  });

  it('有効なチャレンジキーを返す', () => {
    const validKeys = [
      'writeTIL',
      'studyMinutes',
      'readArticle',
      'writeCode',
      'reviewCode',
      'learnNewConcept',
      'shareKnowledge',
      'solveAlgorithm',
      'readDocs',
      'pairProgram',
    ];
    const result = getDailyChallenge(new Date(2026, 5, 10));
    expect(validKeys).toContain(result);
  });
});

function createMockStorage() {
  const store: Record<string, string> = {};
  return {
    store,
    mock: {
      getItem: vi.fn((key: string) => store[key] ?? null),
      setItem: vi.fn((key: string, value: string) => { store[key] = value; }),
      removeItem: vi.fn((key: string) => { delete store[key]; }),
      clear: vi.fn(() => { Object.keys(store).forEach(k => delete store[k]); }),
      length: 0,
      key: vi.fn(() => null),
    } as unknown as Storage,
  };
}

describe('isChallengeCompleted', () => {
  let originalLocalStorage: Storage;

  beforeEach(() => {
    originalLocalStorage = globalThis.localStorage;
    const { mock } = createMockStorage();
    vi.stubGlobal('localStorage', mock);
  });

  afterEach(() => {
    vi.stubGlobal('localStorage', originalLocalStorage);
  });

  it('未完了の場合falseを返す', () => {
    expect(isChallengeCompleted()).toBe(false);
  });

  it('完了済みの場合trueを返す', () => {
    markChallengeCompleted();
    expect(isChallengeCompleted()).toBe(true);
  });

  it('localStorageが使えない場合falseを返す', () => {
    vi.stubGlobal('localStorage', {
      getItem: () => { throw new Error('not available'); },
      setItem: () => { throw new Error('not available'); },
      removeItem: () => {},
      clear: () => {},
      length: 0,
      key: () => null,
    });
    expect(isChallengeCompleted()).toBe(false);
  });
});

describe('markChallengeCompleted', () => {
  let originalLocalStorage: Storage;
  let store: Record<string, string>;

  beforeEach(() => {
    originalLocalStorage = globalThis.localStorage;
    const created = createMockStorage();
    store = created.store;
    vi.stubGlobal('localStorage', created.mock);
  });

  afterEach(() => {
    vi.stubGlobal('localStorage', originalLocalStorage);
  });

  it('localStorageにdoneを書き込む', () => {
    markChallengeCompleted();
    const now = new Date();
    const key = `daily-challenge-${now.getFullYear()}-${now.getMonth() + 1}-${now.getDate()}`;
    expect(store[key]).toBe('done');
  });

  it('古いエントリを削除する', () => {
    const now = new Date();
    const oldDate = new Date(now);
    oldDate.setDate(oldDate.getDate() - 10);
    const oldKey = `daily-challenge-${oldDate.getFullYear()}-${oldDate.getMonth() + 1}-${oldDate.getDate()}`;
    store[oldKey] = 'done';

    markChallengeCompleted();

    expect(store[oldKey]).toBeUndefined();
  });

  it('localStorageが使えない場合エラーを投げない', () => {
    vi.stubGlobal('localStorage', {
      getItem: () => null,
      setItem: () => { throw new Error('not available'); },
      removeItem: () => {},
      clear: () => {},
      length: 0,
      key: () => null,
    });
    expect(() => markChallengeCompleted()).not.toThrow();
  });
});
