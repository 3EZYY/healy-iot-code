// Blueprint §5.2 — WebSocket Hook
// Manages real-time connection to backend WS endpoint with exponential backoff reconnect.

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
