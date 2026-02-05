import { useState } from 'react';
import { Link } from 'react-router-dom';
import { getUsers } from '../api/users';
import type { User } from '../types/user';
import Avatar from '../components/common/Avatar';

export default function SearchPage() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<User[]>([]);
  const [searched, setSearched] = useState(false);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;
    const { data } = await getUsers(query);
    setResults(data || []);
    setSearched(true);
  };

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-2xl font-bold">Search Users</h1>
      <form onSubmit={handleSearch} className="flex gap-2">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search by name or email..."
          className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-md text-white focus:outline-none focus:border-blue-500"
        />
        <button
          type="submit"
          className="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors"
        >
          Search
        </button>
      </form>

      {searched && results.length === 0 && (
        <div className="text-center text-gray-400 py-8">No users found</div>
      )}

      <div className="space-y-2">
        {results.map((user) => (
          <Link
            key={user.id}
            to={`/profile/${user.id}`}
            className="flex items-center gap-3 p-4 bg-gray-900 rounded-lg hover:bg-gray-800 transition-colors"
          >
            <Avatar name={user.name} avatarUrl={user.avatar_url} />
            <div>
              <div className="font-medium">{user.name}</div>
              <div className="text-sm text-gray-400">{user.email}</div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}
