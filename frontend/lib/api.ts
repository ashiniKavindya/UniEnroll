// lib/api.ts
const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
  role: 'student' | 'lecturer' | 'admin';
}

export interface AuthResponse {
  token: string;
  role: string;
  forcePasswordChange?: boolean;
}

export interface UserProfile {
  id: string;
  name: string;
  email: string;
  role: string;
  departmentID?: string;
  yearOfStudy?: number;
  lecturerID?: string;
}

export const authAPI = {
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    const response = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(credentials),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Login failed');
    }
    return response.json();
  },

  register: async (data: RegisterRequest): Promise<{ message: string }> => {
    const response = await fetch(`${API_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Registration failed');
    }
    return response.json();
  },

  getProfile: async (token: string): Promise<UserProfile> => {
    const response = await fetch(`${API_URL}/api/profile`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch profile');
    return response.json();
  },
  changePassword: async (token: string, currentPassword: string, newPassword: string) => {
    const response = await fetch(`${API_URL}/api/change-password`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ currentPassword, newPassword }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to change password');
    }
    return response.json();
  },
};

export interface CreateUserRequest {
  name: string;
  email: string;
  role: 'student' | 'lecturer';
  departmentID?: string;
  yearOfStudy?: number;
  studentNumber?: string;
  specialization?: string;
  officeLocation?: string;
}

export interface CreateUserResponse {
  message: string;
  userID: string;
  email: string;
  tempPassword?: string;
}

export interface ImportStudentsResponse {
  message: string;
  processed: number;
  created: number;
  skipped: number;
  failures: Array<{ row: number; email?: string; error: string }>;
  credentials: Array<{ name: string; email: string; userID: string; tempPassword: string }>;
}

export const adminAPI = {
  createUser: async (token: string, data: CreateUserRequest): Promise<CreateUserResponse> => {
    const response = await fetch(`${API_URL}/api/admin/users`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create user');
    }
    return response.json();
  },

  importStudentsCSV: async (token: string, file: File): Promise<ImportStudentsResponse> => {
    const formData = new FormData();
    formData.append('file', file);

    const response = await fetch(`${API_URL}/api/admin/students/import`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    });

    const data = await response.json();
    if (!response.ok && response.status !== 207) {
      throw new Error(data.error || data.message || 'Failed to import students');
    }

    return data;
  },
};

export const moduleAPI = {
  getAll: async (token: string) => {
    const response = await fetch(`${API_URL}/api/modules`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch modules');
    return response.json();
  },

  create: async (token: string, data: any) => {
    const response = await fetch(`${API_URL}/api/modules`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create module');
    }
    return response.json();
  },

  getById: async (token: string, moduleID: string) => {
    const response = await fetch(`${API_URL}/api/modules/${moduleID}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch module');
    return response.json();
  },

  getByDepartment: async (token: string, departmentID?: string, type?: string) => {
    const params = new URLSearchParams();
    if (departmentID) params.append('departmentID', departmentID);
    if (type) params.append('type', type);
    
    const response = await fetch(`${API_URL}/api/modules/search?${params}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch modules');
    return response.json();
  },
};

export const departmentAPI = {
  getAll: async (token: string) => {
    const response = await fetch(`${API_URL}/api/departments`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch departments');
    return response.json();
  },

  create: async (token: string, data: any) => {
    const response = await fetch(`${API_URL}/api/departments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create department');
    }
    return response.json();
  },
};

export const enrollmentAPI = {
  enroll: async (token: string, studentID: string, moduleID: string) => {
    const response = await fetch(`${API_URL}/api/enrollments`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ studentID, moduleID }),
    });
    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to enroll');
    }
    return response.json();
  },

  getStudentEnrollments: async (token: string, studentID: string) => {
    const response = await fetch(`${API_URL}/api/enrollments/${studentID}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to fetch enrollments');
    return response.json();
  },

  drop: async (token: string, enrollmentID: string) => {
    const response = await fetch(`${API_URL}/api/enrollments/${enrollmentID}/drop`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${token}` },
    });
    if (!response.ok) throw new Error('Failed to drop module');
    return response.json();
  },
};
