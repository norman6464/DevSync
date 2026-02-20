import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { formatDistanceToNow } from '../timeFormat';

describe('formatDistanceToNow', () => {
  const NOW = new Date('2026-02-19T12:00:00Z');

  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(NOW);
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('1分未満の場合「たった今」を返す', () => {
    const result = formatDistanceToNow('2026-02-19T11:59:30Z');
    expect(result).toBe('たった今');
  });

  it('数分前の場合「N分前」を返す', () => {
    const result = formatDistanceToNow('2026-02-19T11:45:00Z');
    expect(result).toBe('15分前');
  });

  it('1分前の場合「1分前」を返す', () => {
    const result = formatDistanceToNow('2026-02-19T11:59:00Z');
    expect(result).toBe('1分前');
  });

  it('59分前の場合「59分前」を返す', () => {
    const result = formatDistanceToNow('2026-02-19T11:01:00Z');
    expect(result).toBe('59分前');
  });

  it('数時間前の場合「N時間前」を返す', () => {
    const result = formatDistanceToNow('2026-02-19T09:00:00Z');
    expect(result).toBe('3時間前');
  });

  it('1時間前の場合「1時間前」を返す', () => {
    const result = formatDistanceToNow('2026-02-19T11:00:00Z');
    expect(result).toBe('1時間前');
  });

  it('数日前の場合「N日前」を返す', () => {
    const result = formatDistanceToNow('2026-02-16T12:00:00Z');
    expect(result).toBe('3日前');
  });

  it('1日前の場合「1日前」を返す', () => {
    const result = formatDistanceToNow('2026-02-18T12:00:00Z');
    expect(result).toBe('1日前');
  });
});
