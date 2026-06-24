import type { Config } from './config.js'

export type ChannelDestination = {
  id: string
  channel: string
  externalId: string
  name: string
  description?: string
  responsibilities?: string
  tags?: string[]
  active?: boolean
}

export type DestinationContext = Pick<
  ChannelDestination,
  'id' | 'channel' | 'externalId' | 'name' | 'description' | 'responsibilities' | 'tags'
>

export type AgentPersonality = {
  assistantName: string
  greetingMorning: string
  greetingAfternoon: string
  greetingEvening: string
  greetingEnabled: boolean
  systemPromptExtra: string
  voiceId: string
  language: string
  timezone: string
}

function headers(cfg: Config): Record<string, string> {
  return {
    'Content-Type': 'application/json',
    'X-Agent-Key': cfg.agentApiKey,
    Authorization: `Bearer ${cfg.agentApiKey}`,
  }
}

async function request<T>(cfg: Config, path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${cfg.managementApiUrl}${path}`, {
    ...init,
    headers: { ...headers(cfg), ...(init?.headers ?? {}) },
  })
  const text = await res.text()
  if (!res.ok) {
    throw new Error(`API ${path} failed: ${res.status} ${text}`)
  }
  if (!text) return undefined as T
  return JSON.parse(text) as T
}

export class ManagementClient {
  constructor(private readonly cfg: Config) {}

  getPersonality() {
    return request<AgentPersonality>(this.cfg, '/v1/internal/agent/personality')
  }

  updatePersonality(patch: Record<string, unknown>) {
    return request<AgentPersonality>(this.cfg, '/v1/internal/agent/personality', {
      method: 'PATCH',
      body: JSON.stringify(patch),
    })
  }

  resetPersonality() {
    return request<AgentPersonality>(this.cfg, '/v1/internal/agent/personality/reset', {
      method: 'POST',
    })
  }

  searchContacts(params: Record<string, string>) {
    const q = new URLSearchParams(params).toString()
    return request<unknown[]>(this.cfg, `/v1/internal/agent/tools/contacts?${q}`)
  }

  getContact(id: string) {
    return request<unknown>(this.cfg, `/v1/internal/agent/tools/contacts/${id}`)
  }

  createContact(body: Record<string, unknown>) {
    return request<unknown>(this.cfg, '/v1/internal/agent/tools/contacts', {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  updateContact(id: string, body: Record<string, unknown>) {
    return request<unknown>(this.cfg, `/v1/internal/agent/tools/contacts/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(body),
    })
  }

  logInteraction(contactId: string, body: Record<string, unknown>) {
    return request<unknown>(this.cfg, `/v1/internal/agent/tools/contacts/${contactId}/interactions`, {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  listContactsDueFollowUp(before?: string) {
    const q = before ? `?before=${encodeURIComponent(before)}` : ''
    return request<unknown[]>(this.cfg, `/v1/internal/agent/tools/contacts/due-follow-up${q}`)
  }

  getContactFinance(id: string) {
    return request<unknown>(this.cfg, `/v1/internal/agent/tools/contacts/${id}/finance`)
  }

  listProjects(params: Record<string, string> = {}) {
    const q = new URLSearchParams(params).toString()
    const suffix = q ? `?${q}` : ''
    return request<unknown[]>(this.cfg, `/v1/internal/agent/tools/projects${suffix}`)
  }

  getProject(id: string) {
    return request<unknown>(this.cfg, `/v1/internal/agent/tools/projects/${id}`)
  }

  createProject(body: Record<string, unknown>) {
    return request<unknown>(this.cfg, '/v1/internal/agent/tools/projects', {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  financeDashboard() {
    return request<unknown>(this.cfg, '/v1/internal/agent/tools/finance/dashboard')
  }

  financeSummary(year?: number, month?: number) {
    const params = new URLSearchParams()
    if (year) params.set('year', String(year))
    if (month) params.set('month', String(month))
    const q = params.toString()
    return request<unknown>(this.cfg, `/v1/internal/agent/tools/finance/summary${q ? `?${q}` : ''}`)
  }

  financeCalendar(year?: number, month?: number) {
    const params = new URLSearchParams()
    if (year) params.set('year', String(year))
    if (month) params.set('month', String(month))
    const q = params.toString()
    return request<unknown>(this.cfg, `/v1/internal/agent/tools/finance/calendar${q ? `?${q}` : ''}`)
  }

  listIncomeSources(params: Record<string, string> = {}) {
    const q = new URLSearchParams(params).toString()
    const suffix = q ? `?${q}` : ''
    return request<unknown[]>(this.cfg, `/v1/internal/agent/tools/finance/income-sources${suffix}`)
  }

  listTransactions(params: Record<string, string> = {}) {
    const q = new URLSearchParams(params).toString()
    const suffix = q ? `?${q}` : ''
    return request<unknown[]>(this.cfg, `/v1/internal/agent/tools/finance/transactions${suffix}`)
  }

  createTransaction(body: Record<string, unknown>) {
    return request<unknown>(this.cfg, '/v1/internal/agent/tools/finance/transactions', {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  createIncomeSource(body: Record<string, unknown>) {
    return request<unknown>(this.cfg, '/v1/internal/agent/tools/finance/income-sources', {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  listSocialPosts(params: Record<string, string> = {}) {
    const q = new URLSearchParams(params).toString()
    const suffix = q ? `?${q}` : ''
    return request<unknown[]>(this.cfg, `/v1/internal/agent/tools/presence/posts${suffix}`)
  }

  listPostTemplates(params: Record<string, string> = {}) {
    const q = new URLSearchParams(params).toString()
    const suffix = q ? `?${q}` : ''
    return request<unknown[]>(this.cfg, `/v1/internal/agent/tools/presence/templates${suffix}`)
  }

  applyPostTemplate(body: { templateSlug: string; projectId: string }) {
    return request<unknown>(this.cfg, '/v1/internal/agent/tools/presence/apply-template', {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  createSocialPost(body: Record<string, unknown>) {
    return request<unknown>(this.cfg, '/v1/internal/agent/tools/presence/posts', {
      method: 'POST',
      body: JSON.stringify(body),
    })
  }

  resolveDestination(channel: string, externalId: string) {
    const q = new URLSearchParams({ channel, externalId }).toString()
    return request<ChannelDestination>(this.cfg, `/v1/internal/agent/tools/messaging/resolve-destination?${q}`)
  }
}
