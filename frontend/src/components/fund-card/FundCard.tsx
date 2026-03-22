import { FundScheme, formatINR } from "@/lib/api-client";
import Link from "next/link";

interface FundCardProps {
  fund: FundScheme;
}

export default function FundCard({ fund }: FundCardProps) {
  return (
    <Link
      href={`/funds/${encodeURIComponent(fund.scheme_code)}`}
      className="block bg-slate-800 hover:bg-slate-700 border border-slate-700 hover:border-slate-500 rounded-xl p-5 transition-all"
    >
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1 min-w-0">
          <p className="text-xs text-blue-400 font-medium mb-1">{fund.amc_code}</p>
          <h3 className="font-semibold text-white text-sm leading-tight line-clamp-2">
            {fund.scheme_name}
          </h3>
        </div>
      </div>

      <div className="flex items-center justify-between text-sm">
        <span className="text-xs bg-slate-700 text-slate-300 px-2 py-1 rounded-full">
          {fund.category}
        </span>
        <div className="text-right">
          <p className="text-xs text-slate-400">NAV</p>
          <p className="font-semibold text-white">{formatINR(fund.nav)}</p>
        </div>
      </div>

      <div className="mt-3 pt-3 border-t border-slate-700 flex justify-between text-xs text-slate-400">
        <span>Min: {formatINR(fund.min_purchase)}</span>
        <span>SIP: {formatINR(fund.min_sip)}/mo</span>
      </div>
    </Link>
  );
}
