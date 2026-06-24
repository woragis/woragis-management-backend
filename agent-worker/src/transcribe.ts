import OpenAI from 'openai'
import { toFile } from 'openai/uploads'

export async function transcribeAudio(
  openai: OpenAI,
  audio: Buffer,
  mimeType = 'audio/ogg',
): Promise<string> {
  const ext = mimeType.includes('mpeg') ? 'mp3' : 'ogg'
  const file = await toFile(audio, `voice.${ext}`, { type: mimeType })
  const result = await openai.audio.transcriptions.create({
    model: 'whisper-1',
    file,
    language: 'pt',
  })
  return result.text.trim()
}
