import Link from "next/link";

export default function StatementsPage() {
  return (
    <div className="min-h-screen bg-slate-900 text-white">
      <nav className="border-b border-slate-700 px-6 py-4 flex items-center justify-between">
        <h1 className="text-xl font-bold">BSE MF Dashboard</h1>
        <div className="flex gap-4 text-sm text-slate-400">
          <Link href="/portfolio" className="hover:text-white">Portfolio</Link>
          <Link href="/funds" className="hover:text-white">Funds</Link>
          <Link href="/sip" className="hover:text-white">SIP</Link>
          <Link href="/transactions" className="hover:text-white">Transactions</Link>
          <Link href="/statements" className="text-white font-medium">Statements</Link>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto px-6 py-8">
        <h2 className="text-2xl font-semibold mb-6">Statements & Reports</h2>

        <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <StatementCard
            title="Capital Gains Statement"
            description="Download your LTCG / STCG statement for income tax filing."
            badge="PDF"
          />
          <StatementCard
            title="CAS (Consolidated Account Statement)"
            description="Complete view of all mutual fund holdings across AMCs via CAMS/KFintech."
            badge="PDF"
          />
          <StatementCard
            title="Transaction Statement"
            description="Detailed list of all buy, sell, and SIP transactions."
            badge="PDF / Excel"
          />
          <StatementCard
            title="TDS Certificate (Form 16A)"
            description="Tax deducted at source certificate for dividend income."
            badge="PDF"
          />
        </div>
      </main>
    </div>
  );
}

function StatementCard({
  title,
  description,
  badge,
}: {
  title: string;
  description: string;
  badge: string;
}) {
  return (
    <div className="bg-slate-800 rounded-xl p-6 flex flex-col justify-between">
      <div>
        <div className="flex items-start justify-between mb-3">
          <h3 className="font-semibold text-white">{title}</h3>
          <span className="text-xs bg-slate-700 text-slate-300 px-2 py-1 rounded font-mono">{badge}</span>
        </div>
        <p className="text-sm text-slate-400">{description}</p>
      </div>
      <button
        disabled
        className="mt-4 w-full py-2 bg-slate-700 hover:bg-slate-600 disabled:opacity-50 rounded-lg text-sm font-medium transition-colors"
        title="Available once you have transactions"
      >
        Download
      </button>
    </div>
  );
}
