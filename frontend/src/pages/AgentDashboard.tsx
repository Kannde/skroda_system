import { useAuth } from '../context/AuthContext'

export default function AgentDashboard() {
  const { user } = useAuth()

  if (user?.role !== 'agent') {
    return <div className="p-10 text-red-500">Access restricted to agents.</div>
  }

  return (
    <div className="max-w-5xl mx-auto px-4 py-10">
      <h1 className="text-2xl font-bold text-gray-900 mb-2">Agent Dashboard</h1>
      <p className="text-gray-500 mb-8">Manage disputes and transactions in your city.</p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
        {[
          { label: 'Open Disputes', value: '—' },
          { label: 'Transactions Handled', value: '—' },
          { label: 'Rating', value: '—' },
        ].map(({ label, value }) => (
          <div key={label} className="bg-white border rounded-xl p-6">
            <p className="text-sm text-gray-500 mb-1">{label}</p>
            <p className="text-3xl font-bold text-gray-900">{value}</p>
          </div>
        ))}
      </div>

      <div className="bg-white border rounded-xl p-6">
        <h2 className="font-semibold text-gray-900 mb-4">Open Disputes</h2>
        <p className="text-gray-400 text-sm">No open disputes assigned to you.</p>
      </div>
    </div>
  )
}
