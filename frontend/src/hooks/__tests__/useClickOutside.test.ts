import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useClickOutside } from '../useClickOutside';

describe('useClickOutside', () => {
  let container: HTMLDivElement;
  let outside: HTMLDivElement;

  beforeEach(() => {
    container = document.createElement('div');
    outside = document.createElement('div');
    document.body.appendChild(container);
    document.body.appendChild(outside);
  });

  const createRef = (el: HTMLDivElement) => ({ current: el });

  it('isOpen=falseの時は外側クリックでonCloseが呼ばれない', () => {
    const onClose = vi.fn();
    renderHook(() => useClickOutside(createRef(container), false, onClose));

    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).not.toHaveBeenCalled();
  });

  it('isOpen=trueで外側クリック時にonCloseが呼ばれる', () => {
    const onClose = vi.fn();
    renderHook(() => useClickOutside(createRef(container), true, onClose));

    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).toHaveBeenCalledOnce();
  });

  it('isOpen=trueでref内部クリック時にonCloseが呼ばれない', () => {
    const onClose = vi.fn();
    renderHook(() => useClickOutside(createRef(container), true, onClose));

    container.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).not.toHaveBeenCalled();
  });

  it('子要素クリック時もonCloseが呼ばれない', () => {
    const child = document.createElement('button');
    container.appendChild(child);
    const onClose = vi.fn();
    renderHook(() => useClickOutside(createRef(container), true, onClose));

    child.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).not.toHaveBeenCalled();
  });

  it('isOpenがtrueからfalseに変わるとリスナーが解除される', () => {
    const onClose = vi.fn();
    const { rerender } = renderHook(
      ({ isOpen }) => useClickOutside(createRef(container), isOpen, onClose),
      { initialProps: { isOpen: true } },
    );

    // isOpen=trueの間は外側クリックで呼ばれる
    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).toHaveBeenCalledOnce();

    onClose.mockClear();
    rerender({ isOpen: false });

    // isOpen=falseになると呼ばれない
    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).not.toHaveBeenCalled();
  });

  it('isOpenがfalseからtrueに変わるとリスナーが登録される', () => {
    const onClose = vi.fn();
    const { rerender } = renderHook(
      ({ isOpen }) => useClickOutside(createRef(container), isOpen, onClose),
      { initialProps: { isOpen: false } },
    );

    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).not.toHaveBeenCalled();

    rerender({ isOpen: true });

    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).toHaveBeenCalledOnce();
  });

  it('refがnullの場合はエラーにならない', () => {
    const onClose = vi.fn();
    renderHook(() => useClickOutside({ current: null }, true, onClose));

    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    // ref.currentがnullの場合、containsチェックが通らないのでonCloseは呼ばれない
    expect(onClose).not.toHaveBeenCalled();
  });

  it('複数回の外側クリックで毎回onCloseが呼ばれる', () => {
    const onClose = vi.fn();
    renderHook(() => useClickOutside(createRef(container), true, onClose));

    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose).toHaveBeenCalledTimes(3);
  });

  it('onCloseコールバックが変更されても正しく動作する', () => {
    const onClose1 = vi.fn();
    const onClose2 = vi.fn();
    const { rerender } = renderHook(
      ({ onClose }) => useClickOutside(createRef(container), true, onClose),
      { initialProps: { onClose: onClose1 } },
    );

    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose1).toHaveBeenCalledOnce();

    rerender({ onClose: onClose2 });

    outside.dispatchEvent(new MouseEvent('mousedown', { bubbles: true }));
    expect(onClose2).toHaveBeenCalledOnce();
  });
});
