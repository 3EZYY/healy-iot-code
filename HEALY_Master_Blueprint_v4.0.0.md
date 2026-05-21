# HEALY — Master Blueprint
## Health Observer Robot: Full-Stack Implementation Guide

---

| Field | Detail |
|---|---|
| **Document Version** | v4.0.0 |
| **Status** | Active — Master Reference |
| **Role Penyusun** | Senior Full-Stack Engineer & Technical Product Manager |
| **Stack** | ESP32 (IoT) + Golang 1.23+ (Backend) + Next.js 14 (Frontend) |
| **Design Philosophy** | Clinical Futurism — Sage Green, Glacial White, Futuristic Medical |
| **Last Updated** | 2025 |
| **Changelog** | v4.0.0: Light/Dark Mode, AI Accent Revision (Teal), Separated History Charts, Interactive AI Chatbot (Context-Aware) |

> **Catatan Arsitektural:** Semua spesifikasi Backend (Golang, Hub, WebSocket, REST API, Supabase, deployment Railway/Render) dari v3.0.0 tetap berlaku tanpa perubahan. v4.0.0 adalah **Frontend-only upgrade** yang dipicu oleh revisi Dosen Penguji. Dokumen ini hanya memuat delta perubahan terhadap v3.0.0. Baca v3.0.0 sebagai fondasi sebelum mengimplementasikan v4.0.0.

---

## Changelog v3.0.0 → v4.0.0

| # | ID | Area | File Terdampak | Jenis | Urgensi |
|---|---|---|---|---|---|
| C1 | F-01 | Design System | `design-tokens.ts`, `tailwind.config.ts`, `globals.css` | MODIFY | HIGH |
| C2 | F-01 | Theme Provider | `app/layout.tsx`, `app/providers.tsx` | MODIFY/NEW | HIGH |
| C3 | F-01 | Semua komponen UI | Seluruh `components/features/*.tsx` | MODIFY | HIGH |
| C4 | F-02 | AI Accent Token | `design-tokens.ts`, `tailwind.config.ts`, `globals.css` | MODIFY | HIGH |
| C5 | F-02 | Komponen AI | `AIInsightCard.tsx`, `NavSidebar.tsx`, `AIChatPanel.tsx` | MODIFY | HIGH |
| C6 | F-03 | History Page | `app/history/page.tsx` | MODIFY | HIGH |
| C7 | F-03 | Chart Components | `TemperatureChart.tsx`, `HeartRateChart.tsx`, `SpO2Chart.tsx` | NEW | HIGH |
| C8 | F-04 | AI Chatbot Panel | `AIChatPanel.tsx`, `useChatbot.ts` | NEW | HIGH |
| C9 | F-04 | Groq Client | `groq-client.ts` | MODIFY | HIGH |
| C10 | F-04 | NavSidebar | `NavSidebar.tsx` | MODIFY | HIGH |
| C11 | F-04 | Types | `src/types/chat.ts` | NEW | MEDIUM |
| C12 | META | Dependency | `package.json` | MODIFY | HIGH |
| C13 | META | Progress Tracker | Blueprint Section 13 | UPDATE | INFO |

---

## Table of Contents

1. [Ringkasan Fitur Baru v4.0.0](#1-ringkasan-fitur-baru-v400)
2. [Dependency Baru](#2-dependency-baru)
3. [Struktur Folder: Delta v3.0 → v4.0](#3-struktur-folder-delta-v30--v40)
4. [F-01: Light/Dark Mode Support](#4-f-01-lightdark-mode-support)
5. [F-02: Revisi AI Accent Color — Clinical Teal](#5-f-02-revisi-ai-accent-color--clinical-teal)
6. [F-03: Pemisahan Grafik Telemetry History](#6-f-03-pemisahan-grafik-telemetry-history)
7. [F-04: Interactive AI Chatbot — Context-Aware](#7-f-04-interactive-ai-chatbot--context-aware)
8. [Design Token Master — Versi Final v4.0.0](#8-design-token-master--versi-final-v400)
9. [groq-client.ts — Revisi Lengkap v4.0.0](#9-groq-clientts--revisi-lengkap-v400)
10. [Known Gotchas v4.0.0](#10-known-gotchas-v400)
11. [Rencana Kerja: Fase 12–15](#11-rencana-kerja-fase-1215)
12. [Codebase Progress Tracker — Updated](#12-codebase-progress-tracker--updated)
13. [Prompt Inisialisasi AI — v4.0.0](#13-prompt-inisialisasi-ai--v400)

---

## 1. Ringkasan Fitur Baru v4.0.0

| ID | Fitur | Scope | Status |
|---|---|---|---|
| F-01 | Light/Dark Mode — `next-themes` | Frontend-only | ✅ DONE |
| F-02 | AI Accent: Violet → Clinical Teal `#0D9488` | Design Token | ✅ DONE (via F-01 globals.css) |
| F-03 | Separated History Charts (3 line chart terpisah) | Frontend-only | TODO |
| F-04 | Interactive AI Chatbot Context-Aware (Groq) | Frontend-only | TODO |

---

## 2. Dependency Baru

Tambahkan ke `frontend/package.json` sebelum memulai implementasi:

```bash
# Install dari direktori /frontend
pnpm add next-themes
pnpm add @radix-ui/react-dialog
```

| Package | Versi Minimum | Kegunaan |
|---|---|---|
| `next-themes` | `^0.3.0` | Theme provider untuk Light/Dark mode di Next.js App Router |
| `@radix-ui/react-dialog` | `^1.1.0` | Modal/Dialog headless untuk AI Chatbot Panel (sudah include di shadcn/ui) |

---

## 3. Struktur Folder: Delta v3.0 → v4.0

Hanya file yang **baru dibuat** atau **dimodifikasi** yang tercantum di sini.

```
frontend/src/
│
├── app/
│   ├── layout.tsx                    ← MODIFY: tambahkan ThemeProvider wrapper
│   ├── providers.tsx                 ← NEW: centralized providers (Theme, Toast, dll)
│   └── history/
│       └── page.tsx                  ← MODIFY: ganti satu chart menjadi 3 chart terpisah
│
├── components/
│   ├── ui/
│   │   └── theme-toggle.tsx          ← NEW: tombol toggle Light/Dark
│   └── features/
│       ├── NavSidebar.tsx            ← MODIFY: tambahkan ChatbotTrigger button
│       ├── AIInsightCard.tsx         ← MODIFY: update warna ai-accent → teal
│       ├── AlertFeed.tsx             ← MODIFY: update warna ai-accent → teal
│       ├── AIChatPanel.tsx           ← NEW: sliding panel chatbot (Phase F-04)
│       ├── TemperatureChart.tsx      ← NEW: standalone chart suhu (Phase F-03)
│       ├── HeartRateChart.tsx        ← NEW: standalone chart BPM (Phase F-03)
│       └── SpO2Chart.tsx             ← NEW: standalone chart SpO2 (Phase F-03)
│
├── hooks/
│   └── useChatbot.ts                 ← NEW: state management chatbot
│
├── types/
│   └── chat.ts                       ← NEW: ChatMessage, ChatRole interfaces
│
└── constants/
    └── design-tokens.ts              ← MODIFY: update aiAccent + tambahkan dark mode tokens
```

---

## 4. F-01: Light/Dark Mode Support

### 4.1 Keputusan Arsitektur

`next-themes` dipilih karena merupakan solusi de-facto untuk tema di Next.js App Router. Library ini menangani: FOUC (Flash of Unstyled Content) prevention via `suppressHydrationWarning`, persistensi tema di `localStorage`, dan sinkronisasi dengan `prefers-color-scheme` sistem operasi.

Strategi implementasi: **CSS Variables sebagai satu-satunya mekanisme warna**. Semua komponen sudah menggunakan CSS variables (`--color-sage`, `--color-bg`, dst). Yang perlu dilakukan hanya mendefinisikan ulang nilai variabel tersebut di dalam selector `.dark`. Tidak ada perubahan di level JSX komponen individual — ini adalah keunggulan arsitektur token yang sudah dibangun di v2.1.0.

### 4.2 Setup: app/providers.tsx (NEW)

Pisahkan semua Provider ke file ini untuk menjaga `layout.tsx` tetap bersih.

```typescript
// src/app/providers.tsx
'use client'

import { ThemeProvider } from 'next-themes'
import { Toaster } from '@/components/ui/toaster'

interface ProvidersProps {
  children: React.ReactNode
}

export function Providers({ children }: ProvidersProps) {
  return (
    <ThemeProvider
      attribute="class"         // Tambahkan class "dark" ke <html> element
      defaultTheme="light"      // Default: light mode
      enableSystem={true}       // Hormati preferensi OS user
      disableTransitionOnChange={false}
    >
      {children}
      <Toaster />
    </ThemeProvider>
  )
}
```

### 4.3 Modifikasi: app/layout.tsx

```typescript
// src/app/layout.tsx — MODIFIKASI
import type { Metadata } from 'next'
import { Providers } from './providers'
import './globals.css'

export const metadata: Metadata = {
  title: 'HEALY — Health Observer Robot',
  description: 'Real-time health monitoring powered by IoT and AI',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="id" suppressHydrationWarning>
      {/*
        suppressHydrationWarning wajib ada.
        next-themes memodifikasi class di <html> setelah hydration,
        yang akan menyebabkan React hydration mismatch tanpa atribut ini.
      */}
      <body className="bg-healy-bg text-healy-graphite font-body antialiased">
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  )
}
```

### 4.4 Modifikasi: globals.css — Dark Mode Variables

Tambahkan blok `.dark` di bawah blok `:root` yang sudah ada. Jangan hapus `:root`.

```css
/* src/app/globals.css — TAMBAHKAN SETELAH BLOK :root */

/* ─── DARK MODE OVERRIDES ─── */
.dark {
  /* Background — dibalik: surface lebih gelap, bukan putih */
  --color-bg:        #0F1A14;   /* Deep Forest Dark — pengganti Glacial White */
  --color-bg-alt:    #162118;   /* Slightly lighter forest — pengganti Soft Mist */
  --color-surface:   #1C2B1F;   /* Card surface gelap */

  /* Brand — sedikit lebih terang agar kontras di dark bg */
  --color-sage:      #5DC491;   /* Sage lebih cerah di dark */
  --color-sage-dark: #4CAF82;   /* Hover — warna sage original */
  --color-mint:      #2D6B4F;   /* Mint lebih redup */

  /* Text — dibalik */
  --color-graphite:  #E8F5EE;   /* Light sage-white untuk body text */
  --color-slate:     #8BAF9A;   /* Muted sage untuk secondary text */

  /* Border — lebih visible di dark bg */
  --color-border:    #2A4035;

  /* Status — tidak berubah, tetap bold */
  --color-warning:   #F5A623;
  --color-critical:  '#E05252';

  /* AI (Clinical Teal) — v4.0 */
  --color-ai-accent: #14B8A6;   /* Teal lebih terang di dark mode */
  --color-ai-surface:#0D2B2A;   /* Dark teal surface */

  /* Device LED — tidak berubah */
  --color-device-on:  #22C55E;
  --color-device-off: #EF4444;
}

/* ─── TRANSITION GLOBAL ─── */
/* Smooth transition saat toggle tema */
*, *::before, *::after {
  transition-property: background-color, border-color, color;
  transition-timing-function: ease;
  transition-duration: 200ms;
}

/* Kecualikan elemen yang tidak butuh transisi */
button, a, input, [role="button"] {
  transition-duration: 150ms;
}
```

### 4.5 Modifikasi: tailwind.config.ts — Dark Mode Strategy

```typescript
// tailwind.config.ts — MODIFIKASI
import type { Config } from 'tailwindcss'

const config: Config = {
  content: ['./src/**/*.{ts,tsx}'],

  // WAJIB: gunakan 'class' strategy agar next-themes bisa toggle
  darkMode: 'class',

  theme: {
    extend: {
      colors: {
        healy: {
          // Semua warna mereferensi CSS variable
          // Saat dark mode aktif, CSS variable di-override oleh selector .dark
          'bg':         'rgb(var(--color-bg) / <alpha-value>)',
          'bg-alt':     'rgb(var(--color-bg-alt) / <alpha-value>)',
          'surface':    'rgb(var(--color-surface) / <alpha-value>)',
          'sage':       'rgb(var(--color-sage) / <alpha-value>)',
          'sage-dark':  'rgb(var(--color-sage-dark) / <alpha-value>)',
          'mint':       'rgb(var(--color-mint) / <alpha-value>)',
          'graphite':   'rgb(var(--color-graphite) / <alpha-value>)',
          'slate':      'rgb(var(--color-slate) / <alpha-value>)',
          'border':     'rgb(var(--color-border) / <alpha-value>)',
          'warning':    '#F5A623',
          'critical':   '#E05252',
          'device-on':  '#22C55E',
          'device-off': '#EF4444',
          // v4.0 — AI Teal (menggantikan Violet)
          'ai-accent':  'rgb(var(--color-ai-accent) / <alpha-value>)',
          'ai-surface': 'rgb(var(--color-ai-surface) / <alpha-value>)',
        },
      },
      // ... fontFamily, borderRadius, boxShadow, animation dari v3.0 tetap sama
    },
  },
  plugins: [],
}

export default config
```

> **Perhatian Implementasi:** Agar sintaks `rgb(var(...) / <alpha-value>)` berfungsi, nilai CSS variable di `:root` dan `.dark` harus dalam format **channel terpisah** (R G B), bukan hex. Contoh: `--color-sage: 76 175 130;` bukan `--color-sage: #4CAF82;`. Konversi semua color variable di `globals.css` ke format ini.

### 4.6 New Component: ThemeToggle (src/components/ui/theme-toggle.tsx)

```typescript
// src/components/ui/theme-toggle.tsx
'use client'

import { useTheme } from 'next-themes'
import { Moon, Sun } from '@phosphor-icons/react'
import { useEffect, useState } from 'react'

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  // Wajib: hindari hydration mismatch
  // Komponen ini hanya boleh render setelah mount di client
  useEffect(() => { setMounted(true) }, [])
  if (!mounted) return <div className="w-8 h-8" />

  const isDark = theme === 'dark'

  return (
    <button
      onClick={() => setTheme(isDark ? 'light' : 'dark')}
      aria-label={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
      className="flex items-center justify-center w-8 h-8 rounded-lg
                 text-healy-slate hover:text-healy-sage hover:bg-healy-bg-alt
                 transition-all duration-150"
    >
      {isDark
        ? <Sun size={18} weight="duotone" />
        : <Moon size={18} weight="duotone" />
      }
    </button>
  )
}
```

**Pasang di:** Dashboard Header (di samping `DeviceLedIndicator`) dan NavSidebar footer.

---

## 5. F-02: Revisi AI Accent Color — Clinical Teal

### 5.1 Keputusan Arsitektur

Warna `aiAccent` diubah dari **Violet `#7C3AED`** menjadi **Clinical Teal `#0D9488`**.

Justifikasi klinis: Teal berada dalam spektrum warna yang diasosiasikan dengan teknologi medis modern, kebersihan, dan ketepatan — konsisten dengan filosofi "Clinical Futurism" HEALY. Violet bersifat terlalu "consumer tech" dan menciptakan disonansi visual terhadap palette Sage Green utama. Teal dan Sage Green berada dalam kelompok warna analog (berdekatan di color wheel), menghasilkan harmoni visual yang lebih kohesif.

### 5.2 Token Lama vs Token Baru

| Token | v3.0.0 (Deprecated) | v4.0.0 (Berlaku) |
|---|---|---|
| `aiAccent` | `#7C3AED` Violet-700 | `#0D9488` Clinical Teal |
| `aiSurface` (Light) | `#F5F3FF` Violet-50 | `#F0FDFC` Teal-50 |
| `aiSurface` (Dark) | `#1E1B4B` | `#0D2B2A` Dark Teal |
| `aiAccent` Dark mode | `#8B5CF6` | `#14B8A6` Teal-400 |

### 5.3 Ganti di: design-tokens.ts

```typescript
// src/constants/design-tokens.ts — GANTI bagian AI token

export const COLORS = {
  // ... semua warna lama tetap sama ...

  // ─── AI TIER — v4.0.0 (menggantikan Violet dari v3.0.0) ───
  aiAccent:  '#0D9488',   // Clinical Teal — aksen utama AI elements
  aiSurface: '#F0FDFC',   // Teal-50 — background card AI (light mode)
  // Dark mode: dihandle oleh CSS variable .dark override

  // Device LED — tidak berubah
  deviceOnline:  '#22C55E',
  deviceOffline: '#EF4444',
} as const
```

### 5.4 File yang Harus Diupdate (Cari & Ganti)

Jalankan find-and-replace berikut di seluruh direktori `src/`:

| Cari | Ganti Dengan |
|---|---|
| `healy-ai-accent` (Tailwind class) | Tidak perlu ganti — CSS variable sudah diupdate |
| `hover:bg-violet-600` | `hover:bg-teal-700` |
| `focus:ring-healy-ai-accent` | Tidak perlu ganti |
| `text-healy-ai-accent` | Tidak perlu ganti |
| Komentar `Violet-700` | Update ke `Clinical Teal` |

File yang perlu diverifikasi manual:
- `src/components/features/AIInsightCard.tsx` — tombol Generate Insight
- `src/components/features/NavSidebar.tsx` — GroqKeySection
- `src/components/features/AlertFeed.tsx` — spinner dan narrative border

---

## 6. F-03: Pemisahan Grafik Telemetry History

### 6.1 Keputusan Arsitektur

History chart sebelumnya menampilkan tiga series (Suhu, BPM, SpO2) dalam satu komponen Recharts. Keputusan pemisahan ini diambil berdasarkan dua alasan teknis:

**Skala Y yang tidak kompatibel:** Suhu (35–40°C), BPM (40–200), dan SpO2 (80–100%) memiliki rentang nilai yang secara fundamental berbeda. Menggabungkannya dalam satu chart — bahkan dengan dual Y-axis — menghasilkan visualisasi yang misleading secara klinis dan sulit dibaca.

**Independensi analisis:** Dokter atau operator perlu menganalisis setiap parameter secara mandiri. Grafik terpisah memungkinkan fokus yang lebih tajam dan perbandingan yang lebih presisi.

### 6.2 Struktur Data: Props Standar Tiap Chart

Ketiga chart menerima props yang identik untuk konsistensi:

```typescript
// src/types/telemetry.ts — TAMBAHKAN

export interface TelemetryChartPoint {
  timestamp: string        // Format: HH:mm atau DD/MM tergantung range
  value: number
  status: SensorStatus     // Untuk color coding titik data
}

export interface ChartConfig {
  range: '1h' | '6h' | '24h' | '7d'
  parameter: 'temperature' | 'bpm' | 'spo2'
}
```

### 6.3 New Component: TemperatureChart.tsx

```typescript
// src/components/features/TemperatureChart.tsx
'use client'

import {
  LineChart, Line, XAxis, YAxis, CartesianGrid,
  Tooltip, ReferenceLine, ResponsiveContainer, Dot
} from 'recharts'
import { TelemetryChartPoint } from '@/types/telemetry'

interface TemperatureChartProps {
  data: TelemetryChartPoint[]
  loading?: boolean
}

// Custom dot: warna berbeda berdasarkan status
const StatusDot = (props: any) => {
  const { cx, cy, payload } = props
  const color = payload.status === 'CRITICAL' ? '#E05252'
              : payload.status === 'WARNING'  ? '#F5A623'
              : '#4CAF82'
  return <Dot cx={cx} cy={cy} r={3} fill={color} stroke="none" />
}

// Custom Tooltip dengan HEALY styling
const CustomTooltip = ({ active, payload, label }: any) => {
  if (!active || !payload?.length) return null
  return (
    <div className="bg-healy-surface border border-healy-border rounded-lg px-3 py-2 shadow-card text-xs">
      <p className="font-mono text-healy-slate mb-1">{label}</p>
      <p className="font-mono font-semibold text-healy-graphite">
        {payload[0].value.toFixed(1)} °C
      </p>
    </div>
  )
}

export function TemperatureChart({ data, loading }: TemperatureChartProps) {
  if (loading) {
    return (
      <div className="h-48 bg-healy-bg-alt rounded-card animate-pulse" />
    )
  }

  return (
    <div className="bg-healy-surface border border-healy-border rounded-card p-4 shadow-card">
      <div className="flex items-center justify-between mb-3">
        <h4 className="font-display font-semibold text-sm text-healy-graphite">
          Temperature Trend
        </h4>
        <span className="text-xs font-mono text-healy-slate">°C</span>
      </div>

      <ResponsiveContainer width="100%" height={160}>
        <LineChart data={data} margin={{ top: 4, right: 8, bottom: 0, left: -20 }}>
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="var(--color-border)"
            opacity={0.5}
          />
          <XAxis
            dataKey="timestamp"
            tick={{ fontSize: 10, fontFamily: 'JetBrains Mono', fill: 'var(--color-slate)' }}
            tickLine={false}
            axisLine={false}
          />
          <YAxis
            domain={[35, 40]}
            tick={{ fontSize: 10, fontFamily: 'JetBrains Mono', fill: 'var(--color-slate)' }}
            tickLine={false}
            axisLine={false}
          />
          <Tooltip content={<CustomTooltip />} />

          {/* Reference lines untuk threshold */}
          <ReferenceLine y={37.5} stroke="#F5A623" strokeDasharray="4 4" strokeWidth={1} />
          <ReferenceLine y={38.5} stroke="#E05252" strokeDasharray="4 4" strokeWidth={1} />

          <Line
            type="monotone"
            dataKey="value"
            stroke="#4CAF82"
            strokeWidth={2}
            dot={<StatusDot />}
            activeDot={{ r: 5, fill: '#4CAF82' }}
            animationDuration={800}
            animationEasing="ease-in-out"
          />
        </LineChart>
      </ResponsiveContainer>

      {/* Legend threshold */}
      <div className="flex gap-4 mt-2">
        <span className="flex items-center gap-1 text-xs text-healy-slate">
          <span className="w-4 border-t border-dashed border-healy-warning" />
          Waspada 37.5°C
        </span>
        <span className="flex items-center gap-1 text-xs text-healy-slate">
          <span className="w-4 border-t border-dashed border-healy-critical" />
          Kritis 38.5°C
        </span>
      </div>
    </div>
  )
}
```

### 6.4 New Component: HeartRateChart.tsx

```typescript
// src/components/features/HeartRateChart.tsx
'use client'

import {
  LineChart, Line, XAxis, YAxis, CartesianGrid,
  Tooltip, ReferenceLine, ResponsiveContainer
} from 'recharts'
import { TelemetryChartPoint } from '@/types/telemetry'

interface HeartRateChartProps {
  data: TelemetryChartPoint[]
  loading?: boolean
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (!active || !payload?.length) return null
  return (
    <div className="bg-healy-surface border border-healy-border rounded-lg px-3 py-2 shadow-card text-xs">
      <p className="font-mono text-healy-slate mb-1">{label}</p>
      <p className="font-mono font-semibold text-healy-graphite">
        {payload[0].value} bpm
      </p>
    </div>
  )
}

export function HeartRateChart({ data, loading }: HeartRateChartProps) {
  if (loading) return <div className="h-48 bg-healy-bg-alt rounded-card animate-pulse" />

  return (
    <div className="bg-healy-surface border border-healy-border rounded-card p-4 shadow-card">
      <div className="flex items-center justify-between mb-3">
        <h4 className="font-display font-semibold text-sm text-healy-graphite">
          Heart Rate Trend
        </h4>
        <span className="text-xs font-mono text-healy-slate">bpm</span>
      </div>

      <ResponsiveContainer width="100%" height={160}>
        <LineChart data={data} margin={{ top: 4, right: 8, bottom: 0, left: -20 }}>
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="var(--color-border)"
            opacity={0.5}
          />
          <XAxis
            dataKey="timestamp"
            tick={{ fontSize: 10, fontFamily: 'JetBrains Mono', fill: 'var(--color-slate)' }}
            tickLine={false}
            axisLine={false}
          />
          <YAxis
            domain={[40, 180]}
            tick={{ fontSize: 10, fontFamily: 'JetBrains Mono', fill: 'var(--color-slate)' }}
            tickLine={false}
            axisLine={false}
          />
          <Tooltip content={<CustomTooltip />} />

          {/* Normal range reference band */}
          <ReferenceLine y={60}  stroke="#4CAF82" strokeDasharray="3 3" strokeWidth={1} opacity={0.5} />
          <ReferenceLine y={100} stroke="#4CAF82" strokeDasharray="3 3" strokeWidth={1} opacity={0.5} />

          <Line
            type="monotone"
            dataKey="value"
            stroke="#F5A623"
            strokeWidth={2}
            dot={false}
            activeDot={{ r: 5, fill: '#F5A623' }}
            animationDuration={800}
            animationEasing="ease-in-out"
          />
        </LineChart>
      </ResponsiveContainer>

      <div className="flex gap-4 mt-2">
        <span className="flex items-center gap-1 text-xs text-healy-slate">
          <span className="w-4 border-t border-dashed border-healy-sage" />
          Normal 60–100 bpm
        </span>
      </div>
    </div>
  )
}
```

### 6.5 New Component: SpO2Chart.tsx

```typescript
// src/components/features/SpO2Chart.tsx
'use client'

import {
  LineChart, Line, XAxis, YAxis, CartesianGrid,
  Tooltip, ReferenceLine, ResponsiveContainer
} from 'recharts'
import { TelemetryChartPoint } from '@/types/telemetry'

interface SpO2ChartProps {
  data: TelemetryChartPoint[]
  loading?: boolean
}

const CustomTooltip = ({ active, payload, label }: any) => {
  if (!active || !payload?.length) return null
  return (
    <div className="bg-healy-surface border border-healy-border rounded-lg px-3 py-2 shadow-card text-xs">
      <p className="font-mono text-healy-slate mb-1">{label}</p>
      <p className="font-mono font-semibold text-healy-graphite">
        {payload[0].value}%
      </p>
    </div>
  )
}

export function SpO2Chart({ data, loading }: SpO2ChartProps) {
  if (loading) return <div className="h-48 bg-healy-bg-alt rounded-card animate-pulse" />

  return (
    <div className="bg-healy-surface border border-healy-border rounded-card p-4 shadow-card">
      <div className="flex items-center justify-between mb-3">
        <h4 className="font-display font-semibold text-sm text-healy-graphite">
          SpO2 Trend
        </h4>
        <span className="text-xs font-mono text-healy-slate">%</span>
      </div>

      <ResponsiveContainer width="100%" height={160}>
        <LineChart data={data} margin={{ top: 4, right: 8, bottom: 0, left: -20 }}>
          <CartesianGrid
            strokeDasharray="3 3"
            stroke="var(--color-border)"
            opacity={0.5}
          />
          <XAxis
            dataKey="timestamp"
            tick={{ fontSize: 10, fontFamily: 'JetBrains Mono', fill: 'var(--color-slate)' }}
            tickLine={false}
            axisLine={false}
          />
          <YAxis
            domain={[80, 100]}
            tick={{ fontSize: 10, fontFamily: 'JetBrains Mono', fill: 'var(--color-slate)' }}
            tickLine={false}
            axisLine={false}
          />
          <Tooltip content={<CustomTooltip />} />

          {/* Threshold reference lines */}
          <ReferenceLine y={95} stroke="#F5A623" strokeDasharray="4 4" strokeWidth={1} />
          <ReferenceLine y={91} stroke="#E05252" strokeDasharray="4 4" strokeWidth={1} />

          <Line
            type="monotone"
            dataKey="value"
            stroke="#0D9488"
            strokeWidth={2}
            dot={false}
            activeDot={{ r: 5, fill: '#0D9488' }}
            animationDuration={800}
            animationEasing="ease-in-out"
          />
        </LineChart>
      </ResponsiveContainer>

      <div className="flex gap-4 mt-2">
        <span className="flex items-center gap-1 text-xs text-healy-slate">
          <span className="w-4 border-t border-dashed border-healy-warning" />
          Waspada &lt;95%
        </span>
        <span className="flex items-center gap-1 text-xs text-healy-slate">
          <span className="w-4 border-t border-dashed border-healy-critical" />
          Kritis &lt;91%
        </span>
      </div>
    </div>
  )
}
```

### 6.6 Modifikasi: app/history/page.tsx

```typescript
// src/app/history/page.tsx — MODIFIKASI: ganti combined chart → 3 chart terpisah

import { TemperatureChart } from '@/components/features/TemperatureChart'
import { HeartRateChart } from '@/components/features/HeartRateChart'
import { SpO2Chart } from '@/components/features/SpO2Chart'

// Sebelumnya: satu <MultiSeriesChart data={historyData} />
// Sekarang:
export default function HistoryPage() {
  const { data, loading, range, setRange } = useHistoryData()

  // Transform data mentah untuk tiap chart
  const tempData = data.map(r => ({
    timestamp: formatTimestamp(r.recorded_at, range),
    value: r.temperature,
    status: r.temp_status,
  }))

  const bpmData = data.map(r => ({
    timestamp: formatTimestamp(r.recorded_at, range),
    value: r.bpm,
    status: r.bpm > 100 ? 'WARNING' : 'NORMAL', // Evaluasi lokal
  }))

  const spo2Data = data.map(r => ({
    timestamp: formatTimestamp(r.recorded_at, range),
    value: r.spo2,
    status: r.spo2_status,
  }))

  return (
    <div className="p-6 space-y-6">
      {/* Filter Bar */}
      <div className="flex items-center gap-2">
        {(['1h', '6h', '24h', '7d'] as const).map(r => (
          <button
            key={r}
            onClick={() => setRange(r)}
            className={`px-3 py-1.5 text-xs font-medium rounded-lg transition-colors ${
              range === r
                ? 'bg-healy-sage text-white'
                : 'border border-healy-border text-healy-graphite hover:bg-healy-bg-alt'
            }`}
          >
            {r === '1h' ? '1 Jam' : r === '6h' ? '6 Jam' : r === '24h' ? '24 Jam' : '7 Hari'}
          </button>
        ))}
      </div>

      {/* Tiga Chart Terpisah — Grid 1 kolom, stacked vertikal */}
      <div className="grid grid-cols-1 gap-4">
        <TemperatureChart data={tempData} loading={loading} />
        <HeartRateChart   data={bpmData}  loading={loading} />
        <SpO2Chart        data={spo2Data} loading={loading} />
      </div>

      {/* Tabel Riwayat — tidak berubah dari v3.0 */}
      {/* ... */}
    </div>
  )
}
```

---

## 7. F-04: Interactive AI Chatbot — Context-Aware

### 7.1 Keputusan Arsitektur

Chatbot diimplementasikan sebagai **sliding panel** yang muncul dari kanan layar, bukan modal atau halaman terpisah. Justifikasi: dashboard tetap visible di background, sehingga operator bisa melihat data sensor sambil bertanya ke AI — ini penting secara klinis.

**Mekanisme context injection:** Setiap kali user mengirim pesan, sistem secara otomatis menyisipkan `system prompt` yang berisi data telemetri terakhir. User tidak perlu menyebutkan angka — bot sudah "tahu" kondisi pasien saat pesan dikirim.

**Arsitektur percakapan:** Multi-turn conversation dengan message history. Backend tidak terlibat — semua state percakapan disimpan di React state (`useState`). Ini konsisten dengan keputusan v3.0 bahwa Groq dipanggil langsung dari browser.

### 7.2 New Types: src/types/chat.ts

```typescript
// src/types/chat.ts — NEW v4.0.0

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
```

### 7.3 New Hook: useChatbot.ts

```typescript
// src/hooks/useChatbot.ts — NEW v4.0.0

import { useState, useCallback, useRef } from 'react'
import { ChatMessage, ChatContext } from '@/types/chat'
import { getStoredGroqKey, validateGroqKey } from '@/lib/groq-client'

const GROQ_API_URL  = 'https://api.groq.com/openai/v1/chat/completions'
const GROQ_MODEL    = 'llama-3.1-8b-instant'
const MAX_HISTORY   = 20  // Batas pesan dalam satu sesi untuk hemat token

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

    const payload = {
      model: GROQ_MODEL,
      messages: [
        { role: 'system', content: buildSystemPrompt(context) },
        ...historyForApi,
        { role: 'user', content: trimmed },
      ],
      max_tokens: 512,
      stream: true,
      temperature: 0.4,
    }

    abortRef.current = new AbortController()

    try {
      const response = await fetch(GROQ_API_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${apiKey}`,
        },
        body: JSON.stringify(payload),
        signal: abortRef.current.signal,
      })

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

    } catch (err: any) {
      if (err.name === 'AbortError') return // User menutup panel, abaikan

      setMessages(prev => prev.map(m =>
        m.id === assistantMsgId
          ? {
              ...m,
              content: `Gagal mendapatkan respons: ${err.message}`,
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
```

### 7.4 New Component: AIChatPanel.tsx

```typescript
// src/components/features/AIChatPanel.tsx — NEW v4.0.0
'use client'

import { useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import {
  X, Robot, PaperPlaneRight, Spinner, Trash, Warning
} from '@phosphor-icons/react'
import { ChatMessage, ChatContext } from '@/types/chat'

interface AIChatPanelProps {
  isOpen:       boolean
  onClose:      () => void
  messages:     ChatMessage[]
  isLoading:    boolean
  inputValue:   string
  onInputChange:(value: string) => void
  onSend:       (input: string) => void
  onClear:      () => void
  context:      ChatContext | null  // Data sensor saat ini
}

export function AIChatPanel({
  isOpen, onClose, messages, isLoading,
  inputValue, onInputChange, onSend, onClear, context
}: AIChatPanelProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null)

  // Auto-scroll ke pesan terbaru
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  const handleSend = () => {
    if (!inputValue.trim() || !context) return
    onSend(inputValue)
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/20 backdrop-blur-sm z-40"
          />

          {/* Sliding Panel */}
          <motion.div
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ type: 'spring', damping: 30, stiffness: 300 }}
            className="fixed right-0 top-0 h-full w-96 max-w-full
                       bg-healy-surface border-l border-healy-border
                       shadow-2xl z-50 flex flex-col"
          >
            {/* Header Panel */}
            <div className="flex items-center justify-between px-4 py-3 border-b border-healy-border">
              <div className="flex items-center gap-2">
                <Robot size={18} weight="duotone" className="text-healy-ai-accent" />
                <span className="font-display font-semibold text-sm text-healy-graphite">
                  HEALY AI Chat
                </span>
                <span className="text-xs bg-healy-ai-accent/15 text-healy-ai-accent
                                 px-2 py-0.5 rounded-full font-mono">
                  Groq
                </span>
              </div>
              <div className="flex items-center gap-1">
                <button
                  onClick={onClear}
                  title="Hapus riwayat chat"
                  className="p-1.5 rounded-lg text-healy-slate hover:text-healy-critical
                             hover:bg-red-50 dark:hover:bg-red-950 transition-colors"
                >
                  <Trash size={14} />
                </button>
                <button
                  onClick={onClose}
                  className="p-1.5 rounded-lg text-healy-slate hover:text-healy-graphite
                             hover:bg-healy-bg-alt transition-colors"
                >
                  <X size={16} />
                </button>
              </div>
            </div>

            {/* Context Badge — tampilkan data sensor aktif */}
            {context && (
              <div className="px-4 py-2 bg-healy-bg-alt border-b border-healy-border">
                <p className="text-xs text-healy-slate font-mono">
                  Konteks: {context.temperature}°C · {context.bpm} bpm · SpO2 {context.spo2}%
                  <span className={`ml-2 font-semibold ${
                    context.overallStatus === 'CRITICAL' ? 'text-healy-critical' :
                    context.overallStatus === 'WARNING'  ? 'text-healy-warning'  :
                    'text-healy-sage'
                  }`}>
                    [{context.overallStatus}]
                  </span>
                </p>
              </div>
            )}

            {!context && (
              <div className="px-4 py-2 bg-healy-bg-alt border-b border-healy-border flex items-center gap-2">
                <Warning size={12} className="text-healy-warning" />
                <p className="text-xs text-healy-slate">
                  Perangkat offline — AI menjawab tanpa data real-time
                </p>
              </div>
            )}

            {/* Message List */}
            <div className="flex-1 overflow-y-auto px-4 py-3 space-y-3">
              {messages.length === 0 && (
                <div className="text-center py-8 space-y-2">
                  <Robot size={32} className="mx-auto text-healy-ai-accent opacity-50" />
                  <p className="text-xs text-healy-slate">
                    Tanyakan kondisi kesehatan pasien, interpretasi data, atau rekomendasi tindakan.
                  </p>
                </div>
              )}

              {messages.map(msg => (
                <div
                  key={msg.id}
                  className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}
                >
                  <div className={`max-w-[85%] rounded-2xl px-3 py-2 text-xs leading-relaxed ${
                    msg.role === 'user'
                      ? 'bg-healy-sage text-white rounded-br-sm'
                      : msg.isError
                        ? 'bg-red-50 dark:bg-red-950 text-healy-critical border border-healy-critical/20 rounded-bl-sm'
                        : 'bg-healy-bg-alt text-healy-graphite border border-healy-border rounded-bl-sm'
                  }`}>
                    <p className="font-body whitespace-pre-wrap">{msg.content}</p>
                    {msg.isStreaming && (
                      <span className="inline-block w-1 h-3 bg-healy-ai-accent ml-0.5 animate-pulse" />
                    )}
                  </div>
                </div>
              ))}

              {isLoading && messages[messages.length - 1]?.content === '' && (
                <div className="flex justify-start">
                  <div className="bg-healy-bg-alt border border-healy-border rounded-2xl
                                  rounded-bl-sm px-3 py-2">
                    <Spinner size={14} className="animate-spin text-healy-ai-accent" />
                  </div>
                </div>
              )}

              <div ref={messagesEndRef} />
            </div>

            {/* Input Area */}
            <div className="px-4 py-3 border-t border-healy-border">
              <div className="flex gap-2">
                <textarea
                  value={inputValue}
                  onChange={e => onInputChange(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder="Tanya tentang kondisi pasien..."
                  rows={2}
                  className="flex-1 text-xs font-body resize-none px-3 py-2 rounded-xl
                             border border-healy-border bg-healy-bg
                             text-healy-graphite placeholder:text-healy-slate
                             focus:outline-none focus:ring-2 focus:ring-healy-ai-accent/30
                             focus:border-healy-ai-accent transition-all"
                />
                <button
                  onClick={handleSend}
                  disabled={!inputValue.trim() || isLoading || !context}
                  className="self-end p-2.5 rounded-xl bg-healy-ai-accent text-white
                             disabled:opacity-40 disabled:cursor-not-allowed
                             hover:bg-teal-700 transition-colors"
                >
                  <PaperPlaneRight size={16} weight="fill" />
                </button>
              </div>
              <p className="text-xs text-healy-slate mt-1.5 text-center">
                Enter untuk kirim · Shift+Enter untuk baris baru
              </p>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}
```

### 7.5 Modifikasi: NavSidebar.tsx — Tambahkan ChatbotTrigger

```typescript
// Tambahkan di bagian atas NavSidebar.tsx — import dan trigger button

import { ChatCircleDots } from '@phosphor-icons/react'

// Props baru untuk NavSidebar
interface NavSidebarProps {
  onOpenChat: () => void   // Callback dari dashboard page
}

// Di dalam JSX sidebar, tambahkan tombol ini di area tools atau footer:
<button
  onClick={onOpenChat}
  className="flex items-center gap-2 w-full px-3 py-2 rounded-lg
             text-healy-graphite hover:bg-healy-bg-alt
             hover:text-healy-sage transition-colors"
>
  <ChatCircleDots size={18} weight="duotone" className="text-healy-ai-accent" />
  <span className="text-sm font-medium">AI Health Chat</span>
  <span className="ml-auto text-xs bg-healy-ai-accent/15 text-healy-ai-accent
                   px-1.5 py-0.5 rounded-full font-mono">AI</span>
</button>
```

### 7.6 Integrasi di dashboard/page.tsx

```typescript
// src/app/dashboard/page.tsx — TAMBAHKAN integrasi chatbot

import { useChatbot } from '@/hooks/useChatbot'
import { AIChatPanel } from '@/components/features/AIChatPanel'
import { ChatContext } from '@/types/chat'

export default function DashboardPage() {
  const { data, conn, deviceOnline } = useTelemetry()

  // Chatbot hook
  const {
    isOpen, openChat, closeChat,
    messages, isLoading,
    inputValue, setInputValue,
    sendMessage, clearHistory,
  } = useChatbot()

  // Rakit context dari data sensor terbaru
  const chatContext: ChatContext | null = data ? {
    temperature:   data.sensor.temperature,
    bpm:           data.sensor.bpm,
    spo2:          data.sensor.spo2,
    tempStatus:    data.status.temperature,
    spo2Status:    data.status.spo2,
    overallStatus: data.status.overall,
    timestamp:     new Date(data.timestamp).toLocaleString('id-ID'),
  } : null

  return (
    <div className="flex h-screen">
      {/* Sidebar — pass onOpenChat prop */}
      <NavSidebar onOpenChat={openChat} />

      {/* Konten utama */}
      <main className="flex-1 overflow-y-auto p-6">
        {/* ... semua komponen dashboard existing ... */}
      </main>

      {/* AI Chat Panel — sliding dari kanan */}
      <AIChatPanel
        isOpen={isOpen}
        onClose={closeChat}
        messages={messages}
        isLoading={isLoading}
        inputValue={inputValue}
        onInputChange={setInputValue}
        onSend={(input) => sendMessage(input, chatContext!)}
        onClear={clearHistory}
        context={chatContext}
      />
    </div>
  )
}
```

---

## 8. Design Token Master — Versi Final v4.0.0

Ini adalah versi **definitif dan menggantikan semua versi sebelumnya** di `design-tokens.ts`.

```typescript
// src/constants/design-tokens.ts — MASTER v4.0.0

export const COLORS = {
  // ─── BACKGROUND ───
  bgPrimary:     '#F8FAFB',   // Glacial White
  bgSecondary:   '#EFF4F2',   // Soft Mist
  surface:       '#FFFFFF',   // Ice White

  // ─── BRAND ───
  brandPrimary:  '#4CAF82',   // Sage Green
  brandSecondary:'#2E8B62',   // Deep Sage
  brandAccent:   '#A8DFCC',   // Mint Glow

  // ─── TEXT ───
  textPrimary:   '#1A2633',   // Graphite
  textSecondary: '#5A7080',   // Slate

  // ─── STATUS ───
  statusNormal:  '#4CAF82',
  statusWarning: '#F5A623',   // Amber Soft
  statusCritical:'#E05252',   // Coral Red

  // ─── BORDER ───
  border:        '#D4E8DF',   // Pale Sage

  // ─── DEVICE LED ───
  deviceOnline:  '#22C55E',   // Green-500
  deviceOffline: '#EF4444',   // Red-500

  // ─── AI TIER — v4.0.0 (Menggantikan Violet v3.0.0) ───
  aiAccent:      '#0D9488',   // Clinical Teal — CHANGED from #7C3AED
  aiSurface:     '#F0FDFC',   // Teal-50 — CHANGED from #F5F3FF

} as const

export const FONTS = {
  display: 'Exo 2',
  body:    'DM Sans',
  mono:    'JetBrains Mono',
} as const

// Threshold constants — single source of truth
export const THRESHOLDS = {
  temp: {
    normalMin: 36.5,
    normalMax: 37.5,
    warnMax:   38.5,
  },
  spo2: {
    normalMin: 95,
    warnMin:   91,
  },
  bpm: {
    normalMin: 60,
    normalMax: 100,
  },
} as const
```

---

## 9. groq-client.ts — Revisi Lengkap v4.0.0

Versi ini menggantikan `groq-client.ts` dari v3.0.0. Penambahan: fungsi `callGroqChat` untuk multi-turn conversation yang digunakan oleh `useChatbot`.

```typescript
// src/lib/groq-client.ts — REVISI LENGKAP v4.0.0

const GROQ_API_URL  = 'https://api.groq.com/openai/v1/chat/completions'
const GROQ_MODEL    = 'llama-3.1-8b-instant'

// ─── KEY MANAGEMENT ───

export function getStoredGroqKey(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('healy_groq_api_key')
}

export function setStoredGroqKey(key: string): void {
  localStorage.setItem('healy_groq_api_key', key.trim())
}

export function removeStoredGroqKey(): void {
  localStorage.removeItem('healy_groq_api_key')
}

export function validateGroqKey(key: string): boolean {
  return /^gsk_[a-zA-Z0-9]{50,}$/.test(key.trim())
}

// ─── INSIGHT (SINGLE-SHOT) — digunakan AIInsightCard & useAutoNarrative ───

export interface GroqInsightRequest {
  temperature: number
  bpm:         number
  spo2:        number
  avgTemp?:    number
  avgBpm?:     number
  avgSpo2?:    number
  status:      string
}

export function buildInsightPrompt(data: GroqInsightRequest): string {
  const avg = data.avgTemp
    ? `Rata-rata 15 menit: Suhu ${data.avgTemp}°C, BPM ${data.avgBpm}, SpO2 ${data.avgSpo2}%.`
    : ''

  return `Kamu adalah asisten medis AI untuk sistem HEALY.

Data biometrik pasien:
- Suhu: ${data.temperature}°C
- BPM: ${data.bpm}
- SpO2: ${data.spo2}%
- Status: ${data.status}
${avg}

Berikan analisis klinis singkat (2–3 kalimat) dalam Bahasa Indonesia.
Format: interpretasi kondisi → potensi risiko (jika ada) → satu rekomendasi.
Bahasa mudah dipahami keluarga pasien. Tanpa disclaimer atau pengantar.`
}

// Single-shot call — digunakan untuk insight & auto-narrative
export async function callGroqInsight(
  prompt:   string,
  apiKey:   string,
  onChunk?: (text: string) => void,
  signal?:  AbortSignal
): Promise<string> {
  const response = await fetch(GROQ_API_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${apiKey}`,
    },
    body: JSON.stringify({
      model:       GROQ_MODEL,
      messages:    [{ role: 'user', content: prompt }],
      max_tokens:  200,
      stream:      !!onChunk,
      temperature: 0.3,
    }),
    signal,
  })

  if (!response.ok) {
    const err = await response.json()
    throw new Error(err.error?.message || `Groq error ${response.status}`)
  }

  if (onChunk && response.body) {
    const reader  = response.body.getReader()
    const decoder = new TextDecoder()
    let fullText  = ''

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
          const text   = parsed.choices?.[0]?.delta?.content || ''
          fullText    += text
          onChunk(fullText)
        } catch { /* skip */ }
      }
    }
    return fullText
  }

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
```

---

## 10. Known Gotchas v4.0.0

Semua Gotcha dari v3.0.0 (01–07) tetap berlaku. Tambahan:

---

### GOTCHA-08: Dark Mode — Hydration Mismatch

**Status:** ACTIVE — wajib diikuti.

**Masalah:** Next.js merender HTML di server tanpa mengetahui preferensi tema user (tersimpan di `localStorage` atau header browser). Jika komponen merender kondisional berdasarkan tema sebelum hydration selesai, React akan throw hydration mismatch error.

**Solusi wajib:**
```typescript
// Untuk komponen yang render berbeda berdasarkan tema:
const [mounted, setMounted] = useState(false)
useEffect(() => { setMounted(true) }, [])
if (!mounted) return <div className="w-8 h-8" /> // Placeholder dimensi sama
```

**Contoh komponen yang wajib pakai pattern ini:** `ThemeToggle`, dan komponen apapun yang menggunakan `useTheme()` secara kondisional.

---

### GOTCHA-09: Chatbot Context — Data Stale saat Panel Tertutup

**Status:** ACTIVE — behavior yang disengaja, bukan bug.

**Masalah:** Ketika user membuka chat panel, lalu data sensor berubah di background, `context` yang diinjeksikan ke system prompt adalah data **saat pesan dikirim**, bukan data saat panel dibuka.

**Ini adalah behavior yang benar:** setiap pesan mendapat snapshot konteks terbaru saat dikirim. Ini lebih akurat dibanding snapshot saat panel dibuka.

**Yang harus diperhatikan:** pastikan `chatContext` di `dashboard/page.tsx` selalu berasal dari state terbaru (`data` dari `useTelemetry`), bukan disimpan di `useRef` atau state stale.

---

### GOTCHA-10: CSS Variable Format — Tailwind v4 Adaptation

**Status:** RESOLVED — Tidak berlaku untuk Tailwind CSS v4.

**Catatan Blueprint Asli:** Dokumen ini awalnya menspesifikasikan Tailwind v3 approach: `tailwind.config.ts` dengan `darkMode: 'class'` dan format `rgb(var(--color-x) / <alpha-value>)` yang membutuhkan CSS variable dalam format channel terpisah (R G B).

**Realita Project:** HEALY menggunakan **Tailwind CSS v4** (`@tailwindcss/postcss` v4, `@import "tailwindcss"`, `@theme` directive di `globals.css`). Pada Tailwind v4:
- Tidak ada `tailwind.config.ts` — konfigurasi sepenuhnya di CSS via `@theme`.
- `darkMode: 'class'` adalah **default** — tidak perlu dikonfigurasi manual.
- Alpha-value opacity modifiers (e.g. `bg-healy-sage/50`) bekerja langsung dengan format hex — **tidak perlu konversi ke R G B channel**.

**Implementasi F-01 menggunakan:** Hex color langsung di `:root` dan `.dark` blocks, dengan `@theme` tokens mereferensi nama yang sama. Ini sepenuhnya kompatibel dan teruji (`next build` berhasil).

---

## 11. Rencana Kerja: Fase 12–15

### FASE 12 — Foundation Dark Mode (F-01)

**Step 12.1** — Install `next-themes`: `pnpm add next-themes`

**Step 12.2** — Konversi semua hex di `globals.css` ke format channel (`76 175 130`)

**Step 12.3** — Buat `src/app/providers.tsx` dengan `ThemeProvider` sesuai Section 4.2

**Step 12.4** — Modifikasi `app/layout.tsx` sesuai Section 4.3 — tambahkan `suppressHydrationWarning`

**Step 12.5** — Update `tailwind.config.ts` — tambahkan `darkMode: 'class'` dan format `rgb(var(...))` sesuai Section 4.5

**Step 12.6** — Buat `src/components/ui/theme-toggle.tsx` sesuai Section 4.6

**Step 12.7** — Tambahkan `<ThemeToggle />` di Dashboard header

**Step 12.8** — Verifikasi visual: toggle tema, pastikan semua card dan komponen merespons

---

### FASE 13 — AI Accent Token Update (F-02)

**Step 13.1** — Update `design-tokens.ts` sesuai Section 8 (aiAccent, aiSurface)

**Step 13.2** — Update CSS variable di `globals.css` (`:root` dan `.dark`)

**Step 13.3** — Cari & ganti `hover:bg-violet-600` → `hover:bg-teal-700` di semua file

**Step 13.4** — Verifikasi visual: AIInsightCard, NavSidebar GroqKeySection, AlertFeed narrative

---

### FASE 14 — Separated History Charts (F-03)

**Step 14.1** — Tambahkan `TelemetryChartPoint` dan `ChartConfig` ke `src/types/telemetry.ts`

**Step 14.2** — Buat `TemperatureChart.tsx` sesuai Section 6.3

**Step 14.3** — Buat `HeartRateChart.tsx` sesuai Section 6.4

**Step 14.4** — Buat `SpO2Chart.tsx` sesuai Section 6.5

**Step 14.5** — Refactor `app/history/page.tsx` sesuai Section 6.6 — hapus komponen chart lama, gunakan 3 chart baru

**Step 14.6** — Verifikasi: filter 1h/6h/24h/7d berfungsi, reference line threshold tampil, animasi draw berjalan

---

### FASE 15 — AI Chatbot Context-Aware (F-04)

**Step 15.1** — Buat `src/types/chat.ts` sesuai Section 7.2

**Step 15.2** — Update `src/lib/groq-client.ts` ke versi Section 9 (tambahkan `callGroqChat`)

**Step 15.3** — Buat `src/hooks/useChatbot.ts` sesuai Section 7.3

**Step 15.4** — Buat `src/components/features/AIChatPanel.tsx` sesuai Section 7.4

**Step 15.5** — Modifikasi `NavSidebar.tsx` — tambahkan `ChatbotTrigger` sesuai Section 7.5

**Step 15.6** — Integrasi di `dashboard/page.tsx` sesuai Section 7.6

**Step 15.7** — Testing skenario:
- Buka chat panel, verifikasi context badge menampilkan data sensor terbaru
- Kirim pertanyaan: *"Apakah suhu pasien normal?"* — verifikasi AI mereferensi data aktual
- Set mock ke CRITICAL, tanya: *"Apa yang harus saya lakukan sekarang?"* — verifikasi AI merespons sesuai kondisi kritis
- Verifikasi streaming berjalan: teks muncul karakter per karakter
- Tutup panel saat streaming berlangsung — verifikasi tidak ada memory leak (AbortController berjalan)
- Verifikasi `clearHistory` menghapus semua pesan

---

## 12. Codebase Progress Tracker — Updated

### Status Backend (Tidak Berubah dari v3.0.0)

Semua file backend berstatus ✅ DONE. Tidak ada perubahan di Fase 12–15.

### Status Frontend — Delta v4.0.0

| File | Status | Keterangan |
|---|---|---|
| `app/providers.tsx` | ✅ DONE | NEW — ThemeProvider wrapper (F-01) |
| `app/layout.tsx` | ✅ DONE | MODIFY — Providers wrapper + suppressHydrationWarning (F-01) |
| `app/globals.css` | ✅ DONE | MODIFY — dark mode `.dark` vars + AI accent Teal + theme transition (F-01, F-02) |
| `components/ui/theme-toggle.tsx` | ✅ DONE | NEW — Light/Dark toggle button, lucide-react Sun/Moon (F-01) |
| `tailwind.config.ts` | ⬜ N/A | Tailwind v4 — tidak ada config file, gunakan `@theme` di globals.css |
| `app/history/page.tsx` | ⬜ TODO | MODIFY — 3 chart terpisah (F-03) |
| `constants/design-tokens.ts` | ✅ DONE | MODIFY — aiAccent teal, master v4.0 (F-02) |
| `types/chat.ts` | ⬜ TODO | NEW — ChatMessage, ChatContext (F-04) |
| `types/telemetry.ts` | ⬜ TODO | MODIFY — tambah TelemetryChartPoint (F-03) |
| `lib/groq-client.ts` | ⬜ TODO | MODIFY — tambah callGroqChat (F-04) |
| `hooks/useChatbot.ts` | ⬜ TODO | NEW — chatbot state management (F-04) |
| `components/features/NavSidebar.tsx` | ⬜ TODO | MODIFY — ChatbotTrigger button (F-04) |
| `components/features/AIInsightCard.tsx` | ✅ DONE | MODIFY — teal accent (F-02) |
| `components/features/AlertFeed.tsx` | ✅ DONE | MODIFY — teal accent (F-02) |
| `components/features/AIChatPanel.tsx` | ⬜ TODO | NEW — sliding chatbot panel (F-04) |
| `components/features/TemperatureChart.tsx` | ⬜ TODO | NEW — standalone temp chart (F-03) |
| `components/features/HeartRateChart.tsx` | ⬜ TODO | NEW — standalone BPM chart (F-03) |
| `components/features/SpO2Chart.tsx` | ⬜ TODO | NEW — standalone SpO2 chart (F-03) |
| `app/dashboard/page.tsx` | ⬜ TODO | MODIFY — integrasi chatbot (F-04) |

---

## 13. Prompt Inisialisasi AI — v4.0.0

### Prompt 0 — Grounding v4.0 (Gantikan prompt v3.0)

```
You are a Senior Full-Stack Engineer on Project HEALY v4.0.0.

I am attaching HEALY_Master_Blueprint_v4.0.0.md.

Read it completely. Confirm understanding:
1. What are the 4 new features in v4.0.0? (F-01 to F-04)
2. Which files are NEW vs MODIFY in the Progress Tracker?
3. What is the CSS variable format required for Tailwind dark mode alpha support?
4. What changed in the AI accent color and why?
5. How does the AI chatbot inject patient context? Where is it assembled?
6. How many separate charts are on the History page, and what are their Y-axis domains?

Do NOT generate code yet.

[ATTACH: HEALY_Master_Blueprint_v4.0.0.md]
```

---

### Prompt F-01 — Dark Mode Implementation

```
Based on HEALY Master Blueprint v4.0 Section 4, implement Light/Dark Mode.

Execute in this order:

Step 1: globals.css
Convert ALL hex color values in :root to channel format (e.g. --color-sage: 76 175 130)
Add .dark {} block with all overrides from Section 4.4
Add global transition block (200ms for bg/border/color)

Step 2: tailwind.config.ts
Add darkMode: 'class' at root level.
Update healy color definitions to use rgb(var(--color-x) / <alpha-value>) format.
Exact config from Section 4.5.

Step 3: src/app/providers.tsx
Full implementation from Section 4.2.
ThemeProvider: attribute="class", defaultTheme="light", enableSystem=true.

Step 4: src/app/layout.tsx
Wrap children with <Providers>.
Add suppressHydrationWarning to <html>.
Section 4.3.

Step 5: src/components/ui/theme-toggle.tsx
Full implementation from Section 4.6.
Use mounted state to prevent hydration mismatch.

Output: All 5 files. No placeholder comments.
Verify: Does each component use healy-* Tailwind tokens? No hardcoded colors.
```

---

### Prompt F-03 — Separated History Charts

```
Based on HEALY Master Blueprint v4.0 Section 6, build the separated chart components.

Step 1: Update src/types/telemetry.ts
Add TelemetryChartPoint and ChartConfig interfaces from Section 6.2.

Step 2: Create TemperatureChart.tsx (Section 6.3)
Y-axis domain: [35, 40]
Reference lines: y=37.5 (warning amber dashed), y=38.5 (critical coral dashed)
Line color: healy-sage (#4CAF82)
Custom StatusDot that colors each point by its status field.
Custom Tooltip with HEALY surface card styling.

Step 3: Create HeartRateChart.tsx (Section 6.4)
Y-axis domain: [40, 180]
Reference lines: y=60 and y=100 (normal range, sage, low opacity)
Line color: healy-warning (#F5A623)
No colored dots (dot={false}) for performance at high data volume.

Step 4: Create SpO2Chart.tsx (Section 6.5)
Y-axis domain: [80, 100]
Reference lines: y=95 (warning amber dashed), y=91 (critical coral dashed)
Line color: healy-ai-accent (Clinical Teal #0D9488)

Step 5: Refactor src/app/history/page.tsx (Section 6.6)
Remove old MultiSeriesChart component reference.
Add three chart components in vertical stack (grid-cols-1).
Data transformation: map TelemetryRecord[] to TelemetryChartPoint[] for each param.

All charts: ResponsiveContainer height=160, animationDuration=800,
loading skeleton: h-48 bg-healy-bg-alt rounded-card animate-pulse.

Output: 5 files with full implementation.
```

---

### Prompt F-04 — AI Chatbot Context-Aware

```
Based on HEALY Master Blueprint v4.0 Section 7 and groq-client.ts Section 9,
implement the Interactive AI Chatbot.

Step 1: src/types/chat.ts — Section 7.2
ChatMessage interface with: id, role, content, timestamp, isStreaming, isError.
ChatContext interface with sensor data fields.

Step 2: src/lib/groq-client.ts — UPDATE to Section 9 full version
Add callGroqChat(messages, apiKey, signal) function that returns raw Response.
Keep all existing functions (callGroqInsight, buildInsightPrompt, etc.) intact.

Step 3: src/hooks/useChatbot.ts — Section 7.3
Implements: isOpen, openChat, closeChat, messages, isLoading, 
            inputValue, setInputValue, sendMessage, clearHistory.
sendMessage: validates key, appends user message, creates streaming assistant placeholder,
streams via callGroqChat, updates message content per chunk, finalizes isStreaming=false.
AbortController: abort on closeChat to prevent memory leaks.
buildSystemPrompt: injects ChatContext data as Bahasa Indonesia system prompt.
MAX_HISTORY = 20.

Step 4: src/components/features/AIChatPanel.tsx — Section 7.4
Sliding panel from right side (framer-motion x: 100% → 0)
Backdrop with blur
Context badge showing latest sensor data
Message bubbles: user=sage-green right, AI=bg-alt left, error=red-tinted
Streaming cursor: inline-block pulse span
Input: textarea rows=2, Enter to send, Shift+Enter newline
Auto-scroll to latest message via useRef

Step 5: Modify NavSidebar.tsx — Section 7.5
Add onOpenChat prop. Add ChatbotTrigger button with ChatCircleDots icon.

Step 6: Modify dashboard/page.tsx — Section 7.6
Instantiate useChatbot. Build chatContext from latest useTelemetry data.
Render AIChatPanel. Pass onOpenChat to NavSidebar.

Output: 6 files. Full implementation.
Test criteria: context badge updates when sensor data changes,
streaming text appears character by character, abort works on panel close.
```

---

*HEALY Master Blueprint v4.0.0 — Revised by Senior Full-Stack Engineer & Technical Product Manager.*
*F-01: Light/Dark Mode | F-02: Clinical Teal AI Accent | F-03: Separated History Charts | F-04: Context-Aware AI Chatbot*
*Dokumen ini adalah master reference aktif. v3.0.0 dinyatakan deprecated.*
