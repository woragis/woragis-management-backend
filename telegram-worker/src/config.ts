export type Config = {
  port: number
  agentWorkerUrl: string
  channelWorkerKey: string
  telegramBotToken: string
  allowedUserIds: Set<number>
}

export function loadConfig(): Config {
  const port = Number(process.env.PORT ?? '3003')
  const agentWorkerUrl = (process.env.AGENT_WORKER_URL ?? '').replace(/\/$/, '')
  const channelWorkerKey = (process.env.CHANNEL_WORKER_KEY ?? process.env.AGENT_API_KEY ?? '').trim()
  const telegramBotToken = (process.env.TELEGRAM_BOT_TOKEN ?? '').trim()
  const allowed = (process.env.TELEGRAM_ALLOWED_USER_IDS ?? '')
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean)
    .map((s) => Number(s))
    .filter((n) => Number.isFinite(n))

  if (!telegramBotToken) throw new Error('TELEGRAM_BOT_TOKEN is required')

  return {
    port,
    agentWorkerUrl,
    channelWorkerKey,
    telegramBotToken,
    allowedUserIds: new Set(allowed),
  }
}
