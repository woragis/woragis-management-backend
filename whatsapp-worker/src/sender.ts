import type { Config } from './config.js'
import { getEnabledGroupJid, logMessage } from './db.js'
import { fetchDispatch, patchWhatsappStatus } from './management-client.js'
import { isConnected, sendGroupMessage } from './baileys/session.js'

function statusPatchForType(type: string): {
  problemSent?: boolean
  discussionSent?: boolean
  solutionSent?: boolean
} {
  switch (type) {
    case 'problem':
      return { problemSent: true }
    case 'discussion':
      return { discussionSent: true }
    case 'solution':
      return { solutionSent: true }
    default:
      return {}
  }
}

export async function runDispatch(cfg: Config, type: string, date?: string): Promise<void> {
  const dispatch = await fetchDispatch(cfg, type, date)
  const groupJid = await getEnabledGroupJid(cfg)

  if (dispatch.skip) {
    console.log(`[dispatch:${type}] skip: ${dispatch.skipReason || 'unknown'}`)
    return
  }
  if (!groupJid) {
    console.error(`[dispatch:${type}] no group configured`)
    return
  }
  if (!isConnected()) {
    console.error(`[dispatch:${type}] whatsapp not connected`)
    return
  }

  try {
    await sendGroupMessage(groupJid, dispatch.message)
    await logMessage(cfg, {
      groupJid,
      templateSlug: dispatch.templateSlug,
      dispatchType: type,
      body: dispatch.message,
      videoId: dispatch.videoId,
      status: 'sent',
    })
    if (dispatch.videoId) {
      await patchWhatsappStatus(cfg, dispatch.videoId, statusPatchForType(type))
    }
    console.log(`[dispatch:${type}] sent to ${groupJid}`)
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err)
    await logMessage(cfg, {
      groupJid,
      templateSlug: dispatch.templateSlug,
      dispatchType: type,
      body: dispatch.message,
      videoId: dispatch.videoId,
      status: 'failed',
      error: message,
    })
    console.error(`[dispatch:${type}] failed:`, message)
    throw err
  }
}

export async function sendRawMessage(
  cfg: Config,
  body: string,
  meta?: { templateSlug?: string; dispatchType?: string; videoId?: string },
): Promise<void> {
  const groupJid = await getEnabledGroupJid(cfg)
  if (!groupJid) throw new Error('no group configured')
  if (!isConnected()) throw new Error('whatsapp not connected')
  await sendGroupMessage(groupJid, body)
  await logMessage(cfg, {
    groupJid,
    templateSlug: meta?.templateSlug || '',
    dispatchType: meta?.dispatchType || 'manual',
    body,
    videoId: meta?.videoId,
    status: 'sent',
  })
}
