# HEALY — Master Blueprint
## Health Observer Robot: Full-Stack Implementation Guide

---

| Field | Detail |
|---|---|
| **Document Version** | v3.0.0 |
| **Status** | Active — Master Reference |
| **Role Penyusun** | Lead Developer / Senior Full-Stack Engineer |
| **Stack** | ESP32 (IoT) + Golang 1.23+ (Backend) + Next.js 14 (Frontend) |
| **Design Philosophy** | Clinical Futurism — Sage Green, Glacial White, Futuristic Medical |
| **Last Updated** | 2025 |
| **Changelog** | v3.0.0: Phase 10 — AI Insight & Device Presence. Phase 11 — Auto-Narrative Alert (inovasi baru). Arsitektur diperbarui, folder struktur diperluas, progress tracker diperbarui. |

---

## Changelog v2.1.0 → v3.0.0

| # | Area | Perubahan | Urgensi |
|---|---|---|---|
| C1 | Backend: hub.go | Broadcast JSON `system` message saat ESP32 connect/disconnect | HIGH |
| C2 | Frontend: useWebSocket.ts | Handle message `type: "system"` — trigger toast + update device state | HIGH |
| C3 | Frontend: Header/Dashboard | Indikator lampu hijau/merah status device secara live | HIGH |
| C4 | Frontend: NavSidebar.tsx | Input field Groq API Key + simpan ke `localStorage` | HIGH |
| C5 | Frontend: AIInsightCard.tsx | Komponen baru — tombol generate insight, tampilkan respons Groq streaming | HIGH |
| C6 | Frontend: AlertFeed.tsx | Update — tampilkan AI narrative di bawah setiap alert CRITICAL | HIGH |
| C7 | Frontend: useAutoNarrative.ts | Hook baru — trigger Groq otomatis saat CRITICAL event masuk (inovasi) | HIGH |
| C8 | Types: telemetry.ts | Tambah `SystemMessage`, `AlertWithNarrative` interfaces | MEDIUM |
| C9 | Arsitektur diagram | Diperbarui dengan Groq AI Cloud dan alur Phase 10–11 | INFO |
| C10 | Gotcha baru | GOTCHA-06: localStorage + Groq API Key exposure | MEDIUM |
| C11 | Progress Tracker | Phase 10 dan 11 ditambahkan sebagai TODO | INFO |

---

## Catatan Lead Dev: Keputusan Arsitektur v3.0

**Mengapa Groq di Frontend (bukan Backend)?**
Keputusan ini disengaja dan rasional. Groq API key bersifat per-user — setiap pengguna memasukkan key mereka sendiri. Jika disimpan di backend, kamu harus membangun sistem key management per-user yang jauh lebih kompleks. Dengan menyimpannya di `localStorage` dan memanggil Groq langsung dari browser, tidak ada data sensitif yang melewati server kamu. Trade-off yang diterima: key bisa dilihat di browser developer tools — ini acceptable untuk konteks academic/demo dan dijelaskan di GOTCHA-06.

**Inovasi Phase 11 — Mengapa Auto-Narrative Alert?**
Phase 10 adalah AI on-demand (user tekan tombol → dapat insight). Ini berguna tapi pasif. Phase 11 membaliknya: sistem yang aktif memberi konteks medis tanpa perlu interaksi user. Ketika SpO2 drop ke 87% pukul 03:00, tidak ada user yang tekan tombol — tapi alert feed akan otomatis berisi: *"SpO2 pasien turun ke 87%, signifikan di bawah batas aman. Potensi hipoksemia — perlu pemeriksaan segera."* Ini adalah perbedaan antara monitoring system dan intelligent health companion.

---

## Table of Contents

1. [Arsitektur Sistem Keseluruhan](#1-arsitektur-sistem-keseluruhan)
2. [Struktur Folder Final (Revised)](#2-struktur-folder-final-revised)
3. [Design System (Frontend)](#3-design-system-frontend)
4. [Backend: Spesifikasi Teknis](#4-backend-spesifikasi-teknis)
5. [Frontend: Spesifikasi Teknis](#5-frontend-spesifikasi-teknis)
6. [Database Schema](#6-database-schema)
7. [Environment Variables](#7-environment-variables)
8. [Deployment Strategy](#8-deployment-strategy)
9. [Known Technical Gotchas](#9-known-technical-gotchas)
10. [Phase 10 — AI Insight & Device Presence](#10-phase-10--ai-insight--device-presence)
11. [Phase 11 — Auto-Narrative Alert (Inovasi)](#11-phase-11--auto-narrative-alert-inovasi)
12. [Rencana Kerja Step-by-Step](#12-rencana-kerja-step-by-step)
13. [Codebase Progress Tracker](#13-codebase-progress-tracker)
14. [Konvensi Kode & Git](#14-konvensi-kode--git)
15. [Prompt Inisialisasi AI (Master Prompt)](#15-prompt-inisialisasi-ai-master-prompt)

---

## 1. Arsitektur Sistem Keseluruhan

```
┌────────────────────────────────────────────────────────────────────────────┐
│                            HEALY SYSTEM v3.0                               │
│                                                                            │
│  ┌──────────┐   I2C/I2S   ┌─────────────────────────────────────────┐    │
│  │  Sensor  │ ──────────▶ │           ESP32 WROOM-32                │    │
│  │ MLX90614 │             │       (IoT Layer / /iot folder)         │    │
│  │ MAX30102  │             │  • Baca sensor (I2C/I2S)               │    │
│  │ INMP441  │             │  • Encode JSON payload                  │    │
│  └──────────┘             └────────────────┬────────────────────────┘    │
│                                            │ WebSocket (JSON)             │
│                                            ▼                              │
│                    ┌───────────────────────────────────────────┐         │
│                    │          Golang Backend                   │         │
│                    │      Railway / Render                     │         │
│                    │                                           │         │
│                    │  • WS Hub: broadcast telemetry data       │         │
│                    │  • WS Hub: broadcast SYSTEM messages      │  ◀─ NEW │
│                    │    {"type":"system","status":"connected"} │         │
│                    │  • Threshold & Alert Engine               │         │
│                    │  • REST API (Gin)                         │         │
│                    └─────────────────┬─────────────────────────┘         │
│                                      │                                    │
│           ┌──────────────────────────┼───────────────────┐               │
│           │ WebSocket (WSS)          │ REST (/api/...)   │               │
│           ▼                          ▼                   ▼               │
│  ┌─────────────────────┐   ┌──────────────────────────────────┐         │
│  │   Next.js 14        │   │  Supabase PostgreSQL (Cloud)     │         │
│  │   Vercel            │   │                                  │         │
│  │                     │   │  telemetry.telemetry_records     │         │
│  │  Phase 10 (NEW):    │   │  public.users                    │         │
│  │  • Device LED       │   │  public.alert_logs               │         │
│  │  • Toast connect    │   │  public.device_settings          │         │
│  │  • AIInsightCard    │   └──────────────────────────────────┘         │
│  │                     │                                                 │
│  │  Phase 11 (NEW):    │            │ Direct REST API call               │
│  │  • useAutoNarrative │            ▼                                    │
│  │  • Narrative in     │   ┌──────────────────────────────────┐         │
│  │    AlertFeed        │   │       Groq AI Cloud              │         │
│  └─────────────────────┘   │   api.groq.com/openai/v1/...    │         │
│           │                │                                  │         │
│           └────────────────▶  • llama-3.1-8b-instant         │         │
│      (from browser,        │  • Streaming / non-streaming     │         │
│       API key in           │  • Medical narrative generation  │         │
│       localStorage)        └──────────────────────────────────┘         │
└────────────────────────────────────────────────────────────────────────────┘
```

**Aliran Data — Phase 10 (Device Presence):**
```
ESP32 connect/disconnect ──▶ hub.go detects event
  ──▶ Hub broadcast {"type":"system","status":"device_connected"} ke semua viewer clients
  ──▶ useWebSocket.ts menerima → set deviceOnline state
  ──▶ Header LED indicator berubah hijau/merah
  ──▶ Toast notification muncul: "Robot HEALY terhubung"
```

**Aliran Data — Phase 10 (AI Insight on-demand):**
```
User klik "Generate AI Insight"
  ──▶ AIInsightCard ambil data dari useTelemetry (latest + rata-rata 15 menit)
  ──▶ Rakit prompt medis
  ──▶ Fetch langsung ke api.groq.com (API key dari localStorage)
  ──▶ Stream respons teks ke dalam card
```

**Aliran Data — Phase 11 (Auto-Narrative Alert):**
```
Telemetry payload masuk dengan overall_status = "CRITICAL"
  ──▶ useAutoNarrative.ts mendeteksi event CRITICAL baru
  ──▶ Otomatis (tanpa user action) panggil Groq API
  ──▶ Groq generate narasi medis 1–2 kalimat dalam Bahasa Indonesia
  ──▶ Narasi di-attach ke AlertFeed entry di bawah status chip
  ──▶ User melihat alert beserta penjelasan klinis langsung
```

---

## 2. Struktur Folder Final (Revised)

```
HEALY-PROJECT/
│
├── .env.example
├── README.md
│
├── iot/                            ← (EXISTING: ESP32 C++/PlatformIO)
│
├── backend/                        ← Golang — Clean Architecture
│   ├── .env.example
│   ├── go.mod
│   ├── cmd/
│   │   └── api/main.go
│   └── internal/
│       ├── delivery/
│       │   ├── http/
│       │   │   ├── handler/
│       │   │   │   ├── auth_handler.go
│       │   │   │   ├── telemetry_handler.go
│       │   │   │   └── settings_handler.go
│       │   │   ├── middleware/jwt_middleware.go
│       │   │   └── router.go
│       │   └── websocket/
│       │       ├── hub.go          ← MODIFIED v3.0: broadcast system messages
│       │       ├── client.go
│       │       └── handler.go
│       ├── usecase/
│       ├── repository/
│       └── domain/
│
├── frontend/                       ← Next.js 14 — App Router
│   ├── .env.example
│   ├── tailwind.config.ts
│   └── src/
│       ├── app/
│       │   ├── layout.tsx
│       │   ├── page.tsx
│       │   ├── login/page.tsx
│       │   ├── dashboard/
│       │   │   ├── layout.tsx
│       │   │   └── page.tsx
│       │   ├── history/page.tsx
│       │   └── settings/page.tsx
│       │
│       ├── components/
│       │   ├── ui/
│       │   └── features/
│       │       ├── SensorCard.tsx
│       │       ├── StatusChip.tsx
│       │       ├── AlertToast.tsx
│       │       ├── AlertFeed.tsx       ← MODIFIED v3.0: tampilkan AI narrative
│       │       ├── NavSidebar.tsx      ← MODIFIED v3.0: input Groq API Key
│       │       ├── SparklineChart.tsx
│       │       ├── ConnectionStatus.tsx
│       │       ├── DeviceLedIndicator.tsx  ← NEW v3.0: lampu hijau/merah
│       │       └── AIInsightCard.tsx       ← NEW v3.0: on-demand AI insight
│       │
│       ├── hooks/
│       │   ├── useWebSocket.ts     ← MODIFIED v3.0: handle system messages
│       │   ├── useTelemetry.ts
│       │   ├── useAuth.ts
│       │   └── useAutoNarrative.ts ← NEW v3.0: auto Groq call on CRITICAL
│       │
│       ├── lib/
│       │   ├── api.ts
│       │   ├── mock-telemetry.ts
│       │   ├── groq-client.ts      ← NEW v3.0: Groq API wrapper
│       │   └── utils.ts
│       │
│       ├── types/
│       │   ├── telemetry.ts        ← MODIFIED v3.0: tambah SystemMessage, AlertWithNarrative
│       │   ├── auth.ts
│       │   └── api.ts
│       │
│       └── constants/
│           ├── design-tokens.ts
│           └── thresholds.ts
│
└── docs/
    ├── PRD_HEALY_Website.md
    ├── HEALY_Master_Blueprint.md
    └── api-spec.md
```

---

## 3. Design System (Frontend)

Tidak ada perubahan pada design system di v3.0. Semua token warna, typography, animasi, dan komponen patterns dari v2.1.0 tetap berlaku.

Penambahan untuk Phase 10–11:

**Warna baru (tambahkan ke design-tokens.ts):**

```typescript
// Tambahkan ke COLORS di design-tokens.ts
deviceOnline:  '#22C55E',   // Green-500 — LED indicator device online
deviceOffline: '#EF4444',   // Red-500 — LED indicator device offline
aiAccent:      '#7C3AED',   // Violet-700 — aksen warna AI/Groq elements
aiSurface:     '#F5F3FF',   // Violet-50 — background AIInsightCard
```

**Tambahkan ke tailwind.config.ts:**

```typescript
healy: {
  // ... existing colors ...
  'device-on':  '#22C55E',
  'device-off': '#EF4444',
  'ai-accent':  '#7C3AED',
  'ai-surface': '#F5F3FF',
}
```

---

## 4. Backend: Spesifikasi Teknis

Semua spesifikasi dari v2.1.0 tetap berlaku. Perubahan hanya pada `hub.go`.

### 4.1 — 4.5: Tidak berubah dari v2.1.0

### 4.6 Modified: WebSocket Hub — Device Presence Broadcast

`hub.go` dimodifikasi untuk memisahkan dua tipe client: **device clients** (ESP32) dan **viewer clients** (Frontend browser). Saat device client connect atau disconnect, Hub melakukan broadcast JSON khusus bertipe `system` ke semua viewer clients.

```go
// internal/delivery/websocket/hub.go — MODIFIKASI v3.0

package websocket

import "encoding/json"

// SystemMessage adalah pesan non-telemetry dari Hub ke Frontend
type SystemMessage struct {
  Type   string `json:"type"`   // selalu "system"
  Status string `json:"status"` // "device_connected" | "device_disconnected"
}

// Hub kini membedakan dua set client
type Hub struct {
  // Viewer clients = browser/frontend yang menonton data
  viewerClients map[*Client]bool

  // Device clients = ESP32 yang mengirim data
  deviceClients map[*Client]bool

  broadcast        chan []byte  // Data telemetri dari ESP32
  registerViewer   chan *Client
  unregisterViewer chan *Client
  registerDevice   chan *Client
  unregisterDevice chan *Client
}

func NewHub() *Hub {
  return &Hub{
    viewerClients:    make(map[*Client]bool),
    deviceClients:    make(map[*Client]bool),
    broadcast:        make(chan []byte, 256),
    registerViewer:   make(chan *Client),
    unregisterViewer: make(chan *Client),
    registerDevice:   make(chan *Client),
    unregisterDevice: make(chan *Client),
  }
}

func (h *Hub) broadcastSystem(status string) {
  msg := SystemMessage{Type: "system", Status: status}
  data, err := json.Marshal(msg)
  if err != nil {
    return
  }
  // Kirim ke semua viewer clients
  for client := range h.viewerClients {
    select {
    case client.send <- data:
    default:
      close(client.send)
      delete(h.viewerClients, client)
    }
  }
}

func (h *Hub) Run() {
  for {
    select {
    case client := <-h.registerViewer:
      h.viewerClients[client] = true
      // Kirim status device saat ini ke viewer yang baru connect
      if len(h.deviceClients) > 0 {
        h.broadcastSystem("device_connected")
      }

    case client := <-h.unregisterViewer:
      if _, ok := h.viewerClients[client]; ok {
        delete(h.viewerClients, client)
        close(client.send)
      }

    case client := <-h.registerDevice:
      h.deviceClients[client] = true
      // Broadcast ke semua viewer: ESP32 baru saja online
      h.broadcastSystem("device_connected")

    case client := <-h.unregisterDevice:
      if _, ok := h.deviceClients[client]; ok {
        delete(h.deviceClients, client)
        close(client.send)
      }
      // Broadcast ke semua viewer: ESP32 offline
      if len(h.deviceClients) == 0 {
        h.broadcastSystem("device_disconnected")
      }

    case message := <-h.broadcast:
      // Forward telemetry data ke semua viewer clients
      for client := range h.viewerClients {
        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(h.viewerClients, client)
        }
      }
    }
  }
}
```

### 4.7 REST API Endpoints (tidak berubah dari v2.1.0)

Semua endpoint dari Blueprint Section 4.6 v2.1.0 tetap berlaku. Tidak ada endpoint baru karena Phase 10–11 sepenuhnya diimplementasikan di Frontend.

---

## 5. Frontend: Spesifikasi Teknis

### 5.1 Updated TypeScript Interfaces (src/types/telemetry.ts)

```typescript
// Tambahkan ke types/telemetry.ts yang sudah ada

// SystemMessage — pesan non-telemetry dari Hub
export interface SystemMessage {
  type: 'system'
  status: 'device_connected' | 'device_disconnected'
}

// Union type untuk semua pesan WebSocket yang mungkin diterima
export type WebSocketMessage = TelemetryPayload | SystemMessage

// Helper type guard
export function isSystemMessage(msg: WebSocketMessage): msg is SystemMessage {
  return (msg as SystemMessage).type === 'system'
}

// Alert dengan AI narrative (untuk Phase 11)
export interface AlertWithNarrative {
  id: string
  timestamp: Date
  alert_type: string
  value: number
  status: SensorStatus
  device_id: string
  // Narrative dari Groq — null saat sedang di-generate, string saat selesai
  narrative: string | null
  narrativeLoading: boolean
}
```

### 5.2 Modified: useWebSocket.ts — Handle System Messages

```typescript
// src/hooks/useWebSocket.ts — MODIFIKASI v3.0
// Tambahkan state deviceOnline dan handler untuk system messages

import { useEffect, useRef, useState, useCallback } from 'react'
import {
  TelemetryPayload,
  ConnectionState,
  SystemMessage,
  isSystemMessage,
  WebSocketMessage
} from '@/types/telemetry'

const RECONNECT_DELAYS = [1000, 2000, 4000, 8000, 16000, 30000]

export function useWebSocket(url: string) {
  const [data, setData] = useState<TelemetryPayload | null>(null)
  const [deviceOnline, setDeviceOnline] = useState(false)           // NEW v3.0
  const [systemMessage, setSystemMessage] = useState<SystemMessage | null>(null) // NEW v3.0
  const [conn, setConn] = useState<ConnectionState>({
    status: 'DISCONNECTED',
    lastUpdate: null,
    retryCount: 0,
  })
  const wsRef = useRef<WebSocket | null>(null)
  const retryRef = useRef(0)

  const connect = useCallback(() => {
    wsRef.current = new WebSocket(url)

    wsRef.current.onopen = () => {
      retryRef.current = 0
      setConn(prev => ({ ...prev, status: 'CONNECTED', retryCount: 0 }))
    }

    wsRef.current.onmessage = (event) => {
      const msg: WebSocketMessage = JSON.parse(event.data)

      // v3.0: Pisahkan handling system message vs telemetry data
      if (isSystemMessage(msg)) {
        setSystemMessage(msg)
        setDeviceOnline(msg.status === 'device_connected')
        return
      }

      // Normal telemetry payload
      setData(msg)
      setConn(prev => ({ ...prev, lastUpdate: new Date() }))
    }

    wsRef.current.onclose = () => {
      const delay = RECONNECT_DELAYS[Math.min(retryRef.current, RECONNECT_DELAYS.length - 1)]
      retryRef.current++
      setConn(prev => ({ ...prev, status: 'RECONNECTING', retryCount: retryRef.current }))
      setDeviceOnline(false) // v3.0: Device assumed offline saat WS terputus
      setTimeout(connect, delay)
    }
  }, [url])

  useEffect(() => {
    connect()
    return () => wsRef.current?.close()
  }, [connect])

  return { data, conn, deviceOnline, systemMessage } // v3.0: return deviceOnline + systemMessage
}
```

### 5.3 New Component: DeviceLedIndicator.tsx

```typescript
// src/components/features/DeviceLedIndicator.tsx

'use client'
import { motion, AnimatePresence } from 'framer-motion'

interface DeviceLedIndicatorProps {
  online: boolean
  className?: string
}

export function DeviceLedIndicator({ online, className = '' }: DeviceLedIndicatorProps) {
  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <div className="relative flex h-3 w-3">
        {online && (
          <motion.span
            className="absolute inline-flex h-full w-full rounded-full bg-healy-device-on opacity-75"
            animate={{ scale: [1, 1.8, 1], opacity: [0.75, 0, 0.75] }}
            transition={{ duration: 1.5, repeat: Infinity, ease: 'easeInOut' }}
          />
        )}
        <span
          className={`relative inline-flex rounded-full h-3 w-3 transition-colors duration-500 ${
            online ? 'bg-healy-device-on' : 'bg-healy-device-off'
          }`}
        />
      </div>
      <span className="text-xs font-mono text-healy-slate">
        {online ? 'Device Online' : 'Device Offline'}
      </span>
    </div>
  )
}
```

**Cara pakai di Dashboard Header:**
```typescript
// Di dalam dashboard/layout.tsx atau page.tsx
const { deviceOnline, systemMessage } = useTelemetry()

// Di JSX Header
<DeviceLedIndicator online={deviceOnline} />
```

### 5.4 New: System Message Toast Handler

Tambahkan logika ini di `dashboard/page.tsx` atau di `useTelemetry.ts` untuk menampilkan toast otomatis saat device connect/disconnect.

```typescript
// Tambahkan di useTelemetry.ts atau dashboard page

import { useEffect } from 'react'
import { toast } from '@/components/ui/use-toast' // shadcn/ui toast

// Di dalam hook atau komponen yang consume useWebSocket
useEffect(() => {
  if (!systemMessage) return

  if (systemMessage.status === 'device_connected') {
    toast({
      title: 'Robot HEALY Terhubung',
      description: 'Perangkat IoT berhasil terdeteksi dan siap mengirim data.',
      variant: 'default',
      // Custom styling: border sage green
      className: 'border-l-4 border-healy-sage',
    })
  } else if (systemMessage.status === 'device_disconnected') {
    toast({
      title: 'Robot HEALY Terputus',
      description: 'Koneksi ke perangkat hilang. Menunggu reconnect...',
      variant: 'destructive',
      className: 'border-l-4 border-healy-critical',
    })
  }
}, [systemMessage])
```

---

## 6. Database Schema

Tidak ada perubahan skema database untuk Phase 10–11. Groq AI dipanggil langsung dari browser dan tidak menyimpan narrative ke database. Jika di masa depan narrative ingin dipersist, tambahkan kolom `ai_narrative TEXT` pada tabel `public.alert_logs`.

---

## 7. Environment Variables

### 7.1 Backend (backend/.env.example) — Tidak berubah

### 7.2 Frontend (frontend/.env.example) — Perbarui

```env
# WebSocket
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws/telemetry

# API
NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Feature flags
NEXT_PUBLIC_USE_MOCK_DATA=true

# Demo credentials
NEXT_PUBLIC_DEMO_USER=admin
NEXT_PUBLIC_DEMO_PASS=healy123

# Groq AI — JANGAN simpan di .env untuk production
# Pengguna memasukkan key mereka sendiri via UI (localStorage)
# Variable ini hanya untuk development/testing lokal
NEXT_PUBLIC_GROQ_API_KEY_DEV=gsk_xxxxxxxxxxxxxxxxxxxx
```

---

## 8. Deployment Strategy

Tidak berubah dari v2.1.0. Frontend di Vercel, Backend di Railway/Render.

---

## 9. Known Technical Gotchas

Semua Gotcha dari v2.1.0 (01–05) tetap berlaku. Tambahan:

---

### GOTCHA-06: Groq API Key di localStorage — Security Trade-off

**Status:** ACTIVE — diketahui dan diterima untuk konteks academic/demo.

**Masalah:** Menyimpan API key di `localStorage` berarti key tersebut:
- Bisa dilihat oleh siapapun yang membuka browser developer tools.
- Bisa dicuri oleh XSS attack (meskipun risiko ini minim jika app tidak punya celah XSS).
- Tidak ter-enkripsi.

**Kenapa tetap dipilih untuk HEALY:**
Ini adalah academic/demo project dengan single-user access. Tidak ada data pasien nyata yang sensitif. Benefit dari simplisitas (tidak perlu backend key management) jauh lebih besar dari risikonya pada konteks ini.

**Mitigasi yang wajib dilakukan:**
```typescript
// Di groq-client.ts, SELALU validasi format key sebelum menyimpan
export function validateGroqKey(key: string): boolean {
  return /^gsk_[a-zA-Z0-9]{50,}$/.test(key)
}

// Jangan pernah log API key ke console
// Jangan pernah kirim API key ke backend kamu sendiri
```

**Jika HEALY berkembang ke production dengan data pasien nyata:**
Pindahkan Groq API key ke environment variable backend dan buat proxy endpoint `POST /api/ai/insight` di Golang. Frontend kirim data sensor, backend yang panggil Groq.

---

### GOTCHA-07: Groq Rate Limiting

**Status:** ACTIVE — perlu dihandle di frontend.

**Masalah:** Groq Free Tier memiliki rate limit. Phase 11 (Auto-Narrative) dipanggil otomatis setiap kali ada CRITICAL event. Jika pasien berada di kondisi kritis berkelanjutan dan data masuk setiap 2 detik, kamu akan hit rate limit dalam menit.

**Solusi — Deduplikasi di useAutoNarrative.ts:**
```typescript
// Gunakan debounce atau cooldown per alert type
// Jangan generate narrative baru jika CRITICAL event sebelumnya
// belum lebih dari 60 detik yang lalu untuk device yang sama
const NARRATIVE_COOLDOWN_MS = 60_000
```

---

## 10. Phase 10 — AI Insight & Device Presence

### 10.1 Overview

Phase 10 menambahkan dua kemampuan baru yang sepenuhnya berjalan di Frontend:
1. **Device Presence** — deteksi dan notifikasi otomatis saat ESP32 connect/disconnect.
2. **AI Insight On-Demand** — user bisa meminta analisis klinis dari Groq AI kapanpun.

### 10.2 Groq API Client (src/lib/groq-client.ts)

```typescript
// src/lib/groq-client.ts

const GROQ_API_URL = 'https://api.groq.com/openai/v1/chat/completions'
const GROQ_MODEL   = 'llama-3.1-8b-instant'  // Cepat, gratis, cukup untuk narasi medis singkat

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
      temperature: 0.3, // Low temperature untuk konsistensi medis
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
```

### 10.3 New Component: AIInsightCard.tsx

```typescript
// src/components/features/AIInsightCard.tsx

'use client'
import { useState } from 'react'
import { Sparkle, Robot, Warning } from '@phosphor-icons/react'
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
      setError('API Key Groq belum dikonfigurasi atau tidak valid. Masukkan key di menu Settings.')
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
        setInsight(streamedText) // Update teks seiring streaming
      })
    } catch (e: any) {
      setError(e.message || 'Gagal menghubungi Groq AI. Periksa koneksi dan API key.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bg-healy-ai-surface border border-healy-border rounded-card p-5 shadow-card">
      {/* Header */}
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Robot size={20} weight="duotone" className="text-healy-ai-accent" />
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
                     px-3 py-1.5 rounded-lg hover:bg-violet-600 disabled:opacity-50
                     disabled:cursor-not-allowed transition-all duration-200"
        >
          <Sparkle size={14} weight="fill" />
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
            <Warning size={16} className="mt-0.5 shrink-0" />
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
```

### 10.4 Modified: NavSidebar.tsx — Groq API Key Input

Tambahkan section ini di bagian bawah NavSidebar, sebelum logout button:

```typescript
// Tambahkan di NavSidebar.tsx — section Groq API Key

import { useState, useEffect } from 'react'
import { Key, Eye, EyeSlash, CheckCircle } from '@phosphor-icons/react'
import { getStoredGroqKey, setStoredGroqKey, validateGroqKey } from '@/lib/groq-client'

function GroqKeySection() {
  const [key, setKey] = useState('')
  const [show, setShow] = useState(false)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    const stored = getStoredGroqKey()
    if (stored) setKey(stored)
  }, [])

  const handleSave = () => {
    if (!validateGroqKey(key)) return
    setStoredGroqKey(key.trim())
    setSaved(true)
    setTimeout(() => setSaved(false), 2000)
  }

  const isValid = validateGroqKey(key)

  return (
    <div className="p-3 border-t border-healy-border">
      <div className="flex items-center gap-1.5 mb-2">
        <Key size={12} className="text-healy-ai-accent" />
        <span className="text-xs font-medium text-healy-slate">Groq API Key</span>
      </div>

      <div className="flex gap-1.5">
        <div className="relative flex-1">
          <input
            type={show ? 'text' : 'password'}
            value={key}
            onChange={e => setKey(e.target.value)}
            placeholder="gsk_..."
            className="w-full text-xs font-mono px-2 py-1.5 pr-7 rounded-lg
                       border border-healy-border bg-healy-surface
                       focus:outline-none focus:ring-2 focus:ring-healy-ai-accent/30
                       focus:border-healy-ai-accent"
          />
          <button
            type="button"
            onClick={() => setShow(!show)}
            className="absolute right-2 top-1/2 -translate-y-1/2 text-healy-slate"
          >
            {show ? <EyeSlash size={12} /> : <Eye size={12} />}
          </button>
        </div>

        <button
          onClick={handleSave}
          disabled={!isValid}
          className="px-2 py-1.5 text-xs rounded-lg bg-healy-ai-accent text-white
                     disabled:opacity-40 disabled:cursor-not-allowed
                     hover:bg-violet-600 transition-colors"
        >
          {saved ? <CheckCircle size={14} weight="fill" /> : 'Simpan'}
        </button>
      </div>

      {key && !isValid && (
        <p className="text-xs text-healy-critical mt-1">Format key tidak valid</p>
      )}
    </div>
  )
}
```

---

## 11. Phase 11 — Auto-Narrative Alert (Inovasi)

### 11.1 Konsep & Nilai Tambah

Phase 10 memberikan AI insight **saat diminta** (on-demand). Phase 11 membaliknya: AI bekerja **secara proaktif** tanpa intervensi user.

Setiap kali sistem mendeteksi status `CRITICAL` baru, dalam background, sistem akan:
1. Mendeteksi bahwa ini adalah CRITICAL event baru (bukan duplikat dalam 60 detik).
2. Memanggil Groq API secara senyap di background.
3. Men-generate narasi medis 1–2 kalimat dalam Bahasa Indonesia.
4. Men-attach narasi tersebut ke alert entry di AlertFeed.

Hasil akhir: Setiap alert di feed bukan hanya berisi chip "KRITIS" — tapi juga penjelasan kontekstual seperti:

> *"SpO2 pasien turun ke 87%, jauh di bawah ambang aman 95%. Kondisi ini dapat mengindikasikan hipoksemia akut — perlu pemeriksaan pernapasan segera."*

Ini secara dramatis meningkatkan nilai informatif dashboard, terutama untuk keluarga pasien yang tidak familiar dengan angka medis.

### 11.2 New Hook: useAutoNarrative.ts

```typescript
// src/hooks/useAutoNarrative.ts

import { useEffect, useRef, useCallback } from 'react'
import { TelemetryPayload, AlertWithNarrative } from '@/types/telemetry'
import {
  getStoredGroqKey,
  validateGroqKey,
  callGroqInsight
} from '@/lib/groq-client'

const NARRATIVE_COOLDOWN_MS = 60_000 // 60 detik antar auto-narrative

interface UseAutoNarrativeOptions {
  onNarrativeReady: (alert: AlertWithNarrative) => void
}

export function useAutoNarrative({ onNarrativeReady }: UseAutoNarrativeOptions) {
  const lastCriticalRef = useRef<number>(0)
  const previousStatusRef = useRef<string>('')

  const buildNarrativePrompt = (data: TelemetryPayload): string => {
    return `Kamu adalah asisten medis AI untuk sistem HEALY.

Kondisi darurat terdeteksi:
- Suhu: ${data.sensor.temperature}°C
- BPM: ${data.sensor.bpm}
- SpO2: ${data.sensor.spo2}%
- Status suhu: ${data.status.temperature}
- Status SpO2: ${data.status.spo2}

Tulis 1-2 kalimat narasi medis darurat dalam Bahasa Indonesia.
Format: interpretasi kondisi + potensi risiko klinis.
Singkat, jelas, mudah dipahami non-medis. Tanpa pengantar.`
  }

  const triggerNarrative = useCallback(async (data: TelemetryPayload) => {
    const now = Date.now()
    const timeSinceLast = now - lastCriticalRef.current

    // Cooldown check — jangan spam Groq API
    if (timeSinceLast < NARRATIVE_COOLDOWN_MS) return

    const apiKey = getStoredGroqKey()
    if (!apiKey || !validateGroqKey(apiKey)) return // Silently skip jika key belum diset

    lastCriticalRef.current = now

    // Buat alert object dulu dengan narrative = null (loading state)
    const newAlert: AlertWithNarrative = {
      id: `alert-${now}`,
      timestamp: new Date(data.timestamp),
      alert_type: `CRITICAL_${data.status.temperature === 'CRITICAL' ? 'TEMP' : 'SPO2'}`,
      value: data.status.temperature === 'CRITICAL'
        ? data.sensor.temperature
        : data.sensor.spo2,
      status: 'CRITICAL',
      device_id: data.device_id,
      narrative: null,       // Belum selesai
      narrativeLoading: true,
    }

    // Immediately notify AlertFeed dengan loading state
    onNarrativeReady(newAlert)

    // Generate narrative di background
    try {
      const prompt = buildNarrativePrompt(data)
      const narrative = await callGroqInsight(prompt, apiKey)
      onNarrativeReady({ ...newAlert, narrative, narrativeLoading: false })
    } catch {
      onNarrativeReady({ ...newAlert, narrative: null, narrativeLoading: false })
    }
  }, [onNarrativeReady])

  const processPayload = useCallback((data: TelemetryPayload | null) => {
    if (!data) return

    const currentStatus = data.status.overall
    const wasNotCritical = previousStatusRef.current !== 'CRITICAL'

    // Hanya trigger saat transisi KE CRITICAL (bukan setiap tick saat sudah critical)
    if (currentStatus === 'CRITICAL' && wasNotCritical) {
      triggerNarrative(data)
    }

    previousStatusRef.current = currentStatus
  }, [triggerNarrative])

  return { processPayload }
}
```

### 11.3 Modified: AlertFeed.tsx — Tampilkan Narrative

```typescript
// src/components/features/AlertFeed.tsx — MODIFIKASI v3.0

import { AlertWithNarrative } from '@/types/telemetry'
import { Warning, Spinner } from '@phosphor-icons/react'
import { motion, AnimatePresence } from 'framer-motion'
import { StatusChip } from './StatusChip'

interface AlertFeedProps {
  alerts: AlertWithNarrative[]
}

export function AlertFeed({ alerts }: AlertFeedProps) {
  if (alerts.length === 0) {
    return (
      <div className="text-center text-healy-slate text-sm py-6">
        Tidak ada alert terbaru
      </div>
    )
  }

  return (
    <div className="space-y-3">
      <AnimatePresence>
        {alerts.slice(0, 5).map((alert) => (
          <motion.div
            key={alert.id}
            initial={{ opacity: 0, x: -12 }}
            animate={{ opacity: 1, x: 0 }}
            exit={{ opacity: 0, x: 12 }}
            className={`rounded-card p-3 border ${
              alert.status === 'CRITICAL'
                ? 'bg-red-50 border-healy-critical/30'
                : 'bg-healy-bg-alt border-healy-border'
            }`}
          >
            {/* Alert Header */}
            <div className="flex items-center justify-between mb-1.5">
              <div className="flex items-center gap-2">
                <Warning
                  size={14}
                  weight="fill"
                  className={alert.status === 'CRITICAL' ? 'text-healy-critical' : 'text-healy-warning'}
                />
                <span className="text-xs font-mono text-healy-slate">
                  {new Date(alert.timestamp).toLocaleTimeString('id-ID')}
                </span>
                <StatusChip status={alert.status} size="sm" />
              </div>
              <span className="text-xs font-mono font-semibold text-healy-graphite">
                {alert.value}
                {alert.alert_type.includes('TEMP') ? '°C' : '%'}
              </span>
            </div>

            {/* AI Narrative — Phase 11 */}
            {alert.status === 'CRITICAL' && (
              <div className="mt-2 pt-2 border-t border-healy-critical/20">
                {alert.narrativeLoading ? (
                  <div className="flex items-center gap-1.5 text-xs text-healy-slate">
                    <Spinner size={12} className="animate-spin text-healy-ai-accent" />
                    <span className="italic">AI sedang menganalisis kondisi...</span>
                  </div>
                ) : alert.narrative ? (
                  <p className="text-xs text-healy-graphite leading-relaxed font-body italic">
                    "{alert.narrative}"
                  </p>
                ) : null}
              </div>
            )}
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  )
}
```

### 11.4 Integrasi di Dashboard Page

```typescript
// Tambahkan di src/app/dashboard/page.tsx

import { useState } from 'react'
import { useAutoNarrative } from '@/hooks/useAutoNarrative'
import { AlertWithNarrative } from '@/types/telemetry'
import { AlertFeed } from '@/components/features/AlertFeed'
import { AIInsightCard } from '@/components/features/AIInsightCard'
import { DeviceLedIndicator } from '@/components/features/DeviceLedIndicator'

export default function DashboardPage() {
  const { data, conn, deviceOnline } = useTelemetry()
  const [alerts, setAlerts] = useState<AlertWithNarrative[]>([])

  // Phase 11: Auto-narrative hook
  const { processPayload } = useAutoNarrative({
    onNarrativeReady: (alert) => {
      setAlerts(prev => {
        // Update existing alert jika sudah ada (narrative selesai di-generate)
        const exists = prev.find(a => a.id === alert.id)
        if (exists) {
          return prev.map(a => a.id === alert.id ? alert : a)
        }
        // Tambah alert baru di depan, max 20 item
        return [alert, ...prev].slice(0, 20)
      })
    }
  })

  // Proses setiap payload baru untuk deteksi CRITICAL
  useEffect(() => {
    processPayload(data)
  }, [data, processPayload])

  // Phase 10: Device presence toast
  useEffect(() => {
    if (!systemMessage) return
    // ... toast logic dari Section 5.4
  }, [systemMessage])

  return (
    <div className="p-6 space-y-6">
      {/* Header dengan Device LED — Phase 10 */}
      <div className="flex items-center justify-between">
        <h1 className="font-display font-bold text-2xl text-healy-graphite">Dashboard</h1>
        <DeviceLedIndicator online={deviceOnline} />
      </div>

      {/* Sensor Cards (tidak berubah) */}
      <div className="grid grid-cols-3 gap-4">
        {/* ... SensorCard components ... */}
      </div>

      {/* Row bawah */}
      <div className="grid grid-cols-2 gap-4">
        {/* Alert Feed dengan Auto-Narrative — Phase 11 */}
        <AlertFeed alerts={alerts} />

        {/* AI Insight On-Demand — Phase 10 */}
        <AIInsightCard
          currentData={data ? {
            temperature: data.sensor.temperature,
            bpm: data.sensor.bpm,
            spo2: data.sensor.spo2,
            status: data.status.overall,
          } : null}
        />
      </div>
    </div>
  )
}
```

---

## 12. Rencana Kerja Step-by-Step

### FASE 0–9 (DONE ✅)

Semua fase dari Blueprint v2.1.0 sudah selesai. Detail ada di progress tracker.

---

### FASE 10 — AI Insight & Device Presence (AKTIF)

**Step 10.1 — Backend: Modifikasi hub.go**
Implementasikan hub.go baru sesuai Section 4.6 dengan pemisahan `viewerClients` dan `deviceClients`. Tambahkan fungsi `broadcastSystem`. Update `handler.go` untuk menggunakan `registerViewer` dan `registerDevice` channel yang tepat berdasarkan path atau header request.

**Step 10.2 — Frontend: Update Types**
Tambahkan `SystemMessage`, `WebSocketMessage`, `isSystemMessage`, dan `AlertWithNarrative` ke `src/types/telemetry.ts`.

**Step 10.3 — Frontend: Modifikasi useWebSocket.ts**
Tambahkan handling untuk `system` messages. Return `deviceOnline` dan `systemMessage` dari hook.

**Step 10.4 — Frontend: DeviceLedIndicator.tsx**
Buat komponen baru sesuai Section 5.3. Pasang di Dashboard header.

**Step 10.5 — Frontend: Toast untuk Device Presence**
Tambahkan `useEffect` di dashboard yang watch `systemMessage` dan trigger toast shadcn/ui.

**Step 10.6 — Frontend: groq-client.ts**
Buat file baru sesuai Section 10.2 lengkap dengan fungsi `getStoredGroqKey`, `setStoredGroqKey`, `validateGroqKey`, `buildInsightPrompt`, dan `callGroqInsight` dengan streaming.

**Step 10.7 — Frontend: Modifikasi NavSidebar.tsx**
Tambahkan `GroqKeySection` component di bagian bawah sidebar.

**Step 10.8 — Frontend: AIInsightCard.tsx**
Buat komponen baru sesuai Section 10.3. Pasang di dashboard page di samping AlertFeed.

---

### FASE 11 — Auto-Narrative Alert (INOVASI)

**Step 11.1 — Frontend: useAutoNarrative.ts**
Buat hook baru sesuai Section 11.2. Perhatikan logic cooldown 60 detik dan deteksi transisi ke CRITICAL.

**Step 11.2 — Frontend: Modifikasi AlertFeed.tsx**
Update komponen untuk menerima `AlertWithNarrative[]` dan menampilkan narrative loading state + text.

**Step 11.3 — Frontend: Integrasi di dashboard/page.tsx**
Hubungkan `useAutoNarrative`, state `alerts`, dan komponen `AlertFeed` sesuai Section 11.4.

**Step 11.4 — Testing Skenario**
- Set Groq API Key di sidebar.
- Tunggu mock generator menghasilkan CRITICAL event.
- Verifikasi: alert muncul di feed dengan loading spinner, lalu diganti narrative dari Groq.
- Tunggu 60 detik, trigger CRITICAL lagi. Verifikasi narrative baru di-generate.
- Trigger CRITICAL dua kali dalam 60 detik. Verifikasi hanya satu narrative yang di-generate (cooldown bekerja).

---

## 13. Codebase Progress Tracker

| Layer | File | Status | Catatan |
|---|---|---|---|
| Domain | `internal/domain/telemetry.go` | ✅ DONE | Synced dengan ESP32 payload |
| Domain | `internal/domain/alert.go` | ✅ DONE | |
| Domain | `internal/domain/user.go` | ✅ DONE | |
| Domain | `internal/domain/settings.go` | ✅ DONE | |
| Repository | `internal/repository/postgres/db.go` | ✅ DONE | pgx/v5, port 6543 |
| Repository | `internal/repository/interfaces/telemetry_repo.go` | ✅ DONE | |
| Repository | `internal/repository/postgres/telemetry_postgres.go` | ✅ DONE | |
| Repository | `internal/repository/interfaces/settings_repo.go` | ✅ DONE | |
| Repository | `internal/repository/postgres/settings_postgres.go` | ✅ DONE | ON CONFLICT upsert |
| Usecase | `internal/usecase/alert_usecase.go` | ✅ DONE | Threshold engine |
| Usecase | `internal/usecase/telemetry_usecase.go` | ✅ DONE | Orchestrator |
| Usecase | `internal/usecase/auth_usecase.go` | ✅ DONE | |
| Delivery | `internal/delivery/websocket/hub.go` | ⬜ TODO | **MODIFY v3.0**: pisah viewerClients/deviceClients, broadcastSystem |
| Delivery | `internal/delivery/websocket/handler.go` | ⬜ TODO | **MODIFY v3.0**: route ke registerViewer atau registerDevice |
| Delivery | `internal/delivery/http/router.go` | ✅ DONE | |
| Delivery | `internal/delivery/http/middleware/jwt_middleware.go` | ✅ DONE | |
| Delivery | `internal/delivery/http/telemetry_handler.go` | ✅ DONE | |
| Delivery | `internal/delivery/http/settings_handler.go` | ✅ DONE | |
| Config | `pkg/config/config.go` | ✅ DONE | |
| Frontend | `tailwind.config.ts` | ⬜ TODO | **MODIFY v3.0**: tambah healy-device-on/off, healy-ai-accent/surface |
| Frontend | `src/constants/design-tokens.ts` | ⬜ TODO | **MODIFY v3.0**: tambah 4 warna baru |
| Frontend | `src/types/telemetry.ts` | ⬜ TODO | **MODIFY v3.0**: tambah SystemMessage, AlertWithNarrative |
| Frontend | `src/hooks/useWebSocket.ts` | ⬜ TODO | **MODIFY v3.0**: handle system messages, return deviceOnline |
| Frontend | `src/hooks/useAutoNarrative.ts` | ⬜ TODO | **NEW v3.0** Phase 11 |
| Frontend | `src/lib/groq-client.ts` | ⬜ TODO | **NEW v3.0** Phase 10 |
| Frontend | `src/components/features/DeviceLedIndicator.tsx` | ⬜ TODO | **NEW v3.0** Phase 10 |
| Frontend | `src/components/features/AIInsightCard.tsx` | ⬜ TODO | **NEW v3.0** Phase 10 |
| Frontend | `src/components/features/AlertFeed.tsx` | ⬜ TODO | **MODIFY v3.0**: tampilkan AI narrative |
| Frontend | `src/components/features/NavSidebar.tsx` | ⬜ TODO | **MODIFY v3.0**: tambah GroqKeySection |
| Frontend | `src/app/dashboard/page.tsx` | ⬜ TODO | **MODIFY v3.0**: integrasi Phase 10 & 11 |
| Frontend | Design tokens + Tailwind | ✅ DONE | Base design system v2.1 |
| Frontend | SensorCard, StatusChip, SparklineChart | ✅ DONE | |
| Frontend | Dashboard, Landing, Login, History, Settings | ✅ DONE | |
| Frontend | History page (Recharts), a11y retrofit | ✅ DONE | |
| Deployment | CORS & Env Preparation | ✅ DONE | Railway + Vercel ready |
| QA | Phase 9: Final QA & SEO Sweep | ✅ DONE | Metadata, OpenGraph dikonfigurasi |

---

## 14. Konvensi Kode & Git

### 14.1 Branch Strategy

```
main              ← Production-ready only
develop           ← Integration branch
feat/phase10-device-presence
feat/phase10-ai-insight-card
feat/phase11-auto-narrative
fix/ws-hub-viewer-device-split
```

### 14.2 Commit Message Format

```
feat(hub): split viewer/device clients + broadcast system messages
feat(websocket): handle system messages + deviceOnline state in useWebSocket
feat(ui): add DeviceLedIndicator with pulse animation
feat(groq): add groq-client with streaming support and key validation
feat(ai): add AIInsightCard with on-demand Groq insight generation
feat(alert): add useAutoNarrative hook with 60s cooldown
feat(alert-feed): display AI narrative below CRITICAL alerts
fix(narrative): deduplicate CRITICAL event triggers using transition detection
```

### 14.3 Naming Conventions (tidak berubah dari v2.1.0)

---

## 15. Prompt Inisialisasi AI (Master Prompt)

### Prompt 0 — Grounding (Wajib di awal setiap session)

```
You are a Senior Full-Stack Engineer on Project HEALY — Health Observer Robot.

I am attaching:
1. PRD_HEALY_Website.md
2. HEALY_Master_Blueprint.md (v3.0.0)

Read both documents completely. Confirm understanding by stating:
1. Tech stack: IoT / Backend / Frontend / Database / Deployment
2. Phase 10: What are the two capabilities? How does Groq key get stored?
3. Phase 11: What triggers Auto-Narrative? What is the cooldown duration?
4. Which files are MODIFY vs NEW in v3.0 Progress Tracker?
5. What is GOTCHA-06 and GOTCHA-07?

Do NOT generate any code yet. Confirmation only.

[ATTACH: PRD_HEALY_Website.md]
[ATTACH: HEALY_Master_Blueprint.md]
```

---

### Prompt Phase 10 — Device Presence & Groq Setup

```
Based on HEALY Master Blueprint v3.0 Section 4.6 and Section 10, 
implement Phase 10: AI Insight & Device Presence.

I will run this in stages. Start with Stage A only.

STAGE A — Backend hub.go modification:
File: internal/delivery/websocket/hub.go

Requirements from Blueprint Section 4.6:
1. Separate viewerClients (map[*Client]bool) from deviceClients (map[*Client]bool)
2. Add 4 channels: registerViewer, unregisterViewer, registerDevice, unregisterDevice
3. Add broadcastSystem(status string) function that marshals SystemMessage{type:"system"} 
   and sends to ALL viewerClients
4. In Run() loop, handle registerDevice: call broadcastSystem("device_connected")
5. In Run() loop, handle unregisterDevice: if deviceClients is empty, 
   call broadcastSystem("device_disconnected")
6. In Run() loop, handle registerViewer: if deviceClients is not empty, 
   immediately send device_connected to just this new client

Also update: internal/delivery/websocket/handler.go
- /ws/device path → use hub.registerDevice
- /ws/telemetry path → use hub.registerViewer

Output: Both files full implementation. No placeholder comments.
```

---

### Prompt Phase 10 — Frontend Groq Integration

```
Based on HEALY Master Blueprint v3.0 Sections 5.1–5.4 and 10.2–10.4,
implement the frontend side of Phase 10.

Implement in order:

Step 1: Update src/types/telemetry.ts
Add SystemMessage, WebSocketMessage union type, and isSystemMessage type guard
exactly as in Blueprint Section 5.1. Also add AlertWithNarrative interface.

Step 2: Update src/hooks/useWebSocket.ts  
Modify to handle system messages. Add deviceOnline and systemMessage to state.
Use isSystemMessage() type guard to route messages. Blueprint Section 5.2.

Step 3: Create src/lib/groq-client.ts
Full implementation from Blueprint Section 10.2:
- getStoredGroqKey / setStoredGroqKey (localStorage key: 'healy_groq_api_key')
- validateGroqKey (regex: /^gsk_[a-zA-Z0-9]{50,}$/)
- buildInsightPrompt (Bahasa Indonesia medical prompt)
- callGroqInsight (supports both streaming and non-streaming via onChunk callback)
Model: llama-3.1-8b-instant, max_tokens: 200, temperature: 0.3

Step 4: Create src/components/features/DeviceLedIndicator.tsx
Blueprint Section 5.3. Use framer-motion for pulse animation on online state.
Use healy-device-on (#22C55E) and healy-device-off (#EF4444) colors.

Step 5: Create src/components/features/AIInsightCard.tsx
Full implementation from Blueprint Section 10.3.
Use streaming mode — update text in real-time as chunks arrive.

Output: All 5 files with complete implementation.
Apply HEALY design tokens throughout. No hardcoded colors.
```

---

### Prompt Phase 11 — Auto-Narrative Alert

```
Based on HEALY Master Blueprint v3.0 Section 11, implement Phase 11.

Prerequisite: groq-client.ts from Phase 10 is already in place.

Step 1: Create src/hooks/useAutoNarrative.ts
Full implementation from Blueprint Section 11.2:
- Track previous status to detect TRANSITION to CRITICAL (not every CRITICAL tick)
- Cooldown: NARRATIVE_COOLDOWN_MS = 60_000
- On CRITICAL transition: immediately call onNarrativeReady with narrativeLoading: true
- Then call callGroqInsight() in background
- On completion: call onNarrativeReady again with narrative text + narrativeLoading: false
- On Groq error: call onNarrativeReady with narrative: null + narrativeLoading: false
- Use Bahasa Indonesia medical prompt from Section 11.2 buildNarrativePrompt

Step 2: Update src/components/features/AlertFeed.tsx
Modify to accept AlertWithNarrative[] instead of old alert type.
Show narrative section below CRITICAL alerts:
- Loading state: Spinner icon + "AI sedang menganalisis kondisi..."
- Loaded state: italic text in quotes
- null narrative (error): render nothing (fail silently)
Use framer-motion AnimatePresence for alert entry animation.

Step 3: Update src/app/dashboard/page.tsx
Wire everything together per Blueprint Section 11.4:
- useState for alerts (AlertWithNarrative[])
- useAutoNarrative hook with onNarrativeReady callback
- useEffect that calls processPayload(data) when data changes
- Pass alerts to AlertFeed component
- Also wire DeviceLedIndicator to deviceOnline from useTelemetry
- Toast for systemMessage device events

Output: Three files. Full implementation. Test by checking mock generator 
triggers CRITICAL event and narrative appears in AlertFeed.
```

---

*HEALY Master Blueprint v3.0.0 — Updated by Lead Developer.*
*Phase 10: AI Insight & Device Presence.*
*Phase 11: Auto-Narrative Alert — inovasi proaktif yang membedakan HEALY dari monitoring system biasa.*
*Dokumen ini adalah satu-satunya master reference yang berlaku.*
*v2.1.0 dinyatakan deprecated.*
