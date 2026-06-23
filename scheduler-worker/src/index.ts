import http from 'node:http'
import { loadConfig } from './config.js'
import { executeJob, fetchDueJobs } from './management-client.js'
import pino from 'pino'

const log = pino({ name: 'scheduler-worker' })

async function tick(cfg: ReturnType<typeof loadConfig>): Promise<void> {
  if (!cfg.workerApiKey) {
    log.warn('WORKER_API_KEY not set; scheduler tick skipped')
    return
  }
  const jobs = await fetchDueJobs(cfg)
  if (jobs.length === 0) return
  log.info({ count: jobs.length }, 'due jobs')
  for (const job of jobs) {
    try {
      const result = await executeJob(cfg, job.id)
      log.info({ jobId: job.id, name: job.name, result }, 'job executed')
    } catch (err) {
      log.error({ err, jobId: job.id, name: job.name }, 'job execution failed')
    }
  }
}

async function main(): Promise<void> {
  const cfg = loadConfig()

  const server = http.createServer((_req, res) => {
    res.writeHead(200, { 'Content-Type': 'application/json' })
    res.end(JSON.stringify({ ok: true, service: 'scheduler-worker' }))
  })
  server.listen(cfg.port, () => log.info({ port: cfg.port }, 'health server listening'))

  if (!cfg.workerApiKey) {
    log.warn('WORKER_API_KEY not set; scheduler will not poll')
  } else {
    log.info({ intervalMs: cfg.pollIntervalMs }, 'scheduler polling started')
    await tick(cfg).catch((err) => log.error({ err }, 'initial tick failed'))
    setInterval(() => {
      tick(cfg).catch((err) => log.error({ err }, 'tick failed'))
    }, cfg.pollIntervalMs)
  }
}

main().catch((err) => {
  log.error(err)
  process.exit(1)
})
