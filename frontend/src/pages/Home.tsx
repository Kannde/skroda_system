import { Link } from 'react-router-dom'
import { ShieldCheck, Users, Zap } from 'lucide-react'

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-emerald-50 to-white">
      <div className="max-w-6xl mx-auto px-4 py-20 text-center">
        <h1 className="text-5xl font-bold text-gray-900 mb-4">
          Trade Across Cities,{' '}
          <span className="text-emerald-600">Without the Risk</span>
        </h1>
        <p className="text-xl text-gray-600 mb-10 max-w-2xl mx-auto">
          Skroda is Africa's escrow platform for peer-to-peer commerce. Funds held
          safely until both parties are satisfied.
        </p>
        <div className="flex gap-4 justify-center">
          <Link
            to="/register"
            className="bg-emerald-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-emerald-700 transition"
          >
            Get Started
          </Link>
          <Link
            to="/login"
            className="border border-emerald-600 text-emerald-600 px-8 py-3 rounded-lg font-semibold hover:bg-emerald-50 transition"
          >
            Log In
          </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-20">
          {[
            { icon: ShieldCheck, title: 'Secure Escrow', desc: 'Funds locked until delivery confirmed' },
            { icon: Users, title: 'Local Agents', desc: 'Trusted agents in every city to mediate' },
            { icon: Zap, title: 'Fast Payouts', desc: 'Instant release via MoMo or bank transfer' },
          ].map(({ icon: Icon, title, desc }) => (
            <div key={title} className="bg-white rounded-2xl p-8 shadow-sm border">
              <Icon className="w-10 h-10 text-emerald-600 mb-4 mx-auto" />
              <h3 className="text-lg font-semibold text-gray-900 mb-2">{title}</h3>
              <p className="text-gray-500">{desc}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
