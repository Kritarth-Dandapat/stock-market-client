import Link from "next/link";

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8 bg-gradient-to-br from-slate-900 to-slate-800 text-white">
      <div className="max-w-2xl w-full text-center space-y-8">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">
            BSE Mutual Fund Dashboard
          </h1>
          <p className="mt-4 text-slate-300 text-lg">
            Invest in mutual funds securely via BSE StAR MF — SEBI regulated.
          </p>
        </div>

        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <Link
            href="/login"
            className="px-8 py-3 bg-blue-600 hover:bg-blue-500 rounded-lg font-semibold transition-colors"
          >
            Login
          </Link>
          <Link
            href="/register"
            className="px-8 py-3 bg-slate-700 hover:bg-slate-600 rounded-lg font-semibold transition-colors"
          >
            Create Account
          </Link>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 mt-12">
          <Feature
            title="SEBI Regulated"
            description="Fully compliant with SEBI Investment Adviser Regulations and PMLA requirements."
          />
          <Feature
            title="BSE StAR MF"
            description="Orders placed directly on BSE's national mutual fund platform with NACH SIP support."
          />
          <Feature
            title="Enterprise Security"
            description="AES-256 PII encryption, JWT auth, rate limiting and immutable audit logs."
          />
        </div>
      </div>
    </main>
  );
}

function Feature({ title, description }: { title: string; description: string }) {
  return (
    <div className="p-6 rounded-xl bg-slate-700/50 border border-slate-600 text-left">
      <h3 className="font-semibold text-white">{title}</h3>
      <p className="mt-2 text-sm text-slate-300">{description}</p>
    </div>
  );
}
