import { AlertWithNarrative } from '@/types/telemetry'
import { motion, AnimatePresence } from 'framer-motion'
import { AlertTriangle, Loader2 } from 'lucide-react'

interface AlertFeedProps {
  alerts: AlertWithNarrative[]
}

export default function AlertFeed({ alerts }: AlertFeedProps) {
  if (alerts.length === 0) {
    return (
      <div className="text-center py-8 text-sm font-body text-healy-slate">
        <AlertTriangle className="w-5 h-5 mx-auto mb-2 opacity-40" />
        Sistem belum mendeteksi peringatan kritis.
      </div>
    )
  }

  return (
    <div className="space-y-4 max-h-100 overflow-y-auto pr-2">
      <AnimatePresence initial={false}>
        {alerts.map((alert) => (
          <motion.div
            key={alert.id}
            initial={{ opacity: 0, y: -8, height: 0 }}
            animate={{ opacity: 1, y: 0, height: 'auto' }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3, ease: 'easeOut' }}
            className="p-4 rounded-xl bg-healy-warning/10 border border-healy-warning/30"
          >
            <div className="flex items-center gap-3 mb-2">
              <div className="w-8 h-8 rounded-full bg-healy-warning/20 flex items-center justify-center shrink-0">
                <AlertTriangle className="w-4 h-4 text-healy-warning" />
              </div>
              <div className="flex-1">
                <h3 className="text-sm font-display font-semibold text-healy-warning">
                  {alert.alert_type} ({alert.value})
                </h3>
                <p className="text-xs font-mono text-healy-slate">
                  {alert.timestamp.toLocaleTimeString()}
                </p>
              </div>
            </div>

            {/* AI Narrative Section */}
            <div className="mt-3 pl-11">
              {alert.narrativeLoading ? (
                <div className="flex items-center gap-2 text-sm text-healy-slate italic">
                  <Loader2 className="w-4 h-4 animate-spin" />
                  AI sedang menganalisis kondisi...
                </div>
              ) : alert.narrative ? (
                <p className="text-sm text-healy-graphite italic border-l-2 border-healy-warning/50 pl-3 py-1">
                  &quot;{alert.narrative}&quot;
                </p>
              ) : null}
            </div>
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  )
}
