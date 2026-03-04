'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { getToken, removeToken } from '@/lib/auth';
import { authAPI, moduleAPI } from '@/lib/api';

interface Module {
  moduleID: string;
  title: string;
  code: string;
  description: string;
  credits: number;
  semester: number;
}

export default function DashboardPage() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [modules, setModules] = useState<Module[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push('/login');
      return;
    }

    const loadData = async () => {
      try {
        const profile = await authAPI.getProfile(token);
        setUser(profile);

        const modulesData = await moduleAPI.getAll(token);
        setModules(modulesData.modules || []);
      } catch (error) {
        console.error('Failed to load data:', error);
        removeToken();
        router.push('/login');
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [router]);

  const handleLogout = () => {
    removeToken();
    router.push('/login');
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-6 flex justify-between items-center">
          <h1 className="text-3xl font-bold text-gray-800">UniEnroll</h1>
          <div className="flex items-center gap-4">
            <span className="text-gray-700">{user?.name}</span>
            <span className="px-3 py-1 bg-blue-600 text-white rounded-full text-sm capitalize">
              {user?.role}
            </span>
            {user?.role === 'admin' && (
              <Link href="/admin" className="text-blue-600 font-semibold hover:underline">
                Admin Panel
              </Link>
            )}
            <button
              onClick={handleLogout}
              className="bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700"
            >
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-gray-600 text-sm">Name</h3>
            <p className="text-2xl font-bold text-gray-800">{user?.name}</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-gray-600 text-sm">Email</h3>
            <p className="text-2xl font-bold text-gray-800">{user?.email}</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow">
            <h3 className="text-gray-600 text-sm">Role</h3>
            <p className="text-2xl font-bold text-gray-800 capitalize">{user?.role}</p>
          </div>
        </div>

        {/* Modules Section */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-bold mb-6">Available Modules</h2>
          
          {modules.length === 0 ? (
            <p className="text-gray-600">No modules available yet.</p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {modules.map((module) => (
                <div
                  key={module.moduleID}
                  className="border border-gray-300 rounded-lg p-4 hover:shadow-lg transition"
                >
                  <h3 className="text-xl font-semibold text-gray-800 mb-2">{module.title}</h3>
                  <p className="text-sm text-gray-600 mb-2">Code: {module.code}</p>
                  <p className="text-gray-700 mb-3">{module.description}</p>
                  <div className="flex justify-between text-sm text-gray-600">
                    <span>{module.credits} Credits</span>
                    <span>Semester {module.semester}</span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
