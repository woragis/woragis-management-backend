import http from 'node:http'
import { loadConfig } from './config.js'
import { startTelegramBot } from './telegram/bot.js'
import pino from 'pino'

const log = pino({ name: 'agent-worker' })

async function main() {
  const cfg = loadConfig()

  const server = http.createServer((_req, res) => {
    res.writeHead(200, { 'Content-Type': 'application/json' })
    res.end(JSON.stringify({ ok: true, service: 'agent-worker' }))
  })
  server.listen(cfg.port, () => log.info({ port: cfg.port }, 'health server listening'))

  startTelegramBot(cfg)
}

main().catch((err) => {
  log.error(err)
  process.exit(1)
})
