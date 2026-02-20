import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import CircleSearchCard from '../CircleSearchCard';
import type { StudyCircle } from '../../../types/studyCircle';

const mockCircle: StudyCircle = {
  id: 1,
  name: 'React Study Group',
  topic: 'React & TypeScript',
  description: 'A study group for React developers.',
  max_members: 20,
  member_count: 8,
  owner_id: 1,
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

const renderWithRouter = (ui: React.ReactElement) =>
  render(<MemoryRouter>{ui}</MemoryRouter>);

describe('CircleSearchCard', () => {
  it('サークル名が表示される', () => {
    renderWithRouter(<CircleSearchCard circle={mockCircle} />);
    expect(screen.getByText('React Study Group')).toBeInTheDocument();
  });

  it('トピックが表示される', () => {
    renderWithRouter(<CircleSearchCard circle={mockCircle} />);
    expect(screen.getByText('React & TypeScript')).toBeInTheDocument();
  });

  it('説明が表示される', () => {
    renderWithRouter(<CircleSearchCard circle={mockCircle} />);
    expect(screen.getByText('A study group for React developers.')).toBeInTheDocument();
  });

  it('メンバー数と上限が表示される', () => {
    renderWithRouter(<CircleSearchCard circle={mockCircle} />);
    expect(screen.getByText('8 / 20')).toBeInTheDocument();
  });

  it('説明が空の場合は表示されない', () => {
    const circleNoDesc = { ...mockCircle, description: '' };
    renderWithRouter(<CircleSearchCard circle={circleNoDesc} />);
    expect(screen.queryByText('A study group for React developers.')).not.toBeInTheDocument();
  });

  it('サークルへのリンクが正しいパスを持つ', () => {
    renderWithRouter(<CircleSearchCard circle={mockCircle} />);
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/study-circles/1');
  });
});
