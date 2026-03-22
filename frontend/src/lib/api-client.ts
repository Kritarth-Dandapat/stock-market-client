import { z } from "zod";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "/api";

/**
 * A type-safe fetch wrapper that routes all requests through the Next.js API
 * routes, keeping the backend URL private from the browser.
 */
async function apiFetch<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE}${path}`;

  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    credentials: "include", // send HTTP-only session cookie
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: "Unknown error" }));
    throw new ApiError(response.status, error?.error ?? "Request failed");
  }

  return response.json() as Promise<T>;
}

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string
  ) {
    super(message);
    this.name = "ApiError";
  }
}

// ── Zod schemas ──────────────────────────────────────────────────────────────

export const FundSchemeSchema = z.object({
  scheme_code: z.string(),
  isin: z.string(),
  scheme_name: z.string(),
  amc_code: z.string(),
  category: z.string(),
  nav: z.number(),
  nav_date: z.string(),
  min_purchase: z.number(),
  min_sip: z.number(),
});
export type FundScheme = z.infer<typeof FundSchemeSchema>;

export const HoldingSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  scheme_code: z.string(),
  scheme_name: z.string(),
  isin: z.string(),
  units: z.number(),
  avg_nav: z.number(),
  current_nav: z.number(),
  invested_amount: z.number(),
  current_amount: z.number(),
  xirr: z.number().optional(),
  updated_at: z.string(),
});
export type Holding = z.infer<typeof HoldingSchema>;

export const PortfolioSummarySchema = z.object({
  total_invested: z.number(),
  current_value: z.number(),
  absolute_return: z.number(),
  xirr: z.number(),
  holdings: z.array(HoldingSchema),
});
export type PortfolioSummary = z.infer<typeof PortfolioSummarySchema>;

export const OrderSchema = z.object({
  id: z.string(),
  user_id: z.string(),
  scheme_code: z.string(),
  scheme_name: z.string(),
  order_type: z.enum(["lumpsum", "sip", "redeem"]),
  amount: z.number(),
  units: z.number().optional(),
  nav: z.number().optional(),
  status: z.enum(["pending", "submitted", "processing", "completed", "failed", "cancelled"]),
  bse_order_id: z.string().optional(),
  payment_id: z.string().optional(),
  failure_reason: z.string().optional(),
  created_at: z.string(),
  updated_at: z.string(),
});
export type Order = z.infer<typeof OrderSchema>;

// ── Auth ─────────────────────────────────────────────────────────────────────

export const authApi = {
  requestOTP: (phone: string) =>
    apiFetch<{ message: string }>("/auth/login/otp/request", {
      method: "POST",
      body: JSON.stringify({ phone }),
    }),

  verifyOTP: (phone: string, otp: string) =>
    apiFetch<{ access_token: string; refresh_token: string; expires_in: number }>(
      "/auth/login/otp/verify",
      {
        method: "POST",
        body: JSON.stringify({ phone, otp }),
      }
    ),

  register: (email: string, phone: string) =>
    apiFetch<{ user_id: string; message: string }>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, phone }),
    }),
};

// ── Funds ─────────────────────────────────────────────────────────────────────

export const fundsApi = {
  list: () => apiFetch<FundScheme[]>("/funds"),
  get: (schemeCode: string) =>
    apiFetch<FundScheme>(`/funds/${encodeURIComponent(schemeCode)}`),
};

// ── Portfolio ─────────────────────────────────────────────────────────────────

export const portfolioApi = {
  getSummary: () => apiFetch<PortfolioSummary>("/portfolio"),
  getHoldings: () => apiFetch<Holding[]>("/portfolio/holdings"),
};

// ── Orders ────────────────────────────────────────────────────────────────────

export const ordersApi = {
  place: (payload: { scheme_code: string; order_type: string; amount: number }) =>
    apiFetch<Order>("/orders", {
      method: "POST",
      body: JSON.stringify(payload),
    }),
  list: () => apiFetch<Order[]>("/orders"),
  get: (id: string) => apiFetch<Order>(`/orders/${encodeURIComponent(id)}`),
};

// ── KYC ──────────────────────────────────────────────────────────────────────

export const kycApi = {
  initiate: (pan: string, fullName: string, dateOfBirth: string) =>
    apiFetch<{ kyc_id: string; pan_masked: string; status: string }>("/kyc/initiate", {
      method: "POST",
      body: JSON.stringify({ pan, full_name: fullName, date_of_birth: dateOfBirth }),
    }),
  getStatus: () =>
    apiFetch<{ user_id: string; status: string }>("/kyc/status"),
};

// ── Formatting helpers ────────────────────────────────────────────────────────

/**
 * Format a number as Indian Rupee using the correct en-IN locale.
 */
export function formatINR(amount: number): string {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    maximumFractionDigits: 2,
  }).format(amount);
}

/**
 * Format a percentage with 2 decimal places and a + prefix for positives.
 */
export function formatPercent(value: number): string {
  const formatted = value.toFixed(2) + "%";
  return value >= 0 ? "+" + formatted : formatted;
}
