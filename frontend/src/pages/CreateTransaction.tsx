import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import api from '../api/client'
import type { Transaction } from '../types'

interface FormData {
  seller_id: string
  title: string
  description: string
  amount: number
  currency: string
  seller_city: string
}

export default function CreateTransaction() {
  const navigate = useNavigate()
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<FormData>({
    defaultValues: { currency: 'NGN' },
  })

  const onSubmit = async (data: FormData) => {
    try {
      const res = await api.post<Transaction>('/transactions', data)
      toast.success('Transaction created!')
      navigate(`/transactions/${res.data.id}`)
    } catch (err: any) {
      toast.error(err.response?.data?.error ?? 'Failed to create transaction')
    }
  }

  return (
    <div className="max-w-xl mx-auto px-4 py-10">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">New Escrow Transaction</h1>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 bg-white rounded-xl border p-6">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Seller ID</label>
          <input {...register('seller_id', { required: 'Seller ID is required' })}
            className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500" />
          {errors.seller_id && <p className="text-red-500 text-xs mt-1">{errors.seller_id.message}</p>}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Title</label>
          <input {...register('title', { required: 'Title is required', minLength: { value: 5, message: 'Min 5 characters' } })}
            className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500" />
          {errors.title && <p className="text-red-500 text-xs mt-1">{errors.title.message}</p>}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea {...register('description')} rows={3}
            className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500" />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Amount</label>
            <input type="number" step="0.01" {...register('amount', { required: 'Required', min: { value: 1, message: 'Must be > 0' } })}
              className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500" />
            {errors.amount && <p className="text-red-500 text-xs mt-1">{errors.amount.message}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Currency</label>
            <select {...register('currency')} className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500">
              <option value="NGN">NGN</option>
              <option value="GHS">GHS</option>
              <option value="KES">KES</option>
              <option value="ZAR">ZAR</option>
              <option value="USD">USD</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Seller's City</label>
          <input {...register('seller_city', { required: 'Required' })}
            className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500" />
          {errors.seller_city && <p className="text-red-500 text-xs mt-1">{errors.seller_city.message}</p>}
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className="w-full bg-emerald-600 text-white py-2.5 rounded-lg font-semibold hover:bg-emerald-700 disabled:opacity-60 transition"
        >
          {isSubmitting ? 'Creating...' : 'Create Transaction'}
        </button>
      </form>
    </div>
  )
}
