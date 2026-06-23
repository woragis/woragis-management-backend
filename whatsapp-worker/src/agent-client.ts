import type { Config } from './config.js'

export type InboundPayload = {
  channel: 'whatsapp'
  externalId: string
  userId?: string
  userName?: string
  text: string
}

function headers(cfg: Config): Record<string, string> {
  const key = cfg.channelWorkerKey
  return {
    'Content-Type': 'application/json',
    'X-Channel-Worker-Key': key,
    Authorization: `Bearer ${key}`,
  }
}

export async function forwardInbound(cfg: Config, payload: InboundPayload): Promise<void> {
  if (!cfg.agentWorkerUrl) return
  const res = await fetch(`${cfg.agentWorkerUrl}/v1/inbound`, {
    method: 'POST',
    headers: headers(cfg),
    body: JSON.stringify(payload),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`agent inbound failed: ${res.status} ${text}`)
  }
}
