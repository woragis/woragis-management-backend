import { Telegraf } from 'telegraf'
import type { Config } from '../config.js'
import { ManagementClient } from '../api-client.js'
import { AgentLoop, type ChatSession } from '../agent.js'
import pino from 'pino'

const log = pino({ name: 'agent-worker' })

export function startTelegramBot(cfg: Config): void {
  const api = new ManagementClient(cfg)
  const agent = new AgentLoop(cfg.openaiApiKey, api)
  const sessions = new Map<number, ChatSession>()

  const bot = new Telegraf(cfg.telegramBotToken)

  bot.use(async (ctx, next) => {
    const userId = ctx.from?.id
    if (!userId) return
    if (cfg.allowedUserIds.size > 0 && !cfg.allowedUserIds.has(userId)) {
      await ctx.reply('Acesso não autorizado.')
      return
    }
    return next()
  })

  bot.command('start', async (ctx) => {
    const userId = ctx.from?.id
    if (!userId) return
    sessions.delete(userId)
    const p = await api.getPersonality()
    await ctx.reply(`Olá! Sou ${p.assistantName}. Envie uma mensagem para começar.`)
  })

  bot.command('reset', async (ctx) => {
    const userId = ctx.from?.id
    if (!userId) return
    sessions.delete(userId)
    await ctx.reply('Conversa reiniciada.')
  })

  bot.on('text', async (ctx) => {
    const userId = ctx.from?.id
    const text = ctx.message.text?.trim()
    if (!userId || !text) return

    let session = sessions.get(userId)
    if (!session) {
      session = { messages: [], greeted: false }
      sessions.set(userId, session)
    }

    try {
      await ctx.sendChatAction('typing')
      const reply = await agent.handleUserMessage(session, text)
      await ctx.reply(reply, { parse_mode: undefined })
    } catch (err) {
      log.error({ err, userId }, 'agent turn failed')
      await ctx.reply('Não consegui processar agora. Tente de novo em instantes.')
    }
  })

  bot.launch().then(() => log.info('telegram bot started'))
  process.once('SIGINT', () => bot.stop('SIGINT'))
  process.once('SIGTERM', () => bot.stop('SIGTERM'))
}
