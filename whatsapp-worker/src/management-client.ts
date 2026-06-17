import type { Config } from './config.js'

export type DispatchResult = {
  videoId?: string
  templateSlug: string
  message: string
  skip?: boolean
  skipReason?: string
}

export type ChannelSettings = {
  timezone: string
  problemPostTime: string
  discussionPostTime: string
  solutionPostTime: string
  weeklySummaryDay: string
  weeklySummaryTime: string
  discussionEnabled: boolean
}

function headers(cfg: Config): Record<string, string> {
  return {
    'Content-Type': 'application/json',
    'X-Worker-Key': cfg.workerApiKey,
    Authorization: `Bearer ${cfg.workerApiKey}`,
  }
}

export async function fetchSettings(cfg: Config): Promise<ChannelSettings> {
  const res = await fetch(`${cfg.managementApiUrl}/v1/internal/content/leetcode/settings`, {
    headers: headers(cfg),
  })
  if (!res.ok) throw new Error(`settings http ${res.status}`)
  return res.json() as Promise<ChannelSettings>
}

export async function fetchDispatch(cfg: Config, type: string, date?: string): Promise<DispatchResult> {
  const params = new URLSearchParams({ type })
  if (date) params.set('date', date)
  const res = await fetch(`${cfg.managementApiUrl}/v1/internal/content/leetcode/dispatch?${params}`, {
    headers: headers(cfg),
  })
  if (!res.ok) throw new Error(`dispatch http ${res.status}`)
  return res.json() as Promise<DispatchResult>
}

export async function patchWhatsappStatus(
  cfg: Config,
  videoId: string,
  patch: { problemSent?: boolean; discussionSent?: boolean; solutionSent?: boolean },
): Promise<void> {
  const res = await fetch(
    `${cfg.managementApiUrl}/v1/internal/content/leetcode/videos/${videoId}/whatsapp-status`,
    {
      method: 'PATCH',
      headers: headers(cfg),
      body: JSON.stringify({
        problemSent: !!patch.problemSent,
        discussionSent: !!patch.discussionSent,
        solutionSent: !!patch.solutionSent,
      }),
    },
  )
  if (!res.ok) throw new Error(`patch status http ${res.status}`)
}
