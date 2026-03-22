"use client";

import { useQuery } from "@tanstack/react-query";
import { ordersApi, formatINR } from "@/lib/api-client";
import Link from "next/link";

const statusColors: Record<string, string> = {
  pending: "text-yellow-400",
  submitted: "text-blue-400",
  processing: "text-blue-300",
  completed: "text-green-400",
  failed: "text-red-400",
  cancelled: "text-slate-400",
};

export default function TransactionsPage() {
  const { data: orders, isLoading, error } = useQuery({
    queryKey: ["orders"],
    queryFn: ordersApi.list,
  });

  return (
    <div className="min-h-screen bg-slate-900 text-white">
      <nav className="border-b border-slate-700 px-6 py-4 flex items-center justify-between">
        <h1 className="text-xl font-bold">BSE MF Dashboard</h1>
        <div className="flex gap-4 text-sm text-slate-400">
          <Link href="/portfolio" className="hover:text-white">Portfolio</Link>
          <Link href="/funds" className="hover:text-white">Funds</Link>
          <Link href="/sip" className="hover:text-white">SIP</Link>
          <Link href="/transactions" className="text-white font-medium">Transactions</Link>
          <Link href="/statements" className="hover:text-white">Statements</Link>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <h2 className="text-2xl font-semibold mb-6">Transactions</h2>

        {isLoading && (
          <div className="space-y-3 animate-pulse">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="h-16 bg-slate-700 rounded-xl" />
            ))}
          </div>
        )}

        {error && (
          <div className="p-4 bg-red-900/50 border border-red-700 rounded-xl text-red-300">
            Failed to load transactions.
          </div>
        )}

        {orders && orders.length === 0 && (
          <div className="text-center py-16 bg-slate-800 rounded-xl">
            <p className="text-slate-400 mb-4">No transactions yet.</p>
            <Link
              href="/funds"
              className="px-6 py-2 bg-blue-600 hover:bg-blue-500 rounded-lg font-medium transition-colors"
            >
              Start Investing
            </Link>
          </div>
        )}

        {orders && orders.length > 0 && (
          <div className="bg-slate-800 rounded-xl overflow-hidden">
            <table className="w-full text-sm">
              <thead className="border-b border-slate-700">
                <tr className="text-slate-400">
                  <th className="text-left p-4">Fund</th>
                  <th className="text-left p-4">Type</th>
                  <th className="text-right p-4">Amount</th>
                  <th className="text-right p-4">Status</th>
                  <th className="text-right p-4">Date</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-700">
                {orders.map((order) => (
                  <tr key={order.id} className="hover:bg-slate-700/50">
                    <td className="p-4">
                      <p className="font-medium text-white">{order.scheme_name || order.scheme_code}</p>
                      {order.bse_order_id && (
                        <p className="text-xs text-slate-500">BSE: {order.bse_order_id}</p>
                      )}
                    </td>
                    <td className="p-4 capitalize text-slate-300">{order.order_type}</td>
                    <td className="p-4 text-right font-medium">{formatINR(order.amount)}</td>
                    <td className={`p-4 text-right capitalize font-medium ${statusColors[order.status] ?? ""}`}>
                      {order.status}
                    </td>
                    <td className="p-4 text-right text-slate-400">
                      {new Date(order.created_at).toLocaleDateString("en-IN")}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </main>
    </div>
  );
}
