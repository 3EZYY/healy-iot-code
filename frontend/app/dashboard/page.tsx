'use client'

import { useState, useRef, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import SensorCard from '@/components/features/SensorCard'
import StatusChip from '@/components/features/StatusChip'
import ConnectionStatus from '@/components/features/ConnectionStatus'
import { DeviceLedIndicator } from '@/components/features/DeviceLedIndicator'
import { AIInsightCard } from '@/components/features/AIInsightCard'
import AlertFeed from '@/components/features/AlertFeed'
import { useTelemetry } from '@/hooks/useTelemetry'
import { useAutoNarrative } from '@/hooks/useAutoNarrative'
import { Activity, Clock, Bell, AlertTriangle } from 'lucide-react'
import { 
  SensorStatus, 
  isSystemMessage, 
  WebSocketMessage 
} from '@/types/telemetry'
import type { GroqInsightRequest } from '@/lib/groq-client'

// ─── Activity Log Entry ───
interface ActivityEntry {
  time: string
  event: string
  status: SensorStatus
}

const MAX_ACTIVITY_LOG = 15

const stagger = {
  hidden: {},
  visible: { transition: { staggerChildren: 0.06 } },
} as const

const fadeUp = {
  hidden: { opacity: 0, y: 16 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease: 'easeOut' as const } },
}

export default function DashboardPage() {
  const [activityLog, setActivityLog] = useState<ActivityEntry[]>([])
  const prevStatusRef = useRef<SensorStatus | null>(null)

  // ─── Activity Log Handlers ───
  const handleNewMessage = useCallback((msg: WebSocketMessage) => {
    if (isSystemMessage(msg)) return

    const now = new Date().toLocaleTimeString()
    const entries: ActivityEntry[] = []

    // Log status transitions
    if (prevStatusRef.current && prevStatusRef.current !== msg.status.overall) {
      entries.push({
        time: now,
        event: `Status changed: ${prevStatusRef.current} → ${msg.status.overall}`,
        status: msg.status.overall,
      })
    }
    prevStatusRef.current = msg.status.overall

    setActivityLog(prev => [...entries, ...prev].slice(0, MAX_ACTIVITY_LOG))
  }, [])

  const handleStatusChange = useCallback((status: string) => {
    if (status === 'CONNECTED') {
      setActivityLog(prev => [{
        time: new Date().toLocaleTimeString(),
        event: `Connection established (${process.env.NEXT_PUBLIC_USE_MOCK_DATA === 'true' ? 'Mock Mode' : 'WebSocket'})`,
        status: 'NORMAL' as SensorStatus,
      }, ...prev].slice(0, MAX_ACTIVITY_LOG))
    }
  }, [])

  const { data, conn, deviceOnline } = useTelemetry(handleNewMessage, handleStatusChange)
  const { alerts } = useAutoNarrative(data)

  // ─── Fallback while waiting for first payload ───
  if (!data) {
    return (
      <div className="max-w-5xl mx-auto flex items-center justify-center min-h-[50vh]">
        <div className="text-center">
          <div className="w-12 h-12 rounded-2xl bg-healy-sage/10 flex items-center justify-center mx-auto mb-4">
            <Activity className="w-6 h-6 text-healy-sage animate-pulse" />
          </div>
          <p className="text-sm font-body text-healy-slate mb-2">Waiting for telemetry data...</p>
          <div className="flex justify-center mb-4">
            <DeviceLedIndicator online={deviceOnline} />
          </div>
          <ConnectionStatus conn={conn} />
        </div>
      </div>
    )
  }

  const groqInsightData: GroqInsightRequest = {
    temperature: data.sensor.temperature,
    bpm: data.sensor.bpm,
    spo2: data.sensor.spo2,
    status: data.status.overall,
  }

  return (
    <motion.div
      variants={stagger}
      initial="hidden"
      animate="visible"
      className="max-w-5xl mx-auto"
    >
      {/* ─── Header ─── */}
      <motion.div variants={fadeUp} className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-2xl font-display font-bold text-healy-graphite">
            Live Monitoring
          </h1>
          <div className="flex items-center gap-3 mt-1">
            <p className="text-sm font-body text-healy-slate">
              Real-time telemetry from {data.device_id}
              {process.env.NEXT_PUBLIC_USE_MOCK_DATA === 'true' && (
                <span className="ml-2 px-2 py-0.5 rounded-md bg-healy-warning/10 text-healy-warning text-xs font-mono">
                  MOCK
                </span>
              )}
            </p>
            <div className="h-4 w-px bg-healy-border" />
            <DeviceLedIndicator online={deviceOnline} />
          </div>
        </div>
        <div className="flex items-center gap-3">
          <StatusChip status={data.status.overall} label={`Overall: ${data.status.overall}`} />
          <ConnectionStatus conn={conn} />
        </div>
      </motion.div>

      {/* ─── AI Insight ─── */}
      <motion.div variants={fadeUp} className="mb-8">
        <AIInsightCard currentData={groqInsightData} />
      </motion.div>

      {/* ─── Sensor Cards Grid ─── */}
      <motion.div variants={fadeUp} className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <SensorCard
          title="Body Temperature"
          value={data.sensor.temperature}
          unit="°C"
          status={data.status.temperature}
          icon="temperature"
        />
        <SensorCard
          title="Heart Rate"
          value={data.sensor.bpm}
          unit="BPM"
          status={data.status.overall === 'CRITICAL' ? 'CRITICAL' : data.status.overall === 'WARNING' ? 'WARNING' : 'NORMAL'}
          icon="bpm"
        />
        <SensorCard
          title="Blood Oxygen"
          value={data.sensor.spo2}
          unit="%"
          status={data.status.spo2}
          icon="spo2"
        />
      </motion.div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* ─── Device Info Panel ─── */}
        <motion.div variants={fadeUp} className="glass-card p-6 lg:col-span-1 h-fit">
          <div className="flex items-center gap-3 mb-4">
            <Activity className="w-5 h-5 text-healy-sage" />
            <h2 className="text-lg font-display font-semibold text-healy-graphite">
              Device Information
            </h2>
          </div>
          <div className="grid grid-cols-2 gap-6">
            {[
              { label: 'Device ID',   value: data.device_id },
              { label: 'Firmware',    value: 'v1.0.0' },
              { label: 'Data Source', value: process.env.NEXT_PUBLIC_USE_MOCK_DATA === 'true' ? 'Mock Generator' : 'WebSocket' },
              { label: 'Last Update', value: conn.lastUpdate ? conn.lastUpdate.toLocaleTimeString() : '—' },
            ].map((info) => (
              <div key={info.label}>
                <span className="text-xs font-body text-healy-slate block mb-1">{info.label}</span>
                <span className="text-sm font-mono text-healy-graphite">{info.value}</span>
              </div>
            ))}
          </div>
        </motion.div>

        {/* ─── Activity Feed (Live) ─── */}
        <motion.div variants={fadeUp} className="glass-card p-6 lg:col-span-2">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center gap-3">
              <Clock className="w-5 h-5 text-healy-sage" />
              <h2 className="text-lg font-display font-semibold text-healy-graphite">
                Live Activity Feed
              </h2>
            </div>
            <span className="text-xs font-mono text-healy-slate">
              {activityLog.length} events
            </span>
          </div>
          <div className="space-y-2 max-h-80 overflow-y-auto pr-2">
            <AnimatePresence initial={false}>
              {activityLog.map((activity, i) => (
                <motion.div
                  key={`${activity.time}-${i}`}
                  initial={{ opacity: 0, y: -8, height: 0 }}
                  animate={{ opacity: 1, y: 0, height: 'auto' }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.2, ease: 'easeOut' as const }}
                  className="flex items-center gap-4 p-3 rounded-xl bg-healy-bg-alt/50 border border-healy-border/30"
                >
                  <span className="text-xs font-mono text-healy-slate w-20 shrink-0">{activity.time}</span>
                  <span className="text-sm font-body text-healy-graphite flex-1">{activity.event}</span>
                  <StatusChip status={activity.status} size="sm" />
                </motion.div>
              ))}
            </AnimatePresence>
            {activityLog.length === 0 && (
              <div className="text-center py-8 text-sm font-body text-healy-slate">
                <Bell className="w-5 h-5 mx-auto mb-2 opacity-40" />
                Waiting for activity...
              </div>
            )}
          </div>
        </motion.div>
      </div>

      {/* ─── Proactive Alerts (Auto-Narrative) ─── */}
      <motion.div variants={fadeUp} className="glass-card p-6 mt-6">
        <div className="flex items-center gap-3 mb-4">
          <AlertTriangle className="w-5 h-5 text-healy-warning" />
          <h2 className="text-lg font-display font-semibold text-healy-graphite">
            Proactive AI Alerts
          </h2>
        </div>
        <AlertFeed alerts={alerts} />
      </motion.div>
    </motion.div>
  )
}
