'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useState } from 'react'
import {
  LayoutDashboard,
  History,
  Settings,
  LogOut,
  Activity,
  Key,
  Check,
  Eye,
  EyeOff,
  MessageCircle
} from 'lucide-react'
import { getStoredGroqKey, setStoredGroqKey, validateGroqKey } from '@/lib/groq-client'

const NAV_ITEMS = [
  { label: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { label: 'History',   href: '/history',   icon: History },
  { label: 'Settings',  href: '/settings',  icon: Settings },
]

interface NavSidebarProps {
  onOpenChat?: () => void
}

export default function NavSidebar({ onOpenChat }: NavSidebarProps) {
  const pathname = usePathname()
  
  const [apiKey, setApiKey] = useState(() => getStoredGroqKey() || '')
  const [showKey, setShowKey] = useState(false)
  const [keySaved, setKeySaved] = useState(false)
  const [keyError, setKeyError] = useState(false)

  // Remove the useEffect that was setting state synchronously
  // useEffect(() => { ... }, []) is no longer needed

  const handleSaveKey = () => {
    if (validateGroqKey(apiKey)) {
      setStoredGroqKey(apiKey)
      setKeySaved(true)
      setKeyError(false)
      setTimeout(() => setKeySaved(false), 2000)
    } else {
      setKeyError(true)
    }
  }

  return (
    <aside
      aria-label="Main navigation sidebar"
      className="
        fixed left-0 top-0 bottom-0
        w-64 bg-healy-surface border-r border-healy-border
        flex flex-col
        z-40
      "
    >
      {/* Logo */}
      <div className="p-6 border-b border-healy-border">
        <Link href="/dashboard" className="flex items-center gap-3 group" aria-label="HEALY Health Observer - Go to Dashboard">
          <div className="
            w-10 h-10 rounded-xl
            bg-healy-sage
            flex items-center justify-center
            shadow-glow
            transition-transform duration-200
            group-hover:scale-105
          ">
            <Activity className="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 className="text-lg font-display font-bold text-healy-graphite tracking-tight">
              HEALY
            </h1>
            <p className="text-[10px] font-body text-healy-slate -mt-0.5">
              Health Observer
            </p>
          </div>
        </Link>
      </div>

      {/* Navigation */}
      <nav className="flex-1 p-4 space-y-1" aria-label="Main navigation">
        {NAV_ITEMS.map((item) => {
          const isActive = pathname?.startsWith(item.href)
          const Icon = item.icon

          return (
            <Link
              key={item.href}
              href={item.href}
              aria-label={`Navigate to ${item.label}`}
              aria-current={isActive ? 'page' : undefined}
              className={`
                flex items-center gap-3 px-4 py-3 rounded-xl
                font-body text-sm font-medium
                transition-all duration-200
                ${
                  isActive
                    ? 'bg-healy-sage/10 text-healy-sage-dark shadow-sm'
                    : 'text-healy-slate hover:bg-healy-bg-alt hover:text-healy-graphite'
                }
              `}
            >
              <Icon className={`w-5 h-5 ${isActive ? 'text-healy-sage' : ''}`} aria-hidden="true" />
              {item.label}
              {isActive && (
               <span className="ml-auto w-1.5 h-1.5 rounded-full bg-healy-sage" />
              )}
            </Link>
          )
        })}
      </nav>

      {/* Chatbot Trigger */}
      <div className="px-4 py-2">
        <button
          onClick={() => {
            if (onOpenChat) onOpenChat()
            else window.dispatchEvent(new CustomEvent('open-ai-chat'))
          }}
          className="flex items-center gap-3 w-full px-4 py-3 rounded-xl
                     text-healy-graphite hover:bg-healy-bg-alt
                     hover:text-healy-sage transition-all duration-200"
        >
          <MessageCircle className="w-5 h-5 text-healy-ai-accent" aria-hidden="true" />
          <span className="font-body text-sm font-medium">AI Health Chat</span>
          <span className="ml-auto text-xs bg-healy-ai-accent/15 text-healy-ai-accent
                           px-2 py-0.5 rounded-full font-mono font-bold">
            AI
          </span>
        </button>
      </div>

      {/* Groq Configuration */}
      <div className="p-4 border-t border-healy-border">
        <div className="mb-4">
          <div className="flex items-center gap-2 mb-2">
            <Key className="w-4 h-4 text-healy-slate" />
            <span className="text-xs font-semibold text-healy-graphite font-display">Groq API Key</span>
          </div>
          <div className="relative flex items-center">
            <input
              type={showKey ? 'text' : 'password'}
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              onBlur={handleSaveKey}
              placeholder="gsk_..."
              className={`w-full text-xs font-mono px-3 py-2 pr-16 border rounded-lg focus:outline-none focus:ring-1 transition-colors ${
                keyError ? 'border-healy-critical focus:ring-healy-critical' : 'border-healy-border focus:ring-healy-sage'
              }`}
            />
            <div className="absolute right-2 flex items-center gap-1">
              <button 
                onClick={() => setShowKey(!showKey)}
                className="p-1 text-healy-slate hover:text-healy-graphite transition-colors"
              >
                {showKey ? <EyeOff size={14} /> : <Eye size={14} />}
              </button>
              {keySaved && <Check size={14} className="text-healy-sage" />}
            </div>
          </div>
          {keyError && <p className="text-[10px] text-healy-critical mt-1">Invalid Groq Key format</p>}
        </div>

        <button
          aria-label="Logout from dashboard"
          className="
            flex items-center gap-3 w-full px-4 py-3 rounded-xl
            font-body text-sm font-medium text-healy-slate
            hover:bg-healy-critical/5 hover:text-healy-critical
            transition-all duration-200
          "
        >
          <LogOut className="w-5 h-5" aria-hidden="true" />
          Logout
        </button>
      </div>
    </aside>
  )
}
