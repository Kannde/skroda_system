import { useParams } from 'react-router-dom'
import { useTransaction } from '../hooks/useTransaction'
import { formatCurrency, formatDate, statusLabel, statusColor } from '../utils/format'

export default function TransactionDetail() {
  const { id } = useParams<{ id: string }>()
  const { transaction, isLoading, error } = useTransaction(id ?? '')

  if (isLoading) return <div className="p-10 text-gray-500">Loading...</div>
  if (error || !transaction) return <div className="p-10 text-red-500">{error ?? 'Not found'}</div>

  return (
    <div className="max-w-2xl mx-auto px-4 py-10">
      <div className="bg-white rounded-xl border p-6">
        <div className="flex items-start justify-between mb-6">
          <div>
            <h1 className="text-xl font-bold text-gray-900">{transaction.title}</h1>
            <p className="text-sm text-gray-500 mt-1">{formatDate(transaction.created_at)}</p>
          </div>
          <span className={`px-3 py-1 rounded-full text-sm font-medium ${statusColor(transaction.status)}`}>
            {statusLabel(transaction.status)}
          </span>
        </div>

        {transaction.description && (
          <p className="text-gray-600 mb-6 text-sm">{transaction.description}</p>
        )}

        <dl className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <dt className="text-gray-500">Amount</dt>
            <dd className="font-semibold text-gray-900">{formatCurrency(transaction.amount, transaction.currency)}</dd>
          </div>
          <div>
            <dt className="text-gray-500">Currency</dt>
            <dd className="font-semibold text-gray-900">{transaction.currency}</dd>
          </div>
          <div>
            <dt className="text-gray-500">Seller City</dt>
            <dd className="font-semibold text-gray-900">{transaction.seller_city}</dd>
          </div>
          {transaction.inspection_ends_at && (
            <div>
              <dt className="text-gray-500">Inspection Ends</dt>
              <dd className="font-semibold text-gray-900">{formatDate(transaction.inspection_ends_at)}</dd>
            </div>
          )}
        </dl>
      </div>
    </div>
  )
}
