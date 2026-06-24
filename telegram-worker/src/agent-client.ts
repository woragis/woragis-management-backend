import type { Config } from './config.js'

export type InboundPayload = {
  channel: 'telegram'
  externalId: string
  userId?: string
  userName?: string
  text: string
  destinationId?: string
}

function headers(cfg: Config): Record<string, string> {
  return {
    'Content-Type': 'application/json',
    'X-Channel-Worker-Key': cfg.channelWorkerKey,
    Authorization: `Bearer ${cfg.channelWorkerKey}`,
  }
}

export async function forwardInbound(cfg: Config, payload: InboundPayload): Promise<void> {
  if (!cfg.agentWorkerUrl) throw new Error('AGENT_WORKER_URL not configured')
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

export async function transcribeViaAgent(cfg: Config, audio: Buffer, mimeType = 'audio/ogg'): Promise<string> {
  if (!cfg.agentWorkerUrl) throw new Error('AGENT_WORKER_URL not configured')
  const res = await fetch(`${cfg.agentWorkerUrl}/v1/transcribe`, {
    method: 'POST',
    headers: headers(cfg),
    body: JSON.stringify({ audioBase64: audio.toString('base64'), mimeType }),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`agent transcribe failed: ${res.status} ${text}`)
  }
  const body = (await res.json()) as { text?: string }
  const text = body.text?.trim()
  if (!text) throw new Error('empty transcription')
  return text
}
