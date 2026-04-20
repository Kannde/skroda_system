import { useNavigate, Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import { useAuth } from '../context/AuthContext'

interface FormData {
  full_name: string
  email: string
  phone: string
  password: string
  role: string
  city: string
  country: string
}

export default function Register() {
  const { register: registerUser } = useAuth()
  const navigate = useNavigate()
  const { register, handleSubmit, formState: { errors, isSubmitting } } = useForm<FormData>({
    defaultValues: { role: 'buyer', country: 'Nigeria' },
  })

  const onSubmit = async (data: FormData) => {
    try {
      await registerUser(data)
      navigate('/dashboard')
    } catch (err: any) {
      toast.error(err.response?.data?.error ?? 'Registration failed')
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4 py-10">
      <div className="w-full max-w-sm bg-white rounded-2xl border p-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">Create Account</h1>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {[
            { name: 'full_name', label: 'Full Name', type: 'text' },
            { name: 'email', label: 'Email', type: 'email' },
            { name: 'phone', label: 'Phone', type: 'tel' },
            { name: 'password', label: 'Password', type: 'password' },
            { name: 'city', label: 'City', type: 'text' },
            { name: 'country', label: 'Country', type: 'text' },
          ].map(({ name, label, type }) => (
            <div key={name}>
              <label className="block text-sm font-medium text-gray-700 mb-1">{label}</label>
              <input type={type} {...register(name as keyof FormData, { required: 'Required' })}
                className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500" />
              {errors[name as keyof FormData] && (
                <p className="text-red-500 text-xs mt-1">{errors[name as keyof FormData]?.message}</p>
              )}
            </div>
          ))}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Role</label>
            <select {...register('role')} className="w-full border rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500">
              <option value="buyer">Buyer</option>
              <option value="seller">Seller</option>
              <option value="agent">Agent</option>
            </select>
          </div>

          <button type="submit" disabled={isSubmitting}
            className="w-full bg-emerald-600 text-white py-2.5 rounded-lg font-semibold hover:bg-emerald-700 disabled:opacity-60 transition">
            {isSubmitting ? 'Creating account...' : 'Create Account'}
          </button>
        </form>
        <p className="text-sm text-center text-gray-500 mt-4">
          Already have an account? <Link to="/login" className="text-emerald-600 hover:underline">Log in</Link>
        </p>
      </div>
    </div>
  )
}
