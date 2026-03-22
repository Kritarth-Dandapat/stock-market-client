"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { kycApi } from "@/lib/api-client";
import Link from "next/link";

const schema = z.object({
  pan: z.string().regex(/^[A-Z]{5}[0-9]{4}[A-Z]$/, "Enter a valid PAN (e.g. ABCDE1234F)"),
  fullName: z.string().min(3, "Full name must be at least 3 characters"),
  dateOfBirth: z.string().regex(/^\d{4}-\d{2}-\d{2}$/, "Enter date as YYYY-MM-DD"),
});

type KYCForm = z.infer<typeof schema>;

export default function KYCPage() {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<{ kyc_id: string; pan_masked: string; status: string } | null>(null);
  const [error, setError] = useState("");

  const { register, handleSubmit, formState: { errors } } = useForm<KYCForm>({
    resolver: zodResolver(schema),
  });

  async function onSubmit(data: KYCForm) {
    setLoading(true);
    setError("");
    try {
      const res = await kycApi.initiate(data.pan, data.fullName, data.dateOfBirth);
      setResult(res);
    } catch (e) {
      setError(e instanceof Error ? e.message : "KYC initiation failed");
    } finally {
      setLoading(false);
    }
  }

  if (result) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-slate-900 p-4">
        <div className="w-full max-w-md bg-slate-800 rounded-2xl p-8 shadow-xl text-center">
          <div className="text-blue-400 text-5xl mb-4">🔍</div>
          <h2 className="text-xl font-bold text-white mb-2">KYC Under Review</h2>
          <p className="text-slate-400 text-sm mb-4">
            Your KYC has been submitted for verification via CKYC registry.
          </p>
          <div className="bg-slate-700 rounded-lg p-4 mb-6 text-left">
            <div className="flex justify-between text-sm">
              <span className="text-slate-400">PAN (masked)</span>
              <span className="text-white font-mono">{result.pan_masked}</span>
            </div>
            <div className="flex justify-between text-sm mt-2">
              <span className="text-slate-400">Status</span>
              <span className="text-yellow-400 capitalize">{result.status}</span>
            </div>
            <div className="flex justify-between text-sm mt-2">
              <span className="text-slate-400">KYC ID</span>
              <span className="text-slate-300 font-mono text-xs">{result.kyc_id}</span>
            </div>
          </div>
          <Link
            href="/portfolio"
            className="block w-full py-3 bg-blue-600 hover:bg-blue-500 rounded-lg font-semibold text-white transition-colors"
          >
            Go to Dashboard →
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-900 p-4">
      <div className="w-full max-w-md bg-slate-800 rounded-2xl p-8 shadow-xl">
        <h1 className="text-2xl font-bold text-white mb-2">KYC Verification</h1>
        <p className="text-slate-400 mb-6 text-sm">
          As per SEBI regulations, all investors must complete KYC before investing.
          Your PAN will be verified via the CKYC registry.
        </p>

        {error && (
          <div className="mb-4 p-3 rounded-lg bg-red-900/50 border border-red-700 text-red-300 text-sm">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm text-slate-300 mb-1">PAN Number</label>
            <input
              {...register("pan")}
              type="text"
              placeholder="ABCDE1234F"
              maxLength={10}
              className="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 uppercase"
              style={{ textTransform: "uppercase" }}
            />
            {errors.pan && (
              <p className="mt-1 text-xs text-red-400">{errors.pan.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm text-slate-300 mb-1">Full Name (as per PAN)</label>
            <input
              {...register("fullName")}
              type="text"
              placeholder="Rajesh Kumar Sharma"
              className="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {errors.fullName && (
              <p className="mt-1 text-xs text-red-400">{errors.fullName.message}</p>
            )}
          </div>

          <div>
            <label className="block text-sm text-slate-300 mb-1">Date of Birth</label>
            <input
              {...register("dateOfBirth")}
              type="date"
              className="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            {errors.dateOfBirth && (
              <p className="mt-1 text-xs text-red-400">{errors.dateOfBirth.message}</p>
            )}
          </div>

          <p className="text-xs text-slate-500">
            Your PAN is encrypted using AES-256 and never stored in plain text. 
            We use CKYC (CVL KRA) for verification.
          </p>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 rounded-lg font-semibold text-white transition-colors"
          >
            {loading ? "Verifying..." : "Submit KYC"}
          </button>
        </form>
      </div>
    </div>
  );
}
