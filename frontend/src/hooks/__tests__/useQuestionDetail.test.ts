import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useQuestionDetail } from '../useQuestionDetail';
import {
  getQuestionById,
  createAnswer,
  setBestAnswer,
  voteQuestion,
  voteAnswer,
} from '../../api/qa';
import type { Question, Answer } from '../../types/qa';
import client from '../../api/client';

vi.mock('../../api/qa', () => ({
  getQuestionById: vi.fn(),
  createAnswer: vi.fn(),
  updateAnswer: vi.fn(),
  deleteAnswer: vi.fn(),
  setBestAnswer: vi.fn(),
  voteQuestion: vi.fn(),
  removeQuestionVote: vi.fn(),
  voteAnswer: vi.fn(),
  removeAnswerVote: vi.fn(),
}));

vi.mock('../../api/client', () => ({
  default: { get: vi.fn() },
}));

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn(), success: vi.fn() },
}));

const mockQuestion = {
  id: 1, title: 'テスト質問', body: '質問内容', user_id: 1,
  vote_count: 3, answer_count: 2, is_solved: false,
  created_at: '2026-02-19', updated_at: '2026-02-19',
} as Question;
const mockAnswers = [
  { id: 10, body: '回答1', user_id: 2, is_best: false },
  { id: 11, body: '回答2', user_id: 3, is_best: false },
];

describe('useQuestionDetail', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('undefinedのidでは質問がnullであること', () => {
    const { result } = renderHook(() => useQuestionDetail(undefined));

    expect(result.current.question).toBeNull();
    expect(result.current.answers).toEqual([]);
    expect(result.current.loading).toBe(false);
  });

  it('質問と回答が正しく取得されること', async () => {
    vi.mocked(getQuestionById).mockResolvedValue({ question: mockQuestion, user_vote: 1 });
    vi.mocked(client.get).mockResolvedValue({ data: mockAnswers });

    const { result } = renderHook(() => useQuestionDetail('1'));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.question).toEqual(mockQuestion);
    expect(result.current.answers).toEqual(mockAnswers);
    expect(result.current.userVote).toBe(1);
  });

  it('回答作成が成功すること', async () => {
    vi.mocked(getQuestionById).mockResolvedValue({ question: mockQuestion, user_vote: 0 });
    vi.mocked(client.get).mockResolvedValue({ data: mockAnswers });

    const newAnswer = { id: 12, body: '新しい回答', user_id: 1, is_best: false } as Answer;
    vi.mocked(createAnswer).mockResolvedValue(newAnswer);

    const { result } = renderHook(() => useQuestionDetail('1'));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let success: boolean | undefined;
    await act(async () => {
      success = await result.current.createAnswer({ body: '新しい回答' });
    });

    expect(success).toBe(true);
    expect(createAnswer).toHaveBeenCalledWith(1, { body: '新しい回答' });
  });

  it('空白の回答作成はスキップされること', async () => {
    vi.mocked(getQuestionById).mockResolvedValue({ question: mockQuestion, user_vote: 0 });
    vi.mocked(client.get).mockResolvedValue({ data: [] });

    const { result } = renderHook(() => useQuestionDetail('1'));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    let success: boolean | undefined;
    await act(async () => {
      success = await result.current.createAnswer({ body: '   ' });
    });

    expect(success).toBeFalsy();
    expect(createAnswer).not.toHaveBeenCalled();
  });

  it('質問投票が正しく動作すること', async () => {
    vi.mocked(getQuestionById).mockResolvedValue({ question: mockQuestion, user_vote: 0 });
    vi.mocked(client.get).mockResolvedValue({ data: [] });
    vi.mocked(voteQuestion).mockResolvedValue({} as unknown as void);

    const { result } = renderHook(() => useQuestionDetail('1'));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.voteQuestion(1);
    });

    expect(voteQuestion).toHaveBeenCalledWith(1, { value: 1 });
  });

  it('ベストアンサー設定が正しく動作すること', async () => {
    vi.mocked(getQuestionById).mockResolvedValue({ question: mockQuestion, user_vote: 0 });
    vi.mocked(client.get).mockResolvedValue({ data: mockAnswers });
    vi.mocked(setBestAnswer).mockResolvedValue({} as unknown as void);

    const { result } = renderHook(() => useQuestionDetail('1'));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.setBestAnswer(10);
    });

    expect(setBestAnswer).toHaveBeenCalledWith(1, 10);
  });

  it('回答投票が正しく動作すること', async () => {
    vi.mocked(getQuestionById).mockResolvedValue({ question: mockQuestion, user_vote: 0 });
    vi.mocked(client.get).mockResolvedValue({ data: mockAnswers });
    vi.mocked(voteAnswer).mockResolvedValue({} as unknown as void);

    const { result } = renderHook(() => useQuestionDetail('1'));

    await vi.waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    await act(async () => {
      await result.current.voteAnswer(10, 1);
    });

    expect(voteAnswer).toHaveBeenCalledWith(1, 10, { value: 1 });
  });
});
