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

## Roadmap consolidado (ordem sugerida)

```text
C1  Contacts core (CRUD + busca)
C2  Contact interactions
F1  Finance ↔ contactId
P1  Agent personality (Postgres + Redis + tools)
A0  Agent API foundation (AGENT_API_KEY, /v1/internal/agent/*)
A1  agent-worker + Telegram texto + tools contacts/finance
A2  Áudio Telegram (Whisper / Realtime)
A3  WhatsApp conversacional
A4  Voz telefônica (Twilio + Realtime)
UI  Frontend Contacts (opcional, paralelo)
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
| Agent orchestrator | Planejado |
| Contacts | Planejado |
| Finance link | Planejado |
| Agent personality | Planejado |
| Implementação | Não iniciado |
