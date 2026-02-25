import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import StepIndicator from '../StepIndicator';

const steps = [
  { label: 'アカウント作成' },
  { label: 'プロフィール設定' },
  { label: '完了' },
];

describe('StepIndicator', () => {
  it('全てのステップラベルが表示される', () => {
    render(<StepIndicator steps={steps} currentStep={0} />);

    expect(screen.getByText('アカウント作成')).toBeInTheDocument();
    expect(screen.getByText('プロフィール設定')).toBeInTheDocument();
    expect(screen.getByText('完了')).toBeInTheDocument();
  });

  it('ステップ番号が表示される', () => {
    render(<StepIndicator steps={steps} currentStep={0} />);

    expect(screen.getByText('1')).toBeInTheDocument();
    expect(screen.getByText('2')).toBeInTheDocument();
    expect(screen.getByText('3')).toBeInTheDocument();
  });

  it('現在のステップがアクティブスタイル', () => {
    const { container } = render(<StepIndicator steps={steps} currentStep={0} />);

    const activeCircles = container.querySelectorAll('.rounded-full.bg-blue-600');
    expect(activeCircles.length).toBe(1);
  });

  it('完了したステップがチェックマーク表示', () => {
    const { container } = render(<StepIndicator steps={steps} currentStep={2} />);

    const checkIcons = container.querySelectorAll('.lucide-check');
    expect(checkIcons.length).toBe(2);
  });

  it('未完了のステップがグレースタイル', () => {
    const { container } = render(<StepIndicator steps={steps} currentStep={0} />);

    const grayCircles = container.querySelectorAll('.rounded-full.bg-gray-700');
    expect(grayCircles.length).toBe(2);
  });

  it('ステップ間の接続線が表示される', () => {
    const { container } = render(<StepIndicator steps={steps} currentStep={0} />);

    const connectors = container.querySelectorAll('.h-0\\.5');
    expect(connectors.length).toBe(2);
  });

  it('完了したステップの接続線がアクティブ', () => {
    const { container } = render(<StepIndicator steps={steps} currentStep={2} />);

    const activeConnectors = container.querySelectorAll('.bg-blue-600.h-0\\.5');
    expect(activeConnectors.length).toBe(2);
  });

  it('カスタムクラス名が適用される', () => {
    const { container } = render(
      <StepIndicator steps={steps} currentStep={0} className="custom-class" />
    );

    expect(container.querySelector('.custom-class')).toBeInTheDocument();
  });

  it('ステップ1つでも表示される', () => {
    render(<StepIndicator steps={[{ label: '唯一' }]} currentStep={0} />);

    expect(screen.getByText('唯一')).toBeInTheDocument();
  });

  it('説明文が表示される', () => {
    const stepsWithDesc = [
      { label: 'ステップ1', description: '詳細説明' },
    ];
    render(<StepIndicator steps={stepsWithDesc} currentStep={0} />);

    expect(screen.getByText('詳細説明')).toBeInTheDocument();
  });
});
