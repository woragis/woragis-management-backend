import type { Config } from './config.js'
import { logMessage } from './db.js'
import { isConnected, sendToJid } from './baileys/session.js'

export async function sendRawMessage(
  cfg: Config,
  body: string,
  target: { externalId: string; destinationId?: string },
  meta?: { templateSlug?: string; dispatchType?: string; videoId?: string },
): Promise<void> {
  if (!target.externalId) throw new Error('externalId is required')
  if (!isConnected()) throw new Error('whatsapp not connected')
  await sendToJid(target.externalId, body)
  await logMessage(cfg, {
    groupJid: target.externalId,
    templateSlug: meta?.templateSlug || '',
    dispatchType: meta?.dispatchType || 'manual',
    body,
    videoId: meta?.videoId,
    status: 'sent',
  })
}
