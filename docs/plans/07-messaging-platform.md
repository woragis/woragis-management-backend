# Plano — Messaging Platform (v2)

Arquitetura de canais burros + cérebro no agente + cron centralizado.

## Serviços

| Serviço | Papel |
|---------|-------|
| **Go API** | Destinos, templates, jobs, deliveries, domínios |
| **scheduler-worker** | Tick → dispara jobs no Go |
| **agent-worker** | Conversa (tools, Whisper) + `automated/compose` (gpt-4o-mini) |
| **whatsapp-worker** | Baileys send/receive, sem cron nem LLM |
| **telegram-worker** | Bot send/receive, sem cron nem LLM |

## Domínio `messaging`

- `ChannelDestination` — grupo/chat (whatsapp jid, telegram chatId)
- `MessageTemplate` — corpo com `{{vars}}`, `composeMode`: static | ai_assisted
- `ScheduledJob` — cron + destino + template/programAction
- `MessageDelivery` — log de envios

## Contratos HTTP

### scheduler-worker → Go

- `GET /v1/internal/scheduler/due`
- `POST /v1/internal/scheduler/jobs/{id}/execute`

### Go → channel-workers

- `POST /v1/send` — `{ externalId, text, destinationId? }`

### channel-workers → agent-worker

- `POST /v1/inbound` — mensagem do usuário

### agent-worker → channel-workers

- `POST /v1/send` — resposta

### Go/agent → automated compose

- `POST /v1/automated/compose` — template + data + destination context

## Fases

- M1 ChannelDestination CRUD
- M2 MessageTemplate
- M3 ScheduledJob + MessageDelivery + scheduler API
- M4 scheduler-worker
- M5 WhatsApp slim (sem GROUP_JID, sem cron)
- M6 telegram-worker
- M7 agent inbound + automated
- M8 WhatsApp inbound
- M9 Twilio opcional
