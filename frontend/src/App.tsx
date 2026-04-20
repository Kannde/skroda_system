import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import { AuthProvider, useAuth } from './context/AuthContext'
import Header from './components/layout/Header'
import Home from './pages/Home'
import Login from './pages/Login'
import Register from './pages/Register'
import Dashboard from './pages/Dashboard'
import CreateTransaction from './pages/CreateTransaction'
import TransactionDetail from './pages/TransactionDetail'
import DisputeCenter from './pages/DisputeCenter'
import AgentDashboard from './pages/AgentDashboard'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth()
  if (isLoading) return null
  return user ? <>{children}</> : <Navigate to="/login" replace />
}

function AppRoutes() {
  return (
    <>
      <Header />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/dashboard" element={<PrivateRoute><Dashboard /></PrivateRoute>} />
        <Route path="/transactions/new" element={<PrivateRoute><CreateTransaction /></PrivateRoute>} />
        <Route path="/transactions/:id" element={<PrivateRoute><TransactionDetail /></PrivateRoute>} />
        <Route path="/disputes" element={<PrivateRoute><DisputeCenter /></PrivateRoute>} />
        <Route path="/agent" element={<PrivateRoute><AgentDashboard /></PrivateRoute>} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
      <Toaster position="top-right" />
    </>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <AppRoutes />
      </AuthProvider>
    </BrowserRouter>
  )
}
