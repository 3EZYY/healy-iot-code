// src/constants/design-tokens.ts
// Blueprint §3.1 — Color Tokens (Source of Truth)

export const COLORS = {
  bgPrimary:     '#F8FAFB',  // Glacial White
  bgSecondary:   '#EFF4F2',  // Soft Mist
  surface:       '#FFFFFF',  // Ice White

  brandPrimary:  '#4CAF82',  // Sage Green
  brandSecondary:'#2E8B62',  // Deep Sage
  brandAccent:   '#A8DFCC',  // Mint Glow

  textPrimary:   '#1A2633',  // Graphite
  textSecondary: '#5A7080',  // Slate

  statusNormal:  '#4CAF82',
  statusWarning: '#F5A623',  // Amber Soft
  statusCritical:'#E05252',  // Coral Red

  border:        '#D4E8DF',  // Pale Sage

  // ─── AI TIER — v4.0.0 (Menggantikan Violet v3.0.0) ───
  aiAccent:      '#0D9488',   // Clinical Teal
  aiSurface:     '#F0FDFC',   // Teal-50
} as const;

// Blueprint §3.2 — Typography
export const FONTS = {
  display: 'Exo 2',          // Hero, H1–H3 (weight: 600, 700, 800)
  body:    'DM Sans',        // Body, label, UI (weight: 400, 500)
  mono:    'JetBrains Mono', // Angka sensor, timestamp (weight: 400)
} as const;

// Status color mapping utility
export const STATUS_COLORS = {
  NORMAL:   { bg: '#E8F5EE', text: '#2E8B62', border: '#A8DFCC' },
  WARNING:  { bg: '#FEF3E0', text: '#C77D00', border: '#F5A623' },
  CRITICAL: { bg: '#FCE8E8', text: '#C73838', border: '#E05252' },
} as const;
