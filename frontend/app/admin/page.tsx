'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { getToken, removeToken } from '@/lib/auth';
import { adminAPI, authAPI, moduleAPI } from '@/lib/api';

interface Module {
  moduleID: string;
  title: string;
  code: string;
  description: string;
  credits: number;
  semester: number;
  type?: string;
  departmentID?: string;
  yearOfStudy?: number;
}

interface ImportSummary {
  message: string;
  processed: number;
  created: number;
  skipped: number;
  failures: Array<{ row: number; email?: string; error: string }>;
  credentials: Array<{ name: string; email: string; userID: string; tempPassword: string }>;
}

export default function AdminPage() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [modules, setModules] = useState<Module[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    title: '',
    code: '',
    departmentID: '',
    courseID: '',
    lecturerID: '',
    type: 'compulsory',
    description: '',
    credits: 3,
    semester: 1,
    yearOfStudy: 1,
    prerequisites: [] as string[],
  });
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState('');
  const [csvFile, setCsvFile] = useState<File | null>(null);
  const [importing, setImporting] = useState(false);
  const [importError, setImportError] = useState('');
  const [importSummary, setImportSummary] = useState<ImportSummary | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      router.push('/login');
      return;
    }

    const loadData = async () => {
      try {
        const profile = await authAPI.getProfile(token);
        
        // Check if user is admin
        if (profile.role !== 'admin') {
          router.push('/dashboard');
          return;
        }
        
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

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'credits' || name === 'semester' || name === 'yearOfStudy' ? parseInt(value) : value,
    }));
  };

  const handleCreateModule = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setMessage('');

    try {
      const token = getToken();
      if (!token) throw new Error('No token found');

      await moduleAPI.create(token, formData);
      
      setMessage('Module created successfully!');
      setFormData({
        title: '',
        code: '',
        departmentID: '',
        courseID: '',
        lecturerID: '',
        type: 'compulsory',
        description: '',
        credits: 3,
        semester: 1,
        yearOfStudy: 1,
        prerequisites: [],
      });
      setShowForm(false);

      // Reload modules
      const modulesData = await moduleAPI.getAll(token);
      setModules(modulesData.modules || []);
    } catch (error) {
      setMessage(error instanceof Error ? error.message : 'Failed to create module');
    } finally {
      setSubmitting(false);
    }
  };

  const handleLogout = () => {
    removeToken();
    router.push('/login');
  };

  const handleImportStudents = async (e: React.FormEvent) => {
    e.preventDefault();
    setImportError('');
    setImportSummary(null);

    if (!csvFile) {
      setImportError('Please choose a CSV file first');
      return;
    }

    try {
      setImporting(true);
      const token = getToken();
      if (!token) throw new Error('No token found');

      const result = await adminAPI.importStudentsCSV(token, csvFile);
      setImportSummary(result);
      setCsvFile(null);
    } catch (error) {
      setImportError(error instanceof Error ? error.message : 'Failed to import students');
    } finally {
      setImporting(false);
    }
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
          <h1 className="text-3xl font-bold text-gray-800">UniEnroll - Admin Panel</h1>
          <div className="flex items-center gap-4">
            <Link href="/profile" className="text-gray-700 hover:text-blue-600 font-semibold">
              {user?.name}
            </Link>
            <span className="px-3 py-1 bg-red-600 text-white rounded-full text-sm capitalize">
              {user?.role}
            </span>
            <Link href="/departments" className="text-blue-600 font-semibold hover:underline">
              Departments
            </Link>
            <Link href="/dashboard" className="text-blue-600 font-semibold hover:underline">
              Dashboard
            </Link>
            <Link href="/profile" className="text-blue-600 font-semibold hover:underline">
              Profile
            </Link>
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
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-2xl font-bold mb-2">Import Students (CSV)</h2>
          <p className="text-sm text-gray-600 mb-4">
            Required columns: <strong>name, email, departmentID</strong>. Optional: yearOfStudy, studentNumber, enrollmentYear.
          </p>

          <form onSubmit={handleImportStudents} className="flex flex-col md:flex-row gap-3 md:items-center">
            <input
              type="file"
              accept=".csv,text/csv"
              onChange={(e) => setCsvFile(e.target.files?.[0] || null)}
              className="block w-full md:w-auto text-sm text-gray-700 file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:bg-blue-600 file:text-white hover:file:bg-blue-700"
            />
            <button
              type="submit"
              disabled={importing}
              className="bg-green-600 text-white px-4 py-2 rounded-lg hover:bg-green-700 disabled:bg-gray-400"
            >
              {importing ? 'Importing...' : 'Upload CSV'}
            </button>
          </form>

          {importError && (
            <div className="mt-4 p-3 rounded bg-red-100 text-red-700">{importError}</div>
          )}

          {importSummary && (
            <div className="mt-4 space-y-3">
              <div className="p-3 rounded bg-green-100 text-green-800">
                {importSummary.message} | Processed: {importSummary.processed}, Created: {importSummary.created}, Skipped: {importSummary.skipped}
              </div>

              {importSummary.failures.length > 0 && (
                <div className="p-3 rounded bg-yellow-100 text-yellow-900">
                  <p className="font-semibold mb-1">Rows with issues:</p>
                  <ul className="text-sm max-h-40 overflow-auto space-y-1">
                    {importSummary.failures.map((failure, index) => (
                      <li key={`${failure.row}-${index}`}>
                        Row {failure.row}{failure.email ? ` (${failure.email})` : ''}: {failure.error}
                      </li>
                    ))}
                  </ul>
                </div>
              )}

              {importSummary.credentials.length > 0 && (
                <div className="p-3 rounded bg-blue-50 border border-blue-200">
                  <p className="font-semibold text-blue-900 mb-1">Temporary credentials (save these now):</p>
                  <div className="max-h-40 overflow-auto text-sm space-y-1 text-blue-900">
                    {importSummary.credentials.map((item) => (
                      <p key={item.userID}>
                        {item.name} ({item.email}) - {item.tempPassword}
                      </p>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Add Module Section */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-2xl font-bold">Modules Management</h2>
            <button
              onClick={() => setShowForm(!showForm)}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700"
            >
              {showForm ? 'Cancel' : 'Add New Module'}
            </button>
          </div>

          {message && (
            <div className={`mb-4 p-3 rounded-lg ${message.includes('successfully') ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
              {message}
            </div>
          )}

          {showForm && (
            <form onSubmit={handleCreateModule} className="bg-gray-50 p-6 rounded-lg space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Module Title *</label>
                  <input
                    type="text"
                    name="title"
                    value={formData.title}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                    placeholder="e.g., Introduction to Computer Science"
                    required
                  />
                </div>

                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Module Code *</label>
                  <input
                    type="text"
                    name="code"
                    value={formData.code}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                    placeholder="e.g., CS101"
                    required
                  />
                </div>

                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Department ID *</label>
                  <input
                    type="text"
                    name="departmentID"
                    value={formData.departmentID}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                    placeholder="e.g., dept-123"
                    required
                  />
                </div>

                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Module Type *</label>
                  <select
                    name="type"
                    value={formData.type}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                    required
                  >
                    <option value="compulsory">Compulsory</option>
                    <option value="elective">Elective</option>
                    <option value="faculty">Faculty (Cross-Department)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Course ID *</label>
                  <input
                    type="text"
                    name="courseID"
                    value={formData.courseID}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                    placeholder="e.g., course-123"
                    required
                  />
                </div>

                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Lecturer ID *</label>
                  <input
                    type="text"
                    name="lecturerID"
                    value={formData.lecturerID}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                    placeholder="e.g., lecturer-456"
                    required
                  />
                </div>

                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Year of Study</label>
                  <input
                    type="number"
                    name="yearOfStudy"
                    value={formData.yearOfStudy}
                    onChange={handleInputChange}
                    min="1"
                    max="5"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Credits</label>
                  <input
                    type="number"
                    name="credits"
                    value={formData.credits}
                    onChange={handleInputChange}
                    min="1"
                    max="10"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-gray-700 font-semibold mb-2">Semester</label>
                  <input
                    type="number"
                    name="semester"
                    value={formData.semester}
                    onChange={handleInputChange}
                    min="1"
                    max="8"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                  />
                </div>
              </div>

              <div>
                <label className="block text-gray-700 font-semibold mb-2">Description</label>
                <textarea
                  name="description"
                  value={formData.description}
                  onChange={handleInputChange}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
                  placeholder="Module description"
                  rows={4}
                />
              </div>

              <button
                type="submit"
                disabled={submitting}
                className="w-full bg-blue-600 text-white py-2 rounded-lg font-semibold hover:bg-blue-700 disabled:bg-gray-400"
              >
                {submitting ? 'Creating...' : 'Create Module'}
              </button>
            </form>
          )}
        </div>

        {/* Modules List */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-bold mb-6">All Modules ({modules.length})</h2>
          
          {modules.length === 0 ? (
            <p className="text-gray-600">No modules yet. Create one to get started!</p>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {modules.map((module) => {
                const getTypeColor = (type?: string) => {
                  switch (type) {
                    case 'compulsory': return 'bg-red-100 text-red-700';
                    case 'elective': return 'bg-blue-100 text-blue-700';
                    case 'faculty': return 'bg-green-100 text-green-700';
                    default: return 'bg-gray-100 text-gray-700';
                  }
                };
                
                const getTypeName = (type?: string) => {
                  switch (type) {
                    case 'compulsory': return 'Compulsory';
                    case 'elective': return 'Elective';
                    case 'faculty': return 'Faculty';
                    default: return 'Module';
                  }
                };
                
                return (
                  <div
                    key={module.moduleID}
                    className="border border-gray-300 rounded-lg p-4 hover:shadow-lg transition"
                  >
                    <div className="flex justify-between items-start mb-2">
                      <h3 className="text-lg font-semibold text-gray-800">{module.title}</h3>
                      {module.type && (
                        <span className={`px-2 py-1 rounded text-xs font-semibold ${getTypeColor(module.type)}`}>
                          {getTypeName(module.type)}
                        </span>
                      )}
                    </div>
                    <p className="text-sm font-mono text-blue-600 mb-2">{module.code}</p>
                    <p className="text-gray-700 text-sm mb-3">{module.description}</p>
                    <div className="flex justify-between text-sm text-gray-600 border-t pt-3">
                      <span>💳 {module.credits} Credits</span>
                      <span>📚 Sem {module.semester}</span>
                    </div>
                    {module.yearOfStudy && (
                      <div className="mt-2 text-xs text-gray-500">
                        Year {module.yearOfStudy}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
