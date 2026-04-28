'use client'

import { type SensorStatus } from '@/types/telemetry'

interface StatusChipProps {
  status: SensorStatus
  label?: string
  size?: 'sm' | 'md'
  pulse?: boolean
}

const STATUS_CONFIG = {
  NORMAL: {
    bg: 'bg-[#E8F5EE]',
    text: 'text-[#2E8B62]',
    dot: 'bg-healy-sage',
    label: 'Normal',
  },
  WARNING: {
    bg: 'bg-[#FEF3E0]',
    text: 'text-[#C77D00]',
    dot: 'bg-healy-warning',
    label: 'Warning',
  },
  CRITICAL: {
    bg: 'bg-[#FCE8E8]',
    text: 'text-[#C73838]',
    dot: 'bg-healy-critical',
    label: 'Critical',
  },
} as const

export default function StatusChip({ status, label, size = 'md', pulse = false }: StatusChipProps) {
  const config = STATUS_CONFIG[status]
  const shouldPulse = pulse || status === 'CRITICAL'

  return (
    <span
      className={`inline-flex items-center gap-1.5 font-body font-medium ${config.bg} ${config.text} ${size === 'sm' ? 'px-2 py-0.5 text-xs rounded-md' : 'px-3 py-1 text-sm rounded-full'}`}
    >
      <span
        className={`inline-block rounded-full ${size === 'sm' ? 'w-1.5 h-1.5' : 'w-2 h-2'} ${config.dot} ${shouldPulse ? 'animate-pulse-critical' : ''}`}
      />
      {label || config.label}
    </span>
  )
}
