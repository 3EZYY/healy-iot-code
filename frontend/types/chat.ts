// src/types/chat.ts — NEW v4.0.0 (Blueprint §7.2)

export type ChatRole = 'user' | 'assistant' | 'system'

export interface ChatMessage {
  id: string
  role: ChatRole
  content: string
  timestamp: Date
  isStreaming?: boolean   // True saat respons AI sedang di-stream
  isError?: boolean       // True saat terjadi error Groq
}

export interface ChatContext {
  temperature: number
  bpm: number
  spo2: number
  tempStatus: string
  spo2Status: string
  overallStatus: string
  timestamp: string
}
