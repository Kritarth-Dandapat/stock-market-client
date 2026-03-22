import { PortfolioSummary, formatINR, formatPercent } from "@/lib/api-client";

interface PortfolioSummaryCardProps {
  summary: PortfolioSummary;
}

export default function PortfolioSummaryCard({ summary }: PortfolioSummaryCardProps) {
  const absoluteReturnPct =
    summary.total_invested > 0
      ? ((summary.current_value - summary.total_invested) / summary.total_invested) * 100
      : 0;

  const isPositive = summary.current_value >= summary.total_invested;

  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <SummaryTile
        label="Current Value"
        value={formatINR(summary.current_value)}
        highlight
      />
      <SummaryTile
        label="Total Invested"
        value={formatINR(summary.total_invested)}
      />
      <SummaryTile
        label="Absolute Return"
        value={formatPercent(absoluteReturnPct)}
        positive={isPositive}
      />
      <SummaryTile
        label="XIRR"
        value={formatPercent(summary.xirr)}
        positive={summary.xirr >= 0}
      />
    </div>
  );
}

function SummaryTile({
  label,
  value,
  highlight,
  positive,
}: {
  label: string;
  value: string;
  highlight?: boolean;
  positive?: boolean;
}) {
  let valueClass = "text-white";
  if (positive === true) valueClass = "text-green-400";
  if (positive === false) valueClass = "text-red-400";
  if (highlight) valueClass = "text-2xl font-bold text-white";

  return (
    <div className="bg-slate-800 rounded-xl p-4">
      <p className="text-xs text-slate-400 mb-1">{label}</p>
      <p className={`font-semibold ${valueClass}`}>{value}</p>
    </div>
  );
}
