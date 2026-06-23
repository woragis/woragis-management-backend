export type Config = {
  port: number
  managementApiUrl: string
  workerApiKey: string
  pollIntervalMs: number
}

export function loadConfig(): Config {
  const port = Number(process.env.PORT ?? '3004')
  const managementApiUrl = (process.env.MANAGEMENT_API_URL ?? 'http://127.0.0.1:8080').replace(/\/$/, '')
  const workerApiKey = (process.env.WORKER_API_KEY ?? '').trim()
  const pollIntervalMs = Number(process.env.POLL_INTERVAL_MS ?? '60000')

  return {
    port,
    managementApiUrl,
    workerApiKey,
    pollIntervalMs: Number.isFinite(pollIntervalMs) && pollIntervalMs > 0 ? pollIntervalMs : 60000,
  }
}
