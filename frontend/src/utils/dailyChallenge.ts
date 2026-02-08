// Daily challenge keys — deterministically selected based on date
const CHALLENGE_KEYS = [
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
] as const;

export type ChallengeKey = (typeof CHALLENGE_KEYS)[number];

/**
 * Get today's challenge key using a date-seed deterministic selection.
 * Same date always returns the same challenge.
 */
export function getDailyChallenge(date: Date = new Date()): ChallengeKey {
  const dateStr = `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`;
  let hash = 0;
  for (let i = 0; i < dateStr.length; i++) {
    hash = (hash << 5) - hash + dateStr.charCodeAt(i);
    hash |= 0; // Convert to 32bit integer
  }
  const index = Math.abs(hash) % CHALLENGE_KEYS.length;
  return CHALLENGE_KEYS[index];
}

function getTodayKey(): string {
  const now = new Date();
  return `daily-challenge-${now.getFullYear()}-${now.getMonth() + 1}-${now.getDate()}`;
}

/**
 * Check if today's daily challenge has been completed.
 */
export function isChallengeCompleted(): boolean {
  try {
    return localStorage.getItem(getTodayKey()) === 'done';
  } catch {
    return false;
  }
}

/**
 * Mark today's daily challenge as completed.
 */
export function markChallengeCompleted(): void {
  try {
    localStorage.setItem(getTodayKey(), 'done');

    // Clean up old entries (keep last 7 days)
    const now = new Date();
    for (let i = 8; i < 60; i++) {
      const past = new Date(now);
      past.setDate(past.getDate() - i);
      const key = `daily-challenge-${past.getFullYear()}-${past.getMonth() + 1}-${past.getDate()}`;
      localStorage.removeItem(key);
    }
  } catch {
    // Silently fail if localStorage is unavailable
  }
}
