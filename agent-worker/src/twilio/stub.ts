import type { Config } from '../config.js'
import pino from 'pino'

const log = pino({ name: 'agent-worker:twilio' })

export function initTwilioStub(cfg: Config): void {
  if (!cfg.twilioAccountSid || !cfg.twilioAuthToken) {
    log.info('Twilio not configured; voice calls disabled')
    return
  }
  log.info(
    { phone: cfg.twilioPhoneNumber || '(unset)' },
    'Twilio credentials present; Realtime voice integration not enabled yet',
  )
}
