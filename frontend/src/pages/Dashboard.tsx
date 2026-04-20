import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { useTransactions } from '../hooks/useTransaction'
import { formatCurrency, formatDate, statusLabel, statusColor } from '../utils/format'

export default function Dashboard() {
  const { user } = useAuth()
  const { transactions, isLoading, error } = useTransactions()

  return (
    <div className="max-w-5xl mx-auto px-4 py-10">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Welcome, {user?.full_name}</h1>
          <p className="text-gray-500 capitalize">{user?.role} · {user?.city}</p>
        </div>
        <Link
          to="/transactions/new"
          className="bg-emerald-600 text-white px-5 py-2 rounded-lg hover:bg-emerald-700 transition text-sm font-medium"
        >
          + New Transaction
        </Link>
      </div>

      {isLoading && <p className="text-gray-500">Loading transactions...</p>}
      {error && <p className="text-red-500">{error}</p>}

      {!isLoading && !error && transactions.length === 0 && (
        <div className="text-center py-20 text-gray-400">
          <p className="text-lg">No transactions yet.</p>
          <Link to="/transactions/new" className="text-emerald-600 hover:underline mt-2 inline-block">
            Create your first transaction
          </Link>
        </div>
      )}

      <div className="space-y-4">
        {transactions.map((tx) => (
          <Link
            key={tx.id}
            to={`/transactions/${tx.id}`}
            className="block bg-white rounded-xl border p-5 hover:shadow-md transition"
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="font-semibold text-gray-900">{tx.title}</p>
                <p className="text-sm text-gray-500 mt-1">{formatDate(tx.created_at)} · {tx.seller_city}</p>
              </div>
              <div className="text-right">
                <p className="font-bold text-gray-900">{formatCurrency(tx.amount, tx.currency)}</p>
                <span className={`inline-block mt-1 px-2 py-0.5 rounded text-xs font-medium ${statusColor(tx.status)}`}>
                  {statusLabel(tx.status)}
                </span>
              </div>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
