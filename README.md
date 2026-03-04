# 🎓 ModuleFlow

**AI-Assisted University Module Enrollment System**

---

## 📌 Overview

ModuleFlow is a modern, AI-assisted university module enrollment system designed to automate and streamline the course registration process. It replaces manual enrollment handling with an intelligent, workflow-driven approach that improves efficiency, transparency, and fairness for students and administrators.

Students submit enrollment requests through a secure web interface. Automated workflows validate academic rules such as prerequisites, module capacity, and eligibility. Exceptional cases are intelligently flagged using AI and routed for administrative review.

---

## 🎯 Problem Statement

Traditional university enrollment systems often suffer from:

* Manual approval processes
* Delayed responses to student requests
* Human error in prerequisite and capacity checks
* Overloaded administrative staff

ModuleFlow addresses these challenges by combining **automation, rule-based validation, and AI-assisted decision support**.

---

## 🚀 Key Features

* Secure student and admin authentication
* Role-based access control (RBAC)
* Automated prerequisite and module capacity checks
* AI-assisted evaluation of special enrollment requests
* Real-time enrollment request tracking
* Email and system notifications
* Transparent approval and rejection logging

---

## 🏗 System Architecture

ModuleFlow follows a **scalable, separation-of-concerns architecture**:

```
Student/Admin
 ↓
Next.js Frontend (Dashboard + Forms)
 ↓
Go Backend API (Gin Framework)
 ↓
Firebase Firestore Database
 ↓
Email Notifications
```

**Roles & Flow:**

1. **Student:** Submits enrollment requests through a secure dashboard
2. **Next.js Frontend:** Handles UI, forms, and client-side validation
3. **Go Backend API:** Processes requests, validates rules, manages authentication with JWT
4. **Firebase Firestore:** Stores all enrollment data, users, modules, and departments
5. **Admin:** Reviews and approves/rejects enrollment requests
6. **Notifications:** Sends email updates to students and admins

---

## 🛠 Tech Stack

| Layer                            | Technology                      | Purpose                                                          |
| -------------------------------- | ------------------------------- | ---------------------------------------------------------------- |
| **Backend**                      | Go 1.23 (Gin Framework)         | RESTful API server for enrollment operations                     |
| **Frontend**                     | Next.js 16 (React 19)           | Secure student/admin UI, form submission, dashboard              |
| **Authentication**               | JWT + NextAuth.js               | Role-based authentication (student/admin) and session management |
| **Database**                     | Firebase Firestore              | Stores students, modules, enrollments, and user data             |
| **Notifications**                | Email (SMTP)                    | Sends enrollment updates and admin alerts                        |
| **Styling & UI**                 | Tailwind CSS 4                  | Modern, responsive interface for dashboards and forms            |

> Optional Enhancements:
>
> * Vercel deployment for Next.js
> * Docker for containerized deployment
> * AI/LLM integration for enrollment recommendations

---

## ⚡ Installation & Setup

### 1️⃣ Clone the repository

```bash
git clone https://github.com/yourusername/UniEnroll.git
cd UniEnroll
```

### 2️⃣ Install Go (Backend)

```bash
# Download and install Go 1.23+
cd /tmp
curl -L https://go.dev/dl/go1.23.0.linux-amd64.tar.gz -o go.tar.gz
sudo tar -C /usr/local -xzf go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version  # Verify installation
```

### 3️⃣ Install Node.js 20+ (Frontend)

```bash
# Install Node.js 20 (required for Next.js 16)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs
node --version  # Should be v20.x or higher
```

### 4️⃣ Set up Firebase

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Create a new project (or use existing)
3. Enable **Firestore Database** in your project
4. Go to **Project Settings** → **Service Accounts** tab
5. Click **"Generate New Private Key"** button
6. Download the JSON file and save as `backend/serviceAccount.json`
7. Note your **Project ID** from the project settings

### 5️⃣ Configure Backend Environment

Create `backend/.env`:

```env
# Firebase Configuration
FIREBASE_PROJECT_ID=your-firebase-project-id

# Server Configuration
PORT=8080

# JWT Secret for authentication
JWT_SECRET=your-jwt-secret-key-change-this-in-production

# Email Configuration (optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

> 🔑 **Important:** Never expose API keys or service account files publicly.

### 6️⃣ Install Backend Dependencies

```bash
cd backend
go mod download
```

### 7️⃣ Install Frontend Dependencies

```bash
cd ../frontend
npm install
```

### 8️⃣ Run the Backend

Open a terminal and run:

```bash
cd backend
export GOOGLE_APPLICATION_CREDENTIALS="$PWD/serviceAccount.json"
go run main.go
```

✅ Backend runs on [http://localhost:8080](http://localhost:8080)

### 9️⃣ Run the Frontend

Open another terminal and run:

```bash
cd frontend
npm run dev
```

✅ Frontend runs on [http://localhost:3000](http://localhost:3000)

### 🔟 Initialize Database (First Time)

Make a POST request to seed initial data:

```bash
curl -X POST http://localhost:8080/init
```

### ✅ Test the System

1. Visit [http://localhost:3000](http://localhost:3000)
2. Register a new account or login
3. Browse modules and departments
4. Submit enrollment requests
5. Test admin features (if you have admin role)

### 🚀 Deployment (Optional)

* **Frontend:** Deploy Next.js to **Vercel**
* **Backend:** Deploy Go to **Google Cloud Run**, **Heroku**, or **Railway**
* **Database:** Already on Firebase (cloud-hosted)

---

## 🔑 API Endpoints

### Authentication
- `POST /auth/login` - User login
- `POST /auth/register` - User registration

### Protected Routes (Requires JWT)
- `GET /api/profile` - Get user profile
- `POST /api/change-password` - Change password
- `GET /api/modules` - Get all modules
- `GET /api/modules/:moduleID` - Get specific module
- `GET /api/modules/search` - Search modules by department
- `GET /api/departments` - Get all departments

### Student Routes
- `POST /api/enrollments` - Enroll in a module
- `GET /api/enrollments/:studentID` - Get student enrollments
- `PUT /api/enrollments/:enrollmentID/drop` - Drop a module

### Admin Routes (Requires admin role)
- `POST /api/admin/create` - Create admin profile
- `POST /api/admin/users` - Create user account
- `POST /api/lecturer/create` - Create lecturer profile
- `POST /api/student/create` - Create student profile
- `POST /api/modules` - Create new module
- `POST /api/departments` - Create new department

### Database Initialization
- `POST /init` - Seed database with initial data

---

## 🔐 Security & Roles

* **Students:** Request modules, view enrollments, track status
* **Lecturers:** Manage course modules
* **Admins:** Full system access, create users, approve/reject requests
* **RBAC:** Role-based access control via JWT middleware
* **Password Security:** Bcrypt hashing for all passwords
* **CORS:** Configured for cross-origin requests between frontend and backend

---

## 🌟 Why ModuleFlow is Valuable

* Solves a **real-world university problem**
* Demonstrates **modern web architecture** with Go backend and Next.js frontend
* Shows **system design, security, and full-stack development skills**
* Uses **Firebase Firestore** for scalable cloud database
* Implements **JWT authentication** and role-based access control
* Ideal for **academic projects, portfolio, or production-ready prototypes**

---

## 📂 Future Enhancements

* Add **multi-semester scheduling & waitlists**
* Integrate **analytics dashboard** for enrollment statistics
* Add **AI/LLM integration** for smart module recommendations
* Implement **real-time notifications** using WebSockets
* Add **mobile-friendly design / PWA support**
* Integrate **n8n workflow automation** for complex approval processes
* Add **payment integration** for enrollment fees

---

