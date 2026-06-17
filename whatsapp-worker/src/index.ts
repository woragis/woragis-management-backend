import { loadConfig } from './config.js'
import { migrate, ensureDefaultGroup } from './db.js'
import { loadPersistedQr, startSession } from './baileys/session.js'
import { createServer } from './server.js'
import { scheduleRefresh, startScheduler } from './scheduler.js'

async function main(): Promise<void> {
  const cfg = loadConfig()
  if (!cfg.workerApiKey) {
    console.warn('warning: WORKER_API_KEY not set')
  }

  await migrate(cfg)
  if (cfg.defaultGroupJid) {
    await ensureDefaultGroup(cfg, cfg.defaultGroupJid)
  }

  await loadPersistedQr(cfg)
  await startSession(cfg)
  await startScheduler(cfg)
  scheduleRefresh(cfg)

  const app = createServer(cfg)
  app.listen(cfg.port, () => {
    console.log(`whatsapp-worker listening on :${cfg.port}`)
  })
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
