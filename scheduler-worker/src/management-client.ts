import type { Config } from './config.js'

export type ScheduledJob = {
  id: string
  name: string
  destinationId: string
  programAction: string
  cronExpr: string
  enabled: boolean
}

function headers(cfg: Config): Record<string, string> {
  const h: Record<string, string> = { 'Content-Type': 'application/json' }
  if (cfg.workerApiKey) {
    h['X-Worker-Key'] = cfg.workerApiKey
    h.Authorization = `Bearer ${cfg.workerApiKey}`
  }
  return h
}

export async function fetchDueJobs(cfg: Config): Promise<ScheduledJob[]> {
  const res = await fetch(`${cfg.managementApiUrl}/v1/internal/scheduler/due`, {
    headers: headers(cfg),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`scheduler due http ${res.status}: ${text}`)
  }
  return res.json() as Promise<ScheduledJob[]>
}

export async function executeJob(cfg: Config, jobId: string): Promise<unknown> {
  const res = await fetch(`${cfg.managementApiUrl}/v1/internal/scheduler/jobs/${jobId}/execute`, {
    method: 'POST',
    headers: headers(cfg),
  })
  const text = await res.text()
  if (!res.ok) {
    throw new Error(`scheduler execute http ${res.status}: ${text}`)
  }
  if (!text) return {}
  return JSON.parse(text) as unknown
}
