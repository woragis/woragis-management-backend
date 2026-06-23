import type { AgentPersonality } from '../api-client.js'

export function pickGreeting(p: AgentPersonality, now = new Date()): string | null {
  if (!p.greetingEnabled) return null
  const hour = zonedHour(now, p.timezone)
  if (hour >= 5 && hour < 12) return p.greetingMorning
  if (hour >= 12 && hour < 18) return p.greetingAfternoon
  return p.greetingEvening
}

function zonedHour(date: Date, timeZone: string): number {
  const parts = new Intl.DateTimeFormat('en-US', {
    timeZone: timeZone || 'America/Sao_Paulo',
    hour: 'numeric',
    hour12: false,
  }).formatToParts(date)
  const hour = parts.find((p) => p.type === 'hour')?.value ?? '0'
  return Number(hour)
}
