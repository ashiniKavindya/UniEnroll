"use client";

import { signIn } from "next-auth/react";

export default function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="w-96 p-6 border rounded-lg shadow">
        <h1 className="text-2xl font-bold mb-4 text-center">ModuleFlow Login</h1>

        <button
          onClick={() =>
            signIn("credentials", {
              email: "student@uni.edu",
              password: "student123",
              callbackUrl: "/dashboard",
            })
          }
          className="w-full mb-3 bg-blue-600 text-white py-2 rounded"
        >
          Login as Student
        </button>

        <button
          onClick={() =>
            signIn("credentials", {
              email: "admin@uni.edu",
              password: "admin123",
              callbackUrl: "/admin",
            })
          }
          className="w-full bg-gray-800 text-white py-2 rounded"
        >
          Login as Admin
        </button>
      </div>
    </div>
  );
}
