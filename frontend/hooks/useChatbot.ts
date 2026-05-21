// src/hooks/useChatbot.ts — NEW v4.0.0 (Blueprint §7.3)

import { useState, useCallback, useRef } from 'react'
import { ChatMessage, ChatContext } from '@/types/chat'
import { getStoredGroqKey, validateGroqKey, callGroqChat } from '@/lib/groq-client'

const MAX_HISTORY = 20  // Batas pesan dalam satu sesi untuk hemat token

function buildSystemPrompt(ctx: ChatContext): string {
  return `Kamu adalah HEALY AI, asisten kesehatan cerdas yang terintegrasi dengan sistem monitoring robot HEALY.

DATA SENSOR PASIEN SAAT INI (diperbarui secara real-time):
- Suhu Tubuh : ${ctx.temperature}°C — Status: ${ctx.tempStatus}
- Detak Jantung : ${ctx.bpm} bpm
- Saturasi Oksigen (SpO2) : ${ctx.spo2}% — Status: ${ctx.spo2Status}
- Status Keseluruhan : ${ctx.overallStatus}
- Waktu Pengukuran : ${ctx.timestamp}

PANDUAN RESPONS:
- Jawab dalam Bahasa Indonesia yang ramah dan mudah dipahami.
- Selalu referensikan data di atas saat relevan dengan pertanyaan.
- Jika status CRITICAL, tegaskan urgensi dengan tenang tanpa menimbulkan kepanikan berlebihan.
- Berikan informasi medis yang akurat namun selalu sarankan konsultasi dokter untuk diagnosis resmi.
- Jangan buat diagnosis — berikan interpretasi dan panduan tindakan awal.
- Jika ditanya hal di luar konteks kesehatan, jawab singkat dan kembalikan ke topik kesehatan pasien.`
}

export function useChatbot() {
  const [isOpen, setIsOpen]           = useState(false)
  const [messages, setMessages]       = useState<ChatMessage[]>([])
  const [isLoading, setIsLoading]     = useState(false)
  const [inputValue, setInputValue]   = useState('')
  const abortRef                      = useRef<AbortController | null>(null)

  const openChat  = useCallback(() => setIsOpen(true),  [])
  const closeChat = useCallback(() => {
    setIsOpen(false)
    // Batalkan request yang sedang berjalan saat panel ditutup
    abortRef.current?.abort()
  }, [])

  const clearHistory = useCallback(() => setMessages([]), [])

  const sendMessage = useCallback(async (
    userInput: string,
    context: ChatContext
  ) => {
    const trimmed = userInput.trim()
    if (!trimmed || isLoading) return

    const apiKey = getStoredGroqKey()
    if (!apiKey || !validateGroqKey(apiKey)) {
      setMessages(prev => [...prev, {
        id: `err-${Date.now()}`,
        role: 'assistant',
        content: 'Groq API Key belum dikonfigurasi. Masukkan key di bagian bawah sidebar.',
        timestamp: new Date(),
        isError: true,
      }])
      return
    }

    // Tambahkan pesan user ke history
    const userMsg: ChatMessage = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: trimmed,
      timestamp: new Date(),
    }

    // Placeholder untuk respons AI (streaming)
    const assistantMsgId = `ai-${Date.now()}`
    const assistantMsg: ChatMessage = {
      id: assistantMsgId,
      role: 'assistant',
      content: '',
      timestamp: new Date(),
      isStreaming: true,
    }

    setMessages(prev => [...prev, userMsg, assistantMsg])
    setInputValue('')
    setIsLoading(true)

    // Rakit message history untuk Groq API
    // Sertakan system prompt + history percakapan (max MAX_HISTORY pesan terakhir)
    const historyForApi = messages
      .slice(-MAX_HISTORY)
      .filter(m => m.role !== 'system' && !m.isError)
      .map(m => ({ role: m.role, content: m.content }))

    const payload = [
      { role: 'system' as const, content: buildSystemPrompt(context) },
      ...historyForApi,
      { role: 'user' as const, content: trimmed },
    ]

    abortRef.current = new AbortController()

    try {
      const response = await callGroqChat(payload, apiKey, abortRef.current.signal)

      if (!response.ok) {
        const errData = await response.json()
        throw new Error(errData.error?.message || `HTTP ${response.status}`)
      }

      // Stream respons
      const reader  = response.body!.getReader()
      const decoder = new TextDecoder()
      let accumulated = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        const chunk = decoder.decode(value)
        const lines = chunk.split('\n').filter(l => l.startsWith('data: '))

        for (const line of lines) {
          const json = line.replace('data: ', '').trim()
          if (json === '[DONE]') continue
          try {
            const parsed = JSON.parse(json)
            const delta  = parsed.choices?.[0]?.delta?.content || ''
            accumulated += delta

            // Update pesan AI secara real-time
            setMessages(prev => prev.map(m =>
              m.id === assistantMsgId
                ? { ...m, content: accumulated }
                : m
            ))
          } catch { /* skip malformed chunk */ }
        }
      }

      // Finalize: hapus flag isStreaming
      setMessages(prev => prev.map(m =>
        m.id === assistantMsgId
          ? { ...m, isStreaming: false }
          : m
      ))

    } catch (err: unknown) {
      if (err instanceof Error && err.name === 'AbortError') return // User menutup panel, abaikan

      const errorMessage = err instanceof Error ? err.message : 'Unknown error'
      setMessages(prev => prev.map(m =>
        m.id === assistantMsgId
          ? {
              ...m,
              content: `Gagal mendapatkan respons: ${errorMessage}`,
              isStreaming: false,
              isError: true,
            }
          : m
      ))
    } finally {
      setIsLoading(false)
    }
  }, [messages, isLoading])

  return {
    isOpen, openChat, closeChat,
    messages, isLoading,
    inputValue, setInputValue,
    sendMessage, clearHistory,
  }
}
