import { Telegraf } from 'telegraf'
import type { Config } from './config.js'
import { forwardInbound } from './agent-client.js'
import pino from 'pino'

const log = pino({ name: 'telegram-worker' })

export function startTelegramBot(cfg: Config): Telegraf {
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

  bot.on('text', async (ctx) => {
    const userId = ctx.from?.id
    const chatId = ctx.chat?.id
    const text = ctx.message.text?.trim()
    if (!userId || chatId == null || !text) return
    if (!cfg.agentWorkerUrl) {
      await ctx.reply('Agente indisponível no momento.')
      return
    }
    try {
      await ctx.sendChatAction('typing')
      await forwardInbound(cfg, {
        channel: 'telegram',
        externalId: String(chatId),
        userId: String(userId),
        userName: ctx.from?.username || ctx.from?.first_name,
        text,
      })
    } catch (err) {
      log.error({ err, userId, chatId }, 'forward inbound failed')
      await ctx.reply('Não consegui encaminhar sua mensagem agora. Tente de novo em instantes.')
    }
  })

  bot.launch().then(() => log.info('telegram bot started'))
  process.once('SIGINT', () => bot.stop('SIGINT'))
  process.once('SIGTERM', () => bot.stop('SIGTERM'))

  return bot
}
