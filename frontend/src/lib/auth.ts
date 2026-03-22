/**
 * auth.ts — Session management utilities.
 *
 * Tokens are stored in HTTP-only, Secure, SameSite=Strict cookies managed by
 * the Next.js API routes (/api/auth/*). This file provides helpers for reading
 * the session state on the client side (without exposing raw tokens).
 */

export interface SessionUser {
  userId: string;
  email: string;
  isKYCVerified: boolean;
}

/**
 * getSession reads the current session from the /api/auth/session endpoint.
 * Returns null if the user is not authenticated.
 */
export async function getSession(): Promise<SessionUser | null> {
  try {
    const res = await fetch("/api/auth/session", { credentials: "include" });
    if (!res.ok) return null;
    return res.json() as Promise<SessionUser>;
  } catch {
    return null;
  }
}

/**
 * logout calls the /api/auth/logout endpoint which clears the session cookie.
 */
export async function logout(): Promise<void> {
  await fetch("/api/auth/logout", { method: "POST", credentials: "include" });
}
