import OpenAI from 'openai'
import type { ChatCompletionMessageParam, ChatCompletionTool } from 'openai/resources/chat/completions.js'
import type { ManagementClient, AgentPersonality, DestinationContext } from './api-client.js'
import { pickGreeting } from './personality/greeting.js'

export const toolDefinitions: ChatCompletionTool[] = [
  {
    type: 'function',
    function: {
      name: 'get_agent_personality',
      description: 'Get current assistant name, greetings, voice, and extra instructions.',
      parameters: { type: 'object', properties: {} },
    },
  },
  {
    type: 'function',
    function: {
      name: 'update_agent_personality',
      description: 'Update assistant personality (partial patch).',
      parameters: {
        type: 'object',
        properties: {
          assistantName: { type: 'string' },
          greetingMorning: { type: 'string' },
          greetingAfternoon: { type: 'string' },
          greetingEvening: { type: 'string' },
          greetingEnabled: { type: 'boolean' },
          voiceId: { type: 'string' },
          systemPromptExtra: { type: 'string' },
        },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'reset_agent_personality',
      description: 'Reset assistant personality to defaults.',
      parameters: { type: 'object', properties: {} },
    },
  },
  {
    type: 'function',
    function: {
      name: 'search_contacts',
      description: 'Search contacts by name, organization, email, or phone.',
      parameters: {
        type: 'object',
        properties: {
          q: { type: 'string' },
          relationship: { type: 'string' },
          organization: { type: 'string' },
          stage: { type: 'string' },
        },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'get_contact',
      description: 'Get a contact by UUID.',
      parameters: {
        type: 'object',
        properties: { id: { type: 'string' } },
        required: ['id'],
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'create_contact',
      description: 'Create a new contact.',
      parameters: {
        type: 'object',
        properties: {
          name: { type: 'string' },
          organization: { type: 'string' },
          roleTitle: { type: 'string' },
          relationship: { type: 'string' },
          email: { type: 'string' },
          phone: { type: 'string' },
          notes: { type: 'string' },
        },
        required: ['name'],
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'update_contact',
      description: 'Update an existing contact.',
      parameters: {
        type: 'object',
        properties: {
          id: { type: 'string' },
          nextFollowUpAt: { type: 'string', description: 'ISO 8601 datetime' },
          stage: { type: 'string' },
          notes: { type: 'string' },
        },
        required: ['id'],
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'log_interaction',
      description: 'Log a touchpoint with a contact.',
      parameters: {
        type: 'object',
        properties: {
          contactId: { type: 'string' },
          type: { type: 'string', enum: ['call', 'meeting', 'message', 'email', 'note'] },
          channel: { type: 'string' },
          summary: { type: 'string' },
          happenedAt: { type: 'string' },
        },
        required: ['contactId', 'type', 'summary'],
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'list_contacts_due_followup',
      description: 'List contacts with follow-up due.',
      parameters: { type: 'object', properties: {} },
    },
  },
  {
    type: 'function',
    function: {
      name: 'get_contact_finance',
      description: 'Income sources and transactions linked to a contact.',
      parameters: {
        type: 'object',
        properties: { id: { type: 'string' } },
        required: ['id'],
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'list_projects',
      description: 'List dev projects.',
      parameters: {
        type: 'object',
        properties: { q: { type: 'string' }, status: { type: 'string' } },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'create_project',
      description: 'Create a new dev project.',
      parameters: {
        type: 'object',
        properties: {
          name: { type: 'string' },
          description: { type: 'string' },
          status: { type: 'string' },
        },
        required: ['name'],
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'finance_dashboard',
      description: 'Finance dashboard overview.',
      parameters: { type: 'object', properties: {} },
    },
  },
  {
    type: 'function',
    function: {
      name: 'finance_summary',
      description: 'Monthly finance summary.',
      parameters: {
        type: 'object',
        properties: { year: { type: 'number' }, month: { type: 'number' } },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'list_transactions',
      description: 'List transactions with optional filters.',
      parameters: {
        type: 'object',
        properties: {
          year: { type: 'number' },
          month: { type: 'number' },
          type: { type: 'string' },
          contactId: { type: 'string' },
        },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'create_transaction',
      description: 'Create a financial transaction. Confirm amount and contact with user first.',
      parameters: {
        type: 'object',
        properties: {
          type: { type: 'string', enum: ['income', 'expense'] },
          amountCents: { type: 'number' },
          description: { type: 'string' },
          contactId: { type: 'string' },
          date: { type: 'string' },
        },
        required: ['type', 'amountCents', 'description'],
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'list_social_posts',
      description: 'List social media post drafts and scheduled posts for LinkedIn, Reddit, or Twitter.',
      parameters: {
        type: 'object',
        properties: {
          projectId: { type: 'string' },
          platform: { type: 'string', enum: ['linkedin', 'reddit', 'twitter'] },
          status: { type: 'string', enum: ['draft', 'scheduled', 'published', 'cancelled'] },
          goal: { type: 'string' },
        },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'list_post_templates',
      description: 'List reusable social post templates with {{variable}} placeholders.',
      parameters: {
        type: 'object',
        properties: {
          platform: { type: 'string' },
          goal: { type: 'string' },
        },
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'apply_post_template',
      description: 'Render a post template for a project (preview only, does not save).',
      parameters: {
        type: 'object',
        properties: {
          templateSlug: { type: 'string' },
          projectId: { type: 'string' },
        },
        required: ['templateSlug', 'projectId'],
      },
    },
  },
  {
    type: 'function',
    function: {
      name: 'create_social_post',
      description: 'Create a social post draft or scheduled post. Never set status to published without explicit user confirmation.',
      parameters: {
        type: 'object',
        properties: {
          projectId: { type: 'string' },
          campaignId: { type: 'string' },
          platform: { type: 'string', enum: ['linkedin', 'reddit', 'twitter'] },
          goal: { type: 'string' },
          status: { type: 'string', enum: ['draft', 'scheduled'] },
          title: { type: 'string' },
          body: { type: 'string' },
          hook: { type: 'string' },
          cta: { type: 'string' },
          templateSlug: { type: 'string' },
          scheduledAt: { type: 'string', description: 'ISO 8601 datetime' },
        },
        required: ['platform', 'body'],
      },
    },
  },
]

export type ChatSession = {
  messages: ChatCompletionMessageParam[]
  greeted: boolean
  destinationContext?: DestinationContext
}

export function buildSystemPrompt(p: AgentPersonality, dest?: DestinationContext): string {
  const extra = p.systemPromptExtra?.trim()
  const lines = [
    `Você é ${p.assistantName}, assistente pessoal operacional do usuário.`,
    'Responda em português brasileiro, de forma clara e objetiva.',
    'Use as tools para dados reais; nunca invente IDs ou valores financeiros.',
    'Antes de ações destrutivas ou financeiras, peça confirmação explícita.',
    'Posts sociais: crie apenas como draft ou scheduled; nunca marque published sem confirmação.',
    'Se houver contatos homônimos, desambigüe por organização e cargo.',
  ]
  if (dest) {
    lines.push(
      `Contexto do canal: conversa no destino "${dest.name}" (${dest.channel}, externalId ${dest.externalId}).`,
    )
    if (dest.description?.trim()) lines.push(`Descrição do destino: ${dest.description.trim()}`)
    if (dest.responsibilities?.trim()) {
      lines.push(`Responsabilidades do destino: ${dest.responsibilities.trim()}`)
    }
    if (dest.tags?.length) lines.push(`Tags: ${dest.tags.join(', ')}`)
  }
  if (extra) lines.push(`Instruções adicionais: ${extra}`)
  return lines.join('\n')
}

export async function runTool(
  api: ManagementClient,
  name: string,
  args: Record<string, unknown>,
): Promise<unknown> {
  switch (name) {
    case 'get_agent_personality':
      return api.getPersonality()
    case 'update_agent_personality':
      return api.updatePersonality(args)
    case 'reset_agent_personality':
      return api.resetPersonality()
    case 'search_contacts':
      return api.searchContacts(stringParams(args, ['q', 'relationship', 'organization', 'stage']))
    case 'get_contact':
      return api.getContact(String(args.id))
    case 'create_contact':
      return api.createContact(args)
    case 'update_contact': {
      const id = String(args.id)
      const { id: _omit, ...body } = args
      return api.updateContact(id, body)
    }
    case 'log_interaction':
      return api.logInteraction(String(args.contactId), {
        type: args.type,
        channel: args.channel,
        summary: args.summary,
        happenedAt: args.happenedAt,
      })
    case 'list_contacts_due_followup':
      return api.listContactsDueFollowUp()
    case 'get_contact_finance':
      return api.getContactFinance(String(args.id))
    case 'list_projects':
      return api.listProjects(stringParams(args, ['q', 'status']))
    case 'create_project':
      return api.createProject(args)
    case 'finance_dashboard':
      return api.financeDashboard()
    case 'finance_summary':
      return api.financeSummary(num(args.year), num(args.month))
    case 'list_transactions':
      return api.listTransactions(stringParams(args, ['type', 'contactId'], numMap(args, ['year', 'month'])))
    case 'create_transaction':
      return api.createTransaction(args)
    case 'list_social_posts':
      return api.listSocialPosts(stringParams(args, ['projectId', 'platform', 'status', 'goal']))
    case 'list_post_templates':
      return api.listPostTemplates(stringParams(args, ['platform', 'goal']))
    case 'apply_post_template':
      return api.applyPostTemplate({
        templateSlug: String(args.templateSlug),
        projectId: String(args.projectId),
      })
    case 'create_social_post':
      return api.createSocialPost(args)
    default:
      throw new Error(`Unknown tool: ${name}`)
  }
}

function stringParams(
  args: Record<string, unknown>,
  keys: string[],
  nums: Record<string, number | undefined> = {},
): Record<string, string> {
  const out: Record<string, string> = {}
  for (const k of keys) {
    if (args[k] != null && args[k] !== '') out[k] = String(args[k])
  }
  for (const [k, v] of Object.entries(nums)) {
    if (v != null) out[k] = String(v)
  }
  return out
}

function numMap(args: Record<string, unknown>, keys: string[]): Record<string, number | undefined> {
  const out: Record<string, number | undefined> = {}
  for (const k of keys) out[k] = num(args[k])
  return out
}

function num(v: unknown): number | undefined {
  if (typeof v === 'number' && Number.isFinite(v)) return v
  return undefined
}

export class AgentLoop {
  private readonly openai: OpenAI

  constructor(
    openaiApiKey: string,
    private readonly api: ManagementClient,
  ) {
    this.openai = new OpenAI({ apiKey: openaiApiKey })
  }

  async handleUserMessage(session: ChatSession, userText: string): Promise<string> {
    const personality = await this.api.getPersonality()
    if (session.messages.length === 0) {
      session.messages.push({
        role: 'system',
        content: buildSystemPrompt(personality, session.destinationContext),
      })
    }
    let prefix = ''
    if (!session.greeted) {
      const greeting = pickGreeting(personality)
      if (greeting) prefix = `${greeting}. `
      session.greeted = true
    }
    session.messages.push({ role: 'user', content: userText })

    for (let step = 0; step < 8; step++) {
      const completion = await this.openai.chat.completions.create({
        model: 'gpt-4o-mini',
        messages: session.messages,
        tools: toolDefinitions,
        tool_choice: 'auto',
      })
      const choice = completion.choices[0]?.message
      if (!choice) throw new Error('Empty OpenAI response')

      session.messages.push(choice)

      const toolCalls = choice.tool_calls ?? []
      if (toolCalls.length === 0) {
        const text = choice.content?.trim() ?? 'Pronto.'
        return prefix + text
      }

      for (const call of toolCalls) {
        if (call.type !== 'function') continue
        const args = JSON.parse(call.function.arguments || '{}') as Record<string, unknown>
        let result: unknown
        try {
          result = await runTool(this.api, call.function.name, args)
        } catch (err) {
          result = { error: err instanceof Error ? err.message : String(err) }
        }
        session.messages.push({
          role: 'tool',
          tool_call_id: call.id,
          content: JSON.stringify(result),
        })
      }
    }

    return prefix + 'Precisei de muitos passos; tente reformular a pergunta.'
  }
}
