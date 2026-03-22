"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { authApi } from "@/lib/api-client";

const schema = z.object({
  email: z.string().email("Enter a valid email address"),
  phone: z.string().regex(/^\+91[6-9]\d{9}$/, "Enter a valid Indian mobile number (+91XXXXXXXXXX)"),
});

type RegisterForm = z.infer<typeof schema>;

export default function RegisterPage() {
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState("");

  const { register, handleSubmit, formState: { errors } } = useForm<RegisterForm>({
    resolver: zodResolver(schema),
  });

  async function onSubmit(data: RegisterForm) {
    setLoading(true);
    setError("");
    try {
      await authApi.register(data.email, data.phone);
      setSuccess(true);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Registration failed");
    } finally {
      setLoading(false);
    }
  }

  if (success) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-slate-900 p-4">
        <div className="w-full max-w-md bg-slate-800 rounded-2xl p-8 shadow-xl text-center">
          <div className="text-green-400 text-5xl mb-4">✓</div>
          <h2 className="text-xl font-bold text-white mb-2">Account Created!</h2>
          <p className="text-slate-400 text-sm mb-6">
            Your account has been registered. Please complete your KYC to start investing.
          </p>
          <Link
            href="/kyc"
            className="block w-full py-3 bg-blue-600 hover:bg-blue-500 rounded-lg font-semibold text-white transition-colors"
          >
            Complete KYC
          </Link>
          <Link href="/login" className="block mt-3 text-sm text-slate-400 hover:text-white">
            Login instead →
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-900 p-4">
      <div className="w-full max-w-md bg-slate-800 rounded-2xl p-8 shadow-xl">
        <h1 className="text-2xl font-bold text-white mb-2">Create Account</h1>
        <p className="text-slate-400 mb-6 text-sm">
          Join BSE Mutual Fund Dashboard to start investing.
        </p>

        {error && (
          <div className="mb-4 p-3 rounded-lg bg-red-900/50 border border-red-700 text-red-300 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm text-slate-300 mb-1">Email Address</label>
            <input
              {...register("email")}
              type="email"
              placeholder="you@example.com"
              className="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {errors.email && (
              <p className="mt-1 text-xs text-red-400">{errors.email.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm text-slate-300 mb-1">Mobile Number</label>
            <input
              {...register("phone")}
              type="tel"
              placeholder="+919876543210"
              className="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {errors.phone && (
              <p className="mt-1 text-xs text-red-400">{errors.phone.message}</p>
            )}
          </div>

          <p className="text-xs text-slate-500">
            By registering you agree to comply with SEBI regulations and our terms of service.
          </p>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 rounded-lg font-semibold text-white transition-colors"
          >
            {loading ? "Creating account..." : "Create Account"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-slate-400">
          Already have an account?{" "}
          <Link href="/login" className="text-blue-400 hover:text-blue-300">
            Login
          </Link>
        </p>
      </div>
    </div>
  );
}
