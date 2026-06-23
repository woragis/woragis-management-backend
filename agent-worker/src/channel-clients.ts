import type { Config } from './config.js'

export type InboundPayload = {
  channel: 'telegram' | 'whatsapp'
  externalId: string
  userId?: string
  userName?: string
  text: string
  destinationId?: string
}

export type ComposeRequest = {
  templateBody: string
  composeMode?: 'static' | 'ai_assisted'
  data?: Record<string, unknown>
  destinationContext?: Record<string, unknown>
}

function channelHeaders(cfg: Config): Record<string, string> {
  return {
    'Content-Type': 'application/json',
    'X-Channel-Worker-Key': cfg.channelWorkerKey,
    Authorization: `Bearer ${cfg.channelWorkerKey}`,
  }
}

export async function sendToChannel(
  cfg: Config,
  channel: InboundPayload['channel'],
  externalId: string,
  text: string,
  destinationId?: string,
): Promise<void> {
  if (channel === 'telegram') {
    if (!cfg.telegramWorkerUrl) throw new Error('TELEGRAM_WORKER_URL not configured')
    const res = await fetch(`${cfg.telegramWorkerUrl}/v1/send`, {
      method: 'POST',
      headers: channelHeaders(cfg),
      body: JSON.stringify({ externalId, text }),
    })
    if (!res.ok) throw new Error(`telegram send failed: ${res.status} ${await res.text()}`)
    return
  }
  if (!cfg.whatsappWorkerUrl) throw new Error('WHATSAPP_WORKER_URL not configured')
  const res = await fetch(`${cfg.whatsappWorkerUrl}/v1/send`, {
    method: 'POST',
    headers: channelHeaders(cfg),
    body: JSON.stringify({ externalId, message: text, destinationId }),
  })
  if (!res.ok) throw new Error(`whatsapp send failed: ${res.status} ${await res.text()}`)
}
