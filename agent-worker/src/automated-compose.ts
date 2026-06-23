import OpenAI from 'openai'

export type ComposeInput = {
  templateBody: string
  composeMode?: 'static' | 'ai_assisted'
  data?: Record<string, unknown>
  destinationContext?: Record<string, unknown>
}

function renderStatic(template: string, data: Record<string, unknown> = {}): string {
  return template.replace(/\{\{\s*([\w.]+)\s*\}\}/g, (_match, key: string) => {
    const value = data[key]
    return value == null ? '' : String(value)
  })
}

export async function composeMessage(openai: OpenAI | null, input: ComposeInput): Promise<string> {
  const mode = input.composeMode ?? 'static'
  const data = input.data ?? {}

  if (mode === 'static' || !openai) {
    return renderStatic(input.templateBody, data)
  }

  const completion = await openai.chat.completions.create({
    model: 'gpt-4o-mini',
    messages: [
      {
        role: 'system',
        content:
          'Compose a concise outbound message from the template and data. Keep placeholders filled. Reply with message text only.',
      },
      {
        role: 'user',
        content: JSON.stringify({
          template: input.templateBody,
          data,
          destinationContext: input.destinationContext ?? {},
        }),
      },
    ],
  })
  return completion.choices[0]?.message?.content?.trim() || renderStatic(input.templateBody, data)
}
