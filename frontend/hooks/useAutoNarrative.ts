import { useState, useEffect, useRef } from 'react'
import { TelemetryPayload, AlertWithNarrative } from '@/types/telemetry'
import { getStoredGroqKey, callGroqInsight } from '@/lib/groq-client'

export function useAutoNarrative(data: TelemetryPayload | null) {
  const [alerts, setAlerts] = useState<AlertWithNarrative[]>([])
  const prevStatusRef = useRef<string | null>(null)

  useEffect(() => {
    if (!data) return

    const currentStatus = data.status.overall

    // Trigger on transition to CRITICAL
    if (currentStatus === 'CRITICAL' && prevStatusRef.current !== 'CRITICAL') {
      const alertId = `alert-${Date.now()}`
      
      // Determine what triggered it for the UI
      let alertType = 'General Critical'
      let alertValue = 0
      
      if (data.status.spo2 === 'CRITICAL') {
        alertType = 'SpO2 Critical'
        alertValue = data.sensor.spo2
      } else if (data.status.temperature === 'CRITICAL') {
        alertType = 'Temperature Critical'
        alertValue = data.sensor.temperature
      } else {
        alertType = 'Heart Rate Critical'
        alertValue = data.sensor.bpm
      }

      const newAlert: AlertWithNarrative = {
        id: alertId,
        timestamp: new Date(),
        alert_type: alertType,
        value: alertValue,
        status: 'CRITICAL',
        device_id: data.device_id,
        narrative: null,
        narrativeLoading: true,
      }

      setAlerts((prev) => [newAlert, ...prev].slice(0, 10)) // Keep last 10 alerts

      // Fetch from Groq without blocking the main thread
      const fetchNarrative = async () => {
        try {
          const apiKey = getStoredGroqKey()
          if (!apiKey) {
            throw new Error('No Groq API Key found in localStorage.')
          }

          const prompt = `PASIEN DALAM KONDISI KRITIS.
Suhu: ${data.sensor.temperature}°C
Detak Jantung: ${data.sensor.bpm} bpm
SpO2: ${data.sensor.spo2}%

Sebagai asisten AI medis, berikan satu paragraf (maksimal 3 kalimat) analisis darurat dan rekomendasi tindakan pertama untuk keluarga. Gunakan Bahasa Indonesia.`

          const result = await callGroqInsight(prompt, apiKey)
          
          setAlerts((prev) => 
            prev.map(a => a.id === alertId ? { ...a, narrative: result, narrativeLoading: false } : a)
          )
        } catch {
          // Fail silently by just turning off the loading state
          setAlerts((prev) => 
            prev.map(a => a.id === alertId ? { ...a, narrative: null, narrativeLoading: false } : a)
          )
        }
      }

      fetchNarrative()
    }

    prevStatusRef.current = currentStatus
  }, [data])

  return { alerts }
}
