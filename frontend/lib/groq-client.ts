const GROQ_API_URL = 'https://api.groq.com/openai/v1/chat/completions'
const GROQ_MODEL   = 'llama-3.1-8b-instant'

export function getStoredGroqKey(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('healy_groq_api_key')
}

export function setStoredGroqKey(key: string): void {
  localStorage.setItem('healy_groq_api_key', key)
}

export function validateGroqKey(key: string): boolean {
  return /^gsk_[a-zA-Z0-9]{50,}$/.test(key.trim())
}

export interface GroqInsightRequest {
  temperature: number
  bpm: number
  spo2: number
  avgTemp?: number
  avgBpm?: number
  avgSpo2?: number
  status: string
  durationMinutes?: number
}

export function buildInsightPrompt(data: GroqInsightRequest): string {
  const avg = data.avgTemp
    ? `Rata-rata 15 menit terakhir: Suhu ${data.avgTemp}°C, BPM ${data.avgBpm}, SpO2 ${data.avgSpo2}%.`
    : ''

  return `Kamu adalah asisten medis AI untuk sistem pemantauan kesehatan HEALY.

Data biometrik pasien saat ini:
- Suhu tubuh: ${data.temperature}°C
- Detak jantung (BPM): ${data.bpm} bpm
- Saturasi oksigen (SpO2): ${data.spo2}%
- Status keseluruhan sistem: ${data.status}
${avg}

Berikan analisis klinis singkat dalam 2–3 kalimat dalam Bahasa Indonesia. 
Fokus pada: (1) interpretasi kondisi saat ini, (2) potensi risiko jika ada, 
(3) satu rekomendasi tindakan. 
Gunakan bahasa yang mudah dipahami keluarga pasien, bukan bahasa medis teknis.
Jangan tambahkan disclaimer atau pengantar panjang.`
}

export async function callGroqInsight(
  prompt: string,
  apiKey: string,
  onChunk?: (text: string) => void
): Promise<string> {
  const response = await fetch(GROQ_API_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${apiKey}`,
    },
    body: JSON.stringify({
      model: GROQ_MODEL,
      messages: [{ role: 'user', content: prompt }],
      max_tokens: 200,
      stream: !!onChunk,
      temperature: 0.3,
    }),
  })

  if (!response.ok) {
    const err = await response.json()
    throw new Error(err.error?.message || `Groq API error: ${response.status}`)
  }

  // Streaming mode
  if (onChunk && response.body) {
    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let fullText = ''

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
          const text = parsed.choices?.[0]?.delta?.content || ''
          fullText += text
          onChunk(fullText)
        } catch {}
      }
    }
    return fullText
  }

  // Non-streaming mode
  const json = await response.json()
  return json.choices?.[0]?.message?.content || ''
}

// ─── CHAT (MULTI-TURN) — digunakan useChatbot ───
// callGroqChat adalah thin wrapper — logika streaming di-handle oleh useChatbot
// karena chatbot membutuhkan kendali penuh atas state per-message

export interface GroqChatMessage {
  role:    'user' | 'assistant' | 'system'
  content: string
}

export async function callGroqChat(
  messages: GroqChatMessage[],
  apiKey:   string,
  signal?:  AbortSignal
): Promise<Response> {
  // Kembalikan raw Response agar caller bisa stream secara fleksibel
  return fetch(GROQ_API_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${apiKey}`,
    },
    body: JSON.stringify({
      model:       GROQ_MODEL,
      messages,
      max_tokens:  512,
      stream:      true,
      temperature: 0.4,
    }),
    signal,
  })
}

