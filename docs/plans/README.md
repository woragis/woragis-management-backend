# Planos — Woragis Management & Agent

Documentação de planejamento para o assistente pessoal, domínio de contatos, finanças e personalidade do agente.

## Índice

| Documento | Conteúdo |
|-----------|----------|
| [01-agent-orchestrator.md](./01-agent-orchestrator.md) | Arquitetura do `agent-worker`, canais, tools, fases A0–A4 |
| [02-contacts-domain.md](./02-contacts-domain.md) | CRM de contatos (investidores, clientes, sócios), fases C1–C2 |
| [03-finance-agent-integration.md](./03-finance-agent-integration.md) | Ligação finanças ↔ contatos ↔ bot, fase F1 |
| [04-agent-personality.md](./04-agent-personality.md) | Nome, saudações, voz; Redis + Postgres; latência |
| [05-entity-catalog.md](./05-entity-catalog.md) | Entidades existentes e o que o agente já pode ler |
| [07-messaging-platform.md](./07-messaging-platform.md) | Canais, templates, jobs, workers |
| [08-projects-taxonomy-and-presence.md](./08-projects-taxonomy-and-presence.md) | Dimensões de projeto e domínio Presence |
| [10-gaps-roadmap.md](./10-gaps-roadmap.md) | Roadmap consolidado dos gaps pós-T3/P4a |

## Roadmap consolidado (ordem sugerida)

```text
C1  Contacts core (CRUD + busca)          ✓
C2  Contact interactions                   ✓
F1  Finance ↔ contactId                   ✓
P1  Agent personality (Postgres + Redis)   ✓
A0  Agent API foundation                   ✓
A1  agent-worker + Telegram texto          ✓ (com resolve destination)
A2  Áudio Telegram (Whisper)               ✓
A3  WhatsApp conversacional                ✓ parcial (inbound + destination context)
A4  Voz telefônica (Twilio)                stub
UI  Frontend Contacts + Personality       ✓
10.x Gaps roadmap (messaging, dashboard, presence, finance alerts)
```

## Princípios

1. **Go API** = fonte da verdade e regras de negócio.
2. **agent-worker** = orquestração, canais, OpenAI; não acessa Postgres direto.
3. **whatsapp-worker** = entrega WhatsApp apenas (Baileys).
4. **Redis** = cache quente para config do agente (personalidade); Postgres = persistência.
5. Tools com escopo — secrets e valores financeiros exigem confirmação.

## Status

| Fase | Status |
|------|--------|
| Agent orchestrator (A0–A3, A2 áudio) | Implementado |
| Contacts (C1–C2) + UI | Implementado |
| Finance link (F1) + alertas dashboard | Implementado |
| Agent personality + UI | Implementado |
| Messaging (resolve destination, sync Telegram, índice único) | Implementado |
| Projects dashboard (maturidade/intent) | Implementado |
| Presence P4b-lite (marcar publicado + URL) | Implementado |
| A4 Twilio voz | Stub apenas |
