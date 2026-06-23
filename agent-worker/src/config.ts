export type Config = {
  port: number
  managementApiUrl: string
  agentApiKey: string
  channelWorkerKey: string
  openaiApiKey: string
  telegramWorkerUrl: string
  whatsappWorkerUrl: string
  twilioAccountSid: string
  twilioAuthToken: string
  twilioPhoneNumber: string
}

export function loadConfig(): Config {
  const port = Number(process.env.PORT ?? '3002')
  const managementApiUrl = (process.env.MANAGEMENT_API_URL ?? 'http://127.0.0.1:8080').replace(/\/$/, '')
  const agentApiKey = (process.env.AGENT_API_KEY ?? '').trim()
  const channelWorkerKey = (process.env.CHANNEL_WORKER_KEY ?? agentApiKey).trim()
  const openaiApiKey = (process.env.OPENAI_API_KEY ?? '').trim()
  const telegramWorkerUrl = (process.env.TELEGRAM_WORKER_URL ?? '').replace(/\/$/, '')
  const whatsappWorkerUrl = (process.env.WHATSAPP_WORKER_URL ?? '').replace(/\/$/, '')
  const twilioAccountSid = (process.env.TWILIO_ACCOUNT_SID ?? '').trim()
  const twilioAuthToken = (process.env.TWILIO_AUTH_TOKEN ?? '').trim()
  const twilioPhoneNumber = (process.env.TWILIO_PHONE_NUMBER ?? '').trim()

  if (!agentApiKey) throw new Error('AGENT_API_KEY is required')

  return {
    port,
    managementApiUrl,
    agentApiKey,
    channelWorkerKey,
    openaiApiKey,
    telegramWorkerUrl,
    whatsappWorkerUrl,
    twilioAccountSid,
    twilioAuthToken,
    twilioPhoneNumber,
  }
}
