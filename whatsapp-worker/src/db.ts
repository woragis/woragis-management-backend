import pg from 'pg'
import type { Config } from './config.js'

const { Pool } = pg

let pool: pg.Pool | null = null

export function getPool(cfg: Config): pg.Pool {
  if (!pool) {
    if (!cfg.databaseUrl) throw new Error('DATABASE_URL is required')
    pool = new Pool({ connectionString: cfg.databaseUrl })
  }
  return pool
}

export async function migrate(cfg: Config): Promise<void> {
  const db = getPool(cfg)
  await db.query(`
    CREATE TABLE IF NOT EXISTS whatsapp_groups (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      jid TEXT NOT NULL UNIQUE,
      name TEXT NOT NULL DEFAULT '',
      enabled BOOLEAN NOT NULL DEFAULT true,
      channel_slug TEXT NOT NULL DEFAULT 'leetcode',
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
    CREATE TABLE IF NOT EXISTS message_log (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      group_jid TEXT NOT NULL,
      template_slug TEXT NOT NULL DEFAULT '',
      dispatch_type TEXT NOT NULL DEFAULT '',
      body TEXT NOT NULL,
      video_id TEXT,
      status TEXT NOT NULL DEFAULT 'sent',
      error_message TEXT,
      sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
    CREATE TABLE IF NOT EXISTS worker_state (
      key TEXT PRIMARY KEY,
      value TEXT NOT NULL DEFAULT '',
      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
  `)
}

export async function ensureDefaultGroup(cfg: Config, jid: string, name = 'LeetCode'): Promise<string> {
  if (!jid) return ''
  const db = getPool(cfg)
  await db.query(
    `INSERT INTO whatsapp_groups (jid, name, enabled, channel_slug)
     VALUES ($1, $2, true, 'leetcode')
     ON CONFLICT (jid) DO UPDATE SET name = EXCLUDED.name, updated_at = NOW()`,
    [jid, name],
  )
  return jid
}

export async function getEnabledGroupJid(cfg: Config): Promise<string | null> {
  const db = getPool(cfg)
  const res = await db.query<{ jid: string }>(
    `SELECT jid FROM whatsapp_groups WHERE enabled = true ORDER BY created_at ASC LIMIT 1`,
  )
  if (res.rows[0]?.jid) return res.rows[0].jid
  if (cfg.defaultGroupJid) return cfg.defaultGroupJid
  return null
}

export async function logMessage(
  cfg: Config,
  entry: {
    groupJid: string
    templateSlug: string
    dispatchType: string
    body: string
    videoId?: string
    status: string
    error?: string
  },
): Promise<void> {
  const db = getPool(cfg)
  await db.query(
    `INSERT INTO message_log (group_jid, template_slug, dispatch_type, body, video_id, status, error_message)
     VALUES ($1, $2, $3, $4, $5, $6, $7)`,
    [
      entry.groupJid,
      entry.templateSlug,
      entry.dispatchType,
      entry.body,
      entry.videoId || null,
      entry.status,
      entry.error || null,
    ],
  )
}

export async function setState(cfg: Config, key: string, value: string): Promise<void> {
  const db = getPool(cfg)
  await db.query(
    `INSERT INTO worker_state (key, value, updated_at) VALUES ($1, $2, NOW())
     ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()`,
    [key, value],
  )
}

export async function getState(cfg: Config, key: string): Promise<string> {
  const db = getPool(cfg)
  const res = await db.query<{ value: string }>(`SELECT value FROM worker_state WHERE key = $1`, [key])
  return res.rows[0]?.value || ''
}
