// Blueprint §5.3 — Mock Telemetry Generator
// Produces realistic payloads for development when backend is unavailable.
// Distribution: 80% normal, 15% warning, 5% critical.

import { TelemetryPayload, SensorStatus } from '@/types/telemetry'

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

  const tempStatus: SensorStatus = temp >= 36.5 && temp <= 37.5 ? 'NORMAL'
                   : temp <= 38.5 ? 'WARNING' : 'CRITICAL'
  const spo2Status: SensorStatus = spo2 >= 95 ? 'NORMAL' : spo2 >= 91 ? 'WARNING' : 'CRITICAL'
  const statuses = [tempStatus, spo2Status]
  const overall: SensorStatus = statuses.includes('CRITICAL') ? 'CRITICAL'
                : statuses.includes('WARNING')  ? 'WARNING' : 'NORMAL'

  return {
    device_id: 'healy-001',
    timestamp: new Date().toISOString(),
    sensor: { temperature: parseFloat(temp.toFixed(1)), bpm, spo2 },
    status: { temperature: tempStatus, spo2: spo2Status, overall },
  }
}
