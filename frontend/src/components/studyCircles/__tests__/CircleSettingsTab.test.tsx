import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import CircleSettingsTab from '../CircleSettingsTab';
import { getUsers } from '../../../api/users';
import type { StudyCircle } from '../../../types/studyCircle';
import type { User } from '../../../types/user';
import type { AxiosResponse } from 'axios';

vi.mock('../../../api/users');
vi.mock('../../../store/authStore', () => ({
  useAuthStore: vi.fn((selector: (s: { user: { id: 99, name: string } }) => unknown) =>
    selector({ user: { id: 99, name: 'オーナー' } }),
  ),
}));

const makeUser = (overrides: Partial<User>): User =>
  ({
    id: 0,
    username: 'user',
    name: 'ユーザー',
    email: 'user@example.com',
    avatar_url: '',
    ...overrides,
  }) as User;

const circle: StudyCircle = {
  id: 1,
  name: 'React 学習会',
  topic: 'React',
  description: '',
  owner_id: 99,
  max_members: 5,
  status: 'active',
  members: [
    { id: 1, circle_id: 1, user_id: 10, user: { id: 10, name: '既存メンバー', avatar_url: '' } },
  ],
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
} as StudyCircle;

const renderTab = (props: Partial<React.ComponentProps<typeof CircleSettingsTab>> = {}) =>
  render(
    <MemoryRouter>
      <CircleSettingsTab
        circle={circle}
        onAddMember={vi.fn().mockResolvedValue(true)}
        onRemoveMember={vi.fn()}
        {...props}
      />
    </MemoryRouter>,
  );

describe('CircleSettingsTab のメンバー追加検索', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(getUsers).mockResolvedValue({
      data: [
        makeUser({ id: 10, name: '既存メンバー' }),
        makeUser({ id: 20, name: '田中太郎' }),
        makeUser({ id: 30, name: '田中次郎' }),
      ],
    } as AxiosResponse<User[]>);
  });

  // useUserSearch の誤配線で、パネルを開くと searchUsers.length の参照でクラッシュしていた。
  it('メンバー追加パネルを開いてもクラッシュしない', async () => {
    const user = userEvent.setup();
    renderTab();

    await user.click(screen.getByRole('button', { name: /メンバーを追加/ }));

    expect(screen.getByPlaceholderText('ユーザーを検索...')).toBeInTheDocument();
  });

  it('検索語を入力すると候補が表示される', async () => {
    const user = userEvent.setup();
    renderTab();

    await user.click(screen.getByRole('button', { name: /メンバーを追加/ }));
    await user.type(screen.getByPlaceholderText('ユーザーを検索...'), '田中');

    expect(await screen.findByText('田中太郎')).toBeInTheDocument();
    expect(screen.getByText('田中次郎')).toBeInTheDocument();
  });

  it('既にメンバーのユーザーは候補に出ない', async () => {
    const user = userEvent.setup();
    renderTab();

    await user.click(screen.getByRole('button', { name: /メンバーを追加/ }));
    await user.type(screen.getByPlaceholderText('ユーザーを検索...'), 'メンバー');

    await screen.findByText('田中太郎');
    // メンバー一覧セクションには表示されるため、追加候補（ボタン）として出ていないことを見る
    expect(screen.queryByRole('button', { name: '既存メンバー' })).not.toBeInTheDocument();
  });

  it('候補をクリックすると onAddMember が呼ばれる', async () => {
    const onAddMember = vi.fn().mockResolvedValue(true);
    const user = userEvent.setup();
    renderTab({ onAddMember });

    await user.click(screen.getByRole('button', { name: /メンバーを追加/ }));
    await user.type(screen.getByPlaceholderText('ユーザーを検索...'), '田中');
    await user.click(await screen.findByText('田中太郎'));

    expect(onAddMember).toHaveBeenCalledWith(20);
  });

  it('検索語が空のときは候補を表示しない', async () => {
    const user = userEvent.setup();
    renderTab();

    await user.click(screen.getByRole('button', { name: /メンバーを追加/ }));

    expect(screen.queryByText('田中太郎')).not.toBeInTheDocument();
  });
});
