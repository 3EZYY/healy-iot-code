# HEALY — Master Blueprint
## Health Observer Robot: Full-Stack Implementation Guide

---

| Field | Detail |
|---|---|
| **Document Version** | v2.1.0 |
| **Status** | Active — Master Reference |
| **Role Penyusun** | Lead Developer / Senior Full-Stack Engineer |
| **Stack** | ESP32 (IoT) + Golang 1.23+ (Backend) + Next.js 14 (Frontend) |
| **Design Philosophy** | Clinical Futurism — Sage Green, Glacial White, Futuristic Medical |
| **Last Updated** | 2025 |
| **Changelog** | v2.1.0: Deployment strategy update, DB migration TimescaleDB → Supabase pg_partman, schema telemetry fix, codebase progress update |

---

## Changelog v2.0.0 → v2.1.0

Dokumen ini diperbarui berdasarkan temuan teknis dari sesi development aktif. Semua perubahan bersifat **breaking** terhadap keputusan di v2.0.0 dan wajib diterapkan segera.

| # | Area | Sebelumnya (v2.0.0) | Sesudahnya (v2.1.0) | Urgensi |
|---|---|---|---|---|
| C1 | Deployment Backend | Tidak dispesifikasikan | Railway atau Render (bukan Vercel) | CRITICAL |
| C2 | Database Engine | TimescaleDB (self-hosted) | Supabase PostgreSQL + pg_partman + pg_cron | CRITICAL |
| C3 | DB Schema | Semua tabel di `public` | Tabel telemetry di schema `telemetry` | HIGH |
| C4 | DB Table Reference | `telemetry_records` | `telemetry.telemetry_records` (schema-qualified) | HIGH |
| C5 | DB Connection | lib/pq, port 5432 | pgx/v5, port 6543 (Transaction Pooler) | HIGH |
| C6 | docker-compose | Diperlukan untuk local DB | Tidak diperlukan (Supabase online) | LOW |
| C7 | Data Retention | Tidak ada policy | Wajib ada retention 30 hari (Supabase Free Tier 500MB) | MEDIUM |
| C8 | Codebase Progress | Semua layer belum dimulai | Domain ✅ Repository ✅ Usecase ⬜ Delivery ⬜ | INFO |

---

## Catatan Lead Dev: Perbaikan dari Saran Sebelumnya

**Yang sudah benar (tetap berlaku):**
- Clean Architecture untuk Golang (cmd, internal, pkg) — tepat.
- Next.js 14 App Router + Tailwind + Shadcn/UI — tepat.
- Monorepo dengan 3 folder terpisah (iot, backend, frontend) — tepat.
- Pemisahan interface repository dari implementasi — sudah benar.

**Perubahan kritikal v2.1.0 (lihat detail di section masing-masing):**
1. Backend TIDAK BOLEH di-deploy ke Vercel — Vercel adalah serverless, WebSocket persistent tidak bisa berjalan di sana.
2. Database diganti dari TimescaleDB lokal ke Supabase cloud dengan strategi partisi menggunakan `pg_partman`.
3. Semua query ke tabel telemetry WAJIB menggunakan nama qualified: `telemetry.telemetry_records`.
4. Driver database diganti dari `lib/pq` ke `pgx/v5` untuk kompatibilitas Supabase Transaction Pooler di port 6543.

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
10. [Rencana Kerja Step-by-Step](#10-rencana-kerja-step-by-step)
11. [Codebase Progress Tracker](#11-codebase-progress-tracker)
12. [Konvensi Kode & Git](#12-konvensi-kode--git)
13. [Prompt Inisialisasi AI (Master Prompt)](#13-prompt-inisialisasi-ai-master-prompt)

---

## 1. Arsitektur Sistem Keseluruhan

```
┌──────────────────────────────────────────────────────────────────────┐
│                          HEALY SYSTEM v2.1                           │
│                                                                      │
│  ┌──────────┐    I2C/I2S    ┌──────────────────────────────────┐    │
│  │  Sensor  │ ────────────▶ │        ESP32 WROOM-32            │    │
│  │ MLX90614 │               │    (IoT Layer / /iot folder)     │    │
│  │ MAX30102  │               │  • Baca sensor (I2C/I2S)        │    │
│  │ INMP441  │               │  • Encode JSON payload           │    │
│  └──────────┘               │  • Kirim via WiFi + mDNS         │    │
│                             └────────────────┬─────────────────┘    │
│                                              │ WebSocket (JSON)      │
│                                              ▼                       │
│                      ┌───────────────────────────────────────┐      │
│                      │     Golang Backend                    │      │
│                      │  Railway / Render (BUKAN Vercel)      │      │
│                      │                                       │      │
│                      │  • WS Server — gorilla/websocket      │      │
│                      │  • Usecase: Threshold & Alert Engine  │      │
│                      │  • REST API (Gin)                     │      │
│                      │  • Write to Supabase PostgreSQL       │      │
│                      └──────────────┬────────────────────────┘      │
│                                     │                                │
│              ┌──────────────────────┼──────────────────┐            │
│              │ WebSocket (WSS)      │ REST (/api/...)  │            │
│              ▼                      ▼                  ▼            │
│   ┌──────────────────┐   ┌────────────────────────────────────┐    │
│   │  Next.js 14      │   │     Supabase PostgreSQL (Cloud)    │    │
│   │  Vercel          │   │                                    │    │
│   │                  │   │  Schema: telemetry                 │    │
│   │  • Landing       │   │  • telemetry.telemetry_records     │    │
│   │  • Dashboard     │   │    (Partitioned daily, pg_partman) │    │
│   │  • History       │   │  Schema: public                    │    │
│   │  • Settings      │   │  • public.users                    │    │
│   └──────────────────┘   │  • public.alert_logs               │    │
│                          │  • public.device_settings          │    │
│                          └────────────────────────────────────┘    │
└──────────────────────────────────────────────────────────────────────┘
```

**Aliran Data Utama:**
```
ESP32 ──(WS JSON)──▶ Golang WS Handler ──▶ Usecase (threshold check)
                 ──▶ Repository (write telemetry.telemetry_records)
                 ──▶ Broadcast ke semua Frontend WS clients
```

---

## 2. Struktur Folder Final (Revised)

```
HEALY-PROJECT/
│
├── .env.example                    ← Referensi env var root (tidak ada docker-compose)
├── README.md
│
├── iot/                            ← (EXISTING: ESP32 C++/PlatformIO)
│   ├── src/
│   │   └── main.cpp
│   ├── include/
│   ├── lib/
│   └── platformio.ini
│
├── backend/                        ← Golang — Clean Architecture
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   │
│   ├── cmd/
│   │   └── api/
│   │       └── main.go             ← Entry point: init server, DB, WS
│   │
│   ├── internal/
│   │   ├── delivery/               ← Transport layer (HTTP + WS)
│   │   │   ├── http/
│   │   │   │   ├── handler/
│   │   │   │   │   ├── auth_handler.go
│   │   │   │   │   ├── telemetry_handler.go
│   │   │   │   │   └── settings_handler.go
│   │   │   │   ├── middleware/
│   │   │   │   │   └── jwt_middleware.go
│   │   │   │   └── router.go
│   │   │   └── websocket/
│   │   │       ├── hub.go
│   │   │       ├── client.go
│   │   │       └── handler.go
│   │   │
│   │   ├── usecase/                ← Business logic (FASE AKTIF SEKARANG)
│   │   │   ├── telemetry_usecase.go
│   │   │   ├── alert_usecase.go    ← Threshold engine
│   │   │   └── auth_usecase.go
│   │   │
│   │   ├── repository/             ← Data access layer
│   │   │   ├── interfaces/         ← Interface definitions
│   │   │   │   ├── telemetry_repo.go
│   │   │   │   └── user_repo.go
│   │   │   └── postgres/           ← Implementasi konkret
│   │   │       ├── db.go           ← (DONE ✅) Connection pool pgx/v5
│   │   │       ├── telemetry_postgres.go
│   │   │       └── user_postgres.go
│   │   │
│   │   └── domain/                 ← (DONE ✅) Entities & DTOs
│   │       ├── telemetry.go
│   │       ├── alert.go
│   │       ├── user.go
│   │       └── settings.go
│   │
│   └── pkg/
│       ├── config/
│       │   └── config.go
│       ├── database/
│       │   └── postgres.go         ← pgx/v5 pool (port 6543)
│       ├── jwt/
│       │   └── jwt.go
│       └── logger/
│           └── logger.go
│
├── frontend/                       ← Next.js 14 — App Router
│   ├── .env.example
│   ├── package.json
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   │
│   └── src/
│       ├── app/
│       │   ├── layout.tsx
│       │   ├── page.tsx            ← Landing page
│       │   ├── login/
│       │   │   └── page.tsx
│       │   ├── dashboard/
│       │   │   ├── layout.tsx
│       │   │   └── page.tsx
│       │   ├── history/
│       │   │   └── page.tsx
│       │   └── settings/
│       │       └── page.tsx
│       │
│       ├── components/
│       │   ├── ui/                 ← Atoms (shadcn/ui base)
│       │   └── features/           ← Organisms (domain-specific)
│       │       ├── SensorCard.tsx
│       │       ├── StatusChip.tsx
│       │       ├── AlertToast.tsx
│       │       ├── AlertFeed.tsx
│       │       ├── NavSidebar.tsx
│       │       ├── SparklineChart.tsx
│       │       └── ConnectionStatus.tsx
│       │
│       ├── hooks/
│       │   ├── useWebSocket.ts
│       │   ├── useTelemetry.ts
│       │   └── useAuth.ts
│       │
│       ├── lib/
│       │   ├── api.ts
│       │   ├── mock-telemetry.ts
│       │   └── utils.ts
│       │
│       ├── types/
│       │   ├── telemetry.ts
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

Source of truth desain. Semua nilai harus diimplementasikan sebagai CSS Variables dan Tailwind config. Tidak boleh ada magic number hex di dalam komponen.

### 3.1 Color Tokens

```typescript
// src/constants/design-tokens.ts

export const COLORS = {
  bgPrimary:     '#F8FAFB',  // Glacial White
  bgSecondary:   '#EFF4F2',  // Soft Mist
  surface:       '#FFFFFF',  // Ice White

  brandPrimary:  '#4CAF82',  // Sage Green
  brandSecondary:'#2E8B62',  // Deep Sage
  brandAccent:   '#A8DFCC',  // Mint Glow

  textPrimary:   '#1A2633',  // Graphite
  textSecondary: '#5A7080',  // Slate

  statusNormal:  '#4CAF82',
  statusWarning: '#F5A623',  // Amber Soft
  statusCritical:'#E05252',  // Coral Red

  border:        '#D4E8DF',  // Pale Sage
} as const;
```

### 3.2 Typography

```typescript
export const FONTS = {
  display: 'Exo 2',          // Hero, H1–H3 (weight: 600, 700, 800)
  body:    'DM Sans',        // Body, label, UI (weight: 400, 500)
  mono:    'JetBrains Mono', // Angka sensor, timestamp (weight: 400)
} as const;
```

### 3.3 Tailwind Config (tailwind.config.ts)

```typescript
import type { Config } from 'tailwindcss'

const config: Config = {
  content: ['./src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        healy: {
          'bg':        '#F8FAFB',
          'bg-alt':    '#EFF4F2',
          'surface':   '#FFFFFF',
          'sage':      '#4CAF82',
          'sage-dark': '#2E8B62',
          'mint':      '#A8DFCC',
          'graphite':  '#1A2633',
          'slate':     '#5A7080',
          'border':    '#D4E8DF',
          'warning':   '#F5A623',
          'critical':  '#E05252',
        },
      },
      fontFamily: {
        display: ['Exo 2', 'sans-serif'],
        body:    ['DM Sans', 'sans-serif'],
        mono:    ['JetBrains Mono', 'monospace'],
      },
      borderRadius: { 'card': '16px' },
      boxShadow: {
        'card':      '0 4px 24px rgba(46, 139, 98, 0.08)',
        'card-hover':'0 8px 32px rgba(46, 139, 98, 0.14)',
        'glow':      '0 0 0 3px rgba(76, 175, 130, 0.20)',
      },
      animation: {
        'pulse-critical': 'pulse 1.2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'fade-up':        'fadeUp 0.3s ease-out forwards',
      },
      keyframes: {
        fadeUp: {
          '0%':   { opacity: '0', transform: 'translateY(12px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}

export default config
```

### 3.4 CSS Variables (globals.css)

```css
@import url('https://fonts.googleapis.com/css2?family=Exo+2:wght@400;600;700;800&family=DM+Sans:wght@400;500&family=JetBrains+Mono:wght@400&display=swap');

:root {
  --color-bg:        #F8FAFB;
  --color-bg-alt:    #EFF4F2;
  --color-surface:   #FFFFFF;
  --color-sage:      #4CAF82;
  --color-sage-dark: #2E8B62;
  --color-mint:      #A8DFCC;
  --color-graphite:  #1A2633;
  --color-slate:     #5A7080;
  --color-border:    #D4E8DF;
  --color-warning:   #F5A623;
  --color-critical:  #E05252;

  --font-display: 'Exo 2', sans-serif;
  --font-body:    'DM Sans', sans-serif;
  --font-mono:    'JetBrains Mono', monospace;

  --radius-card: 16px;
  --shadow-card: 0 4px 24px rgba(46, 139, 98, 0.08);
  --shadow-glow: 0 0 0 3px rgba(76, 175, 130, 0.20);
}

body {
  background-color: var(--color-bg);
  color: var(--color-graphite);
  font-family: var(--font-body);
}
```

### 3.5 Animasi & Motion Principles

| Animasi | Trigger | Durasi | Easing |
|---|---|---|---|
| Page load sections | Mount | 300ms, stagger 80ms | ease-out |
| Angka sensor update | Data baru masuk | 600ms | ease-out counter |
| Alert badge CRITICAL | Status = CRITICAL | Infinite 1.2s | cubic-bezier pulse |
| Card hover elevation | Mouse enter | 200ms | ease |
| Sparkline draw | First render | 800ms | ease-in-out |
| Toast notification | Alert triggered | Slide-in 250ms | ease-out |
| Scroll reveal | Intersection 20% | 400ms, stagger | ease-out |

Library animasi: `framer-motion` untuk React components. CSS keyframes untuk background effects.

### 3.6 Background Texture (Hexagonal Grid)

```css
.hero-bg {
  background-color: var(--color-bg);
  background-image: url("data:image/svg+xml,%3Csvg width='60' height='52' viewBox='0 0 60 52' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M30 0 L60 17.3 L60 34.6 L30 52 L0 34.6 L0 17.3 Z' fill='none' stroke='%234CAF82' stroke-width='0.5' opacity='0.07'/%3E%3C/svg%3E");
  background-size: 60px 52px;
}
```

---

## 4. Backend: Spesifikasi Teknis

### 4.1 Domain Structs (internal/domain/telemetry.go) — STATUS: DONE ✅

Struct ini sudah disync 100% dengan JSON payload ESP32. Jangan ubah field names atau json tags tanpa koordinasi dengan IoT layer.

```go
package domain

import "time"

type SensorStatus string

const (
  StatusNormal   SensorStatus = "NORMAL"
  StatusWarning  SensorStatus = "WARNING"
  StatusCritical SensorStatus = "CRITICAL"
)

type TelemetryPayload struct {
  DeviceID  string     `json:"device_id"`
  Timestamp time.Time  `json:"timestamp"`
  Sensor    SensorData `json:"sensor"`
}

type SensorData struct {
  Temperature float64 `json:"temperature"`
  BPM         int     `json:"bpm"`
  SpO2        int     `json:"spo2"`
}

type TelemetryRecord struct {
  TelemetryPayload
  Status    EvaluatedStatus `json:"status"`
  CreatedAt time.Time       `json:"created_at"`
}

type EvaluatedStatus struct {
  Temperature SensorStatus `json:"temperature"`
  SpO2        SensorStatus `json:"spo2"`
  Overall     SensorStatus `json:"overall"`
}
```

### 4.2 Threshold Engine (internal/usecase/alert_usecase.go) — STATUS: IN PROGRESS ⬜

Ini adalah **tugas aktif saat ini** (Phase 2). Fungsi-fungsi ini adalah pure functions — tidak ada side effect, tidak ada database call. Mudah di-unit test.

```go
package usecase

import "github.com/rafif/healy-backend/internal/domain"

// Threshold defaults — overrideable via device_settings table
const (
  TempNormalMin float64 = 36.5
  TempNormalMax float64 = 37.5
  TempWarnMax   float64 = 38.5  // Above this = CRITICAL

  SpO2NormalMin int = 95
  SpO2WarnMin   int = 91        // Below this = CRITICAL
)

func EvaluateTemperature(temp float64) domain.SensorStatus {
  switch {
  case temp >= TempNormalMin && temp <= TempNormalMax:
    return domain.StatusNormal
  case temp > TempNormalMax && temp <= TempWarnMax:
    return domain.StatusWarning
  default:
    return domain.StatusCritical
  }
}

func EvaluateSpO2(spo2 int) domain.SensorStatus {
  switch {
  case spo2 >= SpO2NormalMin:
    return domain.StatusNormal
  case spo2 >= SpO2WarnMin:
    return domain.StatusWarning
  default:
    return domain.StatusCritical
  }
}

func EvaluateOverall(tempStatus, spo2Status domain.SensorStatus) domain.SensorStatus {
  if tempStatus == domain.StatusCritical || spo2Status == domain.StatusCritical {
    return domain.StatusCritical
  }
  if tempStatus == domain.StatusWarning || spo2Status == domain.StatusWarning {
    return domain.StatusWarning
  }
  return domain.StatusNormal
}

// EvaluatePayload adalah entry point utama — menerima raw payload dari ESP32
// dan mengembalikan TelemetryRecord yang siap disimpan ke DB
func EvaluatePayload(payload domain.TelemetryPayload) domain.TelemetryRecord {
  tempStatus := EvaluateTemperature(payload.Sensor.Temperature)
  spo2Status := EvaluateSpO2(payload.Sensor.SpO2)
  return domain.TelemetryRecord{
    TelemetryPayload: payload,
    Status: domain.EvaluatedStatus{
      Temperature: tempStatus,
      SpO2:        spo2Status,
      Overall:     EvaluateOverall(tempStatus, spo2Status),
    },
  }
}
```

### 4.3 Database Connection (pkg/database/postgres.go) — STATUS: DONE ✅

PERHATIAN: Gunakan `pgx/v5` dan port **6543** (Transaction Pooler Supabase). Bukan `lib/pq` dan bukan port 5432.

```go
package database

import (
  "context"
  "fmt"
  "os"

  "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
  dsn := os.Getenv("DATABASE_URL") // Supabase URI dengan port 6543

  config, err := pgxpool.ParseConfig(dsn)
  if err != nil {
    return nil, fmt.Errorf("failed to parse db config: %w", err)
  }

  config.MaxConns = 10  // Supabase Free Tier: max 15 connections total
  config.MinConns = 2

  pool, err := pgxpool.NewWithConfig(ctx, config)
  if err != nil {
    return nil, fmt.Errorf("failed to create db pool: %w", err)
  }

  return pool, nil
}
```

### 4.4 Repository — Query dengan Schema-Qualified Table Name

WAJIB: Semua query ke tabel telemetry harus menggunakan nama lengkap `telemetry.telemetry_records`. Query tanpa schema prefix akan error karena tabel tidak ada di schema `public`.

```go
// internal/repository/postgres/telemetry_postgres.go

// BENAR:
const insertTelemetry = `
  INSERT INTO telemetry.telemetry_records 
  (device_id, recorded_at, temperature, bpm, spo2, temp_status, spo2_status, overall_status)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

// SALAH — akan error "relation does not exist":
// INSERT INTO telemetry_records ...

const queryHistory = `
  SELECT device_id, recorded_at, temperature, bpm, spo2, 
         temp_status, spo2_status, overall_status
  FROM telemetry.telemetry_records
  WHERE device_id = $1 
    AND recorded_at >= NOW() - $2::interval
  ORDER BY recorded_at DESC
  LIMIT 500
`

const queryLatest = `
  SELECT DISTINCT ON (device_id) 
    device_id, recorded_at, temperature, bpm, spo2,
    temp_status, spo2_status, overall_status
  FROM telemetry.telemetry_records
  WHERE device_id = $1
  ORDER BY device_id, recorded_at DESC
`
```

### 4.5 WebSocket Hub Pattern

```go
// internal/delivery/websocket/hub.go

type Hub struct {
  clients    map[*Client]bool
  broadcast  chan []byte
  register   chan *Client
  unregister chan *Client
}

func (h *Hub) Run() {
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
    case client := <-h.unregister:
      delete(h.clients, client)
      close(client.send)
    case message := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(h.clients, client)
        }
      }
    }
  }
}
```

### 4.6 REST API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/auth/login` | No | Login, return JWT |
| GET | `/api/telemetry/history` | JWT | Riwayat data, param: `range` (1h/6h/24h/7d) |
| GET | `/api/telemetry/latest` | JWT | Data sensor terbaru |
| PUT | `/api/settings/threshold` | JWT | Update threshold |
| GET | `/api/settings/threshold` | JWT | Ambil threshold |
| GET | `/api/device/status` | JWT | Status koneksi ESP32 |
| GET | `/ws/telemetry` | JWT (query param) | WebSocket — Frontend |
| GET | `/ws/device` | Header device_id | WebSocket — ESP32 |

---

## 5. Frontend: Spesifikasi Teknis

### 5.1 TypeScript Interfaces (src/types/telemetry.ts)

Interfaces ini harus identik dengan Go structs di Section 4.1.

```typescript
export type SensorStatus = 'NORMAL' | 'WARNING' | 'CRITICAL'

export interface SensorData {
  temperature: number
  bpm: number
  spo2: number
}

export interface EvaluatedStatus {
  temperature: SensorStatus
  spo2: SensorStatus
  overall: SensorStatus
}

export interface TelemetryPayload {
  device_id: string
  timestamp: string   // ISO 8601
  sensor: SensorData
  status: EvaluatedStatus
}

export interface ConnectionState {
  status: 'CONNECTED' | 'DISCONNECTED' | 'RECONNECTING'
  lastUpdate: Date | null
  retryCount: number
}
```

### 5.2 WebSocket Hook (src/hooks/useWebSocket.ts)

```typescript
import { useEffect, useRef, useState, useCallback } from 'react'
import { TelemetryPayload, ConnectionState } from '@/types/telemetry'

const RECONNECT_DELAYS = [1000, 2000, 4000, 8000, 16000, 30000]

export function useWebSocket(url: string) {
  const [data, setData] = useState<TelemetryPayload | null>(null)
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
      const payload: TelemetryPayload = JSON.parse(event.data)
      setData(payload)
      setConn(prev => ({ ...prev, lastUpdate: new Date() }))
    }

    wsRef.current.onclose = () => {
      const delay = RECONNECT_DELAYS[Math.min(retryRef.current, RECONNECT_DELAYS.length - 1)]
      retryRef.current++
      setConn(prev => ({ ...prev, status: 'RECONNECTING', retryCount: retryRef.current }))
      setTimeout(connect, delay)
    }
  }, [url])

  useEffect(() => {
    connect()
    return () => wsRef.current?.close()
  }, [connect])

  return { data, conn }
}
```

### 5.3 Mock Telemetry Generator (src/lib/mock-telemetry.ts)

```typescript
import { TelemetryPayload } from '@/types/telemetry'

export function generateMockPayload(): TelemetryPayload {
  const rand = Math.random()
  // Distribusi: 80% normal, 15% warning, 5% critical
  const temp = rand < 0.80 ? 36.5 + Math.random() * 1.0
             : rand < 0.95 ? 37.6 + Math.random() * 0.9
             : 38.6 + Math.random() * 1.0

  const spo2 = rand < 0.80 ? Math.floor(95 + Math.random() * 4)
             : rand < 0.95 ? Math.floor(91 + Math.random() * 3)
             : Math.floor(85 + Math.random() * 5)

  const bpm = Math.floor(65 + Math.random() * 30)

  const tempStatus = temp >= 36.5 && temp <= 37.5 ? 'NORMAL'
                   : temp <= 38.5 ? 'WARNING' : 'CRITICAL'
  const spo2Status = spo2 >= 95 ? 'NORMAL' : spo2 >= 91 ? 'WARNING' : 'CRITICAL'
  const statuses = [tempStatus, spo2Status]
  const overall = statuses.includes('CRITICAL') ? 'CRITICAL'
                : statuses.includes('WARNING')  ? 'WARNING' : 'NORMAL'

  return {
    device_id: 'healy-001',
    timestamp: new Date().toISOString(),
    sensor: { temperature: parseFloat(temp.toFixed(1)), bpm, spo2 },
    status: { temperature: tempStatus as any, spo2: spo2Status as any, overall: overall as any },
  }
}
```

---

## 6. Database Schema

Database menggunakan **Supabase PostgreSQL** dengan ekstensi `pg_partman` untuk partisi tabel time-series dan `pg_cron` untuk maintenance otomatis.

PERHATIAN: Supabase Free Tier memberikan 500MB. Dengan frekuensi sensor 1–2 detik (~50.000–100.000 baris/hari), database akan penuh dalam sekitar 30 hari tanpa retention policy.

### 6.1 Setup Ekstensi

```sql
-- Aktifkan di Supabase Dashboard: Extensions > pg_partman, pg_cron
-- Atau via SQL:
CREATE EXTENSION IF NOT EXISTS pg_partman SCHEMA partman;
CREATE EXTENSION IF NOT EXISTS pg_cron;
```

### 6.2 DDL: Schema dan Tabel Utama

```sql
-- Buat schema terpisah untuk data telemetry
CREATE SCHEMA IF NOT EXISTS telemetry;

-- Tabel users (schema public)
CREATE TABLE public.users (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username   VARCHAR(50) UNIQUE NOT NULL,
  password   VARCHAR(255) NOT NULL,   -- bcrypt, cost 12
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tabel parent telemetry (schema telemetry) — dipartisi by day
CREATE TABLE telemetry.telemetry_records (
  id             UUID NOT NULL DEFAULT gen_random_uuid(),
  device_id      VARCHAR(50) NOT NULL,
  recorded_at    TIMESTAMPTZ NOT NULL,
  temperature    DECIMAL(4,1) NOT NULL,
  bpm            SMALLINT NOT NULL,
  spo2           SMALLINT NOT NULL,
  temp_status    VARCHAR(10) NOT NULL,
  spo2_status    VARCHAR(10) NOT NULL,
  overall_status VARCHAR(10) NOT NULL,
  PRIMARY KEY (id, recorded_at)
) PARTITION BY RANGE (recorded_at);

-- Setup pg_partman: buat partisi harian otomatis
SELECT partman.create_parent(
  p_parent_table  => 'telemetry.telemetry_records',
  p_control       => 'recorded_at',
  p_type          => 'range',
  p_interval      => '1 day',
  p_premake       => 7  -- Buat 7 hari ke depan
);

-- Tabel alert logs (schema public)
CREATE TABLE public.alert_logs (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  device_id    VARCHAR(50) NOT NULL,
  alert_type   VARCHAR(30) NOT NULL,
  value        DECIMAL(5,2) NOT NULL,
  status       VARCHAR(10) NOT NULL,
  triggered_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tabel settings per device (schema public)
CREATE TABLE public.device_settings (
  device_id      VARCHAR(50) PRIMARY KEY,
  temp_warn_max  DECIMAL(4,1) DEFAULT 37.5,
  temp_crit_max  DECIMAL(4,1) DEFAULT 38.5,
  spo2_warn_min  SMALLINT DEFAULT 94,
  spo2_crit_min  SMALLINT DEFAULT 90,
  updated_at     TIMESTAMPTZ DEFAULT NOW()
);
```

### 6.3 Index untuk Query Performa

```sql
-- Index untuk query history (paling sering diakses)
CREATE INDEX ON telemetry.telemetry_records (device_id, recorded_at DESC);

-- Index untuk query latest per device
CREATE INDEX ON telemetry.telemetry_records (device_id, recorded_at DESC)
  WHERE overall_status IS NOT NULL;
```

### 6.4 Retention Policy (WAJIB — Free Tier 500MB)

```sql
-- Konfigurasi pg_partman: hapus partisi lebih dari 30 hari
UPDATE partman.part_config
SET retention          = '30 days',
    retention_keep_table = false,
    retention_keep_index = false
WHERE parent_table = 'telemetry.telemetry_records';

-- Setup pg_cron: jalankan maintenance setiap hari pukul 02:00 UTC
SELECT cron.schedule(
  'healy-partman-maintenance',
  '0 2 * * *',
  $$CALL partman.run_maintenance_proc()$$
);
```

### 6.5 Query Referensi Kritis

```sql
-- Optimized query untuk history (gunakan time_bucket jika data banyak)
SELECT 
  DATE_TRUNC('minute', recorded_at) AS bucket,
  AVG(temperature)::DECIMAL(4,1) AS avg_temp,
  AVG(bpm)::INT AS avg_bpm,
  AVG(spo2)::INT AS avg_spo2,
  MAX(overall_status) AS status
FROM telemetry.telemetry_records  -- ← Schema-qualified wajib
WHERE device_id = $1
  AND recorded_at >= NOW() - INTERVAL '24 hours'
GROUP BY bucket
ORDER BY bucket DESC;

-- Latest per device (< 5ms dengan index)
SELECT DISTINCT ON (device_id)
  device_id, recorded_at, temperature, bpm, spo2,
  temp_status, spo2_status, overall_status
FROM telemetry.telemetry_records  -- ← Schema-qualified wajib
WHERE device_id = $1
ORDER BY device_id, recorded_at DESC;
```

---

## 7. Environment Variables

### 7.1 Backend (backend/.env.example)

```env
# Server
APP_PORT=8080
APP_ENV=development

# Supabase PostgreSQL — Transaction Pooler
# PENTING: Gunakan port 6543, bukan 5432
# Format: postgres://user:password@host:6543/dbname?sslmode=require
DATABASE_URL=postgres://postgres.[project-ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres?sslmode=require

# JWT
JWT_SECRET=your-secret-key-minimum-32-characters-long
JWT_EXPIRY_HOURS=24

# CORS — tambahkan URL Railway/Render di production
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://healy.vercel.app

# WebSocket paths
WS_DEVICE_PATH=/ws/device
WS_CLIENT_PATH=/ws/telemetry
```

### 7.2 Frontend (frontend/.env.example)

```env
# WebSocket — ganti dengan URL Railway/Render di production
# Development: ws://localhost:8080/ws/telemetry
# Production:  wss://healy-backend.railway.app/ws/telemetry
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws/telemetry

# API Base URL
NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Toggle mock data (set false untuk connect ke backend nyata)
NEXT_PUBLIC_USE_MOCK_DATA=true

# Demo credentials (hardcoded untuk development)
NEXT_PUBLIC_DEMO_USER=admin
NEXT_PUBLIC_DEMO_PASS=healy123
```

---

## 8. Deployment Strategy

### 8.1 Frontend — Vercel

Next.js 14 di-deploy ke Vercel. Ini tepat karena Next.js dioptimalkan untuk Vercel (SSR/SSG, CDN edge functions, zero-config deployment).

```
Platform : Vercel
Build cmd: pnpm build
Root dir : frontend/
Env vars : Set NEXT_PUBLIC_WS_URL dan NEXT_PUBLIC_API_URL
           ke URL production Railway/Render
```

### 8.2 Backend — Railway atau Render

BACKEND WAJIB DI-DEPLOY KE RAILWAY ATAU RENDER. Alasan teknis yang tidak bisa dikompromikan: Vercel bersifat serverless — fungsinya hidup hanya saat ada request HTTP masuk dan mati setelahnya. WebSocket membutuhkan koneksi TCP yang tetap terbuka (persistent). Dua hal ini tidak kompatibel secara fundamental.

```
Platform  : Railway (direkomendasikan) atau Render
Build cmd : go build -o main ./cmd/api/
Start cmd : ./main
Port      : Sesuai APP_PORT dari env var
Env vars  : Set DATABASE_URL (Supabase) dan semua var dari Section 7.1
```

Perbandingan Railway vs Render untuk kasus ini:

| Aspek | Railway | Render |
|---|---|---|
| Free Tier | $5 credit/bulan | 750 jam/bulan (sleep setelah 15 menit idle) |
| WebSocket | Supported penuh | Supported, tapi tier gratis bisa sleep |
| Cold start | Tidak ada | Ada pada free tier |
| Rekomendasi untuk HEALY | Lebih baik (tidak sleep) | Oke untuk demo |

### 8.3 Database — Supabase

```
Platform  : Supabase (cloud, Free Tier)
Connection: Transaction Pooler, port 6543
Driver    : pgx/v5 (bukan lib/pq)
Limit     : 500MB storage, 2 cores, 1GB RAM shared
```

---

## 9. Known Technical Gotchas

Ini adalah daftar masalah yang SUDAH DIIDENTIFIKASI. Baca sebelum memulai fase development apapun.

---

### GOTCHA-01: Vercel + WebSocket = Tidak Kompatibel

**Status:** RESOLVED — keputusan di Section 8.

**Masalah:** Vercel adalah serverless platform. Setiap function invocation bersifat stateless dan akan timeout. WebSocket membutuhkan koneksi yang hidup terus-menerus.

**Konsekuensi jika diabaikan:** Robot HEALY tidak bisa mengirim data real-time. Koneksi akan putus setiap 10–30 detik (Vercel function timeout).

**Solusi:** Backend Golang di Railway/Render. Frontend Next.js tetap di Vercel.

---

### GOTCHA-02: Schema-Qualified Table Name

**Status:** ACTIVE — wajib diperhatikan setiap menulis query.

**Masalah:** Karena `pg_partman` disetup dengan tabel di schema `telemetry` (bukan `public`), PostgreSQL tidak akan menemukan tabel jika kamu menulis query tanpa schema prefix.

**Kode yang SALAH:**
```go
// Ini akan error: "ERROR: relation 'telemetry_records' does not exist"
row := pool.QueryRow(ctx, "SELECT * FROM telemetry_records WHERE ...")
```

**Kode yang BENAR:**
```go
// Selalu gunakan schema-qualified name
row := pool.QueryRow(ctx, "SELECT * FROM telemetry.telemetry_records WHERE ...")
```

**Cara cek via psql/Supabase dashboard:**
```sql
-- Lihat semua tabel di schema telemetry
SELECT tablename FROM pg_tables WHERE schemaname = 'telemetry';
```

---

### GOTCHA-03: Supabase Free Tier Storage Limit (500MB)

**Status:** ACTIVE — retention policy sudah disetup di Section 6.4, tapi harus dimonitor.

**Masalah:** Supabase Free Tier membatasi storage di 500MB. Dengan data sensor setiap 1–2 detik dan ukuran baris sekitar 100–150 bytes, estimasi:
- 1 hari = ~50.000 baris × 150 bytes = ~7.5MB
- 30 hari = ~225MB
- 60 hari = ~450MB (mendekati batas)

**Solusi yang sudah diimplementasikan:** Retention policy 30 hari via `pg_partman` + `pg_cron` (Section 6.4). Partisi lebih dari 30 hari akan di-drop otomatis.

**Monitoring:** Cek storage usage di Supabase Dashboard > Settings > Usage. Set alert jika mendekati 400MB.

---

### GOTCHA-04: pgx/v5 vs lib/pq — Driver Tidak Sama

**Status:** RESOLVED — db.go sudah menggunakan pgx/v5.

**Masalah:** `lib/pq` tidak mendukung semua fitur Supabase Transaction Pooler dengan baik, terutama pada port 6543. Jika ada yang me-refactor kode database dan mengganti ke `lib/pq`, koneksi bisa gagal.

**Solusi:** Selalu gunakan `pgx/v5`. Pastikan go.mod tidak ada dependency `lib/pq`.

```bash
# Verifikasi tidak ada lib/pq
grep "lib/pq" go.mod  # Harus tidak ada output
```

---

### GOTCHA-05: ESP32 mDNS — Browser Tidak Selalu Resolve

**Status:** ACTIVE — berlaku di environment development.

**Masalah:** `healy-server.local` via mDNS bekerja di OS-level (terminal, Golang), tapi browser Chrome/Firefox kadang tidak me-resolve nama `.local` karena kebijakan DNS mereka.

**Solusi:** Gunakan `NEXT_PUBLIC_WS_URL` di `.env` — set ke IP langsung saat development lokal, dan ke URL Railway/Render di production. Jangan hardcode `healy-server.local` di kode frontend.

---

## 10. Rencana Kerja Step-by-Step

### FASE 0 — Project Init (DONE ✅)

- Monorepo structure dibuat.
- Go modules diinisialisasi.
- Supabase project dibuat, ekstensi pg_partman diaktifkan.
- SQL migration (DDL Section 6.2) sudah dijalankan.
- `.env` dikonfigurasi dengan Supabase URI.
- `code-review-graph` disetup dan dioptimasi.

### FASE 1 — Backend Domain & Repository (DONE ✅)

- `internal/domain/telemetry.go` — structs synced dengan ESP32 payload.
- `internal/repository/postgres/db.go` — connection pool pgx/v5, port 6543.
- `internal/repository/interfaces/` — interface definitions.

### FASE 2 — Backend Usecase Layer (AKTIF SEKARANG ⬜)

Ini adalah fase yang sedang dikerjakan. Target output:

**Step 2.1 — alert_usecase.go (pure functions)**
Implementasikan `EvaluateTemperature`, `EvaluateSpO2`, `EvaluateOverall`, dan `EvaluatePayload` sesuai Section 4.2. Tulis unit test untuk semua fungsi — tidak perlu database.

**Step 2.2 — telemetry_usecase.go (orchestrator)**
Implementasikan flow: terima `TelemetryPayload` → panggil `EvaluatePayload` → save ke `telemetry.telemetry_records` → kirim ke broadcast channel hub.

**Step 2.3 — Integrasi dengan Repository**
Hubungkan usecase dengan repository interface. Pastikan usecase tidak langsung import implementasi postgres — hanya interface.

### FASE 3 — Backend Delivery Layer ⬜

- WebSocket hub + client pump.
- HTTP handlers (auth, telemetry, settings).
- JWT middleware.
- Router setup.

### FASE 4 — Frontend Foundation + Landing Page (DONE ✅)

- Setup design tokens, tailwind config, globals.css.
- Shared components: SensorCard, StatusChip, AlertToast, NavSidebar.
- Landing Page dengan animasi.
- Login Page.

### FASE 5 — Frontend Dashboard Real-Time (DONE ✅)

- `useWebSocket` hook + mock generator.
- Dashboard page: 3 sensor cards, alert feed, connection status.

### FASE 6 — History, Settings, Integrasi E2E ⬜

- History page dengan Recharts.
- Settings page.
- End-to-end test: ESP32 → Backend → Frontend.

### FASE 7 — Polish, QA, Deployment ⬜

- Lighthouse audit.
- Deploy Backend ke Railway.
- Deploy Frontend ke Vercel.
- End-to-end test di environment production.

---

## 11. Codebase Progress Tracker

| Layer | File | Status | Catatan |
|---|---|---|---|
| Domain | `internal/domain/telemetry.go` | ✅ DONE | Synced dengan ESP32 payload |
| Domain | `internal/domain/alert.go` | ✅ DONE | |
| Domain | `internal/domain/user.go` | ✅ DONE | |
| Domain | `internal/domain/settings.go` | ✅ DONE | |
| Repository | `internal/repository/postgres/db.go` | ✅ DONE | pgx/v5, port 6543 |
| Repository | `internal/repository/interfaces/telemetry_repo.go` | ✅ DONE | |
| Repository | `internal/repository/postgres/telemetry_postgres.go` | ✅ DONE | |
| Usecase | `internal/usecase/alert_usecase.go` | ✅ DONE | Threshold engine selesai |
| Usecase | `internal/usecase/telemetry_usecase.go` | ✅ DONE | Orchestrator selesai |
| Usecase | `internal/usecase/auth_usecase.go` | ✅ DONE | |
| Delivery | `internal/delivery/websocket/hub.go` | ✅ DONE | Hub, client, dan handler selesai |
| Delivery | `internal/delivery/http/router.go` | ✅ DONE | Routing dan server selesai |
| Config | `pkg/config/config.go` | ✅ DONE | |
| Frontend | Design tokens + Tailwind | ✅ DONE | Tailwind v4 @theme, globals.css, design-tokens.ts |
| Frontend | SensorCard, StatusChip | ✅ DONE | + NavSidebar, glass-card, framer-motion |
| Frontend | Dashboard page | ✅ DONE | Real-time integration via useTelemetry hook |
| Frontend | Landing page | ✅ DONE | Hero section, features, login page |
| Frontend | `hooks/useWebSocket.ts` | ✅ DONE | Blueprint §5.2, exponential backoff reconnect |
| Frontend | `lib/mock-telemetry.ts` | ✅ DONE | Blueprint §5.3, 80/15/5 distribution |
| Frontend | `hooks/useTelemetry.ts` | ✅ DONE | Bridge hook: mock ↔ WebSocket |
| Frontend | `ConnectionStatus.tsx` | ✅ DONE | Live connection state indicator |

---

## 12. Konvensi Kode & Git

### 12.1 Branch Strategy

```
main              ← Production-ready only
develop           ← Integration branch
feat/usecase-alert
feat/ws-hub
feat/dashboard-realtime
fix/schema-qualified-query
chore/deploy-railway
```

### 12.2 Commit Message Format

```
feat(usecase): implement threshold evaluation engine
feat(repository): add telemetry_postgres with schema-qualified queries
fix(db): switch from lib/pq to pgx/v5 for Supabase compatibility
chore(deploy): configure Railway backend deployment
docs(blueprint): update v2.1.0 deployment strategy and gotchas
```

### 12.3 Naming Conventions

| Layer | Convention | Contoh |
|---|---|---|
| Go files | snake_case | `telemetry_usecase.go` |
| Go types/funcs | PascalCase | `TelemetryPayload`, `EvaluatePayload` |
| Go SQL queries | Konstanta SCREAMING_SNAKE | `insertTelemetry`, `queryHistory` |
| React components | PascalCase | `SensorCard.tsx` |
| React hooks | camelCase + `use` prefix | `useWebSocket.ts` |
| TypeScript types | PascalCase | `SensorStatus`, `TelemetryPayload` |
| Constants | SCREAMING_SNAKE | `TEMP_WARN_MAX` |

---

## 13. Prompt Inisialisasi AI (Master Prompt)

### Prompt 0 — Grounding (Wajib di awal setiap session)

```
You are a Senior Full-Stack Engineer on Project HEALY — Health Observer Robot.

I am attaching:
1. PRD_HEALY_Website.md
2. HEALY_Master_Blueprint.md (v2.1.0)

Read both documents completely. Confirm understanding by stating:
1. Current tech stack: IoT / Backend / Frontend / Database / Deployment
2. Database: What schema holds telemetry data? What is the exact table name?
3. What driver is used for DB connection, and why port 6543?
4. What is the current development phase and the immediate next task?
5. List the 5 Known Gotchas from Blueprint Section 9

Do NOT generate any code yet. Confirmation only.

[ATTACH: PRD_HEALY_Website.md]
[ATTACH: HEALY_Master_Blueprint.md]
```

---

### Prompt Fase 2 — Usecase Layer (AKTIF SEKARANG)

```
Based on HEALY Master Blueprint Section 4.2, implement the Usecase Layer.

Context:
- Domain structs are in internal/domain/telemetry.go (already done, do not modify)
- Repository interface is in internal/repository/interfaces/telemetry_repo.go

Step 1: internal/usecase/alert_usecase.go
Implement ALL functions from Blueprint 4.2:
- EvaluateTemperature(temp float64) domain.SensorStatus
- EvaluateSpO2(spo2 int) domain.SensorStatus  
- EvaluateOverall(tempStatus, spo2Status domain.SensorStatus) domain.SensorStatus
- EvaluatePayload(payload domain.TelemetryPayload) domain.TelemetryRecord
Threshold constants: TempNormalMin=36.5, TempNormalMax=37.5, TempWarnMax=38.5,
SpO2NormalMin=95, SpO2WarnMin=91

Step 2: internal/usecase/alert_usecase_test.go
Write table-driven unit tests for all 4 functions.
Test cases must cover: boundary values, each status category (Normal/Warning/Critical).
No database or external dependencies — pure logic only.

Step 3: internal/usecase/telemetry_usecase.go
Interface: TelemetryUsecase with method ProcessIncoming(ctx, payload TelemetryPayload) error
Implementation:
1. Call EvaluatePayload to get TelemetryRecord
2. Call repository.Save(ctx, record) to persist
3. Marshal record to JSON
4. Send to broadcastChan channel (passed via constructor)
5. If overall status == CRITICAL, call repository.SaveAlert(ctx, alert)

Inject dependencies via constructor — no global variables.

Output: Three files with full implementation. No placeholder comments.
```

---

### Prompt Fase 3 — Backend Delivery Layer

```
Based on HEALY Master Blueprint Section 4.3–4.6, implement the Delivery Layer.

Step 1: internal/delivery/websocket/hub.go + client.go + handler.go
- Hub manages two sets of clients: device clients (ESP32) and viewer clients (Frontend)
- Message from any device client is broadcast to ALL viewer clients
- Client timeout: 60s read, 10s write, 54s ping
- Full implementation from Blueprint Section 4.3

Step 2: internal/delivery/http/router.go
Register all routes from Blueprint Section 4.6:
- POST /api/auth/login (no auth)
- GET /api/telemetry/history (JWT required)
- GET /api/telemetry/latest (JWT required)
- PUT /api/settings/threshold (JWT required)
- GET /api/device/status (JWT required)
- GET /ws/telemetry (JWT via query param)
- GET /ws/device (device_id header)

Step 3: cmd/api/main.go
Wire all dependencies: DB pool → Repository → Usecase → Hub → Handlers → Router.
Start server on APP_PORT from env.

CRITICAL: Do not import any concrete postgres implementation from usecase layer.
Only inject via interfaces.
```

---

*HEALY Master Blueprint v2.1.0 — Updated by Lead Developer.*
*Dokumen ini adalah satu-satunya master reference yang berlaku.*
*Versi sebelumnya (v2.0.0) dinyatakan deprecated pada poin-poin yang tercantum di Changelog.*
