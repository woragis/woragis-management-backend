import makeWASocket, {
  DisconnectReason,
  fetchLatestBaileysVersion,
  useMultiFileAuthState,
  type WASocket,
} from '@whiskeysockets/baileys'
import { Boom } from '@hapi/boom'
import pino from 'pino'
import QRCode from 'qrcode'
import type { Config } from '../config.js'
import { getState, setState } from '../db.js'
import { forwardInbound } from '../agent-client.js'

const logger = pino({ level: process.env.LOG_LEVEL || 'info' })

let sock: WASocket | null = null
let connected = false
let lastQr: string | null = null

export function isConnected(): boolean {
  return connected
}

export function getLastQr(): string | null {
  return lastQr
}

export async function startSession(cfg: Config): Promise<void> {
  const { state, saveCreds } = await useMultiFileAuthState(cfg.sessionDir)
  const { version } = await fetchLatestBaileysVersion()

  sock = makeWASocket({
    version,
    auth: state,
    logger,
    printQRInTerminal: false,
  })

  sock.ev.on('creds.update', saveCreds)

  sock.ev.on('messages.upsert', async ({ messages, type }) => {
    if (type !== 'notify' || !cfg.agentWorkerUrl) return
    for (const msg of messages) {
      if (msg.key.fromMe) continue
      const text =
        msg.message?.conversation ||
        msg.message?.extendedTextMessage?.text ||
        msg.message?.imageMessage?.caption
      if (!text?.trim()) continue
      const externalId = msg.key.remoteJid
      if (!externalId) continue
      const userId = msg.key.participant || msg.key.remoteJid || undefined
      try {
        await forwardInbound(cfg, {
          channel: 'whatsapp',
          externalId,
          userId: userId ?? undefined,
          text: text.trim(),
        })
      } catch (err) {
        logger.error({ err, externalId }, 'forward inbound failed')
      }
    }
  })

  sock.ev.on('connection.update', async (update) => {
    const { connection, lastDisconnect, qr } = update
    if (qr) {
      lastQr = await QRCode.toDataURL(qr)
      await setState(cfg, 'last_qr', lastQr)
    }
    if (connection === 'open') {
      connected = true
      lastQr = null
      await setState(cfg, 'connected_at', new Date().toISOString())
    }
    if (connection === 'close') {
      connected = false
      const code = (lastDisconnect?.error as Boom | undefined)?.output?.statusCode
      const shouldReconnect = code !== DisconnectReason.loggedOut
      if (shouldReconnect) {
        setTimeout(() => startSession(cfg).catch(console.error), 3000)
      }
    }
  })
}

export async function sendGroupMessage(groupJid: string, text: string): Promise<void> {
  return sendToJid(groupJid, text)
}

export async function sendToJid(jid: string, text: string): Promise<void> {
  if (!sock || !connected) throw new Error('whatsapp not connected')
  const target = jid.includes('@') ? jid : `${jid}@g.us`
  await sock.sendMessage(target, { text })
}

export async function loadPersistedQr(cfg: Config): Promise<void> {
  const qr = await getState(cfg, 'last_qr')
  if (qr) lastQr = qr
}
