'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { getToken, removeToken } from '@/lib/auth';
import { authAPI } from '@/lib/api';

interface UserProfile {
  id: string;
  name: string;
  email: string;
  role: string;
  departmentID?: string;
  yearOfStudy?: number;
}

export default function ProfilePage() {
  const router = useRouter();
  const [user, setUser] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push('/login');
      return;
    }

    const loadProfile = async () => {
      try {
        const profile = await authAPI.getProfile(token);
        setUser(profile);
      } catch (error) {
        console.error('Failed to load profile:', error);
        removeToken();
        router.push('/login');
      } finally {
        setLoading(false);
      }
    };

    loadProfile();
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

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-6 flex justify-between items-center">
          <h1 className="text-3xl font-bold text-gray-800">UniEnroll</h1>
          <div className="flex items-center gap-4">
            <Link href="/dashboard" className="text-blue-600 font-semibold hover:underline">
              Dashboard
            </Link>
            {user.role === 'admin' && (
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
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="bg-white rounded-lg shadow-lg overflow-hidden">
          {/* Profile Header */}
          <div className="bg-gradient-to-r from-blue-600 to-purple-600 px-8 py-12">
            <div className="flex items-center gap-6">
              <div className="w-24 h-24 bg-white rounded-full flex items-center justify-center text-4xl font-bold text-blue-600">
                {user.name.charAt(0).toUpperCase()}
              </div>
              <div className="text-white">
                <h1 className="text-3xl font-bold mb-2">{user.name}</h1>
                <p className="text-blue-100 text-lg">{user.email}</p>
                <span className="inline-block mt-2 px-4 py-1 bg-white bg-opacity-20 rounded-full text-sm font-semibold capitalize">
                  {user.role}
                </span>
              </div>
            </div>
          </div>

          {/* Profile Details */}
          <div className="px-8 py-6">
            <h2 className="text-2xl font-bold text-gray-800 mb-6">Profile Information</h2>
            
            <div className="space-y-6">
              {/* Basic Information */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="border-l-4 border-blue-600 pl-4">
                  <label className="text-sm font-semibold text-gray-600 uppercase">User ID</label>
                  <p className="text-lg text-gray-800 font-mono">{user.id}</p>
                </div>
                
                <div className="border-l-4 border-blue-600 pl-4">
                  <label className="text-sm font-semibold text-gray-600 uppercase">Account Type</label>
                  <p className="text-lg text-gray-800 capitalize">{user.role}</p>
                </div>
                
                <div className="border-l-4 border-blue-600 pl-4">
                  <label className="text-sm font-semibold text-gray-600 uppercase">Full Name</label>
                  <p className="text-lg text-gray-800">{user.name}</p>
                </div>
                
                <div className="border-l-4 border-blue-600 pl-4">
                  <label className="text-sm font-semibold text-gray-600 uppercase">Email Address</label>
                  <p className="text-lg text-gray-800">{user.email}</p>
                </div>
              </div>

              {/* Student-Specific Information */}
              {user.role === 'student' && (
                <div className="mt-8 pt-6 border-t border-gray-200">
                  <h3 className="text-xl font-bold text-gray-800 mb-4">Student Information</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {user.departmentID && (
                      <div className="border-l-4 border-green-600 pl-4">
                        <label className="text-sm font-semibold text-gray-600 uppercase">Department ID</label>
                        <p className="text-lg text-gray-800 font-mono">{user.departmentID}</p>
                      </div>
                    )}
                    
                    {user.yearOfStudy && (
                      <div className="border-l-4 border-green-600 pl-4">
                        <label className="text-sm font-semibold text-gray-600 uppercase">Year of Study</label>
                        <p className="text-lg text-gray-800">Year {user.yearOfStudy}</p>
                      </div>
                    )}
                  </div>
                  
                  {(!user.departmentID || !user.yearOfStudy) && (
                    <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
                      <p className="text-yellow-800">
                        <span className="font-semibold">⚠️ Profile Incomplete:</span> Please contact admin to complete your student profile with department and year information.
                      </p>
                    </div>
                  )}
                </div>
              )}

              {/* Admin Information */}
              {user.role === 'admin' && (
                <div className="mt-8 pt-6 border-t border-gray-200">
                  <h3 className="text-xl font-bold text-gray-800 mb-4">Admin Privileges</h3>
                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <ul className="space-y-2 text-gray-700">
                      <li className="flex items-center gap-2">
                        <span className="text-blue-600">✓</span> Manage departments
                      </li>
                      <li className="flex items-center gap-2">
                        <span className="text-blue-600">✓</span> Create and manage modules
                      </li>
                      <li className="flex items-center gap-2">
                        <span className="text-blue-600">✓</span> Create student, lecturer, and admin profiles
                      </li>
                      <li className="flex items-center gap-2">
                        <span className="text-blue-600">✓</span> View all enrollments
                      </li>
                    </ul>
                  </div>
                </div>
              )}

              {/* Lecturer Information */}
              {user.role === 'lecturer' && (
                <div className="mt-8 pt-6 border-t border-gray-200">
                  <h3 className="text-xl font-bold text-gray-800 mb-4">Lecturer Privileges</h3>
                  <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
                    <ul className="space-y-2 text-gray-700">
                      <li className="flex items-center gap-2">
                        <span className="text-purple-600">✓</span> View assigned modules
                      </li>
                      <li className="flex items-center gap-2">
                        <span className="text-purple-600">✓</span> View student enrollments
                      </li>
                      <li className="flex items-center gap-2">
                        <span className="text-purple-600">✓</span> Manage course content
                      </li>
                    </ul>
                  </div>
                </div>
              )}
            </div>

            {/* Action Buttons */}
            <div className="mt-8 pt-6 border-t border-gray-200 flex gap-4">
              <Link
                href="/dashboard"
                className="flex-1 bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 text-center font-semibold"
              >
                Back to Dashboard
              </Link>
              {user.role === 'admin' && (
                <Link
                  href="/admin"
                  className="flex-1 bg-purple-600 text-white px-6 py-3 rounded-lg hover:bg-purple-700 text-center font-semibold"
                >
                  Go to Admin Panel
                </Link>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
