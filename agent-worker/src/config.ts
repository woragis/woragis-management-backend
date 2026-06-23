export type Config = {
  port: number
  managementApiUrl: string
  agentApiKey: string
  openaiApiKey: string
  telegramBotToken: string
  allowedUserIds: Set<number>
}

export function loadConfig(): Config {
  const port = Number(process.env.PORT ?? '3002')
  const managementApiUrl = (process.env.MANAGEMENT_API_URL ?? 'http://127.0.0.1:8080').replace(/\/$/, '')
  const agentApiKey = (process.env.AGENT_API_KEY ?? '').trim()
  const openaiApiKey = (process.env.OPENAI_API_KEY ?? '').trim()
  const telegramBotToken = (process.env.TELEGRAM_BOT_TOKEN ?? '').trim()
  const allowed = (process.env.TELEGRAM_ALLOWED_USER_IDS ?? '')
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
    .map((s) => Number(s))
    .filter((n) => Number.isFinite(n))

  if (!agentApiKey) throw new Error('AGENT_API_KEY is required')
  if (!openaiApiKey) throw new Error('OPENAI_API_KEY is required')
  if (!telegramBotToken) throw new Error('TELEGRAM_BOT_TOKEN is required')

  return {
    port,
    managementApiUrl,
    agentApiKey,
    openaiApiKey,
    telegramBotToken,
    allowedUserIds: new Set(allowed),
  }
}
