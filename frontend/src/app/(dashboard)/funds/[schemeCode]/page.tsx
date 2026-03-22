"use client";

import { use } from "react";
import { useQuery } from "@tanstack/react-query";
import { fundsApi, ordersApi, formatINR } from "@/lib/api-client";
import NAVChart from "@/components/charts/NAVChart";
import { useState } from "react";

export default function FundDetailPage({ params }: { params: Promise<{ schemeCode: string }> }) {
  const { schemeCode } = use(params);
  const [amount, setAmount] = useState("");
  const [placing, setPlacing] = useState(false);
  const [orderResult, setOrderResult] = useState<string>("");

  const { data: fund, isLoading } = useQuery({
    queryKey: ["fund", schemeCode],
    queryFn: () => fundsApi.get(decodeURIComponent(schemeCode)),
  });

  async function handleBuy() {
    const amt = parseFloat(amount);
    if (!fund || isNaN(amt) || amt < (fund.min_purchase ?? 500)) return;
    setPlacing(true);
    try {
      const order = await ordersApi.place({
        scheme_code: fund.scheme_code,
        order_type: "lumpsum",
        amount: amt,
      });
      setOrderResult(`Order placed! ID: ${order.id} — Status: ${order.status}`);
    } catch (e) {
      setOrderResult(e instanceof Error ? e.message : "Order failed");
    } finally {
      setPlacing(false);
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-slate-900 flex items-center justify-center">
        <div className="animate-pulse text-slate-400">Loading fund details...</div>
      </div>
    );
  }

  if (!fund) return null;

  return (
    <div className="min-h-screen bg-slate-900 text-white">
      <main className="max-w-4xl mx-auto px-6 py-8">
        <div className="mb-8">
          <p className="text-sm text-slate-400 mb-1">{fund.amc_code} · {fund.category}</p>
          <h1 className="text-2xl font-bold">{fund.scheme_name}</h1>
          <p className="text-slate-400 text-sm mt-1">Scheme Code: {fund.scheme_code}</p>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-8">
          <Stat label="Current NAV" value={formatINR(fund.nav)} />
          <Stat label="NAV Date" value={fund.nav_date} />
          <Stat label="Min Purchase" value={formatINR(fund.min_purchase)} />
          <Stat label="Min SIP" value={formatINR(fund.min_sip)} />
        </div>

        <div className="bg-slate-800 rounded-xl p-6 mb-6">
          <h3 className="font-semibold mb-4">NAV History</h3>
          <NAVChart schemeCode={fund.scheme_code} />
        </div>

        <div className="bg-slate-800 rounded-xl p-6">
          <h3 className="font-semibold mb-4">Invest Now (Lumpsum)</h3>
          <div className="flex gap-3 flex-col sm:flex-row">
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder={`Min ₹${fund.min_purchase}`}
              min={fund.min_purchase}
              className="flex-1 px-4 py-3 rounded-lg bg-slate-700 border border-slate-600 text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <button
              onClick={handleBuy}
              disabled={placing || !amount}
              className="px-8 py-3 bg-green-600 hover:bg-green-500 disabled:opacity-50 rounded-lg font-semibold transition-colors"
            >
              {placing ? "Placing..." : "Buy"}
            </button>
          </div>
          {orderResult && (
            <p className="mt-3 text-sm text-green-400">{orderResult}</p>
          )}
        </div>
      </main>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: string }) {
  return (
    <div className="bg-slate-800 rounded-xl p-4">
      <p className="text-xs text-slate-400 mb-1">{label}</p>
      <p className="font-semibold text-white">{value}</p>
    </div>
  );
}
