export type UserRole = 'buyer' | 'seller' | 'agent' | 'admin'

export interface User {
  id: string
  full_name: string
  email: string
  phone: string
  role: UserRole
  city: string
  country: string
  is_verified: boolean
  created_at: string
}

export type TransactionStatus =
  | 'pending'
  | 'funded'
  | 'in_progress'
  | 'inspection'
  | 'completed'
  | 'disputed'
  | 'cancelled'
  | 'refunded'

export interface Transaction {
  id: string
  buyer_id: string
  seller_id: string
  agent_id?: string
  title: string
  description: string
  amount: number
  currency: string
  status: TransactionStatus
  buyer_city: string
  seller_city: string
  inspection_ends_at?: string
  created_at: string
  updated_at: string
}

export type DisputeStatus = 'open' | 'under_review' | 'resolved' | 'closed'
export type DisputeResolution = 'buyer' | 'seller' | 'split'

export interface Dispute {
  id: string
  transaction_id: string
  raised_by_id: string
  agent_id?: string
  reason: string
  status: DisputeStatus
  resolution?: DisputeResolution
  resolution_note?: string
  created_at: string
  resolved_at?: string
}

export type PaymentStatus = 'pending' | 'success' | 'failed' | 'refunded'
export type PaymentProvider = 'momo' | 'stripe' | 'manual'

export interface Payment {
  id: string
  transaction_id: string
  payer_id: string
  amount: number
  currency: string
  provider: PaymentProvider
  status: PaymentStatus
  created_at: string
}

export interface AuthResponse {
  token: string
  user: User
}
