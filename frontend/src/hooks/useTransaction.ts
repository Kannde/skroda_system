import { useState, useEffect, useCallback } from 'react'
import type { Transaction } from '../types'
import api from '../api/client'

export function useTransactions() {
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchTransactions = useCallback(async () => {
    try {
      setIsLoading(true)
      const res = await api.get<Transaction[]>('/transactions')
      setTransactions(res.data ?? [])
    } catch {
      setError('Failed to load transactions')
    } finally {
      setIsLoading(false)
    }
  }, [])

  useEffect(() => {
    fetchTransactions()
  }, [fetchTransactions])

  return { transactions, isLoading, error, refetch: fetchTransactions }
}

export function useTransaction(id: string) {
  const [transaction, setTransaction] = useState<Transaction | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) return
    api.get<Transaction>(`/transactions/${id}`)
      .then((res) => setTransaction(res.data))
      .catch(() => setError('Transaction not found'))
      .finally(() => setIsLoading(false))
  }, [id])

  return { transaction, isLoading, error }
}
