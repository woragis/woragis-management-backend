import express from 'express'
import type { Config } from './config.js'
import { getLastQr, isConnected } from './baileys/session.js'
import { runDispatch, sendRawMessage } from './sender.js'
import { patchWhatsappStatus } from './management-client.js'

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
    try {
      const { message, type, videoId, templateSlug } = req.body as {
        message?: string
        type?: string
        videoId?: string
        templateSlug?: string
      }

      if (message) {
        await sendRawMessage(cfg, message, { dispatchType: type, videoId, templateSlug })
        if (videoId && type) {
          const patch =
            type === 'problem'
              ? { problemSent: true }
              : type === 'discussion'
                ? { discussionSent: true }
                : type === 'solution'
                  ? { solutionSent: true }
                  : {}
          if (Object.keys(patch).length) {
            await patchWhatsappStatus(cfg, videoId, patch)
          }
        }
        res.json({ ok: true })
        return
      }

      if (!type) {
        res.status(400).json({ error: 'type or message required' })
        return
      }
      await runDispatch(cfg, type)
      res.json({ ok: true })
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err)
      res.status(500).json({ error: msg })
    }
  })

  return app
}
