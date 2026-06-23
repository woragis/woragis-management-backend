import { loadConfig } from './config.js'
import { createAgentServer, createRuntime } from './server.js'
import { initTwilioStub } from './twilio/stub.js'
import pino from 'pino'

const log = pino({ name: 'agent-worker' })

async function main() {
  const cfg = loadConfig()
  const { sessions, agent, openai } = createRuntime(cfg)

  if (!cfg.openaiApiKey) {
    log.warn('OPENAI_API_KEY not set; conversational inbound disabled')
  }

  initTwilioStub(cfg)

  const server = createAgentServer(cfg, sessions, agent, openai)
  server.listen(cfg.port, () => log.info({ port: cfg.port }, 'agent-worker listening'))
}

main().catch((err) => {
  log.error(err)
  process.exit(1)
})
