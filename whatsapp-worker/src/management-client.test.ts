import assert from 'node:assert/strict'
import http from 'node:http'
import { describe, it } from 'node:test'
import type { Config } from './config.js'
import { fetchDispatch, fetchSettings } from './management-client.js'

function withServer(
  handler: http.RequestListener,
  run: (baseUrl: string) => Promise<void>,
): Promise<void> {
  return new Promise((resolve, reject) => {
    const srv = http.createServer(handler)
    srv.listen(0, async () => {
      const addr = srv.address()
      if (!addr || typeof addr === 'string') {
        srv.close()
        reject(new Error('invalid listen address'))
        return
      }
      const baseUrl = `http://127.0.0.1:${addr.port}`
      try {
        await run(baseUrl)
        resolve()
      } catch (err) {
        reject(err)
      } finally {
        srv.close()
      }
    })
  })
}

describe('management client', () => {
  it('fetches settings and dispatch with worker key', async () => {
    await withServer((req, res) => {
      if (req.url === '/v1/internal/content/leetcode/settings') {
        res.writeHead(200, { 'Content-Type': 'application/json' })
        res.end(
          JSON.stringify({
            timezone: 'America/Sao_Paulo',
            problemPostTime: '09:00',
            discussionPostTime: '17:00',
            solutionPostTime: '22:00',
            weeklySummaryDay: 'sunday',
            weeklySummaryTime: '10:00',
            discussionEnabled: true,
          }),
        )
        return
      }
      if (req.url?.startsWith('/v1/internal/content/leetcode/dispatch')) {
        assert.equal(req.headers['x-worker-key'], 'test-key')
        res.writeHead(200, { 'Content-Type': 'application/json' })
        res.end(JSON.stringify({ message: 'hello group', templateSlug: 'problem_daily' }))
        return
      }
      res.writeHead(404)
      res.end()
    }, async (baseUrl) => {
      const cfg: Config = {
        port: 3001,
        databaseUrl: '',
        managementApiUrl: baseUrl,
        workerApiKey: 'test-key',
        sessionDir: './data/session',
        defaultGroupJid: '',
        timezone: 'America/Sao_Paulo',
      }
      const settings = await fetchSettings(cfg)
      assert.equal(settings.problemPostTime, '09:00')
      const dispatch = await fetchDispatch(cfg, 'problem', '2026-06-16')
      assert.equal(dispatch.message, 'hello group')
    })
  })
})
