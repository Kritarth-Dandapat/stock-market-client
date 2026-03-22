"use client";

import { useQuery } from "@tanstack/react-query";
import { fundsApi, formatINR } from "@/lib/api-client";
import FundCard from "@/components/fund-card/FundCard";
import { useState } from "react";
import Link from "next/link";

export default function FundsPage() {
  const [search, setSearch] = useState("");
  const { data: funds, isLoading } = useQuery({
    queryKey: ["funds"],
    queryFn: fundsApi.list,
  });

  const filtered = funds?.filter((f) =>
    f.scheme_name.toLowerCase().includes(search.toLowerCase()) ||
    f.scheme_code.toLowerCase().includes(search.toLowerCase())
  ) ?? [];

  return (
    <div className="min-h-screen bg-slate-900 text-white">
      <nav className="border-b border-slate-700 px-6 py-4 flex items-center justify-between">
        <h1 className="text-xl font-bold">BSE MF Dashboard</h1>
        <div className="flex gap-4 text-sm text-slate-400">
          <Link href="/portfolio" className="hover:text-white">Portfolio</Link>
          <Link href="/funds" className="text-white font-medium">Funds</Link>
          <Link href="/sip" className="hover:text-white">SIP</Link>
          <Link href="/transactions" className="hover:text-white">Transactions</Link>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 mb-6">
          <h2 className="text-2xl font-semibold">Browse Funds</h2>
          <input
            type="search"
            placeholder="Search by fund name or scheme code..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="w-full sm:w-80 px-4 py-2 rounded-lg bg-slate-700 border border-slate-600 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
          />
        </div>

        {isLoading && (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 animate-pulse">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="h-40 bg-slate-700 rounded-xl" />
            ))}
          </div>
        )}

        {!isLoading && filtered.length === 0 && (
          <div className="text-center py-16 text-slate-400">
            {search ? "No funds match your search." : "No funds available."}
          </div>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map((fund) => (
            <FundCard key={fund.scheme_code} fund={fund} />
          ))}
        </div>
      </main>
    </div>
  );
}
