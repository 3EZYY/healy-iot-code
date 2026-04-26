// Custom hook that abstracts the data source:
// - In mock mode: uses setInterval with generateMockPayload()
// - In production: uses the real useWebSocket hook

import { useState, useEffect, useRef } from 'react'
import { TelemetryPayload, ConnectionState } from '@/types/telemetry'
import { useWebSocket } from '@/hooks/useWebSocket'
import { generateMockPayload } from '@/lib/mock-telemetry'

const MOCK_INTERVAL_MS = 2000 // Blueprint spec: 2-second refresh cycle

interface TelemetrySource {
  data: TelemetryPayload | null
  conn: ConnectionState
}

/**
 * Unified telemetry data provider.
 * Reads NEXT_PUBLIC_USE_MOCK_DATA to decide data source.
 */
export function useTelemetry(): TelemetrySource {
  const useMock = process.env.NEXT_PUBLIC_USE_MOCK_DATA === 'true'
  const wsUrl = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws'

  if (useMock) {
    // eslint-disable-next-line react-hooks/rules-of-hooks
    return useMockTelemetry()
  }

  // eslint-disable-next-line react-hooks/rules-of-hooks
  return useWebSocket(`${wsUrl}/telemetry`)
}

/** Internal mock provider using setInterval */
function useMockTelemetry(): TelemetrySource {
  const [data, setData] = useState<TelemetryPayload | null>(null)
  const [conn, setConn] = useState<ConnectionState>({
    status: 'CONNECTED',
    lastUpdate: null,
    retryCount: 0,
  })
  const intervalRef = useRef<NodeJS.Timeout | null>(null)

  useEffect(() => {
    // Generate initial payload immediately
    const initial = generateMockPayload()
    setData(initial)
    setConn(prev => ({ ...prev, lastUpdate: new Date() }))

    // Then generate every MOCK_INTERVAL_MS
    intervalRef.current = setInterval(() => {
      const payload = generateMockPayload()
      setData(payload)
      setConn(prev => ({ ...prev, lastUpdate: new Date() }))
    }, MOCK_INTERVAL_MS)

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current)
    }
  }, [])

  return { data, conn }
}
