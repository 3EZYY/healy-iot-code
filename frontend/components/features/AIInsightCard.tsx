'use client'
import { useState } from 'react'
import { Sparkle, Bot, AlertTriangle } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'
import {
  getStoredGroqKey,
  validateGroqKey,
  buildInsightPrompt,
  callGroqInsight,
  type GroqInsightRequest
} from '@/lib/groq-client'

interface AIInsightCardProps {
  currentData: GroqInsightRequest | null
}

export function AIInsightCard({ currentData }: AIInsightCardProps) {
  const [insight, setInsight] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleGenerateInsight = async () => {
    const apiKey = getStoredGroqKey()

    if (!apiKey || !validateGroqKey(apiKey)) {
      setError('API Key Groq belum dikonfigurasi atau tidak valid. Masukkan key di menu Settings atau Sidebar.')
      return
    }

    if (!currentData) {
      setError('Tidak ada data sensor. Pastikan perangkat HEALY terhubung.')
      return
    }

    setLoading(true)
    setError(null)
    setInsight('')

    try {
      const prompt = buildInsightPrompt(currentData)
      await callGroqInsight(prompt, apiKey, (streamedText) => {
        setInsight(streamedText) // Update text as it streams
      })
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : 'Gagal menghubungi Groq AI. Periksa koneksi dan API key.'
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bg-healy-ai-surface border border-healy-border rounded-[16px] p-5 shadow-card">
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Bot size={20} className="text-healy-ai-accent" />
          <h3 className="font-display font-semibold text-healy-graphite text-sm">
            AI Health Insight
          </h3>
          <span className="text-xs bg-healy-ai-accent/20 text-healy-ai-accent px-2 py-0.5 rounded-full font-mono">
            Powered by Groq
          </span>
        </div>

        <button
          onClick={handleGenerateInsight}
          disabled={loading}
          className="flex items-center gap-1.5 bg-healy-ai-accent text-white text-xs font-medium
                     px-3 py-1.5 rounded-lg hover:bg-teal-700 disabled:opacity-50
                     disabled:cursor-not-allowed transition-all duration-200"
        >
          <Sparkle size={14} />
          {loading ? 'Menganalisis...' : 'Generate Insight'}
        </button>
      </div>

      {/* Konten */}
      <AnimatePresence mode="wait">
        {error && (
          <motion.div
            initial={{ opacity: 0, y: 4 }}
            animate={{ opacity: 1, y: 0 }}
            className="flex items-start gap-2 text-healy-critical text-sm"
          >
            <AlertTriangle size={16} className="mt-0.5 shrink-0" />
            <p>{error}</p>
          </motion.div>
        )}

        {loading && !insight && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="space-y-2"
          >
            {[100, 75, 90].map((w, i) => (
              <div
               key={i}
                className="h-3 bg-healy-border rounded animate-pulse"
                style={{ width: `${w}%` }}
              />
            ))}
          </motion.div>
        )}

        {insight && (
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="text-sm text-healy-graphite leading-relaxed font-body"
          >
            {insight}
            {loading && (
              <span className="inline-block w-1 h-4 bg-healy-ai-accent ml-1 animate-pulse" />
            )}
          </motion.p>
        )}

        {!loading && !insight && !error && (
          <p className="text-xs text-healy-slate italic">
            Tekan tombol untuk mendapatkan analisis kondisi kesehatan pasien dari AI.
          </p>
        )}
      </AnimatePresence>
    </div>
  )
}
