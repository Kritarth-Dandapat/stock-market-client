"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { authApi } from "@/lib/api-client";

const phoneSchema = z.object({
  phone: z.string().regex(/^\+91[6-9]\d{9}$/, "Enter a valid Indian mobile number (+91XXXXXXXXXX)"),
});
const otpSchema = z.object({
  otp: z.string().length(6, "OTP must be 6 digits").regex(/^\d+$/, "OTP must be numeric"),
});

type PhoneForm = z.infer<typeof phoneSchema>;
type OTPForm = z.infer<typeof otpSchema>;

export default function LoginPage() {
  const [step, setStep] = useState<"phone" | "otp">("phone");
  const [phone, setPhone] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const phoneForm = useForm<PhoneForm>({ resolver: zodResolver(phoneSchema) });
  const otpForm = useForm<OTPForm>({ resolver: zodResolver(otpSchema) });

  async function onRequestOTP(data: PhoneForm) {
    setLoading(true);
    setError("");
    try {
      await authApi.requestOTP(data.phone);
      setPhone(data.phone);
      setStep("otp");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to send OTP");
    } finally {
      setLoading(false);
    }
  }

  async function onVerifyOTP(data: OTPForm) {
    setLoading(true);
    setError("");
    try {
      await authApi.verifyOTP(phone, data.otp);
      window.location.href = "/portfolio";
    } catch (e) {
      setError(e instanceof Error ? e.message : "OTP verification failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-900 p-4">
      <div className="w-full max-w-md bg-slate-800 rounded-2xl p-8 shadow-xl">
        <h1 className="text-2xl font-bold text-white mb-2">Login</h1>
        <p className="text-slate-400 mb-6 text-sm">
          {step === "phone"
            ? "Enter your registered mobile number"
            : `Enter the OTP sent to ${phone}`}
        </p>

        {error && (
          <div className="mb-4 p-3 rounded-lg bg-red-900/50 border border-red-700 text-red-300 text-sm">
            {error}
          </div>
        )}

        {step === "phone" ? (
          <form onSubmit={phoneForm.handleSubmit(onRequestOTP)} className="space-y-4">
            <div>
              <label className="block text-sm text-slate-300 mb-1">Mobile Number</label>
              <input
                {...phoneForm.register("phone")}
                type="tel"
                placeholder="+919876543210"
                className="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              {phoneForm.formState.errors.phone && (
                <p className="mt-1 text-xs text-red-400">
                  {phoneForm.formState.errors.phone.message}
                </p>
              )}
            </div>
            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 rounded-lg font-semibold text-white transition-colors"
            >
              {loading ? "Sending OTP..." : "Send OTP"}
            </button>
          </form>
        ) : (
          <form onSubmit={otpForm.handleSubmit(onVerifyOTP)} className="space-y-4">
            <div>
              <label className="block text-sm text-slate-300 mb-1">One-Time Password</label>
              <input
                {...otpForm.register("otp")}
                type="text"
                inputMode="numeric"
                maxLength={6}
                placeholder="123456"
                className="w-full px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 tracking-[0.5em] text-center text-xl"
              />
              {otpForm.formState.errors.otp && (
                <p className="mt-1 text-xs text-red-400">
                  {otpForm.formState.errors.otp.message}
                </p>
              )}
            </div>
            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 rounded-lg font-semibold text-white transition-colors"
            >
              {loading ? "Verifying..." : "Verify OTP"}
            </button>
            <button
              type="button"
              onClick={() => setStep("phone")}
              className="w-full py-2 text-slate-400 hover:text-white text-sm transition-colors"
            >
              ← Change number
            </button>
          </form>
        )}

        <p className="mt-6 text-center text-sm text-slate-400">
          Don&apos;t have an account?{" "}
          <Link href="/register" className="text-blue-400 hover:text-blue-300">
            Register
          </Link>
        </p>
      </div>
    </div>
  );
}
