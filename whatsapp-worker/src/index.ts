import { loadConfig } from './config.js'
import { migrate } from './db.js'
import { loadPersistedQr, startSession } from './baileys/session.js'
import { createServer } from './server.js'

async function main(): Promise<void> {
  const cfg = loadConfig()

  await migrate(cfg)
  await loadPersistedQr(cfg)
  await startSession(cfg)

  const app = createServer(cfg)
  app.listen(cfg.port, () => {
    console.log(`whatsapp-worker listening on :${cfg.port}`)
  })
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
