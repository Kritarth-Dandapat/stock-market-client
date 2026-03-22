"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { useQuery } from "@tanstack/react-query";
import { formatINR } from "@/lib/api-client";

interface NAVPoint {
  date: string;
  nav: number;
}

// In production this would call /api/funds/:schemeCode/nav-history
async function fetchNAVHistory(schemeCode: string): Promise<NAVPoint[]> {
  // Return placeholder data for demo purposes
  const today = new Date();
  return Array.from({ length: 30 }, (_, i) => {
    const d = new Date(today);
    d.setDate(d.getDate() - (29 - i));
    const base = 50 + Math.random() * 10;
    return {
      date: d.toLocaleDateString("en-IN", { day: "2-digit", month: "short" }),
      nav: parseFloat(base.toFixed(4)),
    };
  });
}

interface NAVChartProps {
  schemeCode: string;
}

export default function NAVChart({ schemeCode }: NAVChartProps) {
  const { data, isLoading } = useQuery({
    queryKey: ["nav-history", schemeCode],
    queryFn: () => fetchNAVHistory(schemeCode),
  });

  if (isLoading) {
    return (
      <div className="h-48 flex items-center justify-center text-slate-400 animate-pulse">
        Loading chart…
      </div>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={200}>
      <LineChart data={data} margin={{ top: 5, right: 5, bottom: 5, left: 5 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#334155" />
        <XAxis
          dataKey="date"
          tick={{ fill: "#94a3b8", fontSize: 11 }}
          tickLine={false}
        />
        <YAxis
          tick={{ fill: "#94a3b8", fontSize: 11 }}
          tickLine={false}
          axisLine={false}
          tickFormatter={(v: number) => `₹${v.toFixed(2)}`}
        />
        <Tooltip
          contentStyle={{ backgroundColor: "#1e293b", border: "1px solid #334155", borderRadius: "8px" }}
          labelStyle={{ color: "#94a3b8", fontSize: 12 }}
          formatter={(value) => {
            const num = typeof value === "number" ? value : 0;
            return [formatINR(num), "NAV"];
          }}
        />
        <Line
          type="monotone"
          dataKey="nav"
          stroke="#3b82f6"
          strokeWidth={2}
          dot={false}
          activeDot={{ r: 4, fill: "#3b82f6" }}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
