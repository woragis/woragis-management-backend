import express from 'express'
import type { Config } from './config.js'
import { getLastQr, isConnected } from './baileys/session.js'
import { sendRawMessage } from './sender.js'

function authorize(cfg: Config, req: express.Request): boolean {
  const header =
    req.header('X-Channel-Worker-Key') ||
    req.header('X-Worker-Key') ||
    req.header('Authorization')?.replace(/^Bearer\s+/i, '')
  if (!header) return !cfg.channelWorkerKey && !cfg.workerApiKey
  return header === cfg.channelWorkerKey || header === cfg.workerApiKey
}

export function createServer(cfg: Config): express.Express {
  const app = express()
  app.use(express.json())

  app.get('/health', (_req, res) => {
    res.json({ ok: true })
  })

  app.get('/v1/status', (_req, res) => {
    res.json({ connected: isConnected() })
  })

  app.get('/v1/qr', (_req, res) => {
    const qr = getLastQr()
    if (!qr) {
      res.status(404).json({ error: 'no qr available' })
      return
    }
    res.json({ qr })
  })

  app.post('/v1/send', async (req, res) => {
    if (!authorize(cfg, req)) {
      res.status(401).json({ error: 'unauthorized' })
      return
    }
    try {
      const { message, externalId, destinationId, type, videoId, templateSlug } = req.body as {
        message?: string
        externalId?: string
        destinationId?: string
        type?: string
        videoId?: string
        templateSlug?: string
      }

      if (!message?.trim()) {
        res.status(400).json({ error: 'message is required' })
        return
      }
      if (!externalId?.trim()) {
        res.status(400).json({ error: 'externalId is required' })
        return
      }

      await sendRawMessage(
        cfg,
        message,
        { externalId: externalId.trim(), destinationId },
        { dispatchType: type, videoId, templateSlug },
      )
      res.json({ ok: true })
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err)
      res.status(500).json({ error: msg })
    }
  })

  return app
}
