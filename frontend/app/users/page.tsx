'use client';

import { useState, useEffect } from 'react';
import { apiClient, type User, type CreateUserRequest } from '@/lib/api';

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newUser, setNewUser] = useState<CreateUserRequest>({
    name: '',
    email: '',
  });

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      setLoading(true);
      const fetchedUsers = await apiClient.getUsers();
      setUsers(fetchedUsers);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch users');
    } finally {
      setLoading(false);
    }
  };

  const handleCreateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await apiClient.createUser(newUser);
      setNewUser({ name: '', email: '' });
      setShowCreateForm(false);
      await fetchUsers();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create user');
    }
  };

  const handleDeleteUser = async (userId: string) => {
    if (!window.confirm('このユーザーを削除しますか？')) {
      return;
    }
    try {
      await apiClient.deleteUser(userId);
      await fetchUsers();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete user');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen grid place-items-center">
        <div className="text-xl">Loading users...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen grid place-items-center">
        <div className="text-center">
          <div className="text-red-500 text-xl mb-4">Error: {error}</div>
          <button
            onClick={fetchUsers}
            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen grid grid-rows-[auto_1fr] p-4 sm:p-8">
      <header className="mb-6 sm:mb-8">
        <h1 className="text-2xl sm:text-4xl font-bold mb-4">Users</h1>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded transition-colors w-full sm:w-auto"
        >
          {showCreateForm ? 'Cancel' : 'Add New User'}
        </button>
      </header>

      <main>
        {showCreateForm && (
          <div className="mb-8 p-4 sm:p-6 border rounded-lg bg-gray-50">
            <h2 className="text-xl sm:text-2xl font-semibold mb-4">Create New User</h2>
            <form onSubmit={handleCreateUser} className="grid gap-4">
              <div>
                <label htmlFor="name" className="block text-sm font-medium mb-1">
                  Name
                </label>
                <input
                  type="text"
                  id="name"
                  value={newUser.name}
                  onChange={(e) =>
                    setNewUser({ ...newUser, name: e.target.value })
                  }
                  className="w-full p-2 sm:p-3 border rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  required
                />
              </div>
              <div>
                <label htmlFor="email" className="block text-sm font-medium mb-1">
                  Email
                </label>
                <input
                  type="email"
                  id="email"
                  value={newUser.email}
                  onChange={(e) =>
                    setNewUser({ ...newUser, email: e.target.value })
                  }
                  className="w-full p-2 sm:p-3 border rounded focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  required
                />
              </div>
              <button
                type="submit"
                className="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded w-fit transition-colors"
              >
                Create User
              </button>
            </form>
          </div>
        )}

        <div className="grid gap-4">
          {users.length === 0 ? (
            <div className="text-gray-500 text-center py-8">
              No users found.
            </div>
          ) : (
            users.map((user) => (
              <div key={user.id} className="p-4 sm:p-6 border rounded-lg bg-white shadow-sm hover:shadow-md transition-shadow">
                <div className="grid grid-cols-1 sm:grid-cols-[1fr_auto] gap-4 items-start">
                  <div>
                    <h3 className="text-lg sm:text-xl font-semibold break-words">{user.name}</h3>
                    <p className="text-gray-600 break-all">{user.email}</p>
                    <div className="text-sm text-gray-400 mt-2">
                      Created: {new Date(user.created_at).toLocaleDateString()}
                    </div>
                  </div>
                  <button
                    onClick={() => handleDeleteUser(user.id)}
                    className="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded text-sm transition-colors w-full sm:w-auto"
                  >
                    削除
                  </button>
                </div>
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}