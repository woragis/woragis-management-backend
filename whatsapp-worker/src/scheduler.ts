import cron from 'node-cron'
import type { Config } from './config.js'
import { cronExpr, weeklyCronExpr } from './cron.js'
import { fetchSettings } from './management-client.js'
import { runDispatch } from './sender.js'

const jobs: cron.ScheduledTask[] = []

function clearJobs(): void {
  for (const job of jobs) {
    job.stop()
  }
  jobs.length = 0
}

export async function startScheduler(cfg: Config): Promise<void> {
  clearJobs()
  const settings = await fetchSettings(cfg)

  jobs.push(
    cron.schedule(cronExpr(settings.problemPostTime), () => {
      runDispatch(cfg, 'problem').catch(console.error)
    }, { timezone: settings.timezone }),
  )

  if (settings.discussionEnabled) {
    jobs.push(
      cron.schedule(cronExpr(settings.discussionPostTime), () => {
        runDispatch(cfg, 'discussion').catch(console.error)
      }, { timezone: settings.timezone }),
    )
  }

  jobs.push(
    cron.schedule(cronExpr(settings.solutionPostTime), () => {
      runDispatch(cfg, 'solution').catch(console.error)
    }, { timezone: settings.timezone }),
  )

  jobs.push(
    cron.schedule(weeklyCronExpr(settings.weeklySummaryDay, settings.weeklySummaryTime), () => {
      runDispatch(cfg, 'weekly').catch(console.error)
    }, { timezone: settings.timezone }),
  )

  console.log(`[scheduler] started (${settings.timezone})`)
}

export function scheduleRefresh(cfg: Config, intervalMs = 15 * 60 * 1000): NodeJS.Timeout {
  return setInterval(() => {
    startScheduler(cfg).catch(console.error)
  }, intervalMs)
}
