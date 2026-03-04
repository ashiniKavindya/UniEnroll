'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { getToken, removeToken } from '@/lib/auth';
import { authAPI, moduleAPI, enrollmentAPI } from '@/lib/api';

interface Module {
  moduleID: string;
  title: string;
  code: string;
  description: string;
  credits: number;
  semester: number;
  type?: string;
  departmentID?: string;
}

interface Enrollment {
  enrollmentID: string;
  moduleID: string;
  status: string;
  grade?: string;
  enrolledAt: string;
}

export default function DashboardPage() {
  const router = useRouter();
  const [user, setUser] = useState<any>(null);
  const [modules, setModules] = useState<Module[]>([]);
  const [enrollments, setEnrollments] = useState<Enrollment[]>([]);
  const [enrolledModuleIds, setEnrolledModuleIds] = useState<Set<string>>(new Set());
  const [loading, setLoading] = useState(true);
  const [enrollingModule, setEnrollingModule] = useState<string | null>(null);
  const [filterType, setFilterType] = useState<string>('all');
  const [departmentID, setDepartmentID] = useState<string>('');

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
        
        // Set department ID for students
        if (profile.role === 'student' && profile.departmentID) {
          setDepartmentID(profile.departmentID);
        }

        const modulesData = await moduleAPI.getAll(token);
        setModules(modulesData.modules || []);

        // If student, load enrollments
        if (profile.role === 'student') {
          try {
            const enrollmentsData = await enrollmentAPI.getStudentEnrollments(token, profile.id);
            setEnrollments(enrollmentsData.enrollments || []);
            const enrolledIds = new Set<string>((enrollmentsData.enrollments || []).map((e: Enrollment) => e.moduleID));
            setEnrolledModuleIds(enrolledIds);
          } catch (err) {
            console.log('No enrollments yet');
          }
        }
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

  const handleEnroll = async (moduleID: string) => {
    if (!user?.id) return;
    
    setEnrollingModule(moduleID);
    try {
      const token = getToken();
      if (!token) return;

      await enrollmentAPI.enroll(token, user.id, moduleID);
      
      // Reload enrollments
      const enrollmentsData = await enrollmentAPI.getStudentEnrollments(token, user.id);
      setEnrollments(enrollmentsData.enrollments || []);
      const enrolledIds = new Set<string>((enrollmentsData.enrollments || []).map((e: Enrollment) => e.moduleID));
      setEnrolledModuleIds(enrolledIds);
      
      alert('Enrolled successfully!');
    } catch (error) {
      alert(error instanceof Error ? error.message : 'Failed to enroll');
    } finally {
      setEnrollingModule(null);
    }
  };

  const handleLogout = () => {
    removeToken();
    router.push('/login');
  };

  const filteredModules = filterType === 'all' 
    ? modules 
    : modules.filter(m => m.type === filterType);

  // Separate modules for students by department
  const departmentModules = user?.role === 'student' && departmentID
    ? filteredModules.filter(m => m.departmentID === departmentID)
    : [];
  
  const electivesFromOtherDepartments = user?.role === 'student' && departmentID
    ? filteredModules.filter(m => 
        m.departmentID !== departmentID && 
        (m.type === 'elective' || m.type === 'faculty')
      )
    : [];

  const modulesToDisplay = user?.role === 'student' 
    ? [] // Will display separately
    : filteredModules;

  const getModuleTypeColor = (type?: string) => {
    switch (type) {
      case 'compulsory': return 'bg-red-100 text-red-700';
      case 'elective': return 'bg-blue-100 text-blue-700';
      case 'faculty': return 'bg-green-100 text-green-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getModuleTypeBadge = (type?: string) => {
    switch (type) {
      case 'compulsory': return 'Compulsory';
      case 'elective': return 'Elective';
      case 'faculty': return 'Faculty';
      default: return 'Module';
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
          <h1 className="text-3xl font-bold text-gray-800">UniEnroll</h1>
          <div className="flex items-center gap-4">
            <Link href="/profile" className="text-gray-700 hover:text-blue-600 font-semibold">
              {user?.name}
            </Link>
            <span className="px-3 py-1 bg-blue-600 text-white rounded-full text-sm capitalize">
              {user?.role}
            </span>
            {user?.role === 'admin' && (
              <Link href="/admin" className="text-blue-600 font-semibold hover:underline">
                Admin Panel
              </Link>
            )}
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

        {/* Student Enrollments */}
        {user?.role === 'student' && enrollments.length > 0 && (
          <div className="bg-white rounded-lg shadow p-6 mb-8">
            <h2 className="text-2xl font-bold mb-6">My Enrolled Modules ({enrollments.length})</h2>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {enrollments.map((enrollment) => {
                const module = modules.find(m => m.moduleID === enrollment.moduleID);
                if (!module) return null;
                
                return (
                  <div
                    key={enrollment.enrollmentID}
                    className="border-l-4 border-blue-600 bg-blue-50 rounded-lg p-4"
                  >
                    <div className="flex justify-between items-start mb-2">
                      <h3 className="text-lg font-semibold text-gray-800">{module.title}</h3>
                      <span className={`px-2 py-1 rounded text-xs font-semibold ${enrollment.status === 'enrolled' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'}`}>
                        {enrollment.status}
                      </span>
                    </div>
                    <p className="text-sm font-mono text-blue-600 mb-2">{module.code}</p>
                    <div className="flex justify-between text-sm text-gray-600">
                      <span>{module.credits} Credits</span>
                      <span>Semester {module.semester}</span>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}

        {/* Modules Section */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold">
              {user?.role === 'student' 
                ? 'My Department Modules' 
                : user?.role === 'lecturer'
                ? 'My Modules'
                : 'Available Modules'}
            </h2>
            
            {/* Filter */}
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value)}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
            >
              <option value="all">All Types</option>
              <option value="compulsory">Compulsory</option>
              <option value="elective">Elective</option>
              <option value="faculty">Faculty</option>
            </select>
          </div>
          
          {/* For Students: Show Department Modules */}
          {user?.role === 'student' && (
            <>
              {!departmentID && (
                <div className="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
                  <p className="text-yellow-800">
                    <span className="font-semibold">⚠️ No Department Assigned:</span> Please contact admin to assign you to a department.
                  </p>
                </div>
              )}
              
              {departmentModules.length === 0 ? (
                <p className="text-gray-600">No modules available in your department yet.</p>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {departmentModules.map((module) => {
                    const isEnrolled = enrolledModuleIds.has(module.moduleID);
                    
                    return (
                      <div
                        key={module.moduleID}
                        className="border border-gray-300 rounded-lg p-4 hover:shadow-lg transition"
                      >
                        <div className="flex justify-between items-start mb-2">
                          <h3 className="text-xl font-semibold text-gray-800">{module.title}</h3>
                          {module.type && (
                            <span className={`px-2 py-1 rounded text-xs font-semibold ${getModuleTypeColor(module.type)}`}>
                              {getModuleTypeBadge(module.type)}
                            </span>
                          )}
                        </div>
                        <p className="text-sm text-gray-600 mb-2">Code: {module.code}</p>
                        <p className="text-gray-700 mb-3">{module.description}</p>
                        <div className="flex justify-between items-center text-sm text-gray-600">
                          <div>
                            <span className="mr-3">{module.credits} Credits</span>
                            <span>Semester {module.semester}</span>
                          </div>
                          
                          <button
                            onClick={() => handleEnroll(module.moduleID)}
                            disabled={isEnrolled || enrollingModule === module.moduleID}
                            className={`px-3 py-1 rounded text-sm font-semibold ${
                              isEnrolled 
                                ? 'bg-gray-200 text-gray-500 cursor-not-allowed' 
                                : 'bg-blue-600 text-white hover:bg-blue-700'
                            }`}
                          >
                            {isEnrolled ? 'Enrolled' : enrollingModule === module.moduleID ? 'Enrolling...' : 'Enroll'}
                          </button>
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </>
          )}

          {/* For Non-Students: Show All Modules */}
          {user?.role !== 'student' && (
            <>
              {modulesToDisplay.length === 0 ? (
                <div className="text-center py-8">
                  {user?.role === 'lecturer' ? (
                    <div className="p-6 bg-yellow-50 border border-yellow-200 rounded-lg">
                      <p className="text-yellow-800 text-lg">
                        <span className="font-semibold">📚 No Modules Assigned</span>
                      </p>
                      <p className="text-yellow-700 mt-2">
                        You haven't been assigned to any modules yet. Please contact the admin to assign you as a lecturer for modules.
                      </p>
                    </div>
                  ) : (
                    <p className="text-gray-600">No modules available yet.</p>
                  )}
                </div>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  {modulesToDisplay.map((module) => (
                    <div
                      key={module.moduleID}
                      className="border border-gray-300 rounded-lg p-4 hover:shadow-lg transition"
                    >
                      <div className="flex justify-between items-start mb-2">
                        <h3 className="text-xl font-semibold text-gray-800">{module.title}</h3>
                        {module.type && (
                          <span className={`px-2 py-1 rounded text-xs font-semibold ${getModuleTypeColor(module.type)}`}>
                            {getModuleTypeBadge(module.type)}
                          </span>
                        )}
                      </div>
                      <p className="text-sm text-gray-600 mb-2">Code: {module.code}</p>
                      <p className="text-gray-700 mb-3">{module.description}</p>
                      <div className="flex justify-between items-center text-sm text-gray-600">
                        <div>
                          <span className="mr-3">{module.credits} Credits</span>
                          <span>Semester {module.semester}</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </>
          )}
        </div>

        {/* Electives from Other Departments (Students Only) */}
        {user?.role === 'student' && electivesFromOtherDepartments.length > 0 && (
          <div className="bg-white rounded-lg shadow p-6 mt-8">
            <h2 className="text-2xl font-bold mb-6">Electives from Other Departments</h2>
            <p className="text-gray-600 mb-4">
              These modules are offered by other departments and available for cross-departmental enrollment.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {electivesFromOtherDepartments.map((module) => {
                const isEnrolled = enrolledModuleIds.has(module.moduleID);
                
                return (
                  <div
                    key={module.moduleID}
                    className="border-2 border-purple-200 bg-purple-50 rounded-lg p-4 hover:shadow-lg transition"
                  >
                    <div className="flex justify-between items-start mb-2">
                      <h3 className="text-xl font-semibold text-gray-800">{module.title}</h3>
                      <div className="flex flex-col gap-1">
                        {module.type && (
                          <span className={`px-2 py-1 rounded text-xs font-semibold ${getModuleTypeColor(module.type)}`}>
                            {getModuleTypeBadge(module.type)}
                          </span>
                        )}
                        <span className="px-2 py-1 rounded text-xs font-semibold bg-purple-100 text-purple-700">
                          Other Dept
                        </span>
                      </div>
                    </div>
                    <p className="text-sm text-gray-600 mb-2">Code: {module.code}</p>
                    {module.departmentID && (
                      <p className="text-xs text-purple-600 mb-2">Department: {module.departmentID}</p>
                    )}
                    <p className="text-gray-700 mb-3">{module.description}</p>
                    <div className="flex justify-between items-center text-sm text-gray-600">
                      <div>
                        <span className="mr-3">{module.credits} Credits</span>
                        <span>Semester {module.semester}</span>
                      </div>
                      
                      <button
                        onClick={() => handleEnroll(module.moduleID)}
                        disabled={isEnrolled || enrollingModule === module.moduleID}
                        className={`px-3 py-1 rounded text-sm font-semibold ${
                          isEnrolled 
                            ? 'bg-gray-200 text-gray-500 cursor-not-allowed' 
                            : 'bg-purple-600 text-white hover:bg-purple-700'
                        }`}
                      >
                        {isEnrolled ? 'Enrolled' : enrollingModule === module.moduleID ? 'Enrolling...' : 'Enroll'}
                      </button>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
