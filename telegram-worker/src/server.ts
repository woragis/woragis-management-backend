import express from 'express'
import type { Telegraf } from 'telegraf'
import type { Config } from './config.js'

function authorize(cfg: Config, req: express.Request): boolean {
  if (!cfg.channelWorkerKey) return false
  const header = req.header('X-Channel-Worker-Key') || req.header('Authorization')?.replace(/^Bearer\s+/i, '')
  return header === cfg.channelWorkerKey
}

export function createServer(cfg: Config, bot: Telegraf): express.Express {
  const app = express()
  app.use(express.json())

  app.get('/health', (_req, res) => {
    res.json({ ok: true })
  })

  app.post('/v1/send', async (req, res) => {
    if (!authorize(cfg, req)) {
      res.status(401).json({ error: 'unauthorized' })
      return
    }
    const { externalId, text, message } = req.body as { externalId?: string; text?: string; message?: string }
    const body = (text ?? message)?.trim()
    const chatId = externalId?.trim()
    if (!chatId || !body) {
      res.status(400).json({ error: 'externalId and text are required' })
      return
    }
    try {
      await bot.telegram.sendMessage(chatId, body)
      res.json({ ok: true })
    } catch (err) {
      const msg = err instanceof Error ? err.message : String(err)
      res.status(500).json({ error: msg })
    }
  })

  return app
}
