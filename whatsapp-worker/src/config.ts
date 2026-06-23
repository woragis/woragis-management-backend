export type Config = {
  port: number
  databaseUrl: string
  managementApiUrl: string
  workerApiKey: string
  agentWorkerUrl: string
  channelWorkerKey: string
  sessionDir: string
  timezone: string
}

export function loadConfig(): Config {
  return {
    port: Number(process.env.PORT || 3001),
    databaseUrl: process.env.DATABASE_URL || '',
    managementApiUrl: (process.env.MANAGEMENT_API_URL || 'http://127.0.0.1:8080').replace(/\/$/, ''),
    workerApiKey: process.env.WORKER_API_KEY || '',
    agentWorkerUrl: (process.env.AGENT_WORKER_URL || '').replace(/\/$/, ''),
    channelWorkerKey: process.env.CHANNEL_WORKER_KEY || process.env.AGENT_API_KEY || '',
    sessionDir: process.env.SESSION_DIR || './data/session',
    timezone: process.env.TZ || 'America/Sao_Paulo',
  }
}
