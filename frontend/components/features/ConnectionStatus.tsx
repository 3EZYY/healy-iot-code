'use client'

import { Wifi, WifiOff, RefreshCw } from 'lucide-react'
import type { ConnectionState } from '@/types/telemetry'

interface ConnectionStatusProps {
  conn: ConnectionState
}

const STATUS_CONFIG = {
  CONNECTED: {
    icon: Wifi,
    label: 'Connected',
    bg: 'bg-healy-sage/10',
    text: 'text-healy-sage-dark',
    dot: 'bg-healy-sage',
    animate: false,
  },
  DISCONNECTED: {
    icon: WifiOff,
    label: 'Disconnected',
    bg: 'bg-healy-critical/10',
    text: 'text-healy-critical',
    dot: 'bg-healy-critical',
    animate: false,
  },
  RECONNECTING: {
    icon: RefreshCw,
    label: 'Reconnecting',
    bg: 'bg-healy-warning/10',
    text: 'text-healy-warning',
    dot: 'bg-healy-warning',
    animate: true,
  },
} as const

export default function ConnectionStatus({ conn }: ConnectionStatusProps) {
  const config = STATUS_CONFIG[conn.status]
  const Icon = config.icon

  return (
    <div className={`
      inline-flex items-center gap-2 px-3 py-1.5 rounded-full
      text-xs font-body font-medium
      ${config.bg} ${config.text}
      transition-colors duration-200
    `}>
      {/* Animated dot */}
      <span className={`
        inline-block w-1.5 h-1.5 rounded-full
        ${config.dot}
        ${conn.status === 'CONNECTED' ? 'animate-pulse' : ''}
      `} />

      {/* Icon */}
      <Icon className={`w-3.5 h-3.5 ${config.animate ? 'animate-spin' : ''}`} />

      {/* Label */}
      <span>{config.label}</span>

      {/* Retry count badge */}
      {conn.retryCount > 0 && conn.status === 'RECONNECTING' && (
        <span className="
          ml-1 px-1.5 py-0.5 rounded-md
          bg-healy-warning/20 text-healy-warning
          text-[10px] font-mono
        ">
          #{conn.retryCount}
        </span>
      )}

      {/* Last update timestamp */}
      {conn.lastUpdate && (
        <span className="ml-1 text-[10px] font-mono opacity-70">
          {conn.lastUpdate.toLocaleTimeString()}
        </span>
      )}
    </div>
  )
}
