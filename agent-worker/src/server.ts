import http from 'node:http'
import OpenAI from 'openai'
import type { Config } from './config.js'
import { ManagementClient } from './api-client.js'
import { AgentLoop, type ChatSession } from './agent.js'
import { composeMessage, type ComposeInput } from './automated-compose.js'
import { sendToChannel, type InboundPayload } from './channel-clients.js'
import pino from 'pino'

const log = pino({ name: 'agent-worker' })

function sessionKey(payload: InboundPayload): string {
  const user = payload.userId || payload.externalId
  return `${payload.channel}:${user}`
}

function authorizeInbound(cfg: Config, req: http.IncomingMessage): boolean {
  if (!cfg.channelWorkerKey) return false
  const header = req.headers['x-channel-worker-key']
  const auth = req.headers.authorization?.replace(/^Bearer\s+/i, '')
  const key = typeof header === 'string' ? header : auth
  return key === cfg.channelWorkerKey
}

function readJson<T>(req: http.IncomingMessage): Promise<T> {
  return new Promise((resolve, reject) => {
    const chunks: Buffer[] = []
    req.on('data', (c) => chunks.push(c))
    req.on('end', () => {
      try {
        const raw = Buffer.concat(chunks).toString('utf8')
        resolve(raw ? (JSON.parse(raw) as T) : ({} as T))
      } catch (err) {
        reject(err)
      }
    })
    req.on('error', reject)
  })
}

export function createAgentServer(
  cfg: Config,
  sessions: Map<string, ChatSession>,
  agent: AgentLoop | null,
  openai: OpenAI | null,
): http.Server {
  return http.createServer(async (req, res) => {
    const send = (status: number, body: unknown) => {
      res.writeHead(status, { 'Content-Type': 'application/json' })
      res.end(JSON.stringify(body))
    }

    if (req.method === 'GET' && req.url === '/') {
      send(200, {
        ok: true,
        service: 'agent-worker',
        openai: !!cfg.openaiApiKey,
        telegram: !!cfg.telegramWorkerUrl,
        whatsapp: !!cfg.whatsappWorkerUrl,
      })
      return
    }

    if (req.method === 'POST' && req.url === '/v1/inbound') {
      if (!authorizeInbound(cfg, req)) {
        send(401, { error: 'unauthorized' })
        return
      }
      if (!agent) {
        send(503, { error: 'OPENAI_API_KEY not configured' })
        return
      }
      try {
        const payload = await readJson<InboundPayload>(req)
        if (!payload.channel || !payload.externalId || !payload.text?.trim()) {
          send(400, { error: 'channel, externalId, and text are required' })
          return
        }
        const key = sessionKey(payload)
        let session = sessions.get(key)
        if (!session) {
          session = { messages: [], greeted: false }
          sessions.set(key, session)
        }
        const reply = await agent.handleUserMessage(session, payload.text.trim())
        await sendToChannel(cfg, payload.channel, payload.externalId, reply, payload.destinationId)
        send(200, { ok: true, reply })
      } catch (err) {
        log.error({ err }, 'inbound failed')
        send(500, { error: err instanceof Error ? err.message : String(err) })
      }
      return
    }

    if (req.method === 'POST' && req.url === '/v1/automated/compose') {
      const header = req.headers['x-agent-key']
      const auth = req.headers.authorization?.replace(/^Bearer\s+/i, '')
      const key = typeof header === 'string' ? header : auth
      if (key !== cfg.agentApiKey) {
        send(401, { error: 'unauthorized' })
        return
      }
      try {
        const body = await readJson<ComposeInput>(req)
        if (!body.templateBody?.trim()) {
          send(400, { error: 'templateBody is required' })
          return
        }
        const message = await composeMessage(openai, body)
        send(200, { message })
      } catch (err) {
        log.error({ err }, 'compose failed')
        send(500, { error: err instanceof Error ? err.message : String(err) })
      }
      return
    }

    send(404, { error: 'not found' })
  })
}

export function createRuntime(cfg: Config): {
  sessions: Map<string, ChatSession>
  agent: AgentLoop | null
  openai: OpenAI | null
} {
  const api = new ManagementClient(cfg)
  const openai = cfg.openaiApiKey ? new OpenAI({ apiKey: cfg.openaiApiKey }) : null
  const agent = openai ? new AgentLoop(cfg.openaiApiKey, api) : null
  return { sessions: new Map(), agent, openai }
}
