export type KnownChat = {
  chatId: string
  name: string
  type: string
  lastSeenAt: string
}

const chats = new Map<string, KnownChat>()

export function registerChat(chatId: string, name: string, type: string): void {
  const id = chatId.trim()
  if (!id) return
  chats.set(id, {
    chatId: id,
    name: name.trim() || id,
    type: type.trim() || 'private',
    lastSeenAt: new Date().toISOString(),
  })
}

export function listKnownChats(): KnownChat[] {
  return [...chats.values()].sort((a, b) => a.name.localeCompare(b.name))
}
