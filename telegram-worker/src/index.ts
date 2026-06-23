import { loadConfig } from './config.js'
import { startTelegramBot } from './bot.js'
import { createServer } from './server.js'
import pino from 'pino'

const log = pino({ name: 'telegram-worker' })

async function main() {
  const cfg = loadConfig()
  const bot = startTelegramBot(cfg)
  const app = createServer(cfg, bot)
  app.listen(cfg.port, () => log.info({ port: cfg.port }, 'telegram-worker listening'))
}

main().catch((err) => {
  log.error(err)
  process.exit(1)
})
