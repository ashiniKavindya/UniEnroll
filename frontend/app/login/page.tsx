"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleLogin(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json();
      if (!res.ok) {
        setError(data.error || "Login failed");
        setLoading(false);
        return;
      }

      // store token (for demo, localStorage). In production prefer httpOnly cookie
      if (data.token) localStorage.setItem("token", data.token);

      // redirect based on role
      if (data.role === "admin") router.push("/admin");
      else router.push("/dashboard");
    } catch (err) {
      setError("Network error");
    } finally {
      setLoading(false);
    }
  }

  async function quickLogin(emailVal: string, passwordVal: string, redirect: string) {
    setEmail(emailVal);
    setPassword(passwordVal);
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${API_URL}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: emailVal, password: passwordVal }),
      });
      const data = await res.json();
      if (!res.ok) {
        setError(data.error || "Login failed");
        setLoading(false);
        return;
      }
      if (data.token) localStorage.setItem("token", data.token);
      router.push(redirect);
    } catch (e) {
      setError("Network error");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="w-96 p-6 border rounded-lg shadow">
        <h1 className="text-2xl font-bold mb-4 text-center">ModuleFlow Login</h1>

        <form onSubmit={handleLogin} className="space-y-3">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full p-2 border rounded"
            required
          />

          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full p-2 border rounded"
            required
          />

          {error && <div className="text-red-600">{error}</div>}

          <button
            type="submit"
            className="w-full mb-3 bg-blue-600 text-white py-2 rounded"
            disabled={loading}
          >
            {loading ? "Logging in..." : "Login"}
          </button>
        </form>

        <div className="mt-4">
          <button
            onClick={() => quickLogin("student@uni.edu", "student123", "/dashboard")}
            className="w-full mb-3 bg-blue-500 text-white py-2 rounded"
          >
            Login as Student
          </button>

          <button
            onClick={() => quickLogin("admin@uni.edu", "admin123", "/admin")}
            className="w-full bg-gray-800 text-white py-2 rounded"
          >
            Login as Admin
          </button>
        </div>
      </div>
    </div>
  );
}
