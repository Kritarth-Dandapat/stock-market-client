"use client";

import { useState } from "react";
import { formatINR } from "@/lib/api-client";
import Link from "next/link";

export default function SIPPage() {
  return (
    <div className="min-h-screen bg-slate-900 text-white">
      <nav className="border-b border-slate-700 px-6 py-4 flex items-center justify-between">
        <h1 className="text-xl font-bold">BSE MF Dashboard</h1>
        <div className="flex gap-4 text-sm text-slate-400">
          <Link href="/portfolio" className="hover:text-white">Portfolio</Link>
          <Link href="/funds" className="hover:text-white">Funds</Link>
          <Link href="/sip" className="text-white font-medium">SIP</Link>
          <Link href="/transactions" className="hover:text-white">Transactions</Link>
          <Link href="/statements" className="hover:text-white">Statements</Link>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <h2 className="text-2xl font-semibold mb-6">SIP Manager</h2>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div>
            <h3 className="text-lg font-medium mb-4">Active SIPs</h3>
            <div className="text-center py-12 bg-slate-800 rounded-xl">
              <p className="text-slate-400 mb-4">No active SIPs.</p>
              <Link
                href="/funds"
                className="px-6 py-2 bg-blue-600 hover:bg-blue-500 rounded-lg font-medium transition-colors"
              >
                Set Up a SIP
              </Link>
            </div>
          </div>

          <div>
            <h3 className="text-lg font-medium mb-4">SIP Calculator</h3>
            <SIPCalculatorWidget />
          </div>
        </div>
      </main>
    </div>
  );
}

function SIPCalculatorWidget() {
  const [monthly, setMonthly] = useState(5000);
  const [years, setYears] = useState(10);
  const [rate, setRate] = useState(12);

  const months = years * 12;
  const r = rate / 12 / 100;
  const futureValue = monthly * ((Math.pow(1 + r, months) - 1) / r) * (1 + r);
  const totalInvested = monthly * months;
  const wealthGained = futureValue - totalInvested;

  return (
    <div className="bg-slate-800 rounded-xl p-6 space-y-5">
      <div>
        <label className="block text-sm text-slate-300 mb-2">
          Monthly Investment: {formatINR(monthly)}
        </label>
        <input
          type="range"
          min={500}
          max={100000}
          step={500}
          value={monthly}
          onChange={(e) => setMonthly(Number(e.target.value))}
          className="w-full accent-blue-500"
        />
        <div className="flex justify-between text-xs text-slate-500 mt-1">
          <span>₹500</span><span>₹1,00,000</span>
        </div>
      </div>

      <div>
        <label className="block text-sm text-slate-300 mb-2">Duration: {years} years</label>
        <input
          type="range"
          min={1}
          max={30}
          value={years}
          onChange={(e) => setYears(Number(e.target.value))}
          className="w-full accent-blue-500"
        />
        <div className="flex justify-between text-xs text-slate-500 mt-1">
          <span>1 yr</span><span>30 yrs</span>
        </div>
      </div>

      <div>
        <label className="block text-sm text-slate-300 mb-2">Expected Return: {rate}% p.a.</label>
        <input
          type="range"
          min={6}
          max={30}
          step={0.5}
          value={rate}
          onChange={(e) => setRate(Number(e.target.value))}
          className="w-full accent-blue-500"
        />
        <div className="flex justify-between text-xs text-slate-500 mt-1">
          <span>6%</span><span>30%</span>
        </div>
      </div>

      <div className="border-t border-slate-700 pt-4 grid grid-cols-3 gap-4">
        <div>
          <p className="text-xs text-slate-400">Invested</p>
          <p className="font-semibold text-white">{formatINR(totalInvested)}</p>
        </div>
        <div>
          <p className="text-xs text-slate-400">Gains</p>
          <p className="font-semibold text-green-400">{formatINR(wealthGained)}</p>
        </div>
        <div>
          <p className="text-xs text-slate-400">Future Value</p>
          <p className="font-semibold text-blue-400">{formatINR(futureValue)}</p>
        </div>
      </div>
    </div>
  );
}
