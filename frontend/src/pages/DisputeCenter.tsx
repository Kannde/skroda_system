import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import api from '../api/client'

interface FormData {
  transaction_id: string
  reason: string
}

export default function DisputeCenter() {
  const { register, handleSubmit, reset, formState: { errors, isSubmitting } } = useForm<FormData>()

  const onSubmit = async (data: FormData) => {
    try {
      await api.post('/disputes', data)
      toast.success('Dispute raised. An agent will review it shortly.')
      reset()
    } catch (err: any) {
      toast.error(err.response?.data?.error ?? 'Failed to raise dispute')
    }
  }

  return (
    <div className="max-w-xl mx-auto px-4 py-10">
      <h1 className="text-2xl font-bold text-gray-900 mb-2">Dispute Center</h1>
      <p className="text-gray-500 mb-6 text-sm">
        Have an issue with a transaction? Raise a dispute and a local agent will mediate.
      </p>

      <form onSubmit={handleSubmit(onSubmit)} className="bg-white rounded-xl border p-6 space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Transaction ID</label>
          <input {...register('transaction_id', { required: 'Required' })}
            className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-red-500" />
          {errors.transaction_id && <p className="text-red-500 text-xs mt-1">{errors.transaction_id.message}</p>}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Reason</label>
          <textarea {...register('reason', { required: 'Required', minLength: { value: 20, message: 'Please provide more detail (min 20 chars)' } })}
            rows={5} placeholder="Describe the issue in detail..."
            className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-red-500" />
          {errors.reason && <p className="text-red-500 text-xs mt-1">{errors.reason.message}</p>}
        </div>

        <button type="submit" disabled={isSubmitting}
          className="w-full bg-red-600 text-white py-2.5 rounded-lg font-semibold hover:bg-red-700 disabled:opacity-60 transition">
          {isSubmitting ? 'Submitting...' : 'Raise Dispute'}
        </button>
      </form>
    </div>
  )
}
