"use client";

import { useQuery } from "@tanstack/react-query";
import { portfolioApi, formatINR, formatPercent } from "@/lib/api-client";
import PortfolioSummaryCard from "@/components/portfolio-summary/PortfolioSummaryCard";
import Link from "next/link";

export default function PortfolioPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ["portfolio"],
    queryFn: portfolioApi.getSummary,
  });

  return (
    <div className="min-h-screen bg-slate-900 text-white">
      <nav className="border-b border-slate-700 px-6 py-4 flex items-center justify-between">
        <h1 className="text-xl font-bold">BSE MF Dashboard</h1>
        <div className="flex gap-4 text-sm text-slate-400">
          <Link href="/portfolio" className="text-white font-medium">Portfolio</Link>
          <Link href="/funds" className="hover:text-white">Funds</Link>
          <Link href="/sip" className="hover:text-white">SIP</Link>
          <Link href="/transactions" className="hover:text-white">Transactions</Link>
          <Link href="/statements" className="hover:text-white">Statements</Link>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <h2 className="text-2xl font-semibold mb-6">My Portfolio</h2>

        {isLoading && (
          <div className="grid grid-cols-1 sm:grid-cols-4 gap-4 animate-pulse">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-24 bg-slate-700 rounded-xl" />
            ))}
          </div>
        )}

        {error && (
          <div className="p-4 bg-red-900/50 border border-red-700 rounded-xl text-red-300">
            Failed to load portfolio. Please try again.
          </div>
        )}

        {data && (
          <>
            <PortfolioSummaryCard summary={data} />

            <div className="mt-8">
              <h3 className="text-lg font-semibold mb-4">Holdings</h3>
              {data.holdings.length === 0 ? (
                <div className="text-center py-16 bg-slate-800 rounded-xl">
                  <p className="text-slate-400 mb-4">No holdings yet.</p>
                  <Link
                    href="/funds"
                    className="px-6 py-2 bg-blue-600 hover:bg-blue-500 rounded-lg font-medium transition-colors"
                  >
                    Browse Funds
                  </Link>
                </div>
              ) : (
                <div className="space-y-3">
                  {data.holdings.map((holding) => (
                    <div
                      key={holding.id}
                      className="bg-slate-800 rounded-xl p-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4"
                    >
                      <div>
                        <p className="font-medium">{holding.scheme_name}</p>
                        <p className="text-sm text-slate-400">{holding.units.toFixed(4)} units @ {formatINR(holding.current_nav)}</p>
                      </div>
                      <div className="sm:text-right">
                        <p className="font-semibold">{formatINR(holding.current_amount)}</p>
                        <p className={`text-sm ${holding.current_amount >= holding.invested_amount ? "text-green-400" : "text-red-400"}`}>
                          {formatPercent(((holding.current_amount - holding.invested_amount) / holding.invested_amount) * 100)}
                        </p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </>
        )}
      </main>
    </div>
  );
}
