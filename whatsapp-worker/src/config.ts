export type Config = {
  port: number
  databaseUrl: string
  managementApiUrl: string
  workerApiKey: string
  sessionDir: string
  defaultGroupJid: string
  timezone: string
}

export function loadConfig(): Config {
  return {
    port: Number(process.env.PORT || 3001),
    databaseUrl: process.env.DATABASE_URL || '',
    managementApiUrl: (process.env.MANAGEMENT_API_URL || 'http://127.0.0.1:8080').replace(/\/$/, ''),
    workerApiKey: process.env.WORKER_API_KEY || '',
    sessionDir: process.env.SESSION_DIR || './data/session',
    defaultGroupJid: process.env.WHATSAPP_GROUP_JID || '',
    timezone: process.env.TZ || 'America/Sao_Paulo',
  }
}
