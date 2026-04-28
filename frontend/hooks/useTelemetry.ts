// Custom hook that abstracts the data source:
// - In mock mode: uses setInterval with generateMockPayload()
// - In production: uses the real useWebSocket hook

import { useState, useEffect, useRef } from 'react'
import { TelemetryPayload, ConnectionState, WebSocketMessage } from '@/types/telemetry'
import { useWebSocket } from '@/hooks/useWebSocket'
import { generateMockPayload } from '@/lib/mock-telemetry'

const MOCK_INTERVAL_MS = 2000 // Blueprint spec: 2-second refresh cycle

interface TelemetrySource {
  data: TelemetryPayload | null
  conn: ConnectionState
  deviceOnline: boolean
}

/**
 * Unified telemetry data provider.
 * Reads NEXT_PUBLIC_USE_MOCK_DATA to decide data source.
 */
export function useTelemetry(
  onMessage?: (data: WebSocketMessage) => void,
  onStatusChange?: (status: ConnectionState['status']) => void
): TelemetrySource {
  const useMock = process.env.NEXT_PUBLIC_USE_MOCK_DATA === 'true'
  const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws'

  if (useMock) {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    return useMockTelemetry(onMessage, onStatusChange)
  }

  // eslint-disable-next-line react-hooks/rules-of-hooks
  const { data, conn, deviceOnline } = useWebSocket(`${wsUrl}/viewer`, onMessage, onStatusChange)
  return { data, conn, deviceOnline }
}

/** Internal mock provider using setInterval */
function useMockTelemetry(
  onMessage?: (data: WebSocketMessage) => void,
  onStatusChange?: (status: ConnectionState['status']) => void
): TelemetrySource {
  const [data, setData] = useState<TelemetryPayload | null>(null)
  const [conn, setConn] = useState<ConnectionState>({
    status: 'DISCONNECTED',
    lastUpdate: null,
    retryCount: 0,
  })
  const intervalRef = useRef<NodeJS.Timeout | null>(null)
  const onMessageRef = useRef(onMessage)
  const onStatusRef = useRef(onStatusChange)

  useEffect(() => {
    onMessageRef.current = onMessage
    onStatusRef.current = onStatusChange
  }, [onMessage, onStatusChange])

  useEffect(() => {
    // Initialize first connection state and payload
    onStatusRef.current?.('CONNECTED')
    setConn({ status: 'CONNECTED', lastUpdate: new Date(), retryCount: 0 })
    
    const initialPayload = generateMockPayload()
    onMessageRef.current?.(initialPayload)
    setData(initialPayload)

    // Generate every MOCK_INTERVAL_MS
    intervalRef.current = setInterval(() => {
      const payload = generateMockPayload()
      onMessageRef.current?.(payload) // Trigger callback
      setData(payload)
      setConn(prev => ({ ...prev, lastUpdate: new Date() }))
    }, MOCK_INTERVAL_MS)

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current)
    }
  }, [])

  return { data, conn, deviceOnline: true } // Mock device is always online
}
