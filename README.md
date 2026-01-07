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
Student
 ↓
Next.js Frontend (Enrollment Form + Dashboard)
 ↓
Next.js API Routes (Validation & Security)
 ↓
n8n Workflow (Automation Engine)
 ↓
AI Decision + Rule Checks
 ↓
Database + Notifications
```

**Roles & Flow:**

1. **Student:** Submits enrollment requests through a secure dashboard
2. **Next.js API Routes:** Validate requests and forward to n8n webhook
3. **n8n Workflow:** Checks prerequisites, capacity, and automates notifications
4. **AI Service:** Evaluates special cases and provides recommendations
5. **Database:** Stores all enrollment data and tracks status
6. **Notifications:** Sends emails or alerts to students and admins

---

## 🛠 Tech Stack

| Layer                            | Technology                      | Purpose                                                          |
| -------------------------------- | ------------------------------- | ---------------------------------------------------------------- |
| **Frontend**                     | Next.js (React)                 | Secure student/admin UI, form submission, dashboard              |
| **Authentication**               | NextAuth.js                     | Role-based authentication (student/admin) and session management |
| **Automation / Workflow Engine** | n8n                             | Orchestrates enrollment workflows, AI evaluation, notifications  |
| **AI / Decision Support**        | OpenAI / LLM                    | Analyzes special cases, flags exceptions, recommends approvals   |
| **Database**                     | PostgreSQL / Supabase           | Stores students, modules, and enrollment requests                |
| **Notifications**                | Email (SMTP) / Telegram / Slack | Sends enrollment updates and admin alerts                        |
| **Styling & UI**                 | Tailwind CSS / Material UI      | Modern, responsive interface for dashboards and forms            |

> Optional Enhancements:
>
> * Vercel deployment for Next.js
> * Supabase Edge Functions for advanced backend logic
> * Docker for containerized deployment

---

## ⚡ Installation & Setup

### 1️⃣ Clone the repository

```bash
git clone https://github.com/yourusername/moduleflow.git
cd moduleflow
```

### 2️⃣ Install frontend dependencies

```bash
npm install
```

### 3️⃣ Set up environment variables

Create a `.env` file in the root directory:

```env
# NextAuth.js
NEXTAUTH_SECRET=your_random_secret_key
NEXTAUTH_URL=http://localhost:3000

# Database
DATABASE_URL=postgres://username:password@localhost:5432/moduleflow_db

# n8n Webhook
N8N_WEBHOOK_URL=http://localhost:5678/webhook/enrollment
N8N_API_KEY=your_n8n_api_key

# OpenAI (AI Evaluation)
OPENAI_API_KEY=your_openai_api_key
```

> 🔑 **Important:** Never expose API keys publicly.

### 4️⃣ Run the Next.js frontend

```bash
npm run dev
```

* Visit [http://localhost:3000](http://localhost:3000)
* Login as **student** or **admin** (create accounts in DB)

### 5️⃣ Run n8n

```bash
n8n
```

* Import the **ModuleFlow workflow** JSON
* Ensure webhooks and AI nodes are configured

### 6️⃣ Initialize the Database

* Use **Supabase / PostgreSQL**
* Run the `db/schema.sql` script to create tables:

  * `students`
  * `modules`
  * `enrollment_requests`

### 7️⃣ Test the System

1. Log in as a student
2. Submit an enrollment request
3. Check:

   * Dashboard for status
   * Admin interface for approval/rejection
   * Email or Telegram notifications

### 8️⃣ Deployment (Optional)

* **Frontend:** Deploy Next.js to **Vercel**
* **n8n:** Deploy to **n8n.cloud** or self-host with Docker
* **Database:** Use **Supabase** or managed PostgreSQL

---

## 🧠 How AI Works in ModuleFlow

AI is used **only for intelligent recommendations**, not final decisions. It:

* Analyzes free-text justification for special enrollment requests
* Flags exceptional cases
* Suggests approval or rejection for admin review

**Example:**

> “I failed this module due to medical reasons and request re-enrollment.”
> AI flags it as **special case → admin review → recommendation generated**

---

## 🔐 Security & Roles

* **Students:** Request modules, track status
* **Admins / Lecturers:** Review, approve, or reject requests
* **System (n8n):** Workflow automation only
* **RBAC:** Role-based access control ensures secure operations

> n8n webhook is **private**, only Next.js API routes can call it.

---

## 🌟 Why ModuleFlow is Valuable

* Solves a **real-world university problem**
* Demonstrates **automation + AI integration**
* Shows **system design, security, and modern frontend skills**
* Ideal for **academic projects, portfolio, or production-ready prototypes**

---

## 📂 Future Enhancements

* Add **multi-semester scheduling & waitlists**
* Integrate **analytics dashboard** for enrollment statistics
* Extend AI to **suggest module combinations** for students based on degree and performance
* Add **mobile-friendly design / PWA support**

---

