import { describe, it, expect, vi } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useSubmitShortcut } from '../useSubmitShortcut';

const createKeyEvent = (overrides: Partial<React.KeyboardEvent> = {}) =>
  ({
    ctrlKey: false,
    metaKey: false,
    key: '',
    preventDefault: vi.fn(),
    ...overrides,
  }) as unknown as React.KeyboardEvent;

describe('useSubmitShortcut', () => {
  it('Ctrl+EnterでcanSubmit=trueの時にonSubmitが呼ばれる', () => {
    const onSubmit = vi.fn();
    const { result } = renderHook(() => useSubmitShortcut(onSubmit, true));

    const e = createKeyEvent({ ctrlKey: true, key: 'Enter' });
    result.current(e);

    expect(onSubmit).toHaveBeenCalledOnce();
    expect(e.preventDefault).toHaveBeenCalledOnce();
  });

  it('Cmd+Enter(metaKey)でもonSubmitが呼ばれる', () => {
    const onSubmit = vi.fn();
    const { result } = renderHook(() => useSubmitShortcut(onSubmit, true));

    const e = createKeyEvent({ metaKey: true, key: 'Enter' });
    result.current(e);

    expect(onSubmit).toHaveBeenCalledOnce();
  });

  it('canSubmit=falseの時はonSubmitが呼ばれない', () => {
    const onSubmit = vi.fn();
    const { result } = renderHook(() => useSubmitShortcut(onSubmit, false));

    const e = createKeyEvent({ ctrlKey: true, key: 'Enter' });
    result.current(e);

    expect(onSubmit).not.toHaveBeenCalled();
    expect(e.preventDefault).toHaveBeenCalledOnce();
  });

  it('Ctrl無しのEnterではonSubmitが呼ばれない', () => {
    const onSubmit = vi.fn();
    const { result } = renderHook(() => useSubmitShortcut(onSubmit, true));

    const e = createKeyEvent({ key: 'Enter' });
    result.current(e);

    expect(onSubmit).not.toHaveBeenCalled();
    expect(e.preventDefault).not.toHaveBeenCalled();
  });

  it('Ctrl+Enter以外のキーではonSubmitが呼ばれない', () => {
    const onSubmit = vi.fn();
    const { result } = renderHook(() => useSubmitShortcut(onSubmit, true));

    const e = createKeyEvent({ ctrlKey: true, key: 'a' });
    result.current(e);

    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('canSubmitがtrueからfalseに変わるとonSubmitが呼ばれなくなる', () => {
    const onSubmit = vi.fn();
    const { result, rerender } = renderHook(
      ({ canSubmit }) => useSubmitShortcut(onSubmit, canSubmit),
      { initialProps: { canSubmit: true } },
    );

    result.current(createKeyEvent({ ctrlKey: true, key: 'Enter' }));
    expect(onSubmit).toHaveBeenCalledOnce();

    onSubmit.mockClear();
    rerender({ canSubmit: false });

    result.current(createKeyEvent({ ctrlKey: true, key: 'Enter' }));
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('onSubmitコールバックが変更されても正しく動作する', () => {
    const onSubmit1 = vi.fn();
    const onSubmit2 = vi.fn();
    const { result, rerender } = renderHook(
      ({ onSubmit }) => useSubmitShortcut(onSubmit, true),
      { initialProps: { onSubmit: onSubmit1 } },
    );

    result.current(createKeyEvent({ ctrlKey: true, key: 'Enter' }));
    expect(onSubmit1).toHaveBeenCalledOnce();

    rerender({ onSubmit: onSubmit2 });

    result.current(createKeyEvent({ ctrlKey: true, key: 'Enter' }));
    expect(onSubmit2).toHaveBeenCalledOnce();
  });
});
