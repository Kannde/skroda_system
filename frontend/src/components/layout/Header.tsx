import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { ShieldCheck } from 'lucide-react'

export default function Header() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/')
  }

  return (
    <header className="bg-white border-b sticky top-0 z-50">
      <div className="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 font-bold text-emerald-600 text-lg">
          <ShieldCheck className="w-5 h-5" />
          Skroda
        </Link>

        <nav className="flex items-center gap-4 text-sm">
          {user ? (
            <>
              <Link to="/dashboard" className="text-gray-600 hover:text-gray-900">Dashboard</Link>
              <Link to="/disputes" className="text-gray-600 hover:text-gray-900">Disputes</Link>
              {user.role === 'agent' && (
                <Link to="/agent" className="text-gray-600 hover:text-gray-900">Agent Hub</Link>
              )}
              <button onClick={handleLogout} className="text-gray-400 hover:text-gray-900">Logout</button>
            </>
          ) : (
            <>
              <Link to="/login" className="text-gray-600 hover:text-gray-900">Log In</Link>
              <Link to="/register" className="bg-emerald-600 text-white px-4 py-1.5 rounded-lg hover:bg-emerald-700 transition">
                Sign Up
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  )
}
