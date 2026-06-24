import { Telegraf, type Context } from 'telegraf'
import type { Config } from './config.js'
import { forwardInbound, transcribeViaAgent } from './agent-client.js'
import { registerChat } from './chat-registry.js'
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
    const chatTitle =
      'title' in ctx.chat! ? ctx.chat.title : ctx.from?.first_name || String(chatId)
    registerChat(String(chatId), chatTitle, ctx.chat?.type ?? 'private')
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

  async function handleVoiceLike(ctx: Context, fileId: string) {
    const userId = ctx.from?.id
    const chatId = ctx.chat?.id
    if (!userId || chatId == null || !fileId) return
    const chatTitle =
      ctx.chat && 'title' in ctx.chat ? ctx.chat.title : ctx.from?.first_name || String(chatId)
    registerChat(String(chatId), chatTitle ?? String(chatId), ctx.chat?.type ?? 'private')
    if (!cfg.agentWorkerUrl) {
      await ctx.reply('Agente indisponível no momento.')
      return
    }
    try {
      await ctx.sendChatAction('typing')
      const link = await bot.telegram.getFileLink(fileId)
      const audioRes = await fetch(link.href)
      if (!audioRes.ok) throw new Error(`download failed: ${audioRes.status}`)
      const audio = Buffer.from(await audioRes.arrayBuffer())
      const text = await transcribeViaAgent(cfg, audio, 'audio/ogg')
      await forwardInbound(cfg, {
        channel: 'telegram',
        externalId: String(chatId),
        userId: String(userId),
        userName: ctx.from?.username || ctx.from?.first_name,
        text,
      })
    } catch (err) {
      log.error({ err, userId, chatId }, 'voice inbound failed')
      await ctx.reply('Não consegui processar o áudio agora. Tente de novo em instantes.')
    }
  }

  bot.on('voice', async (ctx) => {
    const fileId = ctx.message.voice?.file_id
    if (!fileId) return
    await handleVoiceLike(ctx, fileId)
  })

  bot.on('audio', async (ctx) => {
    const fileId = ctx.message.audio?.file_id
    if (!fileId) return
    await handleVoiceLike(ctx, fileId)
  })

  bot.launch().then(() => log.info('telegram bot started'))
  process.once('SIGINT', () => bot.stop('SIGINT'))
  process.once('SIGTERM', () => bot.stop('SIGTERM'))

  return bot
}
